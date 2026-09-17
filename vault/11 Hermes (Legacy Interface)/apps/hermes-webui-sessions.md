# Hermes WebUI — Session Aggregation Architecture

> Last updated: 2026-05-23
> Status: Live

## Overview

The Hermes WebUI (`~/apps/hermes-webui/`) aggregates sessions from **every agent on the OneMind OS fleet** into a single unified session sidebar. You can browse, search, and read transcripts from any agent without leaving the WebUI.

---

## Session Sources

| Source Tag | Label | Storage Path | Schema |
|---|---|---|---|
| `claude_code` | Claude Code | `~/.claude/projects/**/*.jsonl` | Anthropic JSONL — `type` events |
| `api_server` | API (Hermes) | `~/.hermes/sessions.db` (SQLite) | Native Hermes session DB |
| `cli` | CLI (Hermes) | `~/.hermes/sessions.db` (SQLite) | Native Hermes session DB |
| `pi` | Pi Agent | `~/.pi/agent/sessions/**/*.jsonl` | `{type:"message", role, content}` per line |
| `openclaw` | OpenClaw | `~/.openclaw/skill-workshop/*.trajectory.jsonl` | `{type:"model.completed", data:{messagesSnapshot}}` |
| `codex` | Codex | `~/.codex/history.jsonl` | `{session_id, ts, text}` — user prompts only |
| `copilot` | GitHub Copilot | `~/.copilot/session-store.db` (SQLite) | `sessions` + `turns` tables |

---

## Architecture

```
Browser → hermes-webui/server.py (port 8787)
             │
             ├── /api/sessions          → models.py  (Claude Code + native Hermes DB)
             │       └── get_claude_code_sessions()
             │       └── get_hermes_sessions()
             │
             └── /api/external-sessions → external_agent_sessions.py
                     └── get_all_external_agent_sessions()
                             ├── get_pi_sessions()
                             ├── get_openclaw_sessions()
                             ├── get_codex_sessions()
                             └── get_copilot_sessions()
```

**Key design decisions:**
- **Read-only, zero-copy** — WebUI reads agent files directly at request time. No sync job, no DB copy, no writes.
- **Fail-safe** — every parser catches all exceptions silently; a broken agent directory never crashes the session list.
- **No restart required** — adding a new parser is live immediately (Python reads the module fresh per request).
- **Session ID stability** — IDs are `{source}_{sha256(abs_path_or_native_id)[:24]}`, deterministic across restarts.

---

## How Claude Code Sessions Appear

Claude Code sessions come from two places:

1. **Direct Claude Code CLI runs** — `~/.claude/projects/**/*.jsonl` — scraped by `get_claude_code_sessions()` in `models.py`. These show up as `source_tag: "claude_code"`.

2. **Via Synadia harness** — When the `cc` fabric agent (`~/.nats/synadia-agents/agents/claude-code/`) runs, it calls back to the Hermes API (`localhost:8080`), creating `api_server` sessions in `~/.hermes/sessions.db`. These show up as `source_tag: "api_server"`.

---

## Adding a New Agent

Three steps — no routes, no DB changes, no restart:

### 1. Find the storage path

```bash
# Check where the agent writes sessions
ls ~/.myagent/
sqlite3 ~/.myagent/sessions.db ".tables"   # if SQLite
cat ~/.myagent/sessions.jsonl | head -5    # if JSONL
```

### 2. Add a parser to `external_agent_sessions.py`

```python
MYAGENT_SOURCE = "myagent"
MYAGENT_SOURCE_LABEL = "My Agent"

def get_myagent_sessions(sessions_dir=None) -> list:
    # Parse files, build message list, call _session_row()
    ...

def get_myagent_session_messages(sid: str, sessions_dir=None) -> list:
    if not str(sid).startswith(f"{MYAGENT_SOURCE}_"):
        return []
    # Find the right file/row, return messages list
    ...
```

### 3. Wire into the dispatcher

In `get_all_external_agent_sessions()`, add:
```python
(get_myagent_sessions, "MyAgent"),
```

In `get_external_agent_session_messages()`, add:
```python
if sid.startswith(f"{MYAGENT_SOURCE}_"):
    return get_myagent_session_messages(sid)
```

### Message format

Every message in the list must be:
```python
{"role": "user" | "assistant" | "system" | "tool", "content": "...", "timestamp": float}  # timestamp optional
```

### Session row fields (via `_session_row()`)

| Field | Type | Notes |
|---|---|---|
| `session_id` | str | `{source}_{sha256[:24]}` |
| `title` | str | First user message or summary |
| `message_count` | int | Total messages |
| `created_at` | float | Unix timestamp |
| `updated_at` | float | Unix timestamp |
| `source_tag` | str | e.g. `"copilot"` |
| `source_label` | str | e.g. `"GitHub Copilot"` |
| `read_only` | bool | Always `True` for external |

---

## Agent-Specific Schema Notes

### Pi Agent (`~/.pi/agent/sessions/**/*.jsonl`)

```jsonl
{"type": "session", "summary": "...", "title": "..."}
{"type": "message", "message": {"role": "user"|"assistant", "content": [{"type":"text","text":"..."}], "timestamp": "ISO"}}
```

### OpenClaw (`~/.openclaw/skill-workshop/*.trajectory.jsonl`)

```jsonl
{"type": "session.started", "data": {"title": "..."}, "ts": "ISO"}
{"type": "model.completed", "data": {"messagesSnapshot": [{"role":"user","content":"..."}]}, "ts": "ISO"}
```
Note: `messagesSnapshot` contains the full conversation up to that point — the parser deduplicates to avoid repeats from multiple `model.completed` events.

### Codex (`~/.codex/history.jsonl`)

```jsonl
{"session_id": "uuid", "ts": unix_ms, "text": "user prompt"}
```
Codex only stores user prompts, not assistant responses. Sessions show as user-only turns.

### GitHub Copilot (`~/.copilot/session-store.db`)

```sql
sessions (id, summary, cwd, repository, host_type, branch, created_at, updated_at)
turns    (id, session_id, turn_index, user_message, assistant_response, timestamp)
```
Opened read-only via SQLite URI mode (`?mode=ro`).

---

## Live Session Counts (as of 2026-05-23)

| Source | Count |
|---|---|
| OpenClaw | 100 |
| Pi Agent | 5 |
| Codex | 4 |
| GitHub Copilot | 1 |
| **Total external** | **110** |

---

## Files

| Path | Purpose |
|---|---|
| `~/apps/hermes-webui/api/external_agent_sessions.py` | All external parsers + unified dispatcher |
| `~/apps/hermes-webui/api/models.py` | Claude Code + native Hermes session pipeline |
| `~/apps/hermes-webui/api/agent_sessions.py` | `read_importable_agent_session_rows` bridge |
| `~/.hermes/sessions.db` | Native Hermes CLI + API sessions |
| `~/.claude/projects/` | Claude Code native session JSONL files |

---

## Related

- [[Hermes-Overview]] — full Hermes stack overview
- [[openclaw-fleet]] — OpenClaw agent fleet config
- NATS fabric: `agents.prompt.<token>.*` — how Synadia harnesses route work to Claude Code
