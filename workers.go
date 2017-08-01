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

// Matches returns true if the configured repo:branch of the worker
// matches the repo and branch provided as arguments
func (w *SCORSHworker) Matches(repo, branch string) bool {

	for _, r := range w.Repos {
		parts := strings.SplitN(r, ":", 2)
		repoPattern := parts[0]
		branchPattern := parts[1]
		repoMatch, _ := regexp.MatchString(repoPattern, repo)
		branchMatch, _ := regexp.MatchString(branchPattern, branch)
		debug.log("[worker.Matches] repo_match: %s\n", repoMatch)
		debug.log("[worker.Matches] branch_match: %s\n", branchMatch)
		if repoMatch && branchMatch {
			return true
		}
	}
	return false
}

// LoadKeyrings loads the configured keyrings for all the commands
// managed by the worker
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
			kfile := fmt.Sprintf("%s/%s", w.Folder, keyring)
			debug.log("[worker: %s] Trying to open keyring at %s\n", w.Name, kfile)
			f, errFile := os.Open(kfile)
			if errFile != nil {
				log.Printf("[worker] cannot open keyring: %s", errFile)
				_ = f.Close()
			}

			// load the keyring
			kr, errKey := openpgp.ReadArmoredKeyRing(f)

			if errKey != nil {
				log.Printf("[worker] cannot load keyring: %s", errKey)
				_ = f.Close()
				//return fmt.Errorf("Unable to load keyring: ", err_key)
			}
			w.Keys[keyring] = kr
			w.TagKeys[t.Name][keyring] = true
			_ = f.Close()
		}
	}
	return nil
}

// LoadTags loads all the configured commands for the worker
func (w *SCORSHworker) LoadTags() error {

	wTags, err := ioutil.ReadFile(w.Tagfile)
	if err != nil {
		return fmt.Errorf("Cannot read worker config: %s", err)
	}

	err = yaml.Unmarshal(wTags, w)
	//err = yaml.Unmarshal(w_tags, tags)

	if err != nil {
		return fmt.Errorf("Error while reading tags: %s", err)
	}

	return nil
}

//
func runWorker(w *SCORSHworker) {

	var msg SCORSHmsg

	log.Printf("[worker: %s] Started\n", w.Name)

	// notify that we have been started!
	w.StatusChan <- msg

	// This is the main worker loop
	for {
		select {
		case msg = <-w.MsgChan:
			debug.log("[worker: %s] received message %s\n", w.Name, msg.ID)
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

	numWorkers := len(master.Workers)

	// We should now start each worker

	log.Printf("num_workers: %d\n", numWorkers)

	for w := 0; w < numWorkers; w++ {

		worker := &(master.Workers[w])
		// Set the Status and Msg channels
		worker.StatusChan = master.StatusChan
		worker.MsgChan = make(chan SCORSHmsg, 10)

		// Load worker tags from worker.Tagfile
		err := worker.LoadTags()
		if err != nil {
			close(worker.MsgChan)
			return fmt.Errorf("[Starting worker: %s] Unable to load tags: %s", worker.Name, err)
		}

		// Load worker keyrings -- this must be called *after* LoadTags!!!!
		err = worker.LoadKeyrings()
		if err != nil {
			close(worker.MsgChan)
			return fmt.Errorf("[Starting worker: %s] Unable to load keyrings: %s", worker.Name, err)
		}

		// Add the repos definitions to the map master.Repos
		for _, repoName := range worker.Repos {
			master.Repos[repoName] = append(master.Repos[repoName], worker)
		}
		go runWorker(worker)
		<-master.StatusChan
	}
	return nil
}
