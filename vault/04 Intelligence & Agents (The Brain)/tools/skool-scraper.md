---
title: "Skool Scraper — Self-Hosted"
type: project
domain: SO
status: active
created: 2026-05-26
updated: 2026-05-26
tags: [scraping, skool, crawlee, python, apify, self-hosted, leads, community]
---

# Skool Scraper — Self-Hosted

Replaces: https://apify.com/louisdeconinck/skool-scraper ($19/month)
Built with Crawlee for Python. Runs locally on Zeuss-MBP. No Apify account needed.

**Location:** ~/apps/skool-scraper/scraper.py
**Venv:** ~/apps/skool-scraper/.venv (Python 3.13)
**Output:** ~/apps/skool-scraper/output/ (JSON + CSV)

## Activate

```bash
cd ~/apps/skool-scraper
source .venv/bin/activate
```

## Modes & Commands

| Mode | What it does | Cookies? |
|------|-------------|----------|
| discovery | Search Skool groups by keyword | No |
| groups | Scrape a group's info page | No (public) / Yes (private) |
| members | Scrape member list of a group | Yes |

```bash
# Discovery — no login needed
python scraper.py --mode discovery --query "fitness" --max 50

# Single group
python scraper.py --mode groups --url "https://www.skool.com/bbc/about"

# Members (needs cookies)
python scraper.py --mode members --url "https://www.skool.com/bbc/about" --cookies "PASTE_COOKIE_STRING"

# Private group (needs cookies)
python scraper.py --mode groups --url "https://www.skool.com/SLUG/about" --cookies "PASTE_COOKIE_STRING"
```

## Getting Cookies from Chrome

1. Log into skool.com in Chrome
2. DevTools (F12) → Network tab → refresh page
3. Click any request → Request Headers → copy `Cookie:` value
4. Paste as `--cookies "..."` argument
5. Cookies expire — re-grab when auth fails

## What Data It Returns

### Discovery / Groups
- Group name, description
- Member count
- Price (Free / $X/month)
- Group URL

### Members (with cookies)
- Full name
- Role / badge (Admin, Member, etc.)
- Profile URL
- Join date

### Member Profiles (deep scrape — future)
- Bio text
- Social links (Twitter, Instagram, LinkedIn, website)
- Post count
- **No emails** — Skool does not expose emails on profiles

## Email Enrichment Pipeline (future)

Skool has no email field. To get contact info:
```
Skool members scrape
  → names + social handles + websites
    → Apollo / Hunter / Clay enrichment
      → verified emails (~30-40% hit rate)
```

## Build vs Apify

| | Apify Cloud | Self-Hosted |
|--|------------|-------------|
| Cost | $19/month + usage | Free |
| Speed | Fast (managed infra) | Fast (local M2 Max) |
| Privacy | Data goes to Apify | Stays local |
| Control | Limited | Full |
| Maintenance | Automatic | Manual (selectors may break on Skool deploys) |

## Selector Notes

Skool uses styled-components with hashed CSS classes that can change on deploys.
Scraper uses partial match patterns that survive rotations:
- Discovery cards: `[class*="DiscoveryCardLink"]`
- Group name: `[class*="GroupName"]`
- Meta (members + price): `[class*="DiscoveryCardMeta"]`

## Related

- Hermes skill: `crawlee-web-scraping`
- Research: [[ai-pointer-landscape]] (same session — Crawlee context)
- Apify original: https://apify.com/louisdeconinck/skool-scraper
- Crawlee docs: https://crawlee.dev/python/
- Repo: https://github.com/apify/crawlee-python

## Tasks

- [ ] Test members scraper with real Skool cookies
- [ ] Build member profile deep-scraper (bio, socials, website)
- [ ] Add Apollo/Clay enrichment pipeline for emails
