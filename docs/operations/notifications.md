# Portal notifications and alert routing (flagship F3)

Work order: `docs/roadmap/flagship/F3-notifications-alerts.md`. Source of truth for the rules is
one module, `services/onemind-portal/plugins/plugin-onemind-common-backend/src/notificationPolicy.ts`.
Everything below is **source behavior**. Nothing here is live until the F3 image is released
(`pipeline.py publish`, digest-pin PR, Argo) and the alert webhook secret exists.

## What reaches the inbox

| Event | Inbox | Severity | Topic | Scope (dedupe key) |
|---|---|---|---|---|
| Argo sync/operation succeeded | **never** (activity timeline + daily digest) | none | none | none |
| Argo operation failed or errored | yes | high | deployments | `argo:<app>` |
| Argo application Degraded | yes, once per incident | high | deployments | `argo:<app>` |
| Argo application OutOfSync for 15 minutes or more | yes, once per incident | high | deployments | `argo:<app>` |
| Watched backup CronJob suspended | yes, once per incident | high | backups | `backup:backstage/<cronjob>` |
| Backup job failed | yes | high | backups | `backup:backstage/<cronjob>` |
| Backup job succeeded | never (timeline) | none | none | none |
| Alertmanager alert, severity `critical` | yes | critical | alerts | `alert:<fingerprint>` |
| Alertmanager alert, severity `warning` | yes | high | alerts | `alert:<fingerprint>` |
| Alertmanager alert, severity `info` | yes | normal | alerts | `alert:<fingerprint>` |
| Alertmanager alert, any other severity (`none`, Watchdog) | never | none | none | none |
| Daily digest of successful deployments | one per day | normal | digest | `digest:deployments:<YYYY-MM-DD>` |

The same scope on the same producer updates one inbox item instead of adding another (the
notifications backend restores the existing row with the new text and marks it unread).

**Recovery.** When a condition clears, the producer re-sends the same scope with a quiet
`low` "Recovered: ..." or "Resolved: ..." text. Plugins cannot mark a user's item read
(the notifications API only accepts the user), so the updated item shows as an unread low item
until someone reads it. Filter by minimum severity to hide these.

**Which watches exist.** The observer only watches CronJobs named in
`onemind.operationalEvents.backupCronJobs` (default `backup-postgres-backstage`).
`backstage-backup` is suspended on purpose (superseded) and is not watched, so it never alerts.
Telling retired from pending from failing backups across the whole platform is F4's job.

## Alertmanager webhook

### Plugin-first finding (program rule 11)

| Candidate | Result |
|---|---|
| `@backstage/plugin-events-backend` (core, 0.6.x) | **Used.** Native HTTP ingress `POST /api/events/http/<topic>`; per-topic validators guard it (the GitHub/GitLab modules work this way); events service fans out to subscribers across replicas. |
| Alertmanager to notifications or signals module in the Backstage directory, `backstage/community-plugins` or npm | None found. Searched npm for `alertmanager`, `backstage alertmanager`, `events alertmanager`; no events or notifications module exists. |
| `@backstage-community/plugin-grafana` 1.1.0 (2026-09-24) | Reads Grafana alerts for display. Does not push inbox items or signals; not a receiver. Already installed for the Grafana card. |
| `@roadiehq/backstage-plugin-prometheus` 3.4.1 (2026-10-02) | Shows Prometheus graphs and alerts on an entity. Display only. |
| `@pagerduty/backstage-plugin` 0.21.0 | Requires PagerDuty as the middle service. Rejected: adds an external SaaS to a sovereign stack. |

Custom code is limited to the gap, each piece with one reason:
`createAlertmanagerValidator` (bearer check; the native ingress has no auth of its own),
`parseAlertmanagerWebhook` and `alertNotification` (no module maps Alertmanager v4 payloads),
`routeAlerts` plus a small state table (turns Alertmanager's periodic repeats into no-ops and
pairs resolved with firing). There is **no custom Express route**.

