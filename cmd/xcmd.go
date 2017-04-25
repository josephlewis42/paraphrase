package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var CmdXHash = &cobra.Command{
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

var CmdXWinnow = &cobra.Command{
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

var CmdXSim = &cobra.Command{
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

var CmdXNorm = &cobra.Command{
	Use:   "xnorm path [path...]",
	Short: "(read only, debug) Normalizes files like before they're processed",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("You must supply at least one document.")
			fmt.Println()
			cmd.Usage()
			return
		}

		for _, path := range args {
			fmt.Printf("> %s\n", path)
			bytes, err := ioutil.ReadFile(path)

			if err != nil {
				log.Fatal(err)
				continue
			}

			fmt.Println(string(paraphrase.Normalize(bytes)))
		}
	},
}
