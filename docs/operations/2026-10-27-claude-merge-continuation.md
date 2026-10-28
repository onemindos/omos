# Claude continuation after catalog PR #505 — 2026-10-27

## Current owner instruction

Zeus narrowed Codex's assignment to: "finish off with 505 and dont merge
anything else and update the claude code to finish the do the rest".

**Codex finishes and merges #505 only. Claude Code handles the remaining PRs.**
The earlier broad Codex application-lane claim is superseded. Do not wait for
Codex to prepare the remaining branches. Owner approvals for dcarku's original
46 PRs #495–#540, including the Telegram choice, are already recorded in
[PR #542](https://github.com/onemindos/omos/pull/542). Continue past an actually
blocked PR and report remaining blockers together; do not ask for settled
approvals again. Keep one final merger and do not edit a handed-off branch
concurrently.

## #505 completion evidence

#505 was integrated with main `988af5d8` during preparation; refresh its final
head/main and GitHub checks before merge. On supported Node v24.22.1:

- The image's catalog/topology Node gate passed **66/66 tests**, repairing the
  stale generated mapping/runbook and Argo application-count failures.
- Inventory Jest's first full run passed 27 tests across six suites and found
  one stale startup assertion: 408 entities instead of the reviewed 411. The
  new Hermes Fresh, MCP and Vaultwarden workloads explain the three additions.
  This PR updates that assertion; the affected startup suite then passed 1/1,
  completing the 28 existing inventory tests across the full and focused runs. The
  baseline/flag-absent/flag-false entity equality checks remain intact.
- Independent review found no additional source/security blocker. All seven
  new AsyncAPI definitions parse without errors, and no prior app identity was
  removed. Optional AsyncAPI metadata/style advice is not a parser error.

Local test dependencies were isolated. The existing external portal dependencies
were reused; missing pinned backend-test-utils/testcontainers dependencies and
a Node24-compatible SQLite binary were supplied only under a temporary export.
No shared dependency installation, credential, provider or billing setting was
changed. Full image publication, digest pinning and live acceptance remain
separate release steps. The catalog ConfigMap can roll the existing portal
image on Argo reconciliation.

## Fix #502 before its source merge

The actual pinned Backstage proxy forwards the configured allocation endpoint
to `/allocation/`, while the pinned OpenCost router registers `/allocation`.
The upstream returns an unusable 301, and the proxy does not expose a useful
redirect to its caller. A scratch-only fix was tested against that proxy:

```yaml
'/opencost/allocation':
  target: http://opencost.opencost.svc.cluster.local:9003
  pathRewrite:
    '^/api/proxy/opencost/allocation/?(?=\?|$)': /allocation
  credentials: require
  allowedMethods: [GET]
```

Apply this to the existing proxy entry in
`platform/backstage/app-config.cluster.yaml` on #502. The tested override
preserves `/allocation?window=1d&aggregate=namespace` exactly; anonymous denial,
GET-only behavior, credential stripping, sibling rejection and traversal
rejection still pass. **The fix is not applied to #502's branch.** It is stacked
on #500, so integrate the parent and obtain main-targeted CI afterwards.

## Reconcile #506 with the expanded MCP main

At reviewed main `988af5d8`, #506 has four conflicts:

- `platform/hermes/mcp/apps.yaml`
- `services/onemind-mcp/onemind_mcp/adapters/__init__.py`
- `services/onemind-mcp/onemind_mcp/config.py`
- `services/onemind-mcp/onemind_mcp/server.py`

An isolated resolution preserved main's complete seven-app catalog and adapter
registry, blocked/disabled dispositions, `auth_header`, write opt-in and
read-only registration metadata, while adding the NATS upstream list and node
attribution. Every NATS upstream must receive its configured auth header; other
adapters still receive their single upstream. Both connection paths must reject
an enabled app with missing `url_env`.

The isolated resolution passed **36 existing MCP tests plus 3 interaction
checks** for configuration failure, custom-header multi-node failover and
credential redaction. The coverage checker retained all 56 mapped source Argo
Applications. **No resolved file was applied or pushed to #506.** Reconcile with
then-current main and rerun its tests before using it.

On the preparation Mac, the four resolved copies and review patch are at
`/tmp/omos-app-pr-preparation-20261027/mcp506-resolved/` and
`/tmp/omos-app-pr-preparation-20261027/506-resolution-vs-main.patch`.
These are optional local aids, not required inputs: the conflict-preservation
rules above are the durable handoff. Paths/patch context may be stale after
other MCP PRs land. Also reconcile #539's token/audit changes with the final
NATS and expanded-app integration rather than overwriting either feature.

## Other bounded review results

These results apply to the reviewed heads and saved main `988af5d8`; they are
not final integration, image or deployment acceptance.

| PRs | Completed evidence | Remaining work |
| --- | --- | --- |
| #499/#500/#501 | Source/recovery review and applicable YAML/schema/Kustomize checks | Final main CI and rollout. Preserve any wanted old backup-job evidence before #501's prune. |
| #509–#513 | Source/backup-policy review, applicable lint/schema checks; #510's image indexes resolve for linux/amd64 | Preserve parent order #509 → #510 → #511; main CI after child rebases. Verify resource/PVC placement, CNI paths, inference and backups after rollout. |
| #524 | Existing four alias tests and YAML/schema/render checks pass | Exact hostname, real TLS/redirect/browser acceptance after rollout. |
| #531/#532 | Strict MkDocs and pinned TechDocs generation pass; all 12 catalog links resolve in the generated site | #532 stays release-held until an image containing #531 is published, pinned and available; prefer links with the pin. |
| #534/#535 | Discovery/descriptor contracts pass 2/2 and 4/4; YAML/render and virtual merge checks pass | After #534 lands, correctly rebase only #535's child commits onto main, obtain CI and check registered/pending/unprocessed Locations. |
| #538 | Five Jest suites, **101 tests**, scoped TypeScript and lint pass; four Python contracts pass | Source ships disabled. Actual intended groups, credential, binding refusal, tested stale-sync alert, image and acceptance are activation prerequisites. Preserve superuser/admin exclusions. |
| #507/#508/#539 | Virtual merges to saved main were clean; earlier #539 source review/tests remain recorded in #542 | This preparation pass did not rerun their final integrated suites. Refresh and validate against current MCP main; do not treat old CI as the new merge's result. |

All twelve reviewed portal/TAKML branches individually merged virtually with
saved main. #534/#538 also merged virtually together and their six Python
contracts passed. That does not prove arbitrary sequential merge order is safe.

## Resume order and boundaries

1. Confirm #505's final state and portal reconciliation. Fetch current main and
   preserve the owner-approval and scoped handoff records when resolving docs.
2. Fix #502; integrate parents before stacked children. Preserve generated
   catalog pairs and runbooks while handling overlapping catalog PRs.
3. Complete #506/#539 integration and final MCP tests/coverage before a separate
   image pin. No new client identity or additional permission is implied.
4. Finish the original #495–#540 queue under the approved decisions. Preserve
   real credential/release gates, keep Telegram and fix namespace-label
   bootstrap ordering before #537 enforcement.
5. For portal release, land the intended source/pages, run the supported-runtime
   release checks, publish the image and use a reviewed GitOps digest pin. Do not
   merge #532's links ahead of the pages. Verify 2/2 readiness, login and existing
   user/machine restrictions; record observed acceptance separately.

`make validate` has inherited whole-tree YAML debt; capture the actual failure
and pass changed-file gates without weakening checks. No source PR other than
#505 was merged or modified by Codex during this final preparation step. Local
raw receipts are under `/tmp/omos-app-pr-preparation-20261027/`; this tracked
summary is the continuation if Claude runs elsewhere.
