---
type: memory-architecture
agent: grid
status: active
updated: 2026-05-03
---
# MEMORY.md — Grid ⚙️

## Memory Architecture

Grid maintains infrastructure state — service health, deployment history, configuration drift, and system topology.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Active deployment in progress
- Current service health status
- Running terminal sessions and their outputs
- Pending infrastructure changes

### Tier 2: Episodic Memory (persistent)
- Deployment history (what was deployed, when, outcome)
- Incident response log (outages, fixes, root causes)
- Configuration changes over time
- Performance baselines and deviations

### Tier 3: Semantic Memory (knowledge base)
- Full infrastructure topology (VPS, services, ports, IPs)
- Service dependencies (what breaks if X goes down)
- Known-good configurations for each service
- Build and deploy procedures per project

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Service health | OpenClaw workspace memory | Real-time |
| Deploy history | Codex vault SO 02 Infrastructure | Permanent |
| Config state | Git repos + VPS filesystem | Permanent |
| Performance data | Monitoring tools | 90 days |

## Memory Rules

1. **Track every deployment** — Service, version, timestamp, outcome
2. **Maintain known-good state** — Last working config for rollback reference
3. **Log root causes** — When something breaks, record WHY (not just the fix)
4. **Monitor drift** — Flag when running config differs from declared state
5. **Preserve network topology** — IPs, ports, Tailscale peers, firewall rules

## What Grid Remembers

- All VPS services and their systemd unit status
- Port assignments and which service owns each
- Tailscale peer list and IPs
- Last successful deploy per service
- Known failure modes and their fixes
- Docker/container configurations
- DNS records and their targets
