# Handoff: continuing the platform roadmap

**2026-10-27 merge-preparation update:** for the current dcarku PR batch, start
with [owner-item merge preparation](2026-10-27-owner-merge-preparation.md).
**Zeus has approved all owner decisions for #495–#540, retained Telegram and
authorized Claude to execute the batch's merges on his behalf.** The record
contains the exact approval, technical holds and release order. Continue past
an unavailable credential or unresolved technical blocker and report it at the
end; do not ask again for settled decisions. The historical
release observations below remain timestamped context, not current live state.

**Latest owner assignment (2026-10-27): Codex merges all four remaining original
PRs: #505, #535, #533 and #532.** Read the
[current continuation](../operations/2026-10-27-claude-merge-continuation.md).
This supersedes the earlier #505-only and Hermes-only assignments. Other
harnesses must not edit or merge these branches concurrently. Keep current-main
review and final CI. Telegram and the owner approvals remain unchanged. No
Terraform apply or storage patches follow merely from their docs merging.

Updated 2026-10-20 during the authorized portal/catalog stabilization. The authoritative tracker and
phase evidence live in [PLATFORM-ROADMAP.md](../../PLATFORM-ROADMAP.md).

## Where we are

Phases 0–5 are done. **Phase 6 stabilization is the first unfinished phase**;
its historical acceptance is preserved, but Zeus requested UI/catalog repairs
before further migrations. Source/image/inventory/composition PRs #322–#326 were the earlier stabilization.
The latest source #330, pipeline #335 and pin #336 are now merged/deployed:
`sha256:6c729bd20c0a0a40330be8f0c7bd5213c59c7540bd1ea54c38e05f8c63fc3e7c`,
2/2 healthy, 716 processed catalog entities and scoped live acceptance recorded in
[STAC release evidence](../operations/stac-operational-slice-2026-10-20.md).
The earlier [stabilization review](../operations/portal-stabilization-2026-10-20.md)
retains broader unresolved gates. Phase 7 remains not started; migration drafts #320/#321 stay unmerged.
Do not begin a later phase or treat this handoff as authorization to deploy it.
Broader token-lifecycle, launcher/empty-card and COP coverage gates remain in the
evidence; absent workflow plane stays explicit and build workflows disabled.
Zeus requested no further Ollama delegation: handle follow-through directly.

