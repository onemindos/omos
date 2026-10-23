# OneMind stack inventory

> **OpenClaw retired 2026-10-22** (Zeus: not needed any more). Its Application, manifests, namespace, volume and backup source were removed through reviewed Git changes (PRs #379 and the retirement series). Earlier snapshots in this file that mention it are historical.

Tiered inventory of everything running on the DOKS cluster `omos-cloud` (nyc3), what owns it, what state it holds, whether it is backed up, who can reach it, and whether Git or a laptop controls it.

- **Captured:** 2026-10-09, read-only, kube context `do-nyc3-omos-cloud`, Argo CD synced to `main` at `4934f532`.
- **Scope:** all 47 Argo CD Applications (52 on 2026-10-17), all 11 Helm releases first installed by hand (4 adopted into Argo since, see GitOps gaps), and the workloads that nothing in Git or Helm tracks.
- **Status:** a point-in-time snapshot, not a live view. To refresh it, rerun the commands in [How to refresh](#how-to-refresh). Counts in the Summary go stale first.
- **Development phase:** every app is deliberate (system integration). This file ranks apps; it does not recommend removing any.
- **Owner:** there is one operator, `zeus`, for everything. The "Owner" column is therefore omitted from the tables; the first thing to change when a second person joins is to add it.

## Phase 5 addendum — 2026-10-19, 08:07 UTC

This narrow refresh does not recertify the older inventory below. Live cluster:
4 nodes, 70 Argo Applications (68 Synced/Healthy); only pre-existing Matrix
OutOfSync/Healthy and OpenClaw Synced/Progressing exceptions. No Pending pods.
All nodes are below 78% requested CPU/memory. OpenCost duplicate KSM emission is off.
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

## Summary

| Item | Count / value |
|------|---------------|
| Argo CD Applications | 47 (45 Synced and Healthy; `app-of-apps` and `platform-ecosystem-resources` OutOfSync) |
| Helm-CLI releases (not in Git as Applications) | 3 on 2026-10-17: `nats`, `takserver`, `argocd` (11 at capture; `cert-manager`, `ingress-nginx`, `nack`, `tailscale-operator` adopted in PR #176, the four monitoring releases on 2026-10-17) |
| Workloads that neither Argo nor Helm tracks | DOKS system and Tailscale operator proxies. `matrix` (4) and `nats` UIs (5) have Applications defined in the repo (`platform-matrix`, `platform-nats-tools`) that take ownership once synced; until then they are untracked |
| Namespaces | 48 |
| Nodes | 5: 3 x `onemind-hub-core-v2` (8 vCPU / 16 GiB), 2 x `onemind-ecosystem-cpu` (8 vCPU / 32 GiB) |
| Ingress classes | `nginx` (public, LB `143.244.220.155`), `tailscale` (tailnet only) |
| Other external listeners | `tak-tailnet` Tailscale LB (TAK, 8089/8443/8446/64738), `omos-nats` Tailscale LB, OpenChoreo gateway Tailscale LB |
| Container references on mutable tags (`:latest`, bare `busybox`, `run-latest`) | 31 (see [Mutable image tags](#mutable-image-tags)) |

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
| `platform-legacy` | hermes | Hermes "Legacy" agent (`nousresearch/hermes-agent`, NATS leaf sidecar) and bridge secrets | 25 GiB PVC (agent home, state.db) | `VolSync` | `hermes-agent` dashboard tailnet | Argo |

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
| `platform-backstage` | backstage | Backstage developer portal and service catalog | PostgreSQL, 10 GiB | `PG dump` (06:40 UTC), `VolSync` | `onemind-preview` tailnet | Argo |
| `platform-monitoring` | monitoring, nats | Grafana dashboards, rules, ServiceMonitors, `nats-surveyor`, `nats-metrics-scrape` (0 replicas). **Alertmanager has no receiver** | none (stack lives in Helm release `monitoring`) | Git | `grafana` tailnet | Argo |
| `platform-observability` | clickhouse | ClickHouse "Time Machine" (Altinity CHI) | 20 GiB time-series | `VolSync`; native backup Job suspended | internal | Argo |
| `platform-clickhouse-operator` | kube-system | Altinity ClickHouse operator (chart 0.27.3). A second, legacy `clickhouse-operator` Deployment sits at 0 replicas | none | Git | internal | Argo |
| `platform-omos` | omos | `feed-intel`, `feed-satellites`, `omos-tak-projector`, `sink-clickhouse`, `omos-mission-engine` (0 replicas) | none (streams in NATS, data in ClickHouse) | Git | internal | Argo |
| `platform-redpanda-connect` | omos | Redpanda Connect (digest-pinned) | none | Git | internal | Argo |
| `platform-onetak` | onetak | Tak-portal (root host), MediaMTX, geo stack (Tile38), WebODM | about 150 GiB of PVCs | `VolSync` x5 | `onemindos.dev`, geo hosts **public (CF)**; `enroll.tail717ba0.ts.net` **tailnet** | Argo |
| `platform-takml` | takml | TAK-ML vision pipeline | PostgreSQL, 8 GiB | `PG dump` (06:00 UTC), `VolSync` | internal | Argo |
| `platform-n8n` | n8n | n8n workflow automation | PostgreSQL, 15 GiB | `PG dump` (05:20 UTC), `VolSync` | tailnet | Argo |
| `platform-kube` | kube | Headlamp and Kubernetes Dashboard | none | Git | tailnet | Argo |
| `platform-ecosystem-applications` | argocd | Generates the `ecosystem-*` Applications (11) from `platform/ecosystem/applications` | none | Git | internal | Argo |
| `platform-ecosystem-foundation` | argocd | Generates the foundation `ecosystem-*` Applications (10) | none | Git | internal | Argo |
| `platform-ecosystem-resources` | argocd | Ecosystem SealedSecrets, namespaces and shared resources | sealed credentials | Git | internal | Argo. **OutOfSync** (`SharedResourceWarning`: `SealedSecret auth/ghcr` is also owned by `platform-backups`) |

### Argo CD Applications, ecosystem layer

All were `Synced` and `Healthy`. Chart versions are the pinned `targetRevision`. Foundation charts are also checksum-locked in `platform/ecosystem/foundation-lock.json`.

| Application | Namespace | Chart (version) | Purpose | State | Backup | Exposure |
|-------------|-----------|-----------------|---------|-------|--------|----------|
| `ecosystem-cnpg` | cnpg-system | cloudnative-pg 0.29.1 | CloudNativePG operator | none | Git | internal |
| `ecosystem-barman` | cnpg-system | plugin-barman-cloud 0.8.1 | CNPG backup plugin | none | Git | internal |
| `ecosystem-onemind-core` | onemind-core | local `platform/ecosystem/charts/onemind-core` | OneMind core runtime. CNPG cluster `onemind-data` (3 instances, 120 GiB) lives in `onemind-data` | PostgreSQL | `Velero`, `Barman` | internal |
| `ecosystem-openfga` | onemind-core | openfga 0.3.16 | Fine-grained authorization | PostgreSQL (CNPG) | `Velero`, `Barman` | internal |
| `ecosystem-litellm` | onemind-ai | litellm-helm 1.104.2 (OCI) | Model gateway | PostgreSQL (CNPG) | `Velero`, `Barman` | internal. `/metrics` scrape returns 401 (A6) |
| `ecosystem-ollama` | onemind-ai | ollama 1.86.0 | Local model runtime | 20 GiB models | `Velero` (`onemind-ai`), `VolSync` | internal |
| `ecosystem-opensearch` | onemind-search | opensearch 3.9.0 | Search | 20 GiB | `VolSync` | internal |
| `ecosystem-oc-control` | openchoreo-control-plane | openchoreo-control-plane 1.3.1 (+ overrides) | OpenChoreo developer platform | none | Git | gateway LB on **tailnet** |
| `ecosystem-dapr` | dapr-system | dapr 1.18.5 | Dapr runtime | 15 GiB | `VolSync` | internal |
| `ecosystem-keda` | keda | keda 2.21.0 | Event-driven autoscaling | none | Git | internal |
| `ecosystem-kyverno` | kyverno | kyverno 3.9.1 | Policy engine | none | Git | internal |
| `ecosystem-falco` | falco | falco 9.2.0 | Runtime security | none | Git | internal |
| `ecosystem-opencost` | opencost | opencost 2.5.32 | Cost metrics. Duplicates kube-state-metrics series (M2) | none | Git | internal |

### Helm-CLI releases (tier 2)

| Release | Namespace | Chart | Purpose | State | Backup | Exposure | GitOps |
|---------|-----------|-------|---------|-------|--------|----------|--------|
| `monitoring` | monitoring | kube-prometheus-stack-87.16.1 | Prometheus, Alertmanager, Grafana | Prometheus TSDB, Grafana DB, 40 GiB PVCs in the namespace | `VolSync` | Grafana tailnet; Alertmanager internal | Argo (`monitoring-kube-prometheus-stack`, adopted 2026-10-17) |
| `loki` | monitoring | loki-7.0.0 | Log storage. **PVC 89% full** (A4) | 20 GiB | `VolSync` | internal | Argo (`monitoring-loki`, adopted 2026-10-17) |
| `tempo` | monitoring | tempo-1.24.4 | Trace storage | PVC | `VolSync` | internal | Argo (`monitoring-tempo`, adopted 2026-10-17) |
| `otel-collector` | monitoring | opentelemetry-collector-0.165.0 | Telemetry pipeline | none | Git | internal | Argo (`monitoring-otel-collector`, adopted 2026-10-17) |

## Not tracked by Argo or Helm

These run live and are not owned by any Application. Manifests for some of them exist in this repo but nothing applies them.

| Namespace | Workloads | Manifests in Git | Tier | Backup | Exposure |
|-----------|-----------|------------------|------|--------|----------|
| `nats` (UIs) | `heynats`, `nats-dashboard`, `nats-ui` (`:latest`); the agent tools moved to namespace `agent-fabric` (Agent Console, agent-tap, digest-pinned) on 2026-10-17 | `platform/nats-tools/`, Application `platform-nats-tools` (auto-sync, `prune: false`) | 2 | none needed | tailnet |
| `matrix` | Synapse, Element Web, PostgreSQL, Redis | `platform/matrix/argo/`, Application `platform-matrix` (**manual sync** until `Secret/synapse-homeserver` is sealed, see `platform/matrix/SECRETS-TODO.md`; the older `base/` and `overlays/` are not deployed and not synced) | 2 | `PG dump` (05:40 UTC), `VolSync` x3 (80 GiB) | `matrix.onemindos.com`, `element.onemindos.com` **public (open)** |
| `tailscale` | about 24 operator-generated proxy StatefulSets | generated | n/a | none | tailnet |
| `kube-system` | Cilium, CoreDNS, Hubble, metrics-server, DO CSI and agents | DOKS-managed | n/a | none | `hubble` tailnet |

## GitOps gaps

Eleven releases were installed with `helm upgrade --install` (see `platform/install.sh`). **Update 2026-10-17:** four of them (`cert-manager`, `ingress-nginx`, `nack`, `tailscale-operator`) were adopted into Argo in PR #176 and are owned by the `edge-*` Applications. The four monitoring releases (`monitoring`, `loki`, `tempo`, `otel-collector`) were adopted later the same day as `monitoring-*` Applications. The remaining three are Helm CLI only: `nats`, `takserver` and `argocd`. `helm list -A` still shows all eleven, because adoption leaves the release record. Do not run `platform/install.sh`: it would `helm upgrade` the adopted releases too. Current status: `docs/gitops/HELM-ADOPTION.md`.

| Release | Namespace | Chart | Rev | Last updated | Values in Git | Risk if the cluster is rebuilt |
|---------|-----------|-------|-----|--------------|---------------|-------------------------------|
| `argocd` | argocd | argo-cd-10.1.2 | 3 | 2026-07-11 | `platform/values/argocd.yaml` | Argo itself is not declarative yet; `platform-argocd` is defined (manual sync), adoption steps in `docs/gitops/ARGOCD-SELF-MANAGEMENT.md` |
| `cert-manager` | cert-manager | cert-manager-v1.20.3 | 5 | 2026-07-11 | `platform/values/cert-manager.yaml` | No `ClusterIssuer` is deployed (none is needed: public TLS ends at Cloudflare and tailnet TLS is Tailscale's). Add one with its own Cloudflare token the day a direct certificate is needed |
| `ingress-nginx` | ingress-nginx | ingress-nginx-4.15.1 | 1 | 2026-07-04 | `platform/values/ingress-nginx.yaml` | LB IP `143.244.220.155` is bound to this release |
| `loki` | monitoring | loki-7.0.0 | 4 | 2026-07-11 | `platform/values/loki.yaml` | Logs only |
| `monitoring` | monitoring | kube-prometheus-stack-87.16.1 | 8 | 2026-07-17 | `platform/values/kube-prometheus-stack.yaml` | Alert routing (Alertmanager receiver) is not delivered anywhere (A17) |
| `nack` | nats | nack-0.34.0 | 4 | 2026-07-11 | `platform/nats/nack-values.yaml` | Streams are not reconciled until reinstalled |
| `nats` | nats | nats-2.14.2 | 10 | 2026-07-11 | `platform/nats/values.template.yaml` plus an **out-of-band, gitignored** `values.secret.yaml` | Passwords exist only on the machine that last ran Helm; Argo and Helm both own the StatefulSet |
| `otel-collector` | monitoring | opentelemetry-collector-0.165.0 | 3 | 2026-07-17 | `platform/values/otel-collector.yaml` | Telemetry only |
| `tailscale-operator` | tailscale | tailscale-operator-1.98.4 | 3 | 2026-07-11 | `platform/values/tailscale-operator.yaml` | All admin UIs depend on it |
| `takserver` | takserver | takserver-4.6.0 | 34 | 2026-09-29 | `platform/takserver/values-onemind.yaml` and chart `platform/takserver/chart` | The most important workload has no reproducible rebuild path |
| `tempo` | monitoring | tempo-1.24.4 | 4 | 2026-07-17 | `platform/values/tempo.yaml` | Traces only |

Whether the values in Git still match what is deployed was not checked; `helm get values` can contain secrets and was not run. Verify with a values diff before adopting any release.

**Suggested adoption order** (from the cluster review, finding H2): `takserver` and `nats` first, behind fresh volume snapshots; then `monitoring`, `loki`, `tempo`, `otel-collector`, `nack`; then `ingress-nginx`, `cert-manager`, `tailscale-operator`; `argocd` last. Per release: export values, commit them, create an Application with the same chart, version and `releaseName`, sync with `ServerSideApply` and `prune: false`, then confirm a zero diff. To back out, delete the Application with `--cascade=false` (workloads keep running, the Helm release secret still exists). This document does not make any of those changes.

Gaps beyond the eleven releases:

- **Adopted workloads (pending first sync):** the `nats` UIs and `matrix` now have Applications whose manifests were reconciled against the live objects (live won every difference). `platform-nats-tools` is auto-sync with `prune: false`; `platform-matrix` is manual because the live Synapse `homeserver.yaml` holds plaintext database, Redis and registration secrets that cannot be committed. Remaining gaps: the Secrets (`matrix-secrets`, `nats-creds`, `ghcr`) and the `matrix` and `nats` Namespaces are still created out of band; `platform/nats-tools/*.sealedsecret.yaml` and `namespace.yaml` target a `nats-tools` namespace that does not exist and are excluded; the `agenttap` Ingress points at port 8080 while its Service listens on 8088 (unchanged from live); `matrix-secrets` was applied with `kubectl apply`, so its values sit in a `last-applied-configuration` annotation and should be rotated.
- **Two owners of one resource:** `nats` StatefulSet (Helm and `platform-nats`); `SealedSecret auth/ghcr` (`platform-ecosystem-resources` and `platform-backups`).
- **Terraform vs live:** `infra/main.tf` defaults (`onemind-hub`, `s-4vcpu-8gb` x3, `1.36.0-do.2`, pool `onemind-hub-core`) do not match the live cluster (`omos-cloud`, pool `onemind-hub-core-v2` on 8 vCPU / 16 GiB, `1.36.0-do.3`). `infra/ecosystem-workers/` is a separate root for the ecosystem pool. Terraform state was not inspected; do not apply `infra/` without a clean `terraform plan`.

## Exposure map

| Path | What is behind it |
|------|-------------------|
| Public, Cloudflare ranges only (nginx) | `auth.onemindos.dev`, `ws.onemindos.dev` (NATS WebSocket), `onemindos.dev` (tak-portal) |
| Public, **no source restriction** (nginx) | `matrix.onemindos.com`, `element.onemindos.com` |
| Public, control plane | Kubernetes API endpoint, no control-plane firewall (review finding H3) |
| Tailnet only | TAK (`tak`: 8089, 8443, 8446, 64738), TAK Portal (`enroll`), FedHub (`fedhub`), Argo CD, Grafana, Hermes dashboard, n8n, Headlamp, Kubernetes Dashboard, Hubble, Backstage preview, Authentik admin (`onemind-auth`), NATS UIs, `omos-nats` (4222, 7422, 8080), OpenChoreo gateway |

## Mutable image tags

31 container references use `:latest`, a bare image name, or `run-latest` and so are not reproducible. Renovate (`renovate.json`) is deliberately **not** pinning digests for these yet; pinning is a separate, deliberate change because it will touch many manifests at once.

| Namespace | Workloads |
|-----------|-----------|
| nats | `heynats`, `nats-dashboard`, `nats-metrics-scrape`, `nats-ui` (all `ghcr.io/...:latest`), `nats-surveyor:latest` |
| onetak | `nodeodm`, `stac-api`, `stac-browser`, `tak-portal`, `tile38`, `tileserver-gl`, `webodm`, `webodm-db` (and a bare `busybox` init container) |
| omos | `omos-mission-engine:latest` |
| kube | `headlamp:latest` |
| takserver | `murmur:latest` |
| matrix, onemind-search | init containers using bare `busybox` / `busybox:latest` |

Good patterns to copy: NATS server, ingress-nginx, `takserver-plugins` and `redpanda-connect` are digest-pinned; `github-runner` is digest-pinned.

## Plaintext credentials found in tracked files

The new secret scan found nine locations that look like real or example credentials committed in plaintext. Their line numbers are recorded in `.gitleaks-baseline.json` so CI only fails on **new** findings. **They are not fixed by that baseline.** Treat the ones that are real as leaked (they are in Git history), rotate them, and remove them from the files. Values are deliberately not repeated here.

| File | Lines | Likely nature |
|------|-------|---------------|
| `docs/infra/CONSTELLATION-OVERWATCH.md` | 183 | API key literal in an example |
| `platform/hermes/agent/SECRETS.md` | 74, 80 | OAuth client secret in an example command |
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
