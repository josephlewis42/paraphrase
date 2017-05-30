// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string // Software version, auto-populated on build
	Build   string // Software build date, auto-populated on build
	Branch  string // Git branch of the build
)

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
