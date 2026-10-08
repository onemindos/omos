# Contracts, readiness and cost foundations

## Scope

Publish source-verified portal OpenAPI definitions through native API Docs; extend installed Tech Insights with evidence-based checks; connect real DigitalOcean billing only when a constrained source is available. Separate billing totals from namespace allocation; no example cost data.

## Delivery

Use official Backstage documentation and installed package types. Keep private Tailscale access, signed Authentik identity and server-side credential boundaries. Preserve existing favorites and saved layouts. Build and review implementation in this PR; record concrete acceptance and remaining activation requirements before marking ready.

## Acceptance

- TypeScript and production build.
- Appropriate authenticated runtime behavior and source freshness.
- Light/dark and responsive layout.
- No fabricated health, contracts, events or costs.

## Implementation

Native API Docs receives a source-verified OpenAPI 3.0.3 descriptor for the actual status, fabric, cloud and billing routes, mounted through both catalog copies. No TAK or agent payload schema was invented. Tech Insights now checks recorded descriptive metadata and documentation references in addition to ownership. These are explicitly metadata checks, not recovery certification; backups remain timestamped live observations.

User selected DigitalOcean whole-account totals first. The billing adapter queries only GET /v2/customers/my/balance, caches for 15 minutes, validates decimal strings and provider timestamps, and reports configuration-required/unavailable/stale explicitly. It requires an administrator user and denies machine/service access. No raw billing payload, invoices or credentials reach the browser. It does not attribute account costs to OneMind or estimate namespace spending. Native Home widget and Today show the totals with their scope.

Activation requires a dedicated token with only billing:read, sealed as backstage-billing with key DO_BILLING_READ_TOKEN. The existing full doctl credential was not copied into the portal. A local doctl read confirmed the balance endpoint shape; portal activation remains separate. Cost Insights allocation views are deferred until an actual allocation source exists rather than populating its demo client.

Sources: https://docs.digitalocean.com/reference/api/scopes/billing/read/ and https://docs.digitalocean.com/platform/billing/reference/api/ .
