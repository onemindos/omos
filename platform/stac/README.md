# STAC

A fresh STAC (SpatioTemporal Asset Catalog) for OneMind imagery, on the shared `onemind-data` Postgres, filled
automatically from WebODM. Namespace `stac`, Argo CD application `platform-stac` (sync wave 5). Deploy it together
with CloudTAK (PR #390): same wave, same shared Postgres.

It replaces the old `onetak/stac-api` and the TimescaleDB `onetak-geo-db`, which PR #406 removes. Those never worked
(pgstac was never installed in that database). Nothing here depends on them.

## What STAC is for here

STAC is a standard way to **search and catalog** geospatial imagery: "which orthophotos or elevation models cover this
area, taken in this time window", answered as JSON that any tool understands. It does not store or serve pixels; each
Item points at the files (here, Cloud Optimized GeoTIFFs in DigitalOcean Spaces).

- **Hermes** (and other agents) query it to find imagery for a place and time, in-cluster at
  `http://stac-api.stac.svc.cluster.local:8080` (for example `GET /search?bbox=-80.5,27.2,-80.4,27.3&datetime=2026-10-01T00:00:00Z/..`).
- **CloudTAK / TAK** use it to discover which imagery exists and where its files are. Getting that imagery onto a TAK
  map is a separate step, see "How TAK gets the imagery".
- **People** browse it in STAC Browser on the tailnet.

## Components

| Piece | What | Reachable at |
|---|---|---|
| Database | Database `stac` and role `stac` on CloudNativePG cluster `onemind-data` (PostGIS, `btree_gist`, `unaccent`), defined in `platform/ecosystem/resources/` | `onemind-data-rw.onemind-data` (NetworkPolicy allows only the stac pods) |
| `stac-pgstac-migrate` | Argo CD Sync hook Job (`migrate.yaml`): `pypgstac migrate` with `ghcr.io/stac-utils/pgstac-pypgstac` v0.9.12. Installs the pgstac schema, or upgrades it | none |
| `stac-api` | stac-fastapi-pgstac 7.1.0, **read and search only** | `https://stac.tail717ba0.ts.net` (tailnet), `stac-api.stac:8080` (cluster, Hermes) |
| `stac-api-writer` | the same image with the transaction extension on | `stac-api-writer.stac:8080`, only from the sync CronJob |
| `stac-browser` | official `radiantearth/stac-browser` 5.1.0, static files, configured with the tailnet API URL | `https://stac-browser.tail717ba0.ts.net` (tailnet) |
| `webodm-stac-sync` | CronJob every 10 minutes (`webodm-sync.yaml`), code in `services/stac-webodm-sync/` | none |

All images are pinned by tag and digest. No public ingress, no Funnel.

### Read and write are separate Deployments

stac-fastapi-pgstac cannot split reads from writes inside one process: the transaction extension is one switch
(`ENABLE_TRANSACTIONS_EXTENSIONS`) that adds `POST/PUT/PATCH/DELETE` on collections and items, and since 7.0 the app
has a single database connection setting. It also has no authentication. So the split is done with two Deployments of
the same image:

- `stac-api` has the extension **off**. It has no write routes at all (verified: `POST/PUT/DELETE` on `/collections` and
  `POST` on `/collections/x/items` return 405, bulk items 404). This is the one on the tailnet.
- `stac-api-writer` has it **on**, has no Ingress, and its NetworkPolicy admits only pods labelled
  `app.kubernetes.io/component: webodm-sync` in this namespace. The only writer is the sync CronJob.

Both use the primary (`onemind-data-rw`), because pgstac records search statistics with an INSERT. Pointing the read API at
a replica (`onemind-data-ro`) would add database-level read-only, but needs pgstac's `readonly` setting for those
connections; that is possible follow-up hardening, not needed for safety today.

## The WebODM flow

How WebODM stores results (from the WebODM source and the live `onetak` objects): each task writes to the
`webodm-media` volume under `/webodm/app/media/project/<project id>/task/<task id>/assets/`: `odm_orthophoto/odm_orthophoto.tif`,
`odm_dem/dsm.tif`, `odm_dem/dtm.tif`. When a task completes, WebODM itself makes those three files Cloud Optimized
GeoTIFFs (`assure_cogs`). The REST API (`/api/token-auth/` for a JWT, `/api/projects/`, `/api/projects/<p>/tasks/?status=40`,
`/api/projects/<p>/tasks/<t>/download/orthophoto.tif`) lists completed tasks (status 40) with their `available_assets` and
serves the files.

Every 10 minutes `webodm-stac-sync` does this (`services/stac-webodm-sync/sync.py`):

1. Log in to WebODM and make sure the Collection `webodm` exists (created, or its description fields corrected).
2. List completed tasks that have an orthophoto, DSM or DTM. A task is skipped when its STAC Item already exists with the same
   `webodm:fingerprint` (task id, processing time, size, asset list). A task WebODM reprocesses gets a new fingerprint and is redone.
3. For each new or changed task, oldest first, at most 3 per run: download each raster, check it is a COG (convert with GDAL if
   not; normally it is), upload it to Spaces at `webodm/<project id>/<task id>/{orthophoto,dsm,dtm}.tif`, read its footprint.
4. Create or update the Item `webodm/<task id>`:
   - `geometry` and `bbox`: the raster footprint reprojected to WGS 84 (8 points per edge so a rotated UTM edge stays true), bbox
     the union of the task's rasters;
   - `properties.datetime`: the task's `created_at`, in UTC; also `title`, WebODM project and task ids, image count, processing time;
   - assets `orthophoto` (roles `data`, `visual`), `dsm`, `dtm` (role `data`): `href` is the HTTPS URL
     `https://onemind-imagery.nyc3.digitaloceanspaces.com/webodm/...`, media type
     `image/tiff; application=geotiff; profile=cloud-optimized`, `file:size`, `file:checksum`, `proj:code/shape/transform`, `gsd`, and an
     `alternate.s3` link (`s3://onemind-imagery/...`) for tools that read Spaces over the S3 protocol.
5. pgstac keeps the Collection's spatial and temporal extent current (the migration hook turns on `update_collection_extent`).

A failed task is logged and retried next run without blocking the others (the job then exits 1). If WebODM is unreachable
(it is scaled to zero between flights in `platform/onetak/webodm.yaml`) the job logs a warning and exits 0. Until the two
Secrets below exist it logs what is missing and exits 0.

### Why a polling CronJob

| Option | Verdict |
|---|---|
| (c) WebODM webhook | WebODM has no outbound webhook for task completion. The source has none; the only hook is an internal plugin signal (`task_completed`), which would need a custom plugin baked into a WebODM image that is currently `:latest` and unpinned. |
| (b) n8n workflow | n8n is approved for automation, but this job moves multi-GB binaries, validates and converts COGs with GDAL, and signs S3 uploads. n8n here blocks `executeCommand`, has no GDAL, and keeps workflows and credentials in its own database, not Git. Good for what happens **after** an Item exists (notify, route), not for the ETL. |
| **(a) CronJob poll** | Chosen. The only state is the Items themselves, so it is idempotent by task id, recovers from missed runs and from WebODM or STAC being down, backfills existing tasks on the first run, and is a plain Kubernetes Job. Latency is up to 10 minutes. |

The job runs on the upstream `ghcr.io/osgeo/gdal` image (pinned; it has Python, the GDAL bindings and the COG driver) with
`sync.py` mounted from `webodm-sync-code.configmap.yaml`, so there is no image of our own to build and pin. That ConfigMap is
generated: after editing `sync.py` run `python services/stac-webodm-sync/render_configmap.py` and commit both; a test fails if they differ.
The S3 upload uses GDAL's `/vsis3/` handler pointed at the Spaces endpoint. The `AWS_*` variable names are GDAL's names for
S3-protocol settings, not AWS: no AWS service or tool is involved.

## Before the first sync (Zeus)

The API, browser and database come up on their own. The sync needs two secrets that only you can create. Do this on the PR
branch before merging (or any time after; until then the sync does nothing and does not fail).

### 1. Spaces bucket and key

1. DigitalOcean control panel, **Spaces Object Storage**, **Create bucket**: region **NYC3**, name `onemind-imagery`, file
   listing **Restricted**, CDN off. (Another name: change `SPACES_BUCKET` in `webodm-sync.yaml`; the Item hrefs follow it.)
2. **Access Keys**, **Create Access Key**, **Limited Access**, bucket `onemind-imagery`, permission **Read/Write**. Copy the key ID and the
   secret now; the secret is shown once. (The sync never deletes. If a very large upload fails on abort, try Read/Write/Delete.)
3. Seal it (the prompts keep the values out of shell history):

   ```bash
   kubeseal --context do-nyc3-omos-cloud --fetch-cert \
     --controller-name=sealed-secrets-controller --controller-namespace=kube-system > /tmp/pub-cert.pem
   read -r -p 'Spaces access key ID: ' SPACES_KEY
   read -r -s -p 'Spaces secret key: ' SPACES_SECRET; echo
   kubectl create secret generic stac-spaces -n stac \
     --from-literal=access-key="$SPACES_KEY" --from-literal=secret-key="$SPACES_SECRET" \
     --dry-run=client -o json \
     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/stac/stac-spaces.sealedsecret.yaml
   unset SPACES_KEY SPACES_SECRET
   ```

### 2. A WebODM user for the sync

1. In WebODM (Administration, Users) create a user, for example `stac-sync`, not staff. For each project to publish, use the
   project's **Edit Permissions** and give that user **View**. (Using an administrator account also works but gives the cluster
   a password that can change and delete everything in WebODM.)
