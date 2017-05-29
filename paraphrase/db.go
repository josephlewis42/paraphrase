package paraphrase

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/josephlewis42/paraphrase/paraphrase/provider"
	"github.com/josephlewis42/paraphrase/paraphrase/snappyjson"
)

const (
	DbName                 = "paraphrasedb.bolt"
	DocumentBucket         = "documents"
	IndexBucket            = "index"
	SettingsBucket         = "settings"
	FileBucket             = "files"
	MinIndex               = "00000000000000000000"
	MaxIndex               = "99999999999999999999"
	CurrentSettingsVersion = 1 // the version of the settings file, won't match the version of paraphrase
	sha1HexLength          = len("da39a3ee5e6b4b0d3255bfef95601890afd80709")
)

type Settings struct {
	Version         int
	WindowSize      int
	FingerprintSize int
	RobustHash      bool
}

func NewDefaultSettings() Settings {
	var settings Settings

	settings.Version = CurrentSettingsVersion
	settings.WindowSize = 10
	settings.FingerprintSize = 10
	settings.RobustHash = true

	return settings
}

type ParaphraseDb struct {
	directory string
	settings  Settings
	db        *storm.DB
}

// Open or create a new paraphrase database in the given directory
func Open(directory string) (*ParaphraseDb, error) {

	var paraphrase ParaphraseDb
	var err error

	paraphrase.directory = directory
	paraphrase.settings = NewDefaultSettings()

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	dbPath := path.Join(directory, DbName)

	// paraphrase.db, err = storm.Open(dbPath, storm.Codec(snappyjson.Codec))
	paraphrase.db, err = storm.Open(dbPath, storm.Codec(snappyjson.MsgpackCodec))
	// paraphrase.db, err = storm.Open(dbPath)
	if err != nil {
		return nil, err
	}

	err = paraphrase.init()
	if err != nil {
		return nil, fmt.Errorf("Could not init db: %s", err)
	}

	return &paraphrase, nil
}

func (p *ParaphraseDb) init() error {

	err := p.db.Init(&Document{})
	if err != nil {
		return err
	}

	err = p.db.Init(&DocumentData{})

	return err
}

func (p *ParaphraseDb) Close() error {
	return p.db.Close()
}

func (p *ParaphraseDb) AddDocuments(producer provider.DocumentProducer) (added []Document, ok bool) {
	added = make([]Document, 0)
	ok = true

	for key := range producer {
		log.Printf("Adding %s %s\n", key.Namespace(), key.Path())
		body, err := key.Body()

		if err != nil {
			log.Printf("Error getting body of %s: %s", key.Path(), err)
			ok = false
			continue
		}

		doc, err := p.CreateDocument(key.Path(), key.Namespace(), body)
		if err != nil {
			log.Printf("Error saving document %s: %s", key.Path(), err)
			ok = false
			continue
		}

		added = append(added, *doc)
	}

	return added, ok
}

func (p *ParaphraseDb) CreateDocument(path, namespace string, body []byte) (*Document, error) {
	var err error

	doc, docData := NewDocument(path, namespace, body)

	// generate hashes
	doc.Hashes, err = p.WinnowData(body)

	if err != nil {
		return nil, err
	}

	tx, err := p.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// TODO save hashes to db

	err = tx.Save(doc)
	if err != nil {
		return nil, err
	}

	err = tx.Save(docData)
	if err != nil {
		return nil, err
	}

	return doc, tx.Commit()
}

func (p *ParaphraseDb) CountDocuments() (int, error) {
	return p.db.Count(&Document{})
}

// FindDocumentsLike finds documents like the one given.
// * Ids are matched exactly,
// * SHA1s are matched as a prefix (you can give the n characters only)
// * Namespaces are searched like globs
// * Paths are searched like globs
func (p *ParaphraseDb) FindDocumentsLike(query Document) (results []Document, err error) {
	var matchers []q.Matcher

	if query.Id != "" {
		matchers = append(matchers, q.Eq("Id", query.Id))
	}

	if query.Sha1 != "" {
		matchers = append(matchers, q.Re("Sha1", "^"+query.Sha1+".*"))
	}

	if query.Namespace != "" {
		matchers = append(matchers, q.Re("Namespace", GlobToRegexStr(query.Namespace)))
	}

	if query.Path != "" {
		matchers = append(matchers, q.Re("Path", GlobToRegexStr(query.Path)))
	}

	err = p.db.Select(matchers...).Find(&results)

	return results, err
}

func (p *ParaphraseDb) FindDocumentById(id string) (*Document, error) {
	var doc Document
	err := p.db.One("Id", id, &doc)
	return &doc, err
}

func (p *ParaphraseDb) FindDocumentsByIds(ids ...string) (results []Document, err error) {
	err = p.db.Select(q.In("Id", ids)).Find(&results)
	return results, err
}

func (p *ParaphraseDb) FindDocumentsBySha1(sha1 string) (results []Document, err error) {

	if len(sha1) == sha1HexLength {
		err = p.db.Find("Sha1", sha1, &results)

	} else {
		err = p.db.Select(q.Re("Sha1", "^"+sha1+".*")).Find(&results)
	}

	return results, err
}

func (p *ParaphraseDb) FindDocumentDataById(id string) (*DocumentData, error) {
	var doc DocumentData
	err := p.db.One("Id", id, &doc)
	return &doc, err
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
