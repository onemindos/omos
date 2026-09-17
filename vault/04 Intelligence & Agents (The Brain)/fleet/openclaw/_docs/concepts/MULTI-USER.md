---
title: "OpenClaw Multi-User Routing"
type: reference
status: active
created: 2026-04-12
updated: 2026-04-12
owner: zeus
tags: [type/reference, domain/SO, status/active]
ai_summary: "Architecture for multi-user OpenClaw access. One Legacy instance serves the whole family — each person identified by their connection (Telegram phone number or device token), context file loaded, access level enforced."
---

# OpenClaw Multi-User Routing

> One Legacy. One gateway. Every family member gets personalized AI with the right context and the right access.

---

## Architecture

```
Incoming message (Telegram / iOS app / web)
          │
          ▼
    Identify sender
    (phone number / device token / account)
          │
          ▼
    Load user context
    users/{name}.md
          │
          ▼
    Apply access level
    L4 Zeus → full access
    L3 Wife → LE + HP + home
    L2 Kids → age-gated
          │
          ▼
    Legacy responds
    (scoped to user's context)
```

---

## Access Levels

| Level | Who | Scope |
|-------|-----|-------|
| **L4 — Owner** | Zeus | Everything — all agents, all memory, all commands, all domains |
| **L3 — Partner** | Wife | LE domain (home, family, schedule), HP domain (health, finance shared), home automation |
| **L2 — Child** | Kids | Homework help, schedule, Q&A, safe content — NO adult, financial, or infrastructure info |
| **L1 — Observer** | Guests | Query-only — no writes, no personal context |
| **L0 — Entity** | Dog, assets | No chat — registry + telemetry only |

---

## User Registry

| File | Person | Level | Telegram Number |
|------|--------|-------|-----------------|
| `users/zeus.md` | Zeus Delacruz | L4 | {{fill in}} |
| `users/wife.md` | Wife | L3 | {{fill in}} |
| `users/kids.md` | Kids | L2 | {{fill in each}} |
| `users/dog.md` | Dog | L0 | N/A |

---

## How OpenClaw Identifies the Sender

### Telegram (current primary interface)

OpenClaw's Telegram connector provides `msg.from.id` (Telegram user ID) and `msg.from.phone` (if available).

**Option A — Phone number matching** (recommended):
Map each family member's Telegram phone number to their `users/{name}.md` context file.

**Option B — Telegram user ID matching**:
Map their Telegram user ID directly (more reliable — doesn't change when they change phones).

**Configuration** (add to OpenClaw config or a routing SOP in Legacy's workspace):

```yaml
# In Legacy's SOUL.md or a routing rule — pseudo-config
family_routing:
  "+12125550001": users/wife.md      # wife's number
  "+12125550002": users/kids.md      # kid(s) shared
  "+17185550001": users/zeus.md      # Zeus
  default: users/zeus.md             # fallback
```

### onemind-web / gateway API

Sessions already have a `userId` field. Map user IDs to context files via the same table.

### iOS / macOS OpenClaw app

Each device has a unique pairing token. Map device tokens to users in the gateway config.

---

## Legacy System Prompt Injection

When a non-Zeus user messages, Legacy prepends the user context to the session:

```
[CONTEXT INJECT]
You are talking to: {name}
Their access level: L{n}
Their context file: {contents of users/{name}.md}
Do NOT reveal: {blocked topics for this level}
Tone: {from context file}
[END CONTEXT]
```

This happens at session start — Legacy doesn't need to be retrained, just given the right context at runtime.

---

## Implementation Steps

### Phase 1 — Fill in the registry (now)

1. Open `users/wife.md` — add her real name, Telegram number, context
2. Open `users/kids.md` — add each kid's name, age, Telegram number
3. Open `users/dog.md` — add breed, vet info

### Phase 2 — Telegram number routing (manual for now)

Currently: Legacy doesn't auto-route by sender.

**Workaround today**: Have each family member start their Telegram message with their name, OR have Legacy ask "Who am I talking to?" on first message from an unknown number.

**Proper implementation**: Add a routing rule to Legacy's SOUL.md or create a startup hook that checks `msg.from.phone` against this table and injects context.

### Phase 3 — OpenClaw gateway routing (future)

Once OpenClaw supports per-session user context injection natively (or via a plugin), configure it in `config.yaml`:

```yaml
users:
  - phone: "+1..."
    name: "Wife"
    context_file: "users/wife.md"
    access_level: 3
  - phone: "+1..."
    name: "Kid 1"
    context_file: "users/kids.md"
    access_level: 2
```

---

## Privacy Between Family Members

| Scenario | Behavior |
|----------|---------|
| Wife asks about Zeus's business | Legacy says: "That's Zeus's private domain — ask him directly" |
| Kid asks about finances | Legacy says: "That's a grown-up topic — ask Mom or Dad" |
| Kid asks emergency question | Legacy gives home address, calls Zeus + wife immediately |
| Wife asks about home status | Legacy queries Home Assistant, full access |
| Any user asks "who else uses this?" | Legacy does not reveal other users' Telegram numbers or private data |

---

## Related Files

| File | Purpose |
|------|---------|
| `users/zeus.md` | Zeus's context (owner) |
| `users/wife.md` | Wife context + Telegram number |
| `users/kids.md` | Kids context + access rules |
| `interfaces/telegram.md` | Telegram bot config |
| `openclaw/workspaces/legacy/SOUL.md` | Legacy's persona — update with family routing |
| `openclaw/workspaces/legacy/USER.md` | Zeus context (already populated) |

---

*Part of SO — Sovereign Operations. Entity Registry system.*
