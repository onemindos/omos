---
name: OpenClaw — Install & Deploy
updated: 2026-05-10
---

# Install & Deploy

VPS deployment infrastructure for omos-core and omos-sentinel — scripts, systemd services, AWS Lightsail setup, backups, and config templates.

**Official docs:** [docs.openclaw.ai/install](https://docs.openclaw.ai/install/index.md)

## Subfolders

| Folder | Contents |
|---|---|
| [backups/](backups/) | `openclaw.json` config snapshots (omos-core, omos-sentinel) |
| [config-templates/](config-templates/) | `openclaw.json.example`, `env.tpl`, `env.example`, `ambient-overlay.json`, `apply-agent-config.sh` |
| [do-scripts/](do-scripts/) | DigitalOcean droplet scripts: install, setup, firewall, health-check |
| [scripts/](scripts/) | Deployment scripts: bootstrap, onemind-watch, reload, symlinks, workspace-backup, sentinel-setup, cron-reports, voice patch |
| [systemd/](systemd/) | Systemd service/timer files: openclaw.service, onemind-pull.timer, onemind-watch.service, and more |

## Key Official Pages

- [Install overview](https://docs.openclaw.ai/install/index.md)
- [DigitalOcean](https://docs.openclaw.ai/install/digitalocean.md)
- [Linux server (VPS)](https://docs.openclaw.ai/vps.md)
- [Docker](https://docs.openclaw.ai/install/docker.md)
- [Updating](https://docs.openclaw.ai/install/updating.md)

## SSH Access

```bash
ssh -i ~/.ssh/id_ed25519 ubuntu@100.123.126.34  # omos-core via Tailscale
```
