# OneMind Backstage

Private preview: https://onemind-preview.tail717ba0.ts.net
Prepared alias: https://portal.tail717ba0.ts.net (available after the alias rollout; redirects to the preview host; see [Hostnames](#hostnames))
Identity provider: https://onemind-auth.tail717ba0.ts.net

This application is built from upstream Backstage 1.55.3 using the new frontend and backend systems. Source is `services/onemind-portal`. Production catalog files are mounted from the generated ConfigMap in this directory; reviewed inventory changes do not require rebuilding the image. Argo CD manages this directory after merge. Argo uses the existing authenticated repository URL (the old omos-infra alias redirects to omos). Container images are pinned by digest; application changes require a successful image build and a reviewed digest update. The previous Homepage application was retired after browser review. Its source was removed; Git history retains the pre-retirement configuration for recovery.

## Identity and access

Tailscale provides private HTTPS access. Authentik OIDC supplies Backstage identity and membership. The existing `akadmin` account is registered as the initial catalog user and sole allowed Authentik application user. Backstage maps the signed stable subject to `onemind.io/oidc-sub` in the user catalog; it does not trust unverified email as identity. The `onemind-admins` group permits configured Backstage actions; other identified users receive read permissions. No guest login is enabled. Adding a user requires a catalog annotation `onemind.io/oidc-sub` matching the signed subject, Authentik application access, and deliberate group membership. An optional Authentik org sync (off by default; [docs/catalog/authentik-org-sync.md](../../docs/catalog/authentik-org-sync.md)) can publish the members of chosen Authentik groups as catalog users with that annotation; users and groups in `catalog/org.yaml`, including `onemind-admins`, stay Git-owned.

The Kubernetes service account has observation permissions, including pod logs, but cannot read Secrets or modify workloads. Portal actions do not receive a GitHub credential. MCP exposes globally filtered read-only catalog and operational actions. OAuth CIMD uses existing user identity. A dedicated machine token is sealed in backstage-mcp, restricted to the mcp-actions, catalog and OneMind plugins. Global actions filtering excludes mutation tools; a catalog middleware rejects direct machine writes. Permission-restricted service credentials cannot be delegated by this Backstage version; it cannot call scaffolder or cloud summaries. The initial token is stored locally in ~/.config/onemind/backstage-mcp.token (mode 600), never in chat. Use OAuth where supported; rotate the fallback through a reviewed SealedSecret update.

Cluster OIDC requests use the existing Tailscale egress ProxyGroup. CoreDNS rewrites only the identity provider hostname to the operator-managed egress Service while TLS still validates the original hostname. The operator owns that Service's ExternalName, which Argo ignores. CoreDNS must reload newly added custom files; verify resolution and HTTPS after a change.

## Hostnames

Backstage serves one origin, `PORTAL_URL` in `deployment.yaml`. It sets `app.baseUrl`, `backend.baseUrl` and the CORS origin, and the auth backend builds the OIDC callback from it. Each of these takes a single value. A page loaded from a second host would call the API at `backend.baseUrl` cross-origin, which the CORS origin and CSP `connect-src 'self'` block, so a second host cannot serve the portal itself.

`portal` is a tailnet alias (`portal-redirect.yaml`). Its Tailscale Ingress reaches
a small nginx that answers every request except `/healthz` with a 302 to the same
path and query on `PORTAL_URL`; `/healthz` returns 200 for probes. Only that Ingress
proxy can reach it, and it has no egress. Use the canonical host for API, MCP and
OAuth clients. `scripts/tests/test_gitops_portal_alias.py` fails if the redirect
target and `PORTAL_URL` drift apart.

To make `portal` canonical in a separate change, add its
`https://portal.tail717ba0.ts.net/api/auth/oidc/handler/frame` callback to the
Authentik `onemind-backstage` provider first. Then update `PORTAL_URL`, route the
`portal` Ingress to the `backstage` Service, and move the redirect to
`onemind-preview`, targeting `portal`. Update the alias contracts, the provider
bootstrap payload and catalog links in that same change. The current temporary
redirect avoids a permanent redirect back to the old canonical hostname.

