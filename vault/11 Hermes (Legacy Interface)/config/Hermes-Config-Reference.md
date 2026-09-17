---
title: Hermes Config Reference
domain: 11
pillar: SO
status: active
updated: 2026-05-21
tags: [hermes, config, reference]
---

# Hermes Config Reference

## Paths

| Path | Purpose |
|------|---------|
| `~/.hermes/` | Hermes home directory |
| `~/.hermes/config.yaml` | Main configuration |
| `~/.hermes/.env` | Secrets and env vars |
| `~/.hermes/skills/` | Custom skills |
| `~/.hermes/sessions/` | Session database |
| `~/.hermes/memory/` | Persistent memory |
| `~/.hermes/cron/jobs.json` | Scheduled tasks |

## Active Model Config (as of 2026-05-21)

```yaml
model:
  default: us.anthropic.claude-sonnet-4-6
  provider: bedrock

providers:
  bedrock:
    region: us-east-1
  azure:
    base_url: https://onemind-openai.openai.azure.com/openai/v1
```

## API Server Config

```yaml
api_server:
  enabled: true
  port: 8080
  key: <see config.yaml>
  cors_origins: "*"
```

Local endpoint:     http://localhost:8080/v1
Tailscale endpoint: http://omos-mbp.tail717ba0.ts.net:8080/v1
Tailscale IP:       http://100.98.178.23:8080/v1

## Dashboard Config

```yaml
# Dashboard runs at http://localhost:9119
# launchd plist uses: hermes dashboard --tui --no-open --host 127.0.0.1 --port 9119 --insecure --skip-build
# Without --tui, the Chat tab is disabled
```

## Useful CLI Commands

```bash
# Check config
hermes config list

# Set a config value
hermes config set <key> <value>

# Start API server
hermes serve

# Start dashboard
hermes dashboard --tui --no-open --host 127.0.0.1 --port 9119

# List tools
hermes tools

# Run a task
hermes run "your task here"
```

## Related

- [[Hermes-Overview]]
- [[Hermes-API-Server]]
