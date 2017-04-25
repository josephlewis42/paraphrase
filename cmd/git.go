package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var (
	gitCmdPrefix string
)

func init() {
	//DbCmdAdd.Flags().BoolVarP(&addCmdRecursive, "recursive", "r", false, "adds files recursively from given folder(s)")
	DbCmdGit.Flags().StringVar(&gitCmdPrefix, "prefix", "", "adds a prefix to the loaded files")
}

var DbCmdGit = &cobra.Command{
	Use:   "git add [URL]",
	Short: "Add a document to the database or reads from stdin",
	Long: `Adds a document with the given path to the database.
if no paths are specified will read paths from stdin`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				args = append(args, scanner.Text())
			}
		}

		if len(args) == 0 {
			return errors.New("You must supply at least one path.")
		}

		for _, fp := range args {
			fmt.Printf("Adding: %s\n", fp)

			bytes, err := ioutil.ReadFile(fp)

			if err != nil {
				return err
			}

			fakePath := path.Join(addCmdPrefix, fp)

			doc, err := paraphrase.CreateDocumentFromData(fakePath, bytes)

			if err != nil {
				fmt.Printf("Error: %s", err)
				fmt.Println()
				continue
			}

			id, err := db.Insert(doc)

			if err != nil {
				return err
			}

			err = db.InsertDocumentText(id, bytes)
			if err != nil {
				return err
			}

			fmt.Printf("%s got id %d", fp, id)
			fmt.Println()
		}

		return nil
	},
}
