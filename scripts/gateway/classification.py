"""Tool classification files for the tool gateway: schema, validation and the CEL they compile to.

One reviewed file per target: platform/tool-gateway/policy/classification/<target>.yaml

    target: grafana
    read: [search_dashboards, query_*]        # exact names or globs ('*' only)
    write: [create_annotation]
    destructive: [delete_dashboard]           # execute route only (typed `execute <id>`)
    physical: []                              # execute route only (typed `override <id>`)
    guards:                                   # optional: a listed tool that becomes destructive for some arguments
    - {tool: resources_scale, class: destructive, when: 'mcp.tool.arguments.replicas == 0'}
    proposed: []                              # written by classify_tools.py; CI fails until a human empties it

A tool that appears in no list is treated as destructive: it is not exposed on any group route and the
classification check reports it as drift. Design: docs/architecture/onemind-tool-gateway.md section 3.1.
"""
from __future__ import annotations

import re
from pathlib import Path

import yaml

CLASSES = ("read", "write", "destructive", "physical")
EXPOSED = ("read", "write")  # classes allowed on a group route
GATED = ("destructive", "physical")  # classes that only the execute route may ever expose
NAME_RE = re.compile(r"^[A-Za-z0-9_.*-]+$")
TARGET_RE = re.compile(r"^[a-z][a-z0-9]*$")
KEYS = {"target", *CLASSES, "guards", "proposed"}


class PolicyError(Exception):
    pass


def glob_re(pattern: str) -> re.Pattern:
    return re.compile("^" + ".*".join(re.escape(p) for p in pattern.split("*")) + "$")


def matches(name: str, pattern: str) -> bool:
    return bool(glob_re(pattern).match(name))


def parse(path: Path) -> dict:
    """Read and validate one classification file; returns the normalised mapping."""
    where = path.name
    raw = yaml.safe_load(path.read_text())
    if not isinstance(raw, dict):
        raise PolicyError(f"{where}: expected a mapping")
    unknown = set(raw) - KEYS
    if unknown:
        raise PolicyError(f"{where}: unknown keys {sorted(unknown)}")
    target = raw.get("target")
    if not isinstance(target, str) or not TARGET_RE.match(target):
        raise PolicyError(f"{where}: 'target' must match {TARGET_RE.pattern}")
    if where != f"{target}.yaml":
        raise PolicyError(f"{where}: file name must be {target}.yaml")

    out: dict = {"target": target}
    seen: dict[str, str] = {}
    for cls in CLASSES:
        names = raw.get(cls) or []
        if not isinstance(names, list) or not all(isinstance(n, str) for n in names):
            raise PolicyError(f"{where}: '{cls}' must be a list of strings")
        for n in names:
            if not NAME_RE.match(n):
                raise PolicyError(f"{where}: {cls} entry {n!r} must match {NAME_RE.pattern}")
            if n in seen:
                raise PolicyError(f"{where}: {n!r} is listed under both {seen[n]} and {cls}")
            seen[n] = cls
        out[cls] = sorted(names)

    # An exact destructive/physical name must not be swallowed by a read/write glob, and the reverse.
    for gated in (n for c in GATED for n in out[c]):
        if "*" in gated:
            continue
        for exposed in (n for c in EXPOSED for n in out[c]):
            if "*" in exposed and matches(gated, exposed):
                raise PolicyError(f"{where}: {gated!r} is destructive/physical but matches the glob {exposed!r}")
    for gated in (n for c in GATED for n in out[c]):
        if "*" not in gated:
            continue
        for exposed in (n for c in EXPOSED for n in out[c]):
            if "*" not in exposed and matches(exposed, gated):
                raise PolicyError(f"{where}: {exposed!r} is exposed but matches the destructive/physical glob {gated!r}")

    guards = raw.get("guards") or []
    if not isinstance(guards, list):
        raise PolicyError(f"{where}: 'guards' must be a list")
    norm_guards = []
    for g in guards:
        if not isinstance(g, dict) or set(g) != {"tool", "class", "when"}:
            raise PolicyError(f"{where}: each guard needs exactly tool, class and when")
        if g["class"] not in GATED:
            raise PolicyError(f"{where}: guard class must be one of {GATED}")
        if seen.get(g["tool"]) not in EXPOSED or "*" in g["tool"]:
            raise PolicyError(f"{where}: guard tool {g['tool']!r} must be an exact name listed under read or write")
        if not isinstance(g["when"], str) or not g["when"].strip():
            raise PolicyError(f"{where}: guard 'when' must be a CEL string")
        stripped = g["when"].replace("$$", "")
        if "$" in stripped:
            raise PolicyError(f"{where}: bare '$' in guard {g['when']!r}; write '$$'")
        norm_guards.append({"tool": g["tool"], "class": g["class"], "when": g["when"].strip()})
    out["guards"] = sorted(norm_guards, key=lambda g: g["tool"])

    proposed = raw.get("proposed") or []
    if proposed:
        raise PolicyError(f"{where}: 'proposed' is not empty ({sorted(proposed)}): review the tools, put each in its "
                          "final class, then delete the 'proposed' key")
    out["proposed"] = []
    return out


def load_all(directory: Path) -> dict[str, dict]:
    result: dict[str, dict] = {}
    for path in sorted(directory.glob("*.yaml")):
        c = parse(path)
        result[c["target"]] = c
    return result


def _name_term(patterns: list[str]) -> list[str]:
    exact = sorted(p for p in patterns if "*" not in p)
    terms = []
    if exact:
        terms.append("mcp.tool.name in [%s]" % ", ".join('"%s"' % n for n in exact))
    for p in sorted(p for p in patterns if "*" in p):
        rx = ".*".join(part.replace(".", "[.]") for part in p.split("*"))
        terms.append('mcp.tool.name.matches("^%s$$")' % rx)  # '$$' is the gateway's escape for a literal '$'
    return terms


def allow_expr(c: dict) -> str | None:
    """CEL for the tools of this target that a group route may expose, or None when there are none."""
    guarded = {g["tool"]: g["when"] for g in c["guards"]}
    patterns = [n for cls in EXPOSED for n in c[cls] if n not in guarded]
    terms = _name_term(patterns)
    for tool, when in sorted(guarded.items()):
        terms.append('(mcp.tool.name == "%s" && !(%s))' % (tool, when))
    if not terms:
        return None
    return 'mcp.tool.target == "%s" && (%s)' % (c["target"], " || ".join(terms))


def has_gated(c: dict) -> bool:
    return bool(c["destructive"] or c["physical"] or c["guards"])


def classify_name(name: str, c: dict) -> str:
    """Class of a live tool name under this classification; unlisted is destructive (fail closed)."""
    for cls in CLASSES:
        if any(matches(name, p) for p in c[cls]):
            return cls
    return "destructive"
