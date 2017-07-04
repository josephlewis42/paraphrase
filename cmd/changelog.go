// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:     "changelog",
	Short:   "Writes information about changes to the database.",
	Long:    `Writes information about changes to the database.`,
	PreRunE: openDb,
	Run: func(cmd *cobra.Command, args []string) {
		db.WriteChanges(os.Stdout)
	},
}
