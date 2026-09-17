# Postman

> **API development, testing, and documentation platform**

## Overview

Postman is the industry-standard tool for building, testing, and documenting APIs. Used for exploring and debugging all OneMind service APIs — Todoist, Letta, OpenClaw, home automation, and internal VPS services.

## Status

| Attribute   | Value                           |
| ----------- | ------------------------------- |
| Status      | ✅ Active                        |
| Category    | Development                     |
| Type        | API Testing / Development       |
| Integration | Desktop App, CLI, REST          |

## Use Cases in OneMind

1. **API Exploration** — Test Todoist, Letta, Zoho, and all platform APIs
2. **Debugging** — Reproduce and isolate integration issues
3. **Collections** — Maintain saved request collections for each platform
4. **Documentation** — Auto-generate API docs from collections
5. **Newman CLI** — Run API tests in CI pipelines

## Key Features

- Collections (organized API requests)
- Environments (swap dev/prod variables)
- Automated testing (JavaScript assertions)
- Mock servers
- Newman CLI for automation
- OpenAPI import/export

## Quick Start

```bash
# CLI runner
npm install -g newman
newman run collection.json --environment env.json
```

## Links

- Website: [postman.com](https://postman.com)
- Docs: [learning.postman.com](https://learning.postman.com)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
