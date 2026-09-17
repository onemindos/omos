---
skill: entity-system
type: verify
created: 2026-05-04
tags: [type/skill, zone/so, status/active]
---

# VERIFY: Entity System Skill

> Run these tests to confirm the skill is loaded correctly.

## Test 1 — Classification

Given this input:
> "I want to have $50,000 in my investment account by December 2026"

Expected output:
- Type: `goal`
- Pillar: `HP`
- Domain: `30 Finance (Money)`
- Metric: "$50,000 investment balance"
- Horizon: `Q4-2026`

## Test 2 — Routing

Given this input:
> "Notes on how to dispute a credit report error — 5 steps"

Expected output:
- Type: `knowledge`
- Pillar: `HP`
- Domain: `30 Finance (Money)` → `Credit (Credit Profile)`
- Depth: `3`

## Test 3 — Cross-Pillar

Given this input:
> "Weekly date night with my wife — every Friday 7pm"

Expected output:
- Type: `calendar`
- Pillar: `LE`
- Domain: `51 Marriage`
- Recurrence: `weekly`
- Event type: `ritual`

## Test 4 — Inbox routing

Given this input:
> "Some random idea I had about combining fitness tracking with finance metrics"

Expected output:
- Type: `knowledge` or `goal` (either valid)
- Status: `planned`
- Routed to: `99 Inbox (Queue) (The Dock)/`
- ai_summary: describes the idea clearly

## Test 5 — Frontmatter generation

For any of the above, verify the agent generates:
- [ ] `id` field with correct format
- [ ] `type` field matching entity type
- [ ] `pillar` field (SO/HP/LE/GE)
- [ ] `domain` field with number and name
- [ ] `ai_summary` — one sentence, clear
- [ ] `tags` with type/, zone/, domain/, status/ namespaces
- [ ] `depth` set correctly
- [ ] `related` field present (even if empty array)

All 5 tests passing = skill verified.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
