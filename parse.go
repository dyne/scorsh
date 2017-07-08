package main

import (
	"fmt"
	"github.com/KatolaZ/git2go"
	"golang.org/x/crypto/openpgp"
	"os"
	"strings"
	"log"
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


// traverse all the commits between two references, looking for 
func walk_commits(msg SCORSHmsg, keyring openpgp.KeyRing) int {

	fmt.Printf("Inside parse_commits\n")

	reponame := msg.repo
	old_rev := msg.old_rev
	new_rev := msg.new_rev

	repo, err := git.OpenRepository(reponame)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening repository %s (%s)\n",
			reponame, err)
		return SCORSH_ERR_NO_REPO
	}

	old_rev_oid, err := git.NewOid(old_rev)

	oldrev_commit, err := repo.LookupCommit(old_rev_oid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Commit: %s does not exist\n", old_rev)
		return SCORSH_ERR_NO_COMMIT
	}

	new_rev_oid, err := git.NewOid(new_rev)

	newrev_commit, err := repo.LookupCommit(new_rev_oid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Commit: %s does not exist\n", new_rev)
		return SCORSH_ERR_NO_COMMIT
	}

	cur_commit := newrev_commit

	for cur_commit.Id().String() != oldrev_commit.Id().String() {

		commit, err := repo.LookupCommit(cur_commit.Id())
		if err == nil {

			fmt.Printf("%s", CommitToString(commit))
			//signature, signed, err := check_signature(commit, &keyring)
			_, _, err := check_signature(commit, &keyring)
			if err != nil {
				log.Printf("%s\n", SCORSHErr(SCORSH_ERR_SIGNATURE))
				
			}
			cur_commit = commit.Parent(0)
		} else {
			fmt.Printf("Commit %x not found!\n", cur_commit.Id())
			return SCORSH_ERR_NO_COMMIT
		}
	}
	return 0
}
