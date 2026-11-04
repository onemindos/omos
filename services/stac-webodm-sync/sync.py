#!/usr/bin/env python3
"""Publish completed WebODM tasks to the STAC API.

One run of this script (a Kubernetes CronJob, platform/stac/webodm-sync.yaml):

1. logs in to the WebODM REST API and lists every completed task that has an orthophoto, DSM or DTM;
2. for each task that is new or changed since the last run (compared by a fingerprint stored on the STAC Item):
   downloads each raster, makes sure it is a Cloud Optimized GeoTIFF (converting it if not), uploads it to
   a DigitalOcean Spaces bucket, and reads its footprint;
3. creates or updates the STAC Item `webodm/<task id>` through the in-cluster, write-enabled STAC API.

Idempotent: the Item id is the WebODM task UUID and the object keys are derived from it, so a re-run
overwrites the same objects and the same Item. Only the Python standard library is imported at module level;
GDAL (`osgeo`) is imported lazily by `GdalRaster` and `SpacesStorage`, and the image that runs this
(ghcr.io/osgeo/gdal) provides it. That keeps the logic below unit-testable on any machine with plain Python.

Configuration is by environment variable, see `Config.from_env`.
"""
from __future__ import annotations

import datetime as dt
import hashlib
import json
import logging
import os
import re
import shutil
import sys
import tempfile
import time
import urllib.error
import urllib.parse
import urllib.request
from dataclasses import dataclass, field
from typing import Any, Callable, Dict, Iterable, List, Optional, Tuple

LOG = logging.getLogger("stac-webodm-sync")

COLLECTION_ID = "webodm"
STAC_VERSION = "1.1.0"
EXT_PROJECTION = "https://stac-extensions.github.io/projection/v2.0.0/schema.json"
EXT_FILE = "https://stac-extensions.github.io/file/v2.1.0/schema.json"
EXT_ALTERNATE = "https://stac-extensions.github.io/alternate-assets/v1.2.0/schema.json"
COG_MEDIA_TYPE = "image/tiff; application=geotiff; profile=cloud-optimized"

# WebODM status code for a finished task (nodeodm/status_codes.py: QUEUED 10, RUNNING 20, FAILED 30,
# COMPLETED 40, CANCELED 50).
WEBODM_COMPLETED = 40

# WebODM asset file name (Task.ASSETS_MAP key) -> STAC asset key, roles and title.
ASSET_SPECS: Tuple[Tuple[str, str, List[str], str], ...] = (
    ("orthophoto.tif", "orthophoto", ["data", "visual"], "Orthophoto (Cloud Optimized GeoTIFF)"),
    ("dsm.tif", "dsm", ["data"], "Digital surface model (Cloud Optimized GeoTIFF)"),
    ("dtm.tif", "dtm", ["data"], "Digital terrain model (Cloud Optimized GeoTIFF)"),
)

COLLECTION_FIELDS = ("title", "description", "license", "keywords", "providers")


class SyncError(Exception):
    """A task or the whole run could not be synchronised."""


class WebOdmUnavailable(SyncError):
    """WebODM does not answer (for example scaled to zero between flights)."""


# --------------------------------------------------------------------------------------------------
# Configuration
# --------------------------------------------------------------------------------------------------


