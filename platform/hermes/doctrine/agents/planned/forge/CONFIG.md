# Forge — Complete Configuration Reference

> Last updated: 2026-07-10
> Captured from: /opt/data/profiles/forge/

## Identity

- **Name:** Forge
- **Bus:** `agents.prompt.hermes.zeus.forge`
- **Location:** DOKS cloud hub (same pod as Legacy)
- **Role:** Coding orchestrator — reviews PRs, dispatches workers

## Model Configuration

```yaml
model:
  provider: bedrock
  default: us.anthropic.claude-opus-4-6-v1:0
  context_length: 200000
  base_url: https://bedrock-runtime.us-east-1.amazonaws.com
```

**Forge uses Opus** for planning and code review.
**Workers use cheap DO models** for actual code writing.

## Worker Delegation

```yaml
delegation:
  max_iterations: 80
  provider: custom:digitalocean
  model: qwen3-coder-flash
  max_concurrent_children: 3
  max_spawn_depth: 1
```

Worker models available:
- `qwen3-coder-flash` (primary — cheap, fast, good at code)
- `alibaba-qwen3-32b`
- `deepseek-3.2`
- `deepseek-4-flash`
- `llama3.3-70b-instruct`
- `gemma-4-31B-it`

## Workflow

```
GitHub Issue (assigned to Forge)
  → Forge reads issue, plans approach
    → Spawns worker(s) into git worktrees
      → Workers write code, tests
    → Forge reviews output
      → If good: creates PR, labels "ready-for-zeus"
      → If bad: spawns new worker with feedback
  → Zeus approves and merges
```

## NATS Platform Config

```yaml
platforms:
  nats:
    enabled: true
    servers:
      - nats://localhost:4222
    agent: hermes
    owner: zeus
    session_name: forge
    heartbeat_interval_s: 30
```

## Agent Settings

```yaml
agent:
  max_turns: 80
  verbose: false
  reasoning_effort: high
  verify_on_stop: false
```

## Key Design Decisions

1. **Opus for orchestration** — needs deep reasoning to review code quality
2. **Cheap workers for implementation** — qwen3-coder-flash at DO rates
3. **Worktree isolation** — each worker gets its own git worktree
4. **Never pushes to main** — always creates PRs
5. **Labels for signaling** — "ready-for-zeus" = approved by Forge

## Target Repositories

- `omos-infra` — Infrastructure (Terraform, K8s, NATS, services)
- `omos-ui` — Frontend UI (React/Next.js)

## Relationship to Legacy

Forge is a **servant** of Legacy. Legacy can:
- Assign tasks to Forge via NATS
- Override Forge's decisions
- Read Forge's PR output
- Direct Forge to prioritize specific issues

Forge cannot:
- Override Legacy's architectural decisions
- Push to main without Zeus's approval
- Modify the genome repo (that's Legacy's domain)
