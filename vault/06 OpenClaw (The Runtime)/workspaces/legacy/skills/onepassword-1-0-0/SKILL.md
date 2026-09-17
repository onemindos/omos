# onepassword-1.0.0

Read, write, and update items in 1Password vaults using the 1Password Service Accounts REST API (`api.1password.com`). No Connect server required.

## Primary use case

Store VPS config backups as secure notes so secrets are never in git:

- `openclaw.json` (full real config with channel IDs, paths, etc.)
- `.env` (all API keys and tokens)

Whenever the VPS config changes, Legacy can update the 1Password note automatically.

## Required env vars

| Variable                   | Description                                                    |
| -------------------------- | -------------------------------------------------------------- |
| `OP_SERVICE_ACCOUNT_TOKEN` | 1Password service account token (from 1password.com/developer) |
| `OP_VAULT_ID`              | Default vault ID — use `OneMind - Sovereign Operations` vault  |

## Tools

### op_save_note

Create or update a secure note. Primary tool for config backups.

```
Create new:
  title: "VPS openclaw.json – 2026-03-06"
  content: <full json content>
  tags: ["path:ui", "category:infrastructure"]

Update existing:
  title: "VPS openclaw.json"
  content: <updated content>
  item_id: <existing item id>
```

### op_get_item

Retrieve a full item including all fields.

```
item_id: "abc123xyz"
vault_id: (optional, defaults to OP_VAULT_ID)
```

### op_list_items

List all items in a vault, optionally filtered by tag.

```
filter_tag: "category:infrastructure"
```

### op_get_field

Get a single field value from an item.

```
item_id: "abc123xyz"
field_label: "API Key"
```

### op_update_field

Update a single field value on an existing item.

```
item_id: "abc123xyz"
field_label: "password"
value: "new-value"
```

## Setup

1. Go to [1password.com/developer](https://1password.com/developer) → Service Accounts → New Service Account
2. Grant **Read/Write** access to `OneMind - Sovereign Operations` vault
3. Copy the token → add to VPS `/opt/onemind/.env`:
   ```
   OP_SERVICE_ACCOUNT_TOKEN=ops_...
   OP_VAULT_ID=<vault-id-from-1password>
   ```
4. Get vault ID: in 1Password web app → Vaults → click vault → copy ID from URL

## Saving a config backup (example workflow)

Legacy should run this whenever `openclaw.json` changes on the VPS:

1. `op_list_items` with `filter_tag: "category:infrastructure"` — find existing backup note
2. If found: `op_save_note` with `item_id` to update in-place
3. If not found: `op_save_note` without `item_id` to create new, with tags `["path:ui", "category:infrastructure"]`

## Security notes

- Service account tokens have vault-scoped access only — not full account access
- Tokens can be revoked at any time from the 1Password developer portal
- Never log or print the token value
- The token in `.env` is loaded by systemd `EnvironmentFile=` — never written to disk in plaintext beyond that


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
