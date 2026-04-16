package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/abtris/zotero-go-client/zotero"
)

func TestCreateItem(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/users/123/items" {
			wr := zotero.WriteResponse{
				Success: map[string]string{"0": "NEWKEY1"},
			}
			json.NewEncoder(w).Encode(wr)
			return
		}
		http.NotFound(w, r)
	})

	text := callTool(t, session, "create_item", map[string]any{
		"item_type": "book",
		"title":     "New Book",
		"creators":  `[{"creatorType":"author","firstName":"Test","lastName":"Author"}]`,
		"tags":      "tag1, tag2",
	})
	if !strings.Contains(text, "created") {
		t.Errorf("expected 'created' in output, got: %s", text)
	}
	if !strings.Contains(text, "NEWKEY1") {
		t.Errorf("expected 'NEWKEY1' in output, got: %s", text)
	}
}

func TestUpdateItem(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/UPD1" {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			item := zotero.Item{Key: "UPD1", Version: 5, Data: zotero.ItemData{ItemType: "book", Title: "Old Title"}}
			json.NewEncoder(w).Encode(item)
			return
		}
		if r.Method == http.MethodPatch {
			if r.Header.Get("If-Unmodified-Since-Version") != "5" {
				t.Errorf("expected version header 5, got %s", r.Header.Get("If-Unmodified-Since-Version"))
			}
			body, _ := io.ReadAll(r.Body)
			var fields map[string]any
			json.Unmarshal(body, &fields)
			if fields["title"] != "New Title" {
				t.Errorf("expected title 'New Title', got %v", fields["title"])
			}
			item := zotero.Item{Key: "UPD1", Version: 6, Data: zotero.ItemData{ItemType: "book", Title: "New Title"}}
			json.NewEncoder(w).Encode(item)
			return
		}
	})

	text := callTool(t, session, "update_item", map[string]any{
		"item_key": "UPD1",
		"fields":   `{"title":"New Title"}`,
	})
	if !strings.Contains(text, "updated") {
		t.Errorf("expected 'updated' in output, got: %s", text)
	}
}

func TestDeleteItems(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/users/123/items/DEL1" {
			item := zotero.Item{Key: "DEL1", Version: 3}
			json.NewEncoder(w).Encode(item)
			return
		}
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.NotFound(w, r)
	})

	text := callTool(t, session, "delete_items", map[string]any{"item_keys": "DEL1"})
	if !strings.Contains(text, "deleted") {
		t.Errorf("expected 'deleted' in output, got: %s", text)
	}
}

func TestAddTag(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/TAG1" {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			item := zotero.Item{Key: "TAG1", Version: 2, Data: zotero.ItemData{
				ItemType: "book", Title: "Tagged Book",
				Tags: []zotero.Tag{{Tag: "existing"}},
			}}
			json.NewEncoder(w).Encode(item)
			return
		}
		if r.Method == http.MethodPatch {
			item := zotero.Item{Key: "TAG1", Version: 3}
			json.NewEncoder(w).Encode(item)
			return
		}
	})

	text := callTool(t, session, "add_tag", map[string]any{"item_key": "TAG1", "tag": "new-tag"})
	if !strings.Contains(text, "tagged") {
		t.Errorf("expected 'tagged' in output, got: %s", text)
	}
}

func TestAddTagIdempotent(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/TAG2" {
			http.NotFound(w, r)
			return
		}
		item := zotero.Item{Key: "TAG2", Version: 2, Data: zotero.ItemData{
			Tags: []zotero.Tag{{Tag: "already-here"}},
		}}
		json.NewEncoder(w).Encode(item)
	})

	text := callTool(t, session, "add_tag", map[string]any{"item_key": "TAG2", "tag": "already-here"})
	if !strings.Contains(text, "already exists") {
		t.Errorf("expected idempotent message, got: %s", text)
	}
}


