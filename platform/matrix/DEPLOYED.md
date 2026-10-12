# OneMind Matrix Stack — Deployed
# Date: 2026-10-01
# Status: PRODUCTION READY

## Live Services
- Synapse: https://matrix.onemindos.dev (v1.9 compliant)
- Element Web: https://element.onemindos.dev (3 replicas)
- PostgreSQL: 20Gi PVC, collation C
- Redis: 10Gi PVC

## Cluster Info
- Namespace: matrix
- DOKS: do-nyc3-onemindos
- Ingress: 143.244.220.155
- TLS: Cloudflare edge (proxied, zone SSL mode "full"); no origin cert, same as the other onemindos.dev ingresses

## Moved to onemindos.dev (2026-10-12)
server_name is `onemindos.dev` (user IDs: `@zeus:onemindos.dev`, `@legacy:onemindos.dev`). Synapse runs at
`matrix.onemindos.dev`; `onemindos.dev/.well-known/matrix/client` (Ingress/matrix-wellknown) points clients there.
The database was reset for the server_name change (pre-move dump taken; old: 2 users, 1 room). Federation stays off.
User passwords: `kubectl -n matrix get secret matrix-user-credentials`.

## DNS Required (Cloudflare-proxied)
matrix.onemindos.dev → 143.244.220.155
element.onemindos.dev → 143.244.220.155

## GitOps Source
onemindos/omos @ platform/matrix/
