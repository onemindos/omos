---
id: ui-02-20260111053925
title: "DEPLOYMENT ARCHITECTURE DEEP DIVE"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "DEPLOYMENT ARCHITECTURE DEEP DIVE"
---

# DEPLOYMENT ARCHITECTURE DEEP DIVE

**Last Updated:** October 24, 2025  
**Status:** Production Architecture  
**Scope:** Multi-VM deployment with Portainer orchestration and Cloudflare integration

---

## TABLE OF CONTENTS

1. [System Overview](#system-overview)
2. [Infrastructure Architecture](#infrastructure-architecture)
3. [Multi-VM Deployment Strategy](#multi-vm-deployment-strategy)
4. [Portainer Orchestration](#portainer-orchestration)
5. [Cloudflare Integration](#cloudflare-integration)
6. [Service Interconnection](#service-interconnection)
7. [Deployment Workflows](#deployment-workflows)
8. [Security Architecture](#security-architecture)
9. [Monitoring & Observability](#monitoring--observability)
10. [Disaster Recovery](#disaster-recovery)

---

## SYSTEM OVERVIEW

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    CLOUDFLARE EDGE                          │
│  ├─ DNS Resolution                                          │
│  ├─ DDoS Protection                                         │
│  ├─ SSL/TLS Termination                                     │
│  └─ Geographic Load Balancing                               │
└────────────┬────────────────────────────────────────────────┘
             │
             │ HTTPS / Encrypted Traffic
             ▼
┌─────────────────────────────────────────────────────────────┐
│              PRIMARY VM (Ubuntu 22.04)                      │
│  ├─ Portainer Server (Container Orchestration)             │
│  ├─ Docker Daemon (Primary Container Runtime)              │
│  ├─ Core Services Stack                                    │
│  │  ├─ n8n (Automation)                                    │
│  │  ├─ Dify (AI Platform)                                  │
│  │  ├─ PostgreSQL (Primary DB)                             │
│  │  └─ Redis (Cache/Session)                               │
│  ├─ Caddy Reverse Proxy (Internal)                         │
│  └─ Monitoring Stack                                       │
│     ├─ Prometheus                                          │
│     └─ Grafana                                             │
└────────────┬────────────────────────────────────────────────┘
             │
             │ Docker Network (Internal Communication)
             │
             ├──────────────────┬──────────────────┐
             ▼                  ▼                  ▼
    ┌─────────────────┐ ┌──────────────┐ ┌──────────────┐
    │  IoT Agent VM   │ │  LLM VM      │ │  Voice VM    │
    │ (Ubuntu 22.04) │ │(Ubuntu 22.04)│ │(Ubuntu 22.04)│
    │                 │ │              │ │              │
    │ ├─ Portainer   │ │├─ Portainer  │ │├─ Portainer  │
    │ │  Agent       │ ││ Agent       │ ││ Agent       │
    │ │              │ ││             │ ││             │
    │ ├─ Node-RED    │ │├─ Ollama     │ │├─ Frigate    │
    │ ├─ EMQX        │ │├─ LM Studio │ │├─ Whisper    │
    │ ├─ Home        │ │├─ LangChain │ │├─ Piper      │
    │ │  Assistant   │ ││             │ ││             │
    │ └─ Sensors     │ │└─ Embeddings │ │└─ Processing │
    │   Stack        │ │  Database    │ │  Engines     │
    └─────────────────┘ └──────────────┘ └──────────────┘
```

### Core Components

| Component | Purpose | VM Location | Status |
|-----------|---------|-------------|--------|
| **Portainer Server** | Container orchestration | Primary | Active |
| **Portainer Agents** | Remote container mgmt | All Secondary | Active |
| **n8n** | Workflow automation | Primary | Active |
| **Dify** | AI platform & RAG | Primary | Active |
| **Node-RED** | Visual automation | IoT VM | Active |
| **EMQX** | MQTT broker | IoT VM | Active |
| **Home Assistant** | Smart home hub | IoT VM | Active |
| **Ollama** | Local LLM inference | LLM VM | Active |
| **Frigate** | Video AI detection | Voice VM | Active |

---

## INFRASTRUCTURE ARCHITECTURE

### VM Configuration Matrix

#### Primary VM (Control Plane)
```yaml
Name: legacy-ai-primary
Provider: AWS Lightsail
Instance Type: 4GB RAM / 2vCPU
Storage: 60GB SSD
OS: Ubuntu 22.04 LTS
Network: Public IP + Elastic IP
Firewall: UFW (Uncomplicated Firewall)

Services:
  - Portainer CE (Control Plane)
  - Docker Engine (20.10+)
  - Caddy v2 (Reverse Proxy)
  - PostgreSQL (Primary Database)
  - Redis (Cache)
  - Prometheus (Monitoring)
  - Grafana (Dashboards)

Inbound Rules:
  - Port 22/TCP: SSH (Restricted to Known IPs)
  - Port 80/TCP: HTTP (Caddy)
  - Port 443/TCP: HTTPS (Caddy)
  - Port 9000/TCP: Portainer UI (Internal Only)
  - Port 8000/TCP: Portainer Agent (Internal)
```

#### IoT Agent VM
```yaml
Name: legacy-ai-iot
Provider: AWS Lightsail
Instance Type: 2GB RAM / 1vCPU
Storage: 40GB SSD
OS: Ubuntu 22.04 LTS
Network: Private IP (SSH tunneled)

Services:
  - Portainer Agent
  - Docker Engine
  - Node-RED
  - EMQX (MQTT Broker)
  - Home Assistant
  - Sensor Dashboard

Network:
  - Connected to Primary via Docker Overlay Network
  - MQTT: Port 1883/TCP (Internal)
  - HTTP: Port 1880/TCP (Node-RED UI)
```

#### LLM VM
```yaml
Name: legacy-ai-llm
Provider: AWS Lightsail
Instance Type: 8GB RAM / 4vCPU (GPU Available)
Storage: 80GB SSD
OS: Ubuntu 22.04 LTS
Network: Private IP (SSH tunneled)

Services:
  - Portainer Agent
  - Docker Engine
  - Ollama (LLM Inference)
  - LM Studio (Model Management)
  - LangChain Integration
  - Embeddings Database

Optimization:
  - GPU Support: CUDA/cuDNN if GPU available
  - Memory Optimization: Swap configured
  - Model Storage: Dedicated /models partition
```

#### Voice/Vision VM
```yaml
Name: legacy-ai-voice
Provider: AWS Lightsail
Instance Type: 4GB RAM / 2vCPU
Storage: 60GB SSD
OS: Ubuntu 22.04 LTS
Network: Private IP (SSH tunneled)

Services:
  - Portainer Agent
  - Docker Engine
  - Frigate (Video Intelligence)
  - Whisper (Speech Recognition)
  - Piper (Text-to-Speech)
  - Processing Engines

Network:
  - Camera streams: RTSP/RTMP
  - API: Port 5000/TCP
  - UI: Port 5001/TCP
```

### Network Architecture

#### Docker Networking

```
PRIMARY VM:
┌─────────────────────────────────────┐
│ Docker Engine                        │
│                                      │
│ ┌──────────────────────────────────┐│
│ │ Bridge Network (legacy-ai)       ││
│ │                                  ││
│ │ - n8n (8080/TCP)                 ││
│ │ - Dify (80/TCP)                  ││
│ │ - PostgreSQL (5432/TCP)          ││
│ │ - Redis (6379/TCP)               ││
│ │ - Caddy (80, 443/TCP)            ││
│ │                                  ││
│ │ DNS: Docker embedded DNS         ││
│ │ Service discovery: hostname      ││
│ └──────────────────────────────────┘│
│                                      │
│ ┌──────────────────────────────────┐│
│ │ Overlay Network (portainer-net)  ││
│ │ (For multi-host communication)   ││
│ │                                  ││
│ │ - Connects to Agent VMs          ││
│ │ - Encrypted communication        ││
│ │ - Service discovery              ││
│ └──────────────────────────────────┘│
└─────────────────────────────────────┘
```

#### Firewall Rules

```yaml
Primary VM (UFW):
  - ssh: 22/tcp (Restricted)
  - http: 80/tcp (All)
  - https: 443/tcp (All)
  - portainer_internal: 9000/tcp (10.0.0.0/8)
  - portainer_agent: 8000/tcp (10.0.0.0/8)
  - postgres_internal: 5432/tcp (10.0.0.0/8)
  - redis_internal: 6379/tcp (10.0.0.0/8)

Secondary VMs (UFW):
  - ssh: 22/tcp (Primary VM only)
  - portainer_agent: 8000/tcp (Primary VM)
  - docker_internal: All (Docker networks)
  - service_ports: 5000-6000/tcp (Docker services)
```

---

## MULTI-VM DEPLOYMENT STRATEGY

### Deployment Approach

#### 1. **Primary VM Deployment**

```bash
# Phase 1: System Preparation
ubuntu_init.sh                    # OS updates, security hardening
install_docker.sh                 # Docker Engine installation
configure_firewall.sh             # UFW setup

# Phase 2: Core Services
docker compose -f core-stack.yml up -d
  ├─ PostgreSQL (persistent storage)
  ├─ Redis (ephemeral cache)
  ├─ n8n (automation engine)
  └─ Dify (AI platform)

# Phase 3: Orchestration
install_portainer.sh              # Portainer Server setup
configure_portainer.sh            # Server configuration

# Phase 4: Reverse Proxy
install_caddy.sh                  # Caddy web server
configure_caddy.sh                # Routing & SSL rules

# Phase 5: Monitoring
docker compose -f monitoring.yml up -d
  ├─ Prometheus (metrics collection)
  └─ Grafana (visualization)
```

#### 2. **Secondary VM Deployment**

```bash
# Phase 1: System Preparation
ubuntu_init.sh                    # OS setup
install_docker.sh                 # Docker Engine
install_portainer_agent.sh        # Agent setup

# Phase 2: Network Configuration
join_overlay_network.sh           # Connect to primary
verify_connectivity.sh             # Test cross-VM communication

# Phase 3: Service Deployment
# Services deployed via Portainer Web UI or API
docker pull registry/image        # Pull required images
docker compose up -d              # Start service stack
```

### Deployment Orchestration

#### Manual Deployment Flow

```
┌─────────────────────────────┐
│ Pre-Deployment Checklist    │
│ ├─ OS Requirements (22.04)  │
│ ├─ Network Access           │
│ ├─ Storage Available        │
│ ├─ Resources (RAM/CPU)      │
│ └─ Security Groups          │
└────────────┬────────────────┘
             ▼
┌─────────────────────────────┐
│ Primary VM Setup            │
│ ├─ Initialize OS            │
│ ├─ Install Docker           │
│ ├─ Deploy Core Services     │
│ ├─ Install Portainer Server │
│ └─ Configure Caddy          │
└────────────┬────────────────┘
             ▼
┌─────────────────────────────┐
│ Secondary VMs Setup         │
│ ├─ Initialize OS            │
│ ├─ Install Docker           │
│ ├─ Install Agent            │
│ ├─ Join Overlay Network     │
│ └─ Deploy Services          │
└────────────┬────────────────┘
             ▼
┌─────────────────────────────┐
│ Verification & Testing      │
│ ├─ Connectivity Tests       │
│ ├─ Service Health Checks    │
│ ├─ Load Testing             │
│ └─ Security Audit           │
└─────────────────────────────┘
```

### Service Deployment Patterns

#### Pattern 1: Stateful Core Services (Primary VM)
```yaml
# These MUST run on Primary VM
services:
  postgresql:
    constraints:
      node.hostname == primary
    volumes:
      - /data/postgres:/var/lib/postgresql/data
    
  redis:
    constraints:
      node.hostname == primary
    volumes:
      - /data/redis:/data
```

#### Pattern 2: Distributed Agents (All Secondary VMs)
```yaml
# These run on Secondary VMs via Agent
services:
  node-red:
    deploy:
      placement:
        constraints:
          - node.labels.type == iot
```

#### Pattern 3: Flexible Services (Can run anywhere)
```yaml
# These can be deployed based on resource availability
services:
  n8n:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

---

## PORTAINER ORCHESTRATION

### Portainer Architecture

#### Server Installation (Primary VM)

```bash
# Installation Method: Docker Container
docker run -d \
  --name portainer \
  --restart always \
  -p 9000:9000 \
  -p 8000:8000 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v portainer_data:/data \
  portainer/portainer-ce:latest
```

#### Agent Installation (Secondary VMs)

```bash
# Installation Method: Docker Container
docker run -d \
  --name portainer_agent \
  --restart always \
  -p 8000:8000 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/lib/docker/volumes:/var/lib/docker/volumes \
  portainer/agent:latest
```

### Portainer Workflow

#### Environment Management

```
Portainer Server
│
├─ Local Environment (Primary VM)
│  ├─ Docker Containers
│  ├─ Networks
│  ├─ Volumes
│  ├─ Images
│  └─ Registries
│
└─ Remote Environments (Secondary VMs)
   ├─ IoT Agent VM
   │  └─ Local Docker environment
   │
   ├─ LLM VM
   │  └─ Local Docker environment
   │
   └─ Voice/Vision VM
      └─ Local Docker environment
```

#### Stack Deployment via Portainer

```yaml
# Stack: Core Services (Primary)
# Deployment: Via Portainer UI or API
# File: stacks/core-stack.yml

version: '3.8'

services:
  n8n:
    image: n8nio/n8n:latest
    ports:
      - "5678:5678"
    environment:
      - N8N_HOST=n8n.example.com
      - DB_TYPE=postgresdb
      - DB_POSTGRESDB_HOST=postgresql
      - DB_POSTGRESDB_USER=n8n
      - DB_POSTGRESDB_PASSWORD=${N8N_DB_PASSWORD}
      - DB_POSTGRESDB_DATABASE=n8n
    depends_on:
      - postgresql
    networks:
      - legacy-ai

  dify:
    image: langgenius/dify-api:latest
    ports:
      - "5001:5001"
    environment:
      - DATABASE_URL=postgresql://dify:${DIFY_DB_PASSWORD}@postgresql:5432/dify
      - REDIS_URL=redis://redis:6379/0
    depends_on:
      - postgresql
      - redis
    networks:
      - legacy-ai

  postgresql:
    image: postgres:15-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=admin
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    networks:
      - legacy-ai

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    networks:
      - legacy-ai

volumes:
  postgres_data:
  redis_data:

networks:
  legacy-ai:
    driver: bridge
```

### Portainer Management Features

#### Access Control
```
Role-Based Access Control (RBAC):
├─ Administrator
│  └─ Full system access
│
├─ Operator
│  ├─ Container management
│  ├─ Image management
│  └─ Network/volume management
│
├─ Standard User
│  ├─ View resources
│  ├─ View logs
│  └─ Limited operations
│
└─ Read-Only User
   └─ View-only access
```

#### Monitoring Dashboard
```
Metrics Displayed:
├─ Container Status
├─ CPU/Memory Usage
├─ Network I/O
├─ Storage Usage
├─ Resource Allocation
└─ Performance Trends
```

---

## CLOUDFLARE INTEGRATION

### DNS Configuration

#### Domain Setup

```
Domain: example.com
Registrar: (Your current registrar)
Nameservers: Updated to Cloudflare
  - ns1.cloudflare.com
  - ns2.cloudflare.com
```

#### DNS Records

```yaml
# Root Domain
Type: A Record
Name: @
Value: <Primary VM Elastic IP>
Proxy: Cloudflare (Orange Cloud)
TTL: Auto

# Subdomains (via CNAME or A Record)
Type: CNAME
Name: n8n
Value: example.com
Proxy: Cloudflare

Type: CNAME
Name: dify
Value: example.com
Proxy: Cloudflare

Type: CNAME
Name: portainer
Value: example.com
Proxy: Cloudflare

Type: A Record
Name: api
Value: <LLM VM IP>
Proxy: DNS Only (Grey Cloud)
```

### SSL/TLS Configuration

#### Certificate Management

```
Cloudflare SSL Mode: Full (Strict)
├─ Client to Cloudflare: HTTPS
├─ Cloudflare to Origin: HTTPS (Valid Certificate)
└─ Certificate: Cloudflare Managed SSL

Caddy Configuration (On Primary):
├─ Auto HTTPS enabled
├─ DNS challenge via Cloudflare API
├─ Certificate storage: /etc/caddy/certs/
└─ Auto-renewal: Built-in
```

#### Caddy Configuration File

```caddyfile
# Caddyfile for SSL/TLS with Cloudflare

# Global settings
{
    email admin@example.com
    acme_dns cloudflare {
        token {CLOUDFLARE_API_TOKEN}
    }
}

# Primary domain
example.com {
    # Reverse proxy to n8n
    handle /automation* {
        reverse_proxy n8n:5678
    }
    
    # Reverse proxy to Dify
    handle /ai* {
        reverse_proxy dify:5001
    }
    
    # Reverse proxy to other services
    handle /api* {
        reverse_proxy portainer:9000
    }
    
    # Security headers
    header X-Content-Type-Options nosniff
    header X-Frame-Options SAMEORIGIN
    header X-XSS-Protection "1; mode=block"
    header Referrer-Policy strict-origin-when-cross-origin
}

# Additional domain
api.example.com {
    reverse_proxy 10.0.0.5:8000
}
```

### DDoS Protection & WAF

#### Cloudflare DDoS Settings

```yaml
DDoS Protection:
  - Challenge Passage: 1h
  - Sensitivity: High
  - Bot Fight Mode: Enabled
  - Super Bot Fight Mode: Professional
    ├─ Block known bots
    ├─ Challenge suspicious bots
    └─ Block verified bots

WAF Rules:
  - OWASP Core Ruleset: Enabled (High)
  - Rate Limiting:
    ├─ Threshold: 100 requests/minute
    └─ Action: Challenge
  
  - Custom Rules:
    ├─ Block SQL injection attempts
    ├─ Block XSS attempts
    └─ Allow specific IPs (Admin)
```

### Geographic Load Balancing

#### Origin Health Checks

```yaml
Primary VM Health Check:
  - Endpoint: example.com/health
  - Interval: 30 seconds
  - Timeout: 5 seconds
  - Expected Code: 200
  - Consecutive Failures: 2
  
Failover Action:
  - If Primary fails: Route to Secondary
  - Alert: Send notification to admin
  - Auto-recovery: Monitoring enabled
```

---

## SERVICE INTERCONNECTION

### Communication Patterns

#### Internal Communication (Within Primary VM)

```
n8n ←→ PostgreSQL
├─ Method: Direct TCP (5432)
├─ Network: Docker Bridge (legacy-ai)
├─ Auth: Username/Password
└─ Encryption: Optional (Docker network isolation)

n8n ←→ Redis
├─ Method: Direct TCP (6379)
├─ Network: Docker Bridge
├─ Auth: None (Local)
└─ Encryption: None (Local network)

Dify ←→ PostgreSQL
├─ Method: Direct TCP (5432)
├─ Network: Docker Bridge
├─ Auth: Username/Password
└─ Encryption: Optional

Dify ←→ Redis
├─ Method: Direct TCP (6379)
├─ Network: Docker Bridge
├─ Auth: None
└─ Encryption: None
```

#### Cross-VM Communication (Via Docker Overlay)

```
Workflow:
Primary VM (n8n) → Overlay Network → LLM VM (Ollama)

Connection Details:
├─ Network: portainer-net (overlay)
├─ Service Discovery: DNS names
├─ Routing: Embedded routing mesh
├─ Encryption: IPSEC (configurable)
└─ MTU: 1450 bytes (GRE overhead)

Example API Call:
curl http://ollama:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "mistral",
    "prompt": "Why is the sky blue?"
  }'
```

#### External Communication (via Cloudflare)

```
Client → Cloudflare Edge → Caddy Reverse Proxy → Services

Flow:
1. Client resolves example.com
2. Cloudflare resolves to Primary VM IP
3. Cloudflare proxy forwards request
4. Caddy receives request on :443
5. Caddy routes based on path/hostname
6. Service responds through Caddy
7. Response sent through Cloudflare
8. Client receives response
```

### Service Mesh Architecture

#### Without Service Mesh (Current)

```
Direct Communication:
┌──────────┐
│ Service A│
└─────┬────┘
      │
      ├─→ Service B (Docker DNS)
      ├─→ Service C (Docker DNS)
      └─→ Service D (Docker DNS)

Limitations:
- No built-in retries
- No circuit breaking
- Limited traffic shaping
- Manual configuration
```

#### With Linkerd (Future Enhancement)

```
Service Mesh:
┌──────────┐         ┌────────┐
│ Service A│─────────│Linkerd │
└──────────┘         │Proxy   │
                     └────┬───┘
                          │
              ┌───────────┼───────────┐
              ▼           ▼           ▼
         Service B   Service C   Service D
         (Proxy)     (Proxy)     (Proxy)

Benefits:
- Automatic retries
- Circuit breaking
- Load balancing
- Traffic shaping
- Observability
```

---

## DEPLOYMENT WORKFLOWS

### Initial Deployment Workflow

```mermaid
graph TD
    A[Pre-Deployment] --> B{Check Requirements}
    B -->|Pass| C[Create Primary VM]
    B -->|Fail| Z[Exit & Fix Issues]
    
    C --> D[OS Configuration]
    D --> E[Install Docker]
    E --> F[Install Portainer Server]
    
    F --> G[Create Secondary VMs]
    G --> H[Install Portainer Agents]
    H --> I[Configure Overlay Network]
    
    I --> J[Deploy Core Services]
    J --> K[Deploy Secondary Services]
    K --> L[Configure Cloudflare]
    
    L --> M[SSL Certificate Setup]
    M --> N[Configure Caddy]
    
    N --> O[Monitoring Setup]
    O --> P[Testing & Validation]
    
    P --> Q{All Tests Pass?}
    Q -->|Yes| R[Production Ready]
    Q -->|No| S[Debug & Fix]
    S --> P
```

### Scaling Workflow

```
Current State:
├─ Primary VM: n8n, Dify, PostgreSQL, Redis
├─ IoT VM: Node-RED, EMQX, Home Assistant
├─ LLM VM: Ollama, LangChain
└─ Voice VM: Frigate, Whisper, Piper

Scaling Up:
├─ Add API Gateway VM
│  └─ Kong or Traefik for advanced routing
│
├─ Add Database Replication
│  └─ PostgreSQL Replica on Secondary VM
│
├─ Add Cache Layer
│  └─ Redis Sentinel for HA
│
└─ Add Load Balancer
   └─ HAProxy or AWS Network Load Balancer
```

### Backup & Recovery Workflow

```
Backup Schedule:
├─ PostgreSQL: Daily at 02:00 UTC
├─ Redis: Every 4 hours
├─ Volumes: Daily at 03:00 UTC
├─ Configurations: Every 6 hours
└─ Full System: Weekly

Backup Storage:
├─ Local: /backups/local/ (7-day retention)
├─ S3: S3 bucket (30-day retention)
└─ Glacier: Archive (1-year retention)

Recovery Procedure:
1. Stop affected services
2. Restore from backup
3. Verify data integrity
4. Start services
5. Monitor for issues
6. Document incident
```

---

## SECURITY ARCHITECTURE

### Authentication & Authorization

#### SSH Access

```bash
# SSH Key-based Authentication
├─ Admin Key: /home/admin/.ssh/id_rsa
├─ Automation Key: /home/automation/.ssh/id_rsa
└─ Backup Key: /home/backup/.ssh/id_rsa

SSH Configuration (/etc/ssh/sshd_config):
├─ PasswordAuthentication: no
├─ PubkeyAuthentication: yes
├─ AllowUsers: admin automation backup
├─ Port: 22 (Non-standard recommended)
└─ LogLevel: VERBOSE
```

#### Service Authentication

```yaml
PostgreSQL:
  - User: admin (Full access)
  - User: n8n (Limited to n8n DB)
  - User: dify (Limited to dify DB)
  - Passwords: Strong (20+ chars, random)
  - Access: Local socket (127.0.0.1)

Redis:
  - Password: Strong (20+ chars, random)
  - ACL: Basic (users/permissions)
  - Access: Docker network only

n8n:
  - Authentication: Local + OAuth2 (optional)
  - Session: Secure (HttpOnly, Secure flags)
  - RBAC: User-based

Dify:
  - Authentication: User registration + OAuth
  - Session: JWT tokens
  - RBAC: Admin, Team Lead, Member, Viewer
```

### Network Security

#### Firewall Configuration

```bash
# UFW Rules (Primary VM)
sudo ufw allow from any to any port 22 proto tcp comment "SSH"
sudo ufw allow from any to any port 80 proto tcp comment "HTTP"
sudo ufw allow from any to any port 443 proto tcp comment "HTTPS"
sudo ufw allow from 10.0.0.0/8 to any port 9000 proto tcp comment "Portainer UI"
sudo ufw allow from 10.0.0.0/8 to any port 8000 proto tcp comment "Portainer Agent"
sudo ufw allow from 10.0.0.0/8 to any port 5432 proto tcp comment "PostgreSQL"

sudo ufw enable
sudo ufw status verbose
```

#### Network Segmentation

```
Public Network (Internet):
├─ Cloudflare (DNS, DDoS protection)
└─ Primary VM (Ports 80, 443)

Private Network (Docker):
├─ Bridge Network (legacy-ai)
│  ├─ n8n, Dify, PostgreSQL, Redis, Caddy
│  └─ Internal communication only
│
└─ Overlay Network (portainer-net)
   ├─ Primary VM
   ├─ IoT Agent VM
   ├─ LLM Agent VM
   └─ Voice/Vision Agent VM
```

### Secrets Management

#### Environment Variables

```bash
# Primary VM - .env file
POSTGRES_PASSWORD=<random_strong_password>
POSTGRES_USER=admin
N8N_DB_PASSWORD=<random_strong_password>
DIFY_DB_PASSWORD=<random_strong_password>
REDIS_PASSWORD=<random_strong_password>
CLOUDFLARE_API_TOKEN=<cloudflare_token>
JWT_SECRET=<random_jwt_secret>

# Storage: /etc/environment (restricted to root)
# Rotation: Monthly
# Backup: Encrypted (RSA-2048)
```

#### Secret Rotation Process

```
Schedule: Quarterly

Process:
1. Generate new secrets
2. Update services (zero-downtime)
3. Verify functionality
4. Document changes
5. Archive old secrets (encrypted)
6. Audit log entry

Implementation:
├─ PostgreSQL passwords: Via pg_hba.conf update
├─ n8n secrets: Via environment restart
├─ Redis password: Via CONFIG REWRITE
└─ Cloudflare token: Via UI update
```

---

## MONITORING & OBSERVABILITY

### Metrics Collection

#### Prometheus Targets

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  # Docker metrics
  - job_name: 'docker'
    static_configs:
      - targets: ['localhost:9323']
  
  # PostgreSQL metrics
  - job_name: 'postgres'
    static_configs:
      - targets: ['localhost:9187']
  
  # Redis metrics
  - job_name: 'redis'
    static_configs:
      - targets: ['localhost:9121']
  
  # n8n metrics
  - job_name: 'n8n'
    static_configs:
      - targets: ['localhost:5678']
  
  # Portainer metrics
  - job_name: 'portainer'
    static_configs:
      - targets: ['localhost:9000']
```

#### Key Metrics

```
Container Metrics:
├─ Container count (running/stopped/total)
├─ CPU usage (per container)
├─ Memory usage (per container)
├─ Network I/O (bytes in/out)
└─ Disk I/O (read/write)

Application Metrics:
├─ Request latency
├─ Error rate
├─ Throughput (requests/sec)
├─ Queue depth (n8n workflows)
└─ Database connections

Infrastructure Metrics:
├─ VM CPU usage
├─ VM Memory usage
├─ Disk space available
├─ Network bandwidth
└─ System load average
```

### Alerting Configuration

#### Alert Rules

```yaml
# prometheus_rules.yml

# High Memory Usage
- alert: HighMemoryUsage
  expr: container_memory_usage_bytes / container_spec_memory_limit_bytes > 0.8
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "High memory usage detected"
    description: "{{ $labels.container_name }} is using {{ $value }}% of memory"

# Service Down
- alert: ServiceDown
  expr: up{job=~"n8n|dify|postgres"} == 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "Service is down"
    description: "{{ $labels.job }} has been down for more than 1 minute"

# Disk Space Low
- alert: DiskSpaceLow
  expr: node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"} < 0.1
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Low disk space"
    description: "Less than 10% disk space remaining"
```

#### Notification Channels

```
Alerting Destinations:
├─ Email: alerts@example.com
├─ Slack: #infrastructure-alerts
├─ PagerDuty: Critical incidents
├─ Webhook: Custom integrations
└─ SMS: Critical-only (optional)

Severity Levels:
├─ Info: Logged only
├─ Warning: Email + Slack
├─ Critical: Email + Slack + PagerDuty + SMS
└─ Emergency: All channels + Page on-call
```

### Logging Architecture

#### Log Collection

```yaml
Logging Stack:
├─ Source: Container logs (stdout/stderr)
├─ Collection: Docker daemon → syslog
├─ Aggregation: ELK Stack (optional)
│  ├─ Elasticsearch (storage)
│  ├─ Logstash (processing)
│  └─ Kibana (visualization)
└─ Retention: 30 days on disk, 90 days in archive

Log Format (JSON):
{
  "timestamp": "2025-10-24T10:30:00Z",
  "level": "ERROR",
  "service": "n8n",
  "container_id": "abc123",
  "message": "Workflow execution failed",
  "error_code": "ERR_WORKFLOW_001",
  "request_id": "req_123456"
}
```

#### Log Levels

```
DEBUG: Development/troubleshooting
INFO: General operational information
WARNING: Potential issues
ERROR: Error conditions
CRITICAL: System-wide failures

Configuration per Service:
├─ n8n: INFO (Production)
├─ Dify: WARNING (Production)
├─ PostgreSQL: WARNING (Production)
├─ Redis: INFO (Production)
└─ Caddy: WARNING (Production)
```

---

## DISASTER RECOVERY

### Backup Strategy

#### What Gets Backed Up

```
Tier 1 - Critical (Daily, 7-day retention):
├─ PostgreSQL database
├─ Redis snapshots
├─ Portainer configuration
└─ SSL certificates

Tier 2 - Important (Weekly, 30-day retention):
├─ Docker volumes
├─ Configuration files
├─ Service definitions
└─ Monitoring configurations

Tier 3 - Archive (Monthly, 365-day retention):
├─ Full system state
├─ Logs
├─ Audit trails
└─ Documentation
```

#### Backup Automation

```bash
#!/bin/bash
# backup-critical.sh - Runs daily at 02:00 UTC

BACKUP_DIR="/backups/$(date +%Y-%m-%d)"
S3_BUCKET="legacy-ai-backups"

# PostgreSQL backup
pg_dump -h localhost -U admin -d legacy_ai > \
  $BACKUP_DIR/postgres_backup.sql.gz

# Redis backup
redis-cli BGSAVE
cp /var/lib/redis/dump.rdb $BACKUP_DIR/redis_dump.rdb

# Portainer backup
docker exec portainer tar czf - \
  /data/portainer > $BACKUP_DIR/portainer_backup.tar.gz

# Upload to S3
aws s3 sync $BACKUP_DIR/ s3://$S3_BUCKET/$(date +%Y/%m/%d)/

# Clean old backups (7-day retention)
find /backups -mtime +7 -delete
```

### Recovery Procedures

#### Scenario 1: Service Failure

```
Detection:
1. Monitoring alert triggered
2. Health check failure detected
3. Alert sent to operators

Response:
1. SSH into affected VM
2. Check service status: docker ps
3. Review logs: docker logs <service>
4. Restart service: docker restart <service>
5. Verify recovery: Check health endpoint
6. Document incident
```

#### Scenario 2: Database Corruption

```
Detection:
1. PostgreSQL error in logs
2. Failed integrity check
3. Service responds with errors

Recovery Steps:
1. Stop dependent services (n8n, Dify)
2. Restore PostgreSQL from latest backup
   - Restore timestamp: $(date -d '1 hour ago' +%Y-%m-%d_%H:%M:%S)
   - Command: pg_restore -d legacy_ai backup.sql
3. Verify data integrity
4. Start services in dependency order
5. Run validation queries
6. Monitor for issues
```

#### Scenario 3: VM Failure

```
Failure Type: Complete VM outage (primary)

Recovery Plan:
Phase 1: Immediate (0-15 minutes)
├─ Detect outage via monitoring
├─ Failover to backup (if available)
└─ Notify on-call team

Phase 2: Investigation (15-60 minutes)
├─ Check AWS Lightsail console
├─ Review system logs
├─ Determine root cause
└─ Decide: Reboot vs. Restore

Phase 3: Recovery (1-4 hours)
├─ Reboot VM (if software issue)
├─ Restore from AMI (if hardware issue)
├─ Verify services are running
├─ Check data consistency
└─ Run full system tests

Phase 4: Follow-up (Post-incident)
├─ Incident report
├─ Root cause analysis
├─ Preventive measures
└─ Update documentation
```

### RTO/RPO Targets

```
Recovery Time Objective (RTO):

Service Critical Services | RTO Target
├─ n8n                    | 30 minutes
├─ Dify                   | 30 minutes
├─ PostgreSQL             | 1 hour
├─ Redis                  | 15 minutes
└─ Portainer              | 1 hour

Recovery Point Objective (RPO):

├─ PostgreSQL             | 4 hours (6 backups/day)
├─ Redis                  | 1 hour
├─ Service Config         | 6 hours
└─ Full System Snapshot   | 24 hours
```

---

## APPENDIX: TROUBLESHOOTING REFERENCE

### Common Issues & Solutions

#### 1. Services Not Communicating

**Symptom:** "Connection refused" errors between services

```bash
# Check network connectivity
docker network inspect legacy-ai

# Test DNS resolution
docker exec <service> nslookup <target-service>

# Check firewall rules
sudo ufw status verbose

# Restart networking
docker network restart legacy-ai
```

#### 2. Database Connection Issues

**Symptom:** PostgreSQL authentication errors

```bash
# Check PostgreSQL status
docker ps | grep postgres

# View logs
docker logs <postgres-container>

# Verify credentials
docker exec <postgres-container> psql -U admin -d postgres

# Reset password
docker exec <postgres-container> \
  psql -U admin -c "ALTER USER admin WITH PASSWORD 'newpassword';"
```

#### 3. Disk Space Issues

**Symptom:** "No space left on device" errors

```bash
# Check disk usage
df -h
du -sh /var/lib/docker/

# Clean up Docker
docker system prune -a

# Clean up old backups
find /backups -mtime +30 -delete

# Remove dangling volumes
docker volume prune
```

#### 4. Certificate Issues

**Symptom:** SSL certificate errors or expiration warnings

```bash
# Check certificate
openssl x509 -in /etc/caddy/certs/cert.pem -text -noout

# Renew certificate
caddy reload

# Check Cloudflare settings
curl -I https://example.com
```

---

## REFERENCES & RESOURCES

- **Docker Official Documentation:** https://docs.docker.com
- **Portainer Documentation:** https://docs.portainer.io
- **Caddy Documentation:** https://caddyserver.com/docs
- **Cloudflare API Reference:** https://developers.cloudflare.com/api
- **PostgreSQL Documentation:** https://www.postgresql.org/docs
- **Prometheus Documentation:** https://prometheus.io/docs
- **Grafana Documentation:** https://grafana.com/docs

---

## VERSION HISTORY

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-10-24 | System | Initial comprehensive deployment guide |

---

**Last Updated:** October 24, 2025  
**Next Review:** November 24, 2025  
**Status:** PRODUCTION
