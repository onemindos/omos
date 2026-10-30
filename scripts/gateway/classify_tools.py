"""Propose or check a tool classification file for one gateway target.

  # offline, what CI runs: every classification file is valid and matches a target fragment
  python scripts/gateway/classify_tools.py --check

  # live: call tools/list on an MCP endpoint and compare it with platform/tool-gateway/policy/classification/<target>.yaml
  python scripts/gateway/classify_tools.py --target grafana --url http://127.0.0.1:8000/mcp            # report drift, exit 1
  python scripts/gateway/classify_tools.py --target grafana --url http://127.0.0.1:8000/mcp --write    # add new tools
  # credentials are read from the environment, never from argv:
  python scripts/gateway/classify_tools.py ... --header-env Authorization="Bearer {}" --key-env GRAFANA_TOKEN

Point --url at the backend itself (a port-forward) to discover tools: a gateway group route only lists the
tools that are already classified. To check the exposed list of a live route, point --url at the route and
pass --strip-prefix <target>_ (prefixMode always).

New tools are proposed by name (heuristics below; unknown means destructive) and written under the key
`proposed` as well as their class. The renderer and --check fail while `proposed` is non-empty, so a human
must read every proposal, fix the class and delete the key. Nothing here talks to the cluster API.
"""
from __future__ import annotations

import argparse
import json
import os
import re
import sys
import urllib.request
from pathlib import Path

import yaml

sys.path.insert(0, str(Path(__file__).resolve().parent))
import classification as cl  # noqa: E402

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / "platform/tool-gateway"
DIR = BASE / "policy/classification"

# Tokens of a tool name (split on _ - . and camelCase). First match wins in this order, so a name that
# contains both "get" and "delete" is destructive.
PHYSICAL = {"lock", "unlock", "door", "garage", "alarm", "arm", "disarm", "siren", "valve", "cover", "open", "close",
            "takeoff", "land", "arm", "cmd", "actuate", "fire"}
DESTRUCTIVE = {"delete", "remove", "destroy", "drop", "kill", "purge", "truncate", "wipe", "prune", "revoke",
               "terminate", "uninstall", "reset", "restart", "reboot", "shutdown", "stop", "scale", "evict", "drain",
               "kick", "ban", "disable", "deactivate", "unlink", "detach", "rollback", "force", "exec", "execute"}
READ = {"get", "list", "search", "describe", "query", "read", "fetch", "find", "show", "status", "health", "top",
        "logs", "log", "count", "stats", "events", "view", "inspect", "check", "info", "lookup", "resolve", "history",
        "diff", "tree", "browse", "explain", "summary", "overview", "metrics", "ping", "version", "whoami", "me"}
WRITE = {"create", "update", "add", "set", "apply", "patch", "put", "post", "annotate", "label", "tag", "sync",
         "refresh", "run", "publish", "send", "comment", "import", "upload", "write", "edit", "rename", "move", "save",
         "enable", "activate", "start", "push", "merge", "assign", "toggle", "turn", "call", "invoke", "reply"}


def tokens(name: str) -> list[str]:
    name = re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", name)
    return [t for t in re.split(r"[_\-. ]+", name.lower()) if t]


def propose(name: str, annotations: dict | None = None) -> str:
    """Heuristic class for a tool name. A proposal only: a human reviews every one (see module docstring)."""
    toks = set(tokens(name))
    if toks & PHYSICAL:
        return "physical"
    if toks & DESTRUCTIVE:
        return "destructive"
    first = tokens(name)[0] if tokens(name) else ""
    if first in READ or (toks & READ and not toks & WRITE):
        return "read"
    if first in WRITE or toks & WRITE:
        return "write"
    if (annotations or {}).get("readOnlyHint") is True:
        return "read"
    return "destructive"


def _parse_body(body: str) -> dict:
    body = body.strip()
    if body.startswith("{"):
        return json.loads(body)
    for line in body.splitlines():  # text/event-stream
        if line.startswith("data:"):
            return json.loads(line[5:])
    raise RuntimeError(f"unparseable MCP response: {body[:120]!r}")


