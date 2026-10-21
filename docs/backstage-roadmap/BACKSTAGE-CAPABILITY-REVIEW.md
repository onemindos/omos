# OneMind Backstage capability review

Reviewed 2026-10-07 against the official site, next API reference, core feature guides, theme/UI guides, and plugin directory; compared with installed package manifests, backend registration, permissions policy, and the combined #64 source. This is a scoped capability audit, not a claim that every third-party plugin or endpoint has been exhaustively assessed.

## Version rule

The app is Backstage 1.55.3. `/api/next/` documents upcoming APIs: consult it for direction, then verify exports/types in installed packages before use. Do not install next packages merely because they appear in that reference. New frontend extensions and mixed Material UI/Backstage UI are supported; visual customization must cover both systems.

## Coverage and remaining work

| Capability | Source coverage | Remaining operational work |
| --- | --- | --- |
| Software Catalog, users/groups, ownership | Installed; Git-backed inventory and curated registry | Improve API/System/Domain records and explicit ownership as services expand |
| Catalog graph and dependency exploration | Graph installed; Systems workspace #58 | Record verified service-to-service dependencies; no inference from namespace proximity |
| API documentation | API docs package installed | Publish real OpenAPI/AsyncAPI definitions for NATS, agents, TAK and gateways; an installed viewer does not create contracts |
| Kubernetes | Installed with observational access | Validate entity selectors; future home/cloud multi-cluster configuration and credentials |
| Search | PostgreSQL engine, catalog and TechDocs collators registered | Add only authorized sources; no Notion dependency as it is leaving; verify indexing and permissions with real sessions |
| TechDocs | Packaged runbooks and navigation | Continued operational ownership and recovery rehearsal evidence |
| Software Templates | Existing proposal flow; new reviewed workflows #62 | Complete first signed-in draft PR run with narrowly scoped token; proposals do not execute deployment/recovery |
| Auth and permissions | OIDC, backend policy; admin writes, authenticated reads | For community expansion, add resource/team isolation and explicit policies; current all-authenticated-read policy is not tenant isolation |
| Notifications and signals | Installed; bounded inbox/event bridge #57 | Provision restricted producer identity, verify group recipients, connect real Argo/backup/agent producers |
| Service operations | #61 observational adapter and dashboard links | End-to-end runtime acceptance; namespace summaries are not exact per-service metrics |
| Inventory refresh | #55 opt-in reviewed discovery | Dedicated observational runner, discovery token and visibility confirmation |
| Cloud work | #59 paginated Zoho adapter | Read-only provider credentials absent; no live summary yet |
| AI and MCP | Existing restricted MCP; optional evidence assistant #60 | Select/provision model endpoint, confirm performance and evidence quality; no autonomous mutations |
| Offline knowledge | Explicit private export #63 | Validate runtime export; full home-cluster continuity needs compute, identity, DNS, registries and restore procedures |
| Brand and usability | This stacked design PR | Signed-in light/dark, mobile, keyboard and plugin-page acceptance before promotion |
| Reliability scorecards | Not implemented | Define grounded checks for owner/runbook/backup freshness/image pinning and recovery evidence; evaluate Tech Insights after compatibility check |
| CI and PR overview | GitHub links only | Read-only GitHub Actions/PR integration for failed builds and reviews, with restricted credentials |
| FinOps | Not implemented | Start with scoped DO billing/usage and namespace allocation; Cost Insights needs a real cost-data adapter |
| SLOs and incident history | Current health observations and Grafana links | Add service SLOs/history and alert correlation; a successful probe is not historical uptime |
| Audit trail | No new unified audit store in this batch | Retention, redaction and export of authenticated administrative actions before shared/community operation |
| Adoption analytics | Not implemented | Optional self-hosted analytics with minimal data collection; no third-party trackers needed |
| Architecture decisions / technology standards | Docs can host ADRs | Establish ADRs and a small Tech Radar before selecting a dedicated viewer |

## Plugin shortlist and admission criteria

1. **GitHub Actions / PR views:** immediate value for build failures and pending reviews. Require new-frontend support, read-only integration, bounded fetches and backend token handling.
2. **Tech Insights:** readiness scorecards grounded in existing catalog and recovery evidence. Check exact package compatibility, licenses and check definitions before choosing an implementation.
3. **Cost Insights:** useful only after reliable cost ingestion/allocation exists; do not display fabricated estimates.
4. **Grafana / Argo integrations:** existing native #61 read-only adapter covers basic needs; evaluate community plugins for deeper service-specific alert and deployment views without duplicating authentication.
5. **ADR / Tech Radar:** start with owned TechDocs content; add plugins if navigation needs justify maintenance.

The directory explicitly says community plugins have not all been vetted by the Backstage core team. Listing is not compatibility, security, license or maintenance approval. For each candidate check maintained releases, 1.55.3/new frontend support, permissions, backend credential boundaries, network/data handling, transitive dependencies and license. Avoid overlapping chat assistants, writable secrets UIs and generic deploy buttons without a concrete need and deliberate policy.

## Recommended sequence

Finish runtime acceptance and activation of #64; review the branded experience alongside it. Then prioritize real event producers, GitHub build/PR views, readiness scorecards, API contracts and SLO history. Add cost allocation once sources exist. Preserve sovereign hosting and portable Git-backed records. More packages alone do not create more operational capability.

## Official sources

- https://backstage.io/
- https://backstage.io/api/next/
- https://backstage.io/plugins/
- https://backstage.io/docs/golden-path/create-app/custom-theme/
- https://backstage.io/docs/conf/user-interface/
- https://backstage.io/docs/features/software-catalog/descriptor-format/
- https://backstage.io/docs/features/kubernetes/
- https://backstage.io/docs/features/search/
- https://backstage.io/docs/permissions/overview/
