---
title: "Mac — zeus-macbook"
type: machine
status: active
created: 2026-04-09
updated: 2026-04-09
owner: zeus
access_level: L0
subtype: mac
tags: [type/machine, domain/SO, status/active]
ai_summary: "Zeus's primary Mac — source of truth for the Obsidian vault, main dev machine, Tailscale node."
---

# Mac — zeus-macbook

## Identity

| Field | Value |
|-------|-------|
| Name / hostname | fill in |
| Type | mac |
| Make / Model | MacBook Pro — fill in (M-series) |
| OS | macOS — fill in |
| CPU | fill in |
| RAM | fill in |
| Storage | fill in |
| Owned / Leased | owned |

## Network

| Field | Value |
|-------|-------|
| Public IP | dynamic |
| Tailscale IP | fill in |
| Tailscale hostname | fill in |
| Local IP | 192.168.x.x (changes per network) |
| SSH access | not exposed externally |

## Software Stack

| App | Purpose |
|-----|---------|
| Obsidian | Primary vault interface |
| VS Code | Dev + Copilot |
| iTerm2 | Terminal |
| 1Password | Credential manager |
| Tailscale | Mesh VPN |
| Docker | Container workloads |
| Homebrew | Package manager |

## Vault Paths

| Context | Path |
|---------|------|
| Codex vault | `/Users/zeuslegacy/Documents/onemind-codex/` |
| onemind-web | `/Users/zeuslegacy/Documents/onemind-web/` |
| Codex Live | `/Users/zeuslegacy/Documents/onemind-codex-live/` |
| Media (large) | `/Volumes/OM PRIMARY/onemind-media/` |

## Notes

- Source of truth for all vault files — Mac edits propagate to VPS + iPhone via Obsidian Sync (planned)
- Claude Code, GitHub Copilot, Cursor all active
- **Do not run destructive vault operations from VPS without confirming Mac has latest state**


---

## Related

- [[00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)/01 The Fabric|The Fabric]]
