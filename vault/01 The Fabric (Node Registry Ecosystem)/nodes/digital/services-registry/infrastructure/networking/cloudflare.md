---
title: "Cloudflare"
type: reference
status: active
tags: [platform, infrastructure, dns, tunnels, cdn]
---
# Cloudflare

> DNS, tunnels, CDN, and Workers. Zero-trust access without open ports.

## Services Used
- **DNS** — All domains (zeusdelacruz.com, weareonemind.com)
- **Tunnels (cloudflared)** — Expose local services without port forwarding
- **Workers** — Edge compute (API routing, redirects)
- **Zero Trust** — Access policies for internal services

## CLI: cloudflared
```bash
cloudflared tunnel list
cloudflared tunnel run <tunnel-name>
cloudflared tunnel route dns <tunnel> <hostname>
```

## CLI: wrangler (Workers)
```bash
wrangler dev          # Local dev server
wrangler deploy       # Deploy to edge
wrangler tail         # Stream logs
wrangler kv:list      # KV namespace ops
```

## Cost
$0 (free tier covers all current usage)
