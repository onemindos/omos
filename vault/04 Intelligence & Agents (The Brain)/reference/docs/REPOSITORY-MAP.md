---
id: ui-02-20260111053925
title: "🗺️ Legacy AI - Complete Repository Map"
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
ai_summary: "🗺️ Legacy AI - Complete Repository Map"
---

# 🗺️ Legacy AI - Complete Repository Map

**Last Updated**: October 24, 2025  
**Status**: Phase 2A Production + Future Architecture Complete

---

## 📊 Top-Level Structure

```
Legacy-AI/
├── CORE/                  # 🧠 Intelligence System
├── CAPABILITIES/          # 🔧 Digital + Physical Actions
├── ENVIRONMENTS/          # 🌐 Deployment Targets
├── HYBRID-MODES/          # 🔄 Online ↔ Offline
├── OBSERVABILITY/         # 📊 Telemetry & Monitoring
├── COMMAND-CENTER/        # 🎛️ Control & Management
├── KNOWLEDGE/             # 📚 Documentation & Learning
└── PROJECT/               # 📋 Project Management
```

---

## 🧠 CORE - The Intelligence

```
CORE/
├── README.md                              # Complete cognitive architecture
├── docker-compose.yml                     # Main orchestration
├── .env.example                           # Environment template
│
├── brain/                                 # 🧠 REASONING & DECISIONS
│   └── letta/                             # Letta multi-agent system
│       ├── docker-compose.yml
│       ├── config/
│       │   ├── README.md                  # .af file format guide
│       │   └── legacy-ai-main.af          # Main agent config
│       ├── agents/                        # Multi-agent configs
│       │   ├── README.md                  # Agent system overview
│       │   ├── research-agent.af          # (Future) Research specialist
│       │   ├── coding-agent.af            # (Future) Code specialist
│       │   ├── data-analyst.af            # (Future) Data specialist
│       │   ├── creative-agent.af          # (Future) Creative specialist
│       │   ├── project-manager.af         # (Future) PM specialist
│       │   └── orchestrator.af            # (Future) Multi-agent coordinator
│       ├── tools-registry/                # Registered tool configs
│       └── shared-memory/                 # Inter-agent memory
│
├── senses/                                # 👁️👂👄 PERCEPTION
│   ├── livekit/                           # Real-time transport
│   │   ├── config/
│   │   └── agents/
│   ├── vision/                            # 👁️ EYES
│   │   ├── gpt-4o-vision/                 # Cloud vision (online)
│   │   ├── yolov8/                        # Local detection (offline)
│   │   ├── frigate/                       # 🎥 Video surveillance
│   │   │   ├── docker-compose.yml
│   │   │   ├── config.yml
│   │   │   └── models/
│   │   └── brilliant-labs/                # 👓 AR Glasses
│   │       ├── sdk/                       # Brilliant Labs SDK
│   │       ├── frame-config/              # Frame settings
│   │       ├── monocle-config/            # Monocle settings
│   │       └── stream-processing/         # Real-time AR
│   ├── audio/                             # 👂 EARS
│   │   ├── assemblyai/                    # Cloud STT (online)
│   │   ├── whisper/                       # Local STT (offline)
│   │   └── brilliant-labs-audio/          # AR audio input
│   └── voice/                             # 👄 MOUTH
│       ├── cartesia/                      # Cloud TTS (online)
│       ├── piper/                         # Local TTS (offline)
│       └── brilliant-labs-speaker/        # AR audio output
│
├── memory/                                # 💾 LONG-TERM STORAGE
│   ├── supabase-cloud/                    # Conversation logs + RLS
│   │   ├── schema.sql
│   │   └── migrations/
│   ├── postgresql/                        # Relational data
│   │   ├── init-scripts/
│   │   └── backups/
│   ├── redis/                             # Session cache
│   │   └── config/
│   └── vector-store/                      # Embeddings for RAG
│       ├── pinecone/                      # Cloud vector DB
│       └── qdrant/                        # Self-hosted vector DB
│
├── agent/                                 # 🤖 CENTRAL ORCHESTRATION
│   ├── agent.py                           # LiveKit agent entrypoint
│   ├── requirements.txt
│   └── tools/                             # Agent tool configs
│
└── scripts/                               # Core system scripts
    ├── init-legacy-ai.sh
    └── health-check.sh
```

