---
name: hermes-mcp-config
description: "Use when adding, testing or removing MCP servers in Hermes, or when an MCP server stops working."
version: 2.0.0
author: curator
tags: [hermes, mcp, config, integrations]
triggers:
  - add MCP server
  - remove MCP server
  - MCP server not working
  - audit configured MCP servers
---

# Hermes MCP Server Configuration

Rewritten for Hermes v0.22.5 on `onemind-hermes`. The list of servers that are
configured changes; read it from the instance, never from this file:
`hermes mcp list`.

## Where configuration lives

- The durable source of truth is Git (`platform/hermes/`), not edits made in the
  running pod. A change made only with `hermes mcp add` in the pod is lost if the
  volume is rebuilt and is invisible to review. Make the lasting change in Git
  and let Argo CD apply it; use the CLI for tests and short-lived experiments.
- Credentials never go in `config.yaml` or in this repo. Reference an
  environment variable that a SealedSecret supplies (see `platform/hermes/agent/SECRETS.md`).
- Some MCP products embed the key inside the URL path. Treat such a URL as a
  secret: store the whole URL in a SealedSecret and never print it.

## Commands (v0.22.5)

```bash
hermes mcp list                 # configured servers
hermes mcp test <name>          # connection test
hermes mcp add                  # discovery-first add
hermes mcp configure <name>     # toggle individual tools
hermes mcp login <name>         # force re-auth for an OAuth server
hermes mcp reauth --all         # re-auth every OAuth server
hermes mcp remove <name>
hermes mcp catalog              # Nous-approved MCPs for one-click install
hermes config get mcp_servers   # resolved values
```

If `hermes` is not on PATH, use the full path under `/opt/hermes/`.

## Pre-wire test (do this before saving any URL)

POST an MCP `initialize` to the URL first. Only wire URLs that return a `result`.
When a request lists several URLs, test each one independently and wire only the
passing ones: a wrong key embeds silently and breaks the next session.

```bash
curl -s -X POST "<url>" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","method":"initialize","id":1,"params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"hermes","version":"1.0"}}}' \
  --max-time 10
```

- A `result` in the reply passes. "Invalid API Key", 400 or a connection error fails.
- A 401 from an OAuth server is expected before login; use `hermes mcp login`.
- Some servers only accept POST, so a browser GET returning "Method not allowed" is normal.

## Transports

- Streamable HTTP with a bearer token: `url`, `headers.Authorization` from an env var.
- OAuth: set the server to OAuth, then complete the browser flow on first use.
  OAuth needs a browser, so on a headless pod run the login where a browser exists
  (or use a token-based alternative) and check the result with `hermes mcp test`.
- stdio: `command` plus `args`; use the `env` map for variables (not `headers`).
  A source-built server must be built before it is wired.

## Pitfalls

- MCP tools load at session start. After adding a server or rotating a key,
  start a fresh session (`/reset`).
- Killing a stdio server's child process does not reload credentials: the
  watchdog respawns it from the launch arguments captured at session start.
- A hosted MCP that advertises OAuth may still need a bearer token for
  programmatic clients.
- GitHub's hosted MCP (`api.githubcopilot.com`) needs an account with Copilot;
  a personal access token alone is not enough.
- Disabling a skill (`skills.disabled`) is separate from removing an MCP server.
  Clean up both when decommissioning an integration.
- Egress: the pod's NetworkPolicy must allow the MCP host (`platform/hermes/agent/network-policy.yaml`
  or the instance's own `network.yaml`), or the test fails with a timeout.
