// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

const ()

var (
	searchIdParam      int64
	searchDocParam     string
	searchLimit        int
	searchResultFormat string = `
ID:    {{id}}
Path:  {{path}}
SHA1:  {{sha1}}
Score: {{similarity}}

{{body | head 5 | prefix "> "}}

{{repeat 80 "-"}}
`
)

func init() {
	searchCmd.Flags().Int64VarP(&searchIdParam, "id", "i", 0, "search by a document's id")
	searchCmd.Flags().StringVarP(&searchDocParam, "file", "f", "", "search by the text in a given file")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "limit to the top n documents")
	searchCmd.Flags().StringVar(&searchResultFormat, "fmt", searchResultFormat, "The format for searching")

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
		case len(args) != 0 && searchIdParam != 0 && searchDocParam == "":
			return errors.New("You must specify exactly one query, document path or id")

		case len(args) == 1:
			results, err = db.QueryByString(args[0])

		case searchIdParam != 0:
			results, err = db.QueryById(searchIdParam)
			if err != nil {
				return err
			}

		case searchDocParam != "":
			bytes, err := ioutil.ReadFile(searchDocParam)
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

		paraphrase.FormatSearchResults(os.Stdout, results, searchResultFormat, db)

		return err
	},
}
