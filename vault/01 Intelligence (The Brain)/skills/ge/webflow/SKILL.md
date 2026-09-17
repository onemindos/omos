---
name: webflow
description: "Manages Webflow sites, CMS collections, pages, and publishing for the One Mind personal brand via the webflow CLI and REST API v2. Use when publishing content, updating CMS items, deploying site changes, or automating brand content. Maps to GE/75 Brand and GE/78 Marketing."
allowed-tools:
  - Bash(webflow:*)
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Webflow Skill

Manage Webflow sites via the `webflow` CLI and REST API v2. Primary platform for One Mind personal brand site and content engine.

## Codex Mapping

| Webflow Layer  | One Mind Domain                    |
| -------------- | ---------------------------------- |
| Site / Pages   | GE/75 Brand (Personal & Corporate) |
| CMS Blog Posts | GE/78 Marketing (Content & Growth) |
| Product Pages  | GE/77 Products (Offerings & IP)    |
| Forms / Leads  | GE/83 Partnerships (Network)       |

---

## Environment Variables Required

```bash
WEBFLOW_API_TOKEN   # From Webflow dashboard → Account Settings → Integrations → API Access
WEBFLOW_SITE_ID     # From: webflow sites list  OR  API → GET /sites
```

Store both in 1Password vault `00-24 SO (Sovereign Operations)` under item `Webflow`.

---

## CLI Setup

```bash
# Install
npm install -g @webflow/webflow-cli

# Authenticate
webflow auth login               # Opens browser OAuth flow
webflow auth token               # Set token directly
webflow whoami                   # Verify auth

# Set default site (saves to ~/.webflow/config.json)
webflow sites list               # Find your site ID
```

---

## CLI — Quick Reference

### Sites

```bash
webflow sites list               # List all sites + IDs
webflow sites get --siteId $WEBFLOW_SITE_ID   # Site details
```

### Publishing

```bash
# Publish to live domain
webflow sites publish --siteId $WEBFLOW_SITE_ID --domains your-domain.com

# Publish to staging (webflow.io subdomain)
webflow sites publish --siteId $WEBFLOW_SITE_ID
```

### CMS Collections

```bash
webflow collections list --siteId $WEBFLOW_SITE_ID
webflow collections get --siteId $WEBFLOW_SITE_ID --collectionId $COL_ID
```

### CMS Items

```bash
# List items in a collection
webflow items list --siteId $WEBFLOW_SITE_ID --collectionId $COL_ID

# Get single item
webflow items get --siteId $WEBFLOW_SITE_ID --collectionId $COL_ID --itemId $ITEM_ID

# Create item (staged, not live)
webflow items create --siteId $WEBFLOW_SITE_ID --collectionId $COL_ID \
  --fields '{"name":"Post Title","slug":"post-slug","body":"<p>Content</p>"}'

# Publish item live
webflow items publish --siteId $WEBFLOW_SITE_ID --collectionId $COL_ID --itemIds $ITEM_ID
```

---

## API v2 — Base URL

```
https://api.webflow.com/v2
```

### Auth Header

```bash
-H "Authorization: Bearer $WEBFLOW_API_TOKEN"
-H "accept: application/json"
-H "content-type: application/json"
```

---

## API — Common Patterns

### List All Sites

```bash
curl -s "https://api.webflow.com/v2/sites" \
  -H "Authorization: Bearer $WEBFLOW_API_TOKEN" | jq '.sites[] | {id, displayName, shortName}'
```

### List Collections for a Site

```bash
curl -s "https://api.webflow.com/v2/sites/$WEBFLOW_SITE_ID/collections" \
  -H "Authorization: Bearer $WEBFLOW_API_TOKEN" | jq '.collections[] | {id, displayName, slug}'
```

### List CMS Items

```bash
curl -s "https://api.webflow.com/v2/collections/$COLLECTION_ID/items" \
  -H "Authorization: Bearer $WEBFLOW_API_TOKEN" | jq '.items[] | {id, fieldData: .fieldData.name}'
```

### Create CMS Item (Draft)

```bash
curl -s -X POST "https://api.webflow.com/v2/collections/$COLLECTION_ID/items" \
  -H "Authorization: Bearer $WEBFLOW_API_TOKEN" \
  -H "content-type: application/json" \
  -d '{
    "fieldData": {
      "name": "Post Title",
      "slug": "post-slug",
      "post-body": "<p>Your content here.</p>",
      "post-summary": "Short summary for SEO.",
      "_archived": false,
      "_draft": true
    }
  }' | jq '{id: .id, name: .fieldData.name}'
```

### Publish CMS Items Live

```bash
# Publish specific items from draft → live
curl -s -X POST "https://api.webflow.com/v2/collections/$COLLECTION_ID/items/publish" \
  -H "Authorization: Bearer $WEBFLOW_API_TOKEN" \
  -H "content-type: application/json" \
  -d '{"itemIds": ["ITEM_ID_1", "ITEM_ID_2"]}' | jq '.'
```

