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
func FormatDocuments(w io.Writer, docs []Document, templateFormat string, shortSha bool) {
	if templateFormat == "" {
		WriteDocuments(w, docs, shortSha)
		return
	}

	for _, doc := range docs {
		err := RenderDocument(templateFormat, &doc)

		if err != nil {
			log.Println(err)
			break
		}

	}
}

func RenderDocument(templateFormat string, doc *Document) error {

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"body":      func() string { return "BODY_NOT_AVAILABLE_IN_THIS_TEMPLATE" },
		"path":      func() string { return doc.Path },
		"namespace": func() string { return doc.Namespace },
		"id":        func() string { return doc.Id },
		"sha1":      func() string { return doc.Sha1 },
		"date":      func() time.Time { return doc.IndexDate },

		"crlf": func() string { return "\r\n" },
		"tab":  func() string { return "\t" },

		"head":   headFunc,
		"prefix": prefixLines,
		"first":  firstFunc,
		"repeat": repeatText,
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

//
//
// Variables:
//
// 	{{body}} The raw text of the content
// 	{{path}} The internal path of the document, looks like "bar/bazz"
// 	{{namespace}} The starting namespace of the document like "foo"
// 	{{id}} The id of the document
// 	{{sha1}} SHA1 of the body
//
// Formatting Functions:
//
// 	{{VARIABLE | prefixlines ">"}} Prefixes all lines with the given text
// 	{{VARIABLE | head 5}} Only allow the first five lines
// 	{{VARIABLE | first 1024}} Get the first N bytes
// 	{{repeat "=" 10}} Prints the given x times e.g. "=========="
//
//
// "crlf": func() string { return "\r\n" },
// "tab":  func() string { return "\t" },

// Conversion Functions:
//
// 	{{VARIABLE | html}} Escape HTML characters
// 	{{VARIABLE | js}} Escape JavaScript characters
// 	{{VARIABLE | urlquery}} Escape for embedding in URLs
//
