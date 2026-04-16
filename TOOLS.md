# Supported Tools

## Read Tools

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

### `generate_citation`

Generate a formatted citation for a Zotero item using the Zotero API's citeproc-js engine (supports 10,000+ CSL styles).

| Parameter | Type | Required | Description |
|---|---|---|---|
| `item_key` | string | Yes | The Zotero item key |
| `style` | string | No | Citation style: `apa` (default), `chicago`, `mla`, `harvard`, `ieee`, `bibtex`, or any CSL style ID from https://www.zotero.org/styles |
| `locale` | string | No | Locale for citation formatting (e.g. `en-US`, `de-DE`, `fr-FR`) |

## Write Tools

### `create_item`

Create a new item in the Zotero library.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `item_type` | string | Yes | Item type (e.g. `book`, `journalArticle`, `conferencePaper`) |
| `title` | string | Yes | Item title |
| `creators` | string | No | JSON array of creators, e.g. `[{"creatorType":"author","firstName":"Jane","lastName":"Doe"}]` |
| `date` | string | No | Publication date |
| `abstract_note` | string | No | Abstract or note |
| `url` | string | No | URL |
| `tags` | string | No | Comma-separated tags |
| `collections` | string | No | Comma-separated collection keys |

### `update_item`

Partially update an existing Zotero item (patch).

| Parameter | Type | Required | Description |
|---|---|---|---|
| `item_key` | string | Yes | The Zotero item key to update |
| `fields` | string | Yes | JSON object of fields to update, e.g. `{"title":"New Title","date":"2025"}` |

### `delete_items`

Delete one or more items from the Zotero library.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `item_keys` | string | Yes | Comma-separated item keys to delete |

### `add_tag`

Add a tag to a Zotero item (idempotent).

| Parameter | Type | Required | Description |
|---|---|---|---|
| `item_key` | string | Yes | The Zotero item key |
| `tag` | string | Yes | Tag to add |

### `remove_tag`

Remove a tag from a Zotero item.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `item_key` | string | Yes | The Zotero item key |
| `tag` | string | Yes | Tag to remove |

### `add_item_to_collection`

Add an item to a Zotero collection (idempotent).

| Parameter | Type | Required | Description |
|---|---|---|---|
| `item_key` | string | Yes | The Zotero item key |
| `collection_key` | string | Yes | The collection key to add the item to |

### `remove_item_from_collection`

Remove an item from a Zotero collection.

| Parameter | Type | Required | Description |
|---|---|---|---|
| `item_key` | string | Yes | The Zotero item key |
| `collection_key` | string | Yes | The collection key to remove the item from |
