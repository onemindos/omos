# OpenChoreo workflow-plane pilot

Requested by Zeus on 2026-10-20. Priority: **high**. Status: **Not started**.
This brief records future work; it does not install a plane, enable workflows,
grant access, select a paid builder, or authorize an application migration.
Phase 6 stabilization remains the first unfinished phase in the root roadmap.

## Outcome and rationale

Prove one small, first-party application's source-to-image path through
OpenChoreo, with build status/logs available in the existing OneMind Backstage.
Keep production promotion explicit, reproducible and independently recoverable.
This advances the home-centered, AI-readable platform without requiring GitHub
Actions as the final release builder or Docker Desktop as a manual release step.

Do not rebuild every third-party app or replace every existing pipeline. Published
upstream images remain appropriate for many applications. Compare this pilot
against the working external-build path before expanding it. A workflow plane
does not solve the GitHub account billing restriction or make PR checks green.

## Recorded starting point, not a fresh live audit

- The October 20 stabilization evidence records Backstage 1.55.3, OpenChoreo
  1.3.1, two healthy portal replicas and working prebuilt-image deployment.
- `services/onemind-portal/app-config.yaml` has `openchoreo.features.workflows`
  disabled. The absent default ClusterWorkflowPlane is still an explicit catalog
  relationship finding from four inactive builders; do not fake a target entity.
- Argo CD owns platform installation and declarative application intent. Argo
  **Workflows** executes pipelines; it is not a replacement for Argo **CD**.
- Existing human and machine permissions are unchanged. The current portal MCP
  identity is read-only; displaying a workflow must not silently grant execution.
- Capacity, backup schedules and operational exceptions must be re-observed
  before implementation. The Phase 5 capacity snapshot is not today's budget.

See [stabilization evidence](../operations/portal-stabilization-2026-10-20.md)
and [mission](../MISSION.md). No currently deployed app depends on this proposal.

## Proposed responsibility boundary

1. Git retains reviewed source, Dockerfile/build instructions, tests, pipeline
   templates, platform installation and declared environment releases.
2. OpenChoreo's Workflow Plane executes an approved, parameterized pipeline for
   an immutable source revision. Start manually in a development-only pilot.
3. The pipeline runs tests/security gates, builds the intended architecture,
   publishes to a durable OCI registry and records its digest and source revision.
4. A deliberate handoff proposes the digest/release change through a GitOps PR.
   A successful build alone never promotes production or bypasses merge checks.
   Verify the actual v1.3.1 workload-generation/auto-deploy behavior first; if it
   cannot respect this boundary, stop and redesign the handoff before enabling it.
5. After authorized merge, Argo CD reconciles declared intent; OpenChoreo's
   controllers reconcile application resources. Do not make both controllers own
   the same generated child resources.
6. Backstage displays source, run, artifact, release, environment and observed
   deployment state. Human and AI views use the same service identities/evidence.

Workflow definitions and installation belong in Git. Individual `WorkflowRun`
resources are imperative executions, not endlessly reconciled desired state;
do not commit them to the app-of-apps tree. Any future run-trigger API requires
its own scoped authorization, audit and replay/idempotency design.

## Implementation steps: separate reviewed PRs, one at a time

### A. Design, compatibility and capacity

- [ ] Recheck live Phase 6 gates, actual installed CRDs/chart values, controllers,
      Authentik integration, existing secret-store pattern and upstream 1.3.1
      templates. Pin reviewed artifacts; do not copy development/latest samples.
- [ ] Select one small first-party, stateless, non-critical application and a
      development target. Do not make the portal, identity, backups, TAK core or
      the workflow plane itself the first bootstrapping-dependent pilot.
- [ ] Record CPU/memory/ephemeral-storage requests and limits, maximum concurrent
      builds, timeouts, cache/log retention and a measured cost/capacity ceiling.
      Start with concurrency one. Obtain approval before any added paid capacity.
- [ ] Choose placement from evidence: isolated worker capacity in the current
      cluster versus a separate builder cluster/host. Do not name ephemeral nodes.
      Preserve production capacity; a separate namespace is not a security sandbox.
