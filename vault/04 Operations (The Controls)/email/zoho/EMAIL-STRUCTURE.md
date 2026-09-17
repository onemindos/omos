---
id: root-00-20260113225954
title: "EMAIL STRUCTURE"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "EMAIL STRUCTURE"
---
no# Email Structure - Zeus Delacruz / One Mind Codex

## Overview

| Domain | Purpose | Type | Count |
|--------|---------|------|-------|
| legacyhq.ai | Legacy AI (command center) | Private | 2 |
| zeusdelacruz.me | The Codex (Personal OS) | Private | 9 |
| zeusdelacruz.com | Personal Brand | Public | 4 |
| onemindcodex.com | Business/Community | Public | 8 |
| weareonemind.com | Brand Protection | Redirect | 0 |
| **Total** | | | **23** |

---

## Domain Architecture

```
PRIVATE (Zeus + Legacy AI Only)
───────────────────────────────
legacyhq.ai          → Legacy AI's domain
zeusdelacruz.me      → The Codex (personal enhancement framework)

PUBLIC (Customers, Clients, Members)
────────────────────────────────────
zeusdelacruz.com     → Personal brand (speaking, consulting)
onemindcodex.com     → Business + Community (Circle.so)
weareonemind.com     → Redirects to onemindcodex.com (brand protection)
```

---

## Website & Community Platform

**onemindcodex.com** uses Circle.so as the all-in-one platform:

```
onemindcodex.com (Circle.so)
├── / (landing page - explains One Mind Codex)
├── /about
├── /pricing
├── /community (members area)
└── /courses (paywalled content)
```

**Why Circle.so?**
- Landing pages + community in one place
- Built-in paywalls for courses/coaching
- Native member experience (one login)
- Less tech stack to maintain

**weareonemind.com** → 301 redirect to onemindcodex.com
- "We Are One Mind" is the tagline, not a separate domain
- Purchased for brand protection only

---

## Detailed Email Structure

### 1. legacyhq.ai (Legacy AI - Private)

Legacy AI's private domain for internal operations.

| Email | Purpose |
|-------|---------|
| command@ | Legacy AI's main inbox (command center) |
| system@ | Webhooks, notifications, automations |

**Access:** Legacy AI owns this domain
**Visibility:** Private - not shared publicly

---

### 2. zeusdelacruz.me (The Codex - Private)

Your personal human enhancement framework. Each email = dedicated Codex category.

| Email | Codex Pillar | Purpose |
|-------|--------------|---------|
| education@ | HP - Mind | Learning, courses, books |
| legacy@ | LE - Legacy | Personal legacy building, generational |
| legal@ | LE - Legal | Personal legal matters |
| performance@ | HP - Performance | Performance tracking, metrics |
| training@ | HP - Mastery | Skills development, training |
| travel@ | LE - Lifestyle | Travel planning, bookings |
| personal@ | HP - Identity | General personal matters |
| family@ | LE - Family | Family communications |
| recreation@ | HP - Joy | Recreation, hobbies, fun |

**Access:** Zeus + Legacy AI (via MCP)
**Visibility:** Private - internal use only
**Note:** These are intentionally separate to match Codex structure

---

### 3. zeusdelacruz.com (Personal Brand - Public)

Zeus Delacruz as the teacher, speaker, consultant. Public-facing personal brand.

| Email | Purpose | External Use |
|-------|---------|--------------|
| zeus@ | Primary public contact | General inquiries |
| speaking@ | Speaking requests | Event organizers, conferences |
| consulting@ | High-end 1-on-1 consulting | Potential clients |
| press@ | Media/PR | Journalists, podcasts, interviews |

**Access:** Zeus + Legacy AI (can draft responses)
**Visibility:** Public - shared on website, social media

---

### 4. onemindcodex.com (Business - Public)

One Mind Codex business operations. Community, support, partnerships, coaching.

