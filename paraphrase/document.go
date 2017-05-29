package paraphrase

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"text/tabwriter"
	"time"
)

const (
	documentFormatHeader = "ID\tSHA1\tNamespace\tPath"
	documentFormat       = "%s\t%s\t%s\t%s"
	shortShaLen          = 8
)

type TermCountVector map[uint64]int16

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Document struct {
	Id        string `storm:"id,unique"`
	Path      string `storm:"index"`
	Namespace string `storm:"index"`
	IndexDate time.Time
	Sha1      string `storm:"index"`
	Hashes    TermCountVector
}

// Writes the documents in fashion suitable for displaying on-screen
func WriteDocuments(w io.Writer, docs []Document, shortSha bool) {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	fmt.Fprintln(tw, documentFormatHeader)
	for _, doc := range docs {

		sha := doc.Sha1
		if shortSha {
			sha = sha[0:shortShaLen]
		}

		txt := fmt.Sprintf(documentFormat, doc.Id, sha, doc.Namespace, doc.Path)
		fmt.Fprintln(tw, txt)
	}

	tw.Flush()
}

func NewDocument(path, namespace string, body []byte) (*Document, *DocumentData) {
	var doc Document

	docHash := sha1.New()
	docHash.Write(body)

	doc.Id = newDocId()
	doc.Path = path
	doc.Namespace = namespace
	doc.Sha1 = hex.EncodeToString(docHash.Sum(nil))
	doc.IndexDate = time.Now()

	return &doc, NewDocumentData(path, doc.Id, body)
}

type DocumentData struct {
	Id   string `storm:"id,unique"`
	Path string
	Body []byte
}

func NewDocumentData(path, docid string, body []byte) *DocumentData {
	var dd DocumentData

	dd.Id = docid
	dd.Path = path
	dd.Body = body

	return &dd
}

func newDocId() string {
	randint := uint64(rand.Int63())
	return fmt.Sprintf("%x", randint)
}
