# OneMind activity inbox

The home inbox reads the existing authenticated Notifications API; the sidebar supports the full inbox. Notifications are durable in the notifications plugin database, backed up with the other Backstage databases. Existing Signals pushes updates to the full notification UI; the home summary polls every 30 seconds. Retention is 90 days.

## Event bridge

`POST /api/onemind-events/events` accepts only authenticated OneMind administrators or the dedicated external service subject `onemind-event-ingestor`. Existing MCP machine credentials are denied. There is no anonymous webhook and no producer credential is installed by this PR.

Payload fields (additional fields rejected):

```json
{"id":"argocd:platform-backstage:revision123:failed","kind":"deployment.failed","resource":"platform-backstage","occurredAt":"2026-10-07T12:00:00Z"}
```

Supported kinds: deployment.failed, deployment.succeeded, backup.failed, backup.succeeded, agent.completed. The ID must be stable for a source occurrence. Backstage notification scope collapses subsequent reports for the same kind/ID; concurrent producer delivery is not guaranteed exactly once. Recipients are fixed to onemind-admins. Links/severity/text are generated server-side; arbitrary provider payloads, secrets, and arbitrary recipients are not accepted. Timestamp must be within 24 hours; body limit 8 KiB, rate limit 30 requests/minute per authenticated actor per replica. The reverse proxy should add a global limit before exposing a high-volume producer. Events are historical reports, not live health.

## Connecting producers deliberately

1. Create a short-lived dedicated secret through the approved local secret workflow and install as a SealedSecret. Do not reuse MCP or organization-wide credentials.
2. Add `backend.auth.externalAccess` with subject `onemind-event-ingestor` restricted to plugin `onemind-events`; keep all existing restrictions. This credential cannot call the raw notifications API.
3. Configure Argo CD Notifications webhook service against the private endpoint with Authorization Bearer header, trigger only terminal deployment outcomes, and stable application/revision/operation IDs. Configure backup reporters only after their actual job outcome is known, preserving restore verification separately.
4. Have a narrowly scoped NATS consumer translate only approved completed-work subjects into agent.completed; do not ingest raw telemetry or commands. No NATS credentials/subscription are activated here.
5. Verify successful ingestion, retry scope, intended group delivery, anonymous 401, unrelated service 403, invalid payload 400 and rate 429 in staging before connecting sources. No runtime tests were run for this PR.

Producers remain configuration-required until those connections are made. An empty inbox does not mean systems are healthy.

References: https://backstage.io/docs/notifications/ and https://backstage.io/docs/notifications/usage/ . Implementation follows installed notificationService 0.2.30 types and installed Notifications API 0.6.0.
