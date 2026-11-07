"""In-process mission execution.

Bridge: this replaced the Dapr workflow `onemind_mission_v1` (removed 2026-11-07).
Missions run once, in a thread of the API process, with state recorded in the
missions table. There is no durable replay: a mission whose worker dies stays
`running` and is never re-run automatically. Mission orchestration is expected
to move to Hermes, at which point this module goes away.
"""
from . import store
from .capabilities import execute
import json


def run_mission(inp):
    current = store.mission(inp["id"])
    if current["state"] == "completed":
        return current["result"]
    if current["state"] == "cancelled":
        raise RuntimeError("Mission was cancelled")
    # Do not repeat ambiguous external agent work after a worker failure.
    if inp["spec"]["capability"] == "delegate-agent" and store.dispatched(inp["id"]):
        raise RuntimeError("Agent outcome is ambiguous; inspect evidence before issuing a new mission")
    store.set_state(inp["id"], "running")
    try:
        result = json.loads(json.dumps(execute(inp["spec"], inp["id"], inp["actor"]), default=str))
        store.set_state(inp["id"], "completed", result=result)
        return result
    except Exception as exc:
        store.set_state(inp["id"], "failed", error=type(exc).__name__)
        raise