def list_tools(url: str, headers: dict[str, str], timeout: float = 20.0) -> list[dict]:
    """initialize + tools/list over streamable HTTP (JSON or SSE replies), following pagination."""
    base = {"content-type": "application/json", "accept": "application/json, text/event-stream", **headers}

    def post(payload: dict, sid: str | None):
        h = dict(base)
        if sid:
            h["mcp-session-id"] = sid
        req = urllib.request.Request(url, data=json.dumps(payload).encode(), headers=h, method="POST")
        with urllib.request.urlopen(req, timeout=timeout) as r:  # noqa: S310 (operator-supplied URL)
            return r.headers.get("mcp-session-id"), r.read().decode()

    sid, body = post({"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {
        "protocolVersion": "2025-06-18", "capabilities": {}, "clientInfo": {"name": "classify_tools", "version": "1"}}}, None)
    _parse_body(body)
    post({"jsonrpc": "2.0", "method": "notifications/initialized"}, sid)
    tools, cursor, n = [], None, 2
    while True:
        params = {"cursor": cursor} if cursor else {}
        _, body = post({"jsonrpc": "2.0", "id": n, "method": "tools/list", "params": params}, sid)
        res = _parse_body(body)
        if "error" in res:
            raise RuntimeError(f"tools/list failed: {res['error']}")
        tools += res["result"]["tools"]
        cursor = res["result"].get("nextCursor")
        n += 1
        if not cursor:
            return tools


def compare(live: list[str], c: dict) -> tuple[list[str], list[str]]:
    """(new, stale): live names in no list, and exact names in the file that the backend no longer has."""
    new = sorted(n for n in live if not any(cl.matches(n, p) for k in cl.CLASSES for p in c[k]))
    exact = [p for k in cl.CLASSES for p in c[k] if "*" not in p]
    stale = sorted(p for p in exact if p not in live)
    return new, stale


def render_file(c: dict, proposed: list[str]) -> str:
    out = {"target": c["target"]}
    for k in cl.CLASSES:
        out[k] = sorted(c[k])
    if proposed:
        out["proposed"] = sorted(proposed)
    head = (f"# Tool classification for gateway target '{c['target']}'. Reviewed in Git; unlisted = destructive (hidden).\n"
            "# Schema and rules: scripts/gateway/classification.py. Refresh: scripts/gateway/classify_tools.py.\n")
    return head + yaml.safe_dump(out, sort_keys=False, default_flow_style=False, width=120)


def check_offline() -> int:
    from render_gateway_config import load, render, dump, OUT  # type: ignore

    try:
        text = dump(render(*load()))
    except cl.PolicyError as e:
        print(f"classification error: {e}", file=sys.stderr)
        return 1
    if not OUT.is_file() or OUT.read_text() != text:
        print("values.generated.yaml is stale: run python scripts/gateway/render_gateway_config.py", file=sys.stderr)
        return 1
    print("classification files are valid, match their targets, and values.generated.yaml is current")
    return 0


def main(argv: list[str]) -> int:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--check", action="store_true", help="offline validation of all files (CI)")
    ap.add_argument("--target")
    ap.add_argument("--url")
    ap.add_argument("--write", action="store_true", help="add proposals for new tools to the classification file")
    ap.add_argument("--strip-prefix", default="", help="remove this prefix from live tool names")
    ap.add_argument("--header-env", action="append", default=[], metavar="NAME=TEMPLATE",
                    help="header NAME with TEMPLATE, '{}' replaced by the value of --key-env (default Authorization=Bearer {})")
    ap.add_argument("--key-env", help="environment variable holding the credential for --header-env")
    a = ap.parse_args(argv)

    if a.check:
        return check_offline()
    if not (a.target and a.url):
        ap.error("give --check, or --target and --url")

    headers = {}
    if a.key_env:
        secret = os.environ.get(a.key_env)
        if not secret:
            print(f"environment variable {a.key_env} is empty", file=sys.stderr)
            return 2
        for spec in a.header_env or ["Authorization=Bearer {}"]:
            name, _, tmpl = spec.partition("=")
            headers[name] = tmpl.replace("{}", secret)

    path = DIR / f"{a.target}.yaml"
    existing = cl.parse(path) if path.is_file() else {
        "target": a.target, **{k: [] for k in cl.CLASSES}, "proposed": []}
    tools = list_tools(a.url, headers)
    live = {}
    for t in tools:
        n = t["name"][len(a.strip_prefix):] if t["name"].startswith(a.strip_prefix) else t["name"]
        live[n] = t.get("annotations") or {}
    new, stale = compare(sorted(live), existing)

    print(f"{a.target}: {len(live)} live tools, {len(new)} unclassified, {len(stale)} stale")
    for n in new:
        print(f"  new   {n:40s} proposed {propose(n, live[n])}")
    for n in stale:
        print(f"  stale {n} (in the file, not on the backend)")
    if a.write and new:
        props = {k: list(existing[k]) for k in cl.CLASSES}
        for n in new:
            props[propose(n, live[n])].append(n)
        merged = {**existing, **props}
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(render_file(merged, new))
        print(f"wrote {path.relative_to(ROOT)} with a 'proposed' list: review, then delete that key")
        return 1  # still needs review
    return 1 if new or stale else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
