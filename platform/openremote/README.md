# OpenRemote private IoT pilot

Source prepared November1,2026 at Zeus's request, beside an independently managed
ThingsBoard comparison. **Not deployed or adopted.** Manual Argo Application
`platform-openremote` owns this directory; there is no imperative install path.

## Initial scope

The upstream manager, bundled Keycloak and PostgreSQL17/Timescale images are
pinned by immutable multi-architecture digest (Linuxamd64 verified). An unprivileged
nginx proxy serves the manager and `/auth/` under one private origin:
`https://openremote.tail717ba0.ts.net`. Tailscale terminates HTTPS; Keycloak uses its
production `start` command, public HTTPS issuer and forwarded headers. Manager
`OR_DEV_MODE=false` and `OR_SETUP_RUN_ON_RESTART=false`; no demo setup type is selected.
The dedicated bundled Keycloak is required by this upstream stack; existing
Authentik, Hermes and their identities are not changed.

All HTTP listeners are ClusterIP except the Tailscale-owned web proxy. MQTT1883
is an internal manager Service port only; it is **not** exposed by the HTTPS
Ingress and NetworkPolicy does not admit device clients yet. No public MQTT,
HTTP or identity endpoint is created. Required device/broker/HTTP-agent egress
must be reviewed and added explicitly; outbound internet is not broadly enabled.
No NATS, Home Assistant, MCP, hardware, cloud-weather or robot-control integration
is configured. Native rules operate only on assets you deliberately create.

## Before first Argo sync

1. Review capacity on the ecosystem worker and retained-volume cost. This stack
   requests1568Mi and625m across four Pods; limits total less than5Gi. These are
   pilot settings, not measured production sizing or high availability.
2. Choose **fresh evaluation** or a separately reviewed data migration. The saved
   Mac volumes/configuration are retained, not copied into these claims. Do not
   copy live database files between architectures or PostgreSQL versions.
3. Supply private app and distinct backup-repository credentials through the
   [shared preparation guide](../iot-pilot-secrets.md). The empty `sealed/`
   Kustomization intentionally blocks Pod startup until real SealedSecrets exist.
4. Confirm the three image source locks in [image-lock.json](image-lock.json),
   Tailscale hostname availability, namespace policy and PVC classes. Images are
   Linuxamd64-compatible; the workload is selected onto the ecosystem pool.
5. Review/merge this PR, then manually sync `platform-openremote`. This does not
   authorize syncing unrelated Applications or deploying the sibling pilot.
6. Verify PostgreSQL, Keycloak9000 `/auth/health/ready`, manager HTTP and proxy
   probes; inspect sanitized logs. The manager bootstraps the empty schema/realm
   using the sealed administrator password. It must not clean an existing schema.

## First login and owner credentials

Credentials were generated on 2026-11-04 and sealed against the live controller. The
plaintext exists only in an owner-only (0600) JSON file at
`~/.config/onemind/recovery/iot-pilot/openremote-2026-11-04.json` on Zeus's Mac (it also
holds the Restic repositories and passwords); it was never printed or committed. Copy it
to the owner's encrypted credential store and keep an offline copy. To put the
administrator password on the clipboard without displaying it:

```sh
python3 -c "import json,sys; print(json.load(open(sys.argv[1]))[sys.argv[2]][sys.argv[3]], end='')" \
  ~/.config/onemind/recovery/iot-pilot/openremote-2026-11-04.json openremote-env ADMIN_PASSWORD | pbcopy
```
Paste it into the login field, then clear the clipboard.

Open `https://openremote.tail717ba0.ts.net` from a tailnet device and sign in to the
`master` realm as `admin` with that password (the manager and the bundled Keycloak
share it; confirm the username at first login). Then create a named owner user in a
new evaluation realm and change or retire the `admin` password in the store and in
Keycloak. Nothing creates owner accounts beyond this bootstrap administrator.

## Acceptance before real use

- Open the private HTTPS origin on two authorized tailnet clients. Verify login,
  issuer URLs and redirects remain on that origin; HTTP or internal Keycloak URLs
  must not leak into browser redirects. Confirm unauthenticated asset API access
  is denied. Verify self-registration stays disabled on each created realm.
