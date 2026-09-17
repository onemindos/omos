# Context7 MCP

**Type:** Library documentation resolver
**MCP Package:** `@upstash/context7-mcp` (official)
**Website:** context7.com
**Auth:** None required (free) — optional API key for higher limits

## What It Does

Context7 fetches **current, version-accurate library documentation** and injects it into the LLM's context before generating code. Solves the hallucination problem where AI generates code for outdated APIs.

Without Context7: Agent guesses based on training data (may be 1-2 years stale).
With Context7: Agent fetches live docs for the exact library version you're using.

## How It Works

1. Agent calls `resolve-library-id` with library name (e.g., "livekit-agents")
2. Gets back a Context7-compatible library ID
3. Agent calls `get-library-docs` with that ID + a topic
4. Gets back current, focused documentation snippets
5. Uses those snippets to generate accurate code

## Key Tools Exposed

| Tool | Description |
|------|-------------|
| `resolve-library-id` | Find the Context7 ID for a library name |
| `get-library-docs` | Fetch current docs for a library + topic |

## Config

```json
{
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp"]
    }
  }
}
```

No API key needed for basic use.

## Libraries Most Useful for This Stack

- `livekit-agents` — Voice agent framework
- `livekit` — Core LiveKit SDK
- `openai` — Python/JS OpenAI SDK
- `anthropic` — Anthropic SDK
- `nextjs` — Next.js App Router
- `prisma` — Database ORM
- `langchain` — Agent/chain framework

## Stack Use Cases

- **Forge agent:** Coding assistant always uses current API docs — no stale method signatures
- **Grid agent:** Infrastructure code (Terraform, Docker) uses latest provider syntax
- **Legacy agent:** Answers "how do I use X library" with live, accurate examples
- **Any coding task:** Preflight check before generating code for any external library


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
