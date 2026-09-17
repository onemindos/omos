---
title: "CLI — wrangler (Cloudflare CLI)"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/cloudflare, topic/deploy]
ai_summary: "Cloudflare Wrangler CLI quick reference — Pages deploy for codex-live (codex.onemindos.com), Workers, KV storage, and D1 databases."
---

# wrangler — Cloudflare CLI

> Deploy and manage Cloudflare Pages, Workers, KV, and D1. `codex-live` deploys here.
> Install: `npm install -g wrangler`

---

## Auth

```bash
# Login (browser OAuth)
wrangler login

# Check auth status
wrangler whoami

# Logout
wrangler logout
```

---

## Pages (codex-live)

The `onemind-codex-live` Quartz site deploys to Cloudflare Pages.

```bash
# Deploy to Cloudflare Pages
wrangler pages deploy ./public --project-name=onemind-codex

# Build first, then deploy
npx quartz build -d codex && wrangler pages deploy ./public --project-name=onemind-codex

# List Pages projects
wrangler pages project list

# View deployment history
wrangler pages deployment list --project-name=onemind-codex
```

Config: [wrangler.toml](../../../../../onemind-codex-live/wrangler.toml)

Site: [codex.onemindos.com](https://codex.onemindos.com)

---

## Workers

```bash
# Deploy a worker
wrangler deploy

# Run worker locally (dev mode)
wrangler dev

# View worker logs (live tail)
wrangler tail

# List workers
wrangler workers list

# Delete a worker
wrangler delete --name worker-name
```

---

## KV Storage

```bash
# List KV namespaces
wrangler kv:namespace list

# Create namespace
wrangler kv:namespace create MY_NAMESPACE

# Put a value
wrangler kv:key put --namespace-id=<id> "key" "value"

# Get a value
wrangler kv:key get --namespace-id=<id> "key"

# List keys
wrangler kv:key list --namespace-id=<id>
```

---

## D1 (SQLite Database)

```bash
# List D1 databases
wrangler d1 list

# Create database
wrangler d1 create my-db

# Execute SQL
wrangler d1 execute my-db --command "SELECT * FROM table_name;"

# Execute from SQL file
wrangler d1 execute my-db --file=./schema.sql
```

---

## R2 (Object Storage)

```bash
# List buckets
wrangler r2 bucket list

# Create bucket
wrangler r2 bucket create my-bucket

# Upload file
wrangler r2 object put my-bucket/path/file.txt --file ./local-file.txt

# Download file
wrangler r2 object get my-bucket/path/file.txt --file ./local-file.txt
```

---

## Cloudflare Tunnel (cloudflared)

```bash
# Login
cloudflared tunnel login

# Create tunnel
cloudflared tunnel create onemind-tunnel

# Run tunnel
cloudflared tunnel run onemind-tunnel

# List tunnels
cloudflared tunnel list

# Route DNS to tunnel
cloudflared tunnel route dns onemind-tunnel subdomain.onemindos.com
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
