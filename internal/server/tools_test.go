package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/abtris/zotero-go-client/zotero"
)

func TestGetItem(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/ITEM1" {
			http.NotFound(w, r)
			return
		}
		item := zotero.Item{
			Key:     "ITEM1",
			Version: 5,
			Data: zotero.ItemData{
				ItemType: "book",
				Title:    "Test Book",
				Creators: []zotero.Creator{
					{CreatorType: "author", FirstName: "Jane", LastName: "Doe"},
				},
				Date: "2024",
			},
		}
		json.NewEncoder(w).Encode(item)
	})

	text := callTool(t, session, "get_item", map[string]any{"item_key": "ITEM1"})
	if !strings.Contains(text, "Test Book") {
		t.Errorf("expected 'Test Book' in output, got: %s", text)
	}
	if !strings.Contains(text, "Jane") {
		t.Errorf("expected creator 'Jane' in output, got: %s", text)
	}
	if !strings.Contains(text, "ITEM1") {
		t.Errorf("expected key 'ITEM1' in output, got: %s", text)
	}
}

func TestListCollections(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/collections" {
			http.NotFound(w, r)
			return
		}
		colls := []zotero.Collection{
			{Key: "COL1", Data: zotero.CollectionData{Name: "Research"}},
			{Key: "COL2", Data: zotero.CollectionData{Name: "Reading List", ParentCollection: "COL1"}},
		}
		json.NewEncoder(w).Encode(colls)
	})

	text := callTool(t, session, "list_collections", map[string]any{})
	if !strings.Contains(text, "Research") {
		t.Errorf("expected 'Research' in output, got: %s", text)
	}
	if !strings.Contains(text, "Reading List") {
		t.Errorf("expected 'Reading List' in output, got: %s", text)
	}
	if !strings.Contains(text, "COL1") {
		t.Errorf("expected 'COL1' in output, got: %s", text)
	}
}

func TestListCollectionItems(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/collections/COL1/items" {
			http.NotFound(w, r)
			return
		}
		items := []zotero.Item{
			{Key: "I1", Data: zotero.ItemData{ItemType: "book", Title: "Book in Collection"}},
		}
		json.NewEncoder(w).Encode(items)
	})

	text := callTool(t, session, "list_collection_items", map[string]any{
		"collection_key": "COL1",
	})
	if !strings.Contains(text, "Found 1 items") {
		t.Errorf("expected 'Found 1 items', got: %s", text)
	}
	if !strings.Contains(text, "Book in Collection") {
		t.Errorf("expected 'Book in Collection' in output, got: %s", text)
	}
}

func TestListTags(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/tags" {
			http.NotFound(w, r)
			return
		}
		tags := []zotero.TagEntry{
			{Tag: "machine-learning", Meta: &zotero.TagMeta{NumItems: 10}},
			{Tag: "ai", Meta: &zotero.TagMeta{NumItems: 5}},
			{Tag: "research"},
		}
		json.NewEncoder(w).Encode(tags)
	})

	text := callTool(t, session, "list_tags", map[string]any{})
	if !strings.Contains(text, "machine-learning") {
		t.Errorf("expected 'machine-learning' in output, got: %s", text)
	}
	if !strings.Contains(text, "ai") {
		t.Errorf("expected 'ai' in output, got: %s", text)
	}
	if !strings.Contains(text, `"numItems": 10`) {
		t.Errorf("expected numItems 10 in output, got: %s", text)
	}
}

func TestListItemsByTag(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items" {
			http.NotFound(w, r)
			return
		}
		q := r.URL.Query()
		if q.Get("tag") != "ai" {
			t.Errorf("tag = %q, want %q", q.Get("tag"), "ai")
		}
		items := []zotero.Item{
			{Key: "AI1", Data: zotero.ItemData{
				ItemType: "journalArticle",
				Title:    "AI Paper",
				Tags:     []zotero.Tag{{Tag: "ai"}, {Tag: "research"}},
			}},
		}
		json.NewEncoder(w).Encode(items)
	})

	text := callTool(t, session, "list_items_by_tag", map[string]any{"tag": "ai"})
	if !strings.Contains(text, "Found 1 items") {
		t.Errorf("expected 'Found 1 items', got: %s", text)
	}
	if !strings.Contains(text, "AI Paper") {
		t.Errorf("expected 'AI Paper' in output, got: %s", text)
	}
}
