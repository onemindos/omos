# TypeScript

> **Typed superset of JavaScript — primary language of onemind-web**

## Overview

TypeScript adds static types to JavaScript. It's the default language for all new code in onemind-web. Strict mode enforced; `any` is forbidden. Type safety catches bugs at compile time and makes large codebases maintainable. All Next.js, React, and API route code is TypeScript.

## Status

| Attribute   | Value                                  |
| ----------- | -------------------------------------- |
| Status      | ✅ Active                               |
| Category    | Development / Language                 |
| Type        | Typed JavaScript Superset              |
| Integration | Node.js, Next.js, tsc, eslint, tsx     |

## Use Cases in OneMind

1. **onemind-web** — All source code in `src/` is TypeScript `.ts` / `.tsx`
2. **API Routes** — Type-safe HTTP handlers in Next.js
3. **Zustand Stores** — Typed state management for all app state
4. **Type Definitions** — `src/types/index.ts` for shared types
5. **CLI** — `bin/onemind.js` compiled from TypeScript

## Rules in onemind-web

```typescript
// ✅ Correct - use specific types
const agents: Agent[] = []
const data: Record<string, unknown> = {}

// ❌ Never - any defeats the purpose
const data: any = {}

// ✅ Use interfaces for object shapes
interface Session {
  id: string
  agentId: string
  createdAt: Date
}
```

## Config (`tsconfig.json`)

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true
  }
}
```

## Links

- Website: [typescriptlang.org](https://typescriptlang.org)
- Playground: [typescriptlang.org/play](https://typescriptlang.org/play)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
