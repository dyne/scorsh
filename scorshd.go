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

var conf_file = flag.String("c", "./scorsh.cfg", "Configuration file for SCORSH")

func SCORSHerr(err int) error {

	var err_str string

	switch err {
	case SCORSH_ERR_NO_FILE:
		err_str = "Invalid file name"
	case SCORSH_ERR_KEYRING:
		err_str = "Invalid keyring"
	case SCORSH_ERR_NO_REPO:
		err_str = "Invalid repository"
	case SCORSH_ERR_NO_COMMIT:
		err_str = "Invalid commit ID"
	case SCORSH_ERR_SIGNATURE:
		err_str = "Invalid signature"
	default:
		err_str = "Generic Error"
	}
	return fmt.Errorf("%s", err_str)

}

func FindMatchingWorkers(master *SCORSHmaster, msg *SCORSHmsg) []*SCORSHworker {

	var ret []*SCORSHworker

	for idx, w := range master.Workers {
		if w.Matches(msg.Repo, msg.Branch) {
			debug.log("--- Worker: %s matches %s:%s\n", w.Name, msg.Repo, msg.Branch)
			ret = append(ret, &(master.Workers[idx]))
		}
	}
	return ret
}

func Master(master *SCORSHmaster) {

	// master main loop:

	var matching_workers []*SCORSHworker

	matching_workers = make([]*SCORSHworker, len(master.Workers))

	log.Println("[master] Master started ")
	debug.log("[master] StatusChan: %s\n", master.StatusChan)

	for {
		debug.log("[master] Receive loop...\n")
		select {
		case push_msg := <-master.Spooler:
			// here we manage the stuff we receive from the spooler
			debug.log("[master] received message: %s\n", push_msg)
			// - lookup the repos map for matching workers
			matching_workers = FindMatchingWorkers(master, &push_msg)
			debug.log("[master] matching workers: \n%s\n", matching_workers)

			// add the message to WorkingMsg, if it's not a duplicate!
			if _, ok := master.WorkingMsg[push_msg.Id]; ok {
				log.Printf("[master] detected duplicate message %s \n", push_msg.Id)
			} else {
				master.WorkingMsg[push_msg.Id] = 0
				// - dispatch the message to all the matching workers
				for _, w := range matching_workers {
					debug.log("[master] sending msg to worker: %s\n", w.Name)
					// send the message to the worker
					w.MsgChan <- push_msg
					// increase the counter associated to the message
					master.WorkingMsg[push_msg.Id] += 1
					debug.log("[master] now WorkingMsg[%s] is: %d\n", push_msg.Id, master.WorkingMsg[push_msg.Id])
				}
			}
		case done_msg := <-master.StatusChan:
			// Here we manage a status message from a worker
			debug.log("[master] received message from StatusChan: %s\n", done_msg)
			if _, ok := master.WorkingMsg[done_msg.Id]; ok && master.WorkingMsg[done_msg.Id] > 0 {
				master.WorkingMsg[done_msg.Id] -= 1
				if master.WorkingMsg[done_msg.Id] == 0 {
					delete(master.WorkingMsg, done_msg.Id)
					master.Spooler <- done_msg
				}
			} else {
				log.Printf("[master] received completion event for non-existing message name: %s\n", done_msg.Id)
			}
		}
	}
	debug.log("[master] Exiting the for loop, for some mysterious reason...\n")
}

func InitMaster() *SCORSHmaster {

	master := ReadGlobalConfig(*conf_file)

	master.Repos = make(map[string][]*SCORSHworker)
	master.WorkingMsg = make(map[string]int)
	// This is the channel on which we receive acks from workers
	master.StatusChan = make(chan SCORSHmsg)
	// This is the channel on which we exchange messages with the spooler
	master.Spooler = make(chan SCORSHmsg)

	debug.log("[InitMaster] StatusChan: %s\n", master.StatusChan)

	err_workers := StartWorkers(master)
	if err_workers != nil {
		log.Fatal("Error starting workers: ", err_workers)
	} else {
		log.Println("Workers started correctly")
	}
	err_spooler := StartSpooler(master)
	if err_spooler != nil {
		log.Fatal("Error starting spooler: ", err_spooler)
	}
	return master
}

func main() {

	var done chan int

	flag.Parse()

	master := InitMaster()

	go Master(master)

	// wait indefinitely -- we should implement signal handling...
	<-done
}
