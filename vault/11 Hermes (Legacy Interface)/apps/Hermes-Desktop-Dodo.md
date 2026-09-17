---
title: HermesDesktop Workbench (dodo-reach/hermes-desktop)
domain: 11
pillar: SO
status: installed-needs-ssh
updated: 2026-05-21
tags: [hermes, gui, native, macos, workbench, kanban, ssh]
---

# HermesDesktop Workbench — Power GUI

> **Best for:** Power workflows — Kanban, multi-agent orchestration, file editing, usage analytics, terminal. Complements the fathah chat UI.

## What It Is

`dodo-reach/hermes-desktop` is a native macOS app (Swift/SwiftUI, macOS 14+) for managing Hermes instances over SSH. Unlike fathah, it is NOT chat-focused — it's a management workbench.

- GitHub: https://github.com/dodo-reach/hermes-desktop
- Version installed: 0.9.0
- Installed at: `/Applications/HermesDesktop.app`
- Platform: macOS 14+ native (arm64)

## Features

| Feature | Description |
|---------|-------------|
| Kanban | Visual task board for multi-agent workflows |
| Workflows | Build and run agent workflow chains |
| Files | Edit Hermes config, skills, memory, soul files |
| Usage | Token usage analytics and cost tracking |
| Terminal | Embedded terminal connected to Hermes host |
| SSH | Manages remote AND local Hermes instances |

## How to Connect (Local)

This app connects via SSH — even for localhost. Steps:

### 1. Enable Remote Login on Mac
```
System Settings → General → Sharing → Remote Login → ON
```
Allow access for: your user (onemind) or All users

### 2. Test SSH
```bash
ssh localhost echo "ok"
```

### 3. Add Local Connection in HermesDesktop
- Host: `omos-mbp.tail717ba0.ts.net` (or `localhost` if on-device)
- Port: `22`
- User: `onemind`
- Auth: SSH key (recommended) or password
- Hermes path: `/Users/onemind/.hermes`

### SSH Key Setup (if needed)
```bash
# Generate key if you don't have one
ssh-keygen -t ed25519 -C "hermes-desktop-local"
# Authorize it for localhost
cat ~/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

## Quick Launch

```bash
open "/Applications/HermesDesktop.app"
```

## Current Status

- App installed ✓
- SSH to localhost: NOT YET ENABLED
- Action needed: Enable Remote Login in System Settings

## Notes

- macOS 14 (Sonoma) or later required — confirmed compatible with macOS 26.5
- Does NOT use the Hermes API server — communicates via SSH directly
- Can also manage remote VPS Hermes instances (add additional SSH connections)

## Related

- [[Hermes-Overview]]
- [[Hermes-API-Server]]
- [[Hermes-Desktop-Fathah]] — primary day-to-day chat GUI
