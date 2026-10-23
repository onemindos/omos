# OneTAK STAC pilot: observe, diagnose and recover

> **Status: backend removed.** `onetak/stac-api` and `onetak/onetak-geo-db` (with their
> catalog entries, volume and backup source) were removed from Git because STAC was never
> set up: the `stac` database has no `pgstac` schema. STAC Browser has no live catalog
> behind it until a backend is deliberately re-added. The former portal operations code
> has also been removed. The text below describes the former
> setup and is kept as history.

## Ownership and request path

`onetak/stac-api` uses `onetak/onetak-geo-db`; both are owned by Argo application
`platform-onetak`. Never let two controllers own the same workload.

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

- [Backup/recovery runbook](backups.md)

The page exposes no restart, write API, schema-initialization or promotion action.
The existing machine access boundary stays read-only.
