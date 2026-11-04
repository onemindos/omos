# CloudTAK

**Live acceptance verified 2026-10-25 18:23 UTC:** all four workloads are Ready; browser/API/tiles respond over private TLS. Supported first-run setup completed against the existing TAK server. The dedicated CloudTAK login returns an administrator session, its own profile is verified, and the TAK connection reports **live**. Existing CA and hostname checks pass; no test CoT messages or assets were created. See [the signer repair and rollout evidence](../../docs/operations/tak-enrollment-repair-2026-10-25.md).

**Storage provisioned 2026-10-24:** private NYC3 bucket `onemind-cloudtak` exists. The dedicated Read/Write/Delete key is strictly sealed as `cloudtak/cloudtak-spaces` in this directory. Own-bucket listing and private ACL were verified; the temporary provisioning key was removed. This is credential/storage readiness, not application or recovery acceptance. No existing key was changed.

A fresh install of upstream [CloudTAK](https://github.com/dfpc-coe/CloudTAK) (`dfpc-coe/CloudTAK`, pinned
**v13.113.0**, released 2026-10-23) as a browser TAK client. It runs under the Argo CD Application
`platform-cloudtak` (sync wave 5) in namespace `cloudtak`. Dated live observations are recorded above; the rest of this runbook describes its deployment contract.

The branch includes current `main`, including the fresh STAC source in `platform/stac/`. CloudTAK and STAC use separate databases and roles on the existing `onemind-data` cluster; this PR creates no PostgreSQL server or PVC.

**Initial rollout authorized 2026-10-25.** Keep `platform-cloudtak` manual with retry limit zero for the
observed first deployment. The shared CNPG primary and both replicas are healthy; the 2026-10-25 daily backup
completed successfully. Exact pinned CloudTAK migrations and all four image entrypoints passed on an isolated
instance of the exact current PostgreSQL image (see the dated evidence below). Activate the existing CloudTAK
managed role through the shared GitOps application, verify the sealed DSN/password binding without printing it,
then sync this application once. These are deployment prerequisites; local proof is not live application acceptance.

CloudTAK here is **a client only**. TAK Server stays the source of truth: CloudTAK holds a streaming connection to it
per user, shows the map, chat, missions and data packages, and keeps its own small database for profiles, basemaps
and imports. It does not replace TAK Server, and it does not bring its own copies of the shared services below.

This is not the old OneTAK stack. `platform/onetak/` no longer runs a CloudTAK API; nothing here shares code or data
with it.

## What runs

| Deployment | Image (pinned by tag and digest) | Port | Job |
| --- | --- | --- | --- |
| `cloudtak-api` | `ghcr.io/dfpc-coe/cloudtak-api` | 5000 | Web UI (nginx) and API, one process in `CLOUDTAK_Server_Mode=both`. Single replica, `Recreate` |
| `cloudtak-tiles` | `ghcr.io/dfpc-coe/cloudtak-pmtiles` | 5002 | PMTiles tile server for imported data and hosted basemaps |
| `cloudtak-events` | `ghcr.io/dfpc-coe/cloudtak-events` | none | Import worker: KML, GeoJSON, Shapefile, GeoTIFF, MBTiles, data packages into PMTiles |
| `cloudtak-retention` | `ghcr.io/dfpc-coe/cloudtak-retention` | none | Daily purge job (03:00 UTC); does nothing until an admin enables retention in the UI |

The images are public on GHCR, `linux/amd64` only (all cluster nodes are amd64). They are built by upstream's GitHub
Actions on each tag, so there is no build step in this repo. Exact tags and digests are in the Deployments.

## Shared services it uses (not deployed here)

| Service | Where | How CloudTAK reaches it |
| --- | --- | --- |
| TAK Server 5.8 | namespace `takserver` | Set in the UI at first run (see below). 8443 Marti API with the admin client certificate, 8446 WebTAK login, 8089 streaming CoT |
| PostgreSQL + PostGIS | `onemind-data` cluster, namespace `onemind-data` | Database `cloudtak` (`database.yaml`); activate the existing managed role first, `POSTGRES` from `Secret/cloudtak-db` |
| Tile38 | `tile38` in namespace `onetak`, port 9851 | Settings `geofence::enabled` and `geofence::url`, set from env (`CLOUDTAK_Config_geofence_*`) |
| Object storage | DigitalOcean Spaces, nyc3, bucket `onemind-cloudtak` | `ASSET_BUCKET`, `AWS_S3_Endpoint`, `AWS_S3_AccessKeyId`, `AWS_S3_SecretAccessKey`. These are upstream's S3-protocol names; nothing talks to AWS |
| Media server | `onetak-media` (MediaMTX) in namespace `onetak` | **Not connected.** See "Media" |

