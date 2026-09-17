# Building a sovereign knowledge base on OpenClaw

**OpenClaw's markdown-on-disk architecture makes it one of the most knowledge-base-friendly agent frameworks available today**, but it requires careful layering of extensions, cloud sync, and version control to achieve the persistent, offline-accessible, agent-writable "One Mind OS" codex you're designing. The framework's 350K+ GitHub stars and 13,700+ community skills mean nearly every piece of this stack already exists — the challenge is assembling the right combination. Here's the definitive technical map.

## OpenClaw's native workspace is the agent's brain on disk

OpenClaw agents don't live in databases. They live in **plain Markdown files** inside a workspace folder (default `~/.openclaw/workspace/`). Every session, OpenClaw reads these files and assembles the agent's identity, rules, and memory into the system prompt. The files *are* the agent — portable, git-trackable, and editable with any text editor.

The standard directory structure:

```
~/.openclaw/workspace/
├── AGENTS.md          # Operating manual — rules, boot sequence, checklists
├── SOUL.md            # Persona, tone, values, behavioral boundaries
├── USER.md            # Human profile (loaded in DM sessions only)
├── IDENTITY.md        # Agent name, emoji, avatar
├── TOOLS.md           # Environment-specific notes (SSH hosts, devices)
├── HEARTBEAT.md       # Periodic task checklist (keep under 50 lines)
├── BOOT.md            # Startup hook actions
├── BOOTSTRAP.md       # First-run onboarding (delete after setup)
├── MEMORY.md          # Long-term curated memory — durable facts, preferences
├── memory/
│   ├── 2026-04-07.md  # Daily session logs (today + yesterday auto-loaded)
│   └── archive/       # Old logs (>30 days)
├── checklists/        # High-risk operation checklists
├── docs/              # On-demand docs (NOT auto-loaded every turn)
└── skills/            # Installed skills (SKILL.md files)
```

Every bootstrap file (AGENTS.md, SOUL.md, USER.md, IDENTITY.md, TOOLS.md, MEMORY.md, HEARTBEAT.md) is injected into the system prompt at session start, so **file size directly affects token cost**. The combined character limit defaults to **150,000 characters (~50K tokens)**. Files exceeding `bootstrapMaxChars` (default 20,000 chars each) get silently truncated — OpenClaw won't warn you until after truncation happens. Content the agent needs occasionally should go in `docs/` (loaded on demand), not in bootstrap files.

The agent can **natively read and write** any file in the workspace using built-in tools: `read`, `write`, `edit`, and `exec`. Without sandboxing enabled, agents have full host filesystem access. The built-in memory tools `memory_search` (semantic hybrid search) and `memory_get` (targeted file read) provide structured retrieval. The default memory backend uses **SQLite** with `sqlite-vec` for vector similarity and `FTS5` for BM25 keyword matching, stored at `~/.openclaw/memory/{agentId}.sqlite`. Embeddings default to OpenAI `text-embedding-3-small` (1,536 dimensions), with auto-detection for Gemini, Voyage, Mistral, and Ollama providers.

The workspace is **not git-tracked by default**, but the community overwhelmingly recommends initializing it as a git repo. The official docs include instructions for this:

```bash
cd ~/.openclaw/workspace
git init && git add AGENTS.md SOUL.md TOOLS.md IDENTITY.md USER.md HEARTBEAT.md memory/
git commit -m "Add agent workspace"
```

## Memory extensions transform retrieval quality dramatically

The default SQLite memory backend scores **52/100** on the LOCOMO memory benchmark. Upgrading to LanceDB jumps that to **76/100** — the single largest accuracy improvement available. Three extensions matter most for the One Mind OS use case.

**memory-lancedb-pro** (by CortexReach) is the production-grade choice. It replaces the built-in memory with a LanceDB-backed system featuring hybrid retrieval: vector similarity via LanceDB ANN + BM25 full-text search, fused through Reciprocal Rank Fusion, then reranked by cross-encoder (Jina, SiliconFlow, Voyage, or Pinecone). A multi-stage scoring pipeline applies recency boost, Weibull time decay, length normalization, MMR diversity, and importance weighting. The system supports **multi-scope isolation** (global, agent-scoped, project-scoped, user-scoped, custom namespaces), meaning multiple agents can share or isolate memory. It auto-captures preferences, facts, and decisions from conversations at the `agent_end` hook, and auto-recalls relevant memories at `before_agent_start`. Installation is straightforward: clone into `plugins/memory-lancedb-pro`, run `npm install`, and configure in `openclaw.json` with your embedding provider. Requires OpenClaw ≥2026.3.22.

