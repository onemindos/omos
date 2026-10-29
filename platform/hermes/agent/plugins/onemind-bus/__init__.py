"""onemind-bus: Hermes tools over the OneMind NATS bus (FEATURES I2, D5, D11).

Registers named ontology tools (entity and layer reads), policy proposal tools (no confirm
tool exists) and reversible map tools (cmd.map.*). With ONEMIND_BUS_RETIRE_NATS_PUBLISH on
(the default) it also blocks the generic ``nats_publish`` tool of the nats-context plugin
through a ``pre_tool_call`` hook, so the agent loses its bus pipe as soon as this plugin is
enabled. See README.md for configuration and the install steps.
"""
from __future__ import annotations

import logging
import os
from typing import Any, Dict, Optional

from .bus import BusClient, BusConfig
from .tools import SPECS, TOOLSET, Tools, make_handler

log = logging.getLogger("onemind_bus")

# Tools of other plugins that this plugin retires (generic bus pipes, OWASP LLM06).
RETIRED_TOOLS_DEFAULT = ("nats_publish",)
RETIRED_MESSAGE = ("Tool '{tool}' is retired: OneMind bus writes go through onemind_command_propose / "
                   "onemind_task_propose (human confirm) or the onemind_map_* tools.")


def _truthy(v: Optional[str], default: bool) -> bool:
    if v is None or v == "":
        return default
    return v.strip().lower() in ("1", "true", "yes", "on")


def retired_tools(env: Optional[Dict[str, str]] = None) -> frozenset:
    e = os.environ if env is None else env
    if not _truthy(e.get("ONEMIND_BUS_RETIRE_NATS_PUBLISH"), True):
        return frozenset()
    extra = [t.strip() for t in (e.get("ONEMIND_BUS_RETIRED_TOOLS") or "").split(",") if t.strip()]
    return frozenset(RETIRED_TOOLS_DEFAULT) | frozenset(extra)


def make_pre_tool_call(blocked: frozenset):
    def pre_tool_call(tool_name: str = "", **_: Any):
        if tool_name in blocked:
            log.warning('{"event":"onemind_bus_retired_tool_blocked","tool":"%s"}', tool_name)
            return {"action": "block", "message": RETIRED_MESSAGE.format(tool=tool_name)}
        return None
    return pre_tool_call


def _nats_available() -> bool:
    try:
        import nats  # noqa: F401
        return True
    except Exception:
        return False


def build(env: Optional[Dict[str, str]] = None, client: Optional[BusClient] = None) -> Tools:
    e = os.environ if env is None else env
    return Tools(
        client=client or BusClient(BusConfig.from_env(e)),
        agent_id=e.get("ONEMIND_BUS_AGENT_ID") or "agent:hermes/zeus/legacy",
        on_behalf_of=e.get("ONEMIND_BUS_ON_BEHALF_OF") or None,
        model_alias=e.get("ONEMIND_BUS_MODEL_ALIAS") or None,
        layer_snapshot=_truthy(e.get("ONEMIND_BUS_LAYER_SNAPSHOT"), True),
    )


def register(ctx, env: Optional[Dict[str, str]] = None, client: Optional[BusClient] = None) -> Tools:
    tools = build(env, client)
    for name, method, description, params in SPECS:
        ctx.register_tool(
            name=name,
            toolset=TOOLSET,
            schema={"name": name, "description": description, "parameters": params},
            handler=make_handler(tools, name, method),
            check_fn=_nats_available,
            description=description,
        )
    blocked = retired_tools(env)
    if blocked:
        ctx.register_hook("pre_tool_call", make_pre_tool_call(blocked))
    log.info('{"event":"onemind_bus_registered","tools":%d,"retired":%s}', len(SPECS), sorted(blocked))
    return tools
