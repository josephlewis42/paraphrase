package paraphrase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

const (
	DbName         = "paraphrasedb.bolt"
	DocumentBucket = "documents"
	IndexBucket    = "index"
	MinIndex       = "00000000000000000000"
	MaxIndex       = "99999999999999999999"
)

type Document struct {
	Id        uint64 // the ID assigned by the bolt db
	IndexDate string
	Path      string
	Name      string
	Hashes    []uint64
	Meta      map[string]string
}

type ParaphraseDb struct {
	directory string
	db        *bolt.DB
}

// Open or create a new paraphrase database in the given directory
func Open(directory string) (*ParaphraseDb, error) {

	var paraphrase ParaphraseDb
	var err error

	paraphrase.directory = directory

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	dbPath := path.Join(directory, DbName)

	paraphrase.db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = paraphrase.init()
	if err != nil {
		return nil, fmt.Errorf("Could not init db: %s", err)
	}

	return &paraphrase, nil
}

func (db *ParaphraseDb) init() error {
	return db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DocumentBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(IndexBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})
}

func (db *ParaphraseDb) Close() error {
	return db.db.Close()
}

func (db *ParaphraseDb) DocList() ([]string, error) {

	docs := make([]string, 0)

	db.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(DocumentBucket))

		b.ForEach(func(k, v []byte) error {
			docs = append(docs, string(k))
			return nil
		})
		return nil
	})

	return docs, nil
}

func (db *ParaphraseDb) GetDoc(id uint64) (*Document, error) {

	docs, err := db.scanDocs(idToKey(id, ""), 1)

	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, errors.New("No document found with the given id")
	}

	return &docs[0], nil
}

func (db *ParaphraseDb) GetDocsByPath(pathPrefix string) ([]Document, error) {
	docs := make([]Document, 0)

	err := db.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(DocumentBucket))

		return b.ForEach(func(k, v []byte) error {
			key := string(k)
			docPath := key[len(MinIndex)+1:]

			if strings.HasPrefix(docPath, pathPrefix) {
				var doc Document
				err := json.Unmarshal(v, &doc)
				if err != nil {
					return err
				}
				docs = append(docs, doc)
			}
			return nil
		})
	})

	return docs, err
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func (db *ParaphraseDb) GetDocsByHash(hash uint64) ([]uint64, error) {
	keys := make([]uint64, 0)

	err := db.db.View(func(tx *bolt.Tx) error {
		intBucket := tx.Bucket([]byte(IndexBucket))

		hashb := []byte(strconv.Itoa(int(hash)))
		val := intBucket.Get(hashb)

		if val != nil {
			json.Unmarshal(val, &keys)
		}

		return nil
	})

	return keys, err
}

func (db *ParaphraseDb) scanDocs(prefix []byte, limit int) ([]Document, error) {
	docs := make([]Document, 0)

	if limit <= 0 {
		return docs, nil
	}

	err := db.db.View(func(tx *bolt.Tx) error {

		// Assume bucket exists and has keys
		c := tx.Bucket([]byte(DocumentBucket)).Cursor()

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var doc Document
			err := json.Unmarshal(v, &doc)

			if err != nil {
				return err
			}

			docs = append(docs, doc)
			return nil
		}

		return nil
	})

	return docs, err
}

func (db *ParaphraseDb) Insert(doc *Document) (uint64, error) {
	var id uint64
	var err error

	err = db.db.Update(func(tx *bolt.Tx) error {

		docBucket := tx.Bucket([]byte(DocumentBucket))
		intBucket := tx.Bucket([]byte(IndexBucket))

		id, err = docBucket.NextSequence()

		if err != nil {
			return fmt.Errorf("error getting next sequence: %s", err)
		}

		doc.Id = id

		docHash := idToKey(doc.Id, doc.Path)

		buf, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		docBucket.Put(docHash, buf)

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for _, hash := range doc.Hashes {

			hashb := []byte(strconv.Itoa(int(hash)))
			val := intBucket.Get(hashb)

			keys := make([]uint64, 1)
			if val != nil {
				json.Unmarshal(val, &keys)
			}

			keys = append(keys, id)
			val, err = json.Marshal(keys)

			if err != nil {
				return fmt.Errorf("error serializing: %s", err)
			}

			err := intBucket.Put(hashb, val)

			if err != nil {
				return fmt.Errorf("Error storing data: %s", err)
			}
		}

		return nil
	})

	return id, err
}

func idToKey(id uint64, path string) []byte {
	return []byte(strconv.Itoa(int(id)))
}
