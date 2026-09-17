---
id: ui-06-20260111053926
title: "Building a Unified Obsidian Vault: Complete Architecture Guide for Large-Scale Knowledge Management"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/06]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Building a Unified Obsidian Vault: Complete Architecture Guide for Large-Scale K"
---

# Building a Unified Obsidian Vault: Complete Architecture Guide for Large-Scale Knowledge Management

Your goal of creating a comprehensive Obsidian vault that handles thousands of PDFs, integrates with multiple tools, syncs across devices, and serves as an AI-powered knowledge base is entirely achievable with the right architecture. Based on extensive research of 2024-2025 solutions, here's everything you need to know.

## Primary questions answered first

**Obsidian handles large attachment collections remarkably well.** A tested vault with 60,000 PDFs performed reasonably (1-2 second load times per PDF), while another user manages 20,000+ files successfully. The critical finding: Obsidian only indexes filenames of binary files, not their contents, so attachments don't significantly impact vault performance until you actually open them. Individual PDFs larger than 650 pages cause noticeable UI lag, and files over 10MB may slow annotation features, but typical collections perform fine.

**For cloud storage with mixed file types across multiple locations, the answer depends on your platform needs.** Obsidian Sync ($8-16/month) provides the most seamless experience—rock-solid reliability, end-to-end encryption, works flawlessly on all platforms including iOS. For budget-conscious users, Dropbox offers excellent cross-platform support but costs more for storage alone. iCloud works well for Apple-only users but has Windows corruption issues and unreliable sync. The dark horse: Syncthing provides free peer-to-peer sync with local WiFi speed, though it requires more technical setup.

**Files should generally be embedded in the vault for true unified management**, but with strategic exceptions. Keep your core markdown notes and frequently referenced PDFs (under 5MB) inside the vault. Store large media files, video archives, and infrequently accessed documents externally with file:// links or using cloud storage plugins. This hybrid approach maintains the "all in one" experience while keeping sync performant and mobile-friendly.

## The optimal workflow architecture

### Task management with Todoist integration

The Todoist ecosystem for Obsidian has matured significantly. **Ultimate Todoist Sync provides true bi-directional synchronization**—create tasks in Obsidian using checkbox syntax with `#todoist` tags, and they automatically sync to Todoist within 5 minutes. Set due dates with calendar emojis (📅 2025-01-15), priorities with arrows (⏫ urgent, 🔼 high), and link pages to specific Todoist projects. For visualization, install Custom Frames to embed the Todoist web app in your sidebar, giving you the best of both worlds.

Alternatively, the native Obsidian Tasks plugin has evolved into a powerhouse for those ready to move away from external services. It offers custom statuses beyond checked/unchecked, sophisticated query-based dashboards filtering by tags and properties, and recurring task handling with full history logs. Combine it with Calendar and Periodic Notes for time-based planning, or add Kanban for visual project boards. The real advantage: everything stays local with powerful cross-vault querying capabilities.

**Recommended configuration**: Start with Ultimate Todoist Sync if you need mobile notifications and quick capture via Todoist's excellent mobile apps. Gradually shift project planning and detailed task context into Obsidian notes embedded with Todoist queries. Over time, migrate to Tasks plugin if you find you're managing most workflows in Obsidian anyway.

### Migrating your existing content ecosystem

**Apple Notes migration is now seamless** thanks to Obsidian's official Importer plugin, which directly decodes the proprietary SQLite database (macOS only). It preserves everything: formatting, tables, images, scanned documents, PDFs, drawings converted to images, internal note links from iOS 17+, and folder structure. The process takes 10-15 minutes: install Importer plugin, grant permission to access the notes database, and let it convert your library automatically. For 8,000+ notes, expect 20-30 minutes but zero data loss.

**Notion migration requires a two-step dance.** Export your entire workspace as HTML (not Markdown—HTML preserves more content types), then import the ZIP file using the Importer plugin with "Save parent pages in subfolders" enabled. The good news: page content, formatting, links, images, and simple tables transfer beautifully. The challenge: Notion databases don't have direct equivalents. Convert them to folder structures with YAML frontmatter properties, then use Dataview plugin for database-like queries or DB Folder plugin for visual table views. Plan 2-3 hours for migration plus another 2-4 hours adapting database workflows to Obsidian's linking philosophy.

**For your Raindrop collection of 8,000+ bookmarks**, Make It Rain plugin offers the most comprehensive solution. It fetches bookmarks directly via the Raindrop API, preserving titles, descriptions, URLs, tags, cover images, creation dates, and collection hierarchy. The plugin handles pagination automatically—expect 10-15 minutes for a complete import at the API rate limit of 120 requests per minute. Create folder structures matching your Raindrop collections, or flatten everything and rely on tags with Dataview queries. For ongoing sync, the Raindrop Highlights plugin by kaiiiz maintains bidirectional updates, especially valuable if you use Raindrop's annotation features.

