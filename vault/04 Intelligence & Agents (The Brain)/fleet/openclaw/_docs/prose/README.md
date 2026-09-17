---
name: OpenClaw — OpenProse
updated: 2026-05-11
---

# OpenProse

OpenProse is a portable, markdown-first workflow format for orchestrating AI sessions. In OpenClaw it ships as a plugin that installs an OpenProse skill pack plus a `/prose` slash command. Programs live in `.prose` files and can spawn multiple sub-agents with explicit control flow.

**Official docs:** [docs.openclaw.ai/prose](https://docs.openclaw.ai/prose.md)  
**Official site:** [prose.md](https://www.prose.md)

> For deterministic, approval-gated workflows compare with [Lobster](../tools/lobster.md). OpenProse handles multi-agent orchestration and research/synthesis; Lobster handles deterministic step pipelines with approval gates. They pair well together.

---

## What It Can Do

- Multi-agent research + synthesis with explicit parallelism
- Repeatable approval-safe workflows (code review, incident triage, content pipelines)
- Reusable `.prose` programs that run across supported agent runtimes

---

## Install + Enable

Bundled but disabled by default:

```bash
openclaw plugins enable open-prose
# Restart the Gateway after enabling
openclaw gateway restart
```

---

## Slash Command

```bash
/prose help
/prose run <file>         # Run a .prose file
/prose run <handle>       # Run from prose.md registry
/prose run <url>          # Run from direct URL
/prose compile <file>     # Compile without running
/prose examples           # Show example programs
/prose update             # Update the OpenProse plugin
```

---

## Example: Parallel Research + Synthesis

```prose
# Research + synthesis with two agents running in parallel.
input topic: "What should we research?"

agent researcher:
  model: sonnet
  prompt: "You research thoroughly and cite sources."

agent writer:
  model: opus
  prompt: "You write a concise summary."

parallel:
  findings = session: researcher
    prompt: "Research {topic}."
  draft = session: writer
    prompt: "Summarize {topic}."

session "Merge the findings + draft into a final answer."
  context: { findings, draft }
```

---

## File Locations

```
.prose/
├── .env
├── runs/
│   └── {YYYYMMDD}-{HHMMSS}-{random}/
│       ├── program.prose
│       ├── state.md
│       ├── bindings/
│       └── agents/
└── agents/

~/.prose/agents/         ← user-level persistent agents
```

---

## State Modes

| Mode | Notes |
|---|---|
| `filesystem` (default) | `.prose/runs/...` |
| `in-context` | Transient, for small programs |
| `sqlite` (experimental) | Requires `sqlite3` binary |
| `postgres` (experimental) | Requires `psql` + connection string |

---

## OpenClaw Runtime Mapping

| OpenProse concept | OpenClaw tool |
|---|---|
| Spawn session / Task tool | `sessions_spawn` |
| File read/write | `read` / `write` |
| Web fetch | `web_fetch` |

If your tool allowlist blocks these tools, OpenProse programs will fail.

---

## Security

Treat `.prose` files like code. Review before running. Use OpenClaw tool allowlists and approval gates to control side effects.

If a Prose program needs Lobster, allow the `lobster` tool for sub-agents via `tools.subagents.tools`.

---

## One Mind Use Case

OpenProse excels at the intelligence/research layer — spawning Librarian + Pulse in parallel to research a topic, then having Legacy synthesize. Pairs directly with Lobster for the delivery/approval step:

1. `/prose run research-brief.prose` — parallel agents gather and synthesize
2. Lobster pipeline — deterministic delivery with approval gate before posting to Discord

---

## Related

- [tools/lobster.md](../tools/lobster.md)
- [automation/README.md](../automation/README.md)
- [tools/subagents](../tools/README.md)
