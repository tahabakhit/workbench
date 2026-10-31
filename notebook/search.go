package main

import (
	"context"
	"database/sql"
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

const defaultLimit = 20

type searchOptions struct {
	query, repo, kind string
	all               bool
	limit             int
}

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

// searchVault is the search command behind both the CLI and the MCP tool.
func searchVault(ctx context.Context, loc *locations, opts searchOptions) ([]result, error) {
	if strings.TrimSpace(opts.query) == "" {
		return nil, usageError("search needs a query")
	}
	if opts.limit < 1 {
		return nil, usageError("the limit must be at least 1")
	}
	match := matchExpression(opts.query)
	if match == "" {
		return nil, usageError("the query has no words to search for")
	}
	db, err := openFresh(ctx, loc)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return search(ctx, db, match, opts.repo, opts.kind, opts.all, opts.limit)
}

func runSearch(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.SetOutput(stderr)
	loc, err := addPathFlags(fs)
	if err != nil {
		return err
	}
	var opts searchOptions
	fs.StringVar(&opts.repo, "repo", "", "only notes from this `repo`")
	fs.StringVar(&opts.kind, "kind", "", "only notes of this `kind`, such as log or lessons")
	fs.BoolVar(&opts.all, "all", false, "include superseded notes")
	fs.IntVar(&opts.limit, "limit", defaultLimit, "show at most `n` results")
	asJSON := fs.Bool("json", false, "print JSON")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: notebook search <query> [flags]")
		fs.PrintDefaults()
	}
	words, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	opts.query = strings.Join(words, " ")
	results, err := searchVault(context.Background(), loc, opts)
	if err != nil {
		return err
	}

	if *asJSON {
		return printJSON(stdout, results)
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
