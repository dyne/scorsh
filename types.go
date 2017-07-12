package main

import (
	"golang.org/x/crypto/openpgp"
)

const (
	SCORSH_ERR_NO_FILE = -(1 << iota)
	SCORSH_ERR_KEYRING
	SCORSH_ERR_NO_REPO
	SCORSH_ERR_NO_COMMIT
	SCORSH_ERR_SIGNATURE
)

// the SCORSHmsg type represents messages received from the spool and
// sent to workers
type SCORSHmsg struct {
	repo    string
	branch  string
	old_rev string
	new_rev string
}


type SCORSHcmd struct {
	URL  string
	hash string
}

type SCORSHtag struct {
	TagName  string
	Keyrings []string
	Commands []SCORSHcmd
}

// Configuration of a worker
type SCORSHworker_cfg struct {
	Name     string   `yaml:"w_name"`
	Repos    []string `yaml:"w_repos"`
	Folder   string   `yaml:"w_folder"`
	Logfile  string   `yaml:"w_logfile"`
	Tagfile  string   `yaml:"w_tagfile"`
	Keyrings []string `yaml:"w_keyrings"`
}

// State of a worker
type SCORSHworker_state struct {
	Tags map[string]SCORSHtag
	Keys map[string]openpgp.KeyRing
	Chan chan SCORSHmsg
}

// The type SCORSHworker represents the configuration and state of a
// worker
type SCORSHworker struct {
	SCORSHworker_cfg
	SCORSHworker_state
}

// Configuration of the master
type SCORSHmaster_cfg struct {
	Spooldir  string         `yaml:"s_spooldir"`
	Logfile   string         `yaml:"s_logfile"`
	LogPrefix string         `yaml:"s_logprefix"`
	Workers   []SCORSHworker `yaml:"s_workers"`
}

// State of the master
type SCORSHmaster_state struct {
	Spooler chan SCORSHmsg
	Repos   map[string][]*SCORSHworker
}

// The type SCORSHmaster represents the configuration and state of the
// master
type SCORSHmaster struct {
	SCORSHmaster_cfg
	SCORSHmaster_state
}
