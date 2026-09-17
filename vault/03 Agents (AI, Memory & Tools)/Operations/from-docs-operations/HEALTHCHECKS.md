# Docker Health Checks Documentation

## Overview

This document describes the health check configuration for all OneMind OS services. Health checks ensure that services are not just running, but actually functional and ready to accept traffic.

## Health Check Summary

Total services: **18**  
Services with health checks: **16** (89%)  
Services without health checks: **2** (optional infrastructure services)

## Service Health Checks

### Database Services

#### postgres (pgvector/pgvector:pg15)
```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-onemind}"]
  interval: 10s
  timeout: 5s
  retries: 5
```
- **Method**: PostgreSQL readiness check
- **Validates**: Database is accepting connections

#### redis (redis:7-alpine)
```yaml
healthcheck:
  test: ["CMD", "redis-cli", "ping"]
  interval: 10s
  timeout: 5s
  retries: 5
```
- **Method**: Redis PING command
- **Validates**: Redis is responsive

#### timescaledb (timescale/timescaledb:latest-pg15)
```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U onemind"]
  interval: 10s
  timeout: 5s
  retries: 5
```
- **Method**: PostgreSQL readiness check
- **Validates**: TimescaleDB is accepting connections

---

### Core Services

#### nats (nats:latest)
```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "-q", "http://localhost:8222/healthz"]
  interval: 10s
  timeout: 5s
  retries: 5
  start_period: 5s
```
- **Method**: HTTP check on monitoring endpoint
- **Port**: 8222
- **Validates**: NATS server is running and healthy

#### agno (OneMind AgentOS)
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 5
  start_period: 30s
```
- **Method**: HTTP check on health endpoint
- **Port**: 8080
- **Validates**: AgentOS API is responsive
- **Dependencies**: postgres, nats

#### legacy-voice (LiveKit Voice Agent)
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```
- **Method**: HTTP check on health endpoint
- **Port**: 8080
- **Validates**: Voice agent is ready
- **Dependencies**: nats, agno, redis

#### api-gateway (FastAPI Gateway)
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 20s
```
- **Method**: HTTP check on health endpoint
- **Port**: 8000
- **Validates**: Gateway is accepting webhooks
- **Dependencies**: nats

#### action-gateway (Action Executor)
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8001/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 20s
```
- **Method**: HTTP check on health endpoint
- **Port**: 8001
- **Validates**: Action gateway is ready
- **Dependencies**: nats

#### event-processor (Event Handler)
```yaml
healthcheck:
  test: ["CMD", "pgrep", "-f", "python"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 20s
```
- **Method**: Process check
- **Validates**: Event processor is running
- **Dependencies**: nats, agno

---

### UI Services

#### agent-ui (ghcr.io/agno-agi/agent-ui:latest)
```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "-q", "http://localhost:3000"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```
- **Method**: HTTP check on root
- **Port**: 3000
- **Validates**: UI is accessible
- **Dependencies**: agno

#### toc (Tactical Operations Center)
```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "-q", "http://localhost:3002"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 60s
```
- **Method**: HTTP check on root
- **Port**: 3002
- **Validates**: TOC UI is accessible
- **Dependencies**: agno, legacy-voice

#### n8n (n8nio/n8n:latest)
```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "-q", "http://localhost:5678/healthz"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```
- **Method**: HTTP check on healthz endpoint
- **Port**: 5678
- **Validates**: N8N workflow engine is ready

#### uptime-kuma (louislam/uptime-kuma:latest)
```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "-q", "http://localhost:3001"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```
- **Method**: HTTP check on root
- **Port**: 3001
- **Validates**: Monitoring UI is accessible

---

### Infrastructure Services

#### caddy (caddy:2-alpine)
```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "-q", "http://localhost:2019/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 20s
```
- **Method**: HTTP check on admin API health endpoint
- **Port**: 2019
- **Validates**: Caddy reverse proxy is running
- **Dependencies**: api-gateway, n8n, agno