Access is tailnet only, through the Tailscale operator (ingress class `tailscale`, as for n8n and the other admin UIs):

- `https://cloudtak.tail717ba0.ts.net` : web UI and API (`API_URL`)
- `https://cloudtak-tiles.tail717ba0.ts.net` : tile server (`PMTILES_URL`). The browser loads tiles from this second
  origin, so it needs its own name.

There is no public Ingress. `network.yaml` denies all traffic in the namespace and then allows only: DNS; the Tailscale
proxy pods in; API to Postgres 5432, TAK Server 8443/8446/8089, Tile38 9851, and public HTTPS 443 (Spaces, basemap
and search proxying; private, tailnet and metadata ranges excluded); tiles and events to the API and public 443;
retention to the API. No media egress is opened until a compatible media replacement is reviewed.

## Database and credential prerequisites

The shared `onemind-data` Cluster must declare the following managed role before syncing CloudTAK. The existing
sealed `db-cloudtak` password lives with the shared resources; do not generate a new password or grant
superuser, database-creation or role-creation rights.

```yaml
    - name: cloudtak
      ensure: present
      login: true
      superuser: false
      createdb: false
      passwordSecret:
        name: db-cloudtak
```

Sync the owning shared-resource Application and verify CNPG reconciles this role. CloudTAK's wave-2 Database
resource creates only `cloudtak`, retains it on resource deletion and enables PostGIS as the operator before the
wave-3 API starts. The tiles, import and retention workers follow in wave 4; private ingresses follow in wave 5.
The exact-image proof used the same nonsuperuser role privileges with PostGIS preinstalled.

All Spaces prerequisites were completed on 2026-10-24: private NYC3 bucket `onemind-cloudtak`, restricted listing,
CDN off and a dedicated bucket-scoped Read/Write/Delete key. The strictly sealed
`cloudtak/cloudtak-spaces` credential is committed beside the workloads. No manual Secret or new key is required
for this rollout. Check `SealedSecret` reconciliation and `Secret` key names after sync without printing values.

The PR carries sealed `cloudtak-signing` (the JWT secret shared by all four workloads), `cloudtak-db`,
and the role password in `onemind-data/db-cloudtak`. Validate the role/DSN binding at first deployment without
printing the credential values. To rotate the signing secret, reseal `cloudtak-signing` and restart all four Deployments; every
user session and API token is signed with it, so everyone has to log in again.

## First run (existing TAK administrator, in the CloudTAK UI)

After the Application syncs and `cloudtak-api` is Ready (the pinned first start applies 226 database migrations and
loads the default iconsets and basemaps; allow a few minutes), open **https://cloudtak.tail717ba0.ts.net** from a
device on the tailnet. CloudTAK shows **Initial Server Configuration**. Until it is submitted, anyone who can reach
the page can configure it, so do it straight away.

| Field | Value |
| --- | --- |
| Server Name | any label, for example `OneMind TAK` |
| Server Hostname | `takserver` (the app-local alias of the existing canonical TAK proxy) |
| CoT Port | `8089` |
| Marti Port | `8443` |
| WebTAK | `8446` |
| Admin Certificate | the TAK Server admin client certificate (`.p12`) and its password |
| Initial Administrator Username / Password | dedicated TAK login `CloudTAK`; final credential is in strict `cloudtak/cloudtak-tak-login` and private recovery material. It becomes the CloudTAK system administrator |

The app-local `Service/takserver` is an `ExternalName` pointing to
`tak-tailnet.takserver.svc.cluster.local`; it follows the existing `tak-lb-proxy` service for all three ports.
Bootstrap URLs are `ssl://takserver:8089`, `https://takserver:8443` and `https://takserver:8446`.
The short hostname matches the existing TAK certificate SAN. The full canonical service hostname is absent
from that certificate, so use this app-local alias rather than disabling hostname checks or rotating the server
certificate. No platform DNS or NetworkPolicy scope is changed.

