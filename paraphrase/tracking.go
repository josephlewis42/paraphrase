// Copyright 2017 Joseph Lewis III <joseph@josephlewis.net>
// Licensed under the MIT License. See LICENSE file for full details.
package paraphrase

import (
	"fmt"
	"io"
	"log"
	"os/user"
	"text/tabwriter"
	"time"
)

type ChangeLogEntry struct {
	Id     int `storm:"id,increment"`
	User   string
	Date   time.Time
	Change string
}

func (p *ParaphraseDb) logChange(format string, vargs ...interface{}) {
	usr, err := user.Current()

	username := "USER NOT FOUND"
	if err == nil {
		username = usr.Name
	} else {
		log.Printf("Error getting username %v\n", err)
	}

	change := fmt.Sprintf(format, vargs...)

	cle := ChangeLogEntry{0, username, time.Now(), change}

	err = p.db.Save(&cle)
	if err != nil {
		log.Printf("Error writing changelog entry %v\n", err)
	}
}

func (p *ParaphraseDb) WriteChanges(writer io.Writer) {
	var changes []ChangeLogEntry

	p.db.Select().Find(&changes)

	w := new(tabwriter.Writer)
	w.Init(writer, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "Id\tUsername\tDate\tLog")

	for _, change := range changes {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", change.Id, change.User, change.Date, change.Change)
	}

	w.Flush()
}
