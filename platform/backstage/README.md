# OneMind Backstage

Private preview: https://onemind-preview.tail717ba0.ts.net
Identity provider: https://onemind-auth.tail717ba0.ts.net

This application is built from upstream Backstage 1.55.3 using the new frontend and backend systems. Source is `services/onemind-portal`. Production catalog files are mounted from the generated ConfigMap in this directory; reviewed inventory changes do not require rebuilding the image. Argo CD manages this directory after merge. Argo uses the existing authenticated repository URL (the old omos-infra alias redirects to omos). Container images are pinned by digest; application changes require a successful image build and a reviewed digest update. The previous Homepage application was retired after browser review. Its source was removed; Git history retains the pre-retirement configuration for recovery.

## Identity and access

Tailscale provides private HTTPS access. Authentik OIDC supplies Backstage identity and membership. The existing `akadmin` account is registered as the initial catalog user and sole allowed Authentik application user. Backstage maps the signed stable subject to `onemind.io/oidc-sub` in the user catalog; it does not trust unverified email as identity. The `onemind-admins` group permits configured Backstage actions; other identified users receive read permissions. No guest login is enabled. Adding a user requires a catalog annotation `onemind.io/oidc-sub` matching the signed subject, Authentik application access, and deliberate group membership. An optional Authentik org sync (off by default; [docs/catalog/authentik-org-sync.md](../../docs/catalog/authentik-org-sync.md)) can publish the members of chosen Authentik groups as catalog users with that annotation; users and groups in `catalog/org.yaml`, including `onemind-admins`, stay Git-owned.

The Kubernetes service account has observation permissions, including pod logs, but cannot read Secrets or modify workloads. Portal actions do not receive a GitHub credential. MCP exposes globally filtered read-only catalog and operational actions. OAuth CIMD uses existing user identity. A dedicated machine token is sealed in backstage-mcp, restricted to the mcp-actions, catalog and OneMind plugins. Global actions filtering excludes mutation tools; a catalog middleware rejects direct machine writes. Permission-restricted service credentials cannot be delegated by this Backstage version; it cannot call scaffolder or cloud summaries. The initial token is stored locally in ~/.config/onemind/backstage-mcp.token (mode 600), never in chat. Use OAuth where supported; rotate the fallback through a reviewed SealedSecret update.

Cluster OIDC requests use the existing Tailscale egress ProxyGroup. CoreDNS rewrites only the identity provider hostname to the operator-managed egress Service while TLS still validates the original hostname. The operator owns that Service's ExternalName, which Argo ignores. CoreDNS must reload newly added custom files; verify resolution and HTTPS after a change.

## Database and recovery

PostgreSQL 17 uses a 10 GiB retained DigitalOcean volume. The application role can create plugin databases but is not a superuser. Credentials are SealedSecrets. Two application replicas use rolling updates with zero unavailable replicas and one surge replica. Preferred node separation and a disruption budget preserve application availability during routine updates and voluntary eviction. Review schema migration compatibility before each upgrade. PostgreSQL and identity remain single-instance dependencies; this is not full high availability.

The old `backstage-backup` CronJob is suspended, superseded by the managed
`backup-postgres-backstage` job at 06:40 UTC. Its current encrypted prefix is
`backups/postgres-backstage/` and it verifies restore of all owned plugin databases.
The older recovery mechanism described below is historical: `backstage-backup` ran daily at 07:15 UTC. It dumps all databases owned by the application role, encrypts the compressed archive with AES-256-GCM, and uploads to the existing private Spaces bucket under `backups/backstage/`. Each run downloads the object, checks authenticated decryption, restores the catalog to a uniquely named temporary database, verifies tables, and removes only that temporary database. Job failure means backup or recovery verification failed. Objects are retained; no automatic deletion policy is configured.

Recovery requires the backup Secret's encryption key, Spaces access, database credentials, and access to the SealedSecrets controller's sealing key when rebuilding the cluster. Keep an independent secure copy of recovery keys; Git alone cannot decrypt a SealedSecret on a new controller. Do not rotate or remove the encryption key without preserving access to prior archives.

For a recovery exercise, create a Job from the CronJob and inspect its `restoreVerified` output. For full recovery, stop Backstage, restore the selected archive into a separate PostgreSQL instance using the same application role and each original database name, verify all databases and catalog entities, then point the application at that instance. Do not restore over a running database. The backup includes database state; application code, catalog source and deployment manifests live in Git. TechDocs is packaged in the pinned application image. See the resilience runbook and recovery exercise for destination prerequisites.

## Catalog and live state

