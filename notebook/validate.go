package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

const (
	nameRE = `[A-Za-z0-9][A-Za-z0-9._-]*`
	dateRE = `[0-9]{4}-[0-9]{2}-[0-9]{2}`
)

// layout is the Notebook table in contract.md.
var layout = []*regexp.Regexp{
	regexp.MustCompile(`^log/` + dateRE + `-` + nameRE + `-` + nameRE + `\.md$`),
	regexp.MustCompile(`^decisions/` + dateRE + `-` + nameRE + `\.md$`),
	regexp.MustCompile(`^lessons/` + nameRE + `\.md$`),
	regexp.MustCompile(`^corrections\.md$`),
	regexp.MustCompile(`^goals\.md$`),
	regexp.MustCompile(`^plans/` + dateRE + `\.md$`),
	regexp.MustCompile(`^learning/` + nameRE + `/(mission|sources|quiz)\.md$`),
	regexp.MustCompile(`^learning/` + nameRE + `/lessons/` + nameRE + `\.md$`),
}

// updatable lists the files skills keep adding to.
var updatable = regexp.MustCompile(`^(corrections\.md|goals\.md|plans/.*|learning/[^/]+/(mission|sources|quiz)\.md)$`)

var leadingDate = regexp.MustCompile(`^` + dateRE)

func inLayout(rel string) bool {
	if !slices.ContainsFunc(layout, func(re *regexp.Regexp) bool { return re.MatchString(rel) }) {
		return false
	}
	date := leadingDate.FindString(path.Base(rel))
	if date == "" {
		return true
	}
	_, err := time.Parse(time.DateOnly, date)
	return err == nil
}

// problem is one finding. An error refuses a write; a warning does not.
type problem struct{ level, path, message string }

// vaultState is the vault read straight from disk, so validation needs no index.
type vaultState struct {
	texts map[string]string
	notes map[string]note
	dir   *directory
	stems map[string][]string
}

func readVault(vault string) (*vaultState, error) {
	v := &vaultState{texts: map[string]string{}, notes: map[string]note{}}
	if _, err := os.Stat(vault); errors.Is(err, fs.ErrNotExist) {
		v.reindex()
		return v, nil
	}
	files, err := walkVault(vault)
	if err != nil {
		return nil, err
	}
	for rel := range files {
		data, err := os.ReadFile(filepath.Join(vault, filepath.FromSlash(rel)))
		if err != nil {
			return nil, err
		}
		v.texts[rel] = string(data)
		v.notes[rel] = parseNote(rel, data)
	}
	v.reindex()
	return v, nil
}

func (v *vaultState) reindex() {
	v.dir = newDirectory()
	v.stems = map[string][]string{}
	rels := make([]string, 0, len(v.notes))
	for rel := range v.notes {
		rels = append(rels, rel)
	}
	slices.Sort(rels)
	for _, rel := range rels {
		n := v.notes[rel]
		stem := stemOf(rel)
		v.dir.add(rel, n.id, n.title, stem)
		// Topics repeat file names by design: mission.md, 01-intro.md.
		if n.kind != "learning" {
			v.stems[stem] = append(v.stems[stem], rel)
		}
	}
}

func (v *vaultState) put(rel, text string) {
	v.texts[rel] = text
	v.notes[rel] = parseNote(rel, []byte(text))
	v.reindex()
}

// check returns every problem with one note of the vault state.
func (v *vaultState) check(rel string) []problem {
	var found []problem
	fail := func(format string, args ...any) {
		found = append(found, problem{"error", rel, fmt.Sprintf(format, args...)})
	}
	warn := func(format string, args ...any) {
		found = append(found, problem{"warning", rel, fmt.Sprintf(format, args...)})
	}

	if !inLayout(rel) {
		fail("outside the notebook layout")
	}
	n := v.notes[rel]
	text := strings.ReplaceAll(strings.TrimPrefix(v.texts[rel], "\xef\xbb\xbf"), "\r\n", "\n")
	fields, _, ok := splitFrontmatter(text)
	if !ok {
		fail("missing frontmatter")
	} else {
		for _, line := range malformedLines(text) {
			fail("malformed frontmatter line %q", line)
		}
		if date := fields["date"]; date == "" {
			fail("missing date")
		} else if _, err := time.Parse(time.DateOnly, date); err != nil {
			fail("malformed date %q", date)
		}
		if n.kind == "lessons" && fields["scope"] == "" {
			fail("missing scope")
		}
		if n.supersedes != "" {
			if _, ok := v.notes[n.supersedes]; !ok {
				fail("supersedes %s, which does not exist", n.supersedes)
			}
		}
	}
	if n.kind != "learning" {
		for _, other := range v.stems[stemOf(rel)] {
			if other != rel {
				fail("same name as %s", other)
			}
		}
	}
	for _, target := range n.links {
		switch paths := v.dir.resolve(target); len(paths) {
		case 0:
			warn("broken wikilink [[%s]]", target)
		case 1:
		default:
			warn("ambiguous wikilink [[%s]]: %s", target, strings.Join(paths, ", "))
		}
	}
	return found
}

// malformedLines finds frontmatter lines that are no `key: value`. Indented
// lines and comments belong to YAML the tool does not otherwise read.
func malformedLines(text string) []string {
	head, _, _ := strings.Cut(strings.TrimPrefix(text, "---\n")+"\n", "\n---")
	var bad []string
	for _, line := range strings.Split(head, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") ||
			strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "- ") {
			continue
		}
		if !strings.Contains(line, ":") {
			bad = append(bad, line)
		}
	}
	return bad
}

func report(w io.Writer, problems []problem) (failed bool) {
	for _, p := range problems {
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.level, p.path, p.message)
		failed = failed || p.level == "error"
	}
	return failed
}

func runValidate(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	loc, err := addPathFlags(fs)
	if err != nil {
		return err
	}
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: notebook validate [path] [flags]")
		fs.PrintDefaults()
	}
	positional, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(positional) > 1 {
		return usageError("validate takes at most one path")
	}
	if err := checkVault(loc.vault); err != nil {
		return err
	}
	v, err := readVault(loc.vault)
	if err != nil {
		return err
	}

	var rels []string
	if len(positional) == 1 {
		rel := path.Clean(filepath.ToSlash(positional[0]))
		if _, ok := v.notes[rel]; !ok {
			return fmt.Errorf("no note at %s", rel)
		}
		rels = []string{rel}
	} else {
		for rel := range v.notes {
			rels = append(rels, rel)
		}
		slices.Sort(rels)
	}
	var problems []problem
	for _, rel := range rels {
		problems = append(problems, v.check(rel)...)
	}
	if report(stdout, problems) {
		return errFound
	}
	return nil
}

// errFound is a failure already reported in full, so it exits 1 quietly.
var errFound = errors.New("problems found")
