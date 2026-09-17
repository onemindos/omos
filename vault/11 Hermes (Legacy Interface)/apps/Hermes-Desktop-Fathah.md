---
title: Hermes Agent App (fathah/hermes-desktop)
domain: 11
pillar: SO
status: installed
updated: 2026-05-21
tags: [hermes, gui, electron, desktop, chat]
---

# Hermes Agent App — Primary GUI

> **Best for:** Day-to-day chat with Legacy AI, managing skills/memory/tools, visual interface for everything Hermes can do.

## What It Is

`fathah/hermes-desktop` is an Electron app purpose-built for Hermes Agent. It is NOT a generic OpenAI chat wrapper — it was built specifically to expose every Hermes feature through a polished UI.

- GitHub: https://github.com/fathah/hermes-desktop
- Version installed: 0.4.5
- Installed at: `/Applications/Hermes Agent.app`
- Platform: macOS arm64 (M-series)

## Features

| Feature | Description |
|---------|-------------|
| Chat | Full multi-session chat with streaming |
| Skills | Browse, load, create, edit skills visually |
| Memory | View and edit persistent memory |
| Soul Editor | Edit agent identity/soul files |
| Tools | Toggle tools on/off per session |
| Cron Jobs | View and manage scheduled tasks |
| Gateways | Manage Telegram/Discord/etc integrations |
| Claw3d | 3D agent visualization |

## How to Connect

1. Ensure Hermes API server is running (`hermes serve` or dashboard launchd service)
2. Open `/Applications/Hermes Agent.app`
3. On first launch it auto-detects `~/.hermes` — no config needed
3. If prompted for API URL: `http://omos-mbp.tail717ba0.ts.net:8080/v1` (or `http://localhost:8080/v1` if on-device)
5. If prompted for API key: use value from `~/.hermes/config.yaml` → `api_server.key`

## Quick Launch

```bash
open "/Applications/Hermes Agent.app"
```

## Notes

- Auto-detects existing `~/.hermes` install on first launch
- Does NOT need SSH — connects via local API server
- Updates via in-app updater or re-download from GitHub releases

## Related

- [[Hermes-Overview]]
- [[Hermes-API-Server]]
- [[Hermes-Desktop-Dodo]] — complementary power tool
