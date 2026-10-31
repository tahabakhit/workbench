package main

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	forcePushNote = "log/2026-10-29-workbench-hooks-01-force-push.md"
	cleanupNote   = "log/2026-10-30-workbench-branch-cleanup.md"
	backupNote    = "log/2026-10-28-homelab-backup.md"
	// From: printf %s "$forcePushNote" | sha256sum | cut -c1-12
	forcePushID = "aa15fcfe1153"
)

type env struct {
	vault, index string
}

// newEnv copies the fixture vault so tests can change it freely.
func newEnv(t *testing.T) env {
	t.Helper()
	dir := t.TempDir()
	vault := filepath.Join(dir, "vault")
	if err := os.CopyFS(vault, os.DirFS("testdata/vault")); err != nil {
		t.Fatal(err)
	}
	return env{vault: vault, index: filepath.Join(dir, "state", "notebook.db")}
}

func (e env) run(args ...string) (string, string, int) {
	var out, errOut bytes.Buffer
	args = append(args, "--vault", e.vault, "--index", e.index)
	code := run(args, &out, &errOut)
	return out.String(), errOut.String(), code
}

func (e env) search(t *testing.T, args ...string) []result {
	t.Helper()
	out, errOut, code := e.run(append(append([]string{"search"}, args...), "--json")...)
	if code != 0 {
		t.Fatalf("search %v: exit %d: %s", args, code, errOut)
	}
	var results []result
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("search %v: bad JSON %q: %v", args, out, err)
	}
	return results
}

func paths(results []result) []string {
	found := []string{}
	for _, r := range results {
		found = append(found, r.Path)
	}
	return found
}

func wantPaths(t *testing.T, e env, want []string, args ...string) {
	t.Helper()
	if got := paths(e.search(t, args...)); !reflect.DeepEqual(got, want) {
		t.Errorf("search %v:\n got %v\nwant %v", args, got, want)
	}
}

func write(t *testing.T, path, content string, mtime time.Time) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mtime, mtime); err != nil {
		t.Fatal(err)
	}
}

func TestTitleMatchOutranksBodyMatches(t *testing.T) {
	wantPaths(t, newEnv(t), []string{forcePushNote, cleanupNote, backupNote}, "force push")
}

func TestStemmingAndDiacritics(t *testing.T) {
	e := newEnv(t)
	wantPaths(t, e, []string{forcePushNote, cleanupNote, backupNote}, "forcing pushed")
	wantPaths(t, e, []string{backupNote}, "cafe")
}

func TestEveryWordMustMatch(t *testing.T) {
	wantPaths(t, newEnv(t), []string{}, "force sqlite")
}

func TestFilters(t *testing.T) {
	e := newEnv(t)
	wantPaths(t, e, []string{forcePushNote, cleanupNote}, "force push", "--repo", "workbench")
	wantPaths(t, e, []string{"decisions/2026-10-20-use-sqlite.md"}, "sqlite", "--kind", "decisions")
	wantPaths(t, e, []string{}, "sqlite", "--kind", "log")
	wantPaths(t, e, []string{"corrections.md"}, "gnu", "--kind", "corrections")
	wantPaths(t, e, []string{"learning/sh-basics/lessons/01-quoting.md"}, "quoting", "--kind", "learning")
}

func TestLimit(t *testing.T) {
	wantPaths(t, newEnv(t), []string{forcePushNote}, "force push", "--limit", "1")
}

func TestSupersededNotesAreHiddenUnlessAll(t *testing.T) {
	e := newEnv(t)
	wantPaths(t, e, []string{"lessons/one-change-per-commit.md"}, "rebase")
	wantPaths(t, e, []string{"lessons/keep-commits-small.md", "lessons/one-change-per-commit.md"}, "rebase", "--all")
}

func TestSupersedesBySlugAndWikilink(t *testing.T) {
	e := newEnv(t)
	wantPaths(t, e, []string{"lessons/slug-reference.md"}, "zebra")
}

