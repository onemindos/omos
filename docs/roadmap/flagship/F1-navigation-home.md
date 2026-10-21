# F1 — Navigation and home reorganization

Work order for one implementing agent. Branch: `flagship/f1-navigation-home`.
Program index: [README](README.md). Read the rules there first.

## Outcome Zeus sees

- **Sidebar.** A short, grouped sidebar that matches how Zeus works, not an
  alphabetical list of plugin pages.
- **Home.** A compact home that answers three questions at a glance: what is
  broken right now, where my apps are, and what changed recently. No large empty
  areas.
- **Command palette.** ⌘K / Ctrl-K opens one search across apps, runbooks,
  docs, alerts and (once F6 lands) DigitalOcean resources.

## Evidence

- `plugins/plugin-onemind-nav/src/Sidebar.tsx`: after Home/Today/Catalog/Create,
  `nav.rest({ sortBy: 'title' })` dumps every page alphabetically: Announcements,
  APIs, Audit Logs, Catalog Graph, Cost Insights, Delivery Insights, Docs,
  Register Existing Component, Tech Insights, Tech Radar, Visualizer.
- `plugins/plugin-onemind-home/src/HomeLayout.tsx` and the `Directory` launcher:
  130 entries and 25 category chips, plus hosting chips. Below "My Active Stack",
  large blank regions render on desktop (signed-in screenshots, 2026-10-21).
- `/today` already has a good "Needs attention" block. Home does not surface it.

## Build

