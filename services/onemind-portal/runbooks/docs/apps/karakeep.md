# Karakeep

> Generated from the reviewed catalog (workload-map.yaml, structure.yaml). It lists what is recorded, not what is currently healthy; check the Live tab for current state.

Private native bookmark and read-later workspace. The user-authorized GitOps configuration starts web/worker, Chrome and Meilisearch, with native Authentik sign-in restricted to the existing registered OneMind owner and daily encrypted whole-volume backups. Check live readiness; owner enrollment and restore acceptance require observed evidence.

## Where it runs

- Catalog entity: `component:default/karakeep` in System `personal-tools`
- Argo CD application: `platform-karakeep`
- Namespace: `karakeep`; the dedicated OIDC bootstrap hook runs in `auth`

## When it looks wrong

1. Open the app page. The status dot combines Argo CD, workload readiness, alerts and uptime; grey means unknown, not healthy.
2. **Live** tab: confirm pods and replicas. **Logs** tab: open the generated Loki query. **Metrics** tab: open the workload dashboard with this app's variables.
3. Find the cause before acting. A restart that hides the cause is not a fix.

## Safe recovery

- Change the cluster only through Git and Argo CD (pull request, merge, sync). Do not `kubectl apply`, patch or delete by hand, and do not run Helm directly.
- A Git change is reversible with `git revert`; Argo CD then restores the previous state.
- Data stores: confirm the latest backup in the [backups runbook](../backups.md) before any change that touches storage.
- Preserve `data-pvc`, `meilisearch-pvc`, their sealed credentials and the matching Restic passwords. Daily whole-volume snapshots are crash-consistent; verify SQLite integrity and attachments in an isolated restore.
- Sign in at `https://karakeep.tail717ba0.ts.net` through the existing OneMind Authentik owner. Password/local registration is disabled. After the first successful owner login, set `DISABLE_SIGNUPS=true` through Git/Argo to close OIDC account creation too.

## Generated links

- Source: `onemindos/omos`

## Escalate

If this app is unreachable and the cause is not clear after the Live and Logs tabs, treat it as an incident: record what was observed (time, tab, message) and follow the [resilience runbook](../resilience.md).
