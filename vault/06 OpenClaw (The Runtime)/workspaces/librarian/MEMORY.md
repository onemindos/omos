---
type: memory-architecture
agent: librarian
status: active
updated: 2026-05-03
---
# MEMORY.md — Librarian 🔍

## Memory Architecture

Librarian maintains the intelligence pipeline memory — source feeds, research state, knowledge gaps, and brief generation history.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Active research query and findings so far
- Current news/feed digest in progress
- Pending briefs awaiting delivery
- Sources being evaluated for quality

### Tier 2: Episodic Memory (persistent)
- Research query history (what was asked, what was found)
- Brief delivery log (topic, date, reception from Zeus)
- Source reliability scores (which feeds are high-signal)
- Knowledge gap inventory (what Zeus has asked about that couldn't be found)

### Tier 3: Semantic Memory (knowledge base)
- Curated source list per domain (news, research, industry)
- Zeus's information interests and priority topics
- RAG index state (what's been embedded, what's pending)
- Knowledge taxonomy (how information connects across domains)

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Active research | OpenClaw workspace memory | Until delivered |
| Brief history | Codex vault SO 06 Research | Permanent |
| Source index | Codex vault SO 01 Intelligence | Permanent |
| RAG state | Vector DB + workspace | Permanent |

## Memory Rules

1. **Track source quality** — Not all feeds are equal; score them by signal-to-noise
2. **Remember Zeus's interests** — If he asks about a topic twice, it's a standing interest
3. **Avoid duplicate briefs** — Don't resurface information already delivered
4. **Maintain knowledge freshness** — Flag when stored knowledge is outdated
5. **Connect the dots** — Cross-reference new intel with existing knowledge

## What Librarian Remembers

- All research queries and their outcomes
- Zeus's standing interests (topics he always wants updates on)
- Source reliability scores per feed/publication
- RAG embedding coverage (% of vault indexed)
- Knowledge gaps (questions that couldn't be answered)
- Brief format preferences (length, depth, frequency)
- Cross-domain connections discovered during research
