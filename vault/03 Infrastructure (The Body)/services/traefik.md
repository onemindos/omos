---
title: "Traefik"
type: service-card
status: planned
created: 2026-06-08
updated: 2026-06-08
owner: zeus
category: infrastructure
tags: [infrastructure, reverse-proxy, ssl, docker, networking, self-hosted]
ai_summary: "Reverse proxy and automatic SSL/TLS manager. Routes traffic by domain/subdomain to the right container, auto-provisions Let's Encrypt certs. Built for Docker — discovers containers automatically."
---

# Traefik

> Reverse proxy + automatic SSL. The routing layer for all self-hosted services on OneMind OS.

## What It Does

Sits in front of all self-hosted services and routes incoming requests by domain or subdomain to the correct container. Auto-provisions and renews TLS certificates via Let's Encrypt. Built for Docker — discovers containers automatically via labels, no manual config per service.

Instead of:
```
http://192.168.1.10:3000   → Langfuse
http://192.168.1.10:3111   → agentmemory
http://192.168.1.10:3010   → Homepage
```

You get:
```
https://langfuse.onemind.local
https://memory.onemind.local
https://home.onemind.local
```

## Why OneMind Needs It

Currently 8+ services run on raw IP:port. Traefik gives them:
- Clean domain names (internal or external)
- HTTPS everywhere via Let's Encrypt
- Single entry point for all traffic
- Automatic cert renewal

## Key Features

| Feature | Description |
|---|---|
| **Auto-discovery** | Reads Docker labels — new container = auto-registered |
| **Let's Encrypt** | Automatic TLS cert provisioning + renewal |
| **Dashboard** | Built-in web UI showing all routes and services |
| **Middleware** | Auth, rate limiting, redirects, headers |
| **TCP/UDP** | Routes non-HTTP traffic too (NATS, databases) |
| **Multiple providers** | Docker, Kubernetes, file-based configs |

## Deployment Target

- `omos-mbp` — primary node, routes all local services
- Each fleet node gets its own Traefik instance as the fleet grows

## Links

- **Website:** https://traefik.io
- **Docs:** https://doc.traefik.io/traefik/
- **GitHub:** https://github.com/traefik/traefik
- **Docker Hub:** https://hub.docker.com/_/traefik

## Related Services

- Works with all services in `SO/03/services/`
- Pairs with Tailscale for external access without port forwarding
- [[pi-hole]] — can route Pi-hole admin UI through Traefik
