package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var (
	projectBase string
	db          *paraphrase.ParaphraseDb
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

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if db != nil {
			db.Close()
		}
	},
}

func init() {
	RootCmd.AddCommand(DbCmdList, DbCmdGet, DbCmdAdd, CmdReport)

	// commands for debugging
	RootCmd.AddCommand(CmdXNorm, CmdXAdd, CmdXSim, CmdXWinnow, CmdXHash)
	RootCmd.PersistentFlags().StringVarP(&projectBase, "base", "b", ".", "base project directory")
}

func openDb(cmd *cobra.Command, args []string) error {
	var err error
	db, err = paraphrase.Open(projectBase)
	if err != nil {
		return err
	}

	return nil
}

var DbCmdList = &cobra.Command{
	Use:     "list",
	Short:   "List the ids of all documents",
	Long:    `List the ids of all documents`,
	PreRunE: openDb,
	Run: func(cmd *cobra.Command, args []string) {
		docs, _ := db.DocList()

		for _, doc := range docs {
			fmt.Println(doc)
		}
	},
}

var DbCmdGet = &cobra.Command{
	Use:     "get DOCID [DOCID]...",
	Short:   "(read only) Get document info for the given doc id",
	Long:    `Get document info for the given doc id`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("You must supply at least one document id.")
		}

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

		return nil
	},
}

var CmdReport = &cobra.Command{
	Use:     "report docid [docid...]",
	Short:   "Creates similarity reports for the given documents",
	Long:    ``,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("You must supply at least one document.")
		}

		for _, docid := range args {

			id, err := strconv.Atoi(docid)

			if err != nil {
				fmt.Printf("Could not convert %s to a document id.\n", docid)
				continue
			}

			paraphrase.Report(uint64(id), db)
		}

		return nil
	},
}
