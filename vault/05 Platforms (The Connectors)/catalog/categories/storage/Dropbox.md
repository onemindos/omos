# Dropbox

> **Cloud storage platform for vault backup and file access**

## Overview

Dropbox is a cloud storage and sync platform. Planned as Phase 4 of the One Mind sync architecture — a backup mirror for the full onemind-codex vault (all files including PDFs and images) with 180-day version history and REST API access via `dbxcli`.

## Status

| Attribute   | Value                         |
| ----------- | ----------------------------- |
| Status      | ⬜ Planned                    |
| Category    | Storage                       |
| Type        | Cloud Storage / Sync          |
| Integration | dbxcli, REST API, macOS sync  |

## Use Cases in OneMind

1. **Vault Backup** — Full mirror of onemind-codex (PDFs, images, markdown)
2. **180-Day History** — Extended file history beyond Obsidian Sync's 60 days
3. **API Access** — Programmatic file retrieval and upload via agents
4. **Emergency Access** — Web interface if Obsidian Sync is unavailable

## Key Features

- 180-day version history (Plus plan)
- REST API v2 with longpoll for change detection
- `dbxcli` command-line tool
- macOS native Finder integration
- 1 TB (Plus) / 2 TB (Professional)

## Setup Notes

```bash
# Install CLI
brew install dbxcli

# Protect .git from Dropbox sync
xattr -w com.dropbox.ignored 1 /path/to/.git
xattr -w com.dropbox.ignored 1 /path/to/.obsidian

# Basic CLI usage
dbxcli ls /onemind-codex/
dbxcli get /path/to/file
```

## Links

- Website: [dropbox.com](https://dropbox.com)
- API Docs: [dropbox.github.io/dropbox-api-v2-explorer](https://dropbox.github.io/dropbox-api-v2-explorer)
- CLI: [github.com/dropbox/dbxcli](https://github.com/dropbox/dbxcli)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
