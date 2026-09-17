---
name: OpenClaw CLI — Directory
updated: 2026-05-11
---

# `openclaw directory`

ID lookups across supported messaging channels — contacts, groups, and user profiles.

**Official docs:** [docs.openclaw.ai/cli/directory](https://docs.openclaw.ai/cli/directory.md)

---

## Commands

```bash
# Your own profile on a channel
openclaw directory self --channel telegram

# Search contacts/peers
openclaw directory peers list --channel telegram
openclaw directory peers list --channel zalouser --query "name" --limit 50

# Groups
openclaw directory groups list --channel telegram
openclaw directory groups members --channel zalouser --group-id <id>
```

---

## Common Flags

| Flag | Description |
|---|---|
| `--channel <name>` | Required when multiple channel configs exist |
| `--account <id>` | Select account (uses channel default otherwise) |
| `--query <text>` | Filter results |
| `--limit <n>` | Max results |
| `--json` | Machine-readable output |

---

## Platform ID Formats

| Platform | Format |
|---|---|
| Telegram | Chat ID (`-1003905828404`) or `@username` |
| Discord | `channel:<id>` or `user:<id>` |
| Slack | `user:U…` or `channel:C…` |
| WhatsApp | `+15551234567` |
| Matrix | `user:@user:server` |

---

## Usage

Results come from configured allowlists, not live provider databases. Use this to get IDs for `openclaw message send --target`.