func TestNoteDoesNotHideItself(t *testing.T) {
	wantPaths(t, newEnv(t), []string{"lessons/self-reference.md"}, "walrus")
}

func TestSupersededNoteReturnsWhenItsReplacementGoes(t *testing.T) {
	e := newEnv(t)
	e.search(t, "rebase")
	if err := os.Remove(filepath.Join(e.vault, "lessons/one-change-per-commit.md")); err != nil {
		t.Fatal(err)
	}
	wantPaths(t, e, []string{"lessons/keep-commits-small.md"}, "rebase")
}

func TestHiddenFoldersAndOtherFilesAreNotIndexed(t *testing.T) {
	wantPaths(t, newEnv(t), []string{}, "platypus")
}

func TestTitleFallsBackToFileNameAndSkipsCodeFences(t *testing.T) {
	e := newEnv(t)
	got := e.search(t, "quoting")
	if len(got) != 1 || got[0].Title != "01-quoting" {
		t.Errorf("no-heading note: %+v", got)
	}
	got = e.search(t, "sqlite", "--kind", "decisions")
	if len(got) != 1 || got[0].Title != "Use SQLite for the index" {
		t.Errorf("decision note: %+v", got)
	}
}

func TestJSONFields(t *testing.T) {
	got := newEnv(t).search(t, "force push")[0]
	if got.Path != forcePushNote || got.ID != forcePushID || got.Kind != "log" ||
		got.Date != "2026-10-29" || got.Repo != "workbench" ||
		got.Title != "hooks issue 01, ask on force push and remote branch delete" {
		t.Errorf("fields: %+v", got)
	}
	if got.Score <= 0 {
		t.Errorf("score %v: want positive", got.Score)
	}
	if !strings.Contains(got.Snippet, "**force** **push**") {
		t.Errorf("snippet %q: want the match marked", got.Snippet)
	}
}

func TestSnippetIsOneLine(t *testing.T) {
	for _, r := range newEnv(t).search(t, "branches") {
		if strings.ContainsAny(r.Snippet, "\r\n") {
			t.Errorf("snippet %q holds a line break", r.Snippet)
		}
	}
}

func TestPlainOutput(t *testing.T) {
	out, errOut, code := newEnv(t).run("search", "force push")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	first := strings.Fields(strings.SplitN(out, "\n", 2)[0])
	if len(first) != 3 || first[0] != forcePushNote || first[1] != forcePushID {
		t.Errorf("first line: %q", first)
	}
	if !strings.Contains(out, "\n  ") {
		t.Errorf("no indented snippet: %q", out)
	}
}

func TestNoMatchesIsNotAnError(t *testing.T) {
	out, errOut, code := newEnv(t).run("search", "nothingmatchesthis")
	if code != 0 || out != "" || !strings.Contains(errOut, "no matches") {
		t.Errorf("exit %d, out %q, err %q", code, out, errOut)
	}
	out, _, code = newEnv(t).run("search", "nothingmatchesthis", "--json")
	if code != 0 || strings.TrimSpace(out) != "[]" {
		t.Errorf("JSON: exit %d, out %q", code, out)
	}
}

func TestQueryPunctuationIsNotSearchSyntax(t *testing.T) {
	e := newEnv(t)
	for _, q := range []string{"hooks-01", `force"push`, "force*", "title:force", "force OR", "(force", "NEAR(force", "force AND"} {
		if _, errOut, code := e.run("search", q); code != 0 {
			t.Errorf("query %q: exit %d: %s", q, code, errOut)
		}
	}
	wantPaths(t, e, []string{forcePushNote, cleanupNote, backupNote}, "force - push")
}

func TestUsageErrors(t *testing.T) {
	for _, args := range [][]string{{}, {"nonsense"}, {"search"}, {"search", "---"}, {"search", "x", "--limit", "0"}, {"search", "--nope"}} {
		var out, errOut bytes.Buffer
		if code := run(args, &out, &errOut); code != 2 {
			t.Errorf("%v: exit %d, want 2", args, code)
		}
		if errOut.Len() == 0 {
			t.Errorf("%v: no message", args)
		}
	}
}

