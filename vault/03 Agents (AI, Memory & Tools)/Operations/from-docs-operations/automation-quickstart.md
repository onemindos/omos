# Automation Quick Start

## 1. Watchtower (Auto-Update Containers Weekly)

Add to `docker-compose.yml`:

```yaml
watchtower:
  image: containrrr/watchtower:latest
  container_name: onemind-watchtower
  restart: unless-stopped
  environment:
    WATCHTOWER_CLEANUP: "true"
    WATCHTOWER_SCHEDULE: "0 0 3 * * SUN"  # 3 AM Sundays
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  networks:
    - onemind
```

Start:
```bash
docker compose up -d watchtower
```

## 2. N8N Self-Healing Workflow

1. Access N8N: http://YOUR_IP:5678
2. Create workflow: "Service Health Monitor"
3. Add nodes:
   - Cron Trigger (every 5 min)
   - HTTP Request (check health)
   - IF (unhealthy?)
   - Execute Command (restart container)
   - Discord Notification
4. Activate workflow

## 3. Test Auto-Healing

```bash
# Stop a service
docker stop onemind-agno

# Wait 5 minutes
# N8N will detect and restart it
# Check Discord for notification
```

See [../REFERENCE.md](../REFERENCE.md) for complete automation setup.
