# Owner-item merge preparation — 2026-10-27

Zeus asked Codex to resolve the straightforward items from dcarku's 46-PR batch,
prepare them for merge and leave Claude a checked continuation record. This is
source/review preparation. No PR was merged, feature activated, live Secret
changed or deployment synced by this task. Phase 6 remains In progress; Phases
7/8 remain Not started. Follow the existing layer and merge rules in
`PLATFORM-ROADMAP.md`, `CLAUDE.md` and `CONTRIBUTING.md`.

## Decisions and source candidates

The recommendations below implement the user's delegated routine choices. They
do not approve destructive retention, new client access, directory disclosure,
policy enforcement or credential/billing changes.

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
6. For #539, keep the current Hermes-only deployment until additional client
   identities/access are chosen. Later generate distinct tokens without exposing
   them, seal/mount the file, build/pin the MCP image, allow only intended client
   paths and test anonymous denial and per-client audit attribution. All clients
   retain the same existing tool permissions.

## Decisions still needing Zeus

- #514: choose the CoT retention window **before merge**; merging enables actual
  deletion on the next scheduled run, rather than a separate feature toggle.
- #523: confirm Headlamp's shared ServiceAccount may become view-only and decide
  whether portal health-check access stays.
- #533: choose the alerts channel/destination and supply its credentials through
  the secure local workflow. #540 routing alone does not deliver alerts to a person.
- #537: approve the Enforce policies and decide the bootstrap-race disposition.
  Recommendation: resolve namespace-label bootstrap ordering before enforcement;
  no admission-policy change is authorized or implemented here.
- #538: choose directory groups and approve their visibility before enabling
  synchronization. Keep it disabled until its credential, access binding,
  working stale-sync alert and image gates pass. Its larger source change was
  not fully reviewed in this preparation pass.
- Additional #539 clients: identify them before onboarding; source merge does
  not require new clients or credentials.

Local raw logs and bounded review notes are under
`/tmp/omos-pr-merge-prep-20261027/` on the preparation Mac. This tracked document
and the updated PR descriptions carry the durable summary if temporary files
are unavailable to Claude.
