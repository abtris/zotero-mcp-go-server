package server

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/abtris/zotero-go-client/zotero"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type itemSummary struct {
	Key      string           `json:"key"`
	Type     string           `json:"itemType"`
	Title    string           `json:"title"`
	Creators []zotero.Creator `json:"creators,omitempty"`
	Date     string           `json:"date,omitempty"`
	URL      string           `json:"url,omitempty"`
	Tags     []string         `json:"tags,omitempty"`
}

type collectionSummary struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Parent any    `json:"parent,omitempty"`
}

type tagSummary struct {
	Tag      string `json:"tag"`
	NumItems int    `json:"numItems,omitempty"`
}

// formatItemsResult formats a list of items as an MCP tool result.
func formatItemsResult(items []*zotero.Item) *mcp.CallToolResult {
	summaries := make([]itemSummary, len(items))
	for i, item := range items {
		var tags []string
		for _, t := range item.Data.Tags {
			tags = append(tags, t.Tag)
		}
		summaries[i] = itemSummary{
			Key:      item.Key,
			Type:     item.Data.ItemType,
			Title:    item.Data.Title,
			Creators: item.Data.Creators,
			Date:     item.Data.Date,
			URL:      item.Data.URL,
			Tags:     tags,
		}
	}
	header := fmt.Sprintf("Found %s items\n\n", strconv.Itoa(len(summaries)))
	data, _ := json.MarshalIndent(summaries, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: header + string(data)}},
	}
}

// formatCollectionsResult formats a list of collections as an MCP tool result.
func formatCollectionsResult(colls []*zotero.Collection) *mcp.CallToolResult {
	summaries := make([]collectionSummary, len(colls))
	for i, c := range colls {
		summaries[i] = collectionSummary{
			Key:    c.Key,
			Name:   c.Data.Name,
			Parent: c.Data.ParentCollection,
		}
	}
	data, _ := json.MarshalIndent(summaries, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}
}

// formatTagsResult formats a list of tags as an MCP tool result.
func formatTagsResult(tags []*zotero.TagEntry) *mcp.CallToolResult {
	summaries := make([]tagSummary, len(tags))
	for i, t := range tags {
		s := tagSummary{Tag: t.Tag}
		if t.Meta != nil {
			s.NumItems = t.Meta.NumItems
		}
		summaries[i] = s
	}
	data, _ := json.MarshalIndent(summaries, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}
}
