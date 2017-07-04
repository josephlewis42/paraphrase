package provider

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
)

func NewGitProvider(url, namespace string) (DocumentProducer, error) {
	producer := make(DocumentProducer, 5)

	directory, err := ioutil.TempDir("", "paraphrasegit")
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}

	prefixlen := len(absPath)

	log.Printf("Cloning %v to %v\n", url, absPath)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.NoRecurseSubmodules,
		Depth:             1,
	})

	if err != nil {
		return nil, err
	}

	head, _ := r.Head()
	hash := "UNKNOWNHASH"
	if head != nil {
		hash = head.Hash().String()
	}

	trimmedUrl := strings.SplitN(url, "//", 2)[1]

	if namespace == "" {
		namespace = trimmedUrl + " rev: " + hash
	}

	log.Println("Finished clone")

	go func() {
		generatePaths(absPath, namespace, true, prefixlen, producer)
		//os.RemoveAll(directory) // clean up
	}()

	return producer, err
}
