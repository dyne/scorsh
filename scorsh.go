package main

import (
	"errors"
	"log"
	"flag"
)

const (
	SCORSH_ERR_NO_FILE = -(1 << iota)
	SCORSH_ERR_KEYRING
	SCORSH_ERR_NO_REPO
	SCORSH_ERR_NO_COMMIT
	SCORSH_ERR_SIGNATURE
)

type SCORSHconf struct {
	spool string
}



type SCORSHmsg struct {
	repo    string
	branch  string
	old_rev string
	new_rev string
}

var conf_file = flag.String("c", "./scorsh.cfg", "Configuration file for SCORSH")



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



func main() {

	flag.Parse()

	cfg := ReadConfig(*conf_file)

	log.Printf("%s\n", cfg)
	
}