1. **Sidebar groups** (exact order). Keep each existing route and deep link
   working.
   - **Now**: Home, Live Ops (F4 route; until F4 lands, link `/today`), Notifications.
   - **Pillars**: SO, LE, HP, GE (links to the F2 Domain pages; hide the group
     until F2's Domains exist, using a feature flag).
   - **Services**: Catalog, APIs, Catalog Graph, DigitalOcean (F6, flagged).
   - **Insights**: Cost, Delivery, Tech Insights, Tech Radar.
   - **Knowledge**: Docs, Runbooks, Announcements, ADRs.
   - **Change**: Create (templates), Register existing, Retire app (F8, flagged).
   - **Admin** (collapsed, admins only): Audit Logs, Visualizer, DevTools (F9),
     Settings.
   Use `SidebarGroup`/`SidebarSubmenu` from `@backstage/core-components`; do not
   fork the sidebar component. Mobile: the same groups in the bottom bar menu.
2. **Home layout** (`plugin-onemind-home`), top to bottom:
   - **Status strip**: counts of firing alerts, out-of-sync/degraded Argo apps,
     not-ready workloads and suspended backups. Reuse the data already behind
     `/today`'s "Needs attention" (`ActiveStack`, `LiveOverview`). Each count is
     a link, and each shows its observed time.
   - **Favorites row**: starred entities (native `starredEntities` API) and
     pinned shortcuts.
   - **Pillar tiles**: SO / LE / HP / GE with app counts (flagged until F2).
   - **Recently changed**: last 10 deployments and catalog changes, deduplicated.
   - **Launcher**: collapse the 25 category chips into a single filter dropdown
     plus a search box. Keep "Add shortcut" / "Add shared service".
   - Fix the empty-space layout bug: no fixed heights; let content flow.
   - Keep "Customize widgets" (native `CustomHomepageGrid`).
3. **Command palette**: a global ⌘K modal using the existing search backend
   (`@backstage/plugin-search`, Postgres engine). Add result types for runbooks
   and app links. Keyboard navigable. No new search engine.

## Out of scope

Theme colours (F10), catalog data changes (F2), new data sources (F4/F6).

## Files

`plugins/plugin-onemind-nav/**`, `plugins/plugin-onemind-home/**`,
`packages/app/src/App.tsx` (feature registration only), `app-config.yaml`
(feature flags only), and tests next to the changed components.

## Acceptance

- [ ] Every page that existed before is reachable within two clicks; deep links
      and favorites are unchanged.
- [ ] Home's first screen (1440×900) shows the status strip, favorites and the
      launcher search with no blank region taller than 80px.
- [ ] Status strip numbers match `/today` for the same minute.
- [ ] ⌘K finds "backstage", "nats", a runbook title and a Docs page.
- [ ] Keyboard-only and screen-reader labels pass for the sidebar and palette;
      works at 390px width.
- [ ] Component tests for the sidebar grouping and home sections; required
      gates pass.

## Status

**Done (2026-10-21)**

- Sidebar: groups Now, Pillars, Services, Insights, Knowledge, Change, Admin built
  with `SidebarGroup`/`SidebarItem`/`SidebarSubmenu` (no fork). `buildNavModel`
  (pure, `navModel.ts`) resolves the spec against the app's pages; anything not
  placed goes to a collapsed **More** submenu, so every earlier page stays
  reachable. Admin and More are collapsed buttons opening a keyboard-operable
  flyout. Mobile bottom bar shows one entry per group (overlay menu), shrunk to
  fit 390px. Live Ops links `/today`; Notifications keep the live unread badge.
- Flags (`onemind.features.pillars|digitalocean|retireApp`, default false) in
  `app-config.yaml`, schema in `plugin-onemind-nav/config.d.ts`.
- Home: status strip (reuses `/api/onemind/status`; `status.ts` holds the shared
  backup-review rule used by `/today` and the strip), favorites row, pillar
  tiles (flagged), recently changed (notification feed, deployments/catalog
  topics, deduped, max 10), launcher with search + category and hosting
  selects, Active Stack, native Customize widgets. No fixed heights.
- Palette: Cmd/Ctrl-K opens `SearchModal` (`@backstage/plugin-search`) with
  arrow-key navigation. Runbooks and apps are found through the existing
  TechDocs and catalog collators.
- Colours in touched code come from `theme.palette.*` (no new hex).

**Dependencies added:** `@backstage/test-utils` 1.7.22 and
`@testing-library/react` ^14 (both devDependencies, nav plugin, tests only).
No runtime dependency added.

**Plugins considered (rule 11):** `@backstage/plugin-search` SearchModal and
`@backstage/plugin-home` CustomHomepageGrid (used, Backstage core, current);
`@backstage/plugin-catalog` starred entities API (used). No community plugin
offers a grouped sidebar or status strip; custom code is limited to
`navModel.ts` (grouping spec, no plugin exists), `StatusStrip` (reads the
existing status route), `RecentlyChanged` and `PillarTiles` (thin readers of
notifications and catalog APIs).

**Deferred / decisions**
- "Alerts firing" counts failing critical service checks from the existing
  status route; Alertmanager alerts arrive with F3/F4.
- "Catalog changes" in Recently changed need a source (F11); only deployment
  and catalog-topic notifications appear today.
- Pillar and Retire routes (`/catalog/default/domain/<so|le|hp|ge>`,
  `/create/templates/default/retire-app`) are constants in `navModel.ts`; F2/F8
  must confirm them. DevTools and Audit Logs show only if those pages exist.
- Search result types: no custom result items; runbooks and apps already come
  from the existing collators. Alerts and DigitalOcean are not indexed yet.
- Admin visibility is a UI convenience (identity ownership refs); backend
  permissions are unchanged.
- Not verified against the live portal: ids of some plugin pages (matched by id
  or title). Local run used mocked APIs (Brave via Playwright), not real data.

**Gate:** `pipeline.py check --profile portal` run.json = check_passed (only optional whole-tree-diagnostic fails).

**Tests:** nav 15, home 16 new/updated (31 total passing).

**Screenshots (local, not committed):** `<scratch>/agents/f1-shots/after-*.png`
(desktop 1440x900 and phone 390x844, dark and light, plus palette). No "before"
set captured.
