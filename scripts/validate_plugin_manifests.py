"""Validate OneMind plugin manifests (onemind-plugin.yaml, FEATURES K6).

Checks, per manifest:
  * structure against services/schemas/plugins/onemind-plugin.v1.schema.json;
  * every bus subject against the locked NATS grammar (CLAUDE.md "NATS Subject Grammar", plus the
    doctrine forms docs/NATS-DOCTRINE.md section 9 adds: tel.fleet.{class}.{id}.{metric},
    cmd.fleet.{class}.{id}.{verb}, cmd.map.*, evt.map.geofence.{id}.{event}, tak.cot.*);
  * reserved writers: a plugin never publishes ent.* (entity service), cmd.* (policy service, after
    human confirm), evt.so.action.* (audit) or tak.cot.inject / tak.cot.mesh.tx (projector);
  * kinds and blocks agree, and each block's subject is declared in `bus`;
  * community plugins are declarative (no image, no skill).

Usage:
  python scripts/validate_plugin_manifests.py            every tracked onemind-plugin.yaml
  python scripts/validate_plugin_manifests.py FILE...    the given files
Fixtures under scripts/tests/fixtures/plugins/ are invalid on purpose and are skipped in the
default mode; scripts/tests/test_plugin_manifests.py asserts that they fail.
"""
from __future__ import annotations

import json
import re
import subprocess
import sys
from pathlib import Path

import yaml
from jsonschema import Draft202012Validator

ROOT = Path(__file__).resolve().parents[1]
SCHEMA = ROOT / "services/schemas/plugins/onemind-plugin.v1.schema.json"
FIXTURES = "scripts/tests/fixtures/"

TOKEN = re.compile(r"^[A-Za-z0-9_-]+$")
AGENT_VERBS = {"prompt", "hb", "status", "attachments"}
TAK_SUBJECTS = {"tak.cot.inject", "tak.cot.out", "tak.cot.mesh.rx", "tak.cot.mesh.tx", "tak.cot.>"}
RESERVED_PUBLISH = (
    ("ent.", "ent.* has one writer, the entity service"),
    ("cmd.", "cmd.* is published only by the policy service after policy, audit and human confirm"),
    ("evt.so.action.", "evt.so.action.* is the policy service's audit record"),
    ("tak.cot.inject", "only the projector publishes tak.cot.inject"),
    ("tak.cot.mesh.tx", "only the site projector publishes tak.cot.mesh.tx"),
)
VIDEO_CLASS_TYPE = {"drone": "drone", "robot": "robot", "cam": "camera"}


def _lengths(root: str, toks: list[str]) -> set[int] | None:
    """Allowed token counts for a root, given the second and third tokens."""
    second = toks[1] if len(toks) > 1 else ""
    if root == "tel":
        return {4, 5} if second in ("fleet", "*") else {4}
    if root in ("det", "intel"):
        return {4}
    if root == "evt":
        return {4, 5} if toks[1:3] in (["map", "geofence"], ["*", "*"]) else {4}
    if root == "cmd":
        return {4, 5} if second in ("fleet", "map", "*") else {4}
    if root == "ent":
        return {4}
    return None


def check_subject(subject: str) -> str | None:
    """Return why a subject (or subscription pattern) is outside the grammar, or None."""
    toks = subject.split(".")
    if any(t == "" for t in toks):
        return "empty token"
    for i, t in enumerate(toks):
        if t == ">":
            if i != len(toks) - 1:
                return "'>' must be the last token"
        elif t != "*" and not TOKEN.match(t):
            return f"token {t!r} is not a plain token (use '*' for a variable part, never {{id}})"
    root, n = toks[0], len(toks)
    tail = toks[-1] == ">"
    if root == "tak":
        return None if subject in TAK_SUBJECTS else f"tak subjects are limited to {sorted(TAK_SUBJECTS)}"
    if root == "svc":
        return None if (n >= 3 or (tail and n >= 2)) else "svc.{service}.{endpoint}"
    if root == "agents":
        if n >= 2 and toks[1] not in AGENT_VERBS | {"*", ">"}:
            return f"agents verb must be one of {sorted(AGENT_VERBS)}"
        return None if (n >= 5 or tail) else "agents.{verb}.{agent}.{owner}.*"
    allowed = _lengths(root, toks)
    if allowed is None:
        return f"root {root!r} is not in the locked grammar (tel det evt cmd intel ent svc agents tak)"
    if tail:
        return None if n - 1 < max(allowed) else f"{root}.* has at most {max(allowed)} tokens"
    if n not in allowed:
        forms = {
            "tel": "tel.{node}.{class}.{metric} or tel.fleet.{class}.{id}.{metric}",
            "det": "det.{domain}.{type}.{id}",
            "intel": "intel.{domain}.{source}.{type}",
            "evt": "evt.{domain}.{entity}.{event}",
            "cmd": "cmd.{domain}.{entity}.{action} or cmd.fleet.{class}.{id}.{verb}",
            "ent": "ent.{entity_type}.{entity_id}.state",
        }
        return f"{n} tokens; the grammar is {forms[root]}"
    if root == "ent" and toks[3] not in ("state", "*"):
        return "ent.{entity_type}.{entity_id}.state"
    return None