Password-auth SDK calls use ordinary TLS validation. The API mounts the existing private TAK CA from strict
`Secret/cloudtak-tak-ca` and loads it with supported `NODE_EXTRA_CA_CERTS`; no certificate key is included in
this CA secret. This extends trust for this app without disabling certificate/hostname verification or changing
system-wide TLS defaults. The pinned SDK's certificate-auth and CoT paths still disable verification internally;
that upstream behavior is a separate limitation, not a reason to disable validation for password authentication.
Strict chain/hostname preflight must pass for `takserver` before the first-run request.

**Dedicated login provisioned 2026-10-25.** The pre-existing certificate/File principal `CloudTAK` already belonged
to `OneMind` and had no reusable password. The LDAP directory had no user accounts. Under the user's explicit
TAK-discovery/deployment authorization, only `uid=CloudTAK,ou=People,dc=onemindos,dc=dev` and
`cn=OneMind,ou=Group,dc=onemindos,dc=dev` were added. The account uses `inetOrgPerson` plus `takUser`,
`takRole=Team Member`, and the same OneMind group through the existing `memberOf` overlay. It receives no TAK
administrator role. Existing human/File users, CoreConfig and global LDAP policy were preserved.

The initial random password was replaced through authenticated LDAP password modify, respecting the existing
mandatory first change. A final-password LDAP bind, group/policy metadata and the supported TAK
`POST /oauth/token` password grant on 8446 passed; no credential or token was printed. The strict SealedSecret
contains only the final username/password. This verifies the TAK login prerequisite, not CloudTAK first-run acceptance.

**Password rotation.** Existing policy requires at least 12 characters, retains six prior passwords and expires
passwords after 30 days (warning starts seven days before expiry). Before that deadline, use authenticated
`ldappasswd` as this dedicated account to set a fresh random password, verify bind/token login, reseal
`cloudtak-tak-login` and update its restricted recovery copy. Preserve global policy and all human credentials.
CloudTAK does not perform this LDAP password change automatically; fresh username/password login fails after
expiry even if an existing CloudTAK session or issued certificate remains usable. Do not rotate signing secrets
or edit profile authentication rows as a password-change substitute.

**Login.** CloudTAK has no user database of its own. Users sign in with their TAK Server username and password: it
posts them to TAK Server's WebTAK endpoint (8446), and on success downloads that user's client certificate and uses it
for their streaming connection. If TAK Server is backed by LDAP (`tak-ldap`), LDAP users sign in the same way. Optional
OIDC and passkey sign-in exist in **Admin, Config**; they are off. The "external provider" (`provider::*`) is a
COTAK-specific integration and stays unset.

What is set from the environment and what is set in the UI:

- Environment (this repo): database, signing secret, Spaces, `API_URL`/`PMTILES_URL`, and the two Tile38 settings.
  Anything set as `CLOUDTAK_Config_<key>` overwrites the stored value on **every** start, so change those in
  `api.yaml`, not in the UI.
- UI (stored in the database): the TAK Server URLs and admin certificate, system administrators, basemaps, retention,
  login branding, OIDC, search and routing providers.

## Media replacement gate

CloudTAK's video features require the [media-infra wrapper](https://github.com/dfpc-coe/media-infra/tree/v9.12.1),
not plain MediaMTX. The reviewed v9.12.1 wrapper exposes port **9997** (`/path`, `/path/:name`,
`/v3/config/global/get`, `/stream/<name>/index.m3u8`) and uses the same `SigningSecret` to authenticate with CloudTAK.
The current `platform/onetak/media-mtx.yaml` is plain MediaMTX, so `media::url` remains unset and video leases are off.

A direct image swap is unsafe for existing clients. The wrapper changes all RTSP/RTMP/SRT reader and publisher
requests to HTTP authentication through CloudTAK's `/api/video/auth`; it also periodically removes MediaMTX paths
that have no CloudTAK lease (`lib/persist.ts`). The existing takml vision worker uses the unauthenticated named paths
`drone1`, `takcam`, `helmcam` and `maccam`. CloudTAK's video leases have generated identities, so preserving those names
requires an explicit stream mapping or compatible alias/proxy design, plus credentials for every reader/publisher.

Replacement scope before enabling media:

1. Inventory publishers/readers and existing path/recording state; preserve data and agree a maintenance window.
2. Initialize CloudTAK, create corresponding video leases and distribute scoped reader/publisher credentials. Update
   takml `FEEDS` and other clients to the resulting path identities; do not silently disable authentication.
