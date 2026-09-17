---
title: "Contexts"
section: "09-CLI-TOOLING"
type: note
tags: [nats, cli, contexts, connection-profiles]
source: https://docs.nats.io/using-nats/nats-tools/nats_cli#nats-contexts
created: 2026-06-14
---

# Contexts

A **context** is a saved connection profile — server URL, credentials, TLS, JetStream domain — so you don't retype connection flags. The `nats` CLI manages them with `nats context`.

See also: [[NATS/09-CLI-TOOLING/nats-cli|The nats CLI]]

---

## Why

Without contexts every command needs `-s ... --creds ... --js-domain ...`. With contexts you save it once and switch with `--context name` (or set a default).

```shell
nats context add hub \
  --server nats://omos-mbp:4222 \
  --creds /path/zeus.creds \
  --description "OneMind hub"

nats context add hub-sys \
  --server nats://omos-mbp:4222 \
  --user sys --password '***' \
  --description "system account (for server report)"

nats context ls                  # list contexts
nats context select hub          # set the default
nats context info                # show current
nats stream ls --context hub-sys # use a specific context for one command
```

---

## What a Context Stores

| Field | Example |
|-------|---------|
| server URL(s) | `nats://omos-mbp:4222` |
| credentials | creds file, token, or user/password |
| TLS | ca/cert/key paths |
| JetStream domain | `--js-domain hub` |
| nsc account/user | for nsc-managed setups |
| description | human label |

Contexts live under `~/.config/nats/context/*.json`.

---

## Patterns

- **One context per environment/account.** OneMind benefits from a regular-bus context **and** a separate **system-account** context (`hub-sys`) — because `nats server report jetstream` needs SYS creds (verified pitfall). Switch with `--context`.
- **Default + overrides.** `nats context select` sets the everyday default; pass `--context other` to reach elsewhere for one command.
- **Scripts** can pin `--context` so automation always hits the right server/account.

## Sources
- Contexts — https://docs.nats.io/using-nats/nats-tools/nats_cli#nats-contexts
- natscli — https://github.com/nats-io/natscli
