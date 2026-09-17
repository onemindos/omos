# _system/

Operational code for the OneMind Codex. This is the **engine room** — live configs,
deploy scripts, agent workspaces, and development tools.

Each subdirectory corresponds to a knowledge domain in the codex:

| Engine Module | Codex Domain | What's Here |
|---|---|---|
| `infrastructure/` | 01 Infrastructure | VPS, NATS, Docker, deploy scripts, service configs |
| `openclaw/` | 03 Agents | OpenClaw agent workspaces, deploy configs, skill definitions |
| `skill-templates/` | 03 Agents | Reusable skill templates for agent development |
| `legacy-files/` | 03 Agents | Archived legacy agent files |
| `development/` | 10 Development | Migration scripts, cleanup tools |

## Convention

- Knowledge (markdown docs, SOPs, research) lives in the numbered domains
- Code (scripts, configs, workspaces) lives here in `_system/`
- Each domain's `_ops/README.md` links back here for discoverability
- Agents working in a domain should check both the domain AND its `_system/` counterpart

## For Forks

If you're forking the codex as a personal framework, you can safely delete `_system/`
entirely — it contains OneMind-specific operational code, not the framework itself.
