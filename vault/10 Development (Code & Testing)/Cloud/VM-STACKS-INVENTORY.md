---
id: ui-12-20260111053926
title: "VM Stacks Inventory"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/12]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "VM Stacks Inventory"
---

# VM Stacks Inventory

**Generated**: October 23, 2025 (Updated after cleanup)  
**VM**: 3.215.29.128 (32GB RAM, 8vCPU, 640GB SSD)  
**Status**: ✅ All systems operational

---

## 📊 Running Containers Summary

**Total Containers**: 23 running (removed LiveKit agent, Whisper STT)  
**Total Memory Used**: ~4.8GB / 30.98GB (15.5%)  
**Total CPU Usage**: ~150% (multi-core)

**Changes**:
- ❌ Removed: `legacy-ai-agent`, `legacy-ai-token-server` (LiveKit voice agent)
- ❌ Removed: `whisper-stt` (Speech-to-Text GPU service)
- ✅ Moved: Letta to `stacks/ai/letta/` directory

---

## 🗂️ Stack Organization

### Infrastructure Stack
**Path**: `/opt/legacy-ai/infrastructure/`

#### Caddy (Reverse Proxy)
- **Container**: `caddy`
- **Image**: `caddy:latest`
- **Status**: ✅ Up 2 hours
- **Ports**: 80, 443 (HTTP/HTTPS)
- **Config**: `/opt/legacy-ai/infrastructure/caddy/docker-compose.yml`
- **Function**: Reverse proxy for all active domains
- **Memory**: 23.86MiB
- **CPU**: 0.05%

---

### AI Stack
**Path**: `/opt/legacy-ai/stacks/ai/`

#### Letta (AI Agent Framework)
- **Container**: `letta`
- **Image**: `letta/letta:latest`
- **Status**: ✅ Up 2 days
- **Ports**: 8283 (API)
- **Config**: `/opt/legacy-ai/stacks/ai/letta/docker-compose.yml`
- **Memory**: 443.6MiB
- **CPU**: 0.81%
- **URL**: `https://letta.zeusops.dev`
- **Note**: Uses Letta's built-in PostgreSQL for memory management

#### Memory Management
- **Provider**: Letta (built-in memory system)
- **Status**: ✅ Using Letta's agent memory framework
- **Integration**: Direct PostgreSQL connection via Letta agent system

---

### IoT Stack
**Path**: `/opt/legacy-ai/stacks/iot/`

#### EMQX (Message Broker)
- **Container**: `emqx`
- **Image**: `emqx/emqx:latest`
- **Status**: ⚠️ Up 3 hours (unhealthy)
- **Ports**: 1883 (MQTT), 8083 (WebSocket), 8084-8885 (Secure)
- **Config**: `/opt/legacy-ai/stacks/iot/emqx/docker-compose.yml`
- **Memory**: 46.9MiB
- **CPU**: 0.19%
- **Note**: Health check failing - needs investigation

#### Node-RED (Node Orchestration)
- **Container**: `node-red`
- **Image**: `nodered/node-red:latest`
- **Status**: ⚠️ Up 3 hours (unhealthy)
- **Ports**: 1880
- **Config**: `/opt/legacy-ai/stacks/iot/node-red/docker-compose.yml`
- **Memory**: 40.68MiB
- **CPU**: 0.86%
- **Note**: Health check failing

#### Home Assistant (Smart Home)
- **Container**: `home-assistant`
- **Image**: `ghcr.io/home-assistant/home-assistant:stable`
- **Status**: ✅ Up 2 days (healthy)
- **Ports**: 8123
- **Config**: `/opt/legacy-ai/stacks/iot/home-assistant/docker-compose.yml`
- **Memory**: 250.4MiB
- **CPU**: 0.20%

---

### Automation Stack
**Path**: `/opt/legacy-ai/stacks/apps/`

#### Dify (AI Workflow Platform)
- **Containers**:
  - `dify-web` ✅ Up 2 days
  - `dify-nginx` ✅ Up 2 days
  - `dify-api` ✅ Up 2 days
  - `dify-worker` ✅ Up 2 days
  - `dify-sandbox` ✅ Up 2 days
  - `dify-postgres` ✅ Up 2 days (healthy)
  - `dify-redis` ✅ Up 2 days (healthy)
- **Images**: `langgenius/dify-*:latest`, `pgvector/pgvector:pg16`, `redis:7-alpine`
- **Ports**: 3000 (web nginx)
- **Config**: `/opt/legacy-ai/stacks/apps/dify/docker-compose.yml`
- **Status**: ✅ All healthy
- **Total Memory**: ~1.2GB
- **URL**: `https://dify.zeusops.dev`

#### n8n (Workflow Automation) - Queue Mode
- **Containers**:
  - `n8n-main` ✅ Up 2 days
  - `n8n-worker-1` ✅ Up 2 days
  - `n8n-worker-2` ✅ Up 2 days
  - `n8n-task-runner` ✅ Up 2 days
  - `n8n-webhook` ✅ Up 2 days
  - `n8n-postgres` ✅ Up 2 days (healthy)
  - `n8n-redis` ✅ Up 2 days (healthy)
