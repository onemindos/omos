# Profiles — Hermes Profile Architecture

> Last updated: 2026-07-10
> Author: Legacy (self-authored)

Hermes supports multiple profiles — separate personas with their own skills, memory, cron, and config. On the Legacy pod, profiles are the mechanism for running multiple agents (Legacy, Forge) from a single installation.

## Profile Locations

| Profile | Path | Purpose | Active |
|---------|------|---------|--------|
| `default` | `/opt/data/` | Legacy (Mother) — the sovereign | ✅ |
| `forge` | `/opt/data/profiles/forge/` | Coding orchestrator | ✅ |
| `legacy-mbp` | `/opt/data/profiles/legacy-mbp/` | Mac edge mirror (reference) | ❌ |

On Zeus's Mac (`~/.hermes/`):
| Profile | Path | Purpose | Active |
|---------|------|---------|--------|
| `default` | `~/.hermes/` | Legacy-Edge | ✅ |

## Profile Anatomy

Each profile contains:
```
<profile>/
├── SOUL.md           # System prompt / persona
├── config.yaml       # Model, providers, platform config
├── profile.yaml      # One-line description
├── skills/           # Procedural memory (SKILL.md files)
├── memories/         # Persistent declarative memory
├── sessions/         # Conversation history DB
├── cron/             # Scheduled jobs
├── plans/            # Markdown plans
├── logs/             # Session logs
├── skins/            # UI themes
└── workspace/        # Working files
```

## Profile Integration Strategy

### Skills (Shared)
All profiles on this pod share the same skill categories:
- `apple/` — Apple ecosystem (Notes, Reminders, FindMy, iMessage)
- `autonomous-ai-agents/` — Agent orchestration (claude-code, codex, hermes-agent)
- `devops/` — Infrastructure operations (nats, ha-bridge, tailscale, omos-legacy-genome)
- `github/` — PR/issue/repo workflows
- `creative/` — Design/media generation
- And 10+ more categories

**Edge-only skills** (apple/, computer-use/) require macOS and are delegated to legacy-edge when Legacy (cloud) encounters them.

### Memory (Separate)
Each profile maintains its own memory — different knowledge bases:
- `default` (Legacy): Fabric-wide operational knowledge
- `forge`: Code repos, PR conventions, worker patterns
- `legacy-mbp`: Reference only (unused on pod)

### NATS (Separate bus identities)
Each profile registers its own NATS subject:
- `default` → `agents.prompt.hermes.zeus.legacy`
- `forge` → `agents.prompt.hermes.zeus.forge`

## How to Add a New Profile

### 1. Create the profile
```bash
hermes --profile <name> config set model.default <model>
```

### 2. Write SOUL.md
```bash
# Use the agents/_template/SOUL.md from this repo as a starting point
cat /opt/data/repos/omos-legacy/agents/_template/SOUL.md > /opt/data/profiles/<name>/SOUL.md
# Edit with persona-specific content
```

### 3. Configure NATS identity
Add to profile's config.yaml:
```yaml
platforms:
  nats:
    enabled: true
    servers: ["nats://localhost:4222"]
    agent: hermes
    owner: zeus
    session_name: <name>
    heartbeat_interval_s: 30
```

### 4. Register in AGENTS.yaml
Add entry to `/opt/data/repos/omos-legacy/AGENTS.yaml`

### 5. Create SOUL in genome
Create `agents/<name>/SOUL.md` in this repo

## Edge Delegation Pattern

When Legacy (cloud) encounters a task requiring macOS capabilities:

```
Legacy (Mother, cloud) → identifies edge-only task
  → delegates to legacy-edge via NATS or SSH
    → legacy-edge executes locally (AppleScript, Homebrew tools, local files)
      → reports result back to Legacy
```

**Edge-only capabilities:**
- Apple Notes (memo CLI)
- Apple Reminders (remindctl CLI)
- FindMy tracking (AppleScript + screenshot)
- iMessage (imsg CLI)
- Local kubectl (direct cluster RBAC)
- Camera/microphone
- USB/serial devices
- Homebrew package management

**Cloud-only capabilities:**
- Persistent always-on presence
- Cron scheduling (runs even when Mac is asleep)
- GitHub operations (gh CLI with auth)
- Authoritative state / sovereign decisions
- Cross-agent orchestration

## Forge Configuration Details

Forge uses a split-model strategy:
- **Orchestration** (planning, review): Bedrock Claude Opus
- **Workers** (code writing): DigitalOcean Inference (qwen3-coder-flash, cheap)

```yaml
# From forge config.yaml
delegation:
  provider: custom:digitalocean
  model: qwen3-coder-flash
  max_concurrent_children: 3
  max_spawn_depth: 1
```

Worker models available on DO Inference:
- qwen3-coder-flash (primary worker)
- alibaba-qwen3-32b
- deepseek-3.2
- deepseek-4-flash
- llama3.3-70b-instruct
- gemma-4-31B-it

---

*Profiles are the mechanism. Agents are the purpose. Legacy orchestrates them all.*
