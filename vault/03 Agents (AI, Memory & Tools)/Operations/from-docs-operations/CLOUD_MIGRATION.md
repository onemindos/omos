# OneMind OS - Cloud Migration Guide

## Overview

This guide provides step-by-step instructions for migrating OneMind OS between cloud providers (AWS Lightsail, Azure, GCP, etc.) with **minimal downtime**.

---

## Migration Strategy

### Zero-Vendor-Lock Design

OneMind OS is designed to be **100% portable** across cloud providers:

```
✅ Docker Compose (works everywhere)
✅ Standard Ubuntu VM (any provider)
✅ No cloud-specific services (RDS, Lambda, etc.)
✅ Documented port requirements
✅ Automated backup/restore
```

### Migration Scenarios

| From → To | Estimated Time | Downtime |
|-----------|----------------|----------|
| AWS → Azure | 1-2 hours | 15 minutes |
| Azure → GCP | 1-2 hours | 15 minutes |
| Local → Cloud | 1 hour | 5 minutes |
| Cloud → Local | 1 hour | None (parallel) |

---

## Pre-Migration Checklist

### 1. Current System Documentation

Document your current setup:

```bash
# On current server
cd /opt/onemind

# Save configuration
cat containers/cloud/.env > current-config.txt  # Remove before sharing!
docker compose -f containers/cloud/docker-compose.yml config > current-compose.yml

# Document versions
docker compose -f containers/cloud/docker-compose.yml images > current-images.txt
uname -a > system-info.txt

# Test current backup
scripts/backup.sh
```

### 2. Gather Required Information

**Current Server:**
- [ ] IP address
- [ ] SSH key location
- [ ] Domain names (if any)
- [ ] DNS records
- [ ] Firewall rules

**Credentials (from 1Password):**
- [ ] Database passwords
- [ ] API keys
- [ ] Webhook secrets
- [ ] SSL certificates

### 3. Create Final Backup

```bash
# Full backup before migration
ssh onemind-cloud "sudo /opt/onemind/scripts/backup.sh"

# Download backup locally
scp -r onemind-cloud:/backup/onemind/backup_* ./migration-backup/

# Or sync from S3
aws s3 sync s3://onemind-backups/lightsail/ ./migration-backup/
```

---

## Migration Process

### Phase 1: Provision New Server

#### Option A: AWS Lightsail

**Specs:**
- Instance type: Medium (2 vCPU, 4 GB RAM)
- OS: Ubuntu 22.04 LTS
- Storage: 80 GB SSD (minimum)
- Region: us-east-1 (or closest to you)

**Setup:**
```bash
# 1. Create instance via AWS Console or CLI
aws lightsail create-instances \
  --instance-names onemind-os-new \
  --availability-zone us-east-1a \
  --blueprint-id ubuntu_22_04 \
  --bundle-id medium_2_0 \
  --key-pair-name onemind-key \
  --tags key=Project,value=OneMind-OS

# 2. Wait for instance to be running
aws lightsail get-instance --instance-name onemind-os-new

# 3. Get IP address
NEW_IP=$(aws lightsail get-instance --instance-name onemind-os-new \
  --query 'instance.publicIpAddress' --output text)

echo "New server IP: $NEW_IP"

# 4. Open firewall ports
aws lightsail open-instance-public-ports \
  --instance-name onemind-os-new \
  --port-info fromPort=22,toPort=22,protocol=tcp
aws lightsail open-instance-public-ports \
  --instance-name onemind-os-new \
  --port-info fromPort=80,toPort=80,protocol=tcp
aws lightsail open-instance-public-ports \
  --instance-name onemind-os-new \
  --port-info fromPort=443,toPort=443,protocol=tcp
aws lightsail open-instance-public-ports \
  --instance-name onemind-os-new \
  --port-info fromPort=8080,toPort=8080,protocol=tcp
aws lightsail open-instance-public-ports \
  --instance-name onemind-os-new \
  --port-info fromPort=8000,toPort=8000,protocol=tcp
```