- Create a new evaluation realm and a read-only service account. Do not use the
  administrator as the future Hermes/MCP identity. Confirm user/realm boundaries.
- Add one simulated asset and telemetry, restart workloads and verify persistence.
- Confirm no unexpected demo users/assets and no default bootstrap password remain.
- Complete [backup/restore acceptance](#recovery). First sync is not a backup receipt.
- Only then propose specific device ingress/egress and an adapter to existing NATS
  subjects. No new subject grammar or physical actuation is granted by this PR.

## Recovery

`database`10Gi and `storage`5Gi use `do-block-storage-retain`, `Prune=false`, and
separate VolSync Snapshot/Restic sources with separate encrypted repository prefixes.
Temporary movers/cache use the delete class. The database volume contains both
OpenRemote and Keycloak schemas. `storage` also matters: generated identity grant
material, logs and extensions/customized application files may live there.

Scheduled snapshots are **crash-consistent per PVC**, not an atomic database/storage
pair and not a verified application-level backup. For a coordinated recovery point,
stop manager/Keycloak through a reviewed temporary Git replica change, wait for
shutdown, take both volume backups while writers remain stopped, and record their
paired snapshot/repository IDs. Do not merely assume two daily timestamps match.
Before upgrades or important configuration changes, also keep a consistent logical
PostgreSQL dump, roles and matching storage archive in encrypted independent storage.
Use the matching OpenRemote PostgreSQL image/extensions for logical restore; a bare
PostgreSQL image cannot supply all Timescale/PostGIS extensions automatically.

Restore into isolated retained claims and a separate namespace, with matching
image locks and original app/Restic/sealing recovery material. Start database then
Keycloak then manager; test login, service-user identity and sample historical
attributes. Never run clean/setup-on-restart against recovered data. Prove this
procedure before relying on the pilot. No source volume deletion, pruning or live
restore is authorized. Spaces in the same provider/region is not the final independent
offline copy; independent recovery export remains a separate existing backup concern.

## Sources and untested boundaries

- [Official Kubernetes resources](https://github.com/openremote/core/tree/main/kubernetes)
- [Official stack](https://github.com/openremote/core/blob/main/docker-compose.yml)
- [Manager Dockerfile](https://github.com/openremote/core/blob/main/manager/Dockerfile)
- [Bundled Keycloak](https://github.com/openremote/keycloak)
- [Bundled PostgreSQL](https://github.com/openremote/postgresql)
- [AGPL license](https://github.com/openremote/core/blob/main/LICENSE.txt)

Configuration is adapted into plain GitOps resources to avoid the upstream setup
chart's generated plaintext Secret and mutable tags. Runtime routing/identity and
restore acceptance on DOKS remain unexecuted. Cloud hosting does not establish
local/off-grid operation; a later physical edge deployment must keep its own local
identity, rules, storage and device connectivity during WAN loss.

## Authentik sign-in (added 2026-11-08)

OpenRemote's own Keycloak (master realm) now has an identity provider `authentik` (OIDC). The Keycloak login page shows "OneMind (Authentik)"; the Keycloak admin login stays as break-glass and Keycloak is not replaced.

- `auth-provider.yaml` (Sync hook, namespace `auth`) creates the Authentik provider/application `onemind-openremote`, owner-only. Redirect: `https://openremote.tail717ba0.ts.net/auth/realms/master/broker/authentik/endpoint`.
- `idp-config.yaml` (Sync hook, wave 3) signs in to the Keycloak admin API with the existing admin secret (read by the Job only) and creates or updates the provider. The client secret is sealed in `auth/openremote-oidc-bootstrap` and `openremote/openremote-oidc`.
- Keycloak got egress to public HTTPS (Authentik) and an ingress rule for the hook.
- The first Authentik user is created in the master realm with no roles. Sign in once through Authentik, then sign in as the Keycloak admin and grant that user the roles needed (Users -> the user -> Role mapping).
- Rollback: delete the identity provider `authentik` in the master realm (Keycloak admin console), or remove `idp-config.yaml`.
