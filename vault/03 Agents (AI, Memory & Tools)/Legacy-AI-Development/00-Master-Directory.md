---
id: ui-02-20260111053925
title: "Task Engine Master Directory"
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
ai_summary: "Task Engine Master Directory"
---

# Task Engine Master Directory

> **Todoist <-> Obsidian Bidirectional Sync**
> Last Updated: 2024-12-23

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    ONE MIND VAULT                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ 00-09 CMD   │ │ 10-19 HP    │ │ 20-29 LE    │           │
│  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘           │
│         │               │               │                   │
│  ┌──────┴──────┐ ┌──────┴──────┐ ┌──────┴──────┐           │
│  │ 30-39 GE   │ │ 40-49 IT    │ │ Task Views  │           │
│  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘           │
│         └───────────────┼───────────────┘                   │
│                         │                                   │
│                ┌────────▼────────┐                          │
│                │  Todoist Plugin │                          │
│                │  Query Blocks   │                          │
│                └────────┬────────┘                          │
└─────────────────────────┼───────────────────────────────────┘
                          │
                 ┌────────▼────────┐
                 │    TODOIST      │
                 │   (Task API)    │
                 └────────┬────────┘
                          │
         ┌────────────────┼────────────────┐
         ▼                ▼                ▼
    ┌─────────┐     ┌─────────┐     ┌─────────┐
    │ 00-09   │     │ 10-19   │     │ 20-29   │
    │ Command │     │ HP      │     │ LE      │
    └─────────┘     └─────────┘     └─────────┘
         │                │                │
    ┌────┴────┐     ┌────┴────┐     ┌────┴────┐
    │ 30-39   │     │ 40-49   │     │ Inbox   │
    │ GE      │     │ IT      │     │         │
    └─────────┘     └─────────┘     └─────────┘
```

---

## OMCS Project Structure

| Code | Pillar | Color | Projects |
|------|--------|-------|----------|
| 00-09 | Command | Purple | System, Framework, Zeus, Legacy-AI, Atlas, Protocols, Metrics, Resources, Environments, Archive |
| 10-19 | HP | Green | Base, Self, Lifestyle, Mind, Mastery, Performance, Finance, Sovereignty, Relationships |
| 20-29 | LE | Blue | Base, Sanctuary, Family, Legacy, Development, Wellness, Wealth, Legal, Community |
| 30-39 | GE | Orange | Base, Vision, Operations, Strategy, Creation, Analytics, Treasury, Compliance, Partnerships |
| 40-49 | IT | Red | Base, Stack, Workflows, Research, Development, Telemetry, Assets, Security, Integrations |

---

## Quick Access - Today's Tasks

```todoist
filter: "today | overdue"
name: "Today's Focus"
sorting:
  - priority
  - date
groupBy: project
show:
  - due
  - project
  - project
```

---

## All Pillars Dashboard

### 00-09 Command Tasks

```todoist
filter: "##00-09 Command"
name: "Command Layer"
groupBy: project
show:
  - due
  - project
```

### 10-19 HP Tasks

```todoist
filter: "##10-19 HP"
name: "Holistic Performance"
groupBy: project
show:
  - due
  - project
```

### 20-29 LE Tasks

```todoist
filter: "##20-29 LE"
name: "Legacy Evolution"
groupBy: project
show:
  - due
  - project
```

### 30-39 GE Tasks

```todoist
filter: "##30-39 GE"
name: "Generational Endeavor"
groupBy: project
show:
  - due
  - project
```

### 40-49 IT Tasks

```todoist
filter: "##40-49 IT"
name: "Innovative Technology"
groupBy: project
show:
  - due
  - project
```

---

## Priority Views

### P1 Critical

```todoist
filter: "p1"
name: "Critical Priority"
sorting:
  - date
groupBy: project
show:
  - due
  - project
```

### P2 High

```todoist
filter: "p2"
name: "High Priority"
sorting:
  - date
