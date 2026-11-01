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
