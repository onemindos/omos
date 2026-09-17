---
name: OpenClaw CLI — Plugins
updated: 2026-05-11
---

# `openclaw plugins`

Installs, manages, and diagnoses OpenClaw plugins from ClawHub, npm, git, or local paths.

**Official docs:** [docs.openclaw.ai/cli/plugins](https://docs.openclaw.ai/cli/plugins.md)

---

## Commands

```bash
openclaw plugins search <query>           # Search ClawHub for installable packages
openclaw plugins list                     # List installed plugins
openclaw plugins list --enabled           # Only enabled plugins
openclaw plugins list --verbose           # Verbose output
openclaw plugins list --json              # Machine-readable

openclaw plugins install <spec>           # Install plugin
openclaw plugins install <spec> --pin     # Lock exact npm version
openclaw plugins install <spec> --force   # Overwrite existing from different source

openclaw plugins inspect <id>             # Examine plugin details
openclaw plugins inspect <id> --runtime   # Include runtime capabilities

openclaw plugins enable <id>              # Enable plugin
openclaw plugins disable <id>             # Disable plugin

openclaw plugins update <id>              # Update plugin
openclaw plugins update --all             # Update all tracked plugins

openclaw plugins uninstall <id>           # Remove plugin

openclaw plugins doctor                   # Check load errors and incompatibilities
openclaw plugins registry                 # Inspect plugin registry
openclaw plugins registry --refresh       # Rebuild registry
```

---

## Installation Sources

| Source | Syntax | Example |
|---|---|---|
| ClawHub | `clawhub:<package>` | `clawhub:openclaw-codex-app-server@1.2.3` |
| npm (default) | `<package>` | `openclaw-codex-app-server` |
| npm Pack | `npm-pack:<path.tgz>` | Testing managed npm paths |
| Git | `git:<url>[@ref]` | `git:github.com/owner/repo@v1.0.0` |
| Local | `<path>` | `./my-plugin` or `.zip`/`.tgz` archives |

---

## Important Behaviors

**Version pinning:** `--pin` locks exact npm versions. Unversioned ClawHub installs follow newer releases automatically.

**Safety scanning:** Built-in dangerous-code detection. Use `--dangerously-force-unsafe-install` only as a break-glass override for verified false positives.

**Nix mode:** `OPENCLAW_NIX_MODE=1` disables lifecycle mutations — manage plugins through Nix source instead.

**Tracing slow operations:** `OPENCLAW_PLUGIN_LIFECYCLE_TRACE=1` shows phase timing details.

---

## Hook Packs

Hook packs are installed via `openclaw plugins install`, not `hooks install` (legacy name). After install, restart the gateway for changes to take effect.

See also: [`hooks.md`](hooks.md)
