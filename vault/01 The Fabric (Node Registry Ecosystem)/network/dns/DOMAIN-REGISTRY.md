---
title: "Domain Registry"
type: reference
status: active
created: 2026-04-14
updated: 2026-04-14
owner: zeus
tags: [type/reference, domain/SO, domain/02-infrastructure]
ai_summary: "Complete domain portfolio — GoDaddy (22 active), Cloudflare (6 zones), email stack (Zoho), pending decisions, and migration plan."
---

# Domain Registry

> Single source of truth for all domains, DNS, email, and branding decisions.
> Last audited: **2026-04-14** via GoDaddy API + Cloudflare API (Global API Key).

---

## Credentials (Reference Only — stored in 1Password)

| Service | Vault | 1P Item |
|---------|-------|---------|
| GoDaddy API | `00-24 SO` | API key: `h1eEX3WZFqH9_...` / Secret in 1P |
| Cloudflare Global API | `00-24 SO` | Item: `Cloudflare.com` — email: `zeus.vaults@pm.me` |

---

## Cloudflare — 6 Active Zones

All 6 zones use Cloudflare nameservers: `elliot.ns.cloudflare.com` / `mia.ns.cloudflare.com`

### legacyhq.ai
- **Plan:** Free
- **Website:** Zoho Sites (`www` → `zhs.zohosites.com`)
- **Email:** ✅ Zoho Mail active (MX + SPF + DKIM + DMARC)
- **DMARC sink:** `dmarc@onemindcodex.com`
- **Sub-services:** `forms.legacyhq.ai` → Zoho Forms, `help.legacyhq.ai` → Zoho Desk
- **Issues:** Ghost NS records (ns31/ns32.domaincontrol.com) sitting in zone — should be removed

### onemindcodex.com
- **Plan:** Free
- **Website:** GitHub Pages (apex A records → 185.199.x.x), `live.` → CF Pages (`onemind-codex-garden.pages.dev`), `www` → `zeus-delacruz.github.io`
- **Email:** ✅ Zoho Mail (MX + SPF + DKIM) — DMARC sink domain
- **Purpose:** Public Codex site + DMARC reporting address

### onemindos.com
- **Plan:** Free
- **Website:** Apex A → `13.248.243.5` + `76.223.105.230` (GoDaddy parking IPs — no real site at root)
- **pay.onemindos.com** → GoDaddy Payments (`paylinks.commerce.godaddy.com`)
- **Email:** ❌ None configured — **GAP: main platform domain has no email**
- **DMARC:** `p=quarantine` set but no email sending configured
- **Issues:** Apex points to GoDaddy parking page. Need real content or redirect.

### onemindos.dev
- **Plan:** Free
- **Website:** ❌ EMPTY — zero DNS records
- **Email:** ❌ None
- **Status:** Domain owned, in Cloudflare, completely unused — needs a job

### zeusdelacruz.com
- **Plan:** Free
- **Website:** `www` → Webflow (`cdn.webflow.com`), apex → `198.202.211.1` (GoDaddy parking — broken)
- **Email:** ✅ Zoho Mail (MX + SPF + DKIM + DMARC)
- **Webflow verification:** TXT record present
- **Issues:** Apex doesn't point to Webflow — visitors to `zeusdelacruz.com` hit GoDaddy parking instead of site

### zeusdelacruz.me
- **Plan:** Free
- **Website:** `www` CNAME → `zeusdelacruz.me` (loops with no A record — **broken**)
- **Email:** ✅ Zoho Mail (MX + SPF + DKIM + DMARC)
- **Issues:** www is broken. No website here — email-only domain currently.

---

## GoDaddy — 22 Active Domains

### In Cloudflare Already (6)
See Cloudflare section above: `legacyhq.ai`, `onemindcodex.com`, `onemindos.com`, `onemindos.dev`, `zeusdelacruz.com`, `zeusdelacruz.me`

### Remaining on GoDaddy DNS (16)