#### Option B: Azure VM

**Specs:**
- VM size: Standard_B2s (2 vCPU, 4 GB RAM)
- OS: Ubuntu 22.04 LTS
- Disk: 80 GB Premium SSD
- Region: East US

**Setup:**
```bash
# 1. Create resource group
az group create --name onemind-rg --location eastus

# 2. Create VM
az vm create \
  --resource-group onemind-rg \
  --name onemind-os-new \
  --image Ubuntu2204 \
  --size Standard_B2s \
  --admin-username ubuntu \
  --ssh-key-values ~/.ssh/id_rsa.pub \
  --public-ip-address-allocation static

# 3. Get IP address
NEW_IP=$(az vm show -d -g onemind-rg -n onemind-os-new \
  --query publicIps -o tsv)

echo "New server IP: $NEW_IP"

# 4. Open firewall ports
az vm open-port --resource-group onemind-rg --name onemind-os-new --port 22
az vm open-port --resource-group onemind-rg --name onemind-os-new --port 80
az vm open-port --resource-group onemind-rg --name onemind-os-new --port 443
az vm open-port --resource-group onemind-rg --name onemind-os-new --port 8080
az vm open-port --resource-group onemind-rg --name onemind-os-new --port 8000
```

#### Option C: Manual Provisioning (Any Provider)

**Requirements:**
- Ubuntu 22.04 LTS
- 2 vCPU, 4 GB RAM (minimum)
- 80 GB storage
- Static IP
- SSH access

**Firewall ports to open:**
```
22   - SSH
80   - HTTP
443  - HTTPS
8080 - Agno API
8000 - API Gateway
5678 - N8N (optional)
3001 - Uptime Kuma (optional)
9000 - Portainer (optional)
```

---

### Phase 2: Install Base System

**SSH into new server:**
```bash
ssh ubuntu@$NEW_IP
```

**1. Update system:**
```bash
sudo apt update
sudo apt upgrade -y
sudo apt install -y curl git
```

**2. Install Docker:**
```bash
# Official Docker installation
curl -fsSL https://get.docker.com | sh

# Add user to docker group
sudo usermod -aG docker $USER

# Log out and back in for group changes
exit
ssh ubuntu@$NEW_IP
```

**3. Verify Docker:**
```bash
docker --version
docker compose version
```

**4. Create project directory:**
```bash
sudo mkdir -p /opt/onemind
sudo chown $USER:$USER /opt/onemind
cd /opt/onemind
```

---

### Phase 3: Restore Data

**1. Copy backup to new server:**
```bash
# From your local machine
scp -r migration-backup/backup_latest ubuntu@$NEW_IP:/tmp/restore/
```

**2. Clone repository:**
```bash
# On new server
cd /opt/onemind
git clone https://github.com/your-org/onemind-os.git .
```

**3. Restore configuration:**
```bash
cd /opt/onemind

# Extract config backup
tar xzf /tmp/restore/configs.tar.gz

# Or manually create .env
nano containers/cloud/.env
# Paste your credentials from 1Password
```

**4. Restore Docker volumes:**
```bash
# Restore each volume
VOLUMES=(postgres_data agno_data nats_data n8n_data kestra_data timescale_data uptime_data caddy_data)

for volume in "${VOLUMES[@]}"; do
    echo "Restoring $volume..."

    # Create volume
    docker volume create "onemind_${volume}"

    # Restore data
    docker run --rm \
        -v "onemind_${volume}:/target" \
        -v /tmp/restore:/backup:ro \
        alpine tar xzf "/backup/${volume}.tar.gz" -C /target

    echo "✅ $volume restored"
done
```

---

### Phase 4: Start Services

**1. Start containers:**
```bash
cd /opt/onemind
docker compose -f containers/cloud/docker-compose.yml up -d
```

