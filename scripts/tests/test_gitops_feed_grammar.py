"""GitOps contracts for the feed grammar cut-over (FEATURES A1, decision 0004 #18).

Static checks only; the behavioural replay against nats-server is
platform/data/redpanda-connect/test/replay.sh (needs nats-server, nats and redpanda-connect).
"""
from __future__ import annotations

import re
import unittest
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
RPC = ROOT / "platform/data/redpanda-connect"
GRAMMAR = ("feed-mta-grammar.yaml", "feed-satellites-grammar.yaml")
# BUS-DESIGN §3.1 subjects the grammar pipelines may publish (interpolations normalised to {id}).
ALLOWED = {
    "feed-mta-grammar.yaml": {"evt.world.rail.{id}", "tel.feed.mta.status", "evt.sys.mta.rejected"},
    "feed-satellites-grammar.yaml": {"intel.space.celestrak.tle", "tel.feed.celestrak.status", "evt.sys.celestrak.rejected"},
}


def load_docs(path: Path) -> list[dict]:
    return [d for d in yaml.safe_load_all(path.read_text()) if d]


def outputs(node, found=None):
    """Every leaf nats/nats_jetstream output under a pipeline's output section."""
    found = [] if found is None else found
    if isinstance(node, dict):
        for kind in ("nats", "nats_jetstream"):
            if kind in node and isinstance(node[kind], dict) and "subject" in node[kind]:
                found.append((kind, node[kind]))
        for v in node.values():
            outputs(v, found)
    elif isinstance(node, list):
        for v in node:
            outputs(v, found)
    return found


class FeedGrammar(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        docs = load_docs(RPC / "deployment.yaml") + load_docs(RPC / "deployment-grammar.yaml")
        cls.cm = next(d for d in docs if d["kind"] == "ConfigMap" and d["metadata"]["name"] == "redpanda-connect-pipelines")["data"]
        cls.deploys = {d["metadata"]["name"]: d for d in docs if d["kind"] == "Deployment"}
        cls.sat = {d["metadata"]["name"]: d for d in load_docs(ROOT / "platform/omos/feeds.yaml") if d["kind"] == "Deployment"}

    def args(self, name: str) -> list[str]:
        return self.deploys[name]["spec"]["template"]["spec"]["containers"][0]["args"]

    def test_duplicate_feed_hazards_is_removed(self):
        self.assertNotIn("feed-hazards.yaml", self.cm)
        self.assertFalse((RPC / "pipelines/feed-hazards.yaml").exists())
        for name in self.deploys:
            self.assertFalse(any("feed-hazards" in a for a in self.args(name)), name)

    def test_aircraft_and_ships_stay_off(self):
        for name in self.deploys:
            self.assertFalse(any(re.search(r"feed-(aircraft|ships)", a) for a in self.args(name)), name)

    def test_reference_copies_match_the_configmap(self):
        for key in ("feed-mta.yaml", *GRAMMAR):
            self.assertEqual((RPC / "pipelines" / key).read_text(), self.cm[key], key)

    def test_grammar_pipelines_publish_only_grammar_subjects(self):
        for key in GRAMMAR:
            cfg = yaml.safe_load(self.cm[key])
            self.assertNotIn("tak.cot", yaml.safe_dump(cfg["output"]), key)
            outs = outputs(cfg["output"])
            subjects = {re.sub(r"\$\{!.*?\}", "{id}", o["subject"]) for _, o in outs}
            self.assertEqual(subjects, ALLOWED[key], key)
            for kind, o in outs:
                if kind == "nats_jetstream":  # stored records are acked and deduplicated
                    self.assertIn("Nats-Msg-Id", o.get("headers", {}), (key, o["subject"]))

    def test_each_pipeline_runs_in_at_most_one_deployment(self):
        seen: dict[str, str] = {}
        for name in self.deploys:
            for a in self.args(name):
                if a.endswith(".yaml"):
                    self.assertNotIn(a, seen, f"{a} in {seen.get(a)} and {name}")
                    seen[a] = name
        self.assertEqual(self.args("redpanda-connect-grammar")[2:], [f"/config/pipelines/{k}" for k in GRAMMAR])

    def test_grammar_deployment_is_a_reviewed_switch(self):
        d = self.deploys["redpanda-connect-grammar"]
        spec = d["spec"]["template"]["spec"]
        self.assertIn(d["spec"]["replicas"], (0, 1))
        self.assertEqual(d["spec"]["strategy"]["type"], "Recreate")
        self.assertEqual(spec["nodeSelector"], {"onemind.io/compute-role": "ecosystem"})
        c = spec["containers"][0]
        self.assertEqual(c["image"], self.deploys["redpanda-connect"]["spec"]["template"]["spec"]["containers"][0]["image"])
        self.assertIn("@sha256:", c["image"])
        self.assertTrue(c["resources"]["requests"] and c["resources"]["limits"])
        self.assertNotIn("AISSTREAM_API_KEY", yaml.safe_dump(c.get("env", [])))

    def test_legacy_cot_outputs_go_off_only_after_the_new_path_runs(self):
        """Stage order from CUTOVER.md: never turn a legacy CoT output off while its grammar
        replacement is off, or the feed disappears from the bus entirely."""
        grammar_on = self.deploys["redpanda-connect-grammar"]["spec"]["replicas"] >= 1
        mta_legacy_on = "/config/pipelines/feed-mta.yaml" in self.args("redpanda-connect") and \
            self.deploys["redpanda-connect"]["spec"].get("replicas", 1) >= 1
        sat_legacy_on = self.sat["feed-satellites"]["spec"].get("replicas", 1) >= 1
        if not grammar_on:
            self.assertTrue(mta_legacy_on, "feed-mta CoT turned off before feed-mta-grammar runs")
            self.assertTrue(sat_legacy_on, "feed-satellites turned off before feed-satellites-grammar runs")
        self.assertIn("Stage 1", (RPC / "CUTOVER.md").read_text())


if __name__ == "__main__":
    unittest.main()
