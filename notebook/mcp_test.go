package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// TestMain lets the test binary act as the notebook command, so the server
// test can start it as a real child process.
func TestMain(m *testing.M) {
	if os.Getenv("NOTEBOOK_RUN_MAIN") == "1" {
		main()
		return
	}
	os.Exit(m.Run())
}

func startServer(t *testing.T, e env) (*client.Client, context.Context) {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	c, err := client.NewStdioMCPClient(exe, []string{"NOTEBOOK_RUN_MAIN=1"}, "mcp", "--vault", e.vault, "--index", e.index)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	init := mcp.InitializeRequest{}
	init.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	init.Params.ClientInfo = mcp.Implementation{Name: "test", Version: "0"}
	if _, err := c.Initialize(ctx, init); err != nil {
		t.Fatal(err)
	}
	return c, ctx
}

func call(t *testing.T, c *client.Client, ctx context.Context, tool string, args map[string]any) (string, bool) {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Name = tool
	req.Params.Arguments = args
	res, err := c.CallTool(ctx, req)
	if err != nil {
		t.Fatalf("%s: %v", tool, err)
	}
	var text strings.Builder
	for _, content := range res.Content {
		if tc, ok := content.(mcp.TextContent); ok {
			text.WriteString(tc.Text)
		}
	}
	return text.String(), res.IsError
}

func TestServerListsExactlyTheReadTools(t *testing.T) {
	c, ctx := startServer(t, newEnv(t))
	res, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, tool := range res.Tools {
		names = append(names, tool.Name)
		if tool.Annotations.ReadOnlyHint == nil || !*tool.Annotations.ReadOnlyHint {
			t.Errorf("%s is not marked read-only", tool.Name)
		}
	}
	slices.Sort(names)
	if want := []string{"backlinks", "links", "recall", "search"}; !reflect.DeepEqual(names, want) {
		t.Errorf("tools %v, want %v", names, want)
	}
}

func TestServerSearch(t *testing.T) {
	c, ctx := startServer(t, newEnv(t))
	text, isError := call(t, c, ctx, "search", map[string]any{"query": "force push"})
	var got []result
	if err := json.Unmarshal([]byte(text), &got); isError || err != nil {
		t.Fatalf("isError %v, err %v, text %q", isError, err, text)
	}
	if want := []string{forcePushNote, cleanupNote, backupNote}; !reflect.DeepEqual(paths(got), want) {
		t.Errorf("got %v", paths(got))
	}
	text, _ = call(t, c, ctx, "search", map[string]any{"query": "force push", "repo": "homelab", "limit": 1})
	if err := json.Unmarshal([]byte(text), &got); err != nil || !reflect.DeepEqual(paths(got), []string{backupNote}) {
		t.Errorf("filtered: %v, %v", got, err)
	}
	text, _ = call(t, c, ctx, "search", map[string]any{"query": "rebase", "all": true})
	if err := json.Unmarshal([]byte(text), &got); err != nil || len(got) != 2 {
		t.Errorf("all: %v, %v", got, err)
	}
}

func TestServerSeesEditsWhileRunning(t *testing.T) {
	e := newEnv(t)
	c, ctx := startServer(t, e)
	call(t, c, ctx, "search", map[string]any{"query": "force"})
	write(t, e.vault+"/lessons/new.md", "# Aardvark\n", time.Now())
	if text, _ := call(t, c, ctx, "search", map[string]any{"query": "aardvark"}); !strings.Contains(text, "lessons/new.md") {
		t.Errorf("text %q", text)
	}
}

func TestServerRecallLinksBacklinks(t *testing.T) {
	c, ctx := startServer(t, newEnv(t))
	whole, err := os.ReadFile("testdata/vault/" + forcePushNote)
	if err != nil {
		t.Fatal(err)
	}
	if text, isError := call(t, c, ctx, "recall", map[string]any{"note": forcePushNote}); isError || text != string(whole) {
		t.Errorf("recall: isError %v, %q", isError, text)
	}
	if text, isError := call(t, c, ctx, "links", map[string]any{"note": "Link source"}); isError || !strings.Contains(text, `"broken"`) {
		t.Errorf("links: isError %v, %q", isError, text)
	}
	if text, isError := call(t, c, ctx, "backlinks", map[string]any{"note": "keep-commits-small"}); isError || !strings.Contains(text, "lessons/link-source.md") {
		t.Errorf("backlinks: isError %v, %q", isError, text)
	}
}

func TestServerReportsFailuresAsToolErrors(t *testing.T) {
	c, ctx := startServer(t, newEnv(t))
	for tool, args := range map[string]map[string]any{
		"recall":    {"note": "Duplicate title"},
		"links":     {"note": "no-such-note"},
		"backlinks": {},
		"search":    {"query": "---"},
	} {
		if text, isError := call(t, c, ctx, tool, args); !isError {
			t.Errorf("%s: no error: %q", tool, text)
		}
	}
}
