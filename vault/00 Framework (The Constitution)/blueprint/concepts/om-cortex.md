---
title: "OM Cortex"
type: reference
domain: SO
status: pending
created: 2026-05-27
updated: 2026-05-27
tags: [concept, architecture, intelligence, cortex, onemind-os]
---

# OM Cortex

## Concept
The intelligence routing layer of OneMind OS. Not yet solidified — holding the name and direction.

## Working Definition
The "cortex" metaphor: the part of the brain responsible for higher-order processing — decision-making, routing, synthesis. Applied to OneMind OS: the layer that sits above raw agents and infrastructure and makes intelligent decisions about what runs, when, and why.

## Possible Roles
- Central decision hub for agent orchestration
- Context-aware routing (which agent handles what, based on mode/state)
- Intelligence aggregation layer — synthesizes signals from all pillars (SO/HP/LE/GE)
- The "thinking" layer between Zeus's intent and Legacy's execution

## Relationship to Existing Architecture
- Lives inside `SO/02 Intelligence & Agents`
- Could be the formalization of what Legacy already does informally
- Distinct from NATS (messaging) — Cortex is the brain, NATS is the nervous system

## Status
Concept only. No implementation. Name held.

## Next Steps
- [ ] Define boundaries: what does Cortex own vs. what does Legacy own?
- [ ] Map inputs/outputs
- [ ] Determine if this is a named component or a design pattern
- [ ] Explore relationship to OM Cortex → YouInverse → Interdependent (possible connected framework)

## Notes
<!-- Add architecture thoughts, diagrams, decisions here -->
