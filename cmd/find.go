// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"os"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var (
	findShaParam       string
	findIdParam        string
	findPathParam      string
	findNamespaceParam string
	findOutputFormat   string
	findFullSha        bool
	dumpDryRun         bool
)

func init() {
	findCmd.Flags().StringVarP(&findShaParam, "sha", "s", "", "find by sha1 or sha1 prefix")
	findCmd.Flags().StringVarP(&findIdParam, "id", "i", "", "search by a document's id")
	findCmd.Flags().StringVarP(&findPathParam, "path", "p", "", "search by a document's path")
	findCmd.Flags().StringVarP(&findNamespaceParam, "namespace", "n", "", "search by a document's namespace")
	findCmd.Flags().BoolVar(&findFullSha, "full-sha", false, "Show the full sha1 hash")
	findCmd.Flags().StringVar(&findOutputFormat, "fmt", "", "Format the results of the find in a particular way")
}

var findCmd = &cobra.Command{
	Use:     "find [criteria]",
	Short:   "Find documents based on properties",
	Aliases: []string{"ls"},
	Long: `Find documents based on namespace, path, ID and SHA1.

EXAMPLES:

Find a document with an id (or list of them through stdin):

	paraphrase find -i b4e41da
	cat ids.txt | paraphrase find -i

Find a document with a path. * matches any character in a path including /:

	paraphrase find --path "*org/apache/commons/*.java
	find -name *.java | sed -e 's/^../*/' | paraphrase find -p

Find a document with a namespace. * matches any character in a namespace including /:

	paraphrase find --namespace assignment1
	paraphrase find --namespace assignment1 --path student1*.java

Find a document that matches a SHA1, it's prefix or a list in stdin:

	paraphrase find -s 5c410936339270b50362af837f8144f7775f2969
	paraphrase find -s 5c41093633

Change the format of the output.

	cat myids.txt | paraphrase cat --fmt="
		{{id}}\t{{path}}\n{{body | prefix "> "}}\r\n"

FORMATTING OPTIONS:

Variables:

	{{body}} The raw text of the content
	{{path}} The internal path of the document, looks like "/bar/bazz"
	{{namespace}} The starting namespace of the document like "foo"
	{{id}} The id of the document
	{{sha1}} SHA1 of the body
	{{date}} The date and time the document was indexed

Formatting Functions:

	{{VARIABLE | prefix ">"}} Prefixes all lines with the given text
	{{VARIABLE | head 5}} Only allow the first five lines
	{{VARIABLE | first 1024}} Get the first N bytes
	{{repeat 10 "="}} Prints the given x times e.g. "=========="
	{{crlf}} Prints a carriage return line feed CRLF i.e. "\r\n"
	{{tab}} Prints a tab character i.e. "\t"

Conversion Functions:

	{{VARIABLE | html}} Escape HTML characters
	{{VARIABLE | js}} Escape JavaScript characters
	{{VARIABLE | urlquery}} Escape for embedding in URLs

FORMATTING EXAMPLES:

Search Engine Style:

	<a href='http://localhost:8080/{{namespace | urlquery}}/{{path | urlquery}}'>{{id}}</a>{{path | html}}</a>
	<br/>
	<span class='muted'>{{sha1}}</span>
	<code>{{body | first 150 | html}}</code>
	<br/></br>

CLI Style:

	{{repeat 80 "="}}{{crlf}}{{id}} | {{path}}{{crlf}}{{sha1}}{{crlf | repeat 2}}{{body | prefix "> "}}{{crlf}}

`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		var doc paraphrase.Document

		doc.Id = findIdParam
		doc.Path = findPathParam
		doc.Sha1 = findShaParam
		doc.Namespace = findNamespaceParam

		docs, err := db.FindDocumentsLike(doc)

		if err != nil {
			return err
		}

		paraphrase.FormatDocuments(os.Stdout, docs, findOutputFormat, !findFullSha)

		return nil
	},
}
