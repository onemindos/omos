---
id: ui-12-20260111053926
title: "Legacy AI - VM Infrastructure Upgrade Guide"
type: reference
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/12]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Legacy AI - VM Infrastructure Upgrade Guide"
---

# Legacy AI - VM Infrastructure Upgrade Guide

**Date**: November 11, 2025  
**VM**: 3.215.29.128 (Ubuntu 24.04 LTS)  
**Purpose**: Complete system upgrade and optimization

---

## 📋 Overview

This comprehensive upgrade plan covers:

- ✅ **System Updates**: Ubuntu OS, kernel, Docker Engine, Portainer
- ✅ **Container Updates**: All 20 application containers to latest stable versions
- ✅ **Health Check Fixes**: Node-RED and EMQX unhealthy status resolution
- ✅ **Caddy Optimization**: Enhanced SSL configuration, HTTP/3 support, improved routing
- ✅ **Security**: UFW firewall enablement with proper rules
- ✅ **Cleanup**: Unused Docker resources (reclaim ~3.3GB)

---

## 🚀 Quick Start

### Option 1: Automated Deployment (Recommended)

Run from your local machine:

```bash
cd /Users/zeuslegacy/Documents/Legacy-AI
chmod +x scripts/deploy-vm-upgrade.sh
./scripts/deploy-vm-upgrade.sh
```

This will:
1. Connect to the VM via SSH
2. Upload all upgrade scripts
3. Execute the complete upgrade process
4. Apply health check fixes
5. Update Caddy configuration
6. Run validation tests
7. Generate a comprehensive report

**Estimated Time**: 10-15 minutes

### Option 2: Manual Execution

SSH into the VM and run manually:

```bash
# 1. SSH to VM
ssh -i "/path/to/LightsailDefaultKey-us-east-1.pem" ubuntu@3.215.29.128

# 2. Download/upload the upgrade script
# (Copy scripts from repository)

# 3. Run upgrade
sudo bash vm-upgrade-complete.sh

# 4. Fix health checks
bash fix-health-checks.sh

# 5. Update Caddy
sudo cp Caddyfile.optimized /opt/legacy-ai/infrastructure/caddy/Caddyfile
docker exec caddy caddy reload --config /etc/caddy/Caddyfile
```

---

## 📦 What Gets Upgraded

### System Level

| Component | Current | Action |
|-----------|---------|--------|
| **Ubuntu** | 24.04 LTS | Update all packages |
| **Kernel** | 6.14.0 | Update to latest stable |
| **Docker Engine** | 20.10+ | Update to latest |
| **Docker Compose** | v2.40.3 | Update to latest |
| **Portainer** | Latest EE | Pull latest image |

### Container Images

| Service | Current | Action |
|---------|---------|--------|
| **n8n** | latest | Pull & restart (5 containers) |
| **Letta** | Custom build | Rebuild if needed |
| **EMQX** | latest | Pull & restart |
| **Node-RED** | latest | Pull & restart |
| **Home Assistant** | stable | Pull & restart |
| **Caddy** | latest | Pull & restart |
| **PostgreSQL** | 16-alpine | Update if available |
| **Redis** | 7-alpine | Update if available |

### Configurations

- **Caddyfile**: Enhanced with HTTP/3, better headers, logging
- **Health Checks**: Fixed for Node-RED and EMQX
- **Firewall**: UFW enabled with ports 22, 80, 443

---

## 🔍 Pre-Upgrade Checklist

Before running the upgrade:

- [ ] Review current system state
- [ ] Verify backup systems are operational
- [ ] Confirm no critical operations running
- [ ] Schedule maintenance window (if needed)
- [ ] Notify users of potential brief downtime
- [ ] Have rollback plan ready

### Check Current State

```bash
# On VM
ssh -i "key.pem" ubuntu@3.215.29.128

# System info
uname -a
lsb_release -a
docker --version
docker-compose --version

# Container status
docker ps --format 'table {{.Names}}\t{{.Status}}'

# Resource usage
free -h
df -h
docker system df
```

---

## 📝 Upgrade Phases

### Phase 1: Pre-Upgrade Audit & Backup

**What Happens**:
- Creates backup directory: `/home/ubuntu/backups/pre-upgrade-YYYYMMDD_HHMMSS/`
- Documents current system state (OS, Docker, containers, images)
- Backs up all docker-compose.yml files
- Backs up Caddyfile
- Exports all PostgreSQL databases (pg_dumpall)
- Exports all Redis data (dump.rdb)
- Creates backup manifest

**Duration**: ~5 minutes

**Output**: Full system snapshot for rollback

### Phase 2: System-Level Upgrades

