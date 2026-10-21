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

(Implementing agent: write progress, decisions, screenshots and gate results here.)
