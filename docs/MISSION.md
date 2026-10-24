# OneMind mission and durable requirements

Recorded from the user's instructions on 2026-10-08, America/New_York.
This is project memory, not a claim that the target system is already implemented.

## Mission

Build OneMind into the primary workspace and operational workforce for a technology
founder and father of four children establishing sovereign infrastructure for a
homestead. It must combine family and personal applications, founder/developer work,
land management, security, greenhouse and household operations, office systems,
distributed devices, and tactical/field operations.

The intended AI system is continuously available across the connected environment,
able to observe, understand, coordinate work, and take authorized actions through
real application and device interfaces. A chatbot or a directory of links alone
does not satisfy the mission. TAK and real-time operations are central requirements.

The home will become the primary source of infrastructure and operational authority.
DigitalOcean is the current hosting environment; external/cloud infrastructure will
remain useful for selected services, public workloads, reachability and recovery.
The user intends to build a local mesh with distributed nodes and access the
workspace across platforms whenever connected to that mesh. Selected services may
be exposed publicly through a tunnel in the future.

## User requirements

1. Stay with Kubernetes as the long-term application substrate and manage a growing
   collection of apps with reproducible deployment and a useful representation of
   the actual cluster. Kubernetes does not imply every sensor must run Kubernetes.
2. Make local infrastructure the primary home for essential functions. Design for
   cloud/WAN loss; household and site operation must not require a remote dashboard,
   remote identity provider or remote inference service to remain functional.
3. Provide a coherent main workspace for both daily operations and developer work.
   Prefer established open-source capabilities and plugins where they meet the need.
   Preserve the freedom to develop custom services and interfaces without a product
   ceiling or an obligation to rewrite every app into one framework.
4. Treat TAK, live situational awareness, active dashboards, agents, events and
   executable workflows as first-class capabilities, alongside app lifecycle and
   developer tooling. A periodically refreshed inventory is not the live picture.
5. Give AI authenticated access to app contents, actions and relevant device/site
   capabilities, with explicit identities and contracts. An app's built-in AI does
   not by itself provide access to OneMind's agents.
6. Support authenticated use across web, mobile and desktop over the mesh. Verify
   browser reachability, native capabilities, disconnected reading and disconnected
   control separately. A native wrapper around a remote URL is not offline autonomy.
7. Allow deliberate public exposure of selected applications while keeping private
   family, operational and administration functions under their intended access rules.
   Mesh connectivity and login are separate from authorization to operate equipment
   or infrastructure. Do not infer blanket administrator access from the desire for
   unrestricted development freedom.
8. Make recovery, data export, app backups, identity recovery and infrastructure
   portability part of app selection and onboarding, rather than postponing them
   until the migration home.
9. Extend the same operational model toward managing communities, multiple
   locations and small smart cities. Provide a coherent bird's-eye view across
   Kubernetes, AI work, physical assets and live operations; automate useful work
   rather than only expanding the application inventory.
10. Support the long-term direction of robotics, physical AI, drones and autonomy,
    including useful local operation. Preserve extensibility across devices and
    applications. These goals do not establish unlimited capacity, authorize
    autonomous physical actions or imply that Kubernetes supplies real-time vehicle
    control. Recorded from the user's additional instructions on 2026-10-08.

9. Organize the OneMind catalog by **SO — Sovereign Operations**,
   **LE — Legacy Evolution**, **HP — Holistic Performance** and
   **GE — Generational Empire**, alongside technical/system/site views. Shared
   services must remain one canonical entity with reviewed cross-pillar visibility.
   Mission classification does not grant agent/user permissions or move workloads.
   The user reiterated this requirement on 2026-10-21;
   [catalog foundation](roadmap/onemind-catalog-pillars.md) defines the proposed
   implementation and keeps historical numeric agent routing unconfirmed.

The previously recorded personal workspace requirements remain applicable: calendar,
habit tracking, attachments, AI reads/writes, webhooks, a usable mobile client and
verified offline behavior. See `self-hosted-capability-plan.md` for those gates.

## Existing context and decisions still open

- The current portal is Backstage-based OneMind, with Argo CD, Authentik, Tailscale,
  NATS/JetStream, monitoring and application workloads in Kubernetes.
- `onemindos/omos-ui` is the earlier React/NATS operations interface. Evaluate and
  reuse useful contracts and modules; its existence does not prove current deployment,
  production readiness, native releases or complete offline behavior.
- The user asked for due diligence on Backstage versus Baserow and other open-source
  internal tooling before committing to the long-term workspace approach.
- Backstage as the developer/service portal is the current implementation. Its role
  as the sole daily operations interface remains a product/architecture decision.
- The research recommendation is recorded separately in
  `research/2026-10-08-workspace-platform-review.md`. Recommendations there are not
  accepted deployments or permission changes.
- Home hardware, physical network, outage targets, family/worker roles, autonomous
  action policy, recovery targets and exact app inventory still need definition.
- The full CNCF portfolio screening and proposed physical-AI architecture are in
  the CNCF Ecosystem Review, kept outside this repo at
  `~/OneMind-NEW/CNCF Ecosystem Review/2026-10-08-cncf-onemind-ecosystem-review.md`. Its project dispositions
  and staged acceptance gates are recommendations, not deployment decisions.

## How future work should use this memory

Explain how a proposed capability advances the mission and what it adds to the
maintenance burden. Preserve stable app/service identities and interoperable APIs.
Keep authoritative domain state in its owning service; use events and projections
to connect it to the workspace and agents. Do not assume the catalog, the event bus,
the Kubernetes API and an app database are interchangeable sources of truth.

Older architecture documents include historical technology and deployment claims
(Flutter-only surface, SQLite-only storage, bus-only truth, old endpoints). Reconcile
them against this current user mission, current manifests and observed deployment;
do not treat dated plans as proof of live behavior. This memory records goals and
does not silently amend the established NATS subject grammar or accepted ADRs.
