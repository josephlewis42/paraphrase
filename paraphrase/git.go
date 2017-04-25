package paraphrase

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
	git "gopkg.in/src-d/go-git.v4"
)

func Git(url string, matcher string, prefix string, db *ParaphraseDb) error {
	globby, err := glob.Compile(matcher)
	if err != nil {
		return err
	}

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

	filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		if f.Name() == ".git" {
			return filepath.SkipDir
		}

		matched := globby.Match(path)

		if !matched {
			fmt.Printf("Skipping: %s", path)
			fmt.Println()
		} else {
			fmt.Printf("Found: %s", path)
			fmt.Println()
			err := AddFile(path, prefix, db)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}

func AddFile(fp string, prefix string, db *ParaphraseDb) error {
	fmt.Printf("Adding: %s\n", fp)

	bytes, err := ioutil.ReadFile(fp)

	if err != nil {
		return err
	}

	fakePath := fp
	if prefix != "" {
		fakePath = prefix + "/" + fp
	}

	doc, err := CreateDocumentFromData(fakePath, bytes)

	if err != nil {
		return err
	}

	id, err := db.Insert(doc)

	if err != nil {
		return err
	}

	err = db.InsertDocumentText(id, bytes)
	if err != nil {
		return err
	}

	fmt.Printf("%s got id %d (internal path: %s)", fp, id, fakePath)
	fmt.Println()

	return nil
}
