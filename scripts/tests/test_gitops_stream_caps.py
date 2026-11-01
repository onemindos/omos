"""Stream caps (A22): every NACK Stream CR in platform/nats/streams is capped and fits a node.

Static contracts always run. The engine tests translate each Stream CR into a JetStream
stream config, create it on a disposable loopback nats-server (replicas forced to 1) and
show that the limits behave as configured: maxAge ages out, maxBytes with discard: old drops
the oldest, discard: new refuses to store the publish but core subscribers still receive it
(why CMD stays on discard: old, A22-P2), and shrinking a limit on an existing stream (what
NACK does when Argo syncs a changed CR) removes stored data at once.

Engine tests need NATS_SERVER_BIN (or nats-server on PATH) and the nats-py package; they
skip otherwise. No cluster, credentials or production data are used.
"""
from __future__ import annotations

import asyncio
import os
from pathlib import Path
import re
import shutil
import socket
import subprocess
import tempfile
import time
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
STREAMS_DIR = ROOT / "platform/nats/streams"
RULES = ROOT / "platform/monitoring/prometheus-rules-onemind.yaml"

GIB = 1024**3
# Per-node JetStream budget (BUS-DESIGN section 4.1): max_file_store 45 GiB (live varz
# 48318382080), reserved_storage 2001048576 (live jsz nats-0, 2026-10-29), object buckets
# up to 9 GiB and KV buckets up to 2 GiB. Every R3 stream counts fully on every node.
NODE_MAX_STORE = 48318382080
RESERVED = 2001048576
OBJECT_BUDGET = 9 * GIB
KV_BUDGET = 2 * GIB

_UNITS = {"ns": 1e-9, "us": 1e-6, "ms": 1e-3, "s": 1, "m": 60, "h": 3600}


def duration_seconds(value: str) -> float:
    parts = re.findall(r"(\d+(?:\.\d+)?)(ns|us|ms|s|m|h)", value)
    if not parts or "".join(n + u for n, u in parts) != value:
        raise ValueError(f"unsupported duration {value!r}")
    return sum(float(n) * _UNITS[u] for n, u in parts)


def load_crs(kind: str) -> dict[str, dict]:
    crs = {}
    for path in sorted(STREAMS_DIR.glob("*.yaml")):
        for doc in yaml.safe_load_all(path.read_text()):
            if doc and doc.get("kind") == kind:
                crs[path.name] = doc["spec"]
    return crs


STREAMS = load_crs("Stream")
KVS = load_crs("KeyValue")


def by_name(name: str) -> dict:
    return next(spec for spec in STREAMS.values() if spec["name"] == name)


class StreamCapContracts(unittest.TestCase):
    def test_every_stream_has_a_byte_cap(self):
        missing = [f for f, s in STREAMS.items() if int(s.get("maxBytes", -1)) <= 0]
        self.assertEqual(missing, [], "every Stream CR needs maxBytes (BUS-DESIGN 4.1)")

    def test_every_stream_keeps_safety_settings(self):
        for f, s in STREAMS.items():
            with self.subTest(f):
                self.assertTrue(s.get("preventDelete"), "preventDelete keeps data if the CR goes")
                self.assertEqual(s.get("storage"), "file")
                self.assertEqual(s.get("retention"), "limits")
                self.assertIn(s.get("discard"), ("old", "new"))

    def test_duplicate_window_fits_max_age(self):
        # The server rejects a stream whose duplicate window is longer than its max age.
        for f, s in STREAMS.items():
            if "duplicateWindow" in s and "maxAge" in s:
                with self.subTest(f):
                    self.assertLessEqual(duration_seconds(s["duplicateWindow"]), duration_seconds(s["maxAge"]))

    def test_caps_fit_one_node(self):
        streams = sum(int(s["maxBytes"]) for s in STREAMS.values())
        kv_set = sum(int(s.get("maxBytes", 0)) for s in KVS.values())
        self.assertLessEqual(kv_set, KV_BUDGET, "KV caps exceed the 2 GiB KV budget")
        total = streams + OBJECT_BUDGET + KV_BUDGET
        self.assertLessEqual(total, NODE_MAX_STORE - RESERVED,
                             f"streams {streams / GIB:.2f} GiB + objects + KV = {total / GIB:.2f} GiB")

    def test_design_targets(self):
        cot, cmd, tel, det = by_name("COT"), by_name("CMD"), by_name("TEL"), by_name("DET")
        self.assertEqual(duration_seconds(cot["maxAge"]), 72 * 3600)
        self.assertEqual(int(cot["maxBytes"]), 3 * GIB)
        # A22-P2 option b: discard: new does not stop core delivery, so CMD keeps discard: old
        # until every cmd.* consumer reads through a JetStream consumer (option a).
        self.assertEqual(cmd["discard"], "old")
        self.assertEqual(duration_seconds(cmd["maxAge"]), 365 * 86400)
        self.assertEqual(int(cmd["maxBytes"]), 512 * 1024**2)
        self.assertIn("tel.feed.>", tel["subjects"])
        self.assertEqual(duration_seconds(tel["maxAge"]), 72 * 3600)
        self.assertEqual(det["subjects"], ["det.>"])
        self.assertEqual(duration_seconds(det["maxAge"]), 7 * 86400)

    def test_cmd_alert_threshold_tracks_cap(self):
        rules = yaml.safe_load(RULES.read_text())
        alerts = [r for g in rules["spec"]["groups"] for r in g["rules"] if r.get("alert") == "NATSCmdStreamNearFull"]
        self.assertEqual(len(alerts), 1)
        threshold = int(re.search(r">\s*(\d+)", alerts[0]["expr"]).group(1))
        self.assertEqual(threshold, int(int(by_name("CMD")["maxBytes"]) * 0.8))


