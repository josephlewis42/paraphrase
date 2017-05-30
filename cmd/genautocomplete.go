// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var genAutocompleteDirectory string

var genAutocompleteCmd = &cobra.Command{
	Use:   "autocomplete",
	Short: "Generate a bash autocompletion script for Paraphrase",
	Long: `Generate a bash autocompletion script for Paraphrase.

By default, the file is written directly to /etc/bash_completion.d.
for convenience, and the command may need superuser rights, e.g.:

	$ sudo paraphrase gen autocomplete

The default location can be changed using using the --completionfile flag.`,

	RunE: func(cmd *cobra.Command, args []string) error {

		err := cmd.Root().GenBashCompletionFile(genAutocompleteDirectory)

		if err != nil {
			return err
		}

		fmt.Println("Bash completion file for Paraphrase saved to", genAutocompleteDirectory)

		return nil
	},
}

func init() {
	genAutocompleteCmd.PersistentFlags().StringVarP(&genAutocompleteDirectory, "completionfile", "", "/etc/bash_completion.d/paraphrase.sh", "autocompletion file")

	// For bash-completion
	genAutocompleteCmd.PersistentFlags().SetAnnotation("completionfile", cobra.BashCompFilenameExt, []string{})
}