2. Seal it:

   ```bash
   read -r -p 'WebODM sync username: ' WO_USER
   read -r -s -p 'WebODM sync password: ' WO_PASS; echo
   kubectl create secret generic webodm-api -n stac \
     --from-literal=username="$WO_USER" --from-literal=password="$WO_PASS" \
     --dry-run=client -o json \
     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/stac/webodm-api.sealedsecret.yaml
   unset WO_USER WO_PASS
   ```

3. Commit `platform/stac/stac-spaces.sealedsecret.yaml` and `platform/stac/webodm-api.sealedsecret.yaml`. They are encrypted and safe to
   commit. `*.secret.example` files in this directory are templates showing the shape; Argo CD does not apply them.

The database password (`stac-db`, and the role secret in `onemind-data`) was generated at random and sealed in this PR; nobody has seen it.

## Open decisions

- **Bucket privacy.** The bucket is private (file listing restricted, objects not public). The Item hrefs are HTTPS URLs, so a plain
  browser cannot open them and STAC Browser shows no previews. Clients with the Spaces key (or a future tile service) read them over S3.
  The alternative is public-read objects: simpler, but anyone with a URL sees drone imagery of the homestead and operations.
- **Bucket name and region.** `onemind-imagery` in NYC3 is a proposal.
- **Item `datetime`** is the WebODM task's creation time, as asked. The flight time is in the image EXIF and is not in the WebODM API.
- **NATS events** are not published, see below.

