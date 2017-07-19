package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
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

	w.Keys = make(map[string]openpgp.KeyRing, len(w.Keyrings))

	// Open the keyring files
	for _, keyring := range w.Keyrings {
		f, err_file := os.Open(keyring)

		if err_file != nil {
			log.Printf("[worker] cannot open keyring:", err_file)
			f.Close()
			return fmt.Errorf("Unable to open keyring: ", err_file)
		}

		// load the keyring
		kr, err_key := openpgp.ReadArmoredKeyRing(f)

		if err_key != nil {
			log.Printf("[worker] cannot load keyring: ", err_key)
			f.Close()
			return fmt.Errorf("Unable to load keyring: ", err_key)
		}
		w.Keys[keyring] = kr
		f.Close()
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
func Worker(w *SCORSHworker) {

	var msg SCORSHmsg

	log.Printf("[worker: %s] Started\n", w.Name)
	debug.log("[worker: %s] MsgChan: %s\n", w.Name, w.MsgChan)

	// notify that we have been started!
	w.StatusChan <- msg

	// This is the main worker loop
	for {
		select {
		case msg = <-w.MsgChan:
			debug.log("[worker: %s] received message %s\n", w.Name, msg.Id)
			// process message
			err := walk_commits(msg, w)
			if err != nil {
				log.Printf("[worker: %s] error in walk_commits: %s", err)
			}
			debug.log("[worker: %s] Received message: %s", w.Name, msg)
			debug.log("[worker: %s] StatusChan: %s\n", w.Name, w.StatusChan)
			time.Sleep(1000 * time.Millisecond)
			w.StatusChan <- msg
			debug.log("[worker: %s] Sent message back: %s", w.Name, msg)
		}
	}
}

// StartWorkers starts all the workers specified in a given
// configuration and fills in the SCORSHmaster struct
func StartWorkers(master *SCORSHmaster) error {

	num_workers := len(master.Workers)

	// We should now start each worker

	log.Printf("num_workers: %d\n", num_workers)

	for w := 0; w < num_workers; w++ {

		worker := &(master.Workers[w])
		// Set the Status and Msg channels
		worker.StatusChan = master.StatusChan
		worker.MsgChan = make(chan SCORSHmsg, 10)
		// Load worker keyrings
		err := worker.LoadKeyrings()
		if err != nil {
			close(worker.MsgChan)
			return fmt.Errorf("[Starting worker: %s] Unable to load keyrings: %s\n", worker.Name, err)
		}
		// Load worker tags from worker.Tagfile
		err = worker.LoadTags()
		if err != nil {
			close(worker.MsgChan)
			return fmt.Errorf("[Starting worker: %s] Unable to load tags: %s\n", worker.Name, err)
		}
		// Add the repos definitions to the map master.Repos
		for _, repo_name := range worker.Repos {
			master.Repos[repo_name] = append(master.Repos[repo_name], worker)
		}
		go Worker(worker)
		<-master.StatusChan
	}
	return nil
}
