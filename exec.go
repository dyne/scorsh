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

func exec_local_file(cmd_url *url.URL, args, env []string) error {

	cmd := exec.Command(cmd_url.Path, args...)
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
				log.Printf(buff.Text()) // write each line to your log, or anything you need
			}
			err = cmd.Wait()
		}
	}
	return err
}

func check_hash(file, hash string) error {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	hash_bytes := sha256.Sum256(data)
	computed_hash := fmt.Sprintf("%x", string(hash_bytes[:sha256.Size]))
	debug.log("[check_hash] configured hash string: %s\n", hash)
	debug.log("[check_hash] computed hash string: %s\n", computed_hash)
	if computed_hash == hash {
		return nil
	} else {
		return fmt.Errorf("WARNING!!! HASH MISMATCH FOR %s", file)
	}

}

func exec_url(cmd_url *url.URL, args, env []string) error {

	return nil
}

func exec_tag(tag *SCORSHtag_cfg, args []string, env []string) []error {

	var ret []error

	for _, c := range tag.Commands {
		debug.log("[tag: %s] attempting command: %s\n", tag.Name, c.URL)
		cmd_url, err := url.Parse(c.URL)
		if err != nil {
			log.Printf("[tag: %s] error parsing URL: %s", tag.Name, err)
		} else {
			if cmd_url.Scheme == "file" {
				err = nil
				// if a hash is specified, check that it matches
				if c.Hash != "" {
					err = check_hash(cmd_url.Path, c.Hash)
				}
				// if the hash does not match, abort the command
				if err != nil {
					log.Printf("[tag: %s] %s -- aborting command\n", tag.Name, err)
					ret = append(ret, err)
					continue
				} else {
					// finally, the command can be executed
					err = exec_local_file(cmd_url, args, env)
				}

			} else if cmd_url.Scheme == "http" || cmd_url.Scheme == "https" {
				err = exec_url(cmd_url, args, env)
			}
		}
		ret = append(ret, err)
	}
	return ret
}

func set_environment(msg *SCORSHmsg, tag, author, committer string) []string {

	env := os.Environ()
	env = append(env, fmt.Sprintf("SCORSH_REPO=%s", msg.Repo))
	env = append(env, fmt.Sprintf("SCORSH_BRANCH=%s", msg.Branch))
	env = append(env, fmt.Sprintf("SCORSH_OLDREV=%s", msg.Old_rev))
	env = append(env, fmt.Sprintf("SCORSH_NEWREV=%s", msg.New_rev))
	env = append(env, fmt.Sprintf("SCORSH_ID=%s", msg.Id))
	env = append(env, fmt.Sprintf("SCORSH_TAG=%s", tag))
	env = append(env, fmt.Sprintf("SCORSH_AUTHOR=%s", author))
	env = append(env, fmt.Sprintf("SCORSH_COMMITTER=%s", committer))

	return env
}
