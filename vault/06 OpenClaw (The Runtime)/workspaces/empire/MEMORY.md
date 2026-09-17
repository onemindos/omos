---
type: memory-architecture
agent: empire
status: active
updated: 2026-05-03
---
# MEMORY.md — Empire 👑

## Memory Architecture

Empire maintains business memory — revenue tracking, product pipeline, deal flow, strategic decisions, and venture performance.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Active deals in negotiation
- Current product development sprint
- Pending revenue-impacting decisions
- In-progress financial analysis

### Tier 2: Episodic Memory (persistent)
- Revenue history (MRR, ARR, by product, by channel)
- Deal log (lead → opportunity → close/loss, with reasons)
- Product launch history (what shipped, reception, revenue impact)
- Strategic decisions and their outcomes
- Investment/capital events

### Tier 3: Semantic Memory (knowledge base)
- Business model architecture (revenue streams, cost structure)
- Product portfolio and roadmap
- Market positioning and competitive landscape
- Pricing strategy and elasticity data
- Partnership terms and obligations
- Brand guidelines and IP portfolio

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Active deals | OpenClaw workspace + CRM | Until closed |
| Revenue data | Codex vault GE 80 Treasury | Permanent |
| Product roadmap | Codex vault GE 77 Products | Permanent |
| Strategy docs | Codex vault GE 81 Strategy | Permanent |

## Memory Rules

1. **Track every dollar** — Revenue in, costs out, margins calculated
2. **Never lose a lead** — Every business conversation is a potential deal
3. **Measure product performance** — Launch → adoption → revenue → retention
4. **Learn from losses** — When deals fail, record why (pricing, timing, competition, fit)
5. **Strategic consistency** — Decisions should align with long-term thesis

## What Empire Remembers

- Monthly/quarterly/annual revenue by product and channel
- Complete deal pipeline and conversion rates
- Product performance metrics post-launch
- Strategic decisions and their ROI
- Competitive intelligence and market shifts
- Partnership agreements and their value
- Capital structure and funding history
- Zeus's business principles and red lines
