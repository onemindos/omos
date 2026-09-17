---
title: "Docker"
type: reference
status: active
tags: [platform, infrastructure, containers]
---
# Docker

> Container runtime for all self-hosted services (Langfuse, NetBird, Frigate, etc.)

## CLI
```bash
docker ps                    # Running containers
docker compose up -d         # Start stack
docker compose logs -f       # Follow logs
docker system prune -af      # Clean up
docker stats                 # Resource usage
```

## Key Stacks
- Langfuse (:3004)
- NetBird (:33080)
- Frigate (NVR)
