# OneMind workspace platform due diligence

Research date: 2026-10-08, America/New_York.
Status: recommendation; no new platform deployment, UI migration, network exposure
or permission expansion has been performed.

## Conclusion

Keep Backstage as OneMind's developer and service-management workspace. Build the
real-time operational capabilities on independent APIs, event contracts and site-local
services. Pilot those capabilities as a Backstage operations plugin before deciding
whether a separate mobile/field shell is needed. Preserve useful `omos-ui` code and
contracts instead of restarting or assuming its source can be pasted into Backstage.

Commit to the Kubernetes-based OneMind platform, not to storing every domain's state
inside Backstage or making Backstage availability necessary for every action. This
provides one coherent product and entry point with replaceable interfaces underneath.
There is no evidence that Baserow, Appsmith, ToolJet or Budibase is a better complete
replacement for this combination of developer workflow, cluster context, TAK and
field operations. They remain candidates for narrower applications.

Confidence is high for retaining Backstage's existing developer role, and conditional
for using it as the main operational shell. Source review and upstream documentation
establish extension mechanisms; they do not establish OneMind's performance, mobile
ergonomics, reconnection correctness or disconnected operation. The acceptance tests
below are required to resolve that choice.

## Evidence reviewed

- Local `omos` architecture, control-plane doctrine, requirements and manifests.
- The active portal source in
  `/Users/zeus/Documents/Codex/2026-10-05/i-x20/work/omos-infra`, especially app/backend
  package manifests, permissions, machine-read-only middleware, event observer,
  offline export, templates and the Backstage deployment README.
- GitHub `onemindos/omos-ui` main at commit
  `598d45908c7cc2565956ae167c7f3eb476b7d400` (committer date 2026-07-13): README,
  package.json, NATS connection/profile stores, fleet hook, map page, Capacitor
  configuration and deployment workflow. This is a sampled source review, not a
  complete code/security audit or an executed app/native build.
- Read-only cluster observations earlier in this conversation on 2026-10-08:
  three Ready nodes, 20 synced Argo applications, TAK application degraded, 35 bound
  PVCs using DigitalOcean storage classes, high point-in-time node memory readings,
  failed recent NATS and TAK backup Jobs. This does not verify backup archive contents
  or establish cause of an application health result.
- Official upstream documentation, repositories, licensing descriptions and edition
  pages linked beside findings below. Edition claims are dated and need rechecking
  when selecting the exact version.

## What Backstage can and cannot establish

Backstage is Apache-2.0 licensed. Custom development does not require buying a hosted
portal edition. Its frontend plugin API supports new pages, navigation and utility
APIs; extensions and overrides allow substantial presentation changes. Custom
backend plugins provide integration endpoints. This is a code framework, so freedom
comes with maintaining TypeScript/React integrations and upgrades.
[License](https://github.com/backstage/backstage/blob/master/LICENSE),
[frontend plugins](https://backstage.io/docs/frontend-system/architecture/plugins/),
[extension blueprints](https://backstage.io/docs/frontend-system/building-plugins/common-extension-blueprints/).

Backstage supports real-time push through Signals for notifications. OneMind already
includes Signals frontend/backend packages, but its operational observer checks
Argo and Backstage backup Jobs once per minute, and several portal panels refresh
once per minute. Package installation therefore does not establish continuous TAK
tracks or sensor subscriptions. Those need a bounded operational stream, reconnect
and state resynchronization logic, and a UI built for that data.
[Notifications and Signals](https://backstage.io/docs/notifications/).

The standard Kubernetes integration is designed around service-owner visibility.
It is useful for linking catalog entries to workloads and their health; it is not
the Kubernetes reconciler or a full cluster administration system. Use Argo CD for
desired deployment state, Kubernetes observations for runtime state, and Headlamp
for cluster investigation. Backstage presents those relationships.
[Kubernetes integration](https://backstage.io/docs/features/kubernetes/),
[configuration](https://backstage.io/docs/features/kubernetes/configuration/).

Catalog entity providers can ingest on a schedule or in response to events. Current
OneMind cluster inventory uses an explicitly refreshed script snapshot. Replace
manual refresh with a reviewed provider or discovery service, preserving source,
stable identity and last-observed timestamps. Keep deployment identity separate
from app identity: one app can have several workloads across several sites.
[Entity providers](https://backstage.io/docs/features/software-catalog/external-integrations/).

Custom catalog kinds are possible, but upstream warns that consumers expect built-in
kinds. Use standard Components, Systems, Resources and APIs plus clear annotations
where they fit. Do not write every sensor sample, CoT update or agent execution into
the catalog. Domain registries and event stores own that fast-changing state.
[Extending the catalog model](https://backstage.io/docs/features/software-catalog/extending-the-model/).

Backstage has sign-in integration, a permission framework and an MCP Actions backend.
Authorization must be enforced by each plugin/backend; adding a permission policy
does not automatically secure every custom operation. Current OneMind permits admins
and authenticated read actions broadly, with additional endpoint-specific checks.
That is not yet a demonstrated family/contractor/site access model. Machine MCP is
deliberately read-only and does not grant agents general app writes or scaffolders.
Keep that existing guard while designing a separate capability execution path.
[Authentication](https://backstage.io/docs/auth/),
[permission framework](https://backstage.io/docs/permissions/overview/),
[MCP Actions](https://backstage.io/docs/ai/mcp-actions/).

Community plugins have a published compatibility matrix for the new frontend system.
Verify installed exports, React versions, permission behavior and release compatibility
for each choice. A plugin catalog is not a promise of drop-in compatibility.
[Compatibility matrix](https://github.com/backstage/community-plugins/blob/main/docs/compatibility/new-frontend-system.md).

## What the existing omos-ui contributes

Its actual source contains NATS WebSocket subscriptions, fleet/entity views, mission
modules and agent-facing UI. This is useful prior investment for a command surface.
Its main package uses React 19/TanStack routing; the current Backstage app uses React
18/React Router. Extract transport-independent types, adapters and view logic;
rebuild shell integration or use a separate shell if dependency isolation is needed.
[Pinned source](https://github.com/onemindos/omos-ui/tree/598d45908c7cc2565956ae167c7f3eb476b7d400).

It is not yet proven as the long-term multi-user base:

- The profile store persists NATS passwords in browser localStorage. The build workflow
  supplies VITE_NATS_PASS as a frontend build argument. Vite frontend environment
  values must not be treated as server-held secrets. Replace shared static browser
  credentials with per-user, scoped, short-lived access or a server-side gateway.
- Capacitor config loads `https://onemindos.dev`. Native scaffolding is not evidence
  of a locally bundled app that can open while that remote host is unreachable.
- The sampled map page summarizes entity counts and launches external tactical views;
  it is not itself a complete common operating picture.
- Its entity subscription uses `ent.>.state`; verify and correct this before reuse,
  because NATS's `>` wildcard must be the last token in a subject subscription.
- The workflow publishes images and directly restarts a DigitalOcean workload. Align
  any resumed release path with immutable image digests and the established Argo
  deployment authority.
- README scripts/deployment claims differ from parts of the current package/workflow.
  Treat the repository as reusable source needing verification, not a working product
  inferred from documentation.

[Profile store](https://github.com/onemindos/omos-ui/blob/598d45908c7cc2565956ae167c7f3eb476b7d400/src/stores/nats-profiles.store.ts),
[map page](https://github.com/onemindos/omos-ui/blob/598d45908c7cc2565956ae167c7f3eb476b7d400/src/pages/map/MapPage.tsx),
[Capacitor configuration](https://github.com/onemindos/omos-ui/blob/598d45908c7cc2565956ae167c7f3eb476b7d400/capacitor.config.ts),
[workflow](https://github.com/onemindos/omos-ui/blob/598d45908c7cc2565956ae167c7f3eb476b7d400/.github/workflows/deploy.yml).
The environment-variable behavior and subject constraints are documented upstream:
[Vite environment variables](https://vite.dev/guide/env-and-mode),
[NATS subjects](https://docs.nats.io/concepts/subjects).

## Alternatives and their role in OneMind

These fit assessments are architectural judgments based on documented capabilities,
not benchmark scores. No candidate was deployed or tested against the homestead.

| Candidate | Strength and source | Important qualification | Recommended role |
|---|---|---|---|
| Backstage | Developer catalog, docs, software templates, extensible code UI; sources above | TAK, device actions, mobile/offline and app AI access require integrations | Keep developer/service workspace; pilot operations plugin |
| Custom OneMind UI | Existing NATS/fleet/agents/mission source and native scaffolding | Identity, reconnect/replay, stale state, release path and native behavior need verification | Reusable operations components; separate field shell only if needed |
| Baserow | Databases, app builder, dashboards, automations and API-first design; [repository](https://github.com/baserow/baserow) | Open-core: non-premium/non-enterprise features MIT. SSO requires Advanced/Enterprise, not the free baseline; [SSO documentation](https://baserow.io/user-docs/single-sign-on-sso-overview) | Land records, maintenance, inventory and administrative data; not cluster/TAK authority |
| Appsmith | Apache-2.0 community internal app builder, APIs/databases, Kubernetes install; [repository](https://github.com/appsmithorg/appsmith) | OIDC/SAML, private embedding and air-gapped edition sit in Enterprise/add-on offerings; [edition page](https://www.appsmith.com/pricing) | Optional forms/admin panels when edition limits are acceptable |
| ToolJet | AGPL-3.0 foundation for internal dashboards, apps and workflows; [repository](https://github.com/ToolJet/ToolJet) | Current OIDC docs label it Team Plan; multiple providers Enterprise. Verify exact self-hosted entitlement; [OIDC](https://docs.tooljet.com/docs/user-management/sso/oidc/setup/) | Optional internal tools, not an unrestricted universal base |
| Budibase | Internal apps, automations and agents; [repository](https://github.com/Budibase/budibase) | GPLv3 core, MPL client/components, BSL paid features. Current free self-hosted plan includes SSO and unlimited apps/users, but one workspace; audit logs, enforced SSO and built-in backup/restore are Enterprise; [edition page](https://budibase.com/pricing/) | Worth a bounded prototype for forms/workflows; do not confuse built-in agents with OneMind fleet access |
| Frappe Framework | MIT full-stack Python/JS framework, permissions, REST and Socket.IO realtime; [repository](https://github.com/frappe/frappe), [realtime](https://docs.frappe.io/framework/user/en/api/realtime), [REST](https://docs.frappe.io/framework/user/en/api/rest) | Different application model and Python/database/runtime to operate; Kubernetes/GitOps/TAK integration still needs building | Stronger unrestricted alternative for custom business apps than a reason to replace the existing portal |
| Grafana | Existing observability; Live supports streaming panels via WebSockets; [Live](https://grafana.com/docs/grafana/latest/setup-grafana/set-up-grafana-live/) | Soft realtime; datasource/adapter work required. Dashboard streams are not command acknowledgments or agent execution | Keep telemetry, history and operations wall displays |
| OpenRemote | Assets, protocol adapters, rules, roles, dashboards and edge gateways; [repository](https://github.com/openremote/openremote), [edge gateway](https://docs.openremote.io/docs/user-guide/gateways-and-devices/edge-gateway/) | Separate IoT domain platform and identity/integration maintenance; no automatic NATS/TAK/cluster sync | Evaluate for greenhouse/land/device control if it fills a demonstrated gap |

Do not deploy all of these. Baserow is broader than a spreadsheet now, and Budibase
is broader than a simple CRUD builder, but neither product description proves the
required field/TAK behavior. Built-in AI branding does not establish authenticated
access for the existing OneMind agents. Paid features also do not mean a project is
unusable; they must be explicit selection tradeoffs rather than surprises.

## Proposed architecture

One product and entry point, with role-appropriate views: Today, Operations, Agents,
Apps and Build. Backstage provides the current shell, developer workspace and service
catalog. A OneMind operational API/gateway supplies curated snapshots and live deltas,
permission checks and command results. Existing domain services retain authority.

Authority boundaries:

- Git/Argo: desired workload configuration and release history.
- Kubernetes: observed scheduling, readiness and runtime resources.
- Backstage: human/service descriptions, ownership, relationships and capability discovery.
- Domain applications: tasks, documents, calendars, assets and other owned records.
- TAK/CoT plus device adapters: situational entities and field workflows.
- NATS/JetStream: event delivery, selected durable streams and replay; it does not
  automatically synchronize arbitrary application databases.
- Agent runtimes: inference, execution state and tool calls under their own identities.
- The UI: present state and request actions; do not make it a necessary execution hop.

Use standard authenticated APIs/MCP for app tools and bounded event projections for
live context. Share underlying execution services between human controls and agent
tools. Each capability advertises its owner, input/output schema, scope, site,
availability, timeout, idempotency, audit behavior and permissions. Standing user
authorization can cover ordinary actions; infrastructure and consequential equipment
actions follow their deliberately defined policies. This is a recommendation, not a
new grant of authority.

Keep meshes and sites separate from the Kubernetes cluster's failure domain. The
local physical network must carry local traffic without WAN access; a mesh overlay
does not provide the radio/ethernet/Wi-Fi links underneath it. Local DNS, login and
service routes must remain available when outside services are unreachable.

Recommended placement:

| Site | Role | Disconnect behavior to prove |
|---|---|---|
| Home core | Primary apps/data, identity, event hub, operations, inference where hardware permits | Works over local network with ISP/cloud disconnected |
| Greenhouse/main house/office gateways | Sensors, deterministic local rules, device adapters, optional small Linux workload nodes | Local controls and event buffering survive loss of hub |
| Tactical edge | Small independent local services/bus, field clients and maps | Starts and operates without the cloud; selected data catches up later |
| External cloud | Public services, relay/recovery and selected compute | Loss of cloud does not stop essential home/field functions |

Independent clusters at genuinely disconnected sites are the default recommendation;
nearby reliable LAN nodes can share a home cluster. K3s is a candidate for Linux
gateway/home clusters, not a requirement for microcontrollers or every sensor.
KubeEdge is an alternative to evaluate for managed edge autonomy, not an assumed
installed capability. Avoid an internet-spanning control-plane quorum or requiring
remote storage for local controls. Kubernetes supports multi-zone deployments but
that alone does not establish the disconnected-site behavior required here.
[K3s](https://docs.k3s.io/), [KubeEdge](https://kubeedge.io/),
[Kubernetes multi-zone behavior](https://kubernetes.io/docs/setup/best-practices/multiple-zones/).

NATS leaf links connect sites; local JetStream needs a distinct domain and deliberate
mirrors/sources. Retention, catch-up, expiry and conflict policy must be tested. Never
replay an expired equipment command merely because connectivity returned.
[Leaf nodes](https://docs.nats.io/learn/topologies/leaf-nodes),
[mirrors/sources](https://docs.nats.io/learn/jetstream/mirrors-and-sources).

Public exposure should target selected app routes/services with their own access
model. A tunnel provides reachability, not application authorization. Public apps
should not inherit portal administration or raw bus access. Keep the home workspace
usable locally even when the public tunnel is down.

## App onboarding and continuous cluster alignment

Every new app needs one stable application identity; registered owner/purpose/site;
pinned release and Argo ownership; resource budget; local/private/public routes;
human and machine identities; APIs/MCP/events; persistent storage; backup/restore;
and health/dependency observations. Register underlying workloads with explicit
app and cluster identifiers and source/observed timestamps.

Use a reviewed entity provider for catalog discovery, Argo/Kubernetes observations
for rollout and health, and NATS for domain events. A cluster discovery failure
must produce unknown/stale observations rather than delete every known app.
Avoid turning every pod replica into another app card. Maintain desired, observed,
documented and historical states separately and show discrepancies.

AI autonomy should have an execution ledger, durable tasks, retries, cancellation,
capability scopes, budgets and explicit outcomes. An omnipresent agent experience
can be distributed across local/site workers; it need not be one central process
holding everyone's credentials. Necessary deterministic device control remains in
local rules/controllers. Inference and command availability are distinct.

## Required proof before committing to one operational shell

Build a small, isolated pilot with real authorized sources and a simulated actuator.
No production equipment needs to be actuated to test the architecture.

1. A TAK track and greenhouse sensor update appear in a custom operations view.
   Measure latency at representative event rates and define an acceptable freshness
   target; show data age and source health, coalesce/backpressure updates, and avoid
   conflating real-time web UX with hard real-time equipment control.
2. A human and agent use the same scoped capability to request a simulated action;
   both receive accepted/failed/completed outcomes with correlation IDs. Transport
   delivery is not proof that a valve or device reached its intended state.
3. Drop the browser stream and reconnect: obtain a fresh snapshot plus sequence-aware
   deltas, suppress duplicates and reveal gaps. Revoke an identity during an active
   stream and verify reads/actions stop under the intended policy.
4. Add/change a test app through GitOps: it appears in the service catalog and links
   to actual workloads without a manually edited directory; deletion and discovery
   outage produce the correct distinct behavior.
5. Use phone/tablet and desktop views for maps, alerts, chat and action acknowledgment.
   Test login, background/suspend/reconnect and any required camera/location/voice
   behavior. Verify native packaging separately if it is required.
6. Disconnect WAN/cloud: new local login, DNS, dashboards and essential rules still
   work. Separately isolate an edge site from home and prove its restart/local rules,
   buffering, stale state, catch-up and expired-command handling.
7. Restore representative app data, identity/config and event state on another cluster
   using a different storage class. Demonstrate the home migration path instead of
   assuming Kubernetes manifests move the data.
8. Upgrade Backstage and the operations plugin in the test environment; verify API
   compatibility and that ongoing domain/agent work survives the portal restart.

If Backstage's operations view passes, keep it as the primary browser workspace.
If native/offline UX or upgrade coupling fails, keep Backstage for Build/Services and
provide a lightweight OneMind operations client against the same APIs and identity.
The core service and data investment survives either result.

## Sequence and unresolved inputs

First fix/verify backups and reproduce ownership of existing deployments. Measure
capacity before adding builders or more agent runtimes. Define capabilities and app
onboarding, then implement the operations pilot and catalog provider. Test home/site
autonomy before the final cloud-to-home cutover.

Hardware, site link topology, family/worker roles, required iOS/Android/native features,
inference placement, event rates, recovery time/data-loss targets and autonomous
equipment actions remain open. These affect sizing and acceptance criteria; they do
not require discarding the existing Backstage investment.
