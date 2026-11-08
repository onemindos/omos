# OneMind Infra — Claude Context

This repo is the infrastructure-as-code and GitOps source of truth for the OneMind platform: the DOKS (DigitalOcean Kubernetes) cloud hub, platform manifests, NATS stream definitions, TAK Server chart and images, and Terraform for the DigitalOcean foundation.

Read `docs/MISSION.md` before architecture, product, AI or infrastructure decisions. Live state is in `docs/STACK.md` (tiered inventory, generated from the cluster); directory purposes are in `docs/REPO-LAYOUT.md`. Distinguish verified live state from intent: verify against the cluster (read-only) before writing "live" in docs.

The platform plan (four layers, AI trust model and phase tracker) is in `PLATFORM-ROADMAP.md` at the repository root. OpenChoreo has been retired, so its migration phases there are history. Read it before platform or AI-operations work, and update its tracker in the same PR when you complete roadmap work.

## What This Repo Is

If a service runs in the cluster, its manifest should live here, and Argo CD should apply it. Reality is not there yet: 63 Argo Applications exist (verified 2026-11-08; 60 Synced and Healthy, and the three manual-sync `platform-omos-entity`, `platform-policy` and `platform-takserver-tak-client` OutOfSync and Missing). Of the 11 Helm releases first installed by hand, 7 are adopted into Argo: `edge-cert-manager`, `edge-ingress-nginx`, `edge-nack`, `edge-tailscale-operator` (PR #176) and `monitoring-kube-prometheus-stack`, `monitoring-loki`, `monitoring-otel-collector` (2026-10-17). Tempo was removed on 2026-11-03 (`helm list -A` and the cluster no longer show it). **3 are still Helm CLI only**: nats, takserver and argocd (`helm list -A` on 2026-11-08 shows `nats` as `failed` at revision 11 while its three pods run). Never `helm upgrade` an adopted release, and do not run `platform/install.sh`, which still upgrades all of them; see `docs/gitops/HELM-ADOPTION.md`. See "GitOps gaps" in `docs/STACK.md`. Do not add to that list: new workloads get an Argo Application, not `helm install` or `kubectl apply`.

Removed 2026-11-03 to 2026-11-08 (do not redeploy without an owner decision): Dapr (`dapr-system`), LiteLLM and in-cluster Ollama, Matrix (parked in `platform/parked/matrix/`), HeyNATS, `nats-ui`, `nats-dashboard` and `nats-monitor`, Tempo, the `takml-db` PostgreSQL (TAK.ML uses `onemind-data`), the Velero node agent, and the Hermes Matrix gateway. The Argo CD ApplicationSet controller runs at 0 replicas, and so does `takserver-nats-box`. Namespaces were renamed `omos` to `feeds` and `onemind-ai` to `langfuse`; `kube` was not renamed (it still holds Headlamp and the `ai-operator` ServiceAccount, the owner's read-only MCP identity).

The platform is a NATS-native C2/ISR system. NATS is the nervous system, TAK (ATAK/iTAK/WebTAK) is the tactical display, and Hermes ("Legacy") is the AI agent.

## Bus Authority (read first)

**The in-cluster NATS (3-node JetStream, namespace `nats`) is the authority bus. Everything publishes to it.**

- Cloud workloads connect directly to `nats://nats.nats.svc.cluster.local:4222` (a few use `tls://` on the same host:port) with credentials from a Secret (`NATS_USER` / `NATS_PASSWORD`). No sidecar.
- **Hermes is no longer an exception.** `onemind-hermes` (the Legacy persona) is a direct client like every other cloud workload. Its pod has a single `hermes` container (verified 2026-11-08) that connects to `tls://nats.nats.svc.cluster.local:4222` as `onemind-hermes-nats-platform`, with `NATS_USER` / `NATS_PASSWORD` from `omos-nats-creds` and its own agent identity on the fabric (`agents.*.hermes.zeus.legacy`). The old `hermes-legacy` pod, with its `nats-leaf` sidecar on `localhost:4222`, was removed from Git on 2026-11-03. Do not bring a sidecar or leaf back. Hermes traces go to Langfuse at `langfuse-web.langfuse`.
- **Edge and Mac bridges** (MAVLink, Home Assistant, cameras) run a local leaf node that forwards to the cluster through the `omos-nats` Tailscale load balancer (leaf port 7422). The leaf is a forwarding pipe, not a peer bus. If the Mac is offline the cloud bus keeps running.
- When the home server arrives, it becomes another leaf (and eventually the authority). Do not design anything that treats a single laptop's NATS as the source of truth.

```
Edge (Mac / home server)   MAVLink, HA, cameras ── NATS leaf ──►  K8s NATS (3-node)  ◄── all cloud services
                                                  (forward only)       THE BUS            incl. Hermes (onemind-hermes)
```

## Repo Layout (short; full table in `docs/REPO-LAYOUT.md`)

```
platform/      Everything Argo CD syncs. argocd/ holds the Application CRs; the rest is grouped by layer
  argocd/      App-of-apps + 39 platform-* Applications, plus edge-* and monitoring-*
  core/        sealed-secrets controller, GitHub runner, Tailscale policy
  nats/        NATS manifests, TLS, 39 Stream/KeyValue CRs in streams/, Helm values templates
  auth/        Authentik infra (Postgres, Redis, ingress); the chart is in argocd/platform-auth-helm.yaml
  backups/     PG/NATS backup CronJobs, cluster-recovery export, VolSync sources (volsync/)
  hermes/      Hermes agent (onemind-hermes/), shared objects (base/), app tools (mcp/); doctrine/ is history
  security/    Zero-trust layer (platform-security): Kyverno, Falco, OpenFGA apps + labelled namespaces
  takserver/   TAK Server: Helm chart (chart/) + values-onemind.yaml (Helm CLI), Argo overlays deploy/ and fedhub/, runbook
  onetak/ omos/ data/ takml/ n8n/ backstage/ monitoring/ kube/ ops/ parked/ ...   (omos/ deploys to namespace feeds; parked/ is not deployed)
  ecosystem/   Argo-managed "ecosystem" stack: applications/, foundation-applications/, values/, charts/
  values/      Helm values for the add-ons (Argo-adopted: cert-manager, ingress-nginx, tailscale-operator, kube-prometheus-stack, loki, otel-collector; Helm CLI: argocd)
infra/         Terraform: DO VPC + DOKS cluster, NOT matching live, do not apply (infra/README.md); ecosystem-workers/ is a separate root; edge/ holds Mac leaf config (not deployed)
services/ docs/ scripts/                                code (built into images) and bus contracts (services/schemas/), docs, CI helpers
```

Not active or not deployed: `docs/references/` (external project research, moved from the top-level `references/` on 2026-10-17). Do not edit it unless asked. (The `onemind-vision` plugin, installed by hand in the Hermes pod, has its source in `platform/hermes/onemind-hermes/plugins/`.) The old Hermes agent's souls, registry, doctrine and runbook are kept, as history, in `platform/hermes/doctrine/`.

## The Deployment Pattern

1. Application manifests live under `platform/<layer>/` (or `platform/takserver/deploy`) and are applied by an Argo Application in `platform/argocd/` (ecosystem charts: `platform/ecosystem/applications/`). A path with no Application is not deployed.
2. Third-party charts are Argo sources with a pinned `targetRevision`; values live in `platform/ecosystem/values/` or `platform/values/`. Foundation charts are also checksum-locked in `platform/ecosystem/foundation-lock.json`.
3. Images come from GHCR (`ghcr.io/onemindos/...`) or upstream registries; there is no in-cluster registry (Harbor was removed). DOCR is not used. Running images were pinned to their observed digests on 2026-10-23 (`docs/operations/running-image-pins-2026-10-23.md`); a read-only check on 2026-10-27 found no running container on `:latest`. Only parked manifests (`platform/parked/`) still reference `:latest`. Pin a version or digest for anything new.
4. Secrets: SealedSecrets (below). Never `kubectl create secret` and leave it.
5. Every workload needs resource requests and limits. The hub-core node pool is memory-saturated (`docs/operations/cluster-review-2026-10-09.md`, H1); do not add requests casually.
6. Stateful claims use the `do-block-storage-retain` StorageClass: set `storageClassName` on PVCs and StatefulSet `volumeClaimTemplates`, or the chart or operator equivalent, so deleting the claim keeps the volume. Keep a StatefulSet's `persistentVolumeClaimRetentionPolicy` at `Retain`, the Kubernetes default; some charts set `Delete` (Loki's `enableStatefulSetAutoDeletePVC`), which deletes the claims along with the StatefulSet. The class cannot be changed on an existing claim or in a live StatefulSet's templates: apply this to new claims only, and never recreate a live claim to switch it (set its PersistentVolume's reclaim policy to `Retain` instead, with owner approval).
7. Temporary claims stay on `do-block-storage`, so cleanup frees their volumes. On every VolSync ReplicationSource, set `spec.restic.storageClassName` and `spec.restic.cacheStorageClassName` to `do-block-storage` explicitly, as `platform/backups/volsync/sources.yaml` does. Left unset, both fall back to the source claim's class: a retain-class snapshot copy then leaves a billed volume behind after every sync, and a retain-class cache outlives its ReplicationSource. Throwaway restore-drill claims also use `do-block-storage`. A restore that becomes the live claim is not a drill: give it the name, size and class its manifest declares in Git (`docs/portability/REBUILD-FROM-ZERO.md` step 9, `docs/operations/parked-services.md`). If that manifest is a standalone PVC still on `do-block-storage` and the claim does not exist yet, change the manifest to `do-block-storage-retain` before creating it; an existing claim or a live StatefulSet's templates keep their class (item 6).

