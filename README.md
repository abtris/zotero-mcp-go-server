# zotero-mcp-go-server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server for [Zotero](https://www.zotero.org/), written in Go. It provides read-only access to your Zotero library through MCP tools, allowing LLM clients (Claude Desktop, Cursor, etc.) to search, browse, and retrieve your references.

## Prerequisites

- Go 1.24 or later
- A Zotero account with a Web API key

## Getting Credentials

### 1. Zotero API Key

1. Go to [https://www.zotero.org/settings/keys](https://www.zotero.org/settings/keys)
2. Click **Create new private key**
3. Give it a name (e.g. "MCP Server")
4. Under **Personal Library**, check **Allow library access**
5. If you need group access, check the relevant groups under **Group Permissions**
6. Click **Save Key**
7. Copy the generated key — this is your `ZOTERO_API_KEY`

### 2. User ID

Your user ID is displayed on the same [API keys page](https://www.zotero.org/settings/keys) at the top:

> "Your userID for use in API calls is **1234567**"

This is your `ZOTERO_USER_ID`.

### 3. Group ID (optional)

If you want to access a group library instead of your personal library:

1. Go to [https://www.zotero.org/groups](https://www.zotero.org/groups)
2. Open the group you want to access
3. The group ID is the number in the URL: `https://www.zotero.org/groups/1234567/...`

This is your `ZOTERO_GROUP_ID`.

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `ZOTERO_API_KEY` | Yes | Your Zotero Web API key |
| `ZOTERO_USER_ID` | Yes* | Your Zotero user ID |
| `ZOTERO_GROUP_ID` | Yes* | A Zotero group ID |

\* One of `ZOTERO_USER_ID` or `ZOTERO_GROUP_ID` is required. If both are set, the group library is used.

## Installation

### Homebrew (macOS)

```bash
brew install abtris/tap/zotero-mcp-go-server
```

### Docker

```bash
docker pull abtris/zotero-mcp-go-server:latest
```

Run with Docker:

```bash
docker run --rm \
  -e ZOTERO_API_KEY="your-api-key" \
  -e ZOTERO_USER_ID="your-user-id" \
  abtris/zotero-mcp-go-server:latest
```

### Go install

```bash
go install github.com/abtris/zotero-mcp-go-server@latest
```

### Build from source

```bash
git clone https://github.com/abtris/zotero-mcp-go-server.git
cd zotero-mcp-go-server
go build -o zotero-mcp-go-server .
```

## Usage

### Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "zotero": {
      "command": "zotero-mcp-go-server",
      "env": {
        "ZOTERO_API_KEY": "your-api-key",
        "ZOTERO_USER_ID": "your-user-id"
      }
    }
  }
}
```

### Run directly (for development)

```bash
export ZOTERO_API_KEY="your-api-key"
export ZOTERO_USER_ID="your-user-id"
go run .
```

The server communicates over stdio using the MCP protocol.

## Supported Tools

### `search_items`

Search for items in the Zotero library by query string.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `query` | string | Yes | Search query string |
| `limit` | integer | No | Max results, 1–100 (default: 25) |
| `item_type` | string | No | Filter by item type (e.g. `book`, `journalArticle`, `conferencePaper`) |

### `get_item`

Get full details of a specific Zotero item by its key.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `item_key` | string | Yes | The Zotero item key (e.g. `ABCD1234`) |

### `list_collections`

List collections (folders) in the Zotero library.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `limit` | integer | No | Max results, 1–100 (default: 25) |

### `list_collection_items`

List items in a specific Zotero collection.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `collection_key` | string | Yes | The collection key |
| `limit` | integer | No | Max results, 1–100 (default: 25) |

### `list_tags`

List tags in the Zotero library.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `limit` | integer | No | Max results, 1–100 (default: 50) |

### `list_items_by_tag`

List items that have a specific tag.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `tag` | string | Yes | The tag to filter by |
| `limit` | integer | No | Max results, 1–100 (default: 25) |

## Testing and Debugging

### Unit Tests

```bash
go test -v ./...
```

### MCP Inspector

The [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) is an interactive tool for testing and debugging MCP servers. You can use it to connect to this server, list available tools, and call them interactively.

```bash
npx @modelcontextprotocol/inspector go run .
```

This opens a web UI where you can:

- See all registered tools and their schemas
- Call tools with custom arguments and inspect responses
- View the raw JSON-RPC messages exchanged between client and server

Make sure the required environment variables are set before running the inspector:

```bash
export ZOTERO_API_KEY="your-api-key"
export ZOTERO_USER_ID="your-user-id"
npx @modelcontextprotocol/inspector go run .
```

## License

MIT
