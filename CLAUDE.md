# OneMind Infra — Claude Context

This repo is the infrastructure-as-code and GitOps source of truth for the OneMind platform: the DOKS (DigitalOcean Kubernetes) cloud hub, platform manifests, NATS stream definitions, TAK Server chart and images, and Terraform for the DigitalOcean foundation.

Read `docs/MISSION.md` before architecture, product, AI or infrastructure decisions. Live state is in `docs/STACK.md` (tiered inventory, generated from the cluster); directory purposes are in `docs/REPO-LAYOUT.md`. Distinguish verified live state from intent: verify against the cluster (read-only) before writing "live" in docs.

The platform plan (four layers, AI trust model, OpenChoreo migration phases and their tracker) is in `PLATFORM-ROADMAP.md` at the repository root. Read it before platform or AI-operations work, and update its tracker in the same PR when you complete roadmap work.

## Owner-approved merge batch — 2026-10-27

Zeus approved all owner decisions for dcarku's 46 PRs #495–#540, retained
Telegram, and authorized Claude Code to review, fix, validate and execute those
merges on his behalf, including this batch's foundation PRs. See the exact
[approval, decisions and continuation](docs/roadmap/2026-10-27-owner-merge-preparation.md#owner-authorization--2026-10-27)
in PR #542. This satisfies the default human-merge/owner-approval requirement
only for this batch; technical checks, credential and release prerequisites
remain. Do not request the settled approvals again. Hold an actually blocked PR,
continue the queue and report remaining blockers together at the end.

**Latest owner direction (2026-10-27): Hermes is the lead and sole merger.**
Start with #542, #504, #497, #498 and #496 sequentially, reviewing each against
current `main` and verifying its final CI before marking ready and merging.
The existing second harness prepares/tests #505 and application PRs; it does not
merge. These first five are source merges only: do not run Terraform apply or
execute the proposed storage patches. Keep Telegram and the settled approvals.

## What This Repo Is

If a service runs in the cluster, its manifest should live here, and Argo CD should apply it. Reality is not there yet: 59 Argo Applications reconcile from `main` (2026-10-27). Of the 11 Helm releases first installed by hand, 8 are adopted into Argo: `edge-cert-manager`, `edge-ingress-nginx`, `edge-nack`, `edge-tailscale-operator` (PR #176) and `monitoring-kube-prometheus-stack`, `monitoring-loki`, `monitoring-tempo`, `monitoring-otel-collector` (2026-10-17). **3 are still Helm CLI only**: nats, takserver and argocd. Never `helm upgrade` an adopted release, and do not run `platform/install.sh`, which still upgrades all of them; see `docs/gitops/HELM-ADOPTION.md`. Also a few namespaces (`matrix`, the NATS web UIs) have no Application. See "GitOps gaps" in `docs/STACK.md`. Do not add to that list: new workloads get an Argo Application, not `helm install` or `kubectl apply`.

The platform is a NATS-native C2/ISR system. NATS is the nervous system, TAK (ATAK/iTAK/WebTAK) is the tactical display, and Hermes ("Legacy") is the AI agent.

## Bus Authority (read first)

**The in-cluster NATS (3-node JetStream, namespace `nats`) is the authority bus. Everything publishes to it.**

- Cloud workloads connect directly to `nats://nats.nats.svc.cluster.local:4222` (a few use `tls://` on the same host:port) with credentials from a Secret (`NATS_USER` / `NATS_PASSWORD`). No sidecar.
- **Hermes is the exception**: its pod runs a `nats-leaf` sidecar and the container uses `NATS_URL=nats://localhost:4222`, so the agent has its own bus identity. Do not copy that pattern into new services.
- **Edge and Mac bridges** (MAVLink, Home Assistant, cameras) run a local leaf node that forwards to the cluster through the `omos-nats` Tailscale load balancer (leaf port 7422). The leaf is a forwarding pipe, not a peer bus. If the Mac is offline the cloud bus keeps running.
- When the home server arrives, it becomes another leaf (and eventually the authority). Do not design anything that treats a single laptop's NATS as the source of truth.

```
Edge (Mac / home server)   MAVLink, HA, cameras ── NATS leaf ──►  K8s NATS (3-node)  ◄── all cloud services
                                                  (forward only)       THE BUS            incl. Hermes via its own leaf
```

## Repo Layout (short; full table in `docs/REPO-LAYOUT.md`)

```
platform/      Everything Argo CD syncs. argocd/ holds the Application CRs; the rest is grouped by layer
  argocd/      App-of-apps + 25 platform-* Applications
  core/        sealed-secrets controller, GitHub runner, Tailscale policy
  nats/        NATS manifests, TLS, 39 Stream/KeyValue CRs in streams/, Helm values templates
  auth/        Authentik infra (Postgres, Redis, ingress); the chart is in argocd/platform-auth-helm.yaml
  backups/     PG/NATS backup CronJobs, cluster-recovery export, VolSync sources (volsync/)
  hermes/      Hermes agent (agent/) and its sealed secrets (bridge/)
  security/    Zero-trust layer (platform-security): Kyverno, Falco, OpenFGA apps + labelled namespaces
  takserver/   TAK Server: Helm chart (chart/) + values-onemind.yaml (Helm CLI), Argo overlays deploy/ and fedhub/, runbook
  onetak/ omos/ data/ takml/ n8n/ backstage/ monitoring/ kube/ ops/ ...
  ecosystem/   Argo-managed "ecosystem" stack: applications/, foundation-applications/, values/, charts/
  values/      Helm values for the add-ons (Argo-adopted: cert-manager, ingress-nginx, tailscale-operator, kube-prometheus-stack, loki, tempo, otel-collector; Helm CLI: argocd)
infra/         Terraform: DO VPC + DOKS cluster, NOT matching live, do not apply (infra/README.md); ecosystem-workers/ is a separate root; edge/ holds Mac leaf config (not deployed)
services/ docs/ scripts/                                code (built into images) and bus contracts (services/schemas/), docs, CI helpers
```

Not active or not deployed: `docs/references/` (external project research, moved from the top-level `references/` on 2026-10-17). Do not edit it unless asked. (The `onemind-vision` plugin, installed by hand in the Hermes pod, has its source in `platform/hermes/agent/plugins/`.) The Hermes agent's souls, registry, doctrine and runbook are in `platform/hermes/agent/`.

## The Deployment Pattern

1. Application manifests live under `platform/<layer>/` (or `platform/takserver/deploy`) and are applied by an Argo Application in `platform/argocd/` (ecosystem charts: `platform/ecosystem/applications/`). A path with no Application is not deployed.
2. Third-party charts are Argo sources with a pinned `targetRevision`; values live in `platform/ecosystem/values/` or `platform/values/`. Foundation charts are also checksum-locked in `platform/ecosystem/foundation-lock.json`.
3. Images come from GHCR (`ghcr.io/onemindos/...`) or upstream registries; there is no in-cluster registry (Harbor was removed). DOCR is not used. Running images were pinned to their observed digests on 2026-10-23 (`docs/operations/running-image-pins-2026-10-23.md`); a read-only check on 2026-10-27 found no running container on `:latest`. Only parked manifests (`platform/parked/`) still reference `:latest`. Pin a version or digest for anything new.
4. Secrets: SealedSecrets (below). Never `kubectl create secret` and leave it.
5. Every workload needs resource requests and limits. The hub-core node pool is memory-saturated (`docs/operations/cluster-review-2026-10-09.md`, H1); do not add requests casually.

### Public ingress

Public hostnames go through `ingress-nginx` (class `nginx`, LB `143.244.220.155`). DNS is Cloudflare-proxied A records pointing at that LB; ingresses carry the `whitelist-source-range` annotation with Cloudflare ranges. Exceptions today with no source restriction: `matrix.`/`element.onemindos.com`. (`zeusdelacruz.com` retired 2026-10-10.) TAK has no public listener: `tak.onemindos.dev` and the DigitalOcean LB `152.42.155.95` are retired, and EUDs use `tak.tail717ba0.ts.net` over Tailscale. Admin UIs (Argo CD, TAK Portal at `enroll.tail717ba0.ts.net`, FedHub at `fedhub.tail717ba0.ts.net`, Grafana, Hermes dashboard, n8n, Headlamp, Hubble, NATS UIs, Backstage preview) use the `tailscale` ingress class and are tailnet-only (OpenCost too, at `opencost.tail717ba0.ts.net`).

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
tak.cot.inject                       CoT → TAK inject   (NatsBridgePlugin input)
tak.cot.>                            CoT from TAK       (cot-tap output)
```

Full doctrine: `docs/NATS-DOCTRINE.md`.

## NATS Streams

The 39 Stream/KeyValue CRs in `platform/nats/streams/` are applied by the `platform-nats` Application and reconciled into JetStream by the NACK controller (Helm release `nack`, adopted into Argo as `edge-nack`). Do not use the older Python provisioning scripts as the source of truth.

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

Place the file beside the workload that uses it (existing examples: `platform/hermes/bridge/`, `platform/backups/sealed/`, `platform/ecosystem/resources/secrets/`). Argo applies it. Exceptions that are still out of band: the NATS Helm `values.secret.yaml` (gitignored) and the CloudTAK/agent API tokens generated in a UI.

- Never commit plaintext credentials, even test values or "example" commands. CI runs gitleaks and fails on new findings; `.gitleaks-baseline.json` records the nine pre-existing hits (see `docs/STACK.md`), which need rotation, not suppression.
- Hermes: `hermes-env` holds `DIGITALOCEAN_AI_KEY` (the default model provider per the cluster review) and legacy `AWS_*` Bedrock keys. The provider is chosen in the agent's runtime config, not in Git. Do not describe Hermes as "on Bedrock"; remove the AWS keys if confirmed unused.

## Cluster Access

```bash
doctl kubernetes cluster kubeconfig save omos-cloud     # context: do-nyc3-omos-cloud
kubectl --context do-nyc3-omos-cloud get nodes
```

- The Kubernetes API endpoint is **public** (no control-plane firewall; `docs/operations/cluster-review-2026-10-09.md`, H3). Tailscale gives access to admin UIs, not a private API.
- Cluster `omos-cloud` (nyc3, DOKS 1.36): pools `onemind-hub-core-v2` (3 x 8 vCPU / 16 GiB) and `onemind-ecosystem-cpu` (8 vCPU / 32 GiB, 1 node, autoscale 1 to 4; 2026-10-27). `infra/main.tf` does not match and has no state; do not `terraform apply` it (see `infra/README.md`).
- Older local contexts such as `do-nyc3-onemind-hub` may exist. Check the context before running anything.

## TAK Server

TAK Server 5.8 on OneMind hardened images (`ghcr.io/onemindos/takserver-base:*-5.8-onemind`). The core (API, messaging, config, plugins, Ignite, PostgreSQL) is the **Helm release `takserver`** from `platform/takserver/chart` with `platform/takserver/values-onemind.yaml`, applied by hand (revision 34 on 2026-09-29). The `platform-takserver` Application manages only the overlay in `platform/takserver/deploy` (LB proxy, NATS leaf, Mumble, LDAP, secrets) and `platform-fedhub` manages `platform/takserver/fedhub`. Both fedhub and tak-ldap are running. TAK is reachable only over the tailnet: the Tailscale Service `tak-tailnet` (node `tak`, `tak.tail717ba0.ts.net`) exposes 8089, 8443, 8446 and 64738/TCP. The public LB `tak-external` was removed on 2026-10-23, and federation ports 9101/9102 are not exposed.

The NatsBridgePlugin (`services/takserver/plugin-src/`, JAR in `services/takserver/plugin/`) is the CoT-to-NATS wire (`tak.cot.inject` in, `tak.cot.>` out). Images are built by `.github/workflows/build-takserver.yml` on the self-hosted runner. Deploy runbook: `platform/takserver/DEPLOY-RUNBOOK.md`.

## CI (`.github/workflows/quality.yml`)

On every PR to `main`: changed-YAML lint, ShellCheck, Terraform fmt/validate, kubeconform on changed plain manifests (strict, CRD schemas from datreeio), `helm template` of changed Argo Helm Applications and local charts, and gitleaks (working tree against the baseline, and the commits in the PR). A weekly run re-checks the whole tree. Actions are pinned by commit SHA; Renovate (`renovate.json`) proposes weekly updates. Run the same checks locally with `make validate` plus `python scripts/validate_manifests.py --all` and `python scripts/validate_helm_renders.py --all`.

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
- Never point a cloud workload at a Mac or edge NATS; cloud services use the in-cluster address (Hermes' sidecar is the one documented exception)
- Never delete PVCs, namespaces or backup data without a fresh snapshot and explicit approval
- Never add services that duplicate an existing layer; check `docs/STACK.md` first
