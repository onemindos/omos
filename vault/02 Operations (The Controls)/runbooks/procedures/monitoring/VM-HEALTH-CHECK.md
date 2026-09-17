# OneMind OS - VM Health Check & Upgrade Guide

> Last Updated: 2026-01-30
> Systematic checklist for monitoring, cleanup, and upgrade decisions

---

## Quick Health Check (Run Weekly)

```bash
# SSH to VM
ssh zeus@100.102.21.44

# Run all checks at once
echo "=== DISK ===" && df -h / && \
echo "" && echo "=== MEMORY ===" && free -h && \
echo "" && echo "=== CPU LOAD ===" && cat /proc/loadavg && \
echo "" && echo "=== CONTAINERS ===" && docker ps --format 'table {{.Names}}\t{{.Status}}' && \
echo "" && echo "=== DOCKER DISK ===" && docker system df
```

---

## 1. Disk Usage Check

### Commands

```bash
# Overall disk usage
df -h /

# What's using space
sudo du -h --max-depth=1 / 2>/dev/null | sort -hr | head -10

# Docker specific
docker system df
docker system df -v  # Detailed view
```

### Thresholds

| Usage | Status | Action |
|-------|--------|--------|
| < 50% | ✅ Healthy | None needed |
| 50-70% | ⚠️ Monitor | Plan cleanup |
| 70-85% | 🟠 Warning | Run cleanup now |
| > 85% | 🔴 Critical | Immediate cleanup or upgrade |

### Current Specs

| Property | Value |
|----------|-------|
| Disk Size | 64 GB (P6) |
| Tier | Premium SSD |
| Cost | ~$9/month |

### Upgrade Path

| Tier | Size | Cost/mo | When to upgrade |
|------|------|---------|-----------------|
| P6 (current) | 64 GB | $9 | - |
| P10 | 128 GB | $19 | When >50GB used |
| P15 | 256 GB | $38 | Large data/logs |
| P20 | 512 GB | $75 | Heavy workloads |

### Upgrade Command

```bash
# Requires VM stop
az vm deallocate -g ONEMIND-PRODUCTION -n onemind-vm
az disk update -g ONEMIND-PRODUCTION \
  -n onemind-vm_OsDisk_1_9673580f9734496a96322a908d72db22 \
  --size-gb 128
az vm start -g ONEMIND-PRODUCTION -n onemind-vm
# Filesystem expands automatically on boot
```

---

## 2. Memory Check

### Commands

```bash
# Memory overview
free -h

# Per-container memory
docker stats --no-stream --format 'table {{.Name}}\t{{.MemUsage}}\t{{.MemPerc}}'

# Top memory processes
ps aux --sort=-%mem | head -10
```

### Understanding Output

```
              total        used        free      shared  buff/cache   available
Mem:          7.7Gi       2.1Gi       489Mi       73Mi       5.1Gi       5.2Gi
```

- **total**: Total RAM (7.7GB)
- **used**: Actively used (2.1GB)
- **available**: Actually available for new processes (5.2GB) ← This is what matters
- **buff/cache**: Cached data (can be freed if needed)

### Thresholds

| Available | Status | Action |
|-----------|--------|--------|
| > 2 GB | ✅ Healthy | None needed |
| 1-2 GB | ⚠️ Monitor | Check for memory leaks |
| 500MB-1GB | 🟠 Warning | Optimize or upgrade |
| < 500MB | 🔴 Critical | Upgrade VM size |

### Current Specs

| Property | Value |
|----------|-------|
| Total RAM | 8 GB |
| VM Size | Standard_D2s_v3 |
| Cost | ~$96/month |

### Upgrade Path

| VM Size | vCPUs | RAM | Cost/mo | When to upgrade |
|---------|-------|-----|---------|-----------------|
| D2s_v3 (current) | 2 | 8 GB | $96 | - |
| D4s_v3 | 4 | 16 GB | $192 | When <1GB available |
| D8s_v3 | 8 | 32 GB | $384 | Heavy multi-agent |

### Upgrade Command

