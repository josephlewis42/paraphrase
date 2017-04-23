package main

import (
	"fmt"
	"os"

	"github.com/josephlewis42/paraphrase/cmd"
)

var (
	Version string // Software version, auto-populated on build
	Build   string // Software build date, auto-populated on build
	Branch  string // Git branch of the build
)

func main() {
	cmd.Version = Version
	cmd.Build = Build
	cmd.Branch = Branch

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
