---
name: OpenClaw CLI — Skills
updated: 2026-05-10
---

# `openclaw skills`

Manages local skills and integrates with ClawHub for discovery and installation.

**Official docs:** [docs.openclaw.ai/cli/skills](https://docs.openclaw.ai/cli/skills.md)

---

## Commands

```bash
openclaw skills list                           # List all local skills (default subcommand)
openclaw skills list --eligible                # Only skills eligible for current agent
openclaw skills list --verbose                 # Verbose output
openclaw skills list --agent <id>              # Skills visible to a specific agent

openclaw skills search "calendar"              # Search ClawHub
openclaw skills search                         # Browse ClawHub default feed

openclaw skills install <slug>                 # Install latest version
openclaw skills install <slug> --version <v>   # Install specific version
openclaw skills install <slug> --force         # Force overwrite existing

openclaw skills check --agent <id>             # Validate skill visibility to agent
```

---

## Flags

| Flag | Description |
|---|---|
| `--agent <id>` | Target specific agent workspace |
| `--json` | Machine-readable output |
| `--limit <n>` | Restrict search results |
| `--all` | Apply updates across all tracked installations |
| `--force` | Overwrite existing on install |

---

## Our Installed Skills

Skills are installed per-agent workspace. Legacy (main) has the richest set:

**Active (`workspaces/legacy/skills/`):**
`api-gateway`, `github-api`, `obsidian`, `onepassword`, `todoist`, `twilio`, `vision`, `weather`, `rss-digest`, `proactive-agent`, `skill-forge`, and more.

**Archived (`workspaces/legacy/skills-archive/`):**
`browser`, `credit-repair`, `desktop-control`, `uber`, `transit`, `zoho-mail`, and others.

All other agents share: `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync`

---

## ClawHub

Skills registry at [clawhub.openclaw.ai](https://docs.openclaw.ai/clawhub/index.md). See [`../clawhub/README.md`](../clawhub/README.md) for publishing notes.
