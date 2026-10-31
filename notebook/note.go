package main

import (
	"crypto/sha256"
	"encoding/hex"
	"path"
	"strings"
)

// note is what the index keeps from one markdown file.
type note struct {
	id, kind, date, repo, scope string
	supersedes                  string // vault-relative path, or empty
	title, body                 string
}

// noteID is derived from the path alone, so it survives edits and rebuilds.
func noteID(rel string) string {
	sum := sha256.Sum256([]byte(rel))
	return hex.EncodeToString(sum[:])[:12]
}

// parseNote never fails: checking notes is the job of `validate`, and a note
// with odd frontmatter must still be findable.
func parseNote(rel string, data []byte) note {
	text := strings.TrimPrefix(string(data), "\xef\xbb\xbf")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	fields, body := splitFrontmatter(text)
	title, body := splitTitle(body)
	if title == "" {
		title = strings.TrimSuffix(path.Base(rel), ".md")
	}
	return note{
		id:         noteID(rel),
		kind:       kindOf(rel),
		date:       fields["date"],
		repo:       fields["repo"],
		scope:      fields["scope"],
		supersedes: supersededPath(rel, fields["supersedes"]),
		title:      title,
		body:       body,
	}
}

// kindOf is the top folder, or the file's own name for a note at the vault
// root (corrections.md, goals.md).
func kindOf(rel string) string {
	if top, _, nested := strings.Cut(rel, "/"); nested {
		return top
	}
	return strings.TrimSuffix(rel, ".md")
}

// splitFrontmatter reads the flat `key: value` lines the notebook uses. Text
// without a closing delimiter is all body.
func splitFrontmatter(text string) (map[string]string, string) {
	fields := map[string]string{}
	rest, ok := strings.CutPrefix(text, "---\n")
	if !ok {
		return fields, text
	}
	head, body, ok := strings.Cut("\n"+rest, "\n---\n")
	if !ok {
		return fields, text
	}
	for _, line := range strings.Split(head, "\n") {
		if key, value, found := strings.Cut(line, ":"); found {
			fields[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
		}
	}
	return fields, body
}

// splitTitle takes the first level-one heading outside a code fence out of the
// body, so its words are not ranked or shown twice.
func splitTitle(body string) (string, string) {
	lines := strings.Split(body, "\n")
	fenced := false
	for i, line := range lines {
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			fenced = !fenced
		}
		if title, ok := strings.CutPrefix(line, "# "); ok && !fenced {
			return strings.TrimSpace(title), strings.Join(append(lines[:i:i], lines[i+1:]...), "\n")
		}
	}
	return "", body
}

// supersededPath turns a `supersedes` value into a vault-relative path. The
// value can be a path, a wikilink, or a bare name in the note's own folder.
func supersededPath(from, value string) string {
	target := strings.TrimSuffix(strings.TrimPrefix(value, "[["), "]]")
	target, _, _ = strings.Cut(target, "|")
	target = strings.TrimSpace(target)
	if target == "" {
		return ""
	}
	if !strings.HasSuffix(target, ".md") {
		target += ".md"
	}
	if !strings.Contains(target, "/") {
		target = path.Join(path.Dir(from), target)
	}
	return target
}
