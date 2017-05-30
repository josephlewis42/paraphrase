// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	// TODO fill this out so we can have better settings
	// RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes paraphrase",
	Long:  `Sets up paraphrase with some questions and answers`,
	Run: func(cmd *cobra.Command, args []string) {
		// settings := paraphrase.NewDefaultSettings()

	},
}
