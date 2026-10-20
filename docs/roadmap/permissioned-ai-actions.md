# Permissioned AI actions: from observation to controlled operation

Status: **planning only**, requested by Zeus on 2026-10-20. No action executor,
write tool, identity, credential or permission is deployed by this PR. Merging
this document would record the plan, not authorize activation or complete a phase.

Parent: [platform tracker](../../PLATFORM-ROADMAP.md), Phase 6 trust stabilization
and K03/K07 in the [operating structure](../architecture/kubernetes-operating-structure.md).
Preparation can run alongside stabilization; activation needs the affected trust
and app-readiness gates. This does not replace the existing delivery queue.

## Goal and current boundary

Backstage is the human operational front door; authenticated AI tools expose the
same bounded capabilities. Kubernetes, Argo and each application's owning
controller/API remain execution authorities. Visibility is not authorization.

The existing [AI access](../operations/ai-operator-access.md) is read-only. Portal
`backend.actions` filters include read-only and exclude destructive actions;
external machine access and the OpenChoreo host guard impose additional limits.
The existing [release pipeline](../operations/gitops-release-pipeline.md) is an
operator-run local/Linux tool, **not an autonomous agent service**. Its only
checked-in application profile is the portal, a protected foundation workload.
Do not use that profile as the first autonomous write target or pretend it is a
multi-app executor. Local operator credentials used in an authorized coding
session are not the deployed agent's credentials.

Keep existing read-only MCP endpoints, Kubernetes RBAC and portal guards intact.
Do not remove the read-only filter, forward a human's session/token to an agent,
give the portal cluster-admin or use the generic OpenChoreo proxy as a bypass.

## Proposed control path

```text
Human in Backstage / separately authenticated agent tool
  -> backend authorization and typed action request
  -> durable request/approval journal + policy evaluation
  -> isolated, narrowly credentialed executor
  -> validated application PR -> exact-head human review/merge
  -> existing Argo reconciliation -> observed rollout + application checks
  -> durable result, evidence and recovery proposal shown in Backstage/AI tools
```

The executor's privilege is separate from the portal observer identity. Keep UI,
contracts and adapters in the existing monorepo packages; an extra repository or
an additional orchestration platform is not required by this plan. An API/worker
runtime and its ownership, recovery, identities and deployment need their own
implementation review. Reuse the pipeline's validation/artifact/verification
contracts through a reviewed adapter, not arbitrary command strings from an LLM.

NATS may carry requests/progress, but is not an authorization service or evidence
that execution succeeded. This plan defines no subjects or bus contract changes.
Any later transport integration must follow the locked NATS doctrine and prove
identity, ACLs, correlation, replay safety and disconnected behavior separately.

## Authority progression

| Tier | Capability | Approval boundary |
|---|---|---|
| Observe | Read allowed health, inventory and telemetry | Existing authenticated read-only scope; no new grants |
| Propose | Prepare a typed plan and draft application-release PR | Future explicitly approved repo/branch credential; no main write or merge capability |
| Supervised execution | Deploy a reviewed digest through a normal PR merge; observe and verify | Human approval of exact head/diff/target; no automatic merge or force sync |
| Bounded autonomy | Selected app actions under an explicit, expiring preauthorization | Future per-app/action policy after supervised acceptance; not enabled by this plan |
| Foundation/shared changes | Identity, access, networking, security, storage, recovery, controllers, NATS and databases | Preserve current layer rules: propose PRs; human approves/merges; no blanket remediation grant |

A capability such as `propose_release`, `observe_release` or
`propose_rollback` is a **proposed contract name**, not an installed MCP tool.
No arbitrary shell, arbitrary manifest apply, unrestricted URLs, raw kubeconfig,
Secret reads, exec, deletion/pruning, credential rotation, public exposure,
Terraform apply, provider purchases or equipment actuation in the initial pilot.

## Request, approval and execution contract

Before implementing write actions, specify and test:

- Verified caller identity (issuer, audience, expiry and intended agent/user),
  separate executor identity and authorization checked on every backend entry
  point. An `agentId`, catalog annotation, NATS subject or UI button is not proof.
- Versioned action type with schema-validated parameters; allowlisted repository,
  branch policy, manifest path, cluster, environment, namespace, catalog reference,
  actual owning Application/controller and workload identity. Never accept a
  caller-selected shell command, API host or filesystem path.
- Expected base commit, exact candidate commit/diff hash, verified artifact digest
  and source/build evidence, current observed revision/generation and policy
  version. Re-read state immediately before execution; drift invalidates approval.
- A durable approval bound to that request hash, verified approver, exact target,
  action and expiration. No self-approval, replay, transferable chat "yes", or
  approval automatically extended to a changed PR head. Direct API/MCP calls
  enforce the same policy as the human UI.
- Separate propose/check/approve/execute/observe/result states with run ID,
  idempotency key, per-target lock, bounded retries/concurrency/time/spend and a
  durable outcome. Distinguish denied, expired, cancelled, running, succeeded,
  failed and **unknown/needs reconciliation**. Broker delivery is at-least-once;
  do not claim exactly-once effects across GitHub and Kubernetes.
- Lost acknowledgements, worker restarts and duplicate jobs reconcile existing
  branch/PR/merge/run state before retrying. Cancellation cannot undo a merged
  commit; once merged, report the stage honestly and require a new reviewed
  recovery change rather than promising instant reversal.
- Minimal, short-lived backend-only credentials; no existing human/admin token
  reuse. Repository enforcement plus executor path/diff checks must protect main,
  workflow/security/identity files and other targets. Reverify actual branch
  protections/rulesets and merge permissions; do not assume they exist. Stop
  activation if enforceable separation is unavailable. Creating credentials or
  changing repository protections requires its own explicit authorization.
