# Docker

> **Containerization Platform**

## Overview

Docker is the containerization platform for packaging, deploying, and running OneMind OS services in isolated, reproducible environments.

## Status

- **Active**: ✅ Yes
- **Compose**: v2
- **Last Updated**: January 25, 2026

## Use Cases

1. **Local Development** - Consistent dev environment
2. **Deployment** - Production containers
3. **Dependencies** - Isolated services (NATS, Redis, Postgres)
4. **CI/CD** - Build and test pipelines
5. **Microservices** - Service isolation

## Key Concepts

| Concept        | Description                   |
| -------------- | ----------------------------- |
| **Image**      | Blueprint for containers      |
| **Container**  | Running instance of image     |
| **Dockerfile** | Instructions to build image   |
| **Compose**    | Multi-container orchestration |
| **Volume**     | Persistent data storage       |
| **Network**    | Container communication       |

## Installation

### macOS

```bash
# Install Docker Desktop
brew install --cask docker

# Or use Colima (lightweight)
brew install colima docker
colima start
```

### Verify

```bash
docker --version
docker compose version
```

## Dockerfile Examples

### Python Agent

```dockerfile
# agents/legacy/Dockerfile
FROM python:3.12-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY . .

# Run
EXPOSE 8080
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
```

### Multi-Stage Build

```dockerfile
# Build stage
FROM python:3.12-slim AS builder

WORKDIR /app
COPY requirements.txt .
RUN pip install --user --no-cache-dir -r requirements.txt

# Production stage
FROM python:3.12-slim

WORKDIR /app

# Copy dependencies from builder
COPY --from=builder /root/.local /root/.local
ENV PATH=/root/.local/bin:$PATH

# Copy application
COPY . .

EXPOSE 8080
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
```

### Next.js UI

```dockerfile
# ui/Dockerfile
FROM node:20-slim AS builder

WORKDIR /app

# Install dependencies
COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# Build
COPY . .
RUN pnpm build

# Production
FROM node:20-slim

WORKDIR /app

COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

EXPOSE 3000
CMD ["node", "server.js"]
```

## Docker Compose

### Core Services

```yaml
# docker-compose.core.yml
version: "3.8"

services:
  postgres:
    image: pgvector/pgvector:pg16
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: onemind
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: onemind
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "onemind"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes

  nats:
    image: nats:latest
    ports:
      - "4222:4222"
      - "8222:8222"
    command: "-js -m 8222"
    volumes:
      - nats-data:/data

volumes:
  postgres-data:
  redis-data:
  nats-data:
```

### Full Stack

```yaml
# docker-compose.yml
version: "3.8"

services:
  legacy-agent:
    build: ./agents/legacy
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://onemind:${POSTGRES_PASSWORD}@postgres:5432/onemind
      - REDIS_URL=redis://redis:6379
      - NATS_URL=nats://nats:4222
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
      nats:
        condition: service_started

  ui:
    build: ./ui
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_AGENT_API_URL=http://legacy-agent:8080
    depends_on:
      - legacy-agent

include:
  - docker-compose.core.yml
```

## Common Commands

### Container Management

```bash
# Start services
docker compose up -d

# Stop services
docker compose down

# View logs
docker compose logs -f legacy-agent

# Rebuild and start
docker compose up -d --build

# Restart single service
docker compose restart legacy-agent
```

### Image Management

```bash
# Build image
docker build -t legacy-agent:latest ./agents/legacy

# List images
docker images

# Remove image
docker rmi legacy-agent:latest

# Prune unused images
docker image prune -a
```

### Container Inspection

```bash
# List running containers
docker ps

# List all containers
docker ps -a

# Exec into container
docker exec -it legacy-agent sh

# View container logs
docker logs -f legacy-agent

# Inspect container
docker inspect legacy-agent
```

### Volumes

```bash
# List volumes
docker volume ls

# Create volume
docker volume create postgres-data

# Remove volume
docker volume rm postgres-data

# Prune unused volumes
docker volume prune
```

### Networks

```bash
# List networks
docker network ls

# Create network
docker network create onemind-network

# Connect container to network
docker network connect onemind-network legacy-agent
```

## Development Workflow

### Hot Reload Setup

```yaml
# docker-compose.dev.yml
services:
  legacy-agent:
    build:
      context: ./agents/legacy
      dockerfile: Dockerfile.dev
    volumes:
      - ./agents/legacy:/app
      - /app/__pycache__
    command: uvicorn main:app --host 0.0.0.0 --port 8080 --reload
```

### Dev Dockerfile

```dockerfile
# Dockerfile.dev
FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Don't copy code - mounted via volume
EXPOSE 8080

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080", "--reload"]
```

## VS Code Integration

### Dev Containers

```json
// .devcontainer/devcontainer.json
{
  "name": "OneMind Dev",
  "dockerComposeFile": ["../docker-compose.yml", "../docker-compose.dev.yml"],
  "service": "legacy-agent",
  "workspaceFolder": "/app",
  "extensions": ["ms-python.python", "ms-python.vscode-pylance"],
  "forwardPorts": [8080, 3000, 5432, 6379]
}
```

## Registry

### Docker Hub

```bash
# Login
docker login

# Tag image
docker tag legacy-agent:latest username/legacy-agent:v1.0.0

# Push
docker push username/legacy-agent:v1.0.0

# Pull
docker pull username/legacy-agent:v1.0.0
```

### AWS ECR

```bash
# Login to ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin 123456789.dkr.ecr.us-east-1.amazonaws.com

# Push to ECR
docker tag legacy-agent:latest 123456789.dkr.ecr.us-east-1.amazonaws.com/legacy-agent:latest
docker push 123456789.dkr.ecr.us-east-1.amazonaws.com/legacy-agent:latest
```

## Resource Management

### Limits

```yaml
services:
  legacy-agent:
    deploy:
      resources:
        limits:
          cpus: "1.0"
          memory: 1G
        reservations:
          cpus: "0.5"
          memory: 512M
```

### Cleanup

```bash
# Remove all stopped containers
docker container prune

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune

# Remove everything unused
docker system prune -a --volumes
```

## Security

### Non-Root User

```dockerfile
FROM python:3.12-slim

# Create non-root user
RUN useradd -m -s /bin/bash appuser
USER appuser

WORKDIR /home/appuser/app
COPY --chown=appuser:appuser . .

CMD ["python", "main.py"]
```

### Secrets

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
    environment: "API_KEY"
```

## Best Practices

- ✅ Use specific image tags (not `latest` in production)
- ✅ Multi-stage builds for smaller images
- ✅ Use `.dockerignore` to exclude files
- ✅ Run as non-root user
- ✅ Health checks for reliability
- ✅ Use named volumes for persistence
- ✅ Pin package versions
- ✅ Minimize layers (combine RUN commands)

## Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)

---

**Related**: [GitHub](./GitHub.md) | [AWS](./AWS.md) | [Cloudflare](./Cloudflare.md)
