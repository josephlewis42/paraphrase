package paraphrase

import (
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

// Writes the documents in fashion suitable for displaying on-screen
func FormatDocuments(w io.Writer, docs []Document, templateFormat string, shortSha bool, db *ParaphraseDb) {
	if templateFormat == "" {
		WriteDocuments(w, docs, shortSha)
		return
	}

	for _, doc := range docs {
		err := RenderDocument(templateFormat, &doc, db, nil)

		if err != nil {
			log.Println(err)
			break
		}

	}
}

func FormatSearchResults(w io.Writer, docs []SearchResult, templateFormat string, db *ParaphraseDb) {
	for _, doc := range docs {
		extraFuncs := template.FuncMap{
			"similarity": func() float64 { return doc.Similarity() },
		}

		err := RenderDocument(templateFormat, doc.Doc, db, extraFuncs)

		if err != nil {
			log.Println(err)
			break
		}

	}
}

func RenderDocument(templateFormat string, doc *Document, db *ParaphraseDb, extraFuncs template.FuncMap) error {

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"body":      func() string { doc, _ := db.FindDocumentDataById(doc.Id); return string(doc.Body) },
		"path":      func() string { return doc.Path },
		"namespace": func() string { return doc.Namespace },
		"id":        func() int64 { return doc.Id },
		"sha1":      func() string { return doc.Sha1 },
		"date":      func() time.Time { return doc.IndexDate },
		"hashes":    func() map[uint64]int16 { return doc.Hashes },

		"crlf": func() string { return "\r\n" },
		"tab":  func() string { return "\t" },

		"head":   headFunc,
		"prefix": prefixLines,
		"first":  firstFunc,
		"repeat": repeatText,
	}

	if extraFuncs != nil {
		for k, v := range extraFuncs {
			funcMap[k] = v
		}
	}

	tmpl, err := template.New("DocumentTemplate").Funcs(funcMap).Parse(templateFormat)
	if err != nil {
		return err
	}

	// Run the template to verify the output.
	return tmpl.Execute(os.Stdout, doc)
}

// prefix all lines with the given prefix.
func prefixLines(prefix, lines string) string {
	return prefix + strings.Replace(lines, "\n", "\n"+prefix, -1)
}

// gets the first lineCount number of lines of the given text.
func headFunc(lineCount int, text string) string {
	if lineCount <= 0 {
		return ""
	}

	lines := strings.Split(text, "\n")

	return strings.Join(lines[:min(lineCount, len(lines))], "\n")
}

// gets the first N bytes of the given text
func firstFunc(n int, text string) string {
	if n <= 0 {
		return ""
	}

	return text[:min(n, len(text))]
}

// repeats the text n times
func repeatText(n int, text string) string {
	out := ""

	for ; n > 0; n-- {
		out += text
	}

	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
