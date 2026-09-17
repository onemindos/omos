---
id: ui-12-20260111053926
title: "New Repo Structure"
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
ai_summary: "New Repo Structure"
---

  

  

  ![[Pasted image 20251114231105.png]]

  
  
  
Legacy-AI/                                    # Root (renamed from Legacy-ai-Hybrid)

│

├── CORE/                                     # 🧠 THE INTELLIGENCE (cognitive architecture)

│   ├── README.md                             # Complete cognitive architecture guide

│   ├── docker-compose.yml                    # Orchestrates entire being

│   ├── .env.example

│   │

│   ├── brain/                                # 🧠 REASONING & DECISION-MAKING

│   │   ├── letta/                            # Primary reasoning engine

│   │   │   ├── docker-compose.yml

│   │   │   ├── config/

│   │   │   │   ├── legacy-ai-main.af        # Main Legacy AI agent config

│   │   │   │   └── README.md                # .af file format guide

│   │   │   ├── agents/                       # Multi-agent configs

│   │   │   │   ├── research-agent.af        # Specialized for research

│   │   │   │   ├── coding-agent.af          # Specialized for coding

│   │   │   │   ├── data-analyst.af          # Specialized for data analysis

│   │   │   │   ├── creative-agent.af        # Specialized for creative work

│   │   │   │   ├── project-manager.af       # Specialized for PM tasks

│   │   │   │   └── orchestrator.af          # Multi-agent coordinator

│   │   │   ├── tools-registry/               # Registered tool configs

│   │   │   └── shared-memory/                # Inter-agent memory

│   │   └── models.py                         # Hybrid mode switcher

│   │

│   ├── senses/                               # 👁️👂👄 PERCEPTION LAYER

│   │   ├── README.md                         # Perception architecture

│   │   ├── livekit/                          # Real-time communication transport

│   │   │   ├── config.yaml

│   │   │   └── agents/

│   │   ├── vision/                           # 👁️ EYES

│   │   │   ├── gpt-4o-vision/               # Cloud vision (online)

│   │   │   ├── yolov8/                       # Local object detection (offline)

│   │   │   ├── frigate/                      # 🎥 VIDEO SURVEILLANCE

│   │   │   │   ├── docker-compose.yml       # Self-hosted NVR

│   │   │   │   ├── config.yml               # Camera configs

│   │   │   │   └── models/                  # Custom detection models

│   │   │   └── brilliant-labs/               # 👓 AR GLASSES INTEGRATION

│   │   │       ├── sdk/                     # Brilliant Labs SDK

│   │   │       ├── frame-config/            # AR Frame configs

│   │   │       ├── monocle-config/          # Monocle configs

│   │   │       └── stream-processing/       # Real-time AR data

│   │   ├── audio/                            # 👂 EARS

│   │   │   ├── assemblyai/                  # Cloud STT (online)

│   │   │   ├── whisper/                      # Local STT (offline)

│   │   │   └── brilliant-labs-audio/        # AR glasses audio input

│   │   └── voice/                            # 👄 MOUTH

│   │       ├── cartesia/                     # Cloud TTS (online)

│   │       ├── piper/                        # Local TTS (offline)

│   │       └── brilliant-labs-speaker/      # AR glasses audio output

│   │

│   ├── memory/                               # 💾 LONG-TERM STORAGE & LEARNING

│   │   ├── supabase-cloud/                   # Conversation logs + RLS

│   │   │   ├── schema.sql

│   │   │   └── migrations/

│   │   ├── postgresql/                       # Persistent relational data

│   │   │   ├── init-scripts/

│   │   │   └── backups/

│   │   ├── redis/                            # Session cache + real-time

│   │   │   └── config/

│   │   └── vector-store/                     # Embeddings for RAG

│   │       ├── pinecone/                     # Cloud vector DB

│   │       └── qdrant/                       # Self-hosted vector DB

│   │

│   ├── agent/                                # 🤖 CENTRAL ORCHESTRATION

