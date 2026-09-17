---
title: "Context7"
type: reference
domain: SO
status: active
created: 2026-05-28
updated: 2026-05-28
tags: [platform, ai, mcp, documentation]
---

# Context7

> Up-to-date documentation and code examples for any programming library, delivered via MCP.

## Use Case in OneMind OS

- Real-time library docs lookup during coding sessions
- Code examples for any framework (Next.js, Flutter, FastAPI, etc.)
- Replaces searching docs manually — agent queries directly
- Active MCP server in Hermes config

## How It Works

1. **Resolve library** — search for a package name → get Context7 library ID
2. **Query docs** — ask a question with the library ID → get relevant docs + code snippets

## MCP Tools

| Tool | Purpose |
|------|---------|
| `mcp_context7_resolve_library_id` | Find the Context7 ID for a library |
| `mcp_context7_query_docs` | Get docs/examples for a specific library + question |

## Example

```
# Find Flutter library ID
resolve_library_id("Flutter", "how to build a navigation drawer")
→ /flutter/flutter

# Query docs
query_docs("/flutter/flutter", "how to build a navigation drawer")
→ Returns relevant docs + code snippets
```

## Integration

- Configured as MCP server in `~/.hermes/config.yaml`
- No API key needed (free)
- Rate limit: 3 calls per question
