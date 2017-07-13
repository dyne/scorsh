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
	name    string
	repo    string
	branch  string
	old_rev string
	new_rev string
}

type SCORSHcmd struct {
	URL  string `yaml:"c_url"`
	Hash string `yaml:"c_hash"`
}

type SCORSHtag struct {
	Name     string      `yaml:"t_name"`
	Keyrings []string    `yaml:"t_keyrings"`
	Commands []SCORSHcmd `yaml:"t_commands"`
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
	Tags       []SCORSHtag `yaml:"w_tags"`
	Keys       map[string]openpgp.KeyRing
	MsgChan    chan SCORSHmsg
	StatusChan chan SCORSHmsg
}

// The type SCORSHworker represents the configuration and state of a
// worker
type SCORSHworker struct {
	SCORSHworker_cfg   `yaml:",inline"`
	SCORSHworker_state `yaml:",inline"`
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
	Spooler    chan SCORSHmsg
	StatusChan chan SCORSHmsg
	Repos      map[string][]*SCORSHworker
	WorkingMsg map[string]int
}

// The type SCORSHmaster represents the configuration and state of the
// master
type SCORSHmaster struct {
	SCORSHmaster_cfg `yaml:",inline"`
	SCORSHmaster_state
}
