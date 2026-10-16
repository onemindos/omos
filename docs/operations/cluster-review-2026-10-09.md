# OneMind cluster review — 2026-10-09

Scope: DOKS `omos-cloud` (nyc3), ArgoCD, Helm, DigitalOcean account, plus the AI-subscription question.
Mode: **Phase 1 + 2 only (read-only).** No cluster, DO, or Git changes made. Phase 3 is waiting on Zeus's go.
Data collected 2026-10-10 ~01:30 UTC. kube context `do-nyc3-omos-cloud`. Prometheus 7-day window for usage.

---

## 0. Headline

The three Oct-7 problems are **already fixed** (by the `codex/*` PRs #91, #92, #121, #122, #128 merged Oct 8):

| Oct-7 item | Now | Evidence |
|---|---|---|
| `platform-takserver` Degraded | **Synced / Healthy** since 2026-10-08 11:04Z (rev `25131420`) | `argocd app` status, history |
| `nats/backup-nats-streams` never succeeded | **Succeeded** 2026-10-09 03:33Z | job log: `{"kind":"nats","restoreVerified":true,"streams":39}` → `backups/nats-onemind/20261009T033001Z-…/archive.tar.gz.aesgcm` |
| `takserver/backup-takserver-postgres` never succeeded | **Succeeded** 2026-10-09 03:26Z (86 min run) | job log: `{"kind":"postgres","restoreVerified":true,"tables":64}` → `backups/takserver-postgres/20261009T020032Z-…/archive.tar.gz.aesgcm` |

Both backup jobs restore into an in-pod scratch Postgres / NATS and verify before reporting success. I did **not** list the bucket objects directly (that needs the backup credentials; see Phase 3 step 1).

Current open problems are about **memory imbalance on the hub-core pool, non-GitOps Helm releases (TAK core, NATS, monitoring), plaintext tokens in the vault, and ~$190/mo of reclaimable spend.**

---

## 1. Inventory (Phase 1)

### 1.1 DigitalOcean

| Item | Value |
|---|---|
| Account | active, droplet limit 15 |
| Cluster | `omos-cloud` `051f1134…`, nyc3, **1.36.0-do.3**, HA control plane, auto-upgrade on (Sun 07:00, 4h), surge upgrade on, **control-plane firewall: not set** |
| Upgrade available | 1.36.5-do.0 (1.37.0 exists, not yet offered) |
| Pool `onemind-hub-core-v2` | 3 × `s-8vcpu-16gb` ($96) — fixed, no autoscale — 75 days old |
| Pool `onemind-ecosystem-cpu` | 2 × `s-8vcpu-32gb-amd` ($168) — autoscale 2–4 — 26 h old |
| Other droplet | `wp-zeus` s-1vcpu-1gb, **powered off** (still billed) |
| Volumes | 95 PVC-backed, **1,182 GiB**; 48 unattached (246 GiB): 46 × 1 GiB VolSync cache + `onetak/nominatim-data` 100 GiB + `onetak/overpass-data` 100 GiB |
| Volume snapshots | 2 baseline sets (Oct 8 10:56Z, Oct 9 10:27Z), 50 snapshots, **266 GiB** |
| Load balancers | `onemind-tak` 152.42.155.95 (TAK ports 8089/8443/8446/9101/9102/64738, no source ranges); ingress-nginx 143.244.220.155 |
| Managed DBs | none (CNPG / in-cluster Postgres) |
| Container registry (DOCR) | none (GHCR + in-cluster Harbor) |
| Firewalls | DOKS worker + public-access FWs; `onemindos-cloud-fw` attached to nothing |
| Billing | MTD usage **$0.04**; invoices Jul–Sep $0.00 → account is almost certainly on credits. **Verify credit balance/expiry in the DO dashboard** — list price below is what you pay when credits end. |

### 1.2 Nodes (point-in-time + 7-day)

| Node | Pool | CPU used | Mem used | Mem requested | Mem limits |
|---|---|---|---|---|---|
| hub-core-37phka | hub-core | 1.3 cores (16%) | **11.0 GiB (80%)** | **98%** | **274%** |
| hub-core-37phkg | hub-core | 1.2 (15%) | 10.3 GiB (75%) | 87% | 337% |
| hub-core-37phke | hub-core | 1.1 (13%) | 6.7 GiB (49%) | 85% | 226% |
| ecosystem-3oleu2 | ecosystem | 0.9 (11%) | 7.2 GiB (25%) | 57% | 155% |
| ecosystem-3oleul | ecosystem | 0.7 (8%) | 4.0 GiB (13%) | 36% | 143% |

Cluster totals: capacity 40 vCPU / 112 GiB · requested **27.4 cores / ~64 GiB** · **7-day peak 5.7 cores / 42.4 GiB**.

### 1.3 Workloads

- 47 namespaces, 186 deploy/sts/ds, 255 pods. **All deployments/statefulsets/daemonsets fully ready.**
- Pods not Running/Completed: 8, all `Error` pods from the Oct-8 manual backup drills (`backstage`, `n8n`, `nats`, `ops` verify jobs). Later variants passed (`backup-nats-sized-20261008`, `backup-backstage-owned-verify-20261008`). `ops/backup-cluster-recovery-verify-20261008` failed 2/2; the scheduled `ops/backup-cluster-recovery` succeeded 2026-10-09 10:56Z.
- Restarts > 5 (all old, none in the last 24 h except surveyor): authentik-server 18 (23d ago), node-exporter 11, otel-agent 11, nodeodm 11, agenttap 6, heynats 6, **nats-surveyor 12 (last 38h ago)**.
- CronJobs: 13; all non-suspended ones have a lastSuccessfulTime ≤ 24 h. Suspended: `backstage/backstage-backup`, `clickhouse/backup-clickhouse` (retired in #128; ClickHouse now covered by VolSync).
- VolSync: 46 ReplicationSources, **all `Successful`** within the last 24 h.
- PVs: 95, all Bound. No Pending PVCs.

### 1.4 Memory by namespace (7-day peak vs requests)

Requests below are **halved** from Prometheus because OpenCost re-exports kube-state-metrics (finding M2); they match `kubectl describe nodes`.

| Namespace | 7d peak | Requested | Over-request |
|---|---|---|---|
| takserver | 6.0 GiB | 8.1 GiB | +2.1 |
| onemind-ai (ollama, litellm) | 7.4 GiB | 8.5 GiB | +1.1 |
| onetak | 2.7 GiB | 6.5 GiB | **+3.8** |
| onemind-data | 0.6 GiB | 3.0 GiB | **+2.4** |
| onemind-core | 0.3 GiB | 2.0 GiB | **+1.7** |
| monitoring | 3.7 GiB | 2.6 GiB | −1.1 (under-requested) |
| kube-system | 4.1 GiB | 3.3 GiB | −0.8 (DO-managed) |
| nats | 1.5 GiB | 2.2 GiB | +0.7 |
| matrix | 0.3 GiB | 1.5 GiB | +1.2 |
| onemind-search | 1.4 GiB | 2.0 GiB | +0.6 |
| kyverno / spire / openchoreo-cp / registry | 0.2–0.3 GiB each | 0.9–1.1 GiB each | ~+3.5 total |

Top pods by memory now: webodm 1.8 GiB, openclaw 1.7, argocd-app-controller 1.5, opensearch 1.5, fedhub 1.3, hermes-legacy 1.3, clickhouse 1.2, prometheus 1.0.

Containers missing a memory request or limit (running): 50 — 25 Tailscale proxies, 10 DO-managed kube-system, 5 Harbor, 3 SPIRE, `takserver-postgresql` (request, no limit), `takserver-nats` reloader, Prometheus/Alertmanager config-reloaders, velero node-agent, falcoctl, ingress-nginx controller (no limit).

### 1.5 ArgoCD (47 apps, all automated + selfHeal)

Not Synced/Healthy:

| App | State | Cause |
|---|---|---|
| `app-of-apps` | OutOfSync/Healthy | child Applications `platform-ecosystem-applications`, `platform-ecosystem-foundation` differ from Git |
| `ecosystem-litellm` | OutOfSync/Healthy | `Job/onemind-ai/litellm-migrations` (completed Job, immutable spec) |
| `platform-ecosystem-resources` | OutOfSync/Healthy | **SharedResourceWarning**: `SealedSecret/auth/ghcr` owned by both this app and `platform-backups` |

ArgoCD is at `f226db81` = `origin/main` HEAD. Repo-vs-live Application set: no drift (ecosystem-* apps come from `platform/ecosystem/applications/`).

### 1.6 Helm (`helm list -A`)

11 releases installed by **Helm CLI**, not by Argo: `argocd`, `cert-manager`, `ingress-nginx`, `loki`, `monitoring` (kube-prometheus-stack 87.16.1), `nack`, `nats` (2.14.2), `otel-collector`, `tailscale-operator`, **`takserver`** (rev 34, 2026-09-29), `tempo`. None failed/pending.
- `takserver` core deployments (`takserver-api`, messaging, config, plugins, ignite, postgres) carry the Helm release annotation and **no** Argo tracking id. `platform-takserver` (path `tak-server/deploy`) manages only 18 adjacent objects (5 Deployments, 1 StatefulSet, ConfigMaps, SealedSecrets).
- `nats/sts/nats` has **both** a Helm release and Argo tracking (`platform-nats`) → two owners.

### 1.7 Exposure

- Admin UIs (argocd, grafana, headlamp, k8s-dashboard, hubble, n8n, hermes, openclaw, nats UIs, backstage) are Tailscale-ingress only. Good.
- Public nginx ingresses locked to Cloudflare ranges: auth, ws (NATS), onetak, geo-services.
- Public nginx ingresses **without** source restriction: `matrix`/`element`, `enroll.onemindos.dev` (tak-portal), `fedhub.onemindos.dev`, zeus-site.
- TAK LB open to 0.0.0.0/0 on 8089/8443/8446/9101/9102/64738.
- Kubernetes API: public endpoint, no control-plane firewall.

### 1.8 Warning events (last 60, deduped)

- `kube-system/doks-telemetry-config-reloader`: `FailedToRetrieveImagePullSecret onemind-hub` ×1,850 (noise; secret referenced but missing in kube-system).
- VolSync mover pods: transient `FailedScheduling … unbound immediate PVC` (resolve on retry; all ReplicationSources Successful).
- `VolumeFailedDelete pvc-42973aec…` — stale; PV and VolumeAttachment are gone.

### 1.9 Local repo state (`~/omos`)

Checked out on `codex/mission-platform-research` (`aefb075c`), **uncommitted**: `docs/MISSION.md`, `tak-server/build/build-tak-ldap.sh`; untracked `.obsidian/`, `OneMind-OS-Architecture.md`. Not mine — I left them alone. Phase 3 Git work will branch from `origin/main`.

---

## 2. Findings + plan (Phase 2)

### Critical

**C1. Plaintext Notion / Raindrop / 1Password service-account tokens in the vault `CLAUDE.md`**
- Evidence: `…/Dropbox-Onemindos/Zeus Delacruz/OneMind/CLAUDE.md` lines 15, 58, 60. The file is gitignored (`.gitignore:48`) and was never committed to `github.com/Zeus-Delacruz/onemind-vault`, but it is **Dropbox-synced** and loaded into the context of every agent session started in that folder (Claude, Codex, Hermes, OpenClaw, any Ollama Cloud model you point at it).
- Risk: the 1Password SA token can read every item in its vault(s). That's the keys to the cluster's secrets path.
- Fix (you): rotate the three tokens; replace the values with `op://` references or env vars loaded from 1Password at runtime; keep the file gitignored. I won't rotate or touch it without an explicit yes.
- Type: ops (outside the cluster). Rollback: n/a (rotation is one-way; keep the old token alive until consumers are moved).

### High

**H1. Hub-core pool is memory-saturated while the new ecosystem pool idles**
- Evidence: 37phka 80% used / 98% requested / 274% limits; 37phkg 75% / 87% / 337%. Ecosystem nodes 13–25% used.
- Risk: under a spike (TAK import, WebODM job, Ollama load), the kernel OOM-kills or kubelet evicts on hub-core — TAK and NATS live there. New pods can't schedule on hub-core.
- Fix (Git): (a) cut over-requests using the 7-day table (onetak, onemind-data, onemind-core, matrix, kyverno, spire, openchoreo — roughly −14 GiB of requests); (b) raise under-requests for monitoring/prometheus; (c) set `nodeAffinity` so `onemind-*` ecosystem, opensearch, ollama, webodm prefer the `onemind-ecosystem-cpu` pool, keeping TAK/NATS/auth on hub-core; (d) add memory limits where missing on stateful pods (`takserver-postgresql`).
- Rollback: revert the commit; Argo re-syncs the old requests.

**H2. TAK server core, NATS, monitoring, ArgoCD, ingress, cert-manager are Helm-CLI releases outside GitOps**
- Evidence: §1.6. `takserver` rev 34 applied by hand 2026-09-29; NATS has two owners.
- Risk: no reproducible rebuild for the most important workload. A `helm upgrade` from a laptop and an Argo sync can fight over the NATS StatefulSet. Values live wherever the last operator ran `helm`.
- Fix (Git): per release, export current values (`helm get values`), commit them under `platform/<name>/values.yaml`, create an Argo Application with the same chart/version/releaseName, sync with `ServerSideApply` and no prune, then confirm zero diff. Do `takserver` and `nats` first, behind the volume snapshots required by your rules. Leave argocd self-managed last.
- Rollback: delete the new Argo Application with `--cascade=false` (orphans resources and leaves them running); Helm release secret still exists.

**H3. Kubernetes API is reachable from the internet**
- Evidence: `control_plane_firewall: null`.
- Risk: the API is auth-protected, but it's a standing target for every CVE.
- Fix (ops, doctl): enable the DOKS control-plane firewall and allow only your Tailscale exit / home IPs, plus GitHub runner egress if CI talks to the API. Requires knowing every client IP; test with a second kube context first.
- Rollback: `doctl kubernetes cluster update … --enable-control-plane-firewall=false`.

### Medium

**M1. Three ArgoCD apps OutOfSync.** Fix (Git): remove `SealedSecret auth/ghcr` from one of `platform-ecosystem-resources` / `platform-backups`; give `litellm-migrations` the `argocd.argoproj.io/hook: PreSync` + `hook-delete-policy: BeforeHookCreation` annotations (or `ignoreDifferences` on Job spec); re-render `platform/argocd` children and diff. Rollback: revert.

**M2. OpenCost duplicates kube-state-metrics** (`job=opencost` exports 269 identical `kube_pod_container_resource_requests` series). Every Grafana/OpenCost request panel is 2× too high. Fix (Git): disable OpenCost's KSM emission in its values (`opencost.exporter.extraEnv EMIT_KSM_V1_METRICS=false` / `…_ONLY=false`). Rollback: revert.

**M3. 50 containers without memory requests/limits.** Fix (Git): Tailscale `ProxyClass` with resources (covers the 25 proxies); LimitRange defaults in `tailscale`, `onemind-registry`, `spire-system`; explicit limit on `takserver-postgresql`. Rollback: revert.

**M4. Public, unrestricted ingresses** — `enroll.onemindos.dev`, `fedhub.onemindos.dev`, matrix/element. Decide which need to be public. If they're used only by your own EUDs, put them behind the same Cloudflare allowlist as `onemindos.dev`. TAK LB 8446 (cert enrollment) is open to the world by design; keep strong enrollment passwords or move enrollment to Authentik.

**M5. Backup drill gaps.** `ops/backup-cluster-recovery-verify` failed 2/2 on Oct 8 with no passing retry on record. Takserver backup takes 86 min with no `activeDeadlineSeconds` and runs against a live DB (the vacuum-lock messages come from the scratch restore). Fix: rerun the cluster-recovery drill and capture logs; set a deadline of about 3 h on the TAK job.

**M6. Kubernetes 1.36.0 → 1.36.5** will auto-apply in the Sunday window. Check PDBs exist for TAK, NATS, Authentik and ingress before then, or pin the window.

**M7. `doks-telemetry` pull-secret event spam** (1,850 events). Remove the `onemind-hub` imagePullSecret reference from the kube-system default SA, or create it. Low risk.

### Cleanup

| # | Item | Saves | Action type |
|---|---|---|---|
| K1 | `onetak/nominatim-data` + `overpass-data` (2×100 GiB, `do-block-storage-retain`, apps scaled to 0, snapshots ≈0 GiB = essentially empty) | **$20/mo** | Your call — your rules forbid me deleting PVCs. Option: shrink to 20 GiB when you re-enable them. |
| K2 | `wp-zeus` droplet powered off | ~$6/mo | Snapshot + destroy (your yes) |
| K3 | 8 `Error` pods from the Oct-8 drill Jobs | clutter | `kubectl delete job` ×6 (ask) |
| K4 | Two full baseline snapshot sets (266 GiB) | ~$8/mo for the older set | Define retention (keep the newest set plus VolSync). Delete only with your yes. |
| K5 | `onemindos-cloud-fw` attached to nothing | $0 | Attach or delete |
| K6 | Local `~/omos` on a codex branch with uncommitted edits | — | You or Codex commit or stash |

### Cost (list price, before credits)

| Line | Monthly |
|---|---|
| hub-core 3 × s-8vcpu-16gb | $288 |
| ecosystem 2 × s-8vcpu-32gb-amd | $336 |
| HA control plane | $40 |
| 2 load balancers | $24 |
| Block storage 1,182 GiB | $118 |
| Volume snapshots 266 GiB | $16 |
| wp-zeus (off) | $6 |
| **Total** | **≈ $828/mo** |

Usage vs capacity: **peak 5.7 / 40 vCPU (14%) and 42 / 112 GiB (38%).** CPU is heavily over-bought because the 8 vCPU sizes were chosen for RAM.

Right-sizing options, smallest risk first:
1. **After H1's request cuts:** set the ecosystem autoscaler `min_nodes` to 1 → −$168/mo. That leaves 80 GiB raw capacity against 42 GiB peak, with autoscale back to 4 if needed.
2. **Swap hub-core to 3 × `s-4vcpu-16gb-amd` ($84)** → −$36/mo. Same RAM, 12 vCPU against a 5.7-core peak. This means a new pool, cordon and drain, then delete the old pool, behind TAK/NATS snapshots.
3. K1 + K2 + K4 → about −$34/mo.

Options 1 + 3 ≈ **−$200/mo (~24%)** without touching TAK node sizing. I would not drop HA control plane.

### Phase 3 plan (proposed order, waiting on your go)

1. **Backups:** already fixed. Remaining proof: (a) a read-only bucket listing of the two artifacts using the backup credentials from the in-cluster secret (I'll run it as a one-off pod in-cluster so the key never reaches this Mac's shell or chat); (b) a restore of the latest TAK Postgres + NATS archive into a scratch namespace `restore-drill-20261010`, with row and stream counts compared; (c) cleanup K3. Writes: create a scratch namespace and 2 Jobs. **Ask before each.**
2. **platform-takserver:** already Healthy. The real root-cause work is H2 (adopt the `takserver` Helm release into Argo). Snapshot the TAK PVCs first, per your rule 4.
3. **Right-size:** H1 request/affinity commit → watch 48 h → ecosystem `min_nodes` 1 (doctl, ask) → optional hub-core swap.
4. M1, M2, M3 (Git, low risk) in one PR; H3 and M4 need your IP list or decisions.

---

## 3. Actions taken

None that change state. Read-only commands only (doctl get/list, kubectl get/top/describe/logs, helm list/history, argocd status via CRDs, Prometheus queries through the API-server service proxy).

Side effect to note: `hermes --version` in this shell resolves to your `hermes()` zsh function, not the binary. That function `pkill`s any existing hermes / onemind-ai-bridge / copilot-runtime `kubectl port-forward` processes and starts new ones against context `do-nyc3-onemind-hub`. No port-forwards are running now. If you had some running before, they were killed. The real binary is `~/.local/bin/hermes`.

## 4. Rollback notes

No changes were made, so there's nothing to roll back. Per-item rollback is listed with each finding above.

---

## 5. Addendum: app-level review (2026-10-10, read-only)

Sources: Backstage catalog (`platform/backstage/catalog/*.yaml`, since the live catalog API requires sign-in), Alertmanager active alerts, Prometheus, pod logs (last 1 h), workload specs.

| # | Sev | Finding | Evidence | Fix (type) |
|---|---|---|---|---|
| A1 | High | **NATS alerting only covers "down / quorum lost".** Surveyor metrics are collected (`nats_*`), but nothing alerts if CoT stops flowing or a consumer falls behind. (Corrected: an earlier draft said metrics were missing; I had queried the wrong metric names.) | `onemind.nats` rule group has 2 rules | **PR 1:** `TAKCoTFlowStopped`, `NATSConsumerBacklogGrowing`, `NATSSurveyorDown` (Git) |
| A2 | High | **Hermes OOM-killed** 2026-10-09 19:59Z at its 3,200 Mi limit | `PodOOMKilled` alert; 7-day peak = limit | Raise the limit to 4–5 Gi and find the memory spike (Git) |
| A3 | High | **ArgoCD application-controller at 90% of its memory limit** (1.5 Gi) | Prometheus usage/limit | Raise the limit, or enable controller sharding. An OOM here stalls all GitOps (Git, helm values) |
| A4 | High | **Loki volume 89% full** (20 Gi, 744 h retention) | `PVCFillingUp` | Retention 14 d or expand PVC to 50 Gi (Git) |
| A5 | Med | **valhalla-data 99% full** (50 Gi); Valhalla exited 135 in Sep | `KubePersistentVolumeFillingUp` critical | Expand or rebuild tiles for the region you need (Git) |
| A6 | Med | **LiteLLM metrics scrape → 401** every 15 s: `TargetDown` alert + ~170 "Malformed API Key" errors/h | ServiceMonitor `/metrics/` with no auth | Set `require_auth_for_metrics_endpoint: false` or add bearer auth to the ServiceMonitor (Git) |
| A7 | Med | **litellm-migrations Job keeps being recreated** (OutOfSync loop) | Job age 25 s at review time; app OutOfSync on that Job | Argo hook annotations (PreSync + BeforeHookCreation) (Git) |
| A8 | Med | **TAK messaging: about 2 TLS handshakes/min without a client cert** (`PEER_DID_NOT_RETURN_A_CERTIFICATE`) | 120 errors/h, steady | Identify the client (an internal bridge or EUD missing its cert, or a scanner on the open 8089) (ops) |
| A9 | Med | **False-critical alerts**: KubeSchedulerDown, KubeControllerManagerDown, KubeProxyDown (DOKS-managed control plane; Cilium replaces kube-proxy) | Alertmanager | Disable those rules in kube-prometheus-stack values (Git) |
| A10 | Med | **Stale KubeJobFailed alerts** from the Oct-8 drill Jobs | Alertmanager | Delete the 6 failed drill Jobs (ops, ask) |
| A11 | Med | **Single points of failure**: every TAK component, Authentik, ingress-nginx, cert-manager and Hermes run 1 replica | Workload specs | ingress-nginx ×2 and Authentik server ×2 first; TAK HA is a separate design (Git) |
| A12 | Med | **31 workloads on mutable tags** (`:latest` etc.), including all NATS tools, OneTAK geo stack and TAK `*-5.8-onemind` images | Workload specs | Pin digests (Renovate or the existing digest-pin pattern) (Git) |
| A13 | Low | Orphan PDB `hermes/copilot-runtime` (no Deployment). The `hermes()` shell function still port-forwards to it | `kubectl get pdb` | Remove the PDB and the shell-function line (Git + local) |
| A14 | Low | `fedhub-mongodb` CPU-throttled 34% | Prometheus | Raise its CPU limit (Git) |
| A15 | Low | **Backstage catalog drift**: OneTAK marked "⚠️ Degraded" but live Healthy; Ollama "Planned" but `ecosystem-ollama` is deployed; DigitalOcean inference "Evaluation" but it is Hermes's default provider | Catalog vs live | Generate catalog status from Argo health instead of hand-written annotations (Git) |
| A16 | Low | Repo `CLAUDE.md` is stale (Bedrock, TAK 5.8 rows, "no public API endpoint") | File vs live | Refresh (Git) |

### 5.1 Alert delivery and backup location (found 2026-10-10)

| # | Sev | Finding | Evidence | Fix |
|---|---|---|---|---|
| A17 | **Critical** | **Alerts are delivered nowhere.** Alertmanager's only receiver is `null`; the Telegram config in `platform/monitoring/alertmanager-telegram.yaml` is commented out (placeholder chat ID) | Generated Alertmanager config | Pick a channel (Telegram / ntfy / Matrix / email) and wire it in (Git + SealedSecret) |
| A18 | High | **Every backup lives in DigitalOcean Spaces `nyc3`**, the same provider and region as the cluster (Velero `onemind-backups`, platform-backup jobs, VolSync) | Velero BSL, job env | Add an off-provider copy (e.g. Cloudflare R2 or Backblaze B2 now, home NAS later) (Git) |
