package provider

import (
	"regexp"
	"strings"
)

type DocumentProducer chan Document

type BodyFetcher func() ([]byte, error)

// DocumentProvider is a generic way of fetching documents from a variety of
// sources.
type Document struct {
	path      string
	namespace string
	callback  BodyFetcher
}

// Path Grabs the path of the document
func (d *Document) Path() string {
	return d.path
}

// Path Grabs the path of the document
func (d *Document) Namespace() string {
	return d.namespace
}

// Body Fetches the contents of the source. This SHOULD be lazy because
// the caller might skip processing a document
func (d *Document) Body() ([]byte, error) {
	return d.callback()
}

func NewGitProducer(gitUrl string) (DocumentProducer, error) {
	return nil, nil
}

func NewMultiJoinerProducer(producers ...DocumentProducer) DocumentProducer {
	if len(producers) == 0 {
		return nil
	}

	joiner := producers[0]

	for _, producer := range producers[1:] {
		joiner = NewJoinerProducer(joiner, producer)
	}

	return joiner
}

func NewJoinerProducer(a, b DocumentProducer) DocumentProducer {
	joiner := make(DocumentProducer, 10)

	go func() {
		defer close(joiner)

		for {
			select {
			case doc, ok := <-a:
				if ok {
					joiner <- doc
				} else {
					a = nil
				}
			case doc, ok := <-b:
				if ok {
					joiner <- doc
				} else {
					b = nil
				}
			}

			if a == nil && b == nil {
				break
			}
		}
	}()

	return joiner
}

func NewFilterWrapper(glob string, producer DocumentProducer) (DocumentProducer, error) {
	regex, err := GlobToRegex(glob)

	if err != nil {
		return nil, err
	}

	filter := make(DocumentProducer, 10)

	go func() {
		defer close(filter)

		for doc := range producer {
			if regex.Match([]byte(doc.Path())) {
				filter <- doc
			}
		}
	}()

	return filter, nil
}

func GlobToRegex(glob string) (*regexp.Regexp, error) {
	return regexp.Compile(GlobToRegexStr(glob))
}

// GlobToRegexStr converts a basic glob string to a regex
// e.g. "foo*bar.java" to "^foo.*bar\.java$"
// everything that isn't a * gets escaped
func GlobToRegexStr(glob string) string {
	split := strings.Split(glob, "*")

	for i, elem := range split {
		split[i] = regexp.QuoteMeta(elem)
	}

	return "^" + strings.Join(split, ".*") + "$"
}
