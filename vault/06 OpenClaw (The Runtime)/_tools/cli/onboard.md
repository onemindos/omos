---
name: OpenClaw CLI — Onboard
updated: 2026-05-11
---

# `openclaw onboard`

Full guided onboarding for local or remote Gateway setup. Walks through model auth, workspace, gateway, channels, skills, and health in one flow.

**Official docs:** [docs.openclaw.ai/cli/onboard](https://docs.openclaw.ai/cli/onboard.md)

---

## Commands

```bash
openclaw onboard                               # Full interactive onboarding
openclaw onboard --modern                      # Crestodian conversational onboarding preview
openclaw onboard --flow quickstart             # Minimal prompts, auto-generates gateway token
openclaw onboard --flow manual                 # Full prompts for port, bind, and auth
openclaw onboard --flow import                 # Migrate from detected Hermes/Codex state
openclaw onboard --import-from hermes --import-source ~/.hermes
openclaw onboard --skip-bootstrap              # Skip creating workspace bootstrap files
openclaw onboard --mode remote --remote-url wss://gateway-host:18789
```

---

## Non-Interactive Mode

```bash
# GitHub Copilot (our setup)
openclaw onboard --non-interactive \
  --auth-choice github-copilot \
  --accept-risk

# Custom provider
openclaw onboard --non-interactive \
  --auth-choice custom-api-key \
  --custom-base-url "https://llm.example.com/v1" \
  --custom-model-id "foo-large" \
  --custom-api-key "$CUSTOM_API_KEY" \
  --secret-input-mode plaintext \
  --custom-compatibility openai

# Ref mode (stores env-backed refs instead of plaintext keys)
openclaw onboard --non-interactive \
  --auth-choice openai-api-key \
  --secret-input-mode ref \
  --accept-risk

# Gateway token via env ref
export OPENCLAW_GATEWAY_TOKEN="your-token"
openclaw onboard --non-interactive \
  --mode local \
  --auth-choice skip \
  --gateway-auth token \
  --gateway-token-ref-env OPENCLAW_GATEWAY_TOKEN \
  --accept-risk
```

---

## Flow Notes

| Flow | Behavior |
|---|---|
| `quickstart` | Minimal prompts, auto-generates a gateway token |
| `manual` | Full prompts for port, bind, and auth |
| `import` | Runs detected migration provider — requires fresh setup |

---

## Notes

- `--flow import` only runs against a fresh setup; use `openclaw migrate` for dry-run plans outside onboarding
- `--modern` starts Crestodian conversational onboarding preview
- `--skip-bootstrap` sets `agents.defaults.skipBootstrap: true` — skips creating `AGENTS.md`, `SOUL.md`, `TOOLS.md`, `IDENTITY.md`, `USER.md`, `HEARTBEAT.md`, `BOOTSTRAP.md`
- `--json` does NOT imply non-interactive mode — use `--non-interactive` for scripts

---

## Common Follow-Up

```bash
openclaw channels add
openclaw configure
openclaw agents add <name>
```

Use `openclaw setup` for baseline config only. Use `openclaw configure` for targeted changes later.