**Merge and rollout checks:** confirm that no existing tailnet device or Tailscale
Ingress claims `portal`. This redirect keeps Backstage's origin and OIDC callback
on `onemind-preview`; it needs no Authentik provider edit. An additional callback
is a prerequisite only for the separate future canonical-host migration described
above. After Argo reconciliation, check that the alias preserves a harmless path
and query in its redirect, the destination still supports sign-in, and the
operator registered the exact `portal` hostname rather than a suffixed name.
Keep the existing canonical hostname if the alias cannot be registered.

## Database and recovery

PostgreSQL 17 uses a 10 GiB retained DigitalOcean volume. The application role can create plugin databases but is not a superuser. Credentials are SealedSecrets. Two application replicas use rolling updates with zero unavailable replicas and one surge replica. Replicas are required to run on the autoscaled `onemind-ecosystem-cpu` pool (not the memory-saturated hub-core pool); preferred node separation spreads them when that pool has two or more nodes, and a disruption budget preserves application availability during routine updates and voluntary eviction. With a single ecosystem node both replicas share it until the autoscaler adds one. Review schema migration compatibility before each upgrade. PostgreSQL and identity remain single-instance dependencies; this is not full high availability.

Backups come from the managed `backup-postgres-backstage` job (`platform/backups`) at 06:40 UTC.
Its encrypted prefix is `backups/postgres-backstage/` and it verifies restore of all owned
plugin databases. The legacy `backstage-backup` CronJob was suspended on 2026-10-08 and has
been removed from Git. Its `backstage-backup` Secret, image source (`backup-image/`, built by
`build-backstage-backup.yml`) and the `recovery/` exercise are kept: they still decrypt the
older archives, and `platform/backups/tools/mirror.py` reads the Spaces credentials from that Secret.
The description below is historical: `backstage-backup` ran daily at 07:15 UTC. It dumps all databases owned by the application role, encrypts the compressed archive with AES-256-GCM, and uploads to the existing private Spaces bucket under `backups/backstage/`. Each run downloads the object, checks authenticated decryption, restores the catalog to a uniquely named temporary database, verifies tables, and removes only that temporary database. Job failure means backup or recovery verification failed. Objects are retained; no automatic deletion policy is configured.

Recovery requires the backup Secret's encryption key, Spaces access, database credentials, and access to the SealedSecrets controller's sealing key when rebuilding the cluster. Keep an independent secure copy of recovery keys; Git alone cannot decrypt a SealedSecret on a new controller. Do not rotate or remove the encryption key without preserving access to prior archives.

For a recovery exercise on a legacy archive, follow `recovery/README.md` and inspect its `restoreVerified` output. For full recovery, stop Backstage, restore the selected archive into a separate PostgreSQL instance using the same application role and each original database name, verify all databases and catalog entities, then point the application at that instance. Do not restore over a running database. The backup includes database state; application code, catalog source and deployment manifests live in Git. TechDocs is packaged in the pinned application image. See the resilience runbook and recovery exercise for destination prerequisites.

## Catalog and live state

Home reads shared app links from registered Catalog Components with an Open application link or launcher URL annotation, refreshing each minute. Keep Git-owned catalog changes in both apps.yaml copies for mounted entries, or use the shared-service template and register the merged main-branch descriptor URL. The bundled directory JSON supplies legacy category compatibility and offline export defaults; adding a registered shared service no longer requires a frontend rebuild. Personal shortcuts use account-backed UserSettings and remain separate from shared inventory. Cluster catalog identities are an explicit reviewed snapshot from `scripts/snapshot-cluster.cjs`, not automatic discovery. Stable hash suffixes prevent link/workload identity collisions. Rerun and review the snapshot when adding or removing workloads. The snapshot script updates both development and production copies. user enrollment must update both org.yaml copies. The operations panel queries the live cluster every minute and marks unavailable sources and stale results. Existing unrelated unhealthy workloads are shown, not hidden.

## Plugin ecosystem catalog

