package provider

import (
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"

	git "gopkg.in/src-d/go-git.v4"
)

func NewGitProvider(gitUrl, namespace string) (DocumentProducer, error) {
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

	log.Printf("Cloning %v to %v\n", gitUrl, absPath)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               gitUrl,
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

	urlParts, _ := url.Parse(gitUrl)

	if namespace == "" {
		namespace = urlParts.Host + urlParts.Path + " rev: " + hash
	}

	log.Println("Finished clone")

	go func() {
		generatePaths(absPath, namespace, true, prefixlen, producer)
	}()

	return producer, err
}
