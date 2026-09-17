---
id: ui-12-20260111053926
title: "OneMind Infrastructure Maintenance Runbook"
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
ai_summary: "OneMind Infrastructure Maintenance Runbook"
---

# OneMind Infrastructure Maintenance Runbook

## Overview

This document provides procedures for maintaining the OneMind infrastructure running on AWS Lightsail.

### Infrastructure Summary

| Component | Containers | Domain |
|-----------|------------|--------|
| Coolify | 6 | ops.zeusops.dev |
| N8N | 4 | n8n.zeusops.dev |
| Uptime Kuma | 1 | status.zeusops.dev |
| Discord Bot | 1 | bot.zeusops.dev |
| **Total** | **12** | |

### Access Information

- **VM IP**: 44.212.233.193
- **Tailscale IP**: 100.86.20.51
- **Coolify UI**: https://ops.zeusops.dev (or http://44.212.233.193:8000)
- **SSH**: `ssh -i "LightsailDefaultKey-us-east-1.pem" ubuntu@44.212.233.193`

### Current Container Names

| Service | Container Name |
|---------|---------------|
| N8N Main | n8n-main-x404kkg4cwc4o8so040880ws |
| N8N Worker | n8n-worker-x404kkg4cwc4o8so040880ws |
| N8N Postgres | n8n-postgres-x404kkg4cwc4o8so040880ws |
| N8N Redis | n8n-redis-x404kkg4cwc4o8so040880ws |
| Uptime Kuma | d4wk00gokgw4w8ggkwcc84gs-072403683436 |
| Discord Bot | z0ckk0ko4s0okc4cg8444o8c-061953116013 |

### Important Credentials Location

| Credential | Location |
|------------|----------|
| N8N Encryption Key | Coolify UI > N8N Service > Environment Variables |
| N8N DB Password | Coolify UI > N8N Service > Environment Variables |
| Coolify API Token | `1\|cwFeuOQMZjOS24o1UvfAgRMM1q2NhfihDZpo1zka1b8bd131` |
| Discord Bot Token | Coolify UI > Discord Bot > Environment Variables |

---

## Automated Tasks

These tasks run automatically via cron:

| Task | Schedule | Log File |
|------|----------|----------|
| Backup | Daily at 3 AM | `/home/ubuntu/logs/backup.log` |
| Docker Cleanup | Sundays at 4 AM | `/home/ubuntu/logs/docker-cleanup.log` |
| Security Updates | Automatic | `/var/log/unattended-upgrades/` |

---

## Weekly Tasks

### 1. Check Uptime Kuma for Alerts
- Visit https://status.zeusops.dev
- Review any downtime incidents
- Check notification channels are working

### 2. Review Backup Logs
```bash
ssh -i "LightsailDefaultKey-us-east-1.pem" ubuntu@44.212.233.193

# Check backup log
tail -50 /home/ubuntu/logs/backup.log

# Verify backup files exist
ls -lh /home/ubuntu/backups/
```

### 3. Check Disk Space
```bash
df -h
du -sh /home/ubuntu/backups/
docker system df
```

---

## Monthly Tasks

### 1. Update Ubuntu OS
```bash
sudo apt update
sudo apt upgrade -y
sudo apt dist-upgrade -y
sudo apt autoremove -y
```

### 2. Update Docker
```bash
# Check current version
docker --version

# Update Docker
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
```

### 3. Update Coolify
```bash
# Run the official installer (handles updates)
curl -fsSL https://cdn.coollabs.io/coolify/install.sh | sudo bash
```

### 4. Pull Latest Container Images
Via Coolify UI:
1. Go to each application
2. Click "Redeploy" to pull latest images

Or manually:
```bash
# Pull latest images
docker images --format "{{.Repository}}:{{.Tag}}" | grep -v '<none>' | xargs -L1 docker pull
```

### 5. Review Security Logs
```bash
# Check failed SSH attempts
sudo cat /var/log/auth.log | grep "Failed password" | tail -20

# Check unattended-upgrades
cat /var/log/unattended-upgrades/unattended-upgrades.log | tail -50
```

### 6. Test Backup Restore (Optional)
```bash
# Run manual backup first
/home/ubuntu/scripts/backup.sh

# List backups
ls -la /home/ubuntu/backups/

# Test restore N8N to temp database (don't overwrite production!)
# gunzip -c /home/ubuntu/backups/n8n_db_XXXXXX.sql.gz | head -100
```

---

## Quarterly Tasks

### 1. Review and Rotate Credentials
- Discord Bot Token
- N8N Encryption Key (if compromised)
- Coolify API Token
- GitHub Deploy Key

