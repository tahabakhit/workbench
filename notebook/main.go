// Command notebook indexes and searches the markdown vault that workbench
// skills file into. The markdown files are the truth; the index is rebuildable.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const usageText = `usage: notebook <command> [flags]

commands:
  search <query>   ranked snippets
  recall <note>    one whole note
  links <note>     outgoing wikilinks, flagging broken and ambiguous
  backlinks <note> notes that link to it
  file <path>      validate, write and commit a note from stdin
  validate [path]  check one note or the whole vault

run "notebook <command> --help" for flags
`

// usageError is a mistake in the command line, reported with exit code 2.
type usageError string

func (e usageError) Error() string { return string(e) }

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usageText)
		return 2
	}
	var err error
	switch args[0] {
	case "search":
		err = runSearch(args[1:], stdout, stderr)
	case "recall":
		err = runRecall(args[1:], stdout, stderr)
	case "links":
		err = runLinks(args[1:], stdout, stderr)
	case "backlinks":
		err = runBacklinks(args[1:], stdout, stderr)
	case "file":
		err = runFile(args[1:], stdout, stderr)
	case "validate":
		err = runValidate(args[1:], stdout, stderr)
	default:
		err = usageError(fmt.Sprintf("unknown command %q\n%s", args[0], usageText))
	}
	var usage usageError
	switch {
	case err == nil, errors.Is(err, flag.ErrHelp):
		return 0
	case errors.Is(err, errFound):
		return 1
	case errors.Is(err, errBadFlag):
		return 2
	case errors.As(err, &usage):
		fmt.Fprintln(stderr, "notebook:", err)
		return 2
	default:
		fmt.Fprintln(stderr, "notebook:", err)
		return 1
	}
}

// errBadFlag marks a flag error the flag package has already printed.
var errBadFlag = errors.New("bad flag")

// locations are the two locations every command works on.
type locations struct{ vault, index string }

// addPathFlags registers --vault and --index with their defaults.
func addPathFlags(fs *flag.FlagSet) (*locations, error) {
	data := os.Getenv("XDG_DATA_HOME")
	if data == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		data = filepath.Join(home, ".local", "share")
	}
	p := &locations{}
	fs.StringVar(&p.vault, "vault", filepath.Join(data, "workbench", "notebook"), "vault `folder`")
	fs.StringVar(&p.index, "index", filepath.Join(data, "workbench", "state", "notebook.db"), "index `file`")
	return p, nil
}

// parseFlags accepts flags before and after positional arguments, which the
// flag package alone stops reading at the first positional.
func parseFlags(fs *flag.FlagSet, args []string) ([]string, error) {
	var positional []string
	for {
		if err := fs.Parse(args); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return nil, err
			}
			return nil, errBadFlag
		}
		args = fs.Args()
		if len(args) == 0 {
			return positional, nil
		}
		positional = append(positional, args[0])
		args = args[1:]
	}
}
