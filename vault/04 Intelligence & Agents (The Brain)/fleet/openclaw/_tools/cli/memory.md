---
name: OpenClaw CLI — Memory
updated: 2026-05-10
---

# `openclaw memory`

Manages semantic memory indexing, search, and consolidation via the active memory plugin (`memory-core` by default).

**Official docs:** [docs.openclaw.ai/cli/memory](https://docs.openclaw.ai/cli/memory.md)

---

## Commands

### Status & Maintenance
```bash
openclaw memory status           # Memory system health
openclaw memory status --deep    # Probe vector stores + embedding providers
openclaw memory status --fix     # Repair stale locks, normalize metadata
openclaw memory index            # Reindex memory
openclaw memory index --force    # Full reindex
```

### Search
```bash
openclaw memory search "<query>"                      # Semantic search
openclaw memory search "<query>" --max-results 10     # Limit results
openclaw memory search "<query>" --min-score 0.7      # Minimum relevance threshold
```

### Consolidation (Dreaming)
```bash
openclaw memory promote                      # Preview short-term → MEMORY.md candidates
openclaw memory promote --apply              # Write promotions to MEMORY.md
openclaw memory promote-explain "<selector>" # Scoring breakdown for a candidate
openclaw memory rem-harness                  # Preview REM reflections (no write)
```

---

## Dreaming System

Background consolidation runs in three phases:

| Phase | What it does |
|---|---|
| **Light** | Sort and surface short-term memories |
| **REM** | Reflect and synthesize patterns |
| **Deep** | Promote to `MEMORY.md` |

Enable in config: `"dreaming.enabled": true`
Default sweep: `0 3 * * *` (3 AM daily)

### Promotion Scoring Factors
- Recall frequency
- Query relevance
- Query diversity
- Recency
- Cross-day consolidation
- Conceptual richness

Default deep thresholds: `minScore=0.8`, `minRecallCount=3`, `minUniqueQueries=3`

---

## Related

- [concepts/ACTIVE-MEMORY-PLUGIN.md](../concepts/ACTIVE-MEMORY-PLUGIN.md)
- [concepts/NATIVE-MEMORY-ARCHITECTURE.md](../concepts/NATIVE-MEMORY-ARCHITECTURE.md)
- [Official memory overview](https://docs.openclaw.ai/concepts/memory.md)
