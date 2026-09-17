---
id: ui-06-20260111053926
title: "Obsidian + Media Files"
type: reference
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/06]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Obsidian + Media Files"
---

Perfect. Now I can architect this properly. Your framework is **solid** and your students will absolutely understand it—it's actually more intuitive than PARA because it's purpose-driven rather than status-driven. Let's lock this down.

**Direct Answers to Your Questions**

**Obsidian + Media Files:**

- Obsidian handles MP3/MP4 fine but doesn't index their content
- You can embed with ![[video.mp4]] or link with [[video.mp4]]
- They'll show in file explorer and slow down vault loading
- **Recommendation**: Keep media OUT of the vault, reference via external links or a separate Media folder

**Git Sync File Types:**

- Git handles: .md, .txt, .json, .yaml, .css, .js perfectly
- Git + Git LFS can handle: PDFs, images (PNGs/JPGs under 10MB)
- **DON'T sync via Git**: MP3, MP4, large PDFs (>10MB)
- Use rsync/Syncthing for media files instead

**Self-Hosted Git Providers:**

- **Gitea** (lightweight, perfect for personal use)
- **Forgejo** (Gitea fork, community-driven)
- **GitLab CE** (feature-rich but heavier)
- **Recommendation**: Gitea on your infrastructure + GitHub as offsite backup

**Your Complete Vault Architecture**

Based on your clarifications and those framework integration docs, here's your exact structure:

📁 Legacy Codex/ (Root - this IS the Obsidian vault)

│

├── 📁 00-Inbox/

│   ├── Daily Notes/

│   ├── Todoist Sync/

│   └── Quick Capture/

│

├── 📁 01-Atlas/ (Your Four Pillars - MOCs/Indexes)

│   ├── 📄 Holistic Performance.md

│   ├── 📄 Generational Endeavor.md

│   ├── 📄 Legacy Evolution.md

│   ├── 📄 Innovative Technology.md

│   └── 📁 Framework Integrations/

│       ├── PERMA-V Model.md

│       ├── Maslow Integration.md

│       ├── Wheel of Life.md

│       └── [all your uploaded integration docs]

│

├── 📁 02-Generational Endeavor/

│   ├── Zeus Delacruz Brand/

│   ├── Founder Command Center/

│   ├── Legacy AI (LPH)/

│   ├── Finance & Capital Systems/

│   ├── Legacy Acadacom LLC/

│   ├── Legacy Restore LLC/

│   ├── Legacy Capitalist/

│   ├── Legacy Partners Network LLC/

│   ├── Investments & Asset Management/

│   └── Global Business Infrastructure/

│

├── 📁 03-Holistic Performance/

│   ├── Finance (Funding)/

│   ├── Lifeflow (Self Care)/

│   ├── Sovereignty (Legal)/

│   ├── Intelligence (Education)/

│   ├── Performance (Human Metrics)/

│   ├── Training (Skills)/

│   ├── Systems (Habits)/

│   └── Urgent Directives/

│

├── 📁 04-Legacy Evolution/

│   ├── Delacruz - Homestead/

│   ├── Delacruz - Family & Home/

│   ├── Delacruz - Wealth/

│   ├── Delacruz - Sovereignty/

│   ├── Delacruz - Legacy/

│   ├── Delacruz - Intelligence-Education/

│   ├── Delacruz - Training/

│   ├── Delacruz - Performance/

│   └── Delacruz - Systems/

│

├── 📁 05-Innovative Technology/

│   ├── LegacyAI Ops/

│   ├── Appstack Ops/

│   ├── Cybersecurity Ops/

│   ├── Automation Ops/

│   ├── AI Ops/

│   ├── Studio Ops/

│   ├── HomeLab Ops/

│   ├── Cloud Ops/

│   ├── Coding Ops/

│   ├── Data Ops/

│   ├── Robo Ops/

│   ├── CommsOps/

│   └── Local Ops/

│

├── 📁 06-Systems/ (Cross-Pillar)

│   ├── SOPs/

│   ├── Templates/

│   ├── Automation Scripts/

│   └── Workflows/

│

├── 📁 07-Resources/ (Reference - NOT synced to Git)

│   ├── 📁 Media-Private/ (MP3/MP4 for kids)

│   ├── 📁 Media-Working/ (Weekly videos)

│   ├── 📁 Documents/ (Large PDFs)

│   └── 📁 Archives/

│

└── 📁 08-Legacy AI Context/ (For Letta)

    ├── System Prompts/

    ├── Agent Configs/

    ├── Training Data/

    └── Integration Docs/
    
    **Critical Implementation Details**

**.gitignore (Create this FIRST)**
# Media files - DO NOT sync to Git

