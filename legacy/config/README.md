# Agent Configuration Templates

Shared configuration defaults for all servant agents.
Each agent profile extends these with domain-specific settings.

## Base Config

```yaml
model:
  provider: bedrock
  default: us.anthropic.claude-haiku-4-5-20251001-v1:0  # cheap, polling tasks
  context_length: 131072

agent:
  max_turns: 40
  reasoning_effort: medium
  verify_on_stop: false

platforms:
  nats:
    enabled: true
    servers:
      - nats://localhost:4222
    agent: hermes
    owner: zeus
    session_name: <agent-name>
    heartbeat_interval_s: 30

memory:
  memory_enabled: true
  user_profile_enabled: true
  memory_char_limit: 2000
  user_char_limit: 1375

delegation:
  max_iterations: 20
  max_concurrent_children: 1
  max_spawn_depth: 1
```

## Domain Upgrades
```
Oracle:     DeepSeek for fast data processing
Guardian:   Sonnet for threat analysis
Grid:       Haiku for routine health checks
Spartan:    Haiku for daily metrics
Haven:      Haiku for home automation checks
Heritage:   Sonnet for document analysis
Empire:     Sonnet for business strategy
Eden:       Haiku for sensor polling
Trinity:    Sonnet for complex integrations
```