@dataclass
class Config:
    webodm_url: str
    webodm_username: str
    webodm_password: str = field(repr=False)
    stac_url: str
    bucket: str
    endpoint: str
    region: str
    prefix: str
    public_base_url: str
    work_dir: str
    max_tasks: int
    force: bool
    dry_run: bool
    virtual_hosting: bool = True   # SPACES_VIRTUAL_HOSTING / SPACES_HTTPS exist so tests can use a local S3 server
    https: bool = True
    access_key: str = field(default="", repr=False)
    secret_key: str = field(default="", repr=False)

    @classmethod
    def from_env(cls, env: Optional[Dict[str, str]] = None) -> "Config":
        e = os.environ if env is None else env
        bucket = e.get("SPACES_BUCKET", "onemind-imagery")
        region = e.get("SPACES_REGION", "nyc3")
        endpoint = e.get("SPACES_ENDPOINT", "%s.digitaloceanspaces.com" % region)
        return cls(
            webodm_url=e.get("WEBODM_URL", "http://webodm.onetak.svc.cluster.local:8000").rstrip("/"),
            webodm_username=e.get("WEBODM_USERNAME", ""),
            webodm_password=e.get("WEBODM_PASSWORD", ""),
            stac_url=e.get("STAC_URL", "http://stac-api-writer.stac.svc.cluster.local:8080").rstrip("/"),
            bucket=bucket,
            endpoint=endpoint,
            region=region,
            prefix=e.get("SPACES_PREFIX", "webodm").strip("/"),
            public_base_url=e.get(
                "SPACES_PUBLIC_BASE_URL", "https://%s.%s" % (bucket, endpoint)
            ).rstrip("/"),
            work_dir=e.get("WORK_DIR", "/work"),
            max_tasks=int(e.get("MAX_TASKS_PER_RUN", "3")),
            force=e.get("FORCE_RESYNC", "false").lower() == "true",
            dry_run=e.get("DRY_RUN", "false").lower() == "true",
            virtual_hosting=e.get("SPACES_VIRTUAL_HOSTING", "true").lower() == "true",
            https=e.get("SPACES_HTTPS", "true").lower() == "true",
            access_key=e.get("AWS_ACCESS_KEY_ID", ""),
            secret_key=e.get("AWS_SECRET_ACCESS_KEY", ""),
        )

    def missing(self) -> List[str]:
        """Names of required settings that are empty (credentials the operator still has to supply)."""
        out = []
        if not self.webodm_username or not self.webodm_password:
            out.append("WEBODM_USERNAME/WEBODM_PASSWORD (Secret webodm-api)")
        if not self.access_key or not self.secret_key:
            out.append("Spaces access key (Secret stac-spaces)")
        return out


# --------------------------------------------------------------------------------------------------
# HTTP
# --------------------------------------------------------------------------------------------------

# An HTTP function: (method, url, headers, body) -> (status, response headers, body bytes).
# Implemented by `urllib_http`; replaced by a fake in the tests.
HttpFn = Callable[[str, str, Dict[str, str], Optional[bytes]], Tuple[int, Dict[str, str], bytes]]


def urllib_http(method: str, url: str, headers: Dict[str, str], body: Optional[bytes] = None,
                timeout: float = 60.0) -> Tuple[int, Dict[str, str], bytes]:
    req = urllib.request.Request(url, data=body, method=method, headers=headers)
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:  # nosec B310 - cluster-internal URLs
            return resp.status, dict(resp.headers.items()), resp.read()
    except urllib.error.HTTPError as err:
        return err.code, dict(err.headers.items()), err.read()


def _retrying(http: HttpFn, retries: int = 2, pause: float = 2.0) -> HttpFn:
    """Retry connection errors and 502/503/504 a couple of times (all calls here are idempotent)."""

    def call(method, url, headers, body=None):
        last: Optional[BaseException] = None
        for attempt in range(retries + 1):
            try:
                status, hdrs, data = http(method, url, headers, body)
            except (urllib.error.URLError, ConnectionError, TimeoutError, OSError) as err:
                last = err
            else:
                if status not in (502, 503, 504) or attempt == retries:
                    return status, hdrs, data
            if attempt < retries:
                time.sleep(pause * (attempt + 1))
        raise last if last else SyncError("request failed: %s %s" % (method, url))

    return call


def _json(status: int, data: bytes, what: str, ok: Iterable[int] = (200,)) -> Any:
    if status not in tuple(ok):
        raise SyncError("%s: HTTP %d %s" % (what, status, data[:300].decode("utf-8", "replace")))
    return json.loads(data.decode("utf-8")) if data else None


# --------------------------------------------------------------------------------------------------
# WebODM
# --------------------------------------------------------------------------------------------------


