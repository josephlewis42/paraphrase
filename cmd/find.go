package cmd

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"

	"github.com/kennygrant/sanitize"
)

const (
	CmdCatFormat = `{{namespace}} {{path}}{{crlf}}{{body}}{{crlf}}`
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
	RootCmd.AddCommand(cmdFind)
	RootCmd.AddCommand(cmdCat)
	RootCmd.AddCommand(cmdDump)

	cmdFind.Flags().StringVarP(&findShaParam, "sha", "s", "", "find by sha1 or sha1 prefix")
	cmdFind.Flags().StringVarP(&findIdParam, "id", "i", "", "search by a document's id")
	cmdFind.Flags().StringVarP(&findPathParam, "path", "p", "", "search by a document's path")
	cmdFind.Flags().StringVarP(&findNamespaceParam, "namespace", "n", "", "search by a document's namespace")
	cmdFind.Flags().BoolVar(&findFullSha, "full-sha", false, "Show the full sha1 hash")
	cmdFind.Flags().StringVar(&findOutputFormat, "fmt", "", "Format the results of the find in a particular way")

	cmdCat.Flags().StringVarP(&findShaParam, "sha", "s", "", "find by sha1 or sha1 prefix")
	cmdCat.Flags().StringVarP(&findIdParam, "id", "i", "", "search by a document's id")
	cmdCat.Flags().StringVarP(&findPathParam, "path", "p", "", "search by a document's path")
	cmdCat.Flags().StringVarP(&findNamespaceParam, "namespace", "n", "", "search by a document's namespace")

	cmdDump.Flags().StringVarP(&findShaParam, "sha", "s", "", "find by sha1 or sha1 prefix")
	cmdDump.Flags().StringVarP(&findIdParam, "id", "i", "", "search by a document's id")
	cmdDump.Flags().StringVarP(&findPathParam, "path", "p", "", "search by a document's path")
	cmdDump.Flags().StringVarP(&findNamespaceParam, "namespace", "n", "", "search by a document's namespace")
	cmdDump.Flags().BoolVar(&dumpDryRun, "dry", false, "Do a dry run (don't create anything)")
}

var cmdFind = &cobra.Command{
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

var cmdCat = &cobra.Command{
	Use:   "cat [criteria]",
	Short: "Gets the bodies of documents based on their properties",
	Long: `Gets the bodies of documents based on their properties.
This is a special case of the "find" command with the format always set
to ` + CmdCatFormat,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {
		findOutputFormat = CmdCatFormat
		return cmdFind.RunE(cmd, args)
	},
}

var cmdDump = &cobra.Command{
	Use:     "dump [criteria] directory",
	Short:   "Writes the matching docs to a directory",
	Long:    `Writes the matching documents to a directory.`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return errors.New("You must specify one directory to write to")
		}

		var doc paraphrase.Document

		doc.Id = findIdParam
		doc.Path = findPathParam
		doc.Sha1 = findShaParam
		doc.Namespace = findNamespaceParam

		docs, err := db.FindDocumentsLike(doc)

		if err != nil {
			return err
		}

		parent := args[0]

		for _, doc := range docs {
			outpath := filepath.Join(parent, sanitize.Path(doc.Namespace), doc.Path)
			filename := filepath.Base(outpath)
			filedir := filepath.Dir(outpath)

			log.Printf("Writing %s (%s) to %s\n", doc.Id, filename, filedir)

			if dumpDryRun {
				continue
			}

			body, err := db.FindDocumentDataById(doc.Id)
			if err != nil {
				log.Printf("Error getting %s: %s\n", doc.Id, err)
				continue
			}

			// Make parent directory
			err = os.MkdirAll(filedir, 0700)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(outpath, body.Body, 0700)
			if err != nil {
				log.Println(err)
			}
		}

		return nil
	},
}
