// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"errors"
	"fmt"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var (
	projectBase string
	db          *paraphrase.ParaphraseDb

	addMatcher string
)

func init() {
	RootCmd.AddCommand(addCmd)

	RootCmd.AddCommand(findCmd)
	RootCmd.AddCommand(catCmd)
	RootCmd.AddCommand(dumpCmd)

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(licenseCmd)
	RootCmd.AddCommand(GenCmd)

	GenCmd.AddCommand(genmanCmd)
	GenCmd.AddCommand(gendocCmd)
	GenCmd.AddCommand(genAutocompleteCmd)

	// commands for debugging
	// RootCmd.AddCommand(CmdXNorm, CmdXSim, CmdXWinnow, CmdXHash)

	RootCmd.PersistentFlags().StringVarP(&projectBase, "base", "b", ".", "base project directory")
	RootCmd.PersistentFlags().SetAnnotation("base", cobra.BashCompSubdirsInDir, []string{})
}

var RootCmd = &cobra.Command{
	Use:   "paraphrase",
	Short: "Index text and look for duplicated content",
	Long: `Paraphrase looks for duplicated content given collections of text
good if you're looking for plagarism, suspicious copy/pasting, or links
between documents`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if db != nil {
			db.Close()
		}
	},
}

var GenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Auto-generates Paraphrase's documentation",
	Long:  `Auto-generates man pages, markdown docs and bash autocompletion`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func openDb(cmd *cobra.Command, args []string) error {
	var err error
	db, err = paraphrase.Open(projectBase)
	if err != nil {
		return err
	}

	return nil
}

// parseDocIds converts a list of strings to uint64s
// it will process the whole list even if an error is encountered
// so if only one element is bad, the rest will still be returned
// if the total number of successfully parse elements is less than numRequired
// an appropriate error will be returned
func parseDocIds(args []string, numRequired int) ([]string, error) {

	if len(args) < numRequired {
		if numRequired == 1 {
			return args, errors.New("You must supply at least one documents")
		} else {
			return args, errors.New(fmt.Sprintf("You must supply at least %d documents", numRequired))
		}
	}

	return args, nil
}
