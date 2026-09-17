# OneMind OS - Failsafe & High Availability Guide

## Overview

This document covers failsafe strategies for running OneMind OS on a single VM with Docker Compose profiles, plus the path to production-grade high availability.

## Current Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    SINGLE VM (4-8GB RAM)                    │
│                                                             │
│  docker compose --profile core --profile utility up -d      │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  CORE PROFILE (Always Running)                      │   │
│  │  ┌────────┐ ┌───────┐ ┌──────┐ ┌──────┐ ┌────────┐ │   │
│  │  │Postgres│ │ Redis │ │ NATS │ │ Agno │ │   UI   │ │   │
│  │  └────────┘ └───────┘ └──────┘ └──────┘ └────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  UTILITY PROFILE (Optional)                         │   │
│  │  ┌────────┐ ┌─────────┐ ┌─────────┐ ┌───────────┐  │   │
│  │  │ Kestra │ │ Netdata │ │ Grafana │ │VictoriaM. │  │   │
│  │  └────────┘ └─────────┘ └─────────┘ └───────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  NOTIFICATIONS PROFILE (Optional)                    │   │
│  │  ┌──────┐ ┌─────────┐ ┌───────────┐                │   │
│  │  │ Novu │ │ MongoDB │ │ Novu-Web  │                │   │
│  │  └──────┘ └─────────┘ └───────────┘                │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Failsafe Layers

### Layer 1: Container Health Checks & Auto-Restart

Every service must have health checks. Docker will automatically restart unhealthy containers.

```yaml
services:
  postgres:
    image: timescale/timescaledb-ha:pg16
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  nats:
    image: nats:alpine
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8222/healthz"]
      interval: 10s
      timeout: 5s
      retries: 5

  agno:
    image: onemind-agno:latest
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 60s
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      nats:
        condition: service_healthy

  ui:
    image: onemind-ui:latest
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:3000"]
      interval: 30s
      timeout: 10s
      retries: 3
    depends_on:
      agno:
        condition: service_healthy
```

### Layer 2: Resource Limits

Prevent runaway containers from crashing the VM:

```yaml
services:
  agno:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 256M

  postgres:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '0.25'
          memory: 512M

  redis:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          cpus: '0.1'
          memory: 64M
```

### Layer 3: Automated Snapshots

#### Azure VM

```bash
# Use Azure CLI for snapshots
az snapshot create \
  --resource-group onemind-rg \
  --name "onemind-snap-$(date +%Y%m%d)" \
  --source /subscriptions/.../disks/onemind-disk

# List snapshots
az snapshot list --resource-group onemind-rg
```

#### Generic (Any Provider)

```bash
#!/bin/bash
# /opt/scripts/backup.sh - Run via cron

# Stop services gracefully
cd /opt/onemind-os/infra/compose/cloud
docker compose stop

# Backup volumes
tar -czf /backup/onemind-$(date +%Y%m%d-%H%M).tar.gz \
  /var/lib/docker/volumes/onemind_*

# Restart services
docker compose up -d

# Upload to S3/Backblaze
aws s3 cp /backup/onemind-*.tar.gz s3://onemind-backups/

# Cleanup old local backups (keep 7 days)
find /backup -name "onemind-*.tar.gz" -mtime +7 -delete
```

Cron schedule:
```bash
# /etc/cron.d/onemind-backup
0 3 * * * root /opt/scripts/backup.sh >> /var/log/backup.log 2>&1
```

### Layer 4: External Monitoring

Use Uptime Kuma or external service to monitor from outside:

```yaml
# Add to docker-compose.yml
uptime-kuma:
  image: louislam/uptime-kuma
  container_name: uptime-kuma
  ports:
    - "3001:3001"
  volumes:
    - uptime_data:/app/data
  restart: unless-stopped
```

Configure monitors:
- `https://onemindos.dev` - UI availability
- `https://api.onemindos.dev/health` - API health
- `tcp://100.x.x.x:5432` - Database (via Tailscale)

### Layer 5: Alerting

#### Via Novu (Self-hosted)

```python
# When health check fails, notify via multiple channels
async def on_service_down(service: str):
    await novu.trigger(
        workflow_id="service-alert",
        to={"subscriberId": "admin"},
        payload={
            "service": service,
            "severity": "critical",
            "message": f"{service} is down on onemind-cloud"
        }
    )
```

#### Via Apprise (Simpler)

```bash
# Install apprise
pip install apprise

# Send alert
apprise -t "Service Down" -b "Agno is not responding" \
  "tgram://bot_token/chat_id" \
  "slack://webhook_url" \
  "mailto://user:pass@gmail.com"
```

## Recovery Procedures

### Scenario 1: Single Container Crash

Docker auto-restarts. Check logs:
```bash
docker logs onemind-agno --tail 100
docker compose ps
```

### Scenario 2: VM Unresponsive

