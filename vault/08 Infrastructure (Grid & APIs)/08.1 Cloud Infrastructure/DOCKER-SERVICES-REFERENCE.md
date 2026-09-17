# Docker Services Reference

*OneMind Cloud Container Documentation - January 2026*

## Overview

OneMind Cloud runs 13+ Docker containers via Docker Compose. This document details each service, its configuration, health checks, and troubleshooting procedures.

---

## Compose File Location

```
/home/ubuntu/onemind-os/docker_substrate_runtime/cloud/docker-compose.yml
```

---

## Service Categories

```
┌─────────────────────────────────────────────────────────────────┐
│                    SERVICE CATEGORIES                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   AI SERVICES              INFRASTRUCTURE         DATA STORES   │
│   ────────────            ──────────────         ───────────   │
│   • Agno (LEGACY)         • Portainer            • PostgreSQL  │
│   • Letta                 • Uptime Kuma          • Qdrant      │
│   • N8N                   • TOC                  • Redis       │
│   • Kestra                • VNC                  • NATS        │
│   • Actions               • (Caddy - removed)                  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## AI Services

### Agno (LEGACY Brain)

The primary AI agent runtime powering LEGACY.

| Property | Value |
|----------|-------|
| **Container** | onemind-agno |
| **Port** | 8080 |
| **Image** | Custom build |
| **Memory** | 2GB limit |
| **Depends On** | postgres, redis, qdrant |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "python", "-c", "import urllib.request; urllib.request.urlopen('http://localhost:8080/health')"]
  interval: 30s
  timeout: 10s
  retries: 3
```

**Environment Variables:**
- `DATABASE_URL` - PostgreSQL connection
- `REDIS_URL` - Redis connection
- `QDRANT_URL` - Vector DB connection
- `OPENAI_API_KEY` - LLM API key

**Logs:**
```bash
docker logs onemind-agno --tail 100 -f
```

---

### Letta (Memory Server)

Long-term memory and context management for AI agents.

| Property | Value |
|----------|-------|
| **Container** | onemind-letta |
| **Port** | 8283 |
| **Image** | letta/letta:latest |
| **Memory** | 4GB limit |
| **Depends On** | postgres, qdrant |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "python", "-c", "import urllib.request; urllib.request.urlopen('http://localhost:8283/v1/health')"]
  interval: 30s
  timeout: 10s
  retries: 3
```

**Volumes:**
- `./data/letta:/root/.letta` - Persistent memory storage

**Key Endpoints:**
- `GET /v1/health` - Health check
- `GET /v1/agents` - List agents
- `POST /v1/agents/{id}/messages` - Send message

---

### N8N (Workflow Automation)

Visual workflow automation for integrations.

| Property | Value |
|----------|-------|
| **Container** | onemind-n8n |
| **Port** | 5678 |
| **Image** | n8nio/n8n:latest |
| **Memory** | 2GB limit |
| **Depends On** | postgres, redis |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "node", "-e", "require('http').get('http://localhost:5678/healthz', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"]
  interval: 30s
  timeout: 10s
  retries: 3
```

**Volumes:**
- `./data/n8n:/home/node/.n8n` - Workflows and credentials

**Environment Variables:**
- `N8N_HOST` - Hostname
- `N8N_PORT` - Port (5678)
- `N8N_PROTOCOL` - https
- `WEBHOOK_URL` - External webhook URL
- `DB_TYPE` - postgresdb
- `DB_POSTGRESDB_*` - Database connection

---

### Kestra (Data Pipelines)

Event-driven data orchestration platform.

| Property | Value |
|----------|-------|
| **Container** | onemind-kestra |
| **Port** | 8088 (UI), 8081 (health) |
| **Image** | kestra/kestra:latest |
| **Memory** | 2GB limit |
| **Depends On** | postgres |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
  interval: 30s
  timeout: 10s
  retries: 3
```

> **Note:** Kestra health endpoint is on port 8081, not 8088 (UI port).

**Volumes:**
- `./data/kestra:/app/storage` - Flow definitions and data

---

### Actions (Action Gateway)

Custom API gateway for AI actions.

| Property | Value |
|----------|-------|
| **Container** | onemind-actions |
| **Port** | 8001 |
| **Image** | Custom build |
| **Memory** | 1GB limit |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "python", "-c", "import urllib.request; urllib.request.urlopen('http://localhost:8001/health')"]
  interval: 30s
  timeout: 10s
  retries: 3
```

> **Note:** Root endpoint (/) returns 404 - this is expected. Use /health for checks.

---

## Infrastructure Services

### Portainer

Docker management web UI.

| Property | Value |
|----------|-------|
| **Container** | onemind-portainer |
| **Port** | 9000 |
| **Image** | portainer/portainer-ce:latest |
| **Memory** | 256MB limit |

**Health Check:**
```yaml
healthcheck:
  test: ["NONE"]
```

> **Note:** Portainer uses a scratch-based image with no shell tools. Health check disabled.

**Volumes:**
- `/var/run/docker.sock:/var/run/docker.sock` - Docker socket
- `./data/portainer:/data` - Portainer data

---

### Uptime Kuma

Service monitoring and status page.

| Property | Value |
|----------|-------|
| **Container** | onemind-uptime |
| **Port** | 3001 |
| **Image** | louislam/uptime-kuma:latest |
| **Memory** | 512MB limit |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "node", "-e", "require('http').get('http://localhost:3001/api/status', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"]
  interval: 30s
  timeout: 10s
  retries: 3
