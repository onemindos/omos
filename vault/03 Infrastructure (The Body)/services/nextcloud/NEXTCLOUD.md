# Nextcloud — Self-Hosted Cloud Platform

## Overview

| Field | Value |
|-------|-------|
| Version | 33.0.3.2 |
| Host | Mac (Docker) |
| Port | 8880 |
| Local URL | `http://localhost:8880` |
| Tailscale URL | `http://omos-mbp.tail717ba0.ts.net` |
| Admin User | `zeus` / `onemind2026` |
| Compose Path | `~/nextcloud/compose.yml` |
| Status | Active (testing) |

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Docker Compose Stack (~/nextcloud/)                     │
├─────────────────────────────────────────────────────────┤
│  nextcloud-app-1    │ nextcloud:stable │ Port 8880:80   │
│  nextcloud-db-1     │ mariadb:lts      │ Internal       │
│  nextcloud-redis-1  │ redis:alpine     │ Internal       │
└─────────────────────────────────────────────────────────┘
         │
         ├── Tailscale Serve (HTTP) → omos-mbp.tail717ba0.ts.net
         │
         └── Accessible from: Mac, iPhone (via Tailscale), Dell (via Tailscale)
```

## Storage Mounts

| Mount ID | Mount Point | Host Path | Mode | Type |
|----------|-------------|-----------|------|------|
| 4 | `/Documents` | `/Users/onemind/Documents` | Read-only | Personal (zeus) |
| 5 | `/Downloads` | `/Users/onemind/Downloads` | Read-write | Personal (zeus) |
| 6 | `/Local Drive` | `/Users/onemind/nextcloud-external` | Read-write | Personal (zeus) |
| 7 | `/OM PRIMARY` | `/Volumes/OM PRIMARY` | Read-only | Admin (all users) |

### Storage Stats (as of 2026-05-22)

| Mount | Folders | Files | Notes |
|-------|---------|-------|-------|
| Documents | — | ~111K+ | Full Mac Documents folder (236 GB, OneMind vault included) |
| Downloads | 798 | 822 | Mac Downloads folder |
| Local Drive | 1 | 8 | Shared/scratch space |
| OM PRIMARY | 8,949 | 179,320 | External drive (~3.2 TB: onemind-media + onemind-edu) |

> ~116 files on OM PRIMARY have encoding issues (accents: ñ, ö, é) and won't be accessible through Nextcloud.

## Docker Compose

```yaml
services:
  db:
    image: mariadb:lts
    restart: always
    command: --transaction-isolation=READ-COMMITTED
    volumes:
      - db_data:/var/lib/mysql
    environment:
      MYSQL_ROOT_PASSWORD: nextcloud_root_2026
      MYSQL_PASSWORD: nextcloud_pass_2026
      MYSQL_DATABASE: nextcloud
      MYSQL_USER: nextcloud

  redis:
    image: redis:alpine
    restart: always

  app:
    image: nextcloud:stable
    restart: always
    user: root  # Required for macOS volume permission access
    ports:
      - "8880:80"
    depends_on:
      db: { condition: service_healthy }
      redis: { condition: service_healthy }
    volumes:
      - nextcloud_data:/var/www/html
      - /Users/onemind/Documents:/external/documents:ro
      - /Users/onemind/Downloads:/external/downloads
      - /Users/onemind/nextcloud-external:/external/shared
      - /Volumes/OM PRIMARY:/external/om-primary:ro
    environment:
      MYSQL_PASSWORD: nextcloud_pass_2026
      MYSQL_DATABASE: nextcloud
      MYSQL_USER: nextcloud
      MYSQL_HOST: db
      REDIS_HOST: redis
      NEXTCLOUD_ADMIN_USER: zeus
      NEXTCLOUD_ADMIN_PASSWORD: onemind2026
      NEXTCLOUD_TRUSTED_DOMAINS: "localhost 127.0.0.1 onemind-nextcloud.tail717ba0.ts.net"
      PHP_MEMORY_LIMIT: 1024M
      PHP_UPLOAD_LIMIT: 10G

volumes:
  db_data:
  nextcloud_data:
```

## Enabled Apps

| App | Version | Purpose |
|-----|---------|---------|
| files_external | 1.25.1 | External storage (local mounts) |
| tasks | 0.17.1 | CalDAV task management |
| deck | 1.17.1 | Kanban boards |
| contactsinteraction | 1.14.1 | Contact interaction tracking |

## Trusted Domains

```
localhost
127.0.0.1
omos-mbp.tail717ba0.ts.net
omos-mbp
```

## Network Access

| Device | Access Method | URL |
|--------|--------------|-----|
| Mac (local) | Docker port | `http://localhost:8880` |
| Mac (Tailscale) | Tailscale Serve | `http://omos-mbp.tail717ba0.ts.net` |
| iPhone | Tailscale (must be connected) | `http://omos-mbp.tail717ba0.ts.net` |
| Dell | Tailscale | `http://omos-mbp.tail717ba0.ts.net` |

