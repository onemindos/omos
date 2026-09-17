# OneMind OS - Production Best Practices

## Table of Contents

1. [System Hardening](#system-hardening)
2. [Automated Updates](#automated-updates)
3. [Backup & Disaster Recovery](#backup--disaster-recovery)
4. [Monitoring & Alerting](#monitoring--alerting)
5. [Self-Healing Mechanisms](#self-healing-mechanisms)
6. [Resource Management](#resource-management)
7. [Log Management](#log-management)
8. [Security Workflows](#security-workflows)
9. [Deployment Strategy](#deployment-strategy)
10. [Cloud Portability](#cloud-portability)

---

## 1. System Hardening

### Ubuntu Server Security

#### Unattended Upgrades (Automatic Security Updates)

**Install and configure:**
```bash
sudo apt update
sudo apt install -y unattended-upgrades apt-listchanges

# Configure unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

**Configuration file:** `/etc/apt/apt.conf.d/50unattended-upgrades`

```bash
# Edit the config
sudo nano /etc/apt/apt.conf.d/50unattended-upgrades
```

**Recommended settings:**
```
Unattended-Upgrade::Allowed-Origins {
    "${distro_id}:${distro_codename}-security";
    "${distro_id}ESMApps:${distro_codename}-apps-security";
    "${distro_id}ESM:${distro_codename}-infra-security";
};

// Automatically reboot if required
Unattended-Upgrade::Automatic-Reboot "true";
Unattended-Upgrade::Automatic-Reboot-Time "03:00";

// Email notifications
Unattended-Upgrade::Mail "your-email@example.com";
Unattended-Upgrade::MailReport "on-change";

// Remove unused dependencies
Unattended-Upgrade::Remove-Unused-Dependencies "true";
Unattended-Upgrade::Remove-Unused-Kernel-Packages "true";
```

**Auto-update intervals:** `/etc/apt/apt.conf.d/20auto-upgrades`
```
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Download-Upgradeable-Packages "1";
APT::Periodic::AutocleanInterval "7";
APT::Periodic::Unattended-Upgrade "1";
```

#### Firewall Configuration

```bash
# Install ufw
sudo apt install -y ufw

# Default policies
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH (CRITICAL - do this first!)
sudo ufw allow 22/tcp comment 'SSH'

# Allow OneMind OS ports
sudo ufw allow 8080/tcp comment 'Agno API'
sudo ufw allow 8000/tcp comment 'API Gateway'
sudo ufw allow 4222/tcp comment 'NATS Client'
sudo ufw allow 8222/tcp comment 'NATS Monitoring'
sudo ufw allow 5678/tcp comment 'N8N'
sudo ufw allow 3001/tcp comment 'Uptime Kuma'
sudo ufw allow 9000/tcp comment 'Portainer'
sudo ufw allow 80/tcp comment 'HTTP'
sudo ufw allow 443/tcp comment 'HTTPS'

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status verbose
```

#### SSH Hardening

```bash
# Edit SSH config
sudo nano /etc/ssh/sshd_config
```

**Recommended settings:**
```
# Disable root login
PermitRootLogin no

# Use key-based authentication only
PasswordAuthentication no
PubkeyAuthentication yes

# Limit login attempts
MaxAuthTries 3

# Connection timeouts
ClientAliveInterval 300
ClientAliveCountMax 2

# Disable unused features
X11Forwarding no
PermitEmptyPasswords no
```

```bash
# Restart SSH
sudo systemctl restart sshd
```

#### Fail2Ban (Brute Force Protection)

```bash
# Install fail2ban
sudo apt install -y fail2ban

# Create local config
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
sudo nano /etc/fail2ban/jail.local
```

**Configuration:**
```ini
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5
destemail = your-email@example.com
sendername = Fail2Ban
action = %(action_mwl)s

[sshd]
enabled = true
port = 22
logpath = /var/log/auth.log
```

```bash
# Enable and start
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Check status
sudo fail2ban-client status sshd
```

---

## 2. Automated Updates

### Container Image Updates

#### Strategy: Weekly Automated Updates with Watchtower

**Deploy Watchtower container:**
```yaml
# Add to docker-compose.yml
watchtower:
  image: containrrr/watchtower:latest
  container_name: onemind-watchtower
  restart: unless-stopped
  environment:
    WATCHTOWER_CLEANUP: "true"
    WATCHTOWER_POLL_INTERVAL: 604800  # 7 days in seconds
    WATCHTOWER_SCHEDULE: "0 0 3 * * SUN"  # 3 AM every Sunday
    WATCHTOWER_NOTIFICATIONS: "email"
    WATCHTOWER_NOTIFICATION_EMAIL_FROM: "watchtower@onemind.local"
    WATCHTOWER_NOTIFICATION_EMAIL_TO: "${ADMIN_EMAIL}"
    WATCHTOWER_NOTIFICATION_EMAIL_SERVER: "smtp.gmail.com"
    WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PORT: "587"
    WATCHTOWER_NOTIFICATION_EMAIL_SERVER_USER: "${SMTP_USER}"
    WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PASSWORD: "${SMTP_PASSWORD}"
    WATCHTOWER_INCLUDE_STOPPED: "false"
    WATCHTOWER_REVIVE_STOPPED: "false"
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  networks:
    - onemind
```

**Manual update command:**
```bash
# Update all containers NOW
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  containrrr/watchtower --run-once --cleanup
```

#### Alternative: N8N Workflow for Updates

**Workflow: Weekly Container Updates**
1. Trigger: Cron (Sunday 3 AM)
2. Check for image updates
3. Pull new images
4. Restart services one by one
5. Verify health checks
6. Send notification (Discord/Email)
7. Rollback if health checks fail

---

## 3. Backup & Disaster Recovery

### Backup Strategy

#### What to Back Up

**Critical Data (Docker Volumes):**
```
✅ postgres_data      - Agent memory, sessions
✅ agno_data          - Agent state
✅ nats_data          - Event stream
✅ n8n_data           - Workflows
✅ kestra_data        - Workflow definitions
✅ timescale_data     - Historical events
✅ uptime_data        - Monitoring history
✅ caddy_data         - SSL certificates
❌ portainer_data     - UI state (recreatable)
❌ tailscale_data     - VPN state (recreatable)
```

**Configuration Files:**
```
✅ docker-compose.yml
✅ .env (encrypted)
✅ Caddyfile
✅ nats.conf
✅ Custom configs
```

#### Backup Script

Create `/opt/onemind/backup.sh`:

```bash
#!/bin/bash
set -e

# Configuration
BACKUP_DIR="/backup/onemind"
RETENTION_DAYS=30
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH="$BACKUP_DIR/backup_$TIMESTAMP"

# Create backup directory
mkdir -p "$BACKUP_PATH"

echo "🔄 Starting OneMind OS backup..."

# Stop containers (optional - for consistent backups)
# cd /opt/onemind && docker compose stop

# Backup Docker volumes
echo "📦 Backing up Docker volumes..."
for volume in postgres_data agno_data nats_data n8n_data kestra_data timescale_data uptime_data caddy_data; do
    echo "  - $volume"
    docker run --rm \
        -v onemind_${volume}:/source:ro \
        -v "$BACKUP_PATH":/backup \
        alpine tar czf /backup/${volume}.tar.gz -C /source .
done

# Backup configuration files
echo "📄 Backing up configuration files..."
cd /opt/onemind
tar czf "$BACKUP_PATH/configs.tar.gz" \
    containers/cloud/docker-compose.yml \
    containers/cloud/Caddyfile \
    containers/cloud/nats.conf \
    containers/cloud/.env

# Restart containers
# docker compose start

# Create backup manifest
echo "📝 Creating backup manifest..."
cat > "$BACKUP_PATH/manifest.txt" <<EOF
OneMind OS Backup
Date: $(date)
Hostname: $(hostname)
Docker Version: $(docker --version)

Volumes:
$(ls -lh "$BACKUP_PATH"/*.tar.gz)

Total Size: $(du -sh "$BACKUP_PATH" | cut -f1)
EOF

# Cleanup old backups
echo "🧹 Cleaning up old backups (older than $RETENTION_DAYS days)..."
find "$BACKUP_DIR" -name "backup_*" -type d -mtime +$RETENTION_DAYS -exec rm -rf {} \;

echo "✅ Backup complete: $BACKUP_PATH"
echo "📊 Backup size: $(du -sh "$BACKUP_PATH" | cut -f1)"
```

**Make executable and schedule:**
```bash
chmod +x /opt/onemind/backup.sh

# Add to crontab (daily at 2 AM)
crontab -e
```

```cron
# OneMind OS Daily Backup
0 2 * * * /opt/onemind/backup.sh >> /var/log/onemind-backup.log 2>&1
```

#### S3 Backup (Cloud Storage)

**Install AWS CLI:**
```bash
sudo apt install -y awscli
```

**Configure AWS (or use IAM role):**
```bash
aws configure
```

**Sync backups to S3:**
```bash
# Add to backup.sh
echo "☁️ Uploading to S3..."
aws s3 sync "$BACKUP_DIR" s3://onemind-backups/azure/ \
    --storage-class GLACIER_IR \
    --exclude "*" \
    --include "backup_*/*"
```

**S3 Lifecycle Policy (auto-delete old backups):**
```json
{
  "Rules": [
    {
      "Id": "DeleteOldBackups",
      "Status": "Enabled",
      "Prefix": "azure/",
      "Expiration": {
        "Days": 90
      }
    }
  ]
}
```

### Disaster Recovery

#### Recovery Procedure

**1. Provision new VM (same specs)**
```bash
# Azure VM or Azure VM
# Ubuntu 22.04 LTS
# 2 vCPU, 4 GB RAM, 80 GB SSD (minimum)
```

**2. Install Docker**
```bash
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

**3. Restore from backup**
```bash
# Download latest backup from S3
aws s3 sync s3://onemind-backups/azure/backup_latest /tmp/restore

# Extract volumes
cd /opt/onemind
for volume in postgres_data agno_data nats_data n8n_data kestra_data timescale_data uptime_data caddy_data; do
    echo "Restoring $volume..."
    docker volume create onemind_${volume}
    docker run --rm \
        -v onemind_${volume}:/target \
        -v /tmp/restore:/backup:ro \
        alpine tar xzf /backup/${volume}.tar.gz -C /target
done

# Extract configs
tar xzf /tmp/restore/configs.tar.gz
```

**4. Start services**
```bash
docker compose -f containers/cloud/docker-compose.yml up -d
```

**5. Verify**
```bash
docker compose ps
curl http://localhost:8080/health
```

**Recovery Time Objective (RTO):** < 30 minutes
**Recovery Point Objective (RPO):** < 24 hours (daily backups)

---

## 4. Monitoring & Alerting

### Uptime Kuma Configuration

**Access:** http://20.121.38.186:3001

#### Monitors to Add

| Service | Type | Check |
|---------|------|-------|
| Agno API | HTTP | http://localhost:8080/health |
| API Gateway | HTTP | http://localhost:8000/health |
| NATS | HTTP | http://localhost:8222/healthz |
| N8N | HTTP | http://localhost:5678/healthz |
| PostgreSQL | Port | localhost:5432 |
| Disk Space | Script | df -h / |
| Memory Usage | Script | free -m |
| Docker Health | Script | docker ps --format "{{.Status}}" |

#### Notification Channels

**Discord Webhook:**
```
Settings → Notifications → Add Notification
Type: Discord
Webhook URL: https://discord.com/api/webhooks/...
```

**Email:**
```
Type: SMTP Email
Host: smtp.gmail.com
Port: 587
From: alerts@onemind.local
To: your-email@example.com
```

### Prometheus + Grafana (Advanced Monitoring)

**Add to docker-compose.yml:**

```yaml
# Prometheus (metrics collection)
prometheus:
  image: prom/prometheus:latest
  container_name: onemind-prometheus
  restart: unless-stopped
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
    - prometheus_data:/prometheus
  command:
    - '--config.file=/etc/prometheus/prometheus.yml'
    - '--storage.tsdb.retention.time=30d'
  ports:
    - "9090:9090"
  networks:
    - onemind

# Node Exporter (system metrics)
node-exporter:
  image: prom/node-exporter:latest
  container_name: onemind-node-exporter
  restart: unless-stopped
  command:
    - '--path.procfs=/host/proc'
    - '--path.sysfs=/host/sys'
    - '--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)'
  volumes:
    - /proc:/host/proc:ro
    - /sys:/host/sys:ro
    - /:/rootfs:ro
  networks:
    - onemind

# cAdvisor (container metrics)
cadvisor:
  image: gcr.io/cadvisor/cadvisor:latest
  container_name: onemind-cadvisor
  restart: unless-stopped
  volumes:
    - /:/rootfs:ro
    - /var/run:/var/run:ro
    - /sys:/sys:ro
    - /var/lib/docker/:/var/lib/docker:ro
  networks:
    - onemind

# Grafana (visualization)
grafana:
  image: grafana/grafana:latest
  container_name: onemind-grafana
  restart: unless-stopped
  environment:
    GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_PASSWORD}
    GF_INSTALL_PLUGINS: grafana-clock-panel
  volumes:
    - grafana_data:/var/lib/grafana
  ports:
    - "3002:3000"
  networks:
    - onemind
  depends_on:
    - prometheus
```

**prometheus.yml:**
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']

  - job_name: 'cadvisor'
    static_configs:
      - targets: ['cadvisor:8080']

  - job_name: 'agno'
    static_configs:
      - targets: ['agno:8080']
    metrics_path: '/metrics'

  - job_name: 'nats'
    static_configs:
      - targets: ['nats:8222']
    metrics_path: '/metrics'
```

**Grafana Dashboards:**
- Node Exporter Full (ID: 1860)
- Docker and Container Monitoring (ID: 10619)
- Postgres Database (ID: 9628)

---

## 5. Self-Healing Mechanisms

### Docker Auto-Restart Policies

**Already configured:** `restart: unless-stopped`

This handles:
- Container crashes
- Service failures
- System reboots

### Healthcheck-Based Restart

**Add to docker-compose.yml:**

```yaml
# Example service with enhanced health check
agno:
  # ... existing config ...
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 5
    start_period: 30s
  # Docker will automatically restart on health check failure (with restart policy)
```

### N8N Self-Healing Workflow

**Workflow: Service Health Monitor**

1. **Trigger:** Every 5 minutes
2. **Check:** Query all service health endpoints
3. **Detect:** If service unhealthy
4. **Action:** Restart specific container
5. **Verify:** Check if health restored
6. **Alert:** Send notification if restart fails
7. **Escalate:** Try full stack restart if single service restart fails

**N8N Workflow JSON:**
```json
{
  "name": "OneMind Self-Healing",
  "nodes": [
    {
      "type": "n8n-nodes-base.cron",
      "name": "Every 5 Minutes",
      "parameters": {
        "cronExpression": "*/5 * * * *"
      }
    },
    {
      "type": "n8n-nodes-base.httpRequest",
      "name": "Check Agno Health",
      "parameters": {
        "url": "http://agno:8080/health",
        "options": {
          "timeout": 5000
        }
      }
    },
    {
      "type": "n8n-nodes-base.if",
      "name": "Is Unhealthy?",
      "parameters": {
        "conditions": {
          "number": [
            {
              "value1": "={{$json.statusCode}}",
              "operation": "notEqual",
              "value2": 200
            }
          ]
        }
      }
    },
    {
      "type": "n8n-nodes-base.executeCommand",
      "name": "Restart Container",
      "parameters": {
        "command": "docker restart onemind-agno"
      }
    },
    {
      "type": "n8n-nodes-base.discord",
      "name": "Alert Discord",
      "parameters": {
        "text": "🚨 Agno service was unhealthy and has been restarted"
      }
    }
  ]
}
```

### System Resource Monitoring

**Script:** `/opt/onemind/monitor-resources.sh`

```bash
#!/bin/bash

# Thresholds
DISK_THRESHOLD=85
MEMORY_THRESHOLD=90
CPU_THRESHOLD=90

# Check disk usage
DISK_USAGE=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
if [ $DISK_USAGE -gt $DISK_THRESHOLD ]; then
    echo "⚠️ Disk usage: ${DISK_USAGE}%"
    # Clean old Docker data
    docker system prune -af --volumes --filter "until=24h"
fi

# Check memory usage
MEMORY_USAGE=$(free | grep Mem | awk '{print int($3/$2 * 100)}')
if [ $MEMORY_USAGE -gt $MEMORY_THRESHOLD ]; then
    echo "⚠️ Memory usage: ${MEMORY_USAGE}%"
    # Restart high-memory containers
    docker restart onemind-kestra
fi

# Check CPU load
CPU_LOAD=$(uptime | awk -F'load average:' '{print $2}' | awk '{print $1}' | sed 's/,//')
CPU_THRESHOLD_VALUE=2.0  # For 2 vCPU
if (( $(echo "$CPU_LOAD > $CPU_THRESHOLD_VALUE" | bc -l) )); then
    echo "⚠️ High CPU load: $CPU_LOAD"
fi
```

```bash
# Add to crontab (every 15 minutes)
*/15 * * * * /opt/onemind/monitor-resources.sh >> /var/log/onemind-resources.log 2>&1
```

---

## 6. Resource Management

### Container Resource Limits

**Add to docker-compose.yml:**

```yaml
services:
  # Critical services - guaranteed resources
  postgres:
    # ... existing config ...
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M

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

  nats:
    # ... existing config ...
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M

  # Non-critical services - limited resources
  n8n:
    # ... existing config ...
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M

  uptime-kuma:
    # ... existing config ...
    deploy:
      resources:
        limits:
          cpus: '0.25'
          memory: 256M

  kestra:
    # ... existing config ...
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 1G
```

**Resource Allocation Summary (2 vCPU, 4 GB RAM):**

| Service | CPU Limit | Memory Limit | Priority |
|---------|-----------|--------------|----------|
| Postgres | 1.0 | 1 GB | Critical |
| Agno | 1.5 | 2 GB | Critical |
| NATS | 0.5 | 512 MB | Critical |
| Gateway | 0.25 | 256 MB | High |
| Processor | 0.25 | 256 MB | High |
| Actions | 0.25 | 256 MB | High |
| N8N | 0.5 | 512 MB | Medium |
| Kestra | 0.5 | 1 GB | Medium |
| Uptime Kuma | 0.25 | 256 MB | Low |
| **Total** | **~5.5** | **~6 GB** | - |

**Note:** Limits exceed physical resources (overcommit). Docker will share resources dynamically. Critical services have reservations guaranteed.

---

## 7. Log Management

### Log Rotation

**Configure Docker log rotation:**

`/etc/docker/daemon.json`
```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

```bash
# Restart Docker
sudo systemctl restart docker
```

### Centralized Logging (Optional)

**Loki + Promtail (Grafana Logs)**

Add to docker-compose.yml:

```yaml
loki:
  image: grafana/loki:latest
  container_name: onemind-loki
  restart: unless-stopped
  ports:
    - "3100:3100"
  volumes:
    - loki_data:/loki
  networks:
    - onemind

promtail:
  image: grafana/promtail:latest
  container_name: onemind-promtail
  restart: unless-stopped
  volumes:
    - /var/log:/var/log:ro
    - /var/lib/docker/containers:/var/lib/docker/containers:ro
    - ./promtail-config.yml:/etc/promtail/config.yml:ro
  command: -config.file=/etc/promtail/config.yml
  networks:
    - onemind
```

---

## 8. Security Workflows

### Weekly Security Audit (N8N Workflow)

**Workflow: Security Health Check**

1. Check for available Ubuntu security updates
2. Check for Docker image vulnerabilities (Trivy)
3. Review firewall rules
4. Check for failed SSH attempts (fail2ban)
5. Audit open ports
6. Generate security report
7. Send to Discord/Email

### Credential Rotation

**Quarterly credential rotation:**

1. Generate new passwords
2. Update `.env` file
3. Restart affected services
4. Update 1Password vault
5. Test all integrations

**Script:** `/opt/onemind/rotate-credentials.sh`

---

## 9. Deployment Strategy

### Zero-Downtime Deployment

**Strategy: Rolling Updates**

1. Pull new code/images
2. Update one service at a time
3. Wait for health check
4. Continue to next service
5. Rollback if any service fails

**Script:** `/opt/onemind/deploy.sh`

```bash
#!/bin/bash
set -e

echo "🚀 OneMind OS Deployment Starting..."

cd /opt/onemind/containers/cloud

# Pull latest code
echo "📥 Pulling latest changes..."
git pull origin main

# Pull new images (don't restart yet)
echo "🐳 Pulling Docker images..."
docker compose pull

# Update services one by one
SERVICES="postgres nats agno api-gateway event-processor action-gateway n8n"

for service in $SERVICES; do
    echo "🔄 Updating $service..."

    # Restart service
    docker compose up -d --no-deps $service

    # Wait for health check
    echo "⏳ Waiting for $service to be healthy..."
    timeout 60 bash -c "
        while [ \$(docker inspect --format='{{.State.Health.Status}}' onemind-$service 2>/dev/null) != 'healthy' ]; do
            sleep 2
        done
    " || {
        echo "❌ $service failed health check! Rolling back..."
        docker compose up -d --no-deps $service
        exit 1
    }

    echo "✅ $service updated successfully"
done

echo "✅ Deployment complete!"
docker compose ps
```

### Blue-Green Deployment (Future)

For zero-downtime upgrades of the entire stack:

1. Run two identical environments (blue/green)
2. Update green environment
3. Switch Caddy routing to green
4. Monitor for issues
5. Keep blue as rollback option

---

## 10. Cloud Portability

### Infrastructure as Code

**Terraform Module for VM Provisioning**

`/opt/onemind/terraform/vm/main.tf`

```hcl
variable "cloud_provider" {
  type    = string
  default = "aws"  # or "azure", "gcp"
}

variable "vm_size" {
  type = string
  default = "medium"  # small, medium, large
}

locals {
  vm_configs = {
    aws = {
      small  = "t3.small"   # 2 vCPU, 2 GB
      medium = "t3.medium"  # 2 vCPU, 4 GB
      large  = "t3.large"   # 2 vCPU, 8 GB
    }
    azure = {
      small  = "Standard_B2s"   # 2 vCPU, 4 GB
      medium = "Standard_B2ms"  # 2 vCPU, 8 GB
      large  = "Standard_B4ms"  # 4 vCPU, 16 GB
    }
  }
}

# Azure VM
resource "aws_lightsail_instance" "onemind" {
  count             = var.cloud_provider == "aws" ? 1 : 0
  name              = "onemind-os"
  availability_zone = "us-east-1a"
  blueprint_id      = "ubuntu_22_04"
  bundle_id         = local.vm_configs.aws[var.vm_size]

  user_data = file("${path.module}/user-data.sh")

  tags = {
    Project = "OneMind-OS"
  }
}

# Azure VM
resource "azurerm_virtual_machine" "onemind" {
  count               = var.cloud_provider == "azure" ? 1 : 0
  name                = "onemind-os"
  location            = "East US"
  resource_group_name = "onemind-rg"
  vm_size             = local.vm_configs.azure[var.vm_size]

  # ... Azure-specific config ...
}
```

**User Data Script (Cloud-Init):**

`user-data.sh`
```bash
#!/bin/bash

# Install Docker
curl -fsSL https://get.docker.com | sh

# Create project directory
mkdir -p /opt/onemind
cd /opt/onemind

# Clone repository
git clone https://github.com/your-org/onemind-os.git .

# Create .env from secrets
cat > containers/cloud/.env <<EOF
POSTGRES_PASSWORD=${db_password}
API_SECRET_KEY=${api_secret}
# ... other secrets from cloud secrets manager ...
EOF

# Start services
docker compose -f containers/cloud/docker-compose.yml up -d

# Setup backups
cp scripts/backup.sh /opt/onemind/backup.sh
chmod +x /opt/onemind/backup.sh
echo "0 2 * * * /opt/onemind/backup.sh" | crontab -
```

### Deployment Checklist (Any Cloud)

#### Prerequisites
- [ ] Ubuntu 22.04 LTS VM
- [ ] 2+ vCPU, 4+ GB RAM, 80+ GB storage
- [ ] Static IP or domain name
- [ ] Firewall rules configured
- [ ] SSH key access

#### Setup Steps
1. [ ] Provision VM (AWS/Azure/GCP/etc)
2. [ ] SSH into VM
3. [ ] Install Docker
4. [ ] Clone OneMind-OS repo
5. [ ] Configure `.env` with secrets
6. [ ] Run `docker compose up -d`
7. [ ] Configure backups
8. [ ] Configure monitoring
9. [ ] Setup automated updates
10. [ ] Test all services

#### Migration Steps (Cloud → Cloud)
1. [ ] Backup current VM (full backup)
2. [ ] Provision new VM in target cloud
3. [ ] Install Docker
4. [ ] Restore backup
5. [ ] Update DNS/firewall rules
6. [ ] Start services
7. [ ] Verify all services healthy
8. [ ] Update integrations (webhooks)
9. [ ] Decommission old VM

**Estimated Migration Time:** 1-2 hours
**Downtime:** ~15 minutes (DNS propagation)

---

## Implementation Priority

### Phase 1: Critical Stability (Week 1)
1. ✅ Enable unattended-upgrades
2. ✅ Configure UFW firewall
3. ✅ Setup SSH hardening
4. ✅ Add resource limits to docker-compose
5. ✅ Configure Docker log rotation

### Phase 2: Backup & Recovery (Week 2)
1. ✅ Create backup script
2. ✅ Setup cron for daily backups
3. ✅ Configure S3 sync
4. ✅ Test restore procedure
5. ✅ Document DR process

### Phase 3: Monitoring & Alerting (Week 3)
1. ✅ Configure Uptime Kuma monitors
2. ✅ Setup Discord notifications
3. ✅ Deploy Prometheus + Grafana (optional)
4. ✅ Create monitoring dashboards
5. ✅ Setup N8N self-healing workflow

### Phase 4: Automation (Week 4)
1. ✅ Deploy Watchtower for auto-updates
2. ✅ Create deployment script
3. ✅ Setup security audit workflow
4. ✅ Create resource monitoring script
5. ✅ Test zero-downtime deployment

### Phase 5: Cloud Portability (Ongoing)
1. ✅ Create Terraform modules
2. ✅ Document migration process
3. ✅ Test deployment on Azure
4. ✅ Create cloud-agnostic runbook

---

## Maintenance Calendar

### Daily (Automated)
- Backup at 2 AM
- Security updates check
- Resource monitoring

### Weekly (Automated)
- Container image updates (Sunday 3 AM)
- Security audit report
- Health status summary

### Monthly (Manual Review)
- Review logs and alerts
- Check backup integrity
- Update documentation
- Review resource usage

### Quarterly (Manual)
- Credential rotation
- Disaster recovery drill
- Security audit
- Cost optimization review

---

## Quick Reference

### Emergency Commands

```bash
# View all services
docker compose ps

# Restart all services
docker compose restart

# View logs (last 100 lines)
docker compose logs --tail=100

# Restart specific service
docker compose restart agno

# Check service health
docker inspect --format='{{.State.Health.Status}}' onemind-agno

# Emergency stop
docker compose down

# Nuclear restart (loses data in ephemeral containers)
docker compose down && docker compose up -d

# Check disk space
df -h

# Check memory usage
free -m

# Check Docker resource usage
docker stats
```

### Access URLs

| Service | URL |
|---------|-----|
| Agno API | http://20.121.38.186:8080 |
| API Gateway | http://20.121.38.186:8000 |
| N8N | http://20.121.38.186:5678 |
| Uptime Kuma | http://20.121.38.186:3001 |
| Portainer | http://20.121.38.186:9000 |
| Prometheus | http://20.121.38.186:9090 |
| Grafana | http://20.121.38.186:3002 |

---

## Conclusion

This production best practices guide provides a **comprehensive, automated, and self-healing** infrastructure for OneMind OS. The key principles are:

1. **Stability** - Automated updates, health checks, restarts
2. **Resilience** - Backups, disaster recovery, self-healing
3. **Observability** - Monitoring, alerting, logging
4. **Security** - Hardening, audits, credential management
5. **Portability** - Cloud-agnostic, IaC, documented processes

**With these practices, OneMind OS will:**
- ✅ Automatically recover from failures
- ✅ Stay secure with automated updates
- ✅ Preserve data with daily backups
- ✅ Alert you when intervention needed
- ✅ Run efficiently with resource limits
- ✅ Migrate easily between clouds

**Next Steps:** Follow the 5-phase implementation plan over 4 weeks to transform your infrastructure from good to production-grade.
