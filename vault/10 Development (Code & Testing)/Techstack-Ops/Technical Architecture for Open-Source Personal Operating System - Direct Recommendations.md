---
id: ui-12-20260111053926
title: "Technical Architecture for Open-Source Personal Operating System: Direct Recommendations"
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
ai_summary: "Technical Architecture for Open-Source Personal Operating System: Direct Recomme"
---

# Technical Architecture for Open-Source Personal Operating System: Direct Recommendations

## Critical Reality Check First

**You have a major blocker**: The Letta Obsidian plugin is **read-only**. Letta cannot mark tasks complete, create tasks, or modify files. It can only read vault content through its Filesystem API. This fundamentally changes your architecture options.

The Letta plugin is in alpha testing with a warning to "DO NOT USE ON NOTES YOU CANNOT AFFORD TO LOSE." It's designed for fewer than 100 documents, not production task management at scale.

## 1. Hybrid vs All-in-Obsidian: **Use Hybrid Architecture**

**Recommendation: Hybrid approach with Obsidian as central knowledge base**

Keep Obsidian for what it does best—knowledge management, notes, and processed information. Use specialized external tools for features requiring rich UIs, multi-user collaboration, or advanced automation, then sync curated content back to Obsidian.

**Why not all-in-Obsidian**:
- No native AI agent write capabilities (Letta is read-only)
- Performance degrades significantly beyond 10,000 notes
- Limited collaboration features
- No centralized bookmark/article management UI
- Manual archiving and link monitoring
- Task management plugins lack programmatic AI integration

**Hybrid advantages**:
- Specialized tools with superior feature sets
- Better mobile experiences
- Real-time collaboration where needed
- Automatic archiving and preservation
- Obsidian remains single source of truth for processed knowledge
- Data sovereignty maintained through self-hosting

## 2. Recommended Tech Stack

### Core Infrastructure

**Central Knowledge Base**: 
- **Obsidian** - Personal notes, processed information, final outputs
- Vault stored as plain markdown files

**Sync Solution**:
- **Git-based sync** via Obsidian Git plugin (Primary recommendation)
  - Best for AI agent workflows (programmatic control)
  - Complete version history and audit trail
  - Agent can pull → modify → commit → push
  - Works perfectly offline
  - Free (GitHub/GitLab)
- **Alternative**: Obsidian Sync ($4-8/month) if simplicity matters more than AI control

**Task Management** (Working around Letta limitations):
- **Obsidian Tasks plugin** for task storage and human interface
- **Custom automation layer** (not Letta) for AI task operations:
  - Option A: Python script using Obsidian Local REST API plugin
  - Option B: n8n workflows with file system access
  - Option C: Custom Obsidian plugin development
- **Reality**: Native Letta cannot fulfill your task management requirements. You need a different AI integration approach.

**File Storage**:
- **Nextcloud** (self-hosted) replacing Dropbox
  - Full file access for attachments
  - WebDAV integration
  - 2TB+ storage capacity
  - Mobile apps
  - Sync task attachments to Obsidian vault subdirectory

**Automation Platform**:
- **n8n** (self-hosted) for workflow automation
  - Webhook integration via Obsidian Post Webhook plugin
  - Capture workflows from all sources
  - AI processing integration
  - Alternative: Make.com if self-hosting is too complex

### Specialized External Tools

**Bookmark Management**:
- **LinkWarden** (self-hosted)
  - Automatic screenshots, PDF, HTML archiving
  - Optional AI tagging via local Ollama
  - Full REST API for Obsidian sync
  - Wayback Machine integration
  - Replaces: Raindrop, Brave bookmarks

**Read-It-Later**:
- **Wallabag** (self-hosted) 
  - Reader view with annotations
  - Dedicated Obsidian plugin (obsidian-wallabag)
  - Mobile apps for reading
  - Article archiving
  - Replaces: Reader

**Whiteboarding**:
- **Excalidraw plugin** (in Obsidian) for solo work
  - Native integration, files stored in vault
  - Hand-drawn aesthetic
  - Perfect for personal diagrams
- **tldraw** (self-hosted) for team collaboration
  - Real-time multiplayer
  - Export final versions to Obsidian
  - Replaces: Miro for collaborative sessions