### Public ingress

Public hostnames go through `ingress-nginx` (class `nginx`, LB `143.244.220.155`). DNS is Cloudflare-proxied A records pointing at that LB; ingresses carry the `whitelist-source-range` annotation with Cloudflare ranges. No public hostname lacks a source restriction (Matrix is gone: `matrix.` and `element.` were retired 2026-11-07; `zeusdelacruz.com` retired 2026-10-10). Four ingresses use class `nginx` (Authentik, the NATS WebSocket, the OneTAK host and the geo hosts). TAK has no public listener: `tak.onemindos.dev` and the DigitalOcean LB `152.42.155.95` are retired, and EUDs use `tak.tail717ba0.ts.net` over Tailscale. Admin UIs (Argo CD, TAK Portal at `enroll.tail717ba0.ts.net`, FedHub at `fedhub.tail717ba0.ts.net`, Grafana, Hermes dashboard, n8n, Headlamp, Hubble, Backstage portal and preview, Langfuse, Forgejo, Gatus and others) use the `tailscale` ingress class and are tailnet-only (OpenCost too, at `opencost.tail717ba0.ts.net`). The only NATS UI is NUI at `nui.tail717ba0.ts.net`, behind the Authentik proxy outpost `nui-auth`. Of the 26 tailscale-class ingresses, 24 share one Tailscale `ProxyGroup` named `ingress` (2 pods; `platform/core/tailscale-proxygroup-ingress.yaml`) through the annotation `tailscale.com/proxy-group: ingress`; `argocd` and `onemind-auth` keep their own proxy pods. Put new admin UIs in the group.

