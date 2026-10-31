package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// stdin is where `file` reads the note; tests replace it.
var stdin io.Reader = os.Stdin

func runFile(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("file", flag.ContinueOnError)
	fs.SetOutput(stderr)
	loc, err := addPathFlags(fs)
	if err != nil {
		return err
	}
	update := fs.Bool("update", false, "replace an existing corrections, goals, plan or topic file")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: notebook file <path> [flags] < note.md")
		fs.PrintDefaults()
	}
	positional, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(positional) != 1 {
		return usageError("file takes one path")
	}
	rel := positional[0]
	if cleaned := path.Clean(rel); !strings.HasPrefix(cleaned, "../") {
		rel = cleaned
	}
	text, err := io.ReadAll(stdin)
	if err != nil {
		return err
	}

	v, err := readVault(loc.vault)
	if err != nil {
		return err
	}
	old, existed := v.texts[rel]
	refuse := func(message string) error {
		report(stderr, []problem{{"error", rel, message}})
		return errFound
	}
	switch {
	case !inLayout(rel):
		return refuse("outside the notebook layout")
	case *update && !updatable.MatchString(rel):
		return refuse("cannot be updated: only corrections, goals, plans and topic files change after filing")
	case existed && !*update:
		return refuse("already exists")
	}
	v.put(rel, string(text))
	problems := v.check(rel)
	if report(stderr, problems) {
		return errFound
	}
	if existed && old == string(text) {
		fmt.Fprintln(stderr, "unchanged")
		return nil
	}

	if err := commitNote(loc.vault, rel, text, existed); err != nil {
		return err
	}
	fmt.Fprintln(stdout, rel)
	return nil
}

// commitNote writes the note and commits only that path. If any step fails the
// file goes back to how it was, so the vault never holds an uncommitted note.
func commitNote(vault, rel string, text []byte, existed bool) (err error) {
	if err := os.MkdirAll(vault, 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(vault, ".git")); errors.Is(err, fs.ErrNotExist) {
		if err := gitRun(vault, "init", "-q"); err != nil {
			return err
		}
	}
	target := filepath.Join(vault, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	var old []byte
	if existed {
		if old, err = os.ReadFile(target); err != nil {
			return err
		}
	}
	if err := writeAtomic(target, text); err != nil {
		return err
	}
	defer func() {
		if err == nil {
			return
		}
		// Best effort: an unborn branch has nothing to reset to.
		_ = gitRun(vault, "reset", "-q", "--", rel)
		restore := os.Remove(target)
		if existed {
			restore = writeAtomic(target, old)
		}
		err = errors.Join(err, restore)
	}()

	verb := "File"
	if existed {
		verb = "Update"
	}
	if err := gitRun(vault, "add", "--", rel); err != nil {
		return err
	}
	return gitRun(vault, "commit", "-q", "-m", verb+" "+rel, "--", rel)
}

func writeAtomic(target string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(target), ".notebook-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), target)
}

func gitRun(dir string, args ...string) error {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(output.String()))
	}
	return nil
}
