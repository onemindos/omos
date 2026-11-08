# OneMind stack inventory

> **Matrix retired 2026-11-07** (Zeus: 2 users, 0 rooms, no events). Synapse, Element Web, PostgreSQL and Redis in namespace `matrix` were removed. Its manifests moved to `platform/parked/matrix/` (restore steps in its README), the Application `platform-matrix` was removed, and the VolSync sources, the `backup-postgres-matrix` CronJob and the matrix backup credentials were taken out of `platform/backups/`. A final encrypted `pg_dump` and the Synapse signing key are in the backup bucket under `backups/matrix/final-2026-11-07/`; the three VolSync restic repositories stay under `backups/restic/matrix/`. Earlier snapshots in this file that mention Matrix are historical.

> **Kubernetes Dashboard retired 2026-10-27.** `kubernetes-dashboard` v2.7.0 in `kube` ran with `--enable-insecure-login --enable-skip-login` under a cluster-admin ServiceAccount on the tailnet (`k8s-dashboard`); Headlamp covers the use case. Its manifests and catalog entries were removed through a reviewed Git change (`platform-kube` prunes them); the hand-made `ClusterRoleBinding/kubernetes-dashboard-admin` needs a manual delete. Earlier snapshots in this file that mention it are historical.

> **Dapr removed 2026-11-07.** The Dapr control plane (`ecosystem-dapr`, namespace `dapr-system`, 7 CRDs, the sidecar-injector webhook) and the onemind-core Dapr workflow are gone. onemind-core now runs missions in-process (a bridge until mission orchestration moves to Hermes), no sidecar, no durable replay. The `dapr_workflows` database and role in `onemind-data` are set to `ensure: absent` (entries removed in a follow-up PR). The three scheduler PersistentVolumes were Retain and are Released; the DigitalOcean volumes are kept for the owner to delete.

> **LiteLLM and in-cluster Ollama retired 2026-11-06.** Argo Applications `ecosystem-litellm` and `ecosystem-ollama`, their values, network policies, sealed secrets, the model-client bootstrap Job and the VolSync source for PVC `onemind-cpu-model` were removed from Git (agentgateway is the hub; Hermes uses Ollama Cloud). The `litellm` database and role in `onemind-data` are set to `ensure: absent` (CNPG drops them; entries removed in a follow-up PR), the `db-litellm` sealed secret was removed from the repo, and the `onemind-cpu-model` DigitalOcean volume is retained. See `docs/operations/tool-gateway.md`. Rows below that mention them are historical until the cluster is cleaned up.

