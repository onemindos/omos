---
id: ui-12-20260111053926
title: "Gitops Repository Structure   Foundation Setup Plan"
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
ai_summary: "Gitops Repository Structure   Foundation Setup Plan"
---

**GitOps Repository Structure - Foundation Setup Plan**

  
Legacy-AI/

├── .git/

├── .gitignore                    # Root ignore rules

├── .gitattributes                # LFS + merge strategies

├── README.md                     # Main documentation

├── docker-compose.yml            # Root orchestration

│

├── apps/                         # GitOps: Application deployments

│   ├── letta/

│   ├── n8n/

│   ├── home-assistant/

│   ├── emqx/

│   └── node-red/

│

├── infrastructure/               # GitOps: Infrastructure as code

│   ├── caddy/

│   ├── monitoring/

│   ├── scripts/

│   ├── vm/

│   └── kubernetes/              # Future K8s manifests

│

├── vault/                        # Legacy's unified knowledge base

│   ├── 00-Inbox/

│   ├── 01-Atlas/

│   ├── 02-Generational-Entrepreneurship/

│   ├── 03-Holistic-Performance/

│   ├── 04-Legacy-Evolution/

│   ├── 05-Innovative-Technology/

│   ├── 06-Systems/

│   ├── 07-Resources/

│   ├── 08-Legacy-AI-Context/

│   ├── Excalidraw/

│   └── Legacy-AI-Ops/

│

├── .obsidian/                    # Vault workspace (gitignored)

└── docs/                         # Repo documentation

    ├── architecture/

    ├── runbooks/

    └── setup/

  
**Git LFS: Allowed File Types**

✅ **Images (for markdown context & diagrams)**

- *.png - Screenshots, diagrams
- *.jpg / *.jpeg - Photos, documentation images
- *.gif - Animated demos
- *.webp - Optimized web images
- *.svg (large > 100KB) - Vector diagrams

✅ **Documents (selective context)**

- *.pdf (< 10MB) - Reference docs, guides
- *.excalidraw - Diagram source files

✅ **Data (AI training/context only)**

- *.csv (sanitized) - Structured data for AI
- *.json (> 1MB) - Large configs, datasets
- *.parquet - Optimized data formats

✅ **Media (limited, for training)**

- *.mp4 (< 50MB) - Tutorial videos
- *.mov (< 50MB) - Screen recordings
- *.webm - Compressed demos

  

❌ **Files NOT Allowed (Even with LFS)**

**Personal Media**

- *.mp3, *.wav, *.flac - Music files
- Large videos > 100MB
- Personal photos unrelated to docs

**Office Documents (use markdown instead)**

- *.docx, *.xlsx, *.pptx
- *.doc, *.xls, *.ppt

**Binaries & Archives**

- *.zip, *.tar.gz, *.rar
- *.exe, *.dmg, *.app

**Unencrypted Sensitive Data**

- Raw bank statements
- Unredacted financial exports
- PII-heavy datasets

  
  
———————————————————————————————————

  

  

📊 **What Gets Committed**

✅ **Encrypted** (via git-crypt)

