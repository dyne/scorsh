package main

import (
	"bytes"
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
	Name    string `yaml:"m_id"`
	Repo    string `yaml:"m_repo"`
	Branch  string `yaml:"m_branch"`
	Old_rev string `yaml:"m_oldrev"`
	New_rev string `yaml:"m_newrev"`
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
	Name     string      `yaml:"w_name"`
	Repos    []string    `yaml:"w_repos"`
	Folder   string      `yaml:"w_folder"`
	Logfile  string      `yaml:"w_logfile"`
	Tagfile  string      `yaml:"w_tagfile"`
	Keyrings []string    `yaml:"w_keyrings"`
	Tags     []SCORSHtag `yaml:"w_tags"`
}

// State of a worker
type SCORSHworker_state struct {
	Keys       map[string]openpgp.KeyRing
	MsgChan    chan SCORSHmsg
	StatusChan chan SCORSHmsg
}

// The type SCORSHworker represents the configuration and state of a
// worker
type SCORSHworker struct {
	SCORSHworker_cfg `yaml:",inline"`
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

func (cfg *SCORSHmaster) String() string {

	var buff bytes.Buffer

	buff.WriteString("spooldir: ")
	buff.WriteString(cfg.Spooldir)
	buff.WriteString("\nlogfile: ")
	buff.WriteString(cfg.Logfile)
	buff.WriteString("\nlogprefix: ")
	buff.WriteString(cfg.LogPrefix)
	buff.WriteString("\nWorkers: \n")

	for _, w := range cfg.Workers {
		buff.WriteString("---\n  name: ")
		buff.WriteString(w.Name)
		buff.WriteString("\n  repos: ")
		for _, r := range w.Repos {
			buff.WriteString("\n    ")
			buff.WriteString(r)
		}
		buff.WriteString("\n  folder: ")
		buff.WriteString(w.Folder)
		buff.WriteString("\n  logfile: ")
		buff.WriteString(w.Logfile)
		buff.WriteString("\n  tagfile: ")
		buff.WriteString(w.Tagfile)
		buff.WriteString("\n  keyrings: ")
		for _, k := range w.Keyrings {
			buff.WriteString("\n    ")
			buff.WriteString(k)
		}
		buff.WriteString("\n...\n")

	}

	return buff.String()
}

func (msg *SCORSHmsg) String() string {

	var buff bytes.Buffer
	buff.WriteString("\nName: ")
	buff.WriteString(msg.Name)
	buff.WriteString("\nRepo: ")
	buff.WriteString(msg.Repo)
	buff.WriteString("\nBranch: ")
	buff.WriteString(msg.Branch)
	buff.WriteString("\nOld_rev: ")
	buff.WriteString(msg.Old_rev)
	buff.WriteString("\nNew_rev: ")
	buff.WriteString(msg.New_rev)
	return buff.String()

}
