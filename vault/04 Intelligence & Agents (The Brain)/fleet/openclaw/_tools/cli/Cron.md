---
name: OpenClaw CLI — Cron
updated: 2026-05-10
---

# `openclaw cron`

Manages scheduled jobs through the Gateway scheduler.

**Official docs:** [docs.openclaw.ai/cli/cron](https://docs.openclaw.ai/cli/cron.md)

---

## Commands

### List Jobs
```bash
openclaw cron list
openclaw cron list --agent <id>
openclaw cron list --json
```

### View Job Details
```bash
openclaw cron show <job-id>
openclaw cron show <job-id> --json
```

### Create a Job
```bash
openclaw cron add \
  --name "Job Name" \
  --cron "0 7 * * *" \
  --session isolated \
  --message "Your prompt here" \
  --announce \
  --channel telegram \
  --to "123456789"
```

### Edit a Job
```bash
openclaw cron edit <job-id> [flags]
```

### Run a Job Immediately
```bash
openclaw cron run <job-id>
openclaw cron run <job-id> --due
```

### View Run History
```bash
openclaw cron runs --id <job-id> --limit 50
```

---

## Session Options

| Value | Behavior |
|---|---|
| `main` | Binds to the agent's primary session |
| `isolated` | Creates a fresh transcript and session ID per run |
| `current` | Binds to the active session at creation time |
| `session:<id>` | Pins to an explicit persistent session key |

---

## Flags

### Scheduling
| Flag | Description |
|---|---|
| `--cron <expression>` | Recurring schedule (cron syntax) |
| `--at <datetime>` | One-shot execution (UTC default) |
| `--tz <iana>` | Timezone for wall-clock times (e.g. `America/New_York`) |
| `--keep-after-run` | Keep one-shot jobs after they execute |

### Delivery
| Flag | Description |
|---|---|
| `--announce` | Enable fallback delivery (default for `isolated`) |
| `--no-deliver` | Disable runner fallback delivery |
| `--channel <name>` | Target channel (`telegram`, `discord`, `slack`, etc.) |
| `--to <target>` | Recipient identifier (chat ID, user ID) |
| `--thread-id <id>` | Forum topic targeting (Telegram supergroups) |
| `--best-effort-deliver` | Delivery resilience |
| `--failure-alert-include-skipped` | Include skipped runs in failure alerts |

### Agent & Model
| Flag | Description |
|---|---|
| `--agent <id>` | Assign job to a specific agent |
| `--clear-agent` | Remove agent assignment |
| `--model <ref>` | Select a job-specific model |
| `--light-context` | Lightweight bootstrap for isolated jobs |

---

## Examples

### Morning brief (lightweight, no chat delivery)
```bash
openclaw cron add \
  --name "Morning Brief" \
  --cron "0 7 * * *" \
  --session isolated \
  --message "Summarize overnight updates." \
  --light-context \
  --no-deliver
```

### Telegram forum topic delivery
```bash
openclaw cron edit <job-id> \
  --announce \
  --channel telegram \
  --to "-1001234567890" \
  --thread-id 42
```

### Reassign a job to a different agent
```bash
openclaw cron edit <job-id> --agent ops
```

---

## Key Behaviors

**Retry strategy:** Exponential backoff on consecutive errors — 30s → 1m → 5m → 15m → 60m. Resets after a successful run.

**Skipped runs:** Tracked separately. Don't trigger backoff but can opt into failure alerts with `--failure-alert-include-skipped`.

**Silent suppression:** When only a `NO_REPLY` token is returned, nothing posts to chat.

**Model precedence (isolated):** Gateway override → per-job model → stored session override → configured default.

**Migration:** Run `openclaw doctor --fix` to normalize legacy job formats.

---

## Our Cron Jobs

All 19 agents fire daily Discord reports staggered 10:00–11:00pm ET.
- Reports stored at: `/home/ubuntu/onemind/_codex/memory/agents/{id}/reports/`
- Weekly: Spartan, Oracle, Empire (Mon) · Haven (Sun) · Guardian (Wed)
- Monthly: Eden (1st)

Setup script: `install (Install & Deploy)/scripts (Deploy Scripts)/setup-cron-reports.sh`
