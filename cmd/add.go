// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/josephlewis42/paraphrase/paraphrase/provider"
	"github.com/spf13/cobra"
)

const (
	WILDCARD = "*"
)

var (
	addCmdNamespace = time.Now().UTC().Format(time.RFC3339)
	addCmdDryRun    bool
	addCmdMatch     string
)

func init() {
	addCmd.Flags().StringVar(&addCmdNamespace, "namespace", addCmdNamespace, "sets the namespace of the loaded files, by default this will be a timestamp")
	addCmd.Flags().BoolVar(&addCmdDryRun, "dry", false, "list files to add rather than adding them")
	addCmd.Flags().StringVarP(&addCmdMatch, "match", "m", WILDCARD, "only add items matching the given glob")
}

var addCmd = &cobra.Command{
	Use:   "add (-|[PATH]...)",
	Short: "Add a document to the database or reads from stdin (use -)",
	Long: `Adds a document with the given path to the database.
Use add - to read from stdin.`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if len(args) == 0 {
			return errors.New("You must specify at least one file/directory or - to read from stdin")
		}

		log.Printf("Using namespace %s\n", addCmdNamespace)

		var mainProducer provider.DocumentProducer

		if len(args) == 1 && args[0] == "-" {
			mainProducer, err = provider.NewFileListProducer(addCmdNamespace, os.Stdin)
			if err != nil {
				return err
			}

		} else {
			for _, path := range args {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}

				prefixLen := len(absPath)
				if isdir, err := isDirectory(path); err == nil && !isdir {
					// The trailing separator gets removed so we subtract
					// off the length of the file from the whole path instead
					// just in case there's an OS with a funky file separator
					// pattern.
					prefixLen = len(absPath) - len(filepath.Base(absPath))
				}

				log.Printf("Searching recursively in %s\n", absPath)

				tmp := provider.NewTreeWalkerProducer(absPath, addCmdNamespace, true, prefixLen)

				mainProducer = provider.NewJoinerProducer(mainProducer, tmp)
			}
		}

		if addCmdMatch != WILDCARD {
			mainProducer, err = provider.NewFilterWrapper(addCmdMatch, mainProducer)

			if err != nil {
				return err
			}
		}

		if addCmdDryRun {
			mainProducer = provider.NewDummyProducer(mainProducer, os.Stdout)
		}

		db.AddDocuments(mainProducer)

		return nil
	},
}

func currentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
