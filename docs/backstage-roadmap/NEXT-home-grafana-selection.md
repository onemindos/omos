# Personal Grafana dashboards

## Scope

Reuse the installed Grafana API and DashboardViewer. Authenticated dashboard selection saved per account; load charts only when requested; preserve Grafana login and frame policy.

## Delivery

Use official Backstage documentation and installed package types. Keep private Tailscale access, signed Authentik identity and server-side credential boundaries. Preserve existing favorites and saved layouts. Build and review implementation in this PR; record concrete acceptance and remaining activation requirements before marking ready.

## Acceptance

- TypeScript and production build.
- Appropriate authenticated runtime behavior and source freshness.
- Light/dark and responsive layout.
- No fabricated health, contracts, events or costs.

## Implementation

Reuses the installed grafanaApiRef and DashboardViewer. The Viewer-backed list is bounded to 100 dashboards. Up to eight selected exact Grafana dashboard URLs persist through UserSettings, with removal for deleted/inaccessible selections. HTTPS origin and /d/ path are allowlisted. No charts or dashboard-list requests load on the default Home until explicitly requested. Portal frame-src allows only self and the existing Grafana mesh hostname.

Grafana currently sends its default frame denial. Embedded viewing requires a deliberate Grafana setting with a restrictive frame-ancestors policy and a working Grafana session; the UI provides direct authenticated Grafana links. This PR does not disable Grafana authentication or claim embedding acceptance until those prerequisites are addressed. Dashboard selection and direct access work independently.

## Combined-stack correction

The final stack defaults to opening selected dashboards in Grafana until secure embedding is enabled. The opt-in frontend setting onemind.grafanaEmbedding.enabled has an explicit frontend-visible schema and defaults off. This prevents shipping a known-denied iframe as the default chart action. Enable only after Grafana frame-ancestors permits the exact portal origin and a signed-in chart session is accepted.