## How TAK gets the imagery (follow-up, not built)

STAC says where the COGs are; TAK needs pixels as map tiles or files. Options, in the order I would try them:

1. **A TiTiler** (tailnet only, in this namespace or its own) that reads the COGs from Spaces with a read-only key and serves XYZ/WMTS
   tiles for any Item. CloudTAK and ATAK can add a tile URL as a map layer. `titiler-pgstac` can even serve a whole search as one layer.
2. **PMTiles**: an extra sync step renders each orthophoto to a raster PMTiles file in a bucket, and CloudTAK's PMTiles server
   (`cloudtak-tiles`) serves it. Best for offline and field use; costs conversion time and storage.
3. **A TAK data package** (zip with a GeoTIFF or MBTiles and a manifest) pushed to TAK Server, for devices that must work with no network.

## NATS (follow-up)

The architecture rule is that data flows on NATS where it makes sense, and "a new Item exists" is a fact worth publishing. It is not done
here because the NATS users are defined in the Helm values secret that is applied by hand (`platform/nats/OPERATIONS.md`), so a
publish-only user for this job cannot be added cleanly through Git. Also, the subject suggested for this (`onemind.imagery.stac.item.created`)
is outside the locked grammar (`docs/NATS-DOCTRINE.md`). Within it the event is `evt.drone.imagery.item_created` (or `evt.sat.imagery.item_created` for
satellite scenes), which the existing EVT stream already captures, so no stream change is needed. The follow-up: a publish-only NATS
user limited to `evt.drone.imagery.>`, sealed as `stac-nats`, and a best-effort publish after `put_item` in `sync.py`. An n8n workflow can then
subscribe and notify.

