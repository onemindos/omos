"""Tests for scripts/validate_plugin_manifests.py (FEATURES K6).

Acceptance: a new layer (USGS quakes) ships as a plugin with no core change, and a manifest with an
off-grammar subject fails CI.
"""
import json
import subprocess
import sys
import unittest
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))
import validate_plugin_manifests as v  # noqa: E402

SCHEMA = json.loads(v.SCHEMA.read_text())
EXAMPLES = ROOT / "services/plugins/examples"
FIXTURES = ROOT / "scripts/tests/fixtures/plugins"


def load(path):
    return yaml.safe_load(path.read_text())


class Grammar(unittest.TestCase):
    def test_on_grammar(self):
        for s in [
            "tel.drone.x47.position", "tel.fleet.drone.sim-x1.position", "tel.feed.usgs-quakes.status",
            "det.air.visual.track42", "evt.mission.patrol1.completed", "evt.video.sim-x1.online",
            "evt.map.geofence.home.enter", "cmd.drone.x47.intercept", "cmd.fleet.drone.sim-x1.orbit",
            "cmd.map.filter.set", "intel.conflict.acled.event", "ent.ship.mmsi123.state", "ent.*.*.state",
            "svc.hermes.query", "svc.registry.video.get", "svc.meta.task.list", "agents.prompt.hermes.zeus.legacy",
            "tak.cot.inject", "tak.cot.>", "evt.video.>", "tel.fleet.>", "ent.drone.*.state",
        ]:
            self.assertIsNone(v.check_subject(s), s)

    def test_off_grammar(self):
        for s in [
            "fabric.intel.quakes", "intel.quakes.usgs", "omos.air.position", "ent.drone.sim-x1.position",
            "ent.drone.sim-x1", "evt.video.online", "cmd.drone.x47.orbit.now", "tel.drone.x47", "svc.hermes",
            "agents.discover.hermes.zeus.x", "tak.cot.custom", "intel.{domain}.usgs.quake", "intel.>.usgs",
            "det.air.visual.track42.extra", "",
        ]:
            self.assertIsNotNone(v.check_subject(s), s)

    def test_covers(self):
        self.assertTrue(v.covers("cmd.fleet.drone.*.orbit", "cmd.fleet.drone.*.orbit"))
        self.assertTrue(v.covers("evt.video.>", "evt.video.sim-x1.online"))
        self.assertFalse(v.covers("ent.drone.*.state", "ent.event.*.state"))


class Examples(unittest.TestCase):
    def test_every_example_is_valid(self):
        files = sorted(EXAMPLES.glob("*/onemind-plugin.yaml"))
        kinds = set()
        for f in files:
            m = load(f)
            self.assertEqual(v.validate(m, SCHEMA), [], f)
            kinds.update(m["kinds"])
        # One example per kind (K6).
        self.assertEqual(kinds, {"source", "layer", "action", "skill", "video"})

    def test_usgs_layer_needs_no_core_change(self):
        m = load(EXAMPLES / "usgs-quakes/onemind-plugin.yaml")
        self.assertEqual(m["source"]["engine"], "omos-ingest")  # a collector config entry, not new code
        self.assertNotIn("image", m["source"])
        self.assertEqual(m["layer"]["from"], "ent.event.*.state")


    def test_usgs_layer_fields_match_ingest_output(self):
        # The layer must reference fields omos-ingest emits for a USGS feature (#553 geojson
        # parser), or it renders nothing. When the #553 golden is in the tree, check it too.
        layer = load(EXAMPLES / "usgs-quakes/onemind-plugin.yaml")["layer"]
        records = [json.loads((FIXTURES / "usgs-ingest-record.json").read_text())]
        golden = ROOT / "services/bridge/ingest/testdata/golden/geojson-usgs.golden.json"
        if golden.exists():
            records += [r["data"] for r in json.loads(golden.read_text())]
        self.assertEqual(layer["where"], {"kind": "earthquake"})
        self.assertEqual(records[0]["type"], layer["where"]["kind"])
        for rec in records:
            for key in ("sizeField", "colorField", "labelField"):
                field = layer["render"].get(key)
                if field:
                    self.assertIn(field, rec, f"render.{key} {field!r} is not in the ingest record")


