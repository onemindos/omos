# Portal stabilization before further migrations

Authorized by Zeus on 2026-10-20 after the live portal/catalog audit. This is a
Phase 6 follow-up ([PR #322](https://github.com/onemindos/omos/pull/322)), not permission to start Phase 7. Historical October 19 release
acceptance remains valid for its tested scope; it did not establish complete
catalog accuracy or UI quality. Migration drafts #320 and #321 stay unmerged.

## Read-only baseline

Observed 2026-10-20, approximately 02:06–02:17 America/New_York:

- Primary portal: https://onemind-preview.tail717ba0.ts.net/.
- Deployment 2/2 available, zero pod restarts; `platform-backstage` Synced/Healthy.
- Signed-in home loads 130 directory entries. Live overview separately reports
  observations; launcher entries do not imply installed/healthy services.
- STAC Browser Deploy shows `stac-browser-v5` Active in Development. Staging and
  Production explicitly show Not Deployed. No promotion was attempted.
- Freshly reloaded component Overview shows duplicate About/Relations, no
  description, and **Failed to load workflow data**, despite workflows disabled.
- Complete catalog API response: 606 entities, no duplicate entity references,
  reported status errors or orphan annotations. This does not prove processing
  completeness, link reachability or every human role's visibility.
- One missing relation target: `clusterworkflowplane:openchoreo-cluster/default`,
  referenced by four default builders. Read-only Kubernetes inspection confirms
  no ClusterWorkflowPlane exists. This is not fixed by inventing a catalog entity
  or automatically installing build infrastructure.
- Expanded timestamp audit finds 227 October 6 observations: the previously
  reported 155 inventory snapshots plus 72 repository observations. Neither set
  is current live discovery. Seven days is a review threshold, not a health SLA.
- Shared launcher candidates: Argo CD/external Argo CD, Tailscale/external
  Tailscale, Bambu Lab/Bambu Studio. A shared destination is not proof that two
  products should be merged. Preserve stable references until reviewed.
- Permission logs repeatedly report no token/no cached capabilities for the
  signed-in user. Deploy still loads; this is not proof of a login outage.

## First source repair

- Override only OpenChoreo's Component overview layout, retaining its managed-kind
  selector and the existing host/plugin cards. Render native About/Relations once
  instead of embedding a second pair. Non-OpenChoreo entities keep their layout.
- Gate Workflows with the existing upstream feature flag, so disabled workflows
  never mount/fetch. Keep deployment and enabled observability cards. Qualify the
  log-based Runtime Health result rather than treating no logs/errors as readiness.
- Add the pilot description to its owning Component annotation. No Workload,
  release, binding, image, endpoint, ownership or replica change.
- Add `scripts/audit-catalog.cjs` in the portal source. It consumes a **complete**
  `/api/catalog/entities` JSON array from a file or stdin and emits a metadata-only
  report. It does not fetch URLs, read credentials, change files or mutate catalog
  data. URLs are not emitted. A paginated query envelope is rejected; callers must
  not pass one page as a complete catalog. A permission-filtered view may produce
  missing targets; compare with an authorized complete view before diagnosing.
- Report exit codes: 0 no structural/reported errors; 1 duplicate refs, missing
  targets, orphan/status errors or invalid dates; 2 invalid input. Freshness and
  shared launchers remain explicit review findings even when exit code is zero.
- Regression checks run in the production Dockerfile, not only on the developer Mac.

## Permission investigation: preserve fail-closed behavior

OpenChoreo 1.3.1 `AuthzProfileService.getCapabilitiesForUser` first looks up the
user/scope cache. A missing token and expired/absent profile yields empty
capabilities and the observed warning. The policy intersects those results with
OneMind's own restrictions. The upstream service derives cache TTL from token
expiration (five-minute fallback); OneMind stores the original TTL in PostgreSQL.
Browser permission requests supply an IDP token; internal delegated permission
calls may lack it. Sign-in prewarms the cache; a successful browser Deploy request
does not establish that every internal authorization request has a warm profile.

Added a regression using the actual upstream service/policy and shared-cache
adapter: after expiration, a tokenless administrator request is DENY, emits the
warning, and does not fetch capabilities. This validates safe behavior, **not the
specific production request path causing repeated misses**. No TTL increase,
warning suppression, stored raw token, admin bypass or permission change is made.
Ollama DeepSeek-V4.1 Flash provided a bounded read-only review; GPT verified the
installed dependency behavior and implemented/tested the repair independently.

## Acceptance gates (updated after final rollout)

- [x] Authorized review/merge of the Layer 1 portal source PR; immutable Linux image
      build/publication, separate verified-digest GitOps PR and Argo rollout.
- [x] Browser reload on the new image: one About/Relations pair, no disabled
      workflow fetch/card, pilot description, unchanged home/Today/catalog and
      working Deploy/logs/metrics. Check wide and narrow layouts.
- [x] Recheck both replicas, authenticated human reads, machine read-only MCP,
      anonymous denial and machine/platform mutation denial after rollout.
- [ ] Correlate capability warnings with the actual production request paths and
      token-refresh lifecycle without logging credentials. Verify refreshed login
      and expiry across both replicas. Keep denials until authorization is valid.
- [x] Run transactional inventory refresh with all Kubernetes/Tailscale/GitHub
      sources available, review identity/relationship diffs in a separate PR,
      and verify mounted/processed catalog matches after merge. Do not merely
      rewrite timestamps or replace partial discovery results.
- [x] Reconcile launcher/service overlap with stable canonical identities and
      preserved incoming references; distinguish reference/planned/local/cloud
      entries from actual deployments. Review first-party launcher reachability.
- [x] Decide the disabled builders' treatment explicitly: document/defer inactive
      build definitions or provision a real plane under its own approved scope.
      Keep the unresolved relation visible until that decision is implemented.
- [x] Review other OpenChoreo entity-kind layouts for equivalent duplication and
      empty cards; preserve relevant third-party cards and authorization filters.
- [x] Review intended COP coverage. The effective production Argo locator works
      (corrected in the follow-through report); the Grafana dashboard adapter
      returns 401. Review is complete, not a claim that every integration works.

Commands from `services/onemind-portal` (Node 22.23.3, committed Yarn):

```sh
node .yarn/releases/yarn-4.13.0.cjs tsc --incremental false
node .yarn/releases/yarn-4.13.0.cjs workspace app test --watch=false --runInBand ComponentOverview.test.tsx componentOverviewModule.test.ts
node .yarn/releases/yarn-4.13.0.cjs workspace @internal/plugin-onemind-permissions-backend test --watch=false --runInBand
node --test scripts/test-audit-catalog.cjs
node scripts/audit-catalog.cjs /absolute/path/to/complete-entities.json
```

The baseline catalog report intentionally exits 1 for the missing workflow-plane
relation. Do not suppress it to manufacture green acceptance. Repository-wide
`make validate` still fails on unrelated pre-existing YAML formatting debt; report
that separately from changed-file gates and local portal tests.

### Local validation of this repair

- Immutable dependency install, TypeScript, frontend production build, frontend
  lint and changed backend test-file lint passed on Node 22.23.3.
- 50 permission/auth tests, five overview/extension-selector tests and six catalog
  audit tests passed (61 total). These are source checks, not live acceptance.
- Changed YAML lint passed. Kubeconform skipped the two OpenChoreo CRDs because
  published schemas were unavailable; both resources additionally passed actual
  cluster admission with `--dry-run=server`. Nothing was applied.
- No Helm-affecting change. Commit-range and clean-export baseline-aware secret
  scans passed. Full permission-package lint still reports two existing errors in
  unchanged `sharedCache.ts` (`no-use-before-define`, `consistent-return`).
- At creation of source PR #322, no Linux release image had been built/published
  and no deployment digest had changed; subsequent release evidence follows.
  Layer 1 review/merge and the subsequent release/live gates remain required.

GitHub checks on PR #322 did not start: the check annotation reports account
payment/spending-limit failure. Local passes are not GitHub green; no billing or
merge-rule bypass was attempted.

## Repair release (2026-10-20)

Zeus explicitly authorized merging #322 and completing deployment/verification.
Source PR merged at 06:43:55 UTC as
`68a5c612916e8caafa4df945d375dbfc3f7aacb2`.
The pinned Docker Desktop Linux build passed immutable install, TypeScript,
50 security/auth tests, five overview tests, six audit tests, TechDocs and backend
compilation. GHCR publication and remote digest inspection succeeded:

- Image: `ghcr.io/onemindos/onemind-portal@sha256:694dfc6e188f9664c513bf160f4d474080a9d586fe566087c3a5bda9c1aa3cb6`
- Source revision label: `68a5c612916e8caafa4df945d375dbfc3f7aacb2`
- Runtime architecture: Linux/AMD64 (not the Mac's native architecture).
- No backend schema/package cohort or production permission-policy change.
- Rollback: restore digest `sha256:b2fc101eacc504fadff99b638b342ce9bf56782a4901ea49dffe2fdd09f1967a`
  through a reviewed GitOps PR. Do not alter the database or identity.

The image pin is a separate GitOps change. Live rollout/browser/MCP evidence will
be appended after Argo reconciliation; the image build alone is not acceptance.
GitHub jobs still fail before execution because of the account billing restriction.
No billing, branch protection or credentials were changed.

## Live acceptance of the first release and inventory

- #322 source, #323 immutable image pin and #324 catalog reconciliation merged.
- At 06:57–07:03 UTC, Argo reconciled both rollouts; both replicas became ready
  on digest `694dfc6e188f9664c513bf160f4d474080a9d586fe566087c3a5bda9c1aa3cb6`.
- Each replica passed readiness 200, anonymous APIs 401, disabled guest 404,
  machine catalog read 200, catalog/platform/cache/context writes/access 403,
  and MCP enumeration of six explicitly read-only/non-destructive tools.
- Signed-in browser: pilot Overview has one About/Relations pair, the source
  description, no disabled-workflow error; Deploy shows stac-browser-v5 Active;
  Logs shows a fresh HTTP 200 at 06:58:33 UTC; CPU/memory graphs render.
- Processed catalog converged to 715 entities, containing all 664 expected file
  source identities. No duplicate refs, reported processing errors or orphans.
  Both canonical alias annotations appeared after processing, not merely in Git.
- Stale observations fell from 227 to 22, exactly matching the retained absent
  identities marked missing-review-required. No unknown was converted to healthy.
- The only missing relationship target remains the uninstalled workflow plane
  referenced by four inactive default builders. Build workflows stay disabled;
  prebuilt-image GitOps releases do not require pretending that plane exists.
- Fresh browser authentication logged capability pre-caching. Both new replicas
  serve authorized internal catalog reads using the shared profile cache; no
  repeated no-token/no-cache warning was observed during these checks. Expiry
  denial is regression-tested, not suppressed or bypassed; deliberately expiring
  the production identity or claiming every long-running session is tested is
  outside this observation.

## Project/platform overview composition follow-up

The live OneTAK project page still rendered two About/Relations pairs after the
Component repair. A second small frontend module filters only the host duplicates
for the 22 OC-managed kinds with reviewed bespoke upstream overview layouts.
Components keep OneMind's host-owned pair; ordinary Backstage records, managed API
and unknown future kinds retain native cards. User/Group About exclusions remain.
No third-party card, permission, backend, dependency cohort or metadata is removed.

Tests exercise both real host extension filters, all 22 actual upstream layout
selectors, and the pinned OC 1.3.1 version. The targeted release suite is now 100
checks (50 auth/security, 44 frontend, six audit). Source checks/build do not
establish live acceptance; a new immutable image and browser verification follow.
GLM-5.3 Flash's bounded read-only review found no concrete functional defect;
GPT addressed its mapping-coverage gap with the 22 selector tests. Worker reported
76,511 input tokens (45,568 cached) and 1,914 output tokens; no charge is inferred.

## Composition release and visible missing-inventory qualifiers

#325 merged as `de8764e3647bd91613ee3097040184e89918c37e`. Its pinned Linux
production build passed all 100 targeted checks, TypeScript, TechDocs and backend
compilation. Published/remote-verified image:
`ghcr.io/onemindos/onemind-portal@sha256:4a52c3b840fcd2b4e2fc0287929e1704721f7420018f2c09deb8543a6aab74da`.
This image contains the project/platform card composition fix; no plugin cohort,
authorization or database change. Its immediate rollback digest is the first
stabilization release `694dfc6e188f9664c513bf160f4d474080a9d586fe566087c3a5bda9c1aa3cb6`.

The final GitOps change also prefixes the descriptions of exactly the 22 retained
absent inventory records with **Not observed in the latest inventory; review
required.** This is visible in search/About, not only hidden in annotations. It
asserts neither downtime nor deletion and preserves IDs, links and timestamps.
The operator-only refresh script applies this qualifier idempotently on future
missing observations (three standalone regression tests). It is not a runtime
portal plugin and does not require rebuilding the just-verified frontend image.
The mounted catalog ConfigMap carries the reviewed descriptions. A newly observed
identity is rebuilt from actual discovery rather than inheriting a missing flag.

No Ollama worker remains active. Zeus requested no further Ollama use during the
follow-through; all remaining implementation, release and verification are by GPT.
Final Argo/browser/MCP/catalog observations must be recorded after reconciliation.


## Final deployed verification — 2026-10-20, 07:28–07:35 UTC

Source/composition [#325](https://github.com/onemindos/omos/pull/325) and
image/catalog [#326](https://github.com/onemindos/omos/pull/326) are merged.
Argo `platform-backstage` is **Synced/Healthy** at
`76ef8f4e4e8138b5a41fd59fed6dcf486a0d87a9`. Deployment is **2/2 ready**;
`backstage-79456c596c-59lkk` and `backstage-79456c596c-26dnh` have zero restarts
on the final `4a52c3b840fcd2b4e2fc0287929e1704721f7420018f2c09deb8543a6aab74da`
digest. OCI source label matches the #325 merge. Runtime Linux/AMD64 SQLite and
exported JSON smoke tests passed. No database migration or credential change.

- Release build: **100 targeted tests passed** (50 auth/security, 44 frontend,
  six catalog audit), immutable install, TypeScript, frontend/backend and TechDocs.
  Three additional operator-only inventory-description regression tests passed.
- Each final replica: readiness 200; anonymous catalog/OpenChoreo/permission/cache
  APIs 401; guest disabled 404; machine catalog read 200; machine catalog writes,
  OpenChoreo access and capability/context access 403. MCP exposes six explicitly
  read-only/non-destructive tools. A real MCP health invocation returned four
  configured HTTP-200 probes and namespace metrics; this is not all-app health.
- Signed-in final frontend: OneTAK Project and STAC Browser Component each show
  one About/Relations pair; the pilot description renders, disabled Workflow
  card/error is absent. Deploy shows **stac-browser-v5 Active in Development**;
  Staging/Production remain Not Deployed. No promotion or other action invoked.
- Pilot endpoint returned HTTP 200 at **07:30:05 UTC**; its fresh
  `portal-final-gate=20261020` request appeared in the Logs view. CPU/memory usage,
  requests and limits charts rendered at 07:31:53 UTC. Home retained 130 directory
  entries, account favorites and timestamped live overview. These counts are
  different concepts from catalog size or installed applications. Today loaded
  at 07:34:53 UTC and explicitly showed five attention items: Matrix OutOfSync,
  OpenClaw Progressing and three suspended backup schedules (Backstage,
  ClickHouse, offsite-copy). These require separate operational review; portal
  health must not be represented as cluster-wide health or backup readiness.
- Final processed catalog: **715 entities**, including all **664 file-source
  identities**; zero missing expected source refs, duplicate refs, reported
  processing errors or orphans. There are 325 Components, not 715 applications.
  Exactly **22 retained missing-review records** visibly carry the qualified
  description. Canonical external Argo/Tailscale annotations are processed live;
  stable identities, incoming references, URLs and favorites were preserved.
- Build-plane decision: **defer provisioning**. Four inactive builder definitions
  still refer to the absent `clusterworkflowplane:openchoreo-cluster/default`.
  Workflows remain disabled and the structural audit intentionally exits 1 for
  that relation. Prebuilt-image GitOps deployment does not require a build plane.
- Fresh browser authentication pre-cached valid capabilities; final steady-state
  log sampling found no repeated no-token/no-cache warning. A startup catalog
  event-subscription 404 recovered. This is not a universal zero-warning claim.
  Expired tokenless authorization remains DENY in the actual-service regression.

### Scope and outstanding gates

The concrete UI/catalog repair is deployed and its targeted checks pass. **Do not
call the broader Phase 6 stabilization universally complete or start a bulk
migration.** The original remaining reviews are tracked in the
[follow-through](portal-acceptance-followthrough-2026-10-20.md), including explicit
launcher findings, Grafana authentication and live token-expiry/role limits.
Tests cover all 22 actual upstream layout selectors; that is not a human browser
inspection of every entity. No second human role was created or live-tested.

The earlier assumption about a missing production Argo locator is superseded by
the [follow-through review](portal-acceptance-followthrough-2026-10-20.md): it is
configured and its read-only API works. Grafana dashboard search returns 401;
secure embedded viewing is not certified. API/dependency/TechDocs
metadata is incomplete for some apps and must be completed per migration.
Wide desktop and narrow browser observations are scoped smoke checks, not all
pages/devices. The user's active desktop window was left undisturbed during
final checks; an older already-open SPA may need a normal refresh.

GitHub jobs still cannot start due the account billing/spending-limit error;
local validation is **not GitHub green**. Full `make validate` still encounters
pre-existing YAML formatting debt beginning in `platform/onetak/media-mtx.yaml`;
unchanged permission-cache lint debt is recorded above. No billing, protection,
permission or secret settings were changed to bypass these issues.

**No application was migrated by these repairs.** Drafts #320/#321 remain
unmerged. The next controlled migration requires explicit app selection,
provider identity-collision review, ownership/dependency metadata, data backup
and restore validation where applicable, rollback and live acceptance. A healthy
portal is not proof that every app is ready for cutover or offline autonomy.
