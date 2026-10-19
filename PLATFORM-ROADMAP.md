# OneMind platform roadmap: layers, OpenChoreo and AI operations

Living document. Started 2026-10-18. Owner: Zeus. Kept at the repository root so any
agent (Claude Code, Codex or another) can pick up the work where the last one stopped.

This is a plan, not a record of live state. Live state is in [`docs/STACK.md`](docs/STACK.md);
the mission it serves is in [`docs/MISSION.md`](docs/MISSION.md).

> Taking over from another agent? Start with [`docs/roadmap/HANDOFF.md`](docs/roadmap/HANDOFF.md)
> (state, rules, prompt) and the phase briefs in `docs/roadmap/`.

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
| 1 | Claude sees every layer (MCP on the Mac) | Done | [#257](https://github.com/onemindos/omos/pull/257), [#258](https://github.com/onemindos/omos/pull/258) | 2 |
| 2 | OpenChoreo foundations: data plane, gateway, secrets, first template, built-in console | Done | 2a [#262](https://github.com/onemindos/omos/pull/262), 2b [#263](https://github.com/onemindos/omos/pull/263), 2c [#264](https://github.com/onemindos/omos/pull/264), 2d [#265](https://github.com/onemindos/omos/pull/265), 2e [#266](https://github.com/onemindos/omos/pull/266), 2g-1 [#268](https://github.com/onemindos/omos/pull/268), [#269](https://github.com/onemindos/omos/pull/269), 2g-2 [#270](https://github.com/onemindos/omos/pull/270) [#271](https://github.com/onemindos/omos/pull/271) [#272](https://github.com/onemindos/omos/pull/272) [#273](https://github.com/onemindos/omos/pull/273), 2g-3 [#274](https://github.com/onemindos/omos/pull/274), 2f-1 [#275](https://github.com/onemindos/omos/pull/275), 2f-2 [#276](https://github.com/onemindos/omos/pull/276), (this PR) | 3, 4 |
| 3 | Pilot one app as an OpenChoreo Component | In progress | [#279](https://github.com/onemindos/omos/pull/279), [#286](https://github.com/onemindos/omos/pull/286), [#287](https://github.com/onemindos/omos/pull/287), [#288](https://github.com/onemindos/omos/pull/288), [#289](https://github.com/onemindos/omos/pull/289), [#290](https://github.com/onemindos/omos/pull/290), [#291](https://github.com/onemindos/omos/pull/291) | 4 |
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
| 2. `openchoreo-ca` ClusterIssuer | self-signed CA | not created yet | Gateways start with TLS off (the guide's first install also does). Later TLS uses the existing `letsencrypt-dns` issuer |

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
      `*.apps.onemindos.dev` with `letsencrypt-dns` TLS replaces it later. Verified
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
- The Authentik providers made by `auth-provider-bootstrap.yaml` for `workflows`, `langfuse`
  and `registry` were created with an empty grant type list, like `openchoreo` was, and
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
- Evaluation is still open: the pilot's three files total 438 lines, including the
  generated immutable type/workload snapshot. The old browser occupies 47 lines in
  one shared manifest. Do not claim fewer files or easier maintenance without the
  live comparison and a verdict.

- [x] Pick a small stateless Layer 3 app (candidates: `stac-browser`, `agent-console`).
- [x] Write its Component in `platform/openchoreo/projects/`.
      The parallel deployment PR adds a pinned STAC Browser Component, Workload,
      immutable ComponentRelease and development ReleaseBinding. Runtime verification
      and cutover are pending; the existing app has not been removed.
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
| 2026-10-18 | Cilium Gateway API backends need a component-scoped `fromEntities: ingress` policy; keep OpenChoreo's generated policy and add an opt-in trait rather than allowing all ingress. |
| 2026-10-18 | AI read access to the cluster uses ServiceAccount `kube/ai-operator` (built-in `view` plus platform CRDs, no Secrets) through kubernetes-mcp-server in read-only mode. No separate Argo CD MCP. |
| 2026-10-18 | Not adopting kagent, KubeVela, Kratix or Humanitec. HolmesGPT is optional after Phase 5. |

## Incidents found while working the roadmap

| Date | What | Fix |
|---|---|---|
| 2026-10-18 | DOKS replaced the hub-core nodes (v1.36.5). `tak-ldap-0` and `takml-server` stayed Pending, pinned to the removed node `onemind-hub-core-v2-37phkg` (the `takml-server` pin was set live, not in Git). `nats-2` could not schedule: the only node without a NATS pod was full while a new node sat nearly empty. Found by the phase 1 read-only access. | [#259](https://github.com/onemindos/omos/pull/259): node-pool affinity for `tak-ldap`, pod affinity to `takml-vision-worker` for `takml-server`, stale live pin removed. The Hermes pod was rescheduled to make room for `nats-2` (NATS back to 3/3). [#260](https://github.com/onemindos/omos/pull/260): `Recreate` for `takml-db`, whose rolling update deadlocked on its ReadWriteOnce volume. |

Follow-ups: nothing rebalances pods after DOKS replaces nodes. Consider the Kubernetes
descheduler (Layer 1) in phase 5 alongside the capacity work.

## Risks and open questions

- **Maturity:** OpenChoreo is a CNCF Sandbox project with frequent releases and several
  supported lines at once. Pin chart versions and read release notes before upgrades.
- **Backstage version:** the OpenChoreo plugins pin Backstage 1.51.0; our portal runs 1.55.3.
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
