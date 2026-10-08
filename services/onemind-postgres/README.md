# OneMind operational PostgreSQL image

CloudNativePG PostgreSQL 16.15 base is pinned by digest. PostGIS and pgvector use
the base image's PostgreSQL package repository. Apache AGE uses the pinned
PG16-compatible 1.6.0 release-candidate source commit, verified by archive checksum.
This is a prototype extension build and does not modify existing application DBs.

CI starts an isolated PostgreSQL instance and exercises spatial, vector and graph
queries before accepting the image. Production values must use its published
image digest. Application migration and extension schema changes belong to the
owning service. Installing extensions does not populate external source records.
