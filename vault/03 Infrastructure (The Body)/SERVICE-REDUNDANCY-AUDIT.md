---
title: "Service Redundancy Audit — Kill / Keep / Condense"
type: ops-decision
status: awaiting-approval
created: 2026-06-11
updated: 2026-06-11
owner: zeus
category: infrastructure
tags: [audit, redundancy, services, cleanup, launchd, mcp]
decided_by: zeus
---

# Service Redundancy Audit — Kill / Keep / Condense

> **Rule:** Legacy proposes. **Zeus deletes.** Nothing here is removed without your per-line approval.
> Mark each row **KEEP** / **KILL** / **CONDENSE** and hand it back, or say "do the obvious wins."

**Audit date:** 2026-06-11, post-reset, post stack-install.
**Headline:** ~817 MB reclaimable from duplicate MCP processes + 5 failing gateway services creating log/restart churn.

---

## 1. Duplicate MCP processes — the big reclaim (~817 MB)

Each Hermes gateway profile spawns its **own** copy of every MCP server. With 11 gateway profiles loaded, the same MCP runs 4–8×.

| MCP | Instances | RAM | Root cause |
|---|---|---|---|
| `todoist-mcp` | **8** | ~395 MB | One per active profile that lists Todoist |
| `excalidraw-mcp` | 4 | ~150 MB | One per profile w/ excalidraw |
| `context7-mcp` | 4 | ~137 MB | One per profile w/ context7 |
| `telegram/server.py` | 4 | ~135 MB | One per profile w/ telegram |
| **Total** | **20 procs** | **~817 MB** | Per-profile MCP spawn |

**Recommendation — CONDENSE (highest leverage):**
- **Option A (cleanest):** Cut active gateway profiles to the ones you actually use day-to-day. Fewer profiles → proportionally fewer MCP copies. See §2.
- **Option B:** Move shared MCPs (Todoist, Context7, Excalidraw) to a **single shared MCP host** that all profiles connect to, instead of each spawning its own. Bigger refactor, but the right end-state at fleet scale.
- **Decision needed:** A now, B later? Or both?

---

## 2. Hermes gateway profiles — 11 total, 5 failing

`11` gateway services registered. Several are crash-looping (negative = killed, `1` = error exit):

| Service | Status | Meaning | Proposed |
|---|---|---|---|
| `ai.hermes.gateway` | **-15** | SIGTERM (killed/restarting) | INVESTIGATE |
| `ai.hermes.gateway-ceo` | **-9** | SIGKILL (hard-killed, likely OOM/hang) | INVESTIGATE |
| `ai.hermes.gateway-cfo` | **1** | Error exit | INVESTIGATE |
| `ai.hermes.gateway-growth` | **1** | Error exit | INVESTIGATE |
| `ai.hermes.gateway-ops` | **1** | Error exit | INVESTIGATE |
| `gateway-eternal / sentinel / sovereignty / unity / zenith / coder` | 0 | Healthy | KEEP |

**Question for you:** The OneMind doctrine maps 6 profiles (sovereignty, eternal, zenith, unity, sentinel + default). But there are **11** gateways including `ceo / cfo / growth / ops` — these look like **business-pillar (GE) profiles** that may be stale or experimental.
- Which of `ceo / cfo / growth / ops` do you actually use?
- The healthy doctrine profiles (eternal/sentinel/sovereignty/unity/zenith/coder) → **KEEP**.
- The failing ones → fix or cut. Each one you cut also removes its MCP copies (see §1).

---

## 3. Failing non-gateway services

| Service | Status | What it is | Proposed |
|---|---|---|---|
| `com.hermes.dashboard` | **-15** | Old Hermes status dashboard (:9119) | KILL? — superseded by NATS Fabric page + Homepage |
| `com.onemind.hermes-webui` | **-15** | Hermes WebUI launcher | INVESTIGATE — are you using WebUI? |
| `com.onemind.agenttap` | **1** | AgentTap (unknown/experimental) | IDENTIFY — what is this? |

---

## 4. Monitoring stack — now 5 overlapping tools

After today's install you have **five** monitoring surfaces:

| Tool | Port | Role | Proposed |
|---|---|---|---|
| **Prometheus** | 9090 | Metrics TSDB (bus internals via surveyor) | KEEP — doctrine |
| **Grafana** | 3009 | Dashboards | KEEP — doctrine |
| **Netdata** | 19999 | Host metrics (own UI) | CONDENSE → feed into Grafana, drop separate UI |
| **Uptime-Kuma** | 3001 | Uptime/status checks | PICK ONE ↓ |
| **Homepage** | 3010 | Service dashboard / launcher | PICK ONE ↓ |

**Recommendation:** Prometheus + Grafana = the monitoring stack (your doctrine). Point **Netdata** into Grafana as a data source instead of running its own UI. For the "what's up right now" glance, **pick ONE** of Uptime-Kuma or Homepage — running both is redundant.

---

## 5. JetStream health fault (fix, don't remove)

`/healthz` is **amber**: legacy stream `$G > AGENT_EVENTS` "could not be recovered" (empty, orphaned in the GLOBAL account). Harmless to data but keeps the health endpoint red and will confuse the new NATS Fabric page's health banner.
- **Proposed:** delete the orphaned empty stream (it's empty — zero data loss) to restore green health. **Your call** — it's a `nats stream rm` on an empty legacy stream.

---

## Summary — recommended "obvious wins" (pending your OK)

1. **Identify & cut stale GE gateways** (`ceo/cfo/growth/ops` if unused) → reclaims ~400–600 MB via cascading MCP cleanup.
2. **Fix or kill** `com.hermes.dashboard` + `hermes-webui` (both crash-looping, likely superseded).
3. **Identify `agenttap`** — unknown service erroring.
4. **Monitoring:** keep Prometheus+Grafana, fold Netdata in, pick one status page.
5. **Clear the orphaned `$G > AGENT_EVENTS`** empty stream → green health.

> None executed. Mark KEEP/KILL/CONDENSE per item, or say **"do the obvious wins"** and I'll act on 2–5 (holding 1 until you confirm which GE profiles to keep).
