# OneMind platform roadmap: layers, OpenChoreo and AI operations

Living document. Started 2026-10-18. Owner: Zeus. Kept at the repository root so any
agent (Claude Code, Codex or another) can pick up the work where the last one stopped.

2026-10-27 user-directed merge preparation: [portal alias #524](https://github.com/onemindos/omos/pull/524)
keeps `onemind-preview` canonical and only redirects the proposed `portal` hostname.
The current alias requires no new OIDC callback; a later canonical-host change
does. Source/review preparation is separate from hostname registration, Argo
rollout and browser acceptance. Phase 6 remains In progress; Phases 7/8 stay
Not started. No identity-provider mutation or deployment is recorded here.

2026-10-26 user-directed work: launch an independent `hermes-fresh` beside Legacy
through `platform-hermes-fresh`, with new storage and no imported agent state or
operational credentials. Source and isolated image smoke checks are prepared;
cluster rollout, private access and model sign-in require live acceptance. This
does not mark the broader agent migration or Phase 8 complete. See
[`fresh Hermes`](platform/hermes/onemind-hermes/README.md); deployment PR records live evidence.
2026-10-29: renamed to `onemind-hermes` (Application `platform-onemind-hermes`,
persona Legacy as a single soul) with new retain-class storage and hardened
initial config; the `hermes-fresh` resources were removed the same day except the PVC `hermes-fresh-data`, kept pending owner approval.

2026-10-31: Legacy's "How you work" doctrine (one soul, four modes, delegate when big or parallel, checker gate)
and five Codex worker roles (`ops-cluster`, `tak-mission`, `builder`, `researcher`, `checker`) prepared as a
`legacy-roles` delegation backend plus kanban-swarm profiles; applied to the live pod by
`platform/hermes/onemind-hermes/roles/apply.sh` after merge, live verification recorded in the PR.
See the [README section](platform/hermes/onemind-hermes/README.md#how-legacy-works-doctrine-and-worker-roles).

2026-11-01: Legacy's NATS presence (step 7 of the Legacy to onemind-hermes migration) prepared in Git, not applied:
onemind-hermes as a direct client of the cluster bus with no leaf sidecar, `nats-platform` vendored without the
start-up hot-patch, a duplicate-identity guard so it can be enabled before `hermes-legacy` is stopped, and the
`nats-clients` ingress rule. `platforms/tak-platform` is not ported (PR #585 replaces it). Inventory, order,
checks and rollback: [MIGRATION.md](platform/hermes/onemind-hermes/MIGRATION.md#nats-inventory-and-cutover-step-7).
Phase 8 stays Not started; live cutover is recorded in the PR after merge.

2026-10-24 sprint: prepared a dedicated read-only Hermes cluster observer identity
and sanitized readiness collector for the owner's Paperclip status-report task.
Argo activation, RBAC acceptance and the agent-written report remain unverified.

This is a plan, not a record of live state. Live state is in [`docs/STACK.md`](docs/STACK.md);
the mission it serves is in [`docs/MISSION.md`](docs/MISSION.md).

> Taking over from another agent? Start with [`docs/roadmap/HANDOFF.md`](docs/roadmap/HANDOFF.md)
> (state, rules, prompt) and the phase briefs in `docs/roadmap/`.

2026-11-01 owner-requested Zoho onboarding: nine Connections-mode MCP endpoints
are prepared as SealedSecrets with environment-only URL references, exact read
allowlists, and an on-start Hermes CLI reconciler for existing and clean PVCs.
Canonical Projects name: `zoho-projects`; WorkDrive remains separate. No new
workload or scheduler requests. Tests and source delivery do not establish live
unsealing or app authorization; [acceptance and recovery](platform/hermes/onemind-hermes/mcp/zoho/README.md)
remain deployment gates. Broader Phase 8 is not marked complete by this change.

2026-11-01 owner-requested native 1Password access: **source-only draft, NOT
DEPLOYED**. Default-home configuration, checksum-pinned CLI startup bootstrap,
ClickUp environment scoping and proposed wiring are retained only in the non-active
`onepassword/` directory. Production Deployment, defaults and Kustomization are
restored to original base bytes; no installer/bootstrap/source activation is wired. The encrypted bootstrap SealedSecret, corrected Hermes runtime image (the
pinned upstream filter leaks bootstrap to terminal children), and authorized
runtime/recovery checks remain integration gates; no deployed access or broader
Phase 8 completion is claimed. The local hash-locked runtime correction and
synthetic regressions are in
[runtime-fix](platform/hermes/onemind-hermes/onepassword/runtime-fix/README.md).
See [startup prerequisites and recovery](platform/hermes/onemind-hermes/onepassword/README.md).

## How to work on this roadmap (any agent)

1. Read this file, `CLAUDE.md` (operating rules) and `docs/MISSION.md` first.
2. Find the first phase in the [tracker](#phase-tracker) that is not `Done`. Work only on
   that phase unless the tracker says two phases can run in parallel.
3. Every phase is one or more pull requests. Open the PR on a branch named
   `roadmap/phase-<n>-<short-name>`, and in the same PR update this file:
   set the phase status, add the PR link, and tick the finished tasks.
4. A phase is `Done` only when its gate is checked against the live cluster (read-only)
   and the evidence is written under the phase. Do not mark a gate passed from the
   manifests alone.
5. Follow the [layer rules](#the-four-layers-and-the-trust-model): never move a Layer 1
   workload into OpenChoreo, and never give an AI write access above its layer.
6. Record new decisions in the [decision log](#decision-log) with the date.

Status values: `Not started`, `In progress`, `Blocked`, `Done`.

## Why

- One person runs a 4-node DOKS cluster with 60 Argo CD Applications, 43 namespaces and
  about 250 pods. Every app has hand-written Deployment, Service, ingress, secret and
  network YAML, each slightly different.
- The AI operator (Claude Code today) cannot see the cluster in one place: the catalog is in
  Backstage, deploy state in Argo CD, logs and metrics in Grafana.
- OpenChoreo 1.3.1 is installed but only its control plane runs. It has no data,
  workflow or observability plane, and manages no projects or components.

The goal: one platform layer for our own apps, a stable foundation underneath it, and an
AI that can see every layer and act within clear limits.

## The four layers and the trust model

| Layer | What it is | Managed by | What AI may do |
|---|---|---|---|
| **3. Apps** | Software we write or run as products | OpenChoreo Components (Projects group them) | See and act: restart, scale, roll back a release |
| **2. Shared services** | Data, messaging, AI plumbing and observability that apps use | Argo CD; modelled in OpenChoreo as ResourceTypes and observability-plane modules | See, and propose changes as pull requests |
| **1. Foundation** | Everything needed to log in, reach, secure, back up and rebuild the platform, including OpenChoreo itself | Plain Argo CD (break-glass layer) | See only. May draft a pull request; a human merges |
| **0. Kubernetes** | DOKS-managed system namespaces and the Cilium CNI | DigitalOcean | See only |

Rule behind it: the thing that fixes problems must not depend on the thing it is fixing.
If an AI could break Authentik, Tailscale, cert-manager, Velero or the policy engines,
it could lock out both itself and the human who has to recover it.

**Target clarification (2026-10-21):** the table above records the historical
OpenChoreo baseline, not authorization for all listed AI actions. The user's
selected future KubeVela/Crossplane ownership split is recorded in
[decision 0002](docs/architecture/decisions/0002-application-delivery-and-infrastructure-boundaries.md);
implementation/coexistence acceptance and per-app transfers remain pending.
Current foundation/shared-service and machine read-only boundaries are unchanged.

**Bounded owner authorization (2026-10-27):** Zeus approved the owner decisions
for dcarku's 46 PRs **#495–#540**, retained Telegram and authorized Claude Code
to execute these reviewed GitHub merges on his behalf, including foundation
PRs. This satisfies the table's human-merge requirement for this batch only;
it is not standing authority for unrelated changes. The
[approval record and technical merge queue](docs/roadmap/2026-10-27-owner-merge-preparation.md)
keep CI, credential, dependency and rollout gates explicit. Do not ask again for
settled choices; continue past an actually blocked PR and report it at the end.

**Latest merge coordination (2026-10-27): Hermes is the lead and sole merger.**
The initial source-only queue is #542 → #504 → #497 → #498 → #496, with fresh
review/current-main integration and final-head CI for each. The existing second
harness prepares/tests #505 and application changes without merging. This initial
batch does not authorize Terraform apply or execution of the proposed storage
patches; Telegram remains selected. See the updated approval record below.

### Namespace map (live cluster, 2026-10-18)

| Namespace | Main workloads | Layer | Target |
|---|---|---|---|
| `kube-system`, `kube-node-lease`, `kube-public`, `default` | DOKS system, Cilium | 0 | DigitalOcean |
| `argocd` | Argo CD | 1 | Argo CD (Helm CLI today, see `CLAUDE.md`) |
| `auth` | Authentik | 1 | Argo CD |
| `tailscale` | Tailscale operator and proxies | 1 | Argo CD |
| `cert-manager` | cert-manager | 1 | Argo CD |
| `ingress-nginx` | public ingress | 1 | Argo CD |
| `kyverno`, `falco` | policy, runtime security | 1 | Argo CD |
| `velero`, `ops` | Velero, VolSync | 1 | Argo CD |
| `github-runner` | CI runners | 1 | Argo CD |
| `kube` | Headlamp (break-glass UI) | 1 | Argo CD |
| `backstage` | OneMind portal (the human home screen) | 1 | Argo CD |
| `openchoreo-control-plane`, `openchoreo-data-plane`, `openchoreo-workflow-plane`, `openchoreo-observability-plane` | OpenChoreo itself | 1 | Argo CD |
| `takserver` | TAK Server core | 1 for now | Argo CD; Helm CLI release today. Revisit after Phase 7 |
| `nats` | 3-node JetStream bus | 2 | Argo CD; streams as an OpenChoreo ResourceType |
| `cnpg-system`, `onemind-data` | CloudNativePG operator and the `onemind-data` cluster | 2 | Argo CD; databases as an OpenChoreo ResourceType |
| `clickhouse` | ClickHouse | 2 | Argo CD; ResourceType later |
| `onemind-search` | OpenSearch | 2 | Argo CD; also the observability plane's log store |
| `onemind-ai` | LiteLLM, CPU model | 2 | Argo CD. LiteLLM stays the model gateway |
| `dapr-system` | Dapr | 2 | Argo CD |
| `keda` | KEDA | 2 | Argo CD; OpenChoreo autoscaling module |
| `monitoring` | Prometheus, Alertmanager, Grafana, Loki, Tempo, OTel collector | 2 | Argo CD; becomes the observability plane's metrics and events modules |
| `opencost` | OpenCost | 2 | Argo CD; OpenChoreo FinOps module |
| `omos` | feed-intel, feed-satellites, omos-tak-projector, redpanda-connect, sink-clickhouse | 3 | OpenChoreo Project `omos-feeds` |
| `onetak` | 19 geo and mapping apps (Valhalla, Martin, Nominatim, WebODM, STAC, TileServer and others) | 3 | OpenChoreo Project `onetak` |
| `takml` | takml-server, takml-vision-worker, takml-db | 3 | OpenChoreo Project `takml` |
| `onemind-core` | onemind-core, onemind-ingestion, OpenFGA | 3 | OpenChoreo Project `onemind` |
| `agent-fabric` | agent-console, agent-tap | 3 | OpenChoreo Project `agent-fabric` |
| `n8n` | n8n | 3 | OpenChoreo Component |
| `matrix` | Synapse, Element | 3 | OpenChoreo Project `matrix` (Argo Application `platform-matrix` is OutOfSync today) |
| `openclaw` | OpenClaw | 3 | OpenChoreo Component (the OpenClaw Deployment itself is not in Git, see `docs/REPO-LAYOUT.md`) |
| `hermes` | Hermes ("Legacy") agent | 3 | OpenChoreo Component, migrated last |

## Current state (verified read-only, 2026-10-18)

- **OpenChoreo:** chart `openchoreo-control-plane` 1.3.1 via Argo Application
  `ecosystem-oc-control` (values in `platform/ecosystem/values/oc-control.yaml`). Five pods
  run. The data, workflow and observability plane namespaces are empty. Zero planes,
  environments, projects or components are registered. Its built-in Backstage is disabled
  (`backstage.enabled: false`). Identity is already Authentik
  (`onemind-openchoreo` application). The event forwarder already targets our portal at
  `/api/events/http/openchoreo`.
- **Portal:** `services/onemind-portal`, Backstage 1.55.3 on the new frontend and backend
  systems, two replicas, Postgres with encrypted backups, read-only MCP endpoint. Custom
  features are written as proper Backstage plugins and modules but live inside
  `packages/app` and `packages/backend`; `plugins/` is empty. No OpenChoreo plugins
  are installed.
- **AI:** Claude Code is the active operator. Hermes runs but has no model provider
  and no Kubernetes RBAC; it is not a priority.
- **Capacity:** two of the three `onemind-hub-core-v2` nodes have 96% of memory requested.
  The `onemind-ecosystem-cpu` node is at 53%.
- **Gateway:** Cilium provides Gateway API (`cilium`, `onemind-private` classes).
  OpenChoreo's docs assume kgateway.
- **Secrets:** SealedSecrets. External Secrets CRDs exist but no controller runs.
  OpenChoreo planes expect an External Secrets `ClusterSecretStore`.

## Phase tracker

| Phase | Goal | Status | PRs | Parallel with |
|---|---|---|---|---|
| 0 | Blueprint: this document | Done | [#256](https://github.com/onemindos/omos/pull/256) | none |
| 1 | Claude sees every layer (MCP on the Mac) | Done | [#257](https://github.com/onemindos/omos/pull/257), [#258](https://github.com/onemindos/omos/pull/258) | 2 |
| 2 | OpenChoreo foundations: data plane, gateway, secrets, first template, built-in console | Done | 2a [#262](https://github.com/onemindos/omos/pull/262), 2b [#263](https://github.com/onemindos/omos/pull/263), 2c [#264](https://github.com/onemindos/omos/pull/264), 2d [#265](https://github.com/onemindos/omos/pull/265), 2e [#266](https://github.com/onemindos/omos/pull/266), 2g-1 [#268](https://github.com/onemindos/omos/pull/268), [#269](https://github.com/onemindos/omos/pull/269), 2g-2 [#270](https://github.com/onemindos/omos/pull/270) [#271](https://github.com/onemindos/omos/pull/271) [#272](https://github.com/onemindos/omos/pull/272) [#273](https://github.com/onemindos/omos/pull/273), 2g-3 [#274](https://github.com/onemindos/omos/pull/274), 2f-1 [#275](https://github.com/onemindos/omos/pull/275), 2f-2 [#276](https://github.com/onemindos/omos/pull/276), (this PR) | 3, 4 |
| 3 | Pilot one app as an OpenChoreo Component | Done | [#279](https://github.com/onemindos/omos/pull/279), [#286](https://github.com/onemindos/omos/pull/286), [#287](https://github.com/onemindos/omos/pull/287), [#288](https://github.com/onemindos/omos/pull/288), [#289](https://github.com/onemindos/omos/pull/289), [#290](https://github.com/onemindos/omos/pull/290), [#291](https://github.com/onemindos/omos/pull/291), [#292](https://github.com/onemindos/omos/pull/292), [#293](https://github.com/onemindos/omos/pull/293), [#294](https://github.com/onemindos/omos/pull/294) | 4 |
| 4 | Portal refactor: custom features into `plugins/` packages | Done | [#280](https://github.com/onemindos/omos/pull/280) | 2, 3 |
| 5 | Capacity fixes, then the observability plane | Done | [#281](https://github.com/onemindos/omos/pull/281), implementation #300–#310 (individual links and live evidence below) | none |
| 6 | One portal: OpenChoreo plugins in our Backstage, built-in console off | In progress | Historical acceptance #282, #318, #319; reopened for [portal stabilization](docs/operations/portal-stabilization-2026-10-20.md) before migration (#322 source, #323 release, #324 inventory, #325 composition, #326 final image); 2/2 healthy, 715 processed entities; broader acceptance gaps remain documented | none |
| 7 | Migrate Layer 3 apps project by project; NATS and Postgres as ResourceTypes | Not started | | 6 |
| 8 | Expand: home server data plane, Hermes on the same tools | Not started | | none |
| 8a | Tool gateway for Legacy: agentgateway front door, tool groups, execute gate ([design](docs/architecture/onemind-tool-gateway.md), [runbook](docs/operations/tool-gateway.md)) | In progress; verified live read-only 2026-10-30 (Argo Synced and Healthy, backends Ready). **Phase 0 spike and 1 gateway core: Done. Phase 2 policy plane: Done. 3a Kubernetes: Done. 3d data (STAC, onemind-core, ClickHouse, Tile38): Done, keyless. 3e TAK: partial, MediaMTX only; CloudTAK, TAK Server, TAK Portal and NATS publish wait on owner credentials. 3b Argo CD and Grafana, 3c GitHub, DigitalOcean and Cloudflare, 3f Work, 3g Home: wait on owner keys. 4 Execute gate: code merged, not deployed (Matrix bot pending); `/execute/mcp` stays deny-all. 5 Hermes cutover: `ops`, `data` and `tak` wired in `onemind-hermes` (PR open 2026-10-31, live verification after merge); `dev`, `cloud`, `home`, `work` and `execute` pending their backends; cron scripts still use LiteLLM MCP. 6 Retire LiteLLM-as-hub: PR open 2026-10-31 (cron scripts moved to the gateway, LiteLLM `mcp_servers` and the Hermes egress/ingress rules removed; LiteLLM stays the model gateway).** | design [#591](https://github.com/onemindos/omos/pull/591), 0 and 1 [#593](https://github.com/onemindos/omos/pull/593), 2 [#594](https://github.com/onemindos/omos/pull/594), 3e MediaMTX [#595](https://github.com/onemindos/omos/pull/595), 3a [#596](https://github.com/onemindos/omos/pull/596), 3d [#597](https://github.com/onemindos/omos/pull/597) and Tile38 [#598](https://github.com/onemindos/omos/pull/598), 4 [#592](https://github.com/onemindos/omos/pull/592) | 6, 7 |

### Application-lane catalog preparation — 2026-10-27

Zeus's latest instruction assigns Codex all four remaining original PRs:
#505, #535, #533 and #532. Catalog
[PR #505](https://github.com/onemindos/omos/pull/505) was integrated with main `3495a614`; its portal image-build
catalog/topology gate passes all 66 tests on supported Node 24. It repairs the
stale generated mapping/runbook and covers all 60 source Argo Applications,
with a 409-entity static inventory baseline. The
[Claude continuation](docs/operations/2026-10-27-claude-merge-continuation.md)
records 42 of the original 46 PRs merged and the specific technical gates for
the final four. Source preparation does not publish an image
or certify the catalog ConfigMap's later portal rollout. Phase 6 remains In
progress; Phases 7/8 remain Not started.

### Fresh application replacement source (2026-10-24)

[PR #390](https://github.com/onemindos/omos/pull/390) is prepared against the current cleanup baseline, preserving
fresh STAC source and the shared CNPG roles/databases. CloudTAK is a separate TAK client with pinned upstream
images, a held database/role prerequisite on `onemind-data`, and private Tailscale ingress. The Application is manual
with zero retries while shared PostgreSQL recovery and isolated migration validation are pending. Source validation
is not live acceptance. Spaces credentials, first-run TAK configuration and runtime checks remain deployment gates;
imported/hosted tile layers need a trusted internal tile origin, and media-infra replacement needs a reviewed
reader/publisher lease migration. See [CloudTAK's runbook](platform/cloudtak/README.md).

### Migration readiness: core gate versus parallel follow-ups

Zeus requested [a tracked readiness/follow-up checklist](docs/roadmap/migration-readiness-followups.md)
on 2026-10-20, tracked in [PR #331](https://github.com/onemindos/omos/pull/331).
Release/live verification of #330 and the chosen app's security,
dependencies, data/recovery and rollback gates come first. Unrelated launcher
repairs, optional dashboards and the separately planned workflow plane are not
blanket prerequisites for every migration. State/recovery and security failures
still block affected applications; risks require explicit disposition, not a waiver.
Phase 6 remains In progress and Phase 7 Not started until actual acceptance.

Full-stack continuation requested 2026-10-20: **#331 is the resume index** linking
[#333](https://github.com/onemindos/omos/pull/333) (20-category Kubernetes structure,
installed/proposed apps and execution queue),
[#334](https://github.com/onemindos/omos/pull/334) (public ingress lifecycle plan),
[#332](https://github.com/onemindos/omos/pull/332) (staged private Legacy/LibreChat)
and existing [#284](https://github.com/onemindos/omos/pull/284) (home expansion,
to revise against Mac evidence). Planning/read-only preparation may run in
parallel with Phase 6. This does not authorize new controllers, app cutovers,
identity/DNS/storage/provider/billing changes or wider AI permissions. Start with
the core gate and foundation assessment, not a simultaneous platform installation.

## Phases

### Phase 0: Blueprint

Goal: one shared map for humans and agents. No cluster changes.

- [x] Write this roadmap at the repository root.
- [x] Map every namespace to a layer.
- [x] Point `CLAUDE.md`, `AGENTS.md` and `docs/REPO-LAYOUT.md` to it.
- [x] Merge.

Gate: the PR is merged and the tracker shows Phase 0 `Done`.

Evidence: PR #256 merged 2026-10-18 with all CI checks passing (lint, kubeconform, Helm render, gitleaks).

### Phase 1: Claude sees every layer

Goal: Claude Code can answer "what is broken, where and why" across all four layers.
Runs on the Mac (or Docker Desktop), so it adds nothing to cluster memory. All access is
read-only.

- [x] Backstage MCP with the existing machine token (`~/.config/onemind/backstage-mcp.token`).
      Verified 2026-10-18: `tools/list` returns 6 read-only tools.
- [x] Kubernetes MCP server ([containers/kubernetes-mcp-server](https://github.com/containers/kubernetes-mcp-server))
      0.0.67 with `--read-only`, using `kube/ai-operator`: built-in `view` plus
      `ai-operator-read` for platform CRDs, no Secrets
      ([`platform/kube/ai-operator.yaml`](platform/kube/ai-operator.yaml)).
- [x] Setup script for Claude Code and Codex: [`scripts/ai-operator/setup.sh`](scripts/ai-operator/setup.sh).
- [x] Document the setup in [`docs/operations/ai-operator-access.md`](docs/operations/ai-operator-access.md).
- [x] After merge and Argo sync: run the script and pass the gate below.
- [x] Follow-up: add `nodes` to `ai-operator-read` (the built-in `view` role does not include them).
- Moved: OpenChoreo control-plane MCP goes to phase 2. The endpoint works
  (`/mcp` answers 401 and advertises Authentik), but OpenChoreo manages nothing yet, and
  it needs a read-only OpenChoreo role and an Authentik client for the agent.
- Dropped: a separate Argo CD MCP. Argo CD is still a Helm CLI release (foundation), and
  Application sync and health are readable through `onemind-kubernetes`.
- Deferred: the tool gateway (LiteLLM MCP or agentgateway) is only needed for in-cluster
  agents, so it moves to phase 8 with Hermes.

Gate: from a fresh session, Claude lists unhealthy workloads in each layer with the cause,
using only the MCP servers, and every server refuses a write.

Evidence (2026-10-18, after #257 synced):

- `claude mcp list`: `onemind-kubernetes` and `onemind-backstage` connected. Codex config
  has the same two servers.
- `onemind-kubernetes` exposes 14 tools, all read-only. Calling `resources_delete` fails
  with "unknown tool"; reading `secret/kube/ai-operator-token` is forbidden for
  `system:serviceaccount:kube:ai-operator`.
- `onemind-backstage` `onemind.health` returned live checks.
- Read-only queries found real problems, with causes from pod events:
  `nats/nats-2` Pending (bus at 2/3; hub-core nodes lack CPU and memory, autoscaler in
  backoff), `takserver/tak-ldap-0` and `takml/takml-server` Pending (pinned by hostname to
  `onemind-hub-core-v2-37phkg`, a node DOKS replaced about 10:30 UTC), plus
  `falco` and `spire-agent` DaemonSet pods without room on one node.
- Listing nodes was denied; fixed in this PR.

### Phase 2: OpenChoreo foundations

Goal: OpenChoreo can deploy something. All resources in Git, applied by Argo CD.

Follows the official guide
[Run in your environment (v1.3.x)](https://openchoreo.dev/docs/getting-started/try-it-out/on-your-environment/)
step by step, as PRs 2a to 2e. Where this cluster differs from the guide:

| Guide step | Guide | Here | Reason |
|---|---|---|---|
| 1. Gateway API CRDs | v1.5.1 | keep v1.2.1 | Installed and owned by DOKS-managed Cilium; the control plane already runs on them |
| 1. cert-manager | v1.19.4 | keep v1.20.3 | Already installed, newer |
| 1. External Secrets | 2.0.1, `installCRDs=true` | same, as `ecosystem-external-secrets` | Only stale v1beta1 CRDs existed, no controller and no objects |
| 1. kgateway | install | skip; gateways use the Cilium class `onemind-private` (tailnet) | The data-plane chart supports `gateway.gatewayClassName`; the control plane already uses it |
| 1. OpenBao + `ClusterSecretStore` `default` | dev-mode OpenBao | `default` store on ESO's Kubernetes provider, reading the `openchoreo-secrets` namespace, filled by SealedSecrets | The secret-management guide allows any ESO backend named `default`; the guide marks dev-mode OpenBao as not for production |
| 2. `openchoreo-ca` ClusterIssuer | self-signed CA | not created yet | Gateways start with TLS off (the guide's first install also does). Later TLS needs a new ClusterIssuer: `letsencrypt-dns` was removed as unused (its Cloudflare token was invalid) |

- [ ] New tree `platform/openchoreo/` with its Argo Application(s):
      `planes/` (charts and registration), `platform/` (environments, pipeline,
      component types, traits), `projects/` (one folder per project).
- [x] 2a. Secret store: External Secrets 2.0.1 (`ecosystem-external-secrets`) and the
      `ClusterSecretStore` `default` on the Kubernetes provider
      (`platform/openchoreo/prereqs/`, Argo Application `platform-openchoreo-prereqs`).
      SealedSecrets stays the way secrets enter Git. Verified 2026-10-18: ESO pods
      running, `clustersecretstore/default` `Valid`, `Ready=True`.
- [x] 2b. Default resources (guide step 4): label `default` as a control-plane namespace
      and apply the guide's `samples/getting-started/all.yaml`, pinned to v1.3.1, straight
      from the upstream repository (Argo Application `platform-openchoreo-defaults`).
      Creates the `default` project, environments `development`, `staging`,
      `production`, the `default` pipeline, component types (`service`,
      `web-application`, `worker`, `scheduled-task`), resource types (`postgres`,
      `valkey`, `nats`), build workflows and traits. Verified 2026-10-18: Argo
      Application Synced and Healthy; project, three environments, pipeline and 16
      cluster-scoped types present.
- [x] 2c. Data plane (guide step 5, install): chart `openchoreo-data-plane` 1.3.1 with
      `gateway.tls.enabled=false` and the `onemind-private` gateway class, plus the
      `cluster-gateway-ca` ConfigMap copy (Argo Application `platform-openchoreo-data-plane`,
      `platform/openchoreo/data-plane/`). Agent certificate lifetime raised to 10 years,
      because the registration pins its self-signed CA. Verified 2026-10-18: Synced and
      Healthy; agent running; `gateway-default` programmed at tailnet IP 100.110.185.30
      (`openchoreo-data-plane-cilium-gateway-gateway-default.tail717ba0.ts.net`).
- [x] 2d. Register `ClusterDataPlane` `default` (guide step 5, register):
      `platform/openchoreo/data-plane/cluster-data-plane.yaml` with the agent CA,
      `secretStoreRef: default`, and the HTTP listener on the guide's nip.io domain built
      from the gateway's tailnet IP: `apps.openchoreo.100-110-185-30.nip.io` (resolves
      only to a tailnet address, so apps are tailnet-only). A real domain such as
      `*.apps.onemindos.dev` with TLS from a new DNS-01 ClusterIssuer replaces it later. Verified
      2026-10-18: agent connected (tunnel requests return 200), `ClusterDataPlane`
      `Created=True`, the three environments `EnvironmentReady`, and the `default`
      project's bindings `Ready` with namespaces `dp-default-default-{development,staging,production}-*`.
- [x] 2e. First `ClusterComponentType` `onemind-service`
      (`platform/openchoreo/platform/onemind-service.yaml`, Argo Application
      `platform-openchoreo-platform`): the upstream `service` type without its CORS
      filter (unsupported by Gateway API v1.2.1 here), pods prefer the ecosystem pool and
      cannot be pinned to a node, a `strategy` parameter (`Recreate` for ReadWriteOnce
      volumes), default requests 50m / 128Mi. Apps are reached through the tailnet
      gateway, so no per-app Tailscale ingress. Network policy: later, as a project type.
      Verified 2026-10-18: `platform-openchoreo-platform` Synced and Healthy,
      `clustercomponenttype/onemind-service` present.
- [x] 2g. OpenChoreo MCP for Claude (moved from phase 1), in three PRs:
  - [x] 2g-1. Make the control-plane API reachable on the tailnet under its HTTPRoute
        hostname: the gateway Service gets `tailscale.com/hostname: onemind-openchoreo-api`
        and the API advertises `publicUrl` instead of `api.openchoreo.invalid`
        (`platform/ecosystem/values/oc-control.yaml`). Verified 2026-10-18: the device
        `onemind-openchoreo-api` resolves on the tailnet. The rollout then stuck: the
        ecosystem node was at 99% CPU requested and the control plane required that node
        (`nodeSelector`), so the new API pod stayed Pending. Placement changed to a
        preference (#269); the capacity work stays in phase 5.
  - [x] 2g-2. Read-only identity: an Authentik service account `ai-operator` allowed on
        the `onemind-openchoreo` application (hook Job in
        `platform/ecosystem/resources/ai-operator-identity.yaml`; its app password lands
        in Secret `auth/ai-operator-openchoreo`, not in Git), a custom read-only
        `ClusterAuthzRole` `ai-operator-read` (every `:view` action except secrets), and a
        binding for the account. Two steps: the Job prints the account's token claims, then
        the binding is added with the account's `sub`. Verified 2026-10-18: the Job created
        the account and Secret and minted a token (aud `onemind-openchoreo`, lifetime 3600 s,
        no `client_id` claim, so the binding uses `sub`). Finding: the `onemind-openchoreo`
        provider had no grant types at all (`[]`), which Authentik 2026.8 treats as none
        allowed, so its browser login could not work either; the Job now adds
        `client_credentials` (#272). Console login (2f) must also set `authorization_code`
        and `refresh_token`.
  - [x] 2g-3. Mac side: a Claude Code `headersHelper`
        (`scripts/ai-operator/openchoreo-token.sh`) that mints a fresh token per connect,
        registered by `scripts/ai-operator/setup.sh` (also for Codex via
        `OPENCHOREO_MCP_TOKEN`); docs in `docs/operations/ai-operator-access.md`. Verified
        2026-10-18: `claude mcp list` shows `onemind-openchoreo` connected; the account sees
        63 read tools, `list_namespaces` and `list_projects` return data, and
        `create_namespace` is refused (`missing permission "namespace:create"`).
- [x] 2f. Turn on the built-in console (`backstage.enabled: true`), tailnet-only at
      `https://onemind-console.tail717ba0.ts.net` (Tailscale ingress), signing in through
      Authentik, in two PRs:
  - [x] 2f-1. Identities and secrets: the identity Job also creates the service account
        `backstage-catalog` (the console's catalog sync), writes
        `openchoreo-control-plane/backstage-secrets` (generated, never in Git), and gives the
        `onemind-openchoreo` provider the grant types browser login needs
        (`authorization_code`, `refresh_token`); the console URL is added to the provider's
        redirect URIs. Verified 2026-10-18: provider grant types now
        `client_credentials`, `authorization_code`, `refresh_token`; Secret
        `backstage-secrets` created; `backstage-catalog` mints tokens both with user+password
        and with the console's `client_id` + `client_secret` request style (same `sub`).
  - [x] 2f-2. The console itself: chart values, Tailscale ingress, and authorization
        bindings (`backstage-catalog` to `backstage-catalog-reader`, `akadmin` to `admin`).
  - Decision: our 555-entity catalog is not mounted into the temporary console. A
    ConfigMap cannot be shared across namespaces, and phase 6 merges the OpenChoreo plugins
    into our portal, which already holds that catalog. The console shows what OpenChoreo
    manages; the portal remains the full inventory.
  - Verified 2026-10-18: console pod Running on hub-core; `https://onemind-console.tail717ba0.ts.net`
    answers 200 over the tailnet; its sign-in redirects to Authentik's authorization endpoint
    and on to the login flow (no redirect-URI or grant-type error); the catalog sync, using
    the `backstage-catalog` account, processed 34 entities (1 domain, 1 system, 3
    environments, 1 pipeline, 5 cluster component types, 3 resource types, 1 project type,
    5 traits, 1 data plane, 4 workflows). Not done by an agent: typing the `akadmin`
    password. Zeus signs in once to close this item; if it fails, phase 2 reopens.

Gate: the data plane agent is connected, the plane shows healthy in the OpenChoreo API,
and the built-in console signs in through Authentik.

Evidence (2026-10-18): the agent is connected and the control plane routes requests through
it; `ClusterDataPlane` `default` is `Created=True`; the three environments and the `default`
project's bindings are ready; all seven phase-2 Argo Applications are Synced and Healthy;
Claude Code reads OpenChoreo through the MCP server and is refused writes; the console
serves, signs in through Authentik up to the password page, and syncs the catalog. Open: the
interactive `akadmin` sign-in (above).

Findings to act on in later phases:

- The chart's bootstrap bindings use the `groups` claim (`admins`, `developers`,
  `platform-engineers`, `sres`) and `client_id` values (`service_mcp_client`,
  `openchoreo-*-client`). Here only `sub` (users) and `client_id` (service accounts) are
  configured entitlements, and Authentik tokens carry no `client_id`, so none of those match.
  They are inert, including `mcp-tryout-client-binding` (`admin` for `service_mcp_client`).
  Remove them with a `bootstrap` override when roles are reviewed (phase 6).
- The Authentik provider made by `auth-provider-bootstrap.yaml` for `langfuse` (and `registry` and `workflows`,
  deleted with Harbor and Argo Workflows by `authentik-registry-removal.yaml`) was created with an empty grant type list, like `openchoreo` was, and
  cannot complete a browser login until given `authorization_code`.
- The ecosystem node is at about 99% CPU requested; the control plane now prefers it
  without requiring it. Capacity work is phase 5.

### Phase 3: Pilot one app

Goal: prove OpenChoreo saves work before moving anything important.

Progress (2026-10-18):

- Read-only baseline at 21:36 EDT: 64 Argo Applications, 62 Synced/Healthy;
  `platform-matrix` OutOfSync/Healthy and `platform-openclaw` Synced/Progressing
  are existing exceptions. No Pending pods. All phase-2 OpenChoreo Applications
  are Synced/Healthy. The API lists the ready `default` project and no Components.
- Add `platform-openchoreo-projects` and the `onetak` Project with development,
  staging and production bindings, following the v1.3.1 getting-started sample.
  This step only provisions project namespaces; the old STAC workloads stay in place.
  PR #286 merged with all four CI checks passing. Verified at 21:45 EDT through
  read-only MCP: `platform-openchoreo-projects` Synced/Healthy, `onetak` Ready in
  the OpenChoreo API, all three `onetak-*` ProjectReleaseBindings Ready. The
  app-of-apps picked up the commit about six minutes after merge.
- The console's interactive `akadmin` login remains awaiting Zeus's confirmation.
- Candidate inspection: the old STAC Browser pod is Ready with image digest
  `sha256:e6d92b6ba593201b558219d0591fc83820b7f86a316d3fe565a27ad996fc7f63`.
  Its documented tailnet hostname does not resolve from the Mac; migration must
  verify the real browser entry point, runtime configuration and catalog access.
- Pilot deviations from the brief: the image expects `SB_catalogUrl`, `SB_catalogTitle`
  and `SB_pathPrefix`, not the old `CATALOG_URL` / `STAC_BROWSER_TITLE`. A browser
  cannot use `stac-api.onetak.svc.cluster.local`, so the app's nginx proxies catalog
  reads at the same tailnet origin (GET/HEAD, plus POST only to `/search`). Responses
  rewrite the upstream base URL; catalog write methods are refused. The HTTPRoute
  preserves `/stac-browser-http` via the endpoint's `basePath`. These are application
  settings; the shared component type, TAK Server and NATS are untouched.
- PR #287 merged with all four checks passing. At 21:48 EDT, Argo was
  Synced/Healthy and OpenChoreo had rendered the Deployment, Service, HTTPRoute,
  two ConfigMaps and an ingress NetworkPolicy. The runtime pod failed to start:
  Workload file `mountPath` is a directory; the filename is appended by OpenChoreo.
  Correct it to `/etc/nginx/conf.d` in a new immutable `stac-browser-v2` release
  ([#288](https://github.com/onemindos/omos/pull/288)). The old deployment remains Ready; no cutover has happened.
- At 21:55 EDT, v2's pod was Running and its ReleaseBinding had `Ready=True`.
  HTTPRoute was Accepted with ResolvedRefs, but the tailnet request returned 503
  after five seconds: upstream connection timeout. The generated NetworkPolicy
  allows same-namespace pods and `openchoreo.dev/system-component` pods; Cilium
  Gateway API uses the reserved `ingress` identity instead. Add the opt-in
  `cilium-gateway-ingress` ClusterTrait ([#289](https://github.com/onemindos/omos/pull/289)), scoped to the component's full
  pod selectors and TCP 8080. It renders a CiliumNetworkPolicy in the runtime
  namespace; the upstream policy remains. See Cilium's official [Gateway API
  policy reference](https://docs.cilium.io/en/stable/network/servicemesh/gateway-api/gateway-api/#cilium-s-ingress-config-and-ciliumnetworkpolicy).
  Publish immutable release v3 and advance development; verify traffic after sync.
- Concurrent live-state change at 21:55 EDT: the old ecosystem node was cordoned
  and two new ecosystem nodes were Ready on Kubernetes 1.36.5. Many unrelated pods
  were rescheduling; the OpenChoreo API and metrics had transient reachability
  failures. No node, pool or infrastructure changes are part of the pilot PRs.
- At 22:01 EDT, the Cilium policy was `Valid=True`; the tailnet frontend,
  runtime config and all referenced JS/CSS assets returned 200 with the correct
  path prefix. The catalog root returned 200, but its self/collection links used
  `http://stac-api.onetak.svc.cluster.local/` without port 8082. Add that observed
  exact base URL to the proxy rewrites in immutable release v4 ([#290](https://github.com/onemindos/omos/pull/290)). Keep
  external asset links unchanged; collections/search and visual review are pending.
- Catalog `GET /collections` and `POST /search` currently return 500 from the
  existing STAC API. Its logs show `UndefinedFunctionError` for
  `collection_search(jsonb)` and `search(jsonb)`. Add a one-off, Argo-owned
  PostgreSQL client Job ([#291](https://github.com/onemindos/omos/pull/291)) to inspect only schema/function metadata and
  search path in `BEGIN READ ONLY`, with `default_transaction_read_only=on`.
  It references the existing Secret without reading or changing its value.
  Record the result and remove the Job after inspection; do not migrate the
  shared database or remove the old browser based on pod readiness alone.
- The inspection Job completed with `default_transaction_read_only=on` and
  `BEGIN READ ONLY` against database `stac` as role `geo`, PostgreSQL 17.11.
  It found no `pgstac` schema, no PgSTAC tables, and none of `search`,
  `collection_search` or `get_version` in `pgstac`/`public`. Installed extensions
  were `plpgsql` 1.0, `timescaledb` 2.29.1 and `timescaledb_toolkit` 1.24.0;
  PostGIS is not enabled in this database. This does not establish that all user
  data is absent or why the schema is missing. Remove the completed diagnostic
  Job manifest after preserving these results; no database changes are included.
- Follow-up ([#292](https://github.com/onemindos/omos/pull/292)) at 22:58 EDT: the v4 binding still reports
  `ComponentReleaseNotFound`, although its ComponentRelease exists. The official
  [v1.3.1 controller](https://github.com/openchoreo/openchoreo/blob/v1.3.1/internal/controller/releasebinding/controller.go)
  returns without a retry on that condition and does not watch ComponentRelease
  creation. Order Project/Component (wave 0), snapshots (wave 1), then binding
  (wave 2), and advance to a new v5 snapshot to trigger reconciliation. This is a
  GitOps ordering workaround, not evidence of a missing config rollout checksum:
  the existing Deployment already has OpenChoreo's `dp-resource-hash` annotation.
- The shared and frozen gateway trait used an invalid `creates[].id`; the API
  prunes it and defaults `targetPlane` to `dataplane`, causing persistent Argo
  drift. Normalize that field in the shared trait and v3/v4 snapshots to the
  already stored representation. Existing immutable runtime specs remain the same
  after API defaulting. v5 carries v4's workload/proxy configuration unchanged.
  See the official [v1.3.1 ClusterTrait schema](https://github.com/openchoreo/openchoreo/blob/v1.3.1/config/crd/bases/openchoreo.dev_clustertraits.yaml).
- Backup observation: `onetak/volume-50fb74ff16c6` reports a successful VolSync
  snapshot of `onetak-geo-db-data` at `2026-10-18T05:20:18Z`. That is volume-backup
  evidence, not proof of a fresh, application-correct STAC database restore.
  Before any live schema initialization, confirm the intended database and its
  contents, create a fresh PostgreSQL-17-compatible backup, prove an isolated
  restore and test a compatible PgSTAC installation. Live initialization remains
  subject to Zeus's approval; the shared database stays outside OpenChoreo.
- PR #292 merged with all four CI checks passing. Verified live at 23:06 EDT:
  `platform-openchoreo-platform` and `platform-openchoreo-projects` are
  Synced/Healthy at merge `09c0f051`; binding v5 is `Ready=True` with observed
  generation 5. Its healthy resource tree includes the Deployment, Pod, route,
  Service, ConfigMaps and both policies. OpenChoreo's existing hash changed to
  `6469b5cd9` and rolled a new ready Pod, confirming no extra checksum is needed.
  The HTML, runtime config and all six linked assets return 200; the catalog root
  returns 200 with its self link rewritten to the tailnet proxy. OPTIONS and PUT
  on the catalog return 405. Collection listing and limited POST search still
  return 500 from the existing API. The original browser remains available, and
  Argo pruned the completed inspection Job. Visual review remains unverified.
- Next read-only preflight ([#293](https://github.com/onemindos/omos/pull/293)): inspect the existing `geo` and `stac` databases for
  owners/sizes, role flags, available/installed extensions, user schema/relation
  counts, table names and PgSTAC function signatures. This checks for an existing
  installation in another database before selecting any repair. The Argo-owned
  Job uses the same pinned PostgreSQL client, the existing Secret reference,
  read-only session defaults and transactions, bounded queries and no retries.
  The recovery and approval gates are documented in
  [`docs/operations/stac-database-repair.md`](docs/operations/stac-database-repair.md).
  No live database initialization or replacement is part of this diagnostic step.
- PR #293 merged with all four checks passing. At 23:13 EDT,
  `platform-onetak` was Synced/Healthy at `4adc570a` and the preflight Job completed
  successfully with both session and transaction read-only flags `on` in each
  database. `geo` (18 MB) has PostGIS 3.6.4 and public table `spatial_ref_sys`;
  `stac` (10 MB) has no public relations and no enabled PostGIS. Neither database
  has a PgSTAC schema/search functions. PostGIS 3.6.4, btree_gist 1.7, pg_trgm 1.6
  and unaccent 1.1 are available on the server. The existing `geo` role is a
  superuser with database/role creation privileges; the three standard PgSTAC
  roles are absent. These are metadata observations, not application-row or
  historical-data verification. Preserve the detailed evidence in the recovery
  document and remove the completed Job in [#294](https://github.com/onemindos/omos/pull/294).
  No schema, grants, credentials, server
  image or application database target was changed.
- Evaluation and direction: the pre-cutover active Component/Workload, binding
  and v5 snapshot total 492 lines in three files. With five immutable release
  snapshots, the app directory contained seven files and 1,958 lines; shared
  platform/bootstrap files are additional. The old browser occupied 47 lines in
  one shared manifest. Standardized rollout and resource-tree visibility are
  demonstrated benefits, but they did not establish lower maintenance cost.
- On 2026-10-19, Zeus chose option A: keep the pilot and reduce manual release
  YAML before expanding. [#279](https://github.com/onemindos/omos/pull/279) is
  the cutover PR removing the old immutable release snapshots and the legacy
  stac-browser Deployment/Service. No live database initialization was performed.
- On 2026-10-19, Zeus confirmed the OpenChoreo console `akadmin` login works; a
  screenshot shows the Home page with the `default` data plane `Connected`.
- PR #279 merged at `4ecc1b36`. Argo `platform-onetak` auto-synced to that
  revision at 04:02:02 UTC (about six minutes after merge) and pruned the legacy
  `stac-browser` Deployment and Service from namespace `onetak`.
- Live verification at 04:04 UTC: `platform-onetak` and
  `platform-openchoreo-projects` are `Synced/Healthy`; the legacy `stac-browser`
  Deployment/Service no longer exist in `onetak`; no Pending pods cluster-wide;
  the OpenChoreo-managed pod `stac-browser-development-0c6b80f7-756cfdb5ff-r4pr2`
  is `1/1 Running` in namespace `dp-default-onetak-development-b188bf98`; the
  generated HTTPRoute `stac-browser-http-02eaf4a9` reports `Accepted=True` and
  `ResolvedRefs=True`; the nip.io pilot URL returns HTTP 200 and valid HTML.
- No live database initialization was performed. The STAC database repair gates
  in `docs/operations/stac-database-repair.md` remain unapproved, so the catalog
  listing still returns 500 from the existing API; this is a pre-existing
  application-layer issue, not a blocker for the OpenChoreo pilot gate.
- Phase 3 gate is satisfied: one stateless Layer 3 app is managed by OpenChoreo,
  serves traffic through the data-plane gateway, the old manifests are removed,
  and the verdict is recorded.

- [x] Pick a small stateless Layer 3 app (candidates: `stac-browser`, `agent-console`).
- [x] Write its Component in `platform/openchoreo/projects/`.
      The parallel deployment PR adds a pinned STAC Browser Component, Workload,
      immutable ComponentRelease and development ReleaseBinding. Runtime verification
      and cutover are pending; the existing app has not been removed.
- [x] Run it next to the existing deployment, switch traffic, then delete the old manifests.

Gate: the app works, it needs fewer files than before, and the change was easier to make.
If not, stop and record why before Phase 7.

### Phase 4: Portal refactor

Goal: our portal features become portable Backstage plugin packages. No behaviour change.

Progress (2026-10-19): rebased `roadmap/phase-4-portal-refactor` onto `main`; frontend and backend modules extracted into `@internal/plugin-onemind-*` packages, all wired via `packages/backend/src/index.ts` and `packages/app/src/App.tsx`. Phase 4 PR #280 merged; hotfix PR #298 shipped `directory.json` in the production bundle; live gate checks passed.

- [x] Move each frontend feature from `services/onemind-portal/packages/app/src/modules/`
      into its own package under `services/onemind-portal/plugins/`.
- [x] Move each backend feature from `packages/backend/src/` into its own backend plugin
      or module package.
- [x] Tests, TypeScript check, image build and a browser check pass; new image digest pinned.

Evidence:
- `yarn tsc`, `yarn lint:all`, `yarn test:all`, `yarn build:all` pass locally (Node 22.23.3, Yarn 4.13.0).
- CI `Build OneMind Portal` on branch `roadmap/phase-4-portal-refactor` succeeded (run 39205172177); image pushed to `ghcr.io/onemindos/onemind-portal@sha256:3d2c5bdb6cf56c45b051b373c5085d79616af60055420332a452e5b80cc469c3`.
- `platform/backstage/deployment.yaml` updated to the new digest; Argo CD `backstage` Application is OutOfSync and syncing (live read-only verification in progress).
- `scripts/validate_changed_yaml.py` against `origin/main` passed (0 changed YAML files in this branch).
- Full `make validate` still reports pre-existing YAML lint debt in `platform/matrix/*` (unchanged from `origin/main`), so the root Makefile fails at `lint-yaml`; the changed-files gate is clean.

PR: https://github.com/onemindos/omos/pull/280

Gate (2026-10-19): the deployed portal behaves exactly as before.
- CI `Build OneMind Portal` on `main` succeeded after hotfix PR #298 (run 39206336020, SHA `133dfd72`); artifact digest `sha256:8110199af6a56dfd8538cfe2145eba319201128b9dfca555b0c93caa1f6bc6b2`.
- `platform/backstage/deployment.yaml` pinned to that digest in commit `28900eb1`.
- Argo CD `platform-backstage` Application synced revision `28900eb1` at 2026-10-19T05:54:44Z; status `Synced`, health `Healthy`.
- Kubernetes deployment `backstage` rolled out: 2 desired, 2 updated, 2 total, 2 available, 0 unavailable.
- Live gate checks against `https://onemind-preview.tail717ba0.ts.net`:
  - `/.backstage/health/v1/readiness` → HTTP 200 in ~58 ms
  - `/.backstage/health/v1/liveness` → HTTP 200 in ~56 ms
  - `/` → HTTP 200 in ~61 ms (HTML landing page returned)


### Phase 5: Capacity, then observability

**Done — live gates verified 2026-10-19, 08:05–08:07 UTC.** Final evidence/brief:
[#281](https://github.com/onemindos/omos/pull/281).

- [x] Apply measured request/placement fixes; no extra node required.
- [x] Install observability plane 1.3.1 with modules on the existing Prometheus,
      OpenSearch and OTel collector; reuse OpenCost through its FinOps adapter.
- [x] Register ClusterObservabilityPlane/default with its own agent CA and explicitly
      link ClusterDataPlane/default. Agent connected: one authenticated connection.
- [x] Connect Claude Code and Codex to the authenticated Observer MCP for logs,
      metrics, traces, alerts, incidents and costs. Claude reports Connected.

#### Capacity gate (live requested resources)

Effective pod requests include init containers/native sidecars and overhead; completed
Jobs are excluded. No Pending pods at verification. Actual node memory use was
44–65%, distinct from the scheduler request percentages below.

| Node | Memory requested | CPU requested |
|---|---:|---:|
| onemind-ecosystem-cpu-3rhwrp | 21.18 GiB / 75.4% | 6.077 cores / 77.1% |
| onemind-hub-core-v2-3rtf47 | 9.82 GiB / 73.7% | 3.832 cores / 48.6% |
| onemind-hub-core-v2-3rtf4m | 9.38 GiB / 70.4% | 3.917 cores / 49.7% |
| onemind-hub-core-v2-3rtf4q | 10.36 GiB / 77.7% | 5.862 cores / 74.4% |

Hermes, Loki, Prometheus, Grafana and OpenSearch reconciled onto the existing ecosystem
pool through preferred affinity, preserving fallback scheduling. No node-name pins,
new node, descheduler or destructive volume operation was needed. OpenCost's two KSM
emission flags are false. Unavailable DOKS scheduler/controller-manager and Cilium
kube-proxy Down rules are absent; KubePodCrashLooping remains loaded.

#### Observability and human/AI gate

- Core, logs, traces, events and FinOps Argo Applications are Synced/Healthy with
  successful operations; both OpenSearch setup Jobs completed. Core uses chart 1.3.1;
  logs/traces 0.6.0, metrics 0.7.1, events 0.1.2 and FinOps 0.2.0.
- Existing OTel DaemonSet is 4/4 ready. Original Loki/Tempo pipelines remain intact;
  only OpenChoreo-tagged component signals are copied to the existing OpenSearch.
  There is no second Prometheus, OpenSearch or application log/trace collector.
- Pilot: namespace default, project onetak, component stac-browser, environment
  development. Five ordinary HTTP requests returned 200 and produced five real logs.
  The signed-in console visibly rendered all five logs and CPU/memory graphs.
- Authenticated MCP returned the same five logs and eight one-minute points in each
  of six CPU/memory usage/request/limit series. Memory usage was about 7.6 MiB.
  Cost query returned one component allocation. Traces/alerts/incidents queries
  succeeded with empty results: the pilot emits no application spans, and no data
  or incidents were fabricated. Seventeen tools discovered; anonymous calls get 401.
- Observer URL is https://onemind-observer.tail717ba0.ts.net; MCP path /mcp. Console
  clients fetch directly, so registration must not use cluster-only DNS. Issuer,
  audience, read-only authorization and tailnet isolation remain enabled. CORS
  permits the console and portal origins only. Internal webhook port 8081 is not exposed.
- Operator setup and actual Claude MCP connection checks passed for Kubernetes,
  portal, control-plane and Observer servers. Codex uses the existing short-lived
  OPENCHOREO_MCP_TOKEN environment variable; restart clients to load registrations.

#### Recovery gate

Observer SQLite is persistent (live PVC 1 GiB), with a distinct encrypted Restic
repository. Backup completed at 07:52:31 UTC, snapshot 583bae7f (three files).
An Argo-owned isolated restore to a NEW protected PVC completed at 07:59:46 UTC.
The non-root, read-only verifier reported SQLite integrity_check=ok and the
alert_entries/incident_entries tables. Production storage was never overwritten.
The manual source trigger was removed after verification; daily schedule is 04:50
UTC, next run 2026-10-20T04:50:00Z. The restore PVC remains prune/delete protected;
cleanup requires explicit storage-deletion approval.

#### PRs, rollout corrections and validation

- Capacity/resource hygiene: [#300](https://github.com/onemindos/omos/pull/300),
  [#301](https://github.com/onemindos/omos/pull/301), [#302](https://github.com/onemindos/omos/pull/302),
  [#303](https://github.com/onemindos/omos/pull/303), [#304](https://github.com/onemindos/omos/pull/304),
  [#305](https://github.com/onemindos/omos/pull/305).
- Plane/modules and activation: [#306](https://github.com/onemindos/omos/pull/306),
  [#307](https://github.com/onemindos/omos/pull/307).
- Live verification fixes: [#308](https://github.com/onemindos/omos/pull/308) removes
  upstream setup Base64 line wrapping without changing credentials and restricts
  event labels to avoid mapping collisions; [#309](https://github.com/onemindos/omos/pull/309)
  adds isolated recovery and uses Authentik's public token URL so the minted issuer
  matches trust configuration; [#310](https://github.com/onemindos/omos/pull/310)
  registers the browser-reachable private HTTPS Observer URL. URL caches expire in five minutes.
- Local pinned Helm/Kustomize renders, changed YAML/schema gates, live CRD dry runs,
  OTel 0.156.0 validation/fan-out test, MCP JSON/SSE/auth/config-idempotence checks,
  ShellCheck and secret scans passed. Full make validate still fails on pre-existing
  onetak/TAK/Matrix YAML lint debt. Later GitHub CI jobs did not start because of an
  account payment/spending-limit error; no billing settings or merge rules were changed.
- Final cluster snapshot: 70 Argo Applications, 68 Synced/Healthy. The only exceptions
  are pre-existing platform-matrix (OutOfSync/Healthy) and platform-openclaw
  (Synced/Progressing). Existing portal remains 2/2 available, readiness HTTP 200.
  OpenSearch's existing self-signed/demo backend TLS is unchanged; hardening is separate.

Next unfinished phase: **Phase 7**. This does not authorize starting it automatically.

### Phase 6: One portal

Existing pilot integration: [PR #330](https://github.com/onemindos/omos/pull/330), [STAC operational slice](docs/operations/stac-operational-slice-2026-10-20.md) joins identity/dependencies and real evidence; source checks pass,
immutable image/live UI acceptance remains required. STAC collections/search
currently return 500 (no pgstac schema observed); data freshness is unknown.
No data/schema repair or new application migration is performed by this integration.

2026-10-20 follow-through: [PR #329](https://github.com/onemindos/omos/pull/329), [trust/coverage review and launcher repairs](docs/operations/portal-acceptance-followthrough-2026-10-20.md). Refreshed login observed on both replicas; four confirmed launcher fixes;
16 remaining launcher findings are flagged, offsite/Grafana gaps remain explicit.
The existing STAC operational slice is authorized next; this is not a Phase 7 cutover.

**Stabilization reopened 2026-10-20 by Zeus.** Historical release acceptance below
is preserved, but the broader UI/catalog quality audit found integration defects.
Work the [repair and acceptance gates](docs/operations/portal-stabilization-2026-10-20.md)
before Phase 7. Repairs #322–#326 are merged and deployed: 2/2 healthy,
715 processed catalog entities, final browser/MCP gates passed. Remaining broader
acceptance gaps are explicitly recorded; migration remains paused, not implicitly
authorized by the portal repair.

**Done — live acceptance 2026-10-19, 18:32 UTC.** Existing OneMind Backstage is
preserved, not replaced: https://onemind-preview.tail717ba0.ts.net/.

- [x] Backstage 1.55.3 confirmed with pinned OpenChoreo 1.3.1 core/observability
      plugins, Node 22.23.3 and one frontend dependency cohort. This is our tested
      integration, not a claim that upstream's published 1.51 matrix changed.
- [x] Core/catalog/scaffolder/observability plugins installed in the existing host;
      OneMind home, Today, 130 launcher entries, 606 catalog entities, runbooks,
      settings, notifications and operational observations remain available.
- [x] `openchoreo-auth` Authentik login preserves signed-subject mapping to
      `user:default/akadmin` and the existing `onemind-admins` group.
- [x] One intersected permission policy: no admin bypass of upstream authorization;
      host non-admin mutations and machine platform/write access remain denied.
- [x] Built-in console disabled; its Deployment/Service/ingress/Application and
      the memory-only candidate are gone. Reused catalog credential and operator/
      catalog bindings remain. No PVC, database, backup or identity was removed.

**Release:** source `3dd34900a5b455696387394398b72f80ce3376ba`, GHCR digest
`sha256:b2fc101eacc504fadff99b638b342ce9bf56782a4901ea49dffe2fdd09f1967a`.
Main Deployment is **2/2 ready**, Argo Synced/Healthy. App-of-apps, control plane,
observability and identity resources also reconcile Synced/Healthy.

**Live gates:** pilot `stac-browser-v5` Active in Deploy; real Component Definition;
CPU/memory charts and runtime logs including a fresh HTTP 200 at 18:27:08 UTC.
Each new replica passed anonymous API 401, guest 404, machine catalog read 200,
catalog writes/platform/cache/context 403, and six read-only MCP tools. Actual
MCP health reads work. Non-admin human decisions have regression coverage; no
second human login was created or claimed live-tested. No live mutation was used
as a test. Existing PostgreSQL/preferences, integrations, SA/RBAC, PDB and
zero-unavailable rolling strategy remain intact. Today's encrypted backup reports
restoreVerified=true (24 databases/148 tables).

**Pitfalls corrected:** Request-body reuse, exact-origin CORS/CSP, provider owner
name prefix, and a multi-replica authorization cache gap. Capability results now
use the permission plugin DB with SHA-256 keys, upstream expiration and fail-closed
errors; no raw token is stored. Independent PostgreSQL connections and concurrent
initialization/upsert were tested; both production replicas use the shared DB.
The upstream ABAC router remains a version-locked private-API adapter to review
on every upgrade. Catalog sync deliberately polls; the unused webhook forwarder
is disabled rather than posting to a disabled HTTP receiver. Internal Backstage
PostgreSQL events and NATS are unchanged.

**Verification:** 49 tests, immutable Linux build, TypeScript, TechDocs, backend,
JSON exports, AMD64 SQLite, remote digest verification, changed YAML, strict
manifest/schema checks, pinned Helm renders, kustomize and new-secret scans pass.
GitHub Actions still cannot start due account billing; full `make validate` still
hits unrelated old YAML debt. Neither was represented as green or bypassed.

Implementation/rollout PRs: #311–#319. Detailed evidence and rollback:
[Phase 6 integration runbook](docs/operations/phase6-portal-integration.md).
Phases 7 and 8 remain **Not started**; this acceptance does not authorize them.

### Phase 6 release follow-through — local/Linux executor

2026-10-20: Zeus approved the existing verified local/Linux release checks while
GitHub cannot start hosted jobs, and requested a reusable deployment/removal
pipeline. No billing/check-disable/admin-bypass permission is implied. #330's
review added explicit incomplete-state and bounded-request regressions; its old
published candidate is superseded for these changes. Rebuild and remotely verify
the new source before a separate GitOps digest-pin/live acceptance step. Phase 6
remains In progress; no migration or removal is certified by source tests.

Pipeline implementation requested 2026-10-20:
[portable release/lifecycle gates](docs/operations/gitops-release-pipeline.md),
initial portal profile. Stages: local checks, Linux publication, registry/source
proof, separate pin PR, reviewed merge, read-only Argo/rollout verification and
separate semantic acceptance. Protected retirement is plan-only; no app, namespace,
PVC, Secret or backup deletion is executed. First Linux rebuild exposed floating
`npx` TechDocs dependencies; the CLI graph is now locked in workspace Yarn before
retry. Do not promote the failed build's old candidate or call the pipeline proven
end to end until a new immutable artifact and real rollout are verified.

Release execution update 2026-10-20: source PR [#330](https://github.com/onemindos/omos/pull/330)
is merged. Pipeline PR [#335](https://github.com/onemindos/omos/pull/335)
produced and remotely verified Linux/AMD64 source `c9a05f8d0b1fc1809b0496945d99593e884f561d`,
image digest `sha256:6c729bd20c0a0a40330be8f0c7bd5213c59c7540bd1ea54c38e05f8c63fc3e7c`.
All required local gates passed, including 82 security/STAC, 44 layout, 6 catalog,
20 pipeline and 7 release-command tests; the existing whole-tree YAML diagnostic
remains failing and hosted jobs remain unavailable. No rollout is claimed here;
separate pin, Argo and signed-in semantic/security acceptance still follow.

### Phase 7: Migrate Layer 3

Order: `omos-feeds`, `onetak`, `takml`, `onemind`, `agent-fabric`, `n8n`, `matrix`,
`openclaw`, then `hermes` last.

- [ ] One PR per project. Old manifests are deleted only after the new release is healthy.
- [ ] ResourceTypes for NATS streams and CNPG databases so components declare their
      dependencies.
- [ ] Bring `matrix` back in sync and put the OpenClaw Deployment in Git as part of their moves.

Gate per project: healthy in OpenChoreo, visible in the portal, old Argo resources removed.

### Phase 8: Expand

Explicit user scope, 2026-10-27: implement and deploy a small extensible internal
application MCP for the existing Legacy Hermes agent. This request supersedes the
earlier Hermes deferral for this bounded integration only; it does not start the
broader home-runtime or application-migration phases. Source, tests, image build
and manifests are in `services/onemind-mcp/` and `platform/hermes/mcp/`.
[Live acceptance](docs/operations/onemind-mcp-acceptance-2026-10-27.json) passed
2026-10-27: four authenticated read tools and app resource from Legacy, anonymous
rejection, both workloads available and Argo Synced/Healthy. Native Tool Search is
enabled. #487 implements the service; #491 supplies its dedicated pull credential;
#492/#493 make acceptance compatible with the installed SDK. This bounded slice is
delivered; the broader Phase 8/home-runtime work is not marked complete.

- Planning scope requested 2026-10-20: [full Kubernetes operating structure and
  execution queue](docs/architecture/kubernetes-operating-structure.md), backed by
  the [read-only baseline](docs/operations/kubernetes-structure-baseline-2026-10-20.md).
  This is **not** a later-phase deployment or completion. #330/#335/#336
  scoped portal release is verified; #331 records remaining trust/readiness gates; #284 remains the home expansion
  parent and must be revised against the actual Mac storage/power evidence.
- [ ] Decide and verify the home runtime against actual hardware: proposed Mac
      ARM64 Linux VM/K3s pilot; Talos remains a dedicated-host candidate, not a
      requirement to replace macOS. Register approved cluster visibility/GitOps
      and separately prove any second OpenChoreo data-plane integration.
- [ ] Hermes gets a model through LiteLLM and the same read-only tools through the tool gateway.
- [ ] Revisit `takserver` placement.

## Decision log

| Date | Decision |
|---|---|
| 2026-10-18 | Use OpenChoreo for our own apps (Layer 3). Foundation stays on plain Argo CD. |
| 2026-10-18 | Argo CD remains the GitOps engine; OpenChoreo resources are applied by Argo CD, not Flux. |
| 2026-10-18 | Claude Code is the AI operator. Hermes is deferred to Phase 8. |
| 2026-10-18 | Keep the custom portal as the long-term home; use the built-in console only until Phase 6. |
| 2026-10-18 | LiteLLM stays the model gateway. The tool gateway (LiteLLM MCP or agentgateway) is decided in phase 8, when an in-cluster agent needs it. |
| 2026-10-30 | Owner chose **agentgateway** (standalone chart, no CRDs) as Legacy's tool front door, replacing the LiteLLM MCP hub of 2026-10-29; LiteLLM returns to model routing only. Destructive and physical tools sit on a separate `execute` route behind an out-of-band `execute <id>` approval. Design: [onemind-tool-gateway.md](docs/architecture/onemind-tool-gateway.md). |
| 2026-10-18 | Cilium Gateway API backends need a component-scoped `fromEntities: ingress` policy; keep OpenChoreo's generated policy and add an opt-in trait rather than allowing all ingress. |
| 2026-10-18 | AI read access to the cluster uses ServiceAccount `kube/ai-operator` (built-in `view` plus platform CRDs, no Secrets) through kubernetes-mcp-server in read-only mode. No separate Argo CD MCP. |
| 2026-10-18 | Not adopting kagent, KubeVela, Kratix or Humanitec. HolmesGPT is optional after Phase 5. (KubeVela part superseded 2026-10-21, below.) |
| 2026-10-20 | Zeus requested full Kubernetes/Backstage structure and linked PRs. Keep the installed stack and existing phase gates; permit planning and bounded comparison of KubeVela/Crossplane/fleet candidates without deploying or assigning a second owner to live resources. The new queue is a proposal, not a replacement approval or claim of verified integration. |
| 2026-10-21 | Zeus selected Backstage → reviewed Git PR → Argo CD → KubeVela (app lifecycle) + Crossplane (app infrastructure) as the target; Terraform keeps its recorded infrastructure, NATS carries operational events. Recorded in [decision 0002](docs/architecture/decisions/0002-application-delivery-and-infrastructure-boundaries.md). Selection only: nothing installed, no version approved, no ownership transferred from Terraform or OpenChoreo, existing OpenChoreo workloads stay until a separately approved transition. kagent, Kratix and Humanitec stay not adopted. |
| 2026-10-23 | Zeus approved retiring OpenChoreo, then OpenSearch and External Secrets. The portal no longer calls OpenChoreo. Teardown is a stacked PR series (observability modules, then CR-holding apps, then observability plane and data plane, prereqs and gateway, control plane, OpenSearch, External Secrets). Phases 2 to 6 above are historical. Stage 1 (observability modules) removes the logs, traces, FinOps and events apps. Stage 2 removes the custom-resource apps (projects, platform, defaults) including the STAC Browser; the data plane must still be running while they go. Stage 3 removes the observability plane (including the observer alert store and its VolSync source) and the data plane. Stage 4 removes the External Secrets prerequisites (`openchoreo-secrets`, `ClusterSecretStore default`) and the `onemind-private` GatewayClass manifests. Stage 5 removes the control plane, the OpenChoreo identity objects in Authentik, the OpenChoreo secrets and the last namespaces; no OpenChoreo object remains in Git. Stage 6 removes OpenSearch (20 GiB volume, snapshot and restic repository) once nothing queries it. Stage 7 removes External Secrets (operator, 23 CRDs, 2 webhooks); nothing used an ExternalSecret. |

## Incidents found while working the roadmap

| Date | What | Fix |
|---|---|---|
| 2026-10-18 | DOKS replaced the hub-core nodes (v1.36.5). `tak-ldap-0` and `takml-server` stayed Pending, pinned to the removed node `onemind-hub-core-v2-37phkg` (the `takml-server` pin was set live, not in Git). `nats-2` could not schedule: the only node without a NATS pod was full while a new node sat nearly empty. Found by the phase 1 read-only access. | [#259](https://github.com/onemindos/omos/pull/259): node-pool affinity for `tak-ldap`, pod affinity to `takml-vision-worker` for `takml-server`, stale live pin removed. The Hermes pod was rescheduled to make room for `nats-2` (NATS back to 3/3). [#260](https://github.com/onemindos/omos/pull/260): `Recreate` for `takml-db`, whose rolling update deadlocked on its ReadWriteOnce volume. |

Follow-ups: nothing rebalances pods after DOKS replaces nodes. Consider the Kubernetes
descheduler (Layer 1) in phase 5 alongside the capacity work.

## Risks and open questions

- **Maturity:** OpenChoreo is a CNCF Sandbox project with frequent releases and several
  supported lines at once. Pin chart versions and read release notes before upgrades.
- **Backstage version:** upstream documents 1.51.0; our 1.55.3 integration passed Phase 6 live gates. Keep the dependency cohort and private ABAC-router adapter pinned; reverify upgrades.
- **Gateway:** Cilium Gateway API works with the OpenChoreo data plane (phase 2c), but
  its CRDs are v1.2.1 (DOKS-managed) while OpenChoreo targets v1.5.1. The upstream
  `service` component type's CORS filter cannot be used; use `onemind-service`.
  Re-check when DOKS upgrades Cilium or Gateway API.
- **Local-first mission:** `docs/MISSION.md` requires home operation without the cloud.
  Before Phase 8, verify what a home data plane does when the cloud control plane is
  unreachable (running workloads should continue; changes cannot be applied), and decide
  whether the control plane eventually moves home.
- **Capacity:** the observability plane is the heaviest addition; it waits for Phase 5.

## References

- [OpenChoreo architecture](https://openchoreo.dev/docs/overview/architecture/)
- [Run in your environment](https://openchoreo.dev/docs/getting-started/try-it-out/on-your-environment/)
- [Multi-cluster connectivity](https://openchoreo.dev/docs/platform-engineer-guide/multi-cluster-connectivity/)
- [GitOps](https://openchoreo.dev/docs/platform-engineer-guide/gitops/overview/)
- [Backstage plugins](https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/overview/),
  [installing into an existing Backstage](https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/installing-into-existing-backstage/),
  [compatibility matrix](https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/compatibility-matrix/)
- [Built-in console configuration](https://openchoreo.dev/docs/platform-engineer-guide/backstage-configuration/)
- [MCP servers](https://openchoreo.dev/docs/ai/mcp-servers/)
- [Community modules](https://github.com/openchoreo/community-modules)

### Portal GitOps promotion — 2026-10-20

Source #330 and pipeline #335 merged normally using the explicitly approved
local/Linux release evidence (hosted CI unavailable, not green). Promotion
[#336](https://github.com/onemindos/omos/pull/336) was generated by the new pipeline:
only the portal Deployment image changes; no Argo policy, credentials, database,
RBAC, other workload or removal. Verified source `c9a05f8d0b1fc1809b0496945d99593e884f561d`,
digest `sha256:6c729bd20c0a0a40330be8f0c7bd5213c59c7540bd1ea54c38e05f8c63fc3e7c`.

- [x] Publish Linux image, remotely verify immutable source/platform/runtime manifest.
- [x] Run promotion gates and open separate digest-pin PR.
- [x] Verify actual merged revision, Argo Synced/Healthy and two exact-image ready replicas.
- [x] Signed-in component/project Operations, topology, logs/metrics/runbook and scoped security checks.

The old `af7770fb` / `fab652...` candidate in earlier follow-up drafts is superseded
and must not be promoted. STAC collection/search HTTP 500 and data-unknown remain
real backing-service findings, not portal release failures or authorization to
change database schemas. Phase 6 remains In progress; no app migration certified.

Live evidence 2026-10-20, 17:29–17:32 UTC: #336 merged at
`1545b50f4fb5521dbdfc9eb22fa302cf2e4b45d6`; pipeline exact-revision/digest verification
passed, 2/2 ready/updated/available replicas, generation 72. Both pods passed
anonymous/machine denial and six-read-only-MCP-tool checks. Signed-in Component and
Project Operations, real logs/metrics and packaged STAC runbook passed. Catalog
converged to 716 entities, zero duplicates/errors/orphans, expected topology, only
the known inactive workflow-plane target missing. Details and remaining limits:
[STAC release evidence](docs/operations/stac-operational-slice-2026-10-20.md).
Phase 6 is still In progress; STAC schema, broader trust/recovery and production
expiry/second-role coverage are not certified by this release. No migration begun.

### Kubernetes structure review — PR #333, 2026-10-20

[#333](https://github.com/onemindos/omos/pull/333) reconciles the 20-category
[operating structure](docs/architecture/kubernetes-operating-structure.md) with
merged #331/#335/#336 and the fresh 17:40 UTC baseline. Existing integrations,
upstream-supported but untested links and uninstalled candidates are distinguished.
NATS StatefulSet Argo ownership is now observed; do not blindly re-adopt from an
older report. SealedSecrets remains the repository credential pattern despite
ESO controller presence. Portal pipeline proof is retained; other apps require
reviewed profiles and removal remains plan-only. No new image/deployment required.

Next structural implementation: #334 read-only public-ingress parity assessment
and K02 ownership/recovery inventory. K03 will reuse existing onboarding templates
and add a tested app/profile contract, not a second scaffolder. Remaining K00 trust
checks stay open; no Phase 6 completion, Phase 7 migration, controller installation,
OS/storage/provider change or widened AI permission is implied by this plan.


### Native visibility first win — 2026-10-20

Bounded Phase 6 follow-through: restore the installed Kubernetes plugin's controller
matching for Backstage/PostgreSQL and explicitly map the NATS broker to its observed
`platform-nats` Argo Application. Stable entities, ownership and dependency graph
are retained; direct native Kubernetes/Argo/runbook links and refresh regression
checks are added. See [scope and acceptance](docs/operations/native-portal-visibility-2026-10-20.md).
[PR #337](https://github.com/onemindos/omos/pull/337) merged as `e98f1436`.
Scoped production acceptance passed at 18:50–18:54 UTC: native Backstage/DB/NATS
workloads and metrics, correct Argo summaries, source/runbook/dependency links,
read-only MCP metadata and existing denial boundaries. Portal 2/2, generation 73,
same verified digest; DB spec/revision unchanged; NATS 3/3. Native Argo history
revision is distinct from the latest compared Git revision (documented for NATS).
The catalog ConfigMap hash naturally rolls the same image; no database template,
image digest, RBAC, credentials, OpenChoreo ownership or event contracts change.
Phase 6 stays In progress and Phase 7 Not started. After browser acceptance, next
visibility work is the existing Grafana 401, then topology/Prometheus/OpenCost with
separate compatibility and access checks. This does not install new controllers.


### Grafana durability and native visibility — 2026-10-20

User approved a 5 GiB persistent Grafana volume, backed-up data migration, a brief
Grafana maintenance window and a dedicated read-only Backstage credential. The
existing database is on emptyDir; read-only inspection found zero service accounts
and the old portal token returned 401. The native plugin/proxy prefix combination
is being verified rather than blindly rewritten. No AI action permissions change.

Stage A prepares an Argo-owned retained PVC, seed job, encrypted remote-verified
SQLite backup and namespace-bound sealed credentials. Grafana stays on its old
pod until the seed/backup gate passes; the scheduled backup starts suspended.
Stage B is a separate reviewed cutover to the seeded PVC with Recreate and an
init guard, then catalog dashboard mappings and browser verification. Initial
preparation is not proof of migration or fixed authentication. Phase 6 remains
In progress; Phase 7 Not started. PR links and actual acceptance follow below.


Grafana Stage A [#340](https://github.com/onemindos/omos/pull/340) merged as
`9f9358f6`. Retained PVC bound; original database capture quiesced the server for
0.489 seconds, resumed it, then verified encrypted remote download/decrypt/SQLite
restore across 87 tables before seeding the claim. Seed job completed. Original
Grafana remains available until the separately reviewed Stage B cutover.

Stage B prepares the seeded RWO Recreate switch/seed guard, online scheduled and
one-time verified backup, dedicated token reference and three fixed-workload
Prometheus dashboards mapped into Backstage. Existing host extensions are retained
when configuring the native Grafana alert card to require an explicit alert-label
selector; missing Grafana rules must not look like full Prometheus alert coverage.
No iframe/authentication protection is weakened, no new portal image is required,
and actual post-cutover/native browser acceptance remains pending.


Grafana Stage B [#341](https://github.com/onemindos/omos/pull/341) merged as
`dfe7611c`. Argo server-side diff safely blocked the Grafana update: the original
Helm field owner retained `emptyDir`/`rollingUpdate`, creating invalid union fields.
Original Grafana and the verified PVC seed remain intact. Follow-up uses the
supported legacy diff for this one Application and `Replace=true` on the fully
rendered Grafana Deployment only; no Force, PVC deletion, force sync, new owner,
automated-prune change or wider permissions. Full server replace dry-run and local
gates must pass before the normal exact-head PR merge. Acceptance still pending.


Grafana durability/native visibility **accepted 2026-10-20 22:29–22:35 UTC**.
[#342](https://github.com/onemindos/omos/pull/342) merged `f145f14d`; all three
owning apps Synced/Healthy at that revision. Grafana 3/3 on retained 5 GiB PVC,
Deployment UID/image preserved; enabled daily backup plus successful post-cutover
encrypted remote download/decrypt/SQLite restore verification. Portal generation
74, 2/2, same digest; DB unchanged; NATS 3/3. Migrated Viewer identity verified on
both replicas; native Backstage/DB/NATS dashboard links and real CPU/memory panels
passed. Anonymous/machine proxy and Kubernetes mutation denials remain intact.
See [acceptance and remaining limits](docs/operations/grafana-durable-visibility-2026-10-20.md).
Root-folder link presentation, old ClickHouse panels, Prometheus alert integration,
expiry monitoring and wider topology/OpenCost remain separate work. Phase 6 stays
In progress; Phase 7 Not started; #339 action-layer draft is not deployed by this.


### Native Grafana folder links — 2026-10-20

Next bounded Phase 6 fix: provision the three existing service dashboards in a real
`OneMind` folder using the pinned chart's sidecar folder annotation and filesystem
folder provider. Native plugin 1.1.0 otherwise concatenates an absent root-folder
URL into `undefined/`. Dashboard UIDs, queries, selectors, credentials, PVC,
versions and permission boundaries remain unchanged. Unannotated dashboards keep
the root default. Argo owns the brief Grafana restart; no imperative move/apply.
Live folder metadata, signed-in native links, unchanged dashboard inventory and
healthy persistent rollout must pass after merge. Then Prometheus alert visibility
is next, separately scoped; Phase 6 stays In progress and Phase 7 Not started.


Grafana folder fix [#344](https://github.com/onemindos/omos/pull/344) accepted
23:01–23:04 UTC at `48d841b3`: both owning apps Synced/Healthy, Grafana 3/3 on
same PVC/image/Deployment UID; portal remains 2/2, DB unchanged, NATS 3/3.
All 31 dashboard UIDs retained, only three scoped boards moved to `OneMind`;
Viewer metadata and signed-in native folder links/listing verified. Read-only
access/denials preserved. [Evidence](docs/operations/grafana-native-folder-2026-10-20.md).
Prometheus alert integration is the next separately verified task; not installed
by this repair. Phase 6 stays In progress; Phase 7 Not started.


### Native scoped Prometheus alerts — 2026-10-20

User requested the next visibility slice after #344/#345. Prepare the native
Roadie Prometheus alert table with the supported compatibility bridge, a bounded
user-only backend adapter for Backstage/DB/NATS, refresh/timestamps and honest
no-match/partial/unavailable states. Only reviewed fixed GETs/scopes; no generic
PromQL proxy, new credentials, mutation rights or rule changes. Tests join both
portable local/Linux and image gates; source merge is distinct from built image,
GitOps pin and live acceptance. [Scope/gates](docs/operations/prometheus-native-alerts-2026-10-20.md).
Phase 6 stays In progress; Phase 7 Not started. Actual PR/release evidence follows.


Prometheus alert source [#346](https://github.com/onemindos/omos/pull/346) merged
`cbf0e637` (implementation `725c5e4e`). Required local/Linux gates passed;
29 new scoped frontend/backend/route checks plus catalog regression are enforced
in both executors. Real 167-rule source parsed; no selected-scope active match
at preflight, not an all-clear. Linux image published with source/platform/digest
verified: `sha256:e103bf7c0103b615bef2f4b9dfcb603ef952af7d9eb95d498117d1f48a5dc0e6`.
[#347](https://github.com/onemindos/omos/pull/347) promotes that immutable image;
Argo/browser/replica/denial acceptance remains pending this normal merge.
GitHub CI billing block and historical whole-tree YAML debt remain explicit.
No Phase 6 completion, Phase 7 start or wider AI action permission is implied.


Prometheus scoped visibility **accepted 23:30–23:34 UTC** after #347 merge
`f27bb770`: `platform-backstage` Synced/Healthy, generation 76, 2/2 actual imageIDs
match the verified published digest; DB spec/revision unchanged, Grafana retained
storage/template unchanged, NATS 3/3. Native Backstage/DB/NATS alert cards, scoped
metadata/links/timestamps and real pending-alert → no-match transition verified;
user GET 200 observed on both replicas. Anonymous/machine/write denials retained.
[Acceptance and limits](docs/operations/prometheus-native-alerts-2026-10-20.md).
Firing/outage scenarios are fixture evidence, not a production alert-delivery drill.
Next: topology/OpenCost; broaden coverage and accessible native state text in
bounded follow-ups. Phase 6 stays In progress; Phase 7 Not started; #339 unmerged.


### Native dependency and location topology — 2026-10-20

Next requested Phase 6 slice: enrich the installed native graph with reviewed
Backstage runtime/identity/NATS/Grafana/Argo feature dependencies and five existing
PVC catalog identities; fix namespace-group titles without changing entity IDs.
Add bounded forward/reverse/core graph and location links, preserve current
Kubernetes/Grafana/Prometheus/Argo/runbook mappings, and enforce refresh/source
contracts. This is source-reviewed metadata, not live traffic or all-app impact.
No new controller, PVC, permission, image, migration or cloud charge. Mounted
catalog changes naturally roll the same portal image through Argo. [Scope/gates](docs/operations/native-service-topology-2026-10-20.md).
Source merge/Argo/catalog/browser acceptance remain distinct. Phase 6 In progress;
Phase 7 Not started; #339 remains draft. Actual release evidence follows.

Native topology [#349](https://github.com/onemindos/omos/pull/349) **accepted**
2026-10-21 UTC at `82c05d3c`: owning Argo Synced/Healthy, portal generation 77,
2/2 on the same verified image; DB/Grafana specifications unchanged, retained
Grafana storage/backup and NATS 3/3 intact. Both replicas return reviewed catalog
relationships without selected-entity errors. Signed-in core graph shows six
services, five existing PVC records and NYC3; reverse DB graph/location navigation
verified. Existing native Grafana/Prometheus views and read-only denials retained.
[Receipt and limits](docs/operations/native-service-topology-2026-10-20.md).
This is documented dependency metadata, not live traffic/all-app coverage.
**Next: OpenCost**, separately scoped. Phase 6 stays In progress; Phase 7 Not
started; #339 remains draft. Hosted CI billing and whole-tree YAML debt remain.

### Eight-workstream parallel preparation run — 2026-10-21 UTC

User authorized separate PRs and agent distribution for the remaining Backstage
backlog. [Run index, owners, PRs and dependencies](docs/roadmap/backstage-parallel-run.md)
tracks cost, coverage, catalog, trust, lifecycle, actions, agents and home/edge.
Existing #339/#332/#284 are reused. Bounded source/test preparation runs in
parallel; production review/rebase/merge/release remains one PR at a time.
Correct the stale handoff prompt that claimed Phase 6 complete; reconcile scoped
Grafana/Prometheus/topology acceptance without erasing the remaining gaps.
No new permissions, storage, deployment, app cutover or phase completion.
Phase 6 stays In progress; Phases 7/8 Not started.


### Backstage capability and UI/layout foundation — 2026-10-21

User requested a foundation review of the supplied application-flow/plugin list,
not runtime implementation. Master #356 now includes the
[experience foundation](docs/roadmap/backstage-experience-foundation.md): native
feature/source matrix, proposed home/navigation/entity layout, evidence states,
search/docs/events/AI integration gaps, dual MUI/BUI theming, accessibility,
performance and U00–U06 future implementation/acceptance packages.
[Decision 0002](docs/architecture/decisions/0002-application-delivery-and-infrastructure-boundaries.md)
records user-selected KubeVela/Crossplane target boundaries, with Terraform, Argo,
OpenChoreo, domain operators and NATS ownership explicit. No controller install,
OpenChoreo removal, UI/dependency change or new AI permission is included.
The eight existing workstream PRs remain drafts; no new UI implementation PR or
production acceptance is claimed. Phase 6 In progress; Phases 7/8 Not started.


### OneMind mission pillars in the catalog — 2026-10-21

User explicitly requested SO / LE / HP / GE catalog structure as part of the
foundation. [Specification](docs/roadmap/onemind-catalog-pillars.md) records
Sovereign Operations, Legacy Evolution, Holistic Performance and Generational
Empire using the existing doctrine, proposed native Domains/System mappings,
shared/unclassified discovery and stable-reference/permissions/curation gates.
Master #356 and the UI foundation are updated; no live catalog or workload is
reclassified, no new identity permission or NATS routing is granted. Historical
numeric ranges remain unconfirmed, not reordered with the requested UI order.
Phase 6 In progress; Phases 7/8 Not started.

### Flagship portal program — 2026-10-21

Zeus asked for Backstage to become the single, organized front door for every
OneMind service (Kubernetes and not), with live dashboards, DigitalOcean at a
glance, add/remove lifecycle and a brand-correct theme. A signed-in review found
7,750 unread notifications, an alphabetical sidebar, about 70 Systems under three
naming schemes, 6 of 69 Argo apps linked, an empty Cost page, an unrequested
GitHub OAuth pop-up on app pages, a navy/teal theme, and an offsite backup copy
that has never run (two other "suspended" backup jobs are deliberate replacements). Eleven work orders, one draft PR each, for implementing
agents: [index](docs/roadmap/flagship/README.md), F1–F11 =
[#357](https://github.com/onemindos/omos/pull/357)–[#367](https://github.com/onemindos/omos/pull/367).
This supersedes the U00–U06 order in #356. It is Phase 6 follow-through: no
Phase 7 start, no app migration, no AI permission change. Work-order PRs do not
edit this file; the coordinator records each merge here.

### Parallel portal run — coverage

User requested a separate draft workstream for coverage. Bounded offline source
and regression contracts are prepared; see [scope and remaining gates](docs/operations/backstage-coverage-workstream.md).
Tests join existing `test_gitops*.py` discovery; a source pass is not live acceptance.
No new access, storage, controller, deployment or migration. Phase 6 In progress;
Phase 7 Not started. Review/rebase and merge only one accepted PR at a time.

### Parallel portal run — catalog

User requested a separate draft workstream for catalog. Bounded offline source
and regression contracts are prepared; see [scope and remaining gates](docs/operations/backstage-catalog-workstream.md).
Tests join existing `test_gitops*.py` discovery; a source pass is not live acceptance.
No new access, storage, controller, deployment or migration. Phase 6 In progress;
Phase 7 Not started. Review/rebase and merge only one accepted PR at a time.

### Parallel portal run — trust

User requested a separate draft workstream for trust. Bounded offline source
and regression contracts are prepared; see [scope and remaining gates](docs/operations/backstage-trust-workstream.md).
Tests join existing `test_gitops*.py` discovery; a source pass is not live acceptance.
No new access, storage, controller, deployment or migration. Phase 6 In progress;
Phase 7 Not started. Review/rebase and merge only one accepted PR at a time.

### Parallel portal run — cost

User requested a separate draft workstream for cost. Bounded offline source
and regression contracts are prepared; see [scope and remaining gates](docs/operations/backstage-cost-workstream.md).
Tests join existing `test_gitops*.py` discovery; a source pass is not live acceptance.
No new access, storage, controller, deployment or migration. Phase 6 In progress;
Phase 7 Not started. Review/rebase and merge only one accepted PR at a time.

### Flagship wave 1 released — 2026-10-22

Merged and released as one verified portal image: F10 theme [#366](https://github.com/onemindos/omos/pull/366),
F1 navigation and home [#357](https://github.com/onemindos/omos/pull/357), F2 catalog structure
[#358](https://github.com/onemindos/omos/pull/358) (+ image-context test fix [#370](https://github.com/onemindos/omos/pull/370)),
F3 alerts and notification triage [#359](https://github.com/onemindos/omos/pull/359). Source `7e2eea59`;
image `sha256:37f17c546ffa4254363c003b7a8112dd1692567d40bbca92aca1b463616b0c3f`; previous image
`sha256:e103bf7c0103b615bef2f4b9dfcb603ef952af7d9eb95d498117d1f48a5dc0e6` (rollback = reviewed revert of the pin).
Pin PR [#371](https://github.com/onemindos/omos/pull/371) merged at `cf04980a`; `pipeline.py verify` passed (exact
revision, digest, generation, 2/2 replicas). Signed-in checks passed: black/red theme in dark mode, grouped sidebar,
home status strip and favorites, five pillar Domains with 24 Systems, old System URLs redirecting, ⌘K search, Pillar filter.
Catalog data went live first (Argo applied #358 on 10-22 10:25Z, before the image carried the redirects); old System URLs
returned not-found for about 17 minutes in between.

Not yet done: the Pillars sidebar group and home tiles stay behind `onemind.features.pillars` (default off); the
Alertmanager webhook is disabled until Zeus seals `onemind-alert-webhook` (monitoring + backstage) and the receiver is moved out of
`platform/monitoring/pending/`; the 7,750-notification cleanup script has not been run; the home status strip still counts the
two deliberately superseded backup jobs as "to review" (F4); offsite backup copy still has no credentials. Phase 6 In progress;
Phase 7/8 Not started. GitHub-hosted CI remains billing-blocked; whole-tree YAML diagnostic debt unchanged.

### Flagship wave 2 released — 2026-10-22

Released as one verified portal image: F4 Live Operations [#360](https://github.com/onemindos/omos/pull/360), F5 standard app page
and coverage contract [#361](https://github.com/onemindos/omos/pull/361), F11 live inventory [#367](https://github.com/onemindos/omos/pull/367),
light-mode fix [#374](https://github.com/onemindos/omos/pull/374), alert-webhook sealing helper [#375](https://github.com/onemindos/omos/pull/375)
and pillar tile counts [#376](https://github.com/onemindos/omos/pull/376). Source `d3195054`; image
`sha256:5ae28ac6455bba5112090707ce9b9e6811e83a60d6a27a88e5a25444a5905f31`; previous image
`sha256:37f17c546ffa4254363c003b7a8112dd1692567d40bbca92aca1b463616b0c3f` (rollback = reviewed revert of the pin). Pin PR
[#377](https://github.com/onemindos/omos/pull/377) and the config PR [#373](https://github.com/onemindos/omos/pull/373) (pillar navigation and
live inventory enabled through the mounted cluster config) were applied by Argo in one sync, one rollout, 2/2 on the new digest. The
pipeline `verify` stopped on its source-changed guard (the config PR touched one portal test script after the image was built), so the
same checks were made directly: Argo Synced/Healthy at `267831fa`, deployment image equals the release digest, 2/2 ready, no restarts,
no error or warning logs, log line "Live inventory: 175 workloads tracked, 175 defined in catalog/cluster.yaml, 0 published".
Signed-in checks passed: Live Ops (Alertmanager, Argo, workloads, storage, backups), Home status strip with "needs setup 1 · retired 2",
pillar tiles SO 116 · LE 1 · HP 5 · GE 14, the standard app header and Logs/Cost/Runbook tabs, the generic OneMind workload dashboard,
and no GitHub sign-in dialog (Connect GitHub prompt instead).
Before the image: the three backup CronJob annotations, the F5 catalog annotations and the Grafana dashboard were applied by Argo on merge.
Not yet done: offsite backup copy credentials; the 7,757-notification cleanup; alert webhook secrets (run `scripts/ops/seal-alert-webhook.sh`);
F6 DigitalOcean, F7 non-Kubernetes apps and health checks, F8 add/retire templates, F9 plugin pack. Phase 6 In progress; Phase 7/8 Not started.

## User-directed file synchronization — 2026-10-23

Zeus requested immediate bidirectional Mac/Android Obsidian synchronization via
Nextcloud, independently of the current portal phase. This change adds a Layer 2
shared file service under Argo CD; it does not complete or bypass phase gates.

| Work | Status | Evidence |
|---|---|---|
| Nextcloud GitOps deployment and full-library Mac/Android sync | Deployed and healthy; native client setup and acceptance in progress | [Runbook](platform/nextcloud/README.md); [PR #424](https://github.com/onemindos/omos/pull/424); [PR #425](https://github.com/onemindos/omos/pull/425) |

Acceptance requires actual edits flowing in both directions, offline and conflict
checks, a reviewed real-vault upload set, and a coordinated backup/restore drill.
The native Mac upload has started. Full transfer, two-way device acceptance and
recovery verification remain distinct from the observed healthy server deployment.

## Paperclip retired — 2026-11-01

The 2026-10-23 user-directed Paperclip installation was removed by owner decision on
2026-11-01 after a final database dump and volume snapshots. See
[Paperclip: retired](docs/STACK.md#paperclip-retired--2026-11-01) and the
[retirement runbook](docs/operations/paperclip-retirement-2026-11-01.md). Phase 6 stays
In progress; Phases 7/8 stay Not started.

## User-directed personal Vaultwarden — 2026-10-25

Zeus authorized a private personal password vault on the existing Kubernetes stack.
This bounded Layer 2 deployment uses native Kustomize/Argo, one SQLite writer,
protected storage, private HTTPS, closed registration and encrypted backups.
It does not advance the portal migration. Phase 6 stays In progress; Phases 7/8
stay Not started.

| Work | Status | Evidence |
|---|---|---|
| Private personal Vaultwarden | Deployed and healthy; encrypted initial empty-vault archive restored and server started in isolation; owner enrollment/client acceptance pending | [Deployment and acceptance](platform/vaultwarden/README.md); [recovery receipt](platform/vaultwarden/BACKUP.md); [PR #466](https://github.com/onemindos/omos/pull/466); [PR #467](https://github.com/onemindos/omos/pull/467) |

No 1Password migration, agent vault permission, user master password, model billing
change or new Kubernetes authority is included. Owner login, client/passkey sync
and future AI integration require separate authenticated acceptance.

## User-directed policy warning cleanup — 2026-10-26

Zeus requested repair of image-version/resource-budget and telemetry pull-secret
warnings. Existing running app images are pinned; remaining background warnings
include retained old revisions, alongside real budget gaps on TAK's legacy
PostgreSQL/NATS workloads and retention CronJob. Keep all audit/admission checks
for active workloads and preserve rollback history and source data.

| Work | Status | Evidence |
|---|---|---|
| TAK stateful budget safety preflight | Fresh encrypted archives/recovery verified; no-op Argo adoption passed with unchanged pod/PVC UIDs; measured budgets and scoped history-aware scans prepared | `platform/backups/volsync/sources.yaml` |
| Telemetry stale pull-secret reference | Root cause identified; guarded default-service-account repair and sequential provider telemetry pod replacement in verification | [Runbook](docs/operations/2026-10-26-telemetry-pull-secret.md) |

This is bounded operational maintenance. Phase 6 stays In progress; Phases 7/8
stay Not started. Source preparation and successful live verification remain
separate. No node reduction, credential rotation or broad Helm upgrade is included.

Runtime follow-through: TAK budget rollout is healthy and the original image/resource
checks pass on active app controllers/pods. Correct the pre-existing security-parent
ignore/prepatch interaction that left intended audit history guards unapplied; use
explicit observed Kyverno defaults and verify live policy generation and guards.

Final ownership boundary: keep provider-marked DOKS system agents under DigitalOcean
budget management in background scans, while retaining image/admission/other security
checks. Adopt the manually installed unlabelled metrics-server Deployment with a
measured512Mi cap and unchanged image, requests, Services and RBAC. This adds an Argo
ownership record for an existing workload; full live acceptance follows reconciliation.

## User-directed cluster follow-through — 2026-10-26

Zeus requested the remaining cluster findings be worked. This bounded maintenance
keeps Phase 6 In progress and Phases 7/8 Not started.

| Work | Status | Evidence |
|---|---|---|
| Namespace layer-label drift | Verified live: all12 corrected, every policy-covered namespace labelled, original namespace UIDs preserved | `platform/security/shared-namespaces.yaml`; server dry-run preserves namespace UIDs and all existing metadata; no live NetworkPolicy selects this label |
| Failover memory headroom | Measured right-sizing and placement review in progress; keep four nodes | Live KubeMemoryOvercommit, no node pressure |
| Independent backup copy | Destination/credential discovery in progress; do not unsuspend without working recovery credentials | `platform/backups/OFFSITE-COPY.md` |
| Remaining TAK workload ownership and broad drift ignores | Verified live: API/config/Ignite adopted without pod/spec/restart changes; core/TAK broad ignores narrowed; supporting legacy resources remain separate | `docs/gitops/HELM-ADOPTION.md` |

Namespace acceptance follow-through: Dapr's managed-namespace label remained absent
following sync; use a protected explicit Namespace under security, preserving its
VolSync annotation and existing SSA option, rather than claiming source metadata
was applied. Existing namespace UIDs remain unchanged.

Further acceptance: exact-live TAK API/config/Ignite specs pass server dry-run
parity; protected metadata adoption and narrowed TAK ignores are merged for live
verification. Measured Nextcloud OOM corrections are prepared with576Mi additional
reservation; independent copy verification now has fixed-selection regression
coverage. The reviewed five-worker alternative needs an owner cost decision
(+$168/month); independent backup activation needs a destination and credentials.
Neither unresolved dependency is represented as a completed repair.

TAK support follow-through:19 existing Services/RBAC/configuration/helper objects
prepared for protected adoption with server dry-run parity; remove only the
Ignite ServiceAccount's missing legacy pull-secret reference. Sensitive legacy
ConfigMaps/runtime Secrets remain explicitly outside plaintext Git and require
an encrypted-source migration. No full-release Helm operation is authorized.

Live scanner correction: fresh history warnings revealed native-int64 comparisons
in the reports controller were not represented by JSON CLI fixtures. Numeric
normalization and exact pinned native-JMESPath regression are prepared, preserving
all real admission checks and fail-closed unknown history. Actual reports-controller
skip metrics are required for live acceptance; prior CLI-only evidence was incomplete.

Support adoption is verified live: all 19 original resource UIDs preserved, 18 runtime
payloads identical, only the missing Ignite SA pull-secret reference removed; all
running TAK pod UIDs and restart counts unchanged. Nextcloud app/proxy are Ready
with 0 restarts after the measured memory corrections; original data/cache PVC/PV
identities and capacity are preserved. Separate concurrent Hermes work added one
healthy Application; the latest snapshot is 59 Synced/Healthy apps and 84 Bound PVCs.

Scanner live acceptance verified after PR484: generations 3/4 and generated rules
contain numeric normalization; the real reports controller records 9 skips per rule
for retired Matrix ReplicaSets, with its old failure counter stable. No new image/
resource warnings during 150 seconds after sync;815-resource live engine run yields
1446 pass / 999 expected skips / 0 fail / 0 error.153 regression tests and CI pass. The privileged
GitHub DinD runner is a separate real open risk; do not hide it with an audit exception.
Capacity/independent recovery remain pending owner choices. Phase 6 stays In progress.

## User-directed PR merge preparation — 2026-10-27

Zeus initially delegated straightforward owner-item preparation, then approved
all owner decisions and merge execution for dcarku's #495–#540 batch, retaining
Telegram. [Merge queue, approval and evidence](docs/roadmap/2026-10-27-owner-merge-preparation.md)
records the source candidates, approved 24-hour retention/read-only Headlamp/
Telegram/enforcement choices, and #532/#535's remaining technical holds. The
namespace-label bootstrap race must be fixed before enforcement. Technical
credential, CI, image and rollout gates still apply; no approval is a claim that
they passed. No merge or deployment is claimed by this recording task. Phase 6
remains In progress; Phases 7/8 remain Not started.

User-directed MCP coverage follow-through (2026-10-27): expansion prepared from the
initial NATS/Prometheus slice to seven configured app integrations, 12 tools and an
explicit 32-domain access catalog. CI requires coverage dispositions for all 56
Git-declared Applications. Credential-gated and unimplemented app integrations
remain blocked and visible; broad coverage is not claimed complete. No cluster
administration, vault-content permission or physical actuation is added. Runtime
acceptance must follow the pinned-image rollout; wider Phase 8 stays unfinished.

## Nextcloud retirement — 2026-11-01

Zeus authorized retiring Nextcloud and its 200 GiB file/1 GiB cache volumes
to use Dropbox. Deployment manifests are archived under
`platform/parked/nextcloud`; the Argo Application is removed. Execution requires
a completed fresh volume snapshot plus private database/config exports.
Provider deletion and recovery verification must be recorded from live evidence;
this source change alone does not certify either. Preserve device originals and
existing backups. See `platform/parked/nextcloud/RETIREMENT.md`.

## User-directed independent Git recovery — 2026-11-01

- [x] Prepare private Forgejo15LTS, retained single-volume SQLite/Git/LFS storage, restricted bootstrap and manual-sync Argo source.
- [x] Prepare a scoped coordinator that quiesces the writer, exports an encrypted dated VolSync/Restic backup, then restores its prior replica count.
- [ ] Supply strict SealedSecrets, review current capacity/storage/backup window and activate through Argo.
- [ ] Enroll owner, pilot one GitHub pull mirror and verify Git/LFS/wiki plus independent restore before expanding.
- [ ] Inventory/export GitHub PR/issues/releases and required GHCR images separately; automatic discovery and metadata backup are not implemented.

GitHub stays primary. No origin, Actions, Argo source URL or registry cutover.
Source preparation does not claim live deployment, mirror enrollment or recovery acceptance.
See [Forgejo runbook](platform/forgejo/README.md).