func TestHelpSucceeds(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"search", "--help"}, &out, &errOut); code != 0 {
		t.Errorf("exit %d", code)
	}
}

func TestMissingVaultFailsAndCreatesNoIndex(t *testing.T) {
	e := newEnv(t)
	e.vault = filepath.Join(e.vault, "absent")
	_, errOut, code := e.run("search", "force")
	if code != 1 || !strings.Contains(errOut, e.vault) {
		t.Errorf("exit %d: %q", code, errOut)
	}
	if _, err := os.Stat(e.index); !os.IsNotExist(err) {
		t.Errorf("index created: %v", err)
	}
}

func TestDefaultPathsFollowDataHome(t *testing.T) {
	data := t.TempDir()
	t.Setenv("XDG_DATA_HOME", data)
	vault := filepath.Join(data, "workbench", "notebook")
	if err := os.CopyFS(vault, os.DirFS("testdata/vault")); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := run([]string{"search", "force push"}, &out, &errOut); code != 0 {
		t.Fatalf("exit %d: %s", code, errOut.String())
	}
	if !strings.HasPrefix(out.String(), forcePushNote) {
		t.Errorf("out %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(data, "workbench", "state", "notebook.db")); err != nil {
		t.Error(err)
	}
}

func TestEditedNoteIsReindexedByModificationTime(t *testing.T) {
	e := newEnv(t)
	wantPaths(t, e, []string{}, "aardvark")
	path := filepath.Join(e.vault, forcePushNote)
	write(t, path, "---\ndate: 2026-10-29\n---\n\n# Renamed\n\nAn aardvark now.\n", time.Now().Add(time.Hour))
	wantPaths(t, e, []string{forcePushNote}, "aardvark")
	wantPaths(t, e, []string{cleanupNote, backupNote}, "force push")
}

func TestUnchangedModificationTimeIsNotReindexed(t *testing.T) {
	e := newEnv(t)
	e.search(t, "force")
	path := filepath.Join(e.vault, forcePushNote)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	write(t, path, "# Aardvark\n", info.ModTime())
	wantPaths(t, e, []string{}, "aardvark")
}

func TestAddedAndDeletedNotes(t *testing.T) {
	e := newEnv(t)
	e.search(t, "force")
	write(t, filepath.Join(e.vault, "lessons/new.md"), "# Aardvark\n", time.Now())
	if err := os.Remove(filepath.Join(e.vault, forcePushNote)); err != nil {
		t.Fatal(err)
	}
	wantPaths(t, e, []string{"lessons/new.md"}, "aardvark")
	wantPaths(t, e, []string{cleanupNote, backupNote}, "force push")
}

func TestDeletedIndexGivesTheSameResults(t *testing.T) {
	e := newEnv(t)
	for _, args := range [][]string{
		{"search", "force push"},
		{"search", "force push", "--json"},
		{"search", "rebase", "--all", "--json"},
	} {
		before, _, _ := e.run(args...)
		if err := os.Remove(e.index); err != nil {
			t.Fatal(err)
		}
		after, _, _ := e.run(args...)
		if before == "" || before != after {
			t.Errorf("%v:\nbefore %q\nafter  %q", args, before, after)
		}
	}
}

func TestOutdatedIndexSchemaIsRebuilt(t *testing.T) {
	e := newEnv(t)
	e.search(t, "force")
	db, err := openIndex(e.index)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DROP TABLE notes_fts; PRAGMA user_version = 0"); err != nil {
		t.Fatal(err)
	}
	db.Close()
	wantPaths(t, e, []string{forcePushNote, cleanupNote, backupNote}, "force push")
}

func TestConcurrentFirstRuns(t *testing.T) {
	e := newEnv(t)
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, errOut, code := e.run("search", "force push")
			if code != 0 || !strings.HasPrefix(out, forcePushNote) {
				t.Errorf("exit %d, out %q, err %q", code, out, errOut)
			}
		}()
	}
	wg.Wait()
}

