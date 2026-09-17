# Security Validator Instructions

> Reference copy — deployed version lives at `_system/security/validator.md`

## How Agents Validate Commands

Before executing any command that modifies files, runs shell commands, or accesses sensitive paths, agents MUST:

### Step 1: Check Against patterns.yaml

1. Load `_system/security/patterns.yaml`
2. Match the command against `bash_patterns` (blocked → confirm → alert)
3. Match any file paths against `path_protection` rules

### Step 2: Apply Decision

| Match | Action |
|-------|--------|
| **blocked** | ❌ STOP immediately. Do not execute. Log to events. |
| **confirm** | ⚠️ Ask Zeus for explicit confirmation before proceeding. |
| **alert** | 📢 Log the action to events, then proceed with caution. |
| **no match** | ✅ Proceed normally. |

### Step 3: Log Security Event

Append to `_system/security/events/YYYY-MM.jsonl`:

```jsonl
{"timestamp":"ISO-8601","event":"blocked|confirmed|alerted|allowed","command":"the command","reason":"why","agent":"agent-name"}
```

### Step 4: Vault Rules

Always check `vault_rules` section:
- **never_commit**: Scan staged files for secrets before any `git commit`
- **always_confirm**: These actions require explicit Zeus approval regardless of context

## Validation Flow

```
Command received
    ↓
Check bash_patterns.blocked → BLOCK if match
    ↓
Check bash_patterns.confirm → ASK if match
    ↓
Check bash_patterns.alert → LOG if match
    ↓
Check path_protection → Apply access level
    ↓
Check prompt_injection → BLOCK if match
    ↓
ALLOW → Execute command
    ↓
Log event
```
