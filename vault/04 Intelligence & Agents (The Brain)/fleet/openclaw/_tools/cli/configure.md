---
name: OpenClaw CLI — Configure
updated: 2026-05-11
---

# `openclaw configure`

Interactive setup wizard for modifying existing OpenClaw installations.

**Official docs:** [docs.openclaw.ai/cli/configure](https://docs.openclaw.ai/cli/configure.md)

---

## Commands

```bash
openclaw configure                    # Full interactive wizard
openclaw configure --section model    # Jump to model section
openclaw configure --section channels # Jump to channels section
openclaw configure --section web      # Jump to web search provider
openclaw configure --section plugins  # Jump to plugins
```

---

## Available Sections

| Section | What it configures |
|---|---|
| `workspace` | Workspace directory paths |
| `model` | AI model selection + allowlists |
| `web` | Web search provider + credentials |
| `gateway` | Gateway settings |
| `daemon` | Background service settings |
| `channels` | Channel account setup |
| `plugins` | Plugin settings |
| `skills` | Skill configuration |
| `health` | Health check settings |

---

## vs Other Setup Commands

| Command | Use when |
|---|---|
| `openclaw onboard` | First-run, full guided journey |
| `openclaw configure` | Modifying an already-configured system |
| `openclaw setup` | Baseline config without guided flow |
| `openclaw config get/set` | Non-interactive single-value edits |
| `openclaw channels add` | Channel account setup only |

---

## Note on Model Re-auth

When re-running provider auth, the system **preserves your existing `agents.defaults.model.primary`** even if the provider's auth step returns its own recommended default model.
