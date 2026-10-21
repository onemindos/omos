# OneMind flagship portal program

Requested by Zeus on 2026-10-21: make Backstage the single front door for every
OneMind service, Kubernetes and non-Kubernetes. It should be organized, show live
problems in real time, give DigitalOcean at a glance, and support adding and
removing apps. Claude (Opus) wrote these work orders; Sonnet agents implement them,
one PR each.

This program **simplifies** the earlier plan. It replaces the U00–U06 experience
packages from #356 with eleven concrete, buildable work orders. It is Phase 6
("One portal") follow-through. It does not start Phase 7, migrate apps between
controllers, or widen AI permissions.

## What a signed-in review found (2026-10-21)

| Finding | Evidence | Fixed by |
|---|---|---|
| 7,750 unread notifications; every Argo sync success is one | `/notifications`, tab title `(7750)` | F3 |
| Sidebar has no structure: after Home/Today/Catalog/Create, every page is sorted alphabetically | `plugin-onemind-nav/src/Sidebar.tsx` (`nav.rest({ sortBy: 'title' })`) | F1 |
| Home launcher leaves large empty areas below the fold; 130 directory chips with 25 categories | `/` screenshots, `HomeLayout.tsx` | F1 |
| Navy/teal colour scheme; does not match the red/black Legacy brand | `plugin-onemind-theme/src/designModule.tsx` (`#121b29`, `#73d9e4`) | F10 |
| Catalog: about 70 Systems in three naming schemes (`cluster-*`, `fabric-*`, capability names), one `default` Domain, about 280 hash-suffixed Resources | catalog query, 348 non-Component entities | F2 |
| The portal's own entity is `k8s-backstage-backstage-21b99354`; `component:default/backstage` does not exist | `/catalog/default/component/backstage` → not found | F2, F5 |
| Only 6 of 69 Argo CD applications are linked to a catalog entity | catalog query on `argocd/app-name` | F2, F5 |
| App pages pop up an unrequested GitHub OAuth dialog | GitHub Actions/PR cards on every Component | F5 |
| Cost Insights page loads empty | `/cost-insights` | F9 |
| DigitalOcean cost card says it needs a `billing:read` token; no droplet/volume/LB/snapshot view | `/today`, no DO plugin exists in the directory | F6 |
| Three backup jobs suspended since Oct 7–8 (backstage, clickhouse, offsite copy) | `/today` "Needs attention" | F4 shows it; the fix is an ops task outside this program |
| 98 DO volumes (all bound to PVCs) and 58 volume snapshots, no retention view | `doctl` read-only | F6, F8 |

## Work orders

