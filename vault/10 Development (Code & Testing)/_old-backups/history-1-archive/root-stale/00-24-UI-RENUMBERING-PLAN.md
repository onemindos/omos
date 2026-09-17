---
id: root-00-20260113225954
title: "00 24 UI RENUMBERING PLAN"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "00 24 UI RENUMBERING PLAN"
---
# 00-24 SO Renumbering Plan

## Problem

The current folder structure in `00-24 SO/` does NOT match the canonical structure defined in ONEMIND-CODEX.md lines 609-630.

---

## Current vs Should Be

| Current | Should Be (per CODEX) | Action |
|---------|----------------------|--------|
| **00 Framework** | 00 Framework | ✅ CORRECT |
| **01 Command** | 01 Command | ✅ CORRECT |
| **02 Agents** | 02 Agents | ✅ CORRECT |
| ❌ MISSING | **03 Atlas** (Navigation & Maps) | 🆕 CREATE |
| ❌ MISSING | **04 Protocols** (SOPs) | 🔄 RENUMBER from 13 |
| **05 Resources** | **05 Metrics** (Scorecards & Tracking) | 🔄 RENAME/REPURPOSE |
| ❌ MISSING | **06 Resources** (Assets & References) | 🔄 RENUMBER from 05 |
| **07 Infrastructure** | **07 Environments** (Contexts & Locations) | 🔄 REPURPOSE or MOVE |
| **08 Automation** | **08 Inbox** (Unified Capture) | 🔄 RENUMBER from 14 |
| ❌ MISSING | **09 Archive** (Historical) | 🆕 CREATE |
| **10 Security** | **10 Infrastructure** (Servers, networks) | 🔄 RENUMBER from 07 |
| **11 Research** | **11 Automation** (Workflows, n8n) | 🔄 RENUMBER from 08 |
| ❌ MISSING | **12 Development** (Code, repos) | 🔄 RENUMBER from 09 |
| **13 Protocols** | **13 Security** (Cyber protection) | 🔄 RENUMBER from 10 |
| **14 Inbox** | **14 Integrations** (APIs, connections) | 🆕 CREATE or REPURPOSE |
| **15 Future Roadmap** | **15 Research** (R&D, emerging tech) | 🔄 RENUMBER from 11 |
| ❌ MISSING | **16 Assets** (Digital resources) | 🆕 CREATE |
| ❌ MISSING | **17 Telemetry** (Monitoring, logs) | 🆕 CREATE |
| ❌ MISSING | **18-24 Reserved** | Future expansion |

---

## Renumbering Strategy

### Option A: Full Canonical Alignment (Recommended)
Renumber everything to match ONEMIND-CODEX.md exactly.

**Pros:**
- Perfect alignment with canonical source
- Future-proof
- AI agents can rely on CODEX

**Cons:**
- Large renaming operation
- Need to update all references

---

### Option B: Keep Current, Update CODEX
Keep current structure, update CODEX to match reality.

**Pros:**
- No file moves needed
- Less disruption

**Cons:**
- CODEX is the source of truth, shouldn't change it
- Defeats purpose of having canonical structure

---

## Recommendation: Option A (Full Alignment)

**Reason:** ONEMIND-CODEX.md is THE source of truth. The folders should match it, not the other way around.

---

## Execution Plan (Option A)

### Step 1: Create Missing Folders
```bash
cd "00-24 SO (Sovereign Operations)"

# Brain Functions
mkdir -p "03 Atlas (Navigation & Maps)"
mkdir -p "09 Archive (Historical)"

# Tech Functions
mkdir -p "16 Assets (Digital Resources)"
mkdir -p "17 Telemetry (Monitoring & Logs)"
```

### Step 2: Rename in Reverse Order (Avoid Conflicts)
```bash
# Work backwards to avoid numbering conflicts

# Tech Functions (15 down to 10)
mv "11 Research (R&D)" "15 Research (R&D)"
mv "10 Security (Cyber)" "13 Security (Cyber)"
mv "09 Development (Code & Testing)" "12 Development (Code & Testing)"
mv "08 Automation (Workflows)" "11 Automation (Workflows)"
mv "07 Infrastructure (Grid & APIs)" "10 Infrastructure (Grid & APIs)"

# Brain Functions (08 down to 04)
mv "14 Inbox (Queue)" "08 Inbox (Queue)"
mv "05 Resources (Assets & Templates)" "06 Resources (Assets & References)"
mv "13 Protocols (SOPs)" "04 Protocols (SOPs)"
```