func TestRemoveTag(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/RTAG1" {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			item := zotero.Item{Key: "RTAG1", Version: 2, Data: zotero.ItemData{
				Tags: []zotero.Tag{{Tag: "keep"}, {Tag: "remove-me"}},
			}}
			json.NewEncoder(w).Encode(item)
			return
		}
		if r.Method == http.MethodPatch {
			item := zotero.Item{Key: "RTAG1", Version: 3}
			json.NewEncoder(w).Encode(item)
			return
		}
	})

	text := callTool(t, session, "remove_tag", map[string]any{"item_key": "RTAG1", "tag": "remove-me"})
	if !strings.Contains(text, "untagged") {
		t.Errorf("expected 'untagged' in output, got: %s", text)
	}
}

func TestRemoveTagNotFound(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/RTAG2" {
			http.NotFound(w, r)
			return
		}
		item := zotero.Item{Key: "RTAG2", Version: 2, Data: zotero.ItemData{
			Tags: []zotero.Tag{{Tag: "other"}},
		}}
		json.NewEncoder(w).Encode(item)
	})

	text := callTool(t, session, "remove_tag", map[string]any{"item_key": "RTAG2", "tag": "nonexistent"})
	if !strings.Contains(text, "not found") {
		t.Errorf("expected 'not found' message, got: %s", text)
	}
}

func TestAddItemToCollection(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/COLL1" {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			item := zotero.Item{Key: "COLL1", Version: 4, Data: zotero.ItemData{
				Collections: []string{"EXISTING"},
			}}
			json.NewEncoder(w).Encode(item)
			return
		}
		if r.Method == http.MethodPatch {
			item := zotero.Item{Key: "COLL1", Version: 5}
			json.NewEncoder(w).Encode(item)
			return
		}
	})

	text := callTool(t, session, "add_item_to_collection", map[string]any{
		"item_key": "COLL1", "collection_key": "NEWCOL",
	})
	if !strings.Contains(text, "added to collection") {
		t.Errorf("expected 'added to collection', got: %s", text)
	}
}

func TestAddItemToCollectionIdempotent(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/COLL2" {
			http.NotFound(w, r)
			return
		}
		item := zotero.Item{Key: "COLL2", Version: 4, Data: zotero.ItemData{
			Collections: []string{"ALREADY"},
		}}
		json.NewEncoder(w).Encode(item)
	})

	text := callTool(t, session, "add_item_to_collection", map[string]any{
		"item_key": "COLL2", "collection_key": "ALREADY",
	})
	if !strings.Contains(text, "already in collection") {
		t.Errorf("expected idempotent message, got: %s", text)
	}
}

func TestRemoveItemFromCollection(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/RCOL1" {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			item := zotero.Item{Key: "RCOL1", Version: 4, Data: zotero.ItemData{
				Collections: []string{"KEEP", "REMOVE"},
			}}
			json.NewEncoder(w).Encode(item)
			return
		}
		if r.Method == http.MethodPatch {
			item := zotero.Item{Key: "RCOL1", Version: 5}
			json.NewEncoder(w).Encode(item)
			return
		}
	})

	text := callTool(t, session, "remove_item_from_collection", map[string]any{
		"item_key": "RCOL1", "collection_key": "REMOVE",
	})
	if !strings.Contains(text, "removed from collection") {
		t.Errorf("expected 'removed from collection', got: %s", text)
	}
}

func TestRemoveItemFromCollectionNotFound(t *testing.T) {
	session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123/items/RCOL2" {
			http.NotFound(w, r)
			return
		}
		item := zotero.Item{Key: "RCOL2", Version: 4, Data: zotero.ItemData{
			Collections: []string{"OTHER"},
		}}
		json.NewEncoder(w).Encode(item)
	})

	text := callTool(t, session, "remove_item_from_collection", map[string]any{
		"item_key": "RCOL2", "collection_key": "NOTHERE",
	})
	if !strings.Contains(text, "is not in collection") {
		t.Errorf("expected 'not in collection' message, got: %s", text)
	}
}