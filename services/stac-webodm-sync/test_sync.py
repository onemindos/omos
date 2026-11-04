"""Unit tests for sync.py. Standard library only: GDAL, WebODM, STAC and Spaces are replaced by fakes.

Run from the repository root:  python -m unittest discover -s services/stac-webodm-sync -v
"""
import json
import logging
import os
import sys
import tempfile
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import sync  # noqa: E402

logging.disable(logging.CRITICAL)

TASK_ID = "5b1a6f4e-0b0e-4a3c-9a52-6a0f5d6a3c11"


def make_cfg(**over):
    base = dict(
        webodm_url="http://webodm", webodm_username="u", webodm_password="p", stac_url="http://stac",
        bucket="onemind-imagery", endpoint="nyc3.digitaloceanspaces.com", region="nyc3", prefix="webodm",
        public_base_url="https://onemind-imagery.nyc3.digitaloceanspaces.com", work_dir=tempfile.gettempdir(),
        max_tasks=3, force=False, dry_run=False, access_key="k", secret_key="s",
    )
    base.update(over)
    return sync.Config(**base)


def make_task(task_id=TASK_ID, project=7, assets=("orthophoto.tif", "dsm.tif"), **over):
    task = {
        "id": task_id, "project": project, "name": "Barn flight", "status": 40, "pending_action": None,
        "created_at": "2026-10-20T14:30:00.123456-04:00", "processing_time": 123456, "size": 812.5,
        "images_count": 240, "available_assets": list(assets) + ["all.zip"],
    }
    task.update(over)
    return task


INFO = sync.RasterInfo(
    footprint=[[-80.0, 25.0], [-80.0, 25.1], [-79.9, 25.1], [-79.9, 25.0], [-80.0, 25.0]],
    bbox=[-80.0, 25.0, -79.9, 25.1], code="EPSG:32617", shape=[1000, 2000],
    transform=[0.05, 0.0, 500000.0, 0.0, -0.05, 2765000.0], gsd=0.05,
)


class FakeRaster:
    def __init__(self, cog=True):
        self.cog = cog
        self.converted = []

    def is_cog(self, path):
        return self.cog

    def to_cog(self, src, dst):
        self.converted.append(os.path.basename(src))
        with open(src, "rb") as f, open(dst, "wb") as g:
            g.write(f.read() + b"+cog")

    def describe(self, path):
        return INFO


class FakeStorage:
    def __init__(self):
        self.objects = {}

    def upload(self, local, key):
        with open(local, "rb") as f:
            self.objects[key] = f.read()


class FakeWebOdm:
    def __init__(self, projects):
        self.projects_ = projects  # [{"id", "name", "tasks": [...]}]
        self.logins = 0

    def login(self):
        self.logins += 1

    def projects(self):
        return [{"id": p["id"], "name": p["name"]} for p in self.projects_]

    def completed_tasks(self, project_id):
        return [t for p in self.projects_ if p["id"] == project_id for t in p["tasks"] if t["status"] == 40]

    def download(self, project_id, task_id, asset, dest):
        with open(dest, "wb") as f:
            f.write(("raster:%s:%s" % (task_id, asset)).encode())


class FakeStac:
    def __init__(self):
        self.collections = {}
        self.items = {}
        self.writes = []
        self.fail_item_ids = set()

    def get_collection(self, cid):
        return self.collections.get(cid)

    def put_collection(self, collection, exists):
        self.writes.append(("collection", "PUT" if exists else "POST", collection["id"]))
        self.collections[collection["id"]] = collection

    def get_item(self, cid, iid):
        return self.items.get((cid, iid))

    def put_item(self, cid, item, exists):
        if item["id"] in self.fail_item_ids:
            raise sync.SyncError("boom")
        self.writes.append(("item", "PUT" if exists else "POST", item["id"]))
        self.items[(cid, item["id"])] = item