**What Happens**:
- Updates Ubuntu package lists
- Upgrades all system packages
- Updates Docker Engine and CLI
- Updates Docker Compose plugin
- Pulls latest Portainer image
- Configures UFW firewall rules
- Enables firewall

**Duration**: ~3-5 minutes

**Potential Impact**: Brief interruption during Docker restart

### Phase 3: Container Image Updates

**What Happens**:
- Discovers all docker-compose.yml locations
- For each stack:
  - Pulls latest images
  - Performs rolling update (no downtime)
  - Recreates containers with new images
  - Waits for stability

**Duration**: ~5-10 minutes (depends on image sizes)

**Order**:
1. Infrastructure (Caddy, Portainer)
2. Databases (PostgreSQL, Redis)
3. Applications (n8n, Letta)
4. IoT (Node-RED, EMQX, Home Assistant)

### Phase 4: Fix Unhealthy Containers

**What Happens**:
- Identifies containers with failing health checks
- Updates docker-compose.yml files with proper health checks
- Node-RED: Adds curl-based HTTP check
- EMQX: Adds native EMQX control check
- Restarts affected containers
- Monitors health status

**Duration**: ~2-3 minutes

**Expected Result**: All containers healthy

### Phase 5: Caddy SSL & Routing Optimization

**What Happens**:
- Backs up current Caddyfile
- Deploys optimized Caddyfile with:
  - HTTP/3 (QUIC) support
  - Enhanced headers (X-Real-IP, X-Forwarded-For)
  - WebSocket support for all applicable services
  - Increased timeouts for AI operations
  - Per-domain logging
  - GZIP compression
- Reloads Caddy gracefully
- Tests all domains
- Verifies SSL certificates

**Duration**: ~1 minute

**Zero Downtime**: Caddy hot-reloads configuration

### Phase 6: Post-Upgrade Validation & Cleanup

**What Happens**:
- Verifies all containers running
- Checks container health status
- Tests system resources
- Cleans up:
  - Unused Docker images
  - Dangling volumes
  - Orphaned networks
- Generates upgrade report

**Duration**: ~2-3 minutes

**Result**: Clean, optimized system

---

## 🛡️ Safety Features

### Automatic Backups

Every phase creates backups before making changes:

- **System state**: Full documentation
- **Configurations**: All docker-compose.yml files
- **Databases**: Complete exports
- **Caddy**: Configuration backup

### Rollback Capability

If anything goes wrong:

```bash
# On VM
sudo bash rollback-upgrade.sh ~/backups/pre-upgrade-YYYYMMDD_HHMMSS
```

This will:
1. Stop all containers
2. Restore configurations
3. Restore databases
4. Restart all services

### Validation Checks

After each phase:
- Container status verification
- Health check validation
- Resource usage monitoring
- Service availability testing

---

## 📊 Expected Results

### System State (After Upgrade)

```
✅ OS: Ubuntu 24.04 LTS (latest kernel)
✅ Docker: Latest stable version
✅ Containers: 20/20 running
✅ Health: All containers healthy
✅ SSL: All 8 domains with valid certificates
✅ Firewall: Enabled with proper rules
✅ Disk: ~3.3GB reclaimed
```

### Performance Improvements

- **Faster image pulls**: Latest Docker optimizations
- **Better SSL**: HTTP/3 support for improved latency
- **Proper health checks**: Earlier issue detection
- **Optimized routing**: Better headers and timeouts

---

## 🔧 Troubleshooting

### Upgrade Fails During Execution

**Solution**: Check the logs

```bash
# On VM
tail -100 /var/log/vm-upgrade-*.log
```

Common issues:
- **Disk full**: Clean up before retrying
- **Network timeout**: Retry pull operations
- **Permission denied**: Ensure running with sudo

### Container Won't Start After Upgrade

**Solution**: Check individual container logs

```bash
docker logs <container-name>
docker inspect <container-name>
```

Try restarting:
```bash
docker restart <container-name>
```

### Health Checks Still Failing

**Solution**: Manual health check fix

```bash
# SSH to VM
ssh -i "key.pem" ubuntu@3.215.29.128

# Run health check fixer
bash ~/upgrade-scripts/fix-health-checks.sh

# Check status after 2 minutes
docker ps --filter 'health=unhealthy'
```

### Domains Not Responding

**Solution**: Verify Caddy and SSL

```bash
# Check Caddy status
docker logs caddy

# Test Caddy config
docker exec caddy caddy validate --config /etc/caddy/Caddyfile

# Reload Caddy
docker exec caddy caddy reload --config /etc/caddy/Caddyfile

# Check SSL certificates
docker exec caddy caddy list-modules | grep tls
```

### Rollback Needed

**Solution**: Execute rollback script

