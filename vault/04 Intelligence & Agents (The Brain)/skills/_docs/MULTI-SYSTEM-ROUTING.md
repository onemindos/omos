---
name: multi-system-routing
description: "Skill assignment matrix — which of the 67 One Mind skills are available to Claude Code, OpenClaw, and Letta, plus per-agent visibility rules"
type: reference
status: active
created: 2026-05-11
owner: legacy
tags: [type/reference, domain/so-01, stage/implementation]
---

# Multi-System Skill Routing

Complete 67-skill assignment matrix showing which systems have access and per-agent visibility rules.

---

## Quick Reference — Systems

| System | Runtime Skills | Vault-Only | Total Access |
|--------|----------------|-----------|--------------|
| **Claude Code** | 67 (all) | 0 | Full vault |
| **OpenClaw** | 33 | 34 vault-only | Runtime only |
| **Letta** | TBD (pending integration) | TBD | TBD |

---

## Full 67-Skill Matrix

### SO (Sovereign Operations) — 15 Skills

| Skill | Claude Code | OpenClaw | Letta | Agent(s) | Notes |
|-------|-------------|----------|-------|----------|-------|
| om-openclaw | ✅ | ✅ | — | Legacy, Sovereignty | Admin access to gateway |
| code-method | ✅ | ✅ | ✅ | All | Universal operating loop |
| grid | ✅ | ✅ | — | Sovereignty | VPS + infrastructure |
| guardian | ✅ | ✅ | — | Sovereignty | Security + access control |
| guardian-ops | ✅ | ✅ | — | Sovereignty | Security ops automation |
| overwatch | ✅ | ✅ | — | Sovereignty | n8n + workflow automation |
| eagleeye | ✅ | ✅ | — | Sovereignty | Drone + robotics ops |
| trinity | ✅ | ✅ | — | Sovereignty | Deep research + R&D |
| onemind-maintainer | ✅ | ✅ | — | Sovereignty | Vault maintenance |
| onemind-report | ✅ | ✅ | — | Legacy | Executive reporting (promoted) |
| onemind-status | ✅ | ✅ | — | Legacy | System health checks (promoted) |
| om-dashboard | ✅ | ✅ | — | Legacy | Dashboard orchestration (promoted) |
| claude-code-guide | ✅ | — | — | All | Claude Code CLI reference |
| feature-dev | ✅ | — | — | All | Feature architecture guidance |
| superpowers | ✅ | — | — | All | Skills framework reference |

**Runtime Skills in OpenClaw:** 10 (om-openclaw, code-method, grid, guardian, guardian-ops, overwatch, eagleeye, trinity, onemind-maintainer, + 4 promoted)

**Vault-Only Skills:** 4 (claude-code-guide, feature-dev, superpowers, TBD)

---

### HP (Holistic Performance) — 12 Skills

| Skill | Claude Code | OpenClaw | Letta | Agent(s) | Notes |
|-------|-------------|----------|-------|----------|-------|
| spartan | ✅ | ✅ | — | Zenith | Fitness + biometrics |
| oracle | ✅ | ✅ | — | Zenith | Psychology + habits |
| identity | ✅ | ✅ | — | Zenith | Identity + purpose |
| librarian | ✅ | ✅ | — | Zenith | Intelligence pipeline |
| pulse | ✅ | ✅ | — | Zenith | Real-time news + feeds |
| genesis | ✅ | ✅ | — | Zenith | Spiritual practice + meaning |
| statusline-setup | ✅ | — | — | Zenith | Local CLI configuration |
| plan | ✅ | — | — | All | Software architecture planning |
| Explore | ✅ | — | — | All | Codebase search + analysis |
| code-simplifier | ✅ | — | — | All | Code refactoring guidance |
| apex | ✅ | ✅ | ✅ | Eternal | Finance + fundraising |
| code-reviewer | ✅ | — | — | All | Code review + QA |

**Runtime Skills in OpenClaw:** 6 (spartan, oracle, identity, librarian, pulse, genesis)

**Vault-Only Skills:** 6 (statusline-setup, plan, Explore, code-simplifier, code-reviewer, TBD)

---

### LE (Legacy Evolution) — 8 Skills

| Skill | Claude Code | OpenClaw | Letta | Agent(s) | Notes |
|-------|-------------|----------|-------|----------|-------|
| haven | ✅ | ✅ | — | Unity | Family + relationships |
| helix | ✅ | ✅ | — | Unity | Parenting + elder care |
| eden | ✅ | ✅ | — | Unity | Sustainability + ecology |
| nomad | ✅ | ✅ | — | Unity | Off-grid + homestead |
| Plan | ✅ | — | — | All | Implementation planning |
| TeamCreate | ✅ | — | — | All | Team coordination tool |
| TeamDelete | ✅ | — | — | All | Team cleanup tool |
| (reserved) | — | — | — | — | LE 8th skill TBD |

**Runtime Skills in OpenClaw:** 4 (haven, helix, eden, nomad)

