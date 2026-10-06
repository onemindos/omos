# P3 — Service health and metrics

Status: planned; this draft PR contains the implementation plan, not a deployed feature.

## Outcome

Connect availability checks and selected Grafana/Prometheus signals to catalog identities and the attention panel.

## Implementation

- [ ] Review existing monitoring and select critical services before adding probes.
- [ ] Choose the availability monitor using current documentation; use existing monitoring where it suffices.
- [ ] Probe reachable service endpoints from the appropriate network; distinguish authentication responses from outages.
- [ ] Map signals to catalog identities and display freshness, missing data, and investigation links.
- [ ] Keep monitoring credentials server-side and access observational.

## Completion criteria

- [ ] Critical-service checks report real observations with timestamps.
- [ ] An unavailable source never reports healthy.
- [ ] The attention panel links health problems to the correct service and investigation surface.
- [ ] Probe load, permissions, and recovery behavior are reviewed.

## Dependencies and sequencing

P1 and P2 are deployed. P4 will document investigation procedures. Finish before enabling operational AI actions.

## Release evidence

Before marking ready: summarize actual changes, relevant verification, image digest when needed, rollout health, limitations, and recovery path. Keep this PR draft until implementation is reviewable.
