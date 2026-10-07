---
status: accepted
date: 2026-10-07
decision-makers: OneMind administrator
---

# Use established Backstage plugins before custom integrations

## Context and Problem Statement

OneMind combines private Kubernetes infrastructure, cloud applications, documentation and distributed agents. Rebuilding capabilities already offered by maintained plugins increases the code that OneMind must own.

## Decision Outcome

Use existing Backstage plugins wherever they meet the operational need. Check the installed Backstage version, actual plugin exports, permissions, license and hosting requirements. Configure plugin-provided pages, APIs and home widgets before implementing a replacement. Custom code is limited to verified gaps and integration glue.

Keep Authentik as the portal identity provider, Tailscale as private transport and GitOps as the deployment path. GitHub integration uses an organization-owned GitHub App with read-only permissions and selected repository access. Provider credentials remain server-side or in the plugin's documented user authorization flow. Real sources must be connected before claiming live capability.

## Consequences

OneMind benefits from upstream improvements and has fewer custom views to maintain. Plugin upgrades still require compatibility and migration review. Installation alone does not create API contracts, billing inputs, workflow engines or recovery evidence. Existing custom operational views remain until their plugin replacements demonstrate equivalent useful data.