func TestVaultIsNotWritten(t *testing.T) {
	e := newEnv(t)
	before := snapshot(t, e.vault)
	e.search(t, "force push")
	if after := snapshot(t, e.vault); !reflect.DeepEqual(before, after) {
		t.Errorf("vault changed:\nbefore %v\nafter  %v", before, after)
	}
}

func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		files[path] = string(data)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func TestRecall(t *testing.T) {
	e := newEnv(t)
	whole, err := os.ReadFile("testdata/vault/" + forcePushNote)
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{
		forcePushNote,
		strings.TrimSuffix(forcePushNote, ".md"),
		"2026-10-29-workbench-hooks-01-force-push",
		"hooks issue 01, ask on force push and remote branch delete",
		forcePushID,
	} {
		out, errOut, code := e.run("recall", q)
		if code != 0 || out != string(whole) {
			t.Errorf("recall %q: exit %d, err %q, out %q", q, code, errOut, out)
		}
	}
}

func TestRecallTitleInSeveralArguments(t *testing.T) {
	if _, errOut, code := newEnv(t).run("recall", "Keep", "commits", "small"); code != 0 {
		t.Errorf("exit %d: %s", code, errOut)
	}
}

func TestRecallNoMatchListsCandidates(t *testing.T) {
	out, errOut, code := newEnv(t).run("recall", "use sqlite")
	if code != 1 || out != "" {
		t.Errorf("exit %d, out %q", code, out)
	}
	if !strings.Contains(errOut, "no note matches") || !strings.Contains(errOut, "decisions/2026-10-20-use-sqlite.md") {
		t.Errorf("err %q", errOut)
	}
}

func TestRecallTwoMatchesListsBoth(t *testing.T) {
	out, errOut, code := newEnv(t).run("recall", "Duplicate title")
	if code != 1 || out != "" {
		t.Errorf("exit %d, out %q", code, out)
	}
	for _, want := range []string{"more than one", "decisions/2026-10-21-duplicate-a.md", "log/2026-10-21-workbench-duplicate-b.md"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("err %q lacks %q", errOut, want)
		}
	}
}

func TestRecallSeesAnEditAndKeepsSupersededNotes(t *testing.T) {
	e := newEnv(t)
	if _, errOut, code := e.run("recall", "lessons/keep-commits-small.md"); code != 0 {
		t.Errorf("superseded: exit %d: %s", code, errOut)
	}
	write(t, filepath.Join(e.vault, forcePushNote), "# Fresh\n", time.Now().Add(time.Hour))
	if out, _, _ := e.run("recall", "Fresh"); out != "# Fresh\n" {
		t.Errorf("out %q", out)
	}
}

func links(t *testing.T, e env, note string) []link {
	t.Helper()
	out, errOut, code := e.run("links", note, "--json")
	if code != 0 {
		t.Fatalf("links %s: exit %d: %s", note, code, errOut)
	}
	var got []link
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("bad JSON %q: %v", out, err)
	}
	return got
}

func TestLinks(t *testing.T) {
	want := []link{
		{Target: "Duplicate title", Status: "ambiguous", Paths: []string{"decisions/2026-10-21-duplicate-a.md", "log/2026-10-21-workbench-duplicate-b.md"}},
		{Target: "Keep commits small", Status: "ok", Paths: []string{"lessons/keep-commits-small.md"}},
		{Target: "decisions/2026-10-20-use-sqlite.md", Status: "ok", Paths: []string{"decisions/2026-10-20-use-sqlite.md"}},
		{Target: "keep-commits-small", Status: "ok", Paths: []string{"lessons/keep-commits-small.md"}},
		{Target: "later-note", Status: "broken", Paths: []string{}},
		{Target: "lessons/one-change-per-commit", Status: "ok", Paths: []string{"lessons/one-change-per-commit.md"}},
		{Target: "no-such-note", Status: "broken", Paths: []string{}},
	}
	if got := links(t, newEnv(t), "Link source"); !reflect.DeepEqual(got, want) {
		t.Errorf("\n got %+v\nwant %+v", got, want)
	}
}

