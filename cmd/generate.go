package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	RootCmd.AddCommand(genMan)
}

var genMan = &cobra.Command{
	Use:   "genman",
	Short: "(read only) Generate the man page for paraphrase",
	Long:  `Generates man pages for parapharase`,
	Run: func(cmd *cobra.Command, args []string) {

		header := &doc.GenManHeader{
			Title:   "PARAPHRASE",
			Section: "3",
		}
		err := doc.GenManTree(cmd, header, "/tmp")
		if err != nil {
			log.Fatal(err)
		}

		// if len(args) == 0 {
		// 	fmt.Println("You must supply at least one document to hash.")
		// 	fmt.Println()
		// 	cmd.Usage()
		// 	return
		// }
		//
		// for _, path := range args {
		// 	fmt.Printf("> %s\n", path)
		// 	paraphrase.LogFingerprintFile(path)
		// }
	},
}