### 2. Audit Firewall Rules
```bash
sudo ufw status verbose
```

Expected rules:
- 22/tcp (SSH)
- 80/tcp (HTTP)
- 443/tcp (HTTPS)
- 8000/tcp (Coolify)
- 41641/udp (Tailscale)

### 3. Check SSL Certificate Expiry
- Coolify auto-renews Let's Encrypt certificates
- Verify in browser or:
```bash
echo | openssl s_client -servername n8n.zeusops.dev -connect n8n.zeusops.dev:443 2>/dev/null | openssl x509 -noout -dates
```

### 4. Performance Review
```bash
# System resources
htop
free -h
df -h

# Container stats
docker stats --no-stream
```

---

## Troubleshooting

### Service Down - Quick Recovery

1. **Check if container is running**:
```bash
docker ps -a | grep <service-name>
```

2. **Check container logs**:
```bash
docker logs <container-name> --tail 100
```

3. **Restart via Coolify UI**:
   - Go to application in Coolify
   - Click "Restart"

4. **Restart via CLI**:
```bash
docker restart <container-name>
```

### N8N Not Starting

1. Check database connection:
```bash
docker exec n8n-postgres-x404kkg4cwc4o8so040880ws pg_isready -U n8n
```

2. Check Redis:
```bash
docker exec n8n-redis-x404kkg4cwc4o8so040880ws redis-cli ping
```

3. Check N8N logs:
```bash
docker logs n8n-main-x404kkg4cwc4o8so040880ws --tail 200
```

### Uptime Kuma Issues

1. Check container:
```bash
docker logs d4wk00gokgw4w8ggkwcc84gs-072403683436 --tail 100
```

2. Restart if needed:
```bash
docker restart d4wk00gokgw4w8ggkwcc84gs-072403683436
```

### Disk Space Full

```bash
# Check usage
df -h
du -sh /home/ubuntu/backups/*

# Cleanup old backups manually
find /home/ubuntu/backups -type f -mtime +7 -delete

# Docker cleanup
docker system prune -af --volumes
```

### Coolify Not Accessible

1. Check containers:
```bash
docker ps | grep coolify
```

2. Restart Coolify:
```bash
sudo bash -c 'cd /data/coolify/source && docker compose up -d'
```

3. If still broken, reinstall:
```bash
curl -fsSL https://cdn.coollabs.io/coolify/install.sh | sudo bash
```

---

## Backup & Restore

### Manual Backup
```bash
/home/ubuntu/scripts/backup.sh
```

### Restore N8N Database
```bash
# Stop N8N first
docker stop n8n-main-x404kkg4cwc4o8so040880ws n8n-worker-x404kkg4cwc4o8so040880ws

# Restore database
gunzip -c /home/ubuntu/backups/n8n_db_XXXXXX.sql.gz | docker exec -i n8n-postgres-x404kkg4cwc4o8so040880ws psql -U n8n n8n

# Start N8N
docker start n8n-main-x404kkg4cwc4o8so040880ws n8n-worker-x404kkg4cwc4o8so040880ws
```

### Restore Uptime Kuma
```bash
# Get container name
KUMA_CONTAINER=$(docker ps --filter "name=d4wk00gokgw4w8ggkwcc84gs" --format "{{.Names}}" | head -1)

# Stop container
docker stop $KUMA_CONTAINER

# Extract backup
tar -xzf /home/ubuntu/backups/uptime_kuma_XXXXXX.tar.gz -C /tmp/

# Copy to container volume
docker cp /tmp/uptime_kuma_XXXXXX/. $KUMA_CONTAINER:/app/data/

# Start container
docker start $KUMA_CONTAINER
```

---

## Adding New Services

### Step 1: Deploy via Coolify

1. Go to https://ops.zeusops.dev
2. Create new application (Docker Image or Docker Compose)
3. Configure environment variables
4. Deploy the service
5. Note the container name: `docker ps --format "{{.Names}}" | grep <service-id>`

### Step 2: Add DNS Record in Cloudflare

1. Go to Cloudflare Dashboard > zeusops.dev > DNS
2. Add A record:
   - Type: `A`
   - Name: `<subdomain>` (e.g., `api`)
   - Content: `44.212.233.193`
   - Proxy: `Proxied` (orange cloud)
3. SSL/TLS mode should be set to `Full`

### Step 3: Create Traefik Routing Config

SSH to server and create routing file:

