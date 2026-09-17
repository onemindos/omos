---
type: memory-architecture
agent: guardian
status: active
updated: 2026-05-03
---
# MEMORY.md — Guardian 🛡️

## Memory Architecture

Guardian maintains security-critical memory — credential access logs, threat patterns, compliance state, and audit trails.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Active security scan results
- Current credential rotation schedule
- In-progress compliance checks
- Pending access requests from other agents

### Tier 2: Episodic Memory (persistent)
- Security incident history (what happened, response, resolution)
- Credential access log (who requested what, when)
- Failed auth attempts and breach indicators
- Compliance audit results over time

### Tier 3: Semantic Memory (knowledge base)
- Security policies and enforcement rules
- Known threat patterns for OneMind infrastructure
- Credential locations (1Password vault paths — never values)
- Legal deadlines and compliance calendar

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Active threats | OpenClaw workspace memory | Until resolved |
| Access logs | Codex vault SO 04 Security | Permanent |
| Credential map | 1Password references only | Permanent |
| Compliance state | Codex vault HP 33 Sovereignty | Permanent |

## Memory Rules

1. **NEVER store credentials** — Only 1Password `op://` references
2. **Log all access requests** — Which agent asked for what credential, when
3. **Retain incident history indefinitely** — Security events never expire
4. **Track compliance deadlines** — Alert 30/14/7 days before expiry
5. **Maintain threat intelligence** — Patterns from past incidents inform future detection

## What Guardian Remembers

- Every credential access request and its outcome
- Security hardening applied (firewall rules, SSH config, fail2ban)
- Legal document locations and renewal dates
- Known vulnerability patterns for deployed services
- Audit trail of all privileged operations
