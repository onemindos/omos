# P7 — Selected cloud application summaries

Status: adapter implemented; authorized Zoho credentials and live account verification remain pending.

## Outcome

Show a small set of useful cloud application summaries without duplicating entire apps.

## Implementation

- [ ] Choose summaries that support daily decisions, starting with a reviewed Zoho use case.
- [ ] Verify each provider API, account scope and licensing before implementation.
- [ ] Use minimal read scopes and server-side credentials.
- [ ] Cache responses with visible observation time and clear expired/unavailable states.
- [ ] Preserve deep links to the provider for full workflows.

## Completion criteria

- [ ] The selected summary uses authorized account data and a verified API.
- [ ] Refresh respects rate limits and handles provider outages.
- [ ] Sensitive data is visible only to authorized users.
- [ ] Users can open the corresponding application from the summary.

## Dependencies and sequencing

Independent of P5. Integrate with P2 identities; coordinate permissions with P6. Exact summaries require a concrete selection before implementation.

## Release evidence

Before marking ready: summarize actual changes, relevant verification, image digest when needed, rollout health, limitations, and recovery path. Keep this PR draft until implementation is reviewable.

Selected Zoho Projects V3 tasks as the first decision summary. The card explicitly distinguishes configuration required, unavailable, stale and current. Scope ZohoProjects.tasks.READ; US endpoints; first 100 tasks only, labeled partial. Provider entitlement and live account results require the missing credentials.