**PostClaw** replaces memory with PostgreSQL + pgvector. It adds semantic memory (RAG), episodic memory (auto-logged prompts and tool calls), a knowledge graph with typed directed edges, dynamic persona injection, and a sleep cycle for background consolidation every 6 hours. Row-level security provides per-agent isolation. This is the right choice if you're already running Postgres on your DigitalOcean VPS and want multi-agent shared memory with analytics.

**QMD** (by Tobi Lütke of Shopify) is a local search sidecar using BM25 + vector + LLM re-ranking, all on-device (~2GB of GGUF models). Its killer feature: `memorySearch.extraPaths` can index directories **outside** the workspace, and `extraCollections` lets one agent search another agent's session transcripts. For the One Mind OS, QMD enables cross-agent knowledge sharing across your Mason Android nodes and Mac Mini.

The community consensus pattern stacks three layers: a lossless context manager for the current session, LanceDB for cross-session embedding recall, and QMD via MCP for searching everything else.

## MCP integration gives agents structured access to everything

OpenClaw connects to MCP (Model Context Protocol) servers through **MCPorter** — a bundled TypeScript runtime and CLI that translates MCP tool schemas into agent-callable functions. Configuration lives in `openclaw.json`:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/dir"]
    },
    "git": {
      "command": "uvx",
      "args": ["mcp-server-git", "--repository", "/path/to/repo"]
    }
  }
}
```

Three MCP servers are critical for the One Mind OS:

The **Filesystem MCP Server** (`@modelcontextprotocol/server-filesystem`) provides `read_file`, `write_file`, `create_directory`, `list_directory`, `move_file`, `search_files`, `get_file_info`, `read_multiple_files`, and `directory_tree`. Directories passed as command-line arguments define access boundaries — operations outside are blocked. This gives agents structured, permissioned file access to your codex folder.

The **Git MCP Server** (`mcp-server-git`) exposes `git_status`, `git_diff_unstaged/staged`, `git_commit`, `git_add`, `git_reset`, `git_log` (with timestamp filtering), `git_create_branch`, `git_checkout`, `git_show`, `git_init`, and `git_list_branches`. This enables agents to commit their own changes, create branches for experimental work, and inspect version history — all through natural language commands.

**Dropbox MCP Servers** come in two flavors. Dropbox's official remote server at `https://mcp.dropbox.com/mcp` provides OAuth-authenticated read access (ListFolder, GetFileMetadata, GetFileContent up to 5MB, Search). For write access, the community server by ngs (Go-based, installable via Homebrew) supports full CRUD: list, search, download, **upload**, move, copy, delete, create folders, and manage shared links. The `65%` of active OpenClaw skills already wrap MCP servers internally — the ClawHub marketplace has **3,200+ MCP-based skills**.

MCPorter also functions as a package manager: `mcporter search browser automation`, `mcporter install playwright-mcp --target openclaw --configure`. It supports all three MCP transports (Stdio, SSE, WebSocket) and can run as a persistent daemon.

## The CLI gives you 100+ commands and full agent control

OpenClaw's CLI is comprehensive. Core commands: `openclaw onboard` (interactive setup), `openclaw doctor` (50+ health checks across 12 categories), `openclaw gateway start/stop/restart`, `openclaw agents list/add/delete`. Memory management: `openclaw memory status/index/search`. Scheduling: `openclaw cron add/edit/rm/run`. Security: `openclaw security audit --deep --fix`, `openclaw approvals set --policy ask`, `openclaw sandbox explain`.

Agents can be triggered via CLI through cron jobs (`openclaw cron add`), the heartbeat system (reads HEARTBEAT.md at configurable intervals), webhooks, or direct gateway RPC calls. The **Lobster** workflow shell (`openclaw/lobster` repo, 1.1K stars) adds composable typed pipelines — turning skills and tools into chainable automations that OpenClaw can call in one step.

The sandboxing model offers three modes: `off` (default, full host access), `non-main` (sandbox only non-main sessions), and `all` (everything in Docker containers). Workspace access in sandbox mode can be `none`, `ro` (read-only mount), or `rw` (read-write mount). For the One Mind OS, running without sandbox on the Mac Mini gives agents full codex access; the DigitalOcean VPS should use sandbox mode with `rw` workspace access.

## ClawHub has skills for every piece of the knowledge stack

The ClawHub marketplace and awesome-openclaw-skills collection contain relevant skills across every category needed:

