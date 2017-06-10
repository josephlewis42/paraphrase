package paraphrase

import (
	"errors"

	"github.com/asdine/storm"
)

type IndexEntry struct {
	Hash uint64 `storm:"id,unique"`
	Docs map[string]int16
}

func NewIndexEntry(hash uint64) IndexEntry {
	var ie IndexEntry

	ie.Hash = hash
	ie.Docs = make(map[string]int16)

	return ie
}

func (ie *IndexEntry) AddDocument(docId string, frequency int16) {
	if ie.Docs == nil {
		ie.Docs = make(map[string]int16)
	}

	ie.Docs[docId] = frequency
}

func (p *ParaphraseDb) storeHash(tx storm.Node, hash uint64, docId string, count int16) error {
	index, err := p.getIndex(hash)
	if err != nil {
		return err
	}

	index.AddDocument(docId, count)
	return p.db.Save(index)
}

func (p *ParaphraseDb) getIndexOrBlank(hash uint64) (IndexEntry, error) {

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

func (p *ParaphraseDb) getIndex(hash uint64) (IndexEntry, error) {
	var index IndexEntry
	err := p.db.One("Hash", hash, &index)
	return index, err
}

type SearchResult struct {
	Similarity    float64
	simCalculated bool
	Query         *TermCountVector
	Doc           *Document
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

	return p.QueryByVector(vec)
}

func (p *ParaphraseDb) QueryByVector(query TermCountVector) (results []SearchResult, err error) {
	// TODO implement me

	return results, errors.New("Not yet implemented")
}
