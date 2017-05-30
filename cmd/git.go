// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	gitCmdPrefix  string
	gitCmdMatcher string
)

func init() {
	RootCmd.AddCommand(DbCmdGit)

	//DbCmdAdd.Flags().BoolVarP(&addCmdRecursive, "recursive", "r", false, "adds files recursively from given folder(s)")
	DbCmdGit.Flags().StringVar(&gitCmdPrefix, "prefix", "", "adds a prefix to the loaded files")
	DbCmdGit.Flags().StringVar(&gitCmdMatcher, "match", "**", "which files to import from the source, a glob supporting ** and *")
}

var DbCmdGit = &cobra.Command{
	Use:     "addgit [URL]",
	Short:   "Add a document to the database from a git url",
	Long:    `Adds documents to the database with the given git url`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("You must specify one git URL")
		}

		return nil
		// return paraphrase.Git(args[0], gitCmdMatcher, gitCmdPrefix, db)
		//
		//
		// for _, fp := range args {
		// 	fmt.Printf("Adding: %s\n", fp)
		//
		// 	bytes, err := ioutil.ReadFile(fp)
		//
		// 	if err != nil {
		// 		return err
		// 	}
		//
		// 	fakePath := path.Join(addCmdPrefix, fp)
		//
		// 	doc, err := paraphrase.CreateDocumentFromData(fakePath, bytes)
		//
		// 	if err != nil {
		// 		fmt.Printf("Error: %s", err)
		// 		fmt.Println()
		// 		continue
		// 	}
		//
		// 	id, err := db.Insert(doc)
		//
		// 	if err != nil {
		// 		return err
		// 	}
		//
		// 	err = db.InsertDocumentText(id, bytes)
		// 	if err != nil {
		// 		return err
		// 	}
		//
		// 	fmt.Printf("%s got id %d", fp, id)
		// 	fmt.Println()
		// }

		// return nil
	},
}
