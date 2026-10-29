"""Test harness for the onemind-bus plugin.

Every test that needs a bus starts its own disposable nats-server (JetStream, temp store) on a
free port and stops it afterwards. Set NATS_SERVER_BIN or put nats-server on PATH; without it
those tests skip. The plugin is loaded the way Hermes loads a directory plugin
(hermes_cli/plugins.py ``_load_directory_module``): as package ``hermes_plugins.onemind_bus``.
"""
from __future__ import annotations

import importlib.util
import os
import shutil
import socket
import subprocess
import sys
import time
import types
from pathlib import Path

import pytest

PLUGIN_DIR = Path(__file__).resolve().parents[1]
FIXTURES = Path(__file__).resolve().parent / "fixtures"
sys.path.insert(0, str(Path(__file__).resolve().parent))


def load_plugin():
    name = "hermes_plugins.onemind_bus"
    if name in sys.modules:
        return sys.modules[name]
    if "hermes_plugins" not in sys.modules:
        ns = types.ModuleType("hermes_plugins")
        ns.__path__ = []
        sys.modules["hermes_plugins"] = ns
    spec = importlib.util.spec_from_file_location(name, PLUGIN_DIR / "__init__.py",
                                                  submodule_search_locations=[str(PLUGIN_DIR)])
    mod = importlib.util.module_from_spec(spec)
    mod.__package__ = name
    mod.__path__ = [str(PLUGIN_DIR)]
    sys.modules[name] = mod
    spec.loader.exec_module(mod)
    return mod


class FakeCtx:
    """The two PluginContext methods the plugin uses, plus Hermes-style dispatch."""

    def __init__(self):
        self.tools = {}
        self.hooks = {}

    def register_tool(self, name, toolset, schema, handler, check_fn=None, requires_env=None,
                      is_async=False, description="", emoji="", override=False):
        assert name not in self.tools, f"duplicate tool {name}"
        self.tools[name] = {"toolset": toolset, "schema": schema, "handler": handler, "check_fn": check_fn}

    def register_hook(self, name, cb):
        self.hooks.setdefault(name, []).append(cb)

    def blocked(self, tool_name, args=None):
        for cb in self.hooks.get("pre_tool_call", []):
            r = cb(tool_name=tool_name, args=args or {}, task_id="", session_id="", tool_call_id="")
            if isinstance(r, dict) and r.get("action") == "block":
                return r["message"]
        return None

    def dispatch(self, name, args, task_id="session-test", user_task=None):
        import json

        msg = self.blocked(name, args)
        if msg:
            return {"blocked": msg}
        return json.loads(self.tools[name]["handler"](args, task_id=task_id, user_task=user_task))


def nats_server_bin():
    return os.environ.get("NATS_SERVER_BIN") or shutil.which("nats-server")


def free_port():
    with socket.socket() as s:
        s.bind(("127.0.0.1", 0))
        return s.getsockname()[1]


class NatsServer:
    def __init__(self, workdir: Path, config: str = ""):
        self.bin = nats_server_bin()
        self.port = free_port()
        self.workdir = workdir
        workdir.mkdir(parents=True, exist_ok=True)
        conf = workdir / "nats.conf"
        conf.write_text(f'listen: "127.0.0.1:{self.port}"\njetstream {{ store_dir: "{workdir / "js"}" }}\n' + config)
        os.chmod(conf, 0o600)
        self.log = open(workdir / "nats-server.log", "w")
        self.proc = subprocess.Popen([self.bin, "-c", str(conf)], stdout=self.log, stderr=subprocess.STDOUT)
        deadline = time.time() + 10
        while time.time() < deadline:
            try:
                with socket.create_connection(("127.0.0.1", self.port), 0.2):
                    return
            except OSError:
                time.sleep(0.05)
        self.stop()
        raise RuntimeError("nats-server did not start")

    def url(self, user=None, password=None):
        if user:
            return f"nats://{user}:{password}@127.0.0.1:{self.port}"
        return f"nats://127.0.0.1:{self.port}"

    def log_text(self):
        if not self.log.closed:
            self.log.flush()
        return (self.workdir / "nats-server.log").read_text()

    def stop(self):
        if self.proc.poll() is None:
            self.proc.terminate()
            try:
                self.proc.wait(10)
            except subprocess.TimeoutExpired:
                self.proc.kill()
        self.log.close()


@pytest.fixture
def nats_server(tmp_path):
    if not nats_server_bin():
        pytest.skip("nats-server not found (set NATS_SERVER_BIN)")
    srv = NatsServer(tmp_path / "nats")
    yield srv
    srv.stop()


@pytest.fixture
def plugin():
    return load_plugin()


@pytest.fixture
def world(nats_server):
    from stubs import World

    w = World(nats_server.url(), FIXTURES / "entities.json")
    w.start()
    yield w
    w.stop()


@pytest.fixture
def hermes(plugin, nats_server, world):
    """The plugin registered on a fake Hermes context, connected to the disposable server."""
    from stubs import RecordingClient

    env = {"ONEMIND_BUS_URL": nats_server.url(), "ONEMIND_BUS_ON_BEHALF_OF": "user:zeus",
           "ONEMIND_BUS_MAP_ACK_WINDOW_S": "1.0", "ONEMIND_BUS_TIMEOUT_S": "3"}
    client = RecordingClient(plugin.bus.BusConfig.from_env(env))
    ctx = FakeCtx()
    plugin.register(ctx, env=env, client=client)
    ctx.client = client
    yield ctx
    client.close()
