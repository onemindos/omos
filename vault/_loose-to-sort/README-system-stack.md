# Security Validation Gate

> Pre-execution security layer. Agents check commands against these patterns before running anything destructive.

## How It Works

1. Before executing a command, agent checks `patterns.yaml`
2. **BLOCKED** → Refuse immediately, log event
3. **CONFIRM** → Ask Zeus before proceeding, log event
4. **ALERT** → Log and proceed with caution
5. **No match** → Proceed normally

## Files

| File | Purpose |
|------|---------|
| `patterns.yaml` | Security rules — what to block, confirm, alert |
| `validator.md` | Step-by-step validation instructions for agents |
| `events/` | Security event logs (JSONL, one file per month) |
