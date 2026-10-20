# Handoff: continuing the platform roadmap

Updated 2026-10-20 during the authorized portal/catalog stabilization. The authoritative tracker and
phase evidence live in [PLATFORM-ROADMAP.md](../../PLATFORM-ROADMAP.md).

## Where we are

Phases 0–5 are done. **Phase 6 stabilization is the first unfinished phase**;
its historical acceptance is preserved, but Zeus requested UI/catalog repairs
before further migrations. Source/image/inventory/composition PRs #322–#326 are merged and deployed;
final image `sha256:4a52c3b840fcd2b4e2fc0287929e1704721f7420018f2c09deb8543a6aab74da`
is 2/2 healthy, with 715 processed catalog entities and passing targeted live gates;
[stabilization evidence and remaining gates](../operations/portal-stabilization-2026-10-20.md)
are authoritative. Phase 7 remains not started; migration drafts #320/#321 stay unmerged.
Do not begin a later phase or treat this handoff as authorization to deploy it.
Broader token-lifecycle, launcher/empty-card and COP coverage gates remain in the
evidence; absent workflow plane stays explicit and build workflows disabled.
Zeus requested no further Ollama delegation: handle follow-through directly.

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
- Recovery PVC observer-restore-20261019 is intentionally protected/retained. Ask
  before storage cleanup; never delete or restore over production storage.

Latest follow-through: [portal acceptance review](../operations/portal-acceptance-followthrough-2026-10-20.md).
Argo locator works in the production overlay; Grafana dashboard search currently
returns 401. Two suspended backup jobs are intentionally replaced; independent
offsite copy remains unconfigured. Do not infer all-green health from portal readiness.

## Briefs and draft PRs

| Phase | Brief | Goal | Draft PR |
|---|---|---|---|
| 3 | [phase-3-pilot-app.md](phase-3-pilot-app.md) | One small app as an OpenChoreo Component | [#279](https://github.com/onemindos/omos/pull/279) |
| 4 | [phase-4-portal-refactor.md](phase-4-portal-refactor.md) | Portal features into `plugins/` packages | [#280](https://github.com/onemindos/omos/pull/280) |
| 5 | [phase-5-capacity-observability.md](phase-5-capacity-observability.md) | Capacity, then the observability plane | [#281](https://github.com/onemindos/omos/pull/281) |
| 6 | [phase-6-one-portal.md](phase-6-one-portal.md) | OpenChoreo plugins in our portal | [#282](https://github.com/onemindos/omos/pull/282) |
| 7 | [phase-7-migrate-apps.md](phase-7-migrate-apps.md) | Move Layer 3 apps project by project | [#283](https://github.com/onemindos/omos/pull/283) |
| 8 | [phase-8-expand.md](phase-8-expand.md) | Home server and Hermes | [#284](https://github.com/onemindos/omos/pull/284) |

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
