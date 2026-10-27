# Owner-item merge preparation — 2026-10-27

Zeus asked Codex to prepare dcarku's 46-PR batch and leave Claude a checked
continuation record. He subsequently approved all the remaining owner decisions
and retained Telegram. **Owner approval is settled for PRs #495–#540.** Read the
bounded authorization and approved decisions below before continuing. Technical
review, credentials, CI, dependency order and rollout acceptance still apply.
No PR was merged or live configuration changed by this recording task. Phase 6
remains In progress; Phases 7/8 remain Not started.

## Owner authorization — 2026-10-27

In this Codex conversation, after the four remaining owner requirements and the
foundation merge delegation were explained, Zeus stated:

> I APPROVE THEM AND KEEP THE CHANNEL TELEGRAM
>
> I APROVE EVERYTHING ELSE ARE WELL MAKE SURE CLAUDE KNOWS IT SO IT CAN BE READY FOR IT SO ALL OF THE PR'S IT DID IS READY FOR MERGE

Claude Code is authorized to review, fix, validate and execute sequential GitHub
merges on Zeus's behalf for dcarku's **46 PRs #495–#540**, including foundation
PRs in this batch. This is the human approval required by the default foundation
merge rule, not a standing grant for unrelated PRs or infrastructure. Use the
normal Git/Argo delivery path and satisfy the concrete technical gates below.
The handoff PR #542 carries this record.

Do not ask Zeus again for the settled choices. Handle routine fixes, builds,
rebases, supported-runtime tests, credential sealing from existing authorized
sources and release ordering. If one PR cannot proceed because an actual
credential/configuration is unavailable or a technical check cannot be repaired,
hold that PR, continue the rest and report the remaining blockers together at
the end. Approval does not create credentials, waive checks or make a skipped
PR merged. Do not mark every Draft ready before its own technical gates pass.

## Superseding merger assignment — 2026-10-27

Zeus subsequently authorized Hermes to execute the owner-approved merges recorded
in #542 on his behalf and designated **Hermes the lead and sole merger**. This
supersedes the earlier Claude merger assignment above, not the settled approvals
or technical gates. The initial sequence is **#542, #504, #497, #498, #496**.
Review each PR, refresh it against current `main`, resolve issues, verify final
CI, mark ready and merge sequentially with merge commits. Continue past blocked
PRs and report remaining blockers together without asking again for settled
approvals. **Keep Telegram.**

These are source merges: **do not run Terraform apply or execute the proposed
storage patches**. Zeus confirmed an existing agent is already preparing #505;
do not duplicate or interfere with its work. That harness prepares and tests #505
and the application PRs, coordinating dependencies with Hermes; final merges
remain Hermes's responsibility. Source integration is not rollout acceptance.

## Decisions and source candidates

The table preserves the initial source-preparation dispositions. Owner decisions
have now been approved as recorded below; technical integration/release holds
still apply. Do not expand beyond the reviewed PR scopes or change billing or
provider settings to work around release failures.

