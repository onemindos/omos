---
id: so-legacy-access-20260615
title: "Legacy — Device Access Guide"
type: reference
status: active
created: 2026-06-15
modified: 2026-06-15
owner: zeus
tags: [type/reference, status/active, domain/so, ai/legacy, infra/openclaw]
---

# Legacy — Device Access Guide

> How to reach Legacy (OpenClaw gateway) from any device on Tailscale.

---

## The Stack

```
iPhone / iPad / any device
        │  HTTPS (Tailscale cert — valid secure context)
        ▼
omos-mbp.tail717ba0.ts.net:443
        │  Tailscale Serve proxy
        ▼
127.0.0.1:18789  ← OpenClaw gateway (loopback only)
        │
        ▼
LiteLLM :4000 → Bedrock / DigitalOcean / Ollama
```

---

## URLs

| Device | URL | Notes |
|--------|-----|-------|
| Mac (local) | `http://127.0.0.1:18789` | Always works, localhost = secure context |
| iPhone / iPad / any Tailscale device | `https://omos-mbp.tail717ba0.ts.net` | HTTPS required — use this, nothing else |
| LiteLLM UI | `http://localhost:4000/ui` | Mac only |

> **Never use** `http://100.98.178.23:18789` — plain HTTP on a raw IP blocks WebCrypto → device identity fails → endless approval prompts.

---

## Auth Token

The dashboard auto-connects if you're on the Tailscale network. If prompted:

```bash
cat ~/.openclaw/openclaw.json | python3 -c "import sys,json; c=json.load(sys.stdin); print(c['gateway']['auth']['token'])"
```

Append to URL if needed: `https://omos-mbp.tail717ba0.ts.net/#token=<token>`

---

## Architecture

| Service | Host | Tailscale | Port |
|---------|------|-----------|------|
| OpenClaw (Legacy) | Mac | `omos-mbp.tail717ba0.ts.net` | 18789 (loopback) |
| LiteLLM | Mac | `100.98.178.23` | 4000 |
| Hermes | Mac | `100.98.178.23` | 8642 |
| Ollama | Mac | `100.98.178.23` | 11434 |

---

## Why Loopback + Tailscale Serve (not lan bind)

- `bind=lan` (0.0.0.0) exposes the raw port on the Tailscale IP → plain HTTP → browser blocks WebCrypto → no device identity
- `bind=loopback` + Tailscale Serve = HTTPS with a Let's Encrypt cert → secure context → device identity works → no approval loops
- Gateway never exposed to the open internet — Tailscale is tailnet-only

---

## Config Summary

```
gateway.bind = loopback
gateway.tailscale.mode = serve
gateway.controlUi.allowedOrigins = ["https://omos-mbp.tail717ba0.ts.net", "http://127.0.0.1:18789"]
gateway.controlUi.dangerouslyAllowHostHeaderOriginFallback = false
```

---

## Tailscale Serve Routes (active)

```
https://omos-mbp.tail717ba0.ts.net  →  http://127.0.0.1:18789  (OpenClaw)
https://onemind-mac.tail717ba0.ts.net  →  http://127.0.0.1:18789  (alias)
https://oneminds-macbook-pro.tail717ba0.ts.net  →  http://127.0.0.1:18789  (alias)
```

Only `omos-mbp.tail717ba0.ts.net` resolves correctly from all devices. Use that one.

---

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| "control ui requires device identity" | Using HTTP instead of HTTPS | Use `https://omos-mbp.tail717ba0.ts.net` |
| Page doesn't load | Tailscale not connected on device | Connect Tailscale first |
| Keeps asking to approve device | Was connecting via plain HTTP, pending requests queued | Connect via HTTPS — one-time approval only |
| Gateway down | LaunchAgent crashed | `openclaw gateway restart` |
| Mac app can't connect | Use `http://127.0.0.1:18789` on Mac only | Never use LAN IP for Mac app |