class WebOdmClient:
    """The few WebODM REST calls needed. Reference: app/api/urls.py and app/api/tasks.py in OpenDroneMap/WebODM."""

    def __init__(self, base_url: str, username: str, password: str, http: HttpFn = urllib_http):
        self.base = base_url
        self.username = username
        self.password = password
        self.http = _retrying(http)
        self.token: Optional[str] = None

    def _headers(self) -> Dict[str, str]:
        return {"Authorization": "JWT %s" % self.token, "Accept": "application/json"}

    def login(self) -> None:
        body = json.dumps({"username": self.username, "password": self.password}).encode()
        try:
            status, _, data = self.http("POST", self.base + "/api/token-auth/",
                                        {"Content-Type": "application/json"}, body)
        except (urllib.error.URLError, ConnectionError, TimeoutError, OSError) as err:
            raise WebOdmUnavailable("WebODM unreachable at %s: %s" % (self.base, err))
        if status in (502, 503, 504):
            raise WebOdmUnavailable("WebODM answered HTTP %d" % status)
        self.token = _json(status, data, "WebODM login")["token"]

    def _get(self, path: str) -> Any:
        status, _, data = self.http("GET", self.base + path, self._headers(), None)
        if status == 401:  # the 6 hour JWT expired during a long run
            self.login()
            status, _, data = self.http("GET", self.base + path, self._headers(), None)
        return _json(status, data, "WebODM GET " + path)

    def projects(self) -> List[Dict[str, Any]]:
        out: List[Dict[str, Any]] = []
        page = 1
        while True:
            body = self._get("/api/projects/?page=%d" % page)
            if isinstance(body, list):  # unpaginated deployments
                return body
            out.extend(body.get("results", []))
            if not body.get("next"):
                return out
            page += 1

    def completed_tasks(self, project_id: int) -> List[Dict[str, Any]]:
        body = self._get("/api/projects/%s/tasks/?status=%d" % (project_id, WEBODM_COMPLETED))
        return body["results"] if isinstance(body, dict) else body

    def download(self, project_id: int, task_id: str, asset: str, dest: str) -> None:
        """Stream /api/projects/<p>/tasks/<t>/download/<asset> to a file."""
        url = "%s/api/projects/%s/tasks/%s/download/%s" % (self.base, project_id, task_id, asset)
        req = urllib.request.Request(url, headers=self._headers())
        with urllib.request.urlopen(req, timeout=300) as resp, open(dest, "wb") as out:  # nosec B310
            shutil.copyfileobj(resp, out, 1024 * 1024)


# --------------------------------------------------------------------------------------------------
# STAC API (the write-enabled stac-fastapi-pgstac instance)
# --------------------------------------------------------------------------------------------------


class StacClient:
    def __init__(self, base_url: str, http: HttpFn = urllib_http):
        self.base = base_url
        self.http = _retrying(http)

    def _call(self, method: str, path: str, body: Optional[dict] = None) -> Tuple[int, Any]:
        headers = {"Accept": "application/json"}
        payload = None
        if body is not None:
            headers["Content-Type"] = "application/json"
            payload = json.dumps(body).encode()
        status, _, data = self.http(method, self.base + path, headers, payload)
        parsed = None
        if data:
            try:
                parsed = json.loads(data.decode("utf-8"))
            except ValueError:
                parsed = data[:300].decode("utf-8", "replace")
        return status, parsed

    def get_collection(self, collection_id: str) -> Optional[dict]:
        status, body = self._call("GET", "/collections/" + collection_id)
        if status == 404:
            return None
        if status != 200:
            raise SyncError("STAC GET collection: HTTP %d %s" % (status, body))
        return body

    def put_collection(self, collection: dict, exists: bool) -> None:
        if exists:
            status, body = self._call("PUT", "/collections/" + collection["id"], collection)
        else:
            status, body = self._call("POST", "/collections", collection)
        if status not in (200, 201):
            raise SyncError("STAC write collection: HTTP %d %s" % (status, body))

    def get_item(self, collection_id: str, item_id: str) -> Optional[dict]:
        status, body = self._call("GET", "/collections/%s/items/%s" % (collection_id, item_id))
        if status == 404:
            return None
        if status != 200:
            raise SyncError("STAC GET item: HTTP %d %s" % (status, body))
        return body

    def put_item(self, collection_id: str, item: dict, exists: bool) -> None:
        if exists:
            status, body = self._call("PUT", "/collections/%s/items/%s" % (collection_id, item["id"]), item)
        else:
            status, body = self._call("POST", "/collections/%s/items" % collection_id, item)
            if status == 409:  # lost a race with another writer: fall through to update
                status, body = self._call("PUT", "/collections/%s/items/%s" % (collection_id, item["id"]), item)
        if status not in (200, 201):
            raise SyncError("STAC write item: HTTP %d %s" % (status, body))


