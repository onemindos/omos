# Hermes credential cleanup guide

**Status:** ready to do whenever you choose. Nothing here has been done. It is the last thing standing between the Hermes export (`hermes-export-rules.md`) and a first successful run, and it also removes real credentials from the model prompt and from the backups.
**Written:** 2026-10-16. Counts come from a read-only scan of the live pod. No secret values appear in this document.

## Why it matters

The live `config.yaml` (root and profile copies) and the agent-written skills hold real credentials. Skills are loaded into the model prompt when used, everything is on the volume, and the whole volume goes into the nightly backups. The export refuses to run while any remain (that is the safety design working).

## How credentials should live

- **Config:** Hermes expands `${VAR}` and `${env:VAR}` in string values of `config.yaml` when it loads it (`_expand_env_vars` in `hermes_cli/config.py`). So a config value becomes `Bearer ${GITHUB_MCP_TOKEN}` and the real value lives in the environment. The export template already uses the same `${NAME}` form.
- **Environment:** the pod gets its environment from the `hermes-env` Secret (`envFrom` in `deployment.yaml`). The committed form is a SealedSecret, `platform/hermes/bridge/hermes-env.sealedsecret.yaml`. Adding a variable means resealing that file and letting Argo apply it; the pod picks it up on its next restart.
- **Skills:** write "use `$NAME` from the environment" instead of the literal value.
- **Not `op://` yet.** `config.yaml` already has `secrets.onepassword` references, but the `op` CLI is not on Hermes' PATH (it is at `/opt/data/home/.local/bin/op`), so they do not resolve today. Plain environment variables work now.

## What to rotate, and where it is

Counts are from the dry run of `scripts/hermes_identity_export.py` (64 findings). Not every hit is a real credential: some are examples or placeholders. Check each one, and treat the exact lines from the dry run as the checklist.

| Credential | In the config | In the skills | Suggested variable |
|------------|---------------|---------------|--------------------|
| Notion API token | `mcp.servers.notion.args` and `mcp_servers.notion.args` in the root config and in all four overseer profile configs | `notion-workspace` (2) | `NOTION_TOKEN` |
| GitHub Copilot MCP token | `mcp_servers.github.headers.Authorization` | | `GITHUB_MCP_TOKEN` |
| Composio API key | `mcp_servers.composio.headers.x-consumer-api-key` | `composio` (2), `composio-mcp-patterns` (7), `legacy-integrations` (4) | `COMPOSIO_API_KEY` |
| Matrix access token | `platforms.matrix.token` | | `MATRIX_ACCESS_TOKEN` |
| Dashboard login | `dashboard.basic_auth.password_hash` and `.secret` (the username and password are already a Secret: `hermes-dashboard-auth`) | | `DASHBOARD_BASIC_AUTH_SECRET`, `DASHBOARD_BASIC_AUTH_PASSWORD_HASH` |
| WordPress | `mcp_servers.wordpress.headers.Authorization` (server disabled, domain retired) | `wordpress-legacy` (1) | remove the block, or `WORDPRESS_AUTH` |
| GoHighLevel | `mcp_servers.ghl.headers.Authorization` and `mcp_servers.ghl-pro.env.GHL_API_KEY` in the `generational-empire` profile | `legacy-integrations/references/ghl-mcp-pro.md` | `GHL_API_KEY`, `GHL_AUTHORIZATION` |
| Tailscale OAuth client | | `doks-cluster-ops/references/tailscale-egress-device-map.md` (2) | `TS_OAUTH_CLIENT_SECRET` |
| TAK passwords and feed credentials | | `tak-server` (5 plus references), `tak-mission-ops` (3), `tak-mission-runner` (1) | `TAK_CERT_PASSWORD` and others (see also the separate TAK default-password task) |
| NATS, Argo CD, CloudTAK and ETL passwords and tokens | | `legacy-self-awareness` (7), `omos-feed-builder` (1), `hermes-deployment-topology` (1), `doks-cluster-ops` (1) | `NATS_PASSWORD` (already in `hermes-env`), `ARGOCD_PASSWORD`, `CLOUDTAK_TOKEN` |
| Google Workspace OAuth | | `google-workspace/scripts/gws_bridge.py` and `setup.py` | `GOOGLE_OAUTH_CLIENT_SECRET` |
| Telegram bot token | | `legacy-self-awareness/references/hermes-telegram-config.md` | `TELEGRAM_BOT_TOKEN` |
| Not yet identified | | `android-adb-management` (2), `find-zeus` (1) | look at the lines |
| Memory | `memories/MEMORY.md` holds a plaintext WordPress application password | | remove the entry |

About 50 more hits exist in upstream skills (examples in their documentation, such as `api_key=...`). Those are not ours and are not exported; ignore them.

## Procedure

1. **Back up first.** Confirm the last VolSync run is recent (`kubectl get replicationsource -n hermes`), or copy the config and the skills you will edit to `*.pre-cleanup-<date>` next to them.
2. **Rotate each credential at its provider** (Notion integrations page, GitHub, Composio dashboard, Matrix `/_matrix/client/v3/logout` plus a new login, Tailscale OAuth clients, GoHighLevel, the TAK and database passwords). Do the highest-value ones first: Notion, GitHub, Composio, Tailscale.
3. **Reseal the new values into `hermes-env`.** Per `CLAUDE.md`: fetch the cluster public key, build the Secret with the new keys (keep the existing keys such as `DIGITALOCEAN_AI_KEY` and the NATS ones), pipe through `kubeseal --scope namespace-wide`, write `platform/hermes/bridge/hermes-env.sealedsecret.yaml`, and open a PR. Argo applies it.
4. **Edit the live config.** Replace each value with `${NAME}` (the root `config.yaml` and the four profile configs). Hermes keeps a `.bak` copy as it does for other edits.
5. **Edit the skills.** Use the dry-run report to find each line and replace the literal with `$NAME` or a sentence saying where it comes from.
6. **Restart the gateway** (a rollout of `hermes-legacy`, which picks up the new `hermes-env`). Check that the agent registers on the bus and that the MCP servers you use still connect (`hermes config check`).
7. **Delete the plaintext password from memory** and, optionally, from older backups' relevance by relying on the revoked credential.
8. **Run the dry run again** (command below). The goal is `Status: OK`.
9. **Turn the export on:** set the repository variable `HERMES_IDENTITY_EXPORT_ENABLED` to `true` and run the workflow once by hand.

## The dry run

Read-only; it writes nothing to the volume.

```bash
kubectl -n hermes exec -i deploy/hermes-legacy -c hermes -- \
  /opt/hermes/.venv/bin/python - < scripts/hermes_identity_export.py > /tmp/export.tar
tar -xOf /tmp/export.tar EXPORT-REPORT.md
```

A refused run lists each problem as a file and line (never a value). Exit code 2 means refused.

## Effort

About one to two hours: roughly 15 minutes per major provider to rotate, 20 minutes to reseal and apply, 30 minutes to edit config and skills, and a restart. The edits are on the pod, so a careless one can break a skill but not the cluster, and the backups and the `.pre-cleanup` copies make every step reversible.