**Chat with Data** (Replacing ClickUp Brain):
- **Obsidian + Local LLM Setup**:
  - Smart Connections plugin with local embeddings
  - Copilot plugin for AI assistance
  - Or run Obsidian MCP Server with Claude Desktop
- **Alternative**: Custom RAG pipeline with vector database if vault exceeds 10,000 notes

### AI Agent Integration (Addressing Letta Limitations)

**Since Letta cannot write files, implement alternative AI architecture**:

**Option A: Python Automation Layer** (Recommended)
```
Python Script/Service
  ↓
Obsidian Local REST API Plugin
  ↓
Full CRUD on vault (read/write/update tasks)
  ↓
Git commits for audit trail
```

**Capabilities**:
- Read task files, parse markdown task syntax
- Mark tasks complete by updating `- [ ]` to `- [x]`
- Create new tasks in designated files
- Access attachments via filesystem
- Full version control via Git

**Option B: n8n Workflow Automation**
```
Scheduled Trigger / Webhook
  ↓
AI Processing (OpenAI/Anthropic/Local LLM)
  ↓
Task Operations via filesystem or REST API
  ↓
Response back to Obsidian
```

**Option C: Custom Obsidian Plugin**
- Build stateful AI agent directly in Obsidian
- Full API access for task manipulation
- Can integrate any AI model
- Most flexible but requires TypeScript development

**Letta Alternative**: Consider migrating to agents built with LangChain, CrewAI, or AutoGen that have better file system integration and can actually modify markdown files.

## 3. Plugin Stack vs External Tools

### Essential Obsidian Plugins

**Core Functionality**:
1. **Obsidian Tasks** - Task management syntax and queries
2. **Obsidian Git** - Automated Git sync (commit every 10-30 min)
3. **Local REST API** - Enable programmatic vault access
4. **Templater** - Template automation
5. **Dataview** - Task queries and data aggregation

**Capture \u0026 Integration**:
6. **Post Webhook** - Bidirectional n8n integration
7. **Obsidian Web Clipper** (browser extension) - Direct web capture
8. **Wallabag Plugin** - Article sync from Wallabag

**Knowledge Management**:
9. **Excalidraw** - Whiteboarding
10. **Smart Connections** or **Copilot** - Local AI chat
11. **Advanced URI** - Deep linking for automation

**Optional**:
12. **Auto Note Mover** - Automatic organization
13. **RSS Reader** - Feed management (or use external tool)

### External Tools (Use Instead of Plugins)

**Use external tools when**:
- Feature requires rich UI (bookmark browsing)
- Collaboration needed (whiteboarding)
- Mobile-first experience matters (article reading)
- Automatic archival critical (link preservation)
- Multi-user management required

**Keep in Obsidian when**:
- Personal notes and writing
- Information already in markdown
- Integration with other notes important
- Offline access critical
- Simplicity matters more than features

## 4. Sync Solution: Git-Based (with Obsidian Sync as Alternative)

### Primary Recommendation: Git-Based Sync

**Implementation**:
1. Initialize Git repository in Obsidian vault
2. Create GitHub private repository (free for individuals)
3. Install Obsidian Git plugin
4. Configure auto-commit every 10 minutes
5. Add `.gitignore` for `.obsidian/workspace.json` (reduces conflicts)
6. Use Git LFS for large binary attachments

**Why Git wins for AI agents**:
- **Programmatic control**: AI agent runs `git pull` before operations, commits after
- **Audit trail**: Every change tracked with commit messages
- **Conflict resolution**: Standard Git merge tools
- **Rollback safety**: Easy to revert AI mistakes
- **Zero file locking**: Optimistic concurrency
- **Free**: GitHub/GitLab private repos
- **Best offline**: Full history available locally

**AI Agent Workflow with Git**:
```python
# Agent task modification workflow
os.chdir(vault_path)
subprocess.run(['git', 'pull'])  # Get latest
modify_task_file()  # Agent makes changes
subprocess.run(['git', 'add', '.'])
subprocess.run(['git', 'commit', '-m', 'AI: Completed task X'])
subprocess.run(['git', 'push'])
```

**Mobile Consideration**: Use Working Copy app (iOS) or MGit (Android). Alternatively, keep one primary device for AI operations and other devices read-only.

### Alternative: Obsidian Sync

**Use if**:
- Simplicity more important than AI control
- Don't want to manage Git
- Need seamless mobile sync
- Willing to pay $4-8/month

