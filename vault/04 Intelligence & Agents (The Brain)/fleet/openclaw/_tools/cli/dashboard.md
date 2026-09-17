---
name: OpenClaw CLI — Dashboard
updated: 2026-05-11
---

# `openclaw dashboard`

Open the Control UI using your current auth.

**Official docs:** [docs.openclaw.ai/cli/dashboard](https://docs.openclaw.ai/cli/dashboard.md)

```bash
openclaw dashboard            # Open in browser
openclaw dashboard --no-open  # Print URL only (don't open browser)
```

---

## Notes

- Resolves configured `gateway.auth.token` SecretRefs when possible
- TLS-enabled gateways print/open `https://` URLs and connect over `wss://`
- For SecretRef-managed tokens, prints a non-tokenized URL to avoid exposing secrets in terminal output or browser-launch arguments
- If clipboard/browser delivery fails, logs a safe manual-auth hint naming `OPENCLAW_GATEWAY_TOKEN` without printing the token value