class PureHelpers(unittest.TestCase):
    def test_footprint_ring_is_closed_counter_clockwise_and_exact(self):
        gt = (10.0, 1.0, 0.0, 20.0, 0.0, -1.0)  # 4 x 2 raster, origin top-left, north up
        ring = sync.footprint_ring(gt, 4, 2, lambda x, y: (x, y), per_edge=2)
        self.assertEqual(ring[0], ring[-1])
        self.assertEqual(sync.ring_bbox(ring), [10.0, 18.0, 14.0, 20.0])
        area = sum(ring[i][0] * ring[i + 1][1] - ring[i + 1][0] * ring[i][1] for i in range(len(ring) - 1)) / 2
        self.assertGreater(area, 0, "GeoJSON exterior rings must be counter-clockwise")
        self.assertEqual(area, 8.0)

    def test_footprint_ring_applies_the_reprojection_to_every_point(self):
        ring = sync.footprint_ring((0, 1, 0, 0, 0, -1), 2, 2, lambda x, y: (x / 10, y / 10), per_edge=1)
        self.assertEqual(sync.ring_bbox(ring), [0.0, -0.2, 0.2, 0.0])

    def test_union_bbox(self):
        self.assertEqual(sync.union_bbox([[0, 0, 1, 1], [-1, 0.5, 0.5, 3]]), [-1, 0, 1, 3])

    def test_iso_utc_converts_offsets_and_accepts_z_and_naive(self):
        self.assertEqual(sync.iso_utc("2026-10-20T14:30:00.123456-04:00"), "2026-10-20T18:30:00.123456Z")
        self.assertEqual(sync.iso_utc("2026-10-20T18:30:00Z"), "2026-10-20T18:30:00.000000Z")
        self.assertEqual(sync.iso_utc("2026-10-20T18:30:00"), "2026-10-20T18:30:00.000000Z")
        self.assertEqual(sync.iso_utc("2026-10-20T18:30:00.12345Z"), "2026-10-20T18:30:00.123450Z")
        self.assertTrue(sync.iso_utc(None).endswith("Z"))

    def test_task_assets_keeps_only_known_rasters_in_a_fixed_order(self):
        keys = [spec[1] for spec in sync.task_assets(make_task(assets=("dtm.tif", "orthophoto.tif")))]
        self.assertEqual(keys, ["orthophoto", "dtm"])
        self.assertEqual(sync.task_assets(make_task(assets=())), [])

    def test_fingerprint_is_stable_and_tracks_reprocessing(self):
        base = make_task()
        self.assertEqual(sync.fingerprint(base), sync.fingerprint(make_task(name="renamed")))
        self.assertNotEqual(sync.fingerprint(base), sync.fingerprint(make_task(processing_time=999)))
        self.assertNotEqual(sync.fingerprint(base), sync.fingerprint(make_task(assets=("orthophoto.tif",))))

    def test_object_key(self):
        self.assertEqual(sync.object_key("webodm", 7, TASK_ID, "dsm"), "webodm/7/%s/dsm.tif" % TASK_ID)
        self.assertEqual(sync.object_key("", 7, "t", "dsm"), "7/t/dsm.tif")

    def test_config_from_env_defaults_and_missing(self):
        cfg = sync.Config.from_env({"SPACES_BUCKET": "b", "SPACES_REGION": "sfo3"})
        self.assertEqual(cfg.endpoint, "sfo3.digitaloceanspaces.com")
        self.assertEqual(cfg.public_base_url, "https://b.sfo3.digitaloceanspaces.com")
        self.assertEqual(len(cfg.missing()), 2)
        ok = sync.Config.from_env({"WEBODM_USERNAME": "u", "WEBODM_PASSWORD": "p",
                                   "AWS_ACCESS_KEY_ID": "k", "AWS_SECRET_ACCESS_KEY": "s"})
        self.assertEqual(ok.missing(), [])
        self.assertNotIn("'p'", repr(ok))  # secrets stay out of repr
        self.assertNotIn("'k'", repr(ok))

    def test_collection_definition_is_a_valid_looking_stac_collection(self):
        c = sync.collection_definition()
        self.assertEqual((c["type"], c["id"], c["stac_version"]), ("Collection", "webodm", "1.1.0"))
        for required in ("description", "license", "extent", "links"):
            self.assertIn(required, c)


