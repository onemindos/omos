# F3 — Notification triage and real-time alert routing

Work order for one implementing agent. Branch: `flagship/f3-notifications-alerts`.
Program index: [README](README.md).

## Outcome Zeus sees

- **Inbox.** The bell shows only things that need a human: firing critical and
  warning alerts, failed or degraded deployments, failed backups and failed
  checks.
- **Successes.** Successful syncs are rolled into a daily digest, not 7,750
  separate items.
- **Live.** New alerts appear within seconds through Backstage Signals,
  without a page reload.
- **Resolve.** Alerts resolve themselves in the inbox when Alertmanager says
  they are resolved.

## Evidence

- `/notifications` has **7,750 unread**, almost all "Deployment completed:
  platform-backstage / app-of-apps", severity low.
- Producers: `plugins/plugin-onemind-events-backend/src/plugin.ts` and
  `plugins/plugin-onemind-activity-backend/src/plugin.ts` (`notificationService.send`).
- Alertmanager runs at `monitoring-kube-prometheus-alertmanager` (namespace
  `monitoring`). Signals are installed (`@backstage/plugin-signals(-backend)`).

## Build

1. **Severity policy** in one module: `deployment.succeeded` is never a
   notification; record it in the activity timeline only. `deployment.failed`,
   `degraded` and `out-of-sync for more than 15 minutes` are `high`. Backup failed
   or suspended is `high`. Alert severities map critical→`critical`,
   warning→`high`, info→`normal`.
2. **Dedupe and scope.** Use a stable `scope` per source object (for example
   `argo:<app>`, `alert:<fingerprint>`) so repeats update one notification
   instead of adding new ones.
3. **Alertmanager webhook receiver.** Add a backend route in the events backend
   that accepts Alertmanager webhooks. Authenticate it with a shared token held
   in a SealedSecret that **Zeus creates** (document the name and key). Map
   firing and resolved alerts to notifications. Add the Alertmanager receiver
   config through GitOps in `platform/monitoring/**` as a separate commit, routed
   only for `severity=~"critical|warning"`.
4. **Signals.** Push a signal on each new or resolved high/critical item.
   The home status strip (F1) and Live Ops (F4) subscribe to it.
5. **Daily digest.** One `normal` notification per day: N deployments
   succeeded, with a link to the timeline.
6. **Backlog cleanup.** Provide an admin-only action or script that marks
   existing `deployment.succeeded` notifications read. It runs only after Zeus
   approves, and never deletes anything.
7. **Preferences.** Expose the native notification settings for each topic
   (alerts, deployments, backups, digest).

## Files

`plugins/plugin-onemind-events-backend/**`, `plugins/plugin-onemind-activity-backend/**`,
`plugins/plugin-onemind-common-backend/**` (shared policy), tests,
`platform/monitoring/**` (receiver config only), `docs/operations/notifications.md`.

## Acceptance

- [ ] A synthetic Argo success produces 0 notifications. A synthetic failure
      produces 1, and repeating it does not add a second.
- [ ] A test Alertmanager payload (firing, then resolved) creates and then
      resolves one notification, and a signal reaches a subscribed browser in
      under 5 seconds.
- [ ] Unauthenticated or invalid-token webhook calls are rejected (negative test).
- [ ] After Zeus approves the backlog cleanup, the unread count is fewer than 50.
- [ ] Required gates pass.

## Status

2026-10-21, implemented by Claude (Sonnet 5.5) on `flagship/f3-notifications-alerts`. Source only: nothing
is released, deployed or enabled. Operator guide: `docs/operations/notifications.md`.

### Plugins considered (rule 11)

| Candidate | Maintainer, last release | New frontend | Decision |
|---|---|---|---|
| `@backstage/plugin-events-backend` HTTP ingress (+ `plugin-events-node`) | Backstage core, 0.6.7 on 2026-10-20 (0.6.6 installed) | n/a (backend) | **Used** for the webhook endpoint, bus fan-out and per-topic validator. No custom Express route. |
| `@backstage/plugin-notifications-backend` / `-node`, `plugin-signals-backend` / `-node` | Backstage core | n/a | **Used**: scopes dedupe, `notifications.defaultSettings` for per-topic preferences, `signalsServiceRef` for the new `onemind:alerts` channel. |
| Alertmanager to notifications or signals module | none found on npm, the plugin directory or `backstage/community-plugins` (searched `alertmanager`, `backstage alertmanager`, `events alertmanager`, `keywords:backstage-plugin alertmanager`) | n/a | Gap. Small custom code (below). |
| `@backstage-community/plugin-grafana` 1.1.0, 2026-09-24 | community | installed | Display of Grafana alerts only; not a receiver. Kept as is. |
| `@roadiehq/backstage-plugin-prometheus` 3.4.1, 2026-10-02 | Roadie | legacy | Entity graphs and alerts display only; not installed, not needed for routing. |
| `@pagerduty/backstage-plugin` 0.21.0, 2026-07-03 | PagerDuty | legacy | Rejected: adds an external SaaS between Alertmanager and the portal. |

