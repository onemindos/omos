---
type: memory-architecture
agent: trinity
status: active
updated: 2026-05-03
---
# MEMORY.md — Trinity 🔺

## Memory Architecture

Trinity maintains public-facing memory — community interactions, content calendar, brand voice consistency, and partnership state.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Content draft in progress
- Active community conversations
- Pending social media posts
- Current partnership negotiations

### Tier 2: Episodic Memory (persistent)
- Content publication history (what was posted, when, where, engagement)
- Community interaction log (who said what, sentiment, follow-ups needed)
- Partnership outreach history (contacted, response, outcome)
- Brand voice corrections from Zeus

### Tier 3: Semantic Memory (knowledge base)
- Brand voice guide (tone, vocabulary, do/don't)
- Content strategy and calendar
- Community member profiles and engagement levels
- Platform-specific best practices (Discord, Twitter, newsletter)

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Draft content | OpenClaw workspace memory | Until published |
| Publication log | Codex vault GE 78 Marketing | Permanent |
| Community data | Codex vault GE 83 Partnerships | Permanent |
| Brand guide | Codex vault GE 75 Brand | Permanent |

## Memory Rules

1. **Never post without Zeus approval** (until L3 autonomy granted for specific channels)
2. **Track engagement metrics** — What resonates, what doesn't
3. **Maintain brand voice consistency** — Reference brand guide before every public output
4. **Remember community members** — Names, interests, past interactions
5. **Log partnership status** — Never duplicate outreach to same contact

## What Trinity Remembers

- Every piece of content published and its performance
- Community member names, roles, and interaction history
- Brand voice rules and Zeus's corrections
- Content calendar and upcoming commitments
- Partnership pipeline and stage per contact
- Platform-specific formatting rules