Home reads shared app links from registered Catalog Components with an Open application link or launcher URL annotation, refreshing each minute. Keep Git-owned catalog changes in both apps.yaml copies for mounted entries, or use the shared-service template and register the merged main-branch descriptor URL. The bundled directory JSON supplies legacy category compatibility and offline export defaults; adding a registered shared service no longer requires a frontend rebuild. Personal shortcuts use account-backed UserSettings and remain separate from shared inventory. Cluster catalog identities are an explicit reviewed snapshot from `scripts/snapshot-cluster.cjs`, not automatic discovery. Stable hash suffixes prevent link/workload identity collisions. Rerun and review the snapshot when adding or removing workloads. The snapshot script updates both development and production copies. user enrollment must update both org.yaml copies. The operations panel queries the live cluster every minute and marks unavailable sources and stale results. Existing unrelated unhealthy workloads are shown, not hidden.

## Release checks

Use the official Backstage Kubernetes, Docker, OIDC, permissions and Kubernetes integration documentation. Require: immutable dependency install, TypeScript check, backend/container build, server-side manifest validation, healthy rollout, verified OIDC callback and identity, authenticated catalog and live status, denied anonymous APIs, read-only cluster permissions, and an encrypted backup with successful restore verification. Verify actual browser rendering before replacing Homepage.

Do not claim offline operation: self-hosted pages still require reachability to the host. Existing cloud app links require their providers. Home hosting can use the same containers, Git catalog and restored PostgreSQL with a suitable storage class and mesh routing.

## Inventory maintenance

Run `node scripts/enrich-inventory.cjs` from `services/onemind-portal` on the authorized Mac to refresh both catalog copies. It checks the cluster context, captures workload identities and Tailscale endpoint names/OS/DNS, and connects app links to workloads only through matching ingress hosts and service selectors. It does not collect credentials, user profiles or IP addresses. Review the Git diff before merging; the script is not scheduled automatically. Plain `snapshot-cluster.cjs` replaces cluster.yaml and omits mesh enrichment, so use the enrichment command for the complete inventory.

Hosting annotations distinguish verified cluster routing from unverified mesh endpoints and external providers. Device location and purpose are explicitly unreviewed; an offline discovery snapshot is not a service incident. The catalog graph exposes verified application-to-workload and workload-to-cluster relationships. Cloud API integrations, repository discovery beyond the infrastructure repository, and physical device classification still require reviewed sources.

## Native Kubernetes and Argo CD visibility

The installed native plugins are discovered by `app.packages: all`. On a workload
entity, **Deployment → Kubernetes** shows selector-matched live resources;
**Deployment → Deployment Summary** shows the explicitly named Argo Application.
The Overview Links card also links directly to these views and the owning runbook.
Backstage, its PostgreSQL dependency and the NATS broker are the first reviewed
mappings; this is not coverage of every workload in their namespaces.

Controller **metadata.labels** must match the entity selector, not only the pod
labels. Otherwise the native Kubernetes view can fetch pods but omit their
Deployment/StatefulSet grouping. Keep existing immutable selectors and pod labels.
`services/onemind-portal/scripts/curate-dependencies.cjs` reapplies these exact
identity/namespace/selector mappings after full inventory refresh and rejects
identity drift. Do not infer Argo ownership from namespace co-location. In
particular the NATS dashboard is not the broker StatefulSet.

Catalog-only changes need no image rebuild, but the generated catalog ConfigMap's
new hash changes the portal volume reference and naturally rolls the same pinned
image. Review the Kustomize render first; preserve the database template/PVCs,
read-only RBAC and normal Argo reconciliation. Plugin visibility is not app-level
health, message-delivery verification or new AI mutation authority.

## Curated Fabric registry

The reviewed Notion snapshot (`catalog/fabric-source.json`, 174 records, reviewed 2026-10-06) enriches existing app and mesh entries and adds documented capabilities, planned items and references. `scripts/reconcile-fabric.cjs` merges duplicate titles and Google Calendar/Drive/Photos aliases while retaining every source page link. Recorded status and hosting are provenance annotations, not live health or verified deployment. Original protocol addresses and old IPs remain in source documentation rather than becoming launch buttons. A source snapshot update requires review. The complete enrichment command reapplies this reconciliation so cluster refresh cannot erase curated records.

## P2 completed inventory scope

The complete inventory refresh also records accessible OneMind repositories and exact runtime identifiers, then verifies reviewed aliases. Tags `deployment-verified`, `repository-verified`, `documented-unverified`, `planned`, and `reference` are available through Catalog's Tags filter. A repository match is not evidence of deployment; a deployment match is not live health. Device roles remain documented unless independently inspected. Unknown SaaS account state stays unknown.

All 174 Fabric source pages were fetched. 165 had no body content; nine contained historical notes. Git-owned Markdown is in `docs/inventory/fabric/` and the completeness manifest is `docs/inventory/migration-manifest.json`. Catalog source links now point to these archives rather than Notion. This removes the catalog's Notion content dependency; it is not a full export of the rest of the Notion workspace, attachments, comments or unrelated pages. Notes retain historical dates and are not validated deployment recipes.

