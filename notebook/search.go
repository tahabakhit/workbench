package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"strings"
	"unicode"
)

// The title and body weights and the 20-token snippet start from zk's values.
const searchSQL = `
SELECT n.path, n.id, n.kind, n.date, n.repo, n.title,
       -bm25(notes_fts, 500.0, 1.0) AS score,
       snippet(notes_fts, -1, '**', '**', '...', 20)
FROM notes_fts CROSS JOIN notes n ON n.rowid = notes_fts.rowid
WHERE notes_fts MATCH ?1
  AND (?2 = '' OR n.repo = ?2)
  AND (?3 = '' OR n.kind = ?3)
  AND (?4 OR n.path NOT IN (SELECT target FROM supersedes WHERE path <> n.path))
ORDER BY score DESC, n.path
LIMIT ?5`

type result struct {
	Path    string  `json:"path"`
	ID      string  `json:"id"`
	Score   float64 `json:"score"`
	Snippet string  `json:"snippet"`
	Title   string  `json:"title"`
	Kind    string  `json:"kind"`
	Date    string  `json:"date,omitempty"`
	Repo    string  `json:"repo,omitempty"`
}

func runSearch(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.SetOutput(stderr)
	loc, err := addPathFlags(fs)
	if err != nil {
		return err
	}
	repo := fs.String("repo", "", "only notes from this `repo`")
	kind := fs.String("kind", "", "only notes of this `kind`, such as log or lessons")
	all := fs.Bool("all", false, "include superseded notes")
	limit := fs.Int("limit", 20, "show at most `n` results")
	asJSON := fs.Bool("json", false, "print JSON")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: notebook search <query> [flags]")
		fs.PrintDefaults()
	}
	words, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(words) == 0 {
		return usageError("search needs a query")
	}
	if *limit < 1 {
		return usageError("--limit must be at least 1")
	}
	match := matchExpression(strings.Join(words, " "))
	if match == "" {
		return usageError("the query has no words to search for")
	}

	// A missing vault must fail before an empty index is created for it.
	if err := checkVault(loc.vault); err != nil {
		return err
	}
	db, err := openIndex(loc.index)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx := context.Background()
	if err := refresh(ctx, db, loc.vault); err != nil {
		return err
	}
	results, err := search(ctx, db, match, *repo, *kind, *all, *limit)
	if err != nil {
		return err
	}

	if *asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}
	if len(results) == 0 {
		fmt.Fprintln(stderr, "no matches")
		return nil
	}
	for _, r := range results {
		fmt.Fprintf(stdout, "%s  %s  %.3f\n  %s\n", r.Path, r.ID, r.Score, r.Snippet)
	}
	return nil
}

// matchExpression quotes every word so that punctuation in a query is never
// read as FTS5 syntax. Quoted words are ANDed, and a word with no letter or
// digit is dropped because it would be an empty phrase.
func matchExpression(query string) string {
	var terms []string
	for _, word := range strings.Fields(query) {
		if strings.IndexFunc(word, func(r rune) bool {
			return unicode.IsLetter(r) || unicode.IsNumber(r)
		}) < 0 {
			continue
		}
		terms = append(terms, `"`+strings.ReplaceAll(word, `"`, `""`)+`"`)
	}
	return strings.Join(terms, " ")
}

func search(ctx context.Context, db *sql.DB, match, repo, kind string, all bool, limit int) ([]result, error) {
	rows, err := db.QueryContext(ctx, searchSQL, match, repo, kind, all, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := []result{}
	for rows.Next() {
		var r result
		if err := rows.Scan(&r.Path, &r.ID, &r.Kind, &r.Date, &r.Repo, &r.Title, &r.Score, &r.Snippet); err != nil {
			return nil, err
		}
		r.Score = math.Round(r.Score*1000) / 1000
		r.Snippet = strings.Join(strings.Fields(r.Snippet), " ")
		results = append(results, r)
	}
	return results, rows.Err()
}
