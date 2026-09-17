---
name: zoho-mail-1.0.0
description: Read, search, move, and send emails in Zeus's Zoho Mail account via the Zoho MCP endpoint.
metadata: { "clawdbot": { "emoji": "📬" } }
---

# Zoho Mail

Full email management for Zeus's Zoho Mail account (`zeus@zeusdelacruz.com`). Connects via the Zoho MCP HTTP API — no IMAP required.

## Tools

- `zoho_get_account_info()` — Account overview: aliases, send identities, storage, TFA status
- `zoho_list_folders()` — All folders with message counts and unread counts
- `zoho_list_emails(folder, limit?, offset?, unread_only?)` — List emails in a folder
- `zoho_search_emails(query, limit?)` — Search with Zoho syntax
- `zoho_read_email(message_id)` — Read full email content
- `zoho_move_emails(message_ids[], to_folder)` — Move emails to a folder
- `zoho_mark_read(message_ids[], read?)` — Mark emails as read/unread
- `zoho_send_email(to, subject, body, from_identity?, cc?)` — Send email from any identity

## Folders

Use these names in `folder` parameters:

| Name           | Maps To                                  |
| -------------- | ---------------------------------------- |
| `inbox`        | Main INBOX                               |
| `hp`           | 25-49 HP (Holistic Performance)          |
| `ui`           | 00-24 SO (Sovereign Operations)          |
| `le`           | 50-74 LE (Legacy Evolution)              |
| `ge`           | 75-99 GE (Generational Endeavor) |
| `sent`         | Sent                                     |
| `trash`        | Trash                                    |
| `spam`         | Spam                                     |
| `newsletter`   | Newsletter                               |
| `notification` | Notification                             |

Or pass a raw folder ID (from `zoho_list_folders`).

## Send Identities

| Key        | Display Name                  | Sends From            |
| ---------- | ----------------------------- | --------------------- |
| `hp`       | Holistic Performance          | zeus@onemindcodex.com |
| `ui`       | Sovereign Operations          | zeus@onemindcodex.com |
| `le`       | Legacy Evolution              | zeus@onemindcodex.com |
| `ge`       | Generational Endeavor | zeus@onemindcodex.com |
| `personal` | Zeus (Personal Brand)         | zeus@onemindcodex.com |
| `onemind`  | Zeus (OneMind)                | finance@legacyhq.ai   |

Default (no identity): sends from `zeus@zeusdelacruz.com`

## Email Addresses on This Account

- **Primary**: `zeus@zeusdelacruz.com`
- **Aliases**: `zeus@onemindcodex.com`, `hp@zeusdelacruz.me`, `ui@zeusdelacruz.me`, `le@zeusdelacruz.me`, `ge@zeusdelacruz.me`
- **External**: `zeus.vaults@pm.me` (ProtonMail)

## Search Syntax Examples

```
from:oscar subject:credit
has:attachment is:unread
subject:invoice after:2025-01-01
from:zeus@onemindcodex.com in:sent
```

## Examples

- "Do I have any unread emails?" → `zoho_list_emails(folder: 'inbox', unread_only: true)`
- "Check my HP folder" → `zoho_list_emails(folder: 'hp')`
- "Find emails about tax" → `zoho_search_emails(query: 'subject:tax')`
- "Read that invoice email" → `zoho_read_email(message_id: '<id>')`
- "Move these to trash" → `zoho_move_emails(message_ids: [...], to_folder: 'trash')`
- "Send Oscar a reply" → `zoho_send_email(to: 'olarios@consumerdirect.com', subject: '...', body: '...')`

## Environment Variables

- `ZOHO_MCP_URL` — **Required**. Full MCP endpoint URL including API key.
- `ZOHO_ACCOUNT_ID` — Optional. Defaults to Zeus's account ID.

Add to `/opt/onemind/.env.tpl`:

```
ZOHO_MCP_URL=op://aaeyris5z44pdqspk277u2oppi/Zoho MCP/url
```

## Notes

- Filters/rules and blocked senders require Zoho web UI (organization admin scope)
- `zoho_move_emails` is safe — it moves, never deletes
- Trash folder in Zoho is still recoverable; true deletion requires web UI
