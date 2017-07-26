package main

import (
	"fmt"
	"golang.org/x/crypto/openpgp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func (worker *SCORSHworker) Matches(repo, branch string) bool {

	for _, r := range worker.Repos {
		parts := strings.SplitN(r, ":", 2)
		repo_pattern := parts[0]
		branch_pattern := parts[1]
		repo_match, _ := regexp.MatchString(repo_pattern, repo)
		branch_match, _ := regexp.MatchString(branch_pattern, branch)
		debug.log("[worker.Matches] repo_match: %s\n", repo_match)
		debug.log("[worker.Matches] branch_match: %s\n", branch_match)
		if repo_match && branch_match {
			return true
		}
	}
	return false
}

func (w *SCORSHworker) LoadKeyrings() error {

	w.Keys = make(map[string]openpgp.KeyRing)
	w.TagKeys = make(map[string]map[string]bool)

	for _, t := range w.Tags {
		w.TagKeys[t.Name] = make(map[string]bool)

		// Open the keyring files
		for _, keyring := range t.Keyrings {
			if _, ok := w.Keys[keyring]; ok {
				// keyring has been loaded: just add it to the TagKeys map
				w.TagKeys[t.Name][keyring] = true
				continue
			}
			k_file := fmt.Sprintf("%s/%s", w.Folder, keyring)
			debug.log("[worker: %s] Trying to open keyring at %s\n", w.Name, k_file)
			f, err_file := os.Open(k_file)
			if err_file != nil {
				log.Printf("[worker] cannot open keyring: %s", err_file)
				f.Close()
			}

			// load the keyring
			kr, err_key := openpgp.ReadArmoredKeyRing(f)

			if err_key != nil {
				log.Printf("[worker] cannot load keyring: %s", err_key)
				f.Close()
				//return fmt.Errorf("Unable to load keyring: ", err_key)
			}
			w.Keys[keyring] = kr
			w.TagKeys[t.Name][keyring] = true
			f.Close()
		}
	}
	return nil
}

// Still to be implemented
func (w *SCORSHworker) LoadTags() error {

	w_tags, err := ioutil.ReadFile(w.Tagfile)
	if err != nil {
		return fmt.Errorf("Cannot read worker config: %s", err)
	}

	err = yaml.Unmarshal(w_tags, w)
	//err = yaml.Unmarshal(w_tags, tags)

	if err != nil {
		return fmt.Errorf("Error while reading tags: %s", err)
	}

	return nil
}

// FIXME--- still needs some work...
func runWorker(w *SCORSHworker) {

	var msg SCORSHmsg

	log.Printf("[worker: %s] Started\n", w.Name)

	// notify that we have been started!
	w.StatusChan <- msg

	// This is the main worker loop
	for {
		select {
		case msg = <-w.MsgChan:
			debug.log("[worker: %s] received message %s\n", w.Name, msg.Id)
			// process message
			err := walkCommits(msg, w)
			if err != nil {
				log.Printf("[worker: %s] error in walk_commits: %s", err)
			}
			w.StatusChan <- msg
			debug.log("[worker: %s] Sent message back: %s", w.Name, msg)
		}
	}
}

// StartWorkers starts all the workers specified in a given
// configuration and fills in the SCORSHmaster struct
func startWorkers(master *SCORSHmaster) error {

	num_workers := len(master.Workers)

	// We should now start each worker

	log.Printf("num_workers: %d\n", num_workers)

	for w := 0; w < num_workers; w++ {

		worker := &(master.Workers[w])
		// Set the Status and Msg channels
		worker.StatusChan = master.StatusChan
		worker.MsgChan = make(chan SCORSHmsg, 10)

		// Load worker tags from worker.Tagfile
		err := worker.LoadTags()
		if err != nil {
			close(worker.MsgChan)
			return fmt.Errorf("[Starting worker: %s] Unable to load tags: %s\n", worker.Name, err)
		}

		// Load worker keyrings -- this must be called *after* LoadTags!!!!
		err = worker.LoadKeyrings()
		if err != nil {
			close(worker.MsgChan)
			return fmt.Errorf("[Starting worker: %s] Unable to load keyrings: %s\n", worker.Name, err)
		}

		// Add the repos definitions to the map master.Repos
		for _, repo_name := range worker.Repos {
			master.Repos[repo_name] = append(master.Repos[repo_name], worker)
		}
		go runWorker(worker)
		<-master.StatusChan
	}
	return nil
}
