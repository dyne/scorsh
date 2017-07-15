package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
//	"time"
)

// parse a request file and return a SCORSHmessage
func parse_request(fname string, msg *SCORSHmsg) error {

	
	debug.log("[parse_request] message at start: %s\n", msg)
	
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Printf("Unable to open file: %s\n", fname)
		return SCORSHerr(SCORSH_ERR_NO_FILE)
	}

	debug.log("[parse_request] file contains: \n%s\n", data)
	
	debug.log("[parse_request] reading message from file: %s\n", fname)
	
	
	err = yaml.Unmarshal([]byte(data), msg)
	if err != nil {
		return fmt.Errorf("Error parsing message: %s", err)
	}
	debug.log("[parse_request] got message: %s\n", msg)
	return nil
}


func spooler(watcher *fsnotify.Watcher, worker chan SCORSHmsg) {

	log.Println("Spooler started correctly")

	var msg *SCORSHmsg
	msg = new(SCORSHmsg)
	
	for {
		select {
		case event := <-watcher.Events:
			if event.Op == fsnotify.Write {
				//time.Sleep(1000 * time.Millisecond)
				debug.log("[spooler] new file %s detected\n", event.Name)
				err := parse_request(event.Name, msg)
				if err != nil {
					log.Printf("Invalid packet received. [%s]\n", err)
				}
				debug.log("[spooler] read message: %s\n", msg)
				worker <- *msg
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}


func StartSpooler(master *SCORSHmaster) error {

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return fmt.Errorf("Error creating watcher: %s\n", err)
	}

	err = watcher.Add(master.Spooldir)
	if err != nil {
		return fmt.Errorf("Error adding folder: %s\n", err)
	}
	
	go spooler(watcher, master.Spooler)
	
	return nil
}