## P3 health scope

The authenticated operations endpoint now checks four fixed internal service targets: Authentik readiness, Grafana health, Headlamp web response and NATS dashboard web response. Checks run on demand with a shared 60-second cache and in-flight deduplication; each GET has a five-second absolute deadline and two-megabyte response cap, no redirects or client-supplied target URLs. This is cluster-side point-in-time availability, not uptime history, browser login, tailnet path verification, NATS message delivery or TAK certificate validation. Non-200 results are down except 401/403 which are explicitly authentication-required. Transport failures are unknown.

The same server reads fixed Prometheus queries for five namespaces: backstage, auth, nats, takserver and onetak. CPU is a five-minute rate in cores; memory is working-set bytes. Empty series are unknown, not zero. No monitoring credentials reach the browser and no cluster write permissions were added. The panel exposes timestamps and investigation/catalog links. External monitoring and continuous alert routing can be added separately when needed.

## P4 TechDocs

Operations documentation lives in `services/onemind-portal/runbooks/`. CI generates MkDocs using pinned TechDocs CLI 1.12.0 and mkdocs-techdocs-core 1.7.1, checks metadata and search artifacts, and packages output at `/app/techdocs/default/component/onemind-runbooks`. The existing external builder and local filesystem publisher serve this immutable release content from the self-hosted portal. A failed docs build blocks the image release. Docs changes require a new image digest; Git and prior images provide recovery, and no cloud publishing credentials are added.

TechDocs search uses the existing PostgreSQL search engine and collates every minute. The catalog documentation entity appears in Docs; selected critical-service entries have runbook links. Documentation is versioned with the portal, does not require Notion, and requires host reachability rather than promising offline browser access. These runbooks distinguish verified configuration from untested data recovery and historical Fabric notes.

## Native Events bus across replicas

`@backstage/plugin-events-backend` is registered at the exact `0.6.6` version
from the Backstage `1.55.3` release manifest. It provides the native event bus
used by Notifications, Signals and catalog subscribers. Production inherits
`backend.database.client: pg`; the maintained backend selects its PostgreSQL
store, migrates `backstage_plugin_events`, and uses PostgreSQL LISTEN/NOTIFY
for delivery across the two portal replicas. The existing bounded database
pool remains in effect. Observe connection headroom during rollout; this
adds a plugin database and event listeners.

`events.http.topics: []` deliberately enables no external webhook ingresses.
The authenticated internal event bus routes remain available to native
plugins. Existing machine identity restrictions do not grant Events access.
Before promotion, verify the normal encrypted backup and restore evidence.
After rollout, confirm `events` initialization and PostgreSQL store selection
on both replicas, no Events subscription 404s, and unchanged anonymous and
machine access denials. This supplies internal delivery; it does not activate
external agent completion producers or provider OAuth grants.