## NATS Subject Grammar (LOCKED — do not invent new patterns)

```
tel.{node}.{class}.{metric}          telemetry          tel.drone.x47.position
det.{domain}.{type}.{id}             detections         det.air.visual.track42
evt.{domain}.{entity}.{event}        semantic events    evt.mission.patrol1.completed
cmd.{domain}.{entity}.{action}       commands           cmd.drone.x47.intercept
intel.{domain}.{source}.{type}       intelligence       intel.conflict.acled.event
ent.{entity_type}.{entity_id}.state  fused entity state ent.ship.mmsi123.state
svc.{service}.{endpoint}             NATS RPC           svc.hermes.query
agents.{verb}.{agent}.{owner}.*      agent fabric       agents.prompt.hermes.zeus.legacy
tak.cot.inject                       CoT → TAK inject   (planned; no bridge is live)
tak.cot.>                            CoT from TAK       (cot-tap output)
```

Full doctrine: `docs/NATS-DOCTRINE.md`.

## NATS Streams

The Stream/KeyValue CRs in `platform/nats/streams/` are applied by the `platform-nats-streams` Application (sync-wave 2, `prune: false`, no finalizer, ServerSideApply), not by `platform-nats`, which does not recurse into `streams/`. The NACK controller (Helm release `nack`, adopted into Argo as `edge-nack`) reconciles them into JetStream: 38 CRs, all R3 file storage, all with `preventDelete: true`, so deleting a CR never deletes the stream or its data. `kv-feed-config.yaml` is excluded from the Application. Its live bucket is R1 and was created by an app, and it must be matched to `nats kv info feed-config` before adoption. NACK runs with `-crd-connect`: each CR's `spec.servers` (set out of band, never in Git) supplies the connection. A new stream is a new CR file here, never a `nats stream add`. Do not use the older Python provisioning scripts as the source of truth.

