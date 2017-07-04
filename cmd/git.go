// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"errors"

	"github.com/josephlewis42/paraphrase/paraphrase/provider"
	"github.com/spf13/cobra"
)

var (
	gitCmdNamespace string
	gitCmdMatcher   string
)

func init() {

	cmdGit.Flags().StringVar(&gitCmdNamespace, "namespace", "", "set the namespace, by default this will include the URL and revision hash")
	cmdGit.Flags().StringVarP(&gitCmdMatcher, "match", "m", WILDCARD, "only add items matching the given glob")
}

var cmdGit = &cobra.Command{
	Use:     "git [URL]",
	Short:   "Add a document to the database from a git url",
	Long:    `Adds documents to the database with the given git url`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("You must specify one git URL")
		}

		gitProvider, err := provider.NewGitProvider(args[0], gitCmdNamespace)

		if err != nil {
			return err
		}

		if gitCmdMatcher != WILDCARD {
			gitProvider, err = provider.NewFilterWrapper(gitCmdMatcher, gitProvider)

			if err != nil {
				return err
			}
		}

		db.AddDocuments(gitProvider)

		return nil
	},
}
