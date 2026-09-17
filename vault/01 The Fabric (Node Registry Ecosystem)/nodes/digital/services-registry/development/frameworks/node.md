# Node.js

> **JavaScript runtime for server-side code, CLI tools, and build pipelines**

## Overview

Node.js is the server-side JavaScript runtime powering onemind-web (Next.js), CLI tools, build scripts, and npm-based tooling. The foundation of the entire TypeScript/JavaScript ecosystem in One Mind.

## Status

| Attribute   | Value                              |
| ----------- | ---------------------------------- |
| Status      | ✅ Active                           |
| Category    | Development / Runtime              |
| Type        | JavaScript Runtime                 |
| Integration | npm, package.json, CLI, N8N nodes  |

## Use Cases in OneMind

1. **onemind-web** — Next.js 16 app runs on Node.js
2. **CLI Tools** — `onemind` CLI entry point (`bin/onemind.js`)
3. **Build System** — npm scripts for build, test, lint, deploy
4. **N8N** — N8N itself runs on Node.js; custom nodes in JS/TS
5. **Quartz** — codex-live static site generator runs on Node.js

## Version

```bash
node --version    # v20+ LTS recommended
npm --version
# Use nvm to manage versions
nvm install 20 && nvm use 20
```

## Key npm Scripts (onemind-web)

```bash
npm run dev              # Dev server on port 3456
npm run build:ci         # Production build
npm run test:cli         # CLI tests
npm run lint:baseline    # Lint gate
```

## Links

- Website: [nodejs.org](https://nodejs.org)
- npm: [npmjs.com](https://npmjs.com)
- nvm: [github.com/nvm-sh/nvm](https://github.com/nvm-sh/nvm)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
