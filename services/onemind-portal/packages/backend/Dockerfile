# syntax=docker/dockerfile:1
# One production build definition for GitHub Actions, the operator Mac and Linux builders.
# Compile JavaScript on the builder architecture; install runtime native modules on the target.
# Never copy build-stage node_modules into the production image.
FROM --platform=$BUILDPLATFORM node:22.23.3-bookworm-slim AS build
USER root
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && apt-get install -y --no-install-recommends \
      python3 python3-venv libsqlite3-dev make g++ git && \
    rm -rf /var/lib/apt/lists/*
RUN corepack enable && python3 -m venv /opt/techdocs && \
    /opt/techdocs/bin/pip install --no-cache-dir mkdocs-techdocs-core==1.7.1
ENV PATH="/opt/techdocs/bin:${PATH}" \
    NODE_OPTIONS="--max-old-space-size=4096" \
    AUTH_SESSION_SECRET="build-only-not-a-runtime-secret" \
    AUTH_OIDC_CLIENT_ID="build-only" \
    AUTH_OIDC_CLIENT_SECRET="build-only"
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/root/.yarn/berry/cache,sharing=locked \
    yarn install --immutable
RUN yarn tsc
RUN yarn test --watch=false --runInBand plugins/plugin-onemind-permissions-backend/src plugins/plugin-onemind-auth/src/scopedFetch.test.ts plugins/plugin-onemind-auth-backend/src/subject.test.ts
RUN bash scripts/build-techdocs.sh
RUN yarn build:backend

FROM node:22.23.3-bookworm-slim AS runtime

# Install sqlite3 dependencies. You can skip this if you don't use sqlite3 in the image,
# in which case you should also move better-sqlite3 to "devDependencies" in package.json.
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends libsqlite3-dev python3 make g++ && \
    rm -rf /var/lib/apt/lists/*

# From here on we use the least-privileged `node` user to run the backend.
USER node

# This should create the app dir as `node`.
# If it is instead created as `root` then the `tar` command below will fail: `can't create directory 'packages/': Permission denied`.
# If this occurs, then ensure BuildKit is enabled (`DOCKER_BUILDKIT=1`) so the app dir is correctly created as `node`.
WORKDIR /app
ENV NODE_OPTIONS="--no-node-snapshot --max-old-space-size=1536"

# Copy files needed by Yarn
COPY --from=build --chown=node:node /app/.yarn ./.yarn
COPY --from=build --chown=node:node /app/.yarnrc.yml ./
COPY --from=build --chown=node:node /app/backstage.json ./
# Production dependencies must be linked from this stage, not the builder state.
RUN rm -f .yarn/install-state.gz

# This switches many Node.js dependencies to production mode.
ENV NODE_ENV=production

# Copy repo skeleton first, to avoid unnecessary docker cache invalidation.
# The skeleton contains the package.json of each package in the monorepo,
# and along with yarn.lock and the root package.json, that's enough to run yarn install.
COPY --from=build --chown=node:node /app/yarn.lock /app/package.json /app/packages/backend/dist/skeleton.tar.gz ./
RUN tar xzf skeleton.tar.gz && rm skeleton.tar.gz

RUN --mount=type=cache,target=/home/node/.cache/yarn,sharing=locked,uid=1000,gid=1000 \
    yarn workspaces focus --all --production && yarn cache clean

# This will include the examples, if you don't need these simply remove this line
COPY --from=build --chown=node:node /app/templates ./templates
COPY --from=build --chown=node:node /app/catalog ./catalog
COPY --from=build --chown=node:node /app/techdocs ./techdocs
COPY --from=build --chown=node:node /app/runbooks ./runbooks

# Then copy the rest of the backend bundle, along with any other files we might want.
COPY --from=build --chown=node:node /app/packages/backend/dist/bundle.tar.gz /app/app-config*.yaml ./
RUN tar xzf bundle.tar.gz && rm bundle.tar.gz
# Catch the production-only missing-JSON regression and incomplete TechDocs before push.
RUN node -e "for (const name of ['@internal/plugin-onemind-common','@internal/plugin-onemind-common-backend']) require(name+'/directory.json'); const Database=require('better-sqlite3'); const db=new Database(':memory:'); db.prepare('SELECT 1').get(); db.close(); const fs=require('node:fs'); for (const p of ['techdocs/default/component/onemind-runbooks/techdocs_metadata.json','techdocs/default/component/onemind-runbooks/search/search_index.json']) { if (!fs.statSync(p).size) throw new Error('Empty release artifact: '+p); }"

CMD ["node", "packages/backend", "--config", "app-config.yaml", "--config", "app-config.production.yaml"]
