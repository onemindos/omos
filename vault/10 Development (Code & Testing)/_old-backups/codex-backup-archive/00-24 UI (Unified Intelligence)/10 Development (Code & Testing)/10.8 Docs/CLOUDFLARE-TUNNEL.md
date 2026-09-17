---
id: ui-12-20260111053926
title: "Cloudflare Tunnel Setup Guide"
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
ai_summary: "Cloudflare Tunnel Setup Guide"
---

# Cloudflare Tunnel Setup Guide

## Overview

Cloudflare Tunnel allows secure access to your services without opening any ports on your VM.

```
User → zeusops.dev → Cloudflare → Tunnel → Your VM (no open ports!)
```

## Prerequisites

- Cloudflare account
- Domain on Cloudflare (zeusops.dev)
- `cloudflared` installed on VM (done by setup-vm.sh)

## Step 1: Authenticate Cloudflared

```bash
cloudflared tunnel login
```

This opens a browser to authenticate with Cloudflare.

## Step 2: Create Tunnel

```bash
cloudflared tunnel create onemind
```

This creates:
- Tunnel ID (save this!)
- Credentials file: `~/.cloudflared/<TUNNEL_ID>.json`

## Step 3: Create Tunnel Config

Create `/etc/cloudflared/config.yml`:

```yaml
tunnel: YOUR_TUNNEL_ID
credentials-file: /home/ubuntu/.cloudflared/YOUR_TUNNEL_ID.json

ingress:
  # n8n
  - hostname: n8n.zeusops.dev
    service: http://localhost:5678
  
  # Grafana
  - hostname: grafana.zeusops.dev
    service: http://localhost:3000
  
  # Uptime Kuma (Status Page)
  - hostname: status.zeusops.dev
    service: http://localhost:3001
  
  # Langfuse
  - hostname: langfuse.zeusops.dev
    service: http://localhost:3002
  
  # Alertmanager
  - hostname: alerts.zeusops.dev
    service: http://localhost:9093
  
  # Prometheus (optional - usually private)
  # - hostname: prometheus.zeusops.dev
  #   service: http://localhost:9090
  
  # Catch-all (required)
  - service: http_status:404
```

## Step 4: Create DNS Records

```bash
# Create CNAME records pointing to tunnel
cloudflared tunnel route dns onemind n8n.zeusops.dev
cloudflared tunnel route dns onemind grafana.zeusops.dev
cloudflared tunnel route dns onemind status.zeusops.dev
cloudflared tunnel route dns onemind langfuse.zeusops.dev
cloudflared tunnel route dns onemind alerts.zeusops.dev
```

## Step 5: Run as Service

```bash
# Install as system service
sudo cloudflared service install

# Start the service
sudo systemctl start cloudflared

# Enable on boot
sudo systemctl enable cloudflared

# Check status
sudo systemctl status cloudflared
```

## Step 6: Verify

```bash
# Check tunnel is running
cloudflared tunnel info onemind

# Test access
curl -I https://n8n.zeusops.dev
```

## Domains Summary

| Subdomain | Service | Port |
|-----------|---------|------|
| n8n.zeusops.dev | n8n | 5678 |
| grafana.zeusops.dev | Grafana | 3000 |
| status.zeusops.dev | Uptime Kuma | 3001 |
| langfuse.zeusops.dev | Langfuse | 3002 |
| alerts.zeusops.dev | Alertmanager | 9093 |

## Troubleshooting

### Tunnel not connecting
```bash
# Check logs
sudo journalctl -u cloudflared -f

# Test config
cloudflared tunnel --config /etc/cloudflared/config.yml run
```

### DNS not resolving
```bash
# Verify DNS records in Cloudflare dashboard
# Should see CNAME records pointing to <tunnel-id>.cfargotunnel.com
```

### Service unreachable
```bash
# Check if service is running locally
curl http://localhost:5678

# Check Docker containers
docker ps
```

## Security Notes

1. **No open ports**: All traffic goes through Cloudflare
2. **DDoS protection**: Cloudflare handles this
3. **SSL/TLS**: Automatic, managed by Cloudflare
4. **Access control**: Use Cloudflare Access for additional auth (optional)

## Adding New Services

1. Add new ingress rule to `/etc/cloudflared/config.yml`
2. Create DNS route: `cloudflared tunnel route dns onemind newservice.zeusops.dev`
3. Restart cloudflared: `sudo systemctl restart cloudflared`
