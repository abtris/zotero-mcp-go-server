// Package server provides the Zotero MCP server setup and tool registration.
package server

import (
	"github.com/abtris/zotero-go-client/zotero"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "zotero-mcp-server"
	serverVersion = "0.1.0"
)

// New creates and configures a new MCP server with all Zotero tools registered.
func New(client *zotero.Client, lib zotero.LibraryID) *mcp.Server {
	srv := mcp.NewServer(
		&mcp.Implementation{Name: serverName, Version: serverVersion},
		nil,
	)

	registerTools(srv, client, lib)

	return srv
}

// registerTools adds all Zotero tools to the MCP server.
func registerTools(srv *mcp.Server, client *zotero.Client, lib zotero.LibraryID) {
	// Read tools
	addSearchItemsTool(srv, client, lib)
	addGetItemTool(srv, client, lib)
	addListCollectionsTool(srv, client, lib)
	addListCollectionItemsTool(srv, client, lib)
	addListTagsTool(srv, client, lib)
	addListItemsByTagTool(srv, client, lib)
	addGenerateCitationTool(srv, client, lib)

	// Write tools
	addCreateItemTool(srv, client, lib)
	addUpdateItemTool(srv, client, lib)
	addDeleteItemsTool(srv, client, lib)
	addAddTagTool(srv, client, lib)
	addRemoveTagTool(srv, client, lib)
	addAddItemToCollectionTool(srv, client, lib)
	addRemoveItemFromCollectionTool(srv, client, lib)
}
