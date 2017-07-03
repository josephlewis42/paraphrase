// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.

package cmd

import (
	"strconv"

	"github.com/josephlewis42/paraphrase/paraphrase"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

const background = `
Unlike a regular search engine, Paraphrase uses "fingerprints" to look for
similarity between a search and the documents it knows about.

It does this by looking at windows of text, for example a window of length 10
produces 4 windows:

     Hello, world!
    [Hello, wor]
     [ello, worl]
      [llo, world]
       [lo, world!]

These windows get fingerprinted to form a list of document fingerprints. For
example:

	[Hello, wor]	-> 10
	 [ello, worl]	-> 45
	  [llo, world]	-> 23
	   [lo, world!]	-> 0

	[10, 45, 23, 0]

Finally, we need to choose which fingerprints to save. Again we run a window
over the fingerprints and choose the smallest in each window as "defining
characteristics" of the document. For example, a fingerprint window of size 2:

	[10, 45, 23, 0] -> [10, 45], [45, 23], [23, 0] -> [10, 23, 0]
`

func init() {
	// TODO fill this out so we can have better settings
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes paraphrase",
	Long:  `Sets up paraphrase with some questions and answers`,
	RunE: func(cmd *cobra.Command, args []string) error {

		settings := paraphrase.NewDefaultSettings()

		for {
			windowSize := ""
			survey.AskOne(&survey.Input{
				Message: "What size window would you like to use?",
				Help:    `Larger windows mean a smaller database, but requires longer query texts to get matches.`,
				Default: strconv.FormatInt(int64(settings.WindowSize), 10),
			}, &windowSize, nil)
			size, err := strconv.ParseInt(windowSize, 10, 64)
			if err != nil {
				return err
			}
			settings.WindowSize = int(size)

			kgram := ""
			survey.AskOne(&survey.Input{
				Message: "What size k-gram would you like to use?",
				Help:    `Larger k-grams mean more certainty in matches, but may miss small changes.`,
				Default: strconv.FormatInt(int64(settings.FingerprintSize), 10),
			}, &kgram, nil)
			size, err = strconv.ParseInt(kgram, 10, 64)
			if err != nil {
				return err
			}
			settings.FingerprintSize = int(size)

			robustQuestion := &survey.Confirm{
				Message: "Would you like to use robust winnowing?",
				Help: `Robust winnowing can reduce the size of the index if you're planning on indexing
data with many repeated sections but may skew search rankings slightly.`,
				Default: settings.RobustHash,
			}
			survey.AskOne(robustQuestion, &settings.RobustHash, nil)

			correct := true
			lookCorrect := &survey.Confirm{
				Message: "Do these settings look correct?",
				Default: correct,
			}
			survey.AskOne(lookCorrect, &correct, nil)

			if correct {
				break
			}
		}

		// Creates a new database in the given directory with the given settings
		_, err := paraphrase.Create(projectBase, settings)
		return err
	},
}
