---
id: ui-02-20260111053925
title: "Current Services - Active Deployment"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Current Services - Active Deployment"
---

# Current Services - Active Deployment

**Last Updated**: October 23, 2025  
**VM**: 3.215.29.128 | **Total Containers**: 25 running | **Status**: ✅ Production

---

## 📋 Active Services by Stack

### Infrastructure (CRITICAL)

| Service | Container | Port | Access | Status |
|---------|-----------|------|--------|--------|
| **Portainer** | portainer | 9443 | https://command.zeusops.dev | ✅ Running |
| **Caddy** | caddy | 80, 443 | Internal reverse proxy | ✅ Running |

### AI Services

| Service | Container | Port | Access | Status |
|---------|-----------|------|--------|--------|
| **Letta** | letta | 8283 | https://letta.zeusops.dev | ✅ Running |
| **Letta PostgreSQL** | letta-postgres | 5432 | Internal | ✅ Running |

### Automation Stack

| Service | Containers | Port | Access | Status |
|---------|-----------|------|--------|--------|
| **Dify** | dify-api, dify-web, dify-nginx | 3000 | https://dify.zeusops.dev | ✅ Running (7 total) |
| **n8n** | n8n-main, n8n-postgres, n8n-worker | 5678 | https://n8n.zeusops.dev | ✅ Running (7 total) |

### IoT Stack

| Service | Container | Port | Access | Status |
|---------|-----------|------|--------|--------|
| **EMQX** | emqx | 18083 | https://emqx.zeusops.dev | ✅ Running |
| **Node-RED** | node-red | 1880 | https://nodered.zeusops.dev | ✅ Running |
| **Home Assistant** | home-assistant | 8123 | https://ha.zeusops.dev | ✅ Running |

### Memory & Data

| Service | Container | Port | Status |
|---------|-----------|------|--------|
| **Zep Memory** | zep-memory | 8000 | ✅ Running |
| **Zep PostgreSQL** | zep-postgres | 5432 | ✅ Running |

### Speech Services

| Service | Container | Port | Status |
|---------|-----------|------|--------|
| **Whisper STT** | whisper-stt | 9000 | ✅ Running |

---

## 🏗️ Stack Organization

```
legacy-ai-stacks/
├── infrastructure/
│   └── portainer/
│       caddy/
│
├── ai/
│   └── (Letta - deployed via Portainer)
│
├── automation/
│   ├── dify/              (7 containers)
│   └── n8n/               (7 containers)
│
├── iot/
│   ├── emqx/              (1 container)
│   ├── node-red/          (1 container)
│   └── home-assistant/    (1 container)
│
└── apps/
    └── (Legacy configs)
```

---

## 🌐 Domain Configuration (Caddy)

**Primary Access Point**: https://command.zeusops.dev (Portainer)

**Active Domains**:
```
command.zeusops.dev      → Portainer (container management)
letta.zeusops.dev        → Letta (agent framework)
n8n.zeusops.dev          → n8n (workflow automation)
dify.zeusops.dev         → Dify (AI platform)
nodered.zeusops.dev      → Node-RED (automation)
ha.zeusops.dev           → Home Assistant (smart home)
emqx.zeusops.dev         → EMQX (MQTT broker)
```

**Removed Domains**:
```
❌ dashboard.zeusops.dev (removed)
❌ mastra.zeusops.dev    (removed)
```

---

## 📦 Container Status Breakdown

### Healthy Containers (25)
- ✅ Portainer
- ✅ Caddy
- ✅ Dify (7 containers)
- ✅ n8n (7 containers)
- ✅ Letta
- ✅ Letta PostgreSQL
- ✅ EMQX
- ✅ Node-RED
- ✅ Home Assistant
- ✅ Zep Memory
- ✅ Zep PostgreSQL
- ✅ Whisper STT

---

## 🔗 Quick Access Commands

### SSH to VM
```bash
ssh -i "/Users/zeuslegacy/Downloads/LightsailDefaultKey-us-east-1 (1).pem" ubuntu@3.215.29.128
```

### Check All Services
```bash
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
```

### View Specific Service Logs
```bash
docker logs <service-name> --tail 50

# Examples:
docker logs portainer --tail 50
docker logs caddy --tail 50
docker logs letta --tail 50
```

### Restart a Service
```bash
docker restart <service-name>
```

### Check Service Health
```bash
docker ps --format 'table {{.Names}}\t{{.Status}}' | grep -E "healthy|unhealthy|starting"
```

---

## 📊 Resource Usage

**VM Specifications**:
- CPU: 8 vCPU (AWS Lightsail)
- RAM: 32GB
- Disk: 640GB SSD
- OS: Ubuntu 24.04 LTS

**Current Usage** (as of Oct 23):
- Disk: ~60GB / 640GB (10% used)
- Memory: ~6.3GB / 32GB (20% used)
- Load Average: Normal

---

## 🛠️ Deployment Procedures

### Deploying New Service
1. Add docker-compose.yml to appropriate stack directory
2. SSH to VM and deploy: `docker compose up -d`
3. Add domain to Caddy configuration if needed
4. Restart Caddy: `docker restart caddy`
5. Test domain access

### Scaling a Service
1. Edit docker-compose.yml (e.g., change worker count)
2. Redeploy: `docker compose up -d`
3. Verify: `docker ps`

### Backup Services
```bash
./scripts/maintenance/backup-all.sh    # Full backup
./scripts/maintenance/backup-databases.sh # DB only
```

---

## 📝 Notes

- **Portainer** is the primary UI for container management
- **Caddy** automatically handles SSL certificates and renewal
- **Health checks** are configured for all critical services
- **Automatic restart** enabled for all containers
- **Backup automation** runs daily at 2 AM via cron

---

**Next Steps**: 
- See [OPERATIONS.md](../operations/OPERATIONS.md) for daily operations
- See [QUICK-REFERENCE.md](../operations/QUICK-REFERENCE.md) for command reference
- See [TROUBLESHOOTING.md](../operations/stacks/APPS-TROUBLESHOOTING.md) for common fixes


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