**Implementation for AI**:
- Agent still accesses local files directly
- Manually pause sync during batch AI operations
- Monitor sync status via plugin API
- Let Obsidian Sync handle conflict resolution

**Not Recommended**:
- **Nextcloud/Seafile via plugins**: More complex, file locking issues with AI
- **Syncthing**: Excessive conflict files with Obsidian, terrible for AI agents

## 5. Task Management: Native Obsidian Tasks (Not Todoist)

**Recommendation: Obsidian Tasks plugin + Custom AI Integration**

**Why not Todoist integration**:
- Adds external dependency
- Data not in markdown (portability loss)
- Still no native AI agent integration
- Todoist API rate limits
- Subscription cost

**Why Obsidian Tasks**:
- Native markdown syntax: `- [ ] Task description`
- Data stays in vault (data ownership)
- Powerful query language via Dataview
- Recurring tasks via recurrence rules
- Works offline
- Free and open-source

**Task Format**:
```markdown
- [ ] Task description 📅 2025-10-30 ⏫ #homesteading
- [x] Completed task ✅ 2025-10-27
- [ ] Task with recurrence 🔁 every week
```

**AI Integration Implementation**:

Since Letta cannot write files, use **Obsidian Local REST API** with custom script:

```python
import requests

OBSIDIAN_API = "https://127.0.0.1:27123"
API_KEY = "your_api_key"
headers = {"Authorization": f"Bearer {API_KEY}"}

# Read task file
response = requests.get(
    f"{OBSIDIAN_API}/vault/Tasks/daily.md",
    headers=headers,
    verify=False  # Self-signed cert
)
content = response.json()["content"]

# Modify tasks (mark complete)
updated_content = content.replace(
    "- [ ] Review homesteading video",
    "- [x] Review homesteading video ✅ 2025-10-27"
)

# Write back
requests.put(
    f"{OBSIDIAN_API}/vault/Tasks/daily.md",
    headers=headers,
    json={"content": updated_content},
    verify=False
)
```

**Scalability**: For multiple AI agent instances, implement distributed locking via Redis to prevent concurrent modifications.

## 6. Setting Up AI Agent Access Properly

### Letta Current Reality (If You Insist on Using It)

**What Letta CAN Do**:
- Read vault files through Letta Filesystem API
- Semantic search across notes
- Query information
- Provide context-aware responses

**What Letta CANNOT Do**:
- Mark tasks complete
- Create new tasks
- Modify files
- Write to vault

**Setup Process**:
1. Install Letta (self-hosted): `docker run -p 8283:8283 letta/letta`
2. Install `letta-obsidian` plugin from GitHub (manual installation)
3. Configure Letta API endpoint in plugin settings
4. Create Letta agent and attach folder
5. Enable auto-sync (but limit to \u003c100 important documents)

**Use Case**: Treat Letta as **read-only research assistant**, not task manager.

### Recommended AI Architecture (Actually Functional)

**Approach 1: Local REST API + Custom Agent** (Most Flexible)

**Stack**:
- Obsidian Local REST API plugin
- Python service with OpenAI/Anthropic/Local LLM
- Task parsing and modification logic
- Git integration for version control

**Architecture**:
```
Trigger (schedule/webhook)
  ↓
Python AI Agent Service
  ├→ Read tasks via REST API
  ├→ LLM decision making
  ├→ Modify tasks via REST API
  ├→ Access attachments via filesystem
  └→ Git commit changes
```

**Capabilities**:
- Full task CRUD operations
- Attachment access (read task instructions + files)
- Stateful through database or context storage
- Multiple instances via distributed locking
- Scalable to hundreds of agents

**Approach 2: n8n Workflow Automation** (Easier Setup)

**Stack**:
- n8n (self-hosted)
- Post Webhook plugin in Obsidian
- Google Drive/Dropbox for attachments
- OpenAI/Anthropic nodes for AI

**Example Workflow**:
```
Schedule Trigger (daily 8am)
  ↓
Read task file from Google Drive
  ↓
Parse incomplete tasks
  ↓
For each task:
  ├→ Read attachments if present
  ├→ AI analysis (should this be done?)
  ├→ Execute or mark complete
  └→ Update task file
  ↓
Write back to Google Drive
  ↓
Obsidian syncs changes
```

