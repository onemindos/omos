---
id: so-11-email-migration
title: "Email Domain Migration — zeusdelacruz.me + legacyhq.ai → New System"
type: sop
status: active
created: 2026-04-01
modified: 2026-04-01
owner: zeus
tags: [type/sop, status/active, system/email, system/1password]
ai_summary: "Process-of-elimination migration checklist to retire zeusdelacruz.me aliases and human-facing legacyhq.ai addresses. All 263 1P accounts mapped to new target address. Source of truth: BUSINESS-ARCHITECTURE.md (2026-03-21)."
---

# Email Domain Migration

**Goal:** Retire `zeusdelacruz.me` (all aliases) and human-facing `@legacyhq.ai` accounts.
Replace with the clean 2-domain system established in BUSINESS-ARCHITECTURE.md.

**Source of truth:** [BUSINESS-ARCHITECTURE.md](../../75-99 GE (Generational Endeavor)/75 Zeus Delacruz Enterprise (Personal Brand)/BUSINESS-ARCHITECTURE.md)
**Live data:** Pulled from `op` CLI April 1, 2026 — 263 items across 7 vaults.

---

## The New Email System (Target State)

```
zeusdelacruz.com              onemindcodex.com           legacyhq.ai (Legacy AI only)
─────────────────             ─────────────────          ─────────────────────────────
zeus@    Personal brand       zeus@  Business/admin      command@  Legacy AI inbox
me@      All personal life    support@ Customer support  system@   Webhooks/automation
speaking@ Speaking requests   team@  Internal VAs
consulting@ Consulting        partners@ Affiliates
press@   Media/PR

Proton (root of trust — never for regular signups)
──────────────────────────────────────────────────
zeus.vaults@pm.me    1P login, Apple ID, Cloudflare, GitHub, banks
legacywarriorx@proton.me  Secondary secure identity
```

**Key rule:** `me@zeusdelacruz.com` = everything personal. `zeus@zeusdelacruz.com` = everything brand/business.
`@legacyhq.ai` is NOW Legacy AI's private domain only — Zeus does not use it for account signups.

---

## Domain Retirement Plan

| Domain | Status | Action | Timing |
|--------|--------|--------|--------|
| `zeusdelacruz.me` | Keep alive | Retire AFTER all accounts migrated | After migration complete |
| `zeus.ops@legacyhq.ai` | Retire | Migrate all 19 accounts to zeusdelacruz.com | Per checklist below |
| `finance@legacyhq.ai` | Retire | Migrate 6 accounts | Per checklist below |
| `partnernetwork@legacyhq.ai` | Retire | Migrate 7 accounts to partners@onemindcodex.com | Per checklist below |
| `legal@legacyhq.ai` | Retire | Migrate 1 account | Per checklist below |
| `command@legacyhq.ai` | KEEP | Legacy AI's inbox — do not touch | Permanent |
| `system@legacyhq.ai` | KEEP | Webhooks/automation | Permanent |

---

## Migration Checklists

Mark `[x]` when the account email has been updated AND 1Password entry is updated.

---

### BLOCK 1 — `ui@zeusdelacruz.me` → 19 accounts

