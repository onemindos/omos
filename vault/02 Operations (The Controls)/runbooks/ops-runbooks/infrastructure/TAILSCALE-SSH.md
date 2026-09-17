---
id: so-03-20260430044500
title: "Tailscale SSH — Enable & Manage"
type: runbook
pillar: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-04-30
modified: 2026-05-08
owner: zeus
tags: [domain/so, infra/tailscale, infra/ssh, networking]
ai_summary: "Enable and manage Tailscale SSH across the fleet — keyless, identity-based SSH via tailnet"
---

# Tailscale SSH — Enable & Manage

## What It Does

Tailscale SSH replaces traditional SSH key management with identity-based auth via your Tailscale account. No keys to distribute, no port 22 exposed to the internet, access from anywhere on the tailnet.

## Fleet Status

| Machine | Tailscale IP | OS | Tailscale SSH | Notes |
|---------|-------------|-----|---------------|-------|
| omos-core | 100.123.126.34 | Ubuntu 24.04 | ✅ | Lightsail, user: ubuntu |
| omos-sentinel | 100.66.235.38 | Ubuntu | ✅ | Lightsail, user: ubuntu |
| omos-titan | 100.103.197.23 | Ubuntu | ✅ | Lightsail, user: ubuntu |
| omos-dell | 100.66.84.71 | Linux | ✅ | Local workstation |
| omos-ha | 100.64.25.12 | HA OS | ⚠️ | Local LAN only — offline as of 2026-05-08 |

## Enable on Standard Linux (Ubuntu/Debian)

```bash
# SSH into the machine first, then:
sudo tailscale set --ssh --accept-risk=lose-ssh
```

> **Note:** This will kill your current SSH session if connecting via traditional SSH. Reconnect using Tailscale IP afterward.

After running, the old host key changes. From your client:
```bash
ssh-keygen -R <tailscale-ip>
ssh -o StrictHostKeyChecking=accept-new root@<tailscale-ip>
```

## Enable on Home Assistant OS

HA OS runs Tailscale inside a container. The `tailscale` binary is NOT on the host PATH.

### Prerequisites
1. **Advanced SSH & Web Terminal** add-on installed
2. **Protection mode** disabled on the SSH add-on:
   - HA Web UI → Settings → Add-ons → Advanced SSH & Web Terminal
   - Toggle "Protection mode" OFF → Restart add-on

### Steps

```bash
# Find the Tailscale container name
docker ps | grep tailscale

# Run tailscale set inside the container (binary is at /opt/tailscale)
docker exec <container_name> /opt/tailscale set --ssh --accept-risk=lose-ssh
```

Example with actual container name:
```bash
docker exec addon_a0d7b954_tailscale /opt/tailscale set --ssh --accept-risk=lose-ssh
```

## Verify from Client (Mac)

```bash
ssh root@<tailscale-ip> 'echo working'
```

If host key error:
```bash
ssh-keygen -R <tailscale-ip>
ssh -o StrictHostKeyChecking=accept-new root@<tailscale-ip> 'echo working'
```

## Disable Tailscale SSH

```bash
# Standard Linux:
sudo tailscale set --ssh=false

# HA OS:
docker exec <container_name> /opt/tailscale set --ssh=false
```

## Troubleshooting

| Issue | Fix |
|-------|-----|
| `command not found: tailscale` | Binary not on PATH. On HA OS use `docker exec` with `/opt/tailscale` |
| `Permission denied (publickey)` | Tailscale SSH not enabled yet, or connecting via wrong IP |
| `REMOTE HOST IDENTIFICATION HAS CHANGED` | Expected after enabling — run `ssh-keygen -R <ip>` |
| Protection mode error on HA | Disable in HA Web UI → Settings → Add-ons → SSH add-on |
| `OCI runtime exec failed` | Wrong binary path — use `find / -name "tailscale" -type f` inside container |

## Access Control

Manage who can SSH to what via Tailscale admin:
- **ACLs**: https://login.tailscale.com/admin/acls
- **Machines**: https://login.tailscale.com/admin/machines

## Security Notes

- Tailscale SSH uses WireGuard encryption (same as tailnet traffic)
- Auth is tied to your Tailscale identity (Google/GitHub SSO)
- No port 22 exposed publicly — only reachable via tailnet
- Audit logs available in Tailscale admin panel


---

## Related

- [[OMOS-HA-TAILSCALE-FIX]] — Fix runbook for omos-ha add-on crash (ssh:true mismatch)
- [[TAILSCALE-NETWORK]] — Full tailnet map and fleet table
- [[omos-ha]] — Home Assistant edge node reference
- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
