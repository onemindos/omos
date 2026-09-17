# Railway

> **Modern cloud platform for deploying full-stack apps and services**

## Overview

Railway is a developer-friendly PaaS for deploying apps, databases, and services with minimal configuration. Supports Postgres, Redis, MySQL, and custom Docker containers. Git-push-to-deploy workflow. Strong free tier and simple pricing.

## Status

| Attribute   | Value                              |
| ----------- | ---------------------------------- |
| Status      | 🔄 Evaluating                      |
| Category    | Cloud                              |
| Type        | PaaS / Deployment Platform         |
| Integration | CLI, GitHub, REST API              |

## Use Cases in OneMind

1. **App Deployment** — Deploy onemind-web or microservices without VPS management
2. **Database Hosting** — Managed Postgres/Redis as an alternative to self-hosted
3. **Service Orchestration** — Run multiple services with internal networking
4. **Staging Environments** — Separate staging from production (VPS)

## Key Features

- Git-push deployment (GitHub/GitLab)
- Managed Postgres, Redis, MySQL
- Custom Docker container support
- Internal service networking
- Environment variables management
- Usage-based pricing

## Quick Start

```bash
npm i -g @railway/cli
railway login
railway init
railway up
```

## Links

- Website: [railway.app](https://railway.app)
- Docs: [docs.railway.app](https://docs.railway.app)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
