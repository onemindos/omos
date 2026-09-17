# OneMind OS — Homepage Dashboard

## Access

| Method | URL |
|---|---|
| Local | http://localhost:3010 |
| Tailscale (all devices) | http://omos-mbp.tail717ba0.ts.net:3010 |

## Config Location

```
~/homepage/
├── docker-compose.yml
└── config/
    ├── settings.yaml    ← theme, layout, title
    ├── services.yaml    ← service tiles + ping checks
    ├── bookmarks.yaml   ← external link groups
    └── widgets.yaml     ← system stats, clock
```

## Git Repo

```bash
cd ~/homepage
git log        # view change history
git diff       # see uncommitted changes
git add . && git commit -m "description"
```

## Docker Management

```bash
# Status
docker ps --filter name=homepage

# Restart
cd ~/homepage && docker compose restart

# Update to latest image
cd ~/homepage && docker compose pull && docker compose up -d

# Logs
docker logs homepage -f
```

## launchd (auto-start on boot)

```
~/Library/LaunchAgents/com.onemind.homepage.plist
```

```bash
# Reload after editing plist
launchctl unload ~/Library/LaunchAgents/com.onemind.homepage.plist
launchctl load ~/Library/LaunchAgents/com.onemind.homepage.plist
```

## Adding a New Service

Append to `~/homepage/config/services.yaml` under the relevant group:

```yaml
- Service Name:
    href: http://host:port
    description: What it does
    icon: mdi-icon-name
    ping: http://host:port/health
```

No restart needed — Homepage hot-reloads config changes automatically.

## Current Services

### OneMind Core
- Hermes WebUI — localhost:8787
- Hermes Dashboard — localhost:9119
- Hermes API — localhost:8080
- Obsidian Vault — obsidian://open

### Smart Home
- Home Assistant — omos-ha.tail717ba0.ts.net:8123
- OpenClaw — localhost:3000
- NATS Monitoring — localhost:8222

### Cloud & Edge
- Atlas Lightsail — 3.234.227.222
- Synadia NGS — app.ngs.global
- Tailscale Admin — login.tailscale.com/admin

### Dev Tools
- Agent Web UI — localhost:3300
- Hermes Docs
- GitHub
- Todoist

## Why Homepage Over Dashy
- YAML-only config = git-friendly, Hermes can update it programmatically
- 1 open issue vs 59 (Dashy) — cleaner codebase
- Native HA + Docker + NATS integrations
- No database, no GUI state