```bash
# Requires VM stop
az vm deallocate -g ONEMIND-PRODUCTION -n onemind-vm
az vm resize -g ONEMIND-PRODUCTION -n onemind-vm --size Standard_D4s_v3
az vm start -g ONEMIND-PRODUCTION -n onemind-vm
```

---

## 3. CPU Load Check

### Commands

```bash
# Load average (1min, 5min, 15min)
cat /proc/loadavg
uptime

# Real-time CPU per container
docker stats --no-stream --format 'table {{.Name}}\t{{.CPUPerc}}'

# Top CPU processes
top -bn1 | head -15
```

### Understanding Load Average

**Load = processes waiting for CPU**

For a **2-core VM**:

| Load | % Capacity | Status | Meaning |
|------|------------|--------|---------|
| 0.5 | 25% | ✅ Idle | Light usage |
| 1.0 | 50% | ✅ Normal | Moderate usage |
| 1.5 | 75% | ⚠️ Busy | Getting loaded |
| 2.0 | 100% | 🟠 Full | At capacity |
| 3.0+ | 150%+ | 🔴 Overloaded | Processes queuing |

**Rule of thumb**: Load should stay below number of CPU cores (2.0 for this VM)

### Thresholds (2-core VM)

| Load (15min avg) | Status | Action |
|------------------|--------|--------|
| < 1.0 | ✅ Healthy | None needed |
| 1.0 - 1.5 | ⚠️ Monitor | Watch trends |
| 1.5 - 2.0 | 🟠 Warning | Plan upgrade |
| > 2.0 sustained | 🔴 Critical | Upgrade now |

---

## 4. Docker Cleanup

### Check What Can Be Cleaned

```bash
# Overview
docker system df

# Detailed breakdown
docker system df -v
```

### Cleanup Commands

```bash
# Remove unused images (safe)
docker image prune -f

# Remove ALL unused images (more aggressive)
docker image prune -a -f

# Remove unused volumes (careful - check first!)
docker volume ls -f dangling=true
docker volume prune -f

# Remove unused networks
docker network prune -f

# Full cleanup (images, containers, volumes, networks)
docker system prune -a -f --volumes

# Remove specific old images
docker rmi <image-name>
```

### What's Safe to Delete

| Item | Safe? | Notes |
|------|-------|-------|
| Dangling images | ✅ Yes | Orphaned layers |
| Unused images | ⚠️ Check | May need for rollback |
| Build cache | ✅ Yes | Rebuilds as needed |
| Stopped containers | ⚠️ Check | May have logs you need |
| Unused volumes | ❌ Careful | May contain data |

### Cleanup Checklist

- [ ] Check `docker system df` for reclaimable space
- [ ] Remove dangling images: `docker image prune -f`
- [ ] Remove unused images: `docker image prune -a -f`
- [ ] Check for old backup files: `ls -lh ~/backups/`
- [ ] Check log sizes: `sudo du -sh /var/log/*`
- [ ] Clear old logs: `sudo journalctl --vacuum-time=7d`

---

## 5. Log Cleanup

### Check Log Sizes

```bash
# System logs
sudo du -sh /var/log/*

# Journal size
journalctl --disk-usage

# Docker container logs
for c in $(docker ps -q); do
  echo "$(docker inspect $c --format '{{.Name}}'): $(docker inspect $c --format '{{.LogPath}}' | xargs ls -lh 2>/dev/null | awk '{print $5}')"
done
```

### Cleanup Commands

```bash
# Clear old journal logs (keep 7 days)
sudo journalctl --vacuum-time=7d

# Clear old journal logs (keep 500MB)
sudo journalctl --vacuum-size=500M

# Truncate Docker container logs
sudo truncate -s 0 /var/lib/docker/containers/*/*-json.log

# Clear apt cache
sudo apt clean
```

---

## 6. Backup Status Check

### Commands

```bash
# Check Azure Backup status
az backup job list -g ONEMIND-PRODUCTION \
  --vault-name onemind-backup-vault \
  --query "[?properties.operation=='Backup'].{status:properties.status, start:properties.startTime}" \
  -o table

# Check local database backups
ls -lh ~/backups/

# Check backup cron
crontab -l
```

### Expected Status

