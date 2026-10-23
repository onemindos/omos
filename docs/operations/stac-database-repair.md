# STAC database repair: preflight and recovery gates

Status: **superseded.** The STAC API (`onetak/stac-api`), the `onetak-geo-db` server and its
VolSync source were removed from Git by the `cleanup/remove-stac-geodb` change, because STAC was
never set up (no `pgstac` schema, see the evidence below). This page is kept as that evidence;
the repair plan below is no longer pursued and nothing here is authorized.
This is a dependency repair for the phase-3 frontend pilot; the shared PostgreSQL
server remains managed by Argo CD outside OpenChoreo.

## Confirmed evidence

The API in `platform/onetak/stac.yaml` points at database `stac`, role `geo`, on
`onetak-geo-db.onetak.svc.cluster.local`. The first inspection used both
`default_transaction_read_only=on` and `BEGIN READ ONLY` and found PostgreSQL 17.11,
no `pgstac` schema/tables or search functions, and no enabled PostGIS extension.
Installed extensions were `plpgsql` 1.0, `timescaledb` 2.29.1 and
`timescaledb_toolkit` 1.24.0. The existing API returns 500 for collection listing and
search, with missing-function errors. These observations do not establish that
all user data is absent, why initialization is missing, or whether another database
contains a previous installation.

VolSync source `onetak/volume-50fb74ff16c6` reported a successful backup of
`onetak-geo-db-data` at `2026-10-18T05:20:18Z`. That confirms a volume backup, not a
fresh logical STAC backup or a tested application restore. The repository's default
logical-backup image targets PostgreSQL 15; it must not be reused unchanged for
this PostgreSQL 17 server.

## Read-only preflight results (2026-10-18, 23:13 EDT)

PR [#293](https://github.com/onemindos/omos/pull/293) merged with all four CI
checks passing. Argo reconciled merge `4adc570a`; `platform-onetak` is
Synced/Healthy. Job `stac-preflight-20261018` completed successfully at
`2026-10-19T03:13:25Z`, with `default_transaction_read_only` and
`transaction_read_only` both `on` in each connection.

| Observation | `geo` | `stac` |
|---|---|---|
| Owner / reported size | `geo` / 18 MB | `geo` / 10 MB |
| Enabled PostGIS | 3.6.4 | absent |
| Available PostGIS default | 3.6.4 | 3.6.4 |
| PgSTAC schema / search functions | absent | absent |
| Public relations (tables/views/etc.) | 3 | 0 |
| Public tables | `spatial_ref_sys` | none |

`btree_gist` 1.7, `pg_trgm` 1.6 and `unaccent` 1.1 are available but not enabled
in either database. Both have the previously observed TimescaleDB extensions.
The configured `geo` role has superuser, database-creation and role-creation
privileges; none of `pgstac_admin`, `pgstac_read` or `pgstac_ingest` exists.
No role, password or grant was changed.

The preflight found no PgSTAC catalog to recover by simply switching the API to
`geo`. PostGIS availability means the missing extension is not, by itself, a
reason to replace or upgrade the shared server image. No application rows were
read; these results do not explain the historical absence of catalog data or
establish the contents of every database in the instance.

The completed Job is removed after recording this evidence. The queries and
constraints remain reviewable in PR #293; no initialization Job replaces it.

## Diagnostic method

The Argo-owned Job connected separately to `geo` and `stac`, with read-only defaults, explicit
read-only transactions, bounded queries, no retries and a 120-second deadline.
It reads only system-catalog metadata: database owner/size, role flags, installed
and available extensions, schema/relation counts, table names and function
signatures. It does not read application rows, expose credential values, create roles,
change grants or initialize a schema. The password is an existing Secret reference.

If other evidence suggests a usable catalog exists elsewhere, investigate that
configuration before proposing a fresh installation. The intended catalog and
ingestion sources still need to be established; an initialized empty catalog does
not supply imagery by itself.

## Prepare recovery and isolated proof

After inventory, prepare a separate reviewable GitOps change for a fresh backup
and isolated restore. Specify source, destination, image digests, extension versions,
resource requests, storage cost, cleanup and recovery commands. Do not create an
unreviewed new database server or change the shared server image merely to satisfy
an extension dependency. Keep the test instance inaccessible to production clients;
its API must use only the test database. Use no production ingestion writers there.

The backup must match PostgreSQL 17 and preserve the extension requirements and
object ownership needed for recovery. Record timestamp, integrity checks and an
isolated restore result. If using a physical volume backup, verify database crash
recovery and account for the whole shared instance; restoring that volume over the
live server could roll back unrelated databases. Prefer a targeted database recovery
procedure for the eventual STAC change. Preserve existing backup credentials and
archives. Do not put plaintext dumps or credentials in Git or logs.

Select and pin a compatible PgSTAC migration release only after checking the pinned
API dependency and server extension availability. Test the upstream migration on
the isolated restore, and record all created or altered objects and roles. PostGIS
is enabled per database, while [PostgreSQL roles](https://www.postgresql.org/docs/17/database-roles.html)
are shared across the instance;
role changes need a complete impact review. Do not change the shared `geo` role's
password or global search path as a shortcut.

## Live initialization approval

Present the exact proposed migration, verified recovery procedure, isolated test
results, expected locks/downtime and shared-object impact to Zeus. Request approval
immediately before the live persistent change. The diagnostic PR and the general
migration instruction do not authorize that change. Do not place an initialization
Job in an auto-synced path before approval.

## Acceptance and cutover

- A fresh backup restores successfully into the isolated system with expected
  metadata and extension versions. Record actual evidence, not manifest intent.
- The pinned API lists collections and executes a limited geographic search against
  the test database. An empty but valid catalog should return valid empty results;
  that does not establish that imagery has been ingested or recovered.
- Any fixture ingestion runs only against the isolated test catalog, with its
  records identified and removed as part of the test. Do not insert synthetic rows
  into the live catalog to make a gate appear to pass.
- After an approved live repair, verify the API and the OpenChoreo browser end to
  end, including catalog links, collection access, search and browser behavior.
- Keep the original browser Deployment/Service until cutover is verified. Record
  the real maintenance comparison before claiming phase 3 complete or expanding
  the migration.

## Primary references

- [Pinned API source](https://github.com/stac-utils/stac-fastapi-pgstac/tree/725dd57b0f56b83d063a757e90634492ec283a8a)
- [PgSTAC requirements and roles](https://stac-utils.github.io/pgstac/pgstac/)
- [PostgreSQL backup and restore](https://www.postgresql.org/docs/17/backup.html)
- [Repository backup image constraints](../../platform/backups/image/README.md)
