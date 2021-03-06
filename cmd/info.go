// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Short:   "Writes general information about Paraphrase's settings and Database",
	Long:    `Writes general information about Paraphrase's settings and Database`,
	PreRunE: openDb,
	Run: func(cmd *cobra.Command, args []string) {
		db.WriteStats(os.Stdout)
	},
}
