// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var (
	searchIdParam  bool
	searchDocParam bool
	searchLimit    int
)

func init() {
	searchCmd.Flags().BoolVarP(&searchIdParam, "id", "i", false, "search by a document's id")
	searchCmd.Flags().BoolVarP(&searchDocParam, "file", "f", false, "search by the text in a given file")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "limit to the top n documents")

}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search documents matching a query.",
	Long: `Finds documents similar to the one with the given ID, text of a file or query.

EXAMPLES:

Search for documents similar to the one with the given id:

	paraphrase search -i b4e41da

Search for documents matching a query:

	paraphrase search "Hello, world!"

Search for documents similar to the given file:

	paraphrase search -f MyApplication.java

Formatting the search output:

	paraphrase search --fmt="{{id}}\t{{path}}\n{{body | prefix "> "}}\r\n"

` + FormattingOptions,
	Aliases: []string{"q"},
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {
		var results []paraphrase.SearchResult
		var err error

		switch {
		case len(args) == 0:
			return errors.New("You must a specify a query/doc/id")

		case len(args) > 1:
			return errors.New("You must specify only one query/doc/id")

		case searchIdParam && searchDocParam:
			return errors.New("You cannot use both the ID and document flags at the same time.")

		case searchIdParam:
			results, err = db.QueryById(args[0])
			if err != nil {
				return err
			}

		case searchDocParam:
			bytes, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}
			results, err = db.QueryByString(string(bytes))

		default:
			fmt.Println(args[0])
			results, err = db.QueryByString(args[0])
		}

		if err != nil {
			return err
		}

		for _, res := range results {
			fmt.Printf("Result: %s %s %s %f\n", res.Doc.Id, res.Doc.Namespace, res.Doc.Path, res.Similarity())
		}

		return err
	},
}