### Integrating file attachments from existing storage

The hybrid approach proves most practical: **keep essential files in the vault, link to the rest externally.** Configure Obsidian to store new attachments in a centralized `/Attachments` folder using absolute paths (Settings → Files & Links). This prevents broken links when reorganizing and makes backup simpler. For your existing Dropbox library of hundreds of PDFs, you have three solid options.

**Option one**: Move frequently referenced files into the vault's attachment folder. Use Remotely Save plugin to sync the vault to Dropbox, giving you mobile access through Obsidian while keeping files versioned. This works well for collections under 10-20GB.

**Option two**: Keep large files in Dropbox and use External File Embed plugin for relative path links like `![](vault://../Dropbox/Research/paper.pdf)`. On mobile, tap the link to open in the Dropbox app. This keeps your vault lean (critical for Obsidian Sync's 10GB limit) while maintaining access to your full library.

**Option three**: Use Cloud Storage plugin to automatically upload attachments to S3-compatible storage (cheaper than traditional cloud services). The plugin updates links in notes automatically and supports resumable uploads. Best for users with massive media libraries who want to minimize local storage.

**Mobile considerations matter here.** iOS apps can't directly access Dropbox's synced folders due to sandboxing. If you need robust mobile editing with external files, Obsidian Sync for the vault plus Dropbox app for viewing linked documents provides the smoothest experience, despite requiring two apps.

### Version control for markdown plus binary files

**Git integration works beautifully for markdown but requires strategy for binaries.** The Obsidian Git plugin provides automatic commit-and-sync on configurable schedules, with a Source Control View for staging individual files and browsing history. Desktop functionality is complete, but mobile uses isomorphic-git (JavaScript implementation) with significant limitations—no SSH authentication, memory constraints for large repos, and poor merge conflict handling.

**For mixed content vaults, Git LFS (Large File Storage) becomes essential.** Configure it to track `*.pdf`, `*.png`, `*.jpg`, `*.mp4`, and other binary types. This stores actual files on the LFS server while keeping pointers in your Git repository, preventing bloat. GitHub offers 2GB free LFS storage; GitLab is more generous; or self-host with Gitea for unlimited capacity. **Critical limitation**: Git LFS doesn't work on mobile with isomorphic-git, so mobile users should sync attachments separately via Obsidian Sync or cloud storage.

Your `.gitignore` strategy makes or breaks multi-device workflows. **Always exclude** `.obsidian/workspace.json` (causes constant conflicts with window layouts) and `.obsidian/workspace-mobile.json`. Include core settings like `app.json`, `hotkeys.json`, and plugin configurations. For plugin folders, exclude code files (`main.js`, `styles.css`, `manifest.json`) but keep settings (`data.json`). This lets you sync preferences without conflicts from plugin updates.

**VS Code integration enhances the development workflow.** Install the "Open Vault in VSCode" plugin for quick switching, or simply open your vault folder as a VS Code workspace. Use it for advanced search/replace with regex, mass file operations, Git operations with visual merge tools, and formatting with Prettier. The markdown files remain compatible—Obsidian-specific syntax like `[[wikilinks]]` displays as text in VS Code but doesn't break. For AI assistance, Claude Desktop can access your vault via Model Context Protocol (MCP) servers, providing semantic search and content generation with proper permission-based access.

## Future-proofing with AI and semantic capabilities

### Semantic search implementations

**Smart Connections has emerged as the community standard for semantic search**, with zero-setup local AI embeddings that work offline by default. Install it from Community Plugins, let it index your vault using local models (no API keys required), and immediately access semantic search through Smart Chat for conversational queries and Smart View for real-time connection suggestions as you write. The plugin supports 1000+ models—OpenAI for cloud power, Ollama/LM Studio for privacy, or the built-in local model for simplicity.

The technical magic: AI embeddings convert your notes into 768-1536 dimension vectors representing semantic meaning. When you search for "productivity," it finds notes about "efficiency," "time management," and "workflow optimization" even without those exact keywords. Cosine similarity calculations identify related content, with results stored in lightweight JSON files in your vault.

**For advanced users wanting custom RAG implementations**, several paths exist. LlamaIndex provides knowledge graph RAG with structured relationships between concepts—excellent for academic research. LangChain plus ChromaDB offers the most customization, letting you build sophisticated document chunking strategies and custom embedding pipelines. The new Model Context Protocol (MCP) approach creates personal note assistant servers that index your vault in vector databases like Milvus with real-time synchronization to Claude Desktop.

The simplest RAG setup: Smart Connections already provides RAG-like Q&A through Smart Chat, using your vault as context without external setup. For privacy-focused users, combine it with Copilot plugin and Ollama for completely local processing—no cloud dependencies, works offline, handles sensitive information securely.

### Letta AI for stateful interaction

**Letta AI offers something unique: genuine long-term memory rather than just conversation history.** The official Obsidian plugin (actively maintained in 2024-2025) continuously syncs markdown files to Letta's persistent memory system, maintaining directory structure using `__` separators. The agent remembers context across sessions, developing understanding of your vault over time.

The interface provides interactive memory management—create, edit, and delete memory blocks directly in Obsidian, search and attach available blocks to agents, and resolve conflicts from server-side changes. Unlike typical chatbots that forget between sessions, Letta maintains genuine long-term memory, making it ideal for research assistants that remember previous discussions, writing partners maintaining project context, or knowledge curators developing understanding of your organizational patterns.

Setup requires a Letta Cloud account or self-hosted instance. Get your API key from app.letta.com, configure the plugin with your project slug and agent name, and the automatic sync begins. While more complex than simpler plugins, the stateful memory capabilities justify the setup for users managing complex, evolving knowledge bases.

### Canvas for visual organization

**Canvas provides Obsidian's native infinite whiteboard**, released December 2022 and continuously improved. Create unlimited spatial layouts mixing notes (existing or new), text cards without file creation, images, PDFs, videos, audio, and fully interactive embedded web pages. The JSON Canvas format (jsoncanvas.org open standard) makes it future-proof and parseable by other tools.

Real-world use cases span project planning with visual phase layouts, research organization mapping relationships between papers, story plotting for writers (similar to Scrivener's corkboard), and meeting prep with spatially arranged reference materials. Combine Canvas with Dataview for high-level visual organization plus dynamic filtered lists. Layout AI-discovered connections from Smart Connections on canvas to visualize semantic relationships. Embed Canvas views within markdown notes or nest Canvas views within other Canvas views for fractal organization.

**Limitations to know**: No built-in collaboration (though screen.garden plugin adds multiplayer), limited drawing tools compared to dedicated whiteboard apps, and performance can lag with very large canvases (hundreds of cards). For enhanced capabilities, MoreDraw plugin adds AI-assisted diagram creation, chart makers, QR codes, and freehand drawing tools with export to PDF, PNG, and PowerPoint.

## Technical architecture for your unified system

### Recommended sync solution for large vaults with mixed content

**For your needs—multi-device access across multiple physical locations with large mixed-file vaults—here's the optimal architecture:**

**Desktop tier**: Use Obsidian Git with Git LFS for version-controlled markdown and small attachments. Configure `.gitignore` to exclude workspace files and cache. Commit markdown changes frequently (15-30 minute auto-commit intervals) while staging large binary additions manually. Keep vault core under 10GB for mobile compatibility.

**Mobile tier**: Obsidian Sync ($10/month for 10GB or $16/month for 100GB) provides the smoothest mobile experience. Configure selective sync to exclude large attachment folders if needed. The end-to-end encryption, conflict resolution, and background sync (even when Obsidian isn't open) justify the cost for power users.

**Large file tier**: Store infrequently accessed PDFs and all videos in Dropbox or iCloud with file:// links in your notes. Use Cloud Storage plugin for automatic upload of new attachments over 5MB. This three-tier approach keeps your vault lean, sync fast, and access universal.

**Alternative budget approach**: Syncthing for free peer-to-peer sync works remarkably well if you're technically inclined. Install on all devices, configure folder sharing, and enjoy instant sync over local WiFi without cloud costs. Combine with Git for version history. **Limitation**: requires devices online simultaneously for sync, and iOS support through Möbius Sync app remains experimental.

### Vault structure for 10GB+ mixed content

**Start with a minimalist foundation that scales.** Place most notes in the vault root, use tags liberally for multi-dimensional organization, and create folders only for distinct contexts:

```
Vault Root/
├── [Active notes live here directly]
├── Projects/          # Active work
├── Resources/         # Reference materials
├── Attachments/       # All PDFs, images, audio
├── Daily/             # Daily notes (YYYY-MM-DD.md)
├── Templates/         # Note templates
└── Archives/          # Completed projects
```

**Critical configuration**: Settings → Files & Links → Default location for new attachments: "In folder specified below" → `Attachments/`. Set link format to "Absolute path in vault" so attachment links survive reorganization. This prevents the chaos of scattered attachments and broken links that plague large vaults.

**Performance thresholds from real-world testing**: Under 10GB performs smoothly, 10-50GB requires monitoring plugin usage and selective sync, 50GB+ benefits from vault splitting or external linking for largest files. One user successfully manages 936MB including multimedia with excellent performance after initial indexing. The 60,000 PDF test case remained usable, demonstrating Obsidian's scalability.

**For your 8,000+ Raindrop bookmarks**, create dedicated folder like `Bookmarks/` with subfolders mirroring Raindrop collections. Tag each with `#imported` and original Raindrop tags. Use Dataview queries to create topic-based MOCs (Maps of Content) rather than trying to maintain perfect folder hierarchies. This adapts better as your knowledge evolves.

### Essential plugin stack

**Core foundation (install immediately)**:
- **Dataview**: Database queries across vault (3.41M downloads)
- **Templater**: Advanced templates with JavaScript (3.28M downloads)
- **Tasks**: Comprehensive task management (2.89M downloads)
- **Smart Connections**: Semantic search and AI discovery (local-first)

**Workflow enhancement**:
- **Ultimate Todoist Sync**: If continuing with Todoist
- **Calendar + Periodic Notes**: Time-based organization
- **Copilot**: AI text generation and vault Q&A with local model support
- **Obsidian Git**: Version control

**Vault maintenance** (critical at scale):
- **Linter**: Bulk formatting standardization
- **Auto Note Mover**: Tag-based automatic filing
- **Attachment Management**: Organize attachments retroactively
- **Omnisearch + Text Extractor**: Full-text PDF search

**Import and migration**:
- **Official Importer**: Apple Notes, Notion, Evernote
- **Make It Rain**: Raindrop bookmark import
- **Remotely Save**: Cloud sync without Obsidian Sync subscription

**AI and semantic**:
- **Smart Connections**: Core semantic search
- **Text Generator**: Content creation
- **Letta AI**: Stateful agent (advanced users)

**Total plugin count recommendation**: 15-25 plugins maximum for optimal performance. Use Lazy Plugin Loader (new 2025) to reduce startup time if you exceed 20 plugins.

### Performance optimization strategies

**For large PDF collections specifically**: Use symbolic links to keep truly massive reference libraries outside the vault while maintaining access. Enable "exclude attachment folders from search" in settings to speed queries. Consider PDF++ plugin's external file feature for PDFs you need to reference but not annotate. Keep your "working set" of actively annotated PDFs in-vault (perhaps 100-500 files) and archive the rest with file:// links.

**Mobile optimization matters**: Configure Obsidian Sync with selective folder sync to exclude attachment folder on mobile. Use 10-20MB limit for in-vault PDFs on mobile. Link to Dropbox/iCloud for viewing larger files in their respective apps. This keeps your mobile vault under 2-3GB for responsive performance even on older devices.

**Graph view with large vaults**: Exclude the Attachments folder from graph (Settings → Graph view → Files → Add exclusion). Reduce max depth to 2-3 connections. Use local graphs on individual notes rather than the global graph. Consider Juggl plugin for more performant visual graphs with better rendering.

## Addressing your specific concerns

**Will including large folders of PDFs impact performance?** Moderately—the initial indexing takes several minutes for thousands of files, and the vault file count increases load time. But day-to-day editing remains smooth because Obsidian only indexes filenames, not PDF contents. The real performance hit comes from opening many large PDFs simultaneously or using annotation plugins on 650+ page documents. Solution: Keep your active working set manageable (under 1,000 PDFs in active use) and archive the rest with external links.

**Obsidian handles PDF viewing natively** through a rebuilt PDF viewer powered by Mozilla's PDF.js (version 1.3, 2023). You can open PDFs directly in tabs, deep-link to specific pages using `[[file.pdf#page=3]]`, link to text selections and annotations, search within PDFs, and copy text as formatted quotes. Limitations: loads entire PDF into memory (not page-by-page), causing CPU spikes when scrolling in large files. Performance becomes noticeably slower than native PDF readers for documents over 500 pages. Annotations require PDF++ plugin for permanent highlighting visible outside Obsidian.

**Git with binary files works through Git LFS**, which stores large files on separate servers while keeping pointers in your repository. Configure tracking for `*.pdf`, `*.png`, `*.jpg` file types, and Git handles the complexity. Desktop Git clients support this seamlessly. **Critical issue for your workflow**: mobile Git through Obsidian's plugin doesn't support LFS due to isomorphic-git limitations. For mobile users, the recommendation is clear—use Obsidian Sync for your vault including attachments, and separately maintain a Git repository of just markdown files on desktop for version control of your notes. This dual-sync approach provides the best of both worlds.

**Maintaining "all in one" while keeping performance** requires accepting strategic compromises. Keep everything you actively work with in the vault (markdown notes, templates, working PDF collection, frequently used images). Store archives, large media, and infrequently accessed files externally with links. Use Dataview to create unified dashboards that query both in-vault content and track external resources. The psychological benefit of centralization remains—one app, one interface, one search—even if the underlying storage is strategically distributed.

## Implementation roadmap

**Week 1: Foundation (8-10 hours)**
- Set up vault structure with Attachments folder and absolute path configuration
- Install core plugins: Dataview, Templater, Smart Connections, Tasks, Obsidian Git
- Configure sync: Initialize Git repository with .gitignore excluding workspace files, set up Git LFS for PDFs
- Test workflow: Create sample notes, add attachments, commit to Git, verify structure

**Week 2: Migration (12-15 hours)**
- Import Apple Notes using Importer plugin (2-3 hours including verification)
- Migrate Notion content, focusing on restructuring databases into folders + properties (4-6 hours)
- Import Raindrop bookmarks with Make It Rain, organize into folder structure (3-4 hours)
- Review and cleanup: spot-check random items, fix formatting issues, consolidate tags (3-4 hours)

**Week 3: Workflow Integration (8-10 hours)**
- Configure Todoist integration with Ultimate Todoist Sync or migrate to Tasks plugin
- Set up Copilot for AI assistance with local Ollama models for privacy
- Configure Calendar and Periodic Notes for daily planning
- Create template library with Templater for common note types
- Set up Dataview queries for dashboards (active tasks, recent bookmarks, project overviews)

**Week 4: Advanced Features (6-8 hours)**
- Configure Smart Connections for semantic search across vault, let it index completely
- Set up Canvas boards for active projects requiring visual organization
- Install PDF++ for PDF annotation workflow, configure for your research needs
- Set up Remotely Save or Obsidian Sync for mobile access
- Configure VS Code integration for advanced editing and Git operations
- Create MOCs (Maps of Content) for major topic areas

**Month 2: Refinement**
- Monitor performance, optimize plugin usage based on actual needs
- Develop personal workflows and templates based on actual usage patterns
- Set up automated maintenance routines (weekly review template with task rollover)
- Consider Letta AI for advanced stateful agent interaction
- Explore custom Dataview queries for knowledge discovery
- Build RAG implementation if needed for specific AI use cases

**Ongoing Maintenance**
- Weekly: Review and process inbox notes, update task statuses, file unorganized items
- Monthly: Run Orphan Finder to identify unlinked notes, use Linter for formatting cleanup, review and consolidate tags
- Quarterly: Evaluate plugin usage and remove unused, update all plugins, backup vault completely, review and archive completed projects

## Critical success factors

**Start simple and add complexity gradually.** Install 5-7 core plugins initially, use them consistently for 2-4 weeks, then add advanced features as you identify specific needs. Users who install 30 plugins immediately often struggle with conflicts and overwhelming options.

**Backup before major changes.** Before bulk imports, plugin installations, or structure reorganization, create complete vault backups. Git provides excellent version control but keep separate ZIP backups of the entire vault including .obsidian folder for catastrophic recovery.

**Embrace Obsidian's linking philosophy** rather than replicating previous tools exactly. Notion's databases become folders with properties and Dataview queries. Todoist's project hierarchies become tags and links. The transition period feels awkward but results in more flexible, resilient knowledge structures.

**Privacy and sync choices cascade through your architecture.** Choosing local-first AI (Smart Connections + Ollama) enables completely offline work with sensitive information. Choosing Obsidian Sync provides seamless mobile but requires trusting their encryption. Choosing Syncthing maximizes privacy but complicates mobile. Make this decision early as it influences your entire setup.

**Performance monitoring prevents frustration.** Use Task Manager or Activity Monitor to check Obsidian's RAM usage periodically. If it exceeds 4GB, investigate with plugin disabling and vault organization. Keep vault size visible—many users don't realize their vault has grown to 25GB until sync becomes painfully slow.

Your architecture of mixed content with thousands of files, multiple integrations, cross-device access, and AI capabilities is absolutely feasible with 2024-2025 tools. The Obsidian ecosystem has matured remarkably, offering solutions for every component of your workflow. The key lies in thoughtful configuration, strategic sync choices, and incremental implementation rather than attempting everything simultaneously. Start with the foundation, migrate your content systematically, and layer on advanced capabilities as your workflows stabilize. Within 4-6 weeks, you'll have a unified knowledge management system that scales to your needs while maintaining the performance and flexibility that makes Obsidian powerful.