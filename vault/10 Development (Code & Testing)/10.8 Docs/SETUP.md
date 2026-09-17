---
id: ui-12-20260111053926
title: "Initial Setup Guide"
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
ai_summary: "Initial Setup Guide"
---

# Initial Setup Guide

## Prerequisites

- Ubuntu 22.04+ VM (2 vCPU, 4-8GB RAM recommended)
- Docker & Docker Compose installed
- Domain with DNS access (Cloudflare recommended)
- Tailscale account (for private access)

## Step 1: VM Setup

### Install Docker
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Install Docker Compose plugin
sudo apt install docker-compose-plugin -y

# Logout and back in for group changes
exit
```

### Install Tailscale
```bash
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up
```

## Step 2: Clone Repository

```bash
cd ~
git clone https://github.com/yourusername/one-mind-infrastructure.git
cd one-mind-infrastructure
```

## Step 3: Configure Environment

```bash
# Copy example environment
cp .env.example .env

# Generate encryption keys
echo "Generating keys..."
echo "N8N_ENCRYPTION_KEY=$(openssl rand -base64 32)"
echo "N8N_USER_MANAGEMENT_JWT_SECRET=$(openssl rand -base64 32)"
echo "LANGFUSE_SECRET_KEY=$(openssl rand -base64 32)"
echo "LANGFUSE_NEXT_AUTH_SECRET=$(openssl rand -base64 32)"
echo "LANGFUSE_SALT=$(openssl rand -base64 32)"

# Edit .env with your values
nano .env
```

### Critical: Save These Keys!

| Key | Save To |
|-----|---------|
| `N8N_ENCRYPTION_KEY` | Password manager |
| `N8N_USER_MANAGEMENT_JWT_SECRET` | Password manager |
| All database passwords | Password manager |

**If you lose `N8N_ENCRYPTION_KEY`, all n8n credentials are unrecoverable!**

## Step 4: Deploy

```bash
# Make scripts executable
chmod +x scripts/*.sh

# Deploy everything
./scripts/deploy.sh
```

## Step 5: Initial Service Setup

### n8n (Workflow Automation)
1. Go to `http://your-vm-ip:5678`
2. Create owner account
3. Set up your first workflow
4. Create an API key (Settings → API) for backup automation

### Grafana (Dashboards)
1. Go to `http://your-vm-ip:3000`
2. Login: admin / (password from .env)
3. Dashboards are auto-provisioned

### Uptime Kuma (External Monitoring)
1. Go to `http://your-vm-ip:3001`
2. Create admin account
3. Add monitors:
   - Todoist API
   - Letta Cloud
   - LiveKit
   - Your domains

### Langfuse (LLM Tracing)
1. Go to `http://your-vm-ip:3002`
2. Create account
3. Create a project and get API keys for your LLM applications

## Step 6: Domain Setup (Cloudflare Tunnel)

### Option A: Cloudflare Tunnel (Recommended)
```bash
# Install cloudflared
curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o cloudflared
chmod +x cloudflared
sudo mv cloudflared /usr/local/bin/

# Authenticate
cloudflared tunnel login

# Create tunnel
cloudflared tunnel create onemind

# Configure routes
cloudflared tunnel route dns onemind n8n.yourdomain.com
cloudflared tunnel route dns onemind grafana.yourdomain.com
cloudflared tunnel route dns onemind langfuse.yourdomain.com
cloudflared tunnel route dns onemind status.yourdomain.com
```

Create `/etc/cloudflared/config.yml`:
```yaml
tunnel: YOUR_TUNNEL_ID
credentials-file: /root/.cloudflared/YOUR_TUNNEL_ID.json

ingress:
  - hostname: n8n.yourdomain.com
    service: http://localhost:5678
  - hostname: grafana.yourdomain.com
    service: http://localhost:3000
  - hostname: langfuse.yourdomain.com
    service: http://localhost:3002
  - hostname: status.yourdomain.com
    service: http://localhost:3001
  - hostname: alerts.yourdomain.com
    service: http://localhost:9093
  - service: http_status:404
```

```bash
# Run as service
sudo cloudflared service install
sudo systemctl start cloudflared
```

### Option B: Tailscale Funnel
```bash
# Enable funnel for specific ports
tailscale funnel 5678  # n8n
tailscale funnel 3000  # grafana
```

## Step 7: Setup Backups

```bash
# Test backup
./scripts/backup.sh

# Add to cron (daily at 2 AM)
crontab -e
# Add: 0 2 * * * cd /home/ubuntu/one-mind-infrastructure && ./scripts/backup.sh >> /var/log/backup.log 2>&1
```

## Step 8: Configure Monitoring

### Add Monitors in Uptime Kuma
1. n8n: `http://n8n-main:5678`
2. Grafana: `http://grafana:3000`
3. External APIs: Letta Cloud, LiveKit, Todoist

### Import Grafana Dashboards
Pre-configured dashboards are auto-loaded. For additional:
1. Go to Grafana → Dashboards → Import
2. Use IDs: `1860` (Node Exporter), `14282` (cAdvisor)

## Verification Checklist

- [ ] All containers running: `docker ps`
- [ ] Portainer accessible
- [ ] n8n accessible and can create workflows
- [ ] Grafana showing metrics
- [ ] Uptime Kuma monitoring services
- [ ] Backup script works
- [ ] Domain/tunnel configured

## Troubleshooting

### Container won't start
```bash
docker logs <container-name>
```

### Database connection issues
```bash
docker exec -it n8n-postgres psql -U n8n -d n8n -c "SELECT 1"
```

### Network issues
```bash
docker network inspect onemind_network
```
