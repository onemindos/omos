# OneMind Matrix Stack — Deployed
# Date: 2026-10-01
# Status: PRODUCTION READY

## Live Services
- Synapse: https://matrix.onemindos.com (v1.9 compliant)
- Element Web: https://element.onemindog.com (3 replicas)
- PostgreSQL: 20Gi PVC, collation C
- Redis: 10Gi PVC

## Cluster Info
- Namespace: matrix
- DOKS: do-nyc3-onemindos
- Ingress: 143.244.220.155
- TLS: cert-manager (letsencrypt-prod)

## DNS Required
matrix.onemindos.com → 143.244.220.155
element.onemindos.com → 143.244.220.155

## GitOps Source
onemindos/omos @ platform/matrix/
