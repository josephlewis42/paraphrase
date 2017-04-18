package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/josephlewis42/paraphrase/cmd"
)

const (
	WINDOW_SIZE      = 10
	FINGERPRINT_SIZE = 10
)

var (
	whitespace = regexp.MustCompile(`\s*`)
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
