# Portal acceptance follow-through — 2026-10-20

Zeus requested the remaining Phase 6 checks and a complete OneTAK/STAC operational
slice. This first change reviews trust/coverage and repairs confirmed launcher
mappings. It does not migrate apps, enable workflows, rotate credentials, change
permissions or provision resources. Phase 6 remains in progress where noted.

## Backup review: suspension is not automatically a backup failure

Read-only cluster observations at approximately 07:56–08:00 UTC:

- `backstage/backstage-backup` is intentionally superseded by
  `backup-postgres-backstage`. The replacement completed at 06:41:16 UTC;
  its backup container reports authenticated archive recovery,
  `restoreVerified: true`, **26 databases and 156 tables**. No restore was run
  against production. Existing backup keys/storage were not touched.
- `clickhouse/backup-clickhouse` is retired in favor of VolSync.
  `volume-ab4fbd87295d`, source PVC `data-storage-chi-onemind-onemind-0-0-0`,
  last synchronized at 00:41:14 UTC. A filesystem synchronization is not proof
  of application-consistent database restoration; do not imply that it is.
- `ops/backup-offsite-copy` has no successful run and remains suspended awaiting
  the independent destination/credentials described in
  [OFFSITE-COPY.md](../../platform/backups/OFFSITE-COPY.md). This is a real recovery
  gap, not something to fix by blindly unsuspending the CronJob.
- Other active logical backup schedules reported successful October 20 runs.
  The existing Today panel still lists suspended historical schedules; this is
  an interpretation/UI gap, not permission to delete them or hide the offsite gap.

## Login and authorization

- At 08:01:14 UTC, a read-only query of the existing permission database found
  the signed-in administrator's shared profile expiration at 08:34:52 UTC.
  Only expiration metadata was read; no token/profile values were extracted.
- Both replicas subsequently served actual `/api/auth/openchoreo-auth/refresh`
  HTTP 200 and `/api/permission/cache-capabilities` HTTP 200. Sampling
  08:01:29–08:07:59 UTC showed 11 successful refreshes and 13 cache initializations.
  Two 401 refresh responses preceded sign-in recovery; they are not erased from
  the evidence. The existing Sign In flow recovered without a credential change.
- The signed-in pilot Deploy page still showed `stac-browser-v5` Active.
  No repeated no-token/no-cache warning appeared in that sampled interval.
- Added actual upstream-service regressions for renewed-token TTL across shared
  cache adapters, actor/scope isolation, expiry after renewal, and provider failure
  not resurrecting an expired profile. Only the upstream network boundary is
  mocked; production policy, TTL and cache code remain unchanged.
- The expiry cases are automated/isolation evidence, **not deliberately expiring
  the production user or live-testing another human role**. Keep this distinction
  if closing a narrower login-refresh gate. Existing anonymous/machine denials
  remain acceptance requirements after any subsequent runtime release.

## Launcher review and repairs

31 first-party catalog launchers were tested from the operator Mac using ordinary
unauthenticated HEAD, GET fallback for 405, standard TLS verification, no followed
redirects and bounded timeouts. This observes this access path, not universal
application health. The first pass found 11 HTTP-200 endpoints, 18 unresolved DNS
names, one timeout and one untrusted certificate chain.

Four replacement destinations were independently confirmed against actual
Ingress/Tailscale inventory and requests:

| Stable catalog identity | Reviewed destination | Result |
|---|---|---|
| `matrix-synapse` | `https://matrix.onemindos.dev` | 302, no redirect followed |
| `element-chat` | `https://element.onemindos.dev` | 200 |
| `agent-monitor` | `https://agent-tap.tail717ba0.ts.net` | 200 |
| `agent-web-ui` | `https://agent-console.tail717ba0.ts.net` | 200 |

The corresponding app-to-workload mappings for the agent tools now point to the
observed `agent-fabric` Deployments. Both app catalog copies are identical. No
entity ID, favorite, unrelated URL, public exposure or incoming reference is
removed. The remaining **16** findings carry the searchable
`launcher-review-required` tag, path-specific review annotations and observation
time. An unresolved Mac path is not labeled a failed Kubernetes workload.

Remaining review identities: cloudtak, tak-server, onetak, fedhub, clickhouse,
minio, postgis, longhorn, prometheus, nats-surveyor, tileserver-gl, valhalla,
maputnik, webodm, home-assistant and openclaw. TAK's certificate verification was
not bypassed. No DNS record, tunnel, network boundary or application was changed
just to make a link return green. Curated external Argo/Tailscale aliases remain
separate stable records; the duplicate external Kubernetes Dashboard now explicitly
references its canonical Component without changing either identity; Bambu Lab/Studio remain distinct products, not falsely
merged because their destinations overlap.

## Layout and intended COP coverage

The processed catalog has 715 entities and 12 populated OpenChoreo-managed kinds.
Browser samples cover Component, Project/System, Environment, Namespace/Domain,
DeploymentPipeline, ClusterWorkflow, ClusterComponentType, ClusterTraitType,
ClusterDataPlane, ClusterResourceType, ClusterObservabilityPlane and
ClusterProjectType. The remaining supported kinds have no live example; the
existing 22 actual upstream layout-selector regressions cover their composition.
The reviewed platform layouts show one Relations graph and one metadata/details
section; an upstream details section need not have the literal heading About.

Unowned platform metadata, unused gateway modes and empty dependency/API cards
remain visible as absent information, not invented state. Runtime/deployment
health and data freshness are separate. The STAC slice must supply its actual
API/dependencies rather than globally hiding empty cards for all applications.

**Correction to the earlier acceptance note:** the deployed cluster overlay
already configures an Argo app locator, dedicated token and Argo mesh DNS rewrite.
A read-only in-pod request to Argo's application API returned HTTP 200 for
`platform-openchoreo-projects`, Healthy/Synced, at 08:07:12 UTC. The empty locator
in the developer defaults is not the effective production configuration.

Grafana's dashboard-search API returned **401** using the existing configured
adapter credential. Its basic health probe and OpenChoreo CPU/memory charts are
separate functioning paths. Dashboard embedding/authenticated search is therefore
not certified working; no credential was printed, created, rotated or bypassed.

## Validation and disposition

- Added authorization regressions pass with existing policy/host-guard tests.
- Stable catalog references and matching production/development copies checked.
- TypeScript, changed-file lint, changed YAML and secret-scan results are recorded
  in the PR. Whole-tree `make validate` still fails on old unrelated formatting
  beginning in `platform/onetak/media-mtx.yaml`; GitHub billing remains separate.
- **Phase 6 is not marked universally complete:** 16 launcher findings, Grafana
  adapter authentication, historical-backup presentation, and the stated live
  token-expiry/role coverage limit remain explicit. A running portal is not a
  certificate that every registered service, backup or integration is healthy.
- Next is the already-authorized existing STAC pilot integration, still under
  Phase 6 operational coverage. No additional service cutover is implied.
