---
id: todo-vps-creds-1password
tags: [security, credentials, 1password, openclaw, #p2]
owner: Zeus
status: todo
created: 2026-05-12
---

## Migrate VPS Credentials to 1Password op:// References

### Problem

All secrets are stored **plaintext** in two files on omos-core:

| File | Credentials |
|------|-------------|
| `~/.openclaw/openclaw.json` | AWS keys, GitHub PAT, Telegram bot token, Circle token, Home Assistant token, MCP server tokens (Todoist, Zoho Mail) |
| `~/.openclaw/gateway.systemd.env` | Twilio, Azure OpenAI, Azure Speech, Deepgram, LiteLLM, Google, Groq, FAL |

If omos-core is ever compromised, an attacker gets everything from two files.

### Goal

Replace all raw credential values with `op://` secret references per the pattern already prescribed in `CLAUDE.md`:

```
op://aaeyris5z44pdqspk277u2oppi/Item Name/field-name
```

OpenClaw supports secret ref resolution at gateway startup. The `gateway.systemd.env` can be injected via `op run`.

### Steps

1. **Add each credential to 1Password** vault `aaeyris5z44pdqspk277u2oppi` with a clear item name per credential
2. **Replace in `openclaw.json`** — use OpenClaw's secret ref syntax for each sensitive field (check `openclaw config --help` for secret ref format)
3. **Replace in `gateway.systemd.env`** — use `op run -- systemctl --user start openclaw-gateway` or inject via 1Password CLI at service start
4. **Verify** — restart gateway, confirm it resolves all secrets cleanly, no plaintext in either file
5. **Update systemd service** to use `op run` wrapper if needed

### Credentials to migrate

- `OPENCLAW_GATEWAY_TOKEN` — keep in env, low risk, but can reference 1P
- `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` — high priority
- `GH_TOKEN` — high priority
- `AZURE_OPENAI_API_KEY` / `AZURE_OPENAI_VOICE_KEY` / `AZURE_SPEECH_KEY`
- `DEEPGRAM_API_KEY`
- `LITELLM_API_KEY`
- `FAL_KEY`
- `GOOGLE_API_KEY`
- `GROQ_API_KEY`
- `TWILIO_ACCOUNT_SID` / `TWILIO_AUTH_TOKEN`
- Discord bot token (already in 1P — verify openclaw.json references it via op://)
- Telegram bot token (already in 1P — verify openclaw.json references it via op://)
- Circle MCP token
- GitHub MCP PAT
- Todoist API key
- Home Assistant token
- Zoho Mail token
