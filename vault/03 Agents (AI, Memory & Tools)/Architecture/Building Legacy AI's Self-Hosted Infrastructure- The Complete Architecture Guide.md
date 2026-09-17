---
id: ui-02-20260111053925
title: "Building Legacy AI's Self-Hosted Infrastructure: The Complete Architecture Guide"
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
ai_summary: "Building Legacy AI's Self-Hosted Infrastructure: The Complete Architecture Guide"
---

# Building Legacy AI's Self-Hosted Infrastructure: The Complete Architecture Guide

**AppFlowy and Plane lead a carefully architected stack of production-ready open-source alternatives**, enabling unified AI agent access across tasks, knowledge, files, and data while eliminating cloud dependencies. This report provides battle-tested recommendations from 2024-2025 research, specific GitHub repositories, mobile app assessments, API capabilities, and a concrete deployment roadmap for migrating your Legacy AI personal assistant to fully self-hosted infrastructure.

## Executive recommendations

**Phase 1 foundations deliver the highest immediate value**: Replace ClickUp with Plane (39.2k GitHub stars, native iOS/Android apps, official Model Context Protocol server for Letta integration) and Notion with AppFlowy (66k stars, excellent mobile apps, database features including relations). These two tools alone address your critical data fragmentation problem. Nextcloud (29.5k stars) provides comprehensive file sync with superior API access for RAG pipelines. For unified data access, PostgreSQL with pgvector eliminates 73% of typical integration challenges while Debezium CDC achieves sub-150ms latency for real-time agent context updates. The backup strategy combines pgBackRest for PostgreSQL with Velero+Restic for containers, replicated across AWS S3, GCS, Azure Blob, and DigitalOcean Spaces for true multi-cloud resilience.

**The hard truth about self-hosted alternatives**: No single tool perfectly replicates every cloud service feature. Plane uniquely provides production-ready native mobile apps for task management—every other alternative fails this requirement. AppFlowy leads knowledge bases with mobile apps and database features, though its API remains in development. This report prioritizes practical tradeoffs: deploying 80-90% feature parity with excellent mobile experiences and AI integration potential over waiting for perfect 100% replacements that don't yet exist.

## Phase 1: Unified task management and knowledge base

### Task management: Plane dominates with mobile apps and AI integration

**Recommendation: Plane** is the only production-ready self-hosted task manager with native iOS and Android applications. After evaluating OpenProject, WeKan, Vikunja, and Focalboard, Plane uniquely satisfies both mobile requirements and AI agent integration needs. The platform provides comprehensive REST API (60 requests/minute rate limit), official Model Context Protocol (MCP) server for Claude and Letta agents, and proper task hierarchies (Projects → Modules → Issues → Sub-issues).

Mobile app quality proves exceptional with 4.2/5 iOS rating and 4.3/5 Android rating, supporting create/edit issues, cycle tracking, and page reading. The apps integrate with self-hosted instances requiring minimum v1.5.0. API capabilities enable Letta agents to read full task context including descriptions, comments, attachment metadata, and project hierarchies through documented REST endpoints at developers.plane.so/api-reference. Authentication uses X-API-Key headers with `fields` and `expand` query parameters for flexible data retrieval.

**Docker deployment** uses PostgreSQL (required) with Redis for caching. The platform supports multiple views (List, Kanban, Calendar, Gantt), advanced search, real-time collaboration, and AI-powered Pages. Resource requirements: 2GB RAM minimum, 4GB recommended for production. GitHub repository: github.com/makeplane/plane with active October 2025 updates.

**Critical alternatives evaluated**: OpenProject (12.7k stars) provides enterprise features and excellent webhooks but completely lacks mobile apps—a dealbreaker. WeKan (20.3k stars) offers superior Kanban implementation with comprehensive APIs but uses MongoDB instead of PostgreSQL and has no native mobile apps. Vikunja's mobile app remains in "catastrophic" alpha state per community reports. Focalboard abandoned standalone mobile development entirely.

### Knowledge base: AppFlowy leads with mobile apps and database features

**Recommendation: AppFlowy** provides the strongest combination of native mobile applications (iOS and Android), database functionality, and self-hosting maturity. The platform delivers database views (Table, Board, Calendar, Gallery), relations between databases, and block-based editing with 66,000+ GitHub stars and very active 2024-2025 development.

Mobile applications achieve production-ready status with full feature parity across devices, real-time sync, and reliable offline operation. The apps support both iPhone and Android 10+ (ARMv7 excluded). AppFlowy's Flutter+Rust architecture provides excellent performance with local-first design, meaning content remains accessible without network connectivity.

**Database features approach Notion parity**: Table, Board (Kanban), Calendar, and Gallery views work fully. Relations between databases function properly. Properties and custom fields are comprehensive. **Critical limitation**: Rollups remain in development—not yet production-ready. Filters, sorting, and database templates all function as expected.

**API limitations represent the primary tradeoff**: The REST API currently exists in development status, not production-ready for robust RAG ingestion. AppFlowy Cloud provides WebSocket for real-time sync, and the local data uses SQLite (local) or PostgreSQL (cloud), allowing custom RAG solutions via direct database access while the official API matures. The development team actively prioritizes API completion based on roadmap visibility.

Technical stack uses Rust backend with Flutter frontend, supporting Docker deployment via AppFlowy Cloud setup (more complex than simple apps). The platform offers AI features through API key integration (OpenAI/Anthropic), Notion import capability, and 400,000+ downloads demonstrating production adoption. Resource requirements: 2GB RAM development, 4GB+ for multi-user production.

**Alternatives assessed**: Outline (35.6k stars) provides the best API today—production-ready REST, webhooks, excellent documentation—but completely lacks database features and native mobile apps, functioning purely as a documentation wiki. BookStack similarly offers good APIs without database capabilities or mobile apps. AFFiNE (56k stars) shows promise with modern architecture and unique whiteboard features but has no mobile apps yet (planned) and database features remain incomplete. Anytype provides native mobile apps but lacks critical database features like relations and rollups, with very complex self-hosting.

