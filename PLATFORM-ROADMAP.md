# OneMind platform roadmap: layers, OpenChoreo and AI operations

Living document. Started 2026-10-18. Owner: Zeus. Kept at the repository root so any
agent (Claude Code, Codex or another) can pick up the work where the last one stopped.

This is a plan, not a record of live state. Live state is in [`docs/STACK.md`](docs/STACK.md);
the mission it serves is in [`docs/MISSION.md`](docs/MISSION.md).

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

### Namespace map (live cluster, 2026-10-18)

| Namespace | Main workloads | Layer | Target |
|---|---|---|---|
| `kube-system`, `kube-node-lease`, `kube-public`, `default` | DOKS system, Cilium | 0 | DigitalOcean |
| `argocd` | Argo CD | 1 | Argo CD (Helm CLI today, see `CLAUDE.md`) |
| `auth` | Authentik | 1 | Argo CD |
| `tailscale` | Tailscale operator and proxies | 1 | Argo CD |
| `cert-manager` | cert-manager | 1 | Argo CD |
| `ingress-nginx` | public ingress | 1 | Argo CD |
| `cloudflare` | cloudflared (cold standby) | 1 | Argo CD |
| `external-dns` | external-dns | 1 | Argo CD |
| `kyverno`, `falco`, `spire-system` | policy, runtime security, workload identity | 1 | Argo CD |
| `velero`, `ops` | Velero, VolSync | 1 | Argo CD |
| `github-runner` | CI runners | 1 | Argo CD |
| `kube` | Headlamp, Kubernetes Dashboard (break-glass UIs) | 1 | Argo CD |
| `backstage` | OneMind portal (the human home screen) | 1 | Argo CD |
| `openchoreo-control-plane`, `openchoreo-data-plane`, `openchoreo-workflow-plane`, `openchoreo-observability-plane` | OpenChoreo itself | 1 | Argo CD |
| `takserver` | TAK Server core | 1 for now | Argo CD; Helm CLI release today. Revisit after Phase 7 |
| `nats` | 3-node JetStream bus | 2 | Argo CD; streams as an OpenChoreo ResourceType |
| `cnpg-system`, `onemind-data` | CloudNativePG operator and the `onemind-data` cluster | 2 | Argo CD; databases as an OpenChoreo ResourceType |
| `clickhouse` | ClickHouse | 2 | Argo CD; ResourceType later |
| `onemind-search` | OpenSearch | 2 | Argo CD; also the observability plane's log store |
| `onemind-ai` | LiteLLM, CPU model | 2 | Argo CD. LiteLLM stays the model gateway |
| `onemind-registry` | Harbor | 2 | Argo CD |
| `dapr-system` | Dapr | 2 | Argo CD |
| `keda` | KEDA | 2 | Argo CD; OpenChoreo autoscaling module |
| `monitoring` | Prometheus, Alertmanager, Grafana, Loki, Tempo, OTel collector | 2 | Argo CD; becomes the observability plane's metrics and events modules |
| `opencost` | OpenCost | 2 | Argo CD; OpenChoreo FinOps module |
| `omos` | feed-intel, feed-satellites, omos-tak-projector, redpanda-connect, sink-clickhouse | 3 | OpenChoreo Project `omos-feeds` |
| `onetak` | 19 geo and mapping apps (Valhalla, Martin, Nominatim, WebODM, STAC, TileServer and others) | 3 | OpenChoreo Project `onetak` |
| `takml` | takml-server, takml-vision-worker, takml-db | 3 | OpenChoreo Project `takml` |
| `onemind-core` | onemind-core, onemind-ingestion, OpenFGA | 3 | OpenChoreo Project `onemind` |
| `onemind-jobs`, `onemind-media` | empty today | 3 | OpenChoreo Project `onemind` |
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
| 1 | Claude sees every layer (MCP on the Mac) | In progress | [#257](https://github.com/onemindos/omos/pull/257) | 2 |
| 2 | OpenChoreo foundations: data plane, gateway, secrets, first template, built-in console | Not started | | 1, 4 |
| 3 | Pilot one app as an OpenChoreo Component | Not started | | 4 |
| 4 | Portal refactor: custom features into `plugins/` packages | Not started | | 2, 3 |
| 5 | Capacity fixes, then the observability plane | Not started | | none |
| 6 | One portal: OpenChoreo plugins in our Backstage, built-in console off | Not started | | 7 |
| 7 | Migrate Layer 3 apps project by project; NATS and Postgres as ResourceTypes | Not started | | 6 |
| 8 | Expand: home server data plane, Hermes on the same tools | Not started | | none |

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
- [ ] After merge and Argo sync: run the script and pass the gate below.
- Moved: OpenChoreo control-plane MCP goes to phase 2. The endpoint works
  (`/mcp` answers 401 and advertises Authentik), but OpenChoreo manages nothing yet, and
  it needs a read-only OpenChoreo role and an Authentik client for the agent.
- Dropped: a separate Argo CD MCP. Argo CD is still a Helm CLI release (foundation), and
  Application sync and health are readable through `onemind-kubernetes`.
- Deferred: the tool gateway (LiteLLM MCP or agentgateway) is only needed for in-cluster
  agents, so it moves to phase 8 with Hermes.

Gate: from a fresh session, Claude lists unhealthy workloads in each layer with the cause,
using only the MCP servers, and every server refuses a write.

### Phase 2: OpenChoreo foundations

Goal: OpenChoreo can deploy something. All resources in Git, applied by Argo CD.

- [ ] New tree `platform/openchoreo/` with its Argo Application(s):
      `planes/` (charts and registration), `platform/` (environments, pipeline,
      component types, traits), `projects/` (one folder per project).
- [ ] Data plane chart `openchoreo-data-plane` 1.3.1, scheduled on the ecosystem pool.
- [ ] Gateway: try Cilium Gateway API first; fall back to kgateway if the data plane
      needs it. Record the result.
- [ ] Secret store: External Secrets controller plus a `ClusterSecretStore` named `default`.
      Choose the backend (Kubernetes provider or OpenBao) and record it. SealedSecrets
      stays the way secrets enter Git.
- [ ] Register `ClusterDataPlane` `default` (agent CA, secret store, gateway).
- [ ] Environment `prod` and a `DeploymentPipeline`.
- [ ] First `ClusterComponentType` `onemind-service`: resource requests and limits,
      network policy, Tailscale route, standard labels.
- [ ] Turn on the built-in console (`backstage.enabled: true`), tailnet-only, with
      our catalog files mounted through `backstage.extraVolumes` and `appConfig`.

Gate: the data plane agent is connected, the plane shows healthy in the OpenChoreo API,
and the built-in console signs in through Authentik.

### Phase 3: Pilot one app

Goal: prove OpenChoreo saves work before moving anything important.

- [ ] Pick a small stateless Layer 3 app (candidates: `stac-browser`, `agent-console`).
- [ ] Write its Component in `platform/openchoreo/projects/`.
- [ ] Run it next to the existing deployment, switch traffic, then delete the old manifests.

Gate: the app works, it needs fewer files than before, and the change was easier to make.
If not, stop and record why before Phase 7.

### Phase 4: Portal refactor

Goal: our portal features become portable Backstage plugin packages. No behaviour change.

- [ ] Move each frontend feature from `services/onemind-portal/packages/app/src/modules/`
      into its own package under `services/onemind-portal/plugins/`.
- [ ] Move each backend feature from `packages/backend/src/` into its own backend plugin
      or module package.
- [ ] Tests, TypeScript check, image build and a browser check pass; new image digest pinned.

Gate: the deployed portal behaves exactly as before.

### Phase 5: Capacity, then observability

- [ ] Apply the request fixes from `docs/operations/cluster-review-2026-10-09.md` (H1), or
      add an ecosystem node. Target: no node above about 80% of memory requested.
- [ ] Observability plane 1.3.1 with modules wired to the existing Prometheus, OpenSearch
      and OTel collector; OpenCost as the FinOps module.
- [ ] Register `ClusterObservabilityPlane` and link the data plane to it.
- [ ] Connect Claude to the observer MCP (logs, metrics, traces, alerts, incidents, costs).

Gate: an OpenChoreo component shows logs and metrics in the console and through the
observer MCP.

### Phase 6: One portal

- [ ] Align our Backstage with the release line the OpenChoreo plugins support
      (tested on 1.51.0 today; we run 1.55.3), or confirm 1.55 works.
- [ ] Install the core OpenChoreo plugins and the observability pack into our portal,
      using `openchoreo/backstage-plugins` `packages/portal-app` as the reference wiring.
- [ ] Switch sign-in to the `openchoreo-auth` provider against Authentik and keep the
      `onemind.io/oidc-sub` user mapping working.
- [ ] Merge the OpenChoreo permission policy with ours (only one policy can be active).
- [ ] Turn the built-in console off.

Gate: one portal shows our home screens, the full catalog and OpenChoreo tabs.

### Phase 7: Migrate Layer 3

Order: `omos-feeds`, `onetak`, `takml`, `onemind`, `agent-fabric`, `n8n`, `matrix`,
`openclaw`, then `hermes` last.

- [ ] One PR per project. Old manifests are deleted only after the new release is healthy.
- [ ] ResourceTypes for NATS streams and CNPG databases so components declare their
      dependencies.
- [ ] Bring `matrix` back in sync and put the OpenClaw Deployment in Git as part of their moves.

Gate per project: healthy in OpenChoreo, visible in the portal, old Argo resources removed.

### Phase 8: Expand

- [ ] Home server on Talos Linux, registered in Argo CD and as a second OpenChoreo data plane.
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
| 2026-10-18 | AI read access to the cluster uses ServiceAccount `kube/ai-operator` (built-in `view` plus platform CRDs, no Secrets) through kubernetes-mcp-server in read-only mode. No separate Argo CD MCP. |
| 2026-10-18 | Not adopting kagent, KubeVela, Kratix or Humanitec. HolmesGPT is optional after Phase 5. |

## Risks and open questions

- **Maturity:** OpenChoreo is a CNCF Sandbox project with frequent releases and several
  supported lines at once. Pin chart versions and read release notes before upgrades.
- **Backstage version:** the OpenChoreo plugins pin Backstage 1.51.0; our portal runs 1.55.3.
- **Gateway:** Cilium Gateway API with the OpenChoreo data plane is untested.
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
