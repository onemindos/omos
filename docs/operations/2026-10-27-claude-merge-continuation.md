# Claude continuation after catalog PR #505 — 2026-10-27

## Current owner instruction

Zeus's latest instruction is: "merge everything yourself please", followed by
"lets get it out of the way and finish the work no more lagging on it do it".

**Codex owns all four remaining original PRs: #505, #535, #533 and #532.**
This supersedes the earlier #505-only and Hermes-only assignments. Other
harnesses must not edit or merge these branches concurrently.
Owner approvals for dcarku's original 46 PRs #495–#540, including **Telegram**,
are recorded in [PR #542](https://github.com/onemindos/omos/pull/542). Do not ask
again for settled decisions. Technical, credential and release gates still apply;
continue past a blocked PR and report actual blockers together.

## Verified merge queue

GitHub readback on 2026-10-27 (America/New_York), with main at `3495a614`,
confirmed **42 of the original 46 PRs merged**. The only open originals were
#505, #532, #533 and #535. Check live state before starting work: this is an
observation, not a claim about later merges.

| Assignment | Remaining work before merge |
| --- | --- |
| Codex: [#505](https://github.com/onemindos/omos/pull/505) | Integrate current main, preserve catalog mappings and generated data, pass final tests and CI, then merge only this PR. |
| Codex next: [#535](https://github.com/onemindos/omos/pull/535) | #534 has merged; #535 now targets main and all four Quality Gate checks pass. Check authorized catalog Locations registration records, including pending/unprocessed targets, for `services/onemind-portal/catalog-info.yaml`; resolve duplicate registration before merge. Processed entities alone do not complete this check. Rerun focused discovery/descriptor tests against final main. |
| Codex next: [#533](https://github.com/onemindos/omos/pull/533) | Keep Telegram. Obtain the intended dedicated alerts bot token and chat ID through authorized secret access, seal for the correct cluster, validate and merge. Confirm alert delivery after reconciliation. If credentials are unavailable, hold and report exactly what is missing; approval does not manufacture a bot or chat access. |
| Codex release: [#532](https://github.com/onemindos/omos/pull/532) | #531 pages have merged. After #505's build repair, build/publish a portal image containing the intended source/pages, verify it, and use a reviewed GitOps digest pin. Merge the links with that pin or after all 12 pages are verified available. Keep the release hold until then. |

All three other PRs were Draft, targeted main, were mergeable and had four
successful Quality Gate checks at readback. Those checks do not satisfy the
specific registration, credential or image gates above. They require no new
personal preference from Zeus. Bot creation or private credentials can still
require his access if they are not already available to the agent.

## #505 integration evidence

The final preparation integrates main `3495a614`, including the newly merged
NATS-streams Argo application and workload retirements. The catalog maps **60
source Argo Applications**; the static inventory baseline contains **409
Component/Resource entities**. Keep the generated source/mounted catalog copies
identical. Do not restore stale 59-app or 411-entity assertions from an older head.

The image's catalog/topology gate passes **66/66 tests** on supported Node
v24.22.1. Independent source review found no additional blocker: all seven new
AsyncAPI definitions parse without errors and no prior app identity is removed.
Final current-main inventory, YAML, render and GitHub receipts are recorded in
the PR before merge. Refresh the exact head and CI rather than treating this
preparation record as proof of a subsequent build or deployment.

The catalog ConfigMap can roll the existing portal image through Argo CD.
Source integration does not publish an image or complete release acceptance.
After reconciliation verify portal readiness and catalog processing separately.

## Completed work: do not reapply old scratch fixes

#502 and #506 have now merged, as have #499–#504, #507–#531, #534 and
#536–#540. Earlier preparation notes describing their integration as pending
are historical. Do not reapply old temporary patches to current main.

The #502 review caught the pinned OpenCost `/allocation` versus `/allocation/`
proxy redirect issue. The #506 review exercised multi-node NATS failover,
per-upstream auth headers, missing configuration failure and redaction while
preserving expanded app coverage and write opt-in. Check the merged behavior
when doing release acceptance; do not overwrite the later integrated MCP work.

## Release boundaries

Keep all credential material out of Git and tool output. Use Argo-owned GitOps
configuration; do not manually apply workloads, run Terraform apply or execute
proposed storage patches as validation. Image publication, digest pinning,
discovery-provider activation and Authentik group sync activation are distinct
from the source merges. Preserve disabled defaults and existing access rules.

The older open dependency/storage PRs outside #495–#540 are not included in
the 42/46 count or automatically declared ready by the batch approval. Review
their scopes, compatibility and actual CI before recommending another queue.

Local raw receipts are under `/tmp/omos-app-pr-preparation-20261027/`; this
tracked handoff is sufficient if Claude runs elsewhere. The initial preparation was limited to #505; the latest instruction above
authorizes completion of the remaining queue.
