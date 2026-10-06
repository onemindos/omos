# P6 — AI and MCP operational context

Status: planned; this draft PR contains the implementation plan, not a deployed feature.

## Outcome

Expose trustworthy catalog, documentation and health context to authorized AI clients; phase in controlled actions.

## Implementation

- [ ] Use official Backstage MCP documentation and installed package schemas.
- [ ] Start with read-only catalog, documentation and health capabilities.
- [ ] Configure explicit machine authentication and scoped permissions.
- [ ] Return provenance and freshness alongside operational facts.
- [ ] Introduce write actions only through reviewed templates with human approval and audit evidence.

## Completion criteria

- [ ] An authorized client retrieves catalog and operational context.
- [ ] Unauthorized clients and disallowed actions are denied.
- [ ] Secrets are excluded from responses and logs.
- [ ] Stale observations and documented claims remain distinguishable from live facts.

## Dependencies and sequencing

P2 supplies inventory; P3 and P4 supply reliable health and runbooks. Use P5 for any future deployment action.

## Release evidence

Before marking ready: summarize actual changes, relevant verification, image digest when needed, rollout health, limitations, and recovery path. Keep this PR draft until implementation is reviewable.
