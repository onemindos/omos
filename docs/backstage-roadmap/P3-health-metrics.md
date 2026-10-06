# P3 — Service health and metrics

Status: implemented; release verification in progress.

## Outcome

Connect availability checks and selected Grafana/Prometheus signals to catalog identities and the attention panel.

## Implementation

- [x] Review existing monitoring and select critical services before adding probes.
- [x] Choose the availability monitor using current documentation; use existing monitoring where it suffices.
- [x] Probe reachable service endpoints from the appropriate network; distinguish authentication responses from outages.
- [x] Map signals to catalog identities and display freshness, missing data, and investigation links.
- [x] Keep monitoring credentials server-side and access observational.

## Completion criteria

- [x] Critical-service checks report real observations with timestamps.
- [x] An unavailable source never reports healthy.
- [x] The attention panel links health problems to the correct service and investigation surface.
- [x] Probe load, permissions, and recovery behavior are reviewed.

## Dependencies and sequencing

P1 and P2 are deployed. P4 will document investigation procedures. Finish before enabling operational AI actions.

## Release evidence

Before marking ready: summarize actual changes, relevant verification, image digest when needed, rollout health, limitations, and recovery path. Keep this PR draft until implementation is reviewable.

Selected design: on-demand bounded internal GET probes cached for one minute, plus existing Prometheus queries. Four service targets and five namespaces are configured. No continuous uptime history or external user-path monitoring is claimed.
