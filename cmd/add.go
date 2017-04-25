package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var (
	addCmdPrefix string
)

func init() {
	//DbCmdAdd.Flags().BoolVarP(&addCmdRecursive, "recursive", "r", false, "adds files recursively from given folder(s)")
	DbCmdAdd.Flags().StringVar(&addCmdPrefix, "prefix", "", "adds a prefix to the loaded files")

}

var DbCmdAdd = &cobra.Command{
	Use:   "add (-|[PATH]...)",
	Short: "Add a document to the database or reads from stdin (use -)",
	Long: `Adds a document with the given path to the database.
Use add - to read from stdin.`,
	PreRunE: openDb,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("You must specify at least one file or - to read from stdin")
		}

		if len(args) == 1 && args[0] == "-" {

			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				args = append(args, scanner.Text())
			}
			args = args[1:]
		}

		if len(args) == 0 {
			return errors.New("You must supply at least one path.")
		}

		for _, fp := range args {
			fmt.Printf("Adding: %s\n", fp)
			err := paraphrase.AddFile(fp, addCmdPrefix, db)
			if err != nil {
				return err
			}
			//
			// bytes, err := ioutil.ReadFile(fp)
			//
			// if err != nil {
			// 	return err
			// }
			//
			// fakePath := fp
			// if addCmdPrefix != "" {
			// 	fakePath = addCmdPrefix + "/" + fp
			// }
			//
			// doc, err := paraphrase.CreateDocumentFromData(fakePath, bytes)
			//
			// if err != nil {
			// 	fmt.Printf("Error: %s", err)
			// 	fmt.Println()
			// 	continue
			// }
			//
			// id, err := db.Insert(doc)
			//
			// if err != nil {
			// 	return err
			// }
			//
			// err = db.InsertDocumentText(id, bytes)
			// if err != nil {
			// 	return err
			// }
			//
			// fmt.Printf("%s got id %d (internal path: %s)", fp, id, fakePath)
			// fmt.Println()
		}

		return nil
	},
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
