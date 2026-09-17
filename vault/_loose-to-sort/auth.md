---
title: "Auth"
type: reference
domain: SO
status: active
created: 2026-06-07
updated: 2026-06-07
tags: [auth, credentials, secrets, security, access-control]
---

# Auth

Master authentication reference for OneMind OS. All credentials are stored in **1Password** — this document holds locations and access patterns only, never plaintext secrets.

## Credential Store

All secrets live in **1Password**. Never store tokens, keys, passwords, or API secrets in plaintext anywhere in this vault.

---

## Access Map

| Service | 1Password Vault | Notes |
|---------|----------------|-------|
| Tailscale | Personal | Headscale auth key, admin API token |
| GitHub (onemindos) | Personal | SSH key + PAT — see `~/.ssh/` |
| Fastmail | Personal | FASTMAIL_MCP_TOKEN |
| Proton Bridge | Personal | IMAP/SMTP local creds (127.0.0.1:1143/1025) |
| Nous Portal | Personal | OAuth → `~/.hermes/shared/nous_auth.json` |
| Todoist | Personal | TODOIST_API_KEY |
| AWS | Business | Access key + secret — `~/.aws/credentials` |
| Cloudflare | Business | API token |
| NATS | Personal | Credentials file location TBD |
| 1Password | Physical | Master key — hardware-only access |

---

## SSH Keys

```
~/.ssh/          ← all keys live here, NOT in vault
```

Active keys:
- `id_ed25519` — primary GitHub + fleet SSH key
- Additional node keys stored in 1Password under "OneMind Fleet"

---

## Env Files

Sensitive env vars are injected at runtime — never committed to git.

| File | Location | Purpose |
|------|----------|---------|
| `.env` | `~/onemind/` (gitignored) | Vault-level env vars |
| Hermes config | `~/.hermes/config.yaml` | Provider API keys (ref 1Password) |

---

## Access Control Levels

| Level | Who | Scope |
|-------|-----|-------|
| L0 | Zeus only | Master 1Password, hardware keys, legal |
| L1 | Zeus + explicit approval | Financial, infra changes, external comms |
| L2 | Agents (with notify) | Task creation, messaging agents |
| L3 | Agents (full autonomy) | Read, route, search, schedule |

---

## Rotation Policy

- API tokens: rotate on suspected compromise or quarterly
- SSH keys: rotate on device change or fleet node decommission
- Never reuse keys across pillars (personal vs business)

---

## Notes

<!-- Add new service credentials as 1Password references here — location path only, never the secret itself -->