**Target:** All 19 move to `zeus@zeusdelacruz.com` (they're tech/infra tools — brand identity, not personal life)

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | Apify | SO | zeus@zeusdelacruz.com | [ ] |
| 2 | ATOC | SO | zeus@zeusdelacruz.com | [ ] |
| 3 | Clickup | SO | zeus@zeusdelacruz.com | [ ] |
| 4 | Coolify | SO | zeus@zeusdelacruz.com | [ ] |
| 5 | Digital Ocean | SO | zeus@zeusdelacruz.com | [ ] |
| 6 | Elestio | SO | zeus@zeusdelacruz.com | [ ] |
| 7 | Eleven Labs | SO | zeus@zeusdelacruz.com | [ ] |
| 8 | Miro | SO | zeus@zeusdelacruz.com | [ ] |
| 9 | Nabu Casa | SO | me@zeusdelacruz.com | [ ] |
| 10 | NVIDA (NVIDIA) | SO | zeus@zeusdelacruz.com | [ ] |
| 11 | Open Router | SO | zeus@zeusdelacruz.com | [ ] |
| 12 | Postman | SO | zeus@zeusdelacruz.com | [ ] |
| 13 | Railway | SO | zeus@zeusdelacruz.com | [ ] |
| 14 | Raindrop.io | SO | me@zeusdelacruz.com | [ ] |
| 15 | Setapp | SO | me@zeusdelacruz.com | [ ] | ← personal Mac subscription |
| 16 | Supa Base (Supabase) | SO | zeus@zeusdelacruz.com | [ ] |
| 17 | Todoist | SO | me@zeusdelacruz.com | [ ] |
| 18 | Twilio - Zeus | SO | zeus@zeusdelacruz.com | [ ] |
| 19 | Vercel | SO | zeus@zeusdelacruz.com | [ ] |

**Block 1 progress: 0 / 19**

---

### BLOCK 2 — `finance@zeusdelacruz.me` → 18 accounts

**Target:** All personal finance → `me@zeusdelacruz.com`. Crypto that could be business → `zeus@zeusdelacruz.com`.

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | Chex Systems | HP | me@zeusdelacruz.com | [ ] |
| 2 | Chime | HP | me@zeusdelacruz.com | [ ] |
| 3 | Consumerfinance | HP | me@zeusdelacruz.com | [ ] |
| 4 | Credit Karma | GE | me@zeusdelacruz.com | [ ] |
| 5 | Credit Sesame | GE | me@zeusdelacruz.com | [ ] |
| 6 | Early Warning | HP | me@zeusdelacruz.com | [ ] |
| 7 | Equifax | HP | me@zeusdelacruz.com | [ ] |
| 8 | Eversource - Harold | HP | me@zeusdelacruz.com | [ ] |
| 9 | Factor Trust | HP | me@zeusdelacruz.com | [ ] |
| 10 | Kikoff | HP | me@zeusdelacruz.com | [ ] |
| 11 | kraken | GE | me@zeusdelacruz.com | [ ] |
| 12 | microbilt | HP | me@zeusdelacruz.com | [ ] |
| 13 | NCTUE | HP | me@zeusdelacruz.com | [ ] |
| 14 | Paypal - Personal (Zeus) | GE | me@zeusdelacruz.com | [ ] |
| 15 | Smart Credit | GE | me@zeusdelacruz.com | [ ] |
| 16 | Transunion - Zeus | HP | me@zeusdelacruz.com | [ ] |
| 17 | Venmo (Anthony Delacruz) | GE | me@zeusdelacruz.com | [ ] |
| 18 | coinbase.com | GE | me@zeusdelacruz.com | [ ] |

**Block 2 progress: 0 / 18**

> ⚠️ Also fix: `Finance@zeusdelacruz.me` (capital F) → `FTC - IDT` in HP — move to `me@zeusdelacruz.com`

---

### BLOCK 3 — `shopping@zeusdelacruz.me` → 15 accounts

**Target:** All personal shopping → `me@zeusdelacruz.com`

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | Amazon.com | HP | me@zeusdelacruz.com | [ ] |
| 2 | Cooper Hawks Membership | HP | me@zeusdelacruz.com | [ ] |
| 3 | Curaleaf | HP | me@zeusdelacruz.com | [ ] |
| 4 | EBay | HP | me@zeusdelacruz.com | [ ] |
| 5 | Ebay - Zeus *(duplicate — merge with above)* | HP | me@zeusdelacruz.com | [ ] |
| 6 | Fine Fettle Dispensary | HP | me@zeusdelacruz.com | [ ] |
| 7 | Garden Catering | HP | me@zeusdelacruz.com | [ ] |
| 8 | instacart.com | HP | me@zeusdelacruz.com | [ ] |
| 9 | Mighty Kind | HP | me@zeusdelacruz.com | [ ] |
| 10 | Regmovies | HP | me@zeusdelacruz.com | [ ] |
| 11 | Shakeshack | HP | me@zeusdelacruz.com | [ ] |
| 12 | shopify.com | HP | me@zeusdelacruz.com | [ ] |
| 13 | Staples | HP | me@zeusdelacruz.com | [ ] |
| 14 | starbucks | HP | me@zeusdelacruz.com | [ ] |
| 15 | Target | HP | me@zeusdelacruz.com | [ ] |

**Block 3 progress: 0 / 15**

> ⚠️ Also fix: `Shopping@zeusdelacruz.me` (capital S) → `CVS` + `Offer up` in HP — move to `me@zeusdelacruz.com`

---

### BLOCK 4 — `performance@zeusdelacruz.me` → 6 accounts

**Target:** All fitness/wellness → `me@zeusdelacruz.com`

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | AllTrails | HP | me@zeusdelacruz.com | [ ] |
| 2 | doterra.com | HP | me@zeusdelacruz.com | [ ] |
| 3 | Equinox plus | HP | me@zeusdelacruz.com | [ ] |
| 4 | Letsgetchecked | HP | me@zeusdelacruz.com | [ ] |
| 5 | plushcare.com | HP | me@zeusdelacruz.com | [ ] |
| 6 | Rise Above - Floatation | HP | me@zeusdelacruz.com | [ ] |

**Block 4 progress: 0 / 6**

---

### BLOCK 5 — `travel@zeusdelacruz.me` → 5 accounts

**Target:** All travel → `me@zeusdelacruz.com`

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | Airbnb - (Hotel) | HP | me@zeusdelacruz.com | [ ] |
| 2 | amtrack | HP | me@zeusdelacruz.com | [ ] |
| 3 | Jetblue - (Airline) | HP | me@zeusdelacruz.com | [ ] |
| 4 | Lyft | HP | me@zeusdelacruz.com | [ ] |
| 5 | uber.com | HP | me@zeusdelacruz.com | [ ] |

**Block 5 progress: 0 / 5**

---

### BLOCK 6 — `education@zeusdelacruz.me` → 4 accounts

**Target:** Personal courses/content → `me@zeusdelacruz.com`

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | Fund Launch | HP | me@zeusdelacruz.com | [ ] |
| 2 | pateron | HP | me@zeusdelacruz.com | [ ] |
| 3 | takebackyourpower.net | HP | me@zeusdelacruz.com | [ ] |
| 4 | thesovereignsway.com | HP | me@zeusdelacruz.com | [ ] |

**Block 6 progress: 0 / 4**

---

### BLOCK 7 — `social@zeusdelacruz.me` → 4 accounts

**Target:** Brand social accounts → `zeus@zeusdelacruz.com` (brand identity, not personal life)

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | facebook | GE | zeus@zeusdelacruz.com | [ ] |
| 2 | Linkedin | GE | zeus@zeusdelacruz.com | [ ] |
| 3 | Pinterest | GE | zeus@zeusdelacruz.com | [ ] |
| 4 | Tiktok | GE | zeus@zeusdelacruz.com | [ ] |

**Block 7 progress: 0 / 4**

---

### BLOCK 8 — `legal@zeusdelacruz.me` → 4 accounts

**Target:** Personal legal/government → `me@zeusdelacruz.com`

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | Improv Insurance | LE | me@zeusdelacruz.com | [ ] |
| 2 | Notary Live | LE | me@zeusdelacruz.com | [ ] |
| 3 | NYC DMV | LE | me@zeusdelacruz.com | [ ] |
| 4 | Ssn.Gov | LE | me@zeusdelacruz.com | [ ] |

**Block 8 progress: 0 / 4**

---

### BLOCK 9 — Single `zeusdelacruz.me` aliases → 3 accounts

| # | Account | Current Email | 1P Vault | New Email | Done |
|---|---------|--------------|---------|-----------|------|
| 1 | Trustworthy | legacy@zeusdelacruz.me | LE | me@zeusdelacruz.com | [ ] |
| 2 | Everplans | le@zeusdelacruz.me | LE | me@zeusdelacruz.com | [ ] |
| 3 | Instant Gaming | recreation@zeusdelacruz.me | HP | me@zeusdelacruz.com | [ ] |
| 4 | Ubisoft | recreation@zeusdelacruz.me | HP | me@zeusdelacruz.com | [ ] |

**Block 9 progress: 0 / 4**

---

### BLOCK 10 — `zeus.ops@legacyhq.ai` → 19 accounts

**Target decision logic:**
- AI/tech tools → `zeus@zeusdelacruz.com`
- Productivity/personal tools → `me@zeusdelacruz.com`

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | Alibaba | GE | zeus@zeusdelacruz.com | [ ] |
| 2 | Antheropic (Anthropic) | SO | zeus@zeusdelacruz.com | [ ] |
| 3 | AWS | GE | zeus@zeusdelacruz.com | [ ] |
| 4 | Descript | GE | zeus@zeusdelacruz.com | [ ] |
| 5 | Fiverr - Legacy AI | GE | zeus@zeusdelacruz.com | [ ] |
| 6 | Habitica | SO | me@zeusdelacruz.com | [ ] |
| 7 | Heygen | GE | zeus@zeusdelacruz.com | [ ] |
| 8 | infinity - Project Manager | GE | zeus@zeusdelacruz.com | [ ] |
| 9 | Lang | SO | zeus@zeusdelacruz.com | [ ] |
| 10 | Microsoft | GE | zeus@zeusdelacruz.com | [ ] |
| 11 | Ollama | SO | zeus@zeusdelacruz.com | [ ] |
| 12 | Perplexity | SO | me@zeusdelacruz.com | [ ] |
| 13 | quso.ai | GE | zeus@zeusdelacruz.com | [ ] |
| 14 | Runner H | SO | zeus@zeusdelacruz.com | [ ] |
| 15 | Sedo | GE | zeus@zeusdelacruz.com | [ ] |
| 16 | TickTick | SO | me@zeusdelacruz.com | [ ] |
| 17 | Vista Print | GE | zeus@zeusdelacruz.com | [ ] |
| 18 | Envato | _Delete | *(delete item)* | [ ] |
| 19 | Trello | _Delete | *(delete item)* | [ ] |

**Block 10 progress: 0 / 19**

---

### BLOCK 11 — `finance@legacyhq.ai` → 6 accounts

**Target:** Business finance → `zeus@zeusdelacruz.com`

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | Amazon - Legacy AI | GE | zeus@zeusdelacruz.com | [ ] |
| 2 | Amazon Business | GE | zeus@zeusdelacruz.com | [ ] |
| 3 | Founder Pass | GE | zeus@zeusdelacruz.com | [ ] |
| 4 | Stripe - (Restore) (Acadacom) | GE | zeus@zeusdelacruz.com | [ ] |
| 5 | Verizon Home | GE | me@zeusdelacruz.com | [ ] |
| 6 | Whop | GE | zeus@zeusdelacruz.com | [ ] |

**Block 11 progress: 0 / 6**

---

### BLOCK 12 — `partnernetwork@legacyhq.ai` → 7 accounts

**Target:** `partners@onemindcodex.com` (the new partners alias)

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | GoHighlevel Affiliate | GE | partners@onemindcodex.com | [ ] |
| 2 | GoHighlevel Affiliate 2 | GE | partners@onemindcodex.com | [ ] |
| 3 | IDIQ Affiliate | GE | partners@onemindcodex.com | [ ] |
| 4 | Joinsecret | GE | partners@onemindcodex.com | [ ] |
| 5 | Monarch For Professionals | GE | partners@onemindcodex.com | [ ] |
| 6 | Partnerstack | GE | partners@onemindcodex.com | [ ] |
| 7 | Swanwick Affiliate | GE | partners@onemindcodex.com | [ ] |

**Block 12 progress: 0 / 7**

> ⚠️ First: Create `partners@onemindcodex.com` alias in Zoho Mail before migrating these.

---

### BLOCK 13 — `legal@legacyhq.ai` → 1 account

| # | Account | 1P Vault | New Email | Done |
|---|---------|---------|-----------|------|
| 1 | Wyoming LLC Attorney Portal | LE | zeus@zeusdelacruz.com | [ ] |

**Block 13 progress: 0 / 1**

---

### ACCOUNTS NOT MIGRATING (already on correct new system addresses)

These are already using the right emails — no action needed. Just update 1Password tags when you do the tag-cleanup pass.

| Address | # Items | Status |
|---------|--------:|--------|
| `me@zeusdelacruz.com` | 18 | ✅ Already on target domain |
| `zeus@zeusdelacruz.com` | 8 | ✅ Keep |
| `zeus.vaults@pm.me` | 12 | ✅ Keep (root of trust) |
| `legacywarriorx@proton.me` | 1 | ✅ Keep (secondary Proton) |
| `zeus@onemindcodex.com` | 3 | ✅ Keep (AWS Acadacom, Datadog, Synadia — review domain) |

---

## Migration Summary

| Block | Source | Accounts | Target | Status |
|-------|--------|:--------:|--------|--------|
| 1 | ui@zeusdelacruz.me | 19 | zeus@zeusdelacruz.com | ⬜ 0/19 |
| 2 | finance@zeusdelacruz.me | 18 | me@zeusdelacruz.com | ⬜ 0/18 |
| 3 | shopping@zeusdelacruz.me | 15 | me@zeusdelacruz.com | ⬜ 0/15 |
| 4 | performance@zeusdelacruz.me | 6 | me@zeusdelacruz.com | ⬜ 0/6 |
| 5 | travel@zeusdelacruz.me | 5 | me@zeusdelacruz.com | ⬜ 0/5 |
| 6 | education@zeusdelacruz.me | 4 | me@zeusdelacruz.com | ⬜ 0/4 |
| 7 | social@zeusdelacruz.me | 4 | zeus@zeusdelacruz.com | ⬜ 0/4 |
| 8 | legal@zeusdelacruz.me | 4 | me@zeusdelacruz.com | ⬜ 0/4 |
| 9 | legacy@ / le@ / recreation@zeusdelacruz.me | 4 | me@zeusdelacruz.com | ⬜ 0/4 |
| 10 | zeus.ops@legacyhq.ai | 19 | zeus@ or me@zeusdelacruz.com | ⬜ 0/19 |
| 11 | finance@legacyhq.ai | 6 | zeus@zeusdelacruz.com | ⬜ 0/6 |
| 12 | partnernetwork@legacyhq.ai | 7 | partners@onemindcodex.com | ⬜ 0/7 |
| 13 | legal@legacyhq.ai | 1 | zeus@zeusdelacruz.com | ⬜ 0/1 |
| **TOTAL** | | **112** | | **⬜ 0/112** |

---

## Pre-Migration Checklist (Do First)

- [ ] **Create `me@zeusdelacruz.com` alias in Zoho Mail** if not already active (check Zoho → Settings → Email Addresses)
- [ ] **Create `speaking@zeusdelacruz.com`** alias
- [ ] **Create `consulting@zeusdelacruz.com`** alias
- [ ] **Create `partners@onemindcodex.com`** alias (required before Block 12)
- [ ] **Create `support@onemindcodex.com`** alias
- [ ] **Keep `zeusdelacruz.me` DNS/MX active** until 100% of blocks above are checked off — needed to receive confirmation emails during account updates
- [ ] **Verify `zeus@zeusdelacruz.com` Zoho inbox is set as primary** Send Mail As identity

## Post-Migration Checklist (Do After All 112 Accounts Done)

- [ ] Let `zeusdelacruz.me` MX records expire (stop paying for the domain)
- [ ] Update Zoho filters: remove `*@zeusdelacruz.me` routing rules
- [ ] Update 1Password tag `path:ui` → `path:so` (separate cleanup pass)
- [ ] Rename 1P vaults: "00-24 UI" → "00-24 SO", "75-99 GE Entrepreneurship" → "75-99 GE Endeavor"
- [ ] Remove human-facing `@legacyhq.ai` Zoho aliases (keep `command@` and `system@` active for Legacy AI)
- [ ] Run `op item list --format json` final audit to confirm 0 items using old addresses

---

## Process Guidance

**Recommended approach — batch by block:**
1. Open a service's account settings
2. Change email to the target address listed above
3. Confirm the change via the confirmation email (it'll arrive at the old `@zeusdelacruz.me` or `@legacyhq.ai` address — this is why you keep those alive)
4. Update the 1Password entry: change the username field AND update tags if doing a tag pass simultaneously
5. Check off the item

**For sites that won't let you change email easily:**
- Use Cloaked.app to generate a bridge address temporarily
- Or contact support with proof of identity

**Priority order (highest to lowest):**
1. **Block 10** — zeus.ops@legacyhq.ai (19 accounts, most business-critical)
2. **Block 12** — partnernetwork@legacyhq.ai (7 accounts, revenue-related)
3. **Block 1** — ui@zeusdelacruz.me (19 accounts, tech infrastructure)
4. **Block 2** — finance@zeusdelacruz.me (18 accounts, financial security)
5. All others in any order