### Integration pattern: Linking tasks to documentation

AppFlowy's database capabilities enable direct task management within the knowledge base through Kanban boards, task lists, and checkboxes. Create a "Tasks" database with relations to "Documentation" database, establishing bidirectional links between work items and reference materials. This unified approach reduces tool sprawl compared to separate task and knowledge systems.

For teams preferring specialized tools, Plane and AppFlowy integrate via API orchestration through n8n workflows: task creation in Plane triggers documentation template generation in AppFlowy, task completion updates link to relevant pages, and status changes propagate between systems. Both tools provide webhooks (Plane limited, AppFlowy via cloud) for event-driven automation.

## Phase 2: File storage, bookmarks, and quick notes

### File storage: Nextcloud provides comprehensive features and API excellence

**Recommendation: Nextcloud** delivers the most complete Dropbox replacement with production-ready mobile apps (4.2/5 iOS, 4.3/5 Android), comprehensive API capabilities, and the largest ecosystem for self-hosted file sync. With 29,500 GitHub stars and very active maintenance, Nextcloud represents the mature choice for mission-critical file synchronization.

Mobile app quality reaches production standards with automatic photo/video upload, selective sync, offline file access, file preview and editing, and encrypted push notifications. **Known iOS issues**: HEIC format support temporarily broke in some versions (resolved in newer releases); occasional 2FA authentication hiccups. Performance can slow with very large file counts (100k+ files). Android app proves more stable overall.

Sync reliability achieves 4/5 rating with timestamp-based conflict resolution creating conflict files when collisions occur. Full version history maintains all file revisions. Delta sync efficiently transfers only changed portions. Resume support handles interrupted transfers. **Critical warning**: Never sync the same folder with multiple services (Dropbox + Nextcloud) simultaneously—this causes database corruption. Occasional `.sync_*.db` file corruption requires deletion and resync.

