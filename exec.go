package main

import (
	"fmt"
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
		err = cmd.Start()
		if err == nil {
			var output []byte
			_, err := stdout.Read(output)
			if err != nil {
				log.Printf("[%s - stout follows: ]\n%s\n", output)
				err = cmd.Wait()
			}
		}
	}

	return err
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
				err = exec_local_file(cmd_url, args, env)
			} else if cmd_url.Scheme == "http" || cmd_url.Scheme == "https" {
				err = exec_url(cmd_url, args, env)
			}
		}
		ret = append(ret, err)
	}
	return ret
}

func set_environment(msg *SCORSHmsg) []string {

	env := os.Environ()
	env = append(env, fmt.Sprintf("SCORSH_REPO=%s", msg.Repo))
	env = append(env, fmt.Sprintf("SCORSH_BRANCH=%s", msg.Branch))
	env = append(env, fmt.Sprintf("SCORSH_OLDREV=%s", msg.Old_rev))
	env = append(env, fmt.Sprintf("SCORSH_NEWREV_=%s", msg.New_rev))
	env = append(env, fmt.Sprintf("SCORSH_ID_=%s", msg.Id))
	return env
}
