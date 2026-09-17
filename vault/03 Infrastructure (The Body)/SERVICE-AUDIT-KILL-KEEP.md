---
title: "Service Audit — Kill / Keep / Fix Decision List"
type: ops-decision
status: awaiting-approval
created: 2026-06-11
owner: zeus
category: infrastructure
tags: [audit, services, cleanup, launchd, mcp, redundancy]
related: ["ONEMIND-OS-STACK.md"]
---

# Service Audit — Kill / Keep / Fix

> **Rule:** Legacy proposes, Zeus decides. **Nothing is deleted without your explicit per-line approval.**
> Mark each row: ✅ approve · ❌ reject · ⏸️ defer. Then hand it back.
>
> Audit date: 2026-06-11 (post Mac-reset + stack-solidify pass).

---

## A. Process sprawl — `todoist-mcp` ×8 (~400 MB)

**Root cause (confirmed):** every Hermes profile gateway that defines `mcp_servers:` spawns its **own** stdio subprocess of each MCP server. The 4 business profiles — **ceo, cfo, growth, ops** — each redundantly define `todoist`, `context7`, `excalidraw`, `cua-driver`. With main + dashboard + webui also loading MCPs, that's 8 live `todoist-mcp` node procs, 4 `context7`, 4 `excalidraw`.

**Why it matters:** ~400 MB RAM for duplicate Todoist bridges, plus equal duplication for context7/excalidraw. The business profiles likely don't all need Todoist/Excalidraw/Context7 hot 24/7.

| # | Action | Target | Reclaim | Risk | Default |
|---|--------|--------|---------|------|---------|
| A1 | **Remove** `todoist` MCP from ceo/cfo/growth/ops profiles (keep on `default` only) | 4 profile configs | ~240 MB | Low — those profiles rarely create tasks; can route through Legacy | ✅ recommend |
| A2 | **Remove** `excalidraw` MCP from ceo/cfo/growth/ops | 4 profile configs | ~120 MB | Low — diagram tool, niche per-profile | ✅ recommend |
| A3 | **Remove** `context7` MCP from ceo/cfo/growth/ops | 4 profile configs | ~80 MB | Med — docs lookup; useful for coder, not finance | ⏸️ your call |
| A4 | **Keep** `cua-driver` everywhere | — | — | — | keep |

**Net if A1+A2+A3 approved:** ~440 MB reclaimed, 12 fewer node processes.

---

## B. Failing / flapping services

**Important nuance:** `launchctl` status = **last exit code, not current health**. A service with a live PID is running *now*; the code is just the last time it restarted. Breakdown:

| # | Service | Status | What's actually happening | Action | Default |
|---|---------|--------|---------------------------|--------|---------|
| B1 | `ai.hermes.workspace` | **78, no PID** | **Genuinely dead.** Exit 78 = `EX_CONFIG` (bad config). Not running. | Fix config OR remove plist | ⏸️ decide |
| B2 | `com.onemind.crucix` | **78, no PID** | **IDENTIFIED → FIX, don't cut.** Crucix = your "local intelligence engine" (26-29 OSINT sources, live dashboard :3117, public demo crucix.live). Was working (last log: 28/29 sources OK). Dead only because the Mac reset wiped **node v22.22.3** (plist hardcodes it); only v24.15.0 survived. node_modules intact. **Verified: boots clean under v24.** | 1-line plist fix (v22.22.3 → v24.15.0 path) | ✅ fix |
| B3 | `ai.hermes.gateway-ceo` | -9, live PID | SIGKILL on last restart (likely OOM or forced). Running now (8.5h). | Fix root cause if recurring | ⏸️ watch |
| B4 | `gateway-cfo/growth/ops` | 1, live PID | Exit-1 on last restart, running now (8.5h). Tied to MCP load (section A). | Resolve via A1–A3 | ⏸️ watch |
| B5 | `com.onemind.agenttap` | 1, live PID | Exit-1 last restart, running now. :8088. | Diagnose logs | ⏸️ watch |
| B6 | `com.hermes.dashboard`, `hermes-webui`, `ai.hermes.gateway` (main) | -15, live PID | SIGTERM = **normal restart** (we restarted NATS/services this session). Healthy. | None | keep |

**Key:** Only **B1 + B2 are actually down.** The rest are running; their codes are restart artifacts. The business gateways (B3/B4) share the MCP-sprawl root cause — fixing A may clear their instability.

---

## C. Monitoring / dashboard overlap

You now run **5 overlapping observability surfaces**:

| # | Tool | Port | Role | Action | Default |
|---|------|------|------|--------|---------|
| C1 | Prometheus + Grafana | 9090 / 3009 | Bus internals (nats-surveyor) — **stack doctrine** | **Keep** — canonical | keep |
| C2 | Netdata | 19999 | Host metrics (CPU/RAM/disk) | Keep, but **feed into Grafana** instead of separate UI | ⏸️ your call |
| C3 | Uptime-Kuma | 3001 | Uptime/status page | **Pick one** status page | ⏸️ Kuma or Homepage? |
| C4 | Homepage | 3010 | Dashboard / launcher | **Pick one** status page | ⏸️ Kuma or Homepage? |
| C5 | Synadia agent-UI | 2658-spawned | SwarmClaw agent view | Keep (distinct role) | keep |

**Recommendation:** Grafana = metrics truth; collapse C3/C4 to one. No deletion urgency — this is consolidation, not waste.

---

## D. JetStream health fault (fix, not remove)

| # | Issue | Detail | Action | Default |
|---|-------|--------|--------|---------|
| D1 | `/healthz` returns degraded | Legacy empty stream `$G > AGENT_EVENTS` "could not be recovered" — a leftover in the GLOBAL account from before account isolation. Harmless (empty) but keeps health endpoint amber. | Delete the orphaned empty stream (with your OK) | ⏸️ recommend fix |

---

## Summary — what I recommend you approve

- **A1, A2** (remove todoist + excalidraw from 4 business profiles) → ~360 MB, clean, low-risk.
- **B1, B2** investigation (workspace + crucix are actually dead — decide fix vs remove).
- **D1** (clear the orphaned empty stream → health goes green).
- **Defer:** A3, C-series consolidation, B3–B5 (watch after A).

**Nothing here is executed yet. Hand back the rows you approve.**
