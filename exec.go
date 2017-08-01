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

func execTag(tag *SCORSHtagCfg, args []string, env []string) []error {

	var ret []error

	for _, c := range tag.Commands {
		debug.log("[tag: %s] attempting command: %s\n", tag.Name, c.URL)
		cmdURL, err := url.Parse(c.URL)
		if err != nil {
			log.Printf("[tag: %s] error parsing URL: %s", tag.Name, err)
		} else {
			if cmdURL.Scheme == "file" {
				err = nil
				// if a hash is specified, check that it matches
				if c.Hash != "" {
					err = checkHash(cmdURL.Path, c.Hash)
				}
				// if the hash does not match, abort the command
				if err != nil {
					log.Printf("[tag: %s] %s -- aborting command\n", tag.Name, err)
					ret = append(ret, err)
					continue
				} else {
					// finally, the command can be executed
					err = execLocalFile(cmdURL, args, env)
				}

			} else if cmdURL.Scheme == "http" || cmdURL.Scheme == "https" {
				err = execURL(cmdURL, args, env)
			}
		}
		ret = append(ret, err)
	}
	return ret
}

func setEnvironment(msg *SCORSHmsg, tag, author, committer string) []string {

	env := os.Environ()
	env = append(env, fmt.Sprintf("SCORSH_REPO=%s", msg.Repo))
	env = append(env, fmt.Sprintf("SCORSH_BRANCH=%s", msg.Branch))
	env = append(env, fmt.Sprintf("SCORSH_OLDREV=%s", msg.OldRev))
	env = append(env, fmt.Sprintf("SCORSH_NEWREV=%s", msg.NewRev))
	env = append(env, fmt.Sprintf("SCORSH_ID=%s", msg.ID))
	env = append(env, fmt.Sprintf("SCORSH_TAG=%s", tag))
	env = append(env, fmt.Sprintf("SCORSH_AUTHOR=%s", author))
	env = append(env, fmt.Sprintf("SCORSH_COMMITTER=%s", committer))

	return env
}
