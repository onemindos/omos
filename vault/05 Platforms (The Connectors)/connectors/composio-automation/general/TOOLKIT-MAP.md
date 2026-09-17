---
id: composio-toolkit-map
tags: [composio, tools, integrations, social-media, automation, SO]
owner: zeus
status: active
updated: 2026-05-05
---

# Composio Toolkit Map

Full reference of all toolkits available via Composio CLI for One Mind OS agents. Use `composio link <toolkit>` to connect any account, then `composio execute <SLUG>` to run tools.

---

## Currently Connected Accounts

| Toolkit | Status | Notes |
|---------|--------|-------|
| **Dropbox** | ✅ Connected | File storage, search, share |
| **Instagram** | ✅ Connected | Post, carousel, publish |
| **Facebook** | ✅ Connected | Page posts, scheduling |
| **GitHub** | ✅ Connected | Repos, files, commits |
| **Gemini** | ✅ Connected | Video + image generation |
| **YouTube** | ✅ Connected | Upload, update, search |
| **Todoist** | ✅ Connected | Tasks, projects, labels |
| **Strava** | ✅ Connected | Fitness, activities, routes |

---

## Social Media

| Toolkit | Key Slugs | What It Can Do |
|---------|-----------|---------------|
| **Twitter/X** | `TWITTER_CREATION_OF_A_POST` `TWITTER_UPLOAD_MEDIA` `TWITTER_USER_LOOKUP_ME` | Post tweets, threads, upload media, delete posts |
| **Instagram** | `INSTAGRAM_POST_IG_USER_MEDIA` `INSTAGRAM_POST_IG_USER_MEDIA_PUBLISH` `INSTAGRAM_CREATE_CAROUSEL_CONTAINER` | Post photos, carousels, reels; get media info |
| **Facebook** | `FACEBOOK_CREATE_POST` `FACEBOOK_CREATE_PHOTO_POST` `FACEBOOK_GET_SCHEDULED_POSTS` `FACEBOOK_RESCHEDULE_POST` | Page posts, photo posts, schedule, reschedule, publish |
| **LinkedIn** | `LINKEDIN_CREATE_LINKED_IN_POST` `LINKEDIN_REGISTER_IMAGE_UPLOAD` `LINKEDIN_GET_MY_INFO` | Create posts, image uploads, company posts |
| **TikTok** | `TIKTOK_UPLOAD_VIDEO` `TIKTOK_PUBLISH_VIDEO` `TIKTOK_POST_PHOTO` `TIKTOK_FETCH_PUBLISH_STATUS` | Upload + publish video/photo, check status |
| **YouTube** | `YOUTUBE_UPLOAD_VIDEO` `YOUTUBE_UPDATE_VIDEO` `YOUTUBE_UPDATE_THUMBNAIL` `YOUTUBE_SEARCH_YOU_TUBE` | Upload, update metadata, thumbnails, search, delete |
| **Discord** | `DISCORDBOT_LIST_GUILD_CHANNELS` `DISCORD_LIST_MY_GUILDS` `DISCORDBOT_GET_GUILD` | List servers/channels, send messages, manage guild |

### Connect Social Accounts

```bash
composio link twitter
composio link linkedin
composio link tiktok
composio link discord
```

---

## Email & Calendar

| Toolkit | Key Slugs | What It Can Do |
|---------|-----------|---------------|
| **Gmail** | `GMAIL_SEND_EMAIL` `GMAIL_CREATE_EMAIL_DRAFT` `GMAIL_SEND_DRAFT` `GMAIL_FETCH_EMAILS` | Send, draft, fetch, search emails |
| **Outlook** | `OUTLOOK_SEND_EMAIL` `OUTLOOK_CREATE_DRAFT` `OUTLOOK_GET_MAIL_TIPS` | Send, draft, delivery preflight |
| **Resend** | `RESEND_SEND_EMAIL` | Transactional email fallback |
| **Google Calendar** | `GOOGLECALENDAR_CREATE_EVENT` `GOOGLECALENDAR_FIND_FREE_SLOTS` `GOOGLECALENDAR_FIND_EVENT` `GOOGLECALENDAR_PATCH_EVENT` | Create/update events, find availability |
| **Cal.com** | `CAL_GET_AVAILABLE_SLOTS_INFO` `CAL_POST_NEW_BOOKING_REQUEST` | Booking pages, slot availability |
| **Zoom** | `ZOOM_LIST_MEETINGS` `ZOOM_GET_A_MEETING` `ZOOM_UPDATE_A_MEETING` `ZOOM_GET_MEETING_RECORDINGS` | Create/update/list meetings, get recordings |

### Connect Email/Calendar

```bash
composio link gmail
composio link googlecalendar
composio link zoom
```

---

## Productivity & Docs

| Toolkit | Key Slugs | What It Can Do |
|---------|-----------|---------------|
| **Notion** | `NOTION_CREATE_PAGE` `NOTION_UPDATE_PAGE` `NOTION_QUERY_DATABASE` | Create/update pages, query databases |
| **Airtable** | `AIRTABLE_LIST_RECORDS` `AIRTABLE_CREATE_RECORD` `AIRTABLE_UPDATE_RECORD` | Read/write records and tables |
| **Dropbox** | `DROPBOX_SEARCH_FILE_OR_FOLDER` `DROPBOX_READ_FILE` `DROPBOX_CREATE_SHARED_LINK` `DROPBOX_GET_TEMPORARY_LINK` | Search, read, share files |
| **Todoist** | 84 tools — see [[SKILL]] for full reference | Tasks, projects, sections, labels, comments |
| **Slack** | `SLACK_SEND_MESSAGE_AS_USER` `SLACK_LIST_CHANNELS` | Send messages, list channels |

