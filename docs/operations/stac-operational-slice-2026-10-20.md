# OneTAK/STAC operational slice — 2026-10-20

Requested by Zeus after the Phase 6 trust review. This is integration of the
**existing** Development pilot, not a new migration, permission grant, workflow
plane, agent-action capability or database repair.

## Implementation

- Existing `onemind-services` plugin adds its Operations tab to the exact
  OC-managed `component:default/stac-browser` and `system:default/onetak` pages.
  Native namespace operations remain unchanged for other entities.
- The tab joins stable identity, reviewed source, actual dependencies, deployed
  images/readiness, scoped CPU/memory queries, endpoint behavior, data-query
  availability, GitOps configuration, logs/metrics links and a recovery runbook.
- `GET /api/onemind-services/stac` requires the existing authenticated human
  identity. It accepts no target URL, namespace, SQL or mutation instruction.
  No machine scope is widened; this is not the later shared AI-action-contract step.
- Fixed read-only observations use the existing service account and public
  internal metric/API targets. No new Kubernetes RBAC or app credential is added.
  Individual HTTP reads have a five-second absolute deadline and two-megabyte
  response cap, no redirects and no credential forwarding. Overlapping requests
  for the same URL are coalesced; the observation cache is at most 60 seconds.
- Browser traffic is checked through the existing Cilium Gateway service using
  the reviewed Host header. It does not bypass the OC network policy with a
  cross-project direct-Service connection. Client/tailnet reachability is separate.
- Missing/ambiguous workloads stay unknown. Readiness requires observed generation;
  scaled-to-zero and reconciling are explicit. Metrics are per-selected Deployment
  pods; their timestamps are Prometheus query times, not ingestion/scrape proof.
- The curated catalog processor adds Browser → API → database and consumed/provided
  read-only API relations without replacing the OC provider or inventory IDs.
  One new API descriptor documents the subset; it is not full STAC conformance.
- New TechDocs runbook explains diagnosis and backup/migration gates. Catalog
  metadata and status never grant permission to initialize a schema or restart apps.

## Real finding, deliberately not concealed

A read-only prerelease invocation of the compiled observer in the existing portal
pod at 08:28:32 UTC observed:

| Signal | Observation |
|---|---|
| Browser/API/database Deployments | ready; observed generation current |
| Each selected workload CPU/memory | values observed |
| Browser via Gateway | HTTP 200 |
| API landing page | HTTP 200 |
| Collection listing and bounded sorted search | HTTP 500 |
| Data freshness/counts | unknown, not zero or healthy |
| Both owning Argo applications | Healthy/Synced |

API logs reported `UndefinedFunctionError: function search(jsonb) does not exist`.
Read-only database inspection found no `pgstac` namespace in database `stac`.
No schema, data, database identity or backup was altered. A separate repair needs
exact-version/schema compatibility, inventory of existing data, encrypted backup,
isolated restore and reviewed Argo-owned migration/initialization. It must not be
silently bundled into a dashboard change or a pretend empty demo catalog.

This prerelease probe executes read-only observation code, not a deployed frontend
or immutable-image acceptance. Live release evidence must follow.

## Source validation

- Node 22 immutable install and TypeScript passed.
- **26 new targeted tests** cover transport limits/no redirects, source failure,
  mixed ready/500 signals, ambiguity, metric unknown/stale handling, bounded data,
  topology preservation/idempotence and read-only UI/selector behavior.
- Existing release auth/layout/audit tests remain enabled. Dockerfile adds the
  new tests so Linux release gates exercise them, not only the host Mac.
- Changed-file ESLint has no errors. The existing services-to-common frontend
  package-role warning remains; the Material UI v4 test renderer also emits its
  existing findDOMNode deprecation. No dependency-cohort upgrade is mixed in.
- Only already-installed catalog package versions and an existing test-library
  version become explicit workspace dependencies. No new external runtime service.
- Paired catalogs, schema/changed-YAML, commit-range secret scan and `make validate`
  are checked for the PR; full-tree formatting and GitHub billing failures remain
  separately reported, never described as green CI.

## Release/live gate

### Follow-through review (2026-10-20)

Zeus explicitly approved verified local/Linux release gates while GitHub's jobs
cannot start because of its account payment/spending restriction. Do not alter
billing, disable checks, force an administrator merge, or mark those jobs green.

Fresh review reproduced missing incomplete-state warnings for reconciling/zero
replicas and Progressing/OutOfSync Argo, plus a hanging browser request that could
stay loading indefinitely. Added regressions and repairs: non-ready/nonaccepted
states set the incomplete summary; a 15-second client request deadline retains
older evidence only as stale, permits later retries and cancels on unmount.
The observation clock advances during polling and future timestamps are stale.

The previous published candidate at source `af7770fb` / digest
`sha256:fab6522342b40e6ce595ab2e932c84399d3aefec843af45829bade802f955de3`
was independently verified in GHCR as Linux/AMD64 with matching revision/source
labels. **It predates these repairs and must not be promoted as their image.**
A new clean-commit Linux build, source label verification and separate digest-pin
PR are required. This does not fix STAC's backing schema or certify a new rollout.

- [ ] Publish and remotely verify an immutable Linux/AMD64 image with all gates.
- [ ] Pin its digest in a separate GitOps PR; Argo healthy and both replicas ready.
- [ ] Process API/dependency metadata with stable references, no new duplicates,
      errors/orphans or missing relations beyond the known inactive workflow plane.
- [ ] Signed-in Project and Component Operations pages render the joined evidence,
      real 500/data-unknown findings, source/logs/metrics and TechDocs links.
- [ ] Anonymous access denied; machine direct service/platform writes/access remain
      denied; existing six read-only MCP tools remain read-only.

Phase 6 is not declared universally complete by this feature. The follow-through
report retains 16 launcher findings, Grafana adapter 401, offsite-copy gap and
production expiry/second-role coverage limits. Application readiness and recovery
must be proven before a later migration or database change.
