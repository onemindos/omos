---
id: overwatch-watchlist
title: "Watchlist"
type: register
status: active
created: 2026-05-04
modified: 2026-05-04
owner: zeus
tags: [type/register, domain/so, status/active, topic/overwatch]
agent: Legacy
zone: SO
---

# Watchlist

> Active monitoring targets — people, assets, metrics, signals, and systems.
> Everything on this list gets eyes-on during every Overwatch cycle.
> Add items when they matter. Remove or archive when they no longer do.

---

## Watchlist Status Legend

| Status | Meaning |
|--------|---------|
| 🟢 STABLE | Within normal parameters — routine monitoring |
| 🟡 WATCH | Showing variance — elevated attention |
| 🔴 ACTIVE | Requires immediate or daily active monitoring |
| ⚪ PAUSED | Temporarily removed from active watch |
| 🗄️ ARCHIVED | No longer active — historical reference |

---

## Section 1 — Systems & Infrastructure

| ID | Target | What to Watch | Frequency | Status | Last Checked | Notes |
|----|--------|--------------|-----------|--------|-------------|-------|
| SYS-01 | VPS omos-core (100.123.126.34) | CPU, RAM, disk, uptime | 6h | 🟢 STABLE | — | Primary compute node |
| SYS-02 | OpenClaw Gateway | Uptime, response time, auth | 6h | 🟢 STABLE | — | Port 18789, Tailscale served |
| SYS-03 | onemind-watch.service | Running, last commit time | Daily | 🟢 STABLE | — | VPS→GitHub sync |
| SYS-04 | onemind-pull.timer | Running, last pull time | Daily | 🟢 STABLE | — | GitHub→VPS sync |
| SYS-05 | GitHub (OneMind-OS/onemind-codex) | Last push, open issues, PR status | Daily | 🟢 STABLE | — | Source of truth |
| SYS-06 | All 9 Agents | Heartbeat, last response | 6h | 🟢 STABLE | — | Legacy, Spartan, Oracle, Eden, Haven, Empire, Guardian, Grid, Trinity |
| SYS-07 | fail2ban | Active bans, recent attacks | Daily | 🟢 STABLE | — | VPS SSH protection |
| SYS-08 | SSL certificates | Expiry dates | Weekly | 🟢 STABLE | — | Flag 30 days before expiry |

---

## Section 2 — Financial Signals

| ID | Target | What to Watch | Frequency | Status | Last Checked | Notes |
|----|--------|--------------|-----------|--------|-------------|-------|
| FIN-01 | Primary bank account | Balance, incoming, outgoing | Weekly | 🔴 ACTIVE | — | Oracle owns. Alert if < Zeus-defined floor |
| FIN-02 | Monthly income total | Sum across all income streams | Monthly | 🔴 ACTIVE | — | Compare to prior month and target |
| FIN-03 | Monthly expenses | Fixed + variable spend | Monthly | 🔴 ACTIVE | — | Flag if >80% of income |
| FIN-04 | Cash reserve | Total liquid assets | Monthly | 🔴 ACTIVE | — | Must maintain minimum 3-month runway |
| FIN-05 | Investment portfolio | Value, allocation, performance | Monthly | 🟡 WATCH | — | GE 79 domain |
| FIN-06 | Real estate obligations | Payments, maintenance due | Monthly | 🟡 WATCH | — | GE 76, LE 50 |
| FIN-07 | Subscription and recurring spend | All active subscriptions | Monthly | 🟡 WATCH | — | Audit quarterly for unused services |

---

## Section 3 — Health & Performance Signals

| ID | Target | What to Watch | Frequency | Status | Last Checked | Notes |
|----|--------|--------------|-----------|--------|-------------|-------|
| HP-01 | Zeus energy level | Self-reported 1-10 | Daily | 🔴 ACTIVE | — | Trigger TRG-HP-02 if ≤3 for 2 days |
| HP-02 | Sleep quality/duration | Hours, quality rating | Daily | 🔴 ACTIVE | — | Spartan monitors |
| HP-03 | Training consistency | Sessions per week vs. target | Weekly | 🟡 WATCH | — | Flag if 2 consecutive weeks below target |
| HP-04 | Weight/body composition | Trend, not single data point | Weekly | 🟡 WATCH | — | HP 28 domain |
| HP-05 | Mental load / stress level | Self-reported | Weekly | 🔴 ACTIVE | — | High stress + low energy = Culminating Point risk |
| HP-06 | Medication / supplement consistency | Daily | Daily | 🟡 WATCH | — | If applicable |

