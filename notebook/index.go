package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// schemaVersion is raised on any schema change; an index at another version is
// rebuilt, which is always safe because the files are the truth.
const schemaVersion = 2

const schema = `
CREATE TABLE notes (
	rowid INTEGER PRIMARY KEY,
	path TEXT NOT NULL UNIQUE,
	id TEXT NOT NULL,
	kind TEXT NOT NULL,
	date TEXT NOT NULL,
	repo TEXT NOT NULL,
	scope TEXT NOT NULL,
	title TEXT NOT NULL,
	stem TEXT NOT NULL,
	mtime INTEGER NOT NULL
);
CREATE TABLE supersedes (path TEXT NOT NULL, target TEXT NOT NULL);
CREATE INDEX supersedes_target ON supersedes(target);
CREATE TABLE links (path TEXT NOT NULL, target TEXT NOT NULL, PRIMARY KEY (path, target));
CREATE VIRTUAL TABLE notes_fts USING fts5(
	title, body,
	tokenize = 'porter unicode61 remove_diacritics 1'
);
`

// openFresh opens the index and brings it up to date, as every command does
// first. A missing vault fails before an empty index is created for it.
func openFresh(ctx context.Context, loc *locations) (*sql.DB, error) {
	if err := checkVault(loc.vault); err != nil {
		return nil, err
	}
	db, err := openIndex(loc.index)
	if err != nil {
		return nil, err
	}
	if err := refresh(ctx, db, loc.vault); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// openIndex opens the index file, creating or rebuilding its schema as needed.
// Every transaction takes the write lock up front, so concurrent sessions wait
// for each other instead of failing.
func openIndex(file string) (*sql.DB, error) {
	abs, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return nil, err
	}
	dsn := url.URL{
		Scheme:   "file",
		Path:     abs,
		RawQuery: "_pragma=busy_timeout(30000)&_txlock=immediate",
	}
	db, err := sql.Open("sqlite", dsn.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("index %s: %w", file, err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	version, err := userVersion(db)
	if err != nil || version == schemaVersion {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Another process may have finished while this one waited for the lock.
	if version, err = userVersion(tx); err != nil || version == schemaVersion {
		return err
	}
	for _, table := range []string{"notes", "supersedes", "links", "notes_fts"} {
		if _, err := tx.Exec("DROP TABLE IF EXISTS " + table); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(schema); err != nil {
		return err
	}
	if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", schemaVersion)); err != nil {
		return err
	}
	return tx.Commit()
}

type queryRower interface {
	QueryRow(query string, args ...any) *sql.Row
}

func userVersion(q queryRower) (int, error) {
	var version int
	err := q.QueryRow("PRAGMA user_version").Scan(&version)
	return version, err
}

// refresh brings the index in line with the vault: it reindexes every note
// whose modification time differs from the stored one and drops deleted notes.
func refresh(ctx context.Context, db *sql.DB, vault string) error {
	current, err := walkVault(vault)
	if err != nil {
		return err
	}
	stored, err := storedTimes(ctx, db)
	if err != nil {
		return err
	}

	var removed, changed []string
	for rel := range stored {
		if _, ok := current[rel]; !ok {
			removed = append(removed, rel)
		}
	}
	for rel, mtime := range current {
		if prev, ok := stored[rel]; !ok || prev != mtime {
			changed = append(changed, rel)
		}
	}
	if len(removed) == 0 && len(changed) == 0 {
		return nil
	}

	notes := make(map[string]note, len(changed))
	for _, rel := range changed {
		data, err := os.ReadFile(filepath.Join(vault, filepath.FromSlash(rel)))
		if errors.Is(err, fs.ErrNotExist) {
			removed = append(removed, rel)
			continue
		}
		if err != nil {
			return err
		}
		notes[rel] = parseNote(rel, data)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, rel := range removed {
		if err := dropNote(tx, rel); err != nil {
			return err
		}
	}
	for rel, n := range notes {
		if err := dropNote(tx, rel); err != nil {
			return err
		}
		if err := addNote(tx, rel, current[rel], n); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func checkVault(vault string) error {
	if info, err := os.Stat(vault); err != nil || !info.IsDir() {
		return fmt.Errorf("no vault at %s", vault)
	}
	return nil
}

// walkVault lists the vault's markdown files with their modification times.
// Dot folders such as .git and .obsidian are not part of the notebook.
func walkVault(vault string) (map[string]int64, error) {
	files := map[string]int64{}
	err := filepath.WalkDir(vault, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p != vault && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		info, err := d.Info()
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(vault, p)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = info.ModTime().UnixNano()
		return nil
	})
	return files, err
}

func storedTimes(ctx context.Context, db *sql.DB) (map[string]int64, error) {
	rows, err := db.QueryContext(ctx, "SELECT path, mtime FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	times := map[string]int64{}
	for rows.Next() {
		var rel string
		var mtime int64
		if err := rows.Scan(&rel, &mtime); err != nil {
			return nil, err
		}
		times[rel] = mtime
	}
	return times, rows.Err()
}

func dropNote(tx *sql.Tx, rel string) error {
	if _, err := tx.Exec("DELETE FROM notes_fts WHERE rowid = (SELECT rowid FROM notes WHERE path = ?)", rel); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM notes WHERE path = ?", rel); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM links WHERE path = ?", rel); err != nil {
		return err
	}
	_, err := tx.Exec("DELETE FROM supersedes WHERE path = ?", rel)
	return err
}

func addNote(tx *sql.Tx, rel string, mtime int64, n note) error {
	res, err := tx.Exec(
		"INSERT INTO notes (path, id, kind, date, repo, scope, title, stem, mtime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		rel, n.id, n.kind, n.date, n.repo, n.scope, n.title, strings.TrimSuffix(path.Base(rel), ".md"), mtime)
	if err != nil {
		return err
	}
	rowid, err := res.LastInsertId()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO notes_fts (rowid, title, body) VALUES (?, ?, ?)", rowid, n.title, n.body); err != nil {
		return err
	}
	for _, target := range n.links {
		if _, err := tx.Exec("INSERT INTO links (path, target) VALUES (?, ?)", rel, target); err != nil {
			return err
		}
	}
	if n.supersedes != "" {
		_, err = tx.Exec("INSERT INTO supersedes (path, target) VALUES (?, ?)", rel, n.supersedes)
	}
	return err
}