## Verification (after merge and sync)

```bash
K="kubectl --context do-nyc3-omos-cloud -n stac"
$K get pods                                    # stac-api, stac-api-writer, stac-browser Running; migrate Job Completed
$K logs job/stac-pgstac-migrate | tail -3      # "pgstac schema version: 0.9.12"
```

From a device on the tailnet:

```bash
curl -s https://stac.tail717ba0.ts.net/ | head -c 300                      # landing page; links must start with https://stac.tail717ba0.ts.net
curl -s https://stac.tail717ba0.ts.net/collections                         # empty until the first sync, then "webodm"
curl -s -X POST https://stac.tail717ba0.ts.net/search -H 'content-type: application/json' -d '{"limit":5}'
curl -s -o /dev/null -w '%{http_code}\n' -X POST https://stac.tail717ba0.ts.net/collections   # 405: no write routes
```

After a test WebODM task (a small image set is enough) finishes, wait up to 10 minutes or run the job now:

```bash
$K create job --from=cronjob/webodm-stac-sync sync-now
$K logs -f job/sync-now                        # "task <id>: created"
curl -s https://stac.tail717ba0.ts.net/collections/webodm/items | head -c 600
```

Open `https://stac-browser.tail717ba0.ts.net`: the `webodm` collection and the new Item (footprint on the map, three assets) appear.
If the landing-page links show `http://` or a cluster hostname, the Tailscale proxy is not passing `X-Forwarded-Proto`/`Host`, which
stac-fastapi uses to build links; that needs a fix before the browser can follow them.

## Upgrading

- **stac-fastapi-pgstac**: raise the tag and digest in `api.yaml` (both Deployments). Check its changelog for renamed settings (7.0 removed the
  `POSTGRES_*` variables, which is why the old Deployment broke on a newer image) and the supported pgstac range.
- **pgstac**: raise the `pgstac-pypgstac` tag and digest in `migrate.yaml`. Migrations are forward-only: take a Postgres backup first
  (the cluster is backed up by Barman, `platform/ecosystem/resources/postgres-backup-schedule.yaml`). The hook applies them on the next sync.
- **The sync**: edit `sync.py`, run `render_configmap.py`, run `python -m unittest discover -s services/stac-webodm-sync`, commit both files.

## Notes

- The `stac` role has `CREATEROLE` because pgstac creates `pgstac_admin`, `pgstac_ingest` and `pgstac_read` on first install. Those roles are
  cluster-wide and stay if the database is ever dropped.
- `webodm-media` is a 5 GiB volume, small for orthophotos from large flights; if WebODM tasks start failing, check it first.
- Local test used: PostGIS 16, the pinned pypgstac and stac-fastapi images, a fake WebODM serving real GeoTIFFs, and a local S3 server for Spaces:
  fresh install, upgrade from pgstac 0.9.10, and idempotent re-runs all pass. It has not been run against the live cluster, WebODM or Spaces.
