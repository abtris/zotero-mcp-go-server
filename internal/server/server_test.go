package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abtris/zotero-go-client/zotero"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// setupTestServer creates a mock Zotero HTTP server, an MCP server with tools
// registered, and an MCP client session connected via in-memory transport.
func setupTestServer(t *testing.T, handler http.HandlerFunc) *mcp.ClientSession {
	t.Helper()

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	client := zotero.NewClient("test-key", zotero.WithBaseURL(ts.URL))
	lib := zotero.UserLibrary("123")

	srv := New(client, lib)

	ctx := context.Background()
	t1, t2 := mcp.NewInMemoryTransports()

	go func() {
		srv.Connect(ctx, t1, nil)
	}()

	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.1"}, nil)
	session, err := mcpClient.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("connecting MCP client: %v", err)
	}
	t.Cleanup(func() { session.Close() })

	return session
}

// callTool is a helper that calls a tool and returns the text content.
func callTool(t *testing.T, session *mcp.ClientSession, name string, args map[string]any) string {
	t.Helper()
	ctx := context.Background()
	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		t.Fatalf("CallTool(%q): %v", name, err)
	}
	if len(res.Content) == 0 {
		t.Fatalf("CallTool(%q): empty content", name)
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("CallTool(%q): expected TextContent, got %T", name, res.Content[0])
	}
	return tc.Text
}

func TestSearchItems(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items" {
			http.NotFound(w, r)
			return
		}
		q := r.URL.Query()
		if q.Get("q") != "machine learning" {
			t.Errorf("query = %q, want %q", q.Get("q"), "machine learning")
		}
		items := []zotero.Item{
			{Key: "ML1", Data: zotero.ItemData{ItemType: "book", Title: "ML Basics"}},
			{Key: "ML2", Data: zotero.ItemData{ItemType: "journalArticle", Title: "Deep Learning Review"}},
		}
		json.NewEncoder(w).Encode(items)
	})

	text := callTool(t, session, "search_items", map[string]any{"query": "machine learning"})
	if !strings.Contains(text, "Found 2 items") {
		t.Errorf("expected 'Found 2 items' in output, got: %s", text)
	}
	if !strings.Contains(text, "ML Basics") {
		t.Errorf("expected 'ML Basics' in output, got: %s", text)
	}
	if !strings.Contains(text, "Deep Learning Review") {
		t.Errorf("expected 'Deep Learning Review' in output, got: %s", text)
	}
}

func TestSearchItemsWithItemType(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("itemType") != "book" {
			t.Errorf("itemType = %q, want %q", q.Get("itemType"), "book")
		}
		items := []zotero.Item{
			{Key: "B1", Data: zotero.ItemData{ItemType: "book", Title: "A Book"}},
		}
		json.NewEncoder(w).Encode(items)
	})

	text := callTool(t, session, "search_items", map[string]any{
		"query":     "test",
		"item_type": "book",
	})
	if !strings.Contains(text, "Found 1 items") {
		t.Errorf("expected 'Found 1 items', got: %s", text)
	}
}

func TestSearchItemsEmpty(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[]"))
	})

	text := callTool(t, session, "search_items", map[string]any{"query": "nonexistent"})
	if !strings.Contains(text, "Found 0 items") {
		t.Errorf("expected 'Found 0 items', got: %s", text)
	}
}

func TestSearchItemsCustomLimit(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("limit") != "5" {
			t.Errorf("limit = %q, want %q", q.Get("limit"), "5")
		}
		w.Write([]byte("[]"))
	})

	callTool(t, session, "search_items", map[string]any{"query": "test", "limit": 5})
}