NATS_SERVER = os.environ.get("NATS_SERVER_BIN") or shutil.which("nats-server")
try:
    import nats  # type: ignore
    from nats.js import api as jsapi  # type: ignore
    from nats.js.errors import APIError  # type: ignore
except ImportError:  # pragma: no cover - depends on the runner
    nats = None


def to_config(spec: dict, **override):
    """Translate a NACK v1beta2 Stream spec into a JetStream StreamConfig (replicas 1)."""
    cfg = dict(
        name=spec["name"],
        subjects=list(spec["subjects"]),
        retention=jsapi.RetentionPolicy(spec.get("retention", "limits")),
        storage=jsapi.StorageType(spec.get("storage", "file")),
        discard=jsapi.DiscardPolicy(spec.get("discard", "old")),
        max_age=duration_seconds(spec["maxAge"]) if "maxAge" in spec else None,
        max_bytes=int(spec.get("maxBytes", -1)),
        max_msgs_per_subject=int(spec.get("maxMsgsPerSubject", -1)),
        duplicate_window=duration_seconds(spec["duplicateWindow"]) if "duplicateWindow" in spec else None,
        allow_direct=bool(spec.get("allowDirect", False)),
        num_replicas=1,
    )
    cfg.update(override)
    return jsapi.StreamConfig(**cfg)


def free_port() -> int:
    with socket.socket() as s:
        s.bind(("127.0.0.1", 0))
        return s.getsockname()[1]


