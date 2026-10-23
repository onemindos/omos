# CloudTAK

A fresh install of upstream [CloudTAK](https://github.com/dfpc-coe/CloudTAK) (`dfpc-coe/CloudTAK`, pinned
**v13.113.0**, released 2026-10-23) as a browser TAK client. Deployed by the Argo CD Application
`platform-cloudtak` (sync wave 5) into the namespace `cloudtak`.

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
| PostgreSQL + PostGIS | `onemind-data` cluster, namespace `onemind-data` | Database and role `cloudtak` (declared in `platform/ecosystem/resources/`), `POSTGRES` from `Secret/cloudtak-db` |
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
retention to the API.

## Before the first sync: the Spaces key (Zeus)

The pods read the Spaces key from `Secret/cloudtak-spaces`. It is not in Git: it must be created from a key only Zeus
can issue. Do this **on the PR branch before merging**, so everything starts cleanly. If the secret is missing at
merge, the pods wait in `CreateContainerConfigError` and start by themselves once it appears.

1. DigitalOcean control panel, **Spaces Object Storage**, **Create bucket**:
   region **NYC3**, name `onemind-cloudtak`, file listing **Restricted**, CDN off. (If you pick another name, change
   `ASSET_BUCKET` in `config.yaml`.)
2. **Spaces Object Storage**, **Access Keys**, **Create Access Key**: choose **Limited Access**, bucket
   `onemind-cloudtak`, permission **Read/Write/Delete**. Copy the key ID and the secret now; the secret is shown once.
3. Seal it with the cluster's public key (the command prompts for the values, so they never enter shell history):

   ```bash
   kubeseal --context do-nyc3-omos-cloud --fetch-cert \
     --controller-name=sealed-secrets-controller --controller-namespace=kube-system > /tmp/pub-cert.pem
   read -r -p 'Spaces access key ID: ' SPACES_KEY
   read -r -s -p 'Spaces secret key: ' SPACES_SECRET; echo
   kubectl create secret generic cloudtak-spaces -n cloudtak \
     --from-literal=access-key="$SPACES_KEY" --from-literal=secret-key="$SPACES_SECRET" \
     --dry-run=client -o json \
     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/cloudtak/cloudtak-spaces.sealedsecret.yaml
   unset SPACES_KEY SPACES_SECRET
   ```

4. Commit `platform/cloudtak/cloudtak-spaces.sealedsecret.yaml` to the PR branch. It is encrypted and safe to commit.
   Never commit the plain key.

The other two secrets (`cloudtak-signing`, the JWT signing secret shared by all four workloads, and `cloudtak-db`,
plus the role password in `onemind-data/db-cloudtak`) were generated at random and sealed in this PR; nobody has
seen the values. To rotate the signing secret, reseal `cloudtak-signing` and restart all four Deployments; every
user session and API token is signed with it, so everyone has to log in again.

## First run (Zeus, in the CloudTAK UI)

After the Application syncs and `cloudtak-api` is Ready (the first start applies about 230 database migrations and
loads the default iconsets and basemaps; allow a few minutes), open **https://cloudtak.tail717ba0.ts.net** from a
device on the tailnet. CloudTAK shows **Initial Server Configuration**. Until it is submitted, anyone who can reach
the page can configure it, so do it straight away.

| Field | Value |
| --- | --- |
| Server Name | any label, for example `OneMind TAK` |
| Server Hostname | `tak-tailnet.takserver.svc.cluster.local` |
| CoT Port | `8089` |
| Marti Port | `8443` |
| WebTAK | `8446` |
| Admin Certificate | the TAK Server admin client certificate (`.p12`) and its password |
| Initial Administrator Username / Password | an existing TAK Server user (LDAP-backed users work); it becomes the CloudTAK system administrator |

The hostname is the in-cluster name of the `tak-tailnet` Service, which forwards all three ports through the
`tak-lb-proxy` pods. The form takes one hostname for all three ports; the network policy also allows
`takserver-api` and `takserver-messaging` directly, so after the first run you can split them in **Admin, Server**
(`ssl://takserver-messaging-service.takserver.svc.cluster.local:8089`,
`https://takserver-api-service.takserver.svc.cluster.local:8443`,
`https://takserver-api-service.takserver.svc.cluster.local:8446`). CloudTAK does not verify the TAK Server
certificate chain, so the self-signed certificates are accepted.

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

