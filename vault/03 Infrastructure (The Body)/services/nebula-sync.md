---
title: "Nebula Sync"
type: service-card
status: planned
created: 2026-06-08
updated: 2026-06-08
owner: zeus
category: infrastructure
tags: [infrastructure, pi-hole, sync, dns, multi-location, self-hosted, fleet]
ai_summary: "Pi-hole config sync tool. Keeps multiple Pi-hole v6 instances synchronized from a primary node — blocklists, groups, clients, settings. Required for consistent DNS blocking across all OneMind locations."
---

# Nebula Sync

> Syncs Pi-hole configuration across multiple instances. One primary node = all locations stay consistent.

## What It Does

Community tool (not official Pi-hole) that syncs Pi-hole v6 configuration from a designated primary instance to all secondary instances. Blocklists, groups, clients, and settings stay identical across every location automatically.

## Why OneMind Needs It

Pi-hole deployed at House / Barn / Office = 3 separate instances. Without Nebula Sync, managing blocklists and settings on each one manually is error-prone. With Nebula Sync, you update the primary once — every other location syncs automatically.

## How It Works

- Designate one Pi-hole as **primary** (source of truth)
- All others are **secondaries**
- Run Nebula Sync as a Docker container or standalone app
- Configure sync interval (scheduled) or trigger manually (one-off)
- Configured via `nebula-sync.yaml` or environment variables

## What Gets Synced

- Adlists (blocklists)
- Groups
- Clients
- Settings

## Deployment

Runs alongside Pi-hole on the primary node, or as a standalone container anywhere with network access to all Pi-hole instances.

> **Requires Pi-hole v6** — does not work with v5.

## Links

- **GitHub:** https://github.com/lovelaced/nebula-sync
- **Pi-hole Blog Post:** https://pi-hole.net/2023/08/24/nebula-sync/
- **Docker Compose Example:** https://github.com/lovelaced/nebula-sync/blob/main/docker-compose.yml
- **Config Reference:** https://github.com/lovelaced/nebula-sync/blob/main/nebula-sync.yaml

## Related Services

- [[pi-hole]] — required dependency
