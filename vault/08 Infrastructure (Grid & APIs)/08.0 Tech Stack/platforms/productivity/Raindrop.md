# Raindrop.io

> **Bookmark Manager**

## Overview

Raindrop.io is a bookmark manager for organizing web resources, articles, and reference materials. It serves as the external brain for web content.

## Status

- **Active**: ✅ Yes
- **Plan**: Pro
- **Last Updated**: January 25, 2026

## Use Cases

1. **Article Saving** - Save articles for later
2. **Resource Library** - Organize references
3. **Research** - Collect research materials
4. **Inspiration** - Design and ideas
5. **Sharing** - Public collections

## Key Features

- Full-text search
- Automatic tagging
- Permanent copies
- Nested collections
- Browser extensions
- Mobile apps
- Public/private collections

## Collections

| Collection       | Purpose            | Visibility |
| ---------------- | ------------------ | ---------- |
| **Inbox**        | Unsorted saves     | Private    |
| **Reading List** | Articles to read   | Private    |
| **AI/ML**        | AI resources       | Private    |
| **Development**  | Coding resources   | Private    |
| **Design**       | UI/UX inspiration  | Private    |
| **Business**     | Business resources | Private    |
| **Tools**        | Software tools     | Private    |
| **Shared**       | Public resources   | Public     |

## API Integration

### Configuration

```python
import requests

RAINDROP_TOKEN = os.getenv("RAINDROP_TOKEN")

headers = {
    "Authorization": f"Bearer {RAINDROP_TOKEN}",
    "Content-Type": "application/json"
}

base_url = "https://api.raindrop.io/rest/v1"
```

### Get Collections

```python
def get_collections():
    response = requests.get(
        f"{base_url}/collections",
        headers=headers
    )
    return response.json()
```

### Get Bookmarks

```python
def get_bookmarks(collection_id: int = 0, search: str = None, page: int = 0):
    params = {"page": page, "perpage": 50}
    if search:
        params["search"] = search

    response = requests.get(
        f"{base_url}/raindrops/{collection_id}",
        headers=headers,
        params=params
    )
    return response.json()
```

### Create Bookmark

```python
def create_bookmark(url: str, title: str = None, collection_id: int = -1, tags: list = None):
    data = {
        "link": url,
        "collection": {"$id": collection_id},
        "pleaseParse": {}  # Auto-parse metadata
    }
    if title:
        data["title"] = title
    if tags:
        data["tags"] = tags

    response = requests.post(
        f"{base_url}/raindrop",
        headers=headers,
        json=data
    )
    return response.json()
```

### Update Bookmark

```python
def update_bookmark(raindrop_id: int, updates: dict):
    response = requests.put(
        f"{base_url}/raindrop/{raindrop_id}",
        headers=headers,
        json=updates
    )
    return response.json()

# Move to different collection
update_bookmark(123, {"collection": {"$id": 456}})

# Add tags
update_bookmark(123, {"tags": ["ai", "research"]})
```

### Delete Bookmark

```python
def delete_bookmark(raindrop_id: int):
    response = requests.delete(
        f"{base_url}/raindrop/{raindrop_id}",
        headers=headers
    )
    return response.json()
```

### Search

```python
def search_bookmarks(query: str, collection_id: int = 0):
    return get_bookmarks(collection_id, search=query)

# Full-text search
results = search_bookmarks("machine learning tutorial")
```

## Tags

### Common Tags

| Tag            | Purpose            |
| -------------- | ------------------ |
| `#toread`      | Need to read       |
| `#reference`   | Keep for reference |
| `#tutorial`    | Learning material  |
| `#tool`        | Software/service   |
| `#inspiration` | Design/ideas       |
| `#important`   | Priority items     |

### Get Tag Cloud

```python
def get_tags():
    response = requests.get(
        f"{base_url}/tags",
        headers=headers
    )
    return response.json()
```

## Highlights

### Get Highlights

```python
def get_highlights(raindrop_id: int):
    response = requests.get(
        f"{base_url}/raindrop/{raindrop_id}",
        headers=headers
    )
    return response.json().get("item", {}).get("highlights", [])
```

### Create Highlight

```python
def add_highlight(raindrop_id: int, text: str, note: str = None):
    bookmark = get_bookmark(raindrop_id)
    highlights = bookmark.get("item", {}).get("highlights", [])

    highlights.append({
        "text": text,
        "note": note or "",
        "created": datetime.utcnow().isoformat()
    })

    return update_bookmark(raindrop_id, {"highlights": highlights})
```

## Legacy AI Integration

### Save Resource

```python
@Tool
def save_bookmark(url: str, tags: list = None, collection: str = "Inbox") -> str:
    """Save a URL to Raindrop.io for later reading."""

    # Get collection ID
    collections = get_collections()
    collection_id = next(
        (c["_id"] for c in collections.get("items", [])
         if c["title"] == collection),
        -1  # Inbox
    )

    result = create_bookmark(url, collection_id=collection_id, tags=tags)

    return f"Saved: {result.get('item', {}).get('title', url)}"
```

### Search Resources

```python
@Tool
def search_saved_resources(query: str) -> list:
    """Search through saved bookmarks in Raindrop.io."""

    results = search_bookmarks(query)

    return [
        {
            "title": item["title"],
            "url": item["link"],
            "tags": item.get("tags", []),
            "excerpt": item.get("excerpt", "")
        }
        for item in results.get("items", [])[:10]
    ]
```

### Daily Reading List

```python
@Tool
def get_reading_list(limit: int = 5) -> list:
    """Get items from the reading list."""

    # Get "Reading List" collection
    collections = get_collections()
    reading_list_id = next(
        (c["_id"] for c in collections.get("items", [])
         if c["title"] == "Reading List"),
        None
    )

    if not reading_list_id:
        return []

    bookmarks = get_bookmarks(reading_list_id)

    return [
        {"title": item["title"], "url": item["link"]}
        for item in bookmarks.get("items", [])[:limit]
    ]
```

## Browser Extension

### Quick Save

1. Install Raindrop extension
2. Click extension icon on any page
3. Select collection and add tags
4. Save

### Keyboard Shortcut

- **Chrome/Firefox**: `Alt+Shift+S` (customizable)

## Backup

### Export

```python
def export_bookmarks(collection_id: int = 0, format: str = "html"):
    # Export via API or dashboard
    pass
```

### Dashboard Export

1. Go to Settings → Backup
2. Choose format (HTML, CSV, JSON)
3. Download

## Obsidian Integration

### Sync to Obsidian

```python
def sync_to_obsidian(collection_id: int, obsidian_path: str):
    bookmarks = get_bookmarks(collection_id)

    for item in bookmarks.get("items", []):
        filename = f"{slugify(item['title'])}.md"
        content = f"""# {item['title']}

URL: {item['link']}
Tags: {', '.join(item.get('tags', []))}
Saved: {item['created']}

## Notes

{item.get('excerpt', '')}

## Highlights

{format_highlights(item.get('highlights', []))}
"""
        with open(f"{obsidian_path}/{filename}", "w") as f:
            f.write(content)
```

## Best Practices

- ✅ Process Inbox regularly
- ✅ Use consistent tagging
- ✅ Create collection hierarchy
- ✅ Add notes to important saves
- ✅ Use highlights for key info
- ✅ Export backups monthly
- ✅ Share useful collections publicly

## Resources

- [Raindrop.io](https://raindrop.io/)
- [API Documentation](https://developer.raindrop.io/)
- [Browser Extensions](https://raindrop.io/download)

---

**Related**: [Obsidian](./Obsidian.md) | [Todoist](./Todoist.md)
