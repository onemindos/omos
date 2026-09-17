# Webflow API Access - OneMind Integration

> **Status**: CONNECTED
> **Priority**: High - Enables public web presence automation
> **Domain**: UI (00-24) - Infrastructure & APIs
> **Last Verified**: January 15, 2026

---

## Key Insight: Webflow's API Ecosystem

```
★ Insight ─────────────────────────────────────────────────────────────────────
  Webflow's API Ecosystem is Surprisingly Robust:

  1. REST API v2 - Full-featured with excellent docs, 60 req/min rate limit,
     covers Sites, Pages, CMS, E-commerce, Forms, Components

  2. Official MCP Server - Webflow has an official MCP server for AI tools
     like Claude, meaning we can add it directly to Claude Code setup

  3. DevLink CLI - Not a traditional "do everything" CLI, but specifically
     for React component sync (import React → Webflow, export Webflow → React)
─────────────────────────────────────────────────────────────────────────────────
```

---

## Quick Summary

| Aspect | Status |
|--------|--------|
| **API Quality** | Excellent - REST API v2, well-documented |
| **CLI Available** | Yes - `@webflow/webflow-cli` (DevLink) |
| **MCP Server** | Yes - Official Webflow MCP for Claude/AI |
| **SDKs** | JavaScript/TypeScript official SDK |
| **Rate Limit** | 60 requests/minute per API key |

---

## Connected Site Information

### Site Details

