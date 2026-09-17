---
name: OpenClaw — CLI Reference
updated: 2026-05-11
---

# CLI Reference

`openclaw` is the primary CLI for managing the gateway, agents, channels, models, cron jobs, skills, memory, and everything else.

**Official docs:** [docs.openclaw.ai/cli](https://docs.openclaw.ai/cli/index.md)

## Global Flags

| Flag | Description |
|---|---|
| `--dev` | Developer mode |
| `--profile` | Config profile |
| `--container` | Container mode |
| `--no-color` | Disable color output |
| `--update` | Check for updates |
| `-V / --version` | Print version |

## Command Index

### Setup & Onboarding
| Command | What it does | Docs |
|---|---|---|
| `openclaw onboard` | Full first-run guided setup | — |
| `openclaw setup` | Baseline config without guided flow | — |
| `openclaw configure` | Modify existing setup | [configure.md](configure.md) |
| `openclaw doctor` | Diagnose and fix config issues | [doctor.md](doctor.md) |
| `openclaw config` | Non-interactive config get/set/patch | [config.md](config.md) |
| `openclaw migrate` | Import from Claude / Codex / Hermes | [migrate.md](migrate.md) |
| `openclaw update` | Update OpenClaw | [update.md](update.md) |
| `openclaw backup` | Backup configuration | [backup.md](backup.md) |

### Agents & Messaging
| Command | What it does | Docs |
|---|---|---|
| `openclaw agents` | Manage agent workspaces and routing | [agents.md](agents.md) |
| `openclaw agent` | Execute a single agent turn | [agent.md](agent.md) |
| `openclaw message` | Send/edit/delete messages, manage threads | [message.md](message.md) |
| `openclaw acp` | ACP bridge for IDE integration | [acp.md](acp.md) |
| `openclaw mcp` | MCP server mode and config | [mcp.md](00-24%20SO%20(Sovereign%20Operations)/01%20Intelligence%20&%20Agents%20(The%20Brain)/fleet/openclaw/_tools/cli/mcp.md) |
| `openclaw path` | `oc://` workspace file addressing | [path.md](path.md) |

### Health & Monitoring
| Command | What it does | Docs |
|---|---|---|
| `openclaw status` | Full system diagnostics | [status.md](status.md) |
| `openclaw health` | Gateway health check | [health.md](health.md) |
| `openclaw logs` | Tail gateway logs | [logs.md](logs.md) |
| `openclaw sessions` | Session listing and export | [sessions.md](sessions.md) |
| `openclaw system` | System events, heartbeat, presence | [system.md](system.md) |

### Models & Inference
| Command | What it does | Docs |
|---|---|---|
| `openclaw models` | List, configure, auth providers | [Models.md](Models.md) |
| `openclaw infer` | Run inference, generate images/audio/TTS/embeddings | [infer.md](infer.md) |

### Automation
| Command | What it does | Docs |
|---|---|---|
| `openclaw cron` | Scheduled jobs | [Cron.md](Cron.md) |
| `openclaw hooks` | Event-driven hooks | [hooks.md](hooks.md) |
| `openclaw tasks` | Background tasks | [tasks.md](tasks.md) |
| `openclaw webhooks` | Webhook management (Gmail Pub/Sub) | [webhooks.md](webhooks.md) |
| `openclaw approvals` | Exec approval policies | [approvals.md](approvals.md) |

### Knowledge & Skills
| Command | What it does | Docs |
|---|---|---|
| `openclaw skills` | Search, install, manage skills | [skills.md](skills.md) |
| `openclaw memory` | Semantic memory index, search, consolidation | [memory.md](memory.md) |
| `openclaw commitments` | Inferred follow-up memories | [commitments.md](commitments.md) |
| `openclaw wiki` | Memory-wiki vault management | [wiki.md](wiki.md) |

### Plugins & Extensions
| Command | What it does | Docs |
|---|---|---|
| `openclaw plugins` | Install/manage plugins from ClawHub/npm/git/local | [plugins.md](plugins.md) |
| `openclaw clawbot` | Legacy namespace (alias layer) | [clawbot.md](clawbot.md) |

### Security & Secrets
| Command | What it does | Docs |
|---|---|---|
| `openclaw secrets` | SecretRef management, audit, apply, reload | [secrets.md](secrets.md) |
| `openclaw security` | Security audit | [security.md](security.md) |

### Network & Devices
| Command | What it does | Docs |
|---|---|---|
| `openclaw nodes` | Paired remote node management | [nodes.md](nodes.md) |
| `openclaw node` | Deploy headless node agent | [node.md](00-24%20SO%20(Sovereign%20Operations)/01%20Intelligence%20&%20Agents%20(The%20Brain)/fleet/openclaw/_tools/cli/node.md) |
| `openclaw devices` | Paired device management | [devices.md](devices.md) |
| `openclaw channels` | Channel account management | [channels.md](channels.md) |
| `openclaw directory` | ID lookups across channels | [directory.md](directory.md) |
| `openclaw pairing` | DM pairing approval | [pairing.md](pairing.md) |
| `openclaw qr` | Mobile pairing QR code generation | [qr.md](qr.md) |
| `openclaw voicecall` | Phone call management | [voicecall.md](voicecall.md) |
| `openclaw dns` | CoreDNS / Tailscale discovery | [dns.md](dns.md) |
| `openclaw proxy` | Proxy validation and traffic debugging | [proxy.md](proxy.md) |

### Gateway & Runtime
| Command | What it does | Docs |
|---|---|---|
| `openclaw gateway` | Gateway start/stop/reload/service management | [gateway.md](gateway.md) |
| `openclaw sandbox` | Runtime sandbox management | [sandbox.md](sandbox.md) |
| `openclaw tui` | Terminal UI, local embedded mode | [tui.md](tui.md) |
| `openclaw browser` | Browser lifecycle, tabs, automation | [browser.md](browser.md) |
| `openclaw daemon` | Legacy alias for gateway commands | [daemon.md](daemon.md) |

---

## Our Quick-Reference

```bash
# Check system status
openclaw status

# Check gateway health
openclaw health --verbose

# List all agents
openclaw agents list

# List cron jobs
openclaw cron list

# Tail live logs
openclaw logs --follow

# List installed skills
openclaw skills list

# Check memory health
openclaw memory status

# Edit config non-interactively
openclaw config set agents.legacy.model "github-copilot/claude-sonnet-4-6"

# Read a workspace file
openclaw path read oc://legacy/IDENTITY.md

# Show pending device pairing requests
openclaw devices list --pending

# Show system events (last hour)
openclaw system events --since 1h
```
