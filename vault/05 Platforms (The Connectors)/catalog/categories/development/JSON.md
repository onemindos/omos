# JSON

> **JavaScript Object Notation — universal data interchange format**

## Overview

JSON is the standard data format for all API communication, configuration files, and structured data in the One Mind stack. Every REST API returns JSON; all package configs are JSON; Todoist, OpenClaw, N8N, and all platform integrations speak JSON.

## Status

| Attribute   | Value                                    |
| ----------- | ---------------------------------------- |
| Status      | ✅ Active                                 |
| Category    | Development / Data Format               |
| Type        | Data Interchange Format                  |
| Integration | REST APIs, config files, databases, CLI  |

## Use Cases in OneMind

1. **API Communication** — All REST API payloads (Todoist, Anthropic, Letta, etc.)
2. **Configuration** — `package.json`, `tsconfig.json`, `docker-compose` override files
3. **Agent Responses** — OpenClaw agent tool call results are JSON
4. **Storage** — better-sqlite3 stores JSON fields in onemind-web
5. **Vault Frontmatter** — `ecosystem-manifest.json` at codex root

## Key Config Files

| File | Purpose |
| ---- | ------- |
| `package.json` | npm deps + scripts |
| `tsconfig.json` | TypeScript compiler config |
| `components.json` | shadcn/ui components config |
| `ecosystem-manifest.json` | Machine-readable repo map |

## Quick Reference

```json
{
  "name": "onemind",
  "version": "1.0.0",
  "domains": ["SO", "HP", "LE", "GE"],
  "active": true
}
```

## JSON5 / JSONC Note

Some config files (`tsconfig.json`) use JSONC (JSON with comments). Standard `JSON.parse()` won't handle these — use a JSONC parser or strip comments first.

## Links

- Spec: [json.org](https://json.org)
- MDN: [developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
