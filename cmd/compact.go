// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/boltdb/bolt"
	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
)

var compactCmd = &cobra.Command{
	Use:   "compact",
	Short: "Compacts the database keeping the original and the newly compacted one",
	Long:  `Compacts the database keeping the original and the newly compacted one`,
	RunE: func(cmd *cobra.Command, args []string) error {
		source := paraphrase.FindDbPath(projectBase)

		directory, err := ioutil.TempDir("", "compaction")
		if err != nil {
			return err
		}

		//defer os.Remove(directory)

		log.Println("Opening database")
		db, err := bolt.Open(source, 0600, nil)
		if err != nil {
			return err
		}
		defer db.Close()

		tmpFilePath := path.Join(directory, "temp.bolt")
		log.Printf("Creating temporary file %v\n", tmpFilePath)
		tmpFile, err := os.Create(tmpFilePath)
		if err != nil {
			log.Printf("Error %v\n", err)
			return err
		}

		defer tmpFile.Close()
		log.Println("Starting transaction")

		tx, err := db.Begin(false)
		if err != nil {
			return err
		}

		log.Println("Dumping file")
		_, err = tx.WriteTo(tmpFile)
		if err != nil {
			return err
		}

		tmpFile.Close()
		err = db.Close()
		if err != nil {
			return err
		}

		log.Printf("Moving %v to %v.orig\n", source, source)
		err = os.Rename(source, source+".orig")
		if err != nil {
			return err
		}
		log.Printf("Moving %v to %v\n", tmpFilePath, source)

		return os.Rename(tmpFilePath, source)

	},
}
