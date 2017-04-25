package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listCmdMatcher  string
	listResultCount int
)

func init() {
	RootCmd.AddCommand(cmdList)
	cmdList.Flags().StringVar(&listCmdMatcher, "match", "**", "which files to import from the source, a glob supporting ** and *")
	cmdList.Flags().IntVarP(&listResultCount, "num", "n", 100000, "number of results to return")
}

var cmdList = &cobra.Command{
	Use:     "list",
	Short:   "List the documents in the database",
	Long:    `List the documents in the database`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {
		docs, err := db.GetDocsMatching(listCmdMatcher, listResultCount)

		if err != nil {
			return err
		}

		for _, doc := range docs {
			fmt.Printf("%4d\t%s", doc.Id, doc.Path)
			fmt.Println()
		}

		return nil
	},
}
