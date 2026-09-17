---
name: OpenClaw CLI — Security
updated: 2026-05-11
---

# `openclaw security`

Security auditing and optional remediation for OpenClaw deployments.

**Official docs:** [docs.openclaw.ai/cli/security](https://docs.openclaw.ai/cli/security.md)

---

## Commands

```bash
openclaw security audit                 # Basic audit (no plugins, no live probes)
openclaw security audit --deep          # Full audit: live gateway probes + plugin collectors
openclaw security audit --fix           # Apply safe remediations
openclaw security audit --json          # Machine-readable for CI/policy checks
openclaw security audit --password <pw> # Override auth for deep probes
openclaw security audit --token <tok>   # Override auth for deep probes
```

---

## What Audit Checks

| Risk Category | Detail |
|---|---|
| Multi-user ingress | Without proper sandboxing |
| Small models | (<300B params) running with web/browser tools |
| Webhook misconfiguration | Weak tokens, unrestricted paths |
| Docker network modes | Dangerous modes (e.g. host networking) |
| mDNS exposure | Unintended local network discovery |
| Channel allowlists | Mutable identifiers instead of stable IDs |

---

## What `--fix` Does

- Flips `groupPolicy="open"` → `groupPolicy="allowlist"`
- Tightens file permissions

What it does NOT do: rotate credentials, disable tools, rewrite plugins.

---

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Clean — no findings |
| `1` | Findings present |
| `2` | Unresolved SecretRefs (when using `--check`) |

---

## Related

- [`secrets.md`](secrets.md) — credential management
- [`approvals.md`](approvals.md) — exec approval policies
- [`security/README.md`](../security/README.md) — vault security docs
