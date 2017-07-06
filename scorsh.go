package main

import (
	"errors"
	"golang.org/x/crypto/openpgp"
	"log"
	"os"
)

const (
	SCORSH_ERR_NO_FILE = -(1 << iota)
	SCORSH_ERR_KEYRING
	SCORSH_ERR_NO_REPO
	SCORSH_ERR_NO_COMMIT
	SCORSH_ERR_SIGNATURE
)

type SCORSHmsg struct {
	repo    string
	branch  string
	old_rev string
	new_rev string
}

func SCORSHErr(err int) error {

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

	return errors.New(err_str)

}

func SCORSHWorker(keyring string, c_msg chan SCORSHmsg, c_status chan int) {

	// read the worker configuration file

	// Open the keyring file
	f, err := os.Open(keyring)
	defer f.Close()

	if err != nil {
		log.Printf("[worker] cannot open file %s\n", keyring)
		c_status <- SCORSH_ERR_NO_FILE
		return
	}

	// load the keyring
	kr, err := openpgp.ReadArmoredKeyRing(f)

	if err != nil {
		log.Printf("[worker] cannot open keyring %s\n", keyring)
		log.Printf("%s\n", err)
		c_status <- SCORSH_ERR_KEYRING
		return
	}

	// wait for messages from the  c_msg channel

	msg := <-c_msg

	// process message
	ret := walk_commits(msg, kr)

	c_status <- ret

}
