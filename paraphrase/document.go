package paraphrase

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"text/tabwriter"
	"time"

	"github.com/josephlewis42/paraphrase/paraphrase/linalg"
)

const (
	documentFormatHeader = "ID\tSHA1\tNamespace\tPath"
	documentFormat       = "%v\t%v\t%v\t%v"
	shortShaLen          = 8
)

type TermCountVector map[uint64]int16

func (vec TermCountVector) NormalizedTermFrequency() linalg.IFVector {
	tfVector := make(linalg.IFVector)
	totalTerms := 0

	for hash, count := range vec {
		totalTerms += int(count)
		tfVector[hash] = float64(count)
	}

	tfVector.DivF(float64(totalTerms))

	return tfVector
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Document struct {
	Id        int64 `storm:"id,unique"`
	Path      string
	Namespace string
	IndexDate time.Time
	Sha1      string `storm:"index"`
	Hashes    TermCountVector
}

func (d *Document) NormalizedTermFrequency() linalg.IFVector {
	return d.Hashes.NormalizedTermFrequency()
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

	return &doc, NewDocumentData(&doc, body)
}

type DocumentData struct {
	Id        int64 `storm:"id,unique"`
	Path      string
	Namespace string
	IndexDate time.Time
	Body      []byte
}

func (dd *DocumentData) BodySha1() string {
	docHash := sha1.New()
	docHash.Write(dd.Body)
	return hex.EncodeToString(docHash.Sum(nil))
}

func NewDocumentData(doc *Document, body []byte) *DocumentData {
	var dd DocumentData

	dd.Id = doc.Id
	dd.Path = doc.Path
	dd.Namespace = doc.Namespace
	dd.IndexDate = doc.IndexDate
	dd.Body = body

	return &dd
}

func newDocId() int64 {
	v := rand.Int63()

	if v < 0 {
		return -v
	}

	return v
}
