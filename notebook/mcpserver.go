package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// runMCP serves the read commands over stdio. Filing stays on the CLI, so this
// server has no write tool. Each call refreshes the index, like a CLI command.
func runMCP(args []string, stderr io.Writer) error {
	fs := flag.NewFlagSet("mcp", flag.ContinueOnError)
	fs.SetOutput(stderr)
	loc, err := addPathFlags(fs)
	if err != nil {
		return err
	}
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: notebook mcp [flags]")
		fs.PrintDefaults()
	}
	if extra, err := parseFlags(fs, args); err != nil {
		return err
	} else if len(extra) > 0 {
		return usageError("mcp takes no arguments")
	}
	return server.ServeStdio(newServer(loc))
}

func newServer(loc *locations) *server.MCPServer {
	s := server.NewMCPServer("notebook", "1.0.0", server.WithToolCapabilities(false))
	note := mcp.WithString("note", mcp.Required(), mcp.Description("Vault path, title, file name or ID"))
	readOnly := mcp.WithReadOnlyHintAnnotation(true)

	s.AddTool(mcp.NewTool("search",
		mcp.WithDescription("Search the notebook of past work, decisions and lessons. Returns ranked notes with snippets."),
		mcp.WithString("query", mcp.Required(), mcp.Description("Words that must all appear")),
		mcp.WithString("repo", mcp.Description("Only notes from this repository")),
		mcp.WithString("kind", mcp.Description("Only this kind: log, decisions, lessons, plans, learning, corrections or goals")),
		mcp.WithBoolean("all", mcp.Description("Include superseded notes")),
		mcp.WithNumber("limit", mcp.Description("Most results to return, default 20")),
		readOnly,
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		results, err := searchVault(ctx, loc, searchOptions{
			query: req.GetString("query", ""),
			repo:  req.GetString("repo", ""),
			kind:  req.GetString("kind", ""),
			all:   req.GetBool("all", false),
			limit: req.GetInt("limit", defaultLimit),
		})
		return jsonResult(results, err)
	})

	s.AddTool(mcp.NewTool("recall",
		mcp.WithDescription("Read one whole note. Fails if the name matches no note or more than one."),
		note, readOnly,
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		text, err := recallNote(ctx, loc, req.GetString("note", ""))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(text), nil
	})

	s.AddTool(mcp.NewTool("links",
		mcp.WithDescription("List a note's outgoing wikilinks, flagging broken and ambiguous ones."),
		note, readOnly,
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		found, err := linksOf(ctx, loc, req.GetString("note", ""))
		return jsonResult(found, err)
	})

	s.AddTool(mcp.NewTool("backlinks",
		mcp.WithDescription("List the notes that link to a note."),
		note, readOnly,
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sources, err := backlinksOf(ctx, loc, req.GetString("note", ""))
		return jsonResult(sources, err)
	})
	return s
}

func jsonResult(value any, err error) (*mcp.CallToolResult, error) {
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(data)), nil
}
