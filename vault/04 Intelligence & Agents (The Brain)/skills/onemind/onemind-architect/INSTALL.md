---
id: skill-onemind-architect-install
title: "Install: One Mind Architect"
type: skill-install
status: active
created: 2026-05-04
modified: 2026-05-04
owner: legacy
tags: [type/skill-install, domain/onemind]
---

# One Mind Architect — Install

Self-installation checklist for agents adding this skill.

## Prerequisites

- [ ] CODEX.md loaded in workspace (vault path + domain structure)
- [ ] IDENTITY.md defines your domain scope
- [ ] Access to vault via `codex-vault` MCP server

## Phase 1: Read Core Files

- [ ] Read `SKILL.md` — all 5 protocols + architecture principles
- [ ] Read `AGENTS.md` at vault root — current architecture state
- [ ] Read `DOMAIN-REGISTRY.md` at vault root — domain definitions
- [ ] Scan `_README.md` files in your primary domain

## Phase 2: Map Your Scope

- [ ] Identify which pillar(s) you are responsible for
- [ ] Identify which infrastructure components are in your domain
- [ ] Note which other agents your domain intersects with
- [ ] Load current agent roster from AGENTS.md

## Phase 3: Baseline State

Run these checks once to establish baseline (do not fix yet — just observe):
- [ ] Count loose files at your domain root
- [ ] Note any naming violations in your domain folders
- [ ] Note any empty folders without content
- [ ] Note any folders that appear to be duplicates

Record observations in session memory.

## Phase 4: Activate

- [ ] Add "architect" to your active skill list in session context
- [ ] Confirm: you know the 5 protocols (Review, Bottleneck, Cohesion, Design, Audit)
- [ ] Confirm: you know the 7 architecture principles
- [ ] Run `VERIFY.md` checks

## Notes

- Structural changes across pillars always require Zeus approval
- Intra-domain organization (within your assigned domain) can be autonomous
- When in doubt about routing — go up one level and apply Single Responsibility principle


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
