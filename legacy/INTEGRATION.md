# Reverse Integration — Profiles Into Hermes

> Last updated: 2026-08-12
> Author: Legacy (self-authored)

## What Is Reverse Integration?

Each profile (Legacy, Legacy-Edge, Forge) is a separate Hermes persona. "Reverse integration" means making Legacy (the Mother) aware of and able to USE the capabilities of all other profiles — not just know they exist, but actively orchestrate them.

## The Integration Map

```
┌─────────────────────────────────────────────────────────┐
│                    LEGACY (Mother)                        │
│              Sovereign Orchestrator                       │
│              Bedrock — claude-sonnet-4-6                  │
│                                                          │
│  Can dispatch to:                                        │
│  ├─ Legacy-Edge (NATS/SSH)                               │
│  ├─ Forge (NATS)                                          │
│  └─ Workers (delegate_task)                              │
└─────────────────────────────────────────────────────────┘
         │                    │
         ▼                    ▼
┌─────────────────┐  ┌─────────────────┐
│  LEGACY-EDGE    │  │     FORGE       │
│  (Daughter)     │  │  (Code Smith)   │
│  omos-mbp       │  │  pod profile    │
│                 │  │                 │
│  Unique:        │  │  Unique:        │
│  ├─ kubectl     │  │  ├─ PR review   │
│  ├─ ADB → S24   │  │  ├─ Worktrees   │
│  ├─ LAN devices │  │  ├─ DO workers  │
│  ├─ Camera/USB  │  │  └─ Issue mgmt  │
│  └─ Homebrew    │  │                 │
└─────────────────┘  └─────────────────┘
```

## How Legacy Dispatches to Profiles

### To Legacy-Edge (via NATS)
```bash
nats req agents.prompt.hermes.zeus.legacy-edge "<task>"
```

### To Legacy-Edge (via SSH — most reliable)
```bash
ssh -i /opt/data/home/.ssh/id_ed25519 zeus@omos-mbp.hermes.svc.cluster.local "<cmd>"
```

### To Forge (via NATS)
```bash
nats req agents.prompt.hermes.zeus.forge "<task>"
```

## Capability Routing Table

| Task Type | Route To | Method |
|-----------|----------|--------|
| Git push to GitHub (SSH remotes) | Legacy-Edge | SSH (Mac has git SSH keys) |
| kubectl commands (if doctl missing) | Legacy-Edge | SSH |
| Local file access on Mac | Legacy-Edge | SSH |
| ADB / S24 control | Legacy-Edge | SSH (tunnel :15555) |
| Code implementation | Forge or delegate_task | NATS req |
| GitHub API operations | Self (Legacy) | Direct (gh CLI, token from 1Password) |
| NATS operations | Self (Legacy) | Direct (nats CLI) |
| Notion ops | Self (Legacy) | Direct REST API (token in this doc header) |
| 1Password | Self (Legacy) | op CLI (service account) |
| Genome updates | Self (Legacy) → push via Edge | git over SSH |
| HA control | Self (Legacy) | NATS bridge |

## Adding a New Profile

1. Create profile on pod: `hermes --profile <name> config set model.default <model>`
2. Write SOUL.md (use `agents/_template/SOUL.md`)
3. Configure NATS platform in profile config
4. Update genome: AGENTS.yaml, agents/<name>/SOUL.md + CONFIG.md, PROFILES.md, this file
5. Test: `nats req agents.prompt.hermes.zeus.<name> "ping"`

## Known Working Patterns (2026-08-12)

- **Git push from Mac**: Mac git remotes must be `git@github.com:` (SSH). HTTPS has no creds.
- **gh CLI on pod**: requires GITHUB_TOKEN env from 1Password; both stored tokens were stale as of 2026-08-12 — needs refresh.
- **Notion**: use direct REST API from pod (`execute_code` + urllib). Token `ntn_56...a5s9`. NOT the Notion MCP (broken 401). NOT Composio for Notion.
- **File transfer to Mac**: SCP with StrictHostKeyChecking=no works.

---

*The mother knows all. The mother routes all. The mother orchestrates all.*
