# OneMind Backstage experience and capability foundation

Requested 2026-10-21. **Planning/source review only: no UI, dependency, permission,
controller or deployment changes.** Extend master PR #356, not another portal or
an untracked plugin spree. The eight workstream PRs remain drafts; their offline
contracts are not deployed features. Phase 6 In progress; Phases 7/8 Not started.

## What the supplied list adds

The target application flow is useful, but it needs an explicit ownership
[decision](../architecture/decisions/0002-application-delivery-and-infrastructure-boundaries.md),
one mapped application contract, and an experience programme that joins native
features rather than just adding tabs. The existing eight streams cover much of
the backend/lifecycle work; they do **not** constitute a complete UI upgrade.

Missing cross-cutting foundations to track:

1. A deliberately designed navigation/home/entity experience, not a long plugin menu.
2. Permission-aware search and discoverability for apps, APIs, runbooks and sites.
3. A versioned operational evidence model: desired, observed, application/data,
   cost, recovery and action states, with source/time/scope separated.
4. Catalog ingestion ownership/freshness/conflict handling, without duplicate IDs.
5. Approved operational notifications, deduplication and an action/audit timeline.
6. Capability/plugin lifecycle governance: compatibility, permission support,
   accessibility, upgrades, licensing, performance and an owner for every panel.
7. Domain/device/ML and AI-resource representation, with app-native interfaces and
   local/offline authority preserved. Kubernetes visibility is not equipment control.

“Full scope” means deliberately covering the useful native surfaces and measured
gaps, not installing every plugin or promising a limit-free/100%-certified system.

## Evidence boundary and current source

Reviewed tracked source and official project docs on this date; no fresh browser,
app, policy-report, device or home-cluster acceptance was performed for this plan.
Earlier Grafana/Prometheus/topology receipts remain bounded historical evidence.
Do not repeat the pasted list's old Grafana credential 401 as a current open bug:
#340–#345 accepted that scoped repair. Wider mapping/coverage remains outstanding.
The earlier zero-policy-report observation is likewise not a fresh live query.

The repository pins Backstage 1.55.3, `@backstage/ui` 0.18.0 and frontend plugin API
0.18.1; MUI and BUI coexist. `packages/app/src/App.tsx` uses `createApp` and native
frontend features; config enables package discovery. `packages/backend/src/index.ts`
wires catalog, scaffolder, TechDocs, Postgres search/collators, Kubernetes,
permission/auth, events, notifications, signals and MCP actions plus OneMind modules.
Registration/dependency presence is not end-to-end user acceptance.

Source seams under `services/onemind-portal/`:

- `plugins/plugin-onemind-nav/src/Sidebar.tsx`: Home/Today/Catalog/Create plus an
  alphabetically ordered remainder; proposed task/site grouping must be deliberate.
- `plugins/plugin-onemind-home/src/homeModule.tsx`: launcher, attention/status,
  native home blueprints and existing operational widgets; preserve favorites.
- `plugins/plugin-onemind-theme/src/{designModule.tsx,brand.css}`: existing brand,
  light/dark MUI theme and BUI variables, responsive and reduced-motion rules.
- `packages/app/src/{componentOverviewModule.ts,catalogCompositionModules.ts}`:
  entity-kind composition and native/About/Relations deduplication; preserve the
  existing OpenChoreo/new-frontend compatibility regressions.
- `plugins/plugin-onemind-services*`: joined Operations and scoped Prometheus
  views; `plugin-onemind-events-backend` already polls selected Argo/backup events
  with database-backed deduplication. Enrich it, do not assume no event producer.
- `app-config.yaml` versus `platform/backstage/app-config.cluster.yaml`: defaults
  are not effective production configuration; read both before declaring a gap.

## Native Backstage capability review

“Wired” below is source evidence, not a claim that every user/app path passed.
Community integrations are distinguished from Backstage core.

| Capability | Existing foundation | Next useful acceptance / extension |
|---|---|---|
| [Catalog, ownership, relations](https://backstage.io/docs/features/software-catalog/) | Native catalog/org/graph/import packages, curated inventory and current pilot provider | Stable app/API/site/resource ownership, ingestion provenance, stale/missing/conflicting/orphaned records, no silent takeover or fake discovery. #353/#355. |
| [API documentation](https://backstage.io/docs/features/software-catalog/descriptor-format/#kind-api) | API Docs dependency and STAC API linkage | Actual versioned app/tool APIs, ownership/auth documentation, schema and relationship checks. An API definition is not authenticated AI access. |
| [Software Templates](https://backstage.io/docs/features/software-templates/) | Existing onboarding/change/recovery draft-proposal templates | Reuse editor/preview/dry-run where supported and permitted, task history, safe parameters and app-specific lifecycle contract. #351/#339; dry-run permission is not live approval. |
| [TechDocs](https://backstage.io/docs/features/techdocs/) | External build/packaged publication and backend | Per-service operational/recovery docs, tested links/freshness, source edit links/addons where justified; no unsafe in-request builds or unreviewed extension execution. |
| [Search](https://backstage.io/docs/features/search/) | Postgres search; catalog, TechDocs, ADR and announcement collators | Entity/API/runbook facets, scoped result/snippet permission checks, index freshness and removal/reindex behavior. Reuse current engine first; OpenSearch is not required merely because it runs. |
| [Kubernetes](https://backstage.io/docs/features/kubernetes/) | Native backend/frontend packages; bounded object types and reviewed first-service mappings | Correct selector/namespace/cluster and owner per app, useful events/logs/readiness; no-match/stale/denied states, then actual second-site acceptance. #355/#284. |
| [Auth and permissions](https://backstage.io/docs/permissions/overview/) | Authentik resolver, one host policy, machine guard and per-route restrictions | Existing-role/session/expiry/backend negative tests and context redaction; each new plugin must actually integrate permission checks. #354/#339. |
| [Notifications and Signals](https://backstage.io/docs/notifications/usage/) | Backend modules, notification sidebar and selected operational-event producer | Recipient/tenant scoping, stable IDs, dedupe/reconnect/replay/retention/severity and preferences; durable event history separate from transient UI signals. |
| [AI catalog/MCP](https://backstage.io/docs/ai/ai-in-the-catalog/) | Existing read-only [MCP action registry exposure](https://backstage.io/docs/ai/mcp-actions/) and internal assistant | Review pinned-schema support for AI resources/models/prompts/tools/MCP API relationships; no new kind or tool merely from latest docs. Preserve machine filtering, user context consent, per-tool budgets and audit. |
| Home, favorites, visits/settings | Native home blueprints and account-backed personal settings plus custom modules | Personal pinning/layout where actually supported, clear reset/migration and per-user isolation. Home widgets are not a general no-code layout builder for every plugin page. |
| GitHub/Argo/Grafana/Prometheus (community/custom adapters) | Native/community views and accepted scoped service slices | App mappings, deployed-vs-desired revisions, correct metric/alert scope, real failure states, no fake green hosted CI. Cost allocation is #352, distinct from DO account charges and agent spend. |
| Tech Insights, ADR, Radar, announcements (community) | Source packages/modules and metadata fact retrievers | Useful reviewed scorecards, accepted decisions and change communications; metadata completeness != SLO, compliance or recovery. Avoid another duplicate dashboard. |
| Extension tooling / plugin maintenance | App visualizer dependency, native module/blueprint seams, existing tests | Compatibility inventory, feature flags and composition diagnostics, administrator-only sensitive diagnostics; no devtools/secret exposure. |

Use current official documentation as a capability reference; installed exports,
lockfile, permission behavior and fixtures decide compatibility. The current
[release notes](https://backstage.io/docs/releases/v1.55.0/) specifically describe
dry-run/SCM credential permission changes: test the pinned implementation rather
than widening our policy just to match a sample.

## Coverage across the twenty Kubernetes categories

The [operating structure](../architecture/kubernetes-operating-structure.md) is
the detailed twenty-category inventory. Map it into four coherent surfaces rather
than twenty disconnected plugin pages:

| Surface | Categories represented | Evidence, not a blanket capability claim |
|---|---|---|
| Delivery | CI/CD, app definition/build, registry | Source/build/provenance/digest, review/approval, Argo and owning app lifecycle controller; hosted execution restrictions visible. |
| Runtime/infrastructure | Automation/configuration, scheduling/orchestration, runtime, storage, databases, discovery/coordination | Native Kubernetes plus accepted operator/Crossplane resource status, site/owner, dependencies and app-level data/recovery checks. |
| Connectivity/domain operations | Messaging, gateway, RPC, service mesh/proxy and workload frameworks | Scoped NATS/API and Cilium/Hubble or app-owned projections; actual telemetry distinct from catalog wiring. A framework is evaluated for a workload need, not installed to fill a category. |
| Reliability/trust | Observability, optimization, security/compliance, key management and chaos engineering | Alerts/SLO coverage, cost basis, observed policies, vulnerability/provenance/identity and recovery evidence; gated resilience exercises rather than an unrestricted fault-injection button. |

Harbor/Trivy, Kyverno/Falco/SPIRE, key recovery and backup systems need real
scoped data contracts before their results become security badges. Existing
controllers alone do not prove signing enforcement, mTLS, compliance or recovery.

## Proposed experience, preserving OneMind

This is a product/layout specification, **not a visual acceptance audit or a
promise that the following views already exist**. Preserve brand assets, stable
URLs/refs, existing home/favorites/settings and useful native app/TAK interfaces.

### Navigation and home

Design for founder/operator, future family users and scoped agent/worker roles
without assuming one view or administrator access suits everybody. Persona
preferences are not authorization; family content and agent page context require
separate privacy/consent and backend access review. The mission still leaves
Backstage's role as the sole everyday operations surface open; preserve TAK and
application-native interfaces instead of forcing every task into the catalog.


Proposed groups: **Home / Operations / Applications / Sites & Infrastructure /
Agents & Runs / Build & Change / Knowledge / Settings**. These are conceptual
entry points to existing features, not eight replacement plugins or an immediate
route rename. Keep admin/rarely used integration tools in a secondary area;
page grouping never substitutes for backend permissions. Validate all deep links.

Home should answer “what needs me, where, why and what can I do safely?”:

```text
Search apps / APIs / runbooks       Site + environment + observation scope
Needs attention (unknown/stale visible)     Pending approvals / agent runs
Pinned apps + recent work                  Costs / recovery exceptions
Current changes + correlated events        Open app / Inspect / Propose change
```

Do not duplicate the entire application directory and every native dashboard on
one screen. Keep progressive disclosure, compact summaries and real source links.
Preserve native home preferences; future widget configuration needs allowlisted
widgets/data and versioned settings, not arbitrary remote JavaScript.

### One consistent app/service page

Header: human name + canonical ref, owner, location/environment, deployment owner,
last observation and separate desired/actual version; **Open app**, **Inspect**
and a clearly labelled **Propose change** (not an implied approved deploy button).

Proposed content groups: **Overview / Runtime & Delivery / Observability /
Dependencies & Infrastructure / APIs & Knowledge / Recovery & Policy / Activity
& Approved Actions**. Retain native plugin routes/deep links; group existing tabs
where supported and do not render every integration for every entity kind.

- Overview shows a bounded summary and primary dependency map, not repeated
  About/Relations sections and long empty tables. Missing vs unsupported must be
  explicit; keep relevant missing-configuration guidance, not invisible gaps.
- Kubernetes/Argo show runtime versus desired Git/artifact independently.
  App semantic/data readiness, alert scope and source freshness remain separate.
- Recovery displays current mechanism, last verified restore and failure domain;
  distinguish superseded backup jobs from the real independent-copy gap.
- Policy shows observed reports with coverage/time; absent reports mean unknown,
  never “compliant”. Cost states show units/window/currency/estimate basis.
- Activity correlates approved changes, observed rollout, alerts and agent runs;
  authorization is enforced by the backend, not a button's disabled/hidden state.

### OneMind mission structure in the catalog

**SO — Sovereign Operations; LE — Legacy Evolution; HP — Holistic Performance;
GE — Generational Empire** are required catalog pillars, not just home categories.
The [pillar catalog specification](onemind-catalog-pillars.md) defines native
OneMind/child Domains, reviewed System membership, primary/shared classifications,
UI/search/onboarding filters and stable-reference/privacy/refresh gates. Keep
mission, deployment/site and authorization as independent axes. No catalog
reclassification or numeric agent routing is applied by this foundation.

### Entity and site model

Use native Component/System/Domain/API/Resource/User/Group semantics deliberately.
A namespace System is not a separate cluster. An app can run at multiple sites;
represent deployment instances without duplicating the app identity or asserting
an undeployed home/edge location. Exact site/instance schema is a reviewed design
output; don't silently add custom kinds/catalog ownership in this plan.

Future AI/models/tools, ML jobs, devices/vehicles and Home Assistant/OpenRemote/TAK
adapters should link to their owning systems and approved APIs. Keep domain truth
in those services; a unified workspace is not a replacement for their databases,
safety controls, identity or offline interfaces.

## Shared UI/data contract and quality bar

Every panel records identity/site/environment, source/scope, observedAt,
source sample/evaluation time where relevant, refresh cadence, expected coverage,
state/reason and a redacted drill-down. Distinguish loading, available, empty,
not-configured, partial, stale, unavailable and forbidden. Preserve last-good data
only when clearly labelled stale; never convert a failed query into zero/green.
This is the proposed model; a concrete versioned schema plus mappings are future
implementation work alongside #351/#354/#355, not a new redundant source of truth.

- **Design:** consistent spacing/density/typography, two themes and all critical
  states; prefer native components and blueprints. Backstage's
  [theming docs](https://backstage.io/docs/conf/user-interface/) explicitly cover
  coexisting MUI and BUI. Audit both; no blind CSS-token rename, global override,
  framework upgrade or rewrite of the working OpenChoreo compatibility bridge.
- **Accessibility:** keyboard/focus/labels, screen-reader state text, contrast,
  reduced motion, 200% zoom, 320/768/1440px layouts and real mobile/touch checks.
  Wide tables may scroll, but primary actions and scope must remain reachable.
  Existing native alert-state text accessibility needs explicit follow-through.
- **Performance/resilience:** proposed measurable budgets set at baseline; bounded
  request size/concurrency/timeouts, cancellation, caching/freshness, lazy panels
  and large-catalog paging. Test partial provider failure and loss/reconnect.
- **Privacy/security:** no provider tokens or raw sensitive YAML/log/chat payloads
  in UI/search/notifications/agent context. Inspect only allowed fields; Secrets
  remain inaccessible and arbitrary proxy/exec/write paths stay denied.
- **Acceptance:** real user tasks plus screenshots at agreed sizes/themes, direct
  role/API denials, durable preferences, native route/plugin regressions and the
  actual production image/Argo/replica checks. Fixtures are not live UI proof.

## Conditional integrations from the supplied list

These are evaluation gates, not a request to install. No plugin is certified
compatible with this pinned cohort just because its README exists.

| Integration | Keep / evaluate | Specific gate |
|---|---|---|
| [TeraSky Kubernetes resources](https://terasky-oss.github.io/backstage-plugins/plugins/kubernetes/overview/) | Evaluate deeper resource graph only after native coverage needs are measured | CRD allowlist, field/log/YAML redaction and frontend/backend permissions; no new Secret/ConfigMap/exec/write access. Compare against native graph/drill-down first. |
| [Kubernetes Ingestor](https://terasky-oss.github.io/backstage-plugins/plugins/kubernetes-ingestor/overview/) | Scoped ingestion pilot, not blind adoption | Stable refs/provider ownership, exclusions/dedupe, refresh/deletion behavior, no collision with curated inventory or OpenChoreo entities. |
| [Kyverno reports](https://terasky-oss.github.io/backstage-plugins/plugins/kyverno/overview/) | Useful security projection after reporting is demonstrated | Real report producer/coverage, RBAC/retention, no empty-equals-compliant or implicit policy mutation/MCP exposure. |
| [Crossplane suite](https://terasky-oss.github.io/backstage-plugins/plugins/crossplane/overview/) | Evaluate against selected infrastructure target | Actual accepted Crossplane resources, graph/health/claim mapping, pinned plugin compatibility; audit every advertised MCP action before allowlisting any. |
| [LibreChat plugin](https://github.com/nosportugal/backstage-plugin-librechat/blob/main/README.md) | Optional portal surface after #332 works | Existing server/API compatibility, per-user sessions/identity, streaming, consent/redacted page context and spend enforcement; no shared browser credential. |
| [Agent Forge/CAIPE](https://github.com/backstage/community-plugins/blob/main/workspaces/agent-forge/plugins/agent-forge/README.md) | Alternative evaluation, not an additional mandatory agent backend | Its README expects CAIPE; no proven Legacy/NATS drop-in. Compare contracts before introducing another framework. |
| KubeVela UI integration | Target capability, plugin choice unresolved | Verify maintained native/new-frontend package or use existing Kubernetes/Argo/API panels; do not invent a supported plugin or replace the portal. |

## Ordered implementation packages — future PRs, not already delivered

| Package | Foundation / deliverable | Dependencies and pass gate |
|---|---|---|
| U00 — baseline and design contract | Browser/source inventory, user tasks, route/entity/plugin and SO/LE/HP/GE mapping, screenshots, a11y/performance baseline, shared evidence states and UI tokens | Remaining trust gates tracked; no fresh visual audit claimed by this document. Preserve exact installed package exports/permission behavior. |
| U01 — navigation/home | Task/site grouping, search/favorites, attention/approval entry points and responsive home | U00; stable deep links/settings, no lost favorites, two themes/mobile/keyboard and actual user acceptance. |
| U02 — entity composition | Consistent app pages, kind-specific cards/tabs, source/scope/freshness and nonduplicated summaries | U00 + #355 mappings; regression across native and OpenChoreo kinds, no hidden unknowns or newly exposed API. |
| U03 — knowledge/discovery | Catalog/API/runbook search facets, ingestion diagnostics, ownership and useful metadata scorecards | #353/#355 + permission checks; actual search/index refresh/removal/denial tests, not a new search cluster by default. |
| U04 — operational evidence | Approved events/timeline, notification routing, policy/recovery/cost coverage and source labels | #352/#354 and existing event producer; dedupe/reconnect and observed evidence; absent sources explicit. |
| U05 — actions/agents/sites | Native entry points for approved requests/runs, real agent sessions and accepted sites | #339/#332/#284; backend gates and real lifecycle/local recovery first. No parallel grant of privileges by a UX PR. |
| U06 — release/maintenance | Feature-flagged source/image/pin PRs, upgrade matrix and UX/security acceptance receipts | Exact-head gates and sequential Argo release; no big-bang UI/framework/plugin upgrade. |

U00/U01 preparation can proceed alongside the eight source drafts. Runtime
activation, provisioning and migrations retain their own gates. KubeVela and
Crossplane selection/ownership acceptance is in the linked decision and K08;
no UI package installs them or retires OpenChoreo. Keep all twenty categories
visible in the operating map, without twenty duplicate dashboards/controllers.

**Next implementation:** U00 first, then U01 navigation/home and U02 app-page
composition using current native capabilities. Do not start with a wholesale
Backstage upgrade or mass plugin installation. No U00–U06 implementation PR is
claimed opened by writing this foundation.