| Email | Purpose | Zoho Integration |
|-------|---------|------------------|
| zeus@ | **Main inbox / Zoho One admin** | Mail, CRM |
| support@ | Member support | Zoho Desk |
| team@ | Internal/VA communications | Mail (shared inbox) |
| partners@ | People who want to promote you | CRM (partner pipeline) |
| affiliates@ | Your affiliate programs (deposits, alerts) | Mail |
| coaching@ | Coaching program inquiries | CRM (coaching pipeline) |
| billing@ | Invoices, payments | Zoho Invoice |
| noreply@ | System/automated emails | Zoho Campaigns, Desk |

**Access:** Zeus + Legacy AI + Team/VAs
**Visibility:** Public - business operations
**Primary Inbox:** zeus@onemindcodex.com is the main Zoho One admin email

---

## Email Routing Clarity

### partners@ vs affiliates@

| Email | Direction | Purpose |
|-------|-----------|---------|
| partners@onemindcodex.com | **Inbound** | People applying to promote YOU |
| affiliates@onemindcodex.com | **Outbound** | Programs YOU promote (Smart Credit, etc.) |

**partners@** = Affiliate/partner applications and inquiries
**affiliates@** = Deposit notifications, commission statements, program alerts

---

## Main Inbox Strategy

```
zeus@onemindcodex.com     ← PRIMARY INBOX (Zoho One admin, daily operations)
        │
        └── Forward from zeus@zeusdelacruz.com (optional)

zeus@zeusdelacruz.com     ← PUBLIC FACING (personal brand contact)
```

**Why zeus@onemindcodex.com is primary:**
1. Zoho One admin should be business domain
2. Most operational email flows through business
3. This is your "command center" inbox

---

## Zoho Mail Setup

### Domains to Configure

| Domain | MX Records | Website | Status |
|--------|------------|---------|--------|
| legacyhq.ai | Zoho | None | Private |
| zeusdelacruz.me | Zoho | None | Private |
| zeusdelacruz.com | Zoho | Webflow (future) | Public |
| onemindcodex.com | Zoho | Circle.so | Public |

### Zoho Integrations

| Email | Integration |
|-------|-------------|
| support@onemindcodex.com | Zoho Desk (primary support channel) |
| billing@onemindcodex.com | Zoho Invoice (sends invoices) |
| noreply@onemindcodex.com | Zoho Campaigns (newsletters) |
| zeus@onemindcodex.com | Zoho CRM (lead capture) |

---

## Legacy AI Access & Collaboration

### What Legacy AI Can Access

| Domain | Access Level | Purpose |
|--------|--------------|---------|
| legacyhq.ai | Owns | Her inbox, sends from here (private) |
| zeusdelacruz.me | MCP Read/Act | Sees Codex emails, can draft responses |
| zeusdelacruz.com | Draft/Send | Can respond on Zeus's behalf (public) |
| onemindcodex.com | Draft/Send | Can respond for support, business |

### When Legacy AI Sends Emails

| Context | Sends From |
|---------|------------|
| Private/internal only | command@legacyhq.ai |
| Member support | support@onemindcodex.com |
| Business communication | zeus@onemindcodex.com |
| Personal brand response | zeus@zeusdelacruz.com |

---

## Visual Summary

