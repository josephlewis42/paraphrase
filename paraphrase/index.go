package paraphrase

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/asdine/storm"
	"github.com/bradfitz/slice"
	"github.com/josephlewis42/paraphrase/paraphrase/linalg"
)

type IndexEntry struct {
	Hash      uint64 `storm:"id,index"`
	Doc       int64
	Frequency int16
}

func (p *ParaphraseDb) storeHash(tx storm.Node, hash uint64, docId int64, count int16) error {
	return tx.Save(&IndexEntry{hash, docId, count})
}

func (p *ParaphraseDb) getIndex(hash uint64) ([]IndexEntry, error) {
	var index []IndexEntry

	//err := p.db.Select(q.Eq("Hash", hash)).Find(&index)
	err := p.db.Find("Hash", hash, &index)
	return index, err
}

type SearchResult struct {
	Query *TermCountVector
	Doc   *Document
}

func (sr *SearchResult) Similarity() float64 {
	match := 0.0
	mismatch := 0.0

	for k, _ := range sr.Doc.Hashes {
		if _, ok := (*sr.Query)[k]; ok {
			match++
		} else {
			mismatch++
		}
	}

	//sr.Doc.Hashes
	return match / (match + mismatch)
}

func (p *ParaphraseDb) QueryById(id int64) (results []SearchResult, err error) {

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
	matchingDocIds := make(map[int64]bool)

	for hash, _ := range query {
		idx, err := p.getIndex(hash)

		switch err {
		case nil:
			docFrequency := 1 + len(idx)
			idfVector[hash] = 1 + math.Log(count/float64(docFrequency))

			for _, doc := range idx {
				matchingDocIds[doc.Doc] = true
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
			log.Printf("Could not fetch doc %d: %s", id, err)
			continue
		}

		results = append(results, SearchResult{&query, doc})
	}

	slice.Sort(results[:], func(i, j int) bool {
		return results[i].Similarity() > results[j].Similarity()
	})

	return results, nil
}