> **OpenClaw retired 2026-10-22** (Zeus: not needed any more). Its Application, manifests, namespace, volume and backup source were removed through reviewed Git changes (PRs #379 and the retirement series). Earlier snapshots in this file that mention it are historical.

Tiered inventory of everything running on the DOKS cluster `omos-cloud` (nyc3), what owns it, what state it holds, whether it is backed up, who can reach it, and whether Git or a laptop controls it.

- **Captured:** 2026-10-09, read-only, kube context `do-nyc3-omos-cloud`, Argo CD synced to `main` at `4934f532`.
- **Scope:** all 47 Argo CD Applications (52 on 2026-10-17), all 11 Helm releases first installed by hand (4 adopted into Argo since, see GitOps gaps), and the workloads that nothing in Git or Helm tracks.
- **Status:** a point-in-time snapshot, not a live view. To refresh it, rerun the commands in [How to refresh](#how-to-refresh). Counts in the Summary go stale first.
- **Development phase:** every app is deliberate (system integration). This file ranks apps; it does not recommend removing any.
- **Owner:** there is one operator, `zeus`, for everything. The "Owner" column is therefore omitted from the tables; the first thing to change when a second person joins is to add it.

## Refresh 2026-11-07 (read-only)

Hand-edited; this file has no generator, so the sections below were updated one by one against the cluster.
Live on 2026-11-07: 64 Argo Applications (59 Synced and Healthy; `app-of-apps` and `platform-argocd` OutOfSync
and Healthy; manual-sync `platform-omos-entity`, `platform-policy` and `platform-takserver-tak-client`
OutOfSync and Missing), 4 nodes (3 `onemind-hub-core-v2`, 1 `onemind-ecosystem-cpu`, autoscale 1 to 2). The
Karakeep, OpenRemote, ThingsBoard, Forgejo, Gatus and Langfuse sections below describe workloads that are now
running; the Trivy operator is running too. Tempo is gone (removed 2026-11-03), so older rows that mention it are historical.
Helm release `nats` shows status `failed` at revision 11 (2026-11-07 00:05) while its three pods are Running; the
cause is not investigated here.

## Phase 5 addendum — 2026-10-19, 08:07 UTC

This narrow refresh does not recertify the older inventory below. Live cluster:
4 nodes, 70 Argo Applications (68 Synced/Healthy); only pre-existing Matrix
OutOfSync/Healthy and OpenClaw Synced/Progressing exceptions. No Pending pods.
All nodes were below 78% requested CPU/memory at that refresh. On 2026-10-27 hub-core memory requests were 79%, 96% and 71% (read-only `kubectl describe node`). OpenCost duplicate KSM emission is off.
OpenChoreo observability core/modules are deployed on existing backends, with
an authenticated tailnet Observer, real console/MCP pilot logs and metrics, and
verified encrypted SQLite backup/isolated restore. See the timestamped
[Phase 5 evidence](../PLATFORM-ROADMAP.md#phase-5-capacity-then-observability).
Older node counts, request percentages and M2 findings below are historical.

## Phase 6 addendum — 2026-10-19, 18:32 UTC

(Superseded: the OpenChoreo plugins were later removed from the portal; see
`platform/backstage/README.md`.) Existing OneMind Backstage 1.55.3 then hosted pinned
OpenChoreo 1.3.1 core and observability plugins, 2/2 ready and Argo Synced/Healthy. Real Authentik subject
mapping, existing home/Today/catalog/runbooks/settings, pilot Deploy/Definition,
logs/metrics and machine read-only boundaries passed. Shared permission-cache
results use the existing PostgreSQL service; no raw IDP token is stored.
Temporary console and candidate Deployments/routes are removed; their production
replacement is `onemind-preview`, not a second portal. No storage, identity or
backup was deleted. At acceptance: 69 Argo Applications, 68 Healthy; only existing
Matrix OutOfSync/Healthy and OpenClaw Synced/Progressing exceptions. This narrow
observation does not recertify historical inventory.
See [Phase 6 evidence](../PLATFORM-ROADMAP.md#phase-6-one-portal).

## Paperclip: retired — 2026-11-01

Owner decision (2026-11-01): remove Paperclip entirely. `platform-paperclip` and its
manifests are removed from the repository. Final backups taken 2026-11-01 before removal: a
`pg_dump -Fc` of the `paperclip` database (242 tables, 1.3 MB; held locally by the
operator, not committed or uploaded because it may contain user data) and DigitalOcean
snapshots of both volumes (`paperclip-postgres-final-2026-11-01`,
`paperclip-data-final-2026-11-01`, tag `paperclip-retirement`). The two retain-class
volumes are kept until the owner approves their deletion. Live teardown steps:
[retirement runbook](operations/paperclip-retirement-2026-11-01.md). The Hermes-side
pieces (`hermes-paperclip-api`) are retired with the old Hermes, separately.

## Karakeep — prepared 2026-10-24, running 2026-11-07

User-requested bookmark/read-later workspace. Prepared as source on 2026-10-24; deployed on the owner's
authorization of 2026-10-25 ([README](../platform/karakeep/README.md)). Verified live 2026-11-07: `platform-karakeep`
Synced and Healthy with automated sync (prune, selfHeal); web, Chrome and Meilisearch pods Running for 36 hours;
Tailscale ingress `karakeep.tail717ba0.ts.net`; both VolSync sources have synced today. The PVCs (`data-pvc` 20Gi,
`meilisearch-pvc` 5Gi) are on `do-block-storage`, not the retain class. A recorded isolated restore and the
mobile/extension compatibility checks are not recorded here, so treat the restore as unproven.

| Tier | Application | Namespace | State / recovery | Exposure | Status |
| --- | --- | --- | --- | --- | --- |
| 2 | `platform-karakeep` | `karakeep` | Native SQLite/assets 20Gi + Meilisearch 5Gi on `do-block-storage`; two VolSync sources syncing; [isolated restore](../platform/karakeep/BACKUP.md) not recorded | Private `https://karakeep.tail717ba0.ts.net`, app authentication | Auto-sync, Synced/Healthy, 3 pods Running (2026-11-07) |

## Policy service preparation — 2026-10-29

FEATURES I4 (onemind-xr): `svc.policy.*`, the only publisher of `cmd.*` after a human confirm and the single
writer of `evt.so.action.*`. Not deployed: the Application is manual and waits for the owner steps in
[`platform/policy/README.md`](../platform/policy/README.md) (NATS per-service users, OpenFGA model adoption,
two SealedSecrets, pinned image digest).

| Tier | Application | Namespace | State / recovery | Exposure | Status |
| --- | --- | --- | --- | --- | --- |
| 2 | `platform-policy` | `onemind-core` | No PVC; action records in memory (pending actions are forgotten on restart and cannot be confirmed); audit in `EVT` (`evt.so.>`), commands in `CMD` | None (NATS client only; no Service or ingress); ecosystem pool | Manual Argo, not synced; [owner steps](../platform/policy/README.md) pending |

## Langfuse — running 2026-11-07

Self-hosted Langfuse v3 for Hermes traces. Live 2026-11-07: `ecosystem-langfuse` Synced and Healthy, web, worker and Valkey pods Running (10 to 12 hours), Tailscale ingress `onemind-langfuse.tail717ba0.ts.net`. Source: chart `langfuse` 2.2.3 (app 4.75.0),
Application `ecosystem-langfuse`, namespace `onemind-ai`, web and worker one replica each on the ecosystem pool. It reuses the
database `langfuse` in CNPG `onemind-data`, a `langfuse` user and database in ClickHouse `clickhouse/onemind`, Spaces bucket
`onemind-backups` (prefix `langfuse/`), the bundled Valkey (no PVC) and the Authentik provider `onemind-langfuse`. Hermes is wired to it: the
`onemind-hermes` Deployment loads Secret `hermes-langfuse` and installs the pinned Langfuse plugin at start. Trace flow and the ClickHouse/CNPG recovery path are not verified here.

| Tier | Application | Namespace | State / recovery | Exposure | Status |
| --- | --- | --- | --- | --- | --- |
| 2 | `ecosystem-langfuse` | `onemind-ai` | PostgreSQL (CNPG `onemind-data`, `Velero`, `Barman`), ClickHouse database `langfuse` (VolSync of the ClickHouse PVC), event blobs in Spaces; Valkey queue is disposable | Tailnet only: `https://onemind-langfuse.tail717ba0.ts.net`, Authentik sign-in only | Auto-sync, Synced/Healthy (2026-11-07) |

## Summary

Trivy pilot, prepared 2026-10-24 and retargeted 2026-11-04: a bounded Argo-managed
image-vulnerability pilot for `auth,nats,onetak,n8n` under `platform/security/`, with
one concurrent scan Job, namespace-scoped operator permissions and a Grafana dashboard.
Owner approved deploying it 2026-11-04; `platform-security` is Synced and Healthy and `trivy-operator` is Running in `trivy-system` with scan Jobs completing (2026-11-07).
Acceptance is not recorded here, and the private `ghcr.io/onemindos` image
(`tak-portal`) is not scanned. No new PVC, node or public
route is requested. See [capacity, scope and acceptance](operations/trivy-image-scan-pilot.md).

| Item | Count / value |
|------|---------------|
| Argo CD Applications | 47 (45 Synced and Healthy; `app-of-apps` and `platform-ecosystem-resources` OutOfSync) |
| Helm-CLI releases (not in Git as Applications) | 3 on 2026-10-17: `nats`, `takserver`, `argocd` (11 at capture; `cert-manager`, `ingress-nginx`, `nack`, `tailscale-operator` adopted in PR #176, the four monitoring releases on 2026-10-17) |
| Workloads that neither Argo nor Helm tracks | DOKS system and Tailscale operator proxies. The NATS web UI (NUI) is Argo-managed by `platform-nats-tools` |
| Namespaces | 48 |
| Nodes | 5 at capture: 3 x `onemind-hub-core-v2` (8 vCPU / 16 GiB), 2 x `onemind-ecosystem-cpu` (8 vCPU / 32 GiB). 4 on 2026-10-27: the ecosystem pool (autoscale 1–4) was at 1 node |
| Ingress classes | `nginx` (public, LB `143.244.220.155`), `tailscale` (tailnet only) |
| Other external listeners | `tak-tailnet` Tailscale LB (TAK, 8089/8443/8446/64738), `omos-nats` Tailscale LB |
| Container references on mutable tags (`:latest`, bare `busybox`, `run-latest`) | 31 at capture; 0 running on 2026-10-27 (see [Mutable image tags](#mutable-image-tags)) |

## Tier definitions

| Tier | Meaning | What it implies |
|------|---------|-----------------|
| **1** | If it is down or its data is lost, the platform cannot operate or cannot be recovered | Change behind a backup check; no auto-merge of dependency bumps; needs a restore test on record |
| **2** | Everything else: integration and ecosystem work in progress | Fix forward; reasonable to break for a day |

Tier 1 follows the owner's list: TAK Server, NATS, Authentik, ingress, cert-manager, sealed-secrets, backups and ops, Argo CD, Hermes. **Tailscale operator is marked tier 1 here as a proposal** because every admin UI is reachable only through it; change it back if you disagree.

### Column legend

- **Backup:** `PG dump` = scheduled CronJob that uploads an encrypted archive to Spaces and restore-verifies it. `VolSync` = restic copy of the PVC to Spaces. `Velero` = daily namespace backup (06:00 UTC, 720 h retention) for `onemind-core`, `onemind-data`, `onemind-ai`. `Barman` = CNPG continuous backup. `Git` = stateless, rebuilt from this repo. All of these land in DigitalOcean Spaces `nyc3`, the same provider and region as the cluster (report finding A18).
- **Exposure:** `public` = reachable from the internet. `public (CF)` = nginx ingress restricted to Cloudflare source ranges. `public (open)` = nginx ingress with no source restriction. `tailnet` = Tailscale ingress or LB only. `internal` = ClusterIP only.
- **GitOps:** `Argo` = an Argo Application reconciles it from this repo. `Helm CLI` = installed with `helm upgrade` from a laptop (see [GitOps gaps](#gitops-gaps)). `none` = applied by hand, no owner in Git or Helm.

## Tier 1: critical

### Argo CD Applications

| Application | Namespace | Purpose | State / data | Backup | Exposure | GitOps |
|-------------|-----------|---------|--------------|--------|----------|--------|
| `app-of-apps` | argocd | Root app; creates every `platform-*` Application from `platform/argocd` | none | Git | tailnet (Argo UI) | Argo. **OutOfSync** (not yet diagnosed) |
| `platform-core` | kube-system, github-runner, tailscale | Sealed Secrets controller, self-hosted GitHub runner, Tailscale ACL policy | Sealed Secrets private keys (irreplaceable) | `backup-cluster-recovery` (daily 04:00 UTC) exports identity/config and the active sealing keys, encrypted | internal | Argo |
| `platform-takserver` | takserver | TAK Server overlay: LB proxy, NATS leaf, ConfigMaps, Mumble, Ignite, CoT retention. **TAK core is the Helm release `takserver`** | PostgreSQL/PostGIS, Ignite, 37 GiB of PVCs | `PG dump` (02:00 UTC, 86 min), `VolSync` x5 | **tailnet only**: `tak-tailnet` (`tak.tail717ba0.ts.net`) ports 8089, 8443, 8446, 64738/TCP. The public LB `tak-external` was removed 2026-10-23 | Argo for the overlay, **Helm CLI for the core** |
| `platform-nats` | nats | NATS 3-node JetStream cluster manifests, TLS, 39 Stream/KeyValue CRs, WebSocket ingress, Tailscale LB `omos-nats` | JetStream streams, 150 GiB of PVCs | `backup-nats-streams` stream archive (03:30 UTC, restore-verified, 39 streams), `VolSync` x3 | `ws.onemindos.dev` public (CF); `omos-nats` tailnet; 4222 internal | Argo **and** Helm CLI (two owners; the StatefulSet is tracked by both) |
| `platform-auth-helm` | auth | Authentik IdP (chart 2026.8.2) | PostgreSQL 17, Redis, 5 GiB | `PG dump` (05:00 UTC), `VolSync` | `auth.onemindos.dev` public (CF); `onemind-auth` tailnet | Argo |
| `platform-auth-infra` | auth | Authentik Postgres, Redis, secrets, ingress | see above | see above | see above | Argo |
| `platform-backups` | ops, plus Jobs in auth, backstage, matrix, n8n, nats, takml, takserver | Logical PostgreSQL/NATS backup CronJobs and encrypted cluster recovery exports | backup credentials (SealedSecrets) | self (writes to Spaces); drills in `docs/operations/` | internal | Argo. Shares `SealedSecret auth/ghcr` with `platform-ecosystem-resources` |
| `platform-volume-backups` | ops | 46 VolSync `ReplicationSource` objects, all `Successful` at last check | restic repositories in Spaces | n/a | internal | Argo |
| `platform-volsync` | ops | VolSync operator (chart 0.16.0) | none | n/a | internal | Argo |
| `platform-ops` | ops | `health-keeper` read-only cluster observer (every 5 minutes). Logs only, it does not alert | none | Git | internal | Argo |
| `ecosystem-velero` | velero | Velero (chart 12.2.1) and the `velero-onemind-platform` schedule | backups in Spaces | n/a | internal | Argo |
| `platform-hermes-base` | hermes | Namespace `hermes`, the Tailscale egress Services (`omos-ha-edge`, `omos-mbp`, `omos-s24`) and the Secrets `onemind-hermes` mounts (`hermes-dashboard-auth`, `omos-nats-creds`, `hermes-nats-ca`); `prune: false`, no finalizer, objects carry `Prune=false,Delete=false` (`platform/hermes/base/README.md`; added 2026-11-01) | none | n/a | none | Argo |
| `platform-onemind-mcp` | hermes | `onemind-mcp` tool server (the tool gateway's `onemind` target); split out of the retired `platform-legacy` on 2026-11-01, no finalizer | none | n/a | `onemind-mcp.hermes.svc` in-cluster only | Argo |
| `platform-onemind-hermes` | hermes | `onemind-hermes`: new Hermes v0.23.0 instance (persona Legacy, single soul), replaces the `hermes-fresh` launch; runtime acceptance in deployment PR. Holds the Legacy identity `agents.*.hermes.zeus.legacy` as a direct client of the bus (no sidecar; owner-reported 2026-11-03, it answered a real prompt). The old `hermes-legacy` instance and its Application `platform-legacy` were removed from Git on 2026-11-03 (`platform/hermes/onemind-hermes/RETIREMENT-AUDIT.md`); the old 25 GiB volume `hermes-data` is kept as a Retain PV until a separate delete. The `hermes-fresh` resources were removed on 2026-10-29 except the PVC `hermes-fresh-data`, kept pending owner approval (`platform/hermes/onemind-hermes/README.md`) | New 25 GiB `onemind-hermes-data` (`do-block-storage-retain`); no Legacy import | `VolSync` to its own Restic prefix; restore acceptance pending | `onemind-hermes` dashboard tailnet; Webapp (browser Desktop) `onemind-hermes-web` tailnet | Argo |

### Helm-CLI releases (tier 1)

| Release | Namespace | Chart | Purpose | State | Backup | Exposure | GitOps |
|---------|-----------|-------|---------|-------|--------|----------|--------|
| `takserver` | takserver | takserver-4.6.0 (rev 34, 2026-09-29), local chart `platform/takserver/chart` | TAK Server API, messaging, config, plugins, Ignite, PostgreSQL | see `platform-takserver` | see `platform-takserver` | see `platform-takserver` | **Helm CLI** |
| `nats` | nats | nats-2.14.2 (rev 10) | NATS StatefulSet | see `platform-nats` | see `platform-nats` | see `platform-nats` | **Helm CLI** (and Argo tracks the same StatefulSet) |
| `nack` | nats | nack-0.34.0 | JetStream controller that reconciles `platform/nats/streams/` | none | Git | internal | Argo (`edge-nack`, adopted in PR #176) |
| `argocd` | argocd | argo-cd-10.1.2 (app v3.4.4) | GitOps controller | CRDs and Application objects | Git, `backup-cluster-recovery` | tailnet (`argocd.tail717ba0.ts.net`) | **Helm CLI** today; `platform-argocd` (manual sync only, not yet synced) will adopt it after the pending #392 `helm upgrade`, see `docs/gitops/ARGOCD-SELF-MANAGEMENT.md`. Signs in through Authentik (`oidc.config`, group `argocd-admins` = admin) once `platform/values/argocd.yaml` is applied; local `admin` stays as break-glass until SSO is proven; Dex disabled |
| `cert-manager` | cert-manager | cert-manager-v1.20.3 | Certificate issuance | Certificate/Secret objects | `backup-cluster-recovery` | internal | Argo (`edge-cert-manager`, adopted in PR #176). No `ClusterIssuer` exists (the unused `letsencrypt-dns` and its Cloudflare token were removed); the only issuers are namespaced self-signed ones for NATS, CloudNativePG and OpenChoreo, owned by those components |
| `ingress-nginx` | ingress-nginx | ingress-nginx-4.15.1 (app 1.15.1) | Public ingress controller, LB `143.244.220.155`, 1 replica | none | Git | **public** | Argo (`edge-ingress-nginx`, adopted in PR #176) |
| `tailscale-operator` | tailscale | tailscale-operator-1.98.4 | Tailnet ingress, egress and LB proxies (about 24 proxy StatefulSets, operator-created) | proxy node keys | none | tailnet | Argo (`edge-tailscale-operator`, adopted in PR #176) |

## Tier 2: everything else

### Argo CD Applications, platform layer

| Application | Namespace | Purpose | State / data | Backup | Exposure | GitOps |
|-------------|-----------|---------|--------------|--------|----------|--------|
| `platform-fedhub` | takserver | TAK federation hub (`fedhub:5.8`) plus MongoDB, `tak-ldap` | MongoDB | `VolSync` | `fedhub.tail717ba0.ts.net:9100` / `:9443` **tailnet (L4)** | Argo. Health flaps to Progressing |
| `platform-backstage` | backstage | Backstage developer portal and service catalog, 2 replicas. Live since 2026-11-07: the live bridge (`livebridge` backend plugin: read-only NATS subscription plus Kubernetes watch, pushed to the UI as Signals; both replicas logged `livebridge NATS connected`), the Ops Center page `/ops` (CesiumJS globe of CoT tracks, bus stream, namespace health, Argo strip; [README](../services/onemind-portal/plugins/plugin-onemind-ops/README.md)), and the Alertmanager route `AlertmanagerConfig/onemind-backstage` (critical and warning into the `onemind-admins` inbox through `POST /api/events/http/alertmanager`, Falco events included because falcosidekick posts Warning and above to Alertmanager). MCP actions (`/api/mcp-actions/v1`, for operator clients, not Hermes): `catalog_*`, `onemind_health`, `onemind_runbook` and the read-only `livebridge.snapshot_summary` | PostgreSQL, 10 GiB | `PG dump` (06:40 UTC), `VolSync` | `portal` tailnet (`onemind-preview` redirects) | Argo |
| `platform-monitoring` | monitoring, nats | Grafana dashboards, rules, ServiceMonitors, `nats-surveyor`, `nats-leafz-exporter` (per-leaf leafnode metrics, added 2026-10-27 after the capture), `nats-metrics-scrape` (0 replicas). **Alertmanager receiver `onemind-backstage` (portal) added 2026-11-07**; before that its only receiver was `null`, so alerts went nowhere | none (stack lives in Helm release `monitoring`) | Git | `grafana` tailnet | Argo |
| `platform-observability` | clickhouse | ClickHouse "Time Machine" (Altinity CHI), server 26.8.36.3 LTS (upgraded from 24.8 on 2026-11-07 via 25.8 and 26.3; Langfuse v4 needs >= 25.12) | 20 GiB time-series | `VolSync`; native backup Job daily 03:00 UTC (`onemind` + `langfuse`) | internal | Argo |
| `platform-clickhouse-operator` | kube-system | Altinity ClickHouse operator (chart 0.27.3). A second, legacy `clickhouse-operator` Deployment sits at 0 replicas | none | Git | internal | Argo |
| `platform-takserver-tak-client` | takserver | `omos-tak-client` (A21, **proposed, not synced**): TAK Server streaming client bridge; TAK CoT to `tak.cot.out`, `tak.cot.inject` into TAK (queue group shared with NatsBridgePlugin). Replaces the never-deployed cot-tap | none | Git | internal (client of `takserver-messaging:8089` and `nats:4222`) | Argo, **manual sync**; needs two owner-sealed Secrets and a digest pin (`platform/takserver/tak-client/README.md`) |
| `platform-omos` | feeds | `feed-intel`, `feed-satellites`, `omos-tak-projector`, `sink-clickhouse`, `omos-mission-engine` (0 replicas) | none (streams in NATS, data in ClickHouse) | Git | internal | Argo |
| `platform-omos-entity` | feeds | **Proposed, not synced (manual sync).** `omos-entity-service` (FEATURES A4): fuses `COT`/`TEL`/`WORLD`/`INTEL` into `ent.{entity_type}.{entity_id}.state` in the proposed `ENT` stream; one replica on the ecosystem pool, KV lease in the proposed `entity-correlation` bucket | none on disk (state in `ENT`, rebuilt from it on start) | Git; state rebuildable by replay | internal (`svc.entity.get`/`query` on NATS) | Argo, manual sync |
| `platform-redpanda-connect` | feeds | Redpanda Connect (digest-pinned) | none | Git | internal | Argo |
| `platform-onetak` | onetak | Tak-portal (root host), MediaMTX, geo stack (Tile38), WebODM | about 150 GiB of PVCs | `VolSync` x5 | `onemindos.dev`, geo hosts **public (CF)**; `enroll.tail717ba0.ts.net` **tailnet** | Argo |
| `platform-takml` | takml | TAK-ML vision pipeline | PostgreSQL, 8 GiB | `PG dump` (06:00 UTC), `VolSync` | internal | Argo |
| `platform-n8n` | n8n | n8n workflow automation | PostgreSQL, 15 GiB | `PG dump` (05:20 UTC), `VolSync` | tailnet | Argo |
| `platform-kube` | kube | Headlamp | none | Git | tailnet | Argo |
| `platform-ecosystem-applications` | argocd | Generates the `ecosystem-*` Applications (11) from `platform/ecosystem/applications` | none | Git | internal | Argo |
| `platform-ecosystem-foundation` | argocd | Generates the foundation `ecosystem-*` Applications (10) | none | Git | internal | Argo |
| `platform-tool-gateway` | tool-gateway | Legacy's tool front door: agentgateway v1.7.0 standalone chart (`agentgateway-standalone`, no controller, no CRDs) on the ecosystem pool. Nine MCP routes (`/ops`, `/dev`, `/cloud`, `/data`, `/tak`, `/home`, `/work`, `/identity`, `/execute`) and six live targets: `onemind` and `kubernetes` (ops), `stac`, `core` and `clickhouse` (data), `mediamtx` (tak); [backends below](#tool-gateway-backends-tier-2). Per-target classification is the policy plane: destructive tools are hidden on every route and `/execute/mcp` is deny-all until Phase 4. **Hermes is not cut over to it yet** (Phase 5); LiteLLM still serves the MCP hub. Design: [`onemind-tool-gateway.md`](architecture/onemind-tool-gateway.md); runbook: [`tool-gateway.md`](operations/tool-gateway.md). Verified live 2026-10-30 (read-only): Synced and Healthy, 1/1 Ready | none (config is in the repo; sessions are in memory) | Git (SealedSecrets `gateway-clients`, `gateway-target-onemind`, `-core`, `-mediamtx`, `-clickhouse`) | internal: ClusterIP `agentgateway:8080`, ingress only from `onemind-hermes` | Argo |
| `platform-ecosystem-resources` | argocd | Ecosystem SealedSecrets, namespaces and shared resources | sealed credentials | Git | internal | Argo. **OutOfSync** (`SharedResourceWarning`: `SealedSecret auth/ghcr` is also owned by `platform-backups`) |

### Argo CD Applications, ecosystem layer

All were `Synced` and `Healthy`. Chart versions are the pinned `targetRevision`. Foundation charts are also checksum-locked in `platform/ecosystem/foundation-lock.json`.

| Application | Namespace | Chart (version) | Purpose | State | Backup | Exposure |
|-------------|-----------|-----------------|---------|-------|--------|----------|
| `ecosystem-cnpg` | cnpg-system | cloudnative-pg 0.29.1 | CloudNativePG operator | none | Git | internal |
| `ecosystem-barman` | cnpg-system | plugin-barman-cloud 0.8.1 | CNPG backup plugin | none | Git | internal |
| `ecosystem-onemind-core` | onemind-core | local `platform/ecosystem/charts/onemind-core` | OneMind core runtime. CNPG cluster `onemind-data` (3 instances, 120 GiB) lives in `onemind-data` | PostgreSQL | `Velero`, `Barman` | internal |
| `ecosystem-openfga` | onemind-core | openfga 0.3.16 | Fine-grained authorization | PostgreSQL (CNPG) | `Velero`, `Barman` | internal |
| ~~`ecosystem-litellm`~~ | onemind-ai | litellm-helm 1.104.2 (OCI) | **Retired 2026-11-06** (removed from Git). Database and role `litellm` dropped (verified live 2026-11-07) | PostgreSQL (CNPG) | `Velero`, `Barman` | none after cleanup |
| ~~`ecosystem-ollama`~~ | onemind-ai | ollama 1.86.0 | **Retired 2026-11-06** (removed from Git). Volume `onemind-cpu-model` (20 GiB, Retain) kept | 20 GiB models | `Velero` (`onemind-ai`); VolSync source removed | none after cleanup |
| `ecosystem-langfuse` | onemind-ai | langfuse 2.2.3 (app 4.75.0) | LLM observability for Hermes traces (web, worker, bundled Valkey). **Prepared 2026-11-07, not yet synced** | PostgreSQL (CNPG `onemind-data`), ClickHouse `langfuse` database, Spaces prefix `onemind-backups/langfuse/` | `Velero`, `Barman`, `VolSync` (ClickHouse) | tailnet `onemind-langfuse` |
| `ecosystem-kyverno` | kyverno | kyverno 3.9.1 | Policy engine | none | Git | internal |
| `ecosystem-falco` | falco | falco 9.2.0 | Runtime security | none | Git | internal |
| `ecosystem-opencost` | opencost | opencost 2.5.32 | Cost metrics. Duplicates kube-state-metrics series (M2) | none | Git | internal |

### Tool gateway backends (tier 2)

Reached only through `platform-tool-gateway` (namespace `tool-gateway` unless stated); the backend rows share its Argo Application, so they have no Application of their own. All stateless, backed up by Git, `internal` exposure, GitOps `Argo`. "Live" below was checked read-only against `do-nyc3-omos-cloud` on 2026-10-30; tool counts are from the reviewed classification files in `platform/tool-gateway/policy/classification/`. Details: [`tool-gateway-kubernetes.md`](operations/tool-gateway-kubernetes.md), [`tool-gateway-data.md`](operations/tool-gateway-data.md).

| Backend | Where it runs | Tier | Route and tools | Live check |
|---------|---------------|------|-----------------|------------|
| `agentgateway` | Deployment in `tool-gateway`, 1 replica | 2 | all nine routes | Deployment 1/1; image `agentgateway:v1.7.0` by digest |
| `execute-gate` | Deployment and Service in `tool-gateway` (`ghcr.io/onemindos/execute-gate`, digest-pinned), `replicas: 0` until the owner session seals `gate-matrix`, `gate-nats` and `ghcr` | 2 | AuthZEN decision point for `/execute/mcp` (owner types `execute <id>` or `override <id>` in the Matrix room `approvals`); state in NATS KV `exec-approvals` | not running (off by design); until then `/execute` calls are denied (403, fail closed) |
| `mcp-kubernetes` | Deployment and Service in `tool-gateway` (`containers/kubernetes-mcp-server` v0.1.1, ServiceAccount `k8s-operator`, no Secrets in its RBAC) | 2 | `kubernetes` target, `ops`: 10 read, 1 write (`resources_create_or_update`); `pods_delete`, `resources_delete`, `resources_scale` held for `/execute` | Deployment 1/1; target in the rendered config |
| `mcp-clickhouse` | Deployment and Service in `tool-gateway` (`ClickHouse/mcp-clickhouse` 0.7.0), talks to ClickHouse as the user `legacy_ro` | 2 | `clickhouse` target, `data`: 3 read (`list_databases`, `list_tables`, `run_query`) | Deployment 1/1; sealed `mcp-clickhouse-credentials` and `gateway-target-clickhouse` present |
| ClickHouse user `legacy_ro` | declared in the Altinity CHI `onemind` (namespace `clickhouse`, Application `platform-observability`), `readonly=1`, 30 s, 1 GB, 100,000 rows, network `10.0.0.0/8`, password stored as a SHA-256 hash | 2 | used by `mcp-clickhouse` only | CHI `Completed`; `legacy_ro` users and profile settings present in the live spec |
| MediaMTX target | OpenAPI target `mediamtx` onto the Control API of `onetak-media` (namespace `onetak`, port 4000), spec trimmed in Git | 2 | `tak`: 25 read; 8 destructive operations hidden | Service `onetak-media` and its pod Running; ConfigMap `gateway-openapi-mediamtx` |
| STAC target | OpenAPI target `stac` onto `stac-api` (namespace `stac`, 8080), spec trimmed to GET in Git; the `stac-api-writer` Service is not a target | 2 | `data`: 5 read | Service `stac-api`; ConfigMap `gateway-openapi-stac` |
| onemind-core target | OpenAPI target `core` onto `onemind-core` (namespace `onemind-core`, 8080), hand-written read subset | 2 | `data`: 6 read | Service `onemind-core`; ConfigMap `gateway-openapi-core` |
| Tile38 adapter | `tile38` adapter inside `onemind-mcp` (`hermes`), served by the `onemind` target in `ops`; reaches `tile38.onetak:9851` | 2 | 6 read, 1 write (`tile38_set_point`); `tile38_delete_object` and `tile38_drop_collection` are not registered (`allow_destructive` off) | `onemind-mcp` 1/1 ready; `TILE38_URL` set and the app `tile38` has `allow_writes: true` in the live apps ConfigMap; Deployment `tile38` 1/1 |

### Helm-CLI releases (tier 2)

| Release | Namespace | Chart | Purpose | State | Backup | Exposure | GitOps |
|---------|-----------|-------|---------|-------|--------|----------|--------|
| `monitoring` | monitoring | kube-prometheus-stack-87.16.1 | Prometheus, Alertmanager, Grafana | Prometheus TSDB, Grafana DB, 40 GiB PVCs in the namespace | `VolSync` | Grafana tailnet; Alertmanager internal | Argo (`monitoring-kube-prometheus-stack`, adopted 2026-10-17) |
| `loki` | monitoring | loki-7.0.0 | Log storage. **PVC 89% full** (A4) | 20 GiB | `VolSync` | internal | Argo (`monitoring-loki`, adopted 2026-10-17) |
| `otel-collector` | monitoring | opentelemetry-collector-0.165.0 | Telemetry pipeline | none | Git | internal | Argo (`monitoring-otel-collector`, adopted 2026-10-17) |

## Not tracked by Argo or Helm

These run live and are not owned by any Application. Manifests for some of them exist in this repo but nothing applies them.

| Namespace | Workloads | Manifests in Git | Tier | Backup | Exposure |
|-----------|-----------|------------------|------|--------|----------|
| `nats` (UI) | `nui` (NUI 0.10.0, `ghcr.io/nats-nui/nui` digest-pinned, tailnet only at `nui`, behind the Authentik proxy outpost `nui-auth` (admins only); NUI has no login of its own, and its NetworkPolicy admits only the outpost and the Backstage health probe; `/db` is an emptyDir, connection imported from SealedSecret `nui-nats-context`). It replaced `heynats`, `nats-dashboard` and `nats-ui` on 2026-11-07. The agent tools moved to namespace `agent-fabric` (Agent Console, agent-tap, digest-pinned) on 2026-10-17 | `platform/nats-tools/`, Application `platform-nats-tools` (auto-sync, `prune: false`) | 2 | none needed | tailnet |
| ~~`matrix`~~ | Synapse, Element Web, PostgreSQL, Redis | **Retired 2026-11-07.** Manifests parked at `platform/parked/matrix/`; Application `platform-matrix` removed | n/a | Final `pg_dump` and signing key in `backups/matrix/final-2026-11-07/`; VolSync restic repos kept under `backups/restic/matrix/` | none (tailnet hosts `matrix.` and `element.tail717ba0.ts.net` removed) |
| `tailscale` | about 24 operator-generated proxy StatefulSets | generated | n/a | none | tailnet |
| `kube-system` | Cilium, CoreDNS, Hubble, metrics-server, DO CSI and agents | DOKS-managed | n/a | none | `hubble` tailnet |

## GitOps gaps

Eleven releases were installed with `helm upgrade --install` (see `platform/install.sh`). **Update 2026-10-17:** four of them (`cert-manager`, `ingress-nginx`, `nack`, `tailscale-operator`) were adopted into Argo in PR #176 and are owned by the `edge-*` Applications. The four monitoring releases (`monitoring`, `loki`, `tempo`, `otel-collector`) were adopted later the same day as `monitoring-*` Applications. **Update 2026-11-03:** Tempo was removed (no trace backend; the OTel collector now ships logs only), leaving three monitoring Applications. The remaining three are Helm CLI only: `nats`, `takserver` and `argocd`. `helm list -A` still shows all eleven, because adoption leaves the release record. Do not run `platform/install.sh`: it would `helm upgrade` the adopted releases too. Current status: `docs/gitops/HELM-ADOPTION.md`.

| Release | Namespace | Chart | Rev | Last updated | Values in Git | Risk if the cluster is rebuilt |
|---------|-----------|-------|-----|--------------|---------------|-------------------------------|
| `argocd` | argocd | argo-cd-10.1.2 | 3 | 2026-07-11 | `platform/values/argocd.yaml` | Argo itself is not declarative yet; `platform-argocd` is defined (manual sync), adoption steps in `docs/gitops/ARGOCD-SELF-MANAGEMENT.md` |
| `cert-manager` | cert-manager | cert-manager-v1.20.3 | 5 | 2026-07-11 | `platform/values/cert-manager.yaml` | No `ClusterIssuer` is deployed (none is needed: public TLS ends at Cloudflare and tailnet TLS is Tailscale's). Add one with its own Cloudflare token the day a direct certificate is needed |
| `ingress-nginx` | ingress-nginx | ingress-nginx-4.15.1 | 1 | 2026-07-04 | `platform/values/ingress-nginx.yaml` | LB IP `143.244.220.155` is bound to this release |
| `loki` | monitoring | loki-7.0.0 | 4 | 2026-07-11 | `platform/values/loki.yaml` | Logs only |
| `monitoring` | monitoring | kube-prometheus-stack-87.16.1 | 8 | 2026-07-17 | `platform/values/kube-prometheus-stack.yaml` | Alert routing: until 2026-11-07 Alertmanager's only receiver was `null` (A17). Now `AlertmanagerConfig/onemind-backstage` routes critical and warning (not `Watchdog`/`InfoInhibitor`) to the portal; there is no Telegram receiver (on hold) |
| `nack` | nats | nack-0.34.0 | 4 | 2026-07-11 | `platform/nats/nack-values.yaml` | Streams are not reconciled until reinstalled |
| `nats` | nats | nats-2.14.2 | 10 | 2026-07-11 | `platform/nats/values.template.yaml` plus an **out-of-band, gitignored** `values.secret.yaml` | Passwords exist only on the machine that last ran Helm; Argo and Helm both own the StatefulSet |
| `otel-collector` | monitoring | opentelemetry-collector-0.165.0 | 3 | 2026-07-17 | `platform/values/otel-collector.yaml` | Telemetry only |
| `tailscale-operator` | tailscale | tailscale-operator-1.98.4 | 3 | 2026-07-11 | `platform/values/tailscale-operator.yaml` | All admin UIs depend on it |
| `takserver` | takserver | takserver-4.6.0 | 34 | 2026-09-29 | `platform/takserver/values-onemind.yaml` and chart `platform/takserver/chart` | The most important workload has no reproducible rebuild path |

Whether the values in Git still match what is deployed was not checked; `helm get values` can contain secrets and was not run. Verify with a values diff before adopting any release.

**Suggested adoption order** (from the cluster review, finding H2): `takserver` and `nats` first, behind fresh volume snapshots; then `monitoring`, `loki`, `otel-collector`, `nack`; then `ingress-nginx`, `cert-manager`, `tailscale-operator`; `argocd` last. Per release: export values, commit them, create an Application with the same chart, version and `releaseName`, sync with `ServerSideApply` and `prune: false`, then confirm a zero diff. To back out, delete the Application with `--cascade=false` (workloads keep running, the Helm release secret still exists). This document does not make any of those changes.

Gaps beyond the eleven releases:

- **Adopted workloads (pending first sync):** the former `nats` UIs (now NUI) have Applications whose manifests were reconciled against the live objects (live won every difference). `platform-nats-tools` is auto-sync with `prune: false`; `platform-matrix` is manual because the live Synapse `homeserver.yaml` holds plaintext database, Redis and registration secrets that cannot be committed. Remaining gaps: the Secrets (`matrix-secrets`, `nats-creds`, `ghcr`) and the `matrix` and `nats` Namespaces are still created out of band; `platform/nats-tools/*.sealedsecret.yaml` and `namespace.yaml` target a `nats-tools` namespace that does not exist and are excluded; the `agenttap` Ingress points at port 8080 while its Service listens on 8088 (unchanged from live); `matrix-secrets` was applied with `kubectl apply`, so its values sit in a `last-applied-configuration` annotation and should be rotated.
- **Two owners of one resource:** `nats` StatefulSet (Helm and `platform-nats`); `SealedSecret auth/ghcr` (`platform-ecosystem-resources` and `platform-backups`).
- **Terraform vs live:** `infra/main.tf` defaults (`onemind-hub`, `s-4vcpu-8gb` x3, `1.36.0-do.2`, pool `onemind-hub-core`) do not match the live cluster (`omos-cloud`, pool `onemind-hub-core-v2` on 8 vCPU / 16 GiB, `1.36.0-do.3`). `infra/ecosystem-workers/` is a separate root for the ecosystem pool. Terraform state was not inspected; do not apply `infra/` without a clean `terraform plan`.

## Exposure map

| Path | What is behind it |
|------|-------------------|
| Public, Cloudflare ranges only (nginx) | `auth.onemindos.dev`, `ws.onemindos.dev` (NATS WebSocket), `onemindos.dev` (tak-portal) |
| Public, control plane | Kubernetes API endpoint, no control-plane firewall (review finding H3) |
| Tailnet only | Matrix (`matrix`, `element`), TAK (`tak`: 8089, 8443, 8446, 64738), TAK Portal (`enroll`), FedHub (`fedhub`), Argo CD, Grafana, Hermes dashboard, n8n, Headlamp, Hubble, Backstage preview, Authentik admin (`onemind-auth`), NUI (`nui`), `omos-nats` (4222, 7422, 8080), OpenChoreo gateway |

## Mutable image tags

> **Resolved 2026-10-23.** The running references below were pinned to their observed
> image digests ([running-image-pins-2026-10-23.md](operations/running-image-pins-2026-10-23.md)).
> A read-only check on 2026-10-27 (`kubectl get pods -A`, 141 unique container images)
> found none on `:latest` or a bare name. The only remaining `:latest` under `platform/`
> is `platform/parked/omos-mission-engine/`, which has no Application.
> The Kyverno `disallow-latest-image-tag` policy is still in `Audit` mode.

At capture (2026-10-09), 31 container references used `:latest`, a bare image name, or `run-latest` and so were not reproducible:

| Namespace | Workloads |
|-----------|-----------|
| nats | `heynats`, `nats-dashboard`, `nats-metrics-scrape`, `nats-ui` (all `ghcr.io/...:latest`), `nats-surveyor:latest` |
| onetak | `nodeodm`, `stac-api`, `tak-portal`, `tile38`, `tileserver-gl`, `webodm`, `webodm-db` (and a bare `busybox` init container) |
| feeds | `omos-mission-engine:latest` |
| kube | `headlamp:latest` |
| takserver | `murmur:latest` |
| matrix | init containers using bare `busybox` / `busybox:latest` |

Good patterns to copy: NATS server, ingress-nginx, `takserver-plugins` and `redpanda-connect` are digest-pinned; `github-runner` is digest-pinned.

## Plaintext credentials found in tracked files

The new secret scan found nine locations that look like real or example credentials committed in plaintext. Their line numbers are recorded in `.gitleaks-baseline.json` so CI only fails on **new** findings. **They are not fixed by that baseline.** Treat the ones that are real as leaked (they are in Git history), rotate them, and remove them from the files. Values are deliberately not repeated here.

| File | Lines | Likely nature |
|------|-------|---------------|
| `docs/infra/CONSTELLATION-OVERWATCH.md` | 183 | API key literal in an example |
| `platform/hermes/agent/SECRETS.md` | 74, 80 | OAuth client secret in an example command. The file was deleted with `hermes-legacy` on 2026-11-03 but is still in history; rotate the secret |
| `platform/monitoring/deployment.yaml` | 53 | `--password=` argument on a container |
| `services/om_presence.py` | 50, 52 | Default token and API key in code |
| `infra/edge/takserver/deploy/docs/enrollment.md` | 123, 139 | Line 139 is a `curl -u user:password` example; line 123 is probably a false positive |
| `docs/references/ai/orp.md` | 28 | False positive (prose about a keystore) |

`platform/monitoring/` belongs to another change; flag the `--password` argument there rather than editing it in this one.

## How to refresh

```bash
export CTX=do-nyc3-omos-cloud                       # read-only commands only
kubectl --context $CTX -n argocd get applications
kubectl --context $CTX -n argocd get applications -o json   # sources, targetRevision, sync and health
helm --kube-context $CTX list -A
kubectl --context $CTX get ingressclass,ingress -A
kubectl --context $CTX get svc -A | grep -E 'LoadBalancer|NodePort'
kubectl --context $CTX get deploy,sts,ds -A -o json        # images, replicas, ownership annotations
kubectl --context $CTX get pvc -A
kubectl --context $CTX get replicationsources.volsync.backube,schedules.velero.io,cronjob -A
```

Ownership rule used above: a workload carrying `argocd.argoproj.io/tracking-id` is Argo-managed; `meta.helm.sh/release-name` without it is Helm-CLI; neither means untracked.

## Nextcloud file sync — 2026-10-23

Observed live: Argo Application `platform-nextcloud` is Synced/Healthy, Nextcloud
35.0.2 reports installed without maintenance/database upgrade pending, and the
application (2/2) and Redis (1/1) pods are healthy. The retained file PVC is bound
at 200 GiB; the VolSync cache PVC is 1 GiB. Private HTTPS is available at
`https://onemind-files.tail717ba0.ts.net/`. A separate database/role uses existing
CNPG. PRs #424/#425 are merged. The Mac's full-library upload has started;
Samsung FolderSync two-way acceptance and full-transfer completion remain pending.
See [setup, cost and acceptance](../platform/nextcloud/README.md).

## OneMind application MCP — verified 2026-10-27

Tier 2, namespace `hermes`, Deployment `onemind-mcp`, owned by `platform-onemind-mcp` (moved from `platform-legacy` on 2026-11-01).
Stateless private app tools for Legacy: NATS node/JetStream stats and Prometheus
instant queries, with a discoverable app catalog. One available replica, 100m CPU /
128Mi memory requests; no database/PVC. Git source/catalog and SealedSecrets are
recovery artifacts. A dedicated image-pull credential keeps unrelated registry
credentials unchanged. No Kubernetes API privileges are granted. Update 2026-10-30: since the tool-gateway work it also serves the Ollama, OpenCost and Tile38 adapters, and the Tile38 adapter has one write tool (`tile38_set_point`; its two destructive tools are not registered). It is the `onemind` target of `platform-tool-gateway` (see [Tool gateway backends](#tool-gateway-backends-tier-2)).

[Live acceptance](operations/onemind-mcp-acceptance-2026-10-27.json), 04:00 EDT:
all four tools and resource reads pass from Legacy, anonymous access is rejected,
MCP and Legacy are available, and Argo is Synced/Healthy. Extension, rotation and
rollback: [runbook](operations/onemind-mcp.md).

## Nextcloud retirement — 2026-11-01

Zeus authorized retiring Nextcloud and its 200 GiB file/1 GiB cache volumes
to use Dropbox. Deployment manifests are archived under
`platform/parked/nextcloud`; the Argo Application is removed. Execution requires
a completed fresh volume snapshot plus private database/config exports.
Provider deletion and recovery verification must be recorded from live evidence;
this source change alone does not certify either. Preserve device originals and
existing backups. See `platform/parked/nextcloud/RETIREMENT.md`.

## User-requested IoT pilots — prepared November 1, 2026, running 2026-11-07

| Tier | Application | Namespace | State/recovery | Exposure | Status |
|---|---|---|---|---|---|
|2|`platform-openremote`|`openremote`|Dedicated PostgreSQL17/Keycloak, retained10Gi DB plus5Gi manager storage; separate crash-consistent VolSync sources, coordinated recovery pending|Private `openremote.tail717ba0.ts.net` (Tailscale ingress); native login; no external MQTT|Manual sync, Synced/Healthy since 2026-11-05; 4 pods Running (2d16h); retained database and storage claims Bound; VolSync sources syncing; coordinated restore and acceptance not recorded|
|2|`platform-thingsboard`|`thingsboard`|Apache CE4.3.1.8 monolith, retained10Gi PostgreSQL; no-demo guarded schema/admin hook; in-memory queue is an evaluation limitation|Private `thingsboard.tail717ba0.ts.net` (Tailscale ingress); native login; no external MQTT|Manual sync, Synced/Healthy (last sync 2026-11-06); ThingsBoard and PostgreSQL pods Running; VolSync source syncing; acceptance not recorded|

These are fresh evaluation deployments, not migrations of retained Mac volumes. User expressly chose both for
comparison. See the owning READMEs for edition, launch and restore requirements.

## Forgejo secondary forge — launch-ready 2026-11-04, running 2026-11-07

User-directed private GitHub pull-mirror standby. Synced on 2026-11-05 and verified live 2026-11-07: Application
Synced and Healthy (manual sync), pod Running, Tailscale ingress `forgejo.tail717ba0.ts.net`, the retain-class
20Gi claim Bound, and the coordinated backup Job plus VolSync source completing (latest sync today). Restore
acceptance and the owner's first sign-in and mirror pilot are not recorded here. GitHub remains the working forge and Argo source.

| Tier | Application | Namespace | State / recovery | Exposure | Status |
| --- | --- | --- | --- | --- | --- |
| 2 | `platform-forgejo` | `forgejo` | Rootless Forgejo 15.0.11 LTS, SQLite/Git/LFS on one 20Gi `do-block-storage-retain` claim; nightly coordinated stop plus VolSync/Restic export (`backups/restic/forgejo/data`); restore acceptance not recorded | Private `https://forgejo.tail717ba0.ts.net` (Tailscale ingress), local accounts, no SSH or registration; hub pool preferred | Manual Argo, Synced/Healthy, pod Running (2026-11-07); [owner first sign-in and mirror pilot](../platform/forgejo/README.md#owner-first-sign-in) not recorded |

## Gatus status checks — prepared 2026-11-05, running 2026-11-07

Flagship F7. Live HTTP(S) and certificate-expiry checks for apps that run outside Kubernetes (SaaS) and for
cluster apps with no probe of their own. Verified live 2026-11-07: Application Synced and Healthy (manual sync, last sync 2026-11-06), pod Running 38 hours, Tailscale ingress `status.tail717ba0.ts.net`. Endpoints are generated from catalog annotations
(`onemind.io/health-check`, `onemind.io/health-url`); see [the app contract](catalog/app-contract.md).

| Tier | Application | Namespace | State / recovery | Exposure | Status |
| --- | --- | --- | --- | --- | --- |
| 2 | `platform-gatus` | `gatus` | Gatus v5.37.0 pinned by digest, in-memory storage (no PVC, nothing to back up; check history resets on restart, Prometheus keeps it); requests 10m CPU / 48Mi, limits 200m / 128Mi; ecosystem pool preferred | Private `https://status.tail717ba0.ts.net` (Tailscale ingress); `/metrics` scraped by Prometheus; portal reads it through `/api/proxy/gatus/statuses` | Manual Argo, Synced/Healthy, pod Running (2026-11-07); the deliberate failure `test-deliberate-failure` is no longer in the live config (per the [README](../platform/gatus/README.md)) |