### Update CMS Item

```bash
curl -s -X PATCH "https://api.webflow.com/v2/collections/$COLLECTION_ID/items/$ITEM_ID" \
  -H "Authorization: Bearer $WEBFLOW_API_TOKEN" \
  -H "content-type: application/json" \
  -d '{"fieldData": {"name": "Updated Title", "_draft": false}}' | jq '{id, name: .fieldData.name}'
```

### Publish Site (trigger deploy)

```bash
curl -s -X POST "https://api.webflow.com/v2/sites/$WEBFLOW_SITE_ID/publish" \
  -H "Authorization: Bearer $WEBFLOW_API_TOKEN" \
  -H "content-type: application/json" \
  -d '{"publishToWebflowSubdomain": false}' | jq '.'
```

---

## Python — Agent-Ready Pattern

For use in OpenClaw tools or N8N Python nodes:

```python
import os, json, urllib.request

TOKEN = os.environ["WEBFLOW_API_TOKEN"]
SITE_ID = os.environ["WEBFLOW_SITE_ID"]

def wf_get(path):
    req = urllib.request.Request(
        f"https://api.webflow.com/v2{path}",
        headers={"Authorization": f"Bearer {TOKEN}", "accept": "application/json"}
    )
    with urllib.request.urlopen(req, timeout=15) as r:
        return json.loads(r.read())

def wf_post(path, body):
    data = json.dumps(body).encode()
    req = urllib.request.Request(
        f"https://api.webflow.com/v2{path}",
        data=data,
        headers={
            "Authorization": f"Bearer {TOKEN}",
            "content-type": "application/json",
            "accept": "application/json"
        },
        method="POST"
    )
    with urllib.request.urlopen(req, timeout=15) as r:
        return json.loads(r.read())

def wf_patch(path, body):
    data = json.dumps(body).encode()
    req = urllib.request.Request(
        f"https://api.webflow.com/v2{path}",
        data=data,
        headers={
            "Authorization": f"Bearer {TOKEN}",
            "content-type": "application/json",
            "accept": "application/json"
        },
        method="PATCH"
    )
    with urllib.request.urlopen(req, timeout=15) as r:
        return json.loads(r.read())

# Example: publish a blog post from codex content
def publish_post(collection_id, title, slug, body_html, summary):
    item = wf_post(f"/collections/{collection_id}/items", {
        "fieldData": {
            "name": title,
            "slug": slug,
            "post-body": body_html,
            "post-summary": summary,
            "_archived": False,
            "_draft": False
        }
    })
    item_id = item["id"]
    wf_post(f"/collections/{collection_id}/items/publish", {"itemIds": [item_id]})
    return item_id
```

---

## N8N Integration

An N8N HTTP Request node can call the Webflow API directly:

- **Method**: POST / GET / PATCH
- **URL**: `https://api.webflow.com/v2/collections/{{$env.WEBFLOW_COLLECTION_ID}}/items`
- **Auth**: Header Auth → `Authorization: Bearer {{$env.WEBFLOW_API_TOKEN}}`
- **Trigger**: Webhook from codex / Todoist task completion / scheduled publish

Planned N8N flow: `Codex note (Markdown) → Empire agent → Webflow CMS draft → review → publish`

---

## Webhook Events (Inbound)

Configure in Webflow Dashboard → Site Settings → Integrations → Webhooks.

| Event             | Trigger                                     |
| ----------------- | ------------------------------------------- |
| `form_submission` | Lead captured → route to GE/83 Partnerships |
| `site_publish`    | Deploy complete → notify via Telegram       |
| `ecomm_new_order` | Product sold → route to GE/80 Treasury      |

Webhook receiver lives on VPS: `http://104.248.4.212:PORT/webhooks/webflow`

---

## One Mind Brand Content Pipeline

```
GE/77 Products or GE/78 Marketing (Markdown note)
    → Empire agent (OpenClaw) reads + formats
    → wf_post() creates CMS draft
    → Zeus reviews in Webflow editor
    → wf_post(publish) goes live
    → Site publish triggers Cloudflare CDN flush
```

---

## Useful IDs (fill in after setup)

```bash
# Run these once and save to 1Password / codex:
webflow sites list              # → WEBFLOW_SITE_ID
webflow collections list --siteId $WEBFLOW_SITE_ID  # → blog collection ID, etc.
```

| Variable                     | Where to store                           |
| ---------------------------- | ---------------------------------------- |
| `WEBFLOW_API_TOKEN`          | 1Password → Webflow → api-token          |
| `WEBFLOW_SITE_ID`            | 1Password → Webflow → site-id            |
| `WEBFLOW_BLOG_COLLECTION_ID` | 1Password → Webflow → blog-collection-id |

---

## Rate Limits

| Plan         | Requests/min |
| ------------ | ------------ |
| Free / Basic | 60           |
| CMS          | 60           |
| Business+    | 120          |

Add `time.sleep(1)` between bulk operations to stay safe.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