@unittest.skipUnless(NATS_SERVER and nats, "Install nats-server and nats-py (or set NATS_SERVER_BIN)")
class StreamCapEngine(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.scratch = tempfile.TemporaryDirectory(prefix="stream-caps-")
        cls.port = free_port()
        # max_file_store as live (45 GiB): JetStream reserves each stream's maxBytes against it
        # at create time, so creating every CR also checks the per-node budget. Nothing is
        # written up front; the tests store a few hundred KiB.
        conf = Path(cls.scratch.name) / "server.conf"
        conf.write_text(f'listen: "127.0.0.1:{cls.port}"\n'
                        f'jetstream {{ store_dir: "{cls.scratch.name}", max_file_store: {NODE_MAX_STORE} }}\n')
        cls.server = subprocess.Popen(
            [NATS_SERVER, "-c", str(conf)],
            stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        deadline = time.time() + 10
        while time.time() < deadline:
            try:
                socket.create_connection(("127.0.0.1", cls.port), 0.2).close()
                break
            except OSError:
                time.sleep(0.1)
        else:
            raise RuntimeError("nats-server did not start")

    @classmethod
    def tearDownClass(cls):
        cls.server.terminate()
        cls.server.wait(10)
        cls.scratch.cleanup()

    def run_js(self, coro_fn):
        async def main():
            nc = await nats.connect(f"nats://127.0.0.1:{self.port}")
            try:
                js = nc.jetstream()
                for info in await js.streams_info():
                    await js.delete_stream(info.config.name)
                return await coro_fn(js)
            finally:
                await nc.close()
        return asyncio.run(main())

    def test_every_cr_creates_with_its_limits(self):
        async def body(js):
            for f, spec in STREAMS.items():
                await js.add_stream(to_config(spec))
                got = (await js.stream_info(spec["name"])).config
                with self.subTest(f):
                    self.assertEqual(got.subjects, spec["subjects"])
                    self.assertEqual(got.max_bytes, int(spec["maxBytes"]))
                    self.assertEqual(getattr(got.discard, "value", got.discard), spec["discard"])
                    self.assertAlmostEqual(got.max_age, duration_seconds(spec["maxAge"]))
                    self.assertEqual(got.max_msgs_per_subject, int(spec.get("maxMsgsPerSubject", -1)))
                    self.assertEqual(bool(got.allow_direct), bool(spec.get("allowDirect", False)))
                    print(f"  {spec['name']:8} max_age={got.max_age:>10.0f}s max_bytes={got.max_bytes:>11} "
                          f"discard={getattr(got.discard, 'value', got.discard)} per_subject={got.max_msgs_per_subject} direct={got.allow_direct}")
        self.run_js(body)

    def test_cot_max_age_ages_out(self):
        async def body(js):
            await js.add_stream(to_config(by_name("COT"), max_age=2, duplicate_window=1))
            for i in range(20):
                await js.publish("tak.cot.inject", b"<event/>")
            self.assertEqual((await js.stream_info("COT")).state.messages, 20)
            await asyncio.sleep(3.5)
            self.assertEqual((await js.stream_info("COT")).state.messages, 0)
        self.run_js(body)

    def test_cot_shorter_max_age_on_update_removes_older_data(self):
        # What happens live when Argo syncs cot.yaml 720h -> 72h: NACK updates the stream.
        async def body(js):
            await js.add_stream(to_config(by_name("COT"), max_age=720 * 3600))
            for _ in range(50):
                await js.publish("tak.cot.inject", b"old")
            await asyncio.sleep(3)
            for _ in range(5):
                await js.publish("tak.cot.inject", b"new")
            await js.update_stream(to_config(by_name("COT"), max_age=2, duplicate_window=1))
            await asyncio.sleep(1)
            st = (await js.stream_info("COT")).state
            print(f"  COT after max_age shrink: messages={st.messages} first_seq={st.first_seq}")
            self.assertEqual(st.messages, 5)
            self.assertEqual(st.first_seq, 51)
        self.run_js(body)

    def test_cot_byte_cap_discards_oldest(self):
        async def body(js):
            cap = 64 * 1024
            await js.add_stream(to_config(by_name("COT"), max_bytes=cap))
            for _ in range(200):
                await js.publish("tak.cot.inject", b"x" * 1024)  # every publish is acked
            st = (await js.stream_info("COT")).state
            print(f"  COT 64 KiB cap after 200 x 1 KiB: messages={st.messages} bytes={st.bytes} first_seq={st.first_seq}")
            self.assertLessEqual(st.bytes, cap)
            self.assertGreater(st.first_seq, 1)
            self.assertEqual(st.last_seq, 200)
        self.run_js(body)

    def test_smaller_byte_cap_on_update_truncates(self):
        async def body(js):
            await js.add_stream(to_config(by_name("INTEL"), max_bytes=-1))
            for i in range(100):
                await js.publish(f"intel.test.src.event{i}", b"x" * 1024)
            await js.update_stream(to_config(by_name("INTEL"), max_bytes=32 * 1024))
            st = (await js.stream_info("INTEL")).state
            print(f"  INTEL after max_bytes shrink to 32 KiB: messages={st.messages} bytes={st.bytes}")
            self.assertLessEqual(st.bytes, 32 * 1024)
        self.run_js(body)

    def test_cmd_discard_old_keeps_recording_when_full(self):
        async def body(js):
            cap = 8 * 1024
            await js.add_stream(to_config(by_name("CMD"), max_bytes=cap))
            for i in range(50):
                await js.publish(f"cmd.fleet.drone.sim-x1.goto{i}", b"x" * 1024)  # never refused
            st = (await js.stream_info("CMD")).state
            print(f"  CMD discard:old 8 KiB cap: messages={st.messages} first_seq={st.first_seq} last_seq={st.last_seq}")
            self.assertEqual(st.last_seq, 50, "the newest command is always recorded")
            self.assertGreater(st.first_seq, 1, "the oldest records are dropped instead")
            self.assertLessEqual(st.bytes, cap)
        self.run_js(body)

    def test_cmd_discard_new_refuses_when_full(self):
        """Option (a) config: a full discard: new CMD refuses to store the publish."""
        async def body(js):
            cap = 8 * 1024
            await js.add_stream(to_config(by_name("CMD"), max_bytes=cap, discard=jsapi.DiscardPolicy.NEW))
            accepted, refused = 0, None
            for i in range(50):
                try:
                    await js.publish(f"cmd.fleet.drone.sim-x1.goto{i}", b"x" * 1024)
                    accepted += 1
                except APIError as err:
                    refused = err
                    break
            st = (await js.stream_info("CMD")).state
            print(f"  CMD discard:new 8 KiB cap: accepted={accepted} refused err_code={getattr(refused, 'err_code', None)} "
                  f"description={getattr(refused, 'description', None)!r} first_seq={st.first_seq}")
            self.assertIsNotNone(refused, "a full discard:new stream must refuse the publish")
            self.assertEqual(refused.err_code, 10077)  # JSStreamStoreFailedF: maximum bytes exceeded
            self.assertEqual(st.first_seq, 1, "the oldest record is kept")
            self.assertLessEqual(st.bytes, cap)
        self.run_js(body)

    def test_cmd_discard_new_refused_publish_still_reaches_core_subscribers(self):
        """Why discard: new is not fail-closed while bridges subscribe cmd.* on core (A22-P2)."""
        async def body(js):
            nc = js._nc
            await js.add_stream(to_config(by_name("CMD"), max_bytes=400, discard=jsapi.DiscardPolicy.NEW))
            got = []
            async def on_msg(m):
                got.append(m.data)

            sub = await nc.subscribe("cmd.>", cb=on_msg)
            await nc.flush()
            stored = refused = 0
            for i in range(1, 6):
                try:
                    await js.publish("cmd.fleet.drone.sim-x1.hold", f"payload-number-{i}".ljust(100).encode())
                    stored += 1
                except APIError as err:
                    self.assertEqual(err.err_code, 10077)
                    refused += 1
            await nc.flush()
            await asyncio.sleep(0.2)
            await sub.unsubscribe()
            st = (await js.stream_info("CMD")).state
            print(f"  CMD discard:new 400 B cap: stored={stored} refused={refused} "
                  f"stream_messages={st.messages} core_subscriber_received={len(got)}")
            self.assertGreater(refused, 0)
            self.assertEqual(st.messages, stored)
            self.assertEqual(len(got), 5, "core delivery is independent of the stream refusing to store")
        self.run_js(body)

    def test_cmd_update_from_old_to_new_discard(self):
        """The later switch to discard: new (after option a) keeps existing commands."""
        async def body(js):
            await js.add_stream(to_config(by_name("CMD"), max_age=720 * 3600, max_bytes=-1))
            await js.publish("cmd.fleet.drone.sim-x1.rtl", b"{}")
            await js.update_stream(to_config(by_name("CMD"), discard=jsapi.DiscardPolicy.NEW))
            info = await js.stream_info("CMD")
            self.assertEqual(getattr(info.config.discard, "value", info.config.discard), "new")
            self.assertEqual(info.state.messages, 1, "existing commands survive the update")
        self.run_js(body)

    def test_tel_per_subject_history_and_direct_get(self):
        async def body(js):
            await js.add_stream(to_config(by_name("TEL"), max_msgs_per_subject=3))
            for i in range(5):
                await js.publish("tel.feed.nws.status", f'{{"n":{i}}}'.encode())
            await js.publish("tel.fleet.drone.sim-x1.position", b"{}")
            st = (await js.stream_info("TEL")).state
            self.assertEqual(st.messages, 4)
            last = await js.get_last_msg("TEL", "tel.feed.nws.status", direct=True)
            self.assertEqual(last.data, b'{"n":4}')
        self.run_js(body)

    def test_intel_one_hour_dedup(self):
        async def body(js):
            await js.add_stream(to_config(by_name("INTEL")))
            first = await js.publish("intel.seismic.usgs.m4.us7000abcd", b"{}", headers={"Nats-Msg-Id": "usgs-us7000abcd"})
            again = await js.publish("intel.seismic.usgs.m4.us7000abcd", b"{}", headers={"Nats-Msg-Id": "usgs-us7000abcd"})
            self.assertFalse(first.duplicate)
            self.assertTrue(again.duplicate)
            cfg = (await js.stream_info("INTEL")).config
            self.assertEqual(cfg.duplicate_window, 3600)
        self.run_js(body)

    def test_det_captures_detections(self):
        async def body(js):
            await js.add_stream(to_config(by_name("DET")))
            ack = await js.publish("det.air.remoteid.sim-x1", b"{}")
            self.assertEqual(ack.stream, "DET")
        self.run_js(body)


if __name__ == "__main__":
    unittest.main()
