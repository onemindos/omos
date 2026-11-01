# ThingsBoard Community4.3 private comparison pilot

Source prepared November1,2026 at Zeus's request, **not deployed**. This is a
fresh, independently managed comparison beside OpenRemote, not an accepted final
IoT authority or automatic migration of retained Mac data.

## Edition and scope

Pinned upstream `thingsboard/tb-node:4.3.1.8` is the published pre4.4 Community
Apache2.0 line, with a verified Linuxamd64 image. The official CE Docker guide
currently specifies this tag. Its dedicated PostgreSQL17 database persists the
configuration and telemetry; no Kafka or external cloud service is added.

Queue type is explicitly `in-memory`: this is an evaluation monolith. In-flight
messages may be lost on restart; database persistence does not make that queue
durable. Production durability/throughput and backup sizing require a separate
review. Do not silently switch to4.4: the vendor merges editions into BSL
source-available licensing there, with license-key/device limits and online
validation implications. The published4.3LTS security-support horizon is
July20,2027; review supported options before that date.

Private proposed UI/API: `https://thingsboard.tail717ba0.ts.net`, using Tailscale
HTTPS and native application authentication. No public listener is created.
MQTT1883 is ClusterIP only and not admitted from devices by default NetworkPolicy;
the web Ingress does not expose MQTT. Additional transports and cloud/broker/NATS
connectivity are deliberately unconfigured. No device actions, firmware update,
Hermes adapter or changes to existing NATS contracts are included.

## Safe initial schema/account creation

The upstream installer creates a standard system administrator even when demo
loading is off. Exposing that initial account is not acceptable here.

`thingsboard-bootstrap-v1` is an Argo Sync hook at wave1:

1. The PostgreSQL guard accepts an existing OneMind initialization marker without
   reinstalling or changing credentials. An empty database is eligible for install.
   A nonempty database without that marker is rejected, including an interrupted
   install or unreviewed local-data import. **Nothing is dropped or reset.**
2. Only the hook's installer has `INSTALL_TB=true`, with `LOAD_DEMO=false`.
3. Before the application/Ingress waves can start, a transaction replaces the
   sole initial system administrator's email and BCrypt password hash from a
   SealedSecret, then commits the initialization marker. No default demo tenant
   is loaded. Plaintext administrator passwords are not put in Kubernetes Secrets.

The normal server has no install/upgrade flags. Repeated syncs do not reinstall
schema or reset owner credentials. A failed/partial install stops for diagnosis;
review the database and restore/recovery path rather than forcing the guard.
Changing the sealed bootstrap file later is not a password rotation mechanism.

## Before launch and acceptance

- Review available ecosystem capacity. App/database steady requests are1792Mi
  and600m; installer temporarily adds up to2Gi plus the database. This is not HA
  or benchmarked sizing. Account for the sibling OpenRemote and other pilot PRs.
- Choose fresh evaluation versus a separate reviewed migration. No saved Mac
  volume is mounted or imported. Existing retained4.3.1.7 data requires a
  supported, separately validated upgrade/migration; this guard refuses adoption.
- Follow [private credential preparation](../iot-pilot-secrets.md), including a
  real owner email/password, random Base64 JWT signing key and independent Restic
  repository. Never use upstream example passwords.
- Merge reviewed source, then manually sync `platform-thingsboard`. Inspect hook
  completion before application readiness and private HTTPS availability.
- Verify the chosen owner can sign in, upstream default credentials fail, no demo
  tenant/users were loaded and JWT settings are private. Create a tenant, a limited
  tenant identity and one simulator device deliberately after onboarding.
- Submit a small HTTP telemetry sample over the authenticated device endpoint,
  view its dashboard/history, restart and verify persistence. These operations
  remain unexecuted; availability does not prove integration.
- Validate restrictive network behavior, backup results and isolated recovery
  before real sensors, alerts, RPC/commands or long-term retained data.

## Recovery

The10Gi database claim uses `do-block-storage-retain` and `Prune=false`. Its
VolSync Snapshot/Restic source captures the whole PostgreSQL data directory with
WAL, to its own encrypted repository; temporary/cache claims use `do-block-storage`.
This is a crash-consistent physical backup, not a verified logical dump/restore.
Application logs and transient renderer data are emptyDir; important application
entities are in PostgreSQL. Custom files/extensions added later require explicit
backup scope changes.

Before upgrades, preserve a consistent logical dump and database roles plus the
exact image lock, bootstrap marker and independent JWT/SealedSecrets recovery
material. Restore a backup into an isolated matching-version PostgreSQL instance
and test owner login, a sample device and historical telemetry with the app. The
initialization hook accepts only the original marker; it must not reset recovered
credentials. Do not restore a running database directory by file copy or run the
installer on an existing unmarked schema. Record recovery evidence before real use.
The original Mac data and cluster PVCs must not be pruned as part of this pilot.

## Primary sources

- [Official CE Docker installation](https://thingsboard.io/docs/installation/docker/)
- [Pinned4.3.1.8 startup script](https://github.com/thingsboard/thingsboard/blob/v4.3.1.8/msa/tb-node/docker/start-tb-node.sh)
- [Pinned initial-account loader](https://github.com/thingsboard/thingsboard/blob/v4.3.1.8/application/src/main/java/org/thingsboard/server/service/install/DefaultSystemDataLoaderService.java)
- [Pinned schema](https://github.com/thingsboard/thingsboard/blob/v4.3.1.8/dao/src/main/resources/sql/schema-entities.sql)
- [Licensing/support change](https://thingsboard.io/blog/one-thingsboard-source-available/)
- [CE4.3.1.8 license](https://github.com/thingsboard/thingsboard/blob/v4.3.1.8/LICENSE)
