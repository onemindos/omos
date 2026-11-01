# Hosted Whop MCP: startup and private authentication recovery

**Scope: hosted MCP only.** The standalone Whop CLI is not required for current
operations. This change does not install, configure, upgrade, remove or authorize
that CLI. Any already-installed binary and existing CLI state remain untouched.
It does not authorize publication, merge, deployment, purchases, transfers,
payouts or physical actions. Tool discovery is not permission to execute actions;
remote sampling is explicitly disabled.

## What startup owns

`server.json` contains only public metadata: `https://mcp.whop.com/mcp`, OAuth,
300-second connect/login timeouts, and disabled sampling. The hosted MCP is the
existing authenticated integration; no alternate stdio server is introduced.

`initial-config` runs `reconcile.py` on **every** pod start, including an existing
PVC. It calls the supported Hermes CLI:

```sh
hermes config set mcp_servers.whop '<public JSON from server.json>'
```

This is the **Hermes configuration CLI**, not the Whop standalone CLI. Only
`mcp_servers.whop` is Git-owned. The reconciler never replaces the whole server
map, removes a server, starts OAuth, or reads/writes token files. Other MCP
entries, Zoho settings, environment references and all unrelated configuration
are preserved. The existing Zoho startup reconciler and SealedSecrets are
unchanged. Public defaults include the same Whop entry for an empty PVC.

Kustomize generates a hash-suffixed ConfigMap containing just the public manifest
and reconciler; its volume reference changes on edits, triggering a normal
GitOps rollout. No new executable installer, auth sidecar, CLI env override or
credential Secret is added. Initial config ownership repair remains the existing
narrow home/config fix before the supported Hermes launcher drops privilege.

## Private state and encrypted backups

- Hermes home: `/opt/data` on PVC `onemind-hermes-data`.
- Public runtime settings: `/opt/data/config.yaml`.
- Private hosted-MCP OAuth state: `/opt/data/mcp-tokens/whop*.json`.
- Preserve private directory/file ownership and permissions on restore (0700 for
  token directory, 0600 for token files). Never print token files or place them
  in public ConfigMaps, Git, logs, skills or chat.

`../../backup.yaml` declares a VolSync source for the **entire**
`onemind-hermes-data` PVC at 01:25 UTC, Restic repository Secret
`onemind-hermes-restic`, snapshot copy method, retention daily/weekly/monthly/yearly
7/4/3/1. `../../backup.sealedsecret.yaml` holds the encrypted repository access
configuration. Both are already in this application's Kustomization. This change
adds no static token snapshot, backup job or refresh-token Secret.

Git reconstructs public configuration, **not authentication**. The password-free
config-export job is not an OAuth backup. Do not export ephemeral OAuth client
registrations or refresh tokens to Git; any future secret configuration must use
reviewed SealedSecrets. The encrypted complete PVC backup is the recovery source
for the existing private OAuth store.

Manifest coverage is verified locally, **not live backup success or restore
readiness**. Before relying on recovery, an authorized operator must inspect a
recent successful VolSync synchronization, encrypted repository access, the
independently protected Restic password/SealedSecrets controller key, and perform
an isolated restore drill. Keep restore changes GitOps-managed and obtain storage
operation approval; never delete or overwrite production storage as a test.

## Recovery procedures (authorized operator)

### Replacement pod, same PVC

Let Argo recreate the pod. The existing runtime config is not reseeded; only
Git-owned entries reconcile. Tokens remain on the same PVC and no login is
started by this reconciler. Check in the main container:

```sh
HERMES_HOME=/opt/data /opt/hermes/bin/hermes mcp test whop
```

Then read a non-mutating identity/workspace result and confirm the authenticated
account and intended business. Do not equate discovery with authorization for
financial or other state-changing operations. Refresh support does not guarantee
a grant can never expire or be revoked.

### Lost/empty PVC

First restore the encrypted complete PVC through the approved VolSync/Restic
workflow in repository-root `docs/portability/REBUILD-FROM-ZERO.md`. Use the name,
size and storage class in `../../pvc.yaml` for a restored live claim; do not replace
a live claim or deploy a throwaway restore without approval. Quiesce writers,
preserve ownership/modes, and verify completion before starting the gateway.
A normal startup seeds missing public config and reconciles Whop.

**Without usable restored private OAuth state, Git alone cannot recreate the
login.** The owner must authorize a fresh hosted-MCP grant on the intended Hermes
instance. Other integrations need their own credentials recovered too. A backup
from before login does not contain that login.

### Revoked grant, rotated refresh token or unusable restored state

A restore cannot undo provider-side revocation or refresh-token rotation.
Reauthorize **only Whop** through the supported owner-assisted flow:

```sh
HERMES_HOME=/opt/data /opt/hermes/bin/hermes mcp login whop
HERMES_HOME=/opt/data /opt/hermes/bin/hermes mcp test whop
```

Do not use `reauth --all`, delete working token files, or start competing flows.
Complete a pending login on its original host/listener; remote browser callbacks
need supported routing to that listener. Keep callback URLs, codes and state
private. Verify a read-only identity/workspace result afterwards, then confirm a
new encrypted backup includes the new state. Old rotated refresh tokens can be
invalid even when the backup decrypts successfully.

Standalone Whop CLI OAuth, if independently authorized, is a **different grant**
and does not authorize this hosted MCP. That CLI is outside this GitOps change;
there is no CLI login/reinstallation work in this recovery procedure.

## Local checks

```sh
python -m unittest discover -s scripts/tests -p 'test_gitops_*mcp.py' -v
kustomize build platform/hermes/onemind-hermes
make validate
```

The existing CI `test_gitops*.py` quality-gate step discovers the focused Whop
contracts. Tests check exact single-server mutation, suppressed provider errors,
rejected non-public configuration, startup mounts, encrypted backup wiring and
absence of standalone-CLI/auth wiring. When an actual Hermes CLI is present, an
isolated fixture test runs two real config updates, checks unrelated settings
and confirms private token-fixture bytes and modes are unchanged. No test starts
OAuth, uses real credentials, or performs financial actions. Scratch artifacts
are preserved. Authenticated hosted-MCP checks and live disaster restore remain
parent/owner acceptance gates, not claims made by these local tests.
