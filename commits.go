package main

import (
	"fmt"
	"github.com/dyne/git2go.v26"
	"golang.org/x/crypto/openpgp"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
	//	"log"
)

// func commitToString(commit *git.Commit) string {
// 	var ret string
// 	ret += fmt.Sprintf("type: %s\n", commit.Type())
// 	ret += fmt.Sprintf("Id: %s\n", commit.Id())
// 	ret += fmt.Sprintf("Author: %s\n", commit.Author())
// 	ret += fmt.Sprintf("Message: %s\n", commit.Message())
// 	ret += fmt.Sprintf("Parent-count: %d\n", commit.ParentCount())
// 	return ret
// }

// FIXME: RETURN THE ENTITY PROVIDED BY THE CHECK, OR nil
func checkSignature(commit *git.Commit, keyring *openpgp.KeyRing) (signature, signed string, err error) {

	signature, signed, err = commit.ExtractSignature()

	if err == nil {
		_, errSig :=
			openpgp.CheckArmoredDetachedSignature(*keyring, strings.NewReader(signed),
				strings.NewReader(signature))

		if errSig == nil {
			debug.log("[commit: %s] Good signature \n", commit.Id())
			return signature, signed, nil
		}
		err = errSig
	}

	return "", "", err
}

func findScorshMessage(commit *git.Commit) (*clientMsg, error) {

	var commands = new(clientMsg)
	sep := "---\n"

	msg := commit.RawMessage()
	debug.log("[findScorshMessage] found message:\n%s\n", msg)

	// FIXME!!! replace the following with a proper regexp.Match
	idx := strings.Index(msg, sep)

	if idx < 0 {
		return nil, fmt.Errorf("no SCORSH message found")
	}

	err := yaml.Unmarshal([]byte(msg[idx:]), &commands)

	if err != nil {
		// no scorsh message found
		err = fmt.Errorf("unmarshal error: %s", err)
		commands = nil
	} else {
		err = nil
	}

	return commands, nil
}

// return a list of keyring names which verify the signature of a given commit
func getValidKeys(commit *git.Commit, keys *map[string]openpgp.KeyRing) []string {

	var ret []string

	for kname, kval := range *keys {
		_, _, err := checkSignature(commit, &kval)
		if err == nil {
			ret = append(ret, kname)
		}
	}
	return ret
}

func intersectKeys(ref map[string]bool, keys []string) []string {

	var ret []string

	for _, k := range keys {

		if _, ok := ref[k]; ok {
			ret = append(ret, k)
		}
	}
	return ret
}

func findCmdConfig(cmdName string, w *worker) (commandCfg, bool) {

	var cmdNull commandCfg

	for _, c := range w.Commands {
		if c.Name == cmdName {
			return c, true
		}
	}
	return cmdNull, false
}

func getAuthorEmail(c *git.Commit) string {

	sig := c.Author()
	return sig.Email
}

func getCommitterEmail(c *git.Commit) string {

	sig := c.Committer()
	return sig.Email

}

// walkCommits traverses all the commits between two references,
// looking for scorsh commands, and tries to execute those if found
func walkCommits(msg spoolMsg, w *worker) error {

	var cmdMsg *clientMsg
	var cmdStack = make([]command, 0)

	debug.log("[worker: %s] Inside walkCommits\n", w.Name)

	reponame := msg.Repo
	oldRev := msg.OldRev
	newRev := msg.NewRev

	repo, err := git.OpenRepository(reponame)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening repository %s (%s)\n",
			reponame, err)
		return SCORSHerr(errNoRepo)
	}

	oldRevOid, _ := git.NewOid(oldRev)

	oldrevCommit, err := repo.LookupCommit(oldRevOid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Commit: %s does not exist\n", oldRev)
		return fmt.Errorf("%s: %s", SCORSHerr(errNoCommit), oldRev)
	}

	newRevOid, _ := git.NewOid(newRev)

	newrevCommit, err := repo.LookupCommit(newRevOid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Commit: %s does not exist\n", newRev)
		return fmt.Errorf("%s: %s", SCORSHerr(errNoCommit), newRev)
	}

	curCommit := newrevCommit

	// FIXME: replace with a queue of commits
	for curCommit.Id().String() != oldrevCommit.Id().String() {

		commit, err := repo.LookupCommit(curCommit.Id())
		if err == nil {

			// We look for scorsh-commands, and if the commit has any, check
			// if it can be verified by any of the keyrings associated with
			// that specific scorsh-command

			// Check if the commit contains a scorsh message
			cmdMsg, err = findScorshMessage(commit)
			if err == nil {
				//  the commit contains a valid scorsh message
				// 1) get the list of all the keyrings which verify the message signature
				validKeys := getValidKeys(commit, &(w.Keys))
				debug.log("[worker: %s] validated keyrings on commit: %s\n", w.Name, validKeys)

				// 2) then for each command in the message
				for _, c := range cmdMsg.Commands {
					if c.Cmd == "" {
						// The command is empty -- ignore, log, and continue
						log.Printf("[worker: %s] empty command\n", w.Name)
						continue
					}
					// a) check that the command is among those accepted by the worker
					debug.log("[worker: %s] validating command: %s\n", w.Name, c.Cmd)
					var cmd = new(command)
					var goodCmd, goodKeys bool
					cmd.commandCfg, goodCmd = findCmdConfig(c.Cmd, w)
					debug.log("[worker: %s] goodCmd: %s\n", w.Name, goodCmd)

					if !goodCmd {
						debug.log("[worker: %s] unsupported command: %s\n", w.Name, c.Cmd)
						continue
					}

					// b) check that at least one of the accepted command keyrings
					// is in valid_keys
					goodKeys = intersectKeys(w.CommandKeys[c.Cmd], validKeys) != nil
					debug.log("[worker: %s] goodKeys: %s\n", w.Name, goodKeys)

					if !goodKeys {
						debug.log("[worker: %s] no matching keys for command: %s\n", w.Name, c.Cmd)
						continue
					}

					// c) If everything is OK, push the command to the stack
					if goodCmd && goodKeys {
						cmd.setEnvironment(&msg, curCommit.Id().String(), getAuthorEmail(commit), getCommitterEmail(commit))
						cmd.Args = c.Args
						cmdStack = append(cmdStack, *cmd)
					}
				}
			} else {
				log.Printf("[worker: %s] error parsing commit %s: %s", w.Name, curCommit.Id().String(), err)
			}
			// FIXME: ADD ALL THE PARENTS TO THE QUEUE OF COMMITS
			curCommit = commit.Parent(0)
		} else {
			fmt.Printf("Commit %x not found!\n", curCommit.Id())
			return SCORSHerr(errNoCommit)
		}

	}

	// Now we can execute the commands in the stack, in the correct order...
	stackHead := len(cmdStack) - 1
	debug.log("[worker: %s] Executing command stack:\n", w.Name)
	for i := range cmdStack {
		//debug.log("[stack elem: %d] %s\n", i, cmdStack[stackHead-i].String())
		// now we execute the command that emerges from the stack
		cmd := cmdStack[stackHead-i]
		errs := cmd.exec()
		debug.log("[worker: %s] errors in command %s: %s\n", w.Name, cmd.Name, errs)
	}

	return nil
}