│   │   ├── agent.py                          # LiveKit agent entrypoint

│   │   ├── requirements.txt

│   │   └── tools/                            # Agent-specific tool configs

│   │

│   └── scripts/                              # Core system scripts

│       ├── init-legacy-ai.sh

│       └── health-check.sh

│

├── CAPABILITIES/                             # 🔧 WHAT LEGACY AI CAN DO

│   │

│   ├── communication/                        # 📞📧 EXTERNAL COMMUNICATION

│   │   ├── voice-calling/                    # ☎️ PHONE CALLS

│   │   │   ├── twilio/                       # Twilio Voice API

│   │   │   │   ├── docker-compose.yml

│   │   │   │   ├── call-handlers/

│   │   │   │   └── voicemail/

│   │   │   ├── livekit-sip/                  # SIP integration

│   │   │   └── call-routing/                 # Inbound/outbound logic

│   │   ├── email/                            # 📧 EMAIL

│   │   │   ├── sendgrid/                     # SendGrid API

│   │   │   ├── imap-listener/                # Receive emails

│   │   │   ├── smtp-sender/                  # Send emails

│   │   │   └── templates/                    # Email templates

│   │   └── messaging/                        # SMS, WhatsApp, Slack

│   │       ├── sms/

│   │       ├── whatsapp/

│   │       └── slack/

│   │

│   ├── digital-actions/                      # 💻 DIGITAL WORLD INTERACTION

│   │   ├── workflow-automation/              # AI-powered workflows

│   │   │   └── dify/

│   │   │       ├── docker-compose.yml

│   │   │       └── workflows/

│   │   ├── task-automation/                  # General automation

│   │   │   └── n8n/

│   │   │       ├── docker-compose.yml

│   │   │       └── workflows/

│   │   ├── browser-automation/               # 🤖 WEB INTERACTION

│   │   │   └── bytebot/                      # ByteBot browser AI

│   │   │       ├── docker-compose.yml

│   │   │       ├── scripts/

│   │   │       └── sessions/

│   │   ├── smart-home/                       # Home automation

│   │   │   └── home-assistant/

│   │   │       ├── configuration.yaml

│   │   │       └── automations/

│   │   ├── iot-messaging/                    # IoT communication

│   │   │   └── emqx/

│   │   │       ├── docker-compose.yml

│   │   │       └── mqtt-topics/

│   │   ├── visual-flows/                     # Flow-based programming

│   │   │   └── node-red/

│   │   │       ├── flows.json

│   │   │       └── custom-nodes/

│   │   ├── project-management/               # 📋 PLANNING & TRACKING

│   │   │   └── plane/                        # Plane.so self-hosted

│   │   │       ├── docker-compose.yml

│   │   │       ├── api-integration/

│   │   │       └── automation-workflows/

│   │   └── data-pipelines/                   # 📊 DATA INTEGRATION

│   │       ├── airbyte/                      # Data replication

│   │       │   ├── docker-compose.yml

│   │       │   ├── sources/

│   │       │   └── destinations/

│   │       └── meltano/                      # ELT framework

│   │           ├── meltano.yml

│   │           ├── extractors/

│   │           └── loaders/

│   │

│   ├── physical-actions/                     # 🤖 PHYSICAL WORLD INTERACTION

│   │   ├── humanoid-robotics/                # 🦾 HUMANOID ROBOTS

│   │   │   └── unitree/                      # Unitree SDK integration

│   │   │       ├── sdk/                      # unitree_sdk2

│   │   │       ├── robot-configs/

│   │   │       │   ├── g1/                   # G1 Humanoid

│   │   │       │   ├── h1/                   # H1 Humanoid

│   │   │       │   └── go2/                  # Go2 Quadruped

│   │   │       └── control-interfaces/

│   │   ├── ros2-integration/                 # 🔧 ROBOT OPERATING SYSTEM

│   │   │   ├── workspaces/