### Connect Productivity Tools

```bash
composio link notion
composio link airtable
composio link slack
```

---

## Development

| Toolkit | Key Slugs | What It Can Do |
|---------|-----------|---------------|
| **GitHub** | `GITHUB_GET_A_REPOSITORY` `GITHUB_GET_REPOSITORY_CONTENT` `GITHUB_SEARCH_CODE` `GITHUB_LIST_COMMITS` `GITHUB_GET_A_TREE` | Repos, file content, code search, commits, branches |

---

## CRM & Marketing

| Toolkit | Key Slugs | What It Can Do |
|---------|-----------|---------------|
| **HubSpot** | `HUBSPOT_GET_ALL_MARKETING_EMAILS_FOR_A_HUB_SPOT_ACCOUNT` `HUBSPOT_GET_AGGREGATED_STATISTICS` `HUBSPOT_GET_CAMPAIGNS` | Marketing emails, campaigns, analytics |
| **Zoho** | CRM records, contacts, leads | Full CRM suite |
| **Close** | Leads, contacts, activities | Sales CRM |

### Connect CRM

```bash
composio link hubspot
composio link zoho
composio link close
```

---

## E-commerce & Finance

| Toolkit | Key Slugs | What It Can Do |
|---------|-----------|---------------|
| **Shopify** | `SHOPIFY_GET_SHOP_DETAILS` `SHOPIFY_GET_POLICIES` `SHOPIFY_GET_THEMES` | Shop config, products, themes, policies |
| **Stripe** | Payments, invoices, customers | Full payments API |
| **QuickBooks** | Invoicing, accounting, reports | Finance/accounting |

### Connect E-commerce

```bash
composio link shopify
composio link stripe
composio link quickbooks
```

---

## AI & Video Generation

| Toolkit | Key Slugs | What It Can Do |
|---------|-----------|---------------|
| **Gemini** | `GEMINI_GENERATE_VIDEOS` `GEMINI_WAIT_FOR_VIDEO` `GEMINI_GENERATE_IMAGE` | Text-to-video (Veo), image generation, async polling |
| **Shotstack** | `SHOTSTACK_RENDER_VIDEO` `SHOTSTACK_TRANSFER_ASSET` `SHOTSTACK_GET_RENDER_STATUS` | Programmatic video assembly, cuts, captions, overlays |
| **HeyGen** | `HEYGEN_V2_VIDEO_GENERATE` `HEYGEN_V2_TEMPLATES` `HEYGEN_RETRIEVE_VIDEO_STATUS_DETAILS` | AI avatar videos, template-based generation |
| **ImgBB** | `IMGBB_UPLOAD_IMAGE` | Image hosting (used as intermediary for Instagram flows) |

---

## Full Automation Pipeline (Social Media + AI Video)

The complete pipeline — Claude directs, Gemini generates, Shotstack assembles, Composio distributes:

```
Claude (script + scene prompts)
    ↓
GEMINI_GENERATE_VIDEOS (generate clips)
    ↓
GEMINI_WAIT_FOR_VIDEO (poll → get download URL)
    ↓
SHOTSTACK_RENDER_VIDEO (assemble clips, add captions)
    ↓
Post everywhere via Composio:
  TIKTOK_UPLOAD_VIDEO + TIKTOK_PUBLISH_VIDEO
  YOUTUBE_UPLOAD_VIDEO
  INSTAGRAM_POST_IG_USER_MEDIA + INSTAGRAM_POST_IG_USER_MEDIA_PUBLISH
  TWITTER_CREATION_OF_A_POST
  LINKEDIN_CREATE_LINKED_IN_POST
```

### One-Command AI Social Post (Claude → LinkedIn)

```bash
composio run '
  const post = await experimental_subAgent(
    "Write a LinkedIn post about AI productivity. Under 200 characters. No hashtags.",
    { schema: z.object({ text: z.string() }) }
  );
  await execute("LINKEDIN_CREATE_LINKED_IN_POST", { text: post.structuredOutput.text });
  console.log("Posted:", post.structuredOutput.text);
'
```

### Cross-Post to All Social Platforms (parallel)

```bash
composio execute --parallel \
  TWITTER_CREATION_OF_A_POST -d '{"text": "Your post content here"}' \
  LINKEDIN_CREATE_LINKED_IN_POST -d '{"text": "Your post content here"}'
```

---

## Quick Reference — Link Commands

```bash
# Social
composio link twitter
composio link instagram     # already connected
composio link facebook      # already connected
composio link linkedin
composio link tiktok
composio link youtube       # already connected
composio link discord

# Email & Calendar
composio link gmail
composio link outlook
composio link googlecalendar
composio link zoom

# Productivity
composio link notion
composio link airtable
composio link slack

# Development
composio link github        # already connected

# CRM & Marketing
composio link hubspot
composio link zoho

# E-commerce & Finance
composio link shopify
composio link stripe
composio link quickbooks

# AI & Video
composio link gemini        # already connected
composio link shotstack
composio link heygen
```

---

## Discover More

```bash
# Search for any tool by task description
composio search "send a slack message"
composio search "create a notion page"
composio search "upload to youtube"

# Check schema before executing
composio execute <SLUG> --get-schema

# Preview without running
composio execute <SLUG> --dry-run -d '{"param": "value"}'

# List all tools for a connected app
composio tools list twitter --limit 200
```

---

## Related

- [[SKILL]] — Deep reference for Todoist (84 tools) and Strava (34 tools) with full param tables
- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
- [[00-24 SO (Sovereign Operations)/07 Automations (Workflows & Triggers)|07 Automations]]
