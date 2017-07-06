package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

// parse a request file and return a SCORSHmessage
func parse_request(fname string) (SCORSHmsg, error) {

	var ret SCORSHmsg
	_, err := os.Open(fname)
	if err != nil {
		log.Printf("Unable to open file: %s\n", fname)
		return ret, SCORSHErr(SCORSH_ERR_NO_FILE)
	}
	
	return ret, nil
	
}

func spooler(watcher *fsnotify.Watcher, worker chan SCORSHmsg) {

	for {
		select {
		case event := <-watcher.Events:
			if event.Op == fsnotify.Create {
				msg, err := parse_request(event.Name)
				if err != nil {
					log.Printf("Invalid packet received. [%s]\n", err)
				}
				worker <- msg
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}


func main(){

}
