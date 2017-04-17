package cmd

import "github.com/spf13/cobra"

var db = &cobra.Command{
	Use:   "db",
	Short: "Interact directly with the database",
	Long:  `Interact directly with the storage engine`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here

	},
}
