package paraphrase

import (
	"errors"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
	"time"

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

var (
	SettingsNotDefinedErr = errors.New("No settings found. If you meant to create a database run 'paraphrase init'")
	AlreadyInitializedErr = errors.New("It looks like paraphrase has already been initialized.")
)

type Settings struct {
	Version         int `storm:"id,unique"`
	WindowSize      int
	FingerprintSize int
	RobustHash      bool
	CreatedAt       time.Time
}

func NewDefaultSettings() Settings {
	var settings Settings

	settings.Version = CurrentSettingsVersion
	settings.WindowSize = 10
	settings.FingerprintSize = 10
	settings.RobustHash = true
	settings.CreatedAt = time.Now()

	return settings
}

type ParaphraseDb struct {
	directory string
	settings  Settings
	db        *storm.DB
}

// Creates a new database in the given directory with the given settings
func Create(directory string, settings Settings) (*ParaphraseDb, error) {
	db, err := Open(directory)

	switch err {
	case nil:
		return nil, AlreadyInitializedErr
	case SettingsNotDefinedErr:
		db.settings = settings
		return db, db.saveSettings()
	default:
		return nil, err
	}
}

// Open or create a new paraphrase database in the given directory
func Open(directory string) (*ParaphraseDb, error) {

	var paraphrase ParaphraseDb
	var err error

	paraphrase.directory = directory

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
		return nil, fmt.Errorf("Could not open the database: %s", err)
	}

	err = paraphrase.loadSettings()
	if err != nil {
		return &paraphrase, err
	}

	return &paraphrase, nil
}

func (p *ParaphraseDb) init() error {

	err := p.db.Init(&Document{})
	if err != nil {
		return err
	}

	err = p.db.Init(&DocumentData{})

	if err != nil {
		return err
	}

	err = p.db.Init(&IndexEntry{})
	if err != nil {
		return err
	}

	err = p.db.Init(&Settings{})
	if err != nil {
		return err
	}

	return err
}

func (p *ParaphraseDb) Close() error {
	return p.db.Close()
}

func (p *ParaphraseDb) loadSettings() error {
	// load settings
	var settings Settings
	err := p.db.One("Version", CurrentSettingsVersion, &settings)

	switch err {
	case storm.ErrNotFound:
		return SettingsNotDefinedErr
	case nil:
		p.settings = settings

	default:
		return err
	}

	return nil
}

func (p *ParaphraseDb) saveSettings() error {
	return p.db.Save(&p.settings)
}

func (p *ParaphraseDb) GetSettings() Settings {
	return p.settings
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

	for hash, count := range doc.Hashes {
		err := p.storeHash(tx, hash, doc.Id, count)
		if err != nil {
			return nil, err
		}
	}

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

	return results, maskErrNotFound(err)
}

func (p *ParaphraseDb) FindDocumentById(id string) (*Document, error) {
	var doc Document
	err := p.db.One("Id", id, &doc)
	return &doc, err
}

func (p *ParaphraseDb) FindDocumentsBySha1(sha1 string) (results []Document, err error) {

	if len(sha1) == sha1HexLength {
		err = p.db.Find("Sha1", sha1, &results)

	} else {
		err = p.db.Select(q.Re("Sha1", "^"+sha1+".*")).Find(&results)
	}

	return results, maskErrNotFound(err)
}

func (p *ParaphraseDb) FindDocumentDataById(id string) (*DocumentData, error) {
	var doc DocumentData
	err := p.db.One("Id", id, &doc)
	return &doc, err
}

func maskErrNotFound(err error) error {
	if err == storm.ErrNotFound {
		return nil
	}

	return err
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
