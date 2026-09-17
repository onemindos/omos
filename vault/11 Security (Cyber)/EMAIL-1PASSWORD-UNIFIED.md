---
id: so-11-email-1p-unified
title: "Email + 1Password — Unified System"
type: reference
status: active
created: 2026-04-01
modified: 2026-04-01
owner: zeus
tags: [type/reference, status/active, system/email, system/1password, domain/so]
ai_summary: "Single source of truth combining email identity system (Zoho Mail + ProtonMail) with 1Password vault structure — shows how both systems align to the SO/HP/LE/GE framework, documents all active identities, gaps, and cleanup actions."
---

# Email + 1Password — Unified System

> Single source of truth for both systems. Both are intentionally designed to mirror the OneMind SO/HP/LE/GE framework.
> When they drift apart, this document surfaces the gap.

**Platforms:** Zoho One (Zoho Mail) · ProtonMail · 1Password
**Last Audited:** April 1, 2026 — LIVE pull from `op` CLI (263 items)
**Status:** Email infrastructure production-ready · 1P 78% tagged (58 untagged) · Tags still use old `path:ui` (not `path:so`) · Vault names still say "UI / GE Entrepreneurship"

---

## Architecture Overview

The two systems share the same 4-domain spine:

```
┌────────────────────────────────────────────────────────────────────┐
│                  EMAIL SYSTEM          1PASSWORD                   │
│                                                                    │
│  SECURITY ROOT → zeus.vaults@pm.me  →  1P Account Login (pm.me)   │
│                  legacywarriorx@proton.me  →  Personal secure      │
│                                                                    │
│  SO (00-24)  →  ui@zeusdelacruz.me / zeus.ops@legacyhq.ai         │
│                 zeus@zeusdelacruz.com                              │
│              →  1P: OneMind - Sovereign Operations  (89 items)     │
│                                                                    │
│  HP (25-49)  →  hp@zeusdelacruz.me                                │
│                 shopping@ · finance@ · performance@ · travel@      │
│                 education@ · recreation@                           │
│              →  1P: OneMind - Holistic Performance   (104 items)   │
│                                                                    │
│  LE (50-74)  →  le@zeusdelacruz.me                                │
│                 legal@ · legacy@                                   │
│              →  1P: OneMind - Legacy Evolution       (13 items)    │
│                                                                    │
│  GE (75-99)  →  ge@zeusdelacruz.me / partnernetwork@legacyhq.ai  │
│                 finance@legacyhq.ai                                │
│              →  1P: OneMind - Generational Endeavor  (47 items)    │
└────────────────────────────────────────────────────────────────────┘
```

---

## Layer 1 — ProtonMail (Root of Trust)

ProtonMail sits ABOVE Zoho Mail in the security hierarchy. These addresses do NOT receive newsletters or service emails — they unlock everything else.

| Address | Role | What It Protects |
|---------|------|-----------------|
| `zeus.vaults@pm.me` | **Master** — 1P account login | 1Password, Apple ID, GitHub, Discord, OpenAI, Cloudflare, Raycast, X AI |
| `legacywarriorx@proton.me` | Secondary secure identity | Currently: Zeus Proton in 1P (SO vault) |

**Rules:**
- `zeus.vaults@pm.me` is the email you use to LOG INTO 1Password — not for general signups
- ProtonMail accounts are stored in the SO vault in 1P
- Never register new services with Proton addresses unless they require the highest security tier
- ProtonMail is end-to-end encrypted — use it for: 1P account, Apple ID, core identity anchors

---

## Layer 2 — Email Identity System (Zoho Mail)

**Platform:** Zoho One → Zoho Mail
**Primary mailbox:** zeus@zeusdelacruz.com

### Domain Hierarchy