class BuildItem(unittest.TestCase):
    def build(self, existing=None, assets=("orthophoto.tif", "dsm.tif")):
        uploaded = {}
        infos = {}
        for _name, key, _roles, _title in sync.task_assets(make_task(assets=assets)):
            uploaded[key] = {"key": "webodm/7/%s/%s.tif" % (TASK_ID, key), "size": 10, "sha256": "ab" * 32,
                             "converted": False}
            infos[key] = INFO
        return sync.build_item(make_task(assets=assets), {"id": 7, "name": "Homestead"}, uploaded, infos,
                               make_cfg(), existing, "2026-10-23T00:00:00.000000Z")

    def test_geometry_bbox_and_datetime_come_from_the_raster_and_the_task(self):
        item = self.build()
        self.assertEqual(item["id"], TASK_ID)
        self.assertEqual(item["collection"], "webodm")
        self.assertEqual(item["geometry"]["coordinates"], [INFO.footprint])
        self.assertEqual(item["bbox"], INFO.bbox)
        self.assertEqual(item["properties"]["datetime"], "2026-10-20T18:30:00.123456Z")

    def test_assets_have_https_hrefs_roles_media_type_and_checksum(self):
        item = self.build()
        self.assertEqual(sorted(item["assets"]), ["dsm", "orthophoto"])
        ortho = item["assets"]["orthophoto"]
        self.assertEqual(ortho["href"], "https://onemind-imagery.nyc3.digitaloceanspaces.com/webodm/7/%s/orthophoto.tif"
                         % TASK_ID)
        self.assertEqual(ortho["roles"], ["data", "visual"])
        self.assertIn("cloud-optimized", ortho["type"])
        self.assertEqual(ortho["file:checksum"], "1220" + "ab" * 32)
        self.assertEqual(ortho["proj:code"], "EPSG:32617")
        self.assertEqual(ortho["alternate"]["s3"]["href"], "s3://onemind-imagery/webodm/7/%s/orthophoto.tif" % TASK_ID)
        self.assertEqual(item["assets"]["dsm"]["roles"], ["data"])

    def test_dsm_only_task_uses_the_dsm_footprint(self):
        item = self.build(assets=("dsm.tif",))
        self.assertEqual(list(item["assets"]), ["dsm"])
        self.assertEqual(item["geometry"]["coordinates"], [INFO.footprint])

    def test_created_is_kept_from_the_existing_item_and_fingerprint_is_stored(self):
        item = self.build(existing={"properties": {"created": "2026-01-01T00:00:00Z"}})
        self.assertEqual(item["properties"]["created"], "2026-01-01T00:00:00Z")
        self.assertEqual(item["properties"]["updated"], "2026-10-23T00:00:00.000000Z")
        self.assertEqual(item["properties"]["webodm:fingerprint"], sync.fingerprint(make_task()))
        self.assertIn(sync.EXT_PROJECTION, item["stac_extensions"])


class EnsureCollection(unittest.TestCase):
    def test_creates_when_missing_then_leaves_alone(self):
        stac = FakeStac()
        sync.ensure_collection(stac)
        sync.ensure_collection(stac)
        self.assertEqual(stac.writes, [("collection", "POST", "webodm")])

    def test_updates_descriptive_fields_but_keeps_the_extent(self):
        stac = FakeStac()
        stale = sync.collection_definition()
        stale["title"] = "old title"
        stale["extent"] = {"spatial": {"bbox": [[1, 2, 3, 4]]}, "temporal": {"interval": [["2026-01-01T00:00:00Z", None]]}}
        stac.collections["webodm"] = stale
        sync.ensure_collection(stac)
        self.assertEqual(stac.writes, [("collection", "PUT", "webodm")])
        written = stac.collections["webodm"]
        self.assertEqual(written["title"], sync.collection_definition()["title"])
        self.assertEqual(written["extent"]["spatial"]["bbox"], [[1, 2, 3, 4]])