# --------------------------------------------------------------------------------------------------
# Raster and object-storage adapters (GDAL)
# --------------------------------------------------------------------------------------------------


@dataclass
class RasterInfo:
    """What the STAC Item needs to know about one GeoTIFF."""
    footprint: List[List[float]]          # closed ring of [lon, lat] in EPSG:4326
    bbox: List[float]                     # [west, south, east, north]
    code: Optional[str]                   # for example "EPSG:32617"
    shape: List[int]                      # [rows, columns]
    transform: List[float]                # affine a, b, c, d, e, f
    gsd: Optional[float]                  # metres per pixel when the CRS is in metres


class GdalRaster:
    """COG validation, conversion and footprint extraction with the GDAL Python bindings."""

    def __init__(self):
        from osgeo import gdal, osr  # imported lazily so the module loads without GDAL
        gdal.UseExceptions()
        osr.UseExceptions()
        self.gdal = gdal
        self.osr = osr

    def is_cog(self, path: str) -> bool:
        from osgeo_utils.samples.validate_cloud_optimized_geotiff import validate
        _warnings, errors, _details = validate(path, full_check=True)
        return not errors

    def to_cog(self, src: str, dst: str) -> None:
        self.gdal.Translate(dst, src, format="COG", creationOptions=[
            "COMPRESS=DEFLATE", "PREDICTOR=YES", "OVERVIEWS=AUTO", "BIGTIFF=IF_SAFER", "NUM_THREADS=ALL_CPUS",
        ])

    def describe(self, path: str) -> RasterInfo:
        ds = self.gdal.Open(path)
        gt = ds.GetGeoTransform()
        srs_wkt = ds.GetProjection()
        if not srs_wkt:
            raise SyncError("%s has no coordinate reference system" % os.path.basename(path))
        src = self.osr.SpatialReference()
        src.ImportFromWkt(srs_wkt)
        src.SetAxisMappingStrategy(self.osr.OAMS_TRADITIONAL_GIS_ORDER)
        dst = self.osr.SpatialReference()
        dst.ImportFromEPSG(4326)
        dst.SetAxisMappingStrategy(self.osr.OAMS_TRADITIONAL_GIS_ORDER)
        ct = self.osr.CoordinateTransformation(src, dst)
        cols, rows = ds.RasterXSize, ds.RasterYSize
        ring = footprint_ring(gt, cols, rows, lambda x, y: ct.TransformPoint(x, y)[:2])
        code = None
        if src.GetAuthorityName(None) == "EPSG" and src.GetAuthorityCode(None):
            code = "EPSG:%s" % src.GetAuthorityCode(None)
        metres = src.IsProjected() and abs(src.GetLinearUnits() - 1.0) < 1e-9
        return RasterInfo(
            footprint=ring,
            bbox=ring_bbox(ring),
            code=code,
            shape=[rows, cols],
            transform=[gt[1], gt[2], gt[0], gt[4], gt[5], gt[3]],
            gsd=round(abs(gt[1]), 4) if metres else None,
        )


class SpacesStorage:
    """Upload to DigitalOcean Spaces over the S3 protocol using GDAL's /vsis3/ handler.

    The AWS_* names are GDAL's configuration keys for any S3-compatible service; the endpoint is
    Spaces. Nothing here talks to AWS.
    """

    def __init__(self, cfg: Config):
        from osgeo import gdal
        self.gdal = gdal
        self.bucket = cfg.bucket
        gdal.SetConfigOption("AWS_S3_ENDPOINT", cfg.endpoint)
        gdal.SetConfigOption("AWS_VIRTUAL_HOSTING", "TRUE" if cfg.virtual_hosting else "FALSE")
        gdal.SetConfigOption("AWS_REGION", cfg.region)
        gdal.SetConfigOption("AWS_HTTPS", "YES" if cfg.https else "NO")
        gdal.SetConfigOption("AWS_ACCESS_KEY_ID", cfg.access_key)
        gdal.SetConfigOption("AWS_SECRET_ACCESS_KEY", cfg.secret_key)

    def upload(self, local: str, key: str) -> None:
        dst = "/vsis3/%s/%s" % (self.bucket, key)
        ret = self.gdal.CopyFile(local, dst)
        if ret != 0:
            raise SyncError("upload to %s failed: %s" % (dst, self.gdal.GetLastErrorMsg()))


