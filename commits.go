package main

import (
	"fmt"
	"github.com/libgit2/git2go"
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

func findTagConfig(tagName string, w *worker) (*commandCfg, bool) {

	for _, c := range w.Tags {
		if c.Name == tagName {
			return &c, true
		}
	}
	return nil, false
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

	var commands *clientMsg

	debug.log("[worker: %s] Inside walkCommits\n", w.Name)

	reponame := msg.Repo
	oldRev := msg.OldRev
	newRev := msg.NewRev

	repo, err := git.OpenRepository(reponame)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening repository %s (%s)\n",
			reponame, err)
		return SCORSHerr(SCORSH_ERR_NO_REPO)
	}

	oldRevOid, _ := git.NewOid(oldRev)

	oldrevCommit, err := repo.LookupCommit(oldRevOid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Commit: %s does not exist\n", oldRev)
		return SCORSHerr(SCORSH_ERR_NO_COMMIT)
	}

	newRevOid, _ := git.NewOid(newRev)

	newrevCommit, err := repo.LookupCommit(newRevOid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Commit: %s does not exist\n", newRev)
		return SCORSHerr(SCORSH_ERR_NO_COMMIT)
	}

	curCommit := newrevCommit

	// FIXME: replace with a queue of commits
	for curCommit.Id().String() != oldrevCommit.Id().String() {

		commit, err := repo.LookupCommit(curCommit.Id())
		if err == nil {

			// We look for scorsh-commands, and if the commit has any, check
			// if it can be verified by any of the keyrings associated with
			// that specific scorsh-command

			// Check if the commit contains a scorsh command
			commands, err = findScorshMessage(commit)
			if err == nil {
				//  the commit contains a valid scorsh message
				// 1) get the list of all the keyrings which verify the message
				validKeys := getValidKeys(commit, &(w.Keys))
				debug.log("[worker: %s] validated keyrings on commit: %s\n", w.Name, validKeys)

				// 2) then for each tag in the message
				for _, t := range commands.Tags {
					// a) check that the tag is among those accepted by the worker
					tagCfg, goodTag := findTagConfig(t.Tag, w)
					debug.log("[worker: %s] goodTag: %s\n", w.Name, goodTag)

					if !goodTag {
						debug.log("[worker: %s] unsupported tag: %s\n", w.Name, t.Tag)
						continue
					}

					// b) check that at least one of the accepted tag keyrings
					// is in valid_keys
					goodKeys := intersectKeys(w.TagKeys[t.Tag], validKeys) != nil
					debug.log("[worker: %s] goodKeys: %s\n", w.Name, goodKeys)

					if !goodKeys {
						debug.log("[worker: %s] no matching keys for tag: %s\n", w.Name, t.Tag)
						continue
					}

					// c) If everything is OK, execute the tag
					if goodTag && goodKeys {
						env := setEnvironment(&msg, t.Tag, getAuthorEmail(commit), getCommitterEmail(commit))
						errs := execTag(tagCfg, t.Args, env)
						debug.log("[worker: %s] errors in tag %s: %s\n", w.Name, t.Tag, errs)
					}
				}
			} else {
				log.Printf("[worker: %s] error parsing commit %s: %s", w.Name, curCommit.Id().String(), err)
			}
			// FIXME: ADD ALL THE PARENTS TO THE QUEUE OF COMMITS
			curCommit = commit.Parent(0)
		} else {
			fmt.Printf("Commit %x not found!\n", curCommit.Id())
			return SCORSHerr(SCORSH_ERR_NO_COMMIT)
		}
	}
	return nil
}
