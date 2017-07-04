// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var (
	projectBase string
	db          *paraphrase.ParaphraseDb

	addMatcher string
	cpuprofile string
)

func init() {
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(cmdGit)

	RootCmd.AddCommand(findCmd)
	RootCmd.AddCommand(catCmd)
	RootCmd.AddCommand(dumpCmd)
	RootCmd.AddCommand(searchCmd)

	RootCmd.AddCommand(exportCmd)
	RootCmd.AddCommand(importCmd)

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(infoCmd)
	RootCmd.AddCommand(changelogCmd)
	RootCmd.AddCommand(licenseCmd)
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(compactCmd)

	GenCmd.AddCommand(genmanCmd)
	GenCmd.AddCommand(gendocCmd)
	GenCmd.AddCommand(genAutocompleteCmd)

	RootCmd.PersistentFlags().StringVarP(&projectBase, "base", "b", ".", "base project directory")
	RootCmd.PersistentFlags().StringVar(&cpuprofile, "cpuprofile", "", "write cpu profiling info to file")
	RootCmd.PersistentFlags().SetAnnotation("base", cobra.BashCompSubdirsInDir, []string{})
}

var RootCmd = &cobra.Command{
	Use:   "paraphrase",
	Short: "Index text and look for duplicated content",
	Long: `Paraphrase looks for duplicated content given collections of text
good if you're looking for plagarism, suspicious copy/pasting, or links
between documents`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cpuprofile != "" {
			f, err := os.Create(cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if db != nil {
			db.Close()
		}

		if cpuprofile != "" {
			pprof.StopCPUProfile()
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

var (
	queryableShaParam       string
	queryableIdParam        int64
	queryablePathParam      string
	queryableNamespaceParam string
)

func initQueryableCommand(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&queryableShaParam, "sha", "s", "", "find by sha1 or sha1 prefix")
	cmd.Flags().Int64VarP(&queryableIdParam, "id", "i", 0, "search by a document's id")
	cmd.Flags().StringVarP(&queryablePathParam, "path", "p", "", "search by a document's path")
	cmd.Flags().StringVarP(&queryableNamespaceParam, "namespace", "n", "", "search by a document's namespace")
}

func getQuery() paraphrase.Document {
	var query paraphrase.Document

	query.Id = queryableIdParam
	query.Path = queryablePathParam
	query.Sha1 = queryableShaParam
	query.Namespace = queryableNamespaceParam

	return query
}