# --------------------------------------------------------------------------------------------------
# Pure helpers (unit-tested)
# --------------------------------------------------------------------------------------------------


def footprint_ring(gt: Tuple[float, ...], cols: int, rows: int,
                   to_lonlat: Callable[[float, float], Tuple[float, float]], per_edge: int = 8) -> List[List[float]]:
    """Closed [lon, lat] ring around a raster, with `per_edge` points per side so a reprojected edge stays true."""

    def px(c: float, r: float) -> Tuple[float, float]:
        return gt[0] + c * gt[1] + r * gt[2], gt[3] + c * gt[4] + r * gt[5]

    pts: List[Tuple[float, float]] = []
    for i in range(per_edge):                      # top edge, west to east
        pts.append(px(cols * i / per_edge, 0))
    for i in range(per_edge):                      # east edge, north to south
        pts.append(px(cols, rows * i / per_edge))
    for i in range(per_edge):                      # bottom edge, east to west
        pts.append(px(cols * (per_edge - i) / per_edge, rows))
    for i in range(per_edge):                      # west edge, south to north
        pts.append(px(0, rows * (per_edge - i) / per_edge))
    ring = [[round(v, 7) for v in to_lonlat(x, y)] for x, y in pts]
    # GeoJSON wants counter-clockwise exterior rings; the walk above is clockwise in map space.
    ring.reverse()
    ring.append(list(ring[0]))
    return ring


def ring_bbox(ring: List[List[float]]) -> List[float]:
    xs = [p[0] for p in ring]
    ys = [p[1] for p in ring]
    return [min(xs), min(ys), max(xs), max(ys)]


def union_bbox(boxes: List[List[float]]) -> List[float]:
    return [min(b[0] for b in boxes), min(b[1] for b in boxes), max(b[2] for b in boxes), max(b[3] for b in boxes)]


