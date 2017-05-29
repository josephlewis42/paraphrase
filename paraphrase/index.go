package paraphrase

import "github.com/boltdb/bolt"

type indexEntry struct {
	Key  string `storm:"id,unique"`
	Docs map[string]int64
}

func (ie *indexEntry) AddDocument(docId string, frequency int64) {
	if ie.Docs == nil {
		ie.Docs = make(map[string]int64)
	}

	ie.Docs[docId] = frequency
}

func (p *ParaphraseDb) StoreHash(tx bolt.Tx, docId string, count int16) error {
	return nil
}