---

## 🔧 CAPABILITIES - What Legacy AI Can Do

```
CAPABILITIES/
├── README.md                              # Capabilities overview
│
├── communication/                         # 📞📧 EXTERNAL COMMUNICATION
│   ├── voice-calling/                     # ☎️ PHONE CALLS
│   │   ├── twilio/                        # Twilio Voice API
│   │   │   ├── docker-compose.yml
│   │   │   ├── call-handlers/
│   │   │   └── voicemail/
│   │   ├── livekit-sip/                   # SIP integration
│   │   └── call-routing/                  # Call logic
│   ├── email/                             # 📧 EMAIL
│   │   ├── sendgrid/                      # SendGrid API
│   │   ├── imap-listener/                 # Receive emails
│   │   ├── smtp-sender/                   # Send emails
│   │   └── templates/                     # Email templates
│   └── messaging/                         # SMS, WhatsApp, Slack
│       ├── sms/
│       ├── whatsapp/
│       └── slack/
│
├── digital-actions/                       # 💻 DIGITAL WORLD
│   ├── workflow-automation/               # AI workflows
│   │   └── n8n/                          # ✅ PRODUCTION
│   │       ├── docker-compose.yml
│   │       └── workflows/
│   ├── task-automation/                   # General automation
│   │   └── n8n/                           # ✅ PRODUCTION
│   │       ├── docker-compose.yml
│   │       └── workflows/
│   ├── browser-automation/                # 🤖 WEB INTERACTION
│   │   └── bytebot/                       # ByteBot AI
│   │       ├── docker-compose.yml
│   │       ├── scripts/
│   │       └── sessions/
│   ├── smart-home/                        # Home automation
│   │   └── home-assistant/                # ✅ PRODUCTION
│   │       ├── configuration.yaml
│   │       └── automations/
│   ├── iot-messaging/                     # IoT communication
│   │   └── emqx/                          # ✅ PRODUCTION
│   │       ├── docker-compose.yml
│   │       └── mqtt-topics/
│   ├── visual-flows/                      # Flow programming
│   │   └── node-red/                      # ✅ PRODUCTION
│   │       ├── flows.json
│   │       └── custom-nodes/
│   ├── project-management/                # 📋 PLANNING
│   │   └── plane/                         # Plane.so self-hosted
│   │       ├── docker-compose.yml
│   │       ├── api-integration/
│   │       └── automation-workflows/
│   └── data-pipelines/                    # 📊 DATA INTEGRATION
│       ├── airbyte/                       # Data replication
│       │   ├── docker-compose.yml
│       │   ├── sources/
│       │   └── destinations/
│       └── meltano/                       # ELT framework
│           ├── meltano.yml
│           ├── extractors/
│           └── loaders/
│
├── physical-actions/                      # 🤖 PHYSICAL WORLD
│   ├── humanoid-robotics/                 # 🦾 HUMANOID ROBOTS
│   │   └── unitree/                       # Unitree SDK
│   │       ├── sdk/                       # unitree_sdk2
│   │       ├── robot-configs/
│   │       │   ├── g1/                    # G1 Humanoid
│   │       │   ├── h1/                    # H1 Humanoid
│   │       │   └── go2/                   # Go2 Quadruped
│   │       └── control-interfaces/
│   ├── ros2-integration/                  # 🔧 ROS2
│   │   ├── workspaces/
│   │   ├── packages/
│   │   │   ├── legacy_ai_control/
│   │   │   ├── legacy_ai_vision/
│   │   │   └── legacy_ai_navigation/
│   │   ├── launch-files/
│   │   └── docker-compose.yml
│   ├── agricultural-robotics/             # 🌱 FARMBOT
│   │   └── farmbot/
│   │       ├── farmbot-config.json
│   │       ├── sequences/
│   │       ├── plant-database/
│   │       └── api-integration/
│   ├── motion-control/
│   │   ├── arm-control/
│   │   └── locomotion/
│   ├── manipulation/
│   │   ├── object-grasping/
│   │   └── trajectory-planning/
│   ├── navigation/
│   │   ├── path-planning/
│   │   ├── localization/
│   │   └── slam/
│   ├── perception/
│   │   ├── computer-vision/
│   │   ├── lidar-processing/
│   │   └── sensor-fusion/
│   └── hardware-abstraction/
│       ├── motor-drivers/
│       └── sensor-drivers/
│
├── mcp-servers/                           # 🔌 MODEL CONTEXT PROTOCOL
│   ├── README.md                          # MCP configuration
│   ├── github/
│   ├── aws/
│   ├── filesystem/
│   └── ssh/
│
├── custom-skills/                         # 🎯 LEGACY AI SPECIFIC
│   ├── conversation/
│   ├── task-planning/
│   ├── learning/
│   └── decision-making/
│
└── api-integrations/                      # 🌐 EXTERNAL APIS
    ├── openai/
    ├── anthropic/
    ├── livekit/
    └── custom-apis/                       # 🛠️ CUSTOM
        ├── tool-data-api/
        ├── agent-control-api/
        └── integration-api/
```

