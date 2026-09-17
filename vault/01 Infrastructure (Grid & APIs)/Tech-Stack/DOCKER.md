# Docker Reference Guide

> Comprehensive Docker and Docker Compose reference for OneMind OS infrastructure.

---

## Table of Contents

- [Images](#images)
- [Containers](#containers)
- [Docker Compose Services](#docker-compose-services)
- [Restart Policies](#restart-policies)
- [Dependencies (depends_on)](#dependencies-depends_on)
- [Environment Variables](#environment-variables)
- [Naming Conventions](#naming-conventions)
- [Volumes](#volumes)
- [Ports](#ports)
- [Networks](#networks)
- [Health Checks](#health-checks)
- [Resource Limits](#resource-limits)
- [Logging](#logging)
- [Build Configuration](#build-configuration)
- [Secrets Management](#secrets-management)
- [Profiles](#profiles)
- [Common Commands](#common-commands)
- [Troubleshooting](#troubleshooting)

---

## Images

### Pulling Images

```bash
# Pull specific image
docker pull postgres:16-alpine

# Pull with digest (immutable)
docker pull postgres@sha256:abc123...

# Pull all images in compose file
docker compose pull
```

### Listing Images

```bash
# List all images
docker images

# List with filters
docker images --filter "dangling=true"
docker images --format "{{.Repository}}:{{.Tag}} {{.Size}}"
```

### Building Images

```bash
# Build from Dockerfile
docker build -t onemind/legacy-agent:latest .

# Build with build args
docker build --build-arg VERSION=1.0.0 -t onemind/legacy-agent:1.0.0 .

# Multi-platform build
docker buildx build --platform linux/amd64,linux/arm64 -t onemind/legacy-agent:latest .
```

### Image Best Practices

| Practice             | Example                                      |
| -------------------- | -------------------------------------------- |
| Use specific tags    | `postgres:16.1-alpine` not `postgres:latest` |
| Use alpine variants  | Smaller images, faster pulls                 |
| Pin digests for prod | Immutable, reproducible                      |
| Multi-stage builds   | Smaller final images                         |

### Compose Example

```yaml
services:
  legacy-agent:
    image: onemind/legacy-agent:${VERSION:-latest}
```

---

## Containers

### Container Lifecycle

```bash
# Create and start
docker run -d --name legacy-agent onemind/legacy-agent:latest

# Start existing container
docker start legacy-agent

# Stop container (graceful)
docker stop legacy-agent

# Stop with timeout
docker stop -t 30 legacy-agent

# Kill container (force)
docker kill legacy-agent

# Remove container
docker rm legacy-agent

# Remove running container
docker rm -f legacy-agent
```

### Container Inspection

```bash
# List running containers
docker ps

# List all containers (including stopped)
docker ps -a

# Inspect container details
docker inspect legacy-agent

# View container logs
docker logs legacy-agent
docker logs -f --tail 100 legacy-agent

# Execute command in container
docker exec -it legacy-agent bash
docker exec legacy-agent python --version

# View resource usage
docker stats legacy-agent
```

### Container Best Practices

- Run as non-root user when possible
- Use `--init` for proper signal handling
- Set memory and CPU limits
- Use read-only filesystem when possible

---

## Docker Compose Services

### Basic Service Definition

```yaml
services:
  legacy-agent:
    image: onemind/legacy-agent:latest
    container_name: onemind-legacy-agent
    restart: unless-stopped
    environment:
      - LOG_LEVEL=info
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    networks:
      - onemind-network
    depends_on:
      - postgres
      - nats
```

### Service Commands

```bash
# Start all services
docker compose up -d

# Start specific service
docker compose up -d legacy-agent

# Stop all services
docker compose down

# Stop and remove volumes
docker compose down -v

# Restart service
docker compose restart legacy-agent

# View logs
docker compose logs -f legacy-agent

# Scale service
docker compose up -d --scale worker=3
```

---

## Restart Policies

### Available Policies

| Policy           | Behavior                                    |
| ---------------- | ------------------------------------------- |
| `no`             | Never restart (default)                     |
| `always`         | Always restart, including on daemon restart |
| `on-failure`     | Restart only on non-zero exit code          |
| `on-failure:5`   | Restart on failure, max 5 attempts          |
| `unless-stopped` | Like always, but not if manually stopped    |

### Compose Examples

```yaml
services:
  # Production services - always restart
  postgres:
    image: postgres:16-alpine
    restart: unless-stopped

  # Background workers - restart on failure
  worker:
    image: onemind/worker:latest
    restart: on-failure:5

  # One-off jobs - never restart
  migration:
    image: onemind/migration:latest
    restart: "no"
```

### Docker Run

```bash
docker run -d --restart unless-stopped --name legacy-agent onemind/legacy-agent
```

---

## Dependencies (depends_on)

### Basic Dependency

```yaml
services:
  legacy-agent:
    depends_on:
      - postgres
      - nats
```

### Conditional Dependencies

```yaml
services:
  legacy-agent:
    depends_on:
      postgres:
        condition: service_healthy
      nats:
        condition: service_started
      migration:
        condition: service_completed_successfully
```

### Dependency Conditions

| Condition                        | Description                            |
| -------------------------------- | -------------------------------------- |
| `service_started`                | Wait for container to start (default)  |
| `service_healthy`                | Wait for health check to pass          |
| `service_completed_successfully` | Wait for container to exit with code 0 |

### Example with Health Checks

```yaml
services:
  postgres:
    image: postgres:16-alpine
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  legacy-agent:
    depends_on:
      postgres:
        condition: service_healthy
```

---

## Environment Variables

### Methods to Set

#### 1. Inline in Compose

```yaml
services:
  legacy-agent:
    environment:
      - DATABASE_URL=postgresql://user:pass@postgres:5432/db
      - LOG_LEVEL=info
      - DEBUG=false
```

#### 2. Key-Value Format

```yaml
services:
  legacy-agent:
    environment:
      DATABASE_URL: postgresql://user:pass@postgres:5432/db
      LOG_LEVEL: info
      DEBUG: "false"
```

#### 3. Environment File

```yaml
services:
  legacy-agent:
    env_file:
      - .env
      - .env.local
```

#### 4. Variable Substitution

```yaml
services:
  legacy-agent:
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - LOG_LEVEL=${LOG_LEVEL:-info} # Default value
```

### .env File Format

```bash
# .env
DATABASE_URL=postgresql://user:pass@postgres:5432/db
LOG_LEVEL=info
DEBUG=false

# Can reference other variables
API_URL=http://localhost:${API_PORT:-8080}
```

### Best Practices

| Practice             | Description                        |
| -------------------- | ---------------------------------- |
| Use `.env.example`   | Template for required variables    |
| Never commit secrets | Use `.gitignore` for `.env`        |
| Use defaults         | `${VAR:-default}` syntax           |
| Validate at startup  | Fail fast on missing required vars |

---

## Naming Conventions

### Container Names

```yaml
services:
  legacy-agent:
    container_name: onemind-legacy-agent
```

### OneMind Naming Pattern

```
onemind-<service>-<environment>
```

Examples:

- `onemind-legacy-agent-prod`
- `onemind-postgres-dev`
- `onemind-nats-home`

### Project Name

```bash
# Set via environment
COMPOSE_PROJECT_NAME=onemind docker compose up -d

# Set via command line
docker compose -p onemind up -d
```

### Compose File

```yaml
name: onemind

services:
  legacy-agent:
    # ...
```

---

## Volumes

### Volume Types

| Type          | Use Case                           |
| ------------- | ---------------------------------- |
| Named volumes | Persistent data, managed by Docker |
| Bind mounts   | Development, config files          |
| tmpfs         | Temporary data, secrets            |

### Named Volumes

```yaml
services:
  postgres:
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
    driver: local
```

### Bind Mounts

```yaml
services:
  legacy-agent:
    volumes:
      # Short syntax
      - ./config:/app/config:ro
      - ./data:/app/data

      # Long syntax
      - type: bind
        source: ./logs
        target: /app/logs
```

### tmpfs Mounts

```yaml
services:
  legacy-agent:
    tmpfs:
      - /tmp
      - /run
```

### Volume Options

```yaml
volumes:
  postgres_data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: /path/on/host
```

### Volume Commands

```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect postgres_data

# Remove unused volumes
docker volume prune

# Backup volume
docker run --rm -v postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_data.tar.gz /data
```

---

## Ports

### Port Mapping Syntax

```yaml
services:
  legacy-agent:
    ports:
      # HOST:CONTAINER
      - "8080:8080"

      # Bind to specific interface
      - "127.0.0.1:8080:8080"

      # Range of ports
      - "8000-8010:8000-8010"

      # UDP port
      - "514:514/udp"

      # Long syntax
      - target: 8080
        published: 8080
        protocol: tcp
        mode: host
```

### Expose vs Ports

```yaml
services:
  # Ports - accessible from host
  api:
    ports:
      - "8080:8080"

  # Expose - only accessible from other containers
  database:
    expose:
      - "5432"
```

### OneMind Port Assignments

| Service      | Internal Port | External Port |
| ------------ | ------------- | ------------- |
| Legacy Agent | 8080          | 8080          |
| Gateway      | 8000          | 8000          |
| PostgreSQL   | 5432          | 5432          |
| NATS         | 4222          | 4222          |
| NATS HTTP    | 8222          | 8222          |
| Redis        | 6379          | 6379          |
| Qdrant       | 6333          | 6333          |
| LiveKit      | 7880          | 7880          |

---

## Networks

### Network Types

| Type      | Description                   |
| --------- | ----------------------------- |
| `bridge`  | Default, isolated network     |
| `host`    | Use host's network stack      |
| `none`    | No networking                 |
| `overlay` | Multi-host networking (Swarm) |

### Defining Networks

```yaml
services:
  legacy-agent:
    networks:
      - frontend
      - backend

  postgres:
    networks:
      - backend

networks:
  frontend:
    driver: bridge
  backend:
    driver: bridge
    internal: true # No external access
```

### Network Configuration

```yaml
networks:
  onemind-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.28.0.0/16
          gateway: 172.28.0.1
```

### External Networks

```yaml
networks:
  existing-network:
    external: true
    name: my-existing-network
```

### Network Commands

```bash
# List networks
docker network ls

# Inspect network
docker network inspect onemind-network

# Connect container to network
docker network connect onemind-network legacy-agent

# Disconnect
docker network disconnect onemind-network legacy-agent
```

---

## Health Checks

### Health Check Syntax

```yaml
services:
  postgres:
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 60s
      start_interval: 5s
```

### Health Check Parameters

| Parameter        | Description                    | Default  |
| ---------------- | ------------------------------ | -------- |
| `test`           | Command to run                 | Required |
| `interval`       | Time between checks            | 30s      |
| `timeout`        | Max time for check to complete | 30s      |
| `retries`        | Failures before unhealthy      | 3        |
| `start_period`   | Grace period on startup        | 0s       |
| `start_interval` | Interval during start_period   | interval |

### Test Command Formats

```yaml
# Shell command
test: ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"]

# Direct command (no shell)
test: ["CMD", "pg_isready", "-U", "postgres"]

# Disable health check
test: ["NONE"]
```

### Common Health Check Examples

#### PostgreSQL

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-postgres}"]
  interval: 10s
  timeout: 5s
  retries: 5
  start_period: 30s
```

#### Redis

```yaml
healthcheck:
  test: ["CMD", "redis-cli", "ping"]
  interval: 10s
  timeout: 5s
  retries: 5
```

#### HTTP API

```yaml
healthcheck:
  test: ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 60s
```

#### NATS

```yaml
healthcheck:
  test: ["CMD-SHELL", "wget --spider http://localhost:8222/healthz || exit 1"]
  interval: 10s
  timeout: 5s
  retries: 5
```

#### Python App (Custom)

```yaml
healthcheck:
  test:
    [
      "CMD",
      "python",
      "-c",
      "import requests; requests.get('http://localhost:8080/health').raise_for_status()",
    ]
  interval: 30s
  timeout: 10s
  retries: 3
```

### Viewing Health Status

```bash
# View health in container list
docker ps

# Detailed health info
docker inspect --format='{{json .State.Health}}' legacy-agent | jq

# Health check logs
docker inspect --format='{{json .State.Health.Log}}' legacy-agent | jq
```

---

## Resource Limits

### Memory Limits

```yaml
services:
  legacy-agent:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

### CPU Limits

```yaml
services:
  legacy-agent:
    deploy:
      resources:
        limits:
          cpus: "0.50" # 50% of one CPU
        reservations:
          cpus: "0.25"
```

### Combined Example

```yaml
services:
  legacy-agent:
    deploy:
      resources:
        limits:
          cpus: "2"
          memory: 2G
        reservations:
          cpus: "0.5"
          memory: 512M
```

### PIDs Limit

```yaml
services:
  legacy-agent:
    pids_limit: 100
```

---

## Logging

### Logging Drivers

| Driver      | Description                 |
| ----------- | --------------------------- |
| `json-file` | Default, logs to JSON files |
| `local`     | Optimized local logging     |
| `syslog`    | System logging              |
| `journald`  | systemd journal             |
| `none`      | No logging                  |

### Logging Configuration

```yaml
services:
  legacy-agent:
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"
        compress: "true"
```

### Default Logging (docker-compose.yml)

```yaml
x-logging: &default-logging
  driver: json-file
  options:
    max-size: "10m"
    max-file: "3"

services:
  legacy-agent:
    logging: *default-logging
```

---

## Build Configuration

### Build in Compose

```yaml
services:
  legacy-agent:
    build:
      context: ./agents/legacy
      dockerfile: Dockerfile
      args:
        VERSION: ${VERSION:-latest}
        BUILD_DATE: ${BUILD_DATE}
      target: production
      cache_from:
        - onemind/legacy-agent:latest
```

### Multi-Stage Build

```dockerfile
# Dockerfile
FROM python:3.12-slim AS builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --user -r requirements.txt

FROM python:3.12-slim AS production
WORKDIR /app
COPY --from=builder /root/.local /root/.local
COPY . .
CMD ["python", "main.py"]
```

### Build Commands

```bash
# Build all services
docker compose build

# Build specific service
docker compose build legacy-agent

# Build with no cache
docker compose build --no-cache

# Build and start
docker compose up -d --build
```

---

## Secrets Management

### Docker Secrets (Compose)

```yaml
services:
  legacy-agent:
    secrets:
      - db_password
      - api_key

secrets:
  db_password:
    file: ./secrets/db_password.txt
  api_key:
    environment: API_KEY
```

### Accessing Secrets in Container

Secrets are mounted at `/run/secrets/<secret_name>`:

```python
# Python example
with open('/run/secrets/db_password', 'r') as f:
    db_password = f.read().strip()
```

### Best Practices

- Never commit secrets to Git
- Use `.gitignore` for secret files
- Rotate secrets regularly
- Use environment-specific secrets

---

## Profiles

### Defining Profiles

```yaml
services:
  postgres:
    profiles: ["dev", "prod"]

  pgadmin:
    profiles: ["dev"] # Only in dev

  legacy-agent:
    profiles: ["dev", "prod"]
```

### Using Profiles

```bash
# Start dev profile
docker compose --profile dev up -d

# Start multiple profiles
docker compose --profile dev --profile debug up -d

# Set via environment
COMPOSE_PROFILES=dev docker compose up -d
```

---

## Common Commands

### Lifecycle Commands

```bash
# Start services
docker compose up -d

# Stop services
docker compose down

# Restart services
docker compose restart

# Stop without removing
docker compose stop

# Start stopped services
docker compose start
```

### Viewing Commands

```bash
# View running services
docker compose ps

# View logs
docker compose logs -f

# View specific service logs
docker compose logs -f legacy-agent

# View resource usage
docker compose top
```

### Maintenance Commands

```bash
# Pull latest images
docker compose pull

# Rebuild and restart
docker compose up -d --build

# Remove unused resources
docker system prune -a

# Remove unused volumes
docker volume prune
```

### Debugging Commands

```bash
# Execute command in service
docker compose exec legacy-agent bash

# Run one-off command
docker compose run --rm legacy-agent python migrate.py

# View service configuration
docker compose config
```

---

## Troubleshooting

### Common Issues

#### Container Won't Start

```bash
# Check logs
docker compose logs legacy-agent

# Check exit code
docker inspect legacy-agent --format='{{.State.ExitCode}}'

# Check health status
docker inspect legacy-agent --format='{{json .State.Health}}' | jq
```

#### Port Already in Use

```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change the port mapping
```

#### Volume Permission Issues

```bash
# Check container user
docker exec legacy-agent id

# Fix permissions on host
sudo chown -R 1000:1000 ./data
```

#### Container Can't Reach Other Container

```bash
# Check network connectivity
docker exec legacy-agent ping postgres

# Verify same network
docker network inspect onemind-network
```

#### Out of Disk Space

```bash
# Check Docker disk usage
docker system df

# Clean up
docker system prune -a --volumes
```

### Useful Debug Commands

```bash
# Interactive shell in running container
docker exec -it legacy-agent bash

# Network debugging
docker exec legacy-agent nslookup postgres
docker exec legacy-agent curl -v http://postgres:5432

# View all environment variables
docker exec legacy-agent env

# View mounted volumes
docker inspect legacy-agent --format='{{json .Mounts}}' | jq
```

---

## OneMind OS Compose Files

| File                       | Purpose                               | Location               |
| -------------------------- | ------------------------------------- | ---------------------- |
| `docker-compose.core.yml`  | Core services (Postgres, NATS, Redis) | Root                   |
| `docker-compose.yml`       | All services                          | Root                   |
| `cloud/docker-compose.yml` | Cloud production                      | `infra/compose/cloud/` |
| `home/docker-compose.yml`  | Home server                           | `infra/compose/home/`  |

### Quick Start

```bash
# Start core services
docker compose -f docker-compose.core.yml up -d

# Start all services
docker compose up -d

# View status
docker compose ps

# View logs
docker compose logs -f
```

---

## References

- [Docker Compose Specification](https://docs.docker.com/compose/compose-file/)
- [Docker CLI Reference](https://docs.docker.com/engine/reference/commandline/docker/)
- [Best Practices for Dockerfiles](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)

---

_Last updated: January 2026_
