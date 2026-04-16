package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/abtris/zotero-go-client/zotero"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func addCreateItemTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type createArgs struct {
		ItemType     string `json:"item_type" jsonschema:"item type (e.g. book, journalArticle, conferencePaper)"`
		Title        string `json:"title" jsonschema:"item title"`
		Creators     string `json:"creators,omitempty" jsonschema:"JSON array of creators, e.g. [{\"creatorType\":\"author\",\"firstName\":\"Jane\",\"lastName\":\"Doe\"}]"`
		Date         string `json:"date,omitempty" jsonschema:"publication date"`
		AbstractNote string `json:"abstract_note,omitempty" jsonschema:"abstract or note"`
		URL          string `json:"url,omitempty" jsonschema:"URL"`
		Tags         string `json:"tags,omitempty" jsonschema:"comma-separated tags"`
		Collections  string `json:"collections,omitempty" jsonschema:"comma-separated collection keys"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "create_item",
		Description: "Create a new item in the Zotero library",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createArgs) (*mcp.CallToolResult, any, error) {
		item := &zotero.ItemData{
			ItemType:     args.ItemType,
			Title:        args.Title,
			Date:         args.Date,
			AbstractNote: args.AbstractNote,
			URL:          args.URL,
		}
		if args.Creators != "" {
			var creators []zotero.Creator
			if err := json.Unmarshal([]byte(args.Creators), &creators); err != nil {
				return nil, nil, fmt.Errorf("parsing creators JSON: %w", err)
			}
			item.Creators = creators
		}
		if args.Tags != "" {
			for _, t := range strings.Split(args.Tags, ",") {
				t = strings.TrimSpace(t)
				if t != "" {
					item.Tags = append(item.Tags, zotero.Tag{Tag: t})
				}
			}
		}
		if args.Collections != "" {
			for _, c := range strings.Split(args.Collections, ",") {
				c = strings.TrimSpace(c)
				if c != "" {
					item.Collections = append(item.Collections, c)
				}
			}
		}
		wr, _, err := client.Items.Create(ctx, lib, []*zotero.ItemData{item})
		if err != nil {
			return nil, nil, fmt.Errorf("creating item: %w", err)
		}
		details := map[string]any{"success": wr.Success, "failed": wr.Failed}
		key := ""
		for _, v := range wr.Success {
			key = v
			break
		}
		return formatWriteResponse("created", key, details), nil, nil
	})
}

func addUpdateItemTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type updateArgs struct {
		ItemKey string `json:"item_key" jsonschema:"the Zotero item key to update"`
		Fields  string `json:"fields" jsonschema:"JSON object of fields to update, e.g. {\"title\":\"New Title\",\"date\":\"2025\"}"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "update_item",
		Description: "Partially update an existing Zotero item (patch)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateArgs) (*mcp.CallToolResult, any, error) {
		item, _, err := client.Items.Get(ctx, lib, args.ItemKey)
		if err != nil {
			return nil, nil, fmt.Errorf("getting item %s: %w", args.ItemKey, err)
		}
		var fields map[string]any
		if err := json.Unmarshal([]byte(args.Fields), &fields); err != nil {
			return nil, nil, fmt.Errorf("parsing fields JSON: %w", err)
		}
		_, _, err = client.Items.Patch(ctx, lib, args.ItemKey, fields, item.Version)
		if err != nil {
			return nil, nil, fmt.Errorf("updating item %s: %w", args.ItemKey, err)
		}
		return formatWriteResponse("updated", args.ItemKey, fields), nil, nil
	})
}

func addDeleteItemsTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type deleteArgs struct {
		ItemKeys string `json:"item_keys" jsonschema:"comma-separated item keys to delete"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "delete_items",
		Description: "Delete one or more items from the Zotero library",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deleteArgs) (*mcp.CallToolResult, any, error) {
		keys := strings.Split(args.ItemKeys, ",")
		for i := range keys {
			keys[i] = strings.TrimSpace(keys[i])
		}
		// Get version from first item
		item, _, err := client.Items.Get(ctx, lib, keys[0])
		if err != nil {
			return nil, nil, fmt.Errorf("getting item %s for version: %w", keys[0], err)
		}
		if len(keys) == 1 {
			_, err = client.Items.Delete(ctx, lib, keys[0], item.Version)
		} else {
			_, err = client.Items.DeleteMultiple(ctx, lib, keys, item.Version)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("deleting items: %w", err)
		}
		return formatWriteResponse("deleted", strings.Join(keys, ", "), nil), nil, nil
	})
}


func addAddTagTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type addTagArgs struct {
		ItemKey string `json:"item_key" jsonschema:"the Zotero item key"`
		Tag     string `json:"tag" jsonschema:"tag to add"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "add_tag",
		Description: "Add a tag to a Zotero item (idempotent)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addTagArgs) (*mcp.CallToolResult, any, error) {
		item, _, err := client.Items.Get(ctx, lib, args.ItemKey)
		if err != nil {
			return nil, nil, fmt.Errorf("getting item %s: %w", args.ItemKey, err)
		}
		for _, t := range item.Data.Tags {
			if t.Tag == args.Tag {
				return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{
						Text: fmt.Sprintf("Tag %q already exists on item %s", args.Tag, args.ItemKey),
					}},
				}, nil, nil
			}
		}
		tags := append(item.Data.Tags, zotero.Tag{Tag: args.Tag})
		fields := map[string]any{"tags": tags}
		_, _, err = client.Items.Patch(ctx, lib, args.ItemKey, fields, item.Version)
		if err != nil {
			return nil, nil, fmt.Errorf("adding tag to item %s: %w", args.ItemKey, err)
		}
		return formatWriteResponse("tagged", args.ItemKey, map[string]any{"tag": args.Tag}), nil, nil
	})
}

func addRemoveTagTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type removeTagArgs struct {
		ItemKey string `json:"item_key" jsonschema:"the Zotero item key"`
		Tag     string `json:"tag" jsonschema:"tag to remove"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "remove_tag",
		Description: "Remove a tag from a Zotero item",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args removeTagArgs) (*mcp.CallToolResult, any, error) {
		item, _, err := client.Items.Get(ctx, lib, args.ItemKey)
		if err != nil {
			return nil, nil, fmt.Errorf("getting item %s: %w", args.ItemKey, err)
		}
		var filtered []zotero.Tag
		found := false
		for _, t := range item.Data.Tags {
			if t.Tag == args.Tag {
				found = true
			} else {
				filtered = append(filtered, t)
			}
		}
		if !found {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{
					Text: fmt.Sprintf("Tag %q not found on item %s", args.Tag, args.ItemKey),
				}},
			}, nil, nil
		}
		fields := map[string]any{"tags": filtered}
		_, _, err = client.Items.Patch(ctx, lib, args.ItemKey, fields, item.Version)
		if err != nil {
			return nil, nil, fmt.Errorf("removing tag from item %s: %w", args.ItemKey, err)
		}
		return formatWriteResponse("untagged", args.ItemKey, map[string]any{"removed_tag": args.Tag}), nil, nil
	})
}

func addAddItemToCollectionTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type addCollArgs struct {
		ItemKey       string `json:"item_key" jsonschema:"the Zotero item key"`
		CollectionKey string `json:"collection_key" jsonschema:"the collection key to add the item to"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "add_item_to_collection",
		Description: "Add an item to a Zotero collection (idempotent)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args addCollArgs) (*mcp.CallToolResult, any, error) {
		item, _, err := client.Items.Get(ctx, lib, args.ItemKey)
		if err != nil {
			return nil, nil, fmt.Errorf("getting item %s: %w", args.ItemKey, err)
		}
		for _, c := range item.Data.Collections {
			if c == args.CollectionKey {
				return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{
						Text: fmt.Sprintf("Item %s is already in collection %s", args.ItemKey, args.CollectionKey),
					}},
				}, nil, nil
			}
		}
		collections := append(item.Data.Collections, args.CollectionKey)
		fields := map[string]any{"collections": collections}
		_, _, err = client.Items.Patch(ctx, lib, args.ItemKey, fields, item.Version)
		if err != nil {
			return nil, nil, fmt.Errorf("adding item %s to collection: %w", args.ItemKey, err)
		}
		return formatWriteResponse("added to collection", args.ItemKey,
			map[string]any{"collection": args.CollectionKey}), nil, nil
	})
}

func addRemoveItemFromCollectionTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type removeCollArgs struct {
		ItemKey       string `json:"item_key" jsonschema:"the Zotero item key"`
		CollectionKey string `json:"collection_key" jsonschema:"the collection key to remove the item from"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "remove_item_from_collection",
		Description: "Remove an item from a Zotero collection",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args removeCollArgs) (*mcp.CallToolResult, any, error) {
		item, _, err := client.Items.Get(ctx, lib, args.ItemKey)
		if err != nil {
			return nil, nil, fmt.Errorf("getting item %s: %w", args.ItemKey, err)
		}
		var filtered []string
		found := false
		for _, c := range item.Data.Collections {
			if c == args.CollectionKey {
				found = true
			} else {
				filtered = append(filtered, c)
			}
		}
		if !found {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{
					Text: fmt.Sprintf("Item %s is not in collection %s", args.ItemKey, args.CollectionKey),
				}},
			}, nil, nil
		}
		fields := map[string]any{"collections": filtered}
		_, _, err = client.Items.Patch(ctx, lib, args.ItemKey, fields, item.Version)
		if err != nil {
			return nil, nil, fmt.Errorf("removing item %s from collection: %w", args.ItemKey, err)
		}
		return formatWriteResponse("removed from collection", args.ItemKey,
			map[string]any{"removed_collection": args.CollectionKey}), nil, nil
	})
}