**File Access for Attachments**:
- Store attachments in Nextcloud
- Create symlink to Obsidian vault: `vault/attachments → nextcloud/obsidian-files`
- Agent accesses via filesystem path
- Nextcloud handles syncing to other devices

## 7. Best Capture Workflow Architecture

### Multi-Source Capture System

**Architecture Pattern**:
```
Source → Browser/App → n8n Webhook → Processing → Obsidian
```

**Core Tool**: n8n with Obsidian Post Webhook plugin

### Source-Specific Workflows

**Email Capture**:
```
Gmail Trigger (new email with label "capture")
  ↓
Extract: sender, subject, body, attachments
  ↓
Optional: AI summarization
  ↓
Format as markdown with frontmatter
  ↓
Google Drive write (synced folder)
  ↓
Appears in Obsidian/Inbox/
```

**YouTube Videos** (Your Homesteading Example):
```
YouTube URL → Webhook/Airtable
  ↓
YouTube API (get video details)
  ↓
RapidAPI (get transcript)
  ↓
OpenAI (generate summary + action items)
  ↓
Format markdown:
  ---
  title: "Homesteading: Build Chicken Coop"
  url: https://youtube.com/watch?v=...
  tags: [homesteading, todo, video]
  priority: medium
  ---
  
  ## Summary
  [AI summary]
  
  ## Action Items
  - [ ] Buy materials (see timestamp 3:45)
  - [ ] Review coop design (timestamp 8:20)
  
  ## Transcript
  [full transcript]
  ↓
Save to vault via Google Drive or Post Webhook
  ↓
Shows up in Obsidian for review and execution
```

**Social Media**:
```
Twitter/X: Use Tweet to Markdown plugin or Tressel service
Reddit: Manual save or n8n Reddit API workflow
Instagram: Data export import
```

**Articles/Bookmarks**:
```
Option A: Obsidian Web Clipper (browser extension)
  ↓ Direct to vault

Option B: LinkWarden (comprehensive solution)
  ↓ Save to LinkWarden via extension
  ↓ Weekly sync to Obsidian via API
  
Option C: Wallabag (read-later)
  ↓ Save to read
  ↓ After reading/annotating, sync to Obsidian
```

**RSS Feeds**:
```
Option A: RSS Reader plugin in Obsidian
  - Simple, built-in
  - Create notes from interesting items
  
Option B: n8n RSS Trigger
  ↓ New feed items
  ↓ Filter by keywords
  ↓ AI relevance scoring
  ↓ High-scoring items → Obsidian
```

**Newsletters**:
```
Dedicated Gmail label "newsletters"
  ↓
n8n Gmail trigger (new emails)
  ↓
HTML → Markdown conversion
  ↓
AI summary + key points extraction
  ↓
Save to Obsidian/Newsletters/YYYY-MM-DD-{title}.md
```

### Integration Points

**n8n Setup**:
1. Self-host n8n via Docker or use n8n Cloud
2. Create webhook endpoints for each source type
3. Install Post Webhook plugin in Obsidian
4. Configure bidirectional flows

**Post Webhook Configuration**:
```yaml
# In Obsidian Post Webhook settings
Webhook Name: Process with n8n
URL: https://n8n.yourdomain.com/webhook/obsidian-process
Method: POST
Response Handling: Append
```

**Authentication**:
- Use Bearer tokens for n8n webhooks
- HMAC signature verification for security
- HTTPS only (use Caddy or Nginx reverse proxy)

## 8. Potential Pitfalls and Technical Limitations

### Critical Issues to Avoid

**Letta Integration**:
- ⚠️ **Read-only limitation**: Cannot accomplish your task management requirements
- ⚠️ **Alpha status**: Risk of data loss, plugin warns against production use
- ⚠️ **File count limit**: Designed for \u003c100 documents, not full vaults
- ⚠️ **No bidirectional sync**: Changes in Letta don't flow back to Obsidian
- ⚠️ **Server dependency**: Requires always-on Letta server

**Obsidian Performance**:
- Performance degrades significantly beyond 10,000 notes
- Link suggestions take 4 seconds per character at scale
- Initial indexing 30-40 minutes for large vaults
- Heavy plugins (Dataview) compound slowdown
- **Mitigation**: Split into multiple vaults or implement hybrid architecture with metadata database

