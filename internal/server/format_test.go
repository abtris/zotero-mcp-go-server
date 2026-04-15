package server

import (
	"strings"
	"testing"

	"github.com/abtris/zotero-go-client/zotero"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestFormatItemsResultWithTags(t *testing.T) {
	items := []*zotero.Item{
		{
			Key: "T1",
			Data: zotero.ItemData{
				ItemType: "book",
				Title:    "Tagged Book",
				Tags:     []zotero.Tag{{Tag: "tag1"}, {Tag: "tag2"}},
			},
		},
	}
	result := formatItemsResult(items)
	if len(result.Content) == 0 {
		t.Fatal("expected content")
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}
	if !strings.Contains(tc.Text, "Found 1 items") {
		t.Errorf("expected 'Found 1 items' in output, got: %s", tc.Text)
	}
	if !strings.Contains(tc.Text, "tag1") {
		t.Errorf("expected 'tag1' in output, got: %s", tc.Text)
	}
	if !strings.Contains(tc.Text, "tag2") {
		t.Errorf("expected 'tag2' in output, got: %s", tc.Text)
	}
}

func TestFormatItemsResultEmpty(t *testing.T) {
	result := formatItemsResult(nil)
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}
	if !strings.Contains(tc.Text, "Found 0 items") {
		t.Errorf("expected 'Found 0 items', got: %s", tc.Text)
	}
}

func TestFormatCollectionsResult(t *testing.T) {
	colls := []*zotero.Collection{
		{Key: "C1", Data: zotero.CollectionData{Name: "Papers"}},
		{Key: "C2", Data: zotero.CollectionData{Name: "Books", ParentCollection: "C1"}},
	}
	result := formatCollectionsResult(colls)
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}
	if !strings.Contains(tc.Text, "Papers") {
		t.Errorf("expected 'Papers' in output, got: %s", tc.Text)
	}
	if !strings.Contains(tc.Text, "C1") {
		t.Errorf("expected 'C1' in output, got: %s", tc.Text)
	}
}

func TestFormatTagsResult(t *testing.T) {
	tags := []*zotero.TagEntry{
		{Tag: "golang", Meta: &zotero.TagMeta{NumItems: 3}},
		{Tag: "rust"},
	}
	result := formatTagsResult(tags)
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}
	if !strings.Contains(tc.Text, "golang") {
		t.Errorf("expected 'golang' in output, got: %s", tc.Text)
	}
	if !strings.Contains(tc.Text, `"numItems": 3`) {
		t.Errorf("expected numItems 3 in output, got: %s", tc.Text)
	}
	if !strings.Contains(tc.Text, "rust") {
		t.Errorf("expected 'rust' in output, got: %s", tc.Text)
	}
}
