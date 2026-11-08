"""Contracts for the feed source manifests and the feed-control ConfigMaps.

Needs PyYAML and jsonschema (pinned in the Helm render check job of quality.yml).
"""
from __future__ import annotations

import json
import re
import subprocess
import sys
import unittest
from pathlib import Path

import jsonschema
import yaml

ROOT = Path(__file__).resolve().parents[2]
SOURCES_DIR = ROOT / "platform" / "data" / "sources"
SCHEMA_PATH = ROOT / "services" / "schemas" / "feeds" / "source-manifest.v1.schema.json"
RC_DIR = ROOT / "platform" / "data" / "redpanda-connect"
SOURCES_CM = RC_DIR / "feed-sources.configmap.yaml"
CONTROL_CM = RC_DIR / "feed-control.configmap.yaml"
NATS_OUTPUTS = ("nats", "nats_jetstream")
# Variables the engine pod provides to every stream; they are not source credentials.
ENGINE_VARS = {"NATS_URL", "NATS_USER", "NATS_PASSWORD", "NATS_TLS_ENABLED"}
INTERP = re.compile(r"\$\{!.*?\}", re.S)
ENV_REF = re.compile(r"\$\{([A-Za-z_][A-Za-z0-9_]*)(:[^}]*)?\}")
SECRET_LOOKING = re.compile(
    r"(?i)(api[_-]?key|token|secret|password|passwd|authorization|bearer)\"?\s*[:=]\s*[\"']?(?!\$\{)[A-Za-z0-9+/_-]{12,}"
)
KEY_BLOBS = re.compile(r"(AKIA[0-9A-Z]{16}|-----BEGIN [A-Z ]*PRIVATE KEY-----|ghp_[A-Za-z0-9]{20,}|xox[bp]-[A-Za-z0-9-]{10,})")


def load_manifests():
    return [(p, yaml.safe_load(p.read_text())) for p in sorted(SOURCES_DIR.glob("*.yaml"))]


def walk(node):
    """Yield (key, value) for every mapping entry in a nested structure."""
    if isinstance(node, dict):
        for k, v in node.items():
            yield k, v
            yield from walk(v)
    elif isinstance(node, list):
        for item in node:
            yield from walk(item)


def nats_subjects(stream):
    """Every `subject` set on a nats or nats_jetstream output in the stream config."""
    out = []
    for key, value in walk(stream):
        if key in NATS_OUTPUTS and isinstance(value, dict) and "subject" in value:
            out.append(value["subject"])
    return out


def subject_matches(pattern, subject):
    """Token match. `${! ... }` counts as one token; {x} and * in the pattern match one token;
    a trailing > matches one or more."""
    subject = INTERP.sub("X", subject)
    pt, st = pattern.split("."), subject.split(".")
    for i, p in enumerate(pt):
        if p == ">":
            return len(st) > i
        if i >= len(st):
            return False
        if (p.startswith("{") and p.endswith("}")) or p == "*":
            if "." in st[i]:
                return False
            continue
        if p != st[i]:
            return False
    return len(st) == len(pt)


class SubjectMatcherTest(unittest.TestCase):
    def test_matcher(self):
        self.assertTrue(subject_matches("evt.world.rail.{entity}", "evt.world.rail.${! @rail_id }"))
        self.assertTrue(subject_matches("evt.world.rail.{entity}", "evt.world.rail.abc"))
        self.assertFalse(subject_matches("evt.world.rail.{entity}", "evt.world.air.abc"))
        self.assertFalse(subject_matches("evt.world.rail.{entity}", "evt.world.rail.a.b"))
        self.assertFalse(subject_matches("evt.world.rail.{entity}", "tak.cot.inject"))
        self.assertTrue(subject_matches("intel.>", "intel.space.celestrak.tle"))


class FeedSourceManifestTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.schema = json.loads(SCHEMA_PATH.read_text())
        cls.validator = jsonschema.Draft202012Validator(cls.schema)
        cls.manifests = load_manifests()

    def test_schema_itself_is_valid(self):
        jsonschema.Draft202012Validator.check_schema(self.schema)

    def test_there_is_at_least_one_manifest(self):
        self.assertTrue(self.manifests)

    def test_manifests_validate_and_id_matches_filename(self):
        for path, manifest in self.manifests:
            with self.subTest(source=path.name):
                errors = [f"{list(e.path)}: {e.message[:160]}" for e in self.validator.iter_errors(manifest)]
                self.assertEqual(errors, [])
                self.assertEqual(manifest["id"], path.stem)

    def test_schema_rejects_bad_manifests(self):
        good = dict(self.manifests[0][1])
        bad_cases = {
            "ent subject": {"subject": "ent.ship.x.state"},
            "tak subject": {"subject": "tak.cot.inject"},
            "unknown field": {"surprise": 1},
            "bad id": {"id": "Bad_ID"},
            "bad interval": {"schedule": {"interval": "soon"}},
            "bad runner": {"runner": "cron"},
        }
        for name, change in bad_cases.items():
            with self.subTest(case=name):
                self.assertTrue(list(self.validator.iter_errors({**good, **change})))
        no_stream = {k: v for k, v in good.items() if k != "stream"}
        self.assertTrue(list(self.validator.iter_errors(no_stream)), "redpanda-connect needs stream")
        no_interval = {k: v for k, v in good.items() if k != "schedule"}
        self.assertTrue(list(self.validator.iter_errors(no_interval)), "poll needs schedule.interval")

    def test_ids_are_unique(self):
        ids = [m["id"] for _, m in self.manifests]
        self.assertEqual(len(ids), len(set(ids)))

    def test_fallback_points_at_a_known_source(self):
        ids = {m["id"] for _, m in self.manifests}
        for path, m in self.manifests:
            if "fallback" in m:
                with self.subTest(source=path.name):
                    self.assertIn(m["fallback"], ids)

    def test_stream_outputs_stay_inside_the_declared_subject(self):
        for path, m in self.manifests:
            if m["runner"] != "redpanda-connect":
                continue
            names = {m["id"], m["id"].split("-")[0]}
            allowed_extra = [f"tel.feed.{n}.status" for n in names] + [f"evt.sys.{n}.rejected" for n in names]
            subjects = nats_subjects(m["stream"])
            with self.subTest(source=path.name):
                self.assertTrue(subjects, "stream publishes nowhere")
                for s in subjects:
                    self.assertNotIn("tak.cot.inject", s)
                    self.assertFalse(s.startswith(("tak.", "ent.")), s)
                    self.assertTrue(
                        subject_matches(m["subject"], s) or s in allowed_extra,
                        f"{s} is outside {m['subject']} and not a status/rejected subject",
                    )

    def test_no_tak_cot_inject_anywhere_in_a_stream(self):
        for path, m in self.manifests:
            with self.subTest(source=path.name):
                self.assertNotIn("tak.cot.inject", json.dumps(m.get("stream", {})))

    def test_interval_placeholder_is_used_for_poll_sources(self):
        for path, m in self.manifests:
            if m["runner"] == "redpanda-connect" and m["mode"] == "poll":
                with self.subTest(source=path.name):
                    self.assertIn("__INTERVAL__", json.dumps(m["stream"]))

    def test_credentials_are_names_only(self):
        for path, m in self.manifests:
            if m["runner"] != "redpanda-connect":
                continue
            text = path.read_text()
            keys = {c["key"] for c in m.get("credentials", [])}
            with self.subTest(source=path.name):
                self.assertIsNone(KEY_BLOBS.search(text))
                self.assertIsNone(SECRET_LOOKING.search(text), "literal secret-looking value in the file")
                for name, default in ENV_REF.findall(json.dumps(m["stream"])):
                    if name in ENGINE_VARS or default:
                        continue  # engine-provided, or has a non-secret default
                    self.assertIn(name, keys, f"${{{name}}} is not listed in credentials[].key")
                for name in keys:
                    self.assertIn("${" + name, json.dumps(m["stream"]), f"credential {name} is never used")


class FeedControlConfigMapTest(unittest.TestCase):
    def test_generated_sources_configmap_is_current(self):
        result = subprocess.run(
            [sys.executable, str(ROOT / "scripts" / "gen_feed_sources.py"), "--check"],
            capture_output=True,
            text=True,
        )
        self.assertEqual(result.returncode, 0, result.stderr)

    def test_sources_configmap_content(self):
        cm = yaml.safe_load(SOURCES_CM.read_text())
        self.assertEqual(cm["kind"], "ConfigMap")
        self.assertEqual(cm["metadata"]["name"], "feed-sources")
        self.assertEqual(cm["metadata"]["namespace"], "feeds")
        index = json.loads(cm["data"]["_index.json"])
        rc_ids = set()
        for path, m in load_manifests():
            entry = next(e for e in index if e["id"] == m["id"])
            self.assertEqual(entry["default_enabled"], m["default_enabled"])
            self.assertEqual(entry["interval"], m.get("schedule", {}).get("interval"))
            if m["runner"] == "redpanda-connect":
                rc_ids.add(m["id"])
                self.assertEqual(cm["data"][f"{m['id']}.yaml"], path.read_text())
        self.assertEqual({k[:-5] for k in cm["data"] if k != "_index.json"}, rc_ids)
        self.assertEqual([e["id"] for e in index], sorted(e["id"] for e in index))

    def test_control_configmap_equals_pipeline_files(self):
        cm = yaml.safe_load(CONTROL_CM.read_text())
        self.assertEqual(cm["metadata"]["name"], "feed-control-pipelines")
        self.assertEqual(cm["metadata"]["namespace"], "feeds")
        self.assertEqual(sorted(cm["data"]), ["feed-controller.yaml", "feed-seeder.yaml"])
        for name, inline in cm["data"].items():
            with self.subTest(file=name):
                self.assertEqual(inline, (RC_DIR / "pipelines" / name).read_text())

    def test_controller_guards(self):
        ctl = yaml.safe_load((RC_DIR / "pipelines" / "feed-controller.yaml").read_text())
        self.assertEqual(ctl["pipeline"]["threads"], 1)
        self.assertEqual(ctl["input"]["nats_kv"]["bucket"], "feed-control")
        text = (RC_DIR / "pipelines" / "feed-controller.yaml").read_text()
        self.assertIn("_index.json", text)
        self.assertIn("ignore", text)
        self.assertNotIn("tak.cot", text)

    def test_seeder_never_overwrites(self):
        text = (RC_DIR / "pipelines" / "feed-seeder.yaml").read_text()
        seeder = yaml.safe_load(text)
        ops = [p["nats_kv"]["operation"] for p in seeder["pipeline"]["processors"] if "nats_kv" in p]
        self.assertEqual(ops, ["create"])

    def test_kv_cr(self):
        cr = yaml.safe_load((ROOT / "platform" / "nats" / "streams" / "kv-feed-control.yaml").read_text())
        spec = cr["spec"]
        self.assertEqual(cr["kind"], "KeyValue")
        self.assertEqual((spec["bucket"], spec["storage"], spec["replicas"]), ("feed-control", "file", 3))
        self.assertTrue(spec["preventDelete"])
        self.assertNotIn("servers", spec)
        app = (ROOT / "platform" / "argocd" / "platform-nats-streams.yaml").read_text()
        self.assertNotIn("kv-feed-control", app, "the new CR must not be excluded from the Application")


if __name__ == "__main__":
    unittest.main()
