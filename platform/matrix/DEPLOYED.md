# OneMind Matrix Stack — Deployed
# Date: 2026-10-01
# Status: PRODUCTION READY

## Live Services
- Synapse: https://matrix.tail717ba0.ts.net (v1.9 compliant)
- Element Web: https://element.tail717ba0.ts.net (3 replicas)
- PostgreSQL: 20Gi PVC, collation C
- Redis: 10Gi PVC

## Cluster Info
- Namespace: matrix
- DOKS: do-nyc3-onemindos
- Ingress: Tailscale operator (class `tailscale`), tailnet only; no public listener
- TLS: Tailscale-issued certificate

## Moved to onemindos.dev (2026-10-12)
server_name is `onemindos.dev` (user IDs: `@zeus:onemindos.dev`, `@legacy:onemindos.dev`). Synapse runs at
`matrix.onemindos.dev` (historical, retired 2026-10-29, see below).
The database was reset for the server_name change (pre-move dump taken; old: 2 users, 1 room). Federation stays off.
User passwords: `kubectl -n matrix get secret matrix-user-credentials`.

## Tailnet only (2026-10-29)
Matrix left the public internet. `matrix.onemindos.dev`, `element.onemindos.dev` and `onemindos.dev/.well-known/matrix` are gone;
clients use `https://matrix.tail717ba0.ts.net` and `https://element.tail717ba0.ts.net`. `server_name` is still `onemindos.dev`.
Federation is disabled (empty `federation_domain_whitelist`, no federation listener). In-cluster clients use
`http://synapse.matrix.svc.cluster.local:8008`. The Cloudflare records for the .dev hosts should be removed.

## GitOps Source
onemindos/omos @ platform/matrix/