class SyncOnce(unittest.TestCase):
    def setUp(self):
        self.cfg = make_cfg()
        self.stac = FakeStac()
        self.storage = FakeStorage()
        self.raster = FakeRaster()

    def run_sync(self, projects, cfg=None):
        webodm = FakeWebOdm(projects)
        result = sync.sync_once(cfg or self.cfg, webodm, self.stac, self.raster, self.storage)
        return result, webodm

    def test_new_task_is_uploaded_and_published(self):
        result, webodm = self.run_sync([{"id": 7, "name": "Homestead", "tasks": [make_task()]}])
        self.assertEqual((result.created, result.updated, result.skipped, result.failed), ([TASK_ID], [], [], {}))
        self.assertEqual(webodm.logins, 1)
        self.assertEqual(sorted(self.storage.objects), [
            "webodm/7/%s/dsm.tif" % TASK_ID, "webodm/7/%s/orthophoto.tif" % TASK_ID])
        item = self.stac.items[("webodm", TASK_ID)]
        self.assertEqual(item["properties"]["webodm:project_name"], "Homestead")
        self.assertEqual(item["assets"]["dsm"]["file:size"], len(("raster:%s:dsm.tif" % TASK_ID).encode()))
        self.assertEqual(self.raster.converted, [])

    def test_second_run_is_a_no_op(self):
        projects = [{"id": 7, "name": "H", "tasks": [make_task()]}]
        self.run_sync(projects)
        self.storage.objects.clear()
        before = list(self.stac.writes)
        result, _ = self.run_sync(projects)
        self.assertEqual((result.created, result.updated, result.skipped), ([], [], [TASK_ID]))
        self.assertEqual(self.stac.writes, before)
        self.assertEqual(self.storage.objects, {})

    def test_reprocessed_task_updates_the_same_item(self):
        self.run_sync([{"id": 7, "name": "H", "tasks": [make_task()]}])
        result, _ = self.run_sync([{"id": 7, "name": "H", "tasks": [make_task(processing_time=555)]}])
        self.assertEqual(result.updated, [TASK_ID])
        self.assertEqual(self.stac.writes[-1], ("item", "PUT", TASK_ID))
        self.assertEqual(len(self.stac.items), 1)

    def test_force_resyncs_unchanged_tasks(self):
        projects = [{"id": 7, "name": "H", "tasks": [make_task()]}]
        self.run_sync(projects)
        result, _ = self.run_sync(projects, make_cfg(force=True))
        self.assertEqual(result.updated, [TASK_ID])

    def test_non_cog_is_converted_before_upload(self):
        self.raster = FakeRaster(cog=False)
        self.run_sync([{"id": 7, "name": "H", "tasks": [make_task(assets=("dsm.tif",))]}])
        self.assertEqual(self.raster.converted, ["dsm.tif"])
        self.assertTrue(self.storage.objects["webodm/7/%s/dsm.tif" % TASK_ID].endswith(b"+cog"))

    def test_only_completed_tasks_with_rasters_and_no_pending_action(self):
        tasks = [
            make_task("a" * 8, status=20),                       # still running
            make_task("b" * 8, assets=()),                      # nothing we publish
            make_task("c" * 8, pending_action=1),               # being restarted or removed
            make_task("d" * 8),
        ]
        result, _ = self.run_sync([{"id": 7, "name": "H", "tasks": tasks}])
        self.assertEqual(result.created, ["d" * 8])

    def test_backlog_is_limited_per_run_oldest_first(self):
        tasks = [make_task("t%d" % i, created_at="2026-10-0%dT00:00:00Z" % (i + 1)) for i in range(5)]
        result, _ = self.run_sync([{"id": 7, "name": "H", "tasks": list(reversed(tasks))}], make_cfg(max_tasks=2))
        self.assertEqual(result.created, ["t0", "t1"])
        self.assertEqual(result.deferred, 3)

    def test_one_failing_task_does_not_stop_the_others(self):
        self.stac.fail_item_ids.add("bad")
        tasks = [make_task("bad"), make_task("good")]
        result, _ = self.run_sync([{"id": 7, "name": "H", "tasks": tasks}])
        self.assertEqual(result.created, ["good"])
        self.assertEqual(list(result.failed), ["bad"])

    def test_dry_run_writes_nothing(self):
        result, _ = self.run_sync([{"id": 7, "name": "H", "tasks": [make_task()]}], make_cfg(dry_run=True))
        self.assertEqual((result.created, self.storage.objects, list(self.stac.items)), ([], {}, []))

    def test_scratch_directories_are_removed(self):
        with tempfile.TemporaryDirectory() as work:
            self.run_sync([{"id": 7, "name": "H", "tasks": [make_task()]}], make_cfg(work_dir=work))
            self.assertEqual(os.listdir(work), [])


