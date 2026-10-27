# Authentik

The identity provider for the portal, Argo CD, FedHub and the other OIDC applications, in namespace `auth`: two server replicas, one worker, Redis and PostgreSQL 17. Private access is `https://onemind-auth.tail717ba0.ts.net/`. `auth.onemindos.dev` is public through ingress-nginx, restricted to Cloudflare source ranges, for sign-in and callbacks off the tailnet.

Sources: `platform/argocd/platform-auth-helm.yaml`, `platform/auth/`, `docs/networking/NETWORKING.md` (section 2), `platform/backups/volsync/README.md`, `platform/backups/restore-tests/README.md`.

## Ownership

- `platform-auth-helm` renders the pinned `goauthentik/authentik` chart (2026.8.2) with inline values. Every credential comes from the SealedSecret `authentik-secret`.
- `platform-auth-infra` applies `platform/auth`: the `authentik-postgres` StatefulSet (5 GiB claim), the `authentik-redis` Deployment, the server PodDisruptionBudget, the sealed secret and the public Ingress.
- Applications register their OIDC providers through Git-owned resources in `platform/ecosystem/resources/`, for example `argocd-authentik-provider.yaml`.

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud -n argocd get application platform-auth-helm platform-auth-infra
kubectl --context do-nyc3-omos-cloud -n auth get deployments,statefulsets,pods,pvc,ingresses
kubectl --context do-nyc3-omos-cloud -n auth get cronjob backup-postgres-auth
```

Distinguish a hostname or connectivity problem, a login failure and an authorization denial before changing anything; the login, callback and subject-mapping checks are in [Identity and access](../identity.md). Do not bypass email verification, enable guest access or map arbitrary users to administrators.

## Backups and recovery

- `backup-postgres-auth` (namespace `auth`, daily at 05:00 UTC) writes an encrypted logical dump of the Authentik database and restores it into an isolated PostgreSQL 17 sidecar before marking it verified.
- VolSync copies the PostgreSQL volume. Its export has been restored to a new claim, where PostgreSQL booted with 230 application tables. A table inventory proves the database boots; it does not prove login or session recovery.
- Preserve the database storage and `authentik-secret`. Restore to a new claim, never over the running database. A hostname change needs coordinated routing, base URL and redirect URI changes; do not rename partially.

## Verify completion

Complete a real login and confirm the intended identity and permissions; readiness alone is not enough. Check that private HTTPS certificates validate and that anonymous API access is denied.

Generated workload list: [Authentik app runbook](../apps/authentik.md).
