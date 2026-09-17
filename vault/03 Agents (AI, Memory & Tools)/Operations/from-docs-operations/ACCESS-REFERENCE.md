# OneMind OS - Access Reference

> Last Updated: 2026-01-30
> All access points, credentials, and connection information

---

## Quick Access URLs

### Active Services

| Service | URL | Type | Status |
|---------|-----|------|--------|
| **UI Dashboard** | https://command.onemindos.dev | GUI | Active |
| **API** | https://api.onemindos.dev | API | Active |
| **Portainer** | https://portainer.onemindos.dev | GUI | Active |
| **NATS Monitor** | https://nats.onemindos.dev | GUI | Active |

---

## Server Access

### Azure VM (Primary)

| Method | Address | Notes |
|--------|---------|-------|
| Public IP | 20.121.38.186 | Azure VM |
| Tailscale IP | 100.102.21.44 | Internal network |
| SSH | `ssh zeus@100.102.21.44` | Via Tailscale |
| Hostname | onemind-vm | |

**Services Running:**
- Flutter UI (port 3000)
- Agno API (port 8081)
- PostgreSQL + pgvector (port 5432)
- Redis (port 6379)
- NATS (ports 4222, 8222)
- TimescaleDB (port 5433)
- Portainer Business (port 9000) - GitOps auto-deploy via webhook
- Cloudflare Tunnel (cloudflared systemd service)
- Tailscale (tailscaled systemd service)

---

## API Endpoints

### Agno API (api.onemindos.dev)

```bash
# Health check
curl https://api.onemindos.dev/health

# List agents
curl https://api.onemindos.dev/agents

# List sessions
curl https://api.onemindos.dev/sessions

# Get awareness state
curl https://api.onemindos.dev/awareness

# Run agent
curl -X POST https://api.onemindos.dev/agents/legacy/runs \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello"}'
```

---

## Database Access

### PostgreSQL

```bash
# From VM
docker exec -it onemind-postgres psql -U onemind -d onemind

# Remote via Tailscale
psql "postgresql://onemind:PASSWORD@100.102.21.44:5432/onemind"
```

### TimescaleDB

```bash
# From VM
docker exec -it onemind-timescaledb psql -U onemind -d events

# Remote via Tailscale
psql "postgresql://onemind:PASSWORD@100.102.21.44:5433/events"
```

### Redis

```bash
# From VM
docker exec -it onemind-redis redis-cli

# Remote via Tailscale
redis-cli -h 100.102.21.44 -p 6379
```

---

## Credentials Storage

All credentials are stored in 1Password:
- **Vault:** `00-24 SO (Sovereign Operations)`
- **Items:**
  - `OneMind Cloud Server Access` - SSH keys
  - `OneMind Secrets` - API keys and passwords
  - `Cloudflare - OneMind OS` - DNS/Tunnel credentials
  - `Cloudflare Tunnel - OneMind OS` - Tunnel token
  - `AWS Bedrock Credentials` - Claude API access

---

## Port Reference

### Azure VM Ports

| Port | Service | Protocol | Access |
|------|---------|----------|--------|
| 22 | SSH | TCP | Tailscale only |
| 3000 | Flutter UI | TCP | Via Cloudflare tunnel |
| 4222 | NATS Client | TCP | Internal |
| 5432 | PostgreSQL | TCP | Internal |
| 5433 | TimescaleDB | TCP | Internal |
| 6379 | Redis | TCP | Internal |
| 8081 | Agno API | TCP | Via Cloudflare tunnel |
| 8222 | NATS Monitor | TCP | Via Cloudflare tunnel |
| 9000 | Portainer | TCP | Via Cloudflare tunnel |

---

## Cloudflare Tunnel Configuration

### Tunnel Details
- **Tunnel ID:** 7ddfa54f-c672-4dd0-ad50-19b431b5ad72
- **Tunnel Name:** onemind-os
- **Running On:** Azure VM (100.102.21.44)
- **Service:** systemd (cloudflared.service)

### Current Ingress Routes

| Hostname | Service |
|----------|---------|
| api.onemindos.dev | localhost:8081 |
| ui.onemindos.dev | localhost:3000 |
| command.onemindos.dev | localhost:3000 |
| nats.onemindos.dev | localhost:8222 |
| portainer.onemindos.dev | localhost:9000 |

### DNS Zone
- **Domain:** onemindos.dev
- **Zone ID:** 32beabc2a2f2debece7b5c55cfdd9b2a
- **Account ID:** ce13703e0b76938c179ebc51e2945c7e

---

## Tailscale Network

| Device | IP | Notes |
|--------|----|----|
| onemindos (Azure VM) | 100.102.21.44 | Primary server |
| zeuss-macbook-pro | 100.114.142.55 | Dev machine |
| iphone182 | 100.122.203.68 | Mobile |

---

## Service Management

### Cloudflare Tunnel
```bash
# Check status
sudo systemctl status cloudflared

# View logs
journalctl -u cloudflared -f

# Restart
sudo systemctl restart cloudflared
```

### Tailscale
```bash
# Check status
tailscale status

# View logs
journalctl -u tailscaled -f

# Key expiry
tailscale status --json | jq '.Self.KeyExpiry'
```

### Docker Containers
```bash
# View all containers
docker ps -a

# View logs
docker logs -f onemind-agno

# Restart a service
docker restart onemind-agno
```

---

## Useful Commands

```bash
# SSH to Azure VM
ssh zeus@100.102.21.44

# Check all services health
curl https://api.onemindos.dev/health

# View container resource usage
docker stats --no-stream

# Check disk usage
df -h /

# Check memory
free -h
```