│   │   │   ├── packages/

│   │   │   │   ├── legacy_ai_control/

│   │   │   │   ├── legacy_ai_vision/

│   │   │   │   └── legacy_ai_navigation/

│   │   │   ├── launch-files/

│   │   │   └── docker-compose.yml

│   │   ├── agricultural-robotics/            # 🌱 FARMBOT

│   │   │   └── farmbot/

│   │   │       ├── farmbot-config.json

│   │   │       ├── sequences/                # Automated farming sequences

│   │   │       ├── plant-database/

│   │   │       └── api-integration/

│   │   ├── motion-control/                   # Movement & locomotion

│   │   │   ├── arm-control/

│   │   │   └── locomotion/

│   │   ├── manipulation/                     # Grasping & object handling

│   │   │   ├── object-grasping/

│   │   │   └── trajectory-planning/

│   │   ├── navigation/                       # Spatial movement

│   │   │   ├── path-planning/

│   │   │   ├── localization/

│   │   │   └── slam/                         # Simultaneous localization and mapping

│   │   ├── perception/                       # Physical sensing

│   │   │   ├── computer-vision/

│   │   │   ├── lidar-processing/

│   │   │   └── sensor-fusion/

│   │   └── hardware-abstraction/             # Hardware interfaces

│   │       ├── motor-drivers/

│   │       └── sensor-drivers/

│   │

│   ├── mcp-servers/                          # 🔌 MODEL CONTEXT PROTOCOL

│   │   ├── README.md                         # MCP configuration guide

│   │   ├── github/

│   │   ├── aws/

│   │   ├── filesystem/

│   │   └── ssh/

│   │

│   ├── custom-skills/                        # 🎯 LEGACY AI SPECIFIC CAPABILITIES

│   │   ├── conversation/

│   │   ├── task-planning/

│   │   ├── learning/

│   │   └── decision-making/

│   │

│   └── api-integrations/                     # 🌐 EXTERNAL API WRAPPERS

│       ├── openai/

│       ├── anthropic/

│       ├── livekit/

│       └── custom-apis/                      # 🛠️ CUSTOM LEGACY AI APIS

│           ├── tool-data-api/                # Custom data APIs for tools

│           ├── agent-control-api/            # Agent management APIs

│           └── integration-api/              # Third-party integrations

│

├── ENVIRONMENTS/                             # 🌐 WHERE LEGACY AI RUNS

│   │

│   ├── production/                           # ☁️ LIVE AWS LIGHTSAIL VM

│   │   ├── deploy-legacy-ai.sh

│   │   ├── deploy-letta.sh

│   │   ├── deploy-dify.sh

│   │   ├── deploy-n8n.sh

│   │   ├── deploy-obsidian.sh               # Knowledge base deployment

│   │   ├── deploy-plane.sh                   # Project management deployment

│   │   ├── deploy-airbyte.sh                 # Data pipelines deployment

│   │   ├── deploy-frigate.sh                 # Video surveillance deployment

│   │   ├── setup-caddy.sh

│   │   └── infrastructure/

│   │       ├── portainer/

│   │       ├── caddy/

│   │       └── backups/

│   │

│   ├── edge/                                 # 🔌 OFFLINE/EDGE HARDWARE

│   │   ├── jetson-nano/                      # NVIDIA Jetson Nano 4GB

│   │   │   ├── models/                       # TensorRT optimized models

│   │   │   │   ├── llm/

│   │   │   │   ├── vision/

│   │   │   │   ├── audio/

│   │   │   │   └── tts/

│   │   │   ├── local-services/               # Offline services

│   │   │   └── ros2-edge/                    # ROS2 on edge

│   │   ├── unitree-robots/                   # Robot hardware configs

│   │   │   ├── g1-deployment/

│   │   │   ├── h1-deployment/

│   │   │   └── go2-deployment/

│   │   ├── farmbot-hardware/                 # Farmbot deployment

