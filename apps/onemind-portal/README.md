# OneMind portal

Backstage 1.55.3, generated with create-app 0.9.2 using the new frontend system.

## Build

Use Node 22.23.3 and the committed Yarn 4.13.0 release. Run `corepack enable`,
`yarn install --immutable`, `yarn tsc`, and `yarn build:backend`.
The GitHub workflow performs the host build and publishes a commit-tagged image.
Runtime configuration and secrets are injected by Kubernetes, not built into the UI.

`@yarnpkg/core` is pinned to 4.9.1 because published 4.9.2 references a missing
repository-local got patch. This was reproduced during initial installation.
Backstage dependencies are aligned to its official 1.55.3 release manifest.

## Catalog

`node scripts/migrate-homepage.cjs` imports existing services and generates the directory.
`node scripts/snapshot-cluster.cjs` imports workload identities from the verified context.
These are explicit inventory snapshots. Live readiness is queried separately; new workloads
need a refreshed inventory until a periodic discovery provider is installed.

## Authentication and permissions

OIDC is the only configured sign-in provider. The initial catalog identity matches the
existing Authentik `akadmin` record; its existing profile email is preserved, not assumed
to be the owner's personal address. The Authentik application is bound to that user.
Add named users to both Authentik application bindings and the catalog deliberately.
`onemind-admins` can use configured actions. Other authenticated catalog users have read
permissions. Guest auth and the default allow-all policy have been removed.

Catalog MCP actions are installed; unrestricted OAuth client registration and scaffolder
MCP actions are disabled. No cloud mutation credentials are installed in the portal.

## Sources

- https://backstage.io/docs/getting-started/
- https://backstage.io/docs/deployment/docker/
- https://backstage.io/docs/deployment/k8s/
- https://backstage.io/docs/auth/oidc/
- https://backstage.io/docs/features/kubernetes/configuration/
- https://backstage.io/docs/permissions/getting-started/

See `../../platform/backstage/README.md` for rollout and recovery.