- vault/02-Generational-Entrepreneurship/**/Finance/**
- vault/02-Generational-Entrepreneurship/**/Banking/**
- vault/03-Holistic-Performance/**/Health/**
- All *transactions*.csv, *bank*.csv files
- Files with -private or -confidential in name

✅ **Public** (readable by Letta AI)

- vault/01-Atlas/ - Framework & philosophy
- vault/08-Legacy-AI-Context/ - Agent configs, system prompts
- vault/Legacy-AI-Ops/ - Architecture docs
- All *.md files - Training context
- Images (via LFS) - Visual context
- Service configs in apps/ and [infrastructure](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-browser/workbench/workbench.html)

✅ **Via LFS** (efficient storage)

- PNG, JPG, GIF, WebP images
- PDF docs under 10MB
- Excalidraw diagrams
- Small training videos < 50MB

❌ **Never Committed**

- .obsidian/workspace* - Personal settings
- Secrets (.env, .pem, .key)
- Large media files
- Office documents
- OS files ([.DS_Store](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-browser/workbench/workbench.html))
- Node modules, Python cache

  
  
———————————————————

  

✅ **Encryption Key Backed Up**

**Automatic backups created:**

1. **Local backup**: ~/Documents/Backups/Encryption-Keys/legacy-ai-vault-key_20251119.key
2. **iCloud Drive**: ~/Library/Mobile Documents/com~apple~CloudDocs/Secure/legacy-ai-vault-key_20251119.key

**Finder opened** showing your key file - you can now drag it into 1Password!

🔴 **Critical Next Step: Add to 1Password**

Since Finder is now open with the key file visible:

1. Open **1Password**
2. Click **New Item** → **Secure Note**
3. Title: Legacy AI Vault Encryption Key
4. **Drag the file** from Finder into the 1Password note
5. Add tags: encryption, legacy-ai, critical
6. Save

📊 **Backup Status:**

- ✅ Original: ~/legacy-ai-vault-key.key (148 bytes)
- ✅ Documents backup: Created with date stamp
- ✅ iCloud Drive: Syncing to cloud
- ⏳ 1Password: **Do this now** (most important!)
- ⏳ USB backup: Optional but recommended

  
🔐 **What This Key Does**

**Purpose**: Encrypts sensitive files in your Legacy-AI repository so they're secure even if someone gets access to the Git history.

**What gets encrypted** (based on our [.gitattributes](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-browser/workbench/workbench.html) config):

- Financial records: vault/02-Generational-Entrepreneurship/**/Finance/**
- Banking data: vault/02-Generational-Entrepreneurship/**/Banking/**
- Client info: vault/02-Generational-Entrepreneurship/**/Clients/**
- Health records: vault/03-Holistic-Performance/**/Health/**
- Transaction CSVs: **/*transactions*.csv, **/*bank*.csv
- Private files: Any file with -private or -confidential in the name

**What stays unencrypted** (for AI training):

- Framework docs: [01-Atlas](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-browser/workbench/workbench.html)
- AI Context: vault/08-Legacy-AI-Context/
- System architecture: vault/Legacy-AI-Ops/
- All markdown files (*.md)
- Images and diagrams

🎯 **Why You Need This**

Without this key:

- ❌ You **cannot decrypt** encrypted files on a new machine
- ❌ Sensitive data is **permanently locked** if key is lost
- ✅ But hackers also can't read your financial/health data even if they clone the repo

With this key:

- ✅ You can clone the repo on any device
- ✅ Run git-crypt unlock ~/legacy-ai-vault-key.key
- ✅ All encrypted files become readable automatically

🔄 **How It Works**

  
💡 Think of it like:

GitHub repo = Your house

Git-crypt = Safe deposit box inside your house

This key = Only key to open that safe

Encrypted files = Your valuables in the safe

Even if someone breaks into your house (clones your repo), they can't open the safe (decrypt files) without the key!  
  
tl;dr: This key lets you keep sensitive personal/business data in Git alongside your code, encrypted. You've already backed it up to:

  

✅ Local Documents folder

✅ iCloud Drive

⏳ Still need to add to 1Password (most important!)  
  
—————

  

# On your Mac:

cd "~/Documents/Documents Backup/Legacy-AI"

  

# Edit any service

vim apps/letta/docker-compose.yml

  

# Commit and push

git add apps/letta/

git commit -m "feat: update Letta memory config"

git push

  

# Coolify automatically:

# 1. Receives webhook from GitHub

# 2. Pulls latest changes

# 3. Deploys to correct server(s)

# 4. Sends notification (optional)

# 5. Shows deployment status in UI

  # On your Mac:

cd "~/Documents/Documents Backup/Legacy-AI"

  

# Edit any service

vim apps/letta/docker-compose.yml

  

# Commit and push

git add apps/letta/

git commit -m "feat: update Letta memory config"

git push

  

# Coolify automatically:

# 1. Receives webhook from GitHub

# 2. Pulls latest changes

# 3. Deploys to correct server(s)

# 4. Sends notification (optional)

# 5. Shows deployment status in UI

  

# No manual intervention needed!

# No manual intervention needed!

  
Legacy AI Infrastructure

├── Cloud VM (3.215.29.128)

│   ├── Coolify Control Panel ⭐

│   ├── Letta (AI brain)

│   ├── n8n (workflows)

│   └── Caddy (reverse proxy)

│

├── Home Server (future)

│   ├── Coolify Agent

│   ├── Home Assistant

│   ├── Node-RED

│   └── Local data processing

│

└── Edge Devices (future)

    ├── Coolify Agent

    ├── EMQX (MQTT)

    └── Sensor data collectors

  

All managed from ONE UI

All deployed via git push

All monitored in real-time

---

## Related

- [[00-24 SO (Sovereign Operations)/05 Development (Code & Testing) (The Forge)/05 Development|05 Development]]
