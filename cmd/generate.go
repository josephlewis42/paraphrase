// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var genmandir string

// genautocomplete Generate shell autocompletion script for Hugo

var genmanCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man pages for Paraphrase",
	Long: `This command creates man pages for Paraphrase.
By default, it creates them in "./man".`,

	RunE: func(cmd *cobra.Command, args []string) error {
		header := &doc.GenManHeader{
			Section: "1",
			Manual:  "Paraphrase Manual",
			Source:  fmt.Sprintf("Paraphrase %s", Version),
		}

		if err := os.MkdirAll(genmandir, 0777); err != nil {
			return err
		}

		fmt.Printf("Generaintg man pages in %s", genmandir)
		fmt.Println()
		doc.GenManTree(cmd.Root(), header, genmandir)
		return nil
	},
}

func init() {
	genmanCmd.PersistentFlags().StringVar(&genmandir, "dir", "man/", "the directory to write the man pages.")
	genmanCmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})
}