class Fixtures(unittest.TestCase):
    def errors(self, name):
        return "\n".join(v.validate(load(FIXTURES / name / "onemind-plugin.yaml"), SCHEMA))

    def test_off_grammar_fails(self):
        e = self.errors("off-grammar")
        self.assertIn("'fabric.intel.quakes' is off-grammar", e)
        self.assertIn("'intel.quakes.usgs' is off-grammar", e)

    def test_reserved_writers_fail(self):
        e = self.errors("reserved-writer")
        self.assertIn("ent.* has one writer", e)
        self.assertIn("cmd.* is published only by the policy service", e)

    def test_reserved_wildcards_fail(self):
        e = self.errors("reserved-wildcard")
        self.assertIn("'evt.*.*.*' is reserved: evt.so.action.*", e)
        self.assertIn("'tak.cot.>' is reserved: only the projector publishes tak.cot.inject", e)
        self.assertNotIn("intel.hazard.example.quake", e)

    def test_reserved_overlap(self):
        why = {s: "\n".join(e for e in self._pub(s)) for s in [
            "evt.*.*.*", "evt.so.*.x", "evt.*.action.x", "evt.so.>", "evt.>", "ent.*.*.*", "cmd.>",
            "tak.cot.inject", "evt.mission.patrol1.completed", "evt.video.*.online", "intel.hazard.usgs.quake"]}
        for s in ["evt.*.*.*", "evt.so.*.x", "evt.*.action.x", "evt.so.>", "evt.>", "ent.*.*.*", "cmd.>",
                  "tak.cot.inject"]:
            self.assertIn("is reserved", why[s], s)
        for s in ["evt.mission.patrol1.completed", "evt.video.*.online", "intel.hazard.usgs.quake"]:
            self.assertNotIn("is reserved", why[s], s)

    def _pub(self, subject):
        m = load(EXAMPLES / "usgs-quakes/onemind-plugin.yaml")
        m["bus"]["publishes"] = [subject, m["source"]["subject"]]
        return v.validate(m, SCHEMA)

    def test_unsafe_action_fails(self):
        e = self.errors("missing-block")
        self.assertIn("licence", e)
        self.assertIn("'human' was expected", e)

    def test_cli_exit_code(self):
        bad = subprocess.run([sys.executable, str(ROOT / "scripts/validate_plugin_manifests.py"),
                              str(FIXTURES / "off-grammar/onemind-plugin.yaml")], capture_output=True, text=True)
        self.assertEqual(bad.returncode, 1, bad.stdout)
        good = subprocess.run([sys.executable, str(ROOT / "scripts/validate_plugin_manifests.py"),
                               str(EXAMPLES / "usgs-quakes/onemind-plugin.yaml")], capture_output=True, text=True)
        self.assertEqual(good.returncode, 0, good.stdout + good.stderr)

    def test_kind_without_block(self):
        m = load(EXAMPLES / "drone-orbit/onemind-plugin.yaml")
        m["kinds"] = ["action", "layer"]
        self.assertIn("kind 'layer' declared but the layer block is missing", "\n".join(v.validate(m, SCHEMA)))

    def test_real_drone_action_refused(self):
        m = load(EXAMPLES / "drone-orbit/onemind-plugin.yaml")
        m["action"]["subject"] = "cmd.fleet.drone.x47.orbit"
        m["bus"]["subscribes"] = ["cmd.fleet.drone.x47.orbit"]
        self.assertIn("only simulated assets", "\n".join(v.validate(m, SCHEMA)))


if __name__ == "__main__":
    unittest.main()
