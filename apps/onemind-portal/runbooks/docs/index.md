# OneMind Operations Runbooks

These runbooks describe the deployed OneMind configuration reviewed on October 6, 2026. They provide access, investigation and recovery guidance; they do not automatically execute actions.

## Start with the symptom

| Symptom | Runbook |
| --- | --- |
| Dashboard, catalog or documentation unavailable | [Backstage](backstage.md) |
| Login fails or access is denied | [Identity](identity.md) |
| Agent events or message delivery fail | [NATS](nats.md) |
| TAK clients or geospatial services fail | [TAK](tak.md) |
| A health check fails, or metrics are missing | [Monitoring](monitoring.md) |

## Prerequisites

Use an authorized admin workstation connected to the tailnet. Select and confirm Kubernetes context `do-nyc3-omos-cloud`. Use the `onemindos/omos` repository and review the current branch, PR, and Argo application before changing resources. Keep credentials in the approved secret store; these documents contain no credentials.

## Evidence and scope

P3 checks are cluster-side GET observations, cached for a minute, and do not validate Tailscale routes, login flows or NATS message delivery. Inventory verification tags describe the evidence for an entry, not live health. A missing observation is unknown, not healthy.

Historical Fabric notes are preserved under `docs/inventory/fabric/` in Git. Those archived notes include old namespaces, addresses and reference projects; use them as historical context, not current deployment recipes. Notion is not required to read these runbooks.

## Recovery rule

Capture the symptom, timestamp, recent changes and relevant events first. Prefer a reviewed GitOps fix or known-good release over manually fighting Argo reconciliation. Recovery is complete only after verifying both infrastructure health and the affected user workflow.

## Prepare reviewed work

[Guided workflows](workflows.md) create draft proposals for app onboarding,
deployment changes and recovery preparation, with explicit human promotion gates.