## Secrets

All committed secrets are **Bitnami SealedSecrets**; the controller runs in `kube-system` (managed by `platform-core`). Cluster private keys must be backed up independently of the cluster (they are in the encrypted `backup-cluster-recovery` export).

```bash
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem
kubectl create secret generic <name> -n <namespace> --from-literal=KEY=value \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/<layer>/<name>.sealedsecret.yaml
```

Place the file beside the workload that uses it (existing examples: `platform/hermes/base/`, `platform/backups/sealed/`, `platform/ecosystem/resources/secrets/`). Argo applies it. Exceptions that are still out of band: the NATS Helm `values.secret.yaml` (gitignored) and the CloudTAK/agent API tokens generated in a UI.

- Never commit plaintext credentials, even test values or "example" commands. CI runs gitleaks and fails on new findings; `.gitleaks-baseline.json` records the nine pre-existing hits (see `docs/STACK.md`), which need rotation, not suppression.
- Hermes: `onemind-hermes` reads `onemind-hermes-runtime`, the Zoho and gateway Secrets, and the shared `hermes-dashboard-auth`, `omos-nats-creds` and `hermes-nats-ca` in `platform/hermes/base/`. The model provider is chosen in the agent's runtime config, not in Git. The old `hermes-env` (`DIGITALOCEAN_AI_KEY`, legacy `AWS_*` Bedrock keys, 1Password and Todoist tokens) left Git with `hermes-legacy` on 2026-11-03 and its credentials still need rotation (`platform/hermes/onemind-hermes/RETIREMENT-AUDIT.md`, V04). Do not describe Hermes as "on Bedrock".

## Cluster Access

```bash
doctl kubernetes cluster kubeconfig save omos-cloud     # context: do-nyc3-omos-cloud
kubectl --context do-nyc3-omos-cloud get nodes
```

- The Kubernetes API endpoint is **public** (no control-plane firewall; `docs/operations/cluster-review-2026-10-09.md`, H3). Tailscale gives access to admin UIs, not a private API.
- Cluster `omos-cloud` (nyc3, DOKS 1.36): pools `onemind-hub-core-v2` (3 x 8 vCPU / 16 GiB) and `onemind-ecosystem-cpu` (8 vCPU / 32 GiB, 1 node live; 4 nodes in total, verified with `kubectl` and `doctl` 2026-11-08). `infra/main.tf` does not match and has no state; do not `terraform apply` it (see `infra/README.md`).
- Volumes on 2026-11-08 (verified with `kubectl get pvc -A`): NATS 3 x 20 GiB; `onemind-data` (CNPG) 3 x 5 GiB on PostgreSQL 16 (a PostgreSQL 17 upgrade was tried and rolled back on 2026-11-08 because Apache AGE blocks `pg_upgrade`); Hermes 10 GiB; TAK PostgreSQL 10 GiB; ClickHouse 20 GiB on 26.8.36.3 with a 7-day TTL on its system logs and a CronJob that prunes native backups after 14 days; Loki 30 GiB; Prometheus 20 GiB. OneTAK (nodeodm 10, tiles 20, tile38 2), Karakeep (5 + 2), n8n (2 + 3) and Forgejo (5) disks were shrunk on 2026-11-08.
- Monitoring: Grafana signs in through Authentik (the local form stays on) and has the ClickHouse datasource plugin. Falco sends to Alertmanager and Loki. Alerts reach the owner through the Hermes `infra-health` cron (every 15 minutes, delivered on Telegram) and the portal inbox (Alertmanager receiver `onemind-backstage`). The OTel collector ships logs only. Velero backs up the Kubernetes objects of every app namespace daily; its node agent is off, so volume data relies on VolSync and the database backups.
- Older local contexts such as `do-nyc3-onemind-hub` may exist. Check the context before running anything.

