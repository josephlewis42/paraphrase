package paraphrase

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"regexp"
	"time"
)

const (
	WINDOW_SIZE      = 30
	FINGERPRINT_SIZE = 10
)

var (
	whitespace = regexp.MustCompile(`\s*`)
)

func Normalize(document []byte) []byte {
	var empty []byte
	return whitespace.ReplaceAll(document, empty)
}

type Fingerprint struct {
	Fingerprint uint64
	DocumentId  string
	Location    int
}

func FingerprintDocument(document []byte, size int) []Fingerprint {
	fingerprintCount := len(document) - size
	var fingerprints []Fingerprint

	for i := 0; i <= fingerprintCount; i++ {
		var fingerprint Fingerprint
		hash := fnv.New64()
		hash.Write(document[i : i+size])
		fingerprint.Fingerprint = hash.Sum64()
		fingerprint.Location = i
		fingerprints = append(fingerprints, fingerprint)
	}

	return fingerprints
}

func Winnow(fingerprints []Fingerprint, window int, robust bool) []Fingerprint {
	var recorded []Fingerprint

	h := make([]uint64, window)

	for i, _ := range h {
		h[i] = math.MaxUint64
	}

	r := 0   // window right end
	min := 0 // index of min hash

	for _, fingerprint := range fingerprints {
		r = (r + 1) % window // shift window by one
		h[r] = fingerprint.Fingerprint

		if min == r {
			// previous minimum is no longer in this window.
			// scan h leftward starting from r for the rightmost minimal hash.
			// Note min starts with the index of the rightmost hash.
			for i := (r - 1 + window) % window; i != r; i = (i - 1 + window) % window {
				if h[i] < h[min] {
					min = i
				}
			}

			recorded = append(recorded, fingerprint)

		} else {
			// Otherwise, the previous minimum is still in this window. Compare
			// against the new value and update min if necessary.
			if h[r] < h[min] || (!robust && h[r] == h[min]) {
				min = r
				recorded = append(recorded, fingerprint)
			}
		}
	}

	return recorded
}

func LogNormalizeFile(docPath string) {

	bytes, err := ioutil.ReadFile(docPath)

	if err != nil {
		log.Fatal(err)
		return
	}

	norm := Normalize(bytes)

	fmt.Print(string(norm))
}

func LogFingerprintFile(docPath string) {
	bytes, err := ioutil.ReadFile(docPath)

	if err != nil {
		log.Fatal(err)
		return
	}

	norm := Normalize(bytes)
	prints := FingerprintDocument(norm, FINGERPRINT_SIZE)

	fmt.Println("location,fingerprint")
	for _, print := range prints {
		fmt.Printf("%d,%d\n", print.Location, print.Fingerprint)
	}
}

func LogWinnowFile(docPath string) {
	bytes, err := ioutil.ReadFile(docPath)

	if err != nil {
		log.Fatal(err)
		return
	}

	norm := Normalize(bytes)
	prints := FingerprintDocument(norm, FINGERPRINT_SIZE)
	saved := Winnow(prints, WINDOW_SIZE, true)

	for _, print := range saved {
		fmt.Printf("%d\n", print.Fingerprint)
	}
}

func WinnowFile(path string) map[uint64]bool {

	winnowed, err := winnowFile(path)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return winnowed
}

func winnowFile(docPath string) (map[uint64]bool, error) {
	winnowed := make(map[uint64]bool)

	bytes, err := ioutil.ReadFile(docPath)

	if err != nil {
		return winnowed, err
	}

	return WinnowData(bytes)
}

func WinnowData(bytes []byte) (map[uint64]bool, error) {
	winnowed := make(map[uint64]bool)

	norm := Normalize(bytes)
	prints := FingerprintDocument(norm, FINGERPRINT_SIZE)
	saved := Winnow(prints, WINDOW_SIZE, true)

	for _, print := range saved {
		winnowed[print.Fingerprint] = true
	}

	return winnowed, nil
}

func Similarity(doc1, doc2 string) {
	w1 := WinnowFile(doc1)
	w2 := WinnowFile(doc2)

	fmt.Printf("File 1: %s, hashes: %d\n", doc1, len(w1))
	fmt.Printf("File 2: %s, hashes: %d\n", doc2, len(w2))

	if len(w2) < len(w1) {
		tmp := w1
		w1 = w2
		w2 = tmp
	}

	matches := 0

	for key, _ := range w1 {
		if _, ok := w2[key]; ok {
			matches += 1
		}
	}

	fmt.Printf("Matches: %d (%d%%)\n", matches, (matches*100.0)/len(w1))

}

func CreateDocumentFromData(path string, data []byte) (*Document, error) {

	hashesMap, err := WinnowData(data)

	if err != nil {
		return nil, err
	}

	var doc Document

	doc.IndexDate = time.Now().Format(time.RFC3339)
	doc.Path = path
	_, doc.Name = filepath.Split(path)
	doc.Id = 0
	doc.Meta = make(map[string]string)

	doc.Hashes = make([]uint64, len(hashesMap))

	i := 0
	for k := range hashesMap {
		doc.Hashes[i] = k
		i++
	}

	return &doc, nil
}

func Report(documentId uint64, db *ParaphraseDb) error {

	doc, err := db.GetDoc(documentId)

	if err != nil {
		return err
	}

	matches := make(map[uint64]int)
	for _, hash := range doc.Hashes {
		docs, _ := db.GetDocsByHash(hash)

		for _, docId := range docs {
			ct, _ := matches[docId]
			matches[docId] = ct + 1
		}
	}

	// remove self matches
	delete(matches, documentId)

	hashCount := len(doc.Hashes)

	for k, v := range matches {
		fmt.Printf("%d: %d matches (%d%%)", k, v, (v*100.0)/hashCount)
		fmt.Println()
	}

	return nil
}
