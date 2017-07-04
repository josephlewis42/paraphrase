// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
)

func init() {
	initQueryableCommand(dumpCmd)
	dumpCmd.Flags().BoolVar(&dumpDryRun, "dry", false, "Do a dry run (don't create anything)")
}

var dumpCmd = &cobra.Command{
	Use:     "dump [criteria] directory",
	Short:   "Writes the matching docs to a directory",
	Long:    `Writes the matching documents to a directory.`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("You must specify one directory to write to")
		}

		doc := getQuery()

		docs, err := db.FindDocumentsLike(doc)

		if err != nil {
			return err
		}

		parent := args[0]

		for _, doc := range docs {
			outpath := filepath.Join(parent, sanitize.Path(doc.Namespace), doc.Path)
			filename := filepath.Base(outpath)
			filedir := filepath.Dir(outpath)

			log.Printf("Writing %s (%s) to %s\n", doc.Id, filename, filedir)

			if dumpDryRun {
				continue
			}

			body, err := db.FindDocumentDataById(doc.Id)
			if err != nil {
				log.Printf("Error getting %s: %s\n", doc.Id, err)
				continue
			}

			// Make parent directory
			err = os.MkdirAll(filedir, 0700)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(outpath, body.Body, 0700)
			if err != nil {
				log.Println(err)
			}
		}

		return nil
	},
}
