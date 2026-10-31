package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	pushLimit = 2000
	pushLogs  = 3
)

// pushNotes lists paths and titles in the order they are kept: dropping from
// the end removes global lessons, then repo lessons, then logs.
const pushNotes = `
SELECT path, title FROM (
	SELECT n.path, n.title, 0 AS tier, n.date
	FROM notes n
	WHERE n.kind = 'log' AND n.repo = ?1 AND n.path NOT IN (
		SELECT path FROM notes WHERE kind = 'log' AND repo = ?1
		ORDER BY date DESC, path DESC LIMIT -1 OFFSET ?2)
	UNION ALL
	SELECT n.path, n.title, CASE n.scope WHEN ?1 THEN 1 ELSE 2 END, n.date
	FROM notes n
	WHERE n.kind = 'lessons' AND n.scope IN (?1, 'global')
)
WHERE path NOT IN (SELECT target FROM supersedes WHERE supersedes.path <> target)
ORDER BY tier, date DESC, path DESC`

func runPush(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("push", flag.ContinueOnError)
	fs.SetOutput(stderr)
	loc, err := addPathFlags(fs)
	if err != nil {
		return err
	}
	repo := fs.String("repo", "", "the session's `repo`")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: notebook push --repo <name> [flags]")
		fs.PrintDefaults()
	}
	if extra, err := parseFlags(fs, args); err != nil {
		return err
	} else if len(extra) > 0 || *repo == "" {
		return usageError("push takes --repo <name> and nothing else")
	}
	ctx := context.Background()
	db, err := openFresh(ctx, loc)
	if err != nil {
		return err
	}
	defer db.Close()
	lines, err := pushLines(ctx, db, *repo)
	if err != nil {
		return err
	}
	_, err = io.WriteString(stdout, pushText(*repo, lines))
	return err
}

func pushLines(ctx context.Context, db *sql.DB, repo string) ([]string, error) {
	rows, err := db.QueryContext(ctx, pushNotes, repo, pushLogs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lines []string
	for rows.Next() {
		var rel, title string
		if err := rows.Scan(&rel, &title); err != nil {
			return nil, err
		}
		lines = append(lines, "- "+rel+": "+title+"\n")
	}
	return lines, rows.Err()
}

// pushText drops lines from the end until the whole text fits the limit.
func pushText(repo string, lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	head := fmt.Sprintf("## Notebook\n\nNotes for %s. Read one with \"notebook recall <path>\".\n\n", repo)
	for kept := len(lines); ; kept-- {
		text := head
		for _, line := range lines[:kept] {
			text += line
		}
		if dropped := len(lines) - kept; dropped > 0 {
			text += fmt.Sprintf("(%d more dropped)\n", dropped)
		}
		if kept == 0 || utf8.RuneCountInString(text) <= pushLimit {
			return text
		}
	}
}
