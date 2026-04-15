package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/abtris/zotero-go-client/zotero"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func addSearchItemsTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type searchArgs struct {
		Query    string `json:"query" jsonschema:"search query string"`
		Limit    int    `json:"limit,omitempty" jsonschema:"max results (1-100, default 25)"`
		ItemType string `json:"item_type,omitempty" jsonschema:"filter by item type (e.g. book, journalArticle)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "search_items",
		Description: "Search for items in the Zotero library by query string",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args searchArgs) (*mcp.CallToolResult, any, error) {
		opts := []zotero.RequestOption{zotero.WithQuery(args.Query)}
		limit := 25
		if args.Limit > 0 && args.Limit <= 100 {
			limit = args.Limit
		}
		opts = append(opts, zotero.WithLimit(limit))
		if args.ItemType != "" {
			opts = append(opts, zotero.WithItemType(args.ItemType))
		}
		items, _, err := client.Items.List(ctx, lib, opts...)
		if err != nil {
			return nil, nil, fmt.Errorf("searching items: %w", err)
		}
		return formatItemsResult(items), nil, nil
	})
}

func addGetItemTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type getItemArgs struct {
		ItemKey string `json:"item_key" jsonschema:"the Zotero item key"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_item",
		Description: "Get details of a specific Zotero item by its key",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getItemArgs) (*mcp.CallToolResult, any, error) {
		item, _, err := client.Items.Get(ctx, lib, args.ItemKey)
		if err != nil {
			return nil, nil, fmt.Errorf("getting item %s: %w", args.ItemKey, err)
		}
		data, _ := json.MarshalIndent(item, "", "  ")
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}

func addListCollectionsTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type listCollArgs struct {
		Limit int `json:"limit,omitempty" jsonschema:"max results (1-100, default 25)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list_collections",
		Description: "List collections (folders) in the Zotero library",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listCollArgs) (*mcp.CallToolResult, any, error) {
		limit := 25
		if args.Limit > 0 && args.Limit <= 100 {
			limit = args.Limit
		}
		colls, _, err := client.Collections.List(ctx, lib, zotero.WithLimit(limit))
		if err != nil {
			return nil, nil, fmt.Errorf("listing collections: %w", err)
		}
		return formatCollectionsResult(colls), nil, nil
	})
}

func addListCollectionItemsTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type collItemsArgs struct {
		CollectionKey string `json:"collection_key" jsonschema:"the collection key"`
		Limit         int    `json:"limit,omitempty" jsonschema:"max results (1-100, default 25)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list_collection_items",
		Description: "List items in a specific Zotero collection",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args collItemsArgs) (*mcp.CallToolResult, any, error) {
		limit := 25
		if args.Limit > 0 && args.Limit <= 100 {
			limit = args.Limit
		}
		items, _, err := client.Items.ListInCollection(ctx, lib, args.CollectionKey, zotero.WithLimit(limit))
		if err != nil {
			return nil, nil, fmt.Errorf("listing collection items: %w", err)
		}
		return formatItemsResult(items), nil, nil
	})
}

func addListTagsTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type listTagsArgs struct {
		Limit int `json:"limit,omitempty" jsonschema:"max results (1-100, default 50)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list_tags",
		Description: "List tags in the Zotero library",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listTagsArgs) (*mcp.CallToolResult, any, error) {
		limit := 50
		if args.Limit > 0 && args.Limit <= 100 {
			limit = args.Limit
		}
		tags, _, err := client.Tags.List(ctx, lib, zotero.WithLimit(limit))
		if err != nil {
			return nil, nil, fmt.Errorf("listing tags: %w", err)
		}
		return formatTagsResult(tags), nil, nil
	})
}

func addListItemsByTagTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type tagItemsArgs struct {
		Tag   string `json:"tag" jsonschema:"the tag to filter by"`
		Limit int    `json:"limit,omitempty" jsonschema:"max results (1-100, default 25)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list_items_by_tag",
		Description: "List items that have a specific tag",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args tagItemsArgs) (*mcp.CallToolResult, any, error) {
		limit := 25
		if args.Limit > 0 && args.Limit <= 100 {
			limit = args.Limit
		}
		items, _, err := client.Items.List(ctx, lib, zotero.WithTag(args.Tag), zotero.WithLimit(limit))
		if err != nil {
			return nil, nil, fmt.Errorf("listing items by tag: %w", err)
		}
		return formatItemsResult(items), nil, nil
	})
}
