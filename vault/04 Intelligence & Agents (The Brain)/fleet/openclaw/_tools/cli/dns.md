---
name: OpenClaw CLI — DNS
updated: 2026-05-11
---

# `openclaw dns`

DNS helpers for wide-area discovery using Tailscale and CoreDNS.

**Official docs:** [docs.openclaw.ai/cli/dns](https://docs.openclaw.ai/cli/dns.md)

---

## Commands

```bash
openclaw dns setup                              # Plan CoreDNS setup (dry run — no changes)
openclaw dns setup --domain openclaw.internal   # With explicit domain
openclaw dns setup --apply                      # Apply changes (macOS + Homebrew CoreDNS only)
```

---

## What `setup` Shows

- Resolved discovery domain
- Zone file path
- Current Tailnet IPs
- Recommended `openclaw.json` discovery settings
- Tailscale Split DNS nameserver configuration values

---

## Constraints

- `--apply` supports **macOS only**
- Requires **Homebrew CoreDNS** installation
- If `--domain` is omitted, uses `discovery.wideArea.domain` from config
- Automatically bootstraps zone files and ensures CoreDNS import statements exist

---

## Our Setup

We use Tailscale for VPS access (`100.123.126.34` for omos-core). Wide-area DNS-SD discovery is not currently configured.
