// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"errors"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

func init() {
	initQueryableCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import [criteria] [source database]",
	Short: "Imports a set of documents from another Paraphrase database.",
	Long: `Imports a set of documents from another Paraphrase database.

This may be useful if you want to share portions of your database with others.

Share documents in a namespace:

	paraphrase import -n "foss" fromDb.ppdb

Share documents matching a path:

	paraphrase import -p "*.java" fromDb.ppdb
`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("You must specify a new database to import from.")
		}

		dbPath := args[0]

		importDb, err := paraphrase.Open(dbPath)

		if err != nil {
			return err
		}

		query := getQuery()

		return db.ImportDocumentsMatching(importDb, query)
	},
}
