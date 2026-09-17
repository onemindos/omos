# Bash / Shell

> **Unix shell scripting for automation, deployment, and VPS operations**

## Overview

Bash and shell scripting are the glue of the One Mind VPS infrastructure. All deployment scripts, systemd service wrappers, cron jobs, onemind-watch, and maintenance automation run as shell scripts. Default shell on Ubuntu VPS and macOS (zsh is default on Mac but bash scripts run on both).

## Status

| Attribute   | Value                               |
| ----------- | ----------------------------------- |
| Status      | ✅ Active                            |
| Category    | Development / Scripting             |
| Type        | Shell Scripting Language            |
| Integration | bash, zsh, sh, systemd, cron        |

## Use Cases in OneMind

1. **VPS Automation** — Deploy scripts, service control, health checks
2. **onemind-watch.sh** — File watcher for auto git commit + push
3. **Deployment** — `openclaw-reload.sh`, docker compose wrappers
4. **CI/CD** — GitHub Actions shell steps
5. **Vault Maintenance** — File moves, dedup, cleanup scripts

## Key Scripts

| Script | Location | Purpose |
| ------ | -------- | ------- |
| `onemind-watch.sh` | `_system/infrastructure/` | File change → git commit |
| `openclaw-reload.sh` | `_system/openclaw/deploy/` | Deploy OpenClaw to VPS |
| Various | `_system/hooks/` | CODEX lifecycle automation |

## VPS Shell

```bash
# SSH to VPS
ssh root@104.248.4.212

# Check services
systemctl status openclaw n8n

# Deploy
cd /opt/onemind-web && docker compose up -d --build
```

## Best Practices

```bash
#!/usr/bin/env bash
set -euo pipefail          # Exit on error, undefined vars, pipe failures
BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null)
if [[ -z "$BRANCH" ]]; then
  git checkout main         # Recover from detached HEAD
fi
```

## Links

- Bash Manual: [gnu.org/software/bash/manual](https://www.gnu.org/software/bash/manual)
- ShellCheck: [shellcheck.net](https://shellcheck.net) — lint shell scripts