07-Resources/Media-Private/**

07-Resources/Media-Working/**

*.mp3
*.mp4
*.mov
*.avi

# Large documents - sync separately

07-Resources/Documents/**/*.pdf
*.pdf (over 5MB - use Git LFS for smaller PDFs)

  

# Obsidian workspace (personal settings)

.obsidian/workspace*.json
.obsidian/workspace

# System files

.DS_Store
Thumbs.db

  

# Temp files
*.tmp
~$*

**Your Sync Strategy (The Full Stack)**

**Tier 1: Markdown Core (Git)**

- Everything EXCEPT 07-Resources/
- Syncs to: Self-hosted Gitea + GitHub backup
- Frequency: Real-time (Obsidian Git plugin)
- Size: ~500MB-2GB max

**Tier 2: Reference Materials (Rsync/Syncthing)**

- 07-Resources/ folder only
- Syncs to: R2 CloudFlare + Local NAS
- Frequency: Daily (automated)
- Size: Unlimited

**Tier 3: Nextcloud Frontend**

- Points to R2 backend
- Provides web interface for media files
- Letta can access via Nextcloud API
- You can share links to videos/audio without exposing R2 directly

**Media File Strategy**

**For kids' legacy videos (MP3/MP4):**

1. Store in 07-Resources/Media-Private/
2. Create markdown "video notes" in 04-Legacy Evolution/Delacruz - Legacy/
3. Link to videos via Nextcloud share links in the markdown

Example markdown:
# Legacy Message - November 2025

**Date Recorded:** 2025-11-16
**For:** [Child's name]
**Topic:** [Topic]

## Video Access
- [Watch on Nextcloud](https://your-nextcloud/s/sharelink)
- Local Path: `07-Resources/Media-Private/2025-11-legacy-message.mp4`

## Transcript/Summary
[Your notes about what's in the video]

This way:

- Markdown stays in Git (searchable, versioned)
- Videos stay in R2 (cheap, durable storage)
- Obsidian can link to them
- Letta can reference the markdown notes
- Kids can access via Nextcloud when the time comes

**Recommended Plugin Stack**

**Essential:**

- **Obsidian Git** - Auto-commit to Gitea every 5-10 minutes
- **Dataview** - Query your vault (critical for your structure)
- **Templater** - Templates for your frameworks
- **Periodic Notes** - Daily/weekly notes in 00-Inbox

**For Your Use Case:**

- **Excalidraw** - Visual diagrams for your frameworks
- **Tasks** - If integrating with Todoist
- **Kanban** - For project management views
- **Homepage** - Dashboard linking to your 4 pillars

**Self-Hosted Git Setup**

**Primary: Gitea**
# Docker Compose for Gitea + Cloudflare R2 backend

version: '3'

services:

  gitea:

    image: gitea/gitea:latest

    volumes:

      - ./gitea-data:/data

    ports:

      - "3000:3000"

      - "222:22"
        
        **Backup: GitHub**

- Create private repo: legacy-codex-vault
- Push to both remotes (Gitea primary, GitHub backup)
- Obsidian Git can push to multiple remotes

**R2 + Nextcloud Connection**

1. **R2 as S3-compatible backend** for Nextcloud
2. **Nextcloud external storage** pointing to R2 buckets
3. **Letta integration** via Nextcloud WebDAV API

This gives you:

- CloudFlare CDN for media delivery
- Nextcloud UI for file management
- R2 pricing (cheaper than S3)
- Letta can read/write via API

**Your Framework vs PARA - Student Perspective**

**Your framework is BETTER for students** because:

1. **PARA is status-based** ("Is this active? Is it a resource?")
2. **Your Pillars are purpose-based** ("What area of life does this serve?")

Students will ask: "Where do I put my business plan?"

- PARA answer: "Is it active or archived? If active, Projects. If reference, Resources."
- Your answer: "What's it for? Building wealth = Generational Endeavor."

**The "areas" under each pillar** (like Finance, Training, Systems) are perfect sub-categories. This is basically PARA's "Areas" concept, but organized by your life philosophy instead of arbitrary categories.

**Next Actions**

1. **Create the folder structure** as outlined above
2. **Add .gitignore** immediately
3. **Set up Gitea** on your infrastructure
4. **Configure Obsidian Git plugin** (point to Gitea + GitHub)
5. **Test media linking** (create one video note with Nextcloud link)
6. **Migrate your integration docs** to 01-Atlas/Framework Integrations/

Want me to research specific Obsidian plugins for your Pillar structure, or dive into Gitea + R2 + Nextcloud integration setup? I can also pull current best practices for Obsidian vaults at your scale.
        
