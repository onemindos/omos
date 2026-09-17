---
title: "CLI — docker"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/docker, topic/infra]
ai_summary: "Docker and Docker Compose CLI quick reference for the One Mind VPS stack — containers, services, logs, and deployment commands."
---

# docker — Container CLI

> All containers run on VPS `onemind-main` (104.248.4.212). Compose file: `/opt/onemind/docker-compose.yml`.

---

## Running Containers

| Container | Image | Ports |
|-----------|-------|-------|
| `onemind-nats` | `nats:latest` | 4222 (core), 8222 (monitor), 1883 (MQTT), 7422 (leaf), 9222 (WS) |
| `onemind-omnidb` | `timescale/timescaledb:latest-pg16` | 5432 (PostgreSQL + TimescaleDB + pgvector) |
| `onemind-n8n` | `n8nio/n8n:latest` | 5678 |
| `nats-dashboard` | `mdawar/nats-dashboard:latest` | 18300 → 8080 |

onemind-web deploys separately at `/opt/onemind-web/docker-compose.yml`.

---

## Core Commands

```bash
# List running containers
docker ps
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Container logs
docker logs onemind-nats
docker logs onemind-nats -f          # follow/live
docker logs onemind-nats --tail 50   # last 50 lines

# Exec into container
docker exec -it onemind-omnidb bash
docker exec -it onemind-nats sh

# Stop / start a container
docker stop onemind-nats
docker start onemind-nats
docker restart onemind-nats

# Remove container
docker rm onemind-nats

# Inspect container
docker inspect onemind-nats
```

---

## Docker Compose (VPS)

```bash
# Navigate to compose directory
cd /opt/onemind

# Start all services
docker compose up -d

# Stop all services
docker compose down

# Rebuild and restart
docker compose up -d --build

# Restart a single service
docker compose restart nats

# View logs for a service
docker compose logs -f nats
docker compose logs -f n8n

# Check status
docker compose ps
```

### onemind-web (separate compose)

```bash
cd /opt/onemind-web
git pull origin main
docker compose up -d --build
```

Note: Data dir must be owned by UID 1000 (`node` user). If fresh:
```bash
chown -R 1000:1000 /opt/onemind-web/data
```

---

## Images

```bash
# List images
docker images

# Pull latest image
docker pull nats:latest

# Remove unused images
docker image prune
```

---

## Volumes & Networks

```bash
# List volumes
docker volume ls

# Inspect a volume
docker volume inspect onemind_nats_data

# List networks
docker network ls
```

---

## Health Checks

```bash
# Check NATS health
curl http://104.248.4.212:8222/healthz

# Check n8n
curl http://104.248.4.212:5678/healthz

# Check PostgreSQL (via psql — see psql.md)
docker exec -it onemind-omnidb psql -U onemind -c "SELECT 1;"
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