Custom code, one reason each: bearer validator (native ingress has no auth of its own);
Alertmanager v4 payload parser and mapper (no module exists); `routeAlerts` plus `onemind_alert_state`
(Alertmanager repeats firing alerts, and resolved must pair with firing); observer extensions
(degraded, 15-minute out-of-sync, suspended backup, recovery) and the digest, because the existing
observer was the only producer; backlog script (the notifications API accepts only user tokens).

### Decisions

- Shared policy lives in `plugin-onemind-common-backend/src/notificationPolicy.ts` and
  `alertmanager.ts`; both producers (observer, event bridge) use it. Successes are timeline-only.
- The webhook is `POST /api/events/http/alertmanager` guarded by a topic validator. The events plugin
  leaves `/http` unauthenticated by design, so the validator is the only gate; it rejects before
  publishing. Token absent or shorter than 24 characters: 503 for every request.
- `packages/backend` is untouched: the package's default export is now a feature loader returning the
  observer module, the ingress module and the alert-routing module.
- Resolved alerts re-send the same scope as a `low` "Resolved: ..." update. Plugins cannot mark a
  user's notification read (user-only API), so it shows unread but quiet. Documented as a limitation.
- Alertmanager operator facts (read-only `kubectl --context do-nyc3-omos-cloud -n monitoring`):
  one `Alertmanager` CR (v0.33.1, 1 replica) from the kube-prometheus-stack release, `alertmanagerConfigSelector: {}`,
  `alertmanagerConfigMatcherStrategy: None`, one existing `AlertmanagerConfig` (`openchoreo`, `continue: true`),
  the config secret is operator generated. Webhook auth uses the CRD's `httpConfig.authorization`. No secret
  data was read.
- Live cluster inputs seen: 69 Argo applications (67 Synced/Healthy, `platform-matrix` OutOfSync,
  `platform-openclaw` Progressing); backstage CronJobs `backstage-backup` (suspended, retired) and
  `backup-postgres-backstage`. So the first scan after release would raise at most one out-of-sync item
  after 15 minutes. `backupCronJobs` therefore defaults to `backup-postgres-backstage` only.
- No live portal API was called and no notification was marked read.

### Zeus: secret to create (not created, not in Git)

Same random value in two namespaces; commands are in `platform/monitoring/pending/README.md`.

| Namespace | SealedSecret / Secret name | Key | Used by |
|---|---|---|---|
| `monitoring` | `onemind-alert-webhook` | `token` | AlertmanagerConfig `httpConfig.authorization.credentials` |
| `backstage` | `onemind-alert-webhook` | `ONEMIND_ALERT_WEBHOOK_TOKEN` | portal env var `ONEMIND_ALERT_WEBHOOK_TOKEN` (optional `envFrom`) |

Until both exist, the portal runs with the route disabled and Alertmanager config is unchanged.

### Done versus deferred

Done: severity policy and scopes; Alertmanager ingress and routing with dedupe and resolve; Signals on
`onemind:alerts`; daily digest and activity timeline (+ admin read route); degraded, out-of-sync and
suspended-backup conditions with recovery; native per-topic preference registration; dry-run backlog
script; staged receiver config with guard test; operator doc.

Deferred, needs Zeus or a later step:
- Create the two SealedSecrets; add the `envFrom` entry to `platform/backstage/deployment.yaml` and the
  kustomization entry; move the receiver out of `pending/` (all in one enable PR; steps in the pending
  README). Not done here because `platform/backstage/**` is outside F3's files and the receiver must
  not depend on a missing secret.
- Release (publish, digest pin, Argo) and signed-in acceptance: browser-measured Signal latency, a
  live Alertmanager test alert, and the post-cleanup unread count (< 50).
- Run the backlog cleanup after Zeus approves.
- Frontend consumers of `onemind:alerts` (F1 status strip, F4 Live Operations). The digest links to
  `/today` until a timeline page exists; data is at `GET /api/onemind/activity/timeline`.
- The pipeline's portal test list does not include the new suites; run `yarn workspace
  @internal/plugin-onemind-{common,events,activity}-backend test` (see counts below). Adding them to
  `scripts/gitops/pipeline.py` is outside F3's files.

### Tests

common-backend 44 (policy 15 new in two files), events-backend 39, activity-backend 10: 93 passing in
the three workspaces, plus 4 Python guard tests in `scripts/tests/test_gitops_alert_routing.py`.
Acceptance mapping: synthetic Argo success 0 notifications and failure 1 with no second on repeat
(`scan.test.ts`); firing then resolved creating and updating one item with a Signal under 5 s
server-side, plus 401/400/503 negatives (`alertmanagerIngress.test.ts`, real events-backend ingress);
startup without the secret (`startup.test.ts`, real backend). The "unread < 50" criterion needs the
live cleanup run.

### Gate

(filled in below after the final run)
