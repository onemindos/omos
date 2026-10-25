# STAC imagery catalog: observe, diagnose and recover

The current application is `platform-stac` in namespace `stac`. Its PostgreSQL
backend is database/role `stac` on the shared CNPG cluster `onemind-data`.
The former `onetak/stac-api` and `onetak-geo-db` were retired; do not restore their
configuration as a shortcut.

## Request path and acceptance

- `stac-api`: read/search API, tailnet ingress `stac.tail717ba0.ts.net`.
- `stac-api-writer`: transactions enabled, no ingress; only the WebODM sync pod
  selector may reach it through the writer NetworkPolicy.
- `stac-browser`: static browser, tailnet ingress `stac-browser.tail717ba0.ts.net`.
- `webodm-stac-sync`: every-ten-minute ingestion job; it needs existing Spaces
  credentials and a reviewed WebODM source credential. An empty catalog is valid
  API acceptance, but does not prove imagery was ingested.

Check both API Services have ready endpoints. A green browser pod or Argo Healthy
alone is insufficient: during the contained October 24 incident both API
Deployments were absent while the browser remained Ready.

Verify landing, collection listing and bounded search return valid JSON. Read API
collection/item writes must return 405/404. Check generated links use the intended
tailnet origin. Verify actual asset access separately; private Spaces URLs do not
become public merely because an Item contains them.

## Known migration incident and narrow correction

The original migration crashed the shared PostgreSQL backend on October 24, 2026.
STAC sync is manual, the migration Job has zero retries and a ten-minute deadline.
Do not repeat a failed production migration to diagnose it.

Controlled exact-image tests on October 25 isolated the failure to the AGE library
preloaded into the server. PgSTAC's `partition_sys_meta` view uses
`JOIN LATERAL COALESCE(...)`; the pinned AGE parser mishandles that expression
([upstream issue/fix](https://github.com/apache/age/pull/2259)). Without AGE preload
the unmodified migration passed. With preload it produced signal 11. Creating an
AGE extension in the STAC database did not fix it.

The reviewed STAC compatibility wrapper changes that one FROM expression into an
equivalent one-row lateral SELECT. It preserves the original pypgstac CLI and full
transaction, verifies exact SQL checksum/version, and rejects unreviewed upgrades.
It changes no global preload, native library or production role privileges.
The underlying AGE bug remains relevant to other callers using the problematic
expression; any native-library repair needs its own shared-database review.

## Recovery procedure

1. Inspect the failed Job, both container logs and current CNPG/client health.
   Preserve original evidence. Distinguish a SQL error, native crash and missing
   API Deployment; restarting a browser cannot fix a missing schema.
2. Use `platform-stac-validation` only for synthetic tests. It has loopback-only
   PostgreSQL, emptyDir storage, deny-all networking and no production credentials.
   Save logs/status before replacing failed disposable Jobs through Argo CD.
   A ConfigMap update alone does not rerun a failed Job.
3. Verify the exact candidate on both native validation cases and confirm a fresh
   recoverable shared-database backup. Review the scoped migration diff and role
   impact before a single explicit production sync. Never delete production PVCs,
   reset WAL, change shared credentials or disable AGE globally as a shortcut.
4. During sync watch migration completion, all three CNPG instances and shared
   database clients, then API rollout/endpoints and real collection/search paths.
   Stop on failure; do not add retries. Keep automation disabled through acceptance.
5. Record observed timestamps and remaining ingestion/imagery gaps. Retire temporary
   diagnostics only after their evidence is saved; preserve backups and data.

Repository evidence and exact image/SQL pins: `platform/stac/INCIDENT-2026-10-25.md`.
The source policy is Git and Argo CD; do not add imperative or Helm-owned workloads.
