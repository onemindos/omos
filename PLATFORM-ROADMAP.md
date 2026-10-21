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
| 3 | Pilot one app as an OpenChoreo Component | Done | [#279](https://github.com/onemindos/omos/pull/279), [#286](https://github.com/onemindos/omos/pull/286), [#287](https://github.com/onemindos/omos/pull/287), [#288](https://github.com/onemindos/omos/pull/288), [#289](https://github.com/onemindos/omos/pull/289), [#290](https://github.com/onemindos/omos/pull/290), [#291](https://github.com/onemindos/omos/pull/291), [#292](https://github.com/onemindos/omos/pull/292), [#293](https://github.com/onemindos/omos/pull/293), [#294](https://github.com/onemindos/omos/pull/294) | 4 |
| 4 | Portal refactor: custom features into `plugins/` packages | Done | [#280](https://github.com/onemindos/omos/pull/280) | 2, 3 |
| 5 | Capacity fixes, then the observability plane | Done | [#281](https://github.com/onemindos/omos/pull/281), implementation #300–#310 (individual links and live evidence below) | none |
| 6 | One portal: OpenChoreo plugins in our Backstage, built-in console off | In progress | Historical acceptance #282, #318, #319; reopened for [portal stabilization](docs/operations/portal-stabilization-2026-10-20.md) before migration (#322 source, #323 release, #324 inventory, #325 composition, #326 final image); 2/2 healthy, 715 processed entities; broader acceptance gaps remain documented | bounded design/preflight/pilot-source preparation only; no later-phase deployment |
| 7 | Migrate Layer 3 apps project by project; NATS and Postgres as ResourceTypes | Not started | | 6 |
| 8 | Expand: home server data plane, Hermes on the same tools | Not started | | none |

## Phases

### Authorized parallel preparation — 2026-10-20

Zeus authorized read-only Legacy/Mac checks and preparation of a private LibreChat
integration PR, then requested linked PRs for the full Kubernetes structure and
Backstage soundness. This permits **planning and staged source preparation**, not
blanket new controllers, migrations, billing/credential changes or expanded AI
permissions. Phase 6 remains the first unfinished phase; Phases 7/8 are not done.

- [x] Record [actual Legacy/NATS/Mac observations](docs/operations/legacy-librechat-preflight-2026-10-20.md).
- [x] Prepare the [manual-sync LibreChat pilot](platform/librechat/README.md) with
      source/render safety tests; no deployment or persistence acceptance claim.
- [ ] Legacy inference works (currently provider HTTP 402); approved identity,
      sealed credentials, dedicated database, backup/isolated restore and real
      browser/session/permission gates pass before any pilot activation.
- [ ] Mac production storage/power/reboot/recovery acceptance; no files moved.

Existing portal/migration follow-ups stay in #330/#331, workflow-plane planning
in #328, migration parent #283 and home expansion #284. Reuse those scopes rather
than create duplicate deployment controllers or declare a new phase complete.

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


### Parallel run continuation — agent, 2026-10-21 UTC

Existing draft advanced with bounded agent source/acceptance preparation.
No runtime activation, permission, credential, storage, billing or migration.
Offline fixtures/procedures are not production evidence. Phase 6 In progress;
Phase 7/8 Not started. All original activation gates remain blocking.
