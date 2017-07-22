package main

import (
	"bytes"
	"fmt"
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
	Id      string `yaml:"m_id"`
	Repo    string `yaml:"m_repo"`
	Branch  string `yaml:"m_branch"`
	Old_rev string `yaml:"m_oldrev"`
	New_rev string `yaml:"m_newrev"`
	Path    string
}

type SCORSHcmd struct {
	URL  string `yaml:"c_url"`
	Hash string `yaml:"c_hash"`
}

type SCORSHtag_cfg struct {
	Name     string      `yaml:"t_name"`
	Keyrings []string    `yaml:"t_keyrings"`
	Commands []SCORSHcmd `yaml:"t_commands"`
}

// Configuration of a worker
type SCORSHworker_cfg struct {
	Name    string   `yaml:"w_name"`
	Repos   []string `yaml:"w_repos"`
	Folder  string   `yaml:"w_folder"`
	Logfile string   `yaml:"w_logfile"`
	Tagfile string   `yaml:"w_tagfile"`
	//	Keyrings []string        `yaml:"w_keyrings"`
	Tags    []SCORSHtag_cfg `yaml:"w_tags"`
	TagKeys map[string]map[string]bool
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

// client commands

type SCORSHtag struct {
	Tag  string   `yaml:"s_tag"`
	Args []string `yaml:"s_args"`
}

type SCORSHclient_msg struct {
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
	fmt.Fprintf(&buff, "Id: %s\n", msg.Id)
	fmt.Fprintf(&buff, "Repo: %s\n", msg.Repo)
	fmt.Fprintf(&buff, "Branch: %s\n", msg.Branch)
	fmt.Fprintf(&buff, "Old_Rev: %s\n", msg.Old_rev)
	fmt.Fprintf(&buff, "New_rev: %s\n", msg.New_rev)
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

func (msg *SCORSHclient_msg) String() string {

	var buff bytes.Buffer

	for _, t := range msg.Tags {

		fmt.Fprintf(&buff, "s_tag: %s\n", t.Tag)
		for _, a := range t.Args {
			fmt.Fprintf(&buff, "  s_args: %s\n", a)
		}
	}

	return buff.String()

}
