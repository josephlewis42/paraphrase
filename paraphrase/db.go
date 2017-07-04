package paraphrase

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/bradhe/stopwatch"
	"github.com/josephlewis42/paraphrase/paraphrase/provider"
	"github.com/josephlewis42/paraphrase/paraphrase/snappyjson"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	DbExt                  = ".ppdb"
	DbName                 = "paraphrasedb.ppdb"
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
	ImportErr             = errors.New("Errors encountered while importing documents")
	DatabaseDNEErr        = errors.New("It looks like the database does not exist, try running paraphrase init to create it")
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
		db.logChange("Created Database")
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
	dbPath := findDbPath(directory)

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

func findDbPath(directory string) string {

	if path.Ext(directory) == DbExt {
		return directory
	}

	return path.Join(directory, DbName)
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

	err = p.db.Init(&ChangeLogEntry{})
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
	p.logChange("Saved Settings")
	return p.db.Save(&p.settings)
}

func (p *ParaphraseDb) GetSettings() Settings {
	return p.settings
}

// Write information about Paraphrase and the database to an output.
// Output format _may change without warning_.
func (p *ParaphraseDb) WriteStats(writer io.Writer) {

	boltInfo := p.db.Bolt.Info()
	docCount, _ := p.CountDocuments()
	hashCount, _ := p.db.Count(&IndexEntry{})

	settings := []struct {
		Group   string
		Title   string
		Setting interface{}
	}{
		{"Paraphrase Settings", "", ""},
		{"", "Version", p.settings.Version},
		{"", "Window Size", p.settings.WindowSize},
		{"", "Fingerprint Length", p.settings.FingerprintSize},
		{"", "Robust Winnow?", p.settings.RobustHash},
		{"", "Creation Date", p.settings.CreatedAt},
		{"Database Information", "", ""},
		{"", "Page Size", boltInfo.PageSize},
		{"", "Number of Documents", docCount},
		{"", "Number of Distinct Hashes", hashCount},
	}

	w := new(tabwriter.Writer)
	w.Init(writer, 0, 8, 0, '\t', 0)

	for _, setting := range settings {
		fmt.Fprintf(w, "%v\t%v\t%v\n", setting.Group, setting.Title, setting.Setting)
	}
	w.Flush()

}

func (p *ParaphraseDb) AddDocuments(producer provider.DocumentProducer) (added []Document, ok bool) {
	start := stopwatch.Start()
	added = make([]Document, 0)
	ok = true

	for {
		select {
		case key, ok := <-producer:
			if !ok {
				producer = nil
				break
			}

			log.Printf("Adding %s %s\n", key.Namespace(), key.Path())
			log.Println("\tGetting Body")

			body, err := key.Body()

			if err != nil {
				log.Printf("Error getting body of %s: %s", key.Path(), err)
				ok = false
				continue
			}

			log.Println("\tSaving Document")
			doc, err := p.CreateDocument(key.Path(), key.Namespace(), body)
			if err != nil {
				log.Printf("Error saving document %s: %s", key.Path(), err)
				ok = false
				continue
			}

			added = append(added, *doc)
			log.Println("Moving on to next document")
		case <-time.After(time.Second * 1):
			log.Println("Waiting for more documents...")
		}

		if producer == nil {
			break
		}
	}

	watch := stopwatch.Stop(start)
	p.logChange("Added %v documents in %v ms, Had failures? %v", len(added), watch.Milliseconds(), !ok)

	return added, ok
}

func (p *ParaphraseDb) ImportDocumentsMatching(from *ParaphraseDb, query Document) error {
	start := stopwatch.Start()

	docs, err := from.FindDocumentsLike(query)

	if err != nil {
		return err
	}

	fmt.Printf("Starting import of %d documents\n", len(docs))

	bar := pb.StartNew(len(docs))
	var result error

	for _, doc := range docs {
		bar.Increment()

		data, err := from.FindDocumentDataById(doc.Id)
		if err != nil {
			log.Printf("Error getting the data for document: %s %s: %s", doc.Id, doc.Path, err)
			result = ImportErr
			continue
		}

		_, locerr := p.CreateDocument(doc.Path, doc.Namespace, data.Body)

		if locerr != nil {
			log.Printf("Error saving document %s: %s", doc.Path, err)
			result = ImportErr
		}
	}

	bar.FinishPrint("Finished importing")

	watch := stopwatch.Stop(start)
	p.logChange("Imported %v documents matching %v ion %v ms, err %v", len(docs), query, watch.Milliseconds(), result)

	return result
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

	err = tx.Commit()

	p.logChange("Created document %v", doc.Id)

	return doc, err
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