`catalog/plugins.yaml` (generated, ConfigMap `backstage-catalog-plugins` mounted at `/app/catalog-plugins`, registered in `app-config.cluster.yaml`) lists every OneMind plugin, adapter and bridge under the System `onemind-plugin-ecosystem`. Edit `catalog/plugins/inventory.yaml` and run `python scripts/generate_plugin_catalog.py`; see [catalog/plugins/README.md](catalog/plugins/README.md).

## GitHub discovery

`@backstage/plugin-catalog-backend-module-github` is pinned at `0.14.0`, the version in
the Backstage `1.55.3` release manifest. `catalog.providers.github.onemindos` in
`app-config.cluster.yaml` scans the `onemindos` organization every 30 minutes, keeps only
`omos` (`filters.repository: '^omos$'`) and registers
`https://github.com/onemindos/omos/blob/main/catalog-info.yaml` only when the `main`
branch contains that file (`validateLocationsExist`). Forks and archived repositories stay
excluded, and `catalog.rules` applies to everything discovered (no User, Group or
Template). Keep the path literal: a `**/catalog-info.yaml` wildcard would also ingest the
Scaffolder skeletons under `services/onemind-portal/templates/*/skeleton/`.

The provider authenticates with the existing organization-owned read-only GitHub App
(`integrations.github`, sealed in `backstage-plugin-integrations`). No credential, webhook
ingress or event producer is added: a root descriptor merged to `main` is registered on the
next scan, and edits to registered descriptors follow normal catalog processing. The App
installation selects only `omos`
([installation record](../../docs/backstage-roadmap/COMMUNITY-PLUGIN-INSTALLATION.md)),
and that selection is the real limit. The provider's organization listing also returns the
other public `onemindos` repositories, but `@backstage/integration` refuses to read a file
from any repository outside the installation's selection (no `publicAccess` is configured).
A descriptor there would only produce a Location that fails on every processing cycle, so
the filter leaves those repositories out. To add one, select it in the App installation and
widen `filters.repository` in the same change. Anyone who can merge to `main` in `omos` can
add entities within `catalog.rules`, including Locations that point at other URLs. On
2026-10-27 `omos` had no root `catalog-info.yaml` on `main`, so the provider registers
nothing until one is merged. The module also lets Catalog Import analyze GitHub repository
URLs; registration still needs `catalog.location.create` (administrators only).

The block takes effect only once a portal image built from this source is pinned in
`deployment.yaml`; older images ignore it. Rollback: delete `catalog.providers.github`
and let Argo sync (no image build), or drop the module in the next image.

The repository's root `catalog-info.yaml` is a Location listing descriptors kept next to
service source, only for deployed services the static catalog does not define. Today that is
`onemind-data` (Resource: the CloudNativePG cluster that runs the `services/onemind-postgres`
image). Services the static catalog defines (portal, core, TAK Server, TAK Portal,
WebODM-to-STAC sync) stay static. Backstage keeps the first definition of an entity ref and
logs a conflict for any other; it never merges them. To move a service to its descriptor,
remove the static entry and the generator input that recreates it in the change that lists
the descriptor. `services/onemind-portal/catalog-info.yaml` is ready for that and
deliberately unlisted. `scripts/tests/test_gitops_portal_discovery.py` rejects overlaps and
unknown owners, systems or relation targets.

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

## OpenCost allocation proxy

