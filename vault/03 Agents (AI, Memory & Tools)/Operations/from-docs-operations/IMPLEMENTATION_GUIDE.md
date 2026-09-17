# OneMind OS - Production Implementation Guide

## Overview

This guide provides a **week-by-week implementation plan** to transform your OneMind OS infrastructure from development to production-grade with automated maintenance, self-healing, and cloud portability.

**Current State:** ✅ 13/14 services running, health checks configured
**Goal:** 🎯 Production-ready, self-healing, automated infrastructure

**Total Time:** 4 weeks (1-2 hours per week)
**Difficulty:** Medium (mostly scripted automation)

---

## Quick Start (TL;DR)

```bash
# SSH to your server
ssh onemind-cloud

# Navigate to project
cd /opt/onemind

# Pull latest code (includes scripts)
git pull

# Run automated setup
sudo ./scripts/setup-production.sh

# Follow prompts - it will configure everything
```

**That's it!** The setup script handles:
- Security updates
- Firewall
- Backups
- Monitoring
- Resource management

---

## Week 1: Critical Stability

**Goal:** Make system stable, secure, and self-maintaining
**Time:** 1-2 hours
**Status:** 🟢 Can be done in production

### Task 1.1: System Hardening (30 minutes)

**1. Enable automatic security updates:**

```bash
sudo apt update
sudo apt install -y unattended-upgrades apt-listchanges

# Configure (interactive)
sudo dpkg-reconfigure -plow unattended-upgrades
```

**2. Configure UFW firewall:**

```bash
# Install
sudo apt install -y ufw

# Configure (IMPORTANT: Allow SSH first!)
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp comment 'SSH'

# OneMind ports
sudo ufw allow 8080/tcp comment 'Agno API'
sudo ufw allow 8000/tcp comment 'API Gateway'
sudo ufw allow 4222/tcp comment 'NATS Client'
sudo ufw allow 5678/tcp comment 'N8N'
sudo ufw allow 3001/tcp comment 'Uptime Kuma'
sudo ufw allow 9000/tcp comment 'Portainer'
sudo ufw allow 80/tcp comment 'HTTP'
sudo ufw allow 443/tcp comment 'HTTPS'

# Enable
sudo ufw enable
sudo ufw status verbose
```

**3. Install Fail2Ban:**

```bash
sudo apt install -y fail2ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Check status
sudo fail2ban-client status sshd
```

**✅ Verification:**
```bash
# Check firewall
sudo ufw status

# Check fail2ban
sudo systemctl status fail2ban

# Check auto-updates
sudo systemctl status unattended-upgrades
```

### Task 1.2: Docker Log Rotation (10 minutes)

**1. Configure Docker daemon:**

```bash
sudo tee /etc/docker/daemon.json > /dev/null <<EOF
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF
```

**2. Restart Docker:**

```bash
sudo systemctl restart docker

# Verify services come back up
cd /opt/onemind
docker compose -f containers/cloud/docker-compose.yml ps
```

**✅ Verification:**
```bash
# Check Docker config
docker info | grep "Log"

# Logs should now rotate at 10 MB
```

### Task 1.3: Add Resource Limits (20 minutes)

**1. Edit docker-compose.yml:**

