---
domain: SO
subdomain: "02 Infrastructure (VPS, Servers & Networking) (The Body)"
type: index
status: active
updated: 2026-05-03
tags: [domain/so, zone/so, type/index, status/active]
ai_summary: "SO 02 — Full infrastructure layer. Nodes (cloud/on-prem/edge), networking (NATS/Tailscale/DNS), services (LiveKit/n8n/Letta/email/databases), operations (security/runbooks/reference/cicd/monitoring/backup/testing), integrations, platform catalog, scripts, storage, and OpenClaw platform."
agent: Grid
todoist_project_id: 6g6P9p83V2666MP3
zone: SO
review: 2026-Q3
---
# 02 Infrastructure (VPS, Servers & Networking) (The Body)

> *All servers, APIs, cloud services, storage, networking, and integrations.*

## Contents

### openclaw (OpenClaw Platform)/
OpenClaw runtime — gateway configs, workspaces, skills, plugins, deploy scripts.
Subfolders: `deploy (Deployment)`, `guides (Guides)`, `interfaces (Interfaces)`, `nats (NATS Config)`, `OpenClaw-Hosting-Costs (Hosting Costs)`, `openclaw-reference (OpenClaw Reference)`, `OpenClaw-Skills-Catalog (Skills Catalog)`, `OpenClaw-Usecases (Use Cases)`, `patches (Patches)`, `planning (Planning)`, `plugin-templates (Plugin Templates)`, `workspaces (Agent Workspaces)`.

### nodes (Server Nodes)/
| Subfolder | What Lives Here |
|-----------|----------------|
| `cloud (Cloud Nodes)/vps (VPS Configs)/` | VPS configs — main, sentinel, DO deployment, Docker services, media server |
| `cloud (Cloud Nodes)/offline (Offline Operations)/` | Offline operation resilience — services overview, quick reference |
| `cloud (Cloud Nodes)/planning (Cloud Planning)/` | Scaling, GPU analysis, global infrastructure, world grid, AWS quotas |
| `cloud (Cloud Nodes)/providers (Cloud Providers)/` | Provider placeholders (AWS, Azure, DO, GCP, Alibaba) |
| `on-prem (On-Premises Nodes)/home-assistant (Home Assistant)/` | HA integration, voice pipeline, smart home sensors |
| `on-prem (On-Premises Nodes)/media (Media Server)/` | OMOS TV architecture and cast |
| `on-prem (On-Premises Nodes)/iot (IoT Devices)/` | ESP32-CAM glasses parts list |
| `edge (Edge Nodes)/` | Edge compute requirements |

### network (Network Config)/
| Subfolder | What Lives Here |
|-----------|----------------|
| `nats (NATS Config)/` | NATS bus — protocol, events, subjects, websocket, multi-location, boundaries |
| `tailscale (Tailscale Config)/` | Tailscale mesh network config |
| `dns (DNS Config)/` | Cloudflare DNS export |
| `overview (Network Overview)/` | NETWORKING.md (overview) |

### services (Running Services)/
| Subfolder | What Lives Here |
|-----------|----------------|
| `livekit (LiveKit Service)/` | LiveKit self-hosted — setup, voice pipeline, voice-agent code |
| `email (Email Services)/` | Email infrastructure — Zoho config, routing, filters, categories |
| `n8n (N8N Service)/` | n8n automation integration and setup |
| `letta (Letta Service)/` | Letta AI server tool upload checklist |
| `databases (Databases)/` | Database init scripts |

### operations (Ops Procedures)/
| Subfolder | What Lives Here |
|-----------|----------------|
| `security (Security Ops)/` | Security, CORS, Cloudflare setup, fleet hardening, watchdogs, failsafe |
| `runbooks (Runbooks)/` | Activation checklist, quickstarts, migration, maintenance, cleanup |
| `reference (Reference Docs)/` | Infrastructure index, access reference, scaling, AI ops, legacy AI |
| `cicd (CI-CD Pipelines)/` | CI/CD pipelines, GitHub Actions, deployment automation |
| `monitoring (Monitoring)/` | Observability, Datadog, health checks, VM review |
| `backup (Backup Procedures)/` | Backup strategy, disaster recovery, procedures |
| `testing (Testing)/` | Testing guides, evals, guardrails, test results |

### integrations (App Connections)/
| Subfolder | What Lives Here |
|-----------|----------------|
| `messaging (Messaging)/` | Discord, Telegram integrations |
| `tools (Tool Configs)/` | GitHub, n8n, Todoist, Cloudflare, Home Assistant, Nextcloud configs |
| `apis (API Integrations)/` | Composio, atprotocol, Postman, Webflow API docs |
| `overview (Integration Overview)/` | Overview, bridge architecture, integrations index |

### platform-catalog (Platform Catalog)/
| Subfolder | What Lives Here |
|-----------|----------------|
| `platforms (Platform Categories)/` | 19 category subfolders (AI agents, automation, business, cloud, etc.) |
| `data (Platform Data)/` | Database schema, API docs, codex integration, hybrid search |
| `inventory (Platform Inventory)/` | Software catalog, tools by category, ecosystem, cost tracking |
| `decisions (Decision Log)/` | Stack decisions, technical architecture, platform curation |

### scripts (Scripts & Automation)/
| Subfolder | What Lives Here |
|-----------|----------------|
| `nats (NATS Scripts)/config (NATS Config)/` | NATS server configs (hub, leaf, railway) |
| `nats (NATS Scripts)/scripts (NATS Scripts)/` | NATS Python scripts (bus, bridge, ingester, tools) |
| `deploy (Deploy Scripts)/` | Deployment scripts |
| `shell (Shell Scripts)/` | Shell scripts (health check, install, setup, firewall) |

### storage (Storage & Backup)/
| Subfolder | What Lives Here |
|-----------|----------------|
| `cloud (Cloud Storage)/backup (Backup)/` | Cloud-to-cloud backup comparisons, company backup strategy |
| `cloud (Cloud Storage)/sync (Sync Config)/` | GitHub storage, iCloud sync, self-hosted cloud architecture |
| `local (Local Storage)/` | NAS build guide, 6TB drive reference |
| `strategy (Storage Strategy)/` | 3-2-1 backup strategy, holistic storage, data ops, media access |

### Other Top-Level Folders
| Folder | What Lives Here |
|--------|----------------|
| `Tech-Devices (Personal Tech & Setup)/` | Personal tech inventory — Apps, Devices, Subscriptions |
| `vm-upgrades (VM Upgrade Plans)/` | VM upgrade planning and notes |
| `voice (Voice & Audio)/` | Voice system documentation and audio configs |

---

*Part of [SO — Sovereign Operations](../README.md). Route new captures to the correct subfolder; if unsure, use SO 24 Inbox.*

## ⚡ Active Tasks
<!-- Tasks sync to Todoist project: 6g6P9p83V2666MP3 -->

## 📦 Completed Archive