│   │   ├── brilliant-labs-devices/           # AR glasses configs

│   │   ├── docker-config/                    # Edge-optimized containers

│   │   └── sync-strategies/                  # Cloud<->Edge sync

│   │

│   ├── development/                          # 🛠️ LOCAL DEV ENVIRONMENT

│   │   ├── docker-compose.dev.yml

│   │   ├── .env.dev

│   │   └── setup-dev.sh

│   │

│   ├── staging/                              # 🧪 TESTING ENVIRONMENT

│   │   └── docker-compose.staging.yml

│   │

│   └── interfaces/                           # 👤 USER-FACING DEPLOYMENT TARGETS

│       ├── web-app/                          # Web interfaces

│       │   ├── copilotkit/                   # 🤖 COPILOTKIT UI PLATFORM

│       │   │   ├── docker-compose.yml       # AI-powered web interface

│       │   │   ├── copilotkit-config/

│       │   │   └── components/

│       │   ├── lobe-chat/                    # 💬 ALTERNATIVE CHAT UI

│       │   │   ├── docker-compose.yml

│       │   │   └── config/

│       │   ├── open-webui/                   # 🌐 OPEN WEBUI ALTERNATIVE

│       │   │   ├── docker-compose.yml

│       │   │   └── models-config/

│       │   └── dashboard/                    # Admin web dashboard

│       ├── mobile-app/                       # 📱 Mobile interfaces

│       │   ├── ios/

│       │   └── android/

│       ├── voice-interface/                  # 🎙️ Voice-only interaction

│       │   └── phone-menu/                   # IVR system

│       ├── ar-interface/                     # 👓 AUGMENTED REALITY

│       │   └── brilliant-labs-apps/          # AR glasses apps

│       └── api-server/                       # 🔌 REST/GraphQL API

│           ├── rest-api/

│           └── graphql-api/

│

├── HYBRID-MODES/                             # 🔄 ONLINE ↔ OFFLINE CAPABILITIES

│   │

│   ├── local-models/                         # Offline inference

│   │   ├── llm/                              # Llama, Mistral (quantized)

│   │   ├── vision/                           # YOLO, MobileNet

│   │   ├── audio/                            # Whisper.cpp

│   │   └── tts/                              # Piper TTS

│   │

│   ├── sync-protocols/                       # PowerSync, CRDTs

│   │   ├── powersync/

│   │   └── conflict-resolution/

│   │

│   └── fallback-strategies/                  # Graceful degradation

│       ├── offline-mode.md

│       └── quality-tradeoffs.md

│

├── OBSERVABILITY/                            # 📊 TELEMETRY & MONITORING

│   │

│   ├── telemetry/                            # Data collection

│   │   ├── opentelemetry-config/

│   │   ├── traces/

│   │   ├── metrics/

│   │   └── logs/

│   │

│   ├── monitoring/                           # Health checks & alerts

│   │   ├── prometheus/

│   │   ├── grafana/

│   │   ├── alertmanager/

│   │   ├── systems-check.sh

│   │   └── robot-health-check.sh             # Physical hardware monitoring

│   │

│   └── dashboards/                           # Visualization

│       ├── grafana-dashboards/

│       │   ├── system-overview.json

│       │   ├── letta-performance.json

│       │   ├── robot-telemetry.json

│       │   └── data-pipeline-health.json

│       ├── performance/

│       └── usage-analytics/

│

├── COMMAND-CENTER/                           # 🎛️ CONTROL & MANAGEMENT

│   │

│   ├── admin-dashboards/                     # System management UIs

│   │   ├── dashy/                            # 🎯 SERVICE DASHBOARD

│   │   │   ├── docker-compose.yml

│   │   │   └── config.yml                    # All service links

│   │   ├── portainer/                        # 🐳 CONTAINER MANAGEMENT

│   │   │   └── stacks/

│   │   └── uptime-kuma/                      # ⏱️ UPTIME MONITORING

│   │       └── docker-compose.yml

