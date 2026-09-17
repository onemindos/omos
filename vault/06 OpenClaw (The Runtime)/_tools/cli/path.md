---
name: OpenClaw CLI — Path (oc:// Addressing)
updated: 2026-05-11
---

# `openclaw path`

The `oc://` addressing system — inspect, read, and edit files in agent workspaces using a URI scheme.

**Official docs:** [docs.openclaw.ai/cli/path](https://docs.openclaw.ai/cli/path.md)

---

## URI Scheme

```
oc://<agent-id>/<workspace-relative-path>
```

**Examples:**

```
oc://legacy/IDENTITY.md
oc://spartan/MEMORY.md
oc://guardian/SOUL.md
oc://main/USER.md
```

---

## Commands

```bash
# Read a workspace file
openclaw path read oc://legacy/IDENTITY.md
openclaw path read oc://legacy/MEMORY.md --json

# List workspace contents
openclaw path ls oc://legacy/
openclaw path ls oc://legacy/ --recursive

# Write / overwrite a file
openclaw path write oc://legacy/USER.md --content "Updated user context"
openclaw path write oc://spartan/MEMORY.md < updated-memory.md

# Check if a path exists
openclaw path exists oc://legacy/SOUL.md

# Resolve to filesystem path
openclaw path resolve oc://legacy/IDENTITY.md
# → /home/ubuntu/.openclaw/agents/legacy/workspace/IDENTITY.md
```

---

## Workspace Files

Standard files in every agent workspace:

| File | Purpose |
|---|---|
| `IDENTITY.md` | Agent persona and role definition |
| `AGENTS.md` | Universal agent instructions |
| `SOUL.md` | Agent values and behavioral traits |
| `TOOLS.md` | Available tools reference |
| `USER.md` | User context and preferences |
| `MEMORY.md` | Persistent memory index |

---

## Our Workspace Locations

Workspace paths on omos-core:

```
/home/ubuntu/.openclaw/agents/{id}/workspace/    ← newer agents
/home/ubuntu/.openclaw/workspace-{id}/           ← legacy agents
```

Use `openclaw path resolve` to get the exact filesystem path for any agent.

---

## Notes

- `path write` immediately updates the file on disk — the agent picks up changes on next session start
- `oc://` URIs only address agent workspaces, not the full VPS filesystem
- Use `path ls oc://legacy/ --recursive` to audit all workspace files for an agent
