# Alert delivery setup (Telegram + Matrix)

Alertmanager (kube-prometheus-stack release `monitoring`, namespace `monitoring`)
is configured in `platform/values/kube-prometheus-stack.yaml` under
`alertmanager.config`. One receiver, `onemind-notify`, fans out to:

| Channel  | How                                                                                  | Credentials (Secret, namespace `monitoring`)                       |
|----------|--------------------------------------------------------------------------------------|--------------------------------------------------------------------|
| Telegram | Alertmanager native `telegram_configs`, token read from file at send time            | `telegram-alertmanager-creds`, key `bot_token`                     |
| Matrix   | `webhook_configs` -> `matrix-alertmanager-receiver` (this directory) -> Synapse      | `matrix-alert-relay-creds`, keys `user-id`, `access-token`, `room-id` |

Routing: grouped by `alertname` + `namespace`; `critical` repeats every 1h,
`warning` every 6h; `info`, `Watchdog`, `InfoInhibitor` and anything without a
severity go to `null`. A `critical` alert inhibits `warning`/`info` with the
same `alertname` + `namespace`.

## What happens if the Secrets do not exist yet

* Alertmanager still starts: the Telegram Secret is mounted as an **optional**
  volume, and the token file is only read when a notification is sent. Until the
  Secret exists (and `chat_id` is set), Telegram sends fail and are logged and
  counted in `alertmanager_notifications_failed_total`; nothing else breaks.
* The Matrix relay Deployment stays in `CreateContainerConfigError` until its
  Secret exists, so Argo CD will report `platform-monitoring` as Degraded.
  Alertmanager just logs failed webhook notifications.
* Recommended order: do steps 1-3 (seal and commit both SealedSecrets, set
  `chat_id`) in the same PR/branch **before merging**, so delivery works from
  the first sync.

## Required Secrets (exact list)

| Secret                          | Namespace    | Keys                                  | File to commit                                                  |
|---------------------------------|--------------|---------------------------------------|-----------------------------------------------------------------|
| `telegram-alertmanager-creds`   | `monitoring` | `bot_token`                           | `platform/monitoring/telegram-alertmanager-creds.sealedsecret.yaml` |
| `matrix-alert-relay-creds`      | `monitoring` | `user-id`, `access-token`, `room-id`  | `platform/monitoring/matrix-alert-relay-creds.sealedsecret.yaml`    |

The Telegram `chat_id` is not secret; it lives in
`platform/values/kube-prometheus-stack.yaml` (`chat_id: -1  # TODO(user)`).

## 1. Telegram: bot token and chat ID

1. In Telegram open `@BotFather`, send `/newbot`, follow the prompts, copy the
   token.
2. Create a group (or use an existing one), add the bot, and send any message
   in the group (for example `/start@<your_bot_username>`).
3. Look up the chat ID (the token stays in your shell, nothing is printed except
   the ID; group IDs are negative, supergroups look like `-100...`):

   ```bash
   read -rs TG_TOKEN
   curl -s "https://api.telegram.org/bot${TG_TOKEN}/getUpdates" \
     | jq '.result[].message.chat | {id, title, type}'
   unset TG_TOKEN
   ```

4. Put the ID into `platform/values/kube-prometheus-stack.yaml`, replacing the
   `chat_id: -1` placeholder.

## 2. Matrix: bot user, room, access token

Matrix registration is closed, so create the bot with the Synapse admin tool.
The server name is `matrix.onemindos.com` per `platform/matrix/homeserver.yaml`;
confirm it from any existing user ID in Element (`@name:<server_name>`).

```bash
# Create @alertmanager:<server_name> (prompts for a password; not an admin)
kubectl -n matrix exec -it deploy/synapse -- \
  register_new_matrix_user -u alertmanager --no-admin \
  -c /config/homeserver.yaml http://localhost:8008
```

Log in once to obtain an access token (kept in shell variables, never printed):

```bash
read -rs MX_PASSWORD
MX_TOKEN=$(curl -s -X POST https://matrix.onemindos.com/_matrix/client/v3/login \
  -H 'Content-Type: application/json' \
  -d "$(jq -n --arg p "$MX_PASSWORD" \
    '{type:"m.login.password", identifier:{type:"m.id.user", user:"alertmanager"}, password:$p, initial_device_display_name:"alertmanager-relay"}')" \
  | jq -r .access_token)
unset MX_PASSWORD
[ -n "$MX_TOKEN" ] && [ "$MX_TOKEN" != null ] && echo "token captured" || echo "login failed"
```

