package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
)

func execLocalFile(cmdURL *url.URL, args, env []string) error {

	cmd := exec.Command(cmdURL.Path, args...)
	cmd.Env = env
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil
	}

	if err == nil {
		if err = cmd.Start(); err == nil {
			buff := bufio.NewScanner(stdout)
			log.Printf("[%s - stout follows: ]\n", cmd.Path)
			for buff.Scan() {
				log.Print(buff.Text()) // write each line to your log, or anything you need
			}
			err = cmd.Wait()
		}
	}
	return err
}

func checkHash(file, hash string) error {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	hashBytes := sha256.Sum256(data)
	computedHash := fmt.Sprintf("%x", string(hashBytes[:sha256.Size]))
	debug.log("[checkHash] configured hash string: %s\n", hash)
	debug.log("[checkHash] computed hash string: %s\n", computedHash)
	if computedHash == hash {
		return nil
	}
	return fmt.Errorf("WARNING!!! HASH MISMATCH FOR %s", file)
}

func execURL(cmdURL *url.URL, args, env []string) error {

	return nil
}

func execCommand(cmd *commandCfg, args []string, env []string) []error {

	var ret []error

	for _, a := range cmd.Actions {
		debug.log("[command: %s] attempting action: %s\n", cmd.Name, a.URL)
		actionURL, err := url.Parse(a.URL)
		if err != nil {
			log.Printf("[command: %s] error parsing URL: %s", cmd.Name, err)
		} else {
			if actionURL.Scheme == "file" {
				err = nil
				// if a hash is specified, check that it matches
				if a.Hash != "" {
					err = checkHash(actionURL.Path, a.Hash)
				}
				// if the hash does not match, abort the command
				if err != nil {
					log.Printf("[command: %s] %s -- aborting action\n", cmd.Name, err)
					ret = append(ret, err)
					continue
				} else {
					// finally, the command can be executed
					err = execLocalFile(actionURL, args, env)
				}

			} else if actionURL.Scheme == "http" || actionURL.Scheme == "https" {
				err = execURL(actionURL, args, env)
			}
		}
		ret = append(ret, err)
	}
	return ret
}

func setEnvironment(msg *spoolMsg, cmd, author, committer string) []string {

	env := os.Environ()
	env = append(env, fmt.Sprintf("SCORSH_REPO=%s", msg.Repo))
	env = append(env, fmt.Sprintf("SCORSH_BRANCH=%s", msg.Branch))
	env = append(env, fmt.Sprintf("SCORSH_OLDREV=%s", msg.OldRev))
	env = append(env, fmt.Sprintf("SCORSH_NEWREV=%s", msg.NewRev))
	env = append(env, fmt.Sprintf("SCORSH_ID=%s", msg.ID))
	env = append(env, fmt.Sprintf("SCORSH_COMMAND=%s", cmd))
	env = append(env, fmt.Sprintf("SCORSH_AUTHOR=%s", author))
	env = append(env, fmt.Sprintf("SCORSH_COMMITTER=%s", committer))

	return env
}