**2. Monitor startup:**
```bash
# Watch logs
docker compose -f containers/cloud/docker-compose.yml logs -f

# Check status
watch -n 2 'docker compose -f containers/cloud/docker-compose.yml ps'
```

**3. Wait for all services to be healthy:**
```bash
# This may take 2-3 minutes
docker compose -f containers/cloud/docker-compose.yml ps
```

**4. Verify health:**
```bash
# Test Agno API
curl http://localhost:8080/health

# Test API Gateway
curl http://localhost:8000/health

# Test NATS
curl http://localhost:8222/healthz

# Check all services
docker compose -f containers/cloud/docker-compose.yml ps | grep healthy
```

---

### Phase 5: Update DNS & Routing

**1. Test from local machine:**
```bash
# Replace with your new IP
NEW_IP="x.x.x.x"

# Test Agno
curl http://$NEW_IP:8080/health

# Test Gateway
curl http://$NEW_IP:8000/health
```

**2. Update DNS records (if using domain):**
```
# Example DNS records
onemind.example.com      A  → NEW_IP
api.onemind.example.com  A  → NEW_IP
```

**DNS propagation:** 5-15 minutes (varies by TTL)

**3. Update webhook URLs:**

Update these in external services:
- Todoist webhook URL
- GitHub webhook URL
- Zoho webhook URL
- Any other integrations

---

### Phase 6: Production Hardening

**1. Run production setup:**
```bash
cd /opt/onemind
scripts/setup-production.sh
```

**This will configure:**
- Automatic security updates
- Firewall rules
- Fail2Ban
- Docker log rotation
- Automated backups
- Resource monitoring

**2. Test backup:**
```bash
scripts/backup.sh
```

**3. Verify monitoring:**
```bash
# Check Uptime Kuma
open http://$NEW_IP:3001

# Check Portainer
open http://$NEW_IP:9000
```

---

### Phase 7: Parallel Testing

**Before switching traffic, test new server:**

**1. Create test webhooks:**
```bash
# Test API Gateway
curl -X POST http://$NEW_IP:8000/webhooks/todoist \
  -H "Content-Type: application/json" \
  -d '{"event": "test"}'
```

**2. Test Agno API:**
```bash
curl http://$NEW_IP:8080/awareness/mode
curl http://$NEW_IP:8080/context
```

**3. Test N8N workflows:**
- Access http://$NEW_IP:5678
- Trigger test workflow
- Verify external integrations work

**4. Monitor for issues:**
```bash
# Check logs
docker compose logs -f

# Check resource usage
docker stats

# Check disk space
df -h
```

---

### Phase 8: Cutover

**When ready to switch:**

**1. Update main webhook URLs:**
```bash
# Todoist
# Go to Todoist Settings → Integrations → Edit webhook
# Update URL to: http://NEW_IP:8000/webhooks/todoist

# GitHub
# Go to repo Settings → Webhooks → Edit
# Update URL to: http://NEW_IP:8000/webhooks/github
```

**2. Update DNS (if applicable):**
```
onemind.example.com  →  NEW_IP
```

**3. Monitor traffic on new server:**
```bash
docker compose logs -f api-gateway
docker compose logs -f agno
```

**4. Verify events flowing:**
- Check NATS monitoring: http://NEW_IP:8222
- Check Uptime Kuma status
- Test a real task/event

---

### Phase 9: Decommission Old Server

**Wait 24-48 hours** before decommissioning to ensure stability.

**1. Create final backup of old server:**
```bash
ssh OLD_IP "sudo /opt/onemind/scripts/backup.sh"
```

**2. Download any logs needed:**
```bash
scp OLD_IP:/var/log/onemind-*.log ./old-logs/
```

**3. Stop services on old server:**
```bash
ssh OLD_IP
cd /opt/onemind
docker compose down
```

**4. Delete instance:**

**AWS:**
```bash
aws lightsail delete-instance --instance-name onemind-os-old
```

