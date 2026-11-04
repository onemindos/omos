# Hermes Remote/Proxy/Dashboard Env Vars — Authoritative Reference

Extracted from https://hermes-agent.nousresearch.com/docs/reference/environment-variables
(verified 2026-07-11).

## DOES NOT EXIST

| Variable | Status |
|---|---|
| `HERMES_API_URL` | **Not documented, not real.** Zero search results on the official docs site. Do not use. |

---

## Gateway Proxy Mode

| Variable | Description |
|---|---|
| `GATEWAY_PROXY_URL` | URL of a remote Hermes API server to forward messages to (proxy mode). When set, the gateway handles platform I/O only — **all agent work is delegated to the remote server**. Also configurable via `gateway.proxy_url` in config.yaml. |
| `GATEWAY_PROXY_KEY` | Bearer token for authenticating with the remote API server in proxy mode. Must match `API_SERVER_KEY` on the remote host. |

---

## API Server (local OpenAI-compatible endpoint)

| Variable | Description |
|---|---|
| `API_SERVER_ENABLED` | Enable the OpenAI-compatible API server (true/false). Runs alongside other platforms. |
| `API_SERVER_KEY` | Bearer token for API server authentication. **Required whenever the API server is enabled.** |
| `API_SERVER_PORT` | Port for the API server (default: 8642) |
| `API_SERVER_HOST` | Host/bind address for the API server (default: 127.0.0.1). API_SERVER_KEY is still required on loopback; use a narrow API_SERVER_CORS_ORIGINS allowlist for browser access. |
| `API_SERVER_CORS_ORIGINS` | Comma-separated browser origins allowed to call the API server directly (e.g. `http://localhost:3000`). Default: disabled. |
| `API_SERVER_MODEL_NAME` | Model name advertised on `/v1/models`. Defaults to the profile name (or `hermes-agent` for the default profile). Useful for multi-user setups where frontends like Open WebUI need distinct model names per connection. |

---

## Web Dashboard & Hermes Desktop

> "Auth for the web dashboard and for connecting Hermes Desktop to a remote backend. Per the secrets-only convention, credentials belong in ~/.hermes/.env; the OAuth client_id is better set under dashboard.oauth in config.yaml (env wins when set)."

### Basic (username/password) provider

| Variable | Description |
|---|---|
| `HERMES_DASHBOARD_BASIC_AUTH_USERNAME` | Username for bundled username/password auth provider. Activates the provider when set together with a password. Overrides `dashboard.basic_auth.username`. |
| `HERMES_DASHBOARD_BASIC_AUTH_PASSWORD` | Plaintext password for the basic provider (hashed in-memory at load). Wins over a config `password_hash`. Overrides `dashboard.basic_auth.password`. |
| `HERMES_DASHBOARD_BASIC_AUTH_PASSWORD_HASH` | scrypt password hash (preferred — no plaintext at rest). Compute with `python -c "from plugins.dashboard_auth.basic import hash_password; print(hash_password('PW'))"`. Overrides `dashboard.basic_auth.password_hash`. |
| `HERMES_DASHBOARD_BASIC_AUTH_SECRET` | HMAC key (32+ bytes, base64/hex/raw) signing stateless session tokens. **Set explicitly so sessions survive restarts.** Blank → random per-process. Overrides `dashboard.basic_auth.secret`. |
| `HERMES_DASHBOARD_BASIC_AUTH_TTL_SECONDS` | Access-token lifetime (default 12h). Overrides `dashboard.basic_auth.session_ttl_seconds`. |

### Nous Portal OAuth provider

| Variable | Description |
|---|---|
| `HERMES_DASHBOARD_OAUTH_CLIENT_ID` | OAuth client id (`agent:{instance_id}`) for the gated/public dashboard. Activates the Nous provider. Provision with `hermes dashboard register`. Overrides `dashboard.oauth.client_id`. |
| `HERMES_DASHBOARD_PUBLIC_URL` | Complete public URL the dashboard is reached at, for OAuth callback construction behind reverse proxies. Overrides `dashboard.public_url`. |

### Self-hosted OIDC provider

| Variable | Description |
|---|---|
| `HERMES_DASHBOARD_OIDC_ISSUER` | OIDC issuer URL. Required to activate the self-hosted OIDC provider. Overrides `dashboard.oauth.self_hosted.issuer`. |
| `HERMES_DASHBOARD_OIDC_CLIENT_ID` | Public OIDC client id (authorization-code + PKCE). Required to activate. Overrides `dashboard.oauth.self_hosted.client_id`. |
| `HERMES_DASHBOARD_OIDC_SCOPES` | Requested OIDC scopes (default `openid profile email`). Overrides `dashboard.oauth.self_hosted.scopes`. |

### Hermes Desktop (Electron app)

| Variable | Description |
|---|---|
| `HERMES_DESKTOP_REMOTE_URL` | **(Desktop side)** Base URL of the remote backend, e.g. `http://host:9119`. When set, overrides the in-app Gateway URL; you still sign in from the Gateway settings panel. |
| `HERMES_DESKTOP_HERMES` | Desktop backend command override (for packagers/Nix or troubleshooting). |
| `HERMES_DESKTOP_IGNORE_EXISTING` | Set to `1` to make Desktop ignore an existing hermes on PATH during backend resolution. |

---

## Remote Dashboard Setup Recipe

For a remote `hermes dashboard` that Hermes Desktop (or a browser) can connect to:

```ini
# systemd [Service] section
EnvironmentFile=%h/.hermes/.env
ExecStart=/path/to/venv/bin/python -m hermes_cli.main dashboard \
    --host 0.0.0.0 --port 9119 --no-open
```

```bash
# ~/.hermes/.env on the remote host
HERMES_DASHBOARD_BASIC_AUTH_USERNAME=myuser
HERMES_DASHBOARD_BASIC_AUTH_PASSWORD=<from-sealedsecret>
HERMES_DASHBOARD_BASIC_AUTH_SECRET=$(openssl rand -base64 32)
```

**Important gotchas:**
- `--host 0.0.0.0` is required; a loopback bind rejects remote clients at the socket layer regardless of credentials.
- The "remote backend is ready" probe only hits `GET /api/status` (a public endpoint). The live chat WebSocket (`/api/pty`) requires authentication — a "ready" probe success does NOT mean chat will work.
- The remote dashboard is a **full local backend** with its own sessions DB, skills, cron. There is no stateless proxy-only dashboard mode.
- The DNS-rebinding guard requires the `Host` header to match the bind address. Enter the remote URL exactly as the dashboard is reachable.