1. Hard reboot via cloud console
2. Docker will auto-start containers (restart: unless-stopped)
3. Check health: `docker compose ps`

### Scenario 3: Data Corruption

1. Stop services: `docker compose down`
2. Restore from snapshot:
   ```bash
   # Lightsail
   aws lightsail create-instance-from-snapshot \
     --instance-snapshot-name "latest-snapshot" \
     --instance-name onemind-cloud-restored \
     --availability-zone us-east-1a
   ```
3. Update DNS/Tailscale to new IP

### Scenario 4: Full Disaster Recovery

1. Provision new VM
2. Install Docker
3. Clone repo: `git clone https://github.com/user/onemind-os`
4. Restore volumes from S3:
   ```bash
   aws s3 cp s3://onemind-backups/latest.tar.gz /tmp/
   tar -xzf /tmp/latest.tar.gz -C /var/lib/docker/volumes/
   ```
5. Start services: `docker compose --profile core up -d`

## Production Scaling Path

### Phase 1: Single VM (Current)

```
Cost: ~$20/mo
Uptime: 99.5% (manual recovery)
Recovery: 5-15 minutes
```

### Phase 2: Two VMs

```
┌─────────────────┐     ┌─────────────────┐
│    AI VM        │     │   Utility VM    │
│   (Primary)     │◄───►│   (Secondary)   │
│                 │NATS │                 │
│  Core services  │     │  Heavy services │
└─────────────────┘     └─────────────────┘

Cost: ~$40/mo
Uptime: 99.9% (one VM can continue)
Recovery: 1-2 minutes
```

### Phase 3: Multi-Region

```
┌─────────────────┐     ┌─────────────────┐
│   US-East       │     │   EU-West       │
│   (Primary)     │◄───►│   (DR)          │
│                 │NATS │                 │
└─────────────────┘     └─────────────────┘

Cost: ~$80/mo
Uptime: 99.99%
Recovery: Automatic failover
```

## Docker Compose with Profiles

```yaml
# docker-compose.yml with profiles
version: "3.8"

services:
  # ===== CORE PROFILE (Always needed) =====
  postgres:
    profiles: ["core"]
    image: timescale/timescaledb-ha:pg16
    # ... config

  redis:
    profiles: ["core"]
    image: redis:7-alpine
    # ... config

  nats:
    profiles: ["core"]
    image: nats:alpine
    # ... config

  agno:
    profiles: ["core"]
    image: onemind-agno:latest
    # ... config

  ui:
    profiles: ["core"]
    image: onemind-ui:latest
    # ... config

  # ===== UTILITY PROFILE (Optional) =====
  netdata:
    profiles: ["utility"]
    image: netdata/netdata
    # ... config

  victoriametrics:
    profiles: ["utility"]
    image: victoriametrics/victoria-metrics
    # ... config

  grafana:
    profiles: ["utility"]
    image: grafana/grafana
    # ... config

  kestra:
    profiles: ["utility"]
    image: kestra/kestra
    # ... config

  # ===== NOTIFICATIONS PROFILE (Optional) =====
  novu:
    profiles: ["notifications"]
    image: ghcr.io/novuhq/novu/api
    # ... config

  novu-web:
    profiles: ["notifications"]
    image: ghcr.io/novuhq/novu/web
    # ... config

  mongo:
    profiles: ["notifications"]
    image: mongo:6
    # ... config

  # ===== EDGE PROFILE (For Jetsons) =====
  ollama:
    profiles: ["edge"]
    image: ollama/ollama
    # ... config

  frigate:
    profiles: ["edge"]
    image: ghcr.io/blakeblackshear/frigate:stable
    # ... config

  nats-leaf:
    profiles: ["edge"]
    image: nats:alpine
    # ... config
```

### Usage

```bash
# Start core only (minimal)
docker compose --profile core up -d

# Start core + observability
docker compose --profile core --profile utility up -d

# Start everything
docker compose --profile core --profile utility --profile notifications up -d

# On Jetson edge devices
docker compose --profile edge up -d
```

## Monitoring Dashboard

Access points after deployment:

| Service | URL | Purpose |
|---------|-----|---------|
| OneMind UI | https://onemindos.dev | Main application |
| Netdata | http://100.x.x.x:19999 | Real-time metrics |
| Grafana | http://100.x.x.x:3001 | Custom dashboards |
| Uptime Kuma | http://100.x.x.x:3002 | Uptime monitoring |
| Novu | http://100.x.x.x:4200 | Notifications |
| pgAdmin | http://100.x.x.x:5050 | Database admin |

## Checklist

- [ ] All containers have health checks
- [ ] All containers have restart: unless-stopped
- [ ] Resource limits configured
- [ ] Automated snapshots enabled
- [ ] External monitoring configured
- [ ] Alerting configured (Novu/Apprise)
- [ ] Backup script in cron
- [ ] Recovery procedure tested
- [ ] Runbook documented