**Azure:**
```bash
az vm delete --resource-group onemind-rg --name onemind-os-old
```

---

## Rollback Procedure

**If migration fails, rollback is simple:**

**1. Revert DNS (if changed):**
```
onemind.example.com  →  OLD_IP
```

**2. Revert webhook URLs:**
- Point webhooks back to old server

**3. Old server is still running:**
- No data loss
- No downtime (if DNS not changed yet)

**Rollback time:** < 5 minutes

---

## Troubleshooting

### Services not starting

**Check logs:**
```bash
docker compose logs <service-name>
```

**Common issues:**
- Missing .env variables
- Volume permissions
- Port conflicts
- Health check timeouts

### Database connection failures

**Check PostgreSQL:**
```bash
docker exec onemind-postgres pg_isready -U onemind
```

**Reset database (CAUTION: data loss):**
```bash
docker compose down
docker volume rm onemind_postgres_data
docker compose up -d
```

### Network issues

**Check Docker network:**
```bash
docker network ls
docker network inspect onemind
```

**Recreate network:**
```bash
docker compose down
docker network prune
docker compose up -d
```

### Disk space issues

**Check space:**
```bash
df -h
docker system df
```

**Clean up:**
```bash
docker system prune -af --volumes
```

---

## Cloud Provider Comparison

| Feature | AWS Lightsail | Azure VM | GCP Compute Engine |
|---------|---------------|----------|-------------------|
| **Pricing (2vCPU, 4GB)** | $20/month | $30/month | $25/month |
| **Setup Difficulty** | Easy | Medium | Medium |
| **CLI Tool** | `aws lightsail` | `az vm` | `gcloud compute` |
| **Firewall** | Lightsail firewall | NSG | VPC Firewall |
| **Backup** | Snapshots ($0.05/GB) | Snapshots | Persistent Disk Snapshots |
| **Static IP** | Free (1 per instance) | $3.65/month | Free (1 per region) |
| **Migration Ease** | High | High | High |

**Recommendation:** AWS Lightsail for simplicity and cost.

---

## Automation (Future)

**Terraform for multi-cloud:**
```hcl
# Coming soon: /opt/onemind/terraform/
# - AWS module
# - Azure module
# - GCP module
```

**One-command migration:**
```bash
# Vision for future
./scripts/migrate.sh --from aws --to azure
```

---

## Testing Checklist

After migration, verify:

- [ ] All containers running and healthy
- [ ] Agno API responding
- [ ] API Gateway receiving webhooks
- [ ] NATS streaming events
- [ ] N8N workflows executing
- [ ] Backups running on schedule
- [ ] Monitoring alerts working
- [ ] Resource usage normal
- [ ] DNS resolving correctly
- [ ] SSL certificates valid (if using domain)
- [ ] External integrations working

---

## Estimated Costs

### AWS Lightsail
- Instance: $20/month (2 vCPU, 4 GB RAM, 80 GB SSD)
- Static IP: Free
- Backups: ~$5/month (S3 Glacier)
- **Total: ~$25/month**

### Azure
- VM: $30/month (Standard_B2s)
- Static IP: $3.65/month
- Backups: ~$5/month (Blob Storage)
- **Total: ~$39/month**

### Self-Hosted (Raspberry Pi / Mini PC)
- Hardware: $100-300 (one-time)
- Electricity: ~$2/month
- Internet: Existing
- **Total: ~$2/month (after hardware)**

---

## Conclusion

OneMind OS migration is **straightforward and low-risk**:

✅ **No vendor lock-in** - Docker Compose runs anywhere
✅ **Minimal downtime** - Parallel testing before cutover
✅ **Easy rollback** - Old server stays running
✅ **Documented process** - Step-by-step guide
✅ **Tested backups** - Restore procedure verified

**Migration time:** 1-2 hours
**Downtime:** < 15 minutes (only DNS propagation)

For questions or issues during migration, check logs and refer to the troubleshooting section.