class Http(unittest.TestCase):
    def scripted(self, replies):
        calls = []

        def http(method, url, headers, body=None):
            calls.append((method, url, headers, body))
            return replies.pop(0)

        return http, calls

    def test_webodm_login_sends_credentials_and_uses_jwt_header(self):
        http, calls = self.scripted([(200, {}, b'{"token": "T"}'), (200, {}, b'{"results": [], "next": null}')])
        client = sync.WebOdmClient("http://w", "u", "p", http)
        client.login()
        client.projects()
        self.assertEqual(json.loads(calls[0][3]), {"username": "u", "password": "p"})
        self.assertEqual(calls[1][2]["Authorization"], "JWT T")

    def test_webodm_projects_follow_pagination(self):
        http, calls = self.scripted([
            (200, {}, json.dumps({"results": [{"id": 1}], "next": "x"}).encode()),
            (200, {}, json.dumps({"results": [{"id": 2}], "next": None}).encode()),
        ])
        client = sync.WebOdmClient("http://w", "u", "p", http)
        client.token = "T"
        self.assertEqual([p["id"] for p in client.projects()], [1, 2])
        self.assertTrue(calls[1][1].endswith("page=2"))

    def test_webodm_task_list_accepts_plain_list_and_filters_on_status(self):
        http, calls = self.scripted([(200, {}, b'[{"id": "t"}]')])
        client = sync.WebOdmClient("http://w", "u", "p", http)
        client.token = "T"
        self.assertEqual(client.completed_tasks(7), [{"id": "t"}])
        self.assertIn("/api/projects/7/tasks/?status=40", calls[0][1])

    def test_webodm_relogs_in_once_on_401(self):
        http, calls = self.scripted([(401, {}, b"{}"), (200, {}, b'{"token": "T2"}'), (200, {}, b'[]')])
        client = sync.WebOdmClient("http://w", "u", "p", http)
        client.token = "T1"
        client.completed_tasks(1)
        self.assertEqual(calls[2][2]["Authorization"], "JWT T2")

    def test_webodm_unreachable_is_reported_as_unavailable(self):
        def http(method, url, headers, body=None):
            raise ConnectionRefusedError("refused")

        client = sync.WebOdmClient("http://w", "u", "p", http)
        client.http = http  # skip the retry wrapper's sleeps
        with self.assertRaises(sync.WebOdmUnavailable):
            client.login()

    def test_stac_item_post_falls_back_to_put_on_conflict(self):
        http, calls = self.scripted([(409, {}, b'{"code": "ConflictError"}'), (200, {}, b"{}")])
        client = sync.StacClient("http://s", http)
        client.put_item("webodm", {"id": "i"}, exists=False)
        self.assertEqual([c[0] for c in calls], ["POST", "PUT"])
        self.assertTrue(calls[1][1].endswith("/collections/webodm/items/i"))

    def test_stac_get_returns_none_on_404_and_raises_on_errors(self):
        http, _ = self.scripted([(404, {}, b"{}"), (500, {}, b"oops")])
        client = sync.StacClient("http://s", http)
        self.assertIsNone(client.get_item("webodm", "x"))
        with self.assertRaises(sync.SyncError):
            client.get_item("webodm", "x")


class Main(unittest.TestCase):
    def test_exits_zero_and_does_nothing_until_credentials_exist(self):
        old = dict(os.environ)
        try:
            for k in ("WEBODM_USERNAME", "WEBODM_PASSWORD", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"):
                os.environ.pop(k, None)
            self.assertEqual(sync.main(), 0)
        finally:
            os.environ.clear()
            os.environ.update(old)


if __name__ == "__main__":
    unittest.main()
