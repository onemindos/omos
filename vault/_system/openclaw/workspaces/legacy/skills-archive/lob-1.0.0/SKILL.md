---
name: lob
description: Send physical mail (letters, postcards) and verify US addresses via the Lob API.
metadata: { "clawdbot": { "emoji": "✉️" } }
---

# Lob

Send real USPS mail and verify addresses from a chat message.

## Tools

- `lob(action, ...)` — universal mail tool

## Actions

| Action           | What it does                     |
| ---------------- | -------------------------------- |
| `send_letter`    | Mail a physical letter via USPS  |
| `send_postcard`  | Mail a postcard                  |
| `verify_address` | Validate a US mailing address    |
| `get_letter`     | Check delivery status + tracking |
| `list_letters`   | List recent sent letters         |
| `cancel_letter`  | Cancel if not yet printed        |

## Env

`LOB_API_KEY` — get from dashboard.lob.com. Use `test_` prefix for sandbox, `live_` for real mail.

## Examples

- "Send a thank-you letter to [client]"
- "Verify this address: 123 Main St, San Francisco, CA 94102"
- "Mail a check for the invoice"
- "Send a postcard to the family from our trip"
- "What's the delivery status of letter ltr_xxx?"
