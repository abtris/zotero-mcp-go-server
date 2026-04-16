package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/abtris/zotero-go-client/zotero"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// formatCreatorNames returns a formatted author string from a list of creators.
func formatCreatorNames(creators []zotero.Creator) string {
	var names []string
	for _, c := range creators {
		if c.Name != "" {
			names = append(names, c.Name)
		} else if c.LastName != "" {
			names = append(names, c.LastName+", "+c.FirstName)
		}
	}
	return strings.Join(names, "; ")
}

// formatCreatorNamesNatural returns "FirstName LastName" style names.
func formatCreatorNamesNatural(creators []zotero.Creator) string {
	var names []string
	for _, c := range creators {
		if c.Name != "" {
			names = append(names, c.Name)
		} else if c.LastName != "" {
			names = append(names, c.FirstName+" "+c.LastName)
		}
	}
	return strings.Join(names, ", ")
}

// formatCreatorLastNames returns only last names joined by comma.
func formatCreatorLastNames(creators []zotero.Creator) string {
	var names []string
	for _, c := range creators {
		if c.Name != "" {
			names = append(names, c.Name)
		} else if c.LastName != "" {
			names = append(names, c.LastName)
		}
	}
	return strings.Join(names, ", ")
}

// formatCitationAPA formats a simplified APA-style citation.
func formatCitationAPA(item *zotero.Item) string {
	authors := formatCreatorNames(item.Data.Creators)
	if authors == "" {
		authors = "Unknown"
	}
	date := item.Data.Date
	if date == "" {
		date = "n.d."
	}
	return fmt.Sprintf("%s (%s). %s.", authors, date, item.Data.Title)
}

// formatCitationChicago formats a simplified Chicago-style citation.
func formatCitationChicago(item *zotero.Item) string {
	authors := formatCreatorNamesNatural(item.Data.Creators)
	if authors == "" {
		authors = "Unknown"
	}
	date := item.Data.Date
	if date == "" {
		date = "n.d."
	}
	return fmt.Sprintf("%s. %s. %s.", authors, item.Data.Title, date)
}

// formatCitationMLA formats a simplified MLA-style citation.
func formatCitationMLA(item *zotero.Item) string {
	authors := formatCreatorNames(item.Data.Creators)
	if authors == "" {
		authors = "Unknown"
	}
	date := item.Data.Date
	if date == "" {
		date = "n.d."
	}
	return fmt.Sprintf("%s. \"%s.\" %s.", authors, item.Data.Title, date)
}

// formatCitationHarvard formats a simplified Harvard-style citation.
func formatCitationHarvard(item *zotero.Item) string {
	authors := formatCreatorLastNames(item.Data.Creators)
	if authors == "" {
		authors = "Unknown"
	}
	date := item.Data.Date
	if date == "" {
		date = "n.d."
	}
	return fmt.Sprintf("%s (%s) %s.", authors, date, item.Data.Title)
}

// formatCitationIEEE formats a simplified IEEE-style citation.
func formatCitationIEEE(item *zotero.Item) string {
	authors := formatCreatorNamesNatural(item.Data.Creators)
	if authors == "" {
		authors = "Unknown"
	}
	date := item.Data.Date
	if date == "" {
		date = "n.d."
	}
	return fmt.Sprintf("%s, \"%s,\" %s.", authors, item.Data.Title, date)
}

// formatCitationBibTeX formats a simplified BibTeX entry.
func formatCitationBibTeX(item *zotero.Item) string {
	authors := formatCreatorNamesNatural(item.Data.Creators)
	itemType := item.Data.ItemType
	if itemType == "journalArticle" {
		itemType = "article"
	}
	return fmt.Sprintf("@%s{%s,\n  author = {%s},\n  title = {%s},\n  year = {%s}\n}",
		itemType, item.Key, authors, item.Data.Title, item.Data.Date)
}

// formatCitation dispatches to the appropriate citation formatter.
func formatCitation(item *zotero.Item, style string) string {
	switch strings.ToLower(style) {
	case "chicago":
		return formatCitationChicago(item)
	case "mla":
		return formatCitationMLA(item)
	case "harvard":
		return formatCitationHarvard(item)
	case "ieee":
		return formatCitationIEEE(item)
	case "bibtex":
		return formatCitationBibTeX(item)
	default:
		return formatCitationAPA(item)
	}
}

// cslStyleID maps short style names to CSL style IDs used by the Zotero API.
var cslStyleID = map[string]string{
	"apa":     "apa",
	"chicago": "chicago-note-bibliography",
	"mla":     "modern-language-association",
	"harvard": "harvard-cite-them-right",
	"ieee":    "ieee",
}

func addGenerateCitationTool(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	type citationArgs struct {
		ItemKey string `json:"item_key" jsonschema:"the Zotero item key"`
		Style   string `json:"style,omitempty" jsonschema:"citation style: apa (default), chicago, mla, harvard, ieee, bibtex, or any CSL style ID from https://www.zotero.org/styles"`
		Locale  string `json:"locale,omitempty" jsonschema:"locale for citation formatting (e.g. en-US, de-DE, fr-FR)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "generate_citation",
		Description: "Generate a formatted citation for a Zotero item using the Zotero API's citeproc-js engine (supports 10,000+ CSL styles)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args citationArgs) (*mcp.CallToolResult, any, error) {
		style := args.Style
		if style == "" {
			style = "apa"
		}

		// BibTeX is not a CSL style — use our Go formatter
		if strings.ToLower(style) == "bibtex" {
			item, _, err := client.Items.Get(ctx, lib, args.ItemKey)
			if err != nil {
				return nil, nil, fmt.Errorf("getting item %s: %w", args.ItemKey, err)
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: formatCitationBibTeX(item)}},
			}, nil, nil
		}

		// Map short names to CSL style IDs
		if cslID, ok := cslStyleID[strings.ToLower(style)]; ok {
			style = cslID
		}

		var opts []zotero.RequestOption
		opts = append(opts, zotero.WithStyle(style))
		if args.Locale != "" {
			opts = append(opts, zotero.WithLocale(args.Locale))
		}

		bib, _, err := client.Items.GetBibliography(ctx, lib, args.ItemKey, opts...)
		if err != nil {
			return nil, nil, fmt.Errorf("getting bibliography for item %s: %w", args.ItemKey, err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: strings.TrimSpace(bib)}},
		}, nil, nil
	})
}
