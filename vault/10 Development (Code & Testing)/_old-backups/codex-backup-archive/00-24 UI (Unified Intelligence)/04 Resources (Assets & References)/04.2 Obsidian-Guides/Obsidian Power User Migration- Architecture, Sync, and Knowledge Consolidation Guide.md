---
id: ui-06-20260111053926
title: "Obsidian Power User Migration: Architecture, Sync, and Knowledge Consolidation Guide"
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
ai_summary: "Obsidian Power User Migration: Architecture, Sync, and Knowledge Consolidation G"
---

# Obsidian Power User Migration: Architecture, Sync, and Knowledge Consolidation Guide

**Consolidating from ClickUp, Notion, and other platforms into Obsidian is viable and increasingly well-supported**, but requires strategic planning. The key finding: Obsidian now has **official Letta AI integration** (alpha plugin), making it the strongest candidate for AI-augmented knowledge management. However, expect **5-10 hours of migration effort** for power users, and plan for a hybrid approach during transition—keeping ClickUp for team project management while migrating knowledge bases to Obsidian.

---

## Migration validation strongly favors Obsidian with caveats

The case for consolidating into Obsidian centers on **data ownership and AI integration potential**. According to a 2024 survey, **68% of users** who switched from Notion cited data sovereignty and customization as primary motivators. Plain Markdown files mean no vendor lock-in, offline-first operation, and direct access for AI agents.

What you'll lose migrating from **Notion**: native relational databases (Dataview plugin is a workaround), real-time collaboration, team workspaces, and polished mobile experience. From **ClickUp**: Gantt charts, time tracking, team collaboration, workflow automation, and client access controls. The critical insight: many successful power users run **hybrid setups**—Obsidian for knowledge/personal systems, ClickUp for active client work.

**Letta AI integration is production-ready** via the official plugin at `github.com/letta-ai/letta-obsidian`. Features include persistent memory across sessions, automatic vault sync with directory structure preservation, and real-time file change detection. Alternative AI integration uses the **Obsidian MCP Server** (`github.com/cyanheads/obsidian-mcp-server`) which requires the Local REST API plugin and exposes 10+ tools for note reading, searching, and frontmatter management. For VS Code integration, install the **Foam extension** for graph visualization and backlinks, plus the **Open Vault in VS Code** Obsidian plugin for quick context switching.

Migration tools: Use the built-in **Obsidian Importer plugin** for Notion (exports to ZIP, handles cleanup). No direct ClickUp migration tool exists—use the ClickUp Sync plugin for ongoing task sync, but expect manual export/conversion for historical data.

---

## Todoist decisively wins the task plugin comparison for AI integration

For AI agent accessibility, **Todoist is the clear winner**. It offers a fully documented REST API, official MCP server at `ai.todoist.net/mcp`, and TypeScript/Python SDKs. Letta can directly call Todoist's API as custom tools. TickTick's **official API is severely limited**—the TickTickSync Obsidian plugin uses an unofficial, reverse-engineered API described as a "moving target" that may break without warning.

| Feature | Todoist Plugins | TickTickSync | Native Tasks Plugin |
|---------|----------------|--------------|---------------------|
| **Two-way sync** | Partial (varies) | ✅ Full | N/A (local) |
| **Official API** | ✅ Full REST/Sync | ⚠️ Unofficial only | N/A |
| **AI/MCP support** | ✅ Official server | ❌ None | ❌ None |
| **Maintenance** | Active | Active | Very active |
| **Scalability** | High | Medium | High (50k+ tested) |
| **Free tier limits** | 5 projects | 19 lists | Unlimited |

The **Todoist Sync plugin** (`github.com/jamiebrynes7/obsidian-todoist-plugin`) is query-based, materializing Todoist filters into Obsidian code blocks. Note: the "Ultimate Todoist Sync" fork broke in April 2025 due to API changes and has active bug reports. For **pure local task management**, the native **Tasks plugin** (`github.com/obsidian-tasks-group/obsidian-tasks`) has been tested with **50,000+ task cards** and offers complete data ownership.

**Recommendation**: Use Todoist for external task management with AI integration potential; use the native Tasks plugin for vault-embedded tasks that don't need external sync.

---

## Cloud sync platform selection depends entirely on your ecosystem

