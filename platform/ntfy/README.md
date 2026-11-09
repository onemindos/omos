# ntfy

Self-hosted push notifications. Home Assistant, Frigate (through Home
Assistant), Alertmanager, Argo CD, n8n and Hermes publish to topics; the
owner's Android phone receives them over the tailnet.

- URL: `https://ntfy.tail717ba0.ts.net` (tailnet only, Tailscale ingress)
- Image: `binwiederhier/ntfy` v2.37.2, pinned by digest, hub-core only
- Storage: `ntfy-data` 2 GiB (`do-block-storage-retain`) at `/var/lib/ntfy`:
  message cache, user database, attachments (1 GiB total, 15 MiB per file,
  24 h expiry, sized for Frigate snapshots)
- Backup: VolSync source in `backup.yaml`, **paused** until its restic Secret
  is sealed
- Nothing leaves the cluster. There is no upstream relay, no email and no web
  push, and the only egress is DNS.

## Access control

Default access is `deny-all`. Users, tokens and ACLs are provisioned
declaratively (ntfy 2.14+):

| What | Where |
|------|-------|
| Users (bcrypt hashes), tokens and the ACL (`home-assistant` writes `home-alerts`, `cameras`; `alertmanager` writes `cluster`; `hermes` writes `hermes`) | Secret `ntfy-auth`, env `NTFY_AUTH_USERS`, `NTFY_AUTH_TOKENS`, `NTFY_AUTH_ACCESS` |

### TODO(owner): seal `ntfy-auth`

Keys needed, in namespace `ntfy` (beside `deployment.yaml`, add to
`kustomization.yaml`; pattern `platform/karakeep/*.sealedsecret.yaml`).
Generate the values locally with `ntfy user hash` and `ntfy token generate`
(or `docker run --rm -it binwiederhier/ntfy token generate`); never commit
plaintext.

- `NTFY_AUTH_USERS`: comma-separated `user:bcrypt-hash:role`, for
  `admin:<hash>:admin`, `home-assistant:<hash>:user`,
  `alertmanager:<hash>:user`, `hermes:<hash>:user`
- `NTFY_AUTH_TOKENS`: comma-separated `user:tk_<29 chars>:label`, one per
  publisher: `home-assistant`, `alertmanager`, `hermes`
- `NTFY_AUTH_ACCESS`: comma-separated `user:topic:permission`:
  `home-assistant:home-alerts:wo,home-assistant:cameras:wo,alertmanager:cluster:wo,hermes:hermes:wo`.
  It must ship with the users: ntfy refuses to start when an ACL names a
  user that is not provisioned (seen on first sync, 2026-11-08).

Until the Secret exists the pod still starts (the reference is optional), but
nobody can log in or publish.

## Subscribe on Android

1. Install ntfy from F-Droid or Google Play.
2. Settings, Default server (or the first subscription): add
   `https://ntfy.tail717ba0.ts.net`, log in as `admin`. The phone must be on
   the tailnet.
3. Subscribe to the topics below. Enable **Instant delivery** for each
   subscription: the app keeps its own connection open, with no Firebase or
   Google push, so nothing depends on a third party.
4. Allow the app to ignore battery optimisation (Android settings, Apps, ntfy,
   Battery, Unrestricted), or the connection is killed in the background.

Suggested topics: `home-alerts` (Home Assistant), `cameras` (Frigate snapshots,
through Home Assistant), `cluster` (Alertmanager, Argo CD, n8n), `hermes`.
Add a topic and its ACL entry together.

## Publish

```bash
# Replace $NTFY_TOKEN with the publisher's token (never commit it).
curl -H "Authorization: Bearer $NTFY_TOKEN" \
     -H "Title: Front door open" -H "Priority: high" -H "Tags: door" \
     -d "Opened at 22:14" \
     https://ntfy.tail717ba0.ts.net/home-alerts

# From inside the cluster:
curl -H "Authorization: Bearer $NTFY_TOKEN" -d "hello" \
     http://ntfy.ntfy.svc.cluster.local:8080/cluster

# Attachment (Frigate snapshot):
curl -H "Authorization: Bearer $NTFY_TOKEN" -T snapshot.jpg -H "Filename: snapshot.jpg" \
     https://ntfy.tail717ba0.ts.net/cameras
```

- **Home Assistant**: add the built-in *ntfy* integration (URL
  `http://ntfy.ntfy.svc.cluster.local:8080`, access token = the
  `home-assistant` token), then call `notify.send_message` or use the notify
  entity in automations. Frigate alerts reach the phone through a Home
  Assistant automation.
- **Alertmanager**: a `webhook_configs` receiver pointing at
  `http://ntfy.ntfy.svc.cluster.local:8080/cluster`, with
  `http_config.authorization.credentials_file` for the token. Alertmanager
  posts JSON, so ntfy shows the raw payload unless a template or a small
  relay (ntfy supports `?template=` with a Go template) is used.
- **Argo CD / n8n**: Argo CD notifications use a webhook service to the same
  URL; n8n uses an HTTP Request node.
- **Hermes**: plain `curl` or HTTP POST with the `hermes` token.

## Follow-up work (not in this change)

Publishing from each namespace needs its own change:

- Seal `ntfy-auth`. `ntfy-data-restic` is sealed via
  `backup-data.sealedsecret.yaml` and the VolSync source is unpaused
  (reseal the placeholder payloads with kubeseal — see DIAGNOSIS.md).
- `home-assistant`: add an egress rule to `ntfy:8080` (its NetworkPolicy only
  allows public HTTPS today), then add the integration.
- `monitoring`: Alertmanager receiver and an egress rule to `ntfy:8080`.
- `argocd`, `n8n`, `hermes`: egress rule and the notification service or
  node. The ntfy side already admits these namespaces.
- Add an ACL entry and token for any further publisher (`n8n`, `argocd`).
