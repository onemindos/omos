# Phase 4 brief: portal refactor into plugin packages

Status when written (2026-10-18): not started. Can run in parallel with phases 2/3. Tracker: [`PLATFORM-ROADMAP.md`](../../PLATFORM-ROADMAP.md). Read [`HANDOFF.md`](HANDOFF.md) first.

## Goal

Turn the portal's custom features into portable Backstage plugin packages under
`services/onemind-portal/plugins/`, with no behaviour change. This makes phase 6 (adding the
OpenChoreo plugins, or moving to another base portal) a package merge instead of a code merge.

## Starting state

`services/onemind-portal` is Backstage 1.55.3 on the new frontend and backend systems. Features are
already written with the official APIs (`createFrontendPlugin`, `createFrontendModule`,
`createBackendPlugin`, `createBackendModule`) but live inside the app packages:

- Frontend, `packages/app/src/modules/`: `assistant`, `auth.tsx`, `home/` (Launcher, Today,
  ActiveStack, AccessDeployments, LiveOverview, CostPanel, CalendarPanel, FabricPanel, ...),
  `missions`, `nav`, `offline`, `radar`, `services`, `systems`; plus `theme/designModule.tsx`.
  Registered in `packages/app/src/App.tsx`.
- Backend, `packages/backend/src/`: `activity`, `assistant`, `calendar`, `cloud`, `clusterRead`,
  `communityReadOnly`, `context`, `costs`, `health`, `machineReadOnly`, `missions`, `natsSummary`,
  `offline`, `offlinePack`, `operationalEvents`, `operations`, `oidc`, `permissions`,
  `serviceOperations`. Registered in `packages/backend/src/index.ts`.
- `plugins/` is empty (README only). CI image build: `.github/workflows/build-onemind-portal.yml`.
  The deployed image digest is pinned in `platform/backstage/deployment.yaml`.

## Steps

1. Pick the layout. Suggested packages (names `@onemind/...` or `@internal/...`, check existing
   scopes): `plugin-onemind-home` (+ `-common` for shared types), `plugin-onemind-operations`
   and `-backend`, `plugin-onemind-missions` and `-backend`, `plugin-onemind-assistant` and
   `-backend`, `plugin-onemind-offline` and `-backend`, `plugin-onemind-systems`,
   `plugin-onemind-services`, `plugin-onemind-radar`, a `theme` and `nav` module package, and backend
   modules for permissions and read-only machine access.
2. Move one feature at a time, smallest first. Use `yarn new` (Backstage CLI) for scaffolding or
   copy the module layout. Keep exported plugin/module IDs and extension IDs identical so config,
   routes and saved user settings keep working.
3. After each move: `yarn tsc`, `yarn test`, `yarn lint:all`, build the backend and the container
   (`yarn build:all`, `yarn build-image` as CI does), and compare the catalog and operations API
   output before and after.
4. Open one PR per feature. Image rebuild and digest bump is a separate reviewed PR per the portal
   README (application changes need a successful image build and a reviewed digest update).
5. After the last move, check `packages/app/src/App.tsx` and `packages/backend/src/index.ts`
   only wire packages together.

## Gate

The deployed portal behaves as before: sign-in through Authentik, home, Today, Systems, Missions,
Assistant, offline pack, MCP actions (`onemind.health`, `onemind.runbook`, catalog tools) and
`/api/operations`; the catalog counts match; a browser review of the main pages is clean;
the image builds in CI and the new digest is healthy on both replicas.

## Risks

- Cross-module imports inside `modules/home` (shared `colors.ts`, `useDirectory.ts`): extract a
  `-common`/shared package first.
- `directory.json` and `radar/data.json` are bundled data; keep their import paths valid.
- The backend `machineReadOnly` / `communityReadOnly` / `permissions` modules enforce the AI
  read-only rules (token restrictions in `platform/backstage/app-config.cluster.yaml`). Re-verify that
  anonymous requests get 401 and the machine token still cannot mutate.
- Backstage 1.55.3 vs OpenChoreo plugins (tested on 1.51.0): do not upgrade or downgrade here; that is phase 6.

## References

`platform/backstage/README.md`, `platform/backstage/VERIFICATION.md`,
`services/onemind-portal/README.md`, Backstage docs on creating plugins and the new frontend system.
