---
id: root-00-20260113225954
title: "KNOWLEDGE BASE"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "KNOWLEDGE BASE"
---
# ONE MIND CODEX - Knowledge Base

> **Extracted from 469 README files** | Consolidated knowledge to support all pillars

---

## Table of Contents

1. [System Architecture](#1-system-architecture)
2. [Infrastructure & Deployment](#2-infrastructure--deployment)
3. [AI Agents (Letta)](#3-ai-agents-letta)
4. [Automation (n8n)](#4-automation-n8n)
5. [Database (Supabase)](#5-database-supabase)
6. [IoT & Smart Home](#6-iot--smart-home)
7. [Services Quick Reference](#7-services-quick-reference)
8. [Scripts & Commands](#8-scripts--commands)
9. [Security & Credentials](#9-security--credentials)
10. [Pillar-Specific Knowledge](#10-pillar-specific-knowledge)

---

## 1. System Architecture

### Overview

```
ONE MIND CODEX Architecture
============================

                    ┌─────────────────────────────────┐
                    │     Coolify (AWS Lightsail)     │
                    │        Deployment Platform       │
                    └─────────────────────────────────┘
                                    │
        ┌───────────────────────────┼───────────────────────────┐
        │                           │                           │
   ┌────▼────┐               ┌──────▼──────┐              ┌─────▼─────┐
   │   AI    │               │ Automation  │              │    IoT    │
   │  Layer  │               │    Layer    │              │   Layer   │
   └────┬────┘               └──────┬──────┘              └─────┬─────┘
        │                           │                           │
   ┌────▼────┐               ┌──────▼──────┐              ┌─────▼─────┐
   │  Letta  │               │    n8n      │              │   MQTT    │
   │ Agents  │◄──────────────►  Workflows  │◄─────────────►  (EMQX)   │
   └────┬────┘               └──────┬──────┘              └─────┬─────┘
        │                           │                           │
        └───────────────────────────┼───────────────────────────┘
                                    │
                         ┌──────────▼──────────┐
                         │     Supabase        │
                         │  (PostgreSQL + AI)  │
                         └─────────────────────┘
```

### Core Components

| Component | Purpose | Technology |
|-----------|---------|------------|
| **Deployment** | Container orchestration | Coolify on AWS Lightsail |
| **AI Agents** | Intelligent assistants | Letta Cloud |
| **Automation** | Workflow processing | n8n (Queue Mode) |
| **Database** | Data persistence + vectors | Supabase (PostgreSQL) |
| **IoT Hub** | Smart home integration | Home Assistant |
| **Message Broker** | Device communication | EMQX (MQTT) |
| **Reverse Proxy** | SSL/Routing | Caddy |
| **Monitoring** | Observability | Grafana, Prometheus, Langfuse |

### Infrastructure Directory Structure

```
infrastructure/
├── services/
│   ├── discord-bot/          # Discord command center (Coolify/Nixpacks)
│   ├── automation/n8n/       # Workflow automation (Docker)
│   └── observability/
│       ├── monitoring/       # Grafana, Prometheus, Alertmanager
│       ├── uptime-kuma/      # Status page
│       └── langfuse/         # LLM observability
│
├── letta/                    # AI Agent management
│   ├── agents/               # Agent config exports
│   ├── tools/                # Discord, Todoist tools
│   └── scripts/              # Tool attachment scripts
│
├── config/
│   ├── agents/               # Agent YAML configs
│   ├── mcp-servers/          # MCP server configs
│   └── api_endpoints.yaml    # API configurations
│
├── scripts/                  # Deployment & maintenance scripts
├── docs/                     # Operations documentation
├── docker-compose.yml        # Main compose file
└── .env.example              # Environment template
```

---

## 2. Infrastructure & Deployment

### Server Access

| Resource | URL/Address |
|----------|-------------|
| **Coolify Dashboard** | http://44.212.233.193:8000 |
| **Server IP** | 3.215.29.128 |
| **SSH Alias** | `ssh legacy-ai` |

### Service Ports & URLs

| Service | Port | Domain | Notes |
|---------|------|--------|-------|
| Discord Bot | 8100 | - | Letta agent interface |
| n8n | 5678 | n8n.zeusops.dev | Workflow automation |
| Grafana | 3000 | grafana.zeusops.dev | Dashboards |
| Prometheus | 9090 | - | Metrics |
| Uptime Kuma | 3001 | status.zeusops.dev | Status page |
| Langfuse | - | langfuse.zeusops.dev | LLM tracking |
| Home Assistant | 8123 | homeassistant.zeusops.dev | IoT hub |
| Dify | - | dify.zeusops.dev | AI workflows |

### Quick Deployment Commands

```bash
# Copy environment template
cp .env.example .env

# Run Discord Bot locally
cd services/discord-bot
pip install -r requirements.txt
python bot.py

# Deploy via Coolify
# 1. Access Coolify: http://44.212.233.193:8000
# 2. Connect GitHub repo: Zeus-Delacruz/One-Mind-OS
# 3. Configure services
```

### Deployment Scripts

| Script | Purpose |
|--------|---------|
| `deploy.sh` | Deploy all services |
| `backup.sh` | Backup data volumes |
| `restore.sh` | Restore from backup |
| `health-check.sh` | Check service health |
| `setup-vm.sh` | Initial VM setup |

---

## 3. AI Agents (Letta)

### Agent Roster

| Agent | ID | Role | Pillar |
|-------|-----|------|--------|
| **Legacy-AI** | `agent-da585f6f-...` | Central coordinator | All |
| **HP-Conductor** | `agent-1e1155fa-...` | Holistic Performance | HP |
| **LE-Conductor** | `agent-3ef12371-...` | Legacy Evolution | LE |
| **GE-Conductor** | `agent-052cf675-...` | Generational Entrepreneurship | GE |
| **IT-Conductor** | `agent-2655ec04-...` | Innovative Technology | UI |

### Attach Tools to Agents

```bash
cd letta/scripts
python attach_discord_tools.py
```

### Agent Configuration (.af files)

Contents of agent config files:
- Agent persona and personality
- System instructions
- Memory configuration
- Tool access permissions
- Model preferences

---

## 4. Automation (n8n)

### Queue Mode Architecture

```
n8n Queue Mode Components
==========================

┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  n8n-main   │     │  n8n-worker │     │  n8n-worker │
│  (UI/API)   │     │     -1      │     │     -2      │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌──────▼──────┐
                    │   Redis     │
                    │   (Queue)   │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ PostgreSQL  │
                    │ (Database)  │
                    └─────────────┘
```

### n8n Components

- **n8n-main**: Main instance (UI, scheduling, coordination)
- **n8n-worker-1 & 2**: Workers (10 concurrent each)
- **n8n-webhook**: Dedicated webhook processor
- **n8n-task-runner**: Code isolation runtime
- **n8n-postgres**: PostgreSQL database
- **n8n-redis**: Redis queue broker

### Critical Configuration

```bash
# Generate encryption key (MUST be identical across all services)
openssl rand -base64 32

# Generate runner auth token
openssl rand -base64 32
```

### Deployment

```bash
# Prepare server directory
ssh legacy-ai "sudo mkdir -p /opt/legacy-ai/applications/n8n"

# Upload configuration
scp docker-compose.yml legacy-ai:/opt/legacy-ai/applications/n8n/

# Deploy
ssh legacy-ai "cd /opt/legacy-ai/applications/n8n && sudo docker compose up -d"

# Verify (should see 7 containers)
ssh legacy-ai "sudo docker ps | grep n8n"
```

### Scaling Workers

```yaml
# Add to docker-compose.yml
n8n-worker-3:
  image: docker.n8n.io/n8nio/n8n:latest
  container_name: n8n-worker-3
  command: worker --concurrency=10
  # ... (copy environment from worker-1)
```

### Monitor Queue

```bash
# View Redis queue depth
ssh legacy-ai "sudo docker exec n8n-redis redis-cli LLEN bull:n8n:waiting"

# View container stats
ssh legacy-ai "sudo docker stats n8n-main n8n-worker-1 n8n-worker-2"
```

---

## 5. Database (Supabase)

### Why Supabase?

| Feature | Benefit |
|---------|---------|
| **PostgreSQL** | Full SQL power with vector support |
| **Managed** | Automated backups, scaling, security |
| **Authentication** | Built-in user management and JWT |
| **Realtime** | Live subscriptions to changes |
| **REST/GraphQL** | Auto-generated APIs |
| **Vector Support** | Native pgvector for AI embeddings |

### Connection Template

```env
# Supabase Project
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key

# Database
DATABASE_URL=postgresql://postgres:PASSWORD@db.your-project.supabase.co:5432/postgres
DATABASE_POOL_URL=postgresql://postgres:PASSWORD@db.your-project.supabase.co:6543/postgres
```

### Vector Table Setup (for AI)

```sql
-- Enable pgvector
CREATE EXTENSION IF NOT EXISTS vector;

-- Create embeddings table
CREATE TABLE embeddings (
  id bigserial PRIMARY KEY,
  content text NOT NULL,
  embedding vector(1536),
  metadata jsonb,
  created_at timestamp DEFAULT now()
);

-- Create similarity search index
CREATE INDEX ON embeddings USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);
```

### Row-Level Security

```sql
-- Enable RLS
ALTER TABLE your_table ENABLE ROW LEVEL SECURITY;

-- Users can only see their own data
CREATE POLICY "Users can view their own data"
  ON your_table FOR SELECT
  USING (auth.uid() = user_id);
```

### Integration Code

**JavaScript:**
```javascript
import { createClient } from '@supabase/supabase-js'
const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY)
const { data, error } = await supabase.from('table').select('*')
```

**Python:**
```python
from supabase import create_client
supabase = create_client(url, key)
response = supabase.table("table").select("*").execute()
```

---

## 6. IoT & Smart Home

### Home Assistant Architecture

```
IoT Devices/Sensors
    ↓
Mosquitto MQTT Broker
    ↓
Home Assistant (automation hub)
    ↓
├─> Node-RED (complex logic)
├─> n8n (workflow automation)
└─> Dify (AI decision making)
```

### Access

- **Web UI**: https://homeassistant.zeusops.dev
- **Mobile**: Home Assistant Companion (iOS/Android)

### MQTT Discovery Example

```python
import paho.mqtt.publish as publish

config = {
    "name": "Living Room Temperature",
    "state_topic": "home/living-room/temperature",
    "unit_of_measurement": "°F",
    "device_class": "temperature"
}

publish.single(
    "homeassistant/sensor/living-room-temp/config",
    json.dumps(config),
    hostname="mqtt.zeusops.dev",
    auth={'username': 'admin', 'password': 'your-password'}
)
```

### Automation with AI (Dify)

```yaml
automation:
  - alias: "AI Thermostat Control"
    trigger:
      - platform: time_pattern
        minutes: "/30"
    action:
      - service: rest_command.dify_temperature_decision
        response_variable: ai_response
      - service: climate.set_temperature
        data:
          temperature: "{{ ai_response.recommended_temp }}"
```

---

## 7. Services Quick Reference

### Start/Stop Services

```bash
# All services
ssh legacy-ai "cd /opt/legacy-ai && sudo docker compose up -d"
ssh legacy-ai "cd /opt/legacy-ai && sudo docker compose down"

# Individual service
ssh legacy-ai "sudo docker restart n8n-main"
ssh legacy-ai "sudo docker logs home-assistant -f"
```

### Health Checks

```bash
# Check all containers
ssh legacy-ai "sudo docker ps"

# Check specific service
ssh legacy-ai "sudo docker exec n8n-postgres pg_isready -U n8n"

# View resource usage
ssh legacy-ai "sudo docker stats --no-stream"
```

### Backup Commands

```bash
# n8n database
ssh legacy-ai "sudo docker exec n8n-postgres pg_dump -U n8n -d n8n > backup.sql"

# Home Assistant config
ssh legacy-ai "sudo docker run --rm -v home_assistant_config:/config -v $(pwd):/backup alpine tar czf /backup/ha-backup.tar.gz /config"
```

---

## 8. Scripts & Commands

### Environment Setup

```bash
# Copy template
cp .env.example .env

# Generate secrets
openssl rand -base64 32  # Encryption key
python3 -c "import secrets; print(secrets.token_hex(32))"  # Token
```

### Caddy (Reverse Proxy)

```caddyfile
n8n.zeusops.dev {
    reverse_proxy n8n-main:5678
}

homeassistant.zeusops.dev {
    reverse_proxy home-assistant:8123
}
```

### Restart Caddy

```bash
ssh legacy-ai "cd /opt/legacy-ai/infrastructure/caddy && sudo docker compose restart caddy"
```

---

## 9. Security & Credentials

### Required Environment Variables

| Variable | Service | Notes |
|----------|---------|-------|
| `N8N_ENCRYPTION_KEY` | n8n | Must match across all n8n services |
| `N8N_RUNNERS_AUTH_TOKEN` | n8n | Task runner authentication |
| `SUPABASE_SERVICE_KEY` | Supabase | Server-side only, keep secret |
| `LETTA_API_KEY` | Letta | Agent management |
| `DISCORD_TOKEN` | Discord Bot | Bot authentication |
| `TODOIST_TOKEN` | Todoist | Task integration |

### Security Checklist

- [ ] Encryption keys generated and secured
- [ ] Database passwords changed from defaults
- [ ] HTTPS enabled via Caddy
- [ ] RLS policies implemented in Supabase
- [ ] Long-lived access tokens for integrations
- [ ] 2FA enabled on admin accounts
- [ ] Regular credential rotation scheduled

---

## 10. Pillar-Specific Knowledge

### UI (00-24) - Unified Intelligence

**Focus**: Technology, AI, Infrastructure, Automation

**Key Services**:
- Letta Agents (AI assistants)
- n8n (workflow automation)
- Infrastructure monitoring
- Security systems

### HP (25-49) - Holistic Performance

**Focus**: Health, Mind, Skills, Finance, Lifestyle

**Automations**:
- Health tracking integrations
- Habit reminders
- Financial alerts
- Learning schedule management

### LE (50-74) - Legacy Evolution

**Focus**: Family, Home, Heritage, Children

**Integrations**:
- Home Assistant (smart home)
- Family calendar sync
- Security monitoring
- Household automation

### GE (75-99) - Generational Entrepreneurship

**Focus**: Business, Ventures, Credit, Education

**Workflows**:
- Lead management automation
- Revenue tracking
- Client communication
- Business intelligence

---

## Resources

### Documentation Links

- [n8n Queue Mode](https://docs.n8n.io/hosting/scaling/queue-mode/)
- [Supabase Docs](https://supabase.com/docs)
- [Home Assistant](https://www.home-assistant.io/docs/)
- [Letta Documentation](https://docs.letta.com)

### Internal Docs

- Architecture: `00-24 UI/02 Agents/Architecture/`
- Operations: `00-24 UI/10 Infrastructure/`
- Scripts: `00-24 UI/12 Development/`

---

## 11. Cognitive Architecture (Complete Intelligent Being)

### The Vision

This is not just "another AI stack." It's the architecture for building an **intelligent being** with:

- **Perception** (Senses) - Sees, hears, speaks to the world
- **Cognition** (Brain) - Reasons, decides, learns, remembers
- **Action** (Tools) - Executes decisions in the world
- **Memory** (Storage) - Learns from experience, evolves

### System Architecture

```
         🌍 EXTERNAL WORLD 🌍
              ↓ ↑
        ┌─────────────────┐
        │ 👁️👂👄 SENSES   │ (Perception - LiveKit)
        └────────┬────────┘
                 ↓
        ┌─────────────────┐
        │  🧠 BRAIN       │ (Cognition - Letta)
        └────────┬────────┘
                 ↓
        ┌─────────────────┐
        │  🔧 TOOLS       │ (Action - Dify, n8n, HA, Node-RED)
        └────────┬────────┘
                 ↓
        ┌─────────────────┐
        │  💾 MEMORY      │ (Learning - Supabase, PostgreSQL)
        └─────────────────┘
```

### Component Details

| Layer | Components | Purpose |
|-------|------------|---------|
| **Senses** | LiveKit, AssemblyAI, Cartesia, GPT-4o Vision | Eyes, ears, voice |
| **Brain** | Letta (self-editing memory, tool orchestration) | Reasoning, decisions |
| **Tools** | Dify, n8n, Home Assistant, Node-RED, EMQX | Execute in real world |
| **Memory** | Supabase, PostgreSQL, Redis | Store, learn, evolve |

### Performance Targets

**Online Mode (Phase 2A):**
- Voice Latency: <100ms (STT 35ms + LLM 25ms + TTS 30ms)
- Vision Latency: <1000ms
- Memory Access: <20ms

**Offline Mode (Phase 2B - Jetson Nano):**
- Voice Latency: <200ms (Whisper 80ms + Llama 60ms + Piper 50ms)
- Vision FPS: 60fps (YOLOv8-nano)
- Memory Access: <10ms

### Cost Breakdown

| Component | Service | Monthly Cost |
|-----------|---------|--------------|
| Senses (STT) | AssemblyAI | $15 |
| Senses (TTS) | Cartesia | $10 |
| Brain (LLM) | OpenAI GPT-4o | $20 |
| Memory (DB) | Supabase Cloud | $25 |
| **Total** | | **$75/month** |

---

## 12. Capabilities (What Legacy AI Can DO)

### Communication
- `voice-calling/` - Phone calls (Twilio, LiveKit SIP)
- `email/` - Email sending/receiving (SendGrid, IMAP)
- `messaging/` - SMS, WhatsApp, Slack

### Digital Actions
- `workflow-automation/` - AI workflows (Dify)
- `task-automation/` - General automation (n8n)
- `browser-automation/` - Web interaction (ByteBot)
- `smart-home/` - Home automation (Home Assistant)
- `iot-messaging/` - IoT communication (EMQX MQTT)
- `visual-flows/` - Flow programming (Node-RED)
- `data-pipelines/` - Data integration (Airbyte, Meltano)

### Physical Actions (Future)
- `humanoid-robotics/` - Unitree G1/H1/Go2
- `ros2-integration/` - Robot Operating System 2
- `agricultural-robotics/` - Farmbot
- `motion-control/` - Movement and locomotion
- `navigation/` - Path planning and localization

### MCP Servers
- GitHub, AWS, Filesystem, SSH integrations

---

## 13. Environments (Where Legacy AI Runs)

| Environment | Purpose | Hardware |
|-------------|---------|----------|
| **Production** | Live AWS Lightsail VM | Cloud |
| **Edge** | Offline/Edge hardware | Jetson Nano, Unitree robots |
| **Development** | Local dev environment | Local machine |
| **Staging** | Testing before production | Cloud/Local |

### Interfaces
- `web-app/` - CopilotKit, Lobe Chat, Open WebUI
- `mobile-app/` - iOS and Android
- `voice-interface/` - Phone menu/IVR
- `ar-interface/` - Brilliant Labs AR glasses
- `api-server/` - REST and GraphQL APIs

---

## 14. Dify AI Platform

### Components
- **dify-api**: Backend API service
- **dify-worker**: Background task processor (Celery)
- **dify-web**: Frontend application
- **dify-sandbox**: Code execution sandbox
- **dify-postgres**: PostgreSQL with pgvector
- **dify-redis**: Redis for caching/queuing

### Features
- LLM Applications with no-code interface
- Visual workflow builder (400+ integrations)
- RAG with vector search (pgvector)
- Autonomous AI agents with tool use
- Production-ready API

### Deployment

```bash
# Generate secret key
openssl rand -base64 42

# Deploy
ssh legacy-ai "cd /opt/legacy-ai/applications/dify && sudo docker compose up -d"

# Initialize database (REQUIRED!)
ssh legacy-ai "sudo docker compose exec dify-api flask db upgrade"
```

### Access
- **Web UI**: https://dify.zeusops.dev
- **API Docs**: https://dify.zeusops.dev/api/docs

---

## 15. EMQX MQTT Broker

### Connection Details

| Parameter | Value |
|-----------|-------|
| MQTT Port | 1884 |
| WebSocket Port | 8084 |
| MQTT/SSL Port | 8885 |
| Dashboard | https://emqx.zeusops.dev |

### Integration Examples

**Node-RED:**
```
Server: emqx
Port: 1883
Client ID: nodered-client
```

**Python:**
```python
from paho.mqtt import client as mqtt_client
client = mqtt_client.Client('python-client')
client.username_pw_set(username, password)
client.connect('emqx', 1883)
```

### vs Mosquitto

| Feature | EMQX | Mosquitto |
|---------|------|-----------|
| Dashboard | ✅ Full web UI | ❌ None |
| Rule Engine | ✅ Built-in | ❌ External |
| Clustering | ✅ Native | ⚠️ Limited |
| Use Case | Production IoT | Lightweight broker |

---

## 16. Node-RED Visual Automation

### Access
- **Web UI**: https://nodered.zeusops.dev
- **Dashboard**: https://nodered.zeusops.dev/ui

### Architecture

```
IoT Devices → Mosquitto MQTT
                ↓
            Node-RED (automation logic)
                ↓
        ├─> Home Assistant
        ├─> n8n (advanced workflows)
        ├─> Dify (AI processing)
        └─> Database/APIs
```

### Recommended Nodes

```
node-red-dashboard                          # UI dashboard
node-red-contrib-mqtt-broker                # Enhanced MQTT
node-red-contrib-home-assistant-websocket   # Home Assistant
node-red-node-sqlite                        # SQLite database
```

### MQTT Topic Convention

```
home/
  ├── sensors/
  │   ├── temperature/living-room
  │   ├── humidity/bedroom
  │   └── motion/entrance
  ├── actuators/
  │   ├── lights/kitchen
  │   └── thermostat/main
  └── system/
      ├── status
      └── commands
```

---

## 17. Scripts Organization

### New Script Locations

| Script Type | New Location |
|-------------|--------------|
| Deployment | `ENVIRONMENTS/production/` |
| Maintenance | `ENVIRONMENTS/production/infrastructure/` |
| Monitoring | `OBSERVABILITY/monitoring/` |

### Deployment Scripts
- `deploy-legacy-ai.sh`
- `deploy-letta.sh`
- `deploy-n8n.sh`
- `deploy-iot.sh`
- `setup-caddy.sh`

### Maintenance Scripts
- `backup-all.sh`
- `backup-databases.sh`
- `setup-cron-backups.sh`
- `update-memory.sh`

### Monitoring Scripts
- `systems-check.sh`
- `quick-fix.sh`
- `send-alert.sh`

---

## 18. Pillar Folder Templates

### HP - Identity & Purpose (Example: 25)

```
25 Identity/
├── Core Identity/
│   ├── Mission-Statement.md
│   ├── Core-Values.md
│   ├── Personal-Philosophy.md
│   └── Identity-Evolution.md
├── Vision & Goals/
│   ├── Life-Vision.md
│   ├── 5-Year-Vision.md
│   └── 1-Year-Vision.md
├── Journal/
│   ├── Ideas.md
│   ├── Dream Notes.md
│   └── Reflections/
├── Daily Notes/
├── Early Journey/
└── Personal Literature/
    ├── Poetry/
    ├── Reflections/
    └── Stories/
```

### Best Practices
- Review Mission and Values monthly
- Update vision quarterly
- Journal regularly (daily or weekly)
- Track identity evolution at major life transitions

---

**Last Updated**: January 2, 2026
**Source**: Consolidated from 469 README files (comprehensive review)
**Maintained by**: ONE MIND CODEX System
