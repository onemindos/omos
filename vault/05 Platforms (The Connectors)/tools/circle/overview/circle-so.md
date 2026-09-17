# Circle.so MCP

**Type:** Community platform integration
**MCP Package:** `@circle-so/mcp` (or community server)
**Website:** circle.so
**Auth:** Circle API token

## What It Does

Circle.so MCP lets agents read and write to Circle communities — query members, post to spaces, create/update posts, and trigger automations. Key use case: Legacy agent can surface insights from Zeus's communities, post content, or monitor activity.

## Key Tools Exposed

| Tool | Description |
|------|-------------|
| `list_spaces` | Get all spaces in a community |
| `get_posts` | Fetch posts from a space with filters |
| `create_post` | Publish a new post to a space |
| `get_members` | List/search community members |
| `get_member` | Get member profile + activity |
| `add_member` | Add a member to the community |
| `create_comment` | Comment on a post |
| `get_events` | List upcoming/past events |

## Config

```json
{
  "mcpServers": {
    "circle": {
      "command": "npx",
      "args": ["-y", "@circle-so/mcp-server"],
      "env": {
        "CIRCLE_API_TOKEN": "your-token",
        "CIRCLE_COMMUNITY_ID": "your-community-id"
      }
    }
  }
}
```

## Full MCP Capability List

Everything the Circle MCP can do when connected to an agent.

### Posts & Content
- Create, update, delete posts in any space
- Create rich text content (headings, bullets, bold, code blocks, embeds)
- Create image posts, event posts
- Pin posts
- Create/delete comments on posts

### Courses
- Create course sections (modules) inside any course space
- Create lessons inside sections
- Update lesson content
- Track lesson progress per member

### Events
- Create live events with title, date, time, description
- Manage event attendees

### Members
- Update any member profile (headline, bio, social links, location)
- Create/assign member tags
- Add/remove members from spaces or space groups
- Create community segments (filter members by criteria)
- Ban or deactivate members
- Search members by email

### Access & Tiers
- Create access groups (for gating spaces by tier)
- Assign members to access groups
- Remove members from access groups

### Invitations
- Create invitation links (public or restricted)
- Revoke invitation links
- List active invite links

### Community Settings
- Update community name, description, settings
- Update space names, descriptions, display type
- Update space group names

### Organization
- Create topics inside spaces (tags for posts)
- Create/update/delete spaces
- Create/update/delete space groups

### Moderation
- Flag content
- List flagged content

### Analytics
- Gamification leaderboard
- Member activity data

---

## Stack Use Cases

- **Empire agent:** Monitor community health, engagement metrics, member growth
- **Legacy agent:** Answer questions about community members or post on Zeus's behalf
- **Vault agent:** Archive notable discussions into codex
- **Oracle agent:** Research member profiles before outreach

## API Token Location

1Password → "00-24 SO (Sovereign Operations)" vault → Circle API token


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
