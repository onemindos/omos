# OneMind OS Infrastructure Documentation

Complete documentation of OneMind OS infrastructure, resources, networking, CI/CD, and DNS configuration.

---

## Table of Contents

1. [Infrastructure Overview](#infrastructure-overview)
2. [Azure Resources](#azure-resources)
3. [Networking & DNS](#networking--dns)
4. [CI/CD Pipeline](#cicd-pipeline)
5. [VM Specifications](#vm-specifications)
6. [Cost Analysis](#cost-analysis)
7. [Monitoring & Observability](#monitoring--observability)
8. [Security](#security)
9. [Backup & Disaster Recovery](#backup--disaster-recovery)

---

## Infrastructure Overview

**OneMind OS** (stylized as **onemindos**) is a multi-agent AI orchestration platform hosted on Azure infrastructure with Tailscale for secure networking and DNS.

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        OneMind OS Stack                          │
│                                                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │  Frontend   │  │  Backend    │  │  Database   │            │
│  │ Flutter Web │  │  Agno API   │  │ TimescaleDB │            │
│  │  Port 3000  │  │  Port 8080  │  │  Port 5432  │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
│         │                 │                 │                    │
│         └─────────────────┴─────────────────┘                    │
│                           │                                       │
│  ┌─────────────┐  ┌──────┴──────┐  ┌─────────────┐            │
│  │    NATS     │  │    Redis    │  │  Portainer  │            │
│  │  Messaging  │  │   Caching   │  │  Mgmt UI    │            │
│  │  Port 4222  │  │  Port 6379  │  │  Port 9000  │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                   │
│                    Azure VM (Standard_B2ms)                      │
│                    Ubuntu 22.04 LTS - 8GB RAM                    │
└─────────────────────────────────────────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
            ┌───────▼────────┐  ┌──────▼──────┐
            │  Tailscale VPN │  │   Public IP │
            │  onemindos.ts  │  │ 20.121.38.186│
            └────────────────┘  └─────────────┘
```

---

## Azure Resources

All resources are in the **onemind-production** resource group in **eastus** region.

### 1. Virtual Machine (onemind-vm)

**Purpose**: Primary compute resource running all OneMind OS services via Docker. **Performance-optimized for workflows and concurrent agent operations.**

**Current Specifications** (needs resize):
- **Size**: Standard_B2s (2 vCPU, 4GB RAM) - **TEMPORARY**
- **Required**: **8GB RAM minimum** for workflow performance
- **OS**: Ubuntu 22.04.3 LTS (Jammy Jellyfish)
- **Disk**: 30GB Premium SSD
- **Public IP**: 20.121.38.186 (static)
- **Private IP**: 10.0.0.4 (DHCP)
- **Admin User**: zeus
- **SSH Key**: ED25519 (~/.ssh/id_ed25519)

**Recommended VM Sizes for Production Workflows**:

| VM Size | vCPU | RAM | Cost/Month | Best For |
|---------|------|-----|------------|----------|
| **Standard_D2s_v3** ⭐ | 2 | 8GB | ~$70 | **Recommended** - Sustained performance for moderate concurrent workflows |
| **Standard_D2as_v5** | 2 | 8GB | ~$77 | AMD CPUs, excellent network (12.5 Gbps) - Best for network-intensive workflows |
| **Standard_D4s_v3** | 4 | 16GB | ~$140 | Heavy concurrent workloads, many agents running simultaneously |

**Why NOT Standard_B2ms (Burstable)?**
- Burstable VMs use CPU credits - bad for sustained workloads
- Workflows need **consistent CPU performance**
- Concurrent work will exhaust credits, causing throttling
- B-series designed for dev/test, not production workflows

**Why Ubuntu 22.04 instead of 24.04?**
- **Stability**: 22.04 is the current LTS with 5 years support (until April 2027)
- **Docker Compatibility**: Better tested with Docker and Azure provisioning
- **Production Ready**: 24.04 LTS was released April 2024, still maturing
- **Package Availability**: Some dependencies may not have 24.04 packages yet
- **Migration Path**: Can upgrade to 24.04 via `do-release-upgrade` when stable

**Scaling Plan** (Phased Approach):

**Phase 1 - NOW** (Immediate Production Readiness):
- Resize to **Standard_D2s_v3** (2 vCPU, 8GB RAM, ~$76/month)
- Provides sustained performance for moderate workflows
- Eliminates throttling from current B2s
- **This is the minimum for production workflows**

**Phase 2 - FUTURE** (Heavy Concurrent Workloads):
- Upgrade to **Standard_D4s_v3** (4 vCPU, 16GB RAM, ~$146/month)
- When: Ready for many agents running simultaneously
- 2x compute power + 2x memory
- Handles heavy concurrent workflows at scale
- **This is the goal configuration**

**OS Strategy**: Stay on 22.04 for now, upgrade to 24.04 in Q3 2025 after 1 year of stability

**Managed By**: Azure Resource Manager (ARM)

**Tags**:
- Project: OneMind OS
- Environment: Production
- Managed-By: Zeus

**Access**:
```bash
# SSH access
ssh zeus@20.121.38.186

# Via Tailscale
ssh zeus@onemindos
```

**Resize Instructions**: See [infra/deployment/azure/resize-vm.sh](./deployment/azure/resize-vm.sh)

**When to Upgrade from D2s_v3 → D4s_v3**:
- CPU consistently above 70% (check with `htop` or Azure Monitor)
- Memory usage above 6GB consistently
- Workflows taking too long to complete
- Running 5+ agents concurrently regularly
- Planning to implement complex multi-agent teams
- Ready to invest in heavy concurrent operations

---

### 2. Virtual Network (onemind-vmVNET)

**Purpose**: Private network for VM resources, enables communication between Azure resources.

**Specifications**:
- **Address Space**: 10.0.0.0/16 (65,536 IP addresses)
- **Subnet**: default (10.0.0.0/24) - 256 IPs
- **DNS Servers**: Azure-provided DNS
- **Location**: eastus

**Why This Matters**:
- Provides network isolation for security
- Enables internal communication without public internet
- Required for future scaling (adding more VMs, load balancers, etc.)
- Supports Azure Private Link and Service Endpoints

**Current Usage**: Single VM with one NIC attached

**Future Use Cases**:
- Add database in separate subnet
- Add load balancer for high availability
- Connect to Azure services (Key Vault, Storage) privately

---

### 3. Network Interface Card (onemind-vmVMNic)

**Purpose**: Connects the VM to the virtual network, handles all network traffic.

**Specifications**:
- **Private IP**: 10.0.0.4 (dynamic assignment)
- **Public IP**: Associated with onemind-vmPublicIP
- **Network Security Group**: onemind-vmNSG (firewall rules)
- **Accelerated Networking**: Disabled (not available on B-series)

**What It Does**:
- Routes inbound traffic from Public IP to VM
- Routes outbound traffic from VM to internet
- Applies NSG firewall rules
- Handles DNS resolution

**Can You Modify It?**: Generally no - managed by Azure. Changes via VM settings.

---

### 4. Public IP Address (onemind-vmPublicIP)

**Purpose**: Static public IP for accessing OneMind OS from the internet.

**Specifications**:
- **IP Address**: 20.121.38.186
- **Type**: Static (doesn't change on VM restart)
- **SKU**: Standard
- **DNS Name**: Not configured (using Tailscale DNS instead)
- **Cost**: ~$4/month

**Why Static?**:
- DNS records don't break on VM restart
- Firewall rules remain valid
- Consistent access for users and services

**Access Points**:
- Frontend: http://20.121.38.186:3000
- Backend API: http://20.121.38.186:8080
- Portainer: http://20.121.38.186:9000

**Tailscale DNS**: `onemindos.ts.net` or `onemindos` (if MagicDNS enabled)

---

### 5. Network Security Group (onemind-vmNSG)

**Purpose**: Firewall controlling what traffic can reach the VM.

**Inbound Rules** (Priority order, lower = higher priority):

| Priority | Port | Protocol | Source | Purpose |
|----------|------|----------|--------|---------|
| 1000 | 22 | TCP | Any | SSH (default) |
| 1022 | 22 | * | Any | SSH (custom) |
| 1100 | 80 | * | Any | HTTP |
| 1200 | 443 | * | Any | HTTPS |
| 1300 | 8080 | * | Any | Backend API |

**Outbound Rules**: Allow all (default Azure behavior)

**Security Considerations**:
- SSH exposed to internet (consider Tailscale-only access)
- No rate limiting (Azure DDoS Protection Basic is automatic)
- Consider restricting SSH to specific IPs or Tailscale only

**Recommended Improvements**:
```bash
# Allow SSH only from Tailscale network
az network nsg rule update --resource-group onemind-production \
  --nsg-name onemind-vmNSG --name default-allow-ssh \
  --source-address-prefixes 100.64.0.0/10

# Add rate limiting via Azure Firewall (requires Premium tier)
```

---

### 6. OS Disk (onemind-vm_OsDisk_1_...)

**Purpose**: Boot disk containing the operating system and all data.

**Specifications**:
- **Size**: 30GB (Standard_B2ms default)
- **Type**: Premium SSD (P10)
- **Cost**: ~$2.30/month
- **IOPS**: 500 IOPS, 100 MB/s throughput
- **Encryption**: Azure SSE enabled by default

**What's Stored Here**:
- Ubuntu 22.04 OS (~5GB)
- Docker images (~10GB)
- Database data (PostgreSQL/TimescaleDB)
- Application logs
- NATS data

**Monitoring Disk Usage**:
```bash
ssh zeus@20.121.38.186 "df -h /"
```

**Resize Disk** (if needed):
```bash
# Stop VM
az vm deallocate --resource-group onemind-production --name onemind-vm

# Resize disk to 100GB
az disk update --resource-group onemind-production \
  --name onemind-vm_OsDisk_1_<id> --size-gb 100

# Start VM
az vm start --resource-group onemind-production --name onemind-vm

# Expand filesystem
ssh zeus@20.121.38.186 "sudo resize2fs /dev/sda1"
```

---

### 7. Network Watcher (NetworkWatcher_eastus)

**Purpose**: Azure monitoring and diagnostic service for network troubleshooting.

**Features**:
- **IP Flow Verify**: Test if traffic is allowed/denied by NSG
- **Next Hop**: Diagnose routing issues
- **Packet Capture**: Capture network traffic for debugging
- **Connection Monitor**: Monitor endpoint connectivity
- **NSG Flow Logs**: Detailed network traffic logs

**Cost**: Free (included in Azure subscription)

**When To Use**:
- Debugging connectivity issues
- Investigating security incidents
- Monitoring network performance
- Troubleshooting NSG rule issues

**Example Usage**:
```bash
# Test if SSH is allowed from your IP
az network watcher test-ip-flow \
  --resource-group onemind-production \
  --vm onemind-vm \
  --direction Inbound \
  --protocol TCP \
  --local 10.0.0.4:22 \
  --remote <your-ip>:12345
```

---

## Networking & DNS

### Tailscale Integration

**Purpose**: Secure VPN and private DNS for OneMind OS infrastructure.

**Setup**:
1. **Install Tailscale on Azure VM**:
```bash
ssh zeus@20.121.38.186
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up --advertise-tags=tag:server
```

2. **Set Hostname**:
```bash
sudo tailscale set --hostname onemindos
```

3. **Enable MagicDNS** (in Tailscale admin console):
- Go to DNS settings
- Enable MagicDNS
- VM accessible at `onemindos` or `onemindos.ts.net`

**Benefits**:
- Access VM without exposing SSH to internet
- Private DNS (no need for public DNS records)
- Encrypted VPN tunnel
- Multi-device access (laptop, phone, tablet)
- Share access with team members

**Access via Tailscale**:
```bash
# SSH
ssh zeus@onemindos

# Frontend
http://onemindos:3000

# Backend
http://onemindos:8080
```

**Security**: Consider disabling public SSH (port 22) and using Tailscale exclusively.

---

### DNS Configuration

**Public IP DNS**: Not configured (using Tailscale instead)

**Tailscale MagicDNS**:
- Hostname: `onemindos`
- FQDN: `onemindos.ts.net`
- Resolves to Tailscale IP (100.x.x.x)

**Custom Domain** (future):
If you want a custom domain (e.g., `api.onemindos.com`):

1. Register domain (Cloudflare, Namecheap, etc.)
2. Add A record pointing to 20.121.38.186
3. Add HTTPS with Let's Encrypt:
```bash
# Install Caddy (auto HTTPS)
ssh zeus@20.121.38.186
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install caddy

# Configure Caddy
sudo nano /etc/caddy/Caddyfile
```

Caddyfile:
```
api.onemindos.com {
    reverse_proxy localhost:8080
}

app.onemindos.com {
    reverse_proxy localhost:3000
}
```

---

## CI/CD Pipeline

**Status**: ✅ Working

### GitHub Actions Workflow

**File**: `.github/workflows/deploy.yml` (to be created if not exists)

**Triggers**:
- Push to `main` branch
- Manual workflow dispatch

**Pipeline Steps**:

1. **Build Phase**:
   - Checkout code
   - Build Docker images (agno backend, flutter-web frontend)
   - Run tests (pytest for backend, flutter test for frontend)
   - Push images to registry (GitHub Container Registry or Docker Hub)

2. **Deploy Phase**:
   - SSH to Azure VM via Tailscale
   - Pull new images
   - Run database migrations
   - Restart services with `docker compose up -d`
   - Run smoke tests (health checks)

**Example Workflow**:

```yaml
name: Deploy to Azure

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Install Tailscale
        uses: tailscale/github-action@v2
        with:
          authkey: ${{ secrets.TAILSCALE_AUTHKEY }}

      - name: Build Docker images
        run: |
          docker build -t onemindos-agno:latest -f infra/containers/agno/Dockerfile .
          docker build -t onemindos-flutter:latest -f frontend/Dockerfile ./frontend

      - name: Deploy to Azure VM
        run: |
          # Transfer images
          docker save onemindos-agno:latest onemindos-flutter:latest | ssh zeus@onemindos 'docker load'

          # Restart services
          ssh zeus@onemindos 'cd ~/OneMind-OS && docker compose up -d'

          # Health check
          sleep 10
          curl -f http://onemindos:8080/health || exit 1
```

**Secrets Required** (in GitHub repo settings):
- `TAILSCALE_AUTHKEY`: Tailscale authentication key
- `AZURE_SSH_KEY`: SSH private key (if not using Tailscale)

**Future Enhancements**:
- Blue-green deployments (zero downtime)
- Automated rollback on failure
- Slack/Discord notifications
- Performance testing in pipeline

---

## VM Specifications

### Current Configuration (After Resize)

**Recommended VM Size**: Standard_D2s_v3 ⭐
- **vCPUs**: 2 (Intel Xeon Platinum 8272CL)
- **RAM**: 8GB DDR4
- **Temporary Storage**: 16GB (local SSD)
- **Max Data Disks**: 4
- **Max NICs**: 2
- **Network Bandwidth**: 2,000 Mbps (moderate)
- **Sustained Performance**: Yes (no bursting/throttling)
- **CPU Credits**: Not applicable (full performance always available)

**Alternative Option**: Standard_D2as_v5
- **vCPUs**: 2 (AMD EPYC 7763v)
- **RAM**: 8GB DDR4
- **Temporary Storage**: None (remote only)
- **Network Bandwidth**: 12,500 Mbps (excellent) - **Best for network-intensive workflows**

**For Heavy Workloads**: Standard_D4s_v3
- **vCPUs**: 4 (double the compute power)
- **RAM**: 16GB (double the memory)
- **Cost**: ~$140/month

**OS**: Ubuntu 22.04.3 LTS (Jammy Jellyfish)
- Kernel: 5.15 (HWE kernel available)
- Docker: 24.x
- Docker Compose: 2.x

**Disk**:
- **OS Disk**: 30GB Premium SSD (P10)
- **Read/Write**: 500 IOPS, 100 MB/s

**Networking**:
- **Public IP**: 20.121.38.186 (static)
- **Private IP**: 10.0.0.4 (dynamic)
- **Tailscale IP**: 100.x.x.x (assigned by Tailscale)

### Resource Monitoring

**Check CPU/Memory/Disk**:
```bash
ssh zeus@onemindos "htop"  # Interactive (press F10 to quit)
ssh zeus@onemindos "docker stats"  # Container resource usage
ssh zeus@onemindos "df -h /"  # Disk usage
```

**Azure Metrics** (in Azure Portal):
- CPU percentage
- Network in/out
- Disk read/write IOPS
- Memory usage (requires Azure Monitor agent)

---

## Cost Analysis

### Monthly Costs (USD)

**Performance-Oriented Configuration**:

| Resource | Size/Type | Cost/Month | Notes |
|----------|-----------|------------|-------|
| VM (Compute) | Standard_D2s_v3 | ~$70 | 2 vCPU, 8GB RAM, sustained performance ⭐ |
| Public IP | Static Standard | ~$4 | Required for internet access |
| OS Disk | 30GB Premium SSD | ~$2.30 | P10 tier |
| Network Egress | First 100GB free | ~$0-10 | $0.087/GB after 100GB |
| Tailscale | Free tier | $0 | Up to 100 devices |
| **Total** | | **~$76-86/month** | **Recommended for workflows** |

**Alternative Configurations**:

| Configuration | VM Size | Cost/Month | Use Case |
|---------------|---------|------------|----------|
| **Network-Intensive** | Standard_D2as_v5 | ~$83-93 | Better network (12.5 Gbps), AMD CPUs |
| **Heavy Workloads** | Standard_D4s_v3 | ~$146-156 | 4 vCPU, 16GB RAM, many concurrent agents |
| Budget (NOT RECOMMENDED) | Standard_B2ms | ~$66-76 | Burstable - will throttle under load |

**Cost Optimization**:

1. **Reserved Instances** (1-year commitment):
   - Save 30%: ~$49/month for Standard_D2s_v3 (~$588/year prepaid)
   - Save 30%: ~$54/month for Standard_D2as_v5 (~$648/year prepaid)
   - **Best ROI** for production workloads

2. **Reserved Instances** (3-year commitment):
   - Save 50%: ~$35/month for Standard_D2s_v3 (~$1,260 prepaid)
   - **Maximum savings** if committed long-term

3. **Auto-shutdown** (for dev/staging only):
   - Schedule VM shutdown during off-hours
   - Save ~50% if only running 12 hours/day
   - **NOT RECOMMENDED** for production workflows

4. **Spot Instances** (70% cheaper):
   - Risk: Azure can evict VM with 30 seconds notice
   - **NEVER use for production**

**AWS Lightsail Comparison**:
- Lightsail 4GB: $24/month (burstable)
- Lightsail 8GB: $48/month (burstable)
- **Azure Standard_D2s_v3**: ~$76-86/month (+$28-38/month vs Lightsail 8GB)
- **Benefit**: Sustained performance, no throttling, better scaling, enterprise support

**Value Proposition**:
- **Performance**: Standard_D2s_v3 has consistent CPU (no bursting/throttling)
- **Reliability**: Azure SLA 99.9% uptime
- **Workflows**: Critical for concurrent agent operations
- **Cost**: +60% vs Lightsail, but 2x-3x better performance under load

---

## Monitoring & Observability

### Built-in Services

1. **Netdata** (Port 19999):
   - Real-time system monitoring
   - Access: http://20.121.38.186:19999
   - Metrics: CPU, RAM, Disk, Network, Docker containers

2. **Portainer** (Port 9000):
   - Docker container management UI
   - Access: http://20.121.38.186:9000
   - View logs, restart containers, manage images

3. **Azure Monitor**:
   - Built-in Azure metrics (CPU, Network, Disk)
   - Access: Azure Portal > Virtual Machines > onemind-vm > Monitoring

### Application Monitoring

**Backend Logs**:
```bash
ssh zeus@onemindos "cd ~/OneMind-OS && docker compose logs -f agno"
```

**Database Logs**:
```bash
ssh zeus@onemindos "cd ~/OneMind-OS && docker compose logs -f postgres"
```

**All Services**:
```bash
ssh zeus@onemindos "cd ~/OneMind-OS && docker compose logs -f"
```

### Recommended Additions

1. **Arize Phoenix** (already in requirements.txt):
   - LLM tracing and observability
   - Track agent runs, token usage, costs

2. **Sentry** (error tracking):
   ```bash
   # Add to backend requirements.txt
   sentry-sdk[fastapi]>=2.0.0
   ```

3. **Grafana + Prometheus** (metrics dashboard):
   - More advanced than Netdata
   - Custom dashboards
   - Alerting

---

## Security

### Current Security Measures

1. **Network Security**:
   - NSG firewall (only ports 22, 80, 443, 8080 open)
   - SSH key authentication (no password)
   - Tailscale VPN (encrypted tunnel)

2. **VM Security**:
   - Ubuntu security updates (automatic)
   - Docker rootless mode (not enabled - consider)
   - Non-root user (zeus)

3. **Azure Security**:
   - Disk encryption (Azure SSE)
   - Azure Security Center (free tier)
   - DDoS Protection Basic (automatic)

### Recommended Improvements

1. **SSH Hardening**:
```bash
# Disable password authentication
ssh zeus@onemindos
sudo nano /etc/ssh/sshd_config

# Set these:
PasswordAuthentication no
PermitRootLogin no
PubkeyAuthentication yes

# Restart SSH
sudo systemctl restart sshd
```

2. **Firewall (UFW)**:
```bash
# Install and configure UFW
sudo apt install -y ufw
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow from 100.64.0.0/10 to any port 22  # Tailscale only
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp
sudo ufw enable
```

3. **Fail2Ban** (brute force protection):
```bash
sudo apt install -y fail2ban
sudo systemctl enable fail2ban
```

4. **Security Updates**:
```bash
# Enable automatic security updates
sudo apt install -y unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

5. **Environment Secrets**:
- Move .env to Azure Key Vault
- Use Managed Identity for AWS credentials
- Never commit .env to Git

---

## Backup & Disaster Recovery

### Current Backup Strategy

**Database Backups** (manual):
```bash
# Run on VM
ssh zeus@onemindos
cd ~/OneMind-OS
docker exec onemind-timescaledb pg_dump -U onemind -d onemind -F c > backup_$(date +%Y%m%d).dump
```

**Full VM Snapshots** (manual):
```bash
# Stop VM and create snapshot
az vm deallocate --resource-group onemind-production --name onemind-vm
az snapshot create --resource-group onemind-production \
  --name onemind-snapshot-$(date +%Y%m%d) \
  --source onemind-vm_OsDisk_1_<id>
az vm start --resource-group onemind-production --name onemind-vm
```

### Recommended Backup Strategy

1. **Automated Database Backups** (daily):
```bash
# Add to VM crontab
0 2 * * * cd ~/OneMind-OS && docker exec onemind-timescaledb pg_dump -U onemind -d onemind -F c > /backup/onemind_$(date +\%Y\%m\%d).dump

# Retention: Keep last 7 days
0 3 * * * find /backup -name "onemind_*.dump" -mtime +7 -delete
```

2. **Azure Backup** (VM-level):
   - Enable Azure Backup in Azure Portal
   - Daily snapshots with 30-day retention
   - Cost: ~$10-15/month for 30GB disk

3. **Off-site Backup** (to AWS S3):
```bash
# Install AWS CLI
sudo apt install -y awscli

# Upload to S3
aws s3 cp /backup/onemind_$(date +%Y%m%d).dump s3://onemindos-backups/
```

### Disaster Recovery Plan

**Scenario 1: VM Failure**
1. Provision new VM from snapshot
2. Update DNS (Tailscale hostname or public IP)
3. Verify services (should start automatically via docker compose)
4. **RTO**: 30 minutes
5. **RPO**: Last snapshot (max 24 hours)

**Scenario 2: Accidental Data Deletion**
1. Restore database from latest backup
2. **RPO**: Last backup (max 24 hours)

**Scenario 3: Region Failure**
1. Provision VM in different region (e.g., westus2)
2. Restore from off-site backup (S3)
3. Update DNS
4. **RTO**: 2-4 hours
5. **RPO**: Last off-site backup

---

## Quick Reference

### Essential Commands

**SSH to VM**:
```bash
ssh zeus@20.121.38.186
# or via Tailscale
ssh zeus@onemindos
```

**Check Services**:
```bash
ssh zeus@onemindos "cd ~/OneMind-OS && docker compose ps"
```

**View Logs**:
```bash
ssh zeus@onemindos "cd ~/OneMind-OS && docker compose logs -f agno"
```

**Restart Services**:
```bash
ssh zeus@onemindos "cd ~/OneMind-OS && docker compose restart"
```

**Database Backup**:
```bash
ssh zeus@onemindos "docker exec onemind-timescaledb pg_dump -U onemind -d onemind -F c > /tmp/backup.dump"
```

**VM Operations**:
```bash
# Stop VM
az vm stop --resource-group onemind-production --name onemind-vm

# Start VM
az vm start --resource-group onemind-production --name onemind-vm

# Restart VM
az vm restart --resource-group onemind-production --name onemind-vm

# Deallocate (stop billing compute)
az vm deallocate --resource-group onemind-production --name onemind-vm
```

### Access URLs

| Service | Public URL | Tailscale URL |
|---------|-----------|---------------|
| Frontend | http://20.121.38.186:3000 | http://onemindos:3000 |
| Backend API | http://20.121.38.186:8080 | http://onemindos:8080 |
| Portainer | http://20.121.38.186:9000 | http://onemindos:9000 |
| Netdata | http://20.121.38.186:19999 | http://onemindos:19999 |

---

## Change Log

| Date | Change | Author |
|------|--------|--------|
| 2025-01-27 | Initial Azure migration from AWS Lightsail | Zeus |
| 2025-01-27 | VM size requirement updated to 8GB minimum | Zeus |
| 2025-01-27 | Added Tailscale integration for DNS | Zeus |

---

## Support & Maintenance

**Maintained By**: Zeus (@Zeus-Delacruz)

**Related Documentation**:
- [Azure Migration Guide](./deployment/azure/README.md)
- [Migration Plan](./deployment/azure/MIGRATION_PLAN.md)
- [Verification Script](./deployment/azure/verify.sh)
- [Rollback Procedure](./deployment/azure/rollback.sh)

**For Issues**:
- Check logs: `docker compose logs`
- Verify services: `./verify.sh`
- Review [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)

---

**Last Updated**: 2025-01-27
**Document Version**: 1.0


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
