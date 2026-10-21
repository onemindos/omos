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

(Implementing agent: write progress here.)
