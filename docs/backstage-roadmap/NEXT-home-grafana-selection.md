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
