# OneMind Kubernetes operating structure and execution queue

**Proposal / implementation backlog, 2026-10-20. Not a deployment authorization.**
Mission: a home-centered, sovereign community platform with cloud/edge reach,
TAK/ML operations, useful human interfaces and authenticated AI workers. Keep
Kubernetes's native capabilities available; make operating them coherent, not
invisible. AI doing the work does not remove failure domains, reconciliation
conflicts, security boundaries, recovery requirements or compute costs.

## Start here

1. **Continue the remaining Backstage trust gates**, using merged
   [#331](https://github.com/onemindos/omos/pull/331). Source #330, pipeline #335
   and image pin #336 are already merged and their scoped rollout/UI/security
   acceptance is recorded. Do not repeat that release. Live expiry/second-role,
   recovery, launcher and adapter gaps remain separate; STAC repair is not a
   prerequisite for every unrelated app or for this planning PR.
2. **Plan and execute foundation repairs separately**: retired community
   ingress-nginx, GitOps ownership gaps, independently recoverable backups,
   identity/network/policy coverage. Do not hide these behind healthy pod badges.
3. Standardize app onboarding and real operations cards using the installed
   stack; prove one noncritical app end to end before mass migration.
4. Establish the Mac/home failure domain, then test NATS agents and selected
   multi-cluster/edge/ML capabilities. Do not start by installing every CNCF tool.

Current evidence is in the
[timestamped baseline](../operations/kubernetes-structure-baseline-2026-10-20.md).
The root [phase tracker](../../PLATFORM-ROADMAP.md) remains authoritative:
**Phase 6 first unfinished; Phase 7 not started**. This queue supplements it and
does not silently replace the four-layer trust model.

## The structure to build

```text
Human surfaces: Backstage workspace / TAK operations / LibreChat conversation
                         │ authenticated APIs and projections
AI workers: Legacy + specialist harnesses ── NATS discovery/events/request-reply
                         │ scoped tools, approval, evidence and budget controls
Desired state: Git PRs → Argo CD → approved app CRs / native resources / operators
                         │
Sites: DOKS cloud hub ── private networking ── Mac Linux VM / home / field sites
                         │
Local domain systems: Home Assistant / OpenRemote / TAK / apps / device gateways
```

This is a **target**, not a claim that all surfaces, agents or sites are integrated.
Backstage should show every service regardless of deployment owner. Moving an app
under OpenChoreo is **not required just to make it visible**. A catalog identity,
Kubernetes selector/cluster mapping and tested adapters are the integration work.
Backstage's Kubernetes plugin provides service-oriented multi-cluster visibility;
it is not itself the cluster lifecycle manager.
[Backstage Kubernetes](https://backstage.io/docs/features/kubernetes/),
[cluster configuration](https://backstage.io/docs/features/kubernetes/configuration/)

### Ownership: one reconciler per resource/field

| Layer / object | Desired owner | Boundary |
|---|---|---|
| Foundation: identity, private access, recovery, policies, platform controllers | Argo CD + the installed owning controller | AI observes/proposes PRs; human approves foundation changes. Foundation must recover without the application abstraction. |
| Shared data, messaging, registry and telemetry | Argo CD + domain operators | Catalog and APIs expose them without transferring ownership to an app platform. |
| Product apps | Existing Argo resources; OpenChoreo for applications that pass its contract | Argo owns the OpenChoreo intent CR; OpenChoreo owns its generated resources. Never have Argo also apply those generated Deployments. Native CRDs remain possible. |
| External infrastructure | Terraform today; a bounded Crossplane experiment only if approved | Never let Terraform and Crossplane control the same external resource concurrently. |
| Independent clusters | Per-cluster Argo registration/ApplicationSets after acceptance | Not a stretched WAN control plane. Future local recovery cannot depend solely on cloud-hosted GitOps/identity/inference. |
| Selected federated workloads | Karmada only after a dedicated pilot | Argo owns intent in the Karmada control plane, not duplicate generated member-cluster workloads. |
| Edge node/device lifecycle | KubeEdge pilot where disconnected edge nodes need it | Keep device safety and local authority separate from an LLM; normal site K3s clusters do not automatically need EdgeCore. |
| Live domain data and commands | The app/device system itself | NATS connects contracts; Backstage/AI projections do not replace application databases or physical interlocks. |

## Capability coverage: what we already use versus what remains

“Observed” below means matching controllers/workloads were present in the
2026-10-20 baseline, **not** that every app uses them correctly. “Candidate” means
uninstalled/unaccepted here. Recheck versions/compatibility when implementing.

| # | Capability | Use / observed foundation | Next gate or conditional addition |
|---|---|---|---|
| 1 | Streaming / messaging | NATS, JetStream, NACK; 3/3 broker replicas; existing Synadia/Hermes agent integration | Authenticated discovery and actual request/stream/approval tests; replay, ACLs and bounded retention. No Kafka by default. |
| 2 | CI / delivery | GitHub Actions configuration, Argo CD/Rollouts; verified operator-run local/Linux pipeline #335 | Portal publish/promotion/rollout path proven in #336; hosted jobs still unavailable, not green. Other apps need reviewed profiles and semantic gates. No automatic removal executor; workflow-plane plan stays in #328. |
| 3 | App definition / builds | Helm, Kustomize, existing OpenChoreo app/data plane | Standard onboarding contract and reproducible builds. KubeVela comparison is a separate lab, not a second owner for live apps. |
| 4 | Automation / configuration | Argo reconciliation; existing n8n; native Jobs/CronJobs | Durable run IDs, idempotent operations, retry limits and approval records. Cron alone is not durable orchestration. |
| 5 | Observability | Prometheus/Alertmanager, Grafana, OTel, Loki, Tempo, OpenSearch and OpenChoreo Observer adapters | Verify each signal source, permissions, timestamps, semantic health and alert delivery; fix Grafana adapter access before declaring unified telemetry. |
| 6 | Scheduling / orchestration | Native Kubernetes scheduler, KEDA and Argo Rollouts | Pod requests/limits, disruption/spread, priority and safe scaling. Evaluate Kueue for a real queued ML workload, not as a replacement scheduler. |
| 7 | Continuous optimization | OpenCost, metrics and autoscaling APIs | Allocation versus actual agent/model spend; rightsize from evidence, protect reserves, require approval for spend-increasing changes. |
| 8 | Security / compliance | Authentik, RBAC, Kyverno, Falco, SPIRE and OpenFGA controllers | Prove coverage/enforcement, admission denial, least privilege and audit/revocation. Installed engines do not prove compliant apps. |
| 9 | Storage / recovery | DOKS CSI, CNPG/Barman, Velero and VolSync | Application-consistent restore drills and independent copies; a home storage decision before Nextcloud. Do not install Ceph and Longhorn just for coverage. |
| 10 | Chaos engineering | No accepted chaos controller in this baseline | One bounded lab fault plan first (pod loss, WAN loss, disk pressure); evaluate one of Chaos Mesh/Litmus only when rollback and exclusions exist. |
| 11 | App / AI framework | Existing Hermes/Legacy, Synadia agent protocol, Dapr controllers | Select harness by task; retain permission/tool boundaries. Dapr is optional per app, not a compulsory rewrite of the NATS fabric. |
| 12 | Container runtime | DOKS containerd 2.2.3; Docker Desktop currently an image builder on Mac | K3s inside a dedicated ARM64 Linux VM for the proposed home pilot. Podman is a runtime/tool, not a substitute multi-site Kubernetes control plane. |
| 13 | API gateway | Existing Cilium Gateway API; legacy public ingress-nginx; private Tailscale ingress | High-priority public ingress replacement assessment; test actual installed Cilium/Gateway API capabilities and every annotation, not latest-doc assumptions. |
| 14 | RPC / APIs | App HTTP APIs, NATS request/reply and gRPC where a service needs it | Versioned contracts, deadlines, identity and trace propagation. NATS does not automatically translate every third-party API. |
| 15 | Databases | PostgreSQL/CNPG, ClickHouse, OpenSearch; app-owned database instances | Declare the authoritative store and owning operator per app. MongoDB needed for the LibreChat pilot is not provisioned by its staged PR. |
| 16 | Keys / identity | SealedSecrets, ESO, cert-manager and SPIRE controllers present, with distinct roles | SealedSecrets remains the repository credential pattern. ESO presence is not approval to switch stores or claim populated external vaults. Backup sealing keys, prove rotation/local recovery; no catalog tokens. Evaluate OpenBao/KMS only for a demonstrated gap. |
| 17 | Service mesh | Cilium networking/policy and existing SPIRE workload identity | Check concrete mTLS/L7 policy needs. Do not claim universal encryption or install Istio and Linkerd together without an ownership/traffic design. |
| 18 | Registry / supply chain | Harbor (including Trivy workload) and GHCR | Pinned digests, SBOM/provenance/signing verification, vulnerability exceptions and disconnected image availability. Harbor running does not prove admission verification. |
| 19 | Service proxy | Existing Cilium gateway/proxy path and legacy ingress proxy | Reuse the selected gateway; preserve WebSocket/SSE, source restrictions, timeouts and approved public/private exposure. |
| 20 | Discovery / coordination | CoreDNS; managed Kubernetes control-plane state; NATS microservice discovery | Stable app/site/cluster IDs, observed freshness and reconnect reconciliation. Do not add a separate etcd cluster to “complete the category.” |

### What the larger platforms add — and what they do not

- **OpenChoreo:** existing app lifecycle/environment API and Backstage plugins.
  Keep the deployed pilot while testing suitability; do not force databases,
  foundation controllers or TAK into it merely to obtain a UI card.
- **KubeVela:** a legitimate alternative application abstraction with extensible
  definitions, traits and workflows. Its addon ecosystem includes overlapping
  infrastructure. Compare a pinned version on disposable workloads against the
  actual OpenChoreo/native GitOps contract. An addon list is not proof of a
  maintained native integration with our pinned Backstage version. Do not enable
  its Flux addon over Argo-owned resources.
  [KubeVela addon overview](https://kubevela.io/docs/reference/addons/overview/)
- **Crossplane:** custom infrastructure APIs/compositions and external providers;
  useful if reusable database/bucket/environment requests justify it. It can also
  compose Kubernetes resources, so explicitly constrain the experiment's ownership
  rather than assuming it never overlaps an app platform.
  [Composition](https://docs.crossplane.io/latest/composition/),
  [providers](https://docs.crossplane.io/latest/packages/providers/)
- **Karmada:** selected multi-cluster placement/override policies. It does not
  by itself replicate application data, validate cloud independence or make a
  stateful service safely fail over between sites.
  [Karmada concepts](https://karmada.io/docs/core-concepts/concepts/)
- **KubeEdge:** edge-node/device integration and edge autonomy mechanisms, subject
  to actual device/runtime/network tests. It is not needed to make every regular
  home/region cluster visible in Backstage.
  [KubeEdge project](https://github.com/kubeedge/kubeedge)
- **Rancher / Fleet / Harvester / Elemental:** separate cluster UI, delivery,
  virtualization/HCI and machine-OS roles. Reassess when physical hardware and
  fleet provisioning justify them. No Fleet/Argo overlap, no assumption that
  Harvester is the production virtualization answer for an Apple Silicon laptop.
  Existing Headlamp remains the Kubernetes drill-down/break-glass UI.

These are **bounded candidate evaluations**, not rejection because “one person
cannot handle complexity” and not blanket approvals to deploy them all.

## Integration reality: ecosystem support versus OneMind acceptance

| Connection | Evidence / intended contract | State here |
|---|---|---|
| Backstage ↔ OpenChoreo | Installed pinned core/observability plugins, existing identity policy; native entities and Deploy/logs/metrics plus custom joined Operations | Scoped pilot UI and read-only boundaries verified; not universal app integration or machine action access |
| Git → Argo → portal | Reviewed source, Linux artifact receipt, separate digest pin, exact revision/runtime-digest/replica verification | Actually exercised in #335/#336; operator-run, not a new scheduler or unrestricted deployment service |
| Backstage ↔ Kubernetes / Argo | Entity selectors, cluster identity and service-oriented drill-down; owning Argo Application linked | Existing single cloud cluster; prove each new app/cluster mapping rather than assume discovery is full observability |
| NATS ↔ Legacy / other agents | Existing bus/agent workloads and protocol; scoped adapters, request/response/stream/approval/cancel correlation needed | Inventory/discovery is not working inference; earlier provider 402 remains unresolved. App APIs need adapters unless they already speak the contract |
| Crossplane ↔ Argo | Upstream documents annotation tracking and resource-health customization | Documented integration, not installed/tested here; evaluate tracking, health, deletion/orphan and provider scopes in K08 before any adoption |
| KubeVela ↔ current platform | Addons and extensible definitions/workflows exist upstream | Candidate comparison only; no accepted Backstage 1.55.3 plugin or coexistence proof here |
| Karmada / KubeEdge ↔ home/sites | Placement/overrides versus edge-node/device lifecycle are different jobs | Not deployed; independent cluster recovery, ownership, WAN-loss and device tests precede a pilot |
| TAK / Home Assistant / OpenRemote ↔ COP and agents | Native domain UI/API remains authoritative; scoped telemetry and permitted commands projected into the workspace | Per-service/device integration and safety proof still needed; no universal “NATS connects everything” assumption |

Primary references: [Backstage Kubernetes](https://backstage.io/docs/features/kubernetes/),
[OpenChoreo host installation](https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/installing-into-existing-backstage/),
[Crossplane with Argo](https://docs.crossplane.io/latest/guides/crossplane-with-argo-cd/),
[Karmada placement/overrides](https://karmada.io/docs/core-concepts/concepts/).
Upstream documentation is not a version-compatibility test or permission to change
our Argo tracking settings. Inspect actual installed schemas/versions before implementation.

### Working decisions and explicit escape hatches

- Keep Backstage as the extensible service/developer front end; retain TAK and
  application-native surfaces for domain work. No promise that a catalog is the
  sole daily operations UI. Keep plugins in this monorepo unless independent
  ownership/releases justify splitting them.
- Git/Argo remain desired-state authority. Use the existing OpenChoreo contract
  where suitable; native Kubernetes/Helm/operator resources remain first-class
  choices, not second-class apps hidden from the portal.
- Terraform remains current external-infrastructure ownership. Crossplane and
  KubeVela are comparison candidates, not prerequisites or co-owners. Do not adopt
  Flux/Fleet over the same Argo resources merely to obtain another integration.
- Mac Linux VM/K3s is a proposed local pilot, not an installed production cluster
  or a final dedicated-host OS decision. Persistent storage, ARM64 images, local
  bootstrap/identity/registry and outage recovery must be proven independently.
- Agent reads use existing scoped identities. Writes need backend authorization,
  bounded action contracts, approval/evidence and rollback/recovery policy. NATS
  or an MCP endpoint does not create authority; cron is not reconciliation.

### One app lifecycle, regardless of deployment owner

1. **Register:** stable identity, keeper decision, source, owner, dependencies,
   site/cluster/environment, human access and scoped AI API/tool contract.
2. **Validate/build:** pinned inputs, supported architecture, policy/schema and
   secret checks, app tests; immutable source-linked artifact and retained evidence.
3. **Promote:** profile selects exact owner/path/target; normal reviewed Git PR,
   then Argo reconciliation. Operator-generated resources are not separately applied.
4. **Observe:** verify revision, actual runtime image/readiness and real app path;
   check catalog edges, authorization, alerts and data freshness. Unknown stays unknown.
5. **Recover/retire:** tested restore and dependency/retention review first. Changes
   to prune/finalizers/storage/identity require their explicit approval. The current
   pipeline provides retirement inventory only, not deletion or data rollback.

Next K03 deliverable is a reviewed onboarding record/profile adapter for one
selected existing stateless keeper, reusing the installed templates. Do not choose
an app, change controller ownership or start a migration merely to test scaffolding.

## The Backstage contract for every app

Each app should provide one coherent operational page with:

1. Stable catalog identity, owner, role, site, environment and cluster mapping;
   parent system, dependencies, source, runbook and lifecycle stage.
2. Live observed deployment owner/version/digest, requested versus ready replicas,
   generation, recent events and drift. Inventory-only entries are labelled.
3. Useful app-level checks/data freshness alongside pod health; explicitly show
   unknown, stale, unauthorized and unavailable rather than a false green badge.
4. Scoped logs, metrics, traces, alerts and recovery links that work for the
   authenticated human and the appropriately scoped AI.
5. Separate actions for proposed/approved/executed operations; permission checks
   on the backend. A visible button or mesh login is not authorization.
6. The app's real data/tool APIs, NATS bridge where useful, export path, encryption,
   backup owner, last verified restore and rollback/cutover instructions.

The scoped STAC operational-page proof is recorded in #330/#335/#336 and the
[release acceptance](../operations/stac-operational-slice-2026-10-20.md). The actual
API/schema failure remains visible; operational-page acceptance does not mean
STAC data readiness or complete coverage for other apps.
Keep custom plugins in the existing monorepo packages. A repository per plugin
is unnecessary unless independent ownership/release requirements emerge.

## Delivery queue — small PRs with evidence, not one giant migration

| Work package | PR / work scope | Acceptance and dependency |
|---|---|---|
| K00 — trusted starting point | Merged #331/#330/#335/#336; scoped portal release complete, broader trust/recovery gates open | Preserve live catalog/denial evidence; finish remaining live token lifecycle/role and affected-app gates. No Phase 6 “Done” from this structural plan. |
| K01 — public ingress lifecycle | Existing [#334](https://github.com/onemindos/omos/pull/334): read-only assessment first, then separately approved route/capability pilots | Official ingress-nginx retirement recorded; route/annotation compatibility, private/public boundaries, DNS/TLS and rollback proven before cutover. Never migrate Authentik first. |
| K02 — ownership and recovery | NATS StatefulSet now has Argo tracking; audit remaining release/field ownership before assuming a conflict. Inspect TAK core/Argo bootstrap, Matrix drift and OpenClaw; one fix PR each | No blind re-adoption or Helm upgrade of Argo-owned NATS. Inventory field owners/release remnants without fetching Secrets; backup before stateful changes; isolated restore and independent copy. |
| K03 — app onboarding | Extend/reconcile existing templates and the portal-only pipeline profile; do not create a competing scaffolder | Pinned artifact, selector, access, probes, policy, resources, signal links, data/API, AI tools and recovery gates. Works for native apps as well as approved OpenChoreo apps. |
| K04 — agent conversation | [#332](https://github.com/onemindos/omos/pull/332) staged LibreChat → actual Legacy | Inference 402 resolved by approved means; identity/credentials/database/backups and two-conversation/approval tests. No new privilege. |
| K05 — home foundation | Existing [#284](https://github.com/onemindos/omos/pull/284), revised against Mac evidence before implementation | Storage capacity, ARM64 VM/runtime, 24/7 closed-lid/reboot/FileVault operation, private access, GitOps, local identity/recovery and WAN-loss tests. No real file move now. |
| K06 — Nextcloud / knowledge | Per-app PR after K05; retain Obsidian vault and existing Home Assistant | Dedicated data + DB backup/restore, lock/cache/cron/SSO, export, local access and disconnected client/conflict tests. Sync is not backup. Do not put the live server data directory inside a desktop sync folder. |
| K07 — agent fabric + spend | One authenticated worker/adaptor at a time, then metering projection to Agent Console/Backstage; [permissioned-action plan](../roadmap/permissioned-ai-actions.md) stages a GitOps-first supervised pilot without granting writes | Discover, prompt, stream, approval, reconnect, stop semantics and audit correlation. Provider-reported / estimated / unknown spend remain distinct, including delegation. |
| K08 — abstraction comparison | Disposable KubeVela/native/OpenChoreo and bounded Crossplane experiments if approved | Same app/API/secret/policy/recovery tests; actual plugin compatibility, CRD escape hatch and controller ownership. No production replacement on an unexecuted recommendation. |
| K09 — fleet / edge | Two independently recoverable clusters, then a Karmada placement pilot and separate KubeEdge node/device pilot | Placement and storage overrides, disconnection/reboot/rejoin, revocation, no duplicate reconciliation. No stretched cross-region etcd. |
| K10 — ML compute | One real TAK-ML job/serving benchmark | GPU/ARM64 compatibility, priority/reserves/cost, bounded queues and model/artifact provenance. Consider Kueue for admission, KServe for a demonstrated serving need; full Kubeflow is not a prerequisite. |
| K11 — local physical operations | OpenRemote/Home Assistant/TAK adapters with permitted test devices | Device identity, fresh telemetry, explicit command scope, human approval, local safety/interlocks and offline fallback. Transport latency is measured, not promised. |

Implementation priority is **remaining K00 trust gates and K01 assessment**, with
K02 read-only recovery/ownership checks in parallel. The user selected #333 for
structural review; that does not require another portal image or unblock Phase 7. Later packages need their own stated dependencies,
not completion of all twenty categories. User-requested preparation can proceed
in parallel; activating it still requires its acceptance gates.

## Repository and execution rules

- Keep `omos` as the GitOps/source monorepo. Current `platform/<service>/`,
  `platform/argocd/`, `services/` and plugin package paths stay intact now.
- Per-cluster/site overlays can be added by a later layout PR after reviewing
  Argo paths, secrets scope and catalog ownership. Do not move live manifests to
  an unreferenced “organized” folder or combine every application into one file.
- For each implementation: read mission/tracker → collect read-only evidence →
  scoped branch/PR → validation/review → merge when authorized → Argo reconcile →
  live gate → record evidence. No ad hoc unowned Helm releases.
- External provider resources, identity, public exposure, deletion and equipment
  actuation need their specific approvals. AI cron jobs must not become a second
  owner fighting the Kubernetes controllers.
- Measure local bus, LAN, WAN, tool/API and inference latency separately (p50/p95/
  p99). No “microseconds everywhere” guarantee or cloud-free claim while essential
  identity/inference/tools/data remain remote.

### Additional primary references

- [Ingress NGINX retirement](https://kubernetes.io/blog/2025/11/11/ingress-nginx-retirement/)
- [Cilium Gateway API](https://docs.cilium.io/en/stable/network/servicemesh/gateway-api/gateway-api/)
- [Argo cluster generator](https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/Generators-Cluster/)
- [K3s platform requirements](https://docs.k3s.io/installation/requirements)
- [Kueue role and integrations](https://kueue.sigs.k8s.io/docs/overview/)

Latest upstream docs are references, not permission to upgrade DOKS-managed Cilium
or Gateway API CRDs. Pin implementation versions only after compatibility tests.