- Redacted audit/evidence with requester, approver, target, versions, check receipts,
  timestamps, execution identity and outcome; no raw tokens, app payloads or
  personal prompts. Define retention, access, backup/restore and failure handling.
  Unavailable policy, approval or journal storage fails closed for new mutations.
- An operator-controlled stop/revocation path independent of Backstage and the
  executor's own availability. Disable new requests and revoke unused approvals;
  reconcile already-merged/in-flight effects. Prove revocation across replicas,
  queues, cached decisions and reconnects, not only a hidden frontend button.

## First pilot: one application-release PR, not fleet-wide write access

**Target remains unselected.** Select a low-risk stateless keeper only after its
K03 onboarding/profile, owning controller, application-level checks and recovery
path are reviewed. Neither Backstage (foundation), NATS, a database nor the
currently failing STAC data path is automatically chosen by this plan.

1. Register that exact app/environment/action allowlist and reviewed pipeline
   profile. Record exclusions, budget and dependency boundaries.
2. Agent prepares a release proposal using a verified immutable image digest.
   Run the required checks on the exact diff; show evidence and recovery plan in
   a draft PR and Backstage request record. No cluster writes at this stage.
3. Human reviews and merges the exact validated head using normal GitHub rules.
   No executor merge token, admin bypass, disabled checks or payment changes.
   While GitHub CI cannot start, apply the existing explicitly authorized
   local/Linux-check process with receipts; do not label GitHub checks passed.
4. Existing Argo automation reconciles the change. Executor needs read access for
   observation, not force-sync/apply permission. Detect another release winning
   the race, drift, timeout, wrong digest or wrong environment and stop safely.
5. Verify actual merged/observed revision, workload imageID and readiness,
   application behavior, dependency health and data freshness. Healthy pods alone
   are not success. Record the result in the shared journal and show it to the
   initiating human/agent without leaking credentials.
6. Exercise a reviewed rollback proposal to a verified compatible digest. No
   implicit schema/data rollback; stateful changes remain excluded. Human reviews
   the recovery PR and the actual recovered service is verified.

Only after this supervised path passes may a later PR propose expiring, bounded
preauthorization for that app. Restarts/scaling must use its actual owning
controller and explicitly resolve GitOps ownership; no CronJob patch loop
fighting Argo. Shared/foundation layer protections do not disappear when a
stateless application earns limited autonomy.

## Implementation queue and unexecuted acceptance gates

Each row is a later small PR, reviewed one at a time. This planning PR does not
implement any row or authorize the live credential/permission changes it needs.

| Step | Deliverable | Gate before advancing |
|---|---|---|
| A0 | Threat model, typed capability/policy schemas and exact trust/target map | Existing guards identified; affected token-expiry/role gates resolved; least-privilege grant reviewed |
| A1 | Offline policy tests, durable journal/approval model and fake executor | All negative/replay/drift/crash cases below pass without production credentials |
| A2 | One app profile plus draft-PR proposal adapter | Explicitly authorized scoped credential; enforceable protected paths/main; exact-head validation; cannot merge |
| A3 | Backstage request/result surface and separate authenticated agent contract | Same backend denials in UI, REST and tools; second identity and cross-target isolation proven |
| A4 | Human-approved pilot deployment and reviewed recovery drill | Live app/digest/Argo/evidence gates pass; stop, timeout and revocation work; current read-only endpoints remain unchanged |
| A5 | Optional bounded autonomy for the accepted app | Separate approved policy with scope/TTL/budget, retry ceilings, rate limits and verified revocation; no foundation expansion |

Required tests, **not yet implemented or run**:

- [ ] Anonymous, wrong issuer/audience, expired/revoked identity and unknown actions deny.
- [ ] Wrong cluster/app/env/path/controller and attempts to modify protected files deny.
- [ ] Self-approval, stale/expired/reused approval and changed head/digest/base deny.
- [ ] Prompt-injected instructions in logs, catalog, app output or NATS payloads
      cannot alter policy, targets, commands or approvals.
- [ ] Duplicate jobs, concurrent agents, queue redelivery and crashes around PR
      creation/merge/result recording do not create uncontrolled repeated effects.
- [ ] Network/policy/journal failure and stale observations fail closed; resuming
      reconciles actual state rather than blindly rerunning a mutation.
- [ ] Secret values never reach frontend, tool responses, logs, Git or audit receipts.
- [ ] Non-admin human, one authorized agent and another unapproved identity show
      real cross-role/cross-app isolation, not only unit-test assumptions.
- [ ] Exact revision/digest plus application success and supervised recovery pass live.
- [ ] Stop/revocation works across all replicas and queues; existing read-only
      Kubernetes/MCP/host guards and foundation boundaries still deny writes.

Home/edge/physical-control extensions are later work: prove local identity,
recovery and WAN-loss behavior on actual hardware. Initial cloud GitOps requires
its reachable Git authority; do not claim offline execution by queuing commands.
Robotics/vehicles/site equipment retain local safety interlocks and separately
approved action scopes; Kubernetes readiness never authorizes physical actuation.

## Out of scope and related work

No platform replacement, OpenChoreo migration, KubeVela/Crossplane installation,
new build plane, runtime deployment, permission or secret change in this PR.
Existing plans remain separate: [#328](https://github.com/onemindos/omos/pull/328)
(workflow plane), [#332](https://github.com/onemindos/omos/pull/332) (Legacy/LibreChat),
[#334](https://github.com/onemindos/omos/pull/334) (ingress lifecycle), and
[#284](https://github.com/onemindos/omos/pull/284) (home expansion).
No need to deploy all of those to test offline contracts; activation dependencies
must be demonstrated, not assumed. Phase 6 remains In progress, Phase 7 Not started.
