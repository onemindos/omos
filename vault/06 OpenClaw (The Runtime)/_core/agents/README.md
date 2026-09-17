---
title: OpenClaw — Agent Registry
type: reference
status: active
updated: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/agents]
---

# Agent Registry

All 23 agents running on omos-core. Live source: `~/.openclaw/openclaw.json` → `agents.list`.

## Super-Agents (Sonnet 4.6 — full domain orchestrators)

| Agent | ID | Model | Pillar | Workspace | Discord | Telegram |
|---|---|---|---|---|---|---|
| ⚡ Legacy | `main` | Sonnet 4.6 [gh] | SO 00-99 | `workspace/` | `#legacy-chat` | topic 292 |
| 🏔 Zenith | `zenith` | Sonnet 4.6 [gh] | HP 25-49 | `workspace-zenith/` | `#zenith-chat` | topic 318 |
| 🌿 Unity | `unity` | Sonnet 4.6 [gh] | LE 50-74 | `workspace-unity/` | `#unity-chat` | topic 319 |
| 👑 Eternal | `eternal` | Sonnet 4.6 [gh] | GE 75-99 | `workspace-eternal/` | `#eternal-chat` | topic 320 |
| 🌐 Sovereignty | `sovereignty` | Sonnet 4.6 [gh] | SO 00-24 | `workspace-sovereignty/` | `#sovereignty-chat` | topic 321 |
| 👁️ Sentinel | `sentinel` | Sonnet 4.6 [gh] | External | `agents/sentinel/workspace/` | `#sentinel-chat` | DM only |

## Domain Agents (default Sonnet 4.6 — task specialists)

| Agent | ID | Pillar | Specialty | Skills |
|---|---|---|---|---|
| Spartan | `spartan` | HP 28 | Body & fitness | `spartan` |
| Oracle | `oracle` | HP 25-34 | Mind & money | `oracle` |
| Genesis | `genesis` | HP 25-26 | Identity & spirit | `identity` |
| Pulse | `pulse` | HP 27 | Feed & intel | `pulse`, `librarian` |
| Librarian | `librarian` | HP 27 | Research & knowledge | `librarian` |
| Haven | `haven` | LE 50-74 | Family & home | `haven` |
| Helix | `helix` | LE 52-56 | Generational | `helix` |
| Eden | `eden` | LE sustainability | Environment | `eden` |
| Nomad | `nomad` | LE 50 | Sovereign living | `nomad`, `homestead` |
| Empire | `empire` | GE 75-99 | Business strategy | `empire`, `empire-ops` |
| Apex | `apex` | GE 79 | Capital & finance | `apex`, `finance` |
| Creed | `creed` | GE 75-81 | Brand & story | `creed` |
| Trinity | `trinity` | SO 06 | R&D & research | `trinity`, `librarian` |
| Guardian | `guardian` | SO 04 | Security | `guardian`, `guardian-ops` |
| Grid | `grid` | SO 02 | Infrastructure | `grid`, `coding-agent`, `github`, `diffs` |
| Overwatch | `overwatch` | SO 07 | Automation | `overwatch` |
| EagleEye | `eagleeye` | SO 02 | Physical automation | `eagleeye`, `browser-automation` |

## Base Skills (all agents)

Every agent loads: `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync`

Super-agents additionally load: `taskflow-inbox-triage`

## Agent Detail Sheets

- [legacy.md](legacy.md) · [zenith.md](zenith.md) · [unity.md](unity.md) · [eternal.md](eternal.md) · [sovereignty.md](sovereignty.md) · [sentinel.md](sentinel.md)
- Domain agents: [spartan.md](spartan.md) · [oracle.md](oracle.md) · [genesis.md](genesis.md) · [pulse.md](pulse.md) · [librarian.md](librarian.md) · [haven.md](haven.md) · [helix.md](helix.md) · [eden.md](eden.md) · [nomad.md](nomad.md) · [empire.md](empire.md) · [apex.md](apex.md) · [creed.md](creed.md) · [trinity.md](trinity.md) · [guardian.md](guardian.md) · [grid.md](grid.md) · [overwatch.md](overwatch.md) · [eagleeye.md](eagleeye.md)

## Related

- [[../workspaces (Agent Workspaces)/README|Workspaces]] — Workspace file mirrors
- [[../plugin-skills (Custom Skills)/PLUGIN-SKILLS|Plugin Skills]] — Skill catalog
- [[../cron (Scheduled Jobs)/CRON|Cron]] — Scheduled jobs per agent
