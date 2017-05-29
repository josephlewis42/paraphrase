package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/josephlewis42/paraphrase/paraphrase/provider"
	"github.com/spf13/cobra"
)

const (
	WILDCARD = "*"
)

var (
	addCmdNamespace = time.Now().UTC().Format(time.RFC3339)
	addCmdDryRun    bool
	addCmdMatch     string
)

func init() {
	DbCmdAdd.Flags().StringVar(&addCmdNamespace, "namespace", addCmdNamespace, "sets the namespace of the loaded files, by default this will be a timestamp")
	DbCmdAdd.Flags().BoolVar(&addCmdDryRun, "dry", false, "list files to add rather than adding them")
	DbCmdAdd.Flags().StringVarP(&addCmdMatch, "match", "m", WILDCARD, "only add items matching the given glob")
}

var DbCmdAdd = &cobra.Command{
	Use:   "add (-|[PATH]...)",
	Short: "Add a document to the database or reads from stdin (use -)",
	Long: `Adds a document with the given path to the database.
Use add - to read from stdin.`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if len(args) == 0 {
			return errors.New("You must specify at least one file/directory or - to read from stdin")
		}

		log.Printf("Using namespace %s\n", addCmdNamespace)

		var mainProducer provider.DocumentProducer

		if len(args) == 1 && args[0] == "-" {
			mainProducer, err = provider.NewFileListProducer(addCmdNamespace, os.Stdin)
			if err != nil {
				return err
			}

		} else {
			for _, path := range args {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}

				tmp := provider.NewTreeWalkerProducer(absPath, addCmdNamespace, true, len(absPath))

				mainProducer = provider.NewJoinerProducer(mainProducer, tmp)
			}
		}

		if addCmdMatch != WILDCARD {
			mainProducer, err = provider.NewFilterWrapper(addCmdMatch, mainProducer)

			if err != nil {
				return err
			}
		}

		if addCmdDryRun {
			mainProducer = provider.NewDummyProducer(mainProducer, os.Stdout)
		}

		db.AddDocuments(mainProducer)

		return nil
	},
}

func currentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
