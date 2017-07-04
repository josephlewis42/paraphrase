// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"errors"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

func init() {
	initQueryableCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export [criteria] [output]",
	Short: "Create a new Paraphrase DB with a subset of the documents in this one.",
	Long: `Exports a subset of documents as a new Paraphrase database.

This may be useful if you want to share portions of your database with others.

Share documents in a namespace:

	paraphrase export -n "foss" myexport.ppdb

Share documents matching a path:

	paraphrase export -p "github.com/josephlewis42/*" myexport.ppdb
`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("You must specify a new database to export to.")
		}

		dbPath := args[0]

		settings := db.GetSettings()

		exportDb, err := paraphrase.Create(dbPath, settings)

		if err != nil {
			return err
		}

		query := getQuery()

		return exportDb.ImportDocumentsMatching(db, query)
	},
}
