FROM ghcr.io/cloudnative-pg/postgresql:16.15-202610050823-standard-bookworm@sha256:370eae69bfd20329909b1510671a201e566b89d8f9880d09ae4e8f3d9bac64d0 AS age-build
USER root
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential postgresql-server-dev-16 libreadline-dev zlib1g-dev flex bison
ADD --checksum=sha256:46a280ab96a2d3dab0752527f429c512a4f581579425a7a52d8a7946876e661e \
    https://codeload.github.com/apache/age/tar.gz/2db2f060c4c9265a14d40f007eb8c56febf31e4c /tmp/age.tar.gz
RUN mkdir /tmp/age && tar -xzf /tmp/age.tar.gz -C /tmp/age --strip-components=1 \
    && make -C /tmp/age PG_CONFIG=/usr/lib/postgresql/16/bin/pg_config \
    && make -C /tmp/age PG_CONFIG=/usr/lib/postgresql/16/bin/pg_config install

FROM ghcr.io/cloudnative-pg/postgresql:16.15-202610050823-standard-bookworm@sha256:370eae69bfd20329909b1510671a201e566b89d8f9880d09ae4e8f3d9bac64d0
USER root
RUN apt-get update && apt-get install -y --no-install-recommends \
    postgresql-16-postgis-3 postgresql-16-pgvector && rm -rf /var/lib/apt/lists/*
COPY --from=age-build /usr/lib/postgresql/16/lib/age.so /usr/lib/postgresql/16/lib/
COPY --from=age-build /usr/share/postgresql/16/extension/age* /usr/share/postgresql/16/extension/
COPY extension-check.sql /extension-check.sql
USER postgres
