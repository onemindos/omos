# Webflow API v2 — Full Reference

> Quick lookup for endpoints, field names, and common patterns.
> Full docs: https://developers.webflow.com/v2.0.0/reference

---

## Endpoints

### Sites
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sites` | List all authorized sites |
| GET | `/sites/{siteId}` | Get site details |
| POST | `/sites/{siteId}/publish` | Publish site to domains |

### Collections
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sites/{siteId}/collections` | List collections |
| GET | `/collections/{collectionId}` | Get collection schema |

### CMS Items
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/collections/{collectionId}/items` | List items (paginated) |
| GET | `/collections/{collectionId}/items/{itemId}` | Get single item |
| POST | `/collections/{collectionId}/items` | Create item (draft) |
| PATCH | `/collections/{collectionId}/items/{itemId}` | Update item |
| DELETE | `/collections/{collectionId}/items/{itemId}` | Archive item |
| POST | `/collections/{collectionId}/items/publish` | Publish item(s) live |

### Pages
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sites/{siteId}/pages` | List all pages |
| GET | `/pages/{pageId}` | Get page metadata |
| GET | `/pages/{pageId}/dom` | Get page DOM (static content) |
| POST | `/pages/{pageId}/dom` | Update static page content |

### Forms
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sites/{siteId}/forms` | List forms |
| GET | `/forms/{formId}/submissions` | Get form submissions |

### Webhooks
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sites/{siteId}/webhooks` | List webhooks |
| POST | `/sites/{siteId}/webhooks` | Create webhook |
| DELETE | `/sites/{siteId}/webhooks/{webhookId}` | Delete webhook |

---

## Standard CMS Item Fields

These are built-in fields present on every CMS item:

| Field | Type | Notes |
|-------|------|-------|
| `name` | string | Required. Display name. |
| `slug` | string | URL path segment. Auto-generated from name if omitted. |
| `_archived` | boolean | `true` = archived (hidden) |
| `_draft` | boolean | `true` = draft (staged, not live) |

Custom fields depend on your collection schema — fetch schema via `GET /collections/{id}`.

---

## Pagination

All list endpoints return paginated results:

```json
{
  "items": [...],
  "pagination": {
    "limit": 100,
    "offset": 0,
    "total": 342
  }
}
```

To get all items:
```bash
# Page through with offset
curl "https://api.webflow.com/v2/collections/$COL_ID/items?limit=100&offset=0" ...
curl "https://api.webflow.com/v2/collections/$COL_ID/items?limit=100&offset=100" ...
```

Python helper:
```python
def wf_get_all(collection_id):
    items = []
    offset = 0
    while True:
        data = wf_get(f"/collections/{collection_id}/items?limit=100&offset={offset}")
        batch = data.get("items", [])
        items.extend(batch)
        if len(batch) < 100:
            break
        offset += 100
    return items
```

---

## Webhook Payload Examples

### form_submission
```json
{
  "triggerType": "form_submission",
  "payload": {
    "name": "Contact Form",
    "data": {
      "Name": "John Doe",
      "Email": "john@example.com",
      "Message": "I want to work with you."
    },
    "siteId": "..."
  }
}
```

### site_publish
```json
{
  "triggerType": "site_publish",
  "payload": {
    "siteId": "...",
    "publishedAt": "2026-04-10T12:00:00Z"
  }
}
```

---

## CLI Config File

`~/.webflow/config.json` — stores auth token and default site:

```json
{
  "token": "your-api-token",
  "defaultSiteId": "your-site-id"
}
```

---

## Error Codes

| Code | Meaning |
|------|---------|
| 400 | Bad request — check field names / types |
| 401 | Unauthorized — check token |
| 404 | Resource not found — check IDs |
| 429 | Rate limited — slow down, add sleep |
| 500 | Webflow server error — retry |


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
