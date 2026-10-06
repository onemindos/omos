# P8 — Resilience and home-hosting portability

Status: implemented; isolated recovery passed. Rollout verification follows merge.

## Outcome

Reduce dashboard update downtime, verify recovery and prepare the cloud deployment for future home infrastructure.

## Implementation

- [x] Measure current single-replica Recreate behavior and identify a supported safer update strategy.
- [x] Review database availability, storage retention and backup recovery dependencies.
- [x] Perform a scoped recovery exercise using a separate instance; preserve existing keys.
- [x] Document supported configuration for home storage, networking and identity.
- [x] Record realistic offline behavior: cached documentation can work offline, live remote services cannot.

## Completion criteria

- [x] The chosen availability design and operational costs are documented.
- [x] A recovery exercise verifies the selected backup and required keys.
- [x] Home-hosting prerequisites and migration steps are explicit.
- [x] Any availability change preserves identity, permissions and database integrity.

## Dependencies and sequencing

Prioritize after P3/P4 and before wider write automation. No promise of HA or offline operation until implemented and verified.

## Release evidence

Before marking ready: summarize actual changes, relevant verification, image digest when needed, rollout health, limitations, and recovery path. Keep this PR draft until implementation is reviewable.

Recovery evidence: 14 application databases restored to an ephemeral separate PostgreSQL instance; 230 catalog entities from the selected backup. Production storage was untouched. Replica rollout does not eliminate PostgreSQL or identity single points of failure. Home migration is documented, not performed.