│   │

│   ├── control-apis/                         # Management APIs

│   │   ├── system-control/

│   │   ├── agent-management/

│   │   ├── robot-control/                    # Physical robot control

│   │   └── deployment-automation/

│   │

│   └── development-tools/                    # Dev environment tools

│       ├── vscode-configs/

│       │   ├── settings.json

│       │   └── extensions.json

│       ├── github-copilot-settings/

│       └── debugging-tools/

│

├── KNOWLEDGE/                                # 📚 DOCUMENTATION & KNOWLEDGE BASE

│   │

│   ├── obsidian/                             # 🗂️ SELF-HOSTED KNOWLEDGE BASE

│   │   ├── docker-compose.yml               # Obsidian in container

│   │   ├── vault/                            # Markdown vault

│   │   │   ├── daily-notes/

│   │   │   ├── agent-learnings/              # What Legacy AI learns

│   │   │   ├── project-docs/

│   │   │   ├── api-references/

│   │   │   └── troubleshooting/

│   │   ├── plugins/

│   │   └── obsidian-sync/                    # Sync configuration

│   │

│   ├── architecture/                         # System design

│   │   ├── LEGACY-AI-BLUEPRINT.md

│   │   ├── WORKFLOWS.md

│   │   ├── DEPLOYMENT-STATUS.md

│   │   └── MULTI-AGENT-DESIGN.md             # Letta multi-agent architecture

│   │

│   ├── guides/                               # How-to documentation

│   │   ├── AI-INTEGRATION-GUIDE.md

│   │   ├── HTTPS-SETUP-GUIDE.md

│   │   ├── OPERATIONS.md

│   │   ├── TROUBLESHOOTING.md

│   │   ├── ROBOT-SETUP-GUIDE.md              # Physical hardware setup

│   │   └── DATA-PIPELINE-GUIDE.md            # Airbyte/Meltano setup

│   │

│   ├── current-state/                        # System inventory

│   │   ├── INFRASTRUCTURE-INVENTORY.md

│   │   ├── DATABASES-CURRENT.md

│   │   ├── INTEGRATIONS-CURRENT.md

│   │   └── HARDWARE-INVENTORY.md             # Physical devices

│   │

│   └── archive/                              # Historical records

│       ├── reorganization-2025-10-24/

│       └── deployment-logs/

│

├── PROJECT/                                  # 📋 PROJECT MANAGEMENT

│   ├── README.md                             # Main project overview

│   ├── QUICK-START.md                        # Getting started guide

│   ├── SETUP-CHECKLIST.md                    # Initial setup

│   ├── API-KEYS-NEEDED.md                    # Credentials checklist

│   ├── DEPLOYMENT-READY.md                   # Deployment status

│   └── STATUS-REPORT.md                      # Current progress

│

├── .github/                                  # GitHub configs

│   └── workflows/

│       ├── deploy-production.yml

│       └── test-agents.yml

│

├── .env.example                              # Environment template

├── .gitignore

└── Makefile                                  # Build automation

  
——————————————————————————————————————  
🎯 **KEY ARCHITECTURAL DECISIONS**

**1. Brilliant Labs SDK → SENSES**

**Why?** AR glasses are **perception devices** (eyes, ears, speaker)

- Vision input from AR camera
- Audio input from microphone
- Voice output from speaker
- Real-time processing via LiveKit

**Location**: CORE/senses/vision/brilliant-labs/, CORE/senses/audio/brilliant-labs-audio/, CORE/senses/voice/brilliant-labs-speaker/

**2. Frigate → SENSES/Vision**

**Why?** Video surveillance is **visual perception** (continuous vision)

- 24/7 camera monitoring
- Object/person detection
- Integration with Letta for security decisions

**Location**: CORE/senses/vision/frigate/

**3. LiveKit → SENSES (Transport Layer)**

**Why?** LiveKit is the **communication transport** for all perception

