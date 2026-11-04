---
name: hermes-deployment-topology
description: >
  Configure and reason about Hermes Agent deployment topologies — remote backends,
  proxy mode, web dashboard auth, Desktop remote connections, and the API server.
  Use when the user asks about HERMES_API_URL, remote gateways, thin-client dashboard,
  env vars for proxy/delegation, or how sessions/skills/cron are routed in multi-host setups.
triggers:
  - "HERMES_API_URL"
  - "remote gateway"
  - "proxy mode"
  - "thin client dashboard"
  - "hermes dashboard remote"
  - "delegate all operations"
  - "hermes desktop remote"
  - "GATEWAY_PROXY_URL"
  - "API_SERVER_ENABLED"
tags:
  - hermes
  - deployment
  - configuration
  - remote
  - dashboard
  - proxy
---

# Hermes Deployment Topology

Reference for multi-host Hermes deployments, proxy/delegation mode, Dashboard auth,
and the API server. Also clarifies common misconceptions (especially `HERMES_API_URL`).

## ⚠️ Common Misconception: `HERMES_API_URL` Does Not Exist

**`HERMES_API_URL` is not a documented or real env var in Hermes.** It appears nowhere
in the official env var reference. Do not invent or use it. The real variables are below.

See `references/hermes-remote-envvars.md` for the full authoritative table extracted
from the official docs.

---

## Deployment Topologies

### 1. Local-only (default)
`hermes dashboard` starts a local web server at `http://127.0.0.1:9119`.
Everything runs on the same machine. "No data leaves localhost."
The Chat tab spawns `hermes --tui` as a local subprocess.

### 2. Gateway Proxy Mode (messaging platforms → remote agent)
Used when the **messaging gateway** (Telegram/Discord/Slack) should forward all
agent work to a remote Hermes API server.

```bash
# On the gateway host (handles platform I/O only)
GATEWAY_PROXY_URL=http://remote-agent-host:8642
GATEWAY_PROXY_KEY=<same-as-API_SERVER_KEY-on-remote>
```

**Key fact:** When `GATEWAY_PROXY_URL` is set, the gateway handles platform I/O
**only** — all agent work (sessions, tool use, LLM calls, everything) is delegated
to the remote server. This is NOT inference-only delegation.

### 3. Remote Backend for Hermes Desktop (Electron app)
Hermes Desktop normally launches its own local backend, but can attach to a
`hermes dashboard` running on a remote machine.

```bash
# In Desktop's environment, or via Settings → Gateway → Remote gateway
HERMES_DESKTOP_REMOTE_URL=http://<backend-host>:9119
```

**Prerequisite:** A `hermes dashboard` process must be running on the remote host,
bound to a non-loopback address with auth configured:
```bash
# On remote host — run with systemd/tmux
hermes dashboard --host 0.0.0.0 --port 9119 --no-open
```

**The remote dashboard is still a full backend** with its own local sessions DB,
skills, cron, etc. There is no "stateless" or "pure proxy" dashboard mode.

### 4. API Server (OpenAI-compatible endpoint)
Exposes Hermes as an OpenAI-compatible endpoint for external frontends (Open WebUI,
LobeChat, etc.). Uses `POST /v1/chat/completions` with model `hermes-agent`.

```bash
API_SERVER_ENABLED=true
API_SERVER_KEY=<bearer-token>           # required
API_SERVER_PORT=8642                    # default
API_SERVER_HOST=127.0.0.1              # default; set 0.0.0.0 to expose
API_SERVER_CORS_ORIGINS=http://localhost:3000  # for browser access
API_SERVER_MODEL_NAME=hermes-agent     # name advertised on /v1/models
```

---

## Dashboard Authentication Env Vars

Three auth providers ship with the dashboard. For internet-facing or remote installs,
always configure one — a non-loopback bind (`--host 0.0.0.0`) fails closed without auth.

| Provider | Key Variables |
|---|---|
| **Basic (username/password)** | `HERMES_DASHBOARD_BASIC_AUTH_USERNAME`, `HERMES_DASHBOARD_BASIC_AUTH_PASSWORD`, `HERMES_DASHBOARD_BASIC_AUTH_PASSWORD_HASH`, `HERMES_DASHBOARD_BASIC_AUTH_SECRET` |
| **Nous Portal OAuth** | `HERMES_DASHBOARD_OAUTH_CLIENT_ID` (provision with `hermes dashboard register`) |
| **Self-hosted OIDC** | `HERMES_DASHBOARD_OIDC_ISSUER`, `HERMES_DASHBOARD_OIDC_CLIENT_ID`, `HERMES_DASHBOARD_OIDC_SCOPES` |

Set `HERMES_DASHBOARD_BASIC_AUTH_SECRET` explicitly (32+ bytes, `openssl rand -base64 32`)
so sessions survive restarts. Blank → random per-process (logged out on every restart).

---

## Multi-Profile Gateway: Default Profile Must Be in `Serves:`

When `gateway.multiplex_profiles: true`, the gateway serves multiple Hermes profiles from one process. Run `hermes gateway status` and check the `Serves:` line — it lists every profile the gateway is actually handling.

**Pitfall:** If the gateway was previously started from a non-default profile context, the default profile (and its configured platforms) will be absent from `Serves:` even though the gateway is running and other profiles are healthy. The missing profile's platforms silently never connect — no error, just absent from `gateway_state.json`.

**Fix:**
```bash
hermes gateway stop
hermes gateway run   # run from the correct HERMES_HOME context
# then verify:
hermes gateway status  # 'Serves:' must include 'default' (or your expected profile)
cat /opt/data/gateway_state.json | python3 -c \
  "import sys,json; [print(p,v['state']) for p,v in json.load(sys.stdin).get('platforms',{}).items()]"
```

---

## What Is and Isn't Delegated

| Scenario | What's delegated |
|---|---|
| `GATEWAY_PROXY_URL` set | ALL agent work (sessions, tools, LLM, skills, cron) — gateway does platform I/O only |
| Desktop → remote backend | ALL data is on the remote host; Desktop is a thin rendering shell |
| Dashboard web UI | **Nothing** delegated — always local backend, local storage |
| API server (`/v1/chat/completions`) | Agent runs locally; client gets streaming responses |

---

## Quick Reference: All Relevant Env Vars

See `references/hermes-remote-envvars.md` for the complete table with descriptions
quoted directly from the official docs.