| Domain | Type | Purpose | Status |
|--------|------|---------|--------|
| `zeusdelacruz.com` | Public / Brand | Personal brand, manual triage, public contact | ✅ Active |
| `zeusdelacruz.me` | Private / Pillar | Structured alias system — one per OneMind category | ✅ Active |
| `legacyhq.ai` | Business | AI ops, business ops, partner network, finance | ✅ Active |
| `onemindcodex.com` | Business | OneMind community (Circle.so), 2 accounts still here | ⚠️ Review — may retire |
| `legacyrestore.co` | Client | Client-specific communications | ✅ Active |
| `weareonemind.com` | Redirect | Brand protection → 301 to onemindcodex.com | ✅ |

### Alias Registry — zeusdelacruz.me

These route into Zoho Mail via Filters #1–5. **Pillar routing is the primary layer.**

| Alias | OneMind Domain | Zoho Folder | 1P Tag | # 1P Items | Used For |
|-------|---------------|-------------|--------|-----------|---------|
| `ui@zeusdelacruz.me` | SO (00-24) | 00-24 SO | `path:so` ⚠️ | 20 | Tech / infra tools |
| `hp@zeusdelacruz.me` | HP (25-49) | 25-49 HP | `path:hp` | — | Health, learning, lifestyle |
| `le@zeusdelacruz.me` | LE (50-74) | 50-74 LE | `path:le` | 1 | Everplans |
| `ge@zeusdelacruz.me` | GE (75-99) | 75-99 GE | `path:ge` | — | Business (limited) |
| `shopping@zeusdelacruz.me` | HP | 25-49 HP | `path:hp` + `#hp-shopping` | 16+2 | Retail & consumer |
| `finance@zeusdelacruz.me` | HP | 25-49 HP | `path:hp` + `#30-finance` | 18+1 | Personal finance, credit bureaus |
| `performance@zeusdelacruz.me` | HP | 25-49 HP | `path:hp` + `#27-body` | 6 | Fitness, wellness |
| `travel@zeusdelacruz.me` | HP | 25-49 HP | `path:hp` + `#travel` | 5 | Flights, transport |
| `education@zeusdelacruz.me` | HP | 25-49 HP | `path:hp` + `#26-mind` | 5 | Courses, content |
| `social@zeusdelacruz.me` | GE | 75-99 GE | `path:ge` | 4 | LinkedIn, TikTok, FB, Pinterest |
| `recreation@zeusdelacruz.me` | HP | 25-49 HP | `path:hp` | 2 | Gaming |
| `legal@zeusdelacruz.me` | LE | 50-74 LE | `path:le` | 4 | DMV, legal, government |
| `legacy@zeusdelacruz.me` | LE | 50-74 LE | `path:le` | 1 | Trustworthy |

