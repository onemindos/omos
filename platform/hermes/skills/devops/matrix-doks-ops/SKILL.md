---
name: matrix-doks-ops
description: "Operate the Matrix/Synapse stack on DOKS (namespace matrix) and wire the Legacy bot to Matrix, GitOps first."
version: 2.0.0
author: Legacy
license: MIT
platforms: [linux]
metadata:
  hermes:
    tags: [matrix, synapse, element, doks]
    related_skills: [doks-cluster-ops, hermes-mcp-config]
---

# Matrix on DOKS

## Status (verify live; SOUL rule 6)

Namespace `matrix`: Synapse, Element Web (3 replicas), PostgreSQL, Redis, all
managed by the `platform-matrix` Argo Application from `platform/matrix/` in the
infra repo. Server name `onemindos.dev`; Synapse at `matrix.onemindos.dev`, Element
at `element.onemindos.dev`; user IDs are `@zeus:onemindos.dev` (admin) and
`@legacy:onemindos.dev` (the bot). Federation is off. Secrets are SealedSecrets
in Git; credentials are never written into skills, memory or chat.

**Direction: Matrix is moving to tailnet-only (private) access.** Today the hosts
are behind ingress-nginx with Cloudflare. The target is the `tailscale` ingress
class like the other admin UIs, so clients (Element Web, Element X, Element
Desktop) need Tailscale. Do not add or widen public exposure, and do not publish
Matrix hostnames or room IDs outside the owner's channels. Read the current
Ingress objects before stating how it is reachable.

Hermes does not run on Bedrock and Matrix is one channel among others (Telegram is
kept too). Zeus's primary channel is Matrix; the Hermes gateway reaches it with
the bot account.

## Authority

- Read-only checks (`kubectl get`, logs, `whoami` calls) need no approval.
- Changes go through a PR to `platform/matrix/` and Argo CD, not `kubectl apply`
  or `kubectl patch` on live ConfigMaps. Zeus approves and merges.
- Deleting users, rooms, messages, media, PVCs, pods holding volumes, or the
  namespace needs Zeus to write "execute" or "override" for that specific action
  (SOUL rule 4). That includes the old "force-delete the pod and ReplicaSet" trick.
- Opening registration, even briefly, creating users, minting access tokens and
  rotating secrets are owner-led. Propose the exact steps and wait.

## Deploy gotchas (from past incidents)

1. Postgres DB must use `LC_COLLATE 'C' LC_CTYPE 'C' TEMPLATE template0`, or
   Synapse refuses to boot (`IncorrectDatabaseSetup`).
2. Set `PGDATA` to a subdirectory of the mount (`/var/lib/postgresql/data/pgdata`);
   the PVC root contains `lost+found` and initdb fails.
3. Element image is `vectorim/element-web:<tag>`; there is no `elementhq/element-web`.
4. Redis: use `--requirepass`; the config line `require-password` is fatal.
5. Init containers chown the data dirs to UID 991 (Synapse) and 999 (Postgres)
   on block storage.
6. Synapse logging: mount a `log.config` with a stderr StreamHandler and point
   `log_config` at `/config/log.config`; the default file handler crashes on a
   read-only filesystem (`templates/log.config`).
7. Write the whole YAML from a scratch file; do not patch accumulated edits.
   Duplicate keys (such as `enable_registration` twice) resolve to the last one.
8. `registration_shared_secret` is a plain string, not a nested mapping.
9. Do not put an `email:` block in the config unless it is complete
   (`notif_from` is required); Synapse validates it even when disabled.
10. The `matrixdotorg/synapse` image has no `/start.sh`; do not override
    `command:`, set `SYNAPSE_CONFIG_PATH` instead.
11. On rollouts the old pod holds the RWO volume; the new pod stays `Init` until
    the old one is gone. Wait, or ask Zeus for the go-ahead; see `doks-cluster-ops`
    (RWO deadlock) rather than force-deleting.
12. Stateful claims use `do-block-storage-retain`; never recreate a live claim to
    change its class.

`templates/homeserver-clean.yaml` shows the shape of a valid config with
placeholders for every secret. The live file is in Git.

## Creating a user (owner-led)

1. PR sets `enable_registration: true` with a shared secret from a SealedSecret.
2. After Argo syncs and Synapse restarts, run `register_new_matrix_user -c
   /config/homeserver.yaml -u <user> -a` inside the Synapse pod; the password is
   entered interactively by Zeus.
3. A second PR sets `enable_registration: false`. Confirm it applied.
4. The bot's access token goes into the Hermes runtime SealedSecret, never into
   a skill or chat.

## Hermes gateway wiring (Matrix platform)

Configuration for the `onemind-hermes` instance is seeded from Git
(`platform/hermes/onemind-hermes/`), not edited by hand in the pod.

- Platform dependencies are installed with Hermes's own lazy-dep installer
  (`tools.lazy_deps.ensure('platform.matrix')`), not by `pip install matrix-nio`,
  which pulls conflicting pins. Check with `is_available('platform.matrix')`.
- Config keys: `platforms.matrix.enabled`, `.homeserver`, `.user_id`, and the
  access token under **`token`**, not `access_token`. The connection gate reads
  `config.token`; a token stored as `access_token` lands in `extra` and Matrix
  silently never loads. Confirm with `hermes config get platforms.matrix`.
- Verify the token without printing it: `curl -s -o /dev/null -w '%{http_code}'
  -H "Authorization: Bearer $TOKEN" <homeserver>/_matrix/client/v3/account/whoami`
  should return 200.
- After a restart, `gateway_state.json` must list Matrix as `connected`.

Silent-failure checklist (absent from state, no log error): deps not installed;
token stored under `access_token`; gateway stuck draining (`hermes gateway
restart`, then reset the state file only if still stuck); default profile not
served; an earlier manual `pip install` left a conflicting `aiohttp-socks`; token
expired (Zeus re-issues it).

## Direct room with Zeus

Create the room as the bot with `preset: trusted_private_chat`, `is_direct: true`
and Zeus invited. Keep the room ID out of Git; the gateway discovers it.

## Decisions (locked)

- Plain Synapse plus Element Web. ESS Community was rejected (AGPL non-commercial
  terms plus an opinionated auth layer). No SSO layer until a customer needs it.
- Private, invite-only, federation off.
- Cloudflare and DNS changes use the approved Cloudflare path (OAuth MCP or a
  scoped token from a SealedSecret). A global API key is never used or stored.

## Templates

- `templates/homeserver-clean.yaml`: reference shape, secrets as placeholders.
- `templates/log.config`: stderr-only logging for Kubernetes.