`/api/proxy/opencost/allocation` forwards GET requests to the in-cluster OpenCost allocation API (`opencost.opencost.svc.cluster.local:9003/allocation`). Backstage credentials are required and are not forwarded; other OpenCost routes are not reachable through it. It returns raw OpenCost estimates. The Cost page (see [Plugin pack](#plugin-pack-f9)) is its only consumer; `scripts/cost-allocation.cjs` defines the validation (`docs/operations/backstage-cost-workstream.md`). The OpenCost catalog component links to the tailnet UI.

## Plugin pack (F9)

Five additions, one commit each so any can be reverted alone. Exact versions are the ones in the
Backstage 1.55.3 release manifest unless stated. Work order: [F9](../../docs/roadmap/flagship/F9-plugin-pack.md).

| Addition | What a person sees | Package | Boundary |
| --- | --- | --- | --- |
| Cost page | Sidebar Insights > Cost (`/cost`): last 1, 7 or 30 complete UTC days of OpenCost cost per namespace, app and pillar | internal `plugin-onemind-cost` (no new dependency) | Signed-in users; reads only the existing proxy route |
| Shortcuts | "Add shortcut" list in the desktop sidebar | `@backstage-community/plugin-shortcuts` 0.25.1 | Stored per user in the account-backed storage API |
| DevTools | Admin > DevTools, Info tab (host, Node, Backstage and package versions) | `@backstage/plugin-devtools` 0.1.43, `-backend` 0.5.21, `-common` 0.1.26 | Administrators only; raw config and scheduler views are closed |
| Unprocessed entities | Admin > Unprocessed entities: failed and pending catalog entities with errors | `@backstage/plugin-catalog-unprocessed-entities` 0.2.35, `-common` 0.0.17, `@backstage/plugin-catalog-backend-module-unprocessed` 0.6.16 | Administrators only |
| Kyverno policy results | "Policy results (Kyverno)" card on app and workload pages | internal (`operations-backend` route, app card) | Signed-in users; read-only |

**Cost page.** The response is validated with a TypeScript port of `scripts/cost-allocation.cjs`; `allocation.test.ts` runs both on the same payloads and requires identical results, so change them together. Unknown, partial, empty and failed results are labelled and never shown as zero; a namespace with an unusable cost shows "Unknown", and totals become "at least". OpenCost always adds an `__unmounted__` line (load balancers are most of it today); it is shown as "Not owned by a namespace" instead of being dropped. Namespaces are joined to catalog apps through `backstage.io/kubernetes-namespace`; a namespace several apps share is shown as shared and is never split between them, and pillar totals mark assignments that are still `proposed`. The page states that these are estimates and not the DigitalOcean bill. OpenCost does not report a currency, so amounts are labelled USD, which is how DigitalOcean bills. Query: `window=<start>,<end>&aggregate=namespace&accumulate=true` with whole UTC days.

**Shortcuts.** `@backstage/plugin-shortcuts` is deprecated (last release April 2024, moved to community-plugins) and is not in the 1.55.3 manifest. `@backstage-community/plugin-shortcuts` 0.25.1 is the newest release whose dependency ranges the pinned core packages satisfy (0.26.0 would pull a second copy of `core-components`). It has no new-frontend export, so `packages/app/src/shortcutsFeature.ts` converts the legacy plugin to register its API and `plugin-onemind-nav` renders the component. It is hidden on phones, where the bottom bar is already full. This is separate from the shared app links on Home and from Settings.

**DevTools and unprocessed entities.** Every DevTools permission has action `read`, which the default rule would have given to any signed-in user, so the policy lists the administrator-only ones explicitly. `devtools.config` is closed for everyone: the Config view hides only values a plugin's config schema marks `@visibility secret`, so a secret in an unmarked key would be shown. The scheduler view is closed because its trigger would run background jobs. The Config and Scheduled Tasks tabs are also removed in code (`devtoolsFeature.ts`). To re-enable one, remove its entry from `NEVER_ALLOWED_PERMISSIONS` in `policy.ts` and from `devtoolsFeature.ts`, after checking the config schema coverage. The unprocessed-entities plugin ships its page disabled and attaches a copy to DevTools; `unprocessedEntitiesFeature.ts` enables the standalone page and disables the copy. Its delete action stays administrator-only; deleting a row only makes the catalog re-ingest that entity. Negative tests start the real routes behind the real policy (`adminRoutes.test.ts`, `unprocessedRoutes.test.ts`): 403 for a signed-in non-admin and for a service credential such as the MCP token, 403 on `/config` even for an administrator.

**Kyverno policy results.** Options checked on 2026-11-05: Kyverno's Policy Reporter plugin needs the Policy Reporter service as another in-cluster workload (the larger node had about 16% memory requests free); TeraSky `kyverno-policy-reports` 1.21.0 supports the new frontend but needs `core-components` 0.18.15, `frontend-plugin-api` 0.18.2 and `@backstage/ui` 0.19, newer than 1.55.3, and fetches a report per resource UID. Neither is installed. Kyverno's reports controller already writes a namespaced `PolicyReport` per scanned resource (693 on 2026-11-05, no failures), so `GET /api/onemind/policy-reports` in `plugin-onemind-operations-backend` reads `/apis/wgpolicyk8s.io/v1alpha2/policyreports` through the read-only cluster access the Live Operations page uses: signed-in users only, fixed path, at most one cluster list per minute shared by all callers, stale or error reported rather than an all-clear. It resolves each Pod and ReplicaSet to its controller by longest name prefix, so `backstage-postgres` stays apart from `backstage`. The card matches an app's workloads from `onemind.io/kubernetes-workloads` (or the single `onemind.io/kubernetes-workload`), falls back to the whole namespace when none is named, and counts a violation once per workload, policy and rule. Cluster-scoped reports are not read. This needs one RBAC addition: `get`/`list` on `policyreports.wgpolicyk8s.io` in the `backstage-observer` ClusterRole (`rbac.yaml`); `platform-backstage` has automated sync, so it applies when this merges. No Secret access is added. No Trivy Operator work: it runs separately in `trivy-system`.

**Not built.** Bookmarks per entity and the optional Q&A plugin are not part of this pack. Harbor and OpenChoreo are retired and nothing was added for them.

**Bundle size.** See the F9 work order status for the measured change.

## App lifecycle (F8)

Work order: [F8](../../docs/roadmap/flagship/F8-app-lifecycle.md). Built on `feat/backstage-app-lifecycle`; not deployed until a
reviewed image digest update (the templates and backend code are in the image; the sidebar flag `onemind.features.retireApp` is `true`
in the image default).

| Addition | What a person sees | Where |
| --- | --- | --- |
| **Add app** (`add-app`) | Create > "Add app to OneMind": hosting (Kubernetes, DigitalOcean, SaaS, home device), pillar, System, owner, health check, data and backup plan, access. Opens a draft PR. | `services/onemind-portal/templates/app-add/` |
| **Retire app** (`retire-app`) | Create > "Retire app" (also Change > Retire app): pick an app or one workload, see everything attached, choose retention (default: keep the last backup 30 days), open a draft removal PR with a manual teardown checklist. | `templates/app-retire/` |
| **Stack hygiene** (`/stack-hygiene`) | Insights > Stack hygiene: weekly report plus Refresh now. Each finding links to Add or Retire with pre-filled parameters. | `plugin-onemind-services` (page), `plugin-onemind-services-backend/src/lifecycle/` (task) |

**How the PRs are made.** Two scaffolder actions, `onemind:app-add:prepare` and `onemind:app-retire:plan`, read the catalog, the
repository (through the existing GitHub App and URL reader, one download per run) and, for retirement, the cluster (the portal's
existing service account, list calls only). They write file contents into the scaffolder workspace; the template's own
`publish:github:pull-request` step turns that into a draft PR (`filesToDelete` carries the deletions). There is no token field: the
GitHub App credential in `integrations.github` is used. **That App must be allowed to write contents and pull requests on
`onemindos/omos`** (the discovery section above describes it as read-only); until it is, publishing fails with a 403 and nothing is
changed. No step runs `kubectl`, an Argo sync or a DigitalOcean call; a test scans the steps and the lifecycle source for that.

**Add app, Kubernetes.** The PR contains `platform/argocd/platform-NAME.yaml` (manual sync, no `resources-finalizer`), and
`platform/NAME/` with namespace (layer label, Pod Security restricted, `Prune=false,Delete=false`), Deployment (pinned digest,
requests and limits, probes, ecosystem node pool, read-only root), Service, NetworkPolicies, a tailnet or Cloudflare-allowlisted
ingress, a runbook stub and `NAME.secret.example` (key names only, not applied). With a VolSync plan it adds a claim on
`do-block-storage-retain` and `backup.yaml` whose ReplicationSource has a non-empty `trigger.schedule`; the snapshot and cache classes
are `do-block-storage`. It also edits the catalog **inputs**: `workload-map.yaml` (the Argo mapping), `structure.yaml` (app, pillar
override), three preparation exceptions in `coverage-exceptions.yaml` (no live workload yet; delete them after the first sync), the
hard-coded Argo application count in `test-catalog-structure.cjs` and `test-coverage.cjs`, and a `docs/STACK.md` section. Generated
catalog output is not hand-edited: the PR checklist starts with `node scripts/curate-structure.cjs` and the two catalog tests. For
SaaS, DigitalOcean and home-device apps the template appends one Component document to both `apps.yaml` copies and a `structure.yaml`
placement; no Argo files.

**Retire app.** The inventory covers workloads, the Argo application (finalizer and whether the parent prunes), PVCs and bound
volumes, secrets, ingress hosts (Cloudflare and Tailscale names), backup sources, dashboards, alerts and catalog entries. A section the
portal cannot read is `unavailable`, never "none": DigitalOcean snapshots and volumes (F6 is not built), VolumeSnapshots, live Secrets
and ReplicationSources (the `backstage-observer` role has no access, and none was added). The PR deletes the Application file and the
files of its own source directory, removes the app's catalog documents from both copies, its mapping lines, its VolSync
ReplicationSource and coverage exceptions, lowers the Argo count, adds a retirement record under `docs/operations/retirements/`, and
keeps markdown, restic credentials, shared directories and protected platform directories (`platform/argocd`, `security`, `backups`,
and so on), listing each with its reason. Choosing one workload instead of an app removes only manifests named after it. The
checklist encodes the Paperclip, Tempo and KEDA lessons: a finalizer cascades but a `prune: false` parent never removes the
Application, stop traffic and remove the backup source before the namespace (stale VolSync series keep an alert firing), CRDs last,
volumes and snapshots only after the retention date and with owner approval.

**Stack hygiene.** `lifecyclePlugin` (backend plugin `onemind-lifecycle`) runs weekly (Monday 06:00 UTC) and stores the JSON in a
small table in its own plugin database (last 12 kept). `GET /api/onemind-lifecycle/stack-hygiene` needs a signed-in user;
`POST .../refresh` is administrators only (it authorizes `scaffolder.task.create`, which the policy gives to administrators; no new
permission) and is rate-limited to one run a minute. Checks: Argo apps without a catalog entry, catalog apps without a running
workload, workloads with no owner app, no health check (no probe and no `onemind.io/health-check` annotation), a volume with no
VolSync source or backup CronJob, no runbook, images without a digest pin. Orphan volumes and old snapshots are `unavailable` until F6
supplies DigitalOcean data; a check that could not be evaluated says why instead of reporting zero.

**Limits.** The planner trusts the catalog's workload lists; drift shows up in the hygiene report, not in the retire inventory.
Review the first published PRs as drafts.

## P4 TechDocs

Operations documentation lives in `services/onemind-portal/runbooks/`. CI generates MkDocs using pinned TechDocs CLI 1.12.0 and mkdocs-techdocs-core 1.7.1, checks metadata and search artifacts, and packages output at `/app/techdocs/default/component/onemind-runbooks`. The existing external builder and local filesystem publisher serve this immutable release content from the self-hosted portal. A failed docs build blocks the image release. Docs changes require a new image digest; Git and prior images provide recovery, and no cloud publishing credentials are added.

TechDocs search uses the existing PostgreSQL search engine and collates every 15 minutes, plus 30 seconds after each pod start (docs only change with a new image). The catalog documentation entity appears in Docs; selected critical-service entries have runbook links. Documentation is versioned with the portal, does not require Notion, and requires host reachability rather than promising offline browser access. These runbooks distinguish verified configuration from untested data recovery and historical Fabric notes.

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
| `devtools.config`, `devtools.scheduler.read`, `devtools.scheduler.trigger` | DENY for everyone, administrators included |
| Member of `group:default/onemind-admins` | ALLOW everything else |
| Any other signed-in user, `read` actions | ALLOW, except `devtools.administer`, `devtools.info`, `devtools.external-dependencies` and `catalog.entities.unprocessed.read`, which are administrator-only |
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