#### kestra (kestra/kestra:latest)
```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 60s
```
- **Method**: HTTP check on health endpoint
- **Port**: 8080
- **Validates**: Kestra workflow engine is ready
- **Dependencies**: postgres

#### portainer (portainer/portainer-ce:latest)
```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "-q", "http://localhost:9000"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 30s
```
- **Method**: HTTP check on root
- **Port**: 9000
- **Validates**: Portainer UI is accessible

---

### Services Without Health Checks

#### tailscale
- **Reason**: System-level VPN service that doesn't expose a health endpoint
- **Type**: Optional secure access infrastructure

#### novnc
- **Reason**: Optional desktop access service
- **Type**: Development/debugging tool

---

## Health Check Parameters Explained

### `test`
The command to run for health checking. Can be:
- `CMD`: Direct command execution
- `CMD-SHELL`: Shell command execution

### `interval`
Time between health checks (default: 30s)
- Database services: 10s (faster detection)
- Application services: 30s (balanced)

### `timeout`
Maximum time for health check to complete
- Database services: 5s
- Application services: 10s

### `retries`
Number of consecutive failures before marking unhealthy
- Database services: 5 (more tolerance)
- Application services: 3-5

### `start_period`
Grace period after container start before health checks begin
- Fast services (databases): 5s
- Medium services (APIs): 20-30s
- Slow services (UIs, TOC): 40-60s

---

## Dependency Chain

Services wait for their dependencies to be healthy using `condition: service_healthy`:

```
postgres (ready first)
  ├─→ agno
  │    ├─→ legacy-voice
  │    │    └─→ toc
  │    ├─→ event-processor
  │    ├─→ agent-ui
  │    └─→ caddy
  └─→ kestra

redis
  └─→ legacy-voice

nats
  ├─→ agno
  ├─→ api-gateway
  │    └─→ caddy
  ├─→ action-gateway
  └─→ event-processor
```

---

## Testing Health Checks

### View Service Status
```bash
docker compose ps
```

Shows health status for all services:
- `healthy`: Service passed health check
- `unhealthy`: Service failed health checks
- `starting`: Within start_period, health checks not yet enforced

### View Specific Service Health
```bash
docker inspect <container-name> | jq '.[0].State.Health'
```

### Monitor Health Check Logs
```bash
docker compose logs -f <service-name>
```

### Manual Health Check Test
```bash
# Test postgres
docker exec onemind-postgres pg_isready -U onemind && echo "OK"

# Test redis
docker exec onemind-redis redis-cli ping

# Test agno
curl -f http://localhost:8080/health && echo "OK"

# Test nats
wget --spider -q http://localhost:8222/healthz && echo "OK"
```

---

## Troubleshooting

### Service Stuck in "starting" State
- Increase `start_period` if service needs more initialization time
- Check logs: `docker compose logs -f <service>`

### Service Constantly Restarting
- Health check may be too aggressive (reduce retries or increase timeout)
- Service may have actual issues (check application logs)

### Service Marked Unhealthy But Works
- Health check command may be incorrect
- Health endpoint may be slow (increase timeout)

---

## Validation

Run the validation script to verify all health checks:

```bash
cd docker_substrate_runtime/cloud
./validate-healthchecks.sh
```

Expected output:
```
✅ docker-compose.yml syntax is valid
🏥 16 services with health checks configured
✅ All critical services have health checks configured!
```

---

## Best Practices

1. **Always add health checks** to services that expose HTTP endpoints
2. **Use appropriate intervals** - faster for critical services, slower for UI services
3. **Set realistic start_period** - services need time to initialize
4. **Use service_healthy conditions** - ensure proper startup order
5. **Test health checks** before deploying to production
6. **Monitor health status** using `docker compose ps` or Uptime Kuma

---

## References

- [Docker Compose Health Check Documentation](https://docs.docker.com/compose/compose-file/compose-file-v3/#healthcheck)
- [Docker HEALTHCHECK Instruction](https://docs.docker.com/engine/reference/builder/#healthcheck)
- [Best Practices for Health Checks](https://docs.docker.com/engine/reference/builder/#healthcheck)
