"""Read-only check of the server endpoints Hermes Desktop calls, run in the pod as `hermes`.

    desktop/apply.sh --check

It imports the dashboard app (`hermes_cli.web_server`) in a fresh process and calls every parameterless GET
route in-process through ASGI with that process's own session token, against the real /opt/data. That
exercises the route handlers and the plugin routers that mount, not the OAuth/basic-auth gate or the
network path in front of them. Only status codes and the key names of the JSON bodies are printed, never
values. Routes that return credentials, stream, download or start OAuth flows are skipped.

It also reports the voice and media providers Desktop depends on (STT, TTS, image generation), because those
fail quietly: a route answers 200 while no provider is behind it.
"""
import asyncio
import importlib
import json
import os
import re
import sys

os.environ.setdefault("HERMES_HOME", "/opt/data")
DEPS = os.environ.get("ONEMIND_GATEWAY_DEPS", "/opt/data/.local/gateway-deps")
if os.path.isdir(DEPS) and DEPS not in sys.path:
    sys.path.append(DEPS)  # the same thing the nats-platform plugin does at gateway start

import httpx  # noqa: E402
from hermes_cli import web_server as ws  # noqa: E402

SKIP = re.compile(
    r"(/reveal|/raw$|/credentials/pool|/download|/export|/ws$|/pty|/pub$|/events$|/display/ws|/stream"
    r"|/backup|/dump|/oauth|/import|/console$|/hermes/update/check$)"
)
# /api/auth/providers is filled by the launch process (a fresh process has none registered); the SSH ownership
# route is 404 unless Desktop runs the backend over SSH, which this pod never does.
EXPECTED = {"/api/auth/providers", "/api/ssh/ownership"}
TOKENISH =re.compile(r"(eyJ[\w-]+\.[\w-]+\.[\w-]+|[A-Za-z0-9_\-]{32,})")


def shape(body):
    try:
        data = json.loads(body)
    except ValueError:
        return "non-json len=%d" % len(body)
    if isinstance(data, dict):
        return "dict keys=" + ",".join(list(data)[:8])
    if isinstance(data, list):
        return "list len=%d" % len(data)
    return type(data).__name__


async def routes_check():
    app = ws.app
    headers = {ws._SESSION_HEADER_NAME: ws._SESSION_TOKEN}
    paths = sorted({
        r.path for r in app.routes
        if getattr(r, "path", "").startswith("/api/") and "GET" in (getattr(r, "methods", None) or [])
        and "{" not in r.path
    })
    bad = 0
    async with httpx.AsyncClient(transport=httpx.ASGITransport(app=app), base_url="http://127.0.0.1", timeout=30) as c:
        for path in paths:
            if SKIP.search(path):
                continue
            try:
                r = await c.get(path, headers=headers)
            except Exception as exc:  # report, keep going
                print("EXC  %-48s %s" % (path, type(exc).__name__))
                bad += 1
                continue
            if r.status_code == 422:  # needs a query parameter the check does not know
                continue
            if r.status_code == 400 and "required" in r.text:  # same, reported as 400 by these routes
                continue
            if path in EXPECTED:  # answers an error by design in this fresh process
                continue
            if r.status_code >= 400:
                bad += 1
                print("%-4d %-48s %s" % (r.status_code, path, TOKENISH.sub("<R>", r.text[:100].replace("\n", " "))))
    print("routes checked: %d GET routes, %d with errors" % (len(paths), bad))
    return bad


def providers_check():
    problems = 0
    try:
        from hermes_cli.plugins import discover_plugins
        discover_plugins()
    except Exception as exc:
        print("plugin discovery: %s" % type(exc).__name__)
    # STT
    try:
        from tools import transcription_tools as stt
        provider = stt._get_provider(stt._load_stt_config())
        print("stt provider: %s" % provider)
        if provider in (None, "none"):
            problems += 1
    except Exception as exc:
        print("stt provider: check failed (%s)" % type(exc).__name__)
        problems += 1
    # TTS (Edge needs the edge_tts package from the gateway-deps directory)
    try:
        importlib.import_module("edge_tts")
        print("tts edge: edge_tts importable")
    except ImportError:
        print("tts edge: edge_tts NOT importable")
        problems += 1
    # Image generation
    try:
        from agent import image_gen_registry as ig
        active = ig.get_active_provider()
        name = getattr(active, "provider_id", None) or getattr(active, "name", None)
        print("image provider: %s" % name)
        if active is None:
            problems += 1
    except Exception as exc:
        print("image provider: check failed (%s)" % type(exc).__name__)
        problems += 1
    return problems


def main():
    bad = asyncio.run(routes_check())
    problems = providers_check()
    print("result: %s" % ("OK" if not bad and not problems else "%d route error(s), %d provider problem(s)" % (bad, problems)))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
