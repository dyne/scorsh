package main

import (
	"flag"
	"fmt"
	"log"
)



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
	
	for _,w := range master.Workers {
		if w.Matches(msg.repo, msg.branch) {
			ret = append(ret, &w)
		}
	}
	return ret
}


func Master(master *SCORSHmaster) {

	// master main loop:

	var matching_workers []*SCORSHworker
	var push_msg SCORSHmsg
	
	matching_workers = make([]*SCORSHworker, len(master.Workers))

	for {
		select {
		// - receive stuff from the spooler
		case push_msg = <- master.Spooler:
			// - lookup the repos map for matching workers
			matching_workers = FindMatchingWorkers(master, &push_msg)
			// - dispatch the message to all the matching workers
			for _, w := range matching_workers {
				w.Chan <- push_msg
			}
		}
	}
}

func main() {

	flag.Parse()

	master := ReadGlobalConfig(*conf_file)

	err_workers := StartWorkers(master)
	if err_workers != nil {
		log.Fatal("Error starting workers: ", err_workers)
	}
	err_spooler := StartSpooler(master)
	if err_spooler != nil {
		log.Fatal("Error starting spooler: ", err_spooler)
	}
	go Master(master)
}
