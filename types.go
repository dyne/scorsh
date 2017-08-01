package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
)

// error constants
const (
	SCORSH_ERR_NO_FILE = -(1 << iota)
	SCORSH_ERR_KEYRING
	SCORSH_ERR_NO_REPO
	SCORSH_ERR_NO_COMMIT
	SCORSH_ERR_SIGNATURE
)

// SCORSHmsg type represents messages received from the spool and
// sent to workers
type SCORSHmsg struct {
	ID     string `yaml:"m_id"`
	Repo   string `yaml:"m_repo"`
	Branch string `yaml:"m_branch"`
	OldRev string `yaml:"m_oldrev"`
	NewRev string `yaml:"m_newrev"`
	Path   string
}

// SCORSHcmd represents commands configured on the server side
type SCORSHcmd struct {
	URL  string `yaml:"c_url"`
	Hash string `yaml:"c_hash"`
}

// SCORSHtagCfg represents tags configured on the server side
type SCORSHtagCfg struct {
	Name     string      `yaml:"t_name"`
	Keyrings []string    `yaml:"t_keyrings"`
	Commands []SCORSHcmd `yaml:"t_commands"`
}

// SCORSHworkerCfg represents the static configuration of a worker
type SCORSHworkerCfg struct {
	Name    string   `yaml:"w_name"`
	Repos   []string `yaml:"w_repos"`
	Folder  string   `yaml:"w_folder"`
	Logfile string   `yaml:"w_logfile"`
	Tagfile string   `yaml:"w_tagfile"`
	//	Keyrings []string        `yaml:"w_keyrings"`
	Tags    []SCORSHtagCfg `yaml:"w_tags"`
	TagKeys map[string]map[string]bool
}

// SCORSHworkerState represents the runtime state of a worker
type SCORSHworkerState struct {
	Keys       map[string]openpgp.KeyRing
	MsgChan    chan SCORSHmsg
	StatusChan chan SCORSHmsg
}

// SCORSHworker represents the configuration and state of a worker
type SCORSHworker struct {
	SCORSHworkerCfg `yaml:",inline"`
	SCORSHworkerState
}

// SCORSHmasterCfg represents the static configuration of the master
type SCORSHmasterCfg struct {
	Spooldir  string         `yaml:"s_spooldir"`
	Logfile   string         `yaml:"s_logfile"`
	LogPrefix string         `yaml:"s_logprefix"`
	Workers   []SCORSHworker `yaml:"s_workers"`
}

// SCORSHmasterState represents the runtime state of the master
type SCORSHmasterState struct {
	Spooler    chan SCORSHmsg
	StatusChan chan SCORSHmsg
	Repos      map[string][]*SCORSHworker
	WorkingMsg map[string]int
}

// SCORSHmaster represents the configuration and state of the master
type SCORSHmaster struct {
	SCORSHmasterCfg `yaml:",inline"`
	SCORSHmasterState
}

// SCORSHtag is the type of commands sent by clients
type SCORSHtag struct {
	Tag  string   `yaml:"s_tag"`
	Args []string `yaml:"s_args"`
}

// SCORSHclientMsg is the list of commands sent by a client
type SCORSHclientMsg struct {
	Tags []SCORSHtag `yaml:"scorsh"`
}

////////////////////////

func (cfg *SCORSHmaster) String() string {

	var buff bytes.Buffer

	fmt.Fprintf(&buff, "spooldir: %s\n", cfg.Spooldir)
	fmt.Fprintf(&buff, "logfile: %s\n", cfg.Logfile)
	fmt.Fprintf(&buff, "logprefix: %s\n", cfg.LogPrefix)
	fmt.Fprintf(&buff, "Workers: \n")

	for _, w := range cfg.Workers {
		fmt.Fprintf(&buff, "%s", &w)
	}

	return buff.String()
}

func (msg *SCORSHmsg) String() string {

	var buff bytes.Buffer
	fmt.Fprintf(&buff, "Id: %s\n", msg.ID)
	fmt.Fprintf(&buff, "Repo: %s\n", msg.Repo)
	fmt.Fprintf(&buff, "Branch: %s\n", msg.Branch)
	fmt.Fprintf(&buff, "OldRev: %s\n", msg.OldRev)
	fmt.Fprintf(&buff, "Newrev: %s\n", msg.NewRev)
	fmt.Fprintf(&buff, "Path: %s\n", msg.Path)

	return buff.String()

}

func (w *SCORSHworker) String() string {

	var buff bytes.Buffer
	fmt.Fprintf(&buff, "Name: %s\n", w.Name)
	fmt.Fprintf(&buff, "Repos: %s\n", w.Repos)
	fmt.Fprintf(&buff, "Folder: %s\n", w.Folder)
	fmt.Fprintf(&buff, "Logfile: %s\n", w.Logfile)
	fmt.Fprintf(&buff, "Tagfile: %s\n", w.Tagfile)
	//	fmt.Fprintf(&buff, "Keyrings: %s\n", w.Keyrings)

	return buff.String()
}

func (msg *SCORSHclientMsg) String() string {

	var buff bytes.Buffer

	for _, t := range msg.Tags {

		fmt.Fprintf(&buff, "s_tag: %s\n", t.Tag)
		for _, a := range t.Args {
			fmt.Fprintf(&buff, "  s_args: %s\n", a)
		}
	}

	return buff.String()

}