3. Pin a tested media-infra image; share only the required signing secret, set its internal `API_URL`, and allow
   wrapper-to-CloudTAK authentication plus CloudTAK-to-wrapper port 9997 in NetworkPolicies.
4. Preserve existing media Service endpoints during cutover, retain any recording storage with a backup, and verify
   publishing, takml reads, HLS, RTSP and authenticated browser playback. WebRTC additionally needs a reachable ICE
   transport (8189 UDP/TCP or a reviewed TURN service); HTTPS ingress alone does not carry ICE.

This PR leaves the existing shared media server and its clients unchanged. No second media server, recording PVC,
public port, shared-secret rotation or media egress permission is added.

## The AWS-only "layers" feature

CloudTAK's ETL "Layers" create AWS Lambda functions and CloudFormation stacks and use SQS. No AWS credentials are supplied. Public HTTPS egress can reach AWS endpoints as well as Spaces; it is not a provider allowlist. The layer plumbing is switched off or inert:

- `cloudtak-api` starts with `--no-events --no-etl-events` (no Lambda cron invocations, no SQS delivery).
- `AWS_EC2_METADATA_DISABLED=true`, no AWS credentials in any pod, and the network policy blocks the metadata address.
- Creating or running a layer in the UI would fail. Do not use **Layers**; the Layers page may show errors. There is
  no setting that hides the feature.

## Known limitations

- **Imported-data tiles are blocked with this configuration.** The API fetches TileJSON for user-imported assets from `PMTILES_URL` using
  upstream's SSRF-safe fetch, which refuses private and tailnet addresses, and pods cannot reach the tailnet at all.
  The import worker only needs the API and Spaces; the map layer for an imported asset is the part at risk.
  Hosted PMTiles basemap metadata also uses the external tile origin, so those layers need the same reachability
  fix. Public upstream basemaps, chat, missions and the TAK connection do not use this fetch. A complete fix needs a
  pinned upstream change supporting a trusted internal tiles origin while returning the browser-facing origin,
  with scoped SSRF allowance and matching NetworkPolicies. Do not disable `StackName` or SSRF checks to bypass it.
  Acceptance of imported/hosted data layers is blocked until that change is tested.
- Single API replica. Upstream's split mode (`api` + `hub`) exists for scaling, but is not needed here.
- The image entrypoints/workdirs and required environment variables were checked against pinned source and registry
  metadata. Non-root/read-only runtime compatibility still needs a live start; inspect the failing path before making
  a narrowly scoped permission or mount change. The events scratch volume has `fsGroup: 1000` for its non-root writer.

## Resource and recovery budget

The four app pods request **260m CPU / 960 MiB memory**, plus about 50m/128 MiB for the two Tailscale proxies.
They prefer the existing ecosystem node pool. Limits are API 1 CPU/1536 MiB, tiles 500m/512 MiB, events 2 CPU/2 GiB,
and retention 200m/256 MiB. No new node pool or persistent app volume is added. Imports use at most 12 GiB ephemeral
storage per events pod (10 GiB scratch plus 256 MiB `/tmp`), with a 1 GiB scheduling request.

The tiles and retention commands bypass upstream npm scripts that force a 6144 MiB heap; V8 caps are 320/128 MiB.
The API cap is 1024 MiB. Events use a 768 MiB heap per Node isolate and exactly **one** import worker, instead of the
upstream host-CPU default; GDAL and tippecanoe still share the 2 GiB container budget. Large imports may require a
measured resource change. API startup allows 450 seconds for first migrations; tiles has its own startup window.

CNPG's existing cluster-wide WAL/object-store backup covers the new database after creation. Confirm the latest
successful backup before schema upgrades. The new Spaces bucket's object retention/offsite export must be defined
before storing irreplaceable imagery; no independent bucket backup or restore test is provided by this PR.

## Security notes

- **SSRF advisories** (GHSA-vqrw-qphh-p34v / CVE-2026-54546 on `PUT /api/basemap`, GHSA-r95q-fp26-h3hc /
  CVE-2026-55177 on `/api/esri*`): the guards are in v13.113.0 and apply to this deployment. They are skipped when the
  `StackName` environment variable is unset or `test`, which is what upstream's docker-compose does, so
  `cloudtak-api` sets `StackName=cloudtak`. Do not remove it.
- Pods have no Kubernetes API token and reach the internet only on 443.

