"""Versioned workflow; activities record evidence outside deterministic orchestration."""
from dapr.ext.workflow import WorkflowRuntime, DaprWorkflowContext, WorkflowActivityContext
from . import store
from .capabilities import execute

runtime = WorkflowRuntime()


@runtime.activity(name="execute_capability_v1")
def execute_capability(ctx: WorkflowActivityContext, inp):
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
        result = execute(inp["spec"], inp["id"], inp["actor"])
        store.set_state(inp["id"], "completed", result=result)
        return result
    except Exception as exc:
        store.set_state(inp["id"], "failed", error=type(exc).__name__)
        raise


@runtime.workflow(name="onemind_mission_v1")
def mission_workflow(ctx: DaprWorkflowContext, inp):
    return (yield ctx.call_activity(execute_capability, input=inp))
