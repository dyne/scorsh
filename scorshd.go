package main

import (
	"flag"
	"fmt"
	"log"
)

// manage debugging messages

const debug debugging = true

type debugging bool

func (d debugging) log(format string, args ...interface{}) {
	if d {
		log.Printf(format, args...)
	}
}

///////////

var confFile = flag.String("c", "./scorsh.cfg", "Configuration file for SCORSH")

func SCORSHerr(err int) error {

	var errStr string

	switch err {
	case SCORSH_ERR_NO_FILE:
		errStr = "Invalid file name"
	case SCORSH_ERR_KEYRING:
		errStr = "Invalid keyring"
	case SCORSH_ERR_NO_REPO:
		errStr = "Invalid repository"
	case SCORSH_ERR_NO_COMMIT:
		errStr = "Invalid commit ID"
	case SCORSH_ERR_SIGNATURE:
		errStr = "Invalid signature"
	default:
		errStr = "Generic Error"
	}
	return fmt.Errorf("%s", errStr)

}

func findMatchingWorkers(master *SCORSHmaster, msg *SCORSHmsg) []*SCORSHworker {

	var ret []*SCORSHworker

	for idx, w := range master.Workers {
		if w.Matches(msg.Repo, msg.Branch) {
			debug.log("--- Worker: %s matches %s:%s\n", w.Name, msg.Repo, msg.Branch)
			ret = append(ret, &(master.Workers[idx]))
		}
	}
	return ret
}

func runMaster(master *SCORSHmaster) {

	// master main loop:

	log.Println("[master] Master started ")
	debug.log("[master] StatusChan: %s\n", master.StatusChan)

	for {
		debug.log("[master] Receive loop...\n")
		select {
		case pushMsg := <-master.Spooler:
			// here we manage the stuff we receive from the spooler
			debug.log("[master] received message: %s\n", pushMsg)
			// - lookup the repos map for matching workers
			matchingWorkers := findMatchingWorkers(master, &pushMsg)
			debug.log("[master] matching workers: \n%s\n", matchingWorkers)

			// add the message to WorkingMsg, if it's not a duplicate!
			if _, ok := master.WorkingMsg[pushMsg.Id]; ok {
				log.Printf("[master] detected duplicate message %s \n", pushMsg.Id)
			} else {
				master.WorkingMsg[pushMsg.Id] = 0
				// - dispatch the message to all the matching workers
				for _, w := range matchingWorkers {
					debug.log("[master] sending msg to worker: %s\n", w.Name)
					// send the message to the worker
					w.MsgChan <- pushMsg
					// increase the counter associated to the message
					master.WorkingMsg[pushMsg.Id]++
					debug.log("[master] now WorkingMsg[%s] is: %d\n", pushMsg.Id, master.WorkingMsg[pushMsg.Id])
				}
			}
		case doneMsg := <-master.StatusChan:
			// Here we manage a status message from a worker
			debug.log("[master] received message from StatusChan: %s\n", doneMsg)
			if _, ok := master.WorkingMsg[doneMsg.Id]; ok && master.WorkingMsg[doneMsg.Id] > 0 {
				master.WorkingMsg[doneMsg.Id]--
				if master.WorkingMsg[doneMsg.Id] == 0 {
					delete(master.WorkingMsg, doneMsg.Id)
					master.Spooler <- doneMsg
				}
			} else {
				log.Printf("[master] received completion event for non-existing message name: %s\n", doneMsg.Id)
			}
		}
	}
	debug.log("[master] Exiting the for loop, for some mysterious reason...\n")
}

func initMaster() *SCORSHmaster {

	master := readGlobalConfig(*confFile)

	master.Repos = make(map[string][]*SCORSHworker)
	master.WorkingMsg = make(map[string]int)
	// This is the channel on which we receive acks from workers
	master.StatusChan = make(chan SCORSHmsg)
	// This is the channel on which we exchange messages with the spooler
	master.Spooler = make(chan SCORSHmsg)

	debug.log("[InitMaster] StatusChan: %s\n", master.StatusChan)

	errWorkers := startWorkers(master)
	if errWorkers != nil {
		log.Fatal("Error starting workers: ", errWorkers)
	} else {
		log.Println("Workers started correctly")
	}
	errSpooler := startSpooler(master)
	if errSpooler != nil {
		log.Fatal("Error starting spooler: ", errSpooler)
	}
	return master
}

func main() {

	var done chan int

	flag.Parse()

	master := initMaster()

	go runMaster(master)

	// wait indefinitely -- we should implement signal handling...
	<-done
}
