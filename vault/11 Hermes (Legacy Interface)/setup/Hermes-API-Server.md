---
title: Hermes API Server
domain: 11
pillar: SO
status: active
updated: 2026-05-21
tags: [hermes, api, openai-compatible, backend]
---

# Hermes API Server

The API server exposes Hermes as an OpenAI-compatible backend. Any chat UI that supports OpenAI can connect to it.

## Config (set 2026-05-21)

| Key | Value |
|-----|-------|
| Enabled | `true` |
| Port | `8080` |
| Local URL | `http://localhost:8080/v1` |
| Tailscale URL | `http://omos-mbp.tail717ba0.ts.net:8080/v1` |
| Tailscale IP | `http://100.98.178.23:8080/v1` |
| API Key | stored in `~/.hermes/config.yaml` → `api_server.key` |
| CORS | `*` (all origins allowed) |

Config was set via:
```bash
hermes config set api_server.enabled true
hermes config set api_server.port 8080
hermes config set api_server.key <token>
hermes config set api_server.cors_origins "*"
```

## How to Start

The API server starts automatically when Hermes launches with `--serve` flag or via the dashboard:

```bash
hermes serve
# or
hermes dashboard --host 127.0.0.1 --port 9119
```

> Note: The existing dashboard launchd plist uses `--tui --no-open --host 127.0.0.1 --port 9119 --insecure --skip-build`

## Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/health` | Healthcheck |
| GET | `/v1/capabilities` | List models and features |
| POST | `/v1/chat/completions` | Chat (OpenAI-compatible) |
| POST | `/v1/responses` | Responses API |
| POST | `/v1/runs` | Async agent runs |
| GET | `/v1/runs/{id}` | Poll run status |

## Auth

All requests require Bearer token:
```
Authorization: Bearer <api_server.key>
```

## Session Continuity

Pass these headers to maintain a persistent session across requests:
```
X-Hermes-Session-Id: <uuid>
X-Hermes-Session-Key: <session-key>
```

## Connecting a Chat UI

Point any OpenAI-compatible UI at:
- **Base URL:** `http://localhost:8080/v1`
- **API Key:** value from `~/.hermes/config.yaml` → `api_server.key`
- **Model:** `hermes` or leave as default

## Related

- [[Hermes-Overview]]
- [[Hermes-Desktop-Fathah]] — recommended primary GUI
- [[Hermes-Desktop-Dodo]] — power workbench GUI
