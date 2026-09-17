# OneMind Cloud Operations Runbook

*Day-to-Day Procedures and Emergency Response - January 2026*

## Quick Reference

### Connection Methods

```bash
# SSH via Tailscale IP
ssh ubuntu@100.87.87.69

# SSH via Tailscale MagicDNS
ssh ubuntu@onemind-cloud-1.tail717ba0.ts.net

# If SSH fails, use Portainer console
https://portainer.onemindos.dev
```

### Key Paths

| Resource | Path |
|----------|------|
| Docker Compose | `/home/ubuntu/onemind-os/docker_substrate_runtime/cloud/` |
| Compose File | `/home/ubuntu/onemind-os/docker_substrate_runtime/cloud/docker-compose.yml` |
| Environment | `/home/ubuntu/onemind-os/docker_substrate_runtime/cloud/.env` |
| Data Volumes | `/home/ubuntu/onemind-os/docker_substrate_runtime/cloud/data/` |
| Documentation | `/home/ubuntu/onemind-os/docs_memory_banks/` |

---

## Daily Operations

### Morning Health Check

```bash
# 1. Connect to server
ssh ubuntu@100.87.87.69

# 2. Check all containers
docker ps --format 'table {{.Names}}\t{{.Status}}'

# 3. Check for unhealthy containers
docker ps --filter "health=unhealthy"

# 4. Check disk space
df -h | grep -E "^/dev|Filesystem"

# 5. Check memory
free -h

# 6. Check Cloudflare tunnel
sudo systemctl status cloudflared --no-pager

# 7. Check Tailscale
tailscale status
```

### Quick Status One-Liner

```bash
ssh ubuntu@100.87.87.69 "docker ps --format '{{.Names}}: {{.Status}}' && echo '---' && df -h / | tail -1 && free -h | grep Mem"
```

---

## Common Procedures

### Restart a Service

```bash
# Single service restart
docker restart onemind-n8n

# Multiple services
docker restart onemind-agno onemind-letta

# All services (graceful)
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud
docker compose restart
```

### View Logs

```bash
# Live logs (follow)
docker logs onemind-agno -f

# Last 100 lines
docker logs onemind-agno --tail 100

# Logs with timestamps
docker logs onemind-agno --tail 100 -t

# All logs since time
docker logs onemind-agno --since 1h
```

### Stop/Start Services

```bash
# Stop single service
docker stop onemind-voice

# Start single service
docker start onemind-toc

# Stop all (preserves data)
docker compose stop

# Start all
docker compose start
```

### Full Stack Restart

```bash
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud

# Graceful restart (recommended)
docker compose restart

# Full down/up (if issues persist)
docker compose down
docker compose up -d
```

### Check Service Health

```bash
# Health status of specific container
docker inspect onemind-agno --format='{{.State.Health.Status}}'

# All health statuses
docker ps --format '{{.Names}}: {{.Status}}'

# Detailed health info
docker inspect onemind-agno | jq '.[0].State.Health'
```

---

## Update Procedures

### Update Single Container

```bash
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud

# Pull new image
docker compose pull <service-name>

# Recreate container
docker compose up -d <service-name>
```

### Update All Containers

```bash
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud

# Pull all new images
docker compose pull

# Recreate all containers
docker compose up -d

# Clean up old images
docker image prune -f
```

### Update Custom-Built Images

```bash
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud

# Rebuild specific service
docker compose build <service-name>

# Rebuild and restart
docker compose up -d --build <service-name>
```

---

## Backup Procedures

### Backup Database

```bash
# PostgreSQL dump
docker exec onemind-postgres pg_dump -U postgres onemind > /home/ubuntu/backups/postgres_$(date +%Y%m%d).sql

# Compress
gzip /home/ubuntu/backups/postgres_$(date +%Y%m%d).sql
```

### Backup All Volumes

```bash
# Create backup directory
mkdir -p /home/ubuntu/backups/$(date +%Y%m%d)

# Backup all data directories
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud
tar -czvf /home/ubuntu/backups/$(date +%Y%m%d)/data.tar.gz data/
```

### Backup Docker Compose

```bash
cp docker-compose.yml /home/ubuntu/backups/docker-compose.yml.$(date +%Y%m%d)
cp .env /home/ubuntu/backups/.env.$(date +%Y%m%d)
```

---

## Troubleshooting Procedures

### Container Won't Start

```bash
# 1. Check logs for error
docker logs onemind-<name> --tail 50

# 2. Check exit code
docker inspect onemind-<name> --format='{{.State.ExitCode}}'

# 3. Check if dependencies are running
docker ps | grep -E "postgres|redis|qdrant"

# 4. Try manual start with output
docker compose up <service-name>  # (no -d)
```

### Container Unhealthy