**Vault-Only Skills:** 4 (Plan, TeamCreate, TeamDelete, reserved)

---

### GE (Generational Empire) — 25 Skills

| Skill | Claude Code | OpenClaw | Letta | Agent(s) | Notes |
|-------|-------------|----------|-------|----------|-------|
| empire | ✅ | ✅ | ✅ | Eternal | Business strategy |
| creed | ✅ | ✅ | ✅ | Eternal | Brand narrative + storytelling |
| empire-ops | ✅ | ✅ | — | Eternal | Operations automation |
| video-production | ✅ | ✅ | ✅ | Eternal | Content production pipeline |
| finance | ✅ | ✅ | ✅ | Eternal | Business finance |
| finance-ops | ✅ | ✅ | — | Eternal | Finance operations |
| apex | ✅ | ✅ | ✅ | Eternal | Fundraising + investor relations |
| legacy | ✅ | ✅ | — | Legacy | Cross-domain orchestration |
| homestead | ✅ | ✅ | — | Unity | Home + property (promoted) |
| onemind-maintainer | ✅ | ✅ | — | Sovereignty | Vault maintenance (shared) |
| (9 more GE skills) | TBD | — | — | TBD | Under review for scope |

**Runtime Skills in OpenClaw:** 10 (empire, creed, empire-ops, video-production, finance, finance-ops, apex, legacy, homestead, om-dashboard)

**Vault-Only Skills:** TBD

---

### One Mind System — 7 Skills

| Skill | Claude Code | OpenClaw | Letta | Agent(s) | Notes |
|-------|-------------|----------|-------|----------|-------|
| code-method | ✅ | ✅ | ✅ | All | Universal operating loop |
| onemind-report | ✅ | ✅ | — | Legacy | Executive reporting (promoted) |
| onemind-status | ✅ | ✅ | — | Legacy | System health (promoted) |
| om-dashboard | ✅ | ✅ | — | Legacy | Dashboard (promoted) |
| om-openclaw | ✅ | ✅ | — | Legacy, Sovereignty | OpenClaw admin |
| (2 reserved) | — | — | — | — | TBD |

**Runtime Skills in OpenClaw:** 5

**Vault-Only Skills:** 2

---

## Per-Agent Skill Visibility (OpenClaw)

### Legacy (main) — Orchestrator
**Pillar:** SO + cross-domain  
**Skills:**
```
Core: om-openclaw, code-method, onemind-maintainer
Reporting: onemind-report, onemind-status, om-dashboard
Shared: grid, guardian, trinity
```
**Total:** 9 skills

### Zenith — HP Specialist
**Pillar:** HP 25-49  
**Skills:**
```
Core: spartan, oracle, identity, librarian, pulse, genesis
Shared: code-method
```
**Total:** 7 skills

### Unity — LE Specialist
**Pillar:** LE 50-74  
**Skills:**
```
Core: haven, helix, eden, nomad
Shared: code-method, homestead (promoted)
```
**Total:** 6 skills

### Eternal — GE Specialist
**Pillar:** GE 75-99  
**Skills:**
```
Core: empire, creed, empire-ops, video-production, finance, finance-ops, apex, legacy
Shared: code-method
```
**Total:** 9 skills

### Sovereignty — SO Specialist
**Pillar:** SO 00-24  
**Skills:**
```
Core: grid, guardian, guardian-ops, overwatch, eagleeye, trinity, onemind-maintainer
Shared: code-method, om-openclaw
```
**Total:** 9 skills

### Sentinel — Watchdog (Minimal)
**Pillar:** Failover + monitoring  
**Skills:**
```
Core: code-method
Fallback: om-openclaw (if core unreachable)
```
**Total:** 2 skills

---

## Implementation Checklist — OpenClaw Phase 2

- [ ] Create symlinks in `~/.openclaw/plugin-skills/` for all 33 runtime skills
- [ ] Verify symlinks point to vault originals (confirm: `/home/ubuntu/onemind/...`)
- [ ] Restore `skills: []` arrays in `openclaw.json` for all 6 agents (see backup)
- [ ] Add 4 promoted skills to arrays:
  - [ ] `homestead` → Unity
  - [ ] `onemind-report` → Legacy
  - [ ] `onemind-status` → Legacy
  - [ ] `om-dashboard` → Legacy
- [ ] Remove obsolete skill from plugin-skills if any (TBD after cleanup review)
- [ ] Restart gateway: `systemctl --user restart openclaw-gateway`
- [ ] Test: `openclaw skills list` shows 33 skills
- [ ] Test: `openclaw skills info <skill-name>` for each agent's assigned skills
- [ ] Test: Agents load skills on startup (check logs)
- [ ] Verify no duplicate/conflicting skill definitions

---

## Related

- `SKILL-LAYER-ARCHITECTURE.md` — System integration patterns + roadmap
- `../../02 Infrastructure/openclaw/plugin-skills/PLUGIN-SKILLS.md` — OpenClaw documentation
- Backup config: `/home/ubuntu/.openclaw/openclaw.json.bak-mcp-skills`