| ID | PR | Work order | Depends on | Visible result |
|---|---|---|---|---|
| F1 | [#357](https://github.com/onemindos/omos/pull/357) | Navigation and home reorganization | none | Grouped sidebar, compact home with live status strip |
| F2 | [#358](https://github.com/onemindos/omos/pull/358) | Catalog cleanup and SO/LE/HP/GE structure | none | Four pillar Domains, clean System list, human names, noise hidden |
| F3 | [#359](https://github.com/onemindos/omos/pull/359) | Notification triage and real-time alert routing | none | Inbox shows only actionable items; firing alerts arrive live |
| F4 | [#360](https://github.com/onemindos/omos/pull/360) | Live Operations: "what is wrong right now" | F3 (soft) | One page: alerts, Argo drift, failing pods, nodes, storage, backups, uptime |
| F5 | [#361](https://github.com/onemindos/omos/pull/361) | One standard app page for every service | F2 | Same tabs on every app; no GitHub pop-ups; non-K8s variant |
| F6 | [#362](https://github.com/onemindos/omos/pull/362) | DigitalOcean at a glance (custom plugin) | read-only DO token from Zeus | Droplets, DOKS, volumes, snapshots, LBs, firewalls, VPCs, Spaces, uptime, bill |
| F7 | [#363](https://github.com/onemindos/omos/pull/363) | Non-Kubernetes apps: model, health and access | F2 | SaaS/home/device apps managed with the same page and live health |
| F8 | [#364](https://github.com/onemindos/omos/pull/364) | Add and remove apps: lifecycle templates and stack hygiene | F2 | "Add app" and "Retire app" produce reviewed PRs; hygiene report |
| F9 | [#365](https://github.com/onemindos/omos/pull/365) | Plugin pack: verified additions | none | Kyverno policy reports, Harbor images, OpenCost, shortcuts/bookmarks, admin devtools |
| F10 | [#366](https://github.com/onemindos/omos/pull/366) | Theme and visual design system | none | Brand-correct red/graphite theme in dark and light, AA contrast |
| F11 | [#367](https://github.com/onemindos/omos/pull/367) | Live inventory: catalog follows the cluster automatically | F2 | New and removed workloads appear within minutes, attached to their app |

Suggested build order: **F10, F1, F2, F3** in parallel (they touch different
files), then **F4, F5, F11**, then **F6, F7, F8, F9**. Merge one PR at a time; each
rebases on the then-current `main` and reruns the gate.

## Stack additions recommended (evaluate inside the named work order)

- **Gatus** (F7): uptime and status checks as config-as-code for Kubernetes,
  SaaS, home and device endpoints. It exports Prometheus metrics, so alerts reach
  F3/F4 through the existing Alertmanager. Runs in-cluster and later at home.
- **Kyverno Policy Reporter** (F9): Kyverno is already installed; this exposes
  its policy reports so the portal can show them.
- **Trivy Operator** (F9, optional): vulnerability reports per workload, surfaced
  through Policy Reporter. Needs a capacity check first.
- Already installed and only needing portal wiring: Harbor, OpenCost, Loki,
  Tempo, Alertmanager, Velero, Falco.

## Existing drafts

- #356 (plan) can merge as history; this program supersedes its U00–U06 order.
- #351 → folded into F8. #353 and #355 → folded into F2 (reuse their reports as
  F2's before/after evidence). #352 → reused by F9's cost work. #354 → used by
  F4's backup/recovery evidence.
- #339 (AI actions), #332 (LibreChat) and #284 (home cluster) stay parked as
  drafts. They are not part of the front-door work.

Do not close any of them without Zeus's explicit OK.

## Rules every implementing agent follows

1. Read `AGENTS.md`, `CLAUDE.md`, `docs/MISSION.md` and your work order first.
   Work only on your work order's branch, and touch only the files it lists. If
   you need to go wider, stop and write the reason into the PR.
2. Do **not** edit `PLATFORM-ROADMAP.md`. Write progress and evidence into the
   "Status" section at the bottom of your own work-order file. The coordinator
   updates the roadmap tracker when each PR merges. This prevents conflicts.
3. Backstage is **1.55.3 on the new frontend system** (`createApp` from
   `@backstage/frontend-defaults`). Prefer plugins with `/alpha` new-frontend
   exports. Wrap legacy plugins with `@backstage/core-compat-api` only when needed.
   Pin exact versions and keep `yarn.lock` immutable-installable.
4. Never commit secrets. A new credential is a SealedSecret that **Zeus creates**.
   Write the exact secret name, keys and scope into the PR and stop for Zeus.
5. No `kubectl apply/patch/delete`, no Helm installs and no Argo syncs. Git is the
   only way to change the cluster, and Argo applies it after a human merge.
6. Do not widen permissions. Existing read-only AI access stays read-only. Any
   new backend route goes through the existing permission policy, with a
   negative (denied) test.
7. Keep stable entity refs. Renaming is done with `metadata.title`, never by
   changing `metadata.name` of an existing entity, unless the work order says so
   and lists the redirects.
8. Validate before opening the PR for review:

   ```bash
   export PATH="$HOME/.nvm/versions/node/v24.18.0/bin:$PATH"   # Node 22 or 24
   (cd services/onemind-portal && node .yarn/releases/yarn-4.13.0.cjs install --immutable)
   python3 scripts/gitops/pipeline.py check --profile portal --base origin/main
   ```

   All **required** gates must pass. The optional `whole-tree-diagnostic` fails on
   old debt at `platform/onetak/media-mtx.yaml:86`; leave that alone. GitHub CI
   is billing-blocked; never bypass checks or change billing.
9. Releasing is separate and needs Zeus's OK. Code changes need an image:
   `pipeline.py publish`, then a digest-pin PR, then Argo, then signed-in
   acceptance. Catalog-only changes (mounted files) roll with Argo after merge.
10. Every UI change includes before/after screenshots at desktop and phone width,
    in dark and light themes, attached to the PR.