| Domain | Expires | Brand Cluster | Priority |
|--------|---------|---------------|----------|
| `buildingonemind.com` | 2027-03-27 | One Mind | — |
| `onemindclaw.com` | 2027-03-18 | One Mind / Platform | — |
| `onemindframework.com` | 2027-03-26 | One Mind | — |
| `weareonemind.com` | 2027-03-26 | One Mind | — |
| `wearelegacy.ai` | 2027-06-06 | Legacy Brand | — |
| `wearelegacy.shop` | 2026-06-06 | Legacy Brand | — |
| `webuildautomations.com` | 2026-06-28 | Services (old) | Low |
| `wealthwarriorsacademy.com` | 2026-09-28 | Old Brand | TBD — drop? |
| `zeuslegacydelacruz.com` | 2026-10-07 | Personal Brand | Low |
| `zeusops.dev` | 2026-09-01 | Tech / Ops | — |
| **`cielodelacruz.com`** | **⚠️ 2026-04-25** | Family | DECIDE |
| **`delacruzlegacy.com`** | **⚠️ 2026-04-25** | Family | DECIDE |
| **`noahdelacruz.com`** | **⚠️ 2026-04-25** | Family | DECIDE |
| **`solunadelacruz.com`** | **⚠️ 2026-04-25** | Family | DECIDE |
| **`legacyhq.shop`** | **⚠️ 2026-04-25** | Legacy Brand | DECIDE |
| **`legacyhq.store`** | **⚠️ 2026-04-25** | Legacy Brand | DECIDE |

---

## URGENT — Expires April 25, 2026 (11 Days)

Decision needed before Apr 25 or auto-renew kicks in / they lapse.

| Domain | Renew? | Decision |
|--------|--------|----------|
| `cielodelacruz.com` | ❓ | |
| `delacruzlegacy.com` | ❓ | |
| `noahdelacruz.com` | ❓ | |
| `solunadelacruz.com` | ❓ | |
| `legacyhq.shop` | ❓ | |
| `legacyhq.store` | ❓ | |

---

## Available Domains — Confirmed (April 12, 2026)

All confirmed available via whois. Not yet registered.

| Domain | Brand Cluster | Value |
|--------|---------------|-------|
| `thecodexmethod.com` | The CODEX Method | ⭐⭐⭐ strong product name |
| `buildyourcodex.com` | The CODEX Method | ⭐⭐⭐ lead gen / landing |
| `joinonemind.com` | One Mind Community | ⭐⭐⭐ community CTA |
| `onemindhq.com` | One Mind Platform | ⭐⭐⭐ clean redirect/HQ |
| `learnonemind.com` | One Mind Education | ⭐⭐ course platform |
| `onemindbook.com` | One Mind Book | ⭐⭐ manifesto domain |
| `buildingyouronemind.com` | One Mind / Codex | ⭐⭐ |
| `onemindlife.com` | One Mind | ⭐⭐ |
| `legacyclaw.ai` | Legacy + Platform | ⭐⭐ |
| `zeusdelacruz.co` | Personal Brand | ⭐ backup |

Decisions: [ ] = not yet decided

---

## Inactive / Dead Domains

### PENDING_HOLD_REDEMPTION (~$80-200 each to recover)
| Domain | Expired |
|--------|---------|
| `legacyacadacom.co` | ~2026-03-26 |
| `legacyacadacom.com` | ~2026-03-26 |
| `legacycapitalist.co` | ~2026-03-26 |
| `legacypartnernetwork.co` | ~2026-03-26 |
| `legacypartnernetwork.com` | ~2026-03-26 |

**Decision:** Recover or abandon when redemption window closes?

### HELD_EXPIRED (effectively dead)
`legacyplatform.co`, `legacyrestore.co`, `webuildplatforms.co`

### CANCELLED (7 — gone)
`adelacruz.me`, `automatemycreditempire.com`, `ceogrowthaccelerator.com`, `legacyplatform.app`, `legacywarrior.me`, `thecreditexecutive.com` + 1 more

---

## Brand Clusters — Full Map

```
ONE MIND (Platform)
  Active:    onemindos.com · onemindos.dev · onemindcodex.com
             onemindclaw.com · onemindframework.com · weareonemind.com
             buildingonemind.com
  Available: onemindhq.com · onemindlife.com · joinonemind.com
             learnonemind.com · onemindbook.com · buildingyouronemind.com

THE CODEX METHOD (Product)
  Active:    onemindcodex.com (repurposable)
  Available: thecodexmethod.com · buildyourcodex.com

LEGACY / LEGACYHQ (Brand — Agent + Services)
  Active:    legacyhq.ai (main) · wearelegacy.ai · wearelegacy.shop
             legacyhq.shop · legacyhq.store
  Available: legacyclaw.ai

ZEUS DELACRUZ (Personal Brand)
  Active:    zeusdelacruz.com (Webflow site) · zeusdelacruz.me
             zeuslegacydelacruz.com · zeusops.dev
  Available: zeusdelacruz.co

DELACRUZ FAMILY
  Active:    cielodelacruz.com · delacruzlegacy.com
             noahdelacruz.com · solunadelacruz.com
  (all expire Apr 25 — decision needed)
```

