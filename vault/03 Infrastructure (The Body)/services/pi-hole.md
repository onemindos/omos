---
title: "Pi-hole"
type: service-card
status: planned
created: 2026-06-08
updated: 2026-06-08
owner: zeus
category: infrastructure
tags: [infrastructure, dns, ad-blocking, privacy, network, self-hosted, fleet]
ai_summary: "Network-wide DNS ad blocker. Runs as a DNS server on the LAN — blocks ads, trackers, and malware domains before requests leave the network. Covers every device without per-device installs."
---

# Pi-hole

> Network-wide DNS ad blocker. One instance per location blocks ads and trackers for every device on that network.

## What It Does

Acts as the DNS server for your entire network. When any device makes a DNS lookup, Pi-hole checks the domain against blocklists. Ads, trackers, telemetry, and malware domains are blocked at the DNS level — the request never even goes out. No app installs on phones, TVs, or IoT devices required.

## Why OneMind Needs It

3 locations (House / Barn / Office) + a growing fleet of nodes and IoT devices. Pi-hole installed at each location router covers everything on that subnet — DIMOS, OM1, smart home devices, phones, laptops — without touching any individual device.

## Key Features

| Feature | Description |
|---|---|
| **DNS sinkhole** | Blocks at DNS level — no ads ever load |
| **Blocklists** | Community-maintained lists covering millions of domains |
| **Per-device stats** | See exactly which device is making what DNS queries |
| **Allowlisting** | Whitelist domains that get over-blocked |
| **DHCP server** | Can also act as DHCP if needed |
| **Web admin UI** | Dashboard with query logs, stats, blocklist management |
| **API** | REST API for automation and integration |

## Deployment Plan

| Location | Node | Status |
|---|---|---|
| House | TBD | Planned |
| Barn | TBD | Planned |
| Office | TBD | Planned |

> Use [[nebula-sync]] to keep all three instances in sync from a single primary.

## Links

- **Website:** https://pi-hole.net
- **Docs:** https://docs.pi-hole.net
- **GitHub:** https://github.com/pi-hole/pi-hole
- **Docker Hub:** https://hub.docker.com/r/pihole/pihole

## Related Services

- [[nebula-sync]] — keeps multiple Pi-hole instances in sync (required for multi-location)
- [[traefik]] — can route Pi-hole admin UI through Traefik with HTTPS