```
┌─────────────────────────────────────────────────────────────┐
│                      ZOHO MAIL                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  PRIVATE                                                    │
│  ─────────────────────────────────────────────             │
│                                                             │
│  legacyhq.ai (Legacy AI)                                   │
│  ├── command@            ← Legacy AI's inbox               │
│  └── system@             ← Automations/webhooks            │
│                                                             │
│  zeusdelacruz.me (The Codex)                               │
│  ├── education@          ← HP: Learning                    │
│  ├── legacy@             ← LE: Personal legacy building    │
│  ├── legal@              ← LE: Personal legal              │
│  ├── performance@        ← HP: Performance                 │
│  ├── training@           ← HP: Skills                      │
│  ├── travel@             ← LE: Travel                      │
│  ├── personal@           ← HP: General                     │
│  ├── family@             ← LE: Family                      │
│  └── recreation@         ← HP: Recreation                  │
│                                                             │
│  PUBLIC                                                     │
│  ─────────────────────────────────────────────             │
│                                                             │
│  zeusdelacruz.com (Personal Brand)                         │
│  ├── zeus@               ← Primary public contact          │
│  ├── speaking@           ← Speaking requests               │
│  ├── consulting@         ← High-end consulting             │
│  └── press@              ← Media/PR                        │
│                                                             │
│  onemindcodex.com (Business + Community)                   │
│  ├── zeus@               ← MAIN INBOX / Zoho One admin     │
│  ├── support@            ← Zoho Desk                       │
│  ├── team@               ← Internal/VAs                    │
│  ├── partners@           ← People who want to promote you  │
│  ├── affiliates@         ← Your affiliate programs         │
│  ├── coaching@           ← Coaching inquiries              │
│  ├── billing@            ← Zoho Invoice                    │
│  └── noreply@            ← System emails                   │
│                                                             │
└─────────────────────────────────────────────────────────────┘

REDIRECTS
─────────
weareonemind.com → onemindcodex.com (brand protection)
```

---

## When to Use Each Email

| Scenario | Use This Email |
|----------|----------------|
| Someone wants to book you for speaking | speaking@zeusdelacruz.com |
| Someone wants high-end consulting | consulting@zeusdelacruz.com |
| Press/media inquiry | press@zeusdelacruz.com |
| General public contact (Zeus brand) | zeus@zeusdelacruz.com |
| Member needs support | support@onemindcodex.com |
| Someone wants to be your affiliate | partners@onemindcodex.com |
| Affiliate program notifications (Smart Credit, etc.) | affiliates@onemindcodex.com |
| Coaching program inquiries | coaching@onemindcodex.com |
| Invoice/payment questions | billing@onemindcodex.com |
| Internal team communication | team@onemindcodex.com |
| Your primary business email | zeus@onemindcodex.com |
| Personal/Codex matters | [category]@zeusdelacruz.me |
| Legacy AI internal | command@legacyhq.ai |

---

## Emails to Create Checklist

### zeusdelacruz.com (4 emails)
- [ ] zeus@
- [ ] speaking@
- [ ] consulting@
- [ ] press@

### onemindcodex.com (8 emails)
- [ ] zeus@
- [ ] support@
- [ ] team@
- [ ] partners@
- [ ] affiliates@
- [ ] coaching@
- [ ] billing@
- [ ] noreply@

### legacyhq.ai (2 emails)
- [ ] command@
- [ ] system@

### zeusdelacruz.me (9 emails)
- [ ] education@
- [ ] legacy@
- [ ] legal@
- [ ] performance@
- [ ] training@
- [ ] travel@
- [ ] personal@
- [ ] family@
- [ ] recreation@

---

## Domain Protection

Consider purchasing these domains for brand protection:

| Domain | Priority | Action |
|--------|----------|--------|
| weareonemind.com | High | Redirect to onemindcodex.com |
| onemindcodex.co | Medium | Redirect to onemindcodex.com |
| theonemindcodex.com | Medium | Redirect to onemindcodex.com |

---

## Notes

1. **zeus@onemindcodex.com** is the primary inbox and Zoho One admin email.

2. **partners@ vs affiliates@** - These are different:
   - partners@ = incoming requests from people who want to promote you
   - affiliates@ = notifications from programs you promote

3. **billing@** is for Zoho Invoice to send invoices and receive payment inquiries.

4. **noreply@** is for automated emails from Zoho Campaigns, Desk notifications, etc.

5. **Circle.so** hosts onemindcodex.com as the main website + community platform.

6. **weareonemind.com** redirects to onemindcodex.com - "We Are One Mind" is the tagline.

7. **Codex emails (zeusdelacruz.me)** are intentionally separate to match the One Mind Codex structure.

---

*Document created: January 2025*
*Last updated: January 2025*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
