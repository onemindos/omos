---
id: ui-02-20260111053925
title: "🎯 Legacy AI - Unified Development Master"
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
ai_summary: "🎯 Legacy AI - Unified Development Master"
---

# 🎯 Legacy AI - Unified Development Master

**Purpose**: Single comprehensive source of truth for ALL Legacy AI development  
**Created**: October 23, 2025  
**Status**: Architecture Planning + Active Development  
**Scope**: Current infrastructure + Future vision (Letta/LiveKit/Offline AI)  
**Audience**: Personal development | Architectural decision-making

---

## 📑 Quick Navigation

- [Phase 1: Current Infrastructure](#phase-1-current-infrastructure-now)
- [Phase 2: Stateful AI Agent](#phase-2-stateful-ai-agent-q1-2026)
- [Phase 3: Advanced Intelligence](#phase-3-advanced-intelligence-q2-2026)
- [Development Environment](#development-environment)
- [API Integrations](#api-integrations)
- [Deployment & Operations](#deployment--operations)

---

# PHASE 1: CURRENT INFRASTRUCTURE (NOW ✅)

## System Overview

**Status**: ✅ Production Ready | All 25 containers operational | Fully hardened & secured

### VM Specifications
```
Provider: AWS Lightsail (us-east-1)
IP: 3.215.29.128 (Ubuntu 24.04 LTS)
Specs: 32GB RAM | 8vCPU | 640GB SSD
Access: SSH key + Tailscale VPN
Security: Ed25519 SSH, ChaCha20-Poly1305, automatic security updates
```

### Network Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    AWS Lightsail VM                          │
│                  3.215.29.128 (Ubuntu 24.04)                │
│                   32GB RAM | 8vCPU | 640GB SSD              │
└─────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────────┐
│                      Caddy Reverse Proxy                          │
│          (HTTPS termination, SSL/TLS, domain routing)             │
└───────────────────────────────────────────────────────────────────┘

┌────────────────┬──────────────────┬──────────────┬────────────┐
│ INFRASTRUCTURE │   AUTOMATION     │      IoT     │ MANAGEMENT │
│   (35% load)   │    (40% load)    │  (10% load)  │  (5% load) │
├────────────────┼──────────────────┼──────────────┼────────────┤
│ • Caddy        │ • Dify (LLM UI)  │ • EMQX       │ • Portainer│
│ • PostgreSQL   │ • n8n (Workflow) │ • Node-RED   │ • Monitoring
│ • Redis        │ • PostgreSQL     │ • Home Assist│
│ • EMQX         │ • Redis          │              │
└────────────────┴──────────────────┴──────────────┴────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                    Current AI Stack                              │
│         Letta (Agent) | PostgreSQL | Redis | MQTT              │
└──────────────────────────────────────────────────────────────────┘
```

## Active Services (25 Containers)

### 🔧 Infrastructure Layer (CRITICAL)
```
✅ Portainer          (9443)  → Container management UI
✅ Caddy              (80,443) → Reverse proxy + HTTPS
✅ PostgreSQL         (5432)  → Primary database (pgvector enabled)
✅ Redis              (6379)  → Caching layer
```

### 🤖 AI & Agent Layer
```
✅ Letta              (8283)  → Agent framework + persistent memory
✅ Letta PostgreSQL   (5432)  → Letta data store
✅ Whisper STT        (9000)  → Speech-to-text (offline)
```

### ⚙️ Automation Stack
```
✅ Dify API           (8080)  → LLM UI + workflow engine (7 containers)
✅ Dify PostgreSQL    (5432)  → Dify database
✅ Dify Redis         (6379)  → Dify cache
✅ n8n Main           (5678)  → Workflow automation (7 containers)
✅ n8n PostgreSQL     (5432)  → n8n database
✅ n8n Redis          (6379)  → n8n cache
✅ n8n Worker         (worker) → Background task processing (queue mode)
```

### 🏠 IoT & Home Automation
```
✅ EMQX               (1883)  → MQTT broker + dashboard (18083)
✅ Node-RED           (1880)  → Visual automation flows
✅ Home Assistant     (8123)  → Smart home hub
```

### 📊 Monitoring & Maintenance
```
✅ Health check agent → Automated system checks (daily)
✅ Backup system      → 3-tier automated backups (daily/weekly/on-demand)
✅ Update agent       → Security patches (automatic)
```

## 7 Active HTTPS Domains

All domains secured with Let's Encrypt via Caddy:

```
1. https://command.zeusops.dev      → Portainer (container management)
2. https://n8n.zeusops.dev          → n8n (workflow automation)
3. https://dify.zeusops.dev         → Dify (LLM UI & API)
4. https://letta.zeusops.dev        → Letta (AI agents)
5. https://nodered.zeusops.dev      → Node-RED (visual automation)
6. https://ha.zeusops.dev           → Home Assistant (IoT)
7. https://emqx.zeusops.dev         → EMQX (MQTT broker)
```

## Technology Stack - Phase 1

| Layer | Technology | Version | Purpose |
|-------|-----------|---------|---------|
| **Orchestration** | Docker Compose | v2 | Container management |
| **Reverse Proxy** | Caddy | latest | SSL/TLS, routing |
| **Workflow** | n8n | latest | Business automation |
| **LLM UI** | Dify | latest | AI interface |
| **Agents** | Letta | latest | Stateful AI agents |
| **Memory** | Letta | local | Agent memory & persistence |
| **IoT Hub** | Home Assistant | latest | Smart home |
| **Flow Editor** | Node-RED | latest | Visual automation |
| **MQTT** | EMQX | latest | Message broker |
| **Database** | PostgreSQL | latest | Primary storage |
| **Cache** | Redis | latest | Performance layer |
| **Management** | Portainer | latest | Container UI |
| **Speech** | Whisper | latest | STT (offline) |

## Current Development Environment

### Stack Directory Structure
```
legacy-ai-stacks/
├── infrastructure/         # Base services
│   ├── caddy/
│   └── portainer/
├── ai/                     # AI services
│   └── (Letta via Portainer)
├── automation/             # Workflow tools
│   ├── dify/
│   └── n8n-queue-mode/
├── iot/                    # Smart home
│   ├── emqx/
│   ├── node-red/
│   └── home-assistant/
├── database/               # Data storage
│   ├── postgresql/
│   ├── redis/
│   └── letta-memory/
└── apps/                   # Utilities
    └── portainer/
```

### Quick Start - Local Development

```bash
# 1. Clone repository
git clone https://github.com/Legacyx-ai/Legacy-ai-Hybrid.git
cd Legacy-ai-Hybrid

# 2. SSH to production VM
ssh -i "LightsailDefaultKey-us-east-1 (1).pem" ubuntu@3.215.29.128

# 3. Verify services running
docker ps --format 'table {{.Names}}\t{{.Status}}'

# 4. View live logs
docker-compose logs -f [service-name]

# 5. Access any service via domain
# Examples:
# https://command.zeusops.dev (Portainer)
# https://letta.zeusops.dev (Agent)
# https://dify.zeusops.dev (LLM UI)
```

### Environment Variables
```bash
# Core
COMPOSE_PROJECT_NAME=legacy-ai

# Caddy
CADDY_ADMIN=localhost:2019

# Dify
DIFY_PORT=8080
DIFY_DATABASE_URL=postgresql://postgres:password@dify-postgres:5432/dify_db

# n8n
N8N_PORT=5678
N8N_DATABASE_URL=postgresql://postgres:password@n8n-postgres:5432/n8n_db

# Letta
LETTA_PORT=8283
LETTA_DATABASE_URL=postgresql://postgres:password@letta-postgres:5432/letta_db

# EMQX
EMQX_PORT=1883
EMQX_DASHBOARD_PORT=18083

# Redis
REDIS_PORT=6379
REDIS_PASSWORD=secure-password

# See: /legacy-ai-stacks/[service]/.env.example for complete config
```

---

# PHASE 2: STATEFUL AI AGENT (Q1 2026)

## Vision: Multi-Modal, Stateful, Intelligent Agent

**Goal**: Build a "Jarvis-like" AI that knows your routine, predicts needs, has unlimited memory, and evolves with you.

### Phase 2 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              LiveKit Cloud (Voice/Video/Text)               │
│              Real-time communication platform               │
└───────────┬─────────────────────────────────────────────────┘
            │
┌───────────▼─────────────────────────────────────────────────┐
│                   Your LiveKit Agent                        │
│        (Python SDK with business logic)                     │
└───────────┬─────────────────────────────────────────────────┘
            │
        ┌───┴────────────────┬──────────────┬─────────────┐
        │                    │              │             │
┌───────▼────────┐  ┌────────▼─────┐  ┌────▼────┐  ┌────▼────┐
│  Letta Memory  │  │ LLM (Local   │  │ Vision  │  │ Tools   │
│ (Self-hosted)  │  │ or Cloud)    │  │ Engine  │  │ Access  │
│                │  │              │  │         │  │         │
│ - Session      │  │ - Ollama     │  │-Frigate │  │-n8n API │
│ - Episodic     │  │ - GPT-4o     │  │-CompreFace
│ - Semantic     │  │ - Anthropic  │  │         │  │-Home Asst
│ - Procedural   │  │              │  │         │  │-EMQX    │
└────────────────┘  └──────────────┘  └─────────┘  └─────────┘
        ↕                                              ↕
┌───────────────────────────────────────────────────────────────┐
│              Supporting Services (Phase 1 Stack)              │
│   PostgreSQL | Redis | EMQX | Node-RED | Home Assistant      │
└───────────────────────────────────────────────────────────────┘
```

## Phase 2 Technology Decisions

### Memory Layer: Letta (RECOMMENDED ✅)

**Why Letta**:
- ✅ Agent-centric memory management
- ✅ Persistent conversation context
- ✅ Built-in PostgreSQL integration
- ✅ Self-hosted (full control, no vendor lock-in)
- ✅ Open-source (Apache 2.0 license)

**Memory Architecture**:
```
1. Working Memory (RAM)
   └── Current conversation (last 20 messages)

2. Core Memory (Letta)
   ├── User Profile (persistent facts)
   ├── Preferences (communication style, habits)
   └── Active Goals (short-term objectives)

3. Episodic Memory (Letta + Postgres)
   ├── Conversation History (searchable)
   ├── Actions Taken (n8n logs)
   └── Events Experienced (timestamped)

4. Semantic Memory (pgvector)
   ├── Knowledge Base (RAG documents)
   ├── Skills & Procedures (how-to)
   └── Relationships (knowledge graphs)

5. Procedural Memory (Code)
   ├── Tool Definitions
   ├── Automation Rules
   └── Workflows (n8n)
```

**Live Integration Pattern**:
```python
from letta import create_client

client = create_client()
agent = client.get_agent("legacy-ai-agent")

# Automatically preserves:
# - Speaker identification
# - Conversation context
# - User preferences & profile
# - Cross-agent memory sharing
```

### Voice/Video Platform: LiveKit Cloud

**Why LiveKit Cloud** (not self-hosted):
- ✅ Global edge network (low latency)
- ✅ Handles scaling (WebRTC mesh)
- ✅ HTTPS, SIP integration, Egress/Ingress
- ✅ Cost: ~$2-5/month for personal use
- ✅ Can self-host later if needed

**LiveKit Features for Legacy AI**:
```
- Voice Pipeline: STT → LLM → TTS with interruption
- Vision: Video streaming + image analysis
- Text Chat: Messages alongside voice
- Agents Framework: Python SDK for custom logic
- MCP Support: Tool integrations (n8n, APIs)
- Local Models: Ollama, llama.cpp support
- Turn Detection: Smart interruption handling
```

**Cost Analysis**:
```
Personal use (2 hrs/day):
- Connection: 3,600 min/mo × $0.0005 = $1.80
- Bandwidth: 3GB × $0.12 = $0.36
- Total: ~$2-5/month ✅
```

### Offline LLM: Ollama + Llama 3.3 70B

**Setup**:
```bash
# Install
curl -fsSL https://ollama.com/install.sh | sh

# Pull model (40GB quantized)
ollama pull llama3.3:70b-instruct-q4_K_M

# Run
ollama serve  # Starts on localhost:11434
```

**Performance on Your 8vCPU / 32GB VM**:
- Speed: 5-10 tokens/sec (acceptable for voice)
- Context: 128K tokens
- Quality: Near GPT-4 for most tasks

### Offline Speech: Faster-Whisper + Piper

**STT (Faster-Whisper)**:
```python
from faster_whisper import WhisperModel

model = WhisperModel("large-v3", device="cpu", compute_type="int8")
segments, info = model.transcribe("audio.mp3")
# Performance: ~2-3x realtime on 8vCPU
```

**TTS (Piper)**:
```bash
docker run -p 10200:10200 rhasspy/wyoming-piper \
  --voice en_US-lessac-medium
# Performance: Realtime on CPU, natural voices
```

### Multi-Agent Orchestration: LangGraph (OSS)

**Why LangGraph**:
- ✅ FREE (MIT license)
- ✅ Stateful workflow graphs
- ✅ Multi-agent coordination
- ✅ Works with any LLM
- ✅ Self-hosted (no cloud)

**Architecture Pattern**:
```python
from langgraph.graph import StateGraph

workflow = StateGraph(AgentState)

# Add specialized agents
workflow.add_node("voice_agent", voice_handler)
workflow.add_node("vision_agent", vision_handler)
workflow.add_node("automation_agent", automation_handler)

# Route based on task
workflow.add_conditional_edges(
    "voice_agent",
    route_to_agent,  # Decides next agent
    {
        "vision": "vision_agent",
        "automation": "automation_agent",
        "end": END
    }
)

# Persist to Postgres
checkpointer = PostgresSaver.from_conn_string(
    "postgresql://supabase_connection"
)

app = workflow.compile(checkpointer=checkpointer)
```

### Vision & Face Recognition: Frigate + CompreFace

**Frigate NVR**:
- Real-time object detection (100+ FPS)
- MQTT event streaming
- Zone-based tracking
- 24/7 recording

**CompreFace + Double Take**:
- FREE face recognition
- NO cloud dependencies
- REST API + MQTT integration
- Home Assistant automation

**Integration**:
```
Camera Stream → Frigate (detect persons) → Double Take
→ CompreFace (recognize faces) → MQTT event → Home Assistant
→ Legacy Agent ("Zeus arrived home")
```

## Phase 2 Implementation Roadmap

### Week 1-2: Foundation
- [ ] Deploy LiveKit Cloud (request account)
- [x] Set up Letta (self-hosted) with PostgreSQL
- [ ] Install Ollama + Llama 3.3 70B on VM
- [ ] Create basic LiveKit Python agent
- [x] Integrate Letta memory
- [ ] Test voice pipeline

**Deliverable**: Voice agent with memory responding to "Hey Legacy"

### Week 3-4: Vision
- [ ] Deploy Frigate NVR
- [ ] Set up CompreFace + Double Take
- [ ] Configure camera streams
- [ ] Train face recognition
- [ ] Integrate into LiveKit agent

**Deliverable**: Agent recognizes faces and sees environment

### Week 5-6: Offline
- [ ] Deploy Faster-Whisper STT
- [ ] Deploy Piper TTS
- [ ] Configure Ollama as primary LLM
- [ ] Test full offline pipeline

**Deliverable**: Fully functional offline agent

### Week 7-8: Automation
- [ ] Create MCP servers for n8n, Node-RED
- [ ] Build Home Assistant tools
- [ ] Set up EMQX subscriptions
- [ ] Create automation workflows

**Deliverable**: Agent can control homestead

### Week 9-12: Intelligence
- [ ] Deploy LangGraph for orchestration
- [ ] Create specialized agents
- [ ] Implement sleep-time processing
- [ ] Build personality system

**Deliverable**: Intelligent, proactive multi-agent system

## Phase 2 Cost Analysis

| Component | Decision | Cost | Reasoning |
|-----------|----------|------|-----------|
| **Letta** | 🏠 Self-hosted | $0 | Open-source, agent framework |
| **LiveKit** | ☁️ Cloud | $50/mo | Global edge network worth it |
| **Ollama** | 🏠 Self-hosted | $0 | Privacy + offline capability |
| **Voice** | 🔄 Hybrid | $0-20 | Local fallback + optional cloud |
| **Vision** | 🏠 Self-hosted | $0 | Frigate + CompreFace free |
| **LangGraph** | 🏠 Self-hosted | $0 | Open-source orchestration |
| **TOTAL** | Hybrid | ~$75/mo | 97% cheaper than full cloud |

---

# PHASE 3: ADVANCED INTELLIGENCE (Q2 2026)

## Vision Extensions

### Sleep-Time Processing (Optional: Letta)

**When to add**: After Phase 2 core works

**Letta MemGPT Architecture**:
- Self-editing memory (agents manage their own knowledge)
- Background memory organization while offline
- Multi-thread agents (conscious/subconscious)
- Agent-as-a-Service REST API

**Funding**: $10M seed (well-funded, stable)

### Proactive Behavior System

```python
class ProactiveTriggers:
    schedule_check = CronJob("*/30 * * * *")  # Every 30 min
    weather_update = CronJob("0 6 * * *")     # Daily 6am
    health_check = CronJob("0 */4 * * *")     # Every 4 hours

# Examples:
"It's raining, should I close the greenhouse vents?"
"You have a meeting in 30 minutes, want me to prepare?"
"Unknown vehicle detected at the gate"
"Water tank at 30%, schedule refill?"
```

### Cross-Agent Memory Sharing

```python
# Farm Operations Agent stores info
client.add_memory(
    session_id="farm_ops",
    data={"zone_3_status": "watered", "timestamp": "2025-10-22T14:00"}
)

# Scheduling Agent retrieves it
farm_context = client.search(
    query="zone 3 water status",
    session_id="farm_ops"
)
# Uses context for scheduling
```

---

# DEVELOPMENT ENVIRONMENT

## SSH Access

```bash
# Connect to VM
ssh -i "/path/to/LightsailDefaultKey-us-east-1 (1).pem" ubuntu@3.215.29.128

# Verify connection
docker ps --format 'table {{.Names}}\t{{.Status}}'
```

## Local Development Workflow

```bash
# 1. Clone repo
git clone https://github.com/Legacyx-ai/Legacy-ai-Hybrid.git
cd Legacy-ai-Hybrid

# 2. Create feature branch
git checkout -b feature/your-feature-name

# 3. Make changes
# Edit files in legacy-ai-stacks/

# 4. Test locally
docker-compose up -d
curl http://localhost:8080/api/health

# 5. Run tests
pytest tests/unit/
pytest tests/integration/

# 6. Commit
git add .
git commit -m "feat: description of change"

# 7. Push & create PR
git push origin feature/your-feature-name
# Create PR on GitHub

# 8. After merge, deploy to VM
# (Automated or manual via ansible)
```

## Branch Naming Convention

```
feature/add-livekit-integration
feature/improve-memory-search
fix/handle-concurrent-requests
docs/update-api-documentation
chore/upgrade-dependencies
```

---

# API INTEGRATIONS

## Dify API

```python
import requests

API_KEY = "your-dify-api-key"
BASE_URL = "https://dify.zeusops.dev/api/v1"

response = requests.post(
    f"{BASE_URL}/chat-messages",
    headers={
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json"
    },
    json={
        "inputs": {},
        "query": "What is the weather?",
        "response_mode": "streaming",
        "user": "user123"
    }
)

for line in response.iter_lines():
    if line:
        print(line.decode('utf-8'))
```

## n8n API

```bash
# Get API key from https://n8n.zeusops.dev/
# Settings → API Keys

# Trigger workflow
curl -X POST https://n8n.zeusops.dev/webhook/my-webhook \
  -H "Content-Type: application/json" \
  -d '{"data": "value"}'

# Get execution
curl https://n8n.zeusops.dev/api/v1/executions/123 \
  -H "X-N8N-API-KEY: your-api-key"
```

## Letta API

```python
from letta import LocalClient

client = LocalClient()

agent = client.create_agent(
    name="my_agent",
    agent_type="default"
)

response = client.send_message(
    agent_id=agent.agent_id,
    message="Hello, how are you?"
)

print(response.message)
```

## MQTT (EMQX)

```python
import paho.mqtt.client as mqtt

client = mqtt.Client()
client.connect("emqx.zeusops.dev", 1883, 60)

# Subscribe
client.subscribe("home/temperature")

# Publish
client.publish("home/light", "on")

def on_message(client, userdata, msg):
    print(msg.payload.decode())

client.on_message = on_message
client.loop_forever()
```

---

# DEPLOYMENT & OPERATIONS

## Pre-Deployment Checklist

```
✅ All tests passing
✅ Code reviewed
✅ Documentation updated
✅ Environment variables configured
✅ Backup created (automatic)
✅ Security scan passed
```

## Manual Deployment

```bash
# 1. SSH to VM
ssh -i "key.pem" ubuntu@3.215.29.128

# 2. Pull latest
cd /opt/legacy-ai/stacks/
git pull origin main

# 3. Restart service
docker-compose down
docker-compose up -d --pull always

# 4. Verify
docker-compose ps
curl http://localhost:8080/api/health

# 5. Monitor
docker-compose logs -f
```

## Rollback Procedure

```bash
# 1. Identify previous commit
git log --oneline

# 2. Checkout previous version
git checkout <commit-hash>

# 3. Restart
docker-compose up -d --pull always

# 4. Verify
docker ps
curl http://localhost:8080/api/health

# 5. Investigate issue
docker-compose logs
```

## Zero-Downtime Deployment

```bash
# 1. Create new instance
docker-compose up -d service-new

# 2. Verify health
curl http://localhost:8080/health

# 3. Switch traffic (in Caddy)
# Update caddyfile routing to service-new

# 4. Remove old
docker-compose rm service-old

# 5. Restart proxy
docker-compose restart caddy
```

## Testing Strategy

### Unit Tests
```bash
pytest tests/unit/ -v
pytest tests/unit/ --cov=src --cov-report=html
```

### Integration Tests
```bash
pytest tests/integration/ -v
```

### Load Testing
```bash
ab -n 1000 -c 10 http://localhost:8080/api/health
```

### Manual Checklist

```
Dify:
  [ ] Can create workflow
  [ ] Can run workflow
  [ ] API responses valid
  [ ] Error handling works

n8n:
  [ ] Can create workflow
  [ ] Can execute workflow
  [ ] Webhook triggering works

Letta:
  [ ] Can create agent
  [ ] Agent responds to messages
  [ ] Memory persists

Node-RED:
  [ ] Can create flow
  [ ] MQTT message routing works
```

---

# TROUBLESHOOTING

## Service Won't Start

```bash
# Check logs
docker-compose logs service-name

# Check config
docker-compose config

# Validate environment
env | grep SERVICE

# Restart
docker-compose restart service-name
```

## High Memory Usage

```bash
# Check memory
docker stats

# Set limits in docker-compose.yml
services:
  service:
    mem_limit: 2g
    memswap_limit: 2g

# Restart
docker-compose down
docker-compose up -d
```

## Database Connection Errors

```bash
# Verify postgres running
docker-compose ps | grep postgres

# Check logs
docker-compose logs postgres

# Test connection
docker-compose exec postgres psql -U postgres -c "\l"
```

## HTTPS Certificate Issues

```bash
# Verify certificate
curl -vI https://dify.zeusops.dev

# Check Caddy
docker-compose logs caddy

# Force renewal
docker-compose exec caddy caddy reload

# Check files
ls -la /etc/caddy/certs/
```

---

# RESOURCES & REFERENCES

## Essential Documentation
- **[Architecture Details](./architecture/ARCHITECTURE.md)** - System design
- **[Current Services](./architecture/CURRENT-SERVICES.md)** - What's running
- **[Workflows](./architecture/WORKFLOWS.md)** - Common tasks
- **[Backup & Recovery](00-Infrastructure/Guides/BACKUP-AND-RECOVERY.md)** - Disaster recovery
- **[Maintenance Playbook](00-Infrastructure/Guides/MAINTENANCE-PLAYBOOK.md)** - Operations guide

## Service Documentation
- **Dify**: https://docs.dify.ai/
- **n8n**: https://docs.n8n.io/
- **Letta**: https://docs.letta.com/
- **LiveKit**: https://docs.livekit.io/
- **Node-RED**: https://nodered.org/docs/
- **Home Assistant**: https://www.home-assistant.io/docs/
- **EMQX**: https://www.emqx.io/docs/
- **Letta**: https://docs.letta.com/

## Technology Guides
- **Docker Compose**: https://docs.docker.com/compose/
- **Caddy**: https://caddyserver.com/docs/
- **PostgreSQL**: https://www.postgresql.org/docs/
- **Redis**: https://redis.io/documentation/
- **LangGraph**: https://langchain-ai.github.io/langgraph/
- **Ollama**: https://ollama.ai/

## GitHub
- **Repository**: https://github.com/Legacyx-ai/Legacy-ai-Hybrid
- **Issues**: GitHub Issues (feature tracking)
- **Discussions**: GitHub Discussions

---

# DEVELOPMENT PRIORITIES

## ✅ Current Phase: Infrastructure Hardening (COMPLETE)
- [x] System security patches applied (33.7 MB)
- [x] SSH hardened (Ed25519 + ChaCha20-Poly1305)
- [x] Backup automation operational
- [x] All 25 services operational
- [x] Comprehensive documentation created
- [x] Unified development master created

## 🎯 Next Phase: LiveKit Integration (Q1 2026)
- [ ] Deploy LiveKit Cloud account
- [x] Set up Letta memory layer
- [ ] Install Ollama + LLMs
- [ ] Create LiveKit Python agent
- [ ] Test voice pipeline
- [ ] Integrate memory layer
- [ ] Add wake word detection

## 🔜 Future: Vision & Multimodal (Q2 2026)
- [ ] Deploy Frigate NVR
- [ ] Set up CompreFace face recognition
- [ ] Integrate vision into agent
- [ ] Test full multimodal pipeline
- [ ] Add proactive behaviors
- [ ] Implement LangGraph orchestration
- [ ] Deploy edge instances

---

# QUICK REFERENCE

## Daily Commands

```bash
# SSH to VM
ssh -i "key.pem" ubuntu@3.215.29.128

# Check all services
docker ps --format 'table {{.Names}}\t{{.Status}}'

# View logs
docker-compose logs -f [service]

# Restart service
docker-compose restart [service]

# Deploy changes
docker-compose up -d --pull always

# Check backups
ls -la /backups/daily/
```

## Access Services

```
Portainer:      https://command.zeusops.dev
Letta:          https://letta.zeusops.dev
Dify:           https://dify.zeusops.dev
n8n:            https://n8n.zeusops.dev
Node-RED:       https://nodered.zeusops.dev
Home Assistant: https://ha.zeusops.dev
EMQX:           https://emqx.zeusops.dev
```

---

# SIGN-OFF

**Document Status**: ✅ APPROVED  
**Created**: October 23, 2025  
**Master Plan**: UNIFIED DEVELOPMENT REFERENCE  
**Scope**: Current (Phase 1) + Future Vision (Phases 2-3)  
**Maintained By**: Personal Development  

**This is your single source of truth for all Legacy AI development work.**

---

**Ready to start? Pick a task from Phase 2 or Phase 3 and begin development!** 🚀
