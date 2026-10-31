package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const goodLog = "---\ndate: 2026-10-31\nrepo: workbench\n---\n\n# A finished task\n\nDone.\n"

// gitVault is a fixture copy in a temporary repository with one commit.
func gitVault(t *testing.T) env {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	t.Setenv("GIT_AUTHOR_NAME", "Test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.invalid")
	t.Setenv("GIT_COMMITTER_NAME", "Test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.invalid")
	e := newEnv(t)
	git(t, e.vault, "init", "-q")
	git(t, e.vault, "add", ".")
	git(t, e.vault, "commit", "-qm", "Fixture")
	return e
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
	return string(out)
}

func (e env) file(content string, args ...string) (string, string, int) {
	old := stdin
	defer func() { stdin = old }()
	stdin = strings.NewReader(content)
	return e.run(append([]string{"file"}, args...)...)
}

func commits(t *testing.T, e env) string {
	t.Helper()
	return strings.TrimSpace(git(t, e.vault, "rev-list", "--count", "HEAD"))
}

func wantRefused(t *testing.T, e env, content, message string, args ...string) {
	t.Helper()
	before, count := snapshot(t, e.vault), commits(t, e)
	out, errOut, code := e.file(content, args...)
	if code != 1 || !strings.Contains(errOut, message) {
		t.Errorf("file %v: exit %d, out %q, err %q, want %q", args, code, out, errOut, message)
	}
	if !reflect.DeepEqual(before, snapshot(t, e.vault)) || count != commits(t, e) {
		t.Errorf("file %v: a refusal changed the vault", args)
	}
}

