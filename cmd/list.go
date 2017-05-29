package cmd

import (
	"github.com/spf13/cobra"
)

var (
	listCmdMatcher  string
	listResultCount int
)

func init() {
	RootCmd.AddCommand(cmdList)
	cmdList.Flags().StringVarP(&listCmdMatcher, "match", "m", WILDCARD, "only show items matching the given glob")
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List the documents in the database",
	Long: `List documents in the database and filter them by SHA1, name or ID.

EXAMPLES:

Search for a query:

	paraphrase find "public static void main(String[] args)"

Find documents matching stdin:

	paraphrase find < myfile.txt
	cat myfile.txt | paraphrase find

Find a document with an id (or list of them through stdin):

	paraphrase find -i b4e41da
	cat ids.txt | paraphrase find -i

Find a document with a path. * matches any character in a path including /:

	paraphrase find --path "*org/apache/commons/*.java
	find -name *.java | sed -e 's/^../*/' | paraphrase find -p


Find a document that matches a SHA1, it's prefix or a list in stdin:

	paraphrase find -s 5c410936339270b50362af837f8144f7775f2969
	paraphrase find -s 5c41093633
	sha1sum * 2>/dev/null | cut -d " " -f 1 | paraphrase find -s
`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {
		// docs, err := db.FindDocumentsMatching(listCmdMatcher)
		//
		// if err != nil {
		// 	return err
		// }
		//
		// paraphrase.WriteDocuments(os.Stdout, docs, true)
		return nil
	},
}