> HTTPS cert provisioning failed (DNS challenge 500 error). Using HTTP within Tailscale (encrypted tunnel, still secure).

## Mobile Setup (iOS)

### Nextcloud App
- Server: `http://omos-mbp.tail717ba0.ts.net`
- Username: `zeus`
- Password: `onemind2026`
- Requires: Tailscale connected on iPhone

### Native iOS Sync (Calendar + Tasks + Contacts)
- iOS Settings → Calendar → Accounts → Add Account → Other
- **CalDAV Server:** `http://omos-mbp.tail717ba0.ts.net/remote.php/dav`
- **CardDAV Server:** `http://omos-mbp.tail717ba0.ts.net/remote.php/dav`
- Username: `zeus` / Password: `onemind2026`
- Result: Calendar → iOS Calendar, Tasks → iOS Reminders, Contacts → iOS Contacts

## MCP Integration

Nextcloud MCP server (134 tools) configured for all AI agents:

```json
{
  "nextcloud": {
    "command": "uvx",
    "args": ["nextcloud-mcp-server", "run", "--transport", "stdio"],
    "env": {
      "NEXTCLOUD_URL": "http://localhost:8880",
      "NEXTCLOUD_USERNAME": "zeus",
      "NEXTCLOUD_PASSWORD": "onemind2026"
    }
  }
}
```

Configured in:
- `~/.claude/mcp.json` (global — Claude Code)
- `/Users/onemind/Documents/onemind/.claude/mcp.json` (workspace)
- `/Users/onemind/Documents/onemind/.vscode/mcp.json` (VS Code Copilot)

## CalDAV/CardDAV Endpoints

| Protocol | URL |
|----------|-----|
| WebDAV (Files) | `http://localhost:8880/remote.php/dav/files/zeus/` |
| CalDAV (Calendar) | `http://localhost:8880/remote.php/dav/calendars/zeus/` |
| CardDAV (Contacts) | `http://localhost:8880/remote.php/dav/addressbooks/users/zeus/contacts/` |
| Tasks Calendar | `http://localhost:8880/remote.php/dav/calendars/zeus/tasks/` |

## FastMail Import (Completed)

Imported from FastMail (zeuslegacy@fastmail.com):
- **2 contacts** (Leon Green, Legacy Test) → Nextcloud Contacts
- **35 calendar events** (Todoist calendar) → Nextcloud Calendar

## Tasks vs Deck

| Feature | Tasks (CalDAV) | Deck (Kanban) |
|---------|---------------|---------------|
| Sync | Apple Reminders, iOS Calendar, any CalDAV client | Nextcloud app only |
| Best for | Personal to-dos, reminders, quick captures | Project pipelines, visual boards |
| Mobile | Native iOS Reminders (no app needed) | Nextcloud mobile app only |
| API | CalDAV standard (VTODO) | REST API |
| Recommendation | Primary task system (replaces Todoist) | Visual project management |

## Operations

```bash
# Start/restart
cd ~/nextcloud && docker compose up -d

# Stop
cd ~/nextcloud && docker compose down

# Scan new files (after adding content)
docker compose exec app php occ files:scan zeus

# Scan specific mount
docker compose exec app php occ files:scan zeus --path="/zeus/files/OM PRIMARY"

# Check external storage
docker compose exec app php occ files_external:list --all

# View config
docker compose exec app php occ config:system:get trusted_domains

# Logs
docker compose logs app --tail 50
```

## Known Issues

- [ ] HTTPS cert provisioning failed — using HTTP over Tailscale (still encrypted)
- [ ] ~116 files with special characters (ñ, ö, é) inaccessible on OM PRIMARY
- [ ] macOS hidden system folders (.Spotlight-V100, .Trashes) error during scan (harmless)
- [ ] Container runs as `root` (required for macOS volume permissions — Docker Desktop limitation)
- [ ] OM PRIMARY only available when external drive is plugged in
- [ ] Full Documents scan takes 2+ minutes (236 GB, 111K+ files)

## Decision Log

| Date | Decision | Reason |
|------|----------|--------|
| 2026-05-22 | Port 8880 (not 8080) | NATS using 8080 |
| 2026-05-22 | Container runs as root | macOS 700 permissions on ~/Documents |
| 2026-05-22 | HTTP (not HTTPS) via Tailscale | DNS challenge failed; Tailscale tunnel is already encrypted |
| 2026-05-22 | Documents mounted read-only | Safety — prevent accidental modification of vault |
| 2026-05-22 | OM PRIMARY mounted read-only | Safety — 3.2 TB media library |
| 2026-05-22 | Tasks for daily to-dos, Deck for projects | Tasks syncs natively to iOS; Deck is visual-only |
