package provider

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func NewTreeWalkerProducer(directory, namespace string, ignoreHidden bool, prefixlen int) DocumentProducer {
	producer := make(DocumentProducer, 5)

	go generatePaths(directory, namespace, ignoreHidden, prefixlen, producer)

	return producer
}

func generatePaths(root, namespace string, ignoreHidden bool, prefixlen int, output DocumentProducer) {
	defer close(output)

	if !isDirectory(root) {
		output <- Document{root[prefixlen:], namespace, readFileCallback(root)}
		return
	}

	filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if ignoreHidden && strings.HasPrefix(f.Name(), ".") {
			if f.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if f.IsDir() {
			return nil
		}

		output <- Document{path[prefixlen:], namespace, readFileCallback(path)}

		return nil
	})
}

func readFileCallback(path string) BodyFetcher {
	return func() ([]byte, error) {
		return ioutil.ReadFile(path)
	}
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
