---
id: ui-02-20260111053925
title: "Nextcloud + Obsidian + Nas"
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
ai_summary: "Nextcloud + Obsidian + Nas"
---

![[The Nextcloud + Obsidian + NAS Architecture.png]]![[How Nexterowd Makes Las HaDDen.png]]**Virtual Files Mode:**

- Files appear in Finder/Explorer but aren't downloaded
- Double-click = downloads on-demand
- Right-click → "Make available offline" for specific files
- Identical to Dropbox Smart Sync behavior

  
  
**2. OBSIDIAN LINKS THAT NEVER BREAK**

Here's the key: Use **relative paths** within your vault:

  
# In your Obsidian note (anywhere in vault)

  

# These ALWAYS work because they're relative to vault root:

[[00_Daily/2025-01-15]]

![[03_Attachments/diagram.png]]

[Q3 Report](../Active_Files/Current_Projects/Q3_2024.pdf)

  

# For Archive files (even if not downloaded):

[Historical Data](../Archive/Legacy_Evolution/2020_analysis.pdf)

```

  

**When you move files:**

- Move within Nextcloud web interface → paths update automatically

- Move within Nextcloud desktop → paths stay consistent

- Obsidian's "Automatically update internal links" keeps wiki-links working

- File:// URLs are relative to vault root → portable across machines

  

**3. UNIFIED BACKUP STRATEGY**

```

BACKUP LAYERS (3-2-1 Rule)

═══════════════════════════════════════════════════════

  

LAYER 1: Primary Storage

├── Nextcloud Server (Portainer)

│   ├── Live data in PostgreSQL

│   └── Files on NAS + S3

│

LAYER 2: NAS Replication

├── Synology/QNAP at Farm

│   ├── Nextcloud External Storage mount

│   └── Scheduled rsync from server

│

LAYER 3: External Backup

├── External HDD (weekly rotation)

│   ├── Nextcloud backup command

│   └── Borgbackup for versioning

│

LAYER 4: Cloud Backup

└── Wasabi S3 (separate bucket from primary)

    ├── Nextcloud's built-in backup

    └── Lifecycle policies for retention  
  
  
**4. NAS INTEGRATION (Multiple Options)**

  

**Option A: Nextcloud External Storage**

```

Nextcloud → Settings → External Storage

├── Add NFS/SMB mount to your NAS

├── Make Archive/ folder point to NAS

└── Nextcloud manages sync between NAS + S3

**Option B: Direct NFS Mount**

  

  

bash

_# Mount NAS directly in Nextcloud Docker_

docker-compose.yml:

  volumes:

    - nfs-share:/var/www/html/data/Archive:shared

```

  
  ![[5. PROGRESSIVE DOWNLOAD WORKFLOW 4.png]]

**Option C: Nextcloud Sync Client on NAS**

```

_# If your NAS supports Docker_

Install Nextcloud client in Docker on NAS

Sync Archive/ folder bidirectionally

	Acts as local cache of S3 data
Nextcloud Workflow App automates this:  
# .nextcloud/workflows/convert_and_archive.yaml

- name: Auto-convert PDFs

  events: file_created

  conditions:

    - path: Active_Files/Processing/*.pdf

  actions:

    - run_script: /scripts/extract_text.py

    - move_file: Archive/Processed/

    - tag: converted

```

  

### **THE COMPLETE INTEGRATION PICTURE**

```

YOUR MULTI-LOCATION SETUP

═══════════════════════════════════════════════════════

  

Office (Mac Desktop)

├── Nextcloud Client: Full sync

├── Obsidian: Points to ~/Nextcloud/Legacy_Codex_Command/Obsidian_Vault/

├── Working Files: Active_Files/ synced

└── Archive: Virtual files (on-demand)

  

Home (Mac Laptop)  

├── Nextcloud Client: Full sync

├── Obsidian: Points to ~/Nextcloud/Legacy_Codex_Command/Obsidian_Vault/

├── Working Files: Active_Files/ synced

└── Archive: Virtual files (on-demand)

  

Farm/Warehouse (NAS)

├── Nextcloud External Storage: Archive/ stored here

├── S3 Sync: Backs up to Wasabi nightly

└── Offline Access: Full copy locally

  

Mobile (iOS)

├── Nextcloud App: Browse all files

├── Obsidian Mobile: Vault synced via Nextcloud

└── Download on demand: Tap to open PDFs  
  
  
**LINK INTEGRITY ACROSS MOVES**

**How Nextcloud + Obsidian maintain links:**

1. **Within Obsidian vault** (wiki-links):

- Obsidian tracks files by unique ID
- Moving files auto-updates [[wiki-links]]
- Works even if Nextcloud isn't syncing yet

3. **To files outside vault** (markdown links):

- Use relative paths: [Doc](../../Archive/file.pdf)
- When you move in Nextcloud, path structure persists
- Search/replace if you restructure folders

5. **Absolute paths via Nextcloud**:

- Nextcloud generates WebDAV URLs
- These URLs never change: https://nextcloud.legacyplatform.com/remote.php/dav/files/mrdelacruz/Archive/file.pdf
- Link via URL = always works, even from mobile

  
Pro Move - Use Dataview for Dynamic Links:  
<!-- In Obsidian -->

```dataview

TABLE file.link AS "Document", file.mtime AS "Modified"

FROM "Archive"

WHERE contains(file.name, "Q3") AND contains(file.name, "2024")

SORT file.mtime DESC

```  
  
This queries Nextcloud in real-time, so even if files move, the query finds them.

**MY VERDICT: YOUR ARCHITECTURE IS PERFECT FOR NEXTCLOUD**

**What you get that Dropbox/iCloud can't provide:**

✅ **Selective sync** - Virtual files mode, identical to Smart Sync ✅ **Unified structure** - Everything under one roof ✅ **NAS integration** - External storage or direct mount ✅ **External backup** - Built-in backup tools + scripts ✅ **No broken links** - Relative paths + WebDAV URLs + Obsidian intelligence ✅ **Progressive workflow** - Workflow automation for convert/archive ✅ **Infinite scale** - S3 backend grows forever ✅ **Sovereignty** - You own every component ✅ **Cost control** - Wasabi: $5.99/TB vs iCloud $10/TB ✅ **AI integration** - Direct API access for Letta agents

The migration path:
![[PHASE 1 Setup Nextcloud.png]]