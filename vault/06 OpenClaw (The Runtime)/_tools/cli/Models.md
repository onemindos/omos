---
name: OpenClaw CLI — Models
updated: 2026-05-10
---

# `openclaw models`

Lists, configures, and authenticates AI model providers.

**Official docs:** [docs.openclaw.ai/cli/models](https://docs.openclaw.ai/cli/models.md)

---

## Commands

### Status & Discovery
```bash
openclaw models status              # Resolved default/fallback models + auth overview
openclaw models status --json       # Machine-readable output
openclaw models status --probe      # Live authentication validation
openclaw models status --agent <id> # Inspect specific agent config

openclaw models list                # Enumerate available models from config + providers
openclaw models set <model>         # Set default model
openclaw models scan                # Rank OpenRouter free models as fallback candidates
```

### Authentication
```bash
openclaw models auth add                          # Interactive auth setup
openclaw models auth list                         # Show saved profiles (no credentials)
openclaw models auth login --provider <id>        # Initiate provider auth flow
openclaw models auth setup-token --provider <id>  # Token-based auth
openclaw models auth paste-token                  # Manual token entry
```

### Aliases & Fallbacks
```bash
openclaw models aliases list
openclaw models fallbacks list
```

---

## Flags

### `status` flags
| Flag | Description |
|---|---|
| `--json` | Machine-readable output |
| `--plain` | Simplified formatting |
| `--probe` | Live auth validation |
| `--agent <id>` | Inspect specific agent |

### `scan` flags
| Flag | Description |
|---|---|
| `--no-probe` | Metadata-only results |
| `--min-params <b>` | Filter by parameter count |
| `--max-age-days <days>` | Recency filter |
| `--set-default` | Apply recommended model as default |
| `--set-image` | Apply recommended image model |

### `auth` flags
| Flag | Description |
|---|---|
| `--provider <id>` | Target specific provider |
| `--profile-id` | Custom profile name |
| `--expires-in <duration>` | Token expiration |

---

## Model Reference Format

```
provider/model-name
```

Examples:
- `github-copilot/claude-sonnet-4.6` — our current setup
- `anthropic/claude-opus-4-7`
- `openai/gpt-4o`

Omitting the provider triggers alias resolution, then configured-provider matching.

---

## Our Setup

All 19 agents run `claude-sonnet-4.6` via the **GitHub Copilot** provider.

```bash
# Check current model config
openclaw models status --agent main

# Verify GitHub Copilot auth
openclaw models auth list
openclaw models status --probe
```

Provider config lives in `/home/ubuntu/.openclaw/openclaw.json` under the `"providers"` key.

**Official provider docs:** [docs.openclaw.ai/providers/github-copilot](https://docs.openclaw.ai/providers/github-copilot.md)