**Obsidian Sync ($4-8/month)** remains the most reliable option with native integration, end-to-end encryption, and selective plugin sync. However, several alternatives work well for specific use cases.

### Cross-platform (macOS, iOS, Windows, Linux) recommendations

**For full cross-platform including Linux**: Obsidian Sync is the path of least friction. Budget alternative: **Remotely Save plugin + Dropbox** (free 2GB tier) or **S3 bucket** (<$1/month with technical setup).

**Apple ecosystem only**: iCloud works "mostly fine" but watch for sync delays, files stuck uploading, and the 10-30 second "Waiting for iCloud to synchronize" message on iOS. Critical issue: iCloud de-prioritizes Obsidian folders and downloads files on-demand, breaking search. Workaround: disable "Optimize Mac Storage."

**Privacy-focused**: **Syncthing** (peer-to-peer, free) is excellent for desktop sync with full plugin settings sync. iOS limitation: requires **Möbius Sync** ($5 one-time) which cannot run continuously in background.

**Dropbox critical warning**: Does NOT sync files with certain emoji characters (specifically emojis starting with hexadecimal "f0" like 📚, 📖). Use the "Filename Emoji Remover" plugin if choosing Dropbox.

| Platform | iOS Native | Conflict Handling | Version History | Best For |
|----------|-----------|-------------------|-----------------|----------|
| Obsidian Sync | ✅ | Auto-merge | ✅ Extensive | Hassle-free setup |
| iCloud | ✅ | Silent overwrite | 30 days | Apple-only, small vaults |
| Dropbox | ❌ (needs plugin) | Creates duplicates | 30-365 days | Desktop + Remotely Save |
| Syncthing | ⚠️ (limited) | Duplicates | None | Privacy, desktop sync |
| S3 | ✅ (via plugin) | Smart (PRO) | Configurable | Technical users, cheap |

---

## Obsidian Publish enables same-vault selective publishing

**Yes, you can publish from the same vault** you use privately—nothing publishes automatically. When clicking "Publish changes," you explicitly select each note. The security model relies on manual selection plus frontmatter flags.

**Real Obsidian Publish sites**:
- `publish.obsidian.md/hub` — Obsidian community documentation
- `notes.nicolevanderhoeven.com` — PKM workflows, custom domain
- `help.obsidian.md` — Official documentation
- `untested.sonnet.io` — Technical blog via Vercel proxy

**Pricing**: $8-10/month per site, 4GB storage, with 40% discount for students/educators.

### Best free alternative: Quartz

**Quartz** (`github.com/jackyzha0/quartz`) offers graph view, backlinks, full-text search, and wikilink support—free hosting on GitHub Pages, Cloudflare Pages, or Vercel. Setup requires Node.js and Git knowledge. Use the **ExplicitPublish** filter so only notes with `publish: true` frontmatter deploy.

For **Circle.so integration**: Create Obsidian Publish site with custom domain, set password protection, share password in Circle private spaces. Limitation: no SSO integration; password protection is site-wide, not per-note.

**Recommended architecture for public/private separation**:
```
vault/
├── Public/        # Publishable content
│   ├── Blog/
│   └── Wiki/
├── Private/       # Never published
│   ├── Journal/
│   └── Work/
└── Attachments/
    ├── public-*/  # Prefixed for publishing
    └── private-*/ # Never synced
```

---

## Vault scalability peaks at 10,000-20,000 markdown files

Obsidian has **no hard file count limit**—stress tests confirm 20,000 markdown files "works okay" on fast desktops. One user tested **280,986 files** (with significant degradation). However, **total file count matters more than markdown count**: 2,700 markdown files with 570,000 attachments becomes unusable.

### Performance degradation order

1. **Initial loading/cache** — first to suffer
2. **Search and [[]] link completion** — 4 seconds per keystroke reported at 10k+ files
3. **Graph view** — CPU/GPU intensive
4. **Startup time** — affected by plugin count

**Memory patterns**: Base ~70-100MB per vault; large vaults (20k notes + 50k images) consume **6-10GB RAM**. V8 heap limit of ~4GB means you may need to split vaults.

**Plugin load impact**: **60+ active plugins** noticeably slows command palette; **100+ plugins** causes "excessive RAM usage." Known heavy plugins: **Dataview** (up to 200% CPU on large vaults), Excalidraw, Omnisearch. Use Settings → General → Advanced → **Startup Time overlay** to identify slow plugins.

