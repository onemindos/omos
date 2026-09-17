---
title: Hermes — Legacy Interface
domain: 11
pillar: SO
status: active
updated: 2026-05-21
tags: [hermes, ai, legacy, api, gui, infrastructure]
---

# Hermes — Legacy Interface

> Hermes Agent is the local AI runtime that powers Legacy AI. This domain documents the full stack: core agent, API server, and both GUI front-ends.

## What Is Hermes

Hermes Agent (`hermes-agent`) is an open-source, self-hostable AI agent framework built by Nous Research. It runs locally on Zeuss-MBP and serves as the backend brain for Legacy AI.

- GitHub: https://github.com/nousresearch/hermes-agent
- Docs: https://hermes-agent.nousresearch.com/docs
- Install location: `~/.hermes/`
- Config: `~/.hermes/config.yaml`
- Env: `~/.hermes/.env`

## Architecture

```
┌─────────────────────────────────────────┐
│            Zeus (You)                   │
│   CLI / Telegram / GUI front-end        │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Hermes API Server               │
│   http://localhost:8080/v1              │
│   OpenAI-compatible REST API            │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Hermes Agent Core               │
│   Models: Claude (Bedrock) + Azure GPT  │
│   Tools: 100+ (terminal, browser, etc.) │
│   Memory: ~/.hermes/memory/             │
│   Skills: ~/.hermes/skills/             │
└─────────────────────────────────────────┘
```

## Tailscale Access

| Service | Tailscale URL |
|---------|--------------|
| Hermes API | `http://omos-mbp.tail717ba0.ts.net:8080/v1` |
| Hermes Dashboard | `http://omos-mbp.tail717ba0.ts.net:9119` |
| Tailscale IP | `100.98.178.23` |
| MagicDNS hostname | `omos-mbp.tail717ba0.ts.net` |

## Related Notes

- [[Hermes-API-Server]] — API server config & endpoints
- [[Hermes-Desktop-Fathah]] — Primary GUI (Hermes Agent app)
- [[Hermes-Desktop-Dodo]] — Power GUI (HermesDesktop workbench)
- [[Hermes-Config-Reference]] — Full config reference
