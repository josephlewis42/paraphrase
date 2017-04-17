package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "paraphrase",
	Short: "Index text and look for duplicated content",
	Long: `Paraphrase looks for duplicated content given collections of text
good if you're looking for plagarism, suspicious copy/pasting, or links
between documents`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	RootCmd.AddCommand(DbCmdList, DbCmdGet, DbCmdHash,
		DbCmdWinnow, CmdAbout, CmdSim, CmdDumpFile, DbCmdAdd, CmdReport)
}

var DbCmdList = &cobra.Command{
	Use:   "list",
	Short: "List the ids of all documents",
	Long:  `List the ids of all documents`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here

		db, err := paraphrase.Open(".")

		if err != nil {
			panic(err)
		}

		defer db.Close()

		docs, _ := db.DocList()

		for _, doc := range docs {
			fmt.Println(doc)
		}
	},
}

var DbCmdGet = &cobra.Command{
	Use:   "get docid [docid ...]",
	Short: "(read only) Get document info for the given doc id",
	Long:  `Get document info for the given doc id`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("You must supply at least one document id.")
			fmt.Println()
			cmd.Usage()
			return
		}

		db, err := paraphrase.Open(".")

		if err != nil {
			panic(err)
		}

		defer db.Close()

		for _, docid := range args {

			id, _ := strconv.Atoi(docid)
			doc, err := db.GetDoc(uint64(id))

			if err != nil {
				fmt.Printf("Error document %s does not exist.\n", docid)
				continue
			}

			b, err := json.MarshalIndent(doc, "", "    ")

			fmt.Println(string(b))

		}

		// Do Stuff Here
	},
}

var DbCmdAdd = &cobra.Command{
	Use:   "add path [path ...]",
	Short: "Add a document to the database",
	Long:  `Adds a document with the given path to the database`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("You must supply at least one document id.")
			fmt.Println()
			cmd.Usage()
			return
		}

		db, err := paraphrase.Open(".")

		if err != nil {
			panic(err)
		}

		defer db.Close()

		if len(args) == 0 {
			fmt.Println("You must supply at least one document.")
			fmt.Println()
			cmd.Usage()
			return
		}

		for _, path := range args {
			fmt.Printf("Adding: %s\n", path)

			doc, err := paraphrase.CreateDocumentFromFile(path)

			if err != nil {
				fmt.Printf("Error: %s", err)
				fmt.Println()
				continue
			}

			id, err := db.Insert(doc)

			if err != nil {
				fmt.Printf("Error: %s", err)
				fmt.Println()
				continue
			}

			fmt.Printf("%s got id %d", path, id)
			fmt.Println()
		}

		// Do Stuff Here
	},
}

var DbCmdHash = &cobra.Command{
	Use:   "xhash path [path...]",
	Short: "(read only, debug) Print the hashes for a document",
	Long: `Calculates the hashes for the given document and prints them on the
screen. Mostly useful for testing.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("You must supply at least one document to hash.")
			fmt.Println()
			cmd.Usage()
			return
		}

		for _, path := range args {
			fmt.Printf("> %s\n", path)
			paraphrase.LogFingerprintFile(path)
		}
	},
}

var DbCmdWinnow = &cobra.Command{
	Use:   "xwinnow path [path...]",
	Short: "(read only, debug) Print the winnowed hashes",
	Long:  `Calculates the hashes for the given document and winnows them.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("You must supply at least one document to winnow.")
			fmt.Println()
			cmd.Usage()
			return
		}

		for _, path := range args {
			fmt.Printf("> %s\n", path)

			paraphrase.LogWinnowFile(path)
		}
	},
}

var CmdSim = &cobra.Command{
	Use:   "xsim path1 path2",
	Short: "(read only, debug) Calculates the similarity of two documents",
	Long:  `Calculates the similarity of two documents using winnowed hashes.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			fmt.Println("You must supply two documents to compare.")
			fmt.Println()
			cmd.Usage()
			return
		}

		paraphrase.Similarity(args[0], args[1])
	},
}

var CmdDumpFile = &cobra.Command{
	Use:   "xadd path [path...]",
	Short: "(read only, debug) Dry run of an add.",
	Long: `Prepeares a document for insertion, but prints it out rather than
adding it to the database.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("You must supply at least one document.")
			fmt.Println()
			cmd.Usage()
			return
		}

		for _, path := range args {
			fmt.Printf("Preparing: %s\n", path)

			doc, err := paraphrase.CreateDocumentFromFile(path)

			if err != nil {
				fmt.Printf("Error: %s", err)
				continue
			}

			out, _ := json.MarshalIndent(doc, "", "    ")
			fmt.Println(string(out))
		}

	},
}

var CmdAbout = &cobra.Command{
	Use:   "about",
	Short: "About this application and how it works",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`TODO add text here`)
	},
}

var CmdReport = &cobra.Command{
	Use:   "report docid [docid...]",
	Short: "Creates similarity reports for the given documents",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("You must supply at least one document.")
			fmt.Println()
			cmd.Usage()
			return
		}

		db, err := paraphrase.Open(".")

		if err != nil {
			panic(err)
		}

		defer db.Close()

		for _, docid := range args {

			id, err := strconv.Atoi(docid)

			if err != nil {
				fmt.Printf("Could not convert %s to a document id.\n", docid)
				continue
			}

			paraphrase.Report(uint64(id), db)
		}

	},
}

//
// var InfoCmd = &cobra.Command {
// 	Use:	"info"
// }