### Security considerations

Plugins inherit Obsidian's full file system access—no sandboxing. Community plugin review is not a security audit. Notable vulnerabilities (all patched): path handling file exfiltration (fixed v1.2.8), RCE via obsidian:// hook (fixed v0.15.5). **Encryption options**: OS-level (FileVault, BitLocker), container (Cryptomator), plugin (Cryptsidian for AES-256), or Obsidian Sync's end-to-end encryption.

---

## Vault-to-vault sync works through Git-based selective publishing

Automatic sync between private and public vaults is achievable via **Quartz + frontmatter filtering**. Add `publish: true` to notes, configure `ExplicitPublish` filter in `quartz.config.ts`, and only flagged notes deploy.

**Obsidian Git plugin** (`github.com/Vinzent03/obsidian-git`, 9.1k+ stars) provides automatic commit-and-sync:
- Recommended settings: 5-10 minute backup interval, pull on startup enabled
- Mobile support is unstable (uses JavaScript Git reimplementation, no SSH)
- Workspace files cause conflicts—add to `.gitignore`

### Essential configuration files

**.gitignore (standard)**:
```gitignore
.obsidian/workspace.json
.obsidian/workspace-mobile.json
.obsidian/cache
.trash/
.DS_Store
```

**.gitattributes (with LFS)**:
```gitattributes
*.md text eol=lf
*.png filter=lfs diff=lfs merge=lfs -text
*.jpg filter=lfs diff=lfs merge=lfs -text
*.pdf filter=lfs diff=lfs merge=lfs -text
*.mp3 filter=lfs diff=lfs merge=lfs -text
*.mp4 filter=lfs diff=lfs merge=lfs -text
```

**Git LFS thresholds**: Use for binary files >50KB that change frequently. GitHub free tier: 1GB storage, 1GB bandwidth. GitLab offers 10GB free with unlimited bandwidth.

---

## Binary file handling requires dedicated attachment strategy

**Folder structure recommendation**:
```
vault/
├── 00-Inbox/
├── 01-Daily/
├── 02-Projects/
│   └── project-name/assets/
├── 03-Areas/
├── 04-Resources/
├── 05-Archive/
├── Attachments/
│   ├── images/
│   ├── pdfs/
│   ├── audio/
│   └── video/
└── Templates/
```

Configure Settings → Files & Links → Default location for new attachments to **"In the folder specified below"** → `Attachments`. This centralizes media for easier Git LFS tracking and exclusion.

**Vault size thresholds**:
| Size | Performance | Recommendation |
|------|-------------|----------------|
| <1GB | Excellent | No special handling |
| 1-5GB | Good | Use LFS for binaries |
| 5-10GB | Moderate | LFS + exclude from mobile |
| >10GB | Degraded | Split vault or external linking |

For large media, use **external linking** instead of embedding: `[Watch video](file:///path/to/video.mp4)` rather than `![[video.mp4]]`. This keeps the vault lean while maintaining references.

---

## Conclusion: Actionable implementation path

**Week 1**: Set up Obsidian with Letta AI plugin and Obsidian Git. Create folder structure with clear public/private separation. Configure `.gitignore` and `.gitattributes` with LFS.

**Week 2-3**: Migrate Notion using built-in Importer. Run Notion-to-Obsidian-Converter script to clean UIDs and fix links. Keep ClickUp active for team/client work.

**Week 4**: Implement Todoist sync for external task management. Set up native Tasks plugin for vault-embedded tasks. Configure Quartz for public knowledge sharing with `publish: true` frontmatter filtering.

**Ongoing**: Use Obsidian Sync or Syncthing (Apple ecosystem) / Remotely Save + S3 (cross-platform) for sync. Monitor vault size—split if exceeding 10k files or 5GB. Limit active plugins to under 50 for optimal performance.

**Key tool links**:
- Letta-Obsidian: `github.com/letta-ai/letta-obsidian`
- Obsidian MCP Server: `github.com/cyanheads/obsidian-mcp-server`
- Todoist MCP: `ai.todoist.net/mcp`
- Obsidian Git: `github.com/Vinzent03/obsidian-git`
- Quartz: `github.com/jackyzha0/quartz`
- Remotely Save: Search Obsidian Community Plugins