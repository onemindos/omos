# Live Home overview

## Scope

Reuse existing status, Notifications/Signals and GitHub integrations. Surface compact availability, deployment, backup and unread-event summaries with timestamps, partial/stale states and expandable detail. Avoid redundant polling and heavy default charts.

## Delivery

Use official Backstage documentation and installed package types. Keep private Tailscale access, signed Authentik identity and server-side credential boundaries. Preserve existing favorites and saved layouts. Build and review implementation in this PR; record concrete acceptance and remaining activation requirements before marking ready.

## Acceptance

- TypeScript and production build.
- Appropriate authenticated runtime behavior and source freshness.
- Light/dark and responsive layout.
- No fabricated health, contracts, events or costs.

## Implementation

Default Home now shows timestamped, minute-refresh workload, Argo, backup schedule and critical-service summaries. Failed source lists are unknown rather than zero. Previous observations are explicitly stale after refresh failure. Inbox loads only when expanded; full Signals-backed Notifications remains available. Native Most Visited joins clocks and recent visits in the fluid layout. No chart iframes load by default.
