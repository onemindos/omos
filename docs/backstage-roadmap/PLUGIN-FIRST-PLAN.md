# OneMind plugin-first implementation plan

Reviewed 2026-10-07 against Backstage 1.55.3, React 18, React Router 6 and the new frontend system. This is the installation shortlist, not a claim that the packages are installed or release-qualified. Source versions below are upstream package manifests; select published, compatible versions through dependency resolution before changing the lockfile.

## Installation order

| Order | Plugin/package | Upstream source version | Admission evidence | Required setup and intended result |
| --- | --- | --- | --- | --- |
| 1 | `@backstage-community/plugin-github-actions` | 1.5.0 | New frontend `/alpha` documented; React 18 peers; Apache-2.0 | GitHub OAuth provider and project-slug annotations. Build history, failures and logs. Keep Authentik login; GitHub authorization is an additional integration. Restrict grants and inspect mutation controls before enabling. |
| 1 | `@roadiehq/backstage-plugin-github-pull-requests` | 3.9.1 | New frontend `/alpha` documented; React 18 peers; Apache-2.0 | Same GitHub authorization. Entity PR lists plus existing home widgets for requested reviews and open PRs. Filter to onemindos. |
| 2 | `@backstage-community/plugin-grafana` | 1.1.0 | New frontend dashboard/alert extensions documented; React 18 peers; Apache-2.0 | Read-only Grafana access through authenticated backend proxy, entity dashboard/alert mappings. Embed only where Grafana authentication and frame policy permit it; never disable authentication to make embedding work. |
| 2 | `@backstage-community/plugin-argocd` | 3.1.3 | New frontend documented; React 18 peers; Apache-2.0 | Existing Kubernetes integration plus documented Argo custom-resource fetching, selectors and `argocd.view.read`. Deployment lifecycle and summary. Compare against Roadie's alternative before resolving package versions. No sync/delete authority. |
| 3 | `@backstage-community/plugin-tech-insights` plus backend/check modules | 1.5.0 frontend | New frontend scorecards documented; React 18 peers; Apache-2.0 | Configure fact collection and checks: owner, runbook and source metadata first. Backup/recovery checks require actual dated evidence. Unknown facts must not become passes. |
| 4 | `@backstage-community/plugin-adr` plus ADR backend/search integration | 0.31.0 frontend | New frontend/search extensions documented; React 18 peers; Apache-2.0 | Git-owned MADR records, reader integration and `backstage.io/adr-location`. Searchable architecture decisions. Mounted catalog locations require explicit resolvable ADR URLs. |
| 5 | `@backstage-community/plugin-tech-radar` | 1.23.0 | New frontend API blueprint guidance; React 18 peers; Apache-2.0 | Curated Git-owned technology data. Approved, experimental and retiring technologies. Verify page registration in selected published release. |
| 6 | `@backstage-community/plugin-cost-insights` | 0.32.0 | New frontend documented; React 18 peers; Apache-2.0 | Supply real billing/allocation API implementation. DigitalOcean costs and namespace allocation are separate sources; do not ship example/mock costs as live data. |

All first-stage packages remain candidates until published manifests, Backstage dependencies, backend permissions and builds are checked. Upstream new-frontend support does not prove compatibility with this exact lockfile. Preserve notices and review transitive licenses before redistribution.

## Coverage of priorities 1–9

| Priority | Reuse first | Remaining configuration or gap |
| --- | --- | --- |
| 1 Today home | Installed Home + GitHub PR home widgets + Notifications + Tech Insights | Compose existing widgets into an action-oriented home; connect real sources. |
| 2 Live events | Installed Notifications/Signals, Argo CD, Grafana | Configure real producers and verify recipients. Plugin display integration is not necessarily an event producer. |
| 3 GitHub work hub | GitHub Actions + Roadie Pull Requests | OAuth, repository mappings and organization filters. |
| 4 Zoho tasks/calendar | Evaluate directory and maintainer repositories before expanding existing adapter | No suitable plugin verified in this review. Existing tasks adapter still needs credentials; calendar coverage is absent. Do not assume no plugin exists. |
| 5 Knowledge search | Installed Search/TechDocs + ADR plugin | Curated portable content; ADR collator. Obsidian export requires scoped source selection. |
| 6 Device/mesh inventory | Installed Catalog/Kubernetes; evaluate Topology for visualization | Tailscale/device discovery remains a data-source concern. A graph viewer does not discover devices. |
| 7 Readiness | Tech Insights | Evidence-based facts and checks, permissions, collection freshness. |
| 8 TAK/agents | Installed API Docs/Templates/Notifications; evaluate Argo Workflows only if its engine is actually adopted | No TAK/NATS operational plugin verified here. Inspect protocol/API contracts and plugin candidates before custom adapters. Agent Forge is a candidate for later review, not approved. |
| 9 Costs | Cost Insights | Real billing and allocation sources; no fabricated estimates. |

## Integration rules

1. Build in focused PRs: GitHub authorization/work views; Grafana/Argo; scorecards; ADR/radar; cost ingestion/UI.
2. Reuse plugin UI, APIs, backend modules and home widgets. Add only documented extension/configuration glue where needed.
3. Retain existing operational views until plugin equivalents demonstrate the same useful data and access boundaries. Retire duplicate custom UI in that replacement PR.
4. Keep secrets server-side except user OAuth flows specifically required by plugins. Preserve signed Authentik identity, private Tailscale access, observational cluster permissions and machine guards.
5. Published package resolution and compilation come before release. Before promotion, explicitly verify real signed-in behavior, data freshness, light/dark design, access denial and rollout within authorized scope.
6. No new package, credential, Argo permission or deployment has been changed by this planning document.

## Maintainer sources

- [Plugin directory](https://backstage.io/plugins/)
- [Community new frontend compatibility table](https://github.com/backstage/community-plugins/blob/main/docs/compatibility/new-frontend-system.md)
- [GitHub Actions](https://github.com/backstage/community-plugins/blob/main/workspaces/github/plugins/github-actions/README.md)
- [Roadie Pull Requests](https://github.com/RoadieHQ/roadie-backstage-plugins/blob/main/plugins/frontend/backstage-plugin-github-pull-requests/README.md)
- [Grafana](https://github.com/backstage/community-plugins/blob/main/workspaces/grafana/plugins/grafana/README.md)
- [Community Argo CD](https://github.com/backstage/community-plugins/blob/main/workspaces/argocd/plugins/argocd/README.md)
- [Roadie Argo CD alternative](https://roadie.io/backstage/plugins/argo-cd/)
- [Tech Insights](https://github.com/backstage/community-plugins/blob/main/workspaces/tech-insights/plugins/tech-insights/README.md)
- [ADR](https://github.com/backstage/community-plugins/blob/main/workspaces/adr/plugins/adr/README.md)
- [Tech Radar](https://github.com/backstage/community-plugins/blob/main/workspaces/tech-radar/plugins/tech-radar/README.md)
- [Cost Insights](https://github.com/backstage/community-plugins/blob/main/workspaces/cost-insights/plugins/cost-insights/README.md)