| Property | Value |
|----------|-------|
| **Site Name** | Zeus Delacruz |
| **Site ID** | `696800f815781f16cea347f8` |
| **Workspace ID** | `685b06eb7ad095055c3819f0` |
| **Primary Domain** | [zeusdelacruz.com](https://zeusdelacruz.com) |
| **WWW Domain** | [www.zeusdelacruz.com](https://www.zeusdelacruz.com) |
| **Timezone** | America/New_York |
| **Created** | January 14, 2026 |

### Current Pages

| Page | URL | Status |
|------|-----|--------|
| Home | `/` | Published |
| Style Guide | `/style-guide` | Published |

### CMS Collections

| Collection | ID | Fields |
|------------|----|--------|
| **Blog Posts** | `69699b32876328bdfb93f48e` | Name, Slug, Post Body, Summary, Featured Image, Published Date, Author |
| **Projects** | `69699b4b1f86ef1c5a5d2e72` | Name, Slug, Description, Summary, Project URL, GitHub URL, Tech Stack, Thumbnail, Featured |
| **Services** | `69699b642c258aee48f1c3da` | Name, Slug, Description, Summary, Price, Icon, Features, CTA Link, Featured |

### API Credentials

**Site Token**: Stored securely (generated January 15, 2026)
**Token Type**: Site API Token (not Workspace token)
**Scopes**: Full read/write access

> **Security Note**: The API token should be stored in environment variables or 1Password, not in documentation. Add to `.env`:
> ```
> WEBFLOW_API_TOKEN=your-token-here
> WEBFLOW_SITE_ID=696800f815781f16cea347f8
> ```

---

## Integration Options

```
┌───────────────┬────────────────────────────────────────────────────────────────┐
│    Option     │                           Best For                             │
├───────────────┼────────────────────────────────────────────────────────────────┤
│ MCP Server    │ Direct Claude Code access to Webflow - "Zeus, publish this     │
│               │ blog post" and Claude does it immediately                      │
├───────────────┼────────────────────────────────────────────────────────────────┤
│ Letta Tools   │ Legacy AI can manage sites autonomously - scheduled publishes, │
│               │ content updates, proactive suggestions                         │
├───────────────┼────────────────────────────────────────────────────────────────┤
│ N8N Workflows │ Event-driven automation - form submission → Todoist task,      │
│               │ new CMS item → Telegram notification, scheduled publishes      │
└───────────────┴────────────────────────────────────────────────────────────────┘
```

**Recommendation**: Start with **MCP Server** for immediate Claude access, then add **N8N** for automated workflows.

---

## Real-Life OneMind Integration Examples

### UI Domain (00-24): Tech & Infrastructure

#### Example 1: Claude-Powered Content Publishing
```
Zeus: "Claude, publish my latest OneMind architecture post to the blog"

Flow:
1. Claude reads Obsidian note from vault
2. Transforms markdown → Webflow-compatible rich text
3. Creates CMS item via API
4. Publishes site
5. Returns live URL

Result: Blog post live in <30 seconds without leaving terminal
```

#### Example 2: Portfolio Project Updates
```
Zeus: "Update my portfolio with the 3 latest GitHub projects"

Flow:
1. Claude fetches recent GitHub repos (already have gh CLI)
2. Extracts project name, description, tech stack, stars
3. Creates/updates CMS items in "Projects" collection
4. Publishes changes

Result: Portfolio always reflects latest work automatically
```

---

### HP Domain (25-49): Personal Growth

#### Example 3: Fitness Journey Blog Auto-Publishing
```
Trigger: Weekly HP review completed in Obsidian

Flow:
1. N8N detects new file in "25-49 HP/Reflections/"
2. Extracts key metrics (weight, workouts, habits)
3. Generates public-safe summary (no private data)
4. Creates blog post in "Fitness Journey" collection
5. Publishes on Sunday evening

Result: Automatic weekly fitness update without manual work
```

#### Example 4: Skills Portfolio
```
Zeus: "Add my new AWS certification to the skills page"

Flow:
1. Claude creates CMS item in "Certifications" collection
2. Adds: name, date, badge image URL, verification link
3. Updates "Skills" page component
4. Publishes

Result: Professional profile always current
```

---

### LE Domain (50-74): Legacy & Family

#### Example 5: Family Chronicle Publishing
```
Scenario: Publishing approved family stories to private family site

Flow:
1. Zeus marks story as "approved for web" in Obsidian
2. N8N detects tag change
3. Transforms to family-site format
4. Creates CMS item with access controls
5. Notifies family via Telegram: "New story published!"

Result: Family history preserved and shared securely
```

#### Example 6: Memorial/Tribute Pages
```
Zeus: "Create a tribute page for Grandpa's 90th birthday"

Flow:
1. Claude gathers content from LE vault
2. Creates dedicated page with photos, stories, timeline
3. Sets up family member access
4. Sends invites to contribute

Result: Living memorial that family can add to
```

---

### GE Domain (75-99): Business & Wealth

#### Example 7: Client Inquiry → CRM → Response
```
Trigger: Form submission on business site

Flow:
1. Webflow webhook fires on form submit
2. N8N receives: name, email, inquiry type, message
3. Creates Zoho CRM lead (already integrated)
4. Creates Todoist task: "Follow up with [Name]"
5. Sends auto-response email
6. Notifies Zeus via Telegram

Result: No lead falls through cracks, instant response impression
```

#### Example 8: Service/Product Launch
```
Zeus: "Launch the new consulting package on the site"

Flow:
1. Claude reads service definition from GE vault
2. Creates CMS items: service details, pricing, FAQs
3. Updates pricing page components
4. Creates Stripe payment link (if integrated)
5. Publishes and announces

Result: New offering live across all touchpoints
```

#### Example 9: Case Study Auto-Generation
```
Trigger: Project marked "completed" in Todoist/vault

Flow:
1. N8N detects project completion
2. Legacy AI generates case study draft from project notes
3. Sends to Zeus for approval via Telegram
4. On approval: creates CMS item, adds to portfolio
5. Publishes

Result: Portfolio grows automatically with each completed project
```

---

### Cross-Domain: The Full Stack

#### Example 10: "Morning Briefing" Site Widget
```
Daily automation:

1. N8N runs at 6 AM
2. Gathers: today's tasks (Todoist), calendar events (Zoho), weather
3. Updates "Daily Focus" component on personal dashboard site
4. Visible on any device via Webflow-hosted page

Result: Personal command center accessible anywhere
```

#### Example 11: Unified Content Pipeline
```
Content creation workflow:

1. Write in Obsidian (source of truth)
2. Tag with #publish-blog or #publish-portfolio
3. N8N monitors for publish tags
4. Routes to appropriate Webflow collection
5. Applies domain-specific formatting
6. Publishes and removes tag
7. Adds "published" frontmatter with URL

Result: Write once, publish everywhere, tracked in vault
```

---

## API Capabilities Reference

### What We Can Automate

| Category | Capabilities | OneMind Use |
|----------|--------------|-------------|
| **Sites** | List, get details, publish | Multi-site management |
| **Pages** | List, get/update metadata, update content | Landing pages, static content |
| **CMS Collections** | List, create, manage fields | Blog, portfolio, services |
| **CMS Items** | Full CRUD, bulk operations, publish | Content management |
| **Components** | List, get/update content | Reusable design elements |
| **E-commerce** | Products, orders, inventory | If selling services/products |
| **Forms** | List forms, get submissions | Lead capture, inquiries |
| **Users** | Manage members | Gated content, client portals |
| **Webhooks** | Subscribe to events | Real-time automations |

### API Technical Details

| Aspect | Details |
|--------|---------|
| **Base URL** | `https://api.webflow.com/v2` |
| **Auth** | Bearer token (`Authorization: Bearer TOKEN`) |
| **Format** | JSON request/response |
| **Rate Limits** | 60 requests/minute per API key |
| **Publish Limit** | 1 site publish per minute |
| **Bulk Operations** | Up to 100 items per request |

---

## Implementation Plan

### Phase 1: Foundation (Day 1) ✅ COMPLETE
- [x] Zeus provides API token
- [x] Test API connection
- [x] List all sites and collections
- [x] Document site structure
- [x] Create CMS collections (Blog Posts, Projects, Services)
- [x] Store token in 1Password (UI vault)
- [x] Add credentials to letta-tools/.env

### Phase 2: MCP Integration (Day 1-2)
- [ ] Add Webflow MCP server to Claude Code config
- [ ] Test direct publishing from Claude
- [ ] Create publish workflow commands

### Phase 3: Letta Tools (Day 2-3)
- [ ] Create `tools/webflow/` in letta-tools
- [ ] Implement: `cms.py`, `sites.py`, `pages.py`
- [ ] Attach tools to Legacy AI agent
- [ ] Test autonomous content management

### Phase 4: N8N Automations (Day 3-5)
- [ ] Set up Webflow webhooks for forms
- [ ] Create form → Zoho CRM workflow
- [ ] Create Obsidian → Webflow publish workflow
- [ ] Create scheduled content workflows

### Phase 5: Advanced (Ongoing)
- [ ] Cross-platform content sync
- [ ] Analytics integration
- [ ] A/B testing automation
- [ ] Multi-language support

---

## Technical Implementation

### MCP Server Setup

Add to `~/.claude/mcp.json`:
```json
{
  "mcpServers": {
    "webflow": {
      "command": "npx",
      "args": ["-y", "@anthropic-ai/webflow-mcp"],
      "env": {
        "WEBFLOW_API_TOKEN": "your-token-here"
      }
    }
  }
}
```

### Letta Tools Structure

```
00-24 SO (Sovereign Operations)/letta-tools/tools/webflow/
├── __init__.py
├── cms.py          # CMS collection/item operations
├── sites.py        # Site listing, publishing
├── pages.py        # Page content management
└── forms.py        # Form submission handling
```

### N8N Workflow Templates

```
Webflow Workflows:
├── form-to-crm.json           # Form → Zoho CRM Lead
├── form-to-todoist.json       # Form → Todoist Task
├── obsidian-to-blog.json      # Vault note → Blog post
├── github-to-portfolio.json   # GitHub repo → Portfolio item
└── scheduled-publish.json     # Time-based publishing
```

---

## Example API Calls

### List Sites
```bash
curl https://api.webflow.com/v2/sites \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### List CMS Collections
```bash
curl https://api.webflow.com/v2/sites/{site_id}/collections \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Create CMS Item
```bash
curl -X POST https://api.webflow.com/v2/collections/{collection_id}/items \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fieldData": {
      "name": "My Blog Post",
      "slug": "my-blog-post",
      "post-body": "<p>Content here...</p>",
      "published-date": "2026-01-15T00:00:00Z"
    }
  }'
```

### Publish Site
```bash
curl -X POST https://api.webflow.com/v2/sites/{site_id}/publish \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"publishToWebflowSubdomain": true}'
```

---

## Resources

- [Webflow API Documentation](https://developers.webflow.com/data/reference/rest-introduction)
- [Webflow DevLink CLI](https://developers.webflow.com/devlink/reference/overview)
- [Authentication Guide](https://developers.webflow.com/data/reference/authentication)
- [Rate Limits](https://developers.webflow.com/data/reference/rate-limits)
- [CMS API Reference](https://developers.webflow.com/data/reference/cms/collections/list)
- [Webhooks Guide](https://developers.webflow.com/data/docs/working-with-webhooks)

---

## Next Steps (Backlog)

### Immediate Opportunities

| Priority | Integration | Description | Effort |
|----------|-------------|-------------|--------|
| **High** | Obsidian → Blog | N8N workflow monitors vault for `#publish-blog` tag, transforms markdown → Webflow rich text, creates CMS item | Medium |
| **High** | GitHub → Portfolio | Auto-sync latest repos to Projects collection (name, description, tech stack, stars) | Medium |
| **Medium** | MCP Server | Add to Claude Code config for direct "publish this" commands | Low |
| **Medium** | Form → CRM | Webflow form webhook → N8N → Zoho CRM lead + Todoist task | Medium |
| **Low** | Letta Tools | Python tools for Legacy AI autonomous content management | High |

### Quick Wins (Can Do Now)

1. **Test Blog Post Creation**
   ```bash
   curl -X POST "https://api.webflow.com/v2/collections/69699b32876328bdfb93f48e/items/live" \
     -H "Authorization: Bearer $WEBFLOW_API_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"fieldData": {"name": "Test Post", "slug": "test-post", "post-body": "<p>Hello from OneMind!</p>"}}'
   ```

2. **Design CMS Pages in Webflow**
   - Create Blog Post template page
   - Create Project detail page
   - Create Services listing page
   - Add Collection Lists to homepage

3. **Add MCP Server to Claude**
   ```json
   // Add to ~/.claude/mcp.json
   "webflow": {
     "command": "npx",
     "args": ["-y", "@webflow/mcp-server"],
     "env": { "WEBFLOW_API_TOKEN": "op://00-24 SO (Sovereign Operations)/Webflow Site API Token/credential" }
   }
   ```

### Future Ideas

- **Weekly Digest**: Auto-generate "This Week on Zeus Delacruz" blog post from activity
- **Portfolio Analytics**: Track which projects get the most views
- **A/B Testing**: Test different service descriptions/pricing
- **Multi-language**: Add Spanish version of key pages
- **Client Portal**: Gated content for consulting clients
- **Newsletter Integration**: Sync blog posts to email list

---

## Changelog

| Date | Change |
|------|--------|
| 2026-01-15 | Initial setup: API connected, 3 CMS collections created, token secured |

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  "The web is just another surface for OneMind. What you write in your      │
│   vault should flow effortlessly to wherever it needs to be seen."         │
│                                                    — OneMind Philosophy    │
└─────────────────────────────────────────────────────────────────────────────┘
```