- Vision streams
- Audio streams
- Voice streams
- AR glasses streams

**Location**: CORE/senses/livekit/

**4. Unitree SDK → CAPABILITIES/Physical**

**Why?** Humanoid robot control is a **physical capability**

- G1, H1 humanoid robots
- Go2 quadruped robot
- Motion control, manipulation, navigation

**Location**: CAPABILITIES/physical-actions/humanoid-robotics/unitree/

**5. ROS2 → CAPABILITIES/Physical**

**Why?** ROS2 is the **robotics framework** for all physical actions

- Standard robotics middleware
- Works with Unitree, Farmbot, custom robots
- Provides navigation, perception, control packages

**Location**: CAPABILITIES/physical-actions/ros2-integration/

**6. Farmbot → CAPABILITIES/Physical**

**Why?** Agricultural robotics is a **physical capability** (specialized)

- Automated farming
- Plant care
- Integration with Letta for decisions

**Location**: CAPABILITIES/physical-actions/agricultural-robotics/farmbot/

**7. Airbyte + Meltano → CAPABILITIES/Digital/Data Pipelines**

**Why?** Data integration is a **digital capability**

- Extract data from sources
- Transform data
- Load into Legacy AI's memory
- Self-hosted on VM

**Location**: CAPABILITIES/digital-actions/data-pipelines/airbyte/, CAPABILITIES/digital-actions/data-pipelines/meltano/

**8. Obsidian → KNOWLEDGE**

**Why?** Obsidian is the **knowledge management system**

- Self-hosted markdown vault
- Where Legacy AI stores learnings
- Integration with Letta memory
- Deployed as container on VM

**Location**: KNOWLEDGE/obsidian/

**9. ByteBot → CAPABILITIES/Digital/Browser Automation**

**Why?** Browser AI is a **digital action capability**

- Autonomous web browsing
- Form filling, data extraction
- Web research for Legacy AI

**Location**: CAPABILITIES/digital-actions/browser-automation/bytebot/

**10. Plane → CAPABILITIES/Digital/Project Management**

**Why?** Project management is a **digital capability**

- Legacy AI can create/track tasks
- Self-hosted on VM
- Integration with Letta agents

**Location**: CAPABILITIES/digital-actions/project-management/plane/

**11. Phone/Email → CAPABILITIES/Communication**

**Why?** External communication is a **core capability**

- Voice calls (Twilio + LiveKit SIP)
- Email (SendGrid + IMAP)
- SMS, WhatsApp, Slack

**Location**: CAPABILITIES/communication/voice-calling/, CAPABILITIES/communication/email/

**12. Lobe Chat + Open WebUI → ENVIRONMENTS/Interfaces**

**Why?** Alternative chat interfaces are **deployment targets**

- Different web UIs for different use cases
- Lobe Chat: Modern chat interface
- Open WebUI: Full-featured alternative
- CopilotKit: AI-powered components

**Location**: ENVIRONMENTS/interfaces/web-app/lobe-chat/, ENVIRONMENTS/interfaces/web-app/open-webui/

