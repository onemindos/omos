#!/usr/bin/env python3
"""Generate the feed ConfigMaps from Git sources.

Reads platform/data/sources/*.yaml, validates each against
services/schemas/feeds/source-manifest.v1.schema.json, and writes:

  platform/data/redpanda-connect/feed-sources.configmap.yaml
      ConfigMap feed-sources (namespace feeds): one key <id>.yaml per redpanda-connect source
      (the manifest verbatim) plus _index.json.
  platform/data/redpanda-connect/feed-control.configmap.yaml
      ConfigMap feed-control-pipelines (namespace feeds): the controller and seeder streams,
      copied from platform/data/redpanda-connect/pipelines/.
  platform/data/redpanda-connect/deployment-grammar.yaml
      Only the pod annotation onemind.io/feed-config-sha256: a sha256 of the data of both
      ConfigMaps. The controller reads _index.json when its stream starts, so a manifest change
      must restart the pod; the changed annotation does that.

Usage:
  python scripts/gen_feed_sources.py           write both files
  python scripts/gen_feed_sources.py --check   exit 1 if either file is stale

Needs PyYAML and jsonschema.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from pathlib import Path

import jsonschema
import yaml

ROOT = Path(__file__).resolve().parent.parent
SOURCES_DIR = ROOT / "platform" / "data" / "sources"
SCHEMA_PATH = ROOT / "services" / "schemas" / "feeds" / "source-manifest.v1.schema.json"
RC_DIR = ROOT / "platform" / "data" / "redpanda-connect"
SOURCES_CM = RC_DIR / "feed-sources.configmap.yaml"
CONTROL_CM = RC_DIR / "feed-control.configmap.yaml"
GRAMMAR_DEPLOY = RC_DIR / "deployment-grammar.yaml"
HASH_LINE = re.compile(r'^(\s*onemind\.io/feed-config-sha256: )"[0-9a-f]{64}"$', re.M)
CONTROL_FILES = ("feed-controller.yaml", "feed-seeder.yaml")

SOURCES_HEADER = """\
# GENERATED FILE. Do not edit by hand.
# Source: platform/data/sources/*.yaml
# Regenerate: python scripts/gen_feed_sources.py   (CI runs it with --check)
"""

CONTROL_HEADER = """\
# GENERATED FILE. Do not edit by hand.
# Source: platform/data/redpanda-connect/pipelines/{feed-controller,feed-seeder}.yaml
# Regenerate: python scripts/gen_feed_sources.py   (CI runs it with --check)
# The Argo Application excludes pipelines/, so this ConfigMap is what actually deploys.
"""


def load_schema() -> dict:
    return json.loads(SCHEMA_PATH.read_text())


def load_sources() -> list[tuple[Path, dict]]:
    """Return (path, manifest) for every source file, validated, sorted by id."""
    validator = jsonschema.Draft202012Validator(load_schema())
    out = []
    for path in sorted(SOURCES_DIR.glob("*.yaml")):
        manifest = yaml.safe_load(path.read_text())
        errors = sorted(validator.iter_errors(manifest), key=lambda e: list(e.path))
        if errors:
            msg = "; ".join(f"{'/'.join(map(str, e.path)) or '<root>'}: {e.message[:200]}" for e in errors)
            raise SystemExit(f"{path.relative_to(ROOT)}: schema errors: {msg}")
        if manifest["id"] != path.stem:
            raise SystemExit(f"{path.relative_to(ROOT)}: id {manifest['id']!r} must equal the file name")
        out.append((path, manifest))
    out.sort(key=lambda pm: pm[1]["id"])
    return out


def block(key: str, text: str) -> str:
    if not text.endswith("\n"):
        text += "\n"
    lines = [f"  {key}: |"]
    for line in text.splitlines():
        lines.append(("    " + line) if line else "")
    return "\n".join(lines) + "\n"


def configmap(header: str, name: str, items: list[tuple[str, str]]) -> str:
    head = (
        "---\n"
        "apiVersion: v1\n"
        "kind: ConfigMap\n"
        "metadata:\n"
        f"  name: {name}\n"
        "  namespace: feeds\n"
        "  labels:\n"
        "    app.kubernetes.io/part-of: omos-feeds\n"
        "data:\n"
    )
    return header + head + "".join(block(k, v) for k, v in sorted(items))


def render_sources() -> str:
    items = []
    index = []
    for path, manifest in load_sources():
        index.append(
            {
                "id": manifest["id"],
                "runner": manifest["runner"],
                "default_enabled": manifest["default_enabled"],
                "interval": manifest.get("schedule", {}).get("interval"),
            }
        )
        if manifest["runner"] == "redpanda-connect":
            items.append((f"{manifest['id']}.yaml", path.read_text()))
    index.sort(key=lambda s: s["id"])
    items.append(("_index.json", json.dumps(index, separators=(",", ":"), sort_keys=True) + "\n"))
    return configmap(SOURCES_HEADER, "feed-sources", items)


def render_control() -> str:
    items = [(n, (RC_DIR / "pipelines" / n).read_text()) for n in CONTROL_FILES]
    return configmap(CONTROL_HEADER, "feed-control-pipelines", items)


def config_hash(sources_text: str, control_text: str) -> str:
    """sha256 over the data of both ConfigMaps (not their comments), in a stable order."""
    data = [yaml.safe_load(t)["data"] for t in (sources_text, control_text)]
    return hashlib.sha256(json.dumps(data, sort_keys=True, separators=(",", ":")).encode()).hexdigest()


def render_grammar_deployment(digest: str) -> str:
    text = GRAMMAR_DEPLOY.read_text()
    if not HASH_LINE.search(text):
        raise SystemExit(f"{GRAMMAR_DEPLOY.relative_to(ROOT)}: no onemind.io/feed-config-sha256 annotation line")
    return HASH_LINE.sub(lambda m: f'{m.group(1)}"{digest}"', text)


def outputs() -> list[tuple[Path, str]]:
    sources, control = render_sources(), render_control()
    return [
        (SOURCES_CM, sources),
        (CONTROL_CM, control),
        (GRAMMAR_DEPLOY, render_grammar_deployment(config_hash(sources, control))),
    ]


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__.split("\n\n")[0])
    ap.add_argument("--check", action="store_true", help="exit 1 if a generated file is stale")
    args = ap.parse_args()
    stale = []
    for path, text in outputs():
        current = path.read_text() if path.exists() else None
        if args.check:
            if current != text:
                stale.append(path)
        elif current != text:
            path.write_text(text)
            print(f"wrote {path.relative_to(ROOT)}")
    if stale:
        for p in stale:
            print(f"stale: {p.relative_to(ROOT)} (run python scripts/gen_feed_sources.py)", file=sys.stderr)
        return 1
    if args.check:
        print("feed ConfigMaps and the pod config hash are up to date")
    return 0


if __name__ == "__main__":
    sys.exit(main())
