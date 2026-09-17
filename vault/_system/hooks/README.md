# Lifecycle Hooks

> Event-driven automation mapped to CODEX workflow stages. Agents check this registry and execute hooks at the right moment.

## How It Works

1. Agent reads `hooks.yaml` at session start
2. At each lifecycle event, agent executes the registered hooks
3. Hooks are markdown instruction files — agents read and follow them
4. Each hook has a `blocking` flag (must complete before proceeding) and `fail_behavior`

## Events

| Event | When | Hooks |
|-------|------|-------|
| `session-start` | Agent begins interaction | LoadContext, LoadMemory |
| `capture` | New information enters | ClassifyInput |
| `organize` | Content being routed | RouteToDomain |
| `direct` | Prioritizing work | CheckGoals |
| `execute-pre` | Before running a command | SecurityValidate |
| `execute-post` | After action completes | LogAction |
| `multiply` | During reflection | ExtractLearning |
| `session-end` | Agent session ending | SaveMemory, CaptureSignal |

## Files

- `hooks.yaml` — Hook registry and configuration
- Individual hook `.md` files — Execution instructions
