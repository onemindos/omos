# P8 — Resilience and home-hosting portability

Status: planned; this draft PR contains the implementation plan, not a deployed feature.

## Outcome

Reduce dashboard update downtime, verify recovery and prepare the cloud deployment for future home infrastructure.

## Implementation

- [ ] Measure current single-replica Recreate behavior and identify a supported safer update strategy.
- [ ] Review database availability, storage retention and backup recovery dependencies.
- [ ] Perform a scoped recovery exercise using a separate instance; preserve existing keys.
- [ ] Document supported configuration for home storage, networking and identity.
- [ ] Record realistic offline behavior: cached documentation can work offline, live remote services cannot.

## Completion criteria

- [ ] The chosen availability design and operational costs are documented.
- [ ] A recovery exercise verifies the selected backup and required keys.
- [ ] Home-hosting prerequisites and migration steps are explicit.
- [ ] Any availability change preserves identity, permissions and database integrity.

## Dependencies and sequencing

Prioritize after P3/P4 and before wider write automation. No promise of HA or offline operation until implemented and verified.

## Release evidence

Before marking ready: summarize actual changes, relevant verification, image digest when needed, rollout health, limitations, and recovery path. Keep this PR draft until implementation is reviewable.
