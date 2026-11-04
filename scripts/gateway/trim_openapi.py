#!/usr/bin/env python3
"""Trim an upstream OpenAPI 3 document to an allow list of operations and write it as a ConfigMap.

agentgateway turns every operation of an OpenAPI target into an MCP tool and cannot filter them, so the spec
that is mounted into the gateway pod is trimmed in Git: only the operations we mean to expose (read, write and
the destructive ones the execute route will gate) are kept, and unreferenced schemas are dropped.

  python scripts/gateway/trim_openapi.py --name mediamtx --source /path/to/upstream-openapi.yaml \\
      --upstream "bluenviron/mediamtx v1.21.2 api/openapi.yaml" --server http://onetak-media.onetak.svc.cluster.local:4000 \\
      --classification platform/tool-gateway/policy/classification/mediamtx.yaml \\
      --out platform/tool-gateway/manifests/openapi-mediamtx.configmap.yaml

The kept operations are exactly the names in the classification file (every class), so the spec and the policy
cannot drift: scripts/tests/test_gitops_tool_gateway_openapi.py checks both ways.

The ConfigMap `gateway-openapi-<name>` carries the spec under key `openapi.yaml`; platform/tool-gateway/values.yaml
mounts it into the gateway pod at /specs/<name>/. Review the diff of the operation list like any other policy change.
"""
from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

import yaml

METHODS = ("get", "put", "post", "delete", "patch", "head", "options")
REF_RE = re.compile(r"^#/components/([A-Za-z]+)/(.+)$")


def refs(node, out: set[tuple[str, str]]) -> None:
    if isinstance(node, dict):
        ref = node.get("$ref")
        if isinstance(ref, str):
            m = REF_RE.match(ref)
            if m:
                out.add((m.group(1), m.group(2)))
        for v in node.values():
            refs(v, out)
    elif isinstance(node, list):
        for v in node:
            refs(v, out)


def trim(spec: dict, ops: list[str], server: str) -> dict:
    wanted = set(ops)
    paths: dict = {}
    found: set[str] = set()
    for path, item in spec["paths"].items():
        kept = {}
        for key, val in item.items():
            if key in METHODS and val.get("operationId") in wanted:
                # agentgateway builds the tool description from `description` only; many specs fill just `summary`.
                if not (val.get("description") or "").strip() and val.get("summary"):
                    val = {**val, "description": val["summary"]}
                kept[key] = val
                found.add(val["operationId"])
            elif key not in METHODS:  # path-level parameters
                kept[key] = val
        if any(k in METHODS for k in kept):
            paths[path] = kept
    missing = wanted - found
    if missing:
        raise SystemExit(f"operations not in the source spec: {sorted(missing)}")

    components = spec.get("components", {})
    needed: set[tuple[str, str]] = set()
    refs(paths, needed)
    todo = set(needed)
    while todo:  # transitive closure over components
        kind, name = todo.pop()
        sub: set[tuple[str, str]] = set()
        refs(components.get(kind, {}).get(name), sub)
        for r in sub - needed:
            needed.add(r)
            todo.add(r)
    kept_components: dict = {}
    for kind, name in sorted(needed):
        kept_components.setdefault(kind, {})[name] = components[kind][name]

    out = {
        "openapi": spec["openapi"],
        "info": spec["info"],
        "servers": [{"url": server}],
        "security": [],
        "paths": paths,
    }
    if kept_components:
        out["components"] = kept_components
    return out


def main() -> None:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--name", required=True, help="target name; ConfigMap is gateway-openapi-<name>")
    ap.add_argument("--source", required=True, type=Path, help="upstream OpenAPI file (yaml or json)")
    ap.add_argument("--upstream", required=True, help="provenance, e.g. 'org/repo vX.Y.Z path'")
    ap.add_argument("--server", required=True, help="in-cluster base URL written into servers[0]")
    ap.add_argument("--classification", required=True, type=Path,
                    help="classification file; its exact tool names (all classes) are the operationIds to keep")
    ap.add_argument("--out", required=True, type=Path)
    a = ap.parse_args()

    raw = a.source.read_bytes()
    spec = yaml.safe_load(raw)
    cls = yaml.safe_load(a.classification.read_text())
    ops = sorted({n for k in ("read", "write", "destructive", "physical") for n in (cls.get(k) or [])})
    if any("*" in o for o in ops):
        raise SystemExit("OpenAPI classification entries must be exact operationIds, not globs")
    trimmed = trim(spec, ops, a.server)
    text = yaml.safe_dump(trimmed, sort_keys=False, width=120)
    cm = {
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "name": f"gateway-openapi-{a.name}",
            "namespace": "tool-gateway",
            "annotations": {
                "argocd.argoproj.io/sync-wave": "-1",
                "onemind.io/openapi-source": a.upstream,
            },
        },
        "data": {"openapi.yaml": text},
    }

    class Block(str):
        pass

    def block(dumper, data):
        return dumper.represent_scalar("tag:yaml.org,2002:str", data, style="|")

    yaml.add_representer(Block, block)
    cm["data"]["openapi.yaml"] = Block(text)
    header = (f"# GENERATED by scripts/gateway/trim_openapi.py: {a.upstream}, trimmed to {len(ops)} operations.\n"
              "# Do not edit by hand; regenerate and review the operation list.\n")
    a.out.write_text(header + yaml.dump(cm, sort_keys=False, width=1000))
    print(f"{a.out}: {len(ops)} operations, {len(text)} bytes", file=sys.stderr)
    print(json.dumps(sorted(ops)))


if __name__ == "__main__":
    main()
