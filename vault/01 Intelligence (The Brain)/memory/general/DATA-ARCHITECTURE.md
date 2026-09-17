---
id: root-00-20260113225954
title: "DATA ARCHITECTURE"
type: reference
domain: SO
area: "00 Framework"
agent: LEGACY
status: active
created: 2026-01-13
modified: 2026-04-30
owner: zeus
tags: [type/reference, status/active, domain/so, system/openclaw]
related:
  - OPENCLAW-FRAMEWORK-STRUCTURES.md
  - HIVE-CAPTURE-SYSTEM.md
todoist_project_id: "6g6P9mmQ3Pf36qjG"
ai_summary: "Master reference for accessing Zeus's data ecosystem — structured APIs + unstructured pipeline feeding OpenClaw agents."
---
# Zeus Data Architecture

**Created:** 2026-01-12
**Updated:** 2026-04-30
**Owner:** LEGACY + OpenClaw Agents
**Purpose:** Master reference for accessing Zeus's data ecosystem

---

## Architecture Overview

```
STRUCTURED DATA                    UNSTRUCTURED DATA
(Real-time API)                    (Vector Pipeline)
      │                                   │
      ▼                                   ▼
┌─────────────────┐              ┌─────────────────┐
│  Direct API     │              │  Dropbox MCP    │──► Quick Access
│  Calls          │              │  (Dash AI)      │
└───────┬─────────┘              └───────┬─────────┘
        │                                │
        │                                ▼
        │                        ┌─────────────────┐
        │                        │  rclone         │
        │                        │  (hourly sync)  │
        │                        └───────┬─────────┘
        │                                ▼
        │                        ┌─────────────────┐
        │                        │  Unstructured   │
        │                        │  (Docker)       │
        │                        └───────┬─────────┘
        │                                ▼
        │                        ┌─────────────────┐
        │                        │  Qdrant         │──► Deep Search
        │                        │  (Docker)       │
        │                        └───────┬─────────┘
        │                                │
        ▼                                ▼
┌──────────────────────────────────────────────────┐
│              OPENCLAW AGENTS                      │
│   (LEGACY + Spartan/Oracle/Coach/Eden/Empire)    │
└──────────────────────────────────────────────────┘
```

---

## PART 1: Data Sources & Connections

### 1.1 Unstructured Data (Documents)
**Primary Source:** Dropbox (500GB)

| Method | Purpose | Status |
|--------|---------|--------|
| Dropbox MCP | Quick semantic search via Dash AI | 📋 Planned |
| Qdrant Pipeline | Deep contextual search with filters | 📋 Planned |

**File Types:** PDF, DOCX, PPTX, TXT, MD, XLSX, CSV, Images, Code files

**Sync:** rclone hourly → `/data/dropbox` → Unstructured.io → Qdrant

### 1.2 Structured Data (SaaS)

| Platform | Tool | Status |
|----------|------|--------|
| Zoho CRM | `zoho_crm` | ✅ Built |
| Zoho Desk | `zoho_desk` | ✅ Built |
| Zoho Calendar | `zoho_calendar` | ✅ Built |
| Todoist | `todoist` | ✅ Built |
| GitHub | `github` | ✅ Built |
| Circle.so | `circle` | ✅ Built |
| Home Assistant | `home_assistant` | ✅ Built |

### 1.3 Future: Stripe, Wearables, Gmail, Slack

---

## PART 2: Tools Reference

### Active Tools (8)
- `zoho_crm` - Leads, Contacts, Deals, Accounts
- `zoho_desk` - Support tickets
- `zoho_calendar` - Events, availability
- `todoist` - Tasks, projects
- `github` - Vault files
- `circle` - Community members
- `home_assistant` - Smart home
- `web_search` - Internet search

### Planned Tools (2)
- `dropbox_mcp` - Quick file access
- `rag_search` - Deep semantic search

---

## PART 3: Agent Access & Real-World Examples

### 3.1 Agent Access Matrix

