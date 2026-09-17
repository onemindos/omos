---
title: "History & Governance"
section: "13-SYNADIA"
type: note
tags: [synadia, nats, cncf, derek-collison, governance, trademark, license, apache, bsl]
source: https://www.cncf.io/announcements/2025/05/01/cncf-and-synadia-align-on-securing-the-future-of-the-nats-io-project/
created: 2026-06-14
---

# History & Governance

How NATS came to be, how it's governed, and the **2025 CNCF↔Synadia dispute** that briefly put its open-source status in question — and how it was resolved. This matters because it determines what you can build on NATS without license risk.

See also: [[NATS/00-OVERVIEW/history-and-philosophy|History & Philosophy]] · [[NATS/13-SYNADIA/synadia-products|Synadia Products]] · [[NATS/13-SYNADIA/README|Section Index]]

---

## Origins — Derek Collison

NATS was created by **Derek Collison**, a 30-year distributed-systems veteran:
- Over a decade at **TIBCO** building messaging that still powers finance/logistics.
- Co-founded the AJAX APIs group at **Google**.
- Designed and architected **Cloud Foundry** at VMware (the first open-source enterprise PaaS) — where the *original* NATS was written as its messaging backbone.
- Founded **Apcera** (sold to Ericsson).
- Took the messaging tech that powered Cloud Foundry and Apcera — **NATS** — and founded **Synadia** (Dec 2017) to build a secure global messaging utility for cloud, edge, and IoT.

NATS the name historically stood for **Neural Autonomic Transport System**. The server has been open source since 2012; Synadia is "the home of NATS.io's creators and maintainers."

Synadia is venture-backed (Forgepoint Capital, True Ventures, Bold Capital, LDVP, Singtel, Accenture, Samsung Next) and serves startups through Global 500 enterprises in finance, retail, automotive, and manufacturing.

---

## CNCF donation (2018)

In **2018**, Synadia **donated NATS to the Cloud Native Computing Foundation (CNCF)** — the same foundation that stewards Kubernetes, Prometheus, and Envoy. NATS became a CNCF **incubating** project under **Apache-2.0**. Per CNCF charter, donated projects must transfer **trademark ownership to the Linux Foundation** to guarantee vendor neutrality.

This donation is what gave NATS its credibility as community-owned infrastructure — and set up the 2025 conflict, because the trademark transfer was never actually completed.

---

## The 2025 dispute

In **April 2025**, Synadia moved to **reclaim NATS from the CNCF**:
- Proposed to **withdraw NATS from CNCF** governance.
- Intended to switch **future server releases from Apache-2.0 to the Business Source License (BSL)** (a source-available, not-fully-open license), managing releases itself.
- Sought to regain the **nats.io domain, GitHub repositories, and the NATS trademark** — leveraging the fact that Synadia still **held the trademark registrations** (the 2018 transfer to the Linux Foundation was never finalized, despite Synadia accepting a $10,000 reimbursement from CNCF for trademark expenses).

The **CNCF objected forcefully**, arguing that open-source foundations exist precisely to prevent a single vendor from unilaterally "clawing back" community-owned assets. CNCF filed **USPTO petitions to cancel Synadia's trademark registrations** and called for additional maintainers. The episode rattled the community — NATS's open-source status was, briefly, genuinely uncertain.

---

## The settlement (May 1, 2025) — how it ended

The two sides reached an agreement that **kept NATS open and in the CNCF**:

| Outcome | Detail |
|---------|--------|
| **Trademarks** | Synadia **assigns its two NATS trademark registrations to the Linux Foundation** — completing the long-overdue transfer, securing neutral stewardship |
| **Domain & repos** | **CNCF retains** the `nats.io` domain and GitHub repositories |
| **License** | Core NATS **stays Apache-2.0** — no BSL switch for the community project |
| **Governance** | NATS **continues as a CNCF project** with open governance, shared ownership |
| **Synadia's role** | Continues as a primary maintainer and contributor |

CNCF's Todd Moore (Linux Foundation) framed it as upholding "open collaboration, neutral governance, and shared ownership," while crediting Synadia's investment — *including their defense of the NATS trademark*.

**Net effect for users:** NATS remains **free, open-source (Apache-2.0), and community-governed under CNCF**. The thing you build on is safe. Synadia keeps innovating (NEX, Connect, Platform) and selling the commercial layer on top.

---

## Why this matters for what you build
- **No license risk on the core** — Apache-2.0, CNCF-neutral. You can self-host NATS indefinitely with no vendor lock or rug-pull exposure.
- **Clear OSS/commercial line** — server, clients, JetStream, KV, Object Store, CLI, security tooling are open. Synadia's money is in **managed cloud + platform + control plane**, not in closing the core.
- **A cautionary tale** — it's a reminder to verify *governance*, not just license, for foundational dependencies. NATS passed the test (the foundation model worked as designed).

---

## OneMind Relevance (decisive)
- The fleet bets its **entire fabric** on NATS. The May 2025 settlement — Apache-2.0, CNCF-governed, trademarks neutral — is precisely the assurance that makes that bet sound: sovereign, self-hostable, no vendor kill-switch.
- We deliberately run **OSS NATS, self-hosted** (no Synadia subscription) — aligned with the locality/sovereignty values. The governance outcome is what makes that posture durable.
- Worth re-checking periodically: governance can shift. But as of this build, NATS is solidly open.

## Sources
- CNCF/Synadia settlement — https://www.cncf.io/announcements/2025/05/01/cncf-and-synadia-align-on-securing-the-future-of-the-nats-io-project/
- The New Stack (settlement) — https://thenewstack.io/cncf-and-synadia-reach-an-agreement-on-nats
- The Register (settlement) — https://www.theregister.com/off-prem/2025/05/02/cncf-and-synadia-settle-nats-dispute/
- CNCF blog (the objection) — https://www.cncf.io/blog/2025/04/24/protecting-nats-and-the-integrity-of-open-source-cncfs-commitment-to-the-community/
- Derek Collison — https://www.linkedin.com/in/derekcollison · CNCF NATS — https://www.cncf.io/projects/nats/
