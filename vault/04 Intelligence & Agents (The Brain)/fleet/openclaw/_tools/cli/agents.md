---
name: OpenClaw CLI — Agents
updated: 2026-05-10
---

# `openclaw agents`

Manages agent workspaces, routing bindings, and identity.

**Official docs:** [docs.openclaw.ai/cli/agents](https://docs.openclaw.ai/cli/agents.md)

---

## Commands

```bash
openclaw agents list                                         # List all agents
openclaw agents add <name> --workspace <dir>                 # Create a new agent
openclaw agents delete <id>                                  # Delete agent (moves files to Trash)
openclaw agents bindings                                     # View all routing bindings
openclaw agents bindings --json                              # Machine-readable bindings

# Bind agent to a channel
openclaw agents add work --workspace ~/.openclaw/workspace-work --bind telegram:ops

# Set agent identity
openclaw agents set-identity --agent main --avatar avatars/openclaw.png
```

---

## Routing Bindings

Bindings pin inbound channel traffic to specific agents. Three scopes:

| Scope | Example | Behavior |
|---|---|---|
| Channel-only | `telegram` | All Telegram traffic → this agent |
| Account-specific | `telegram:ops` | Only `ops` account traffic → this agent |
| Wildcard | `accountId: "*"` | All accounts on the channel → this agent |

Upgrade a channel-only binding to account-scoped by rebinding with an explicit account ID.

---

## Identity Config

Each agent can have: name, emoji, avatar, theme. Stored in `IDENTITY.md` or set via CLI:

```bash
openclaw agents set-identity --agent spartan --avatar avatars/spartan.png
```

---

## Our Agents (19)

| ID | Name | Pillar |
|---|---|---|
| `main` | Legacy | SO — Orchestrator |
| `spartan` | Spartan | HP |
| `oracle` | Oracle | HP |
| `genesis` | Genesis | HP |
| `pulse` | Pulse | HP |
| `librarian` | Librarian | HP |
| `eden` | Eden | LE |
| `haven` | Haven | LE |
| `helix` | Helix | LE |
| `nomad` | Nomad | LE |
| `empire` | Empire | GE |
| `apex` | Apex | GE |
| `creed` | Creed | GE |
| `trinity` | Trinity | SO |
| `guardian` | Guardian | SO |
| `grid` | Grid | SO |
| `overwatch` | Overwatch | SO |
| `eagleeye` | Eagleeye | SO |
| `sentinel` | Sentinel | SO |

> Note: `main` is reserved — cannot be deleted or reused as a new agent ID.
