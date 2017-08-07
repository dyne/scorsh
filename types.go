package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
)

// error constants
const (
	errNoFile = -(1 << iota)
	errKeyring
	errNoRepo
	errNoCommit
	errSignature
)

// spoolMsg type represents messages received from the spool and
// sent to workers
type spoolMsg struct {
	ID     string `yaml:"m_id"`
	Repo   string `yaml:"m_repo"`
	Branch string `yaml:"m_branch"`
	OldRev string `yaml:"m_oldrev"`
	NewRev string `yaml:"m_newrev"`
	Path   string
}

// An action represents a script of a  command configured on the server side
type action struct {
	URL  string `yaml:"a_url"`
	Hash string `yaml:"a_hash"`
}

// commandCfg represents a command configured on the server side
type commandCfg struct {
	Name     string   `yaml:"c_name"`
	Keyrings []string `yaml:"c_keyrings"`
	Actions  []action `yaml:"c_actions"`
}

// workerCfg represents the static configuration of a worker
type workerCfg struct {
	Name    string   `yaml:"w_name"`
	Repos   []string `yaml:"w_repos"`
	Folder  string   `yaml:"w_folder"`
	LogFile string   `yaml:"w_logfile"`
	CfgFile string   `yaml:"w_cfgfile"`
	//	Keyrings []string        `yaml:"w_keyrings"`
	Commands    []commandCfg `yaml:"w_commands"`
	CommandKeys map[string]map[string]bool
}

// workerState represents the runtime state of a worker
type workerState struct {
	Keys       map[string]openpgp.KeyRing
	MsgChan    chan spoolMsg
	StatusChan chan spoolMsg
}

// worker represents the configuration and state of a worker
type worker struct {
	workerCfg `yaml:",inline"`
	workerState
}

// masterCfg represents the static configuration of the master
type masterCfg struct {
	Spooldir  string   `yaml:"s_spooldir"`
	LogFile   string   `yaml:"s_logfile"`
	LogPrefix string   `yaml:"s_logprefix"`
	Workers   []worker `yaml:"s_workers"`
}

// masterState represents the runtime state of the master
type masterState struct {
	Spooler    chan spoolMsg
	StatusChan chan spoolMsg
	Repos      map[string][]*worker
	WorkingMsg map[string]int
}

// master represents the configuration and state of the master
type master struct {
	masterCfg `yaml:",inline"`
	masterState
}

// clientCmd is the type of commands sent by clients
type clientCmd struct {
	Cmd  string   `yaml:"s_cmd"`
	Args []string `yaml:"s_args"`
}

// clientMsg is the list of commands sent by a client
type clientMsg struct {
	Commands []clientCmd `yaml:"scorsh"`
}

////////////////////////

func (cfg *master) String() string {

	var buff bytes.Buffer

	fmt.Fprintf(&buff, "spooldir: %s\n", cfg.Spooldir)
	fmt.Fprintf(&buff, "logfile: %s\n", cfg.LogFile)
	fmt.Fprintf(&buff, "logprefix: %s\n", cfg.LogPrefix)
	fmt.Fprintf(&buff, "Workers: \n")

	for _, w := range cfg.Workers {
		fmt.Fprintf(&buff, "%s", &w)
	}

	return buff.String()
}

func (msg *spoolMsg) String() string {

	var buff bytes.Buffer
	fmt.Fprintf(&buff, "Id: %s\n", msg.ID)
	fmt.Fprintf(&buff, "Repo: %s\n", msg.Repo)
	fmt.Fprintf(&buff, "Branch: %s\n", msg.Branch)
	fmt.Fprintf(&buff, "OldRev: %s\n", msg.OldRev)
	fmt.Fprintf(&buff, "Newrev: %s\n", msg.NewRev)
	fmt.Fprintf(&buff, "Path: %s\n", msg.Path)

	return buff.String()

}

func (w *worker) String() string {

	var buff bytes.Buffer
	fmt.Fprintf(&buff, "Name: %s\n", w.Name)
	fmt.Fprintf(&buff, "Repos: %s\n", w.Repos)
	fmt.Fprintf(&buff, "Folder: %s\n", w.Folder)
	fmt.Fprintf(&buff, "LogFile: %s\n", w.LogFile)
	fmt.Fprintf(&buff, "CfgFile: %s\n", w.CfgFile)
	//	fmt.Fprintf(&buff, "Keyrings: %s\n", w.Keyrings)

	return buff.String()
}

func (msg *clientMsg) String() string {

	var buff bytes.Buffer

	for _, c := range msg.Commands {

		fmt.Fprintf(&buff, "s_cmd: %s\n", c.Cmd)
		for _, a := range c.Args {
			fmt.Fprintf(&buff, "  s_args: %s\n", a)
		}
	}

	return buff.String()

}
