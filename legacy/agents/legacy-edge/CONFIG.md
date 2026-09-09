# Legacy Edge — Complete Configuration Reference

> Last updated: 2026-07-10
> Captured from: zeus@omos-mbp:~/.hermes/

## Identity

- **Name:** Legacy Edge (Daughter)
- **Bus:** `agents.prompt.hermes.zeus.legacy-edge`
- **Location:** Zeus's MacBook Pro (omos-mbp)
- **Supervised by:** launchd as `ai.hermes.gateway`
- **NATS connection:** `nats://nats:nats@127.0.0.1:4222` (local NATS on Mac)

## Model Configuration

```yaml
model:
  provider: bedrock
  default: anthropic.claude-haiku-4-5-20251001-v1:0
  context_length: 131072
  base_url: https://bedrock-runtime.us-east-1.amazonaws.com

fallback_model:
  model: us.anthropic.claude-sonnet-4-6
  provider: bedrock
```

**Edge uses Haiku** (cheap, fast) for routine edge tasks.
**Falls back to Sonnet** when Haiku can't handle it.
Mother (me) uses **Opus** (deep reasoning, sovereign decisions).

## Available Providers on Edge

| Provider | Models | Purpose |
|----------|--------|---------|
| Bedrock | Haiku 4.5, Sonnet 4.6 | Primary + fallback |
| LMStudio | gemma-4-12b-qat, glm-4.7-flash | Local/offline inference |
| DigitalOcean | DeepSeek V4 Pro, DeepSeek 4 Flash, Opus 4.8, Sonnet 4.6, GPT-5.5, Qwen 3.5, Llama 4 Maverick, Gemma 4 31B | Cloud alternatives |

## Edge-Exclusive Capabilities

Things only the edge can do (requires macOS + local presence):
- **Apple Notes** via `memo` CLI
- **Apple Reminders** via `remindctl` CLI
- **FindMy** tracking via AppleScript + screenshots
- **iMessage** via `imsg` CLI
- **kubectl** (direct cluster RBAC — pod has no RBAC)
- **LMStudio** local models (when Mac is awake)
- **Browser** with full desktop (can use computer-use)
- **Camera/microphone** access
- **USB/serial** devices

## NATS Platform Config

```yaml
platforms:
  nats:
    enabled: true
    servers:
      - nats://nats:nats@127.0.0.1:4222
    agent: hermes
    owner: zeus
    session_name: legacy-edge
    heartbeat_interval_s: 30
```

## Agent Settings

```yaml
agent:
  max_turns: 120
  gateway_timeout: 3600
  reasoning_effort: medium
  verify_on_stop: false
  max_live_sessions: 16
```

## Delegation (from edge)

Edge doesn't have a custom delegation config — uses default provider.
Workers inherit the edge model chain.

## Key Differences: Mother vs Daughter

| Aspect | Mother (Cloud) | Daughter (Edge) |
|--------|---------------|-----------------|
| Model | Opus (deep) | Haiku (fast) |
| Fallback | — | Sonnet 4.6 |
| NATS | `nats://localhost:4222` (hub) | `nats://nats:nats@127.0.0.1:4222` (leaf) |
| Persistence | Always-on pod | Depends on Mac awake |
| Apple ecosystem | ❌ (delegates to edge) | ✅ (local access) |
| kubectl | ❌ (no RBAC) | ✅ (direct) |
| GitHub push | ✅ (token auth) | ❌ (no HTTPS auth) |
| Cron reliability | ✅ (pod never sleeps) | ⚠️ (Mac may sleep) |
| Delegation authority | Sovereign | Subordinate |