**Sync Conflicts**:
- Multiple devices editing simultaneously → conflicts
- Git conflicts require manual resolution knowledge
- Obsidian Sync occasionally overwrites newer data
- **Mitigation**: Designate primary device for AI operations, implement file locking

**Concurrent Access**:
- Multiple AI agents writing simultaneously → corruption
- No built-in file locking in markdown
- Race conditions with automation
- **Mitigation**: Distributed locking via Redis, single-writer pattern, or optimistic locking with version checks

**Webhook Reliability**:
- Network failures cause missed captures
- n8n workflow failures lose data
- No automatic retry for failed webhooks
- **Mitigation**: Implement idempotency keys, dead letter queues, exponential backoff retry

**Mobile Limitations**:
- Git on mobile requires external apps (Working Copy, MGit)
- Heavy vaults drain battery
- Sync delays cause confusion
- Limited plugin support on mobile
- **Mitigation**: Keep mobile mostly read-only, primary work on desktop

**Plugin Compatibility**:
- Updates break workflows without warning
- Plugin incompatibilities cause crashes
- Community plugins not officially supported
- **Mitigation**: Pin plugin versions in production, test updates in separate vault

**Backup Failures**:
- Sync is NOT backup
- Obsidian has no native backup system
- Corrupted Git repos lose history
- **Mitigation**: Automated encrypted backups (Restic/Duplicati) following 3-2-1 rule, regular restore testing

**Automation Complexity**:
- Over-engineering leads to brittle systems
- Too many plugins create conflicts
- Debugging automation workflows is hard
- **Mitigation**: Start simple, document everything, version control all configs

**Data Lock-In Risks**:
- Obsidian-specific syntax reduces portability
- Plugin-dependent workflows break if plugin abandoned
- Binary attachments not markdown-portable
- **Mitigation**: Use standard markdown syntax, minimize proprietary features, regular exports to HTML/PDF

## 9. Scalability and Future-Proofing

### Scalability Architecture

**Current Scale (Launch)**:
- File-based Obsidian vault (\u003c 5,000 notes)
- Git sync for version control
- Single AI agent instance on primary device
- n8n automation on single server

**Medium Scale (1-2 years)**:
- Obsidian vault reaches 5,000-10,000 notes
- **Add**: PostgreSQL metadata database
  - Full-text search index
  - Task metadata
  - Cached content for fast queries
- Vault remains source of truth, database is cache
- 2-3 AI agent instances with distributed locking

**Architecture Evolution**:
```
Current (MVP):
Obsidian Vault (files) ← AI Agent ← n8n

Medium (Scale):
                    ┌→ PostgreSQL (metadata/search)
                    │
Obsidian Vault ←───┼─ AI Agents (3x) ←─ n8n
                    │
                    └→ Redis (locks/cache)

Large (Enterprise):
                              ┌→ PostgreSQL + pgvector
                              │
Obsidian Vault (files) ←─────┼─ AI Orchestrator
                              │   ↓
                              └─ Agent Pool (10+) + Redis
                                  ↓
                              Event Queue (RabbitMQ)
```

**Large Scale (Theoretical - 50K+ notes)**:
- Split vault into domain-specific vaults
- Vector database (Pinecone/Weaviate) for semantic search
- Agent orchestration layer (supervisor-worker pattern)
- Event-driven architecture (RabbitMQ/Kafka)
- Full observability stack (Prometheus/Grafana)

### Future-Proofing Strategies

**Data Portability**:
- Plain text markdown (universally readable)
- Minimize Obsidian-specific syntax
- Regular exports to standard formats (HTML, PDF)
- Keep raw data in Nextcloud (independent of Obsidian)

**Technology Independence**:
- Use open standards (REST APIs, OAuth, Webhooks)
- Self-host everything (no SaaS dependencies)
- Containerize all services (Docker)
- Document all integration points

**Modularity**:
- Each component replaceable independently
- Use hexagonal architecture (ports and adapters)
- Decouple file storage from business logic
- Standard interfaces between components

**Version Control Everything**:
- Vault via Git
- Automation configs in Git
- Infrastructure as Code (Terraform/Docker Compose)
- Regular config backups

**Monitoring and Observability**:
- Structured logging (JSON format)
- Metrics collection (Prometheus)
- Error alerting (Grafana/PagerDuty)
- Trace AI agent actions for debugging

