# OneMind operational PostgreSQL image

CloudNativePG PostgreSQL 17.11 (standard bookworm) base is pinned by tag and digest. PostGIS and
pgvector use the base image's PostgreSQL package repository. TimescaleDB 2.26.4 comes from
Timescale's apt repository (key pinned by checksum, package versions pinned). Apache AGE uses the
PG17 1.6.0 release-candidate source tag (`PG17/v1.6.0-rc0`), verified by archive checksum.
The PG16 image (digest in git history) stays available for rollback.

CI starts an isolated PostgreSQL instance (with `timescaledb` preloaded) and exercises spatial,
vector, graph and hypertable queries before accepting the image. Production values must use its
published image digest. Application migration and extension schema changes belong to the owning
service. Installing extensions does not populate external source records.
