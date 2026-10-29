"""Map-action acks across accounts (proposal I2-P2).

In the A19 layout (omos PR #562, platform/nats/identities/identities.yaml) AGENTS exports
``cmd.map.>`` to CLIENT as a *stream*. A stream import carries the message but not a usable reply
path: a client's ack to the reply subject lands in CLIENT, so Hermes never hears that the map
applied the action, and the D11 rule "report only actions that succeeded" makes every map tool
report failure. With CLIENT exporting ``cmd.map.>`` as a *service* (response_type: stream) that
AGENTS imports, the acks come back. Both layouts run here on a disposable server.
"""
from __future__ import annotations

import asyncio
import json
import secrets
import threading

import pytest

from conftest import FakeCtx, NatsServer, load_plugin, nats_server_bin

P = load_plugin()
pytestmark = pytest.mark.skipif(not nats_server_bin(), reason="nats-server not found")

HERMES_PERMS = """permissions: {
  publish: { allow: ["agents.*.hermes.>", "cmd.map.>", "svc.entity.>", "svc.layer.list", "svc.layer.meta",
                     "svc.policy.command.propose", "svc.policy.task.propose", "svc.policy.action.get"] }
  subscribe: { allow: ["agents.prompt.hermes.>", "_INBOX.>"] } }"""
ZEUS_PERMS = """permissions: {
  publish: { allow: ["svc.>", "cmd.map.>"] }, subscribe: { allow: ["cmd.map.>", "_INBOX.>"] }, allow_responses: true }"""

LAYOUTS = {
    # as in PR #562 today
    "a19_stream_import": ('exports: [ { stream: "cmd.map.>", accounts: [CLIENT] } ]',
                          'imports: [ { stream: { account: AGENTS, subject: "cmd.map.>" } } ]'),
    # proposal I2-P2
    "i2p2_service_import": ('imports: [ { service: { account: CLIENT, subject: "cmd.map.>" } } ]',
                            'exports: [ { service: "cmd.map.>", accounts: [AGENTS], response_type: stream } ]'),
}


def conf(layout, pw):
    agents, client = LAYOUTS[layout]
    return f"""
accounts {{
  AGENTS {{ users: [ {{ user: hermes, password: "{pw['hermes']}", {HERMES_PERMS} }} ]
    {agents} }}
  CLIENT {{ users: [ {{ user: zeus, password: "{pw['zeus']}", {ZEUS_PERMS} }} ]
    {client} }}
}}
"""


def run_layout(tmp_path, layout):
    pw = {"hermes": secrets.token_hex(12), "zeus": secrets.token_hex(12)}
    srv = NatsServer(tmp_path / layout, conf(layout, pw))
    loop = asyncio.new_event_loop()
    threading.Thread(target=loop.run_forever, daemon=True).start()
    received = []

    async def client():
        import nats

        nc = await nats.connect(srv.url("zeus", pw["zeus"]))

        async def cb(m):
            received.append(m.subject)
            p = json.loads(m.data)
            if m.reply:
                await nc.publish(m.reply, json.dumps({"v": 1, "id": p["id"], "ok": True, "client": "zeus-web",
                                                      "applied": {"camera": {"lat": p["lat"], "lon": p["lon"]}}}).encode())

        await nc.subscribe("cmd.map.>", cb=cb)
        await nc.flush()
        return nc

    znc = asyncio.run_coroutine_threadsafe(client(), loop).result(5)
    ctx = FakeCtx()
    tools = P.register(ctx, env={"ONEMIND_BUS_URL": srv.url("hermes", pw["hermes"]), "ONEMIND_BUS_MAP_ACK_WINDOW_S": "0.8"})
    try:
        return ctx.dispatch("onemind_map_fly_to", {"lat": 30.27, "lon": -97.74}), received
    finally:
        tools.client.close()
        asyncio.run_coroutine_threadsafe(znc.close(), loop).result(5)
        loop.call_soon_threadsafe(loop.stop)
        srv.stop()


def test_a19_stream_import_delivers_but_loses_the_ack(tmp_path):
    r, received = run_layout(tmp_path, "a19_stream_import")
    assert received == ["cmd.map.viewport.set"]  # the client got it and applied it ...
    assert r["ok"] is False and r["error"] == "no_client_ack"  # ... but Hermes cannot know


def test_i2p2_service_import_returns_the_ack(tmp_path):
    r, received = run_layout(tmp_path, "i2p2_service_import")
    assert received == ["cmd.map.viewport.set"]
    assert r["ok"] is True and r["applied_by"] == ["zeus-web"]
