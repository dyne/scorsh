package main

import (
	"fmt"
	"github.com/KatolaZ/git2go"
	"github.com/go-yaml/yaml"
	"golang.org/x/crypto/openpgp"
	"log"
	"os"
	"strings"
	//	"log"
)

func CommitToString(commit *git.Commit) string {

	var ret string

	ret += fmt.Sprintf("type: %s\n", commit.Type())
	ret += fmt.Sprintf("Id: %s\n", commit.Id())
	ret += fmt.Sprintf("Author: %s\n", commit.Author())
	ret += fmt.Sprintf("Message: %s\n", commit.Message())
	ret += fmt.Sprintf("Parent-count: %d\n", commit.ParentCount())

	return ret
}

// FIXME: RETURN THE ENTITY PROVIDED BY THE CHECK, OR nil
func check_signature(commit *git.Commit, keyring *openpgp.KeyRing) (signature, signed string, err error) {

	signature, signed, err = commit.ExtractSignature()

	if err == nil {
		_, err_sig :=
			openpgp.CheckArmoredDetachedSignature(*keyring, strings.NewReader(signed),
				strings.NewReader(signature))

		if err_sig == nil {
			fmt.Printf("Good signature \n")
			return signature, signed, nil
		}
		err = err_sig
	}

	return "", "", err
}

func find_scorsh_message(commit *git.Commit) (string, error) {

	sep := "---\n"

	msg := commit.RawMessage()
	debug.log("[find_scorsg_msg] found message:\n %s\n", msg)

	// FIXME!!! replace the following with a proper regexp.Match
	idx := strings.Index(msg, sep)

	return msg[idx:], nil
}

// return a list of keyring names which verify the signature of a given commit
func get_valid_keys(commit *git.Commit, keys *map[string]openpgp.KeyRing) []string {

	var ret []string

	for k_name, k_val := range *keys {
		_, _, err := check_signature(commit, &k_val)
		if err == nil {
			ret = append(ret, k_name)
		}
	}
	return ret
}

func intersect_keys(ref map[string]bool, keys []string) []string {

	var ret []string

	for _, k := range keys {

		if _, ok := ref[k]; ok {
			ret = append(ret, k)
		}
	}
	return ret
}

func find_tag_config(tag_name string, w *SCORSHworker) (*SCORSHtag_cfg, bool) {

	for _, c := range w.Tags {
		if c.Name == tag_name {
			return &c, true
		}
	}
	return nil, false
}

// traverse all the commits between two references, looking for scorsh
// commands
// fixme: we don't have just one keyring here....
func walk_commits(msg SCORSHmsg, w *SCORSHworker) error {

	var tags SCORSHclient_msg
	var commit_msg string

	fmt.Printf("Inside parse_commits\n")

	reponame := msg.Repo
	old_rev := msg.Old_rev
	new_rev := msg.New_rev

	repo, err := git.OpenRepository(reponame)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening repository %s (%s)\n",
			reponame, err)
		return SCORSHerr(SCORSH_ERR_NO_REPO)
	}

	old_rev_oid, err := git.NewOid(old_rev)

	oldrev_commit, err := repo.LookupCommit(old_rev_oid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Commit: %s does not exist\n", old_rev)
		return SCORSHerr(SCORSH_ERR_NO_COMMIT)
	}

	new_rev_oid, err := git.NewOid(new_rev)

	newrev_commit, err := repo.LookupCommit(new_rev_oid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Commit: %s does not exist\n", new_rev)
		return SCORSHerr(SCORSH_ERR_NO_COMMIT)
	}

	cur_commit := newrev_commit

	for cur_commit.Id().String() != oldrev_commit.Id().String() {

		commit, err := repo.LookupCommit(cur_commit.Id())
		if err == nil {

			//debug.log("commit: %s", CommitToString(commit))
			// We should look for scorsh-tags, and if the commit has any,
			// check if it can be verified by any of the keyrings associated
			// with that specific scorsh-tag

			// check if the commit contains a scorsh command

			commit_msg, err = find_scorsh_message(commit)
			if err != nil {
				log.Printf("[worker: %s] %s\n", w.Name, SCORSHerr(SCORSH_ERR_SIGNATURE))
			}

			// Check if is the comment contains a valid scorsh message
			err = yaml.Unmarshal([]byte(commit_msg), &tags)

			if err != nil {
				// no scorsh message found
				log.Printf("[worker: %s] no scorsh message found: %s", err)
			} else {
				// there is a scorsh message there so

				// 1) get the list of all the keys which verify the message
				valid_keys := get_valid_keys(commit, &(w.Keys))
				debug.log("[worker: %s] validated keyrings on commit: %s\n", w.Name, valid_keys)

				// 2) then for each tag in the message
				for _, t := range tags.Tags {
					// a) check that the tag is among those accepted by the worker
					tag_cfg, good_tag := find_tag_config(t.Tag, w)
					debug.log("[worker: %s] good_tag: %s\n", w.Name, good_tag)

					if !good_tag {
						debug.log("[worker: %s] unsupported tag: %s\n", w.Name, t.Tag)
						continue
					}

					// b) check that at least one of the accepted tag keys is in valid_keys
					good_keys := intersect_keys(w.TagKeys[t.Tag], valid_keys) != nil
					debug.log("[worker: %s] good_keys: %s\n", w.Name, good_keys)

					if !good_keys {
						debug.log("[worker: %s] no matching keys for tag: %s\n", w.Name, t.Tag)
						continue
					}

					// c) If everything is OK, execute the tag
					if good_tag && good_keys {
						env := set_environment(&msg)
						errs := exec_tag(tag_cfg, t.Args, env)
						debug.log("[worker: %s] errors in tag %s: %s\n", w.Name, t.Tag, errs)
					}
				}
			}

			cur_commit = commit.Parent(0)
		} else {
			fmt.Printf("Commit %x not found!\n", cur_commit.Id())
			return SCORSHerr(SCORSH_ERR_NO_COMMIT)
		}
	}
	return nil
}
