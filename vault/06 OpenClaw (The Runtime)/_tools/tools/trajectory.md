---
name: OpenClaw — Trajectory (Flight Recorder)
updated: 2026-05-11
---

# Trajectory Bundles

Trajectory capture is OpenClaw's per-session **flight recorder**. It records a structured timeline for each agent run, then `/export-trajectory` packages the current session into a redacted support bundle.

**Official docs:** [docs.openclaw.ai/tools/trajectory](https://docs.openclaw.ai/tools/trajectory.md)

Use it to answer:
- What prompt, system prompt, and tools were sent to the model?
- Which tool calls led to this answer?
- Did the run time out, abort, compact, or hit a provider error?
- Which model, plugins, skills, and runtime settings were active?
- What usage and prompt-cache metadata did the provider return?

---

## Quick Start

```bash
/export-trajectory
# or
/trajectory

# With custom output name
/export-trajectory bug-1234
```

Exports to: `.openclaw/trajectory-exports/openclaw-trajectory-<session>-<timestamp>/`

CLI equivalent:
```bash
openclaw sessions export-trajectory --session-key "agent:main:telegram:direct:123" --workspace .
```

---

## Access

Owner command — sender must pass command authorization and owner checks for the channel.

In group chats, OpenClaw sends the approval prompt and export result to the owner **privately**, not back to the shared room.

---

## What Gets Recorded

Runtime events captured automatically:

| Event | Description |
|---|---|
| `session.started` | Session open |
| `trace.metadata` | Runtime metadata |
| `context.compiled` | System prompt + tool compilation |
| `prompt.submitted` | Full prompt sent to model |
| `model.fallback_step` | Fallback reason, next model, chain position |
| `model.completed` | Model response received |
| `trace.artifacts` | Final artifacts |
| `session.ended` | Session close |

Transcript events also reconstructed: user messages, assistant messages, tool calls, tool results, compactions, model changes, labels.

---

## Bundle Files

| File | Contents |
|---|---|
| `manifest.json` | Schema, source files, event counts |
| `events.jsonl` | Ordered runtime + transcript timeline |
| `session-branch.json` | Redacted active transcript branch |
| `metadata.json` | OC version, OS/runtime, model, config, plugins, skills |
| `artifacts.json` | Final status, errors, usage, compaction count |
| `prompts.json` | Submitted prompts and prompt-building details |
| `system-prompt.txt` | Latest compiled system prompt |
| `tools.json` | Tool definitions sent to the model |

---

## Privacy & Limits

Trajectory bundles are for debugging — not public posting. OpenClaw auto-redacts:
- Credentials and secret-like payload fields
- Image data
- Local state paths
- Workspace paths (replaced with `$WORKSPACE_DIR`)
- Home directory paths

Size limits:
- Runtime sidecar files: 50 MiB
- Session files: 50 MiB
- Runtime events: 200,000
- Total exported events: 250,000

**Review bundles before sharing.** Redaction is best-effort.

---

## Capture Location

Runtime trajectory events written alongside the session file:
```
.trajectory.jsonl
.trajectory-path.json  ← pointer file
```

Custom directory:
```bash
export OPENCLAW_TRAJECTORY_DIR=/var/lib/openclaw/trajectories
```

Disable capture:
```bash
export OPENCLAW_TRAJECTORY=0
```

---

## Troubleshooting

| Issue | Fix |
|---|---|
| No runtime events in export | Ensure `OPENCLAW_TRAJECTORY=0` is NOT set |
| Command rejects output path | Use relative names like `bug-1234`, not `/tmp/...` |
| Size error | Session or sidecar exceeded safety limits; start fresh session |

---

## One Mind Use Case

When an agent on omos-core produces a wrong result or takes an unexpected action, trajectory gives full forensic visibility — exact prompts sent, exact tool calls made, exact model responses, timing, errors. This is the primary debugging tool for the fleet.

```bash
# From omos-core, in the agent session that misbehaved:
/export-trajectory incident-2026-05-11
# Then inspect the bundle in .openclaw/trajectory-exports/
```
