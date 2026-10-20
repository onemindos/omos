# OneTAK STAC pilot: observe, diagnose and recover

The existing Development pilot is `component:default/stac-browser` in Project
`onetak`. Its Operations tab joins reviewed topology with timestamped read-only
observations. It does not migrate or mutate a service.

## Ownership and request path

STAC Browser is managed by OpenChoreo from the Git-owned component/release/binding
under `platform/openchoreo/projects/onetak/stac-browser`, reconciled by Argo
`platform-openchoreo-projects`. The browser proxies read-only catalog access to
`onetak/stac-api`, which uses `onetak/onetak-geo-db`, both still owned by
`platform-onetak`. Generated data-plane namespace names are discovered, not pinned
as lasting service identities. Never let two controllers own the same workload.

The catalog contract `api:default/onetak-stac-readonly` documents only the reviewed
landing/collection/search subset. It is not a full upstream conformance claim.

## Do not confuse these signals

- Kubernetes readiness and observed generation: container/controller state.
- API landing HTTP 200: the landing document responds; it may not query the DB.
- Collections/search HTTP 200 plus valid JSON: those particular read paths work.
- Item `datetime`: acquisition/phenomenon time, not ingestion completion time.
- Prometheus values: scoped Deployment-pod queries, with query evaluation times;
  not proof of raw scrape recency, data delivery or an ingestion SLA.
- Argo Healthy/Synced: shared configuration application, not every app transaction.
- Volume sync or a ready DB pod: not an application-consistent restore test.

The browser probe uses the existing Cilium Gateway service and reviewed HTTPRoute
hostname. It does not bypass the data-plane network policy with direct cross-project
service access, and it is not an end-user tailnet-path test.

All unavailable sources stay unknown. No ingestion SLA is currently configured.
Collection results can be paginated; the one-item sorted search is bounded.
No object/asset URLs, coordinates, credentials or raw upstream error bodies are
returned by the operational observation endpoint.

## Finding observed 2026-10-20

The landing page returned 200 while collections/search returned 500. API logs
reported `UndefinedFunctionError: function search(jsonb) does not exist`.
Read-only inspection of database `stac` found no `pgstac` namespace. This explains
why a green pod is insufficient; it is not authorization to install schema blindly.

Before a schema repair: inventory existing tables/ownership and exact image/schema
compatibility, preserve an encrypted backup, verify isolated restoration, and
review an Argo-owned initialization/migration Job in a separate PR. Do not run
`DROP`, restore over production, switch the DB, or ingest invented demo data to
manufacture a healthy result. Recheck collections and search after any approved
repair, then test the real user browser path. Retain a rollback/recovery plan;
an application rollback alone does not necessarily reverse a database migration.

## Useful surfaces

- [Pilot Deploy](/catalog/default/component/stac-browser/environments)
- [Pilot Logs](/catalog/default/component/stac-browser/runtime-logs)
- [Pilot Metrics](/catalog/default/component/stac-browser/metrics)
- [STAC API catalog](/catalog/default/component/k8s-onetak-stac-api-746ca5f2)
- [Database catalog](/catalog/default/component/k8s-onetak-onetak-geo-db-a22783b9)
- [Backup/recovery runbook](backups.md)

The page exposes no restart, write API, schema-initialization or promotion action.
The existing machine access boundary stays read-only; this page does not grant it
new direct access to the OpenChoreo or service-operations REST proxies.