- **File management**: `file-system-manager` (CRUD with configurable path permissions), `file-search` (content-aware semantic indexing), `file-organizer` (auto-sort into typed subdirectories), Fast.io (persistent cloud workspaces with RAG via "Intelligence Mode," multi-agent file locks, 19 tools)
- **Knowledge base / RAG**: `openclaw-rag-skill` by wmantly (ChromaDB + all-MiniLM-L6-v2, fully local, multi-source indexing), `knowledge-base` (built-in ingestion + vector store), RAGFlow integration (enterprise deep document understanding), `RAG System Builder` (FAISS + sentence-transformers, local-only)
- **Document handling**: Built-in PDF tool, `pdf-tools` by cmpdchtr (pdfplumber/PyPDF2/reportlab — extract, edit, merge, split), `book-reader` (epub/pdf/txt with progress tracking), `office-document-editor` (DOCX/PPTX/Excel with tracked changes, ~7,100 downloads), `md2pdf-converter` (Pandoc + WeasyPrint)
- **Git**: Official `github` skill (OAuth + `gh` CLI), `doro-git-essentials`, GitClaw (auto-backup agent workspace to GitHub), `arc-skill-gitops` (automated deployment + rollback)
- **Obsidian**: Official `obsidian` skill (obsidian-cli for search, create, move, delete notes), **ObsidianClaw** plugin (chat sidebar in Obsidian, agent reads/writes vault files directly, WebSocket to gateway), `obsidian-vault` skill (RAG over personal notes)
- **Cloud storage**: `openclaw-workspace-sync` by ashbrener (bidirectional sync with Dropbox, Google Drive, OneDrive, S3 via rclone — three modes: mailbox, mirror, bisync, with AES-256 encrypted backups)
- **Second brain**: `2nd-brain` (organized brain/ directory with categories), MemClaw (project-scoped workspaces with isolated memories)

## Dropbox works well as sync layer but Git needs careful handling

Dropbox is a strong candidate for syncing the One Mind codex across devices. On macOS, the desktop app creates a local `~/Dropbox` folder that any process (including OpenClaw agents) can read and write directly — **no API or MCP server needed for basic operations**. Files written there sync automatically. Smart Sync (Plus and above) lets you mark folders as "Available Offline" or "Online-Only" per folder, which is exactly what the "One Mind Media" folder needs: set it to online-only, and videos appear in Finder as placeholders consuming minimal disk space, downloading on demand when opened.

On iPhone, individual files can be made available offline on any plan, but **folder-level offline requires a paid plan** (Plus or higher). Offline files update over Wi-Fi by default. The Dropbox app's offline system is separate from the Files app — you cannot access Dropbox files offline through the iOS Files app, only through Dropbox's own "Available Offline" feature.

Dropbox's built-in version history provides **30 days on free/Plus, 180 days on Professional/Essentials, and 365 days on Business Advanced/Enterprise**. The Rewind feature can revert entire folders or your whole account to a specific point in time. One critical caveat: version history only tracks versions within the lookback window from when each version was *created* — if a file sits unchanged for 31 days and then gets overwritten, the original is unrecoverable on a Plus plan.

**Do not put a `.git` directory inside a Dropbox-synced folder.** Dropbox syncs files individually without understanding Git's internal data structures. Concurrent changes to `.git/objects` and `.git/refs` during Git operations cause repository corruption. The safe approach is `--separate-git-dir`:

```bash
git init --separate-git-dir=$HOME/.gitroots/codex.git ~/Dropbox/OneMinds/codex
```

This stores working files (your markdown codex) in Dropbox for sync while keeping the `.git` database outside Dropbox entirely. Alternatively, `git-remote-dropbox` (by anishathalye) uses the Dropbox API as a proper Git remote transport, avoiding filesystem sync conflicts.

For CLI access, Dropbox has no official macOS CLI. The best third-party option is **Maestral** (maestral.app) — a lightweight open-source Dropbox client with powerful CLI tools supporting `.gitignore`-style exclusion patterns, version viewing/restoring, and sync monitoring. `dbxcli` (under the Dropbox GitHub org, Go-based) also works but is sporadically maintained.

## How the four cloud storage options actually compare

For the specific requirements of the One Mind OS — offline Mac and iPhone access, CLI tools, version history, Git compatibility, MCP availability, and AI agent read/write:

**Dropbox** scores highest overall. It offers a real local folder, the best version history (up to 365 days), official and community MCP servers, `dbxcli`/Maestral for CLI, and excellent Python SDK. Main downsides: most expensive at $11.99/month for 2TB Plus, Git requires careful `.git` separation, and the newer File Provider on macOS has reported bugs.

**Google Drive** is the strongest runner-up. The `gws` CLI launched March 2026 with **native MCP server support built in** — a game-changer for agent integration. Must use Mirror mode (not Stream) for local file access and Git compatibility, which doubles storage usage. Best free tier at 15GB. Official MCP server (`@modelcontextprotocol/server-gdrive`) plus multiple community servers with full CRUD. Pricing matches iCloud at $9.99/month for 2TB.

**iCloud Drive** offers the best Apple ecosystem integration — seamless Files app access on iPhone, native Obsidian Sync alternative, $9.99/month for 2TB. But it's the worst for developers: **no public API**, no official CLI, no Drive-specific MCP server, and aggressive file eviction that makes Git compatibility the poorest of all four services. iCloud creates numbered duplicate files (e.g., `HEAD 2`) that corrupt Git refs.

