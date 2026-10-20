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

## Remaining acceptance gates (must not call this complete)

- [ ] Human review/merge of the Layer 1 portal source PR; immutable Linux image
      build/publication, separate verified-digest GitOps PR and Argo rollout.
- [ ] Browser reload on the new image: one About/Relations pair, no disabled
      workflow fetch/card, pilot description, unchanged home/Today/catalog and
      working Deploy/logs/metrics. Check wide and narrow layouts.
- [ ] Recheck both replicas, authenticated human reads, machine read-only MCP,
      anonymous denial and machine/platform mutation denial after rollout.
- [ ] Correlate capability warnings with the actual production request paths and
      token-refresh lifecycle without logging credentials. Verify refreshed login
      and expiry across both replicas. Keep denials until authorization is valid.
- [ ] Run transactional inventory refresh with all Kubernetes/Tailscale/GitHub
      sources available, review identity/relationship diffs in a separate PR,
      and verify mounted/processed catalog matches after merge. Do not merely
      rewrite timestamps or replace partial discovery results.
- [ ] Reconcile launcher/service overlap with stable canonical identities and
      preserved incoming references; distinguish reference/planned/local/cloud
      entries from actual deployments. Review first-party launcher reachability.
- [ ] Decide the disabled builders' treatment explicitly: document/defer inactive
      build definitions or provision a real plane under its own approved scope.
      Keep the unresolved relation visible until that decision is implemented.
- [ ] Review other OpenChoreo entity-kind layouts for equivalent duplication and
      empty cards; preserve relevant third-party cards and authorization filters.
- [ ] Verify intended COP coverage. The dedicated Argo plugin has no configured
      app locator; existing OneMind Argo summaries are a separate working path.
      Grafana links are available; secure embedded viewing remains unconfigured.

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
- No Linux release image was built/published and no deployment digest changed.
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
