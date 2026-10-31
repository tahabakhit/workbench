package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var (
	wikilink  = regexp.MustCompile(`\[\[([^\[\]\n]+)\]\]`)
	codeFence = regexp.MustCompile("(?ms)^(```|~~~).*?(^(```|~~~)[^\n]*$|\\z)")
	codeSpan  = regexp.MustCompile("`[^`\n]*`")
)

// parseLinks returns the raw targets of a note's wikilinks, without alias or
// heading. Text in code is not a link.
func parseLinks(body string) []string {
	body = codeFence.ReplaceAllString(body, "")
	body = codeSpan.ReplaceAllString(body, "")
	var targets []string
	for _, m := range wikilink.FindAllStringSubmatch(body, -1) {
		target, _, _ := strings.Cut(m[1], "|")
		target, _, _ = strings.Cut(target, "#")
		target = strings.TrimSpace(strings.TrimSuffix(target, `\`))
		if target != "" && !slices.Contains(targets, target) {
			targets = append(targets, target)
		}
	}
	return targets
}

// directory answers which notes a name refers to.
type directory struct {
	ids   map[string]string
	paths map[string]bool
	names map[string][]string
}

func newDirectory() *directory {
	return &directory{ids: map[string]string{}, paths: map[string]bool{}, names: map[string][]string{}}
}

func (d *directory) add(rel, id, title, stem string) {
	d.ids[id] = rel
	d.paths[rel] = true
	d.names[title] = append(d.names[title], rel)
	if stem != title {
		d.names[stem] = append(d.names[stem], rel)
	}
}

func loadDirectory(ctx context.Context, db *sql.DB) (*directory, error) {
	rows, err := db.QueryContext(ctx, "SELECT path, id, title, stem FROM notes ORDER BY path")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	d := newDirectory()
	for rows.Next() {
		var rel, id, title, stem string
		if err := rows.Scan(&rel, &id, &title, &stem); err != nil {
			return nil, err
		}
		d.add(rel, id, title, stem)
	}
	return d, rows.Err()
}

// resolve returns every note a name could mean: a path match alone if there is
// one, otherwise all notes with that exact title or file name. Only a result of
// length one is a resolved link; nothing here ever picks the first.
func (d *directory) resolve(name string) []string {
	withExt := name
	if !strings.HasSuffix(withExt, ".md") {
		withExt += ".md"
	}
	if d.paths[withExt] {
		return []string{withExt}
	}
	found := append([]string{}, d.names[name]...)
	if stem := strings.TrimSuffix(name, ".md"); stem != name {
		found = append(found, d.names[stem]...)
	}
	slices.Sort(found)
	return slices.Compact(found)
}

// findNote turns a command argument (path, title, file name or ID) into the
// one note it means, and fails naming the candidates otherwise.
func findNote(ctx context.Context, db *sql.DB, name string) (string, error) {
	d, err := loadDirectory(ctx, db)
	if err != nil {
		return "", err
	}
	if rel, ok := d.ids[name]; ok {
		return rel, nil
	}
	found := d.resolve(name)
	switch len(found) {
	case 1:
		return found[0], nil
	case 0:
		var similar []result
		if match := matchExpression(name); match != "" {
			if similar, err = search(ctx, db, match, "", "", true, 5); err != nil {
				return "", err
			}
		}
		msg := fmt.Sprintf("no note matches %q", name)
		if len(similar) > 0 {
			msg += "; similar:"
		}
		for _, r := range similar {
			msg += "\n  " + r.Path
		}
		return "", fmt.Errorf("%s", msg)
	default:
		return "", fmt.Errorf("more than one note matches %q:\n  %s", name, strings.Join(found, "\n  "))
	}
}

// link is one outgoing wikilink. Status is ok, broken or ambiguous; Paths holds
// the note it resolves to, or the candidates when ambiguous.
type link struct {
	Target string   `json:"target"`
	Status string   `json:"status"`
	Paths  []string `json:"paths"`
}

func outgoingLinks(ctx context.Context, db *sql.DB, rel string) ([]link, error) {
	d, err := loadDirectory(ctx, db)
	if err != nil {
		return nil, err
	}
	targets, err := queryStrings(ctx, db, "SELECT target FROM links WHERE path = ? ORDER BY target", rel)
	if err != nil {
		return nil, err
	}
	found := []link{}
	for _, target := range targets {
		paths := d.resolve(target)
		status := "ambiguous"
		switch len(paths) {
		case 0:
			status = "broken"
		case 1:
			status = "ok"
		}
		found = append(found, link{Target: target, Status: status, Paths: paths})
	}
	return found, nil
}

// backlinksTo lists the notes with a link that resolves to rel alone.
func backlinksTo(ctx context.Context, db *sql.DB, rel string) ([]string, error) {
	d, err := loadDirectory(ctx, db)
	if err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(ctx, "SELECT path, target FROM links ORDER BY path")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sources := []string{}
	for rows.Next() {
		var source, target string
		if err := rows.Scan(&source, &target); err != nil {
			return nil, err
		}
		if paths := d.resolve(target); len(paths) == 1 && paths[0] == rel {
			sources = append(sources, source)
		}
	}
	return slices.Compact(sources), rows.Err()
}

func queryStrings(ctx context.Context, db *sql.DB, query string, args ...any) ([]string, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := []string{}
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, rows.Err()
}

// withNote opens a fresh index and finds the one note a name means.
func withNote(ctx context.Context, loc *locations, name string, use func(db *sql.DB, rel string) error) error {
	if strings.TrimSpace(name) == "" {
		return usageError("a note is needed")
	}
	db, err := openFresh(ctx, loc)
	if err != nil {
		return err
	}
	defer db.Close()
	rel, err := findNote(ctx, db, name)
	if err != nil {
		return err
	}
	return use(db, rel)
}

func recallNote(ctx context.Context, loc *locations, name string) (string, error) {
	var text []byte
	err := withNote(ctx, loc, name, func(_ *sql.DB, rel string) (err error) {
		text, err = os.ReadFile(filepath.Join(loc.vault, filepath.FromSlash(rel)))
		return err
	})
	return string(text), err
}

func linksOf(ctx context.Context, loc *locations, name string) ([]link, error) {
	var found []link
	err := withNote(ctx, loc, name, func(db *sql.DB, rel string) (err error) {
		found, err = outgoingLinks(ctx, db, rel)
		return err
	})
	return found, err
}

func backlinksOf(ctx context.Context, loc *locations, name string) ([]string, error) {
	var sources []string
	err := withNote(ctx, loc, name, func(db *sql.DB, rel string) (err error) {
		sources, err = backlinksTo(ctx, db, rel)
		return err
	})
	return sources, err
}

// noteArgs parses the arguments of a command that takes one note.
func noteArgs(name string, args []string, stderr io.Writer, asJSON *bool) (*locations, string, error) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	loc, err := addPathFlags(fs)
	if err != nil {
		return nil, "", err
	}
	if asJSON != nil {
		fs.BoolVar(asJSON, "json", false, "print JSON")
	}
	fs.Usage = func() {
		fmt.Fprintf(stderr, "usage: notebook %s <path, title, file name or ID> [flags]\n", name)
		fs.PrintDefaults()
	}
	words, err := parseFlags(fs, args)
	if err != nil {
		return nil, "", err
	}
	return loc, strings.Join(words, " "), nil
}

func runRecall(args []string, stdout, stderr io.Writer) error {
	loc, name, err := noteArgs("recall", args, stderr, nil)
	if err != nil {
		return err
	}
	text, err := recallNote(context.Background(), loc, name)
	if err != nil {
		return err
	}
	_, err = io.WriteString(stdout, text)
	return err
}

func runLinks(args []string, stdout, stderr io.Writer) error {
	var asJSON bool
	loc, name, err := noteArgs("links", args, stderr, &asJSON)
	if err != nil {
		return err
	}
	found, err := linksOf(context.Background(), loc, name)
	if err != nil {
		return err
	}
	if asJSON {
		return printJSON(stdout, found)
	}
	if len(found) == 0 {
		fmt.Fprintln(stderr, "no links")
	}
	for _, l := range found {
		line := l.Status + "\t" + l.Target
		if len(l.Paths) > 0 {
			line += "\t" + strings.Join(l.Paths, ", ")
		}
		fmt.Fprintln(stdout, line)
	}
	return nil
}

func runBacklinks(args []string, stdout, stderr io.Writer) error {
	var asJSON bool
	loc, name, err := noteArgs("backlinks", args, stderr, &asJSON)
	if err != nil {
		return err
	}
	sources, err := backlinksOf(context.Background(), loc, name)
	if err != nil {
		return err
	}
	if asJSON {
		return printJSON(stdout, sources)
	}
	if len(sources) == 0 {
		fmt.Fprintln(stderr, "no backlinks")
	}
	for _, source := range sources {
		fmt.Fprintln(stdout, source)
	}
	return nil
}

func printJSON(w io.Writer, value any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}
