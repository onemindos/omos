# Phase 7 brief: migrate Layer 3 apps

Status when written (2026-10-18): not started. Only after phase 3 says yes. Tracker: [`PLATFORM-ROADMAP.md`](../../PLATFORM-ROADMAP.md). Read [`HANDOFF.md`](HANDOFF.md) first.

## Goal

Move our own applications (Layer 3) into OpenChoreo one project at a time, and model shared
infrastructure (NATS streams, Postgres databases) as OpenChoreo resource types so components declare
their dependencies. Foundation (Layer 1) stays on plain Argo CD.

## Order and scope

| Project | Namespaces and workloads | Notes |
|---|---|---|
| `omos-feeds` | `omos`: feed-intel, feed-satellites, omos-tak-projector, redpanda-connect, sink-clickhouse | needs NATS and ClickHouse |
| `onetak` | `onetak`: the 19 geo and mapping apps (Valhalla, Martin, Nominatim, WebODM, STAC, TileServer, Tile38, Overpass, ...) | databases (PostGIS, webodm-db, MinIO) are Layer 2, keep them as Argo-managed services |
| `takml` | `takml`: takml-server, takml-vision-worker | `takml-db` is a database; `takml-models` and `takml-server-jar` are ReadWriteOnce volumes shared by server and worker: both pods must share a node (podAffinity, see `platform/takml/deployment.yaml`) and `Recreate` |
| `onemind` | `onemind-core`: onemind-core, onemind-ingestion, OpenFGA (OpenFGA is security: keep it Layer 1/2) | `onemind-jobs`, `onemind-media` are empty |
| `agent-fabric` | `agent-console`, `agent-tap` | |
| `n8n` | `n8n` | its Postgres stays |
| `matrix` | Synapse, Element | Argo app `platform-matrix` is OutOfSync; fix it first |
| `openclaw` | OpenClaw | the Deployment is not in Git today; put it in Git first |
| `hermes` | `hermes` | last; it does the fixing. Phase 8 gives it a model |

Not moved: `takserver`, `nats`, `auth`, `tailscale`, `monitoring`, `velero`, `kyverno`, `falco`, `spire-system`,
`argocd`, `openchoreo-*`, `backstage`, `kube`, `external-dns`, `cloudflare`, `ingress-nginx`, `cert-manager`.

## Per-project procedure (one PR per project)

1. Read every manifest of the project; list volumes, secrets, env, ingress, cross-namespace calls,
   NetworkPolicies, ServiceAccounts/RBAC, jobs and cronjobs.
2. Write Components/Workloads under `platform/openchoreo/projects/<project>/` (layout from phase 3).
   Use `onemind-service`; add component types (`worker`, `scheduled-task` exist upstream) or traits
   (`persistent-volume-claim`, `health-probes`, `horizontal-pod-autoscaler`) where needed. Set
   `strategy: Recreate` for ReadWriteOnce volumes.
3. Secrets: SealedSecrets stay; reference them through OpenChoreo `SecretReference`/ExternalSecret
   only when the component needs per-environment values. Warning: the `default` ClusterSecretStore
   exposes the whole `openchoreo-secrets` namespace to any ExternalSecret in the cluster; do not store
   platform-internal secrets there.
4. Run the new component beside the old one if possible, switch consumers (Service names, ingress), watch
   it for a day, then delete the old manifests in a follow-up PR. Keep rollback simple: the old
   manifests stay in Git history.
5. Stateful workloads: migrate data deliberately (Velero, VolSync and PG-dump backups exist; take one
   first and verify the restore path).

## Resource types

- NATS: model streams and KV as an OpenChoreo `ClusterResourceType` backed by the existing
  `jetstream.nats.io` Stream/KeyValue CRs (`platform/nats/streams/`). The upstream defaults include a `nats`
  resource type to start from.
- Postgres: a resource type over CNPG (`postgresql.cnpg.io`), as `postgres` in the upstream defaults.
- Give `matrix` and `openclaw` a GitOps owner as part of their move.

## Gate (per project)

Healthy in OpenChoreo and visible in the portal, old Argo resources removed, backups still cover its data,
no new Pending pods or capacity regressions.

## Risks

- Capacity (phase 5 first). Cross-namespace DNS names change when a workload moves to `dp-*` namespaces.
- Tailnet-only app hostnames (`*.apps.openchoreo.<ip>.nip.io`): public apps (`tak.onemindos.dev`, portals) keep their
  ingress until a real domain and TLS (`letsencrypt-dns`) are set up for the data-plane gateway.
- Do not touch TAK Server or NATS without asking Zeus.