### Request path

```
Alertmanager  --Bearer token-->  POST /api/events/http/alertmanager   (events-backend ingress)
   validator: token present and >= 24 chars? header matches? body <= 512 KiB? payload shape valid?
   -> events service topic "alertmanager"
   -> subscriber (module of plugin "onemind"): one inbox item per fingerprint + Signal
```

Responses: `202` accepted; `401` missing or wrong token; `413` too large; `400` malformed;
`503` webhook not enabled (token not configured). The body is JSON-parsed by the events backend
before the validator runs (its 5 MB limit applies); no event is published unless the validator
passes.

### The secret (Zeus creates it; nothing is created by this PR)

The portal reads `onemind.alerts.webhookToken` from `${ONEMIND_ALERT_WEBHOOK_TOKEN}`
(`services/onemind-portal/app-config.yaml`). The variable is **optional**: when it is absent
the config key is dropped, the route answers 503 to everything and the portal starts normally
(tested against the real backend in `startup.test.ts` and against the real `app-config.yaml`
in `alertmanagerIngress.test.ts`).

| Namespace | Secret | Key | Consumer |
|---|---|---|---|
| `backstage` | `onemind-alert-webhook` | `ONEMIND_ALERT_WEBHOOK_TOKEN` | portal pod env (optional `envFrom`) |
| `monitoring` | `onemind-alert-webhook` | `token` | AlertmanagerConfig `httpConfig.authorization` |

Both hold the same random value (64 hex characters recommended; the route also requires at
least 24). `scripts/ops/seal-alert-webhook.sh` generates it in memory and seals it for both
namespaces with `kubeseal` (convention in `CLAUDE.md`); it is a dry run unless given `--apply`
and refuses to overwrite. The ordered steps, verification and rollback are in
`platform/monitoring/pending/README.md`.

### Alertmanager side (staged, not deployed)

Alertmanager is run by the operator from the `monitoring-kube-prometheus-stack` Application with
`alertmanagerConfigSelector: {}` and `alertmanagerConfigMatcherStrategy: None`, so every
`AlertmanagerConfig` in `monitoring` is merged into the live config (today only `openchoreo`).
The receiver is
`platform/monitoring/pending/alertmanagerconfig-onemind-backstage.yaml.disabled`:
`severity=~critical|warning`, excluding Watchdog and InfoInhibitor, `continue: true` so all
other routing is unchanged, one webhook call per alert, resolved notifications on.

It is deliberately outside anything Argo syncs:

- `.disabled` suffix: the `platform-monitoring` source reads `platform/monitoring` recursively
  but only `.yaml/.yml/.json/.jsonnet` files.
- `git grep` finds no kustomization, Application or manifest that names it.
- `scripts/tests/test_gitops_alert_routing.py` (runs in the pipeline gate) fails if it becomes
  reachable while the secret manifests are missing.

Enabling is a separate PR once the secret exists; steps are in `platform/monitoring/pending/README.md`.

## Signals

High and critical items also push a Signal. Channel `onemind:alerts`, delivered to the verified
members of `group:default/onemind-admins`:

```ts
// plugin-onemind-common-backend: AlertSignal, ALERT_SIGNAL_CHANNEL
{ action: 'alert.firing' | 'alert.resolved',
  source: 'alertmanager' | 'argocd' | 'backup',
  key: string,        // the notification scope
  severity: 'critical' | 'high' | 'normal' | 'low',
  title: string, at: string /* ISO time */ }
```

The home status strip (F1) and Live Operations (F4) subscribe with `useSignal('onemind:alerts')`
from `@backstage/plugin-signals-react`. The notifications backend additionally emits its own
`notifications` channel on every send, which the existing inbox UI already uses to refresh.
Signal publication needs no extra cluster permission.

## Daily digest and activity timeline

