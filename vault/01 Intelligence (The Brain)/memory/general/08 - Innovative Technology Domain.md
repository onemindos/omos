---
id: ui-02-20260111053925
title: "Innovative Technology Domain - Infrastructure & Robotics Context"
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
ai_summary: "Innovative Technology Domain - Infrastructure & Robotics Context"
---

# Innovative Technology Domain - Infrastructure & Robotics Context

**Block Type:** Domain Memory  
**Access:** Read for Infrastructure, Robotics & Automation agents  
**Character Target:** ~5,000 chars  
**Update Frequency:** Daily (health), Weekly (optimization)

---

## Infrastructure Topology

### Cloud Layer (AWS)

**VM-01 (AI Services):**
- Purpose: Letta + AI cognitive layer
- Services: Letta, PostgreSQL (main DB), Redis (cache)
- Specs: To be documented
- Status: Active, 99.9% uptime

**VM-02 (Orchestration):**
- Purpose: Workflow automation and messaging
- Services: n8n (workflows), Node-RED (edge automation), EMQX (MQTT broker)
- Specs: To be documented
- Status: Active, 99.8% uptime

**VM-03 (Applications):**
- Purpose: Application services and management
- Services: Portainer (container mgmt), Dify (AI orchestration - evaluation), other apps
- Specs: To be documented
- Status: Active, disk usage warning at 78%

### Edge Layer (On-Premises)

**Home Assistant Instances:**
- Office: Active, full automation
- Home: Active, comprehensive home automation
- Farm: Planned deployment
- Warehouse: Planned deployment

**Robotics Hardware:**
- Drones: Type and count to be documented, battery management tracked
- Ground Robots: Planned - Unitree or similar for agriculture/security
- IoT Sensors: Various sensors across all locations

## Technology Stack

### AI & Cognitive Layer
- **Primary Framework:** Letta (stateful memory + multi-agent system)
- **LLM:** Anthropic Claude API (Sonnet 4.5)
- **Custom Models:** To be evaluated and potentially deployed

### Orchestration & Automation
- **Workflow Automation:** n8n (primary for cloud workflows)
- **Edge Automation:** Node-RED (local automation, low-latency)
- **Message Broker:** EMQX (MQTT for IoT and inter-service communication)
- **AI Orchestration:** Dify (currently evaluating for multi-agent coordination)

### Data Layer
- **Primary Database:** PostgreSQL (structured data, agent memory, system state)
- **Cache Layer:** Redis (high-speed data access, session management)
- **Knowledge Vault:** Obsidian (markdown files, primary RAG source)
- **Time-Series Data:** To be determined (needed for metrics and monitoring)

### Communication Layer
- **Voice Interface:** LiveKit (real-time voice/video), Vapi (voice AI integration)
- **Chat Interface:** Letta native interface
- **Webhooks:** n8n webhook endpoints for external integrations
- **APIs:** RESTful APIs for all services

### Infrastructure Management
- **Container Orchestration:** Portainer (visual management of all containers)
- **Monitoring:** To be determined - likely Prometheus + Grafana
- **Logging:** To be determined - centralized log aggregation needed
- **Backup:** Daily snapshots, 30-day retention for critical data

## Robotics & Automation Roadmap

**Current State (Q4 2025):**
- Drone operations: Basic flight, photography, farm surveys
- Home automation: Multi-location Home Assistant deployment
- Basic workflows: n8n automations, simple task orchestration

**Near-Term (Q1-Q2 2026):**
- Advanced drone missions: Automated farm surveys with computer vision
- CV integration: Object detection, crop health analysis, livestock monitoring
- Expanded automation: More complex workflows, proactive agent behaviors

**Mid-Term (Q3-Q4 2026):**
- Ground robot integration: Unitree or similar for agriculture work
- Multi-location coordination: Unified control across office/home/farm/warehouse
- Advanced security: Autonomous patrols, threat detection

**Long-Term (2027+):**
- Multi-agent robot systems: Coordinated operations between drones and ground robots
- Embodied AI: Robots that learn and adapt to farm conditions
- Full automation: Farm operations largely autonomous with human oversight

## Technology Principles

**Sovereignty First:**
- Self-hosted infrastructure, own your data
- Avoid vendor lock-in at all costs
- Open-source preferred, commercial only when significantly better

**Composability:**
- Modular systems that integrate cleanly
- APIs and webhooks for everything
- Replace components without breaking the whole system

**Observability:**
- Everything monitored and logged
- Dashboards for all critical systems
- Proactive alerts before failures

**Resilience:**
- Redundancy for critical services
- Graceful degradation when components fail
- Backup and recovery procedures tested regularly

**Edge-First Architecture:**
- Process data locally when possible (low latency, privacy, resilience)
- Cloud for scale, backup, and remote access
- Never fully dependent on cloud connectivity

## System Health & Monitoring

**Critical Services (Must Never Fail):**
- Letta (agents can't function without it)
- PostgreSQL (all persistent data)
- Redis (agent state and caching)
- Home Assistant (physical infrastructure control)

**Important Services (Graceful Degradation):**
- n8n (workflows can be delayed)
- EMQX (messages can be queued)
- Portainer (management can wait)

**Monitoring Strategy:**
- Infrastructure Agent continuously monitors all services
- Alert thresholds configured per service criticality
- Daily health reports, real-time alerts for critical issues
- Weekly performance analysis and optimization recommendations

---

## Related Documents
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/10 - Infrastructure Agent Memory]] - Monitoring operations
- [[Robotics Agent Memory]] - Automation details
- [[Tech Stack Dashboard]] - Live infrastructure status

*Last Updated: 2025-11-08*