Create the room in Element (https://element.onemindos.com) as yourself, e.g.
`#alerts:onemindos.com`:

* turn **end-to-end encryption OFF** (the relay cannot send encrypted events);
* invite `@alertmanager:<server_name>` (the relay joins the room on first send);
* copy the internal room ID from Room settings -> Advanced (`!abc...:<server_name>`).

## 3. Seal and commit both Secrets

Follows the "Secrets Convention" in the repo `CLAUDE.md` (namespace-wide scope).

```bash
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem

# Telegram (paste the bot token when prompted; nothing is echoed)
read -rs TG_TOKEN
printf '%s' "$TG_TOKEN" \
  | kubectl create secret generic telegram-alertmanager-creds -n monitoring \
      --from-file=bot_token=/dev/stdin --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/monitoring/telegram-alertmanager-creds.sealedsecret.yaml
unset TG_TOKEN

# Matrix (MX_TOKEN from step 2; set MX_USER and MX_ROOM_ID)
MX_USER='@alertmanager:matrix.onemindos.com'
MX_ROOM_ID='!REPLACE_ME:matrix.onemindos.com'
kubectl create secret generic matrix-alert-relay-creds -n monitoring \
  --from-literal=user-id="$MX_USER" \
  --from-literal=access-token="$MX_TOKEN" \
  --from-literal=room-id="$MX_ROOM_ID" \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/monitoring/matrix-alert-relay-creds.sealedsecret.yaml
unset MX_TOKEN
```

Commit the two `*.sealedsecret.yaml` files plus the `chat_id` edit. Argo CD app
`platform-monitoring` applies the SealedSecrets (its directory exclude `*.secret*`
only matches plain `*.secret.yaml` files, not `*.sealedsecret.yaml`).

## 4. Verify

```bash
# Secrets decrypted by the sealed-secrets controller (names only, no values)
kubectl -n monitoring get secret telegram-alertmanager-creds matrix-alert-relay-creds

# Relay running and logged in
kubectl -n monitoring rollout restart deploy/matrix-alertmanager-receiver   # only if it started before the Secret existed
kubectl -n monitoring logs deploy/matrix-alertmanager-receiver

# Alertmanager accepted the config
kubectl -n monitoring logs alertmanager-monitoring-kube-prometheus-alertmanager-0 -c alertmanager | grep -i "completed loading"
```

Send a test alert through the real route (creates a short-lived alert; resolves
itself after 5 minutes):

```bash
kubectl -n monitoring port-forward svc/monitoring-kube-prometheus-alertmanager 9093:9093 &
curl -s -XPOST localhost:9093/api/v2/alerts -H 'Content-Type: application/json' -d \
  "[{\"labels\":{\"alertname\":\"TestAlert\",\"severity\":\"warning\",\"namespace\":\"monitoring\"},\"annotations\":{\"summary\":\"Test alert from ALERTING-SETUP.md\"},\"endsAt\":\"$(date -u -v+5M +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -d '+5 min' +%Y-%m-%dT%H:%M:%SZ)\"}]"
kill %1
```

Expect a message in the Telegram group and the Matrix room after `group_wait`
(30s).

## Rotating credentials

Re-run step 3 for the rotated Secret. Telegram picks up the new file without a
restart. The Matrix relay reads its Secret at container start, so run
`kubectl -n monitoring rollout restart deploy/matrix-alertmanager-receiver`.

## After merge: one-time cleanup of control-plane leftovers

The control-plane scrapers (`kubeScheduler`, `kubeControllerManager`,
`kubeProxy`, `kubeEtcd`) are disabled in the chart values, but the objects that
the old hand-run Helm release created stay in the cluster (Argo CD runs with
`prune: false` and never owned them). Until they are removed the
`KubeSchedulerDown` / `KubeControllerManagerDown` / `KubeProxyDown` alerts keep
firing. Run once, after the Argo apps are Synced:

```bash
kubectl -n monitoring delete prometheusrule \
  monitoring-kube-prometheus-kube-scheduler.rules \
  monitoring-kube-prometheus-kubernetes-system-scheduler \
  monitoring-kube-prometheus-kubernetes-system-controller-manager \
  monitoring-kube-prometheus-kubernetes-system-kube-proxy \
  monitoring-kube-prometheus-etcd
kubectl -n monitoring delete servicemonitor \
  monitoring-kube-prometheus-kube-scheduler \
  monitoring-kube-prometheus-kube-controller-manager \
  monitoring-kube-prometheus-kube-proxy \
  monitoring-kube-prometheus-kube-etcd
kubectl -n monitoring delete configmap \
  monitoring-kube-prometheus-scheduler \
  monitoring-kube-prometheus-controller-manager \
  monitoring-kube-prometheus-proxy \
  monitoring-kube-prometheus-etcd
kubectl -n kube-system delete service \
  monitoring-kube-prometheus-kube-scheduler \
  monitoring-kube-prometheus-kube-controller-manager \
  monitoring-kube-prometheus-kube-proxy \
  monitoring-kube-prometheus-kube-etcd
```
