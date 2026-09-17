---
name: OpenClaw — Agent Workspaces
updated: 2026-05-11
---

# Agent Workspaces

Per-agent workspace files for agents running on omos-core. Each agent folder contains its identity, soul, memory, tools, and config files.

## Active Agent Workspaces (23 agents)

| Agent | Folder | Pillar |
|---|---|---|
| Legacy | [legacy (Legacy Workspace)/](legacy%20(Legacy%20Workspace)/) | SO — Orchestrator |
| Zenith | [zenith (Zenith Workspace)/](zenith%20(Zenith%20Workspace)/) | HP 25-49 |
| Unity | [unity (Unity Workspace)/](unity%20(Unity%20Workspace)/) | LE 50-74 |
| Eternal | [eternal (Eternal Workspace)/](eternal%20(Eternal%20Workspace)/) | GE 75-99 |
| Sovereignty | [sovereignty (Sovereignty Workspace)/](sovereignty%20(Sovereignty%20Workspace)/) | SO 00-24 |
| Sentinel | [sentinel (Sentinel Workspace)/](sentinel%20(Sentinel%20Workspace)/) | SO — Watchdog |
| Spartan | [spartan (Spartan Workspace)/](spartan%20(Spartan%20Workspace)/) | HP 28 |
| Oracle | [oracle (Oracle Workspace)/](oracle%20(Oracle%20Workspace)/) | HP 25-34 |
| Genesis | [genesis (Genesis Workspace)/](genesis%20(Genesis%20Workspace)/) | HP 25-26 |
| Pulse | [pulse (Pulse Workspace)/](pulse%20(Pulse%20Workspace)/) | HP 27 |
| Librarian | [librarian (Librarian Workspace)/](librarian%20(Librarian%20Workspace)/) | HP 27 |
| Haven | [haven (Haven Workspace)/](haven%20(Haven%20Workspace)/) | LE 50-74 |
| Helix | [helix (Helix Workspace)/](helix%20(Helix%20Workspace)/) | LE 52-56 |
| Eden | [eden (Eden Workspace)/](eden%20(Eden%20Workspace)/) | LE |
| Nomad | [nomad (Nomad Workspace)/](nomad%20(Nomad%20Workspace)/) | LE 50 |
| Empire | [empire (Empire Workspace)/](empire%20(Empire%20Workspace)/) | GE 75-99 |
| Apex | [apex (Apex Workspace)/](apex%20(Apex%20Workspace)/) | GE 79 |
| Creed | [creed (Creed Workspace)/](creed%20(Creed%20Workspace)/) | GE 75-81 |
| Trinity | [trinity (Trinity Workspace)/](trinity%20(Trinity%20Workspace)/) | SO 06 |
| Guardian | [guardian (Guardian Workspace)/](guardian%20(Guardian%20Workspace)/) | SO 04 |
| Grid | [grid (Grid Workspace)/](grid%20(Grid%20Workspace)/) | SO 02 |
| Overwatch | [overwatch (Overwatch Workspace)/](overwatch%20(Overwatch%20Workspace)/) | SO 07 |
| EagleEye | [eagleeye (Eagleeye Workspace)/](eagleeye%20(Eagleeye%20Workspace)/) | SO 02 |

## The 7 Workspace Files

| File | Purpose | Who writes it |
|------|---------|---------------|
| `IDENTITY.md` | Agent name, role, personality, operating style | Zeus / Claude Code |
| `SOUL.md` | Core values, mission, principles that never change | Zeus |
| `TOOLS.md` | Available tools reference — what the agent can do | Zeus / OpenClaw |
| `HEARTBEAT.md` | Periodic self-check prompt — keeps agent calibrated | Zeus |
| `AGENTS.md` | Fleet roster, routing rules, standing orders | Zeus / Claude Code |
| `USER.md` | Zeus's profile, preferences, communication style | Zeus |
| `MEMORY.md` | Evolving memory — what the agent has learned about Zeus | Agent (auto-written) |

## Backup & Version Control

All workspace files are auto-committed and pushed to GitHub by `watch.sh`.  
Repo: `github.com/OneMind-OS/openclaw-core-workspaces`  
See [[backup-repo|Backup Repo]] for health checks and restore procedures.

## Editing Workspace Files

1. Edit in Obsidian (vault) or via SSH on VPS
2. If editing in vault: git sync pushes to `/home/ubuntu/onemind/` — then manually copy to workspace: see `03 Command/ops-runbooks/maintenance/AGENT-WORKSPACE-SYNC.md`
3. If editing on VPS: `watch.sh` auto-commits and pushes to GitHub within ~3s
4. Reload agent after significant changes: `openclaw agent reload {id}`

## Legacy Workspace Folders (old domain agent architecture)

These folders exist in the vault from before the 6-agent consolidation. The workspace files are preserved for reference.

| Agent | Folder |
|---|---|
| Spartan | [spartan (Spartan Workspace)/](spartan%20(Spartan%20Workspace)/) |
| Oracle | [oracle (Oracle Workspace)/](oracle%20(Oracle%20Workspace)/) |
| Empire | [empire (Empire Workspace)/](empire%20(Empire%20Workspace)/) |
| Haven | [haven (Haven Workspace)/](haven%20(Haven%20Workspace)/) |
| Guardian | [guardian (Guardian Workspace)/](guardian%20(Guardian%20Workspace)/) |
| Grid | [grid (Grid Workspace)/](grid%20(Grid%20Workspace)/) |
| Trinity | [trinity (Trinity Workspace)/](trinity%20(Trinity%20Workspace)/) |
| Eden | [eden (Eden Workspace)/](eden%20(Eden%20Workspace)/) |
| Apex | [apex (Apex Workspace)/](apex%20(Apex%20Workspace)/) |
| Creed | [creed (Creed Workspace)/](creed%20(Creed%20Workspace)/) |

## Backup Repository

All workspace files are automatically version-controlled.

| Field | Value |
| --- | --- |
| **Repo** | `github.com/OneMind-OS/openclaw-core-workspaces` |
| **Watcher** | `~/.openclaw/watch.sh` (inotifywait, 3s debounce) |
| **Service** | `onemind-watch.service` (systemd) |
| **Status** | ✅ Active — last verified 2026-05-11 |

See [backup-repo.md](backup-repo.md) for full documentation and runbooks.

## Official Docs

- [Agent workspace](https://docs.openclaw.ai/concepts/agent-workspace.md)
- [SOUL.md template](https://docs.openclaw.ai/reference/templates/SOUL.md)
- [IDENTITY template](https://docs.openclaw.ai/reference/templates/IDENTITY.md)
- [HEARTBEAT.md template](https://docs.openclaw.ai/reference/templates/HEARTBEAT.md)
