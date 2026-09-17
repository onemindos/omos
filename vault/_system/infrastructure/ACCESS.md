# ACCESS.md — OneMind OS Quick Access Reference

All URLs, ports, and access points for the OneMind OS infrastructure.

---

## VPS Services (onemind-main)

| Service | Internal Port | Tailscale HTTPS | Purpose |
|---------|--------------|----------------|---------|
| **OpenClaw Gateway** | `:18789` | https://onemind-main.tail717ba0.ts.net | Agent runtime (WebSocket + REST API) |
| **OpenClaw Studio** | `:3000` | https://onemind-main.tail717ba0.ts.net:3443 | Native agent dashboard (chat, approvals, cron) |
| **SwarmClaw** | `:3456` | https://onemind-main.tail717ba0.ts.net:8443 | Agent orchestration (tasks, chatrooms, MCP, workflows) |

> All Tailscale URLs are **tailnet-only** — requires Tailscale on the connecting device.

### SwarmClaw Access Key

```
49359e4f6cd67002f2e2b9b2491a6dbb
```

### VPS SSH

```bash
ssh root@104.248.4.212        # Main VPS (onemind-main)
```

- Tailscale IP: `100.77.52.93`
- OS: Ubuntu 24.04 | 16GB RAM | 4 cores
- OpenClaw version: `2026.3.13 (61d171a)`

---

## Sentinel VPS (onemind-sentinel)

| Service | Internal Port | Tailscale HTTPS | Purpose |
|---------|--------------|----------------|---------|
| **OpenClaw Sentinel** | `:19789` | https://onemind-sentinel.tail717ba0.ts.net | Security & monitoring agent |

```bash
ssh root@159.203.122.175      # Sentinel VPS
```

- Tailscale IP: `100.67.223.5`
- OpenClaw version: `2026.3.13 (61d171a)`

---

## Local Development

| Service | URL | Start Command |
|---------|-----|--------------|
| SwarmClaw (dev) | http://localhost:3456 | `git clone https://github.com/One-Mind-OS/swarmclaw.git && cd swarmclaw && npm run dev` |
| Chowder-iOS | Xcode | `git clone https://github.com/One-Mind-OS/chowder-iOS.git` → Open in Xcode |
| OpenClaw Gateway (local) | ws://localhost:18789 | `openclaw gateway` |

---

## Agent Roster (9 agents on onemind-main)

| Agent | ID | Model | Heartbeat |
|-------|----|-------|-----------|
| ⚡ Legacy | `legacy` | azure-openai/gpt-5.2 | 1h |
| ⚔️ Spartan | `spartan` | azure-openai/gpt-4o | 2h |
| 🔮 Oracle | `oracle` | azure-openai/gpt-4o | 3h |
| 🌳 Eden | `eden` | azure-openai/gpt-4o | 3h |
| 👑 Empire | `empire` | azure-openai/gpt-5.2 | 2h |
| 🌿 Haven | `haven` | azure-openai/gpt-4o | 3h |
| 🛡️ Guardian | `guardian` | azure-openai/gpt-4o | 2h |
| ⚙️ Grid | `grid` | azure-openai/gpt-4o | 3h |
| 🔺 Trinity | `trinity` | azure-openai/gpt-4o | 3h |

---

## Systemd Services (onemind-main)

| Service | Unit File | Status Command |
|---------|----------|---------------|
| OpenClaw Gateway | `openclaw.service` | `systemctl status openclaw` |
| OpenClaw Studio | `openclaw-studio.service` | `systemctl status openclaw-studio` |
| SwarmClaw | `swarmclaw.service` | `systemctl status swarmclaw` |
| Codex Watch | `codex-watch.service` | `systemctl status codex-watch` |
| Codex Pull | `codex-pull.timer` | `systemctl status codex-pull` |
| Hot Reload | `openclaw-reload.service` | `systemctl status openclaw-reload` |

```bash
# Quick status check (all services)
ssh root@104.248.4.212 "systemctl status openclaw openclaw-studio swarmclaw codex-watch codex-pull --no-pager -l"

# Restart a service
ssh root@104.248.4.212 "systemctl restart openclaw"

# View logs
ssh root@104.248.4.212 "journalctl -u openclaw -f"
ssh root@104.248.4.212 "journalctl -u swarmclaw -f"
ssh root@104.248.4.212 "journalctl -u openclaw-studio -f"
```

---

## GitHub

| Resource | URL |
|----------|-----|
| Monorepo | https://github.com/One-Mind-OS/onemindos-v0.5 |

---

## Key Paths

| Context | Path |
|---------|------|
| VPS monorepo | `/home/onemind/onemind-os/` |
| VPS gateway config | `/home/onemind/.openclaw/openclaw.json` |
| VPS legacy workspace | `/home/onemind/.openclaw/workspace` → symlink to `onemind-os/openclaw/workspaces/legacy` |
| VPS codex | `/home/onemind/onemind-os/codex/` |
| VPS secrets template | `/opt/onemind/.env.tpl` (op:// references) |
| VPS static env | `/opt/onemind/.env.sa` |
| Local monorepo | `/Users/zeuslegacy/Documents/onemindos-v0.5/` |
| Local codex | `/Users/zeuslegacy/Documents/onemindos-v0.5/codex/` |
