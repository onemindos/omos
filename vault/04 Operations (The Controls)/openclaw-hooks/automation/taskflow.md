---
name: OpenClaw — Task Flow
updated: 2026-05-11
---

# Task Flow

Task Flow is the durable multi-step workflow orchestration layer that sits above background tasks. It manages flows with their own state, revision tracking, and sync semantics across gateway restarts.

**Official docs:** [docs.openclaw.ai/automation/taskflow](https://docs.openclaw.ai/automation/taskflow.md)

---

## When to Use Task Flow

| Scenario | Use |
|---|---|
| Single background job | Plain task (`openclaw tasks`) |
| Multi-step pipeline (A → B → C) | Task Flow — managed mode |
| Observe externally created tasks | Task Flow — mirrored mode |
| One-shot reminder | Cron job |
| Multi-step pipeline that must survive restarts | Task Flow + Lobster |

---

## Sync Modes

### Managed Mode

Task Flow owns the lifecycle end-to-end. It creates tasks as flow steps, drives them to completion, and advances flow state automatically.

```
Flow: weekly-report
Step 1: gather-data    → task created → succeeded
Step 2: generate-report → task created → succeeded
Step 3: deliver         → task created → running
```

### Mirrored Mode

Task Flow observes externally created tasks (from cron jobs, CLI, etc.) and keeps flow state in sync without taking ownership of task creation. Example: three independent cron jobs that together form a "morning ops" routine — a mirrored flow tracks their collective progress without controlling when or how they run.

---

## Reliable Scheduled Workflow Pattern

For recurring workflows (market intelligence briefs, daily digests, weekly reports), layer four components:

1. **Cron job** — timing
2. **Persistent session** — build on prior context across runs
3. **Lobster** — deterministic steps, approval gates, resume tokens
4. **Task Flow** — tracks the multi-step run across child tasks, waits, retries, and restarts

Example cron shape:

```bash
openclaw cron add \
  --name "Morning Intelligence Brief" \
  --cron "0 7 * * 1-5" \
  --tz "America/New_York" \
  --session session:market-intel \
  --message "Run the market-intel Lobster workflow. Verify source freshness before summarizing." \
  --announce \
  --channel discord \
  --to "channel:CHANNEL_ID"
```

Use `session:` (not `isolated`) when the recurring workflow needs deliberate history, previous run summaries, or standing context.

### Lobster Workflow File for the Pattern

```yaml
name: market-intel-brief
steps:
  - id: preflight
    command: market-intel check --json
  - id: collect
    command: market-intel collect --json
    stdin: $preflight.json
  - id: summarize
    command: market-intel summarize --json
    stdin: $collect.json
  - id: approve
    command: market-intel deliver --preview
    stdin: $summarize.json
    approval: required
  - id: deliver
    command: market-intel deliver --execute
    stdin: $summarize.json
    condition: $approve.approved
```

### Preflight Checks (Recommended)

Before any summarization step, check:
- Browser availability + profile choice
- API credentials and quota for each source
- Network reachability for required endpoints
- Required tools enabled (`lobster`, `browser`, `llm-task`)
- Failure destination configured so preflight failures are visible

### Data Provenance Fields

Every collected item should carry:

```json
{
  "sourceUrl": "https://example.com/report",
  "retrievedAt": "2026-04-24T12:00:00Z",
  "asOf": "2026-04-24",
  "title": "Example report",
  "content": "..."
}
```

Reject or mark stale items before summarization. LLM steps should receive only structured JSON and preserve provenance fields in output.

---

## Durable State

Each flow persists its own state and tracks revisions so progress **survives gateway restarts**. The flow registry uses SQLite with bounded write-ahead-log maintenance — long-running gateways don't retain unbounded `registry.sqlite-wal` sidecar files.

Revision tracking enables conflict detection when multiple sources attempt to advance the same flow concurrently.

---

## Cancel Behavior

`openclaw tasks flow cancel` sets a **sticky cancel intent** on the flow. Active tasks within the flow are cancelled, and no new steps are started. The cancel intent persists across restarts — a cancelled flow stays cancelled even if the gateway reboots before all child tasks have terminated.

---

## CLI Commands

```bash
# List active and recent flows
openclaw tasks flow list
openclaw tasks flow list --json

# Inspect a specific flow
openclaw tasks flow show <id>

# Cancel a running flow and its active tasks
openclaw tasks flow cancel <id>
```

Note: `openclaw flows` is a redirect — flow commands are subcommands of `openclaw tasks`, not a standalone command.

---

## How Flows Relate to Tasks

Flows **coordinate** tasks, they don't replace them. A single flow may drive multiple background tasks over its lifetime.

- `openclaw tasks` — inspect individual task records
- `openclaw tasks flow` — inspect the orchestrating flow

---

## Related

- [automation/README.md](README.md)
- [tools/lobster.md](../tools/lobster.md)
- [cli/tasks.md](../cli/tasks.md)
- [cli/Cron.md](../cli/Cron.md)