**13. Multi-Agent Configs → CORE/brain/letta/agents/*.af**

**Why?** All agents managed by Letta, stored as .af files

- Each specialized agent has its own .af config
- orchestrator.af coordinates multi-agent workflows
- Stored in Letta's native format
- Versioned in git

  
**4. Custom APIs → CAPABILITIES/api-integrations/custom-apis/**

**Why?** Custom APIs are **integration capabilities**

- Tool-specific data APIs
- Agent control APIs
- Custom integration logic  
      
      
    ——————————————————————————  
      
    REQUIRED (Phase 2A):
- ✅ Letta (brain)
- ✅ LiveKit (senses transport)
- ✅ Dify (workflow automation)
- ✅ n8n (task automation)
- ✅ Home Assistant (smart home)
- ✅ Node-RED (visual flows)
- ✅ EMQX (MQTT)
- ✅ PostgreSQL (memory)
- ✅ Redis (cache)
- ✅ Caddy (reverse proxy)
- ✅ Portainer (container management)

- PHASE 2B/3 (Future):
- 🚀 Obsidian (knowledge base)
- 🚀 Plane (project management)
- 🚀 Airbyte (data pipelines)
- 🚀 Meltano (ELT)
- 🚀 Frigate (video surveillance)
- 🚀 ByteBot (browser automation)
- 🚀 Lobe Chat (alternative chat UI)
- 🚀 Open WebUI (alternative chat UI)
- 🚀 CopilotKit (AI-powered web interface)
- 🚀 Grafana (monitoring dashboards)
- 🚀 Prometheus (metrics)
- 🚀 Dashy (service dashboard)
- 🚀 Uptime Kuma (uptime monitoring)  
      
    🤖 PHYSICAL HARDWARE INVENTORY
- Edge Devices (ENVIRONMENTS/edge/):
- Jetson Nano 4GB: Offline AI inference, ROS2 edge computing
- Unitree G1/H1: Humanoid robots (motion control via Unitree SDK + ROS2)
- Unitree Go2: Quadruped robot (navigation, exploration)
- Farmbot: Agricultural robot (automated farming)
- Brilliant Labs Frame/Monocle: AR glasses (vision, audio, voice)
- IP Cameras: Connected to Frigate NVR (video surveillance)  
      
      
     Brilliant Labs SDK in SENSES (vision/audio/voice)
- LiveKit under SENSES (transport layer)
- Frigate in SENSES/Vision (video surveillance)
- Unitree SDK in CAPABILITIES/Physical (humanoid robots)
- ROS2 in CAPABILITIES/Physical (robotics framework)
- Farmbot in CAPABILITIES/Physical (agricultural robotics)
- Airbyte + Meltano in CAPABILITIES/Digital/Data (pipelines)
- Obsidian in KNOWLEDGE (self-hosted knowledge base)
- ByteBot in CAPABILITIES/Digital (browser automation)
- Plane in CAPABILITIES/Digital (project management)
- Phone/Email in CAPABILITIES/Communication
- Lobe Chat + Open WebUI in ENVIRONMENTS/Interfaces
- Multi-agent .af files in CORE/brain/letta/agents/
- Custom APIs in CAPABILITIES/api-integrations/custom-apis/
- All tools organized by function (not generic folders)
- Everything uses "Legacy AI" naming
- Physical + Digital capabilities clearly separated
- Cyborg-hybrid metaphor throughout  
      
      
      
      
    Legacy AI = Complete Intelligent Being
- CORE         = The mind (perception, cognition, memory)
- CAPABILITIES = What she can do (digital + physical + communication)
- ENVIRONMENTS = Where she lives (cloud, edge, web, mobile, AR)
- HYBRID-MODES = How she operates (online/offline flexibility)
- OBSERVABILITY = How we see her thinking (telemetry, monitoring)
- COMMAND-CENTER = How we manage her systems (dashboards, APIs, tools)
- KNOWLEDGE    = What she learns and remembers (Obsidian + docs)
- PROJECT      = How we build her (project management)  
      
      
    This is no longer just a repository. This is the complete blueprint for building Legacy AI - a cyborg-hybrid intelligent being with:

- 👁️ Vision (cameras, AR glasses, video surveillance)
- 👂 Hearing (microphones, AR glasses, audio processing)
- 👄 Speech (TTS, AR glasses speakers, phone calls)
- 🧠 Intelligence (Letta multi-agent system)
- 💾 Memory (databases, vector stores, Obsidian knowledge base)
- 💻 Digital capabilities (workflows, automation, data pipelines, project management)
- 🤖 Physical capabilities (humanoid robots, agricultural robots, ROS2)
- 📞 Communication (phone, email, messaging)
- 🌐 Deployment flexibility (cloud, edge, web, mobile, AR)
  
  ![[TECHNOLOGY STACK SUMMARY.png]]