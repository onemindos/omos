---
title: "OpenClaw — Updates Management"
type: sop
status: active
created: 2026-04-12
owner: zeus
tags: [type/sop, domain/SO, domain/openclaw, status/active]
ai_summary: "How to update OpenClaw on the VPS — check version, npm update, rollback, service restart, health check."
---

# OpenClaw — Updates Management

How to check and update OpenClaw on the VPS (`104.248.4.212`).

---

## Check Current Version

```bash
ssh root@104.248.4.212 'openclaw --version'
# or via npm:
ssh root@104.248.4.212 'npm list -g openclaw'
```

---

## Update OpenClaw

```bash
# 1. Update via npm
ssh root@104.248.4.212 'npm update -g openclaw'

# 2. Verify version bumped
ssh root@104.248.4.212 'openclaw --version'

# 3. Restart the service
ssh root@104.248.4.212 'systemctl restart openclaw'

# 4. Check it came back healthy
ssh root@104.248.4.212 'systemctl status openclaw --no-pager -l | head -20'
```

---

## Check Service Health After Update

```bash
# Status
ssh root@104.248.4.212 'systemctl status openclaw --no-pager -l'

# Last 50 log lines
ssh root@104.248.4.212 'journalctl -u openclaw -n 50 --no-pager'

# Live logs (stream)
ssh root@104.248.4.212 'journalctl -u openclaw -f'
```

---

## Rollback to Previous Version

```bash
# Install a specific version
ssh root@104.248.4.212 'npm install -g openclaw@<previous-version>'

# Restart
ssh root@104.248.4.212 'systemctl restart openclaw'
```

To find previous version number: check npm registry or `journalctl` logs for the last known good version string.

---

## Update Plugins / Extensions

```bash
# List loaded plugins
ssh root@104.248.4.212 'openclaw plugins list'

# Update a specific plugin
ssh root@104.248.4.212 'openclaw plugins update <plugin-name>'

# Or update all extensions via npm (they live in openclaw's node_modules)
ssh root@104.248.4.212 'npm update -g openclaw'
```

---

## VPS Service File

Location: `/etc/systemd/system/openclaw.service`

```bash
# View service definition
ssh root@104.248.4.212 'cat /etc/systemd/system/openclaw.service'

# After editing service file — reload daemon
ssh root@104.248.4.212 'systemctl daemon-reload && systemctl restart openclaw'
```

---

## Related

- [DEPLOYMENT.md](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/openclaw/DEPLOYMENT.md) — pushing workspace/config changes
- [VPS-SERVICES.md](VPS-SERVICES.md) — all running services on VPS
- [OC VPS STATUS.md](OC%20VPS%20STATUS.md) — live snapshot of VPS state