Successful Argo operations and backup jobs are written to `onemind_activity_timeline` (30 days)
by the observer, and external success events posted to `POST /api/onemind-events/events` are
recorded there too (accepted with `202 {delivery:"timeline", notified:false}`). Once per UTC day,
after `onemind.operationalEvents.digestHourUtc` (default 13, about 9 am New York), one `normal`
"Daily digest: N deployments succeeded" item lists the busiest applications and links to `/today`.
No successes in 24 hours means no digest.

`GET /api/onemind/activity/timeline?limit=50` returns the recorded entries. It requires a
signed-in `onemind-admins` member; service tokens and ordinary users get 401/403 (tests included).

## Preferences

`notifications.defaultSettings` in `app-config.yaml` registers the topics `alerts`,
`deployments`, `backups`, `digest` (origin `plugin:onemind`) and `deployments`, `backups`,
`agents` (origin `plugin:onemind-events`) with the native Notifications settings page, so each
person can mute any topic from User settings, Notifications, before the first item arrives.
Defaults are all enabled.

## Backlog cleanup (not run)

The existing unread "Deployment completed: ..." items (7,750 on 2026-10-21) can be marked read
with `services/onemind-portal/plugins/plugin-onemind-events-backend/scripts/mark-deployment-notifications-read.cjs`.
It is **admin-only by construction**: it acts with your own Backstage user token on your own
notifications, so it cannot touch anyone else's. It **only marks read**; it has no delete call,
and it never marks failures (title not `Deployment completed:`, or severity above `low`).

Run only after Zeus approves, from a checkout, with the user token copied from a signed-in
browser (any `/api/...` request's `Authorization: Bearer` header):

```bash
cd services/onemind-portal/plugins/plugin-onemind-events-backend/scripts
export PORTAL_URL=https://onemind-preview.tail717ba0.ts.net     # the portal base URL
read -rs PORTAL_TOKEN && export PORTAL_TOKEN                    # paste the token; not echoed
node mark-deployment-notifications-read.cjs                     # dry run: prints matched/kept counts
node mark-deployment-notifications-read.cjs --apply --expect=<matched count from the dry run>
unset PORTAL_TOKEN
```

`--apply` refuses to run unless `--expect` equals the count the same run just found. Expected
result: unread count falls below 50 (what remains are failures and other topics). Repeat for
each administrator account if more than one has a backlog.

## Tests

| Where | What |
|---|---|
| `plugin-onemind-common-backend/src/notificationPolicy.test.ts`, `alertmanager.test.ts` | Severity mapping, scope stability, payload parsing, link and text hygiene |
| `plugin-onemind-events-backend/src/scan.test.ts`, `conditions.test.ts` | Synthetic Argo success gives 0 notifications; failure gives 1 and a repeat gives no second; degraded, 15-minute out-of-sync, suspended backups, recovery, digest |
| `.../alertRouting.test.ts` | Firing then resolved on one scope, repeats ignored, Watchdog dropped |
| `.../alertmanagerIngress.test.ts` | Real events-backend ingress: firing then resolved with Signal within 5 s; 401 no or wrong token; 400 malformed; 503 token absent; real `app-config.yaml` loads without the env var |
| `.../startup.test.ts` | A real backend (events, signals, notifications plus our modules) starts without the secret and enforces it with it |
| `.../timelineRoute.test.ts`, `plugin-onemind-activity-backend/src/plugin.test.ts` | Admin only; anonymous, service and ordinary-user denials; successes never notify |
| `.../backlogScript.test.ts` | Dry run writes nothing; apply needs the expected count; no delete; predicate parity |
| `scripts/tests/test_gitops_alert_routing.py` | Receiver cannot reach Argo while its secret is missing |

## Not done here

- Browser-measured Signal latency and the live "unread below 50" result: needs the release and
  the cleanup run. The 5-second check above is server-side (HTTP to Signal publish).
- Frontend consumers of `onemind:alerts` (F1, F4).
- Enabling the receiver and creating the secrets (Zeus).
