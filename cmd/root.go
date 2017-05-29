package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var (
	projectBase string
	db          *paraphrase.ParaphraseDb

	Version      string // Software version, auto-populated on build
	Build        string // Software build date, auto-populated on build
	Branch       string // Git branch of the build
	showAllFiles bool
	addMatcher   string
)

func init() {
	RootCmd.AddCommand(DbCmdGet, DbCmdAdd, versionCmd)
	RootCmd.AddCommand(cmdDocText)

	// commands for debugging
	// RootCmd.AddCommand(CmdXNorm, CmdXSim, CmdXWinnow, CmdXHash)

	RootCmd.PersistentFlags().StringVarP(&projectBase, "base", "b", ".", "base project directory")
}

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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints version information to stdout",
	Long:  `Prints the build, version and branch to stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s", Version)
		fmt.Println()
		fmt.Printf("Build: %s", Build)
		fmt.Println()
		fmt.Printf("Branch: %s", Branch)
		fmt.Println()

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

var DbCmdGet = &cobra.Command{
	Use:     "get DOCID [DOCID]...",
	Short:   "(read only) Get document info for the given doc id",
	Long:    `Get document info for the given doc id`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {
		docIds, err := parseDocIds(args, 1)

		if err != nil {
			return err
		}

		for _, id := range docIds {
			doc, err := db.FindDocumentById(id)

			if err != nil {
				fmt.Printf("Error document %s does not exist.\n", id)
				continue
			}

			b, err := json.MarshalIndent(doc, "", "    ")

			fmt.Println(string(b))

		}

		return nil
	},
}

var cmdDocText = &cobra.Command{
	Use:     "doctext docid [docid...]",
	Short:   "gets the text of the given document(s)",
	Long:    ``,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {
		docIds, err := parseDocIds(args, 1)

		if err != nil {
			return err
		}

		for _, id := range docIds {
			dd, err := db.FindDocumentDataById(id)

			if err != nil {
				return err
			}

			fmt.Println(string(dd.Body))
		}

		return nil
	},
}

// parseDocIds converts a list of strings to uint64s
// it will process the whole list even if an error is encountered
// so if only one element is bad, the rest will still be returned
// if the total number of successfully parse elements is less than numRequired
// an appropriate error will be returned
func parseDocIds(args []string, numRequired int) ([]string, error) {

	if len(args) < numRequired {
		if numRequired == 1 {
			return args, errors.New("You must supply at least one documents")
		} else {
			return args, errors.New(fmt.Sprintf("You must supply at least %d documents", numRequired))
		}
	}

	return args, nil
}