func TestFileWritesOneCommitHoldingOnlyTheNote(t *testing.T) {
	e := gitVault(t)
	// An unrelated staged change must stay out of the commit.
	if err := os.WriteFile(filepath.Join(e.vault, "stray.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, e.vault, "add", "stray.txt")
	rel := "log/2026-10-31-workbench-finished-task.md"
	out, errOut, code := e.file(goodLog, rel)
	if code != 0 || strings.TrimSpace(out) != rel {
		t.Fatalf("exit %d, out %q, err %q", code, out, errOut)
	}
	if got := commits(t, e); got != "2" {
		t.Errorf("commits: %s", got)
	}
	if got := strings.TrimSpace(git(t, e.vault, "show", "--name-only", "--format=%s", "HEAD")); got != "File "+rel+"\n\n"+rel {
		t.Errorf("commit: %q", got)
	}
	data, err := os.ReadFile(filepath.Join(e.vault, rel))
	if err != nil || string(data) != goodLog {
		t.Errorf("content %q, err %v", data, err)
	}
	wantPaths(t, e, []string{rel}, "finished task")
}

func TestFileCreatesTheRepositoryAndFolder(t *testing.T) {
	e := gitVault(t)
	e.vault = filepath.Join(t.TempDir(), "new", "vault")
	rel := "decisions/2026-10-31-first.md"
	if _, errOut, code := e.file(goodLog, rel); code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if got := commits(t, e); got != "1" {
		t.Errorf("commits: %s", got)
	}
}

func TestFileRefusesPathsOutsideTheLayout(t *testing.T) {
	e := gitVault(t)
	for _, rel := range []string{
		"notes.md", "log/x.md", "log/2026-10-31-norepo.md", "log/2026-13-45-repo-slug.md",
		"log/sub/2026-10-31-repo-slug.md", "plans/today.md", "lessons/a/b.md", "learning/t/other.md",
		"../outside.md", "/etc/passwd.md", "lessons/../../outside.md", "lessons/x.txt", "lessons/.hidden.md",
	} {
		wantRefused(t, e, goodLog, "outside the notebook layout", rel)
	}
}

func TestFileAcceptsEveryLayoutPath(t *testing.T) {
	e := gitVault(t)
	lesson := "---\ndate: 2026-10-31\nscope: global\n---\n\n# L\n"
	for rel, content := range map[string]string{
		"goals.md":                      goodLog,
		"plans/2026-10-31.md":           goodLog,
		"lessons/new-lesson.md":         lesson,
		"learning/t/mission.md":         goodLog,
		"learning/t/sources.md":         goodLog,
		"learning/t/quiz.md":            goodLog,
		"learning/t/lessons/01-a.md":    goodLog,
		"learning/u/lessons/01-a.md":    goodLog,
		"learning/u/mission.md":         goodLog,
		"log/2026-10-31-my.repo-a-b.md": goodLog,
	} {
		if _, errOut, code := e.file(content, rel); code != 0 {
			t.Errorf("%s: exit %d: %s", rel, code, errOut)
		}
	}
}

func TestFileRefusesBadFrontmatter(t *testing.T) {
	e := gitVault(t)
	rel := "log/2026-10-31-workbench-x.md"
	for content, message := range map[string]string{
		"# No frontmatter\n":                  "missing frontmatter",
		"":                                    "missing frontmatter",
		"---\ndate: 2026-10-31\n# Unclosed\n": "missing frontmatter",
		"---\nrepo: x\n---\n# T\n":            "missing date",
		"---\ndate: soon\n---\n# T\n":         "malformed date",
		"---\ndate: 2026-02-30\n---\n# T\n":   "malformed date",
		"---\ndate: 2026-10-31\nno colon here\n---\n# T\n": "malformed frontmatter",
	} {
		wantRefused(t, e, content, message, rel)
	}
	wantRefused(t, e, "---\ndate: 2026-10-31\n---\n# T\n", "missing scope", "lessons/no-scope.md")
}

func TestFileRefusesMissingSupersedesTarget(t *testing.T) {
	e := gitVault(t)
	wantRefused(t, e, "---\ndate: 2026-10-31\nscope: global\nsupersedes: lessons/absent.md\n---\n# T\n",
		"supersedes lessons/absent.md", "lessons/replacement.md")
	good := "---\ndate: 2026-10-31\nscope: global\nsupersedes: keep-commits-small\n---\n# T\n"
	if _, errOut, code := e.file(good, "lessons/replacement.md"); code != 0 {
		t.Errorf("exit %d: %s", code, errOut)
	}
}

func TestFileRefusesAnExistingName(t *testing.T) {
	e := gitVault(t)
	wantRefused(t, e, goodLog, "already exists", forcePushNote)
	// The same name in another folder would make wikilinks to it ambiguous.
	wantRefused(t, e, "---\ndate: 2026-10-31\nscope: global\n---\n# T\n",
		"same name as decisions/2026-10-20-use-sqlite.md", "lessons/2026-10-20-use-sqlite.md")
}

func TestFileWarnsOnBrokenAndAmbiguousLinksAndStillFiles(t *testing.T) {
	e := gitVault(t)
	rel := "log/2026-10-31-workbench-links.md"
	content := goodLog + "\n[[no-such-note]] [[Duplicate title]] [[keep-commits-small]]\n"
	_, errOut, code := e.file(content, rel)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	for _, want := range []string{"broken wikilink [[no-such-note]]", "ambiguous wikilink [[Duplicate title]]"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("err %q lacks %q", errOut, want)
		}
	}
	if strings.Contains(errOut, "keep-commits-small") {
		t.Errorf("err %q warns about a good link", errOut)
	}
}

func TestFileUpdate(t *testing.T) {
	e := gitVault(t)
	first := "---\ndate: 2026-10-31\n---\n\n# Goals\n"
	if _, errOut, code := e.file(first, "goals.md"); code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	wantRefused(t, e, first+"more\n", "already exists", "goals.md")
	if _, errOut, code := e.file(first+"more\n", "goals.md", "--update"); code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if got := git(t, e.vault, "log", "-1", "--format=%s"); got != "Update goals.md\n" {
		t.Errorf("subject %q", got)
	}
	count := commits(t, e)
	if _, errOut, code := e.file(first+"more\n", "goals.md", "--update"); code != 0 {
		t.Errorf("same content: exit %d: %s", code, errOut)
	}
	if commits(t, e) != count {
		t.Error("same content made a commit")
	}
}

func TestUpdateCoversOnlyLivingFiles(t *testing.T) {
	e := gitVault(t)
	wantRefused(t, e, goodLog, "cannot be updated", forcePushNote, "--update")
	wantRefused(t, e, goodLog, "cannot be updated", "learning/t/lessons/01-a.md", "--update")
	for _, rel := range []string{"corrections.md", "plans/2026-10-31.md", "learning/t/quiz.md", "learning/t/sources.md", "learning/t/mission.md"} {
		if _, errOut, code := e.file(goodLog, rel, "--update"); code != 0 {
			t.Errorf("%s: exit %d: %s", rel, code, errOut)
		}
	}
}

func failCommits(t *testing.T, e env) {
	t.Helper()
	hook := filepath.Join(e.vault, ".git", "hooks", "pre-commit")
	if err := os.MkdirAll(filepath.Dir(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hook, []byte("#!/bin/sh\necho hook said no >&2\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestFailedCommitLeavesNothingBehind(t *testing.T) {
	e := gitVault(t)
	failCommits(t, e)
	rel := "log/2026-10-31-workbench-unlucky.md"
	_, errOut, code := e.file(goodLog, rel)
	if code != 1 || !strings.Contains(errOut, "hook said no") {
		t.Errorf("exit %d: %q", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(e.vault, rel)); !os.IsNotExist(err) {
		t.Errorf("file left: %v", err)
	}
	if got := git(t, e.vault, "status", "--porcelain"); got != "" {
		t.Errorf("status %q", got)
	}
}

func TestFailedCommitRestoresAnUpdatedFile(t *testing.T) {
	e := gitVault(t)
	rel := "plans/2026-10-31.md"
	if _, errOut, code := e.file(goodLog, rel); code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	failCommits(t, e)
	if _, _, code := e.file(goodLog+"more\n", rel, "--update"); code != 1 {
		t.Errorf("exit %d", code)
	}
	data, err := os.ReadFile(filepath.Join(e.vault, rel))
	if err != nil || string(data) != goodLog {
		t.Errorf("content %q, err %v", data, err)
	}
	if got := git(t, e.vault, "status", "--porcelain"); got != "" {
		t.Errorf("status %q", got)
	}
}

func TestFailedFirstCommitInANewRepository(t *testing.T) {
	e := gitVault(t)
	e.vault = t.TempDir()
	git(t, e.vault, "init", "-q")
	failCommits(t, e)
	if _, _, code := e.file(goodLog, "goals.md"); code != 1 {
		t.Errorf("exit %d", code)
	}
	if got := git(t, e.vault, "status", "--porcelain"); got != "" {
		t.Errorf("status %q", got)
	}
}

func TestFileUsageErrors(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"file"}, &out, &errOut); code != 2 {
		t.Errorf("no path: exit %d", code)
	}
	if code := run([]string{"file", "a.md", "b.md"}, &out, &errOut); code != 2 {
		t.Errorf("two paths: exit %d", code)
	}
}

func TestValidateVaultReportsEachProblemByPath(t *testing.T) {
	e := newEnv(t)
	out, _, code := e.run("validate")
	want := []string{
		"error\tcorrections.md\tmissing frontmatter",
		"warning\tlessons/link-source.md\tambiguous wikilink [[Duplicate title]]: decisions/2026-10-21-duplicate-a.md, log/2026-10-21-workbench-duplicate-b.md",
		"warning\tlessons/link-source.md\tbroken wikilink [[later-note]]",
		"warning\tlessons/link-source.md\tbroken wikilink [[no-such-note]]",
		"error\tlearning/sh-basics/lessons/01-quoting.md\tmalformed date \"not-a-date\"",
		"error\tlearning/sh-basics/lessons/01-quoting.md\tmalformed frontmatter line \"broken line without a colon\"",
	}
	got := strings.Split(strings.TrimSpace(out), "\n")
	if code != 1 || !sameItems(got, want) {
		t.Errorf("exit %d:\n got %q\nwant %q", code, got, want)
	}
}

func sameItems(a, b []string) bool {
	seen := map[string]int{}
	for _, s := range a {
		seen[s]++
	}
	for _, s := range b {
		seen[s]--
	}
	for _, n := range seen {
		if n != 0 {
			return false
		}
	}
	return true
}

func TestValidateOneNote(t *testing.T) {
	e := newEnv(t)
	if out, errOut, code := e.run("validate", forcePushNote); code != 0 || out != "" {
		t.Errorf("good note: exit %d, out %q, err %q", code, out, errOut)
	}
	out, _, code := e.run("validate", "corrections.md")
	if code != 1 || !strings.Contains(out, "missing frontmatter") {
		t.Errorf("bad note: exit %d, out %q", code, out)
	}
	if out, _, code := e.run("validate", "lessons/link-source.md"); code != 0 || !strings.Contains(out, "warning") {
		t.Errorf("warnings alone: exit %d, out %q", code, out)
	}
	if _, errOut, code := e.run("validate", "lessons/absent.md"); code != 1 || !strings.Contains(errOut, "absent") {
		t.Errorf("absent note: exit %d, err %q", code, errOut)
	}
}

func TestValidateReportsPathsOutsideTheLayoutAndSharedNames(t *testing.T) {
	e := newEnv(t)
	write(t, filepath.Join(e.vault, "README.md"), goodLog, testNow)
	write(t, filepath.Join(e.vault, "lessons/2026-10-20-use-sqlite.md"), "---\ndate: 2026-10-31\nscope: x\n---\n# T\n", testNow)
	out, _, _ := e.run("validate")
	for _, want := range []string{
		"error\tREADME.md\toutside the notebook layout",
		"error\tlessons/2026-10-20-use-sqlite.md\tsame name as decisions/2026-10-20-use-sqlite.md",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("out %q lacks %q", out, want)
		}
	}
}

func TestValidateDoesNotWriteTheVault(t *testing.T) {
	e := newEnv(t)
	before := snapshot(t, e.vault)
	e.run("validate")
	if !reflect.DeepEqual(before, snapshot(t, e.vault)) {
		t.Error("vault changed")
	}
}
