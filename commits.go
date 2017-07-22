package main

import (
	"fmt"
	"github.com/KatolaZ/git2go"
	"golang.org/x/crypto/openpgp"
	"gopkg.in/yaml.v2"
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
			debug.log("[commit: %s] Good signature \n", commit.Id())
			return signature, signed, nil
		}
		err = err_sig
	}

	return "", "", err
}

func find_scorsh_message(commit *git.Commit) (string, error) {

	sep := "---\n"

	msg := commit.RawMessage()
	debug.log("[find_scorsg_msg] found message:\n%s\n", msg)

	// FIXME!!! replace the following with a proper regexp.Match
	idx := strings.Index(msg, sep)

	if idx < 0 {
		return "", fmt.Errorf("no SCORSH message found\n")
	}
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

func get_author_email(c *git.Commit) string {

	sig := c.Author()
	return sig.Email
}

func get_committer_email(c *git.Commit) string {

	sig := c.Committer()
	return sig.Email

}

// walk_commits traverses all the commits between two references,
// looking for scorsh commands, and tries to execute those if found
func walk_commits(msg SCORSHmsg, w *SCORSHworker) error {

	var tags SCORSHclient_msg
	var commit_msg string

	debug.log("[worker: %s] Inside parse_commits\n", w.Name)

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

	// FIXME: replace with a queue of commits
	for cur_commit.Id().String() != oldrev_commit.Id().String() {

		commit, err := repo.LookupCommit(cur_commit.Id())
		if err == nil {

			// We look for scorsh-tags, and if the commit has any, check if
			// it can be verified by any of the keyrings associated with
			// that specific scorsh-tag

			// Check if the commit contains a scorsh command
			commit_msg, err = find_scorsh_message(commit)
			if err == nil {
				// Check if is the comment contains a valid scorsh message
				err = yaml.Unmarshal([]byte(commit_msg), &tags)

				if err != nil {
					// no scorsh message found
					err = fmt.Errorf("unmarshal error: %s", err)
				} else {
					// there is a scorsh message there so....

					// 1) get the list of all the keyrings which verify the message
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

						// b) check that at least one of the accepted tag keyrings
						// is in valid_keys
						good_keys := intersect_keys(w.TagKeys[t.Tag], valid_keys) != nil
						debug.log("[worker: %s] good_keys: %s\n", w.Name, good_keys)

						if !good_keys {
							debug.log("[worker: %s] no matching keys for tag: %s\n", w.Name, t.Tag)
							continue
						}

						// c) If everything is OK, execute the tag
						if good_tag && good_keys {
							env := set_environment(&msg, t.Tag, get_author_email(commit), get_committer_email(commit))
							errs := exec_tag(tag_cfg, t.Args, env)
							debug.log("[worker: %s] errors in tag %s: %s\n", w.Name, t.Tag, errs)
						}
					}
				}
			} else {
				log.Printf("[worker: %s] error parsing commit %s: %s", w.Name, cur_commit.Id().String(), err)
			}
			// FIXME: ADD ALL THE PARENTS TO THE QUEUE OF COMMITS
			cur_commit = commit.Parent(0)
		} else {
			fmt.Printf("Commit %x not found!\n", cur_commit.Id())
			return SCORSHerr(SCORSH_ERR_NO_COMMIT)
		}
	}
	return nil
}