## Upgrading

1. Read the upstream release notes between the pinned tag and the target (`CHANGELOG.md` in the repo). Look for new
   or renamed environment variables and for database migrations, which are forward-only.
2. Get the new multi-arch digests for all four images (all four must move to the same tag):

   ```bash
   for n in cloudtak-api cloudtak-pmtiles cloudtak-events cloudtak-retention; do
     T=$(curl -s "https://ghcr.io/token?scope=repository:dfpc-coe/$n:pull" | python3 -c 'import sys,json;print(json.load(sys.stdin)["token"])')
     printf '%s ' "$n"
     curl -s -D - -o /dev/null -H "Authorization: Bearer $T" \
       -H 'Accept: application/vnd.oci.image.index.v1+json' \
       "https://ghcr.io/v2/dfpc-coe/$n/manifests/vX.Y.Z" | grep -i '^docker-content-digest'
   done
   ```

3. Edit the `image:` line in `api.yaml`, `tiles.yaml`, `events.yaml` and `retention.yaml` to
   `ghcr.io/dfpc-coe/<image>:vX.Y.Z@sha256:<digest>`. Never `:latest`.
4. Open a PR. After review and an authorized manual sync, Argo restarts `cloudtak-api` (Recreate, so a short outage); it runs the migrations on start, so
   watch `kubectl -n cloudtak logs deploy/cloudtak-api` until `/api` answers. Roll back by reverting the PR, but
   note that a migration that already ran is not undone; take a Postgres backup first when the notes mention one.

## Isolated image proof — 2026-10-25

This proof used Docker Desktop on the operator Mac, running the production `linux/amd64` images under emulation
on its arm64 host. The Docker network was internal, with no published ports and no connection to production
PostgreSQL, TAK Server or Spaces. Only disposable credentials and data were used.

- PostgreSQL: `ghcr.io/onemindos/onemind-postgres` at digest
  `sha256:35302d70df121bc86b370ebad800d85401edcdb890f8146cd3db690826a9925d`, matching the live
  `onemind-data` ImageCatalog and `status.image`; PostgreSQL 16.15, `shared_preload_libraries=age`,
  `restart_after_crash=off`, TLS enabled, `shared_buffers=256MB`, `max_connections=150`.
- Bootstrap: synthetic `cloudtak` login with `NOSUPERUSER NOCREATEDB NOCREATEROLE`, an empty database it owns,
  and PostGIS installed by the test administrator before app startup. Local-only host authentication was trusted
  on the internal Docker network; no production credential was copied into the fixture.
- API: the exact pinned image and manifest command, `StackName=cloudtak` (SSRF guards retained), `both` mode,
  1 CPU and 1536 MiB memory cap. All **226** migrations completed. The API returned 200/version 13.113.0,
  `/api/server` correctly returned `unconfigured`, and the web UI returned 200. Seeding produced 5 basemaps,
  10 iconsets and 3677 icons.
- Repeat startup: restarting only the API returned to 200; migration and seed counts stayed unchanged, and the
  database postmaster start time stayed unchanged. No native PostgreSQL crash or container OOM occurred.
- Tiles, events and retention: the exact pinned images and manifest commands ran as UID/GID 1000 with read-only
  root filesystems, all capabilities dropped, no privilege escalation and writable temporary/scratch mounts.
  Each used its manifest CPU/memory cap. `/tiles` returned 200, the import worker polled the API without errors
  using exactly one worker, and retention started its daily scheduler.

This establishes local migration and container compatibility. It does not verify native amd64 cluster behavior,
Kubernetes policies/probes, Spaces reads/writes, private HTTPS, a TAK login or live CoT. Keep the runtime checks
below. Imported/hosted tile layers, video and independent bucket recovery retain their separate limitations.

## Checks after a sync

```bash
kubectl --context do-nyc3-omos-cloud -n cloudtak get pods
kubectl --context do-nyc3-omos-cloud -n cloudtak logs deploy/cloudtak-api | tail -30
```

Look for the API reporting its Tile38 geofence connection, no `ECONNREFUSED` to Postgres, and no S3 `AccessDenied`. Confirm all four pods stay Ready without restarts, the two private HTTPS origins work,
and a TAK user can log in, see CoT/chat/mission data, and reconnect after an API restart. Check both database and
bucket writes before declaring acceptance. Imported/hosted tile layers and media keep their gates above.