## Media (open decision)

CloudTAK's video features do not work with plain MediaMTX. Its setting `media::url` must point at the `media-infra`
wrapper (`ghcr.io/dfpc-coe/media-infra`), which exposes an API on port **9997** with routes CloudTAK calls
(`/path`, `/path/:name`, `/v3/config/global/get`, `/stream/<name>/index.m3u8`) and which authenticates MediaMTX
clients by calling back to CloudTAK (`/api/video/auth`). The shared `onetak-media` is plain `bluenviron/mediamtx`
with no auth, no 9997 and none of those routes (it replaced media-infra on 2026-09-29), so `media::url` is left unset
and CloudTAK's video lease feature is simply off. Everything else works.

Options, if video leases are wanted:

1. **Swap `onetak-media` to media-infra** (`ghcr.io/dfpc-coe/media-infra`, pin by digest). One server for everyone, but
   media-infra turns on MediaMTX HTTP auth, so every RTSP/RTMP/SRT publisher and reader, including the takml vision
   worker, would have to authenticate through CloudTAK's lease system. That is the change the cluster moved away from.
   `platform/onetak/` is being edited by another change set; do not do this in this PR.
2. **Run media-infra for CloudTAK only** in this namespace, leaving `onetak-media` untouched. No impact on existing
   feeds, but it is a second media server, which this install deliberately avoids, and it needs the same shared
   `SigningSecret` and public ports for RTSP/WebRTC.
3. **Do nothing** (current). The rest of CloudTAK works; only its own video lease feature is off.

Whichever is chosen, set `CLOUDTAK_Config_media_url` on `cloudtak-api` and keep the 9997 egress rule in `network.yaml`.

## The AWS-only "layers" feature

CloudTAK's ETL "Layers" create AWS Lambda functions and CloudFormation stacks and use SQS. Nothing here has AWS
credentials, and the pods cannot reach AWS. The layer plumbing is switched off or inert:

- `cloudtak-api` starts with `--no-events --no-etl-events` (no Lambda cron invocations, no SQS delivery).
- `AWS_EC2_METADATA_DISABLED=true`, no AWS credentials in any pod, and the network policy blocks the metadata address.
- Creating or running a layer in the UI would fail. Do not use **Layers**; the Layers page may show errors. There is
  no setting that hides the feature.

## Known limitations

- **Imported-data tiles may not load.** The API fetches TileJSON for user-imported assets from `PMTILES_URL` using
  upstream's SSRF-safe fetch, which refuses private and tailnet addresses, and pods cannot reach the tailnet at all.
  The import worker only needs the API and Spaces; the map layer for an imported asset is the part at risk.
  Hosted basemaps, the map, chat, missions and the TAK connection are unaffected. A one-line upstream fix would allow
  the tiles host; until then, treat imports as untested. Check this after the first import.
- Single API replica. Upstream's split mode (`api` + `hub`) exists for scaling, but is not needed here.
- Hardening of the three non-root pods (read-only root filesystem, uid 1000) is untested against the live images.
  If one fails to start, relax `securityContext` in that Deployment first.

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
4. Open a PR. After merge, Argo restarts `cloudtak-api` (Recreate, so a short outage); it runs the migrations on start, so
   watch `kubectl -n cloudtak logs deploy/cloudtak-api` until `/api` answers. Roll back by reverting the PR, but
   note that a migration that already ran is not undone; take a Postgres backup first when the notes mention one.

## Checks after a sync

```bash
kubectl --context do-nyc3-omos-cloud -n cloudtak get pods
kubectl --context do-nyc3-omos-cloud -n cloudtak logs deploy/cloudtak-api | tail -30
```

Look for the API reporting its Tile38 geofence connection, no `ECONNREFUSED` to Postgres, and no S3 `AccessDenied`.