Add resource limits to critical services. See [PRODUCTION_BEST_PRACTICES.md](PRODUCTION_BEST_PRACTICES.md#6-resource-management) for full configuration.

**Quick example for Agno:**

```yaml
agno:
  # ... existing config ...
  deploy:
    resources:
      limits:
        cpus: '1.5'
        memory: 2G
      reservations:
        cpus: '0.75'
        memory: 1G
```

**2. Apply changes:**

```bash
cd /opt/onemind
docker compose -f containers/cloud/docker-compose.yml up -d
```

**✅ Verification:**
```bash
# Check resource usage
docker stats

# All containers should be within limits
```

---

## Week 2: Backup & Disaster Recovery

**Goal:** Never lose data, fast recovery
**Time:** 1 hour
**Status:** 🟢 Safe to implement

### Task 2.1: Setup Automated Backups (30 minutes)

**1. Create backup directory:**

```bash
sudo mkdir -p /backup/onemind
sudo chown $USER:$USER /backup/onemind
```

**2. Test backup script:**

```bash
cd /opt/onemind
./scripts/backup.sh
```

**Expected output:**
```
✅ Backup complete: /backup/onemind/backup_YYYYMMDD_HHMMSS
📊 Backup size: ~2-5 GB
```

**3. Schedule daily backups:**

```bash
# Add to crontab
crontab -e

# Add this line (daily at 2 AM):
0 2 * * * /opt/onemind/scripts/backup.sh >> /var/log/onemind-backup.log 2>&1
```

**4. Verify crontab:**

```bash
crontab -l
```

**✅ Verification:**
```bash
# Check backup exists
ls -lh /backup/onemind/

# Check backup contents
ls /backup/onemind/backup_*/

# You should see:
# - postgres_data.tar.gz
# - agno_data.tar.gz
# - nats_data.tar.gz
# - configs.tar.gz
# - manifest.txt
```

### Task 2.2: Setup S3 Backup Sync (20 minutes)

**1. Install AWS CLI:**

```bash
sudo apt install -y awscli
aws --version
```

**2. Configure AWS credentials:**

```bash
aws configure
# Enter:
# - AWS Access Key ID
# - AWS Secret Access Key
# - Default region: us-east-1
# - Default output: json
```

**3. Create S3 bucket:**

```bash
# Via CLI
aws s3 mb s3://onemind-backups

# Or use AWS Console
```

**4. Enable S3 sync in backup script:**

Edit `/opt/onemind/scripts/backup.sh`:

```bash
# Change this line:
ENABLE_S3_SYNC=false

# To:
ENABLE_S3_SYNC=true
S3_BUCKET="s3://onemind-backups/azure"
```

**5. Test S3 sync:**

```bash
./scripts/backup.sh
```

**✅ Verification:**
```bash
# Check S3 bucket
aws s3 ls s3://onemind-backups/azure/ --recursive

# You should see backup files
```

### Task 2.3: Test Disaster Recovery (10 minutes)

**Practice recovery on test container:**

```bash
# Stop a service
docker compose -f containers/cloud/docker-compose.yml stop agno

# Remove its volume
docker volume rm onemind_agno_data

# Restore from backup
LATEST_BACKUP=$(ls -t /backup/onemind/ | head -1)
docker volume create onemind_agno_data
docker run --rm \
  -v onemind_agno_data:/target \
  -v "/backup/onemind/$LATEST_BACKUP":/backup:ro \
  alpine tar xzf /backup/agno_data.tar.gz -C /target

# Restart service
docker compose -f containers/cloud/docker-compose.yml start agno

# Verify health
docker compose ps | grep agno
```

**✅ Verification:**
```bash
# Agno should be healthy
curl http://localhost:8080/health
```

---

## Week 3: Monitoring & Alerting

**Goal:** Know when things go wrong, before users do
**Time:** 1-2 hours
**Status:** 🟢 Safe to implement

### Task 3.1: Configure Uptime Kuma (20 minutes)

**1. Access Uptime Kuma:**
```
http://YOUR_SERVER_IP:3001
```

**2. Create monitors for each service:**

| Monitor Name | Type | Check Interval | URL/Host |
|--------------|------|----------------|----------|
| Agno API | HTTP | 60s | http://localhost:8080/health |
| API Gateway | HTTP | 60s | http://localhost:8000/health |
| NATS | HTTP | 60s | http://localhost:8222/healthz |
| N8N | HTTP | 60s | http://localhost:5678/healthz |
| PostgreSQL | Port | 60s | localhost:5432 |

**3. Setup Discord webhook notification:**

- Settings → Notifications → Add Notification
- Type: Discord
- Webhook URL: (from Discord channel settings)
- Send to: All monitors

**✅ Verification:**
- Trigger a test notification
- All monitors should show green

### Task 3.2: Setup Resource Monitoring (20 minutes)

**1. Make monitor script executable:**

```bash
cd /opt/onemind
chmod +x scripts/monitor-resources.sh
```

**2. Test the script:**

```bash
./scripts/monitor-resources.sh
```

**Expected output:**
```
💾 Disk usage: XX% (XX GB available)
🧠 Memory usage: XX% (XX GB/XX GB)
⚡ CPU load: X.XX
🐳 All containers healthy
```

**3. Schedule monitoring (every 15 minutes):**

```bash
crontab -e

# Add this line:
*/15 * * * * /opt/onemind/scripts/monitor-resources.sh >> /var/log/onemind-resources.log 2>&1
```

**4. Configure Discord webhook in environment:**

```bash
# Edit .bashrc or .profile
echo 'export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/..."' >> ~/.bashrc
source ~/.bashrc
```

**✅ Verification:**
```bash
# Wait 15 minutes, then check log
tail -50 /var/log/onemind-resources.log

# Should see monitoring output
```

### Task 3.3: Setup Prometheus + Grafana (Optional, 30 minutes)

**For advanced metrics visualization.**

**1. Add services to docker-compose.yml:**

See [PRODUCTION_BEST_PRACTICES.md](PRODUCTION_BEST_PRACTICES.md#4-monitoring--alerting) for full configuration.

**2. Create prometheus.yml:**

```bash
mkdir -p /opt/onemind/containers/cloud/prometheus
nano /opt/onemind/containers/cloud/prometheus/prometheus.yml
```

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']

  - job_name: 'containers'
    static_configs:
      - targets: ['cadvisor:8080']
```

**3. Start monitoring stack:**

```bash
docker compose -f containers/cloud/docker-compose.yml up -d
```

**4. Access Grafana:**
```
http://YOUR_SERVER_IP:3002
Default login: admin/admin
```

**5. Add Prometheus data source:**
- Configuration → Data Sources → Add Prometheus
- URL: http://prometheus:9090
- Save & Test

**6. Import dashboards:**
- Node Exporter Full (ID: 1860)
- Docker Monitoring (ID: 10619)

**✅ Verification:**
- Open Grafana
- View Node Exporter dashboard
- See real-time metrics

---

## Week 4: Automation & Self-Healing

**Goal:** Zero-touch maintenance
**Time:** 1 hour
**Status:** 🟡 Test on staging first

### Task 4.1: Setup Watchtower (Container Auto-Updates) (15 minutes)

**1. Add Watchtower to docker-compose.yml:**

```yaml
watchtower:
  image: containrrr/watchtower:latest
  container_name: onemind-watchtower
  restart: unless-stopped
  environment:
    WATCHTOWER_CLEANUP: "true"
    WATCHTOWER_SCHEDULE: "0 0 3 * * SUN"  # 3 AM Sundays
    WATCHTOWER_INCLUDE_STOPPED: "false"
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  networks:
    - onemind
```

**2. Start Watchtower:**

```bash
cd /opt/onemind
docker compose -f containers/cloud/docker-compose.yml up -d watchtower
```

**3. Test manual update:**

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  containrrr/watchtower --run-once --cleanup
```

**✅ Verification:**
```bash
# Check Watchtower logs
docker logs onemind-watchtower

# Updates will happen automatically every Sunday at 3 AM
```

### Task 4.2: Create N8N Self-Healing Workflow (30 minutes)

**1. Access N8N:**
```
http://YOUR_SERVER_IP:5678
```

**2. Create new workflow: "OneMind Self-Healing"**

**Nodes:**

1. **Cron Trigger**
   - Every 5 minutes

2. **HTTP Request (Check Agno)**
   - Method: GET
   - URL: http://agno:8080/health

3. **IF (Is Unhealthy?)**
   - Condition: Status Code ≠ 200

4. **Execute Command (Restart)**
   - Command: `docker restart onemind-agno`

5. **Wait (2 seconds)**

6. **HTTP Request (Verify)**
   - Check health again

7. **Discord Notification**
   - "🔄 Agno was unhealthy and has been restarted"

**3. Repeat for other services:**
- API Gateway
- NATS
- Event Processor

**4. Activate workflow**

**✅ Verification:**
- Stop a service manually: `docker stop onemind-agno`
- Wait 5 minutes
- Service should auto-restart
- Check Discord for notification

### Task 4.3: Setup Zero-Downtime Deployment (15 minutes)

**1. Make deploy script executable:**

```bash
cd /opt/onemind
chmod +x scripts/deploy.sh
```

**2. Test deployment:**

```bash
# Dry run (check what would happen)
./scripts/deploy.sh

# This will:
# 1. Pull latest code
# 2. Pull new images
# 3. Update services one by one
# 4. Wait for health checks
# 5. Rollback on failure
```

**3. Configure Discord notifications:**

```bash
# Add to .bashrc
echo 'export DISCORD_WEBHOOK_URL="..."' >> ~/.bashrc
source ~/.bashrc
```

**✅ Verification:**
- Run `./scripts/deploy.sh`
- All services should update with zero downtime
- Discord notification on completion

---

## Ongoing Maintenance

### Daily (Automated)
```
✅ 2:00 AM - Backup to local disk
✅ 2:30 AM - Sync to S3
✅ 3:00 AM - Security updates check
✅ Every 15 min - Resource monitoring
✅ Every 5 min - Self-healing checks
```

### Weekly (Automated)
```
✅ Sunday 3 AM - Container image updates (Watchtower)
```

### Monthly (Manual)
```
📅 Review logs: /var/log/onemind-*.log
📅 Check disk usage: df -h
📅 Review security: sudo fail2ban-client status
📅 Test backup restore
📅 Update documentation
```

### Quarterly (Manual)
```
📅 Rotate credentials
📅 Security audit
📅 Disaster recovery drill
📅 Review resource usage and costs
```

---

## Quick Reference Commands

### Daily Operations

```bash
# View service status
docker compose -f containers/cloud/docker-compose.yml ps

# View logs (all services)
docker compose -f containers/cloud/docker-compose.yml logs -f

# View specific service logs
docker compose -f containers/cloud/docker-compose.yml logs -f agno

# Restart service
docker compose -f containers/cloud/docker-compose.yml restart agno

# Check resource usage
docker stats

# Check disk space
df -h
```

### Deployment

```bash
# Zero-downtime update
cd /opt/onemind && ./scripts/deploy.sh

# Quick restart (minimal downtime)
docker compose -f containers/cloud/docker-compose.yml restart

# Full rebuild
docker compose -f containers/cloud/docker-compose.yml down
docker compose -f containers/cloud/docker-compose.yml up -d --build
```

### Backup & Restore

```bash
# Manual backup
/opt/onemind/scripts/backup.sh

# List backups
ls -lh /backup/onemind/

# Restore specific service
LATEST_BACKUP=$(ls -t /backup/onemind/ | head -1)
docker compose stop <service>
docker volume rm onemind_<service>_data
docker volume create onemind_<service>_data
docker run --rm \
  -v onemind_<service>_data:/target \
  -v "/backup/onemind/$LATEST_BACKUP":/backup:ro \
  alpine tar xzf /backup/<service>_data.tar.gz -C /target
docker compose start <service>
```

### Troubleshooting

```bash
# Check health status
docker inspect --format='{{.State.Health.Status}}' onemind-agno

# View container details
docker inspect onemind-agno

# Check resource limits
docker inspect onemind-agno | grep -A 10 Resources

# Clean up disk space
docker system prune -af --volumes

# View system logs
journalctl -u docker -f

# Check firewall
sudo ufw status verbose

# Check fail2ban
sudo fail2ban-client status sshd
```

---

## Success Metrics

After implementing all weeks, you should have:

### Week 1 ✅
- [x] Automatic security updates enabled
- [x] Firewall configured with minimal ports
- [x] Fail2Ban protecting SSH
- [x] Docker logs rotating at 10 MB
- [x] Resource limits preventing runaway containers

### Week 2 ✅
- [x] Daily automated backups
- [x] S3 cloud backup sync
- [x] Tested disaster recovery
- [x] 30-day backup retention
- [x] RTO < 30 minutes, RPO < 24 hours

### Week 3 ✅
- [x] Uptime Kuma monitoring all services
- [x] Discord alerts on failures
- [x] Resource monitoring every 15 minutes
- [x] Auto-cleanup on disk/memory pressure
- [x] Optional: Grafana dashboards

### Week 4 ✅
- [x] Watchtower auto-updating containers weekly
- [x] N8N self-healing workflows
- [x] Zero-downtime deployment script
- [x] Automated recovery from failures

### Overall System Health
```
System Status: 🟢 PRODUCTION READY

✅ Stability: Services restart on failure
✅ Security: Automatic updates + firewall + fail2ban
✅ Resilience: Daily backups + disaster recovery tested
✅ Observability: Monitoring + alerting + dashboards
✅ Automation: Self-healing + auto-updates
✅ Portability: Cloud-agnostic Docker Compose

Uptime Target: 99.5% (43.8 hours downtime/year)
Recovery Time: < 30 minutes
Data Loss Risk: < 24 hours (daily backups)
```

---

## Need Help?

### Documentation
- [PRODUCTION_BEST_PRACTICES.md](PRODUCTION_BEST_PRACTICES.md) - Complete reference
- [CLOUD_MIGRATION.md](CLOUD_MIGRATION.md) - Move between clouds
- [HEALTHCHECKS.md](../containers/cloud/HEALTHCHECKS.md) - Health check reference

### Scripts
- `scripts/setup-production.sh` - Automated setup
- `scripts/backup.sh` - Backup automation
- `scripts/deploy.sh` - Zero-downtime deployment
- `scripts/monitor-resources.sh` - Resource monitoring

### Monitoring
- Uptime Kuma: http://YOUR_IP:3001
- Portainer: http://YOUR_IP:9000
- Grafana: http://YOUR_IP:3002 (if enabled)

### Logs
- `/var/log/onemind-backup.log` - Backup logs
- `/var/log/onemind-resources.log` - Resource monitoring
- `docker compose logs -f` - Application logs

---

## Next Steps

After completing this 4-week implementation:

1. **Monitor** - Watch metrics for 1-2 weeks
2. **Tune** - Adjust resource limits based on usage
3. **Optimize** - Identify bottlenecks and optimize
4. **Scale** - Add more resources if needed
5. **Document** - Keep runbooks updated

**You now have a production-grade, self-healing infrastructure! 🎉**
