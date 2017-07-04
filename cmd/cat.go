// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import "github.com/spf13/cobra"

const (
	catCmdFormat = `{{namespace}} {{path}}{{crlf}}{{body}}{{crlf}}`
)

func init() {
	initQueryableCommand(catCmd)
}

var catCmd = &cobra.Command{
	Use:   "cat [criteria]",
	Short: "Gets the bodies of documents based on their properties",
	Long: `Gets the bodies of documents based on their properties.
This is a special case of the "find" command with the format always set
to ` + catCmdFormat,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {
		findOutputFormat = catCmdFormat
		return findCmd.RunE(cmd, args)
	},
}
