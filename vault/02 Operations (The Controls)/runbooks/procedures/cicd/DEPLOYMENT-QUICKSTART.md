# Zero-Downtime Deployment Quick Start

## Usage

```bash
# Deploy updates
cd /opt/onemind
./production/4-deployment/deploy.sh
```

The script will:
1. Pull latest code from git
2. Pull new Docker images
3. Update services one by one
4. Wait for health checks
5. Rollback on failure
6. Send Discord notification

## Features

- ✅ Rolling updates (one service at a time)
- ✅ Health check verification
- ✅ Automatic rollback
- ✅ Discord notifications
- ✅ No downtime

## Manual Deployment (if script fails)

```bash
cd /opt/onemind
git pull
docker compose -f containers/cloud/docker-compose.yml pull
docker compose -f containers/cloud/docker-compose.yml up -d
docker compose ps
```

See [../REFERENCE.md](../REFERENCE.md) for full deployment strategies.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
