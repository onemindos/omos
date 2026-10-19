# Handoff: continuing the platform roadmap

Written 2026-10-18 after phase 2. Any AI agent (Claude Code, Codex, another) or person can
continue from here. The plan, tracker and decisions live in
[`PLATFORM-ROADMAP.md`](../../PLATFORM-ROADMAP.md); this file says where things stand and how to work.

## Where we are

Phases 0, 1 and 2 are done (PRs #256 to #277). Next is **phase 3** (pilot one app in OpenChoreo).
Phases 3 to 8 each have a brief in this folder and a draft PR carrying it. Open items:

- Zeus has not yet signed in to the OpenChoreo console as `akadmin`
  (`https://onemind-console.tail717ba0.ts.net`). If it fails, reopen phase 2.
- `platform-matrix` is OutOfSync and `platform-openclaw` was still progressing at the last check.
- The ecosystem node sits near 99% CPU requested (phase 5 work).

## Briefs and draft PRs

| Phase | Brief | Goal |
|---|---|---|
| 3 | [phase-3-pilot-app.md](phase-3-pilot-app.md) | One small app as an OpenChoreo Component |
| 4 | [phase-4-portal-refactor.md](phase-4-portal-refactor.md) | Portal features into `plugins/` packages |
| 5 | [phase-5-capacity-observability.md](phase-5-capacity-observability.md) | Capacity, then the observability plane |
| 6 | [phase-6-one-portal.md](phase-6-one-portal.md) | OpenChoreo plugins in our portal |
| 7 | [phase-7-migrate-apps.md](phase-7-migrate-apps.md) | Move Layer 3 apps project by project |
| 8 | [phase-8-expand.md](phase-8-expand.md) | Home server and Hermes |

Each draft PR adds only its brief. Work the phase on that PR's branch (or a branch from it),
keep the PR a draft until the gate passes, then mark it ready and merge. Every PR also updates the
tracker in `PLATFORM-ROADMAP.md`.

## Access (set up in phase 1 and 2)

- Read-only cluster and portal access for AI agents: `scripts/ai-operator/setup.sh`, documented in
  [`docs/operations/ai-operator-access.md`](../operations/ai-operator-access.md). MCP servers:
  `onemind-kubernetes`, `onemind-backstage`, `onemind-openchoreo` (63 read tools; writes refused).
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