> ⚠️ **`ui@` alias issue:** The underlying filter (#3) and 1P tag still say `#ui`. After the SO rename, this alias routes correctly but the label name and tag need updating. The alias email itself (`ui@`) can stay as-is for backward compat — just update the Zoho label and 1P tags.

### Account Registry — legacyhq.ai

| Address | Purpose | # 1P Items | 1P Vault |
|---------|---------|-----------|---------|
| `zeus.ops@legacyhq.ai` | AI tools, business ops | 19 | SO + GE (split — see cleanup) |
| `partnernetwork@legacyhq.ai` | Affiliate accounts | 7 | GE |
| `finance@legacyhq.ai` | Amazon Business, Stripe, Verizon | 6 | GE |
| `legal@legacyhq.ai` | Wyoming LLC Attorney Portal | 1 | LE |

### Accounts to Migrate OFF Legacy Email Addresses

These external accounts use non-system emails and should be migrated during next cleanup window:

| Current Identity | # Accounts | Migrate To | Notes |
|-----------------|-----------|-----------|-------|
| `executivebycredit@gmail.com` | 6 | `education@zeusdelacruz.me` or `ge@` | Old credit course sites |
| `legacyplatformholdings@gmail.com` | 4 | `finance@zeusdelacruz.me` (Hushed, Credit Strong), keep LPH Gmail as login | LPH entity accounts |
| `zeus@onemindcodex.com` | 2 | `zeus.ops@legacyhq.ai` | If onemindcodex.com domain retires |
| `Founder@legacyplatform.co` | 1 | TBD | If legacyplatform.co retires |
| `delacruzonline1@gmail.com` | 1 | `finance@zeusdelacruz.me` | creditsecret.org |

---

## Layer 3 — 1Password Vault System

**Account:** zeus.vaults@pm.me @ my.1password.com
**Live pull:** April 1, 2026 · **Total:** 263 items

> ⚠️ **Vault names have NOT been updated.** Still show old naming (UI, GE Entrepreneurship). The January 2026 restructuring doc described target names — they were never applied.

### Vault Registry (Live Names)

| Current Vault Name | ID | Target Name | Items | Purpose |
|-------------------|----|-------------|------:|--------|
| **00-24 UI (Unified Intelligence)** | `aaeyris5z44pdqspk277u2oppi` | 00-24 SO (Sovereign Operations) | ~120 | AI/tech/infra logins + API keys |
| **25-49 HP (Holistic Performance)** | `azs7ujqwshhbikjwlx6qwelj4q` | *(keep as-is)* | ~104 | Personal life — finance, health, shopping, travel |
| **50-74 LE (Legacy Evolution)** | `sjnxwfebtg7aym5cw4egcz2gv4` | *(keep as-is)* | ~13 | Legal, family, estate |
| **75-99 GE (Generational Entrepreneurship)** | `vf7fxndqwwthdqivvx5q6bwr3e` | 75-99 GE (Generational Endeavor) | ~47 | Business, ventures, affiliate, revenue |
| **Legacy AI** | `4allf4osg3kfequ22fb3zhuwnm` | *(keep or populate)* | 0 | Still empty — AI API keys could live here |
| **Personal** | `vos3svqijr7xouspkgu7qgqlau` | *(keep as-is)* | 3 | Off-system items (identity docs, API key) |
| **_Delete** | `vcwyk7vivme2dvtlh4xfi6nia4` | *(purge)* | 2 | Envato, Trello — ready to delete |

### Tag System (2-Tag Model)

Every item should have BOTH tags:

```
path:[domain]  +  category:[function]

Examples:
  path:so  +  category:infrastructure   → VPS, DNS, server
  path:so  +  category:intelligence     → OpenAI, Claude, Perplexity
  path:hp  +  category:finance          → Chime, Venmo, Coinbase
  path:ge  +  category:partnernetwork   → GoHighLevel, Partnerstack
  path:le  +  category:legal            → DMV, SSN, Notary
```

> ⚠️ **Tag drift**: Tag `path:ui` was used before the SO rename. Any items tagged `path:ui` should be updated to `path:so`. Run: `op item list --tags "path:ui" --format json` to find them.

### What Lives in Each Vault

**SO Vault** (89 items) — Tech & Operations
- AI agents, APIs (Anthropic, OpenAI, X AI, ElevenLabs, etc.)
- Infrastructure (Digital Ocean, Cloudflare, Railway, Vercel, Coolify)
- Dev tools (GitHub, Postman, Apify, Supabase, Twilio)
- Core identity anchors (1P account, Apple ID, Discord, Obsidian)
- Agent platforms (Habitica, TickTick, Todoist, Clickup)
- *Should also contain*: ProtonMail accounts, Letta, Home Assistant

**HP Vault** (104 items) — Personal Life
- Finance & credit (Chime, Bank of America, Wells Fargo, Coinbase, Kraken, PayPal)
- Credit bureaus (Equifax, Experian, TransUnion, Chex Systems, etc.)
- Shopping (Amazon, eBay, Target, Starbucks, etc.)
- Health (Equinox, AllTrails, Oura Ring, PlushCare, etc.)
- Travel (Airbnb, JetBlue, Lyft, Uber, Amtrak)
- Education (Udemy, Skool, Substack, Fund Launch)
- Government (SSA, IRS, DMV NYC — note: DMV is LE-ish but stored here)

**LE Vault** (13 items) — Foundation & Legacy
- Legal (Improv Insurance, Notary Live, NYC DMV, SSN.Gov)
- Estate planning (Everplans, Trustworthy)
- Insurance, legal portals

**GE Vault** (47 items) — Ventures & Revenue
- Business ops (Stripe, Whop, Amazon Business, Microsoft)
- Content/Branding (Descript, Heygen, Sedo, VistaPrint, Webflow)
- Partners/Affiliates (GoHighLevel x2, Partnerstack, IDIQ, Swanwick)
- Social brand (Instagram, Social Media Codex accounts)

---

## Unified Identity Reference Card

| Identity | Type | Used For | 1P Vault | Count |
|----------|------|---------|---------|-------|
| `zeus.vaults@pm.me` | Proton (master) | 1P login · Apple ID · GitHub · Discord · OpenAI · Cloudflare · X AI · Raycast | SO | 13 |
| `legacywarriorx@proton.me` | Proton (secondary) | Zeus Proton | SO | 1 |
| `zeus@zeusdelacruz.com` | Public brand | Circle, Zoho One, Mercury, Webflow, Payoneer, AT&T, EventBrite | SO + GE | 8 |
| `zeus.ops@legacyhq.ai` | Business ops | AI tools (Anthropic, Ollama) + biz (AWS, Microsoft, Sedo, Fiverr) | SO + GE | 19 |
| `partnernetwork@legacyhq.ai` | Affiliate | GoHighLevel affiliates, Partnerstack, IDIQ, Joinsecret | GE | 7 |
| `finance@legacyhq.ai` | Business finance | Amazon Business, Stripe, Verizon, Whop, Founder Pass | GE | 6 |
| `legal@legacyhq.ai` | Business legal | Wyoming LLC Attorney Portal | LE | 1 |
| `ui@zeusdelacruz.me` | SO pillar alias | ATOC, Apify, Bouncie, DigitalOcean, ElevenLabs, Miro, etc. | SO | 20 |
| `finance@zeusdelacruz.me` | HP personal finance | Banks, credit bureaus, crypto | HP | 18+1 |
| `shopping@zeusdelacruz.me` | HP shopping | Amazon, eBay, Target, Starbucks, etc. | HP | 16+2 |
| `hp@zeusdelacruz.me` | HP general | Misc HP accounts (Skool, Udemy, Robinhood, Oura) | HP | 14 |
| `performance@zeusdelacruz.me` | HP fitness | Equinox, AllTrails, FloatSpa, PlushCare | HP | 6 |
| `travel@zeusdelacruz.me` | HP travel | JetBlue, Airbnb, Lyft, Amtrak, Uber | HP | 5 |
| `education@zeusdelacruz.me` | HP learning | Substack, Patreon, Fund Launch | HP | 5 |
| `social@zeusdelacruz.me` | GE social brand | LinkedIn, TikTok, Facebook, Pinterest | GE | 4 |
| `legal@zeusdelacruz.me` | LE personal legal | NYC DMV, SSN.Gov, Notary Live | LE | 4 |
| `legacy@zeusdelacruz.me` | LE | Trustworthy | LE | 1 |
| `le@zeusdelacruz.me` | LE | Everplans | LE | 1 |
| `recreation@zeusdelacruz.me` | HP gaming | Instant Gaming, Ubisoft | HP | 2 |

---

## What ProtonMail Does That Zoho Can't

| Scenario | Use ProtonMail | Use Zoho Mail |
|----------|---------------|--------------|
| Account recovery email for 1Password | zeus.vaults@pm.me | — |
| Apple ID / major identity anchors | zeus.vaults@pm.me | — |
| Signing up for new services | — | Use appropriate alias |
| Business communications | — | zeus.ops@legacyhq.ai |
| Receiving newsletter / service emails | — | Domain alias |
| End-to-end encrypted correspondence | legacywarriorx@proton.me | — |
| Webhook / automation notifications | — | system@legacyhq.ai (future) |

---

## Cleanup Action Items

### 🔴 Critical (System Integrity)

- [ ] **`me@zeusdelacruz.com` — 18 accounts using an unregistered email** — This email is NOT in the Zoho alias registry. It appears to be an old personal address used before the structured alias system was built. All 18 accounts need migration:
  - HP accounts → `hp@zeusdelacruz.me`: Monarch Money, MyFico, Udemy, Skool, Substack, Deep Learning AI, Strava, Schedule35, Airsoft GI, Pirateship.com, Walmart, growyourownvegetables.org, The Airbnbae Academy, Use Origin, ouraring.com, Robinhood (Anthony Delacruz)
  - LE accounts → `le@zeusdelacruz.me`: Bouncie, Everplans
  - OR: Register `me@zeusdelacruz.com` as an alias in Zoho Mail (simpler short-term fix)

- [ ] **`path:ui` tags are stale — 42 items tagged, 0 tagged `path:so`** — The SO rename was never applied to tags. Batch update all `path:ui` → `path:so`. Command: `op item list --tags "path:ui" --format json | jq -r '.[].id'` then update each.

- [ ] **Vault names never updated** — Still say "00-24 UI (Unified Intelligence)" and "75-99 GE (Generational Entrepreneurship)". Rename via 1Password Settings → Vaults to match framework.

- [ ] **Zoho filter labels still say "#ui"** — After PR rename, update filter #3 in Zoho Mail from "3 - Route UI" to "3 - Route SO"; rename the `#ui` label to `#so`.

- [ ] **Fix capitalization in 1P** — `Shopping@zeusdelacruz.me` → `shopping@` (CVS, Offer Up) and `Finance@zeusdelacruz.me` → `finance@` (FTC - IDT). They fail to match Zoho filters with capital letters.

- [ ] **Data corruption** — `OneMind Home (Verizon)` has `category:infrastructure,path:ui` stored as its USERNAME field. The tags string was accidentally saved there. Fix: edit the item, clear the username field, set proper tags.

### 🟡 Medium (Organization)

- [ ] **58 untagged items** (up from 14 in January — lots of drift from new API keys added). Tag backlog concentrated in the UI vault. Priority untagged items:
  - `Anthropic`, `OpenAI`, `ElevenLabs`, `GitHub`, `Discord`, `Cloudflare - OneMind OS`, `Azure OpenAI`, `Synadia Cloud` — all need `path:ui` + appropriate category
  - API credential items (AWS Bedrock, Cartesia, Deepgram, LiveKit, OmniDB, etc.) — tag as `path:ui` + `category:infrastructure` or `category:intelligence`

- [ ] **Purge `_Delete` vault** — Only 2 items remain: Envato and Trello (`zeus.ops@legacyhq.ai`). Both are confirmed stale. Delete items, then remove vault.

- [ ] **Populate or remove Legacy AI vault** — Still empty. Either move AI API keys here (Anthropic, OpenAI, ElevenLabs, Deepgram, Azure OpenAI, Cartesia) OR remove the vault. Keeping an empty vault is noise.

- [ ] **Split `zeus.ops@legacyhq.ai` between vaults** — Currently spans both UI and GE. AI tools already in UI vault (Antheropic, Habitica, Lang, Ollama, Perplexity, Runner H, TickTick) are correctly placed. Business tools (Alibaba, AWS, Descript, Fiverr, Heygen, Microsoft, Sedo, VistaPrint, infinity, quso.ai) are in GE vault — that's also correct. No vault moves needed, but the same email for both domains is worth noting.

- [ ] **Move Bouncie** — Currently in LE vault (tagged `category:family,path:le`). That's actually correct for a vehicle tracker used for family/home purposes. ✅ Leave it.

- [ ] **Add missing usernames to API credential items** — Items like `Airbyte`, `TestSprite`, `Composio`, `Context7`, `Docker`, `Raindrop`, `Tavily`, `Traintime`, `TSA PreCheck` have no username. Add the login email for each.

### 🟢 Low (Long-term Migration)

- [ ] **Migrate `executivebycredit@gmail.com`** (6 accounts: geru.com, milliondollarads.com, secretstrilogy.com, socialmediacodex.com, thegrownetwork.com, videoschallenge.com) → `education@zeusdelacruz.me`

- [ ] **Migrate `legacyplatformholdings@gmail.com`** (4 accounts) — Hushed + Credit Strong → `finance@zeusdelacruz.me`; LPH Gmail becomes login-only

- [ ] **Review `zeus@onemindcodex.com`** — Now has 3 items: AWS Acadacom, Datadog, Synadia Cloud. If onemindcodex.com domain retires, migrate to `zeus.ops@legacyhq.ai`

- [ ] **Remove duplicate** — `creditlaunchsecrets.com` stored twice under `Ajhailey@theclassof700.com` — delete one

- [ ] **Consolidate eBay** — `EBay` and `Ebay - Zeus` (same `shopping@` address) → merge into one entry

- [ ] **Review external credentials** — Items under `ahailey`, `Jeremylewis357@gmail.com`, `Josh@getmormedia.com`, `eflores1524@gmail.com`, `Luzmdelacruz@icloud.com` — audit and clean or document why they're here

### 🔵 Future Enhancements

- [ ] **Add `me@zeusdelacruz.com` as a Zoho alias** pointing to the catch-all folder, while migrating accounts to proper aliases
- [ ] **Add `so@zeusdelacruz.me` alias** — redirect to SO folder; use `so@` for NEW signups; leave `ui@` for existing accounts (backward compat)
- [ ] **Add `system@legacyhq.ai`** as dedicated webhook/automation inbox
- [ ] **Letta agent integration** — query 1P by tag via `op` CLI for "what are my finance credentials?" lookups
- [ ] **N8N auto-tagging workflow** — new 1P items automatically get path/category tags

---

## Quick Commands

```bash
# Find all items still tagged path:ui (need renaming to path:so)
op item list --tags "path:ui" --format json | jq '.[].title'

# Count items per vault
op vault list --format json | jq '.[].name'

# List Z Delete vault
op item list --vault "Z Delete" --format json | jq '.[].title'

# Find items with no username
op item list --format json | jq '.[] | select(.username == null or .username == "") | .title'

# Find untagged items
op item list --format json | jq '.[] | select(.tags == null or .tags == []) | .title'

# Rename path:ui tag to path:so on an item
op item edit "ITEM_TITLE_OR_ID" --tags "path:so,category:EXISTING_CATEGORY"
```

---

## Reference Links

- Full 1P restructuring doc: [1PASSWORD-VAULT-RESTRUCTURING.md](./Credentials/1PASSWORD-VAULT-RESTRUCTURING.md)
- 1P cleanup audit: [1PASSWORD-CLEANUP.md](./1PASSWORD-CLEANUP.md)
- Email system master: [EMAIL-SYSTEM-MASTER.md](../01 Infrastructure (Grid & APIs)/Email/Zoho-MAIL-Config/EMAIL-SYSTEM-MASTER.md)
- Email filter specs: [EMAIL-SMART-FILTERS-DUAL-TAG.md](../01 Infrastructure (Grid & APIs)/Email/Zoho-MAIL-Config/EMAIL-SMART-FILTERS-DUAL-TAG.md)
- Secrets management: [SECRETS-MANAGEMENT.md](./Accounts-Recovery/SECRETS-MANAGEMENT.md)
