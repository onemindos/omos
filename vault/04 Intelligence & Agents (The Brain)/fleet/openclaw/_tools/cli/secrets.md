---
name: OpenClaw CLI — Secrets
updated: 2026-05-11
---

# `openclaw secrets`

Manages SecretRefs and maintains runtime credential integrity.

**Official docs:** [docs.openclaw.ai/cli/secrets](https://docs.openclaw.ai/cli/secrets.md) · [gateway/secrets](https://docs.openclaw.ai/gateway/secrets.md)

---

## Commands

### Audit
```bash
openclaw secrets audit         # Scan for security issues
openclaw secrets audit --check # Returns exit code 1 for findings, 2 for unresolved refs
```

Audit flags:
- Plaintext secrets in config
- Unresolved SecretRefs
- Credential shadowing
- Legacy residues

### Configure (interactive)
```bash
openclaw secrets configure     # Interactive TTY wizard — requires terminal
```
Guides through: provider setup → credential mapping → preflight validation → apply.

### Apply
```bash
openclaw secrets apply             # Execute saved plan from configure
openclaw secrets apply --dry-run   # Validate without writing
openclaw secrets apply --allow-exec # Allow provider command execution
```

### Reload
```bash
openclaw secrets reload   # Re-resolve SecretRefs and update runtime snapshot via gateway RPC
```
Preserves last-known-good state if resolution fails.

---

## Recommended Workflow

```
1. openclaw secrets audit          → identify issues
2. openclaw secrets configure      → build a plan
3. openclaw secrets apply --dry-run → validate
4. openclaw secrets apply          → apply changes
5. openclaw secrets audit          → verify clean
6. openclaw secrets reload         → update runtime
```

---

## Safety Features

- Plaintext secrets are scrubbed after successful application
- No rollback backups (relies on preflight validation instead)
- Exec providers have strict path security (symlink + ACL restrictions)
- Changes are atomic — partial updates don't occur on failure

---

## Our Setup

All credentials stored in **1Password vault** `aaeyris5z44pdqspk277u2oppi`.
Reference format: `op://aaeyris5z44pdqspk277u2oppi/Item Name/field-name`

Never store secret values in markdown files.

See also: [`security/README.md`](../security/README.md)