| Agent | Role | CRM | Desk | Calendar | Tasks | Vault | Docs |
|-------|------|-----|------|----------|-------|-------|------|
| **LEGACY** | Master Orchestrator | ✅ | ✅ | ✅ | ✅ | ✅ | 📋 |
| **Spartan / Coach** | HP — Personal Performance | ❌ | ❌ | ✅ | ✅ | ✅ | 📋 |
| **Empire / Oracle** | GE — Business/Revenue | ✅ | ✅ | ✅ | ✅ | ✅ | 📋 |
| **Eden / Haven** | LE — Family/Legacy | ❌ | ❌ | ✅ | ✅ | ✅ | 📋 |
| **Grid / Guardian** | SO — Infrastructure/Security | ❌ | ❌ | ❌ | ✅ | ✅ | 📋 |

**Legend:** ✅ Active | ❌ No Access | 📋 Planned

### 3.2 Agent Tool Distribution

**LEGACY (Master Orchestrator):**
- All 8 active tools
- Full data ecosystem access
- Routes queries to domain agents via NATS JetStream
- Runs on omos-core VPS (100.123.126.34)

**Spartan / Coach (HP — Personal):**
- `zoho_calendar` - Schedule management
- `todoist` - Task tracking
- `github` - Vault documentation
- Focus: Health, fitness, sleep, energy, habits

**Empire / Oracle (GE — Business):**
- `zoho_crm` - Sales pipeline
- `zoho_desk` - Customer support
- `zoho_calendar` - Business meetings
- `todoist` - Business tasks
- `github` - Business documentation
- Focus: Revenue, deals, customers, growth

**Eden / Haven (LE — Family/Legacy):**
- `zoho_calendar` - Family schedule
- `todoist` - Family/legacy tasks
- `github` - Knowledge documentation
- Focus: Family, homestead, relationships, heritage

**Grid / Guardian (SO — Infrastructure/Security):**
- `todoist` - Technical tasks
- `github` - System documentation
- Focus: Servers, NATS, OpenClaw platform, security

### 3.3 Real-World Query Examples

#### Example 1: Business Deal Research
```
User: "Find all contracts with ABC Corp mentioning payment terms"

Legacy AI Flow:
1. Receives query
2. Recognizes: needs document search + CRM check
3. Uses rag_search → Finds contracts (planned)
4. Uses zoho_crm → Gets deal history
5. Combines context
6. Returns: "Found 3 contracts, 2 open deals worth $125K"
```

#### Example 2: Daily Performance Check
```
User: "What's my performance status today?"

HP Conductor Flow:
1. Receives routed query from Legacy
2. Checks calendar → zoho_calendar (2 meetings)
3. Checks tasks → todoist (5 due, 3 overdue)
4. Checks vault → github (performance doc)
5. Returns: "2 meetings, 5 tasks due. Sleep logged: 7.2hrs"
```

#### Example 3: Revenue Review
```
User: "Show me this month's revenue and pipeline"

GE Conductor Flow:
1. Receives routed query from Legacy
2. Queries CRM → zoho_crm (closed deals this month)
3. Queries pipeline → zoho_crm (open opportunities)
4. Checks vault → github (revenue goals)
5. Returns: "Closed: $45K, Pipeline: $180K, Goal: $100K"
```

#### Example 4: Quick File Access
```
User: "Find the ABC contract I saved yesterday"

Legacy AI Flow:
1. Receives query
2. Recognizes: quick file lookup
3. Uses dropbox_mcp → Fast Dash AI search (planned)
4. Returns: "Found: /Contracts/ABC-Corp-2026.pdf"
```

### 3.4 Search Capability Matrix

| Search Type | Tool | Speed | Use Case |
|-------------|------|-------|----------|
| Quick file lookup | `dropbox_mcp` | <1s | Recent files, known names |
| Semantic search | `rag_search` | <100ms | "Find docs about X" |
| Keyword search | `rag_search` | <100ms | Exact term matching |
| Filtered search | `rag_search` | <100ms | "PDFs from 2024" |
| CRM query | `zoho_crm` | <500ms | Deals, contacts, pipeline |
| Task query | `todoist` | <200ms | Due dates, projects |
| Vault search | `github` | <1s | Documentation |
| Web search | `web_search` | 2-5s | Current internet info |

### 3.5 Implementation Timeline