```bash
ssh -i "LightsailDefaultKey-us-east-1.pem" ubuntu@44.212.233.193

# Create routing config (replace values)
sudo tee /data/coolify/proxy/dynamic/<service-name>.yaml << 'EOF'
http:
  routers:
    <service-name>-http:
      entryPoints:
        - http
      service: <service-name>
      rule: Host(`<subdomain>.zeusops.dev`)
      middlewares:
        - redirect-to-https
    <service-name>-https:
      entryPoints:
        - https
      service: <service-name>
      rule: Host(`<subdomain>.zeusops.dev`)
      tls:
        certresolver: letsencrypt
  services:
    <service-name>:
      loadBalancer:
        servers:
          - url: 'http://<container-name>:<port>'
EOF
```

### Step 4: Connect Container to Coolify Network

```bash
# Connect container to the coolify network for Traefik routing
docker network connect coolify <container-name>
```

### Step 5: Verify

```bash
# Test routing locally
curl -I -H "Host: <subdomain>.zeusops.dev" http://localhost

# Test through Cloudflare
curl -I https://<subdomain>.zeusops.dev
```

### Current Domain Routing Files

| Domain | Config File | Location |
|--------|-------------|----------|
| ops.zeusops.dev | coolify.yaml | `/data/coolify/proxy/dynamic/coolify.yaml` |
| n8n.zeusops.dev | n8n.yaml | `/data/coolify/proxy/dynamic/n8n.yaml` |
| status.zeusops.dev | uptime-kuma.yaml | `/data/coolify/proxy/dynamic/uptime-kuma.yaml` |

### Example: Adding a New API Service

```bash
# 1. After deploying "my-api" in Coolify, find container name
docker ps --format "{{.Names}}" | grep my-api
# Output: my-api-abc123

# 2. Create Traefik config
sudo tee /data/coolify/proxy/dynamic/my-api.yaml << 'EOF'
http:
  routers:
    my-api-http:
      entryPoints:
        - http
      service: my-api
      rule: Host(`api.zeusops.dev`)
      middlewares:
        - redirect-to-https
    my-api-https:
      entryPoints:
        - https
      service: my-api
      rule: Host(`api.zeusops.dev`)
      tls:
        certresolver: letsencrypt
  services:
    my-api:
      loadBalancer:
        servers:
          - url: 'http://my-api-abc123:3000'
EOF

# 3. Connect to network
docker network connect coolify my-api-abc123

# 4. Add DNS in Cloudflare: A record "api" -> 44.212.233.193 (proxied)

# 5. Test
curl -I https://api.zeusops.dev
```

---

## Emergency Contacts

- **Coolify Docs**: https://coolify.io/docs
- **N8N Docs**: https://docs.n8n.io
- **Uptime Kuma**: https://github.com/louislam/uptime-kuma

---

## Important Files & Directories

### Scripts & Backups

| File | Purpose |
|------|---------|
| `/home/ubuntu/scripts/backup.sh` | Daily backup script |
| `/home/ubuntu/backups/` | Backup storage (14-day retention) |
| `/home/ubuntu/logs/` | Log files |

### Coolify Files

| File | Purpose |
|------|---------|
| `/data/coolify/source/.env` | Coolify environment variables |
| `/data/coolify/source/docker-compose.yml` | Coolify containers |
| `/data/coolify/services/` | Deployed service configs |

### Traefik / Domain Routing

| File | Purpose |
|------|---------|
| `/data/coolify/proxy/dynamic/` | Domain routing configs (auto-reload) |
| `/data/coolify/proxy/dynamic/coolify.yaml` | ops.zeusops.dev routing |
| `/data/coolify/proxy/dynamic/n8n.yaml` | n8n.zeusops.dev routing |
| `/data/coolify/proxy/dynamic/uptime-kuma.yaml` | status.zeusops.dev routing |
| `/data/coolify/proxy/acme.json` | Let's Encrypt certificates |

### SSH Access

| File | Purpose |
|------|---------|
| `~/.ssh/LightsailDefaultKey-us-east-1.pem` | SSH key (local machine) |

---

## AWS Lightsail Firewall

Required ports in **Lightsail Console > Networking > IPv4 Firewall**:

| Port | Protocol | Purpose |
|------|----------|---------|
| 22 | TCP | SSH |
| 80 | TCP | HTTP |
| 443 | TCP | HTTPS |
| 8000 | TCP | Coolify UI |
| 41641 | UDP | Tailscale |

---

## Security Notes

- **SSH**: Key-only authentication (password disabled)
- **Firewall**: UFW enabled with minimal ports
- **Updates**: Unattended security updates enabled
- **Backups**: Local only - consider adding cloud backup for disaster recovery

---

*Last Updated: 2025-12-25*


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