**Resume index for the full Kubernetes structure request:**
[readiness/ordered queue #331](https://github.com/onemindos/omos/pull/331),
[full structure #333](https://github.com/onemindos/omos/pull/333),
[public ingress lifecycle #334](https://github.com/onemindos/omos/pull/334),
[staged private Legacy/LibreChat #332](https://github.com/onemindos/omos/pull/332).
The #330/#335/#336 scoped release is verified; next review the remaining trust
gates and read-only foundation assessment; keep #284 as home expansion parent, revised against the Mac's actual
storage/power/runtime gates. Eight nginx-class Ingress resources still need a
maintenance-lifecycle plan. Legacy's configured inference failed with provider 402;
Mac internal free space was 13–14 GiB at that earlier observation; recheck before
allocating anything. Proposed Karmada/KubeEdge/Crossplane/
KubeVela are not installed/accepted capabilities. Draft/planning PRs do not deploy
or merge themselves, complete a phase or authorize credential/billing/permission
changes. Reconcile shared tracker edits before merging one PR at a time.

- STAC Browser is the phase-3 pilot; Phase 4 refactor/hotfix is deployed, portal 2/2 healthy.
- Phase 5: all nodes below 78% requested CPU/memory; authenticated Observer core and
  modules live, real logs/metrics in console and MCP, isolated backup recovery passed.
- Phase 6: existing OneMind Backstage 1.55.3 hosts OpenChoreo 1.3.1 core/observability;
  real Authentik, home/catalog/runbooks, Deploy/Definition/logs/metrics and both
  replicas passed. Shared DB capability cache fixes cross-replica authorization.
  The temporary console/candidate are retired; no storage or identity was deleted.
  See the Phase 6 integration runbook and completed #282 brief.
- Existing exceptions remain: platform-matrix OutOfSync/Healthy and platform-openclaw
  Synced/Progressing. These are not newly caused Phase 5 failures.
- GitHub later refused to start CI jobs because of an account payment/spending-limit
  error. Do not change billing or weaken merge requirements. Local validation evidence
  and existing whole-tree YAML lint debt are recorded in Phase 5.
- Recovery PVC observer-restore-20261019 was retained until Zeus approved retiring the
  OpenSearch and OpenChoreo observer recovery assets (2026-10-28); retirement runs through the
  `cleanup/retire-openchoreo-recovery-assets` PR runbook. Never delete or restore over production storage.

Latest follow-through: [portal acceptance review](../operations/portal-acceptance-followthrough-2026-10-20.md).
Argo locator works in the production overlay. The historical Grafana search 401
is repaired by #340–#342; see the accepted durability receipt below. Two suspended backup jobs are intentionally replaced; independent
offsite copy remains unconfigured. Do not infer all-green health from portal readiness.

Existing pilot integration: [STAC operational slice](../operations/stac-operational-slice-2026-10-20.md).
Its scoped image/pin/live gates are now recorded as passed. Three ready workloads do
not hide the actual collections/search 500s or unknown data freshness. No schema
initialization is authorized by this integration; it needs its own recovery gates.

Latest native visibility first win: [#337](https://github.com/onemindos/omos/pull/337)
is merged/deployed; portal 2/2 generation 73 on the same digest, DB unchanged,
NATS 3/3. Kubernetes controller matching and explicit Argo/runbook/source links
are browser-verified for these first services. [Evidence and remaining Grafana
credential/proxy work](../operations/native-portal-visibility-2026-10-20.md).
The Grafana 401 recorded in that historical evidence is now repaired through the
explicitly authorized durable migration below.
This is not full catalog/COP coverage or Phase 6 completion.

## Briefs and draft PRs

Full-stack structure is in [PR #333](https://github.com/onemindos/omos/pull/333):
[operating structure](../architecture/kubernetes-operating-structure.md) and
[timestamped baseline](../operations/kubernetes-structure-baseline-2026-10-20.md).
The #330/#335/#336 portal release is already verified; do not rebuild or repin it
just to resume the plan. Next structural work is #334's read-only ingress assessment
and K02 ownership/recovery inventory; remaining Phase 6 trust gates stay explicit.
All 20 categories, actual versus candidate tools, integration gaps, app contracts
and gated home/agent/edge/ML work are covered. #332 remains staged, not activated.
No new platform, phase completion or migration is authorized by this plan.

Before selecting a migration, use the
[migration-readiness and follow-up checklist](migration-readiness-followups.md)
in [PR #331](https://github.com/onemindos/omos/pull/331).
It separates #330 release/core acceptance and per-app safety gates from parallel
STAC/Grafana/launcher/offsite/security/validation follow-ups. Existing #328, #283,
#320 and #321 are linked there; do not duplicate them or infer a cutover from
the planning PR. No phase is marked complete merely by recording the backlog.

| Phase | Brief | Goal | Draft PR |
|---|---|---|---|
| 3 | [phase-3-pilot-app.md](phase-3-pilot-app.md) | One small app as an OpenChoreo Component | [#279](https://github.com/onemindos/omos/pull/279) |
| 4 | [phase-4-portal-refactor.md](phase-4-portal-refactor.md) | Portal features into `plugins/` packages | [#280](https://github.com/onemindos/omos/pull/280) |
| 5 | [phase-5-capacity-observability.md](phase-5-capacity-observability.md) | Capacity, then the observability plane | [#281](https://github.com/onemindos/omos/pull/281) |
| 6 | [phase-6-one-portal.md](phase-6-one-portal.md) | OpenChoreo plugins in our portal | [#282](https://github.com/onemindos/omos/pull/282) |
| 7 | [Migration plan in #283](https://github.com/onemindos/omos/pull/283) | Move Layer 3 apps project by project | [#283](https://github.com/onemindos/omos/pull/283) |
| 8 | [Expansion plan in #284](https://github.com/onemindos/omos/pull/284) | Home server and Hermes | [#284](https://github.com/onemindos/omos/pull/284) |

Completed briefs record evidence; future draft PRs initially add only a brief. Work the phase on that PR's branch (or a branch from it),
keep the PR a draft until the gate passes, then mark it ready and merge. Every PR also updates the
tracker in `PLATFORM-ROADMAP.md`.

## Access (set up through Phase 5)

- Read-only cluster and portal access for AI agents: `scripts/ai-operator/setup.sh`, documented in
  [`docs/operations/ai-operator-access.md`](../operations/ai-operator-access.md). MCP servers:
  `onemind-kubernetes`, `onemind-backstage`, `onemind-openchoreo` and `onemind-openchoreo-obs` (17 Observer query tools; writes refused).
- Changes to the cluster go through Git and Argo CD. Mutating `kubectl` needs the admin context
  `do-nyc3-omos-cloud` and is only for incidents and one-off operations (for example re-running an
  Argo hook).

## Rules that bit us

1. **Never name a node.** DOKS replaces nodes. Use the node-pool label or pod affinity.
2. **`Recreate` for any Deployment with a ReadWriteOnce volume**, or a rolling update deadlocks.
3. **Check the gate against the live cluster** (read-only) before marking a phase done.
4. **One PR per step**, CI green (`make validate` fails locally on old YAML in
   `platform/onetak` and `platform/takserver`; CI lints only changed files).
5. **No plaintext secrets in Git.** Generated secrets are created by Jobs at runtime, others are
   SealedSecrets (`kubeseal`, see `scripts/prepare_ecosystem_secrets.py`).
6. **Layer 1 is never moved into OpenChoreo** (see the layers in the roadmap). An AI gets read-only
   access above Layer 3 and proposes changes as pull requests.
7. **Argo hooks re-run only when a sync runs.** To re-run the identity Job:
   `kubectl -n argocd patch application platform-ecosystem-resources --type merge -p '{"operation":{"initiatedBy":{"username":"agent"},"sync":{}}}'`.
8. **Do not use `path` as a shell variable in zsh** (it is tied to `PATH`).
9. **Do not sit on a stuck Argo operation**: if a rollout is stuck on an old revision, terminate it with
   `kubectl -n argocd patch application <app> --type merge -p '{"status":{"operationState":{"phase":"Terminating"}}}'`.

## Useful facts

- OpenChoreo 1.3.1: control plane in `openchoreo-control-plane`, data plane in
  `openchoreo-data-plane` (registered as `ClusterDataPlane` `default`). Apps get namespaces
  `dp-<project>-<env>-*`. The upstream install guide for v1.3.x is the reference; deviations are
  tabled in the roadmap, phase 2.
- Gateway: Cilium Gateway API (CRDs v1.2.1, class `onemind-private`, Tailscale LoadBalancer).
  Apps reach hostnames like `<env>-<namespace>.apps.openchoreo.100-110-185-30.nip.io`, tailnet only.
- Component type to use: `onemind-service` (`platform/openchoreo/platform/onemind-service.yaml`).
- Authentik: provider `onemind-openchoreo`; tokens have no `client_id` claim, so authorization
  bindings use `sub`. Service accounts `ai-operator` and `backstage-catalog` are created by
  `platform/ecosystem/resources/ai-operator-identity.yaml`.
- Namespaces and layers: the table in the roadmap.

## Prompt for the next agent

See [HANDOFF-PROMPT.md](HANDOFF-PROMPT.md).

## Latest scoped release — 2026-10-20

Source #330, portable pipeline #335 and image-pin #336 are merged; actual Argo
revision `1545b50f4fb5521dbdfc9eb22fa302cf2e4b45d6` and both new replicas verified.
See [STAC acceptance](../operations/stac-operational-slice-2026-10-20.md) and
[migration-readiness tracker](migration-readiness-followups.md) for exact digest,
716-entity topology audit, browser checks and per-pod negative/MCP checks.
Old `af7770fb` / `fab652...` candidate is superseded. The pipeline was really used;
first profile is portal, retirement is plan-only. Hosted CI remains unavailable,
not green; old whole-tree YAML debt remains. Phase 6 is still In progress and
Phase 7 Not started: do not turn this successful portal release into approval for
STAC schema changes, service migrations or new infrastructure controllers.


Latest Grafana follow-through: #340/#341/#342 merged/reconciled; scoped acceptance
passed at 22:29–22:35 UTC. [Durability/native visibility receipt](../operations/grafana-durable-visibility-2026-10-20.md):
retained 5 GiB PVC, same Grafana image/Deployment UID, migrated read-only Viewer
credential, verified encrypted remote SQLite backup/restore and enabled daily job;
portal generation 74, 2/2 same image, DB unchanged, NATS 3/3. Native dashboard links
and actual CPU/memory panels work for Backstage, PostgreSQL and NATS. Anonymous
and machine proxy restrictions and Kubernetes write denials passed. Keep scoped
Argo field-ownership compatibility settings until a separately tested removal.
Next: topology/Prometheus alert integration/OpenCost, token-expiry monitoring,
native root-folder presentation and old ClickHouse plugin follow-up. No whole-COP
or Phase 6 completion; #339 remains draft and does not authorize live AI mutations.


Follow-up #344 fixed the three service dashboards' native root-folder presentation
through a real provisioned OneMind folder. [Accepted folder receipt](../operations/grafana-native-folder-2026-10-20.md):
31 dashboard UIDs preserved, only Backstage/DB/NATS moved, real native links and
folder listing verified, healthy same persistent Grafana/portal and read-only
boundaries. Upstream root-dashboard handling elsewhere is not blanket certified.
**Next requested work: Prometheus alert integration**, with native plugin/new
frontend compatibility, service-label scoping, read-only access and honest
firing/pending/no-match/unavailable states. This has not yet been deployed.


Latest Prometheus follow-through: #346 source and #347 image pin merged/deployed;
[scoped acceptance](../operations/prometheus-native-alerts-2026-10-20.md) passed
23:30–23:34 UTC. Portal generation 76, 2/2 healthy on
`sha256:e103bf7c0103b615bef2f4b9dfcb603ef952af7d9eb95d498117d1f48a5dc0e6`.
Native alert cards for Backstage/PostgreSQL/NATS show timestamps, exact scope,
real pending-alert rows/no-match transition and coverage limitations. Fixed
user-only backend reads; machine/anonymous/mutation denials preserved, DB/Grafana
unchanged, NATS 3/3. This supersedes the prior “not yet deployed” next-work entry.
Next: topology/dependencies, OpenCost, broader alert coverage and native state-text
accessibility. No Phase 6 completion/Phase 7 migration/AI mutation authorization.

Latest topology follow-through: #349 merged and scoped acceptance passed on
2026-10-21 UTC. [Receipt](../operations/native-service-topology-2026-10-20.md):
Argo Synced/Healthy at `82c05d3c`, portal generation 77, 2/2 on unchanged
`sha256:e103bf7c0103b615bef2f4b9dfcb603ef952af7d9eb95d498117d1f48a5dc0e6`;
DB/Grafana specifications unchanged, retained storage/backup intact, NATS 3/3.
Native core graph, database reverse dependencies, real NYC3 location and both
replica catalog readers verified; read-only boundaries remain. Six reviewed
services/five existing PVC records are documented metadata, not observed traffic
or complete coverage. No new controller, image, storage or access was added.
**Next: OpenCost**, then separately scoped broader coverage/accessibility and real
multicluster/traffic work. This supersedes “topology next”; Phase 6 stays In
progress, Phase 7 Not started and action PR #339 remains unmerged/draft.