---

## Email Stack — Current State (Zoho Mail)

| Domain | Email Working | SPF | DKIM | DMARC |
|--------|--------------|-----|------|-------|
| `legacyhq.ai` | ✅ | ✅ | ✅ | ✅ |
| `zeusdelacruz.com` | ✅ | ✅ | ✅ | ✅ |
| `zeusdelacruz.me` | ✅ | ✅ | ✅ | ✅ |
| `onemindcodex.com` | ✅ | ✅ | ✅ | sink |
| `onemindos.com` | ❌ **NOT SET UP** | ❌ | ❌ | partial |
| `onemindos.dev` | ❌ | ❌ | ❌ | ❌ |
| All others | ❌ | ❌ | ❌ | ❌ |

**DMARC reporting sink:** All reports flow to `dmarc@onemindcodex.com`

### Email Gaps / Decisions
- `onemindos.com` has no email — the flagship platform domain should have at minimum `hello@onemindos.com`
- No canonical "primary" email defined — scattered across 4 domains
- Recommended primary: `zeus@legacyhq.ai` (brand) + `zeus@zeusdelacruz.com` (personal)
- Platform contact: `hello@onemindos.com` (needs setup)

---

## Known DNS Issues (Fix When Touching Each Domain)

| Domain | Issue | Fix |
|--------|-------|-----|
| `legacyhq.ai` | Ghost NS records (ns31/ns32.domaincontrol.com) in CF zone | Delete those NS records from CF |
| `onemindos.com` | Apex → GoDaddy parking IPs (13.248.243.5, 76.223.105.230) | Point to real server or Cloudflare redirect |
| `onemindos.dev` | Zero DNS records — completely empty | Assign purpose + add records |
| `zeusdelacruz.com` | Apex → GoDaddy parking (198.202.211.1) | Add A record → Webflow or redirect to www |
| `zeusdelacruz.me` | www CNAME → zeusdelacruz.me with no A record (looping) | Add A record or point to real host |

---

## Cloudflare Migration — Remaining 16 Domains

All 22 active GoDaddy domains should eventually move to Cloudflare NS for unified DNS management.

**Priority order:**
1. Core platform domains: `onemindclaw.com`, `weareonemind.com`, `buildingonemind.com`
2. Legacy brand: `wearelegacy.ai`, `wearelegacy.shop`
3. Personal brand: `zeusops.dev`, `zeuslegacydelacruz.com`
4. Family (if renewed): `cielodelacruz.com`, `noahdelacruz.com`, `solunadelacruz.com`, `delacruzlegacy.com`
5. Low-priority / old: `wealthwarriorsacademy.com`, `webuildautomations.com`

**Migration process per domain:**
1. Add domain in Cloudflare → import existing DNS records
2. Cloudflare provides NS pair (elliot/mia.ns.cloudflare.com)
3. Update nameservers at GoDaddy → Settings → Nameservers → Custom
4. Propagation: 5 min–48 hours
5. Verify in Cloudflare dashboard (status goes from Pending → Active)

---

## Open Decisions Log

| # | Decision | Due | Status |
|---|----------|-----|--------|
| 1 | Renew 6 domains expiring Apr 25 — which to keep? | Apr 24 | ❌ Pending |
| 2 | Which available domains to register? | — | ❌ Pending |
| 3 | Set up email on onemindos.com | — | ❌ Pending |
| 4 | Fix onemindos.com apex (parking → real site) | — | ❌ Pending |
| 5 | Fix zeusdelacruz.com apex (parking → Webflow) | — | ❌ Pending |
| 6 | Assign onemindos.dev a purpose | — | ❌ Pending |
| 7 | Recover or abandon 5 PENDING_HOLD_REDEMPTION domains | — | ❌ Pending |
| 8 | Migrate remaining 16 domains to Cloudflare NS | — | ❌ In progress |
| 9 | Define canonical primary email address | — | ❌ Pending |
