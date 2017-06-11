package paraphrase

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/asdine/storm"
	"github.com/josephlewis42/paraphrase/paraphrase/linalg"
)

type IndexEntry struct {
	Hash uint64 `storm:"id"`
	Docs map[string]int16
}

func NewIndexEntry(hash uint64) *IndexEntry {
	var ie IndexEntry

	ie.Hash = hash
	ie.Docs = make(map[string]int16)

	return &ie
}

func (ie *IndexEntry) AddDocument(docId string, frequency int16) {
	if ie.Docs == nil {
		ie.Docs = make(map[string]int16)
	}

	ie.Docs[docId] = frequency
}

func (p *ParaphraseDb) storeHash(tx storm.Node, hash uint64, docId string, count int16) error {
	index, err := p.getIndexOrBlank(hash)
	if err != nil {
		return err
	}

	index.AddDocument(docId, count)
	return tx.Save(index)
}

func (p *ParaphraseDb) getIndexOrBlank(hash uint64) (*IndexEntry, error) {

	index, err := p.getIndex(hash)
	if err == nil {
		return index, nil
	}

	if err == storm.ErrNotFound {
		ie := NewIndexEntry(hash)
		return ie, nil
	}

	return index, err
}

func (p *ParaphraseDb) getIndex(hash uint64) (*IndexEntry, error) {
	var index IndexEntry
	err := p.db.One("Hash", hash, &index)
	return &index, err
}

type SearchResult struct {
	Query *TermCountVector
	Doc   *Document
}

func (sr *SearchResult) Similarity() float64 {
	return 1.0
}

func (p *ParaphraseDb) QueryById(id string) (results []SearchResult, err error) {

	doc, err := p.FindDocumentById(id)
	if err != nil {
		return nil, err
	}

	return p.QueryByVector(doc.Hashes)
}

func (p *ParaphraseDb) QueryByString(query string) (results []SearchResult, err error) {
	vec, err := p.WinnowData([]byte(query))

	if err != nil {
		return results, err
	}

	if len(vec) == 0 {
		return results, errors.New("Query was not long enough to search.")
	}

	return p.QueryByVector(vec)
}

func (p *ParaphraseDb) QueryByVector(query TermCountVector) (results []SearchResult, err error) {

	countI, err := p.CountDocuments()

	if err != nil {
		return results, err
	}

	count := float64(countI)

	idfVector := make(linalg.IFVector)
	matchingDocIds := make(map[string]bool)

	for hash, _ := range query {
		idx, err := p.getIndex(hash)
		fmt.Println(hash)

		switch err {
		case nil:
			docFrequency := 1 + len(idx.Docs)
			idfVector[hash] = 1 + math.Log(count/float64(docFrequency))

			for id, _ := range idx.Docs {
				matchingDocIds[id] = true
			}

		case storm.ErrNotFound:
			fmt.Printf("Nothing found for %d\n", hash)

			continue // a query might not have any matching documents

		default:
			return results, err
		}
	}

	for id, _ := range matchingDocIds {
		doc, err := p.FindDocumentById(id)
		if err != nil {
			log.Printf("Could not fetch doc %s: %s", id, err)
			continue
		}

		results = append(results, SearchResult{&query, doc})
	}

	return results, nil
}