**OneDrive** is best on pricing at $6.99/month for 1TB (includes Office apps), with a strong Microsoft Graph API and official MCP servers. However, CLI tooling is weak (rclone only), and Files On-Demand breaks Git similarly to other services.

The developer consensus across all services: **do not put Git repos in cloud sync folders**. The recommended pattern separates concerns: content files (markdown, configs) live in the cloud-synced folder for device access, `.git/` lives outside it via `--separate-git-dir`, and a Git remote on GitHub/Gitea handles backup and cross-device version control.

## The recommended knowledge base architecture for One Mind OS

Based on community production setups (GodMode, Zakk's Life OS, Nathan's Reef agent), the optimal stack layers five components:

**Layer 1 — Workspace structure.** Your codex markdown files (SOUL.md, HEARTBEAT.md, AGENTS.md, custom knowledge files) live in a dedicated workspace directory. Point `agents.defaults.workspace` in `openclaw.json` to your codex path. Use the standard file-per-concern discipline: AGENTS.md for operating rules, SOUL.md for identity, MEMORY.md for curated long-term facts, `memory/YYYY-MM-DD.md` for daily logs, `docs/` for reference material loaded on demand.

**Layer 2 — LanceDB for retrieval.** Install `memory-lancedb-pro` as the memory backend. It provides hybrid vector + BM25 search with cross-encoder reranking, automatic fact extraction from conversations, and multi-scope isolation so your Mason Android nodes, Mac Mini, and VPS agents can share or isolate memory contexts. Enable `memoryFlush` with a **40K token reserve** (not the default 20K) to ensure context survives compaction.

**Layer 3 — Git for version control.** Initialize git in the workspace with auto-commit via OpenClaw Mission Control (commits every 30 minutes with descriptive messages, provides rollback API) or a HEARTBEAT.md cron entry running `git add -A && git commit -m "auto-$(date +%Y%m%d-%H%M)" && git push`. Use Gitea on your DigitalOcean VPS as the private remote — Nathan's Reef pattern adds Woodpecker CI and TruffleHog pre-push hooks for secret scanning. The Git MCP Server gives agents direct version control capabilities.

**Layer 4 — Dropbox for device sync.** Use Dropbox as the sync layer with `--separate-git-dir` keeping `.git/` outside the synced folder. Mark the codex folder as "Available Offline" on both Mac and iPhone (requires Plus plan). Set the "One Mind Media" folder to Online-Only/Smart Sync for download-on-demand of large video files. Configure the Filesystem MCP Server pointing to the local Dropbox path for structured agent access, and add a Dropbox MCP server (ngs/dropbox-mcp-server) for API-level operations when the desktop client isn't running.

**Layer 5 — Obsidian as the human interface.** Install ObsidianClaw to get a chat sidebar in your vault — the agent reads and writes vault files directly, changes appear instantly. The official `obsidian` skill provides programmatic search, create, and move operations via `obsidian-cli`. For RAG over the full vault, `openclaw-rag-skill` (ChromaDB, fully local) or the built-in `knowledge-base` skill index your markdown files with vector search. Structure notes using Zettelkasten (atomic notes, cross-linked, tagged) following Zakk's pattern of 350+ literature notes.

**Orchestration layer — Hooks and cron.** The `afterTurn` hook (available since v2026.3.7's ContextEngine plugin interface) can trigger git commits after every conversation turn. HEARTBEAT.md should include: memory index coverage checks, git backup sync, system health monitoring. Nightly cron jobs handle index rebuilds, memory consolidation (the experimental "dreaming" sweep), and cross-agent knowledge reconciliation via QMD `extraCollections`.

## One critical insight from production deployments

The single most important lesson from community production setups is this: **automate fact extraction, don't rely on the agent to remember**. Dave Swift's widely-cited observation — "OpenClaw's promise of persistent, forever memory is, frankly, not true out of the box" — captures a real limitation. The fix is layered automation: `memoryFlush` captures context before compaction, `afterTurn` hooks save decisions immediately, LanceDB's `autoCapture` embeds every meaningful exchange, and nightly "dreaming" sweeps consolidate short-term signals into long-term memory. The read-write loop (agent reads context → does work → files findings back) is what makes a knowledge base compound rather than stagnate.

With this stack — OpenClaw workspace + LanceDB Pro + Git/Gitea + Dropbox sync + Obsidian interface + automated extraction hooks — you get agents that autonomously read and write files, version control with per-commit rollback, offline access on Mac and iPhone, CLI and MCP integration, and a separate media folder with download-on-demand. Jensen Huang called OpenClaw "the operating system for personal AI." Your One Mind OS is the kind of sovereign system that architecture was designed for.