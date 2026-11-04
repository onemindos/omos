# Backstage trust and recovery receipt workstream

Status: offline implementation for human review. No deployment, migration,
credential, role, access, storage or billing change. Phase 6 remains In progress;
Phase 7 remains Not started. The coordinator owns the shared roadmap update.

## Reconciled evidence and remaining gates

The October 19 [integration receipt](phase6-portal-integration.md) records a
successful historical release, including isolated fixture checks and selected
production checks. Its opening “Phase 6 complete” is historical scope, superseded
by the current [roadmap](../../PLATFORM-ROADMAP.md) and
[readiness tracker](../roadmap/migration-readiness-followups.md). It does not
certify every current session, role, application or recovery path.

| Area | Existing evidence in source | Remaining boundary |
|---|---|---|
| Release and denials | Scoped release receipts record exact revisions, two replicas, anonymous denial, machine catalog reads and write denials; [latest topology receipt](native-service-topology-2026-10-20.md) preserves those boundaries | Re-observe against the reviewed target; a source test or old receipt does not establish current runtime acceptance |
| Session lifecycle | [Follow-through](portal-acceptance-followthrough-2026-10-20.md) records refresh 200 on both replicas, recovery after 401 and upstream/cache regression cases | Controlled isolated expiry/refresh/provider-failure tests exist; production expiry and full cross-replica lifecycle coverage remain unproven. Do not deliberately expire the current user |
| Roles | Actual policy intersects upstream authorization and denies non-admin mutation; actor/scope tests exist | Another existing authorized human role has not been live-tested. Do not manufacture users/roles to tick a box |
| Portal backup | Replacement `backup-postgres-backstage` reports isolated authenticated archive restoration; historical suspended job is superseded | Fresh exact-scope restore evidence and independently recoverable identity/key inputs must be reviewed; a schedule or sync is insufficient |
| Grafana | [Durable acceptance](grafana-durable-visibility-2026-10-20.md): retained PVC, migrated Viewer identity, enabled daily backup, encrypted remote download/decrypt/isolated SQLite restore across 87 tables. [Folder fix](grafana-native-folder-2026-10-20.md) accepted separately | Earlier 401/emptyDir findings are repaired in those receipts. Token expiry monitoring/rotation and independent recovery remain separate; durable storage does not remove lifecycle obligations |
| Independent recovery | [Backup review](portal-acceptance-followthrough-2026-10-20.md) records the suspended offsite copier with no successful run; [offsite design](../../platform/backups/OFFSITE-COPY.md) describes required destination/access preparation | Same-provider/region Spaces restoration is valuable but does not prove independent offsite or universal hardware/WAN-loss recovery. No destination, keys or storage may be created by this work |

Existing source regressions live in
`services/onemind-portal/plugins/plugin-onemind-permissions-backend/src/`:
`sharedCache.test.ts` covers cross-replica adapters, renewed TTL, actor/scope
isolation, expiry and provider/database failure; `policy.test.ts` and
`hostGuard.test.ts` cover authorization intersections and backend boundaries.
These are source/fixture evidence, not additional live role or expiry acceptance.
This work does not rerun or replace the portal's release test cohort.

## Receipt v1

`scripts/gitops/portal_acceptance.py` accepts one local, redacted JSON file.
It reads no referenced evidence, environment credentials or tokens and performs
no network, subprocess or write operation. Evidence references are opaque: URLs
are not fetched and paths are not opened. Never place tokens, signed URLs,
credential values, raw session material or personal payloads in a receipt.

Exact fields (unknown fields are invalid at every object level):

- Root: `schemaVersion` = `onemind.portal-acceptance.v1`; `commit` = full lowercase
  40-character Git SHA of the reviewed GitOps revision; `observedAt` = RFC3339
  timestamp with timezone; `scope` = `portal-trust-recovery`; `checks` = array.
- Each check: `name` from the table below, `status` exactly `passed`, `failed` or
  `not-tested`, its own `observedAt`, and `evidence`.
- Each evidence object: `kind` exactly `observation` or `fixture`, nonempty `ref`
  (a redacted archive/document identifier, optionally with an anchor), and
  `sha256` = `sha256:` followed by 64 lowercase hex characters identifying the
  referenced artifact bytes. For `not-tested`, reference the hashed limitation
  record explaining the unexecuted check. A hash is required even then.

Do not rewrite historical receipts to refresh timestamps. Preserve an immutable
copy and its referenced artifacts under the operator's reviewed retention/access
process. Each check is scoped to the root commit and scope; historical evidence
from another release belongs in a separate receipt. A release check's artifact
must bind source/image digest, actual GitOps revision, replica identities and
observations. This validator does not infer those facts from a Git SHA.

| Required check name | Evidence needed before recording an observation pass |
|---|---|
| `release.exact-revision-replicas` | Reviewed revision, owning Argo app, actual imageIDs and both ready replicas |
| `security.anonymous-denied` | Anonymous protected API denials on both replicas |
| `security.machine-read-only` | Existing authenticated machine read success, MCP allowlist and write/platform/proxy denials |
| `security.backend-mutations-denied` | Direct backend and Kubernetes boundary denials, not merely hidden buttons |
| `session.refresh-recovery-both-replicas` | Refresh/sign-in recovery correlated across both replicas without recording tokens |
| `session.expiry-fail-closed-both-replicas` | Authorized observed expiry boundary on both replicas; isolated expiry regression remains `fixture` |
| `session.provider-failure-fail-closed` | Failure cannot resurrect expired authorization; isolated provider-failure tests remain `fixture` |
| `roles.existing-human-isolation` | Reviewed existing human role differences and direct backend isolation; no artificial enrollment |
| `roles.actor-scope-isolation` | Cross-actor/scope boundary observations; mocked adapter tests remain `fixture` |
| `recovery.portal-isolated-restore` | Authenticated archive download/decrypt and isolated database restore verification |
| `recovery.grafana-durable-restore` | Durable database/identity plus remote isolated restore evidence; not just PVC Bound |
| `recovery.independent-backup-restore` | Independently recoverable copy and isolated restore with explicit failure-domain limits |
| `recovery.identity-key-availability` | Redacted human verification of independently available recovery inputs; never key contents |
| `recovery.reviewed-rollback-compatibility` | Exact compatible image/config recovery plan and data/schema limits |

This is a deliberately named trust/recovery review bundle, not every Phase 6
criterion or a universal per-app migration policy. An independent-backup gap
keeps this bundle incomplete. The readiness tracker still requires the selected
app's actual dependency, security, state and rollback assessment; optional unrelated
dashboard/launcher work is not automatically a blocker for every app.

For example, this partial receipt is intentionally **incomplete** (replace the
illustrative digest with the hash of a real redacted limitation record):

```json
{
  "schemaVersion": "onemind.portal-acceptance.v1",
  "commit": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "observedAt": "2026-10-21T01:00:00Z",
  "scope": "portal-trust-recovery",
  "checks": [{
    "name": "roles.existing-human-isolation",
    "status": "not-tested",
    "observedAt": "2026-10-21T01:00:00Z",
    "evidence": {
      "kind": "fixture",
      "ref": "redacted-limitations.json#existing-human-role",
      "sha256": "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    }
  }]
}
```

## Validation and review behavior

```sh
python3 scripts/gitops/portal_acceptance.py /absolute/path/receipt.json \
  --expected-commit FULL_REVIEWED_GITOPS_SHA
python3 -m unittest discover -s scripts/tests -p 'test_gitops_portal_acceptance.py' -v
```

The expected commit is supplied explicitly by the reviewer; no Git/cluster query
is made. Mismatch is incomplete. Review freshness defaults to 24 hours, optionally
`--max-age-hours 1..168`; this is a review threshold, not a health SLA. Each check
has its own freshness bound, so a new envelope cannot launder old observations.
Future observations, checks later than the envelope, malformed/timezone-free
timestamps, duplicate JSON fields/checks, unknown gates, booleans as statuses,
non-finite JSON numbers and missing/invalid evidence references are invalid.

Exit 0 means all named checks have fresh supplied observation passes at the
expected commit. Exit 1 means failed or incomplete: any explicit failed check
dominates, otherwise any missing, not-tested, stale or fixture check stays
incomplete. Exit 2 means invalid JSON/schema/input. Output lists failed, missing
and incomplete check names plus commit/freshness reasons. Supplied reference
values are never echoed. A canonical receipt SHA-256 binds the supplied record
for review; artifact hashes bind referenced bytes. Neither proves authenticity,
artifact availability, truth, retention or that a fixture was honestly labelled.
`evidenceContentsVerified` is always false: the human must inspect retained
artifacts and their hashes. There is no automatic artifact fetch or signature
verification, and no claim of tamper-proof storage.

Even a complete observation receipt always returns `humanReviewRequired: true`,
`phase6Done: false`, `migrationAuthorized: false`, `actionAuthorized: false`.
No fixture or missing check can produce blanket Phase 6 Done. This standalone
review tool is not wired into promotion, merge, runtime policy or action execution.
Do not convert exit 0 into an approval or migration permission.

## Draft #339: first bounded offline implementation recommendation

Read-only inspection of [draft #339](https://github.com/onemindos/omos/pull/339)
(`codex/permissioned-ai-actions-plan`) finds planning-only A0–A5, with no selected
pilot or deployed executor. Its proposed contract requires exact-head/hash-bound
approvals, expiry/revocation, target isolation, audit and failure reconciliation.
This work does not modify that draft or approve its activation.

First task: implement A0's versioned request/approval schema and a pure offline
authorization evaluator with a fake clock and fixed target fixtures. Bind request,
base/head/diff/artifact digest, target and policy version; deny unknown actions,
protected foundation targets, self-approval, changed heads, stale/expired/revoked
or reused approvals and cross-target replay. Treat unavailable policy/journal or
unknown reconciliation state as denial. Tests should prove identical decisions
for UI/API/tool-shaped inputs without connecting an executor, GitHub, NATS or a
cluster. In-memory fakes do not prove A1 durable crash/restart semantics; journal
persistence and those tests remain a separate step. No new identity, credential,
merge right or live role/expiry exercise is required to implement that contract.
