---
name: api-integration
description: >
  [REPLACE THIS] Integrates with [Service Name] API. Use when the user needs
  to [list specific actions]. Supports [list key operations like: create, read,
  update, delete, list, search].
metadata:
  author: zeus
  version: "1.0.0"
---

# [Service Name] Integration

## When to Use

Activate this skill when the user:

- Asks to [create / list / update / delete] [resource type]
- Mentions "[service name]", "[resource type]", or "[key phrase]"
- Wants to [describe primary use case]

## Environment Variables Required

| Variable             | Description                                |
| -------------------- | ------------------------------------------ |
| `MY_SERVICE_URL`     | Base URL (e.g., `https://api.example.com`) |
| `MY_SERVICE_API_KEY` | API key or Bearer token                    |

## Available Tools

| Tool         | Description                          |
| ------------ | ------------------------------------ |
| `api_list`   | List resources with optional filters |
| `api_get`    | Get a single resource by ID          |
| `api_create` | Create a new resource                |
| `api_update` | Update an existing resource          |
| `api_delete` | Delete a resource                    |

## Quick Reference

### List all items

```
api_list({ limit: 20 })
```

### List with filter

```
api_list({ filter: "status=active", limit: 10 })
```

### Get a specific item

```
api_get({ id: "item-id-here" })
```

### Create a new item

```
api_create({
  name: "Item Name",
  description: "What this item is for",
  status: "active"
})
```

### Update an item

```
api_update({
  id: "item-id-here",
  name: "Updated Name",
  status: "inactive"
})
```

### Delete an item

```
api_delete({ id: "item-id-here" })
```

## Error Handling

- **401 Unauthorized**: `MY_SERVICE_API_KEY` is invalid or expired
- **404 Not Found**: The `id` does not exist — use `api_list` to find valid IDs
- **429 Rate Limited**: Too many requests — wait a moment and retry
- **Timeout**: The API did not respond within 30 seconds — retry once

## Notes

- IDs are strings, not integers
- Pagination: use `limit` and `offset` for large result sets
- All timestamps are ISO 8601 UTC


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