def iso_utc(value: Optional[str]) -> str:
    """Normalise a WebODM timestamp (ISO 8601, usually with a numeric offset) to `...Z`."""
    if not value:
        return dt.datetime.now(dt.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%fZ")
    text = value.strip()
    if text.endswith("Z"):
        text = text[:-1] + "+00:00"
    # Python before 3.11 only parses 3 or 6 fractional digits: normalise to 6.
    text = re.sub(r"\.(\d+)", lambda m: "." + (m.group(1) + "000000")[:6], text, count=1)
    parsed = dt.datetime.fromisoformat(text)
    if parsed.tzinfo is None:
        parsed = parsed.replace(tzinfo=dt.timezone.utc)
    return parsed.astimezone(dt.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%fZ")


def task_assets(task: Dict[str, Any]) -> List[Tuple[str, str, List[str], str]]:
    """The ASSET_SPECS present on a WebODM task."""
    available = set(task.get("available_assets") or [])
    return [spec for spec in ASSET_SPECS if spec[0] in available]


def fingerprint(task: Dict[str, Any]) -> str:
    """Changes whenever WebODM reprocesses the task or its outputs change."""
    basis = {
        "id": task["id"],
        "project": task["project"],
        "status": task.get("status"),
        "processing_time": task.get("processing_time"),
        "size": task.get("size"),
        "assets": sorted(spec[0] for spec in task_assets(task)),
    }
    return hashlib.sha256(json.dumps(basis, sort_keys=True).encode()).hexdigest()


def object_key(prefix: str, project_id: Any, task_id: str, stac_key: str) -> str:
    return "/".join(p for p in (prefix, str(project_id), task_id, stac_key + ".tif") if p)


def collection_definition() -> dict:
    return {
        "type": "Collection",
        "stac_version": STAC_VERSION,
        "id": COLLECTION_ID,
        "title": "WebODM drone imagery",
        "description": ("Orthophotos and elevation models (DSM, DTM) produced by the OneMind WebODM "
                        "photogrammetry pipeline from drone flights. One Item per WebODM task; the assets are "
                        "Cloud Optimized GeoTIFFs in DigitalOcean Spaces."),
        "license": "proprietary",
        "keywords": ["drone", "orthophoto", "dsm", "dtm", "webodm", "opendronemap"],
        "providers": [{"name": "OneMind", "roles": ["producer", "processor", "host"]}],
        "extent": {"spatial": {"bbox": [[-180.0, -90.0, 180.0, 90.0]]}, "temporal": {"interval": [[None, None]]}},
        "links": [],
    }


def build_item(task: Dict[str, Any], project: Dict[str, Any], uploaded: Dict[str, Dict[str, Any]],
               infos: Dict[str, RasterInfo], cfg: Config, existing: Optional[dict], now: str) -> dict:
    """The STAC Item for a task.

    `uploaded` maps the STAC asset key to {"key", "size", "sha256", "converted"}; `infos` maps it to RasterInfo.
    The footprint comes from the orthophoto when there is one, else the first elevation model.
    """
    reference = next(k for k in ("orthophoto", "dsm", "dtm") if k in infos)
    ref = infos[reference]
    assets: Dict[str, dict] = {}
    for _name, stac_key, roles, title in ASSET_SPECS:
        if stac_key not in uploaded:
            continue
        up, info = uploaded[stac_key], infos[stac_key]
        asset: Dict[str, Any] = {
            "href": "%s/%s" % (cfg.public_base_url, up["key"]),
            "type": COG_MEDIA_TYPE,
            "title": title,
            "roles": roles,
            "file:size": up["size"],
            "file:checksum": "1220" + up["sha256"],  # multihash: sha2-256
            "proj:shape": info.shape,
            "proj:transform": info.transform,
            "alternate": {"s3": {"href": "s3://%s/%s" % (cfg.bucket, up["key"]),
                                 "title": "S3 protocol (DigitalOcean Spaces)"}},
        }
        if info.code:
            asset["proj:code"] = info.code
        if info.gsd:
            asset["gsd"] = info.gsd
        assets[stac_key] = asset

    props: Dict[str, Any] = {
        "datetime": iso_utc(task.get("created_at")),
        "title": task.get("name") or ("WebODM task %s" % task["id"][:8]),
        "created": (existing or {}).get("properties", {}).get("created", now),
        "updated": now,
        "webodm:task_id": task["id"],
        "webodm:project_id": task["project"],
        "webodm:project_name": project.get("name"),
        "webodm:images_count": task.get("images_count"),
        "webodm:processing_time_ms": task.get("processing_time"),
        "webodm:fingerprint": fingerprint(task),
    }
    if ref.code:
        props["proj:code"] = ref.code
    if ref.gsd:
        props["gsd"] = ref.gsd
    return {
        "type": "Feature",
        "stac_version": STAC_VERSION,
        "stac_extensions": [EXT_PROJECTION, EXT_FILE, EXT_ALTERNATE],
        "id": task["id"],
        "collection": COLLECTION_ID,
        "geometry": {"type": "Polygon", "coordinates": [ref.footprint]},
        "bbox": union_bbox([i.bbox for i in infos.values()]),
        "properties": props,
        "assets": assets,
        "links": [],
    }


def sha256_file(path: str) -> str:
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


# --------------------------------------------------------------------------------------------------
# The sync
# --------------------------------------------------------------------------------------------------


@dataclass
class Result:
    created: List[str] = field(default_factory=list)
    updated: List[str] = field(default_factory=list)
    skipped: List[str] = field(default_factory=list)
    failed: Dict[str, str] = field(default_factory=dict)
    deferred: int = 0


def ensure_collection(stac: StacClient) -> None:
    """Create the `webodm` Collection, or bring its descriptive fields back to the definition."""
    desired = collection_definition()
    existing = stac.get_collection(COLLECTION_ID)
    if existing is None:
        LOG.info("creating collection %s", COLLECTION_ID)
        stac.put_collection(desired, exists=False)
        return
    if any(existing.get(f) != desired[f] for f in COLLECTION_FIELDS):
        LOG.info("updating collection %s", COLLECTION_ID)
        merged = dict(existing)
        for f in COLLECTION_FIELDS:
            merged[f] = desired[f]
        stac.put_collection(merged, exists=True)


def process_task(task: Dict[str, Any], project: Dict[str, Any], cfg: Config, webodm: WebOdmClient,
                 stac: StacClient, raster: Any, storage: Any, existing: Optional[dict]) -> str:
    """Upload one task's rasters and write its Item. Returns "created" or "updated"."""
    uploaded: Dict[str, Dict[str, Any]] = {}
    infos: Dict[str, RasterInfo] = {}
    work = tempfile.mkdtemp(prefix="task-", dir=cfg.work_dir)
    try:
        for asset_file, stac_key, _roles, _title in task_assets(task):
            local = os.path.join(work, asset_file)
            LOG.info("task %s: downloading %s", task["id"], asset_file)
            webodm.download(task["project"], task["id"], asset_file, local)
            converted = False
            if not raster.is_cog(local):
                LOG.info("task %s: %s is not a COG, converting", task["id"], asset_file)
                cog = os.path.join(work, "cog-" + asset_file)
                raster.to_cog(local, cog)
                os.replace(cog, local)
                converted = True
            infos[stac_key] = raster.describe(local)
            key = object_key(cfg.prefix, task["project"], task["id"], stac_key)
            size = os.path.getsize(local)
            digest = sha256_file(local)
            LOG.info("task %s: uploading %s (%d bytes) to %s", task["id"], asset_file, size, key)
            storage.upload(local, key)
            uploaded[stac_key] = {"key": key, "size": size, "sha256": digest, "converted": converted}
            os.remove(local)  # one raster on disk at a time
    finally:
        shutil.rmtree(work, ignore_errors=True)

    now = dt.datetime.now(dt.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%fZ")
    item = build_item(task, project, uploaded, infos, cfg, existing, now)
    stac.put_item(COLLECTION_ID, item, exists=existing is not None)
    return "updated" if existing is not None else "created"


def sync_once(cfg: Config, webodm: WebOdmClient, stac: StacClient, raster: Any, storage: Any) -> Result:
    result = Result()
    webodm.login()
    ensure_collection(stac)
    pending: List[Tuple[Dict[str, Any], Dict[str, Any], Optional[dict]]] = []
    for project in webodm.projects():
        for task in webodm.completed_tasks(project["id"]):
            if task.get("pending_action") is not None or not task_assets(task):
                continue
            existing = stac.get_item(COLLECTION_ID, task["id"])
            if (existing and not cfg.force
                    and existing.get("properties", {}).get("webodm:fingerprint") == fingerprint(task)):
                result.skipped.append(task["id"])
                continue
            pending.append((task, project, existing))

    # Oldest first, a few per run, so a first run over a large backlog finishes within the job's deadline.
    pending.sort(key=lambda p: p[0].get("created_at") or "")
    result.deferred = max(0, len(pending) - cfg.max_tasks)
    for task, project, existing in pending[: cfg.max_tasks]:
        if cfg.dry_run:
            LOG.info("dry run: would sync task %s (%s)", task["id"], task.get("name"))
            continue
        try:
            outcome = process_task(task, project, cfg, webodm, stac, raster, storage, existing)
            (result.created if outcome == "created" else result.updated).append(task["id"])
            LOG.info("task %s: %s", task["id"], outcome)
        except Exception as err:  # keep going: one bad task must not block the rest
            LOG.error("task %s failed: %s", task["id"], err)
            result.failed[task["id"]] = str(err)
    return result


def main() -> int:
    logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s", stream=sys.stdout)
    cfg = Config.from_env()
    missing = cfg.missing()
    if missing:
        LOG.warning("not configured yet, nothing to do. Missing: %s. See platform/stac/README.md.",
                    "; ".join(missing))
        return 0
    os.makedirs(cfg.work_dir, exist_ok=True)
    try:
        raster = GdalRaster()
        storage = SpacesStorage(cfg)
        result = sync_once(cfg, WebOdmClient(cfg.webodm_url, cfg.webodm_username, cfg.webodm_password),
                           StacClient(cfg.stac_url), raster, storage)
    except WebOdmUnavailable as err:
        LOG.warning("%s; WebODM is probably scaled down. Will retry on the next schedule.", err)
        return 0
    except SyncError as err:
        LOG.error("%s", err)
        return 1
    LOG.info("done: %d created, %d updated, %d unchanged, %d failed, %d deferred to the next run",
             len(result.created), len(result.updated), len(result.skipped), len(result.failed), result.deferred)
    return 1 if result.failed else 0


if __name__ == "__main__":
    sys.exit(main())