---

## 🌐 ENVIRONMENTS - Where Legacy AI Runs

```
ENVIRONMENTS/
├── README.md                              # Environments overview
│
├── production/                            # ☁️ AWS LIGHTSAIL VM
│   ├── deploy-legacy-ai.sh                # ✅ Main deployment
│   ├── deploy-letta.sh                    # ✅ Brain deployment
│   ├── deploy-n8n.sh                      # ✅ n8n deployment
│   ├── deploy-iot.sh
│   ├── deploy-obsidian.sh                 # Knowledge base
│   ├── deploy-plane.sh                    # Project management
│   ├── deploy-airbyte.sh                  # Data pipelines
│   ├── deploy-frigate.sh                  # Video surveillance
│   ├── setup-caddy.sh                     # ✅ Reverse proxy
│   └── infrastructure/
│       ├── portainer/
│       ├── caddy/
│       └── backups/
│           ├── backup-all.sh
│           ├── backup-databases.sh
│           ├── setup-cron-backups.sh
│           └── setup-cron.sh
│
├── edge/                                  # 🔌 OFFLINE/EDGE
│   ├── jetson-nano/                       # Jetson Nano 4GB
│   │   ├── models/
│   │   │   ├── llm/                       # Local language models
│   │   │   ├── vision/                    # Local vision models
│   │   │   ├── audio/                     # Local STT models
│   │   │   └── tts/                       # Local TTS models
│   │   ├── local-services/
│   │   └── ros2-edge/                     # ROS2 on edge
│   ├── unitree-robots/                    # Robot hardware
│   │   ├── g1-deployment/
│   │   ├── h1-deployment/
│   │   └── go2-deployment/
│   ├── farmbot-hardware/                  # Farmbot configs
│   ├── brilliant-labs-devices/            # AR glasses
│   ├── docker-config/                     # Edge containers
│   └── sync-strategies/                   # Cloud ↔ Edge sync
│
├── development/                           # 🛠️ LOCAL DEV
│   ├── docker-compose.dev.yml
│   ├── .env.dev
│   └── setup-dev.sh
│
├── staging/                               # 🧪 TESTING
│   └── docker-compose.staging.yml
│
└── interfaces/                            # 👤 USER-FACING
    ├── web-app/
    │   ├── copilotkit/                    # 🤖 AI UI Platform
    │   │   ├── docker-compose.yml
    │   │   ├── copilotkit-config/
    │   │   └── components/
    │   ├── lobe-chat/                     # 💬 Chat UI
    │   │   ├── docker-compose.yml
    │   │   └── config/
    │   ├── open-webui/                    # 🌐 Alternative UI
    │   │   ├── docker-compose.yml
    │   │   └── models-config/
    │   └── dashboard/
    ├── mobile-app/
    │   ├── ios/
    │   └── android/
    ├── voice-interface/
    │   └── phone-menu/                    # IVR system
    ├── ar-interface/
    │   └── brilliant-labs-apps/           # AR apps
    └── api-server/
        ├── rest-api/
        └── graphql-api/
```

