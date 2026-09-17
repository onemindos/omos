---
name: rss-digest
description: Fetch and digest RSS feeds from Substack newsletters and other sources.
homepage: https://substack.com
metadata: { "clawdbot": { "emoji": "📰" } }
---

# RSS Digest

Fetch, manage, and digest RSS feeds — primarily Substack newsletters. Replaces email newsletter delivery with on-demand or scheduled digests.

## Tools

- `fetch_rss_digest(feed_name?, hours?, max_per_feed?)` — Fetch recent articles from all or one feed. Returns title, date, summary, and link grouped by category. Default: last 24 hours, 5 per feed.
- `list_rss_feeds()` — List all configured feeds with names, URLs, and categories.
- `add_rss_feed(name, url?, category?)` — Add a new feed. For Substack, just provide the name — URL auto-generates. Validates the feed is reachable before adding.
- `remove_rss_feed(name)` — Remove a feed by name.

## Pre-loaded Feeds (14)

### Homestead & Gardening

- `houseofgreen` — House of Green
- `holisticfarming` — Holistic Farming
- `charlesdowding` — Charles Dowding
- `acornlandlabs` — Acorn Land Labs
- `peasantpolis` — Peasant Polis

### Parenting

- `dralizapressman` — Dr. Aliza Pressman
- `theparentingguide` — The Parenting Guide
- `screenlessdad` — Screenless Dad

### Health

- `dremilyleeming` — Dr. Emily Leeming
- `vajenda` — The Vajenda (Dr. Jen Gunter)
- `joshturknett` — Josh Turknett (Brainjo)

### Philosophy & Mindset

- `stoicwisdoms` — Stoic Wisdoms
- `johannadebiase` — Johanna de Biase
- `yungpueblo` — Yung Pueblo

## Examples

- "Give me today's RSS digest" → `fetch_rss_digest()`
- "What did House of Green post this week?" → `fetch_rss_digest(feed_name: "houseofgreen", hours: 168)`
- "Show me all my feeds" → `list_rss_feeds()`
- "Add Tim Ferriss Substack to my feeds" → `add_rss_feed(name: "timferriss", category: "productivity")`
- "Remove the vajenda feed" → `remove_rss_feed(name: "vajenda")`

## Usage Patterns

### Daily Morning Briefing

Call `fetch_rss_digest()` with default 24h window. Summarize the top articles by category for Zeus.

### Weekly Catch-up

Call `fetch_rss_digest(hours: 168)` for a full week's content. Good for weekend reading lists.

### Category Focus

Fetch a single feed to deep-dive: `fetch_rss_digest(feed_name: "charlesdowding")`.

## Notes

- Feeds are persisted in `feeds.json` in the skill directory. Adding/removing feeds is permanent.
- Substack RSS is free and public — no API key needed.
- Non-Substack RSS/Atom feeds work too — just provide the full URL.
- Articles are sorted newest-first within each feed.
- Content snippets are plain text (HTML stripped) and truncated to 500 chars.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
