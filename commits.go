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

// return a list of keyring names which verify the signature of this commit
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

func exec_tag(tag SCORSHtag, valid_keys []string) error {

	return nil
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
				debug.log("validated keyrings on commit: %s\n", valid_keys)
				// 2) Try to execute each of the tag included in the message

				for _, t := range tags.Tags {
					err = exec_tag(t, valid_keys)
					if err != nil {
						log.Printf("[worker: %s] unable to execute tag: %s : %s", w.Name, t.Tag, err)
					} else {
						log.Printf("[worker: %s] tag %s executed\n", w.Name, t.Tag)
					}
				}

			}

			//signature, signed, err := check_signature(commit, &w.Keys)
			//_, _, err := check_signature(commit, w.keys)

			cur_commit = commit.Parent(0)
		} else {
			fmt.Printf("Commit %x not found!\n", cur_commit.Id())
			return SCORSHerr(SCORSH_ERR_NO_COMMIT)
		}
	}
	return nil
}