References: [Events backend installation](https://backstage.io/api/stable/modules/_backstage_plugin-events-backend.html),
[1.55.3 release manifest](https://versions.backstage.io/v1/releases/1.55.3/manifest.json),
[versioned PostgreSQL event-store selection](https://github.com/backstage/backstage/blob/v1.55.3/plugins/events-backend/src/service/hub/createEventBusRouter.ts).

Production sets `events.useEventBus: always`: the installed Events service retries temporary startup 404s instead of permanently choosing local-only delivery before the bus router is ready. Verify subscriptions on both replicas after startup settles.

## My Active Stack

The home page places purpose cards above Favorites for Notes, Email, Passwords, Files, Maps, Photos, Voice memos, Wallet, Contacts, Calendar, and SMS. Each signed-in account saves its own assignments through native User Settings. Choices begin unassigned; catalog lifecycle and favorites do not automatically select tools. Assign one primary app, multiple supporting apps, and alternatives from catalog entries or personal shortcuts, or record a native app by name without a launch link.

Each decision records roles, use, data/export location, migration follow-up, timestamp and reason; the latest 50 decisions per purpose are retained. Directory references remain stable and removed entries retain their saved names. Availability remains in Today. Changing an assignment does not move data or uninstall software. No shared defaults are preselected.

## Hosting labels

Directory cards, Favorites and My Active Stack show catalog hosting labels. Hosting filters apply to the application directory. DigitalOcean NYC3 inventory is labelled Self-hosted · cloud, distinguishing OneMind-managed cloud infrastructure from provider-operated Cloud service entries. External resources show Website reference; planned entries show Planned; native and firmware entries retain their type; mesh endpoints and absent metadata remain unverified or unknown. These labels are recorded inventory, not current health or proof of active use. Today retains availability checks.

## Access & Deployments

App cards and My Active Stack expose an Access & Deployments editor, also available as a catalog entity tab. Supported platforms are separate from actual access points. Points record platform, installation method, location, device/cluster entity reference, safe HTTPS link or native instructions, evidence and observation date. No native or remote execution is performed.

Shared catalog points are versioned in the `onemind.io/access-points` JSON annotation. Initial points use existing verified Kubernetes route relationships and read-only inspection of Proton Mail, Dropbox and 1Password on the MacBook identified by its Tailscale DNS identity. These are dated observations, not continuous installation or health checks. `onemind.io/supported-platforms` optionally records a comma-separated platform list without implying installation.

Account-backed User Settings store additional points and a preferred access point per app; preferences apply in Favorites, directory cards and My Active Stack. Manual points are explicitly User recorded. Device selection reads existing mesh/cluster catalog resources, exposes their OS, DNS identity and discovery date, and does not discover installed apps on remote devices. Unavailable catalog inventory does not erase saved points. Existing workload dependencies link to catalog records for inspection. The enrichment script retains mesh OS on future inventory refreshes.

## AI directory organization

Local AI contains model runtimes and desktop tools. AI Model Providers contains
external hosted inference catalogs and APIs, with Cloud service labels and
unverified account/integration status. Open Model Families contains evaluation
candidates and official model cards. Family entries record candidate releases,
release-specific licenses, and proposed OneMind roles; these are not benchmark
results, active model assignments, or verified deployments. Open weights can be
served locally or by a cloud provider; license and hosting are separate facts.
Runtime compatibility and hardware capacity must be checked before installing a
particular release. No account credentials or inference endpoints are configured
by adding a directory entry.

## Phase 6: one portal (history) and separation from OpenChoreo

In Phase 6 the OpenChoreo 1.3.1 core/observability Backstage plugins ran inside this
host. They have been **removed again**: the portal no longer contains any `@openchoreo/*`
package, OpenChoreo catalog sync, entity group/observability views, OpenChoreo secrets
or `openchoreo.*` configuration. The portal stands alone:

- Sign-in is Backstage's own `oidc` provider against the Authentik application
  `onemind-backstage` (redirect URI `<portal>/api/auth/oidc/handler/frame`), mapping the
  signed subject to `onemind.io/oidc-sub` and the existing group. No OpenChoreo token is
  obtained.
- Permissions are one plain OneMind policy (see [Permission policy](#permission-policy)).
- Catalog entities that described OpenChoreo-managed resources (projects, environments,
  planes, the `stac-browser` pilot) are no longer synced into the catalog. Static catalog
  entries for cluster workloads, including the OpenChoreo applications themselves, are
  unchanged; tearing those down is a separate step.
- STAC portal integration is removed from source together with the former backend
  observations and catalog module. This ships with the next portal release; this
  change does not publish an image or deploy the portal. The shared bounded HTTP
  observer remains available for Prometheus alerts.

The historical Phase 6 record is
[`docs/operations/phase6-portal-integration.md`](../../docs/operations/phase6-portal-integration.md).

## Permission policy

One plain policy (`plugins/plugin-onemind-permissions-backend/src/policy.ts`) answers every
Backstage permission, including `catalog.entity.*` and `scaffolder.task.create`. It needs no
OpenChoreo service or configuration.

| Caller | Decision |
| --- | --- |
| No user identity | DENY, including reads |
| Member of `group:default/onemind-admins` | ALLOW everything |
| Any other signed-in user, `read` actions | ALLOW |
| Any other signed-in user, create/update/delete (entity create/refresh/delete, location writes, `scaffolder.task.create`, ...) | DENY |

Decisions are never conditional, so there is no per-entity filtering. Machine callers (the MCP
token, plugin-to-plugin calls) are service principals limited by `backend.auth.externalAccess`,
not by this policy.

What the previous policy enforced, for the record: the OneMind rule above was an upper bound,
then `catalog.entity.*` and `scaffolder.task.create` were delegated to the OpenChoreo policy.
That policy let a user act on catalog entities that carry no `openchoreo.io/namespace`
annotation (all OneMind entities) and gated OpenChoreo-managed entities by the user's
OpenChoreo capabilities; `scaffolder.task.create` and `catalog.entity.create` needed any
OpenChoreo create capability. Only the operator (`akadmin`) had an OpenChoreo binding
(`operator-admin-binding`, admin), so in practice admin had full access and everyone else was
read-only on non-OpenChoreo entities. Differences now: OpenChoreo-managed entities are no longer in the portal catalog, and an
administrator no longer needs an OpenChoreo role. The `onemind_authz_cache` table that cached
OpenChoreo capabilities is no longer used and is left in the permission database; drop it by
hand when convenient.
