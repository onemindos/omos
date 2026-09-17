---
id: ui-12-20260111053926
title: "Operations Manual"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/12]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Operations Manual"
---

# Operations Manual

Complete guide for maintaining, accessing, and securing your One Mind Infrastructure.

## Table of Contents

1. [Access Methods](#access-methods)
2. [Security](#security)
3. [Daily Operations](#daily-operations)
4. [Deploying New Apps](#deploying-new-apps)
5. [Backup & Restore](#backup--restore)
6. [Troubleshooting](#troubleshooting)
7. [GitOps Automation](#gitops-automation)

---

## Access Methods

### 1. SSH via Tailscale (Recommended for Admin)

```bash
# From your Mac (after Tailscale is set up on both)
ssh ubuntu@<tailscale-hostname>
# or
ssh ubuntu@100.x.x.x  # Tailscale IP
```

**Why Tailscale?**
- No open ports on the internet
- Works from anywhere (home, coffee shop, phone)
- MagicDNS: access by hostname instead of IP
- Never expires (unlike SSH keys on public IP)

### 2. SSH via Public IP (Backup)

```bash
ssh -i ~/.ssh/lightsail-key.pem ubuntu@<public-ip>
```

**Note:** Only use if Tailscale is down. Requires port 22 open.

### 3. AWS Lightsail Console

1. Go to https://lightsail.aws.amazon.com
2. Click your instance
3. Click "Connect using SSH" button

**Use for:** Emergency access if locked out of SSH

### 4. VS Code Remote SSH

Add to `~/.ssh/config`:
```
Host onemind
    HostName <tailscale-ip-or-hostname>
    User ubuntu
    IdentityFile ~/.ssh/lightsail-key.pem
```

Then: `code --remote ssh-remote+onemind /home/ubuntu/one-mind-infrastructure`

---

## Security

### Firewall (UFW)

```bash
# Check status
sudo ufw status

# Current rules (set by launch script):
# - 22/tcp (SSH)
# - 41641/udp (Tailscale)
# All other incoming traffic blocked
```

### Fail2Ban

Protects against brute-force SSH attacks.

```bash
# Check banned IPs
sudo fail2ban-client status sshd

# Unban an IP
sudo fail2ban-client set sshd unbanip <IP>
```

### SSH Keys

- **Never use password authentication**
- Store private key securely (1Password, encrypted disk)
- Rotate keys periodically

### Secrets Management

| Secret | Location | Backup To |
|--------|----------|-----------|
| SSH Private Key | `~/.ssh/lightsail-key.pem` | Password manager |
| N8N_ENCRYPTION_KEY | `.env` on VM | Password manager |
| Tailscale Auth Key | Tailscale admin console | N/A (regenerate) |
| Cloudflare Tunnel Token | `/etc/cloudflared/` | Password manager |

### Security Checklist

- [ ] Password auth disabled in SSH
- [ ] Root login disabled
- [ ] UFW enabled with minimal ports
- [ ] Fail2Ban running
- [ ] Tailscale connected
- [ ] All services behind Cloudflare Tunnel (no direct port exposure)
- [ ] Encryption keys backed up to password manager

---

## Daily Operations

### Check System Health

```bash
# Quick health check
./scripts/health-check.sh

# Container status
docker ps

# Resource usage
htop

# Disk space
df -h

# Docker disk usage
docker system df
```

### View Logs

```bash
# All container logs
docker compose logs -f

# Specific service
docker logs n8n-main -f --tail 100
docker logs prometheus -f --tail 100

# System logs
sudo journalctl -u docker -f
```

### Restart Services

```bash
# Restart everything
./scripts/deploy.sh

# Restart single service
docker compose restart n8n-main

# Restart with fresh pull
docker compose pull n8n-main && docker compose up -d n8n-main
```

### Update Containers

```bash
# Manual update (with backup first)
./scripts/update.sh

# Update single container
docker compose pull grafana && docker compose up -d grafana
```

### Cleanup

```bash
# Remove unused images/containers/volumes
docker system prune -a

# Remove only dangling images
docker image prune

# Check what's using disk space
ncdu /
```

---

## Deploying New Apps

### Step 1: Create App Structure

```bash
./scripts/add-app.sh <category> <app-name>
# Example:
./scripts/add-app.sh automation my-new-app
```

This creates:
```
apps/<category>/<app-name>/
├── docker-compose.yml
└── README.md
```

### Step 2: Configure Docker Compose

Edit `apps/<category>/<app-name>/docker-compose.yml`:

```yaml
services:
  my-new-app:
    image: someimage:latest
    container_name: my-new-app
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - MY_VAR=${MY_VAR}
    volumes:
      - my_new_app_data:/data
    networks:
      - onemind

volumes:
  my_new_app_data:
    name: my_new_app_data

networks:
  onemind:
    external: true
    name: onemind_network
```

### Step 3: Add to Master Compose

Edit `docker-compose.yml`:

```yaml
include:
  # ... existing includes
  - apps/<category>/<app-name>/docker-compose.yml
```

### Step 4: Add Environment Variables

Edit `.env`:
```bash
# My New App
MY_VAR=some_value
```

### Step 5: Add Cloudflare Tunnel Route

Edit `/etc/cloudflared/config.yml` on VM:
```yaml
ingress:
  # ... existing routes
  - hostname: my-new-app.zeusops.dev
    service: http://localhost:8080
```

Then:
```bash
cloudflared tunnel route dns onemind my-new-app.zeusops.dev
sudo systemctl restart cloudflared
```

### Step 6: Deploy

```bash
./scripts/deploy.sh my-new-app
```

### Step 7: Commit to Git

```bash
git add .
git commit -m "Add my-new-app"
git push
```

---

## Backup & Restore

### Manual Backup

```bash
./scripts/backup.sh
```

Creates:
- `backups/n8n_db_YYYYMMDD.sql.gz` - Database dump
- `workflows/*.json` - n8n workflow exports

### Restore

```bash
./scripts/restore.sh              # Latest backup
./scripts/restore.sh 20241219_120000  # Specific backup
```

### Critical Files to Backup Externally

| File | What It Contains |
|------|------------------|
| `.env` | All secrets and config |
| `/etc/cloudflared/` | Tunnel credentials |
| `backups/*.sql.gz` | Database dumps |

### Disaster Recovery

If VM is completely lost:

1. Create new Lightsail instance with launch script
2. SSH in and connect Tailscale
3. Clone repo: `git clone <repo> ~/one-mind-infrastructure`
4. Restore `.env` from password manager
5. Restore database: `./scripts/restore.sh`
6. Setup Cloudflare Tunnel (re-authenticate)
7. Deploy: `./scripts/deploy.sh`

**Recovery Time:** ~30 minutes

---

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker logs <container-name>

# Check compose config
docker compose config

# Try recreating
docker compose up -d --force-recreate <service>
```

### Out of Disk Space

```bash
# Find what's using space
ncdu /

# Clean Docker
docker system prune -a --volumes

# Check log sizes
du -sh /var/lib/docker/containers/*/
```

### Can't Connect to Service

```bash
# Check if container is running
docker ps | grep <service>

# Check if port is listening
sudo ss -tlnp | grep <port>

# Check Cloudflare Tunnel
sudo systemctl status cloudflared
cloudflared tunnel info onemind
```

### High Memory Usage

```bash
# Check per-container memory
docker stats --no-stream

# Reduce n8n workers if needed (edit compose, reduce to 1 worker)
```

### Tailscale Disconnected

```bash
sudo tailscale status
sudo tailscale up
```

---

## GitOps Automation

### How It Works (When Enabled)

```
You push to GitHub → GitHub Actions → SSH to VM → Deploy
```

### Enable Workflows

```bash
# Rename back to .yml
cd .github/workflows
mv deploy.yml.disabled deploy.yml
mv backup.yml.disabled backup.yml
mv update.yml.disabled update.yml
git add . && git commit -m "Enable workflows" && git push
```

### Required GitHub Secrets

Go to GitHub repo → Settings → Secrets → Actions:

| Secret | Value |
|--------|-------|
| `VM_HOST` | Your VM's public IP or Tailscale hostname |
| `VM_USER` | `ubuntu` |
| `VM_SSH_KEY` | Contents of your private key file |

### Workflow Schedule

| Workflow | Trigger |
|----------|---------|
| `deploy.yml` | Push to main (auto-deploy) |
| `backup.yml` | Daily 2 AM UTC |
| `update.yml` | Sunday 3 AM UTC |

### Manual Trigger

1. Go to GitHub → Actions
2. Select workflow
3. Click "Run workflow"

---

## Quick Reference

### Common Commands

```bash
# Deploy all
./scripts/deploy.sh

# Check health
./scripts/health-check.sh

# View logs
docker compose logs -f

# Restart service
docker compose restart <service>

# SSH to VM
ssh ubuntu@<tailscale-ip>

# Tailscale status
tailscale status
```

### Service URLs

| Service | URL |
|---------|-----|
| n8n | https://n8n.zeusops.dev |
| Grafana | https://grafana.zeusops.dev |
| Uptime Kuma | https://status.zeusops.dev |
| Langfuse | https://langfuse.zeusops.dev |
| Alertmanager | https://alerts.zeusops.dev |

### File Locations on VM

| Path | Contents |
|------|----------|
| `~/one-mind-infrastructure/` | Main repo |
| `~/one-mind-infrastructure/.env` | Secrets |
| `~/one-mind-infrastructure/backups/` | Local backups |
| `/etc/cloudflared/` | Tunnel config |
| `/var/lib/docker/volumes/` | Container data |


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
