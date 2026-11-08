# PostgreSQL 17 operational image (CNPG standard bookworm base, 17.11 built 2026-10-19).
FROM ghcr.io/cloudnative-pg/postgresql:17.11-202610190957-standard-bookworm@sha256:0d0e407ca1d986022651671e830022ccfe496994a8b99fe925b5ae3f0264cb1a AS age-build
USER root
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential postgresql-server-dev-17 libreadline-dev zlib1g-dev flex bison
# Apache AGE 1.6.0-rc0 for PostgreSQL 17 (tag PG17/v1.6.0-rc0). Same extension version as the PG16 image,
# so existing databases upgrade in place without ALTER EXTENSION age UPDATE.
ADD --checksum=sha256:bc9b013b3449dfa99fb665282f0f0151084f5252ecc74943f8ec20f0a6b832b4 \
    https://codeload.github.com/apache/age/tar.gz/54905a09bf8462f22a87c3adfd2ab5752e5c1e71 /tmp/age.tar.gz
RUN mkdir /tmp/age && tar -xzf /tmp/age.tar.gz -C /tmp/age --strip-components=1 \
    && make -C /tmp/age PG_CONFIG=/usr/lib/postgresql/17/bin/pg_config \
    && make -C /tmp/age PG_CONFIG=/usr/lib/postgresql/17/bin/pg_config install

FROM ghcr.io/cloudnative-pg/postgresql:17.11-202610190957-standard-bookworm@sha256:0d0e407ca1d986022651671e830022ccfe496994a8b99fe925b5ae3f0264cb1a
USER root
# TimescaleDB 2.26.4 (the version OpenRemote runs) from Timescale's apt repository; the signing key is checksum-pinned.
ADD --checksum=sha256:eb53bd3dc31e06fbf1ed769139f12b1fffaf1958a96c9b4f9727c19aada83321 \
    https://packagecloud.io/timescale/timescaledb/gpgkey /tmp/timescaledb.asc
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates gnupg \
    && gpg --dearmor -o /usr/share/keyrings/timescaledb.gpg /tmp/timescaledb.asc \
    && echo "deb [signed-by=/usr/share/keyrings/timescaledb.gpg] https://packagecloud.io/timescale/timescaledb/debian/ bookworm main" \
        > /etc/apt/sources.list.d/timescaledb.list \
    && apt-get update && apt-get install -y --no-install-recommends \
    postgresql-17-postgis-3 postgresql-17-pgvector \
    timescaledb-2-postgresql-17=2.26.4~debian12-1709 timescaledb-2-loader-postgresql-17=2.26.4~debian12-1709 \
    && rm -rf /var/lib/apt/lists/* /tmp/timescaledb.asc
COPY --from=age-build /usr/lib/postgresql/17/lib/age.so /usr/lib/postgresql/17/lib/
COPY --from=age-build /usr/share/postgresql/17/extension/age* /usr/share/postgresql/17/extension/
COPY extension-check.sql /extension-check.sql
USER postgres