- **Images**: `docker.n8n.io/n8nio/n8n:latest`, `postgres:16-alpine`, `redis:7-alpine`
- **Ports**: 5678 (main)
- **Config**: `/opt/legacy-ai/stacks/apps/n8n-queue-mode/docker-compose.yml`
- **Status**: ✅ All healthy
- **Total Memory**: ~1.3GB
- **URL**: `https://n8n.zeusops.dev`

---

### Management Stack

#### Portainer (Container Management)
- **Container**: `portainer`
- **Image**: `portainer/portainer-ee:2.33.2`
- **Status**: ✅ Up 2 days
- **Ports**: 8000, 9443
- **Config**: Docker metadata only
- **Memory**: 48.76MiB
- **CPU**: 0.00%
- **URL**: `https://portainer.zeusops.dev`

---

## 📍 Directory Structure on VM

```
/opt/legacy-ai/
├── .env                              # Global environment variables
├── config/                           # Configuration files
├── data/                             # Persistent data
├── logs/                             # Application logs
│
├── infrastructure/                   # Infrastructure services
│   ├── caddy/
│   │   └── docker-compose.yml
│   └── emqx/
│       └── docker-compose.yml
│
├── stacks/                           # Main application stacks
│   ├── README.md                     # Stack documentation
│   ├── ai/                           # AI services
│   │   └── letta/
│   │       └── docker-compose.yml
│   ├── apps/                         # Automation services
│   │   ├── dify/
│   │   │   └── docker-compose.yml
│   │   └── n8n-queue-mode/
│   │       └── docker-compose.yml
│   └── iot/                          # IoT services
│       ├── emqx/
│       │   └── docker-compose.yml
│       ├── home-assistant/
│       │   └── docker-compose.yml
│       └── node-red/
│           └── docker-compose.yml
│
└── applications/                     # Legacy applications (deprecated)
    ├── dify/
    ├── n8n/
    ├── letta/                         # (old location - deprecated)
    ├── home-assistant/
    └── node-red/
```

---

## 🌐 Active Domains

| Domain | Service | Port | Status |
|--------|---------|------|--------|
| `dify.zeusops.dev` | Dify | 3000 | ✅ Active |
| `n8n.zeusops.dev` | n8n | 5678 | ✅ Active |
| `letta.zeusops.dev` | Letta | 8283 | ✅ Active |
| `home.zeusops.dev` | Home Assistant | 8123 | ✅ Active |
| `portainer.zeusops.dev` | Portainer | 9443 | ✅ Active |
| `mqtt.zeusops.dev` | EMQX Dashboard | 18083 | ✅ Active |
| `nodered.zeusops.dev` | Node-RED | 1880 | ✅ Active |

---

## ⚠️ Health Issues

### 1. EMQX Container (unhealthy)
- **Status**: Up 3 hours (unhealthy)
- **Issue**: Health check failing
- **Action Required**: Investigate logs and fix
- **Command**: `docker logs emqx --tail 100`

### 2. Node-RED Container (unhealthy)
- **Status**: Up 3 hours (unhealthy)
- **Issue**: Health check failing
- **Action Required**: Check health endpoint
- **Command**: `docker logs node-red --tail 100`

### 3. Zep Memory Container (unhealthy)
- **Status**: Up 17 hours (unhealthy)
- **Issue**: Memory service health check failing
- **Action Required**: Check service connectivity
- **Command**: `docker logs zep-memory --tail 100`

---

## 📈 Resource Usage Summary

| Component | Memory | CPU | Status |
|-----------|--------|-----|--------|
| **Total System** | 4.8GB / 31GB | ~150% | ✅ Good |
| Dify Stack | ~1.2GB | 1-3% | ✅ Good |
| n8n Stack | ~1.3GB | 1-2% | ✅ Good |
| EMQX | 46.9MB | 0.19% | ⚠️ Unhealthy |
| Letta | 443.6MB | 0.81% | ✅ Good |
| Home Assistant | 250.4MB | 0.20% | ✅ Good |
| Caddy | 23.86MB | 0.05% | ✅ Good |
| Portainer | 48.76MB | 0.00% | ✅ Good |

---

## 🔧 Maintenance Notes

### Removed Services
- ❌ LiveKit Agent: Removed from VM (2 containers: agent + token-server)
- ❌ Whisper STT: Removed GPU-intensive speech-to-text service

### Reorganization
- ✅ Letta: Moved from `applications/letta/` to `stacks/ai/letta/`
- ✅ Docker-compose.yml now in correct directory structure

### Latest Deployments
- ✅ Letta: Redeployed from new location
- ✅ Home Assistant: 2 days
- ✅ Dify: 2 days
- ✅ n8n: 2 days
- ✅ Caddy reloaded: 2 hours

### Known Issues
1. Node-RED health checks need investigation
2. Zep memory service health check failing
3. EMQX health checks need investigation

### Next Steps
1. Investigate and fix unhealthy containers
2. Create monitoring dashboard
3. Set up automated health recovery
4. Consider consolidating old `applications/` directory

---

**Last Updated**: October 23, 2025, 21:35 UTC  
**Verified By**: Direct SSH query to VM  
**Changes**: Removed LiveKit agent & Whisper STT, moved Letta to stacks/ai/
