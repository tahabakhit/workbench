package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestPushOrderAndContent(t *testing.T) {
	e := newEnv(t)
	out, errOut, code := e.run("push", "--repo", "workbench")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	want := `## Notebook

Notes for workbench. Read one with "notebook recall <path>".

- log/2026-10-30-workbench-branch-cleanup.md: Clean up stale branches
- log/2026-10-29-workbench-hooks-01-force-push.md: hooks issue 01, ask on force push and remote branch delete
- lessons/link-source.md: Link source
- lessons/self-reference.md: Self reference
- lessons/slug-reference.md: Slug reference
- lessons/one-change-per-commit.md: One change per commit
`
	if out != want {
		t.Errorf("\n got %q\nwant %q", out, want)
	}
}

func TestPushKeepsTheThreeLatestLogs(t *testing.T) {
	e := newEnv(t)
	for day := 1; day <= 5; day++ {
		name := fmt.Sprintf("log/2026-11-0%d-workbench-x.md", day)
		write(t, filepath.Join(e.vault, name), fmt.Sprintf("---\ndate: 2026-11-0%d\nrepo: workbench\n---\n# Log %d\n", day, day), testNow)
	}
	out, _, _ := e.run("push", "--repo", "workbench")
	for _, want := range []string{"Log 5", "Log 4", "Log 3"} {
		if !strings.Contains(out, want) {
			t.Errorf("lacks %s: %q", want, out)
		}
	}
	for _, not := range []string{"Log 2", "force push", "dropped"} {
		if strings.Contains(out, not) {
			t.Errorf("holds %s: %q", not, out)
		}
	}
}

func TestPushLeavesOutSupersededNotesAndOtherRepos(t *testing.T) {
	out, _, _ := newEnv(t).run("push", "--repo", "workbench")
	for _, not := range []string{"keep-commits-small", "bare-slug-target", "homelab"} {
		if strings.Contains(out, not) {
			t.Errorf("holds %s", not)
		}
	}
}

func addLessons(t *testing.T, e env, prefix, scope string, count int) {
	t.Helper()
	long := strings.Repeat("words ", 30)
	for i := range count {
		name := fmt.Sprintf("lessons/%s-%02d.md", prefix, i)
		write(t, filepath.Join(e.vault, name), fmt.Sprintf("---\ndate: 2026-10-10\nscope: %s\n---\n# %s %02d %s\n", scope, prefix, i, long), testNow)
	}
}

func TestPushDropsGlobalLessonsFirst(t *testing.T) {
	e := newEnv(t)
	addLessons(t, e, "wide", "global", 20)
	out, _, _ := e.run("push", "--repo", "workbench")
	if utf8.RuneCountInString(out) > 2000 {
		t.Errorf("%d characters", utf8.RuneCountInString(out))
	}
	if !strings.Contains(out, "lessons/self-reference.md") || !strings.Contains(out, "force push") {
		t.Errorf("dropped the wrong kind: %q", out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	last := lines[len(lines)-1]
	kept := strings.Count(out, "lessons/wide-")
	if want := fmt.Sprintf("(%d more dropped)", 21-kept); last != want {
		t.Errorf("last line %q, want %q", last, want)
	}
}

func TestPushDropsRepoLessonsBeforeLogs(t *testing.T) {
	e := newEnv(t)
	addLessons(t, e, "wide", "global", 3)
	addLessons(t, e, "near", "workbench", 30)
	out, _, _ := e.run("push", "--repo", "workbench")
	if utf8.RuneCountInString(out) > 2000 {
		t.Errorf("%d characters", utf8.RuneCountInString(out))
	}
	if strings.Contains(out, "lessons/wide-") || strings.Contains(out, "one-change-per-commit") {
		t.Errorf("a global lesson survived: %q", out)
	}
	if !strings.Contains(out, "force push") || !strings.Contains(out, "stale branches") {
		t.Errorf("a log was dropped: %q", out)
	}
	if !strings.Contains(out, "more dropped)") {
		t.Errorf("no dropped line: %q", out)
	}
}

func TestPushDropsLogsLast(t *testing.T) {
	e := newEnv(t)
	long := strings.Repeat("words ", 200)
	write(t, filepath.Join(e.vault, "log/2026-11-01-workbench-long.md"), "---\ndate: 2026-11-01\nrepo: workbench\n---\n# "+long+"\n", testNow)
	write(t, filepath.Join(e.vault, "log/2026-11-02-workbench-long.md"), "---\ndate: 2026-11-02\nrepo: workbench\n---\n# "+long+"\n", testNow)
	out, _, _ := e.run("push", "--repo", "workbench")
	if n := utf8.RuneCountInString(out); n > 2000 {
		t.Errorf("%d characters", n)
	}
	if !strings.Contains(out, "log/2026-11-02") || strings.Contains(out, "log/2026-11-01") {
		t.Errorf("wrong log kept: %q", out)
	}
}

func TestPushEmptyResultPrintsNothing(t *testing.T) {
	e := newEnv(t)
	write(t, filepath.Join(e.vault, "lessons/one-change-per-commit.md"), "---\ndate: 2026-10-01\nscope: other\n---\n# Moved\n", testNow.Add(1))
	write(t, filepath.Join(e.vault, "lessons/keep-commits-small.md"), "---\ndate: 2026-10-01\nscope: other\n---\n# Moved too\n", testNow.Add(1))
	if out, _, code := e.run("push", "--repo", "nothing-here"); code != 0 || out != "" {
		t.Errorf("exit %d, out %q", code, out)
	}
}

func TestPushNeedsARepo(t *testing.T) {
	var out, errOut strings.Builder
	if code := run([]string{"push"}, &out, &errOut); code != 2 {
		t.Errorf("exit %d", code)
	}
}
