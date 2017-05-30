// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var gendocdir string

func init() {
	gendocCmd.PersistentFlags().StringVar(&gendocdir, "dir", "docs/", "the directory to write the doc")
	gendocCmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})
}

var gendocCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate Markdown documentation for Pandoc.",
	Long: `Generate Markdown documentation for Pandoc.
It creates one Markdown file per command.`,

	RunE: func(cmd *cobra.Command, args []string) error {

		if err := os.MkdirAll(gendocdir, 0777); err != nil {
			return err
		}

		fmt.Printf("Generating markdown pages in %s", gendocdir)
		fmt.Println()
		doc.GenMarkdownTree(cmd.Root(), gendocdir)

		return nil
	},
}