### Step 3: Create/Reorganize Special Cases

**05 Metrics (NEW):**
- Create folder
- Move scorecard/tracking content here
- This is for PI/LI/WI scores, KPIs

**07 Environments (REPURPOSE):**
- Create folder for contexts/locations
- Move relevant content from old Infrastructure

**14 Integrations (CREATE):**
- Create folder for API connections
- Move integration docs from Infrastructure

---

## What Goes Where (After Renumbering)

### Brain Functions (00-09)

**00 Framework** - ✅ Already correct
- Core doctrine, principles, architecture docs

**01 Command** - ✅ Already correct
- Human interface, dashboard, command center

**02 Agents** - ✅ Already correct
- AI agents, memory blocks, skills, tools

**03 Atlas** - 🆕 CREATE
- Navigation maps, vault structure guides
- "How to find anything" reference
- Move: Navigation docs from Command

**04 Protocols** - 🔄 RENAME from 13
- Standard Operating Procedures
- Runbooks, checklists, how-tos

**05 Metrics** - 🆕 CREATE
- Scorecards (PI, LI, WI, Unified)
- KPI tracking
- Dashboard data sources

**06 Resources** - 🔄 RENUMBER from 05
- Assets, templates, references
- Keep current content

**07 Environments** - 🆕 CREATE/REPURPOSE
- Contexts (home, office, mobile)
- Locations, devices
- Environment-specific configs

**08 Inbox** - 🔄 RENUMBER from 14
- Unified capture point
- Unprocessed items
- HIVE protocol entry point

**09 Archive** - 🆕 CREATE
- Historical records
- Completed projects
- Deprecated content

---

### Tech Functions (10-24)

**10 Infrastructure** - 🔄 RENUMBER from 07
- Servers, networks, homelab
- Docker, VMs, cloud resources

**11 Automation** - 🔄 RENUMBER from 08
- N8N workflows
- Scripts, cron jobs
- Workflow automation

**12 Development** - 🔄 RENUMBER from 09
- Code repositories
- Development projects
- Technical creation

**13 Security** - 🔄 RENUMBER from 10
- Cybersecurity
- Credentials, secrets
- Backup strategies

**14 Integrations** - 🆕 CREATE
- API connections
- Webhooks
- Service integrations (Todoist, GitHub, etc.)

**15 Research** - 🔄 RENUMBER from 11
- R&D, emerging tech
- Experimental projects
- Research notes

**16 Assets** - 🆕 CREATE
- Digital resources
- Media files
- Reusable components

**17 Telemetry** - 🆕 CREATE
- System monitoring
- Logs, metrics
- Observability (Grafana, Prometheus)

**18-24 Reserved** - Future expansion

---

## Migration Checklist

### Pre-Migration
- [ ] Backup entire vault (git commit)
- [ ] Document current state
- [ ] Verify no active work in folders being moved

### Execution
- [ ] Create missing folders (03, 05, 07, 09, 14, 16, 17)
- [ ] Rename folders in reverse order (15→11, 13→10, etc.)
- [ ] Move content to repurposed folders
- [ ] Remove empty folders
- [ ] Update INDEX.md references
- [ ] Update AI agent navigation protocols

### Post-Migration
- [ ] Verify all folders match CODEX
- [ ] Test AI agent navigation
- [ ] Update any hardcoded paths in scripts
- [ ] Commit changes to git

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking links | Medium | Use Obsidian's link update feature |
| Path conflicts | High | Rename in reverse order |
| Lost content | Critical | Git commit before starting |
| Agent confusion | Medium | Update agent memory after migration |

---

## Time Estimate

- **Planning & Backup:** 15 minutes
- **Folder operations:** 30 minutes
- **Content reorganization:** 45 minutes
- **Verification:** 30 minutes
- **Total:** ~2 hours

---

## Decision Required

**Should we:**
1. ✅ **Execute full renumbering** (aligns with CODEX)
2. ❌ **Keep current structure** (updates CODEX instead)
3. ⏸️ **Defer renumbering** (address later)

**Recommendation: Option 1** - Execute full renumbering now while momentum is high and we're already reorganizing.

---

*This plan ensures 00-24 SO matches the canonical structure in ONEMIND-CODEX.md (lines 609-630).*
