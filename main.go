package main

import (
	"fmt"
	"os"

	"github.com/josephlewis42/paraphrase/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
