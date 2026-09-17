---
title: "Forgejo"
type: reference
status: planned
tags: [platform, development, source-control, self-hosted, foss, federation]
---
# Forgejo

> Community-governed, federated Git forge. Fork of Gitea — guaranteed FOSS forever. Lightweight, self-hostable, sovereignty-first.

## Use Case in OneMind OS
- Sovereign code hosting — no corporate dependency
- Federation via ActivityPub (repos discoverable across instances)
- Self-hosted on any node (runs on a Pi, ~100MB RAM)
- Webhooks → NATS for CI/CD events

## Key Features
| Feature | Notes |
|---------|-------|
| Git hosting | Repos, PRs, issues, wiki, releases |
| Federation | ActivityPub — cross-instance discovery |
| CI/CD | Forgejo Actions (GitHub Actions compatible YAML) |
| Container Registry | OCI/Docker images per-repo |
| Packages | npm, PyPI, Go, Maven, Cargo |
| Mirror | Pull/push mirror from GitHub/GitLab |
| OAuth2 | Built-in provider + consumer |
| Lightweight | Single binary, SQLite or Postgres |

## Deploy
```bash
# Docker
docker run -d \
  --name forgejo \
  -p 3000:3000 -p 222:22 \
  -v forgejo-data:/data \
  codeberg.org/forgejo/forgejo:latest

# Or binary
wget https://codeberg.org/forgejo/forgejo/releases/download/v*/forgejo-*-linux-arm64
chmod +x forgejo-*
./forgejo-* web
```

## CLI
```bash
# Uses tea CLI (Gitea/Forgejo compatible)
tea login add --url https://git.local --token TOKEN
tea repo list
tea pr create --title "..." --base main
tea issue create --title "..."
```

## Why Forgejo over Gitea
- Non-profit governance (Codeberg e.V.) — can't go proprietary
- Federation support (ActivityPub)
- Same codebase, same performance, same compatibility
- Codeberg.org proves it at scale (100k+ users)

## Integration Points
- Webhooks → NATS subjects for CI/CD events
- Forgejo Actions = GitHub Actions syntax (reuse existing workflows)
- Mirror from GitHub for backup/sovereignty
