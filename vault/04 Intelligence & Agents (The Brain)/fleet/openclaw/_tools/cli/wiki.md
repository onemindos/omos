---
name: OpenClaw CLI — Wiki
updated: 2026-05-11
---

# `openclaw wiki`

Manages a memory-wiki vault — a compiled knowledge system with search, provenance tracking, and contradiction detection.

**Official docs:** [docs.openclaw.ai/cli/wiki](https://docs.openclaw.ai/cli/wiki.md)

---

## Commands

```bash
openclaw wiki status          # Check vault initialization and health
openclaw wiki doctor          # Run diagnostics on configuration

openclaw wiki search "<q>"    # Query with mode-specific ranking
openclaw wiki get <page>      # Retrieve a page with provenance details

openclaw wiki ingest <source> # Ingest notes or external sources
openclaw wiki compile         # Rebuild indexes and digests for machine consumption

openclaw wiki lint            # Detect structural issues, contradictions, stale content
openclaw wiki apply           # Apply mutations safely without manual editing
```

---

## Key Principle

> "Use `wiki search` + `wiki get` when provenance and page identity matter."

Use `wiki apply` for structured mutations over manual editing — maintains vault integrity.

---

## Integration Points

- **Memory plugins** — bridge import
- **Obsidian** — local vault management
- **Search backends** — shared or local

---

## Related

- [`memory.md`](memory.md) — semantic memory index (different system)
- [Official memory-wiki plugin](https://docs.openclaw.ai/plugins/memory-wiki.md)