def covers(pattern: str, subject: str) -> bool:
    """True if NATS pattern `pattern` matches every subject `subject` can match."""
    p, s = pattern.split("."), subject.split(".")
    for i, tok in enumerate(p):
        if tok == ">":
            return len(s) > i
        if i >= len(s):
            return False
        if tok != "*" and tok != s[i]:
            return False
    return len(p) == len(s)


def _declared(bus: dict, keys: tuple[str, ...], subject: str) -> bool:
    return any(covers(p, subject) for k in keys for p in bus.get(k, []) or [])


def validate(manifest: dict, schema: dict) -> list[str]:
    errors = [f"schema: {'/'.join(map(str, e.absolute_path)) or '(root)'}: {e.message}"
              for e in sorted(Draft202012Validator(schema).iter_errors(manifest), key=str)]
    if errors:
        return errors
    bus = manifest.get("bus", {})
    for key in ("publishes", "subscribes", "requests", "serves"):
        for s in bus.get(key, []) or []:
            why = check_subject(s)
            if why:
                errors.append(f"bus.{key}: {s!r} is off-grammar: {why}")
            if key in ("requests", "serves") and not s.startswith("svc."):
                errors.append(f"bus.{key}: {s!r} must be a svc.* RPC subject")
    for s in bus.get("publishes", []) or []:
        for prefix, why in RESERVED_PUBLISH:
            if s.startswith(prefix) or (s.endswith(">") and prefix.startswith(s[:-1])):
                errors.append(f"bus.publishes: {s!r} is reserved: {why}")
    kinds = set(manifest["kinds"])
    for kind in ("source", "layer", "action", "skill", "video"):
        if kind in kinds and not manifest.get(kind):
            errors.append(f"kind {kind!r} declared but the {kind} block is missing")
        if manifest.get(kind) and kind not in kinds:
            errors.append(f"{kind} block present but {kind!r} is not in kinds")
    src = manifest.get("source")
    if src:
        if not _declared(bus, ("publishes",), src["subject"]):
            errors.append(f"source.subject {src['subject']!r} is not in bus.publishes")
        if src["engine"] == "image" and not src.get("image"):
            errors.append("source.engine image needs source.image")
        if src["engine"] == "omos-ingest" and not (src.get("parser") and src.get("url")):
            errors.append("source.engine omos-ingest needs parser and url")
    layer = manifest.get("layer")
    if layer:
        why = check_subject(layer["from"])
        if why:
            errors.append(f"layer.from {layer['from']!r} is off-grammar: {why}")
        elif not _declared(bus, ("subscribes",), layer["from"]):
            errors.append(f"layer.from {layer['from']!r} is not in bus.subscribes")
    action = manifest.get("action")
    if action:
        subj = action["subject"]
        why = check_subject(subj)
        if why or not subj.startswith("cmd."):
            errors.append(f"action.subject {subj!r} must be an on-grammar cmd.* subject ({why or 'not cmd.*'})")
        elif not _declared(bus, ("subscribes",), subj):
            errors.append(f"action.subject {subj!r} is not in bus.subscribes (the plugin consumes approved commands)")
        sid = subj.split(".")
        if len(sid) == 5 and sid[1] == "fleet" and sid[3] not in ("*",) and not sid[3].startswith("sim-"):
            errors.append(f"action.subject {subj!r}: only simulated assets (sim-*) in this run")
    skill = manifest.get("skill")
    if skill:
        for tool in skill["tools"]:
            rpc = tool["rpc"]
            if not rpc.startswith("svc.") or check_subject(rpc):
                errors.append(f"skill tool {tool['name']!r}: rpc {rpc!r} must be an on-grammar svc.* subject")
            elif not _declared(bus, ("requests", "serves"), rpc):
                errors.append(f"skill tool {tool['name']!r}: rpc {rpc!r} is not in bus.requests or bus.serves")
    video = manifest.get("video")
    if video:
        cls, eid = video["path"].split("/")
        want = f"ent.{VIDEO_CLASS_TYPE[cls]}.{eid}.state"
        if video["entity"] != want:
            errors.append(f"video.entity {video['entity']!r} does not match path {video['path']!r} (want {want!r})")
    if manifest["trust"] == "community":
        if src and src.get("image"):
            errors.append("community plugins are declarative: no source.image")
        if skill:
            errors.append("community plugins are declarative: no skill code")
    return errors


def tracked_manifests() -> list[Path]:
    out = subprocess.check_output(["git", "ls-files", "-z", "*onemind-plugin.yaml"], cwd=ROOT)
    return [ROOT / p for p in out.decode().split("\0") if p and not p.startswith(FIXTURES)]


def main(argv: list[str]) -> int:
    schema = json.loads(SCHEMA.read_text())
    files = [Path(a) for a in argv] if argv else tracked_manifests()
    if not files:
        print("no onemind-plugin.yaml found", file=sys.stderr)
        return 1
    failed = 0
    for f in files:
        try:
            manifest = yaml.safe_load(f.read_text())
        except yaml.YAMLError as exc:
            print(f"FAIL {f}: {exc}")
            failed += 1
            continue
        errs = validate(manifest if isinstance(manifest, dict) else {}, schema)
        rel = f.resolve().relative_to(ROOT) if f.resolve().is_relative_to(ROOT) else f
        if errs:
            failed += 1
            print(f"FAIL {rel}")
            for e in errs:
                print(f"  - {e}")
        else:
            print(f"ok   {rel} ({manifest['id']} {manifest['version']}, kinds {','.join(manifest['kinds'])})")
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