- [ ] Choose durable registry, push/pull identities and retention. GHCR/Harbor are
      candidates, not a new registry decision. Never use ephemeral ttl.sh for a
      production artifact. Verify independent pull/recovery and home portability.
- [ ] Capture exactly how completed artifacts enter the existing GitOps release
      path without direct production writes or competing reconciliation.

### B. GitOps-owned plane and security

- [ ] Add a pinned Argo-owned workflow-plane installation and registration,
      reviewed network/identity/secret references and minimal service accounts.
      No unmanaged Helm install or imperative workload apply.
- [ ] Permit only required source/registry/artifact endpoints; restrict workflow
      authorship separately from run triggering. No cluster-admin, host Docker
      socket, blanket privileged builders, production secret mounts or credentials
      exposed to browser clients/build logs. Explicitly review any exception.
- [ ] Keep fork/untrusted-PR builds away from publishing credentials and trusted
      caches. Separate trusted release builds from untrusted verification work.
- [ ] Build-plane outage cannot prevent the existing portal, Argo recovery,
      authentication or deployed applications from functioning.

### C. One end-to-end application and human/AI visibility

- [ ] Run locked dependency install, tests, secret scanning and agreed image
      vulnerability policy before publication. Record SBOM/provenance; evaluate
      signature verification before promotion without treating signing as a given.
- [ ] Build/verify Linux architecture for the selected target; document ARM64
      versus AMD64 requirements for future Mac/home deployments. An image that
      builds on one architecture does not prove multi-architecture compatibility.
- [ ] Connect the existing portal's workflow UI only after the real plane and
      authorization are working. Test both replicas, browser login/refresh,
      run/log visibility, anonymous denial and current machine write denial.
- [ ] Verify the approved workflow is actually allowed by the selected
      ComponentType. Do not simply turn on every default builder or template.
- [ ] Demonstrate a passing build, failing test that blocks publication/promotion,
      cancelled/timed-out run and duplicate-trigger handling. Verify source and
      remote digest before creating the release PR; deploy only to the approved
      development environment and verify workload/readiness plus app behavior.
- [ ] Capture comparable duration, resource use, upkeep and failure recovery
      versus the existing builder. Decide expand, revise or defer from evidence.

### D. Recovery and portable operations

- [ ] Disable triggers/feature flag and restore the external-build release path
      without changing deployed workloads, production data or identity.
- [ ] Rebuild a disposable builder from Git with restored secret references;
      retained OCI artifacts must remain usable if build caches/logs are lost.
- [ ] Define source mirrors, registry mirrors, dependencies, identity and recovery
      prerequisites for a future home/offline builder. Self-hosted computation is
      not proof of offline operation while GitHub/registries remain dependencies.
- [ ] Record source/chart revisions, registry digest, permissions/negative tests,
      observed live results and rollback before marking the pilot accepted.

## Non-goals and expansion rules

- No automatic migration of #320/#321 or broad switch from GitHub Actions.
- No removal of PR checks, security gates, human Layer 1 review or backup gates.
- No new NATS subject/event contract or agent-autonomy permission in this pilot.
- General ETL, geospatial processing or ML workflows are possible later but require
  their own resource, data and authorization contracts. Physical device control
  and real-time safety loops stay in their owning services/controllers.
- Keep portal plugins in the existing `services/onemind-portal/plugins/` workspace;
  reuse its workflow integration rather than forking upstream or adding another UI.

## References checked 2026-10-20

Upstream documentation describes capabilities, not validation of our deployment.
Reconcile it with the exact installed/pinned 1.3.1 source before implementation.

- [OpenChoreo workflows: Argo engine, definitions and imperative runs](https://openchoreo.dev/docs/platform-engineer-guide/workflows/overview/)
- [CI governance and allowed workflows](https://openchoreo.dev/docs/platform-engineer-guide/workflows/ci-governance/)
- [Architecture: optional workflow plane and external/hybrid CI](https://openchoreo.dev/docs/overview/architecture/)
- [Plane deployment topology](https://openchoreo.dev/docs/platform-engineer-guide/deployment-topology/)
- [Existing OneMind GitOps and review rules](../../CONTRIBUTING.md)
