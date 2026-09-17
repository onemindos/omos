---
name: OpenClaw — Lobster (Workflow Shell)
updated: 2026-05-11
---

# Lobster

Lobster is a workflow shell that lets OpenClaw run multi-step tool sequences as a single, deterministic operation with explicit approval checkpoints.

**Official docs:** [docs.openclaw.ai/tools/lobster](https://docs.openclaw.ai/tools/lobster.md)

> Lobster is one authoring layer above detached background work. For flow orchestration above individual tasks, see [Task Flow](../automation/README.md) (`openclaw tasks flow`).

---

## Why Use Lobster

Without Lobster, complex workflows require many sequential LLM tool calls — each costing tokens, each requiring the model to re-orchestrate. Lobster moves that orchestration into a typed runtime:

- **One call instead of many** — OpenClaw runs one Lobster tool call and gets a structured result
- **Approvals built in** — side effects (send email, post message) halt the workflow until explicitly approved
- **Resumable** — halted workflows return a token; approve and resume without re-running prior steps
- **Deterministic + auditable** — pipelines are data (YAML/JSON), so they're easy to log, diff, replay, and review
- **Safety enforced** — timeouts, output caps, sandbox checks enforced by the runtime, not each script

---

## Enable Lobster

Lobster is an **optional plugin tool** — not enabled by default.

Recommended (additive, preserves core tools):

```json
{
  "tools": {
    "alsoAllow": ["lobster"]
  }
}
```

Per-agent (to enable for a specific agent only):

```json
{
  "agents": {
    "list": [
      {
        "id": "legacy",
        "tools": {
          "alsoAllow": ["lobster"]
        }
      }
    ]
  }
}
```

> Do NOT use `tools.allow: ["lobster"]` unless you intend restrictive allowlist mode — it will strip core tools.

---

## How It Works

Lobster runs workflows **in-process** inside the gateway — no external subprocess. The embedded runner executes the pipeline and returns a JSON envelope directly.

If a pipeline pauses for approval, the tool returns a `resumeToken` to continue later.

---

## Inline Pipeline (tool call)

```json
{
  "action": "run",
  "pipeline": "inbox list --json | inbox categorize --json | inbox apply --json | approve --preview-from-stdin --limit 5 --prompt 'Apply changes?'",
  "timeoutMs": 30000
}
```

Resume after approval:

```json
{
  "action": "resume",
  "token": "<resumeToken>",
  "approve": true
}
```

---

## Workflow Files (.lobster)

Pipelines can be defined as YAML files with `name`, `args`, `steps`, `env`, `condition`, and `approval` fields:

```yaml
name: inbox-triage
args:
  tag:
    default: "family"
steps:
  - id: collect
    command: inbox list --json
  - id: categorize
    command: inbox categorize --json
    stdin: $collect.stdout
  - id: approve
    command: inbox apply --approve
    stdin: $categorize.stdout
    approval: required
  - id: execute
    command: inbox apply --execute
    stdin: $categorize.stdout
    condition: $approve.approved
```

Run a workflow file:

```json
{
  "action": "run",
  "pipeline": "/path/to/inbox-triage.lobster",
  "argsJson": "{\"tag\":\"family\"}"
}
```

- `stdin: $step.stdout` passes a prior step's output
- `condition` / `when` gates steps on `$step.approved`

---

## Tool Parameters

### `run`

| Field | Description | Default |
|---|---|---|
| `action` | `"run"` | required |
| `pipeline` | Inline pipeline string or `.lobster` file path | required |
| `cwd` | Relative working directory (stays within gateway dir) | — |
| `timeoutMs` | Abort if exceeded | 20000 |
| `maxStdoutBytes` | Abort if output exceeds this size | 512000 |
| `argsJson` | JSON string of args passed to workflow file | — |

### `resume`

| Field | Description |
|---|---|
| `action` | `"resume"` |
| `token` | Resume token from a `needs_approval` response |
| `approve` | `true` to continue, `false` to cancel |

---

## Output Envelope

Lobster returns a JSON envelope with one of three statuses:

| Status | Meaning |
|---|---|
| `ok` | Finished successfully |
| `needs_approval` | Paused — `requiresApproval.resumeToken` needed to resume |
| `cancelled` | Explicitly denied or cancelled |

```json
{
  "ok": true,
  "status": "needs_approval",
  "output": [{ "summary": "5 need replies, 2 need action" }],
  "requiresApproval": {
    "type": "approval_request",
    "prompt": "Send 2 draft replies?",
    "items": [],
    "resumeToken": "..."
  }
}
```

---

## LLM Steps (llm-task)

For structured LLM steps inside a workflow, enable the optional `llm-task` plugin and call it from Lobster:

```json
{
  "plugins": { "entries": { "llm-task": { "enabled": true } } },
  "agents": {
    "list": [{ "id": "legacy", "tools": { "alsoAllow": ["llm-task"] } }]
  }
}
```

> **Important limitation:** The embedded Lobster runner does NOT reliably support `openclaw.invoke` for nested calls. Use `llm-task` as a direct tool call outside Lobster, or use the standalone Lobster CLI when `openclaw.invoke` nesting is needed.

---

## Known Limitation

The bundled Lobster plugin runs **in-process** inside the gateway. In embedded mode, `openclaw.invoke` does not automatically inherit gateway URL/auth context for nested OpenClaw CLI calls. This pattern is **not currently reliable** in the embedded runner:

```
openclaw.invoke --tool llm-task --action json --args-json '{ ... }'
```

Use direct `llm-task` tool calls outside Lobster for now.

---

## Safety

- Local in-process only — no network calls from the plugin itself
- Sandbox-aware — disabled when tool context is sandboxed
- Timeouts and output caps enforced by the embedded runner
- No secrets management — Lobster calls OpenClaw tools that handle OAuth

---

## Troubleshooting

| Error | Fix |
|---|---|
| `lobster timed out` | Increase `timeoutMs` or split a long pipeline |
| `lobster output exceeded maxStdoutBytes` | Raise `maxStdoutBytes` or reduce output size |
| `lobster returned invalid JSON` | Ensure pipeline runs in tool mode and prints only JSON |
| `lobster failed` | Check gateway logs for embedded runner error details |

---

## Install Standalone CLI (Optional)

The embedded runner ships with the Lobster plugin — no separate binary needed for standard use.

For development or external pipelines, install the standalone CLI:

```bash
# From https://github.com/openclaw/lobster
# Ensure `lobster` is on PATH
```

---

## Related

- [Tools overview](README.md)
- [Automation & cron](../automation/README.md)
- [cli/hooks.md](../cli/hooks.md)
- [cli/tasks.md](../cli/tasks.md)
