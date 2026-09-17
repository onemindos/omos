---
name: likekit-platform
description: Interact with Likekit platform API. Use when you need to manage Likekit resources, integrations, or workflows.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
---

# Likekit Skill

Interact with the Likekit platform.

## Environment Variables

```bash
LIKEKIT_API_KEY      # Your Likekit API key
LIKEKIT_BASE_URL     # API base URL (if not default)
```

## Base URL

```
https://api.likekit.com  # Update with actual base URL
```

## Authentication

```bash
Authorization: Bearer $LIKEKIT_API_KEY
# or
X-API-Key: $LIKEKIT_API_KEY
```

## Operations

### List Resources
```bash
curl -s "https://api.likekit.com/resources" \
  -H "Authorization: Bearer $LIKEKIT_API_KEY" | jq
```

### Get Resource
```bash
curl -s "https://api.likekit.com/resources/{id}" \
  -H "Authorization: Bearer $LIKEKIT_API_KEY" | jq
```

### Create Resource
```bash
curl -s -X POST "https://api.likekit.com/resources" \
  -H "Authorization: Bearer $LIKEKIT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Resource Name",
    "data": {}
  }'
```

### Update Resource
```bash
curl -s -X PATCH "https://api.likekit.com/resources/{id}" \
  -H "Authorization: Bearer $LIKEKIT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Name"
  }'
```

### Delete Resource
```bash
curl -s -X DELETE "https://api.likekit.com/resources/{id}" \
  -H "Authorization: Bearer $LIKEKIT_API_KEY"
```

---

## TODO: Customize This Skill

This is a template skill. Update with actual Likekit API documentation:

1. Correct base URL
2. Authentication method
3. Available endpoints
4. Request/response formats
5. Platform-specific operations

Refer to Likekit documentation for accurate API details.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