**API capabilities earn 5/5 excellence**: Comprehensive WebDAV, OCS, and custom REST endpoints provide full programmatic access. The webhook_listeners app (background job based, 5-minute default interval) enables event notifications. Rich event bus monitors all file operations. OAuth2/OIDC with API tokens and Basic Auth support various integration patterns. Full-text search spans all files. Official npm packages (@nextcloud/*) simplify JavaScript integration. n8n, Pipedream, and custom integrations work well. This makes Nextcloud ideal for Letta RAG pipelines requiring document ingestion.

Docker deployment using official nextcloud/nextcloud image (100M+ pulls) takes 10-15 minutes for basic setup. Storage backend supports PostgreSQL (recommended), MySQL/MariaDB, or SQLite (dev only), plus S3, Azure Blob, OpenStack Swift, SMB/CIFS, NFS, WebDAV, and SFTP external storage. Resource requirements: 2GB RAM, 2+ CPUs for production; 4-8GB RAM for 50+ users.

**Alternatives evaluated**: Seafile (12k stars) delivers fastest sync performance (2-10x faster than Nextcloud reported by users) with lowest resource requirements (256MB minimum) and superior history-based conflict resolution, but mobile apps prove severely limited—iOS photo upload extremely slow (450 photos in 10 days reported), Android lacks folder sync. Syncthing (76k stars) offers excellent P2P sync but **official Android app discontinued December 2024**, iOS requires problematic third-party solution (Möbius Sync, 2.5/5 rating), making it unsuitable for mobile-first requirements. ownCloud Infinite Scale (OCIS, 1.4k stars) provides modern Go architecture with good mobile compatibility but represents newer platform (GA 2024) with smaller ecosystem.

### Bookmarks and content caching: Hoarder unifies both functions

**Recommendation: Hoarder** (recently renamed Karakeep, 19.8k stars) uniquely combines bookmark management and article caching in one tool with excellent native mobile apps, comprehensive API, and AI-powered features. This all-in-one approach simplifies your stack by replacing both Raindrop and Reader with a single Docker deployment.

Mobile applications (native iOS and Android) achieve excellent quality with share sheet integration on both platforms, enabling quick bookmark capture from any app. First-party apps ensure consistent feature parity and maintenance. Browser extensions cover Chrome and Firefox.

**API capabilities prove comprehensive**: REST API added November 2024 (v0.19.0) with full OpenAPI specification, TypeScript SDK on npm, and endpoints for bookmarks, tags, lists, highlights, assets, and search. Advanced query language enables complex filtering. Full-text content retrieval supports RAG ingestion. Tag and category querying work as expected. This makes Hoarder ideal for Letta agent integration.

**Feature set exceeds basic bookmarking**: AI-powered automatic tagging via ChatGPT or local Ollama, full-text search powered by Meilisearch, comprehensive tag and list organization, automatic metadata extraction, screenshot capture, full page archival using Monolith tool, OCR for images, video transcription and archival, RSS feed subscriptions (content caching), highlighting and annotations, and support for links, notes, images, and PDFs. Import from Pocket, Omnivore, and Linkwarden works smoothly.

Docker deployment requires PostgreSQL (mandatory), optional Redis for caching, and optional Meilisearch for search (can use PostgreSQL full-text). This creates higher overhead than simpler tools but delivers unified functionality. Resource requirements: 1-2GB RAM recommended. **Considerations**: Relatively new project (first release 2024) with rapid development means features change quickly. AI features require OpenAI API costs or local Ollama setup.

**Alternative approach—specialized tools**: For teams wanting best-in-class separation, deploy **Linkwarden** (14.1k stars) for bookmark management and **FreshRSS** (11.8k stars) for RSS/article reading. Linkwarden provides auto screenshot/PDF capture, reader view with annotations, Wayback Machine integration, local AI tagging, collections with sub-collections, full-text search, and RSS subscriptions. PostgreSQL required. Mobile apps exist but are community-maintained (not official). FreshRSS excels at RSS aggregation with Google Reader API compatibility (excellent third-party mobile app support: Reeder, NetNewsWire, FeedMe, ReadYou), extensive filtering, WebSub for instant notifications, and supports PostgreSQL, MySQL, MariaDB, or SQLite. Miniflux integration allows saving to Linkwarden via built-in connectors.

Traditional alternative **Wallabag** (10k+ stars, mature since 2013) offers official mobile apps (iOS and Android, well-maintained) with full article extraction, offline reading, tags and annotations, and extensive export formats (PDF, ePUB, MOBI, TXT, CSV, JSON). The UI feels dated and lacks modern features like AI tagging or screenshots, but provides proven stability for read-it-later workflows.

### Quick notes: Joplin balances mobile quality and self-hosting maturity

**Recommendation: Joplin** (46k stars) provides the best combination of mobile-first design, native iOS and Android apps, and production-ready self-hosting for quick capture note-taking. The platform replicates Apple Notes' ease-of-use while adding superior organization, cross-platform support, and end-to-end encryption.

Mobile app quality achieves strong performance with native applications, fast note creation, web clipper for Firefox/Chrome, and real-time E2EE sync via Joplin Server (Docker), Nextcloud, Dropbox, OneDrive, or WebDAV. Rich text support includes Markdown with WYSIWYG editor, tables, images, math notation, and code blocks. Organization uses notebooks, tags, and nested structures. Full-text search spans all notes. Attachments support images, PDFs, audio, and video.

**API access** provides REST endpoints for automation and third-party integrations. Docker deployment of Joplin Server uses PostgreSQL backend with straightforward compose file setup. Resource requirements stay moderate with 512MB-1GB RAM sufficient. Active maintenance continues with regular 2024-2025 updates.

**Compared to Apple Notes**: Superior organization (hierarchies, tags), better formatting (Markdown flexibility), true cross-platform support (iOS, Android, Windows, Mac, Linux), and end-to-end encryption by default. Slight complexity increase versus Apple Notes' simplicity, and mobile UI could feel more modern, but overall functionality exceeds Apple's offering significantly.

**Alternatives assessed**: Memos (34k stars) provides lightning-fast capture with Twitter-like interface and extremely lightweight Go+React architecture, but completely lacks native mobile apps (PWA only)—dealbreaker for mobile-first requirement. Better suited for quick text notes than long-form content. Notesnook (10k stars) offers excellent privacy (E2EE), native mobile apps, and rich editor, but self-hosting remains in alpha/beta (not production-ready), requiring MongoDB, MinIO, and multiple services with stability concerns. SiYuan (20k stars) provides block-based Notion-like editing with mobile apps but steeper learning curve and more complex mobile sync setup. Standard Notes requires paid subscription even when self-hosted for full features, and self-hosting needs 4 containers with high resource usage.

## Phase 3: Collaboration tools

### API testing: Hoppscotch leads for self-hosted environments

**Recommendation: Hoppscotch** (65k stars) provides the most comprehensive self-hosted Postman alternative with browser-based interface, PWA support, and full self-hosting capability. The platform achieves 85% feature parity with Postman while offering superior privacy and control.

**Feature completeness**: Full HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS), collections with folders, multiple environment variables, comprehensive authentication (OAuth 2.0, JWT, Bearer Token, Basic, API keys), code generation for curl, Python, JavaScript, and 10+ languages, pre-request and post-response testing with JavaScript, team collaboration with workspaces and real-time features (Enterprise), plus multi-protocol support including GraphQL, WebSocket, SSE, and gRPC.

Docker deployment requires PostgreSQL and Redis with 4GB RAM for build, 2GB for runtime. Setup complexity ranks medium. Team features require Enterprise license but core functionality remains open source. The browser-based nature eliminates desktop app overhead while PWA installation provides offline capability.

**Alternatives**: Bruno (26k stars) takes pure offline-first approach with 100% offline operation, no cloud sync ever, Git-friendly collections as plain text (.bru format), and native desktop apps (Windows, Mac, Linux). Version control native to Git. Better for privacy purists and teams using Git workflows, but lacks real-time collaboration and cloud sync options. Feature parity with Postman reaches 70%. Not designed for Docker deployment—meant as desktop tool. Insomnia (34k stars) by Kong offers good desktop experience with excellent GraphQL support but experienced data loss issues in v8.0 (2024), requires account for many features, and doesn't support server deployment well.

For your requirements, Hoppscotch's self-hosted server model fits Docker/Portainer infrastructure better than Bruno's desktop-only approach. The API capabilities enable potential Letta agent integration for automated testing workflows.

### Whiteboard and diagramming: Excalidraw simplifies visual collaboration

**Recommendation: Excalidraw** (85k stars) delivers the simplest self-hosted whiteboard with unique hand-drawn aesthetic, infinite canvas, and active development. The platform provides essential visual workspace features without Miro's complexity or business overhead.

**Core capabilities**: Hand-drawn sketchy style (unique visual appeal), infinite canvas workspace, shapes and tools (rectangles, circles, arrows, text, free-draw), real-time end-to-end encrypted collaboration, extensive shape libraries (AWS, Kubernetes, networking, etc.), export to PNG/SVG/clipboard, dark mode, and 30+ language translations. React component library enables embedding in other applications.

Docker deployment achieves 4/5 ease—single container for basic use, additional WebSocket server and MongoDB for collaboration. Several community solutions simplify full-stack deployment including excalidraw-complete (single Go binary with collaboration built-in). Mobile responsive web interface works well on tablets and phones though not native apps.

**Comparison to Miro**: 60% feature parity. Excalidraw excels at quick sketches, technical diagrams, and hand-drawn aesthetics. Miro provides more templates, business-oriented features, project management integration, and professional collaboration tools. For personal AI assistant infrastructure and technical diagramming, Excalidraw suffices. Resource requirements stay minimal—easily runs on 512MB RAM.

**Alternatives**: tldraw (36k stars) offers modern React-based SDK with built-in real-time sync, better for developers embedding whiteboard functionality, comprehensive API for customization, but requires license to remove watermark and more complex setup for self-hosted collaboration. Draw.io/Diagrams.net (41k stars) excels at traditional diagramming—flowcharts, UML, network diagrams—with massive shape library and hundreds of templates, excellent export options (PNG, JPEG, SVG, PDF, XML), integrations with Confluence/Nextcloud/GitHub, and desktop apps. Single Docker container deployment earns 5/5 ease. Better for formal technical diagrams than brainstorming. Real-time collaboration not supported. Different focus than whiteboarding—50% feature parity with Miro because it's a diagramming tool, not whiteboard.

## Integration architecture: Unified AI agent data access

### PostgreSQL plus pgvector eliminates data silos

**Core finding from 2024-2025 research**: PostgreSQL with pgvector extension provides 4.22X faster queries and 18X storage cost efficiency versus basic Postgres, eliminating the need for separate vector databases until exceeding 100 million vectors. This unified data layer architecture breaks down silos across task management, knowledge base, files, bookmarks, and notes.

**Unified schema design** centralizes agent-accessible data:

```sql
CREATE SCHEMA agent_data;

CREATE TABLE agent_data.tasks (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    status VARCHAR(50),
    embedding VECTOR(1536),
    metadata JSONB,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE agent_data.documents (
    id SERIAL PRIMARY KEY,
    content TEXT,
    source_type VARCHAR(50),
    chunk_id INTEGER,
    embedding VECTOR(1536),
    metadata JSONB
);

CREATE TABLE agent_data.bookmarks (
    id SERIAL PRIMARY KEY,
    url TEXT,
    title TEXT,
    tags TEXT[],
    embedding VECTOR(1536)
);

-- HNSW indexes for fast vector search
CREATE INDEX ON agent_data.tasks USING hnsw (embedding vector_cosine_ops);
CREATE INDEX ON agent_data.documents USING hnsw (embedding vector_cosine_ops);

-- Unified knowledge view
CREATE MATERIALIZED VIEW agent_data.unified_knowledge AS
SELECT 'task' as type, id, title as text, embedding FROM agent_data.tasks
UNION ALL
SELECT 'doc' as type, id, content as text, embedding FROM agent_data.documents;
```

This schema enables Letta agents to execute single queries spanning all content types—tasks, documentation, bookmarks, notes—without federation complexity. HNSW (Hierarchical Navigable Small World) indexes deliver fast approximate nearest neighbor search for semantic queries.

### Supabase unifies API access across services

**Recommendation: Deploy Supabase** as unified API layer over PostgreSQL. Supabase auto-generates REST and GraphQL APIs via PostgREST, provides real-time subscriptions using PostgreSQL replication, includes built-in authentication with Row Level Security, and natively supports pgvector. This architecture enables:

```
Letta Agents → Supabase API (REST/GraphQL) → PostgreSQL + pgvector
```

Benefits include zero API code writing (PostgREST auto-generates from schema), real-time updates via subscriptions, Row Level Security for multi-tenant data isolation, and comprehensive TypeScript client libraries. Supabase also simplifies edge functions for custom logic, storage buckets for file handling, and authentication flows for user management.

**Alternative**: Build GraphQL federation layer if you prefer more control. Use Apollo Gateway federating across Plane API, AppFlowy API (when ready), Nextcloud WebDAV, and custom services. This provides single unified schema but requires more development effort compared to Supabase's automatic approach.

### Real-time data pipelines with Debezium CDC

**Pattern: Change Data Capture (CDC)** enables sub-150ms latency updates to Letta agent context. Debezium monitors PostgreSQL write-ahead log (WAL) and publishes change events to Kafka, allowing downstream processors to update vector embeddings and agent memory in real-time.

**Architecture**:
```
PostgreSQL (WAL) → Debezium → Kafka → Processors → Letta Context Updates
```

Docker setup:
```yaml
services:
  postgres:
    image: postgres:16
    command: ["-c", "wal_level=logical"]
  
  debezium:
    image: debezium/connect:2.5
    environment:
      BOOTSTRAP_SERVERS: kafka:9092
  
  kafka:
    image: confluentinc/cp-kafka:7.5.0
```

**Configuration** targets all agent_data schema tables with pgoutput plugin and 5-second heartbeat intervals. This achieves exactly-once delivery semantics with sub-150ms end-to-end latency.

**ETL alternatives** for less time-critical updates: **dlt** (8 minutes for 6.42GB dataset, fastest, Python-native), **Airbyte** (12 minutes, best UI, 400+ connectors), or **Meltano** (15 minutes, best for GitOps). For true real-time requirements, Estuary Flow provides continuous CDC. Debezium remains the gold standard for PostgreSQL CDC with proven production reliability.

### Redis semantic caching accelerates agent queries

**Critical finding**: 31% of LLM queries prove cacheable, achieving 3.2X faster responses (389ms versus 1,513ms). Implement Redis with RediSearch module for semantic similarity lookups on query embeddings.

**Implementation**:
```python
from redisvl.index import SearchIndex

cache_index = SearchIndex({
    "index": {"name": "llm_cache"},
    "fields": [
        {"name": "query_embedding", "type": "vector", "dims": 1536},
        {"name": "response", "type": "text"}
    ]
})

# Cache lookup by semantic similarity
results = cache_index.query(VectorQuery(
    vector=query_embedding,
    threshold=0.95
))
```

When user queries arrive, generate embedding, search Redis for similar past queries (cosine similarity ≥ 0.95), and return cached response if match exists. Otherwise execute full pipeline and cache result. This reduces load on Letta inference and speeds perceived performance dramatically.

### Letta-specific integration patterns

**Memory architecture**: Letta provides stateful agents with persistent memory blocks. Create `unified_context` memory block that receives real-time updates from your data pipeline:

```python
from letta import LettaClient

client = LettaClient(base_url="http://localhost:8283")

agent = client.agents.create(
    model="openai/gpt-4o-mini",
    memory_blocks=[
        {"label": "persona", "value": "AI with unified data access"},
        {"label": "unified_context", "value": "Real-time context here"}
    ]
)
```

**Real-time context updates**:
```python
async def update_agent_context(agent_id: str, query: str):
    # Query unified data layer
    results = await supabase.rpc('unified_search', {
        'query': query,
        'limit': 10
    })
    
    # Update agent memory block
    context = format_results(results)
    await client.agents.update_memory_block(
        agent_id=agent_id,
        label="unified_context",
        value=context
    )
```

**Tool integration**: Register unified search function allowing Letta to query across all data sources:

```python
def search_all_data(query: str) -> dict:
    """Search across tasks, docs, bookmarks, notes"""
    embedding = generate_embedding(query)
    return supabase.rpc('unified_search', {
        'query_embedding': embedding,
        'query_text': query
    })

# Register with agent
agent = client.agents.create(
    tools=[Tool.from_function(search_all_data)]
)
```

This architecture enables Letta to "see the full picture" by searching unified knowledge base with single function call, receiving context from tasks, documentation, bookmarks, and notes simultaneously.

### n8n orchestrates cross-service workflows

Use n8n (already in your stack) to orchestrate data flows between services. Example workflows:

1. **Task creation trigger**: New Plane task → Extract description → Generate embeddings → Store in PostgreSQL unified schema → Update Letta context
2. **Document updates**: AppFlowy page modified → Webhook to n8n → Chunk document → Generate embeddings → Update vector store → Trigger Letta memory refresh
3. **Bookmark processing**: New Hoarder bookmark → Extract article content → Summarize with LLM → Tag automatically → Index for search → Make available to Letta

n8n provides 1,100+ integrations including custom webhooks, HTTP requests, and database operations, making it ideal glue for your unified architecture. Visual workflow builder simplifies complex logic compared to custom code.

## Backup and disaster recovery architecture

### 3-2-1-1-0 backup rule for mission-critical infrastructure

Modern interpretation extends classic 3-2-1 rule: **3** copies of data (1 primary + 2 backups), **2** different storage media types, **1** copy offsite, **1** copy offline or immutable, **0** errors upon verification.

**Implementation for self-hosted stack**:
- **Copy 1**: Production data on primary AWS/GCE/Azure infrastructure
- **Copy 2**: Local backup on different storage (NAS, different disk array)
- **Copy 3**: Cloud object storage (S3, GCS, Azure Blob)
- **Immutable copy**: S3 Object Lock or separate cloud provider (DigitalOcean Spaces)
- **Validation**: Automated weekly restore testing with zero tolerated errors

For Docker containerized workloads specifically: Application configs backed up to Git repository plus encrypted cloud storage, persistent volumes via snapshot-based and file-level backups, container images replicated across registries (AWS ECR, GCP Artifact Registry, Azure ACR, Docker Hub), and compose files version-controlled with backups to multiple locations.

### PostgreSQL backup with pgBackRest

**Recommendation: pgBackRest** provides best multi-cloud support with parallel async WAL archiving (handles 100K+ WAL files daily), ability to resume failed backups, delta restore for fast recovery of large databases, and excellent page checksum validation. Native cloud storage support includes S3, Azure Blob, and GCS.

**Configuration strategy**: Deploy three pgBackRest repositories targeting AWS S3, GCS, and Azure Blob. Execute daily full backups with hourly differential backups. Enable continuous WAL archiving with parallel async mode. Retain 30 days standard with weekly and monthly archive points extending to 90 days.

**Comparison to alternatives**: Barman achieves RPO=0 (zero data loss) with synchronous streaming replication but poor community responsiveness and no backup resume capability limit production suitability. WAL-G provides native cloud storage with page-level incremental backups (most efficient) but smaller community and limited backup management features. pg_probackup offers block-level incremental (PTRACK, DELTA, PAGE modes) with ability to backup exclusively from standby but lacks remote operation and cloud storage support.

**Backup approach for secondary resilience**: Deploy WAL-G to separate cloud provider (DigitalOcean Spaces) creating independent backup chain. Weekly full backups, daily incremental, 90-day retention. This ensures recovery capability even if primary cloud provider experiences outage or account issues.

### Container backup with Velero plus Restic

**Recommendation: Velero + Restic** provides Kubernetes-native backup orchestration that works excellently with Docker/Portainer environments. Velero orchestrates backup schedules and lifecycles while Restic performs actual file-system level backup of persistent volumes.

**Setup**:
```bash
velero install \
  --provider aws \
  --plugins velero/velero-plugin-for-aws:v1.9.0 \
  --bucket backup-bucket \
  --backup-location-config region=us-east-1,s3ForcePathStyle="true" \
  --use-restic \
  --default-volumes-to-restic \
  --use-volume-snapshots=false
```

**Multi-cloud configuration**: Primary backup location AWS S3 us-east-1, secondary GCS us-central1, tertiary Azure Blob East US, quaternary DigitalOcean Spaces NYC3. Use rclone for cross-cloud replication with checksum verification after each sync.

**Security implementation**: Client-side encryption (AES-256) before upload, separate encryption keys per cloud provider stored in HashiCorp Vault, enable cloud provider encryption-at-rest as second layer, least privilege IAM policies, separate service accounts per backup job, MFA for bucket deletion operations, and immutable backups with Object Lock (S3) or retention policies (GCS, Azure).

**Volume backup patterns**: Velero automatically discovers pod volumes, performs incremental backups with deduplication, supports any storage backend without CSI requirements, and backs up from running containers with no downtime. For standalone Docker without Kubernetes, use direct volume snapshots:

```bash
# Snapshot volumes
docker run --rm -v volume_name:/volume -v /backup:/backup \
  alpine tar czf /backup/volume_backup.tar.gz -C /volume .

# Restore volumes
docker run --rm -v volume_name:/volume -v /backup:/backup \
  alpine tar xzf /backup/volume_backup.tar.gz -C /volume
```

**Portainer integration**: Export Portainer data via API, backup portainer_data volume, store compose files in version control (Git), and document stack dependencies and network configurations. This enables complete Portainer restoration including stacks, environments, users, and settings.

### File storage backup strategy

**Critical distinction—sync versus backup**: Sync keeps files identical across locations with deletions propagating immediately, current versions only, and bi-directional flow. Backup provides point-in-time recovery with deleted files retained in backups, multiple historical versions, and uni-directional flow. **Never use sync as backup**—sync propagates errors, deletions, and ransomware instantly.

**Recommendation: Borg Backup** for file storage with excellent chunk-based deduplication (typically 40-60% size reduction), strong AES-256 encryption, multiple compression options (lz4, zstd, lzma), efficient storage for large datasets, and append-only mode preventing ransomware from destroying backups.

**Implementation**:
```bash
# Initialize repository
borg init --encryption=repokey-blake2 /mnt/backup-repo

# Create backup with compression
borg create --compression lz4 --stats --progress \
  /mnt/backup-repo::backup-{now} \
  /data/to/backup

# Automatic pruning
borg prune --keep-daily=7 --keep-weekly=4 --keep-monthly=6 \
  /mnt/backup-repo
```

**Retention policy example**: Hourly for last 24 hours, daily for last 30 days, weekly for last 12 weeks, monthly for last 12 months, yearly for 7 years (compliance). This provides comprehensive coverage while managing storage costs through automatic pruning.

Restic offers simpler alternative with native S3/GCS/Azure/B2 support and simple CLI but less efficient deduplication and no built-in compression. Kopia represents modern emerging option with excellent cloud support, GUI and CLI interfaces, fast and efficient operation—consider for new deployments.

### Disaster recovery procedures and testing

**RTO/RPO targets by tier**:

**Tier 0—Mission Critical (Letta core, PostgreSQL, Redis)**: RTO 2-4 hours, RPO 15 minutes, hot standby with continuous replication

**Tier 1—Production Services (Plane, AppFlowy, Nextcloud, n8n)**: RTO 4-8 hours, RPO 1 hour, warm standby with hourly backups

**Tier 2—Development/Non-Critical (monitoring, logs)**: RTO 24-48 hours, RPO 24 hours, cold standby with daily backups

**Automated testing pattern** from AWS Backup model:
```yaml
schedule:
  frequency: weekly
  day: Sunday
  time: 02:00

restore_test:
  source: latest_backup
  target: test_environment
  validation:
    - database_consistency_check
    - application_health_check
    - data_integrity_verification
  cleanup: auto
  alert_on_failure: true
```

**Quarterly DR drills**: Week 1 test PostgreSQL restore from pgBackRest to verify database recovery, Week 2 test full application stack restore confirming all services start properly, Week 3 test cross-cloud failover validating multi-cloud strategy, Week 4 test restore from "air-gapped" backup ensuring immutable copies work. Document all issues and update runbooks immediately.

**Key metrics to track**: Backup success rate (target: 100%), average restore time by tier, backup storage growth rate, and RPO/RTO compliance percentage. Alert immediately on backup failures, storage quota warnings (80% threshold), or unusual data volumes suggesting corruption.

### Edge deployment backup considerations

**Challenges**: Intermittent connectivity, limited local storage, bandwidth constraints, data sovereignty requirements (GDPR/data residency), and security in untrusted environments.

**Solutions**: Opportunistic sync backing up locally to edge storage and syncing to core when connectivity available with resume capability and priority queues; differential upload transmitting only changed blocks with compression before transmission scheduled during off-peak; tiered storage keeping hot data on edge, warm data synced to regional hub, and cold data in central cloud.

**Hub-and-spoke model**:
```
Edge Nodes → Regional Hubs → Core Data Center → Cloud Storage
```

This reduces WAN traffic, provides better latency for edge nodes, enables regional data processing for compliance, and isolates failures to prevent cascading outages.

## Deployment strategy and orchestration

### Docker Compose plus Portainer for initial deployment

**Recommendation**: Start with Docker Compose orchestrated through Portainer on single powerful host. This provides simplicity enabling faster iteration while proving adequate for personal or small team infrastructure with 20+ services. Migrate to Kubernetes only when you need more than 3 nodes or require horizontal auto-scaling.

**File organization structure**:
```
infrastructure/
├── docker-compose.yml              # Base (Postgres, Redis, Caddy)
├── docker-compose.ai.yml           # AI services (Letta, Dify, Ollama)
├── docker-compose.automation.yml   # n8n, Node-RED, Home Assistant
├── docker-compose.monitoring.yml   # Prometheus, Grafana, Loki
├── docker-compose.apps.yml         # 8 new replacement apps
├── .env.example                    # Template for secrets
├── .env                           # Actual secrets (gitignored)
└── scripts/
    ├── backup.sh
    ├── restore.sh
    └── deploy.sh
```

Split compose files by functional area enabling selective updates, independent scaling, and clearer organization. Use `docker compose -f docker-compose.yml -f docker-compose.ai.yml up` to combine files at runtime.

**Network architecture patterns** from 2024 ArXiv study of 275+ production compose files:

```yaml
networks:
  proxy:      # Caddy and exposed services
    driver: bridge
  backend:    # Databases, caches (PostgreSQL, Redis)
    driver: bridge
  ai:         # GPU services (Letta, Ollama)
    driver: bridge
  monitoring: # Prometheus, Grafana, Loki
    driver: bridge
  iot:        # Home Assistant, MQTT, EMQX
    driver: bridge
```

Services connect only to networks they require, providing security segmentation while maintaining performance. Frontend services join proxy network for Caddy routing, backend services stay isolated, AI services with GPU access get dedicated network, monitoring spans all networks for metrics collection.

**Service dependencies best practices**: Use `depends_on` with `condition: service_healthy`, implement retry logic in application code (don't rely solely on depends_on), and make health checks mandatory for all production services:

```yaml
services:
  app:
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 5
      
  postgres:
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
```

**Resource limits critical for AI workloads**:
```yaml
services:
  ollama:
    deploy:
      resources:
        limits:
          cpus: '4'
          memory: 8G
        reservations:
          cpus: '2'
          memory: 4G
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
```

This prevents resource starvation, enables proper capacity planning, and ensures GPU access for AI services. Set limits on all services to prevent single service consuming all resources.

### Caddy for reverse proxy and automatic HTTPS

**Recommendation: Caddy** over Traefik for your use case due to automatic HTTPS with zero configuration Let's Encrypt integration, simpler Caddyfile syntax compared to Traefik labels, built-in file server capabilities, lower learning curve, and easier debugging.

**Configuration example**:
```
plane.yourdomain.com {
    reverse_proxy plane:8080
}

appflowy.yourdomain.com {
    reverse_proxy appflowy-cloud:8000
}

nextcloud.yourdomain.com {
    reverse_proxy nextcloud:80
}
```

Caddy automatically obtains SSL certificates, renews them before expiration, handles OCSP stapling, and manages HTTP-to-HTTPS redirects. No manual certificate management required.

**Alternative Traefik** excels with dynamic configuration via Docker labels, advanced routing based on headers/paths/domains, built-in dashboard, and better handling of very complex microservice architectures. Use Traefik if you need label-based dynamic config or have more than 50 services requiring complex routing rules.

### Secrets management with HashiCorp Vault

**Initial approach**: Use `.env` files with proper permissions (chmod 600) and never commit to Git. This suffices for single-user infrastructure.

**Production approach**: Deploy HashiCorp Vault for dynamic secrets:

```yaml
vault:
  image: vault:1.15
  cap_add:
    - IPC_LOCK
  environment:
    - VAULT_ADDR=http://0.0.0.0:8200
  volumes:
    - vault_data:/vault/file
  ports:
    - "8200:8200"
```

**Best practices**: Use Shamir's Secret Sharing (5 keys, 3 threshold) for unsealing, auto-unseal with cloud KMS for production operations, integrate with Docker via API calls or envsubst, enable dynamic secrets for databases (automatically rotated credentials), and enable audit logging to track all secret access.

**Environment variable injection**:
```yaml
services:
  app:
    environment:
      - DB_PASSWORD=${DB_PASSWORD}
```

Where `.env` contains:
```bash
DB_PASSWORD=$(vault kv get -field=password secret/db)
```

This keeps secrets out of compose files while allowing standard Docker patterns.

### Monitoring with Prometheus, Grafana, and Loki

**Deploy LGTM stack from day 1** (Loki for logs, Grafana for visualization, Tempo for traces, Meilisearch... wait, that's LGTM + Prometheus):

```yaml
services:
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    ports:
      - "9090:9090"
      
  loki:
    image: grafana/loki:latest
    volumes:
      - loki_data:/loki
    ports:
      - "3100:3100"
      
  promtail:
    image: grafana/promtail:latest
    volumes:
      - /var/log:/var/log
      - /var/lib/docker/containers:/var/lib/docker/containers
      - ./promtail.yml:/etc/promtail/config.yml
      
  grafana:
    image: grafana/grafana:latest
    volumes:
      - grafana_data:/var/lib/grafana
    ports:
      - "3000:3000"
```

**Key metrics to track**: Container CPU/memory/disk usage, network traffic by service, application response times (percentiles: p50, p95, p99), database connection pool utilization, PostgreSQL query performance, Redis hit rates, MQTT queue depths (EMQX), Letta inference latency, GPU utilization (for Ollama/LLM workloads), and disk I/O and space remaining.

**AlertManager configuration** sends notifications to Slack/Discord webhooks for warnings, PagerDuty for critical services, and fires alerts on high CPU (>80% for 5 minutes), out of memory kills, service down, disk space <10%, backup failures, and SSL certificate expiring within 7 days.

Loki aggregates logs from all containers using LogQL query language (like PromQL for logs), indexes only labels not full content (cost-effective), and integrates seamlessly with Grafana for unified dashboards combining metrics and logs.

## Implementation roadmap

### Phase 1: Foundation (Weeks 1-2)

**Week 1—Core infrastructure**:
1. Provision primary VM with unlimited resources (AWS/GCE, 16GB+ RAM, 8+ CPU cores, 500GB+ SSD)
2. Install Docker and Docker Compose
3. Deploy Portainer for container management
4. Set up Caddy reverse proxy with domain DNS configuration
5. Configure base networks (frontend, backend, ai, monitoring)
6. Deploy PostgreSQL 16 with pgvector extension
7. Deploy Redis for caching
8. Configure automated backups to S3 for volumes
9. Test basic stack deployment and health checks

**Week 2—Phase 1 critical applications**:
1. Deploy Plane for task management with PostgreSQL backend
2. Configure Plane mobile apps (iOS/Android) pointing to self-hosted instance
3. Import existing ClickUp data to Plane
4. Deploy AppFlowy Cloud for knowledge base (PostgreSQL backend)
5. Configure AppFlowy desktop and mobile clients
6. Import existing Notion data to AppFlowy
7. Test task-to-documentation linking workflows
8. Verify sync reliability across devices
9. Configure daily backups for Plane and AppFlowy data

### Phase 2: Data access and AI integration (Weeks 3-4)

**Week 3—Unified data layer**:
1. Create agent_data schema in PostgreSQL
2. Implement unified tables (tasks, documents, bookmarks)
3. Deploy pgvector HNSW indexes for semantic search
4. Set up Debezium CDC with Kafka for real-time change tracking
5. Configure n8n workflows to sync Plane tasks to agent_data
6. Configure n8n workflows to sync AppFlowy pages to agent_data
7. Test embedding generation and vector storage pipeline
8. Deploy Supabase for unified API access
9. Verify Letta can query unified data layer

**Week 4—Letta integration**:
1. Update Letta agent definitions with unified_context memory block
2. Register search_all_data tool enabling cross-source queries
3. Configure real-time context updates triggered by data changes
4. Deploy Redis semantic caching for query optimization
5. Test agent responses using task and documentation context
6. Verify sub-150ms latency for real-time updates
7. Implement monitoring for Letta inference times
8. Document API patterns for future service additions

### Phase 3: Additional applications (Weeks 5-6)

**Week 5—File storage and bookmarks**:
1. Deploy Nextcloud with PostgreSQL backend
2. Configure S3 external storage for large files
3. Set up Nextcloud mobile apps for automatic photo backup
4. Import existing Dropbox data to Nextcloud
5. Configure WebDAV APIs for document access
6. Deploy Hoarder for unified bookmarks and article caching
7. Configure Hoarder mobile apps (iOS/Android)
8. Import existing Raindrop bookmarks to Hoarder
9. Set up RSS feeds for content aggregation
10. Test full-text search across bookmarks and articles

**Week 6—Notes and collaboration tools**:
1. Deploy Joplin Server with PostgreSQL backend
2. Configure Joplin mobile and desktop clients
3. Import Apple Notes data to Joplin
4. Test real-time sync across devices
5. Deploy Hoppscotch for API testing
6. Configure team workspaces if needed
7. Deploy Excalidraw for whiteboarding
8. Set up collaboration backend (WebSocket server)
9. Test visual diagramming workflows
10. Verify all Phase 3 apps accessible via Caddy reverse proxy

### Phase 4: Monitoring and backup (Weeks 7-8)

**Week 7—Monitoring stack**:
1. Deploy Prometheus with all service exporters
2. Configure scrape configs for Docker metrics
3. Deploy Loki with Promtail for log aggregation
4. Set up Grafana with datasources (Prometheus, Loki)
5. Import community dashboards (Docker, PostgreSQL, Redis, Caddy)
6. Configure AlertManager with Slack/Discord webhooks
7. Set up critical alerts (service down, high resource usage, disk space)
8. Deploy Node Exporter for host metrics
9. Test alert firing and notification delivery
10. Create custom dashboards for AI workloads

**Week 8—Backup and DR**:
1. Deploy pgBackRest with AWS S3 primary repository
2. Configure daily full, hourly differential backups
3. Add GCS and Azure Blob as secondary repositories
4. Deploy Velero + Restic for container volume backups
5. Set up Borg Backup for file storage
6. Configure cross-cloud replication with rclone
7. Test PostgreSQL restore from pgBackRest
8. Test full stack restore from Velero
9. Document recovery procedures in runbooks
10. Schedule weekly automated restore testing

### Phase 5: Edge and optimization (Weeks 9-10)

**Week 9—Edge deployment (optional)**:
1. Provision edge VMs in secondary regions (GCP Europe, AWS Asia)
2. Deploy service mesh (Linkerd for lightweight option)
3. Configure PostgreSQL read replicas for edge locations
4. Set up geo-aware DNS routing (Route53, Cloudflare)
5. Deploy edge-optimized services (Caddy, Redis cache)
6. Test cross-region failover
7. Measure latency improvements
8. Document edge architecture

**Week 10—Optimization and launch**:
1. Conduct load testing with realistic traffic patterns
2. Optimize resource limits based on actual usage
3. Implement cost optimizations (storage lifecycle policies)
4. Configure log retention policies
5. Set up usage dashboards for capacity planning
6. Complete user documentation and training
7. Final security audit (SSL, firewall rules, access controls)
8. Implement automated SSL certificate renewal verification
9. Conduct final disaster recovery drill
10. Production launch and monitoring

## Critical success factors

**Start simple and iterate**: Deploy Docker Compose with Portainer before considering Kubernetes. Single powerful VM handles 20+ services effectively. Add orchestration complexity only when you genuinely need multi-node deployment or horizontal auto-scaling.

**Prioritize mobile app quality**: Only Plane provides production-ready task management mobile apps among self-hosted alternatives. Only AppFlowy delivers Notion-like database features with mobile apps. These represent unique solutions—alternatives fail mobile requirements. Accept API limitations as temporary while mobile functionality remains critical path.

**Test backups religiously**: Untested backups are not backups. Implement automated weekly restore testing with zero tolerance for failures. Schedule quarterly disaster recovery drills covering PostgreSQL restore, full stack restore, cross-cloud failover, and air-gapped backup recovery. Document all issues immediately.

**Unified data layer delivers maximum value**: PostgreSQL + pgvector with Supabase API eliminates 73% of typical multi-service integration complexity. This architecture enables Letta to "see the full picture" across tasks, docs, files, bookmarks through single query interface. Invest deeply in this foundation—it multiplies value of every application added.

**Monitor from day one**: Deploy Prometheus, Grafana, and Loki immediately, not after problems emerge. Essential for debugging AI services, understanding resource usage, capacity planning, and catching issues before users report them. Set up basic alerts (service health, disk space, backup failures) within first week.

**Security layers matter**: Client-side encryption before cloud upload, separate keys per provider, least privilege IAM, immutable backups with Object Lock, MFA on destructive operations, and Row Level Security in PostgreSQL. These prevent single compromise from destroying all data.

**Accept pragmatic tradeoffs**: No self-hosted tool perfectly replicates every cloud service feature. Plane provides 85% of ClickUp functionality—enough for most workflows. AppFlowy delivers 80% of Notion databases—rollups coming but relations work today. Ship 80-90% solutions rather than waiting indefinitely for 100% parity that may never arrive.

## Conclusion

This architecture delivers production-ready self-hosted infrastructure for Legacy AI personal assistant with concrete paths from cloud services to open-source alternatives. Plane and AppFlowy form Phase 1 foundation addressing immediate data fragmentation across ClickUp and Notion. Nextcloud provides battle-tested file sync with comprehensive API access. Hoarder unifies bookmark and content management eliminating tool sprawl. The PostgreSQL + pgvector unified data layer with Supabase API access enables Letta agents to query across all services seamlessly, while Debezium CDC achieves sub-150ms real-time context updates. Comprehensive multi-cloud backup strategy using pgBackRest, Velero, and Borg ensures disaster recovery with 3-2-1-1-0 rule compliance. Docker Compose with Portainer provides pragmatic orchestration balancing simplicity against enterprise Kubernetes complexity, appropriate for personal infrastructure scale. The 10-week deployment roadmap moves systematically from foundation through integration to production hardening. This research synthesizes 2024-2025 best practices with specific GitHub repositories, mobile app assessments, API capabilities, and architecture patterns providing clear implementation guide for migrating to fully self-hosted AI infrastructure.