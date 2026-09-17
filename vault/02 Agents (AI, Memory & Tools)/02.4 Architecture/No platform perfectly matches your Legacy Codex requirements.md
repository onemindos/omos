---
id: ui-02-20260111053925
title: "No platform perfectly matches your Legacy Codex requirements"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "No platform perfectly matches your Legacy Codex requirements"
---

# No platform perfectly matches your Legacy Codex requirements

**None of the 9 evaluated platforms fully support all critical requirements**—specifically, no platform integrates with PowerSync for offline-first sync. However, NocoDB offers the best Supabase backend compatibility, while several platforms provide excellent MCP support for AI agents. The recommended approach is a layered architecture combining NocoDB for web UI, custom PowerSync implementation for mobile, and direct Supabase API access for AI agents.

This matters because your vision of a unified Supabase database with pgvector, offline-first mobile access, and AI agent integration cannot be achieved with any single no-code platform. You'll need to either compromise on requirements or build a custom hybrid solution. Among the options, **two platforms (Cal.com and Appsmith) are immediately disqualified** as they're fundamentally the wrong tool type—Cal.com is a scheduling platform, not a project management system, and Appsmith is a UI builder requiring you to build everything from scratch.

The core challenge is architectural: PowerSync requires direct PostgreSQL replication access, while these platforms abstract database access through their own layers. Additionally, most platforms either require their own database schemas (Plane.so, Teable, NocoBase) or use incompatible technologies (Anytype uses IPFS, Aitable.ai uses MySQL). Only NocoDB was specifically designed to connect to external databases as a frontend layer, making it the sole viable option for true Supabase backend integration.

## Platform comparison matrix

| Platform | Supabase Backend | PowerSync | MCP Integration | Feature Completeness | Mobile App | Production Status | Overall Fit |
|----------|-----------------|-----------|-----------------|---------------------|------------|------------------|-------------|
| **NocoDB** | ✅ Excellent | ❌ No | ✅ Yes (2025) | Moderate | ❌ No | ✅ Stable | 🟡 Best option |
| **Plane.so** | ❌ No | ❌ No | ✅ Excellent | Excellent | ✅ Yes | ✅ Stable | 🟡 Feature-rich |
| **AppFlowy** | ⚠️ Fixed schema | ❌ No | ⚠️ Zapier only | Moderate | ✅ Yes | ⚠️ Mostly stable | 🟡 Offline-first |
| **Anytype** | ❌ IPFS only | ❌ No | ✅ Excellent | Limited PM | ✅ Yes | ⚠️ Beta | 🔴 Wrong arch |
| **Teable** | ❌ Own PG | ❌ No | ✅ Community | Good | ❌ No | ⚠️ Alpha | 🔴 Immature |
| **NocoBase** | ⚠️ Paid plugin | ❌ No | ❌ Experimental | Excellent | ⚠️ Limited | ✅ Stable | 🔴 Poor fit |
| **Aitable.ai** | ❌ MySQL only | ❌ No | ✅ Official | Good | ⚠️ Unclear | ✅ Stable | 🔴 Wrong DB |
| **Appsmith** | ✅ Yes | ❌ No | ❌ No | N/A - UI builder | ❌ No | ✅ Stable | ❌ Disqualified |
| **Cal.com** | ✅ Yes | ❌ No | ✅ Excellent | N/A - Scheduling | ❌ No | ✅ Stable | ❌ Disqualified |

## The Supabase backend challenge

**NocoDB is the only platform that truly supports using Supabase as an external backend.** This is NocoDB's core design—it acts as a spreadsheet-like UI layer on top of existing databases. You connect it to your Supabase PostgreSQL instance, specify your schema, and NocoDB generates views, APIs, and webhooks for your existing tables. The platform maintains its own metadata database separately but your actual data stays in Supabase where pgvector can access it.

In contrast, most platforms require ownership of the database schema. Plane.so uses Django ORM with specific migrations that expect full control over table structure. Teable and NocoBase both manage PostgreSQL internally through their own ORM layers, preventing external database use. AppFlowy technically supports Supabase but only with **AppFlowy's predefined schema**—tables like `af_collab_update`, `af_user`, and `af_collab_snapshot` that don't match your custom structure. This means you'd need to restructure all data to fit AppFlowy's CRDT-based collaborative editing format, losing the flexibility to design your own schema for pgvector embeddings.

The worst incompatibilities come from platforms using entirely different database technologies. Anytype stores data in a proprietary IPFS network with Protobuf serialization—innovative for privacy and decentralization but completely incompatible with PostgreSQL. Aitable.ai requires MySQL (not PostgreSQL) plus Redis, with no support for external databases at all. These architectural decisions make Supabase integration fundamentally impossible without rebuilding the platforms.

NocoBase deserves special mention as a near-miss. It does support external PostgreSQL connections, but this capability is locked behind the **"External Data Source - PostgreSQL" commercial plugin**. The community edition under AGPL-3.0 license only allows NocoBase to manage its own database. Even with the paid plugin, NocoBase would create an abstraction layer over your Supabase data rather than working as a true external frontend like NocoDB.

## The PowerSync impossibility

**No platform supports PowerSync integration, and this limitation is architectural rather than a missing feature.** PowerSync works by reading PostgreSQL's Write-Ahead Log (WAL) using logical replication, transforming changes into a format suitable for SQLite, and syncing to client devices. This requires direct database-level access and specific replication slot configurations that none of these no-code platforms expose.

The fundamental problem is that these platforms sit between your application and the database, handling all queries through their own API layers. NocoDB queries your Supabase database on behalf of users but doesn't provide the low-level replication access PowerSync needs. Plane.so and others manage their own databases entirely, making it impossible for PowerSync to even see your data. Even if you configured PowerSync to read from Plane.so's PostgreSQL instance, it would sync Plane's internal schema (django_migrations, plane_workspace, etc.) rather than a clean data model suitable for mobile apps.

AppFlowy and Anytype both offer their own offline-first architectures as alternatives. **AppFlowy uses Yjs CRDT (Conflict-free Replicated Data Types)** for collaborative editing with offline support, syncing between local SQLite and cloud PostgreSQL automatically. Anytype goes further with a peer-to-peer sync protocol called "any-sync" built on a private IPFS network with end-to-end encryption. These are sophisticated offline solutions but completely different from PowerSync's approach and incompatible with existing PowerSync-based mobile apps.

The only viable path forward is building a **custom mobile app using PowerSync SDK** directly against your Supabase database, independent of any no-code platform. You'd use React Native or Flutter with PowerSync's client libraries, configure sync rules in Supabase, and build your own mobile UI. The no-code platform (like NocoDB) would serve as the web administration interface only, with mobile users accessing data through a completely separate app that happens to share the same Supabase backend.

## MCP support across platforms varies widely

**Three platforms offer excellent official MCP servers**: Plane.so, Anytype, and Cal.com. Plane.so's `@makeplane/plane-mcp-server` package provides over 30 specialized tools for project management, including creating issues, managing cycles, updating projects, and adding comments. Anytype's `@anyproto/anytype-mcp` enables AI agents to perform workspace searches, manage spaces and members, handle objects and lists, and work with properties and tags through natural language. Cal.com's `@calcom/cal-mcp` gives AI agents complete access to the scheduling API, though this is only relevant if you need scheduling functionality within your project management system.

**NocoDB recently added MCP support in January 2025** with version 2025.09.0, marking a significant enhancement for AI integration. The MCP server runs as part of NocoDB's backend, exposing endpoints for conversational queries like "show me all customers who haven't paid in 30 days." It supports instant insights, bulk operations, smart data cleanup (finding duplicates and inconsistencies), and full CRUD operations through AI prompts. The implementation works with Claude Desktop, Cursor, and Windsurf, requiring only an API token for authentication. Community members have also built alternative MCP servers for NocoDB, providing multiple integration paths.

Teable and Aitable.ai have **community-built MCP servers** rather than official implementations. For Teable, developers created `teable-mcp-sse` using OpenAPI MCP Server with Server-Sent Events transport, and `teable-mcp-server` for direct TypeScript-based interaction. Aitable.ai's official MCP server at `github.com/apitable/aitable-mcp-server` provides read and write access to datasheets with tools for listing spaces, searching nodes, reading records, fetching field schemas, creating records, and uploading attachments.

**AppFlowy's MCP integration is limited to Zapier**, which adds complexity and dependency on a third-party service. You'd configure Zapier's MCP server as a bridge, giving AI agents access to AppFlowy through Zapier's API layer. This works but introduces latency and potential reliability issues compared to native MCP servers. NocoBase has an experimental repository at `github.com/nocobase/mcp-server-nocobase` but it's minimal with no releases or documentation—essentially unusable in production. Appsmith isn't designed to be an MCP server at all since it's a UI builder rather than a data platform.

## Top three platforms ranked by suitability

### NocoDB emerges as the best partial solution

NocoDB wins primarily because **it's the only platform designed to work with your existing Supabase database as an external data source**. Installation takes 1-2 hours using Docker, then you add your Supabase connection string, specify your schema (like "public" or a custom schema name), and NocoDB immediately provides a spreadsheet-like interface for all your tables. You can create multiple views (Grid, Kanban, Gallery, Calendar) on the same data, configure webhooks that trigger on record changes, and use the REST API for external integrations. The January 2025 MCP update means AI agents can now query and manipulate data through natural language, satisfying your requirement for hybrid human-UI/AI-API access.

However, NocoDB has significant gaps for project management. It's a **database frontend, not a task management system**—there are no built-in concepts of projects, recurring tasks, reminders, or sprints. You'd need to build this structure yourself by creating tables for tasks with date fields, status columns, and relationship fields. Recurring tasks would require external automation (using webhooks to trigger n8n or Zapier workflows). The platform also lacks a native mobile app, offering only a responsive web interface that requires constant internet connectivity. For true offline-first mobile access, you'd need to build separate mobile apps using PowerSync SDK while keeping NocoDB as the web admin interface.

The **implementation path** requires configuring NocoDB's external data source with your Supabase connection details, including SSL settings for Supabase's pooler. You'd design your task/document schema directly in Supabase (ensuring it supports your pgvector embedding needs), then have NocoDB connect and generate the UI layer. For AI agent access, you'd either use NocoDB's MCP server or build a custom MCP wrapper around Supabase's PostgREST API for more direct database access. Total setup time ranges from 45-90 hours depending on whether you need the custom PowerSync mobile layer, with most complexity in building the separate mobile app rather than configuring NocoDB itself.

### Plane.so offers the best feature set but requires dual databases

**Plane.so is a purpose-built project management platform** with issues, sprints (called "cycles"), modules for feature grouping, roadmaps, and real-time collaboration—exactly what you'd expect from a modern PM tool. The `@makeplane/plane-mcp-server` is exceptionally well-designed with over 30 specialized tools, and native iOS/Android mobile apps provide good user experience. The platform is production-ready with 39,500+ GitHub stars, active development, and enterprise customers. If your primary requirement were project management features rather than database flexibility, Plane.so would be the clear winner.

The critical dealbreaker is **Plane.so cannot use Supabase as its backend**. It requires its own PostgreSQL database with Django-specific schema including tables for workspaces, projects, issues, cycles, modules, and dozens of other models. The architecture expects full control over migrations, using Django ORM to handle all database operations. You cannot point Plane.so at an existing Supabase database and expect it to work—it needs to create and manage its own schema entirely.

The only workaround is building a **custom bidirectional sync service** that maintains data consistency between Supabase (your source of truth with pgvector) and Plane.so's PostgreSQL instance. This sync service would consume webhooks from both systems, translate data models between schemas, handle conflict resolution when the same record changes in both systems, and manage eventual consistency. Implementation would take 80-120 hours just for the sync layer, introduces significant latency, creates a single point of failure, and adds maintenance burden. This architectural complexity defeats your goal of having a unified database, essentially giving you two sources of truth that must be kept synchronized. Not recommended unless you're willing to abandon the Supabase-as-single-backend requirement.

### AppFlowy provides offline-first but with schema rigidity

**AppFlowy excels at offline-first architecture** with native mobile apps for iOS and Android that work fully offline, automatically syncing when connectivity returns. The platform uses Yjs CRDT for conflict-free collaborative editing, meaning multiple users can edit the same document simultaneously without conflicts. For a team needing offline access and Notion-like document editing, AppFlowy delivers well. The 66,000+ GitHub stars indicate strong community support, and the AGPL-3.0 license allows self-hosting with full control.

However, AppFlowy's Supabase integration is **rigidly structured around AppFlowy's predefined schema**. The official AppFlowy-Supabase repository includes a setup tool that creates specific tables: `af_collab_update` for collaboration updates, `af_user` for user profiles, `af_collab_statistics` for metrics, and `af_collab_snapshot` for recovery snapshots. These tables store data in CRDT format optimized for real-time collaborative editing, not as clean relational data suitable for direct SQL queries or pgvector embeddings. You cannot use your custom schema—you must restructure all data to fit AppFlowy's model.

This means **losing pgvector flexibility** since your data is encoded in Yjs format within AppFlowy's collaboration tables rather than stored as clean columns you can directly embed and search. AI agents accessing data through AppFlowy would work with the collaborative document structure, not your custom task/document schema. MCP integration exists only through Zapier, adding another dependency layer. The platform also lacks recurring tasks (on roadmap for Q4 2025) and native reminders, limiting project management capabilities. Setup takes 26-49 hours including data migration, but you'd be accepting significant architectural constraints that conflict with your pgvector-enabled AI vision.

## Immediately disqualified platforms

**Cal.com should not have been included in this evaluation**—it's a scheduling and calendar booking platform like Calendly, not a project management system. While it has excellent PostgreSQL/Supabase support and an official MCP server (`@calcom/cal-mcp`), its features are entirely focused on appointment scheduling, availability management, team booking pages, and calendar integrations. There's no task tracking, no document management, no file storage, no project organization—none of the capabilities needed for Legacy Codex. Cal.com might be useful as a supplementary scheduling tool within a larger project management ecosystem, but it cannot serve as the core platform.

**Appsmith is similarly misaligned** as it's a low-code UI builder for creating custom admin panels and dashboards, not a complete application platform. You'd use Appsmith to build interfaces on top of existing databases (it does connect to Supabase well), but you're starting with a blank canvas. There are no pre-built task management features, no knowledge base, no mobile apps, no offline support—you'd need to construct everything from scratch using widgets and custom JavaScript. While powerful for building bespoke internal tools, it requires hundreds of hours of development work to create what other platforms provide out-of-the-box. It's the wrong tool type for your use case.

## Platforms with critical architectural gaps

**Anytype uses a proprietary IPFS-based architecture** that's fundamentally incompatible with PostgreSQL. Data is stored in a private IPFS network using Protobuf serialization with peer-to-peer sync via the "any-sync" protocol. This innovative approach provides excellent privacy (zero-knowledge encryption) and local-first operation, with official MCP support for AI agents. However, it's designed as a personal knowledge management tool (competing with Obsidian and Roam Research) rather than project management, and it's still in beta/alpha with acknowledged quality issues. Most critically, there's no path to integrate with Supabase—the architectures are completely different at the foundational level.

**Teable manages its own PostgreSQL instance** via ORM and cannot connect to external databases like Supabase. Despite being positioned as an "Airtable alternative" and using PostgreSQL internally, it abstracts database access through its application layer. The platform offers 7 view types, API-first design, and community MCP servers, but it's still releasing alpha builds (v1.8.0-alpha) with mixed production readiness signals. Documentation from 2024 warned it was in "technical preview" and unsuitable for production. Even if stability improves, the inability to use external Supabase as a backend disqualifies it from your architecture.

**NocoBase requires a commercial plugin** called "External Data Source - PostgreSQL" to connect to external databases—the free AGPL-3.0 community edition only works with its own managed database. Even with the paid plugin, NocoBase doesn't have working MCP support (only an experimental repository with no releases). The platform does offer excellent features including Kanban, Gantt charts, powerful BPM workflows, and comprehensive automation with 22+ workflow plugins, making it production-ready for general use. However, the commercial plugin requirement and missing MCP integration make it a poor fit compared to NocoDB's free external database support and official MCP server.

**Aitable.ai (also called APITable)** uses MySQL plus Redis as its backend, making it incompatible with your PostgreSQL/Supabase requirement at the database technology level. PowerSync specifically requires PostgreSQL's logical replication features that don't exist in MySQL. While the platform offers an official MCP server, 7 view types, API-first design, and AGPL-3.0 licensing, the database mismatch is insurmountable. Migration would mean moving all data from Supabase PostgreSQL to Aitable.ai's MySQL instance, losing your existing infrastructure and likely requiring schema changes.

## Recommended architecture: layered approach

Given that no platform meets all requirements, the most practical solution is **accepting a multi-component architecture** rather than seeking a single unified platform. Use NocoDB as your web administration interface, connecting directly to Supabase as an external data source. This gives desktop users a rich spreadsheet-like UI with Kanban boards, Calendar views, and Grid views for managing tasks and documents. Configure NocoDB's MCP server or build a lightweight custom MCP wrapper to enable AI agent access to the same Supabase data.

For mobile offline-first access, **build custom native apps using PowerSync SDK** with React Native or Flutter. These apps would use PowerSync's client libraries to sync between Supabase PostgreSQL and local SQLite databases on iOS/Android devices. The sync happens at the database level (via PostgreSQL logical replication), completely independent of NocoDB. You'd design a mobile-optimized UI tailored for task management on the go, with full offline capability and conflict-free sync when connectivity returns. This separates concerns: web interface via NocoDB, mobile via custom PowerSync apps, both accessing the same Supabase backend.

**Estimated development effort** breaks down as follows: NocoDB setup and view configuration takes 2-4 hours, custom MCP server (if needed instead of using NocoDB's built-in server) requires 10-20 hours, PowerSync configuration on Supabase takes 4-8 hours, and building the mobile apps requires 150-300 hours for a production-ready iOS/Android solution. Total effort ranges from 166-332 hours, with the vast majority spent on mobile app development. This seems like significant work, but remember you're getting exactly the architecture you specified—unified Supabase backend with pgvector, offline-first mobile sync, AI agent access via MCP, and a rich web UI.

The alternative is **building everything custom from scratch** on Supabase without any no-code platform. You'd create a Next.js web application, implement your own task management UI, build mobile apps with PowerSync, and create an MCP server wrapping Supabase's API. This provides maximum flexibility and no platform limitations, but requires 600-1000+ hours of development. The NocoDB-based approach saves 400-600+ hours by using NocoDB's pre-built web UI components while still achieving your core architectural requirements.

## Migration path from existing tools

Your mention of migrating data from ClickUp, Notion, Raindrop, and Dropbox requires careful planning regardless of platform choice. **ClickUp exports** provide JSON or CSV files with tasks, subtasks, custom fields, comments, and attachments. You'd need to transform this into your Supabase schema, mapping ClickUp's structure to your task model. NocoDB can import CSV files directly, but you'll likely need a custom migration script to handle relationships, attachments, and complex field types properly.

**Notion exports** come as Markdown files with a folder structure representing the page hierarchy, plus CSV exports for databases. Notion's block-based content needs transformation into whatever format you choose for document storage—possibly rich text fields in Supabase tables or separate document storage. If using AppFlowy, its Notion import feature could help, but remember you'd be locked into AppFlowy's schema. For NocoDB, you'd import document content into text fields or use a separate document management system.

**Raindrop.io bookmarks** export as HTML or CSV with URLs, titles, descriptions, tags, and collections. This is relatively straightforward to import into Supabase tables, either through NocoDB's CSV import or via SQL scripts. The challenge is maintaining the tagging and collection structure in your new system. **Dropbox file migration** is the most storage-intensive, requiring a decision about whether to store files in Supabase Storage, AWS S3, or another service. NocoDB supports S3-compatible storage, so you could migrate Dropbox files to S3 and reference them through NocoDB's attachment fields.

Expect the full migration process to take **40-80 hours** depending on data volume and complexity. Building custom ETL (Extract, Transform, Load) scripts is likely necessary to preserve relationships and metadata. Test migrations thoroughly with a subset of data before running the full import to catch schema mismatches and edge cases.

## Timeline and effort estimates

For a **quick proof-of-concept** using NocoDB with Supabase (web-only, no mobile), you could have a working system in **1-2 weeks** with:
- Supabase schema design: 8-16 hours
- NocoDB setup and configuration: 4-8 hours  
- View creation and customization: 8-16 hours
- MCP server configuration: 4-8 hours
- Basic data migration: 16-32 hours
- Testing and refinement: 8-16 hours

This gives you a functional web-based system where humans can manage tasks/documents through NocoDB's UI and AI agents can access data via MCP, all backed by your Supabase database with pgvector support.

For the **full solution with offline-first mobile**, expect **3-6 months** including:
- Everything from proof-of-concept: 48-96 hours
- PowerSync configuration and testing: 8-16 hours
- Mobile app development (iOS/Android): 150-300 hours
- Custom MCP server (if needed): 10-20 hours
- Full data migration with validation: 40-80 hours
- End-to-end testing: 40-80 hours
- Bug fixes and polish: 40-80 hours

Total effort: **336-672 hours** (8-17 weeks at 40 hours/week).

## Key technical decisions to make immediately

**First, confirm whether PowerSync offline-first mobile is truly required** or if responsive web access is sufficient. If most users work from desktop/laptop with reliable internet, you could ship much faster with just NocoDB web UI and delay mobile app development. This reduces timeline from months to weeks and effort from 600+ hours to under 100 hours.

**Second, decide on pgvector usage patterns**. Will AI agents need to search embeddings frequently, or is this an occasional operation? If frequent, ensure your custom MCP server (or NocoDB's MCP) efficiently handles pgvector similarity searches. You might need custom endpoints specifically for semantic search that aren't provided by generic database MCP servers.

**Third, define your task management requirements precisely**. Do you actually need recurring tasks, dependencies, Gantt charts, and sprint planning? Or would simpler task tracking with due dates and status suffice? The complexity of your PM needs determines whether you should build more custom functionality on top of NocoDB or consider a hybrid approach where Plane.so handles project management and syncs to Supabase.

**Fourth, evaluate whether you can accept NocoDB's limitations** for web UI or need more sophisticated project management features. If team collaboration, sprint planning, roadmaps, and issue tracking are critical, the NocoDB-only approach may feel too basic. In that case, using Plane.so as the PM interface with a sync service to Supabase becomes more justified despite the complexity.

## Final recommendation

**Start with NocoDB connected to Supabase** as your web interface and direct Supabase API access (via PostgREST or custom MCP server) for AI agents. This gets you operational quickly with minimal investment. Use NocoDB for the administration tasks it handles well—data viewing, editing, basic automation via webhooks. Build lightweight custom tools where NocoDB falls short, such as a simple web page for recurring task management or a custom dashboard for pgvector similarity searches.

Once the web system proves valuable, **then invest in mobile apps** using PowerSync if offline-first access becomes critical. By delaying mobile development until you've validated the core system, you avoid wasting months building apps that might need significant changes based on user feedback. The PowerSync mobile work is substantial enough that you want to be certain it's necessary before committing resources.

This staged approach requires **48-96 hours for initial web system**, then **300-500 additional hours for mobile** if you proceed with that phase. It provides the fastest path to validating your Legacy Codex concept while keeping the door open for the full offline-first architecture you envision. Most importantly, your data lives entirely in Supabase throughout this evolution, so you're never locked into a platform's proprietary database—you maintain full control and flexibility to change direction as requirements evolve.