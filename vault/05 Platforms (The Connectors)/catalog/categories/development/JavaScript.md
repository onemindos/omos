# JavaScript

> **The language of the web — client-side scripting and Node.js runtime**

## Overview

JavaScript is the universal web language. In One Mind, it powers onemind-web client-side code, Next.js API routes, N8N custom function nodes, Quartz (codex-live) configuration, and all browser interactions. TypeScript is preferred but JavaScript remains foundational.

## Status

| Attribute   | Value                               |
| ----------- | ----------------------------------- |
| Status      | ✅ Active                            |
| Category    | Development / Language              |
| Type        | Programming Language                |
| Integration | Browser, Node.js, V8, Deno, Bun     |

## Use Cases in OneMind

1. **onemind-web** — Client-side React components (TSX/JSX compiles to JS)
2. **N8N Custom Code** — Function nodes in N8N use JS/TS
3. **Quartz Config** — `quartz.config.ts` / `quartz.layout.ts` in codex-live
4. **CLI Tools** — `bin/onemind.js` CLI entry point
5. **Build Scripts** — `scripts/` directory in onemind-web

## Key Patterns in Stack

```javascript
// ES Module syntax (used throughout)
import { something } from './module.js'
export const handler = async (req, res) => { ... }

// Async/await pattern
const result = await fetch('/api/agents')
const agents = await result.json()
```

## Note

TypeScript is strongly preferred over plain JS for all new code in onemind-web. Plain JS is used for legacy configs and scripts that don't need type safety overhead.

## Links

- MDN: [developer.mozilla.org/en-US/docs/Web/JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript)
- ECMAScript: [tc39.es](https://tc39.es)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