```bash
# List available backups
ls -la ~/backups/

# Execute rollback
sudo bash rollback-upgrade.sh ~/backups/pre-upgrade-YYYYMMDD_HHMMSS
```

---

## 📈 Post-Upgrade Tasks

### Immediate (Within 1 Hour)

- [ ] Verify all domains accessible
- [ ] Test critical workflows in n8n
- [ ] Test Dify AI functionality
- [ ] Check Letta agent responses
- [ ] Monitor container logs for errors
- [ ] Verify Home Assistant automations

### Short-term (Within 24 Hours)

- [ ] Monitor system resources
- [ ] Check automated backups running
- [ ] Verify SSL certificate renewal
- [ ] Test all integrations
- [ ] Update documentation

### Long-term (Within 1 Week)

- [ ] Performance baseline comparison
- [ ] Review security logs
- [ ] Optimize resource allocation if needed
- [ ] Plan next maintenance window

---

## 📁 File Locations

### On VM

```
/home/ubuntu/
├── backups/
│   └── pre-upgrade-YYYYMMDD_HHMMSS/
│       ├── system-info/
│       ├── docker-configs/
│       ├── databases/
│       ├── caddy/
│       ├── BACKUP_INFO.txt
│       └── UPGRADE_REPORT.txt
├── upgrade-scripts/
│   ├── vm-upgrade-complete.sh
│   ├── fix-health-checks.sh
│   └── Caddyfile.optimized
└── /var/log/
    └── vm-upgrade-YYYYMMDD_HHMMSS.log

/opt/legacy-ai/
├── infrastructure/
│   └── caddy/
│       ├── Caddyfile
│       ├── Caddyfile.backup-YYYYMMDD
│       └── docker-compose.yml
└── stacks/
    ├── ai/
    ├── apps/
    └── iot/
```

### In Repository

```
Legacy-AI/
├── scripts/
│   ├── vm-upgrade-complete.sh
│   ├── deploy-vm-upgrade.sh
│   ├── fix-health-checks.sh
│   └── rollback-upgrade.sh
├── ENVIRONMENTS/production/infrastructure/caddy/
│   └── Caddyfile.optimized
└── KNOWLEDGE/docs/operations/
    └── VM-UPGRADE-GUIDE.md (this file)
```

---

## 🔐 Security Considerations

### Firewall Rules

After upgrade, UFW will be enabled with:

```
Port 22/tcp  - SSH (Required)
Port 80/tcp  - HTTP (Caddy, redirects to HTTPS)
Port 443/tcp - HTTPS (Caddy reverse proxy)
Port 443/udp - HTTP/3 QUIC
```

All other ports blocked by default.

### SSL/TLS

- **Automatic HTTPS**: Let's Encrypt certificates
- **Auto-renewal**: Caddy handles automatically
- **HTTP/3**: Enabled for performance
- **Security headers**: Added to all proxied services

### Best Practices

- ✅ Keep SSH key secure (600 permissions)
- ✅ Use key-based authentication only
- ✅ Regular backups (automated daily)
- ✅ Monitor security logs
- ✅ Keep systems updated

---

## 📞 Support & References

### Related Documentation

- **Current Services**: `KNOWLEDGE/docs/current-state/VM-STACKS-INVENTORY.md`
- **Maintenance**: `KNOWLEDGE/docs/guides/MAINTENANCE-PLAYBOOK.md`
- **Security**: `KNOWLEDGE/docs/current-state/VM-SECURITY-HARDENING-COMPLETE.md`
- **Backups**: `KNOWLEDGE/docs/guides/BACKUP-AND-RECOVERY.md`

### Quick Links

- **Portainer**: https://command.zeusops.dev
- **n8n**: https://n8n.zeusops.dev
- **Dify**: https://dify.zeusops.dev
- **Letta**: https://letta.zeusops.dev

### Emergency Contacts

If critical issues occur:
1. Check logs: `/var/log/vm-upgrade-*.log`
2. Review backup: `~/backups/pre-upgrade-*/`
3. Execute rollback if needed
4. Document issue for future reference

---

## ✅ Success Criteria

The upgrade is considered successful when:

- [ ] All 20 containers running
- [ ] No unhealthy containers
- [ ] All 8 domains responding (HTTP 200/307)
- [ ] SSL certificates valid
- [ ] Firewall enabled and configured
- [ ] System resources optimal (<50% CPU, <60% memory)
- [ ] No errors in container logs
- [ ] All applications functional

---

## 📝 Version History

| Date | Version | Changes |
|------|---------|---------|
| 2025-11-11 | 1.0 | Initial comprehensive upgrade plan |

---

**Last Updated**: November 11, 2025  
**Maintained By**: Legacy AI DevOps Team
