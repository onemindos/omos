# OneMind operational PostgreSQL image

CloudNativePG PostgreSQL 17.11 (standard bookworm) base is pinned by tag and digest. PostGIS and
pgvector use the base image's PostgreSQL package repository. TimescaleDB 2.26.4 comes from
Timescale's apt repository (key pinned by checksum, package versions pinned). Apache AGE uses the
PG17 1.6.0 release-candidate source tag (`PG17/v1.6.0-rc0`), verified by archive checksum.
The PG16 image (digest in git history) stays available for rollback.

## In-place upgrade notes (16 to 17, 2026-11-08)

`onemind-data` was upgraded with CNPG's offline `pg_upgrade`. Two things blocked it and must be
handled again on the next major upgrade:

- **Apache AGE**: its catalog tables use `regclass`/`regnamespace` columns, which `pg_upgrade`
  refuses. AGE had no graphs, so `DROP EXTENSION age` ran in `onemind` before the upgrade and
  `CREATE EXTENSION age` after. If graphs exist by then, use a logical dump and restore instead.
- **pgstac**: generated columns call functions that use PostGIS without a schema. PostgreSQL 17
  inlines them under an empty `search_path` during restore, so `pg_upgrade` failed. The five
  functions (`collection_geom`, `collection_datetime`, `collection_enddatetime`,
  `collection_base_item`, `search_hash`) now carry `SET search_path = pgstac, public`. A pgstac
  migration can reset this; check it before any future upgrade.

Test first: a schema-only `pg_dump` of every database restored into the new image (scratch pod,
`ON_ERROR_STOP=1`) reproduces both failures in minutes. TimescaleDB is created per database
(`onemind` has it). Barman writes the PG17 archive under server name `onemind-data-pg17`.

CI starts an isolated PostgreSQL instance (with `timescaledb` preloaded) and exercises spatial,
vector, graph and hypertable queries before accepting the image. Production values must use its
published image digest. Application migration and extension schema changes belong to the owning
service. Installing extensions does not populate external source records.