**Long-Term Data Strategy**:
- Encrypted backups with 3-2-1 rule (3 copies, 2 media, 1 offsite)
- Test restore procedures quarterly
- Plan for Obsidian obsolescence (export paths)
- Keep human-readable formats

**Community Participation**:
- Contribute to open-source projects you depend on
- Fork critical plugins for insurance
- Document your architecture publicly
- Build community around your stack

### Technology Refresh Plan

**Quarterly Reviews**:
- Check for security updates
- Review plugin health (last commit date, open issues)
- Performance benchmarking
- Backup restore testing

**Annual Major Reviews**:
- Evaluate new tools (better alternatives?)
- Architecture fitness assessment
- Scale testing
- Disaster recovery drill

**Abandonment Signals to Watch**:
- Plugin not updated in 12+ months
- Developer announces deprecation
- Major security issues unfixed
- Community migration away from tool

**Exit Strategies**:
- Obsidian → Alternative: Logseq, Foam, or plain text editors
- LinkWarden → Alternative: LinkAce, linkding
- Wallabag → Alternative: Omnivore fork, self-hosted reader
- n8n → Alternative: Zapier, Make, Pipedream (with lock-in cost)
- Git → Alternative: Always portable (just files)

## Final Recommendations Summary

### ✅ DO THIS

1. **Use hybrid architecture**: Obsidian (central KB) + specialized tools (LinkWarden, Wallabag, Nextcloud)
2. **Sync with Git**: Best for AI agents, version control, and offline work
3. **Replace Letta for task management**: Use Obsidian Local REST API + custom Python/n8n agent
4. **Obsidian Tasks plugin**: Not Todoist (data ownership, no external dependency)
5. **n8n for capture workflows**: Post Webhook plugin for integration
6. **Self-host everything**: LinkWarden, Wallabag, Nextcloud, n8n, Git (Gitea/GitLab)
7. **Start simple**: File-based system, add database only when performance degrades
8. **Automated encrypted backups**: Restic/Duplicati with 3-2-1 rule
9. **Version control**: Git for vault, configs, and infrastructure

### ❌ DON'T DO THIS

1. **Don't use Letta for task management**: It's read-only and alpha quality
2. **Don't use Syncthing**: Conflict file hell with Obsidian
3. **Don't use Todoist integration**: Data lock-in, no better AI access
4. **Don't keep everything in Obsidian**: Specialized tools have superior UIs for specific functions
5. **Don't run multiple sync solutions**: Pick one (Git recommended)
6. **Don't ignore backups**: Sync ≠ backup
7. **Don't over-engineer early**: MVP first, scale later
8. **Don't skip testing**: Test restore, test conflicts, test automation

### Concrete Next Steps

**Week 1: Foundation**
1. Set up Git repo for Obsidian vault
2. Install Obsidian Git plugin, test on secondary device
3. Install Obsidian Tasks plugin, migrate existing tasks
4. Set up Nextcloud (self-hosted or paid hosting)
5. Create backup strategy (Restic + automated testing)

**Week 2: External Tools**
1. Deploy LinkWarden (Docker Compose)
2. Deploy Wallabag (Docker Compose)
3. Install Obsidian Web Clipper extension
4. Migrate bookmarks from Raindrop to LinkWarden
5. Test article workflow: save to Wallabag, read, sync to Obsidian

**Week 3: Automation**
1. Deploy n8n (Docker or n8n Cloud)
2. Install Post Webhook plugin in Obsidian
3. Create first workflow: Email → Obsidian capture
4. Build YouTube capture workflow (your homesteading example)
5. Test webhook security (Bearer tokens, HTTPS)

**Week 4: AI Integration**
1. Install Obsidian Local REST API plugin
2. Build Python script for task operations (read, mark complete, create)
3. Test with single agent instance
4. Implement Git commit after agent changes
5. Create monitoring (log all agent actions)

**Month 2: Scale and Refine**
1. Add distributed locking if running multiple agent instances
2. Build metadata database if vault \u003e 5,000 notes
3. Set up monitoring (Prometheus/Grafana)
4. Optimize automation workflows based on usage
5. Document entire system for future you

This architecture gives you data ownership, AI agent integration (properly implemented), scalable capture workflows, and a clear path from MVP to production scale. The hybrid approach balances Obsidian's simplicity for knowledge management with specialized tools where they excel, all while maintaining self-hosted data sovereignty.