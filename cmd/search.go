package cmd

//
// var (
// 	resultCount int
// )
//
// func init() {
// 	RootCmd.AddCommand(cmdSearchId, cmdSearchText)
//
// 	cmdSearchId.Flags().IntVarP(&resultCount, "num", "n", 100, "number of search results")
// 	cmdSearchText.Flags().IntVarP(&resultCount, "num", "n", 100, "number of search results")
//
// }
//
// var cmdSearchId = &cobra.Command{
// 	Use:     "searchdoc docid [docid...]",
// 	Short:   "Performs a search for documents matching the one with the given id",
// 	Long:    ``,
// 	PreRunE: openDb,
// 	RunE: func(cmd *cobra.Command, args []string) error {
//
// 		docIds, err := parseDocIds(args, 1)
//
// 		if err != nil {
// 			return err
// 		}
//
// 		for _, id := range docIds {
// 			doc, matches, err := db.SearchDoc(id, resultCount)
//
// 			if err != nil {
// 				return err
// 			}
//
// 			fmt.Printf("Search results for %d (%s)", doc.Id, doc.Path)
// 			fmt.Println()
//
// 			for _, result := range matches {
//
// 				fmt.Printf("Id: %4d Matches: %3d Rank: %5.2f Path: %s", result.Doc.Id, result.Matches, result.Rank, result.Doc.Path)
// 				fmt.Println()
// 			}
//
// 		}
//
// 		return nil
//
// 	},
// }
//
// var cmdSearchText = &cobra.Command{
// 	Use:   "searchtext [TERM]",
// 	Short: "Performs a search for documents matching the given string",
// 	Long: `Performs a search for documents matching the given string.
// Note that this isn't the same as doing a full text search. The documents
// returned may not be complete and the string MUST be greater than the window
// length specified when setting up paraphrase.`,
// 	PreRunE: openDb,
// 	RunE: func(cmd *cobra.Command, args []string) error {
//
// 		if len(args) != 1 {
// 			return errors.New("You must specifiy a single search query")
// 		}
// 		norm := paraphrase.Normalize([]byte(args[0]))
// 		prints := paraphrase.FingerprintDocument(norm, 10)
//
// 		hashes := make([]uint64, 0)
// 		for _, print := range prints {
// 			hashes = append(hashes, print.Fingerprint)
// 		}
//
// 		matches, err := db.Search(hashes, resultCount)
// 		if err != nil {
// 			return err
// 		}
//
// 		fmt.Printf("Search '%s' was turned into hashes", args[0])
// 		fmt.Println()
//
// 		for _, hash := range hashes {
// 			fmt.Printf("\t%d", hash)
// 			fmt.Println()
// 		}
//
// 		fmt.Println("Results:")
//
// 		for _, result := range matches {
// 			fmt.Printf("Id: %4d Matches: %3d Rank: %3.2f Path: %s", result.Doc.Id, result.Matches, result.Rank, result.Doc.Path)
// 			fmt.Println()
// 		}
//
// 		return nil
//
// 	},
// }