| PR | Disposition | Next action |
| --- | --- | --- |
| [#539](https://github.com/onemindos/omos/pull/539) | Prepare for source merge. Existing Hermes authentication stays compatible. No new client is activated. | Review the exact head and CI, then merge through the normal workflow. Leave token generation, image pinning and new-client network access for a separate activation change. |
| [#524](https://github.com/onemindos/omos/pull/524) | Prepare the private `portal` redirect alias, keeping `onemind-preview` canonical. No extra OIDC callback is required. | Recheck the exact hostname, review final CI, merge, then observe Argo and test the real alias/TLS/browser path. |
| [#534](https://github.com/onemindos/omos/pull/534) | Prepare source merge with discovery restricted to `onemindos/omos`, branch `main`, literal root descriptor, existing read-only App and catalog rules. | Keep the existing config-with-source design. Its merge rolls the current portal image; discovery starts with a later pinned image containing the module. No wider organization/App scope is approved. |
| [#535](https://github.com/onemindos/omos/pull/535) | Keep Draft: integration hold, rather than a new owner preference. | After #534, rebase only the child commits onto current main, retarget to main, run CI and repeat the registered-Location check. Never merge this child into #534's feature branch. |
| [#532](https://github.com/onemindos/omos/pull/532) | Keep Draft: documentation release hold. | Land #531's pages and the build prerequisite, publish/pin an image containing them, then merge the catalog links. Prefer including the links in the digest-pin PR so they cannot precede the pages. |

All five descriptions/titles are updated to distinguish source readiness from
integration/release holds. #524's README and roadmap are corrected in
`f1611953` and `fdf5e3f0`; the latter also documents the `/healthz` exception and
the complete separate canonical-host migration. Runtime manifests are unchanged
from the original alias implementation.

#524, #534 and #539 were marked Ready for review after their final heads passed
all four Quality Gate checks: lint/validation, Kubernetes schema validation,
Helm rendering and secret scanning. #532/#535 remain Draft. This status records
source-review readiness, not a merge approval or deployed acceptance.

## Checks performed

Reviewed and tested heads:

| PR | Head | Fresh focused evidence |
| --- | --- | --- |
| #539 | `15e7ab91b1dec3a7f187396126a252515f94528d` | 34 tests passed in an isolated Python 3.12.14 environment using the branch's pinned development requirements. Authentication, legacy compatibility, audit isolation, redaction and source-only activation boundary reviewed. |
| #534 | `b8b8c495d1a5e8263d8f35f4587b1627a309e67f` | 2 discovery contracts passed; YAML lint and diff check passed. Anchored OMOS scope, existing App restrictions and catalog allowed kinds reviewed. |
| #535 | `30d1e6ec34291b8a07b1a26c6c7abebf61f69073` | 4 descriptor/discovery contracts passed; YAML lint and diff check passed. No CI is reported while it targets the feature branch. |
| #524 | `fdf5e3f0` | 4 alias contracts passed; changed-YAML lint passed; kubeconform validated all 4 new manifest resources. Kustomize rendered and the generated ConfigMap reference resolved. |

The alias was also exercised in Docker using its exact pinned nginx image:
non-root uid 101, read-only root, no capabilities, no network and no host ports.
The health endpoint returned 200 and a browser path/query returned 302 to the
unchanged canonical origin with both preserved. The test container was removed.
This proves the container behavior, not live TLS, CNI routing or browser login.
The PR commit range passed gitleaks; no baseline or validation rule was relaxed.
An independent source review found no blocker; its small documentation correction
was incorporated before the final push.

Read-only preflight found:

- No exact `portal` device name among 43 devices visible to the local tailnet
  identity. Visibility is ACL-scoped, not an authoritative admin inventory.
- No `portal` hostname claim among 33 live ingresses in
  `do-nyc3-omos-cloud`.
- All 8 processed catalog Location entities were returned, none targeting
  `services/onemind-portal/catalog-info.yaml`. No
  `component:default/onemind-portal` entity was returned. No registration was
  changed. Repeat before #535 merges; pending/unprocessed registrations were not
  certified absent.

Before each merge, fetch again and verify CI belongs to that exact final head.
Ready-for-review status is not production acceptance or a substitute for review.

## Existing validation and release blockers

`make validate` was run and fails inherited whole-tree YAML debt, beginning at
`platform/onetak/media-mtx.yaml:86`, with additional findings in unchanged files
such as `platform/takml/ghcr.sealedsecret.yaml` and Matrix manifests. Changed-file
checks pass. Do not call the whole-tree gate green or weaken it to merge.

The main-checkout portal command
`node --test scripts/test-audit-catalog.cjs scripts/test-topology.cjs` returned
62 passes and 3 failures at main `f0f9b464`:

1. The hard-coded source Argo Application count expects 56; Git has 59.
2. The generated inventory placement bundle is stale.
3. The generated Karakeep runbook differs from the catalog inputs.

That run used host Node v26.5.1, not the documented release runtime; repeat the
release checks with the supported runtime. The static mismatches are concrete
and the existing #505 is the intended repair candidate. Review and verify its
exact rebased head rather than disabling the tests or duplicating its work.
No portal image was built or published by this preparation task. The existing
PR descriptions record hosted image-workflow billing restrictions; this task did
not recheck billing or change it. Use the documented manual release path if
necessary, without bypassing tests.

## Claude continuation order

1. Fetch, preserve unrelated local work and read the five updated PR descriptions.
   Review/merge one candidate at a time through the normal workflow. For #524,
   recheck hostname availability and observe the alias after Argo reconciliation.
   For #534, observe the same-image ConfigMap rollout and verify that the running
   image still lacks the provider before describing discovery as inactive.
2. Review #505 or its equivalent build repair. Land #531 before the documentation
   image build. Keep catalog copies and generated artifacts consistent while
   integrating the other overlapping catalog PRs.
3. Record #534's final feature-branch head before merging or deleting its branch.
   After it lands, fetch main and rebase only #535's child commits onto current
   main (account for squash merging), push with lease, retarget to main and
   obtain fresh CI. Check registration records, including pending/unprocessed
   targets via the authorized catalog Locations API, before merging; the
   processed-entity observation above does not complete that check. Inspect any
   conflict before unregistering anything.
4. Use `services/onemind-portal/scripts/RELEASE.md` and `release-image.sh` to run
   the supported release gates and publish the reviewed source. Prepare a
   separate GitOps digest-pin PR. Include #532's links with that pin, or land
   them only after the new pages are verified available. If the links are folded
   into the digest-pin PR, close #532 as superseded after that PR lands.
5. After image reconciliation, verify 2/2 portal readiness, normal sign-in,
   unchanged user/machine authorization, provider task success/no access errors
   and the new runbook pages. Record actual receipts; source merges alone do not
   complete these steps.
6. For #539, derive intended existing client identities from the actual approved
   client configuration, then generate distinct tokens without exposing them,
   seal/mount the file, build/pin the MCP image, allow only intended client paths
   and test anonymous denial and per-client audit attribution. Source merge does
   not require a new client. Leave the existing Hermes deployment intact if no
   additional client configuration is present; do not invent client identities
   or grant additional tool permissions.

## Approved decisions — do not request these approvals again

- **#514: 24-hour CoT retention is approved**, including deletion of older
  accumulated history on the scheduled run after merge. Verify the documented
  recovery/preflight and bounded SQL behavior before merging; preserve backups.
- **#523: the shared Headlamp ServiceAccount may become view-only; retain portal
  health-check ingress.** Keep user-owned admin access, merge #503 first and
  complete the child rebase/main-CI gates.
- **#533: keep Telegram.** Channel approval is complete. Use the intended
  existing alerts destination and valid bot credentials from authorized secure
  configuration, seal them and verify decryption and actual delivery. If the
  destination or token is unavailable, record that concrete credential/setup
  blocker and continue other PRs. Do not merge an unconfigured receiver: the
  required SealedSecret and delivery prerequisites remain. #540 routing alone
  does not deliver alerts to a person.
- **#537: Enforce is approved for `disallow-latest-image-tag` and
  `require-namespace-layer-label`. Fix the bootstrap namespace-labeling race
  before enforcement**, as recommended. Resolve the implementation through the
  scoped GitOps workflow, then obtain main-targeted CI and admission/rebuild
  verification. The other two policies stay Audit.
- **#538: source merge and the documented directory-sync scope/visibility are
  approved.** Use the existing intended portal-access groups from actual
  configuration; the reviewed users' names, emails and published group
  memberships may be visible as documented. Preserve the superuser-group
  exclusions. Do not invent a group list or broaden portal access. Source can
  merge disabled; enabling still needs valid credentials, verified bindings,
  a working stale-sync alert, the image release and actual acceptance. Keep it
  disabled if those inputs are unavailable. Its larger source change still needs
  full source review by the integrating agent.
- **#539: per-client credentials for the intended existing MCP clients are
  approved.** Discover their actual configuration, preserve current tool
  permissions and complete the implementation/acceptance steps above. No
  additional personal approval is needed for the source merge.

For the other PRs in #495–#540, owner authorization to review, repair and merge is
also complete. Draft state, red/missing checks, conflicts, parent dependencies,
unavailable credentials and image/rollout ordering are implementation gates for
Claude to resolve; they are not a reason to ask again for these owner choices.

Local raw logs and bounded review notes are under
`/tmp/omos-pr-merge-prep-20261027/` on the preparation Mac. This tracked document
and the updated PR descriptions carry the durable summary if temporary files
are unavailable to Claude.
