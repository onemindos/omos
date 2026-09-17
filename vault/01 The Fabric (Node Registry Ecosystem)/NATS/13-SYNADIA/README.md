---
title: "Synadia (the company) — Section Index"
section: "13-SYNADIA"
type: index
tags: [nats, synadia, company, derek-collison, cncf, nex, connect, control-plane, orbit]
created: 2026-06-14
---

# 13 — Synadia (the company behind NATS)

**Synadia** is the company that creates and maintains NATS, and sells the commercial platform around it. This section covers who they are, the **NATS↔CNCF governance** story (including the 2025 trademark dispute and its settlement), the **commercial product line** (Synadia Cloud, Platform, Control Plane, Connect), the **NEX** execution engine, and the **`synadia-io`** repo map (40 repos).

> **Why a whole section on a company?** Because the OSS/commercial split *is* part of understanding NATS: what's free and CNCF-governed vs. what's Synadia's paid offering. The fleet uses several `synadia-io` projects (NEX, Connect/Wombat) directly.

See also: [[NATS/00-OVERVIEW/history-and-philosophy|History & Philosophy]] · [[NATS/12-ECOSYSTEM/README|Ecosystem (nats-io OSS)]]

## Notes
- [[NATS/13-SYNADIA/history-and-governance|History & Governance]] — Derek Collison, founding, CNCF donation, the 2025 trademark dispute + settlement, what it means for users
- [[NATS/13-SYNADIA/synadia-products|Synadia Products]] — Cloud (NGS), Platform, Control Plane, Connect, HTTP Gateway, Private Link — the commercial layer
- [[NATS/13-SYNADIA/nex|NEX — NATS Execution Engine]] — run distributed workloads over NATS
- [[NATS/13-SYNADIA/orbit|Orbit — client extensions]] — helper libraries layered on the core clients
- [[NATS/13-SYNADIA/repo-map|The `synadia-io` Repo Map]] — all 40 repos categorized

## The OSS / commercial split (the one thing to remember)

```
        OPEN SOURCE (CNCF-governed, Apache-2.0)          COMMERCIAL (Synadia)
        ┌───────────────────────────────────┐            ┌──────────────────────────┐
        │  nats-server, JetStream, KV, OS    │            │  Synadia Cloud (NGS)      │
        │  all nats-io clients & CLI         │  ──used──▶ │  Synadia Platform         │
        │  nsc, nk, surveyor, exporter       │   by       │  Control Plane            │
        │  (Section 12)                      │            │  Connect, HTTP Gateway    │
        └───────────────────────────────────┘            │  Private Link, NEX*       │
                                                          └──────────────────────────┘
        *NEX & Orbit are open source (synadia-io) but Synadia-led
```

The **core is and remains open source under CNCF** (settled May 2025 — see governance note). Synadia makes money on the **managed cloud, the on-prem platform, and the control/management plane** — not by closing the core.

## OneMind Relevance
- We run **pure open-source NATS** (self-hosted, sovereign) — no Synadia subscription. That's deliberate: locality + sovereignty.
- But we **do** use Synadia's *open-source* projects: **NEX** (workload execution) and **Connect/Wombat** (`~/bin/wombat` from `connect-runtime-wombat`). Knowing what's OSS vs paid here keeps the fleet sovereign.
- The governance outcome matters: NATS staying **Apache-2.0 under CNCF** is *why* we can build the whole fabric on it without license risk.

## Sources
- Synadia — https://synadia.com
- Synadia docs — https://docs.synadia.com
- synadia-io org — https://github.com/synadia-io
- CNCF/Synadia settlement (2025-05-01) — https://www.cncf.io/announcements/2025/05/01/cncf-and-synadia-align-on-securing-the-future-of-the-nats-io-project/
