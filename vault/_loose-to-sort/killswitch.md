---
title: "Killswitch"
type: reference
domain: SO
status: active
created: 2026-06-07
updated: 2026-06-07
tags: [killswitch, emergency, system-control, sovereignty]
---

# Killswitch

Emergency stop protocol for OneMind OS. This document defines the conditions, authority, and procedures for halting any agent, service, or the full system.

## Authority

Zeus Delacruz is the sole operator. No agent or external party may invoke a killswitch action without his explicit directive.

## Scope

| Level | Target | Action |
|-------|--------|--------|
| L1 | Single agent | Halt process, drain queue |
| L2 | Agent fleet | Stop all Hermes profiles + OpenClaw |
| L3 | NATS bus | Shutdown NATS server — severs all inter-agent comms |
| L4 | Full system | L3 + all LaunchAgents disabled + network isolation |

## Procedures

### Single Agent (L1)
```bash
# Hermes profile
hermes stop

# OpenClaw
# Via OpenClaw admin panel or process kill
```

### NATS Bus (L3)
```bash
# Stop NATS LaunchAgent
launchctl unload ~/Library/LaunchAgents/nats-server.plist
```

### Full System (L4)
1. L3 — stop NATS
2. Disable all LaunchAgents: agentmemory, Homepage, NATS, Hermes gateway, Langfuse
3. Network isolation if required (Tailscale exit node or physical disconnect)
4. Document the incident in `SO/05 Operations (The Controls)/`

## Recovery

After any killswitch event, document what triggered it and what was halted before re-enabling services. Recovery order: NATS → agentmemory → Hermes gateway → agent fleet → OpenClaw.

## Notes

<!-- Specific trigger conditions, agent-level stop signals, and NATS subject contracts for graceful shutdown go here -->