## TAK Server

TAK Server 5.8 on OneMind hardened images (`ghcr.io/onemindos/takserver-base:*-5.8-onemind`). The core (API, messaging, config, plugins, Ignite, PostgreSQL) is the **Helm release `takserver`** from `platform/takserver/chart` with `platform/takserver/values-onemind.yaml`, applied by hand (revision 34 on 2026-09-29). The `platform-takserver` Application manages only the overlay in `platform/takserver/deploy` (LB proxy, NATS leaf, Mumble, LDAP, secrets) and `platform-fedhub` manages `platform/takserver/fedhub`. Both fedhub and tak-ldap are running. TAK is reachable only over the tailnet: the Tailscale Service `tak-tailnet` (node `tak`, `tak.tail717ba0.ts.net`) exposes 8089, 8443, 8446 and 64738/TCP. The public LB `tak-external` was removed on 2026-10-23, and federation ports 9101/9102 are not exposed.

Since 2026-10-29 the CA-key, server-keystore and LDAP bind passwords are Secrets (`tak-ca-credentials`, `tak-server-stores`, `tak-ldap-creds`); `CoreConfig.xml` is the password-free template `platform/takserver/deploy/configmap-core-config-v2.yaml` and the config/messaging init containers fill it in (Helm stayed at revision 34; see `platform/takserver/DEPLOY-RUNBOOK.md`, "Password supply").

The NatsBridgePlugin (`services/takserver/plugin-src/`, JAR in `services/takserver/plugin/`) was meant to be the CoT-to-NATS wire (`tak.cot.inject` in, `tak.cot.>` out), but it is **not installed**: the live `takserver-plugins` pod holds no NATS plugin, and `tak-nats-leaf` has 0 client connections (`/varz`, 2026-11-08). So the plugin wire carries nothing today. PR #739 (feed-platform design, still open) is the plan, and the TAK-client bridge `platform-takserver-tak-client` (manual sync, not synced) is the candidate. Storage: TAK PostgreSQL is 10 GiB. Images are built by `.github/workflows/build-takserver.yml` on the self-hosted runner. Deploy runbook: `platform/takserver/DEPLOY-RUNBOOK.md`.

## CI (`.github/workflows/quality.yml`)

`quality.yml` and `ecosystem-charts.yml` run on GitHub-hosted runners (`ubuntu-latest`); the self-hosted runner in `github-runner` keeps the cluster-facing and image-build jobs (`build-*.yml`, inventory refresh, MCP). On every PR to `main`: changed-YAML lint, ShellCheck, Terraform fmt/validate, kubeconform on changed plain manifests (strict, CRD schemas from datreeio), `helm template` of changed Argo Helm Applications and local charts, and gitleaks (working tree against the baseline, and the commits in the PR). A weekly run re-checks the whole tree. Actions are pinned by commit SHA; Renovate (`renovate.json`) proposes weekly updates. Run the same checks locally with `make validate` plus `python scripts/validate_manifests.py --all` and `python scripts/validate_helm_renders.py --all`.

## Adding a New Service (checklist)

- [ ] Image built and pushed to GHCR with a pinned tag (not `:latest`)
- [ ] Manifests under the right `platform/<layer>/` directory, with requests and limits
- [ ] Argo Application created or extended in `platform/argocd/` (sync wave set); not applied by hand
- [ ] SealedSecret committed beside it; no plaintext
- [ ] Ingress: `nginx` class with the Cloudflare allowlist for public, or `tailscale` class for admin
- [ ] Backup decided: PVCs get a VolSync source (`platform/backups/volsync/`), databases a dump job or CNPG
- [ ] Row added to `docs/STACK.md` with a tier
- [ ] NATS subjects follow the grammar above

## What NOT to Do

- Never commit secrets or credentials (even test values) to this repo
- Never deploy imperatively without a matching manifest here; never add a new Helm-CLI release
- Never create a new NATS subject pattern outside the grammar above
- Never use DOCR; images are on GHCR
- Never point a cloud workload at a Mac or edge NATS; cloud services use the in-cluster address (there is no exception any more)
- Never delete PVCs, namespaces or backup data without a fresh snapshot and explicit approval
- Never add services that duplicate an existing layer; check `docs/STACK.md` first
