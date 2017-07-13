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
)

func (worker *SCORSHworker) Matches(repo, branch string) bool {
	
	for _, r := range worker.Repos {
		parts := strings.SplitN(r, ":", 2)
		repo_pattern := parts[0]
		branch_pattern := parts[1]
		repo_match, _ := regexp.MatchString(repo_pattern, repo)
		branch_match, _ := regexp.MatchString(branch_pattern, branch)
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
	if err != nil{
		log.Printf("[worker:%s] Cannot read worker config: ", w.Name, err)
		return err
	}
	
	err = yaml.Unmarshal(w_tags, w.Tags)

	if err != nil {
		log.Printf("[worker:%s] Error while reading tags: ", w.Name, err)
		return err
	}

	
	return nil
}

// FIXME--- STILL UNDER HEAVY WORK...
func SCORSHWorker(w *SCORSHworker) {


	// This is the main worker loop
	for {
		select {
		case msg := <-w.MsgChan:
			// process message
			err := walk_commits(msg, w)
			if err != nil {
				log.Printf("[worker: %s] error in walk_commits: %s", err)
			}
		}
	}
}

// StartWorkers starts all the workers specified in a given
// configuration and fills in the SCORSHmaster struct
func StartWorkers(master *SCORSHmaster) error {

	num_workers := len(master.Workers)
	
	// We should now start each worker

	for w:=1; w<num_workers; w++ {
		
		worker := & (master.Workers[w])
		// Set the Status and Msg channels
		worker.StatusChan = master.StatusChan
		worker.MsgChan = make(chan SCORSHmsg)
		// Load worker keyrings
		err := worker.LoadKeyrings()
		if err != nil {
			log.Printf("[worker: %s] Unable to load keyrings (Exiting): %s\n", worker.Name, err)
			close(worker.MsgChan)
			return err
		}
		// Load worker tags from worker.Tagfile
		err = worker.LoadTags()
		if err != nil {
			log.Printf("[worker: %s] Unable to load tags (Exiting): %s\n", worker.Name, err)
			close(worker.MsgChan)
			return err
		}
		// Add the repos definitions to the map master.Repos
		for _, repo_name := range worker.Repos {
			master.Repos[repo_name] = append(master.Repos[repo_name], worker)
		}
	}
	return nil
}