```bash
# 1. Check health log
docker inspect onemind-<name> | jq '.[0].State.Health.Log'

# 2. Test endpoint manually
docker exec onemind-<name> curl -s http://localhost:PORT/health

# 3. Check if service is actually responding
curl http://100.87.87.69:PORT/health
```

### High Memory Usage

```bash
# 1. Check per-container memory
docker stats --no-stream

# 2. Find the culprit
docker stats --no-stream --format "{{.Name}}: {{.MemUsage}}"

# 3. Restart memory-heavy container
docker restart onemind-<name>

# 4. If system-wide issue
sudo sync && echo 3 | sudo tee /proc/sys/vm/drop_caches
```

### High Disk Usage

```bash
# 1. Check disk usage
df -h

# 2. Check Docker disk usage
docker system df

# 3. Clean up unused images
docker image prune -f

# 4. Clean up unused volumes (CAREFUL!)
docker volume prune -f

# 5. Clean up build cache
docker builder prune -f

# 6. Find large files
sudo du -sh /home/ubuntu/* | sort -h
```

### Network Issues

```bash
# 1. Check Cloudflare tunnel
sudo systemctl status cloudflared
sudo journalctl -u cloudflared --tail 50

# 2. Check Tailscale
tailscale status

# 3. Check Docker networks
docker network ls
docker network inspect onemind-network

# 4. Test internal DNS
docker exec onemind-agno ping onemind-postgres
```

### Service Not Accessible via Cloudflare

```bash
# 1. Check tunnel status
sudo systemctl status cloudflared

# 2. Restart tunnel
sudo systemctl restart cloudflared

# 3. Check if service responds locally
curl http://localhost:PORT

# 4. Check DNS propagation
dig service.onemindos.dev
```

---

## Emergency Procedures

### Full System Freeze

```bash
# 1. SSH to server (if possible)
ssh ubuntu@100.87.87.69

# 2. If Docker is unresponsive
sudo systemctl restart docker

# 3. If that doesn't work
sudo reboot

# 4. After reboot
cd /home/ubuntu/onemind-os/docker_substrate_runtime/cloud
docker compose up -d
```

### Database Corruption

```bash
# 1. Stop all dependent services
docker stop onemind-agno onemind-letta onemind-n8n onemind-kestra

# 2. Restore from backup
docker exec -i onemind-postgres psql -U postgres < /home/ubuntu/backups/postgres_YYYYMMDD.sql

# 3. Restart services
docker start onemind-agno onemind-letta onemind-n8n onemind-kestra
```

### Disk Full Emergency

```bash
# 1. Quick cleanup
docker system prune -af

# 2. Remove old logs
sudo truncate -s 0 /var/log/*.log

# 3. Find and remove large files
sudo find /var -type f -size +100M -delete

# 4. Clear Docker logs
sudo truncate -s 0 $(docker inspect --format='{{.LogPath}}' onemind-agno)
```

### Security Incident

```bash
# 1. Isolate - stop all services
docker compose down

# 2. Block external access
sudo systemctl stop cloudflared

# 3. Check for unauthorized access
last -20
cat /var/log/auth.log | tail -50

# 4. Review container logs
for c in $(docker ps -aq); do docker logs $c --tail 100 > /tmp/logs_$(docker inspect --format='{{.Name}}' $c).txt; done
```

---

## Monitoring Endpoints

### Internal Health Checks

| Service | Health URL |
|---------|------------|
| Agno | http://localhost:8080/health |
| Letta | http://localhost:8283/v1/health |
| N8N | http://localhost:5678/healthz |
| Kestra | http://localhost:8081/health |
| Actions | http://localhost:8001/health |
| Uptime Kuma | http://localhost:3001/api/status |
| Qdrant | http://localhost:6333/ |
| NATS | http://localhost:8222/ |

### External Monitoring

- **Uptime Kuma Dashboard**: https://status.onemindos.dev
- **Portainer**: https://portainer.onemindos.dev
- **NATS Monitor**: https://nats.onemindos.dev

---

## Scheduled Maintenance

### Weekly Tasks

- [ ] Check disk usage and clean if > 80%
- [ ] Review unhealthy container alerts
- [ ] Check for pending updates
- [ ] Verify backups are running

### Monthly Tasks

- [ ] Update all container images
- [ ] Test backup restoration
- [ ] Review Cloudflare access logs
- [ ] Rotate any temporary credentials

### Quarterly Tasks

- [ ] Full security audit
- [ ] Performance optimization review
- [ ] Documentation update
- [ ] Disaster recovery drill

---

## Contact Information

### Escalation Path

1. **Self-Service**: Use this runbook
2. **Dashboard Access**: https://portainer.onemindos.dev
3. **AWS Console**: https://lightsail.aws.amazon.com
4. **Cloudflare Dashboard**: https://dash.cloudflare.com

### Credentials Location

All credentials stored in 1Password:
- **Vault**: OneMind - Unified Intelligence
- **Item**: "Cloudflare - OneMind OS"
- **Item**: "AWS Lightsail"


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