---

## Section 4 — Brand & Business Signals

| ID | Target | What to Watch | Frequency | Status | Last Checked | Notes |
|----|--------|--------------|-----------|--------|-------------|-------|
| BRD-01 | Content publishing cadence | Last publish date | Weekly | 🟡 WATCH | — | TRG-GE-05 fires at 14 days |
| BRD-02 | Social media follower growth | Trend across platforms | Weekly | 🟡 WATCH | — | Empire monitors |
| BRD-03 | Email list / community size | Growth rate, churn | Weekly | 🟡 WATCH | — | If applicable |
| BRD-04 | onemind-codex GitHub stars | Growth trend | Weekly | 🟡 WATCH | — | Leading indicator of community interest |
| BRD-05 | Competitor activity (PKM/AI space) | Major announcements, launches | Weekly | 🟡 WATCH | — | Second Brain, Fabric, Notion AI, Obsidian |
| BRD-06 | Inbound opportunities | Partnerships, collab requests, media | Daily | 🟡 WATCH | — | Don't let these go stale > 48h |

---

## Section 5 — Relationships & Network

| ID | Target | What to Watch | Frequency | Status | Last Checked | Notes |
|----|--------|-------------|-----------|--------|-------------|-------|
| REL-01 | Inner circle (key relationships) | Last meaningful contact | Weekly | 🔴 ACTIVE | — | Haven monitors. Flag if >14 days no contact |
| REL-02 | Key partnerships / collaborators | Relationship health, last contact | Weekly | 🟡 WATCH | — | Empire + Haven |
| REL-03 | Mentors / advisors | Last contact | Monthly | 🟡 WATCH | — | Flag if >30 days |
| REL-04 | Family obligations | Scheduled events, commitments | Weekly | 🔴 ACTIVE | — | Haven. Non-negotiable. |

---

## Section 6 — Legal & Compliance

| ID | Target | What to Watch | Deadline | Status | Last Checked | Notes |
|----|--------|--------------|----------|--------|-------------|-------|
| LEG-01 | Business entity filings | Annual report, renewal dates | Varies | 🟡 WATCH | — | GE 78 Holdings. 30-day advance flag. |
| LEG-02 | Tax obligations | Quarterly and annual deadlines | Varies | 🟡 WATCH | — | Oracle + Guardian. Never miss. |
| LEG-03 | Active contracts / agreements | Expiry, renewal, obligation dates | Varies | 🟡 WATCH | — | Flag 60 days before any expiry |
| LEG-04 | IP / trademark status | Registration status, renewal dates | Varies | 🟡 WATCH | — | Brand and product IP |
| LEG-05 | Property legal status | Mortgage, insurance, property tax | Annually | 🟡 WATCH | — | LE 50, GE 76 |

---

## Section 7 — Intelligence Targets (OSINT)

| ID | Target | What to Watch | Frequency | Status | Notes |
|----|--------|--------------|-----------|--------|-------|
| INT-01 | AI model releases (NVIDIA, Anthropic, OpenAI, Google) | New models, capability announcements | Weekly | 🟢 STABLE | Relevant to SO 01 models registry |
| INT-02 | PKM / second brain community | New frameworks, tools, influencers | Weekly | 🟡 WATCH | Competitive intelligence for One Mind |
| INT-03 | AI agent framework space (OpenClaw, LangChain, etc.) | Updates, new players | Weekly | 🟡 WATCH | SO 01 + infrastructure |
| INT-04 | Real estate market (Zeus's target markets) | Prices, rates, inventory | Monthly | 🟡 WATCH | GE 76, LE 50 |
| INT-05 | Regulatory changes (AI, crypto, business) | New laws, compliance requirements | Monthly | 🟡 WATCH | Guardian + Oracle |

---

## SITREP Log

A record of all SITREPs delivered during active Overwatch. Update after each delivery.

| Serial | Date/Time | State | Delivered | Key Flags | Zeus Response |
|--------|-----------|-------|-----------|-----------|--------------|
| — | — | — | — | — | — |

---

## Incident Log

A record of all T3+ trigger events and how they were resolved.

| Date | Trigger ID | Description | Severity | Response | Resolution | AAR Done |
|------|-----------|-------------|----------|----------|-----------|---------|
| — | — | — | — | — | — | — |

---

*Review and update watchlist at each sprint boundary.*
*Add items when a new domain, asset, or relationship becomes active.*
*Archive items when they are no longer relevant — never delete.*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
