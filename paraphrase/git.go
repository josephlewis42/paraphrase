package paraphrase

import (
	"fmt"
	"io/ioutil"
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

func Git(url string) error {
	directory, err := ioutil.TempDir("", "paraphrasegit")
	if err != nil {
		return err
	}

	defer os.RemoveAll(directory) // clean up

	fmt.Printf("git clone %s %s --recursive", url, directory)
	fmt.Println()

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		return err
	}

	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	fmt.Println(commit)

	return nil
}