| Phase | Target | Components | Status |
|-------|--------|------------|--------|
| **Phase 1** | Complete | 8 Direct API tools | ✅ Done |
| **Phase 2** | Week 1 | Dropbox MCP connection | 📋 Ready |
| **Phase 3** | Week 2 | Qdrant + Unstructured setup | 📋 Planned |
| **Phase 4** | Week 3 | rag_search tool creation | 📋 Planned |
| **Phase 5** | Week 4 | Agent tool distribution | 📋 Planned |

---

## PART 4: Operations, Security & Reference

### 4.1 Pipeline Components (Self-Hosted)

| Component | Purpose | Location | Docker |
|-----------|---------|----------|--------|
| **rclone** | Hourly Dropbox sync | VM cron | No |
| **Unstructured** | Document processing | Docker | Yes |
| **Qdrant** | Vector storage/search | Docker | Yes |
| **OpenClaw** | Agent platform + routing | omos-core VPS | Yes |
| **NATS JetStream** | Agent-to-agent messaging | omos-core VPS | Yes |

**Docker Compose:**
```yaml
services:
  qdrant:
    image: qdrant/qdrant:latest
    ports: ["6333:6333"]
    volumes: [qdrant_data:/qdrant/storage]
    
  unstructured:
    image: quay.io/unstructured-io/unstructured-api:latest
    ports: ["8000:8000"]
```

### 4.2 Security

**API Authentication:**
- Zoho: OAuth2 (auto-refresh)
- GitHub: Personal Access Token
- Dropbox: OAuth2
- Home Assistant: Long-lived token

**Access Control:**
- Agent permissions enforced by OpenClaw (L0-L3 autonomy levels)
- No API keys exposed to agents
- Sub-agents have restricted tool access via NATS subject namespacing

### 4.3 Monitoring

**Daily (Automated):**
- Dropbox sync status (rclone logs)
- Qdrant health + vector count
- API availability checks

**Weekly (Review):**
- Document processing lag
- Search latency (<100ms target)
- Failed job count

**Monthly (Audit):**
- Tool usage statistics
- Storage utilization
- Data completeness

### 4.4 Troubleshooting

| Problem | Check | Fix |
|---------|-------|-----|
| No search results | Qdrant running? | `docker ps`, restart if needed |
| CRM error | Token expired? | Refresh OAuth in Letta |
| Sync stale | Cron running? | Check `crontab -l` |
| Slow queries | Index health | Rebuild Qdrant index |

### 4.5 Cost Analysis

**Current (Working):**
| Service | Monthly |
|---------|---------|
| Zoho Bundle | ~$40 |
| Todoist | Free |
| GitHub | Free |
| Circle.so | ~$99 |
| Home Assistant | Free |
| **Subtotal** | **~$140** |

**Planned (Self-Hosted):**
| Service | Monthly |
|---------|---------|
| AWS EC2 (t3.large) | ~$70 (credits) |
| Unstructured | $0 (Docker) |
| Qdrant | $0 (Docker) |
| rclone | $0 (open source) |
| **Subtotal** | **~$70** |

**Total:** ~$210/mo (mostly covered by AWS credits)

### 4.6 Quick Command Reference

**For Legacy AI:**
```
"Find open deals" → zoho_crm (routes to Empire/Oracle via NATS)
"Open tickets?" → zoho_desk (routes to Empire via NATS)
"Today's schedule?" → zoho_calendar
"What's due?" → todoist
"Find vault file X" → github
"Community members?" → circle
"Office lights on" → home_assistant
"Search web for X" → web_search
"Find doc about X" → rag_search (planned)
"Find file X" → dropbox_mcp (planned)
```

### 4.7 Roadmap

**Q1 2026:**
- [x] Direct API tools (8)
- [ ] Dropbox MCP deployment
- [ ] Qdrant + Unstructured setup
- [ ] rag_search tool
- [ ] Agent tool distribution

**Q2 2026:**
- [ ] Stripe integration
- [ ] Wearable metrics
- [ ] Advanced RAG filtering

**Q3 2026:**
- [ ] Email integration
- [ ] Multi-agent collaboration
- [ ] Predictive analytics

---

## Document History

| Date | Change |
|------|--------|
| 2026-01-12 | Initial comprehensive document |

---

**Next Review:** 2026-02-12

*Last Updated: 2026-01-12 - Complete*


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
