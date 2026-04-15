package main

import (
	"context"
	"log"
	"os"

	"github.com/abtris/zotero-go-client/zotero"
	"github.com/abtris/zotero-mcp-go-server/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	apiKey := os.Getenv("ZOTERO_API_KEY")
	userID := os.Getenv("ZOTERO_USER_ID")
	groupID := os.Getenv("ZOTERO_GROUP_ID")

	if apiKey == "" {
		log.Fatal("ZOTERO_API_KEY environment variable is required")
	}
	if userID == "" && groupID == "" {
		log.Fatal("ZOTERO_USER_ID or ZOTERO_GROUP_ID environment variable is required")
	}

	client := zotero.NewClient(apiKey)

	var lib zotero.LibraryID
	if groupID != "" {
		lib = zotero.GroupLibrary(groupID)
	} else {
		lib = zotero.UserLibrary(userID)
	}

	srv := server.New(client, lib)

	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}