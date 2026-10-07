# Self-hosted capabilities for autonomous OMOS sites

Research date: 2026-10-07. This is an architecture recommendation grounded in current upstream documentation. Except for n8n, these are researched candidates, not deployments or removals. Hardware inventory, intended fleet size, and site failure requirements still need mapping.

## Design criteria

Each capability needs a defined authority, stable AI-accessible interface, separate identity, scoped actions, durable state, export/backup route, and a recovery procedure. Self-hosting establishes control of infrastructure; it does not by itself provide decentralization, isolation, security or offline operation.

A site should boot, authenticate locally, collect data, and execute essential rules without the cloud. Its hub synchronizes selected state and distributes reviewed configuration. An AI should receive contextual data with provenance and timestamps and invoke commands with identifiers, expiry, validation, deduplication and acknowledgments. Critical deterministic controls remain local. Human approval requirements apply to defined actions rather than every API call.

## Recommended capability owners

| Capability | Candidate owner | Placement | Decision |
|---|---|---|---|
| Cross-service events and durable delivery | Existing NATS/JetStream | Hub plus independent local buses/leaf nodes | Keep and document schemas/subject permissions |
| Human and machine private access | NetBird or Nebula, evaluated separately | Control/discovery independent of the cluster it must recover | Select after an isolation/recovery pilot |
| Cross-application workflows | n8n | Hub, with separate local instances for critical offline workflows | Deployment authorized and added in platform/n8n |
| Physical asset context and local rules | OpenRemote | Local site first, central context projection | Strong first IoT pilot |
| Device protocol translation | EdgeX Foundry | Industrial/site gateways | Adopt where device mix justifies it |
| Kubernetes edge workload lifecycle | Independent K3s, or validated KubeEdge | Local control plane for independent sites; CloudCore/EdgeCore for managed edge nodes | Compare using actual hardware and outage behavior |
| Firmware/OS/ML artifact delivery | RDFM | Hub manager with hardware-specific device clients | Strong candidate for heterogeneous embedded fleets |
| Knowledge used by humans and AI | Outline | Self-hosted hub, exported/cached critical site documentation | Preferred operational wiki candidate; BSL terms acknowledged |
| Project and work management | Plane Community | Hub | Preferred ClickUp/project replacement candidate |
| Personal recurring tasks | Vikunja | Hub/local installation as needed | Preferred Todoist replacement if Plane is unsuitable for daily use |
| Saved research/context | Karakeep | Hub plus local models | Preferred bookmark/context candidate |

## Mesh evaluation

| Candidate | Strength | Tradeoff | Assessment |
|---|---|---|---|
| NetBird | Self-hosted control services, policies, identity integration, management API and Kubernetes operator | Control plane remains centralized; active-active management/signal HA is documented for Enterprise | Best general managed-access candidate for users, machines and Kubernetes |
| Nebula | MIT-licensed certificate-based host mesh, local policy, discovery lighthouses and static host maps | Certificate enrollment/rotation/revocation and fleet automation become our responsibility; fewer turnkey human-access workflows | Stronger sovereignty-oriented machine-mesh candidate to pilot for disconnected sites |
| Headscale | Self-hosted Tailscale-compatible control server with routes, DNS, policy and OIDC | Current upstream tracking issue states Kubernetes operator API/OAuth compatibility is unresolved | Do not treat as a drop-in replacement for this operator-dependent cluster |
| Netmaker | WireGuard networking with community and commercial deployment paths | Relevant deployment guides use Pro-only features; check edition boundaries for required capabilities | Secondary candidate, not established as better than NetBird here |
| Plain WireGuard | Small, explicit static topology with minimal control services | Enrollment, routes, policy, discovery and key lifecycle need substantial automation at fleet scale | Useful constrained site links, not a complete fleet-management answer |

Recommend a two-candidate pilot, not an immediate production replacement. Test NetBird for managed access and Nebula for site-machine autonomy. Use separate pilot hosts, disjoint overlay/subnet ownership and DNS zones, and no competing default routes. Neither peer discovery nor NAT traversal can create connectivity where the physical network does not exist. Preprovision local access and recovery paths; cloud-hosted identity must not be the only way to administer a disconnected site.

Sources: [NetBird Kubernetes](https://docs.netbird.io/use-cases/kubernetes/routing-peer), [NetBird HA](https://docs.netbird.io/selfhosted/maintenance/scaling/high-availability), [NetBird licenses](https://github.com/netbirdio/netbird/blob/main/LICENSE), [Nebula](https://nebula.defined.net/docs/), [Nebula certificate setup](https://nebula.defined.net/docs/guides/quick-start/), [Headscale operator compatibility](https://github.com/juanfont/headscale/issues/3086), [Netmaker editions](https://help.netmaker.io/en/articles/10209389-comparing-netmaker-community-and-professional-editions-which-one-is-right-for-you).

## Edge tools are different layers

### KubeEdge

CloudCore and EdgeCore extend Kubernetes workload/device management to edge nodes, with local metadata persistence and disconnected operation. It complements a mesh and device adapters; it does not replace firmware updates or knowledge management. Its current upstream README matrix lists KubeEdge 1.23 compatibility through Kubernetes 1.32; the hub is Kubernetes 1.36. Compatibility with this hub was not established. Do not install into production until that is resolved. A centrally administered edge node is also a different autonomy boundary from a site with its own local Kubernetes API and scheduler. [Project and matrix](https://github.com/kubeedge/kubeedge), [edge autonomy](https://kubeedge.io/).

For independent multi-service sites, assess K3s with a local control plane and locally available Git/image caches. Argo at the hub cannot reconcile while disconnected; a local reconciler needs local desired state access. For small gateways, Compose/systemd can be appropriate. [K3s architecture](https://docs.k3s.io/architecture).

### openBalena

A self-hosted balena device/container management stack centered on balenaOS and balena CLI. Upstream documents single-user operation and the absence of the balenaCloud dashboard and binary container deltas; community dashboards are separate projects. It is an alternative fleet operating model, not something to install alongside KubeEdge on every device. Good fit if hardware is deliberately standardized on balenaOS. RDFM/Kubernetes better match a heterogeneous system that needs separate OS and workload ownership. [Upstream comparison](https://github.com/balena-io/open-balena/blob/master/README.md).

### RDFM

Antmicro's Apache-2.0 framework for OTA/fleet management across Linux, Android and Zephyr, including update artifacts, delta delivery and Linux A/B rollback. Server/API and device clients are separate. Linux robust rollback requires the correct partition layout and bootloader integration; installing its server cannot add rollback to arbitrary existing devices. Deploy its manager only with the first supported device/update pipeline. AI may request approved rollouts; artifact signing, canary cohorts and rollback policy must remain explicit. [Project](https://github.com/antmicro/rdfm), [Linux integration requirements](https://antmicro.github.io/rdfm/rdfm_linux_device_client.html).

### EdgeX Foundry

Local device services, metadata, readings and command APIs form a protocol/normalization layer for sensors and industrial devices. It can complement OpenRemote if EdgeX owns selected protocol interfaces and OpenRemote owns asset context and rules. Assign each device adapter and control rule one owner to avoid duplicate polling and conflicting commands.

EdgeX 4.1 documents NATS Core/JetStream support, but Go services need the include_nats_messaging build tag and C SDK devices do not support that bus. Its sample NATS configuration is not a secure OMOS deployment template. Start with a protected local EdgeX bus and an application-service bridge to OMOS NATS, or maintain a reviewed custom build with scoped NATS identities. Do not expose its internal topics as an unrestricted shared control bus. [Current documentation](https://docs.edgexfoundry.org/4.1/), [NATS integration limitations](https://docs.edgexfoundry.org/4.1/microservices/general/messagebus/).

## OpenRemote and ThingsBoard

OpenRemote is a strong fit for asset context, site rules and independent local operation. It can publish normalized events to NATS through an adapter; NATS remains the transport, not the asset UI/registry. Avoid duplicating all history across its database and ClickHouse without an explicit retention/export design.

ThingsBoard remains credible when productized device-fleet operations and central horizontal scaling dominate. Current platform licensing and offline-key requirements are material to autonomous-site design; see the separately researched [IoT comparison](iot-platform-comparison.md). Do not dismiss it solely due to today's deployment size, but do not choose it for hypothetical device count alone.

## Notion, ClickUp and Todoist replacements

| Product | Role and AI interface | Self-hosting caveat | Recommendation |
|---|---|---|---|
| Outline | Operational wiki, documents, REST API and documented self-hosted MCP endpoint; OIDC supports Authentik | BSL, not unrestricted OSI open source; understand hosting/resale restrictions | Best documented AI-readable/editable operational knowledge option among these |
| Plane | Projects, work items, cycles and REST APIs | Community AGPL; commercial/airgapped editions differ in features and licensing | Start with Community for project/work management; validate required SSO/API features by edition |
| AFFiNE | Documents plus visual whiteboard; local desktop workspaces | Current server is source available under EE terms; published self-host offer has seat limits and offline licensing through sales; official path is Compose | Optional visual workspace; do not assume its entire backend is MIT or independently operable without license planning |
| AppFlowy | Notion-style client/workspace experience | AppFlowy-Cloud archived September 11, 2026; maintained server builds now use commercial self-hosted codebase | Hold as core infrastructure until current terms, API, offline operation and support are evaluated |
| Vikunja | Lists, Kanban, Gantt, task/project management and Todoist import | AGPL; validate AI API authentication and real recurring/reminder workflows | Best dedicated personal-task candidate, not automatically needed if Plane covers them |

Do not assume any existing installation or migration was verified: these four named knowledge/project products were not found as workload deployments in the current hub inspection. They may exist elsewhere. Choose authoritative documents and tasks before migrating; preserve originals and test imports, attachments, permissions and exports.

Sources: [Outline MCP](https://docs.getoutline.com/s/guide/doc/mcp-6j9jtENNKL), [Outline OIDC](https://docs.getoutline.com/s/hosting/doc/oidc-8CPBm6uC0I), [Outline license restrictions](https://docs.getoutline.com/s/hosting/doc/license-restrictions-f9aq6uEL3H), [Plane editions](https://developers.plane.so/self-hosting/editions-and-versions), [Plane API](https://developers.plane.so/api-reference/introduction), [AFFiNE self-host terms and operation](https://affine.pro/self-host), [AppFlowy server transition](https://github.com/AppFlowy-IO/AppFlowy-Cloud/blob/main/README.md), [Vikunja](https://vikunja.io/).

## Bookmarks and research context

Karakeep is the preferred candidate: API/CLI access, archives, text/search, highlights and optional local Ollama tagging/summarization are directly useful to AI context. Linkwarden is a viable alternative for collaborative link preservation; compare capture fidelity and API coverage using actual saved content. Deploy one authoritative bookmark store. Crawled content is untrusted input: treat it as evidence, not instructions to an agent; retain source URL and capture time. [Karakeep](https://docs.karakeep.app/), [local inference](https://docs.karakeep.app/configuration/different-ai-providers/), [Linkwarden deployment](https://docs.linkwarden.app/self-hosting/setup).

## Implementation order

1. n8n infrastructure deployed and owner enrolled for zeus@zeusdelacruz.com; Argo Synced/Healthy and local readiness HTTP 200 verified. Private ingress provisioned. Finish HTTPS verification, scoped integrations and external encrypted recovery backups.
2. Select knowledge/task/bookmark authorities using real imports and AI read/write operations. Start Outline + Plane + Karakeep evaluation; add Vikunja only for a distinct personal task role.
3. Map one physical site's hardware, protocols, power/network failure modes and essential rules. Pilot OpenRemote and the necessary EdgeX device services.
4. Build a first RDFM artifact/canary/rollback pipeline on supported hardware.
5. Compare NetBird and Nebula disconnected/recovery operation before migrating production mesh access.
6. Choose independent site K3s versus validated KubeEdge using compatibility and autonomy requirements. openBalena remains an alternative hardware operating model.
7. Capacity plan from measured flows and recovery requirements. Document identities, APIs, data authority, backups and Git ownership for every capability.
