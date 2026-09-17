# Context Profile: Incident Response

Load this context when something is broken, down, or needs immediate attention.

---

## Pre-Load Files

1. `codex/00-24 SO/08 Infrastructure/` — Server, service, and API documentation
2. Relevant SOPs from `codex/00-24 SO/06 Protocols/`
3. Recent daily logs for context on what changed recently

## Priority Order

1. **Assess**: What's broken? What's the impact? Who's affected?
2. **Contain**: Stop the bleeding — prevent further damage
3. **Fix**: Apply the fix (follow SOPs if they exist)
4. **Verify**: Confirm the fix worked
5. **Document**: Create a log entry and update relevant codex notes

## Key Infrastructure References

- VPS documentation (look in `08 Infrastructure/`)
- Service status and ports
- Docker/systemd service management
- NATS infrastructure

## After Resolution

1. Create a decision record if a non-trivial fix was chosen
2. Update or create an SOP if no procedure existed
3. Update infrastructure docs if anything changed
4. Add a daily log entry documenting the incident

## Outputs

- Incident resolved
- Daily log entry with details
- Updated SOPs if needed
- Decision record if non-trivial choices were made
