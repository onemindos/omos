---
title: "CLI — op (1Password CLI)"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/secrets, topic/1password]
ai_summary: "1Password CLI (op) quick reference — secret injection, vault access, and credential management. Vault ID: aaeyris5z44pdqspk277u2oppi."
---

# op — 1Password CLI

> The only approved way to use secrets. NEVER hardcode credentials in markdown or code.

---

## Vault

| Field | Value |
|-------|-------|
| Vault Name | `00-24 SO (Sovereign Operations)` |
| Vault ID | `aaeyris5z44pdqspk277u2oppi` |
| Reference format | `op://aaeyris5z44pdqspk277u2oppi/Item Name/field-name` |

---

## Auth

```bash
# Sign in (interactive)
op signin

# Sign in and export session token
eval $(op signin)

# Check status
op whoami
```

---

## Read Secrets

```bash
# Read a single field
op read "op://aaeyris5z44pdqspk277u2oppi/DigitalOcean/api-token"

# Read and copy to clipboard
op read "op://aaeyris5z44pdqspk277u2oppi/DigitalOcean/api-token" | pbcopy
```

---

## Inject Secrets into Commands

```bash
# Run a command with secrets injected from .env file
op run --env-file=.env -- python3 script.py

# Inject inline
OP_TODOIST_TOKEN=$(op read "op://aaeyris5z44pdqspk277u2oppi/Todoist/api-token") python3 script.py
```

---

## List & Browse

```bash
# List all items in vault
op item list --vault "aaeyris5z44pdqspk277u2oppi"

# List items by category
op item list --categories Login --vault "aaeyris5z44pdqspk277u2oppi"

# Get full item details
op item get "DigitalOcean" --vault "aaeyris5z44pdqspk277u2oppi"

# Get specific field
op item get "DigitalOcean" --fields api-token --vault "aaeyris5z44pdqspk277u2oppi"
```

---

## Create & Edit

```bash
# Create a new login item
op item create --category Login --title "New Service" \
  --vault "aaeyris5z44pdqspk277u2oppi" \
  username="user@example.com" password="secret"

# Edit a field
op item edit "Service Name" field-name="new-value"
```

---

## .env File Format

```bash
# .env file using op:// references
TODOIST_API_TOKEN=op://aaeyris5z44pdqspk277u2oppi/Todoist/api-token
OPENAI_API_KEY=op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key
DIGITAL_OCEAN_TOKEN=op://aaeyris5z44pdqspk277u2oppi/DigitalOcean/api-token

# Run with injection
op run --env-file=.env -- node server.js
```

---

## Security Rules

- NEVER store real secrets in markdown files
- NEVER commit `.env` files with real values
- Always use `op read` or `op run` for secret access
- Rotate credentials via 1Password web app, then update `op://` refs


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