| Backup | Frequency | Location | Check |
|--------|-----------|----------|-------|
| Azure VM Backup | Daily 3pm UTC | Recovery Vault | Azure Portal |
| PostgreSQL | Daily 3am UTC | ~/backups/ | `ls -lh ~/backups/postgres_*` |
| TimescaleDB | Daily 3am UTC | ~/backups/ | `ls -lh ~/backups/timescale_*` |

---

## 7. Container Health Check

### Commands

```bash
# All container status
docker ps -a --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'

# Check for unhealthy containers
docker ps --filter health=unhealthy

# Container resource usage
docker stats --no-stream

# View container logs
docker logs --tail 50 onemind-agno
```

### Expected Healthy State

| Container | Expected Status |
|-----------|-----------------|
| onemind-agno | Up (healthy) |
| onemind-flutter-web | Up (healthy) |
| onemind-postgres | Up (healthy) |
| onemind-redis | Up (healthy) |
| onemind-nats | Up (healthy) |
| onemind-timescaledb | Up (healthy) |
| onemind-watchtower | Up (healthy) |
| onemind-portainer | Up |

### Restart Unhealthy Container

```bash
docker restart <container-name>
```

---

## 8. Network/Connectivity Check

### Commands

```bash
# Check Cloudflare tunnel
sudo systemctl status cloudflared

# Check Tailscale
tailscale status

# Check Tailscale key expiry
tailscale status --json | jq '.Self.KeyExpiry'

# Test API endpoint
curl -s https://api.onemindos.dev/health | jq .
```

### Expected Status

| Service | Status |
|---------|--------|
| cloudflared | active (running) |
| tailscaled | active (running) |
| API health | {"status": "healthy"} |

---

## 9. Security Check

### Commands

```bash
# Check for available updates
sudo apt update && apt list --upgradable

# Check failed SSH attempts
sudo grep "Failed password" /var/log/auth.log | tail -10

# Check listening ports
sudo ss -tlnp

# Check firewall status
sudo ufw status
```

### Security Checklist

- [ ] SSH only via Tailscale (port 22 not public)
- [ ] No sensitive ports exposed publicly
- [ ] Regular security updates applied
- [ ] Strong passwords on all services

---

## Weekly Maintenance Checklist

### Quick Check (5 minutes)

- [ ] Disk usage < 70%: `df -h /`
- [ ] Memory available > 2GB: `free -h`
- [ ] Load average < 2.0: `uptime`
- [ ] All containers healthy: `docker ps`
- [ ] No backup failures: Check Azure Portal

### Monthly Cleanup (15 minutes)

- [ ] Run Docker cleanup: `docker system prune -f`
- [ ] Clear old logs: `sudo journalctl --vacuum-time=7d`
- [ ] Remove old backups (>30 days): Auto via script
- [ ] Check for system updates: `sudo apt update`
- [ ] Review Portainer for container trends: https://portainer.onemindos.dev

---

## One-Liner Health Check

Copy and run this for a quick status:

```bash
ssh zeus@100.102.21.44 'echo "=== DISK ===" && df -h / | tail -1 && echo "=== MEM ===" && free -h | grep Mem && echo "=== LOAD ===" && uptime && echo "=== DOCKER ===" && docker ps --format "{{.Names}}: {{.Status}}" | head -5 && echo "=== BACKUPS ===" && ls -1t ~/backups/*.sql.gz 2>/dev/null | head -2'
```

---

## Upgrade Decision Matrix

| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| Disk > 80% full | Images/logs | Cleanup first, then upgrade |
| Memory < 1GB available | Container bloat | Restart containers, then upgrade |
| Load > 2.0 sustained | CPU bound | Upgrade to D4s_v3 |
| Slow API responses | Multiple issues | Check all metrics |
| Container OOM killed | Memory | Increase container limits or upgrade |

---

## Quick Reference: Current VM Specs

| Resource | Current | Limit |
|----------|---------|-------|
| **VM Size** | Standard_D2s_v3 | - |
| **vCPUs** | 2 | - |
| **RAM** | 8 GB | - |
| **Disk** | 64 GB (P6) | - |
| **Cost** | ~$105/month | - |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