```

> **Note:** Root (/) returns 302 redirect. Use /api/status for health checks.

**Volumes:**
- `./data/uptime-kuma:/app/data` - Monitoring config

---

### TOC (Table of Contents)

Dashboard and navigation interface.

| Property | Value |
|----------|-------|
| **Container** | onemind-toc |
| **Port** | 3002 |
| **Image** | Custom build |
| **Memory** | 512MB limit |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:3002/"]
  interval: 30s
  timeout: 10s
  retries: 3
```

---

### VNC (noVNC)

Browser-based remote desktop access.

| Property | Value |
|----------|-------|
| **Container** | onemind-vnc |
| **Port** | 6080 |
| **Image** | theasp/novnc:latest |
| **Memory** | 1GB limit |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:6080/"]
  interval: 30s
  timeout: 10s
  retries: 3
```

---

## Data Services

### PostgreSQL

Primary relational database.

| Property | Value |
|----------|-------|
| **Container** | onemind-postgres |
| **Port** | 5432 |
| **Image** | postgres:15-alpine |
| **Memory** | 1GB limit |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres"]
  interval: 10s
  timeout: 5s
  retries: 5
```

**Volumes:**
- `./data/postgres:/var/lib/postgresql/data`

**Environment Variables:**
- `POSTGRES_USER` - postgres
- `POSTGRES_PASSWORD` - (from .env)
- `POSTGRES_DB` - onemind

---

### Qdrant

Vector database for AI embeddings.

| Property | Value |
|----------|-------|
| **Container** | onemind-qdrant |
| **Port** | 6333 (HTTP), 6334 (gRPC) |
| **Image** | qdrant/qdrant:latest |
| **Memory** | 2GB limit |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:6333/"]
  interval: 30s
  timeout: 10s
  retries: 3
```

**Volumes:**
- `./data/qdrant:/qdrant/storage`

---

### Redis

Cache and pub/sub messaging.

| Property | Value |
|----------|-------|
| **Container** | onemind-redis |
| **Port** | 6379 |
| **Image** | redis:7-alpine |
| **Memory** | 512MB limit |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "redis-cli", "ping"]
  interval: 10s
  timeout: 5s
  retries: 5
```

**Volumes:**
- `./data/redis:/data`

---

### NATS

High-performance message queue.

| Property | Value |
|----------|-------|
| **Container** | onemind-nats |
| **Port** | 4222 (client), 8222 (monitor) |
| **Image** | nats:latest |
| **Memory** | 256MB limit |

**Health Check:**
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8222/"]
  interval: 30s
  timeout: 10s
  retries: 3
```

**Key Endpoints:**
- Port 4222 - Client connections
- Port 8222 - Monitoring/metrics UI

---

## Stopped/Removed Services

### Voice (Stopped)
- **Status:** Stopped due to Python event loop bug
- **Error:** `RuntimeError: This event loop is already running`
- **Fix:** Requires code changes to async handling

### Agent-UI (Stopped)
- **Status:** Not configured
- **Port:** 3000
- **Note:** Would serve command.onemindos.dev

### Caddy (Removed)
- **Status:** Removed as redundant
- **Reason:** Cloudflare Tunnel handles all routing
- **Note:** Was adding unnecessary complexity

---

## Health Check Patterns

Different images require different health check approaches:

| Image Type | Health Check Method | Example |
|------------|---------------------|---------|
| **Python-based** | `python -c "import urllib..."` | agno, letta, actions |
| **Node.js-based** | `node -e "require('http')..."` | n8n, uptime-kuma |
| **Alpine/curl** | `curl -f http://...` | kestra, qdrant |
| **Scratch images** | `NONE` | portainer |
| **Database** | Native CLI | postgres, redis |

### Why Not Just Use curl?

Many slim Docker images don't include curl or wget to reduce size:
- Alpine images: Often have curl
- Distroless: No shell at all
- Scratch: Literally empty

Using the language runtime (Python/Node) that's already in the image is more reliable.

---

## Common Operations

### View All Container Status
```bash
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
```

### Check Specific Container Health
```bash
docker inspect --format='{{.State.Health.Status}}' onemind-agno
```

### View Container Logs
```bash
docker logs onemind-n8n --tail 100 -f
```

### Restart Single Service
```bash
docker restart onemind-agno
```

### Restart All Services
```bash
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud
docker compose restart
```

### Full Stack Restart
```bash
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud
docker compose down
docker compose up -d
```

### Check Resource Usage
```bash
docker stats --no-stream
```

---

## Troubleshooting

### Container Won't Start
1. Check logs: `docker logs onemind-<name>`
2. Check dependencies: `docker ps` (are postgres/redis running?)
3. Check disk space: `df -h`
4. Check memory: `free -h`

### Container Unhealthy
1. Check health command: `docker inspect onemind-<name> | grep -A 20 Health`
2. Run health check manually inside container
3. Check if service actually responds on expected port

### Container Restart Loop
1. Check exit code: `docker inspect onemind-<name> --format='{{.State.ExitCode}}'`
2. Check logs for crash reason
3. Common causes: missing env vars, port conflicts, dependency not ready

### Port Already in Use
```bash
# Find what's using the port
sudo lsof -i :8080
# Or
sudo netstat -tlnp | grep 8080
```
