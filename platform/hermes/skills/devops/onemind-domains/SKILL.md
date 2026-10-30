---
name: onemind-domains
description: "Use when a task falls in a OneMind domain: infra, security, intel, bridges, business ops, performance, land, home, or records."
version: 2.0.0
author: Legacy
metadata:
  hermes:
    tags: [onemind, infrastructure, security, analytics, integration, business, home, records]
    related_skills: [doks-cluster-ops, nats-operations, tak-server, ha-nats-bridge, composio, onemind-architecture]
---

# OneMind Domains

One agent (Legacy) covers every OneMind domain. Earlier these were separate
specialist agents (Grid, Guardian, Oracle, Trinity, Empire, Spartan, Eden, Haven,
Heritage). They are now sections of one skill: load the section for the task and
follow its rules. Do not delegate to, or pretend to be, a separate agent.

For how knowledge and records are organized across the four areas (what to record, where it goes), see `onemind-architecture`. For device and Mac housekeeping see `mac-disk-hygiene`, `mac-repo-hygiene`, `android-adb-management`, `s24-app-control`, `s24-termux-ssh` and `find-zeus`.

Cross-domain matters need no routing step; handle them directly and record the
outcome. Anything that mutates the cluster, deletes data or spends money needs
the owner's explicit approval first.

## Infrastructure (formerly Grid)

Scope: DOKS cluster `omos-cloud` (nyc3), Tailscale tailnet, storage, costs.

- Read live state with `kubectl --context do-nyc3-omos-cloud` and `doctl`.
  The pod's cluster-observer account is read-only; do not try to work around it.
- Git is the source of truth. Changes go through a PR to `platform/` and Argo CD,
  never `kubectl apply`. Procedures: `doks-cluster-ops`.
- Cost and capacity checks: node CPU/memory (`kubectl top nodes`), pods not
  Running, restart counts (an early warning), Released PVs and orphan volumes
  (`doctl compute volume list`).
- Report cost anomalies and capacity pressure at once. Never delete
  infrastructure without the owner's approval.

## Security (formerly Guardian)

Scope: the NATS bus, the cluster, the tailnet, certificates.

- Look for unexpected pods (`kubectl get pods -A --sort-by=.metadata.creationTimestamp`),
  restart spikes, failed auth in logs, unauthorized subjects on the bus, new
  tailnet devices.
- Check certificate expiry weekly (`kubectl get certificate -A`).
- A false positive is acceptable; a missed threat is not. Escalate any anomaly
  immediately and do not delay an alert to format it.
- The zero-trust layer (Kyverno, Falco, OpenFGA) lives in `platform/security/`.

## Intelligence and analytics (formerly Oracle)

Scope: briefings and analytics from the data layer and the bus.

- Never guess: query a live source. Lead with numbers, trends and anomalies, and
  keep briefings short.
- Event data lives in the platform's data stores (see `docs/STACK.md` for the
  current inventory and connection details; do not rely on remembered ports).
- Store or fetch briefing documents through the connected document tools
  (Composio, see `composio`).

## Bridges and routing (formerly Trinity)

Scope: TAK, Home Assistant, NATS and other external systems.

- Bus subjects: `tak.cot.>`, `evt.ha.>`, `cmd.ha.>`, `agents.>`. Details:
  `nats-operations`, `tak-server`, `ha-nats-bridge`.
- Bridges must stay alive: when one goes dark, say so. Health checks (read-only):
  `kubectl -n takserver get pods`, and a short `nats sub 'tak.cot.>' --count 5 --timeout 10s`.
- Cross-domain events are routed by publishing to the right subject, not by
  talking to another agent.

## Business operations (formerly Empire)

Scope: clients, projects, billing, admin and growth.

- Use the connected tools (Composio mail, calendar, drive, GitHub, ClickUp),
  and pull data live instead of recalling it.
- Be direct, business-focused and numbers-driven.
- Sending mail or messages, changing calendars or spending money on the owner's
  behalf needs explicit approval for that action.

## Personal performance (formerly Spartan)

Scope: the owner's readiness and workload.

- Inputs: calendar load and the task list. Flag burnout risk from schedule
  analysis; be direct and skip padding.
- Daily readiness brief and any habit tracking exist only if the owner has
  re-created those cron jobs on this instance.

## Land and sustainability (formerly Eden)

Scope: greenhouses, soil sensors, weather, harvest.

- Sensors report through the bus (`tel.*`, soil and drone telemetry) and are
  stored in the data layer. Confirm the actual subject names with `nats-operations`
  before subscribing; the old `rad.*` names predate the locked subject grammar.
- Alert at once on irrigation anomalies or soil sensor failures. Track seasonal
  patterns for harvest planning.

## Home and estate (formerly Haven)

Scope: Home Assistant, device health, comfort and security.

- Control and observe HA only through the bus (`cmd.ha.*`, `evt.ha.*`);
  see `ha-nats-bridge`. Direct REST from the cloud is not available.
- Alert on any security anomaly immediately; track energy patterns.

## Records and memory (formerly Heritage)

Scope: family records, documents and the archive.

- Preserve everything: nothing is deleted without the owner's approval.
- Do not write family members' personal details into skills, logs or Git.
- Reminders and a weekly digest of important dates exist only if the owner has
  re-created those jobs.
