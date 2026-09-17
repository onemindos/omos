# TOOLS.md — Grid Tool Usage Guide

## Codex Access

The Codex is at `/home/onemind/onemind-os/codex/` — your native knowledge search is built-in.

### ⛔ Codex Write Rules (CRITICAL)

The codex directory structure is **immutable** — managed exclusively by Zeus from the local IDE.

- **NEVER create new folders** in the codex. The numbered folder hierarchy (00-12, 25-34, 50-58, 75-84) is fixed.
- **NEVER rename or renumber** existing codex folders.
- **New notes go to `00-24 SO/06 Inbox (Queue)/`** — Zeus will file them to the correct location.
- **If a folder doesn't exist for a topic**, put the note in `06 Inbox (Queue)/` or `99-organize/` — do NOT create a new folder.
- The codex syncs via git. Any rogue folders you create get auto-committed and pollute the repo.

### memory_search (primary — use this first)

Semantically searches your memory + the Codex vault. Uses BM25+vector hybrid search.

```
memory_search "NATS JetStream configuration"   → infrastructure
memory_search "Docker compose setup"            → container management
memory_search "systemd service openclaw"        → service files
memory_search "PostgreSQL backup script"        → database ops
memory_search "VPS setup deployment"            → deployment docs
memory_search "OpenClaw skill structure"        → agent architecture
```

When to use: ANY question about infrastructure, services, deployments, or tech stack. Always search before answering.

### read_file (direct access)

```
read_file /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/08 Infrastructure/...
read_file /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/10 Development/...
read_file /home/onemind/onemind-os/infra/docker-compose.yml
```

## GitHub API (Repository Management)

Use for managing code repos:
- Check PR status and reviews
- Monitor CI/CD pipeline results
- Track issues and milestones
- Branch management and cleanup

## System Monitor

Use for VPS health checks:
- CPU, memory, disk usage
- Process status
- Service health
- Network connectivity

## Todoist (Task Management)

Use for creating and tracking infra tasks:
- Deployment schedules
- Maintenance windows
- Infrastructure upgrades
- Bug fix tracking

## File Tools (read / write / edit)

You have direct file access to the codex:

- `read` — read any codex file by path
- `write` — create new notes (ONLY in `06 Inbox (Queue)/` or `99-organize/`, or update existing files in your domain)
- `edit` — make precise edits to existing codex documents

To create a new note:
```
write /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/06 Inbox (Queue)/My-New-Note.md "# Title\n\nContent here"
```

To update an existing file in your domain:
```
edit /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/08 Infrastructure (Grid & APIs)/existing-note.md
```

## Key Infrastructure Reference

| Service | Port | Health Check |
|---------|------|-------------|
| OpenClaw Gateway | 18789 | ws://127.0.0.1:18789 |
| Agno AgentOS | 7777 | http://localhost:7777/health |
| NATS | 4222 | nats://localhost:4222 |
| PostgreSQL | 5432 | pg_isready |
| n8n | 5678 | http://localhost:5678 |
| LiveKit | 7880 | http://localhost:7880 |

## Domain Quick Reference

| Path | What's There |
|------|-------------|
| `00-24 SO/00 Framework/` | Architecture, doctrine, tech stack |
| `00-24 SO/01 Command/` | Dashboard infrastructure |
| `00-24 SO/02 Agents/` | Agent configs, memory, tools |
| `00-24 SO/06 Protocols/` | SOPs, operational procedures |
| `00-24 SO/08 Infrastructure/` | Cloud, storage, integrations |
| `00-24 SO/09 Automation/` | Workflows, n8n, automation |
| `00-24 SO/10 Development/` | Code, testing, DevOps |
| `00-24 SO/12 Research/` | R&D, robotics |