---

## 🔄 HYBRID-MODES - Online ↔ Offline

```
HYBRID-MODES/
├── README.md                              # Hybrid modes overview
│
├── local-models/                          # Offline inference
│   ├── llm/                               # Llama, Mistral (quantized)
│   ├── vision/                            # YOLO, MobileNet
│   ├── audio/                             # Whisper.cpp
│   └── tts/                               # Piper TTS
│
├── sync-protocols/                        # Cloud ↔ Edge sync
│   ├── powersync/                         # PowerSync config
│   └── conflict-resolution/               # CRDT-based
│
└── fallback-strategies/                   # Graceful degradation
    ├── offline-mode.md
    └── quality-tradeoffs.md
```

---

## 📊 OBSERVABILITY - Telemetry & Monitoring

```
OBSERVABILITY/
├── README.md                              # Observability overview
│
├── telemetry/                             # Data collection
│   ├── opentelemetry-config/
│   ├── traces/
│   ├── metrics/
│   └── logs/
│
├── monitoring/                            # Health checks
│   ├── prometheus/
│   ├── grafana/
│   ├── alertmanager/
│   ├── systems-check.sh                   # ✅ System health
│   ├── quick-fix.sh
│   ├── send-alert.sh
│   └── robot-health-check.sh              # Physical hardware
│
└── dashboards/                            # Visualization
    ├── grafana-dashboards/
    │   ├── system-overview.json
    │   ├── letta-performance.json
    │   ├── robot-telemetry.json
    │   └── data-pipeline-health.json
    ├── performance/
    └── usage-analytics/
```

---

## 🎛️ COMMAND-CENTER - Control & Management

```
COMMAND-CENTER/
├── README.md                              # Command center overview
│
├── admin-dashboards/                      # System management UIs
│   ├── dashy/                             # 🎯 Service dashboard
│   │   ├── docker-compose.yml
│   │   └── config.yml
│   ├── portainer/                         # 🐳 Container management
│   │   └── stacks/
│   └── uptime-kuma/                       # ⏱️ Uptime monitoring
│       └── docker-compose.yml
│
├── control-apis/                          # Management APIs
│   ├── system-control/
│   ├── agent-management/
│   ├── robot-control/
│   └── deployment-automation/
│
└── development-tools/                     # Dev environment
    ├── vscode-configs/
    │   ├── settings.json
    │   └── extensions.json
    ├── github-copilot-settings/
    └── debugging-tools/
```

---

## 📚 KNOWLEDGE - Documentation & Learning

```
KNOWLEDGE/
├── README.md                              # Knowledge overview
│
├── obsidian/                              # 🗂️ KNOWLEDGE BASE
│   ├── docker-compose.yml                 # Self-hosted Obsidian
│   ├── vault/                             # Markdown vault
│   │   ├── daily-notes/
│   │   ├── agent-learnings/               # What Legacy AI learns
│   │   ├── project-docs/
│   │   ├── api-references/
│   │   └── troubleshooting/
│   ├── plugins/
│   └── obsidian-sync/
│
└── docs/                                  # Technical docs (from /docs)
    ├── architecture/                      # System design
    │   ├── LEGACY-AI-BLUEPRINT.md
    │   ├── WORKFLOWS.md
    │   ├── DEPLOYMENT-STATUS.md
    │   └── ARCHITECTURE.md
    ├── guides/                            # How-to
    │   ├── AI-INTEGRATION-GUIDE.md
    │   ├── HTTPS-SETUP-GUIDE.md
    │   ├── OPERATIONS.md
    │   ├── TROUBLESHOOTING.md
    │   ├── BACKUP-SETUP.md
    │   └── LETTA-CONNECTION-INFO.md
    ├── current-state/                     # System inventory
    │   ├── INFRASTRUCTURE-INVENTORY.md
    │   ├── DATABASES-CURRENT.md
    │   ├── INTEGRATIONS-CURRENT.md
    │   ├── AWS-VM-CURRENT.md
    │   ├── PORTAINER-CURRENT.md
    │   ├── DOMAINS-ACCESS-CURRENT.md
    │   └── QUICK-REFERENCE.md
    └── archive/                           # Historical
        └── reorganization-2025-10-24/
```

---

## 📋 PROJECT - Project Management

```
PROJECT/
├── README.md                              # Main overview (from root)
├── REPOSITORY-GUIDE.md                    # Navigation guide
├── ARCHITECTURE-COMPLETE-MAP.md           # Visual maps
├── QUICK-START.md                         # Getting started
├── SETUP-CHECKLIST.md                     # Setup steps
├── API-KEYS-NEEDED.md                     # Credentials
├── DEPLOYMENT-READY.md                    # Deployment status
└── STATUS-REPORT.md                       # Progress tracking
```

---

## 🎯 Repository Philosophy

```
Legacy AI = Complete Intelligent Being

CORE             = The mind (perception, cognition, memory)
CAPABILITIES     = What she can do (digital + physical + communication)
ENVIRONMENTS     = Where she lives (cloud, edge, web, mobile, AR)
HYBRID-MODES     = How she adapts (online/offline flexibility)
OBSERVABILITY    = How we see her thinking (telemetry, monitoring)
COMMAND-CENTER   = How we manage her systems (dashboards, APIs, tools)
KNOWLEDGE        = What she learns (Obsidian + documentation)
PROJECT          = How we build her (project management)
```

---

## ✅ Current Implementation Status

### ✅ **Production (Phase 2A)**
- CORE/brain/letta/ (Letta reasoning)
- CORE/senses/livekit/ (LiveKit transport)
- CORE/memory/ (Supabase, PostgreSQL, Redis)
- CAPABILITIES/digital-actions/ (n8n, HA, Node-RED, EMQX)
- ENVIRONMENTS/production/ (AWS Lightsail VM)

### 🚧 **In Progress (Phase 2B)**
- CORE/senses/brilliant-labs/ (AR glasses)
- CORE/senses/vision/frigate/ (Video surveillance)
- ENVIRONMENTS/edge/jetson-nano/ (Edge deployment)
- HYBRID-MODES/local-models/ (Offline models)
- CAPABILITIES/digital-actions/browser-automation/ (ByteBot)
- CAPABILITIES/digital-actions/project-management/ (Plane)
- CAPABILITIES/digital-actions/data-pipelines/ (Airbyte, Meltano)

### 📋 **Planned (Phase 3 & Robotics)**
- CORE/brain/letta/agents/ (Multi-agent system)
- CAPABILITIES/communication/ (Phone, email, messaging)
- CAPABILITIES/physical-actions/ (Unitree, ROS2, Farmbot)
- ENVIRONMENTS/interfaces/ (CopilotKit, Lobe Chat, Open WebUI, mobile)
- OBSERVABILITY/ (Prometheus, Grafana, OpenTelemetry)
- COMMAND-CENTER/ (Dashy, Portainer, Uptime Kuma)
- KNOWLEDGE/obsidian/ (Self-hosted knowledge base)

---

**This is the complete map of Legacy AI - a cyborg-hybrid intelligent being with digital and physical capabilities, operating across cloud and edge environments, with full observability and control.** 🚀


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