groupBy: project
show:
  - due
  - project
```

---

## Project Directory

### 00-09 Command (Meta)

| # | Project | Sections | Vault Path |
|---|---------|----------|------------|
| 00 | System | Api-integrations, Apps, Backups, Config, Docs, MCP-servers, Scripts, Workflows | [[00 System (Core Infrastructure)]] |
| 01 | Framework | Artificial Cognition, Cognitive Authority, Human Cognition, Legacy Codex, One Mind, Values | [[01 Framework (Unified Governance)]] |
| 02 | Zeus | Quick Capture, Dashboards, Templates | [[02 Zeus (Hub & Dashboard)]] |
| 03 | Legacy-AI | Architecture, Interfaces, Letta Agents, AI Communications, Observability, Physical-actions, Digital-actions, Control-APIs, Custom-skills | [[03 Legacy-AI (Co-Pilot)]] |
| 04 | Atlas | Codex Directory, Databases, Diagrams, Launch Pad, Locations, Maps, Operations, Systems | [[04 Atlas (Navigation)]] |
| 05 | Protocols | Decision Making, SOPs, Checklists | [[05 Protocols (Procedures)]] |
| 06 | Metrics | KPIs, Dashboards, Reports | [[06 Metrics (Health)]] |
| 07 | Resources | Templates, Assets, References | [[07 Resources (Assets)]] |
| 08 | Environments | Haven - USA, Remote, Travel | [[08 Environments (Locations)]] |
| 09 | Archive | Completed, Deprecated, Reference | [[09 Archive]] |

### 10-19 HP (Holistic Performance)

| # | Project | Sections | Vault Path |
|---|---------|----------|------------|
| 10 | Base | Inbox, Systems, Public Profile | [[10 Base (Personal Hub)]] |
| 11 | Self | Daily Notes, Journal, Vision Board, Readings, Personal Literature | [[11 Self (Identity & Purpose)]] |
| 12 | Lifestyle | Core Routines, Community, Relationships, Creative Expression, Entertainment, Mindset | [[12 Lifestyle (Daily Operations)]] |
| 13 | Mind | AI, Business, Finance, Crypto, Cybersecurity, Books & Research | [[13 Mind (Education & Learning)]] |
| 14 | Mastery | Certifications, Creative Skills, Language Learning, Martial Arts, Music, Physical Training | [[14 Mastery (Skills & Training)]] |
| 15 | Performance | Biometrics, Daily Routines, Fitness Tracking, Habit Tracking, Nutrition, Recovery | [[15 Performance (Body & Biology)]] |
| 16 | Finance | Banking, Budget, Credit, Investments, Taxes, Bills & Subscriptions | [[16 Finance (Money & Assets)]] |
| 17 | Sovereignty | Legal, Emergency Prep, Digital Security, Insurance, Identity, Asset Protection | [[17 Sovereignty (Privacy & Protection)]] |
| 18 | Relationships | Family, Friends, Professional, Community | [[18 Relationships (Social Network)]] |

### 20-29 LE (Legacy Evolution)

| # | Project | Sections | Vault Path |
|---|---------|----------|------------|
| 20 | Base | Inbox, Archive, Family Hub | [[20 Base (Family Hub)]] |
| 21 | Sanctuary | Home Maintenance, Gardening, Food & Recipes, Emergency Prep, Property, Vehicles | [[21 Sanctuary (Home & Property)]] |
| 22 | Family | Children, Calendar, Goals, Relationships, Travel | [[22 Family (Relationships & Roles)]] |
| 23 | Legacy | Heritage, History, Culture, Videos, Patents | [[23 Legacy (Heritage & History)]] |
| 24 | Development | Education, Skills, Library | [[24 Development (Education & Growth)]] |
| 25 | Wellness | Nutrition, Medical, Fitness | [[25 Wellness (Health & Medical)]] |
| 26 | Wealth | Family Office, Estate Planning, Investments, Asset Protection, Generational Wealth | [[26 Wealth (Family Office & Estate)]] |
| 27 | Legal | Trusts, Insurance, Child Support, Family Sovereignty | [[27 Legal (Protection & Documents)]] |
| 28 | Community | Extended Family, Network, Giving | [[28 Community (Extended Circle)]] |

### 30-39 GE (Generational Endeavor)

| # | Project | Sections | Vault Path |
|---|---------|----------|------------|
| 30 | Base | Inbox, Presentations, Quick Actions | [[30 Base (Business Hub)]] |
| 31 | Vision | Brand, Mission, Strategy, Scaling | [[31 Vision (Brand & Mission)]] |
| 32 | Operations | Entities, Infrastructure, Processes, Team | [[32 Operations (Execution & Entities)]] |
| 33 | Strategy | Planning, Intelligence, Roadmaps | [[33 Strategy (Planning & Intelligence)]] |
| 34 | Creation | Products, Content, R&D, Patents | [[34 Creation (Products & Content)]] |
| 35 | Analytics | Dashboards, Reports, KPIs | [[35 Analytics (Business Intelligence)]] |
| 36 | Treasury | Revenue, Expenses, Real Estate, Investments | [[36 Treasury (Finance & Investments)]] |
| 37 | Compliance | Legal, Contracts, Licensing, Tax | [[37 Compliance (Legal & Contracts)]] |
| 38 | Partnerships | Clients, Vendors, Network, Affiliates | [[38 Partnerships (Clients & Network)]] |

### 40-49 IT (Innovative Technology)

| # | Project | Sections | Vault Path |
|---|---------|----------|------------|
| 40 | Base | Inbox, Telemetry, Quick Actions | [[40 Base (Technology Hub)]] |
| 41 | Stack | Cloud, Local, DevOps, Software | [[41 Stack (Infrastructure & Tools)]] |
| 42 | Workflows | N8N, Automations, Communications | [[42 Workflows (Automation & Processes)]] |
| 43 | Research | Emerging Tech, Robotics, AI | [[43 Research (R&D & Emerging Tech)]] |
| 44 | Development | Projects, Code, Models, Platforms | [[44 Development (Code & AI)]] |
| 45 | Telemetry | Monitoring, Logs, Alerts | [[45 Telemetry (Monitoring & Logs)]] |
| 46 | Assets | Accounts, Domains, Credentials | [[46 Assets (Digital Resources)]] |
| 47 | Security | 2FA, Recovery, Protocols, Audits | [[47 Security (Cyber Protection)]] |
| 48 | Integrations | APIs, MCP, Webhooks | [[48 Integrations (APIs & Connections)]] |

---

## Agent Labels

| Label | Purpose | Color |
|-------|---------|-------|
| @mother-ai | Mother AI orchestration | Purple |
| @legacy-ai | Legacy AI co-pilot | Purple |
| @hp-conductor | HP pillar conductor | Green |
| @le-conductor | LE pillar conductor | Blue |
| @ge-conductor | GE pillar conductor | Orange |
| @it-conductor | IT pillar conductor | Red |
| @approval-needed | Requires human review | Yellow |
| @blocked | Waiting on dependency | Grey |
| @urgent | Time-sensitive | Red |

---

## Sync Configuration

### Plugin Settings

```yaml
# Obsidian Todoist Plugin
API Token: (stored in .obsidian/todoist-token)
Auto-refresh: 300 seconds
Default filter: "today | overdue"
```

### Query Block Reference

```todoist
filter: "<todoist filter>"
name: "<display name>"
sorting:
  - date | priority | dateAdded
  - ascending | descending
groupBy: project | section | due | labels | priority
show:
  - due | date
  - deadline
  - description
  - project
  - labels
autorefresh: <seconds>
```

---

## Related Files

- [[todoist-config]] - API configuration
- [[todoist_setup.py]] - Structure setup script
- [[Todoist Integration Architecture]] - System design