func TestLinkResolvesOnceItsTargetIsAdded(t *testing.T) {
	e := newEnv(t)
	status := func() string {
		for _, l := range links(t, e, "Link source") {
			if l.Target == "later-note" {
				return l.Status + " " + strings.Join(l.Paths, ",")
			}
		}
		return "missing"
	}
	if got := status(); got != "broken " {
		t.Fatalf("before: %q", got)
	}
	write(t, filepath.Join(e.vault, "lessons/later-note.md"), "# Later\n", time.Now())
	if got := status(); got != "ok lessons/later-note.md" {
		t.Errorf("after: %q", got)
	}
}

func TestLinkStaysUnresolvedWhenASecondMatchAppears(t *testing.T) {
	e := newEnv(t)
	write(t, filepath.Join(e.vault, "log/keep-commits-small.md"), "# Other\n", time.Now())
	for _, l := range links(t, e, "Link source") {
		if l.Target == "keep-commits-small" && l.Status != "ambiguous" {
			t.Errorf("%+v", l)
		}
	}
}

func TestPathTierBeatsNameTier(t *testing.T) {
	e := newEnv(t)
	write(t, filepath.Join(e.vault, "lessons/name.md"), "# lessons/keep-commits-small\n", time.Now())
	if out, _, code := e.run("recall", "lessons/keep-commits-small"); code != 0 || !strings.Contains(out, "# Keep commits small") {
		t.Errorf("exit %d, out %q", code, out)
	}
}

func TestLinksPlainOutput(t *testing.T) {
	out, errOut, code := newEnv(t).run("links", "Link source")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	for _, want := range []string{
		"broken\tno-such-note\n",
		"ambiguous\tDuplicate title\tdecisions/2026-10-21-duplicate-a.md, log/2026-10-21-workbench-duplicate-b.md\n",
		"ok\tkeep-commits-small\tlessons/keep-commits-small.md\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("out %q lacks %q", out, want)
		}
	}
}

func TestBacklinks(t *testing.T) {
	e := newEnv(t)
	out, errOut, code := e.run("backlinks", "keep-commits-small")
	if code != 0 || out != "lessons/link-source.md\n" {
		t.Errorf("exit %d, err %q, out %q", code, errOut, out)
	}
	out, _, code = e.run("backlinks", "keep-commits-small", "--json")
	var got []string
	if err := json.Unmarshal([]byte(out), &got); code != 0 || err != nil || !reflect.DeepEqual(got, []string{"lessons/link-source.md"}) {
		t.Errorf("JSON: exit %d, err %v, out %q", code, err, out)
	}
}

func TestAmbiguousLinkIsNoBacklink(t *testing.T) {
	out, errOut, code := newEnv(t).run("backlinks", "decisions/2026-10-21-duplicate-a.md")
	if code != 0 || out != "" || !strings.Contains(errOut, "no backlinks") {
		t.Errorf("exit %d, err %q, out %q", code, errOut, out)
	}
}

func TestLinksAndBacklinksNeedExactlyOneNote(t *testing.T) {
	e := newEnv(t)
	for _, cmd := range []string{"links", "backlinks", "recall"} {
		if _, _, code := e.run(cmd, "Duplicate title"); code != 1 {
			t.Errorf("%s two matches: exit %d", cmd, code)
		}
		if _, _, code := e.run(cmd, "no-such-note"); code != 1 {
			t.Errorf("%s no match: exit %d", cmd, code)
		}
		var out, errOut bytes.Buffer
		if code := run([]string{cmd}, &out, &errOut); code != 2 {
			t.Errorf("%s no argument: exit %d", cmd, code)
		}
	}
}

func TestDeletedIndexKeepsLinks(t *testing.T) {
	e := newEnv(t)
	before := links(t, e, "Link source")
	if err := os.Remove(e.index); err != nil {
		t.Fatal(err)
	}
	if after := links(t, e, "Link source"); !reflect.DeepEqual(before, after) {
		t.Errorf("\nbefore %+v\nafter  %+v", before, after)
	}
}
