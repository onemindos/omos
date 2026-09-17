---
name: circle
description: Manage Circle.so community platform. Use when you need to create posts, manage members, moderate spaces, or query community activity.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Circle.so Integration Skill

Interact with Circle.so community platform API.

## Environment Variables

```bash
CIRCLE_API_KEY       # API key from Circle settings
CIRCLE_COMMUNITY_ID  # Your community ID
```

## Base URL

```
https://app.circle.so/api/v1
```

## Authentication

```bash
Authorization: Token $CIRCLE_API_KEY
```

---

## Operations

### Spaces

**List spaces:**
```bash
curl -s "https://app.circle.so/api/v1/spaces?community_id=$CIRCLE_COMMUNITY_ID" \
  -H "Authorization: Token $CIRCLE_API_KEY" | jq '.[] | {id, name, slug}'
```

**Get space:**
```bash
curl -s "https://app.circle.so/api/v1/spaces/{space_id}" \
  -H "Authorization: Token $CIRCLE_API_KEY"
```

### Posts

**List posts in space:**
```bash
curl -s "https://app.circle.so/api/v1/posts?space_id={space_id}&community_id=$CIRCLE_COMMUNITY_ID" \
  -H "Authorization: Token $CIRCLE_API_KEY" | jq '.[] | {id, name, created_at}'
```

**Create post:**
```bash
curl -s -X POST "https://app.circle.so/api/v1/posts" \
  -H "Authorization: Token $CIRCLE_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "community_id": "'"$CIRCLE_COMMUNITY_ID"'",
    "space_id": "{space_id}",
    "name": "Post Title",
    "body": "Post content in markdown...",
    "is_pinned": false,
    "is_comments_enabled": true
  }'
```

**Update post:**
```bash
curl -s -X PUT "https://app.circle.so/api/v1/posts/{post_id}" \
  -H "Authorization: Token $CIRCLE_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Title",
    "body": "Updated content"
  }'
```

**Delete post:**
```bash
curl -s -X DELETE "https://app.circle.so/api/v1/posts/{post_id}" \
  -H "Authorization: Token $CIRCLE_API_KEY"
```

### Members

**List members:**
```bash
curl -s "https://app.circle.so/api/v1/community_members?community_id=$CIRCLE_COMMUNITY_ID" \
  -H "Authorization: Token $CIRCLE_API_KEY" | jq '.[] | {id, name, email}'
```

**Get member:**
```bash
curl -s "https://app.circle.so/api/v1/community_members/{member_id}" \
  -H "Authorization: Token $CIRCLE_API_KEY"
```

**Invite member:**
```bash
curl -s -X POST "https://app.circle.so/api/v1/community_members" \
  -H "Authorization: Token $CIRCLE_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "community_id": "'"$CIRCLE_COMMUNITY_ID"'",
    "email": "newuser@example.com",
    "name": "New User",
    "space_ids": [123, 456]
  }'
```

### Comments

**List comments on post:**
```bash
curl -s "https://app.circle.so/api/v1/comments?post_id={post_id}" \
  -H "Authorization: Token $CIRCLE_API_KEY"
```

**Create comment:**
```bash
curl -s -X POST "https://app.circle.so/api/v1/comments" \
  -H "Authorization: Token $CIRCLE_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "post_id": "{post_id}",
    "body": "Comment text..."
  }'
```

---

## One Mind Integration

Circle community maps to GE pillar:
- **Community** → `40 Community (Extended Network)` or `83 Partnerships`
- **Courses** → `76 Legacy Acadacom (Education)`

## Tips

- Business plan required for full API access
- Posts support markdown formatting
- Use `is_pinned: true` for announcements
- Rate limits: ~100 requests/minute


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
