import asyncio
import contextlib
import hashlib
import hmac
import json
import os
import uuid
from contextlib import asynccontextmanager
from datetime import datetime, timezone
import httpx
from fastapi import Depends, FastAPI, Header, HTTPException, Query
from pydantic import BaseModel, Field, field_validator
from dapr.ext.workflow import DaprWorkflowClient
from . import store
from .capabilities import CAPABILITIES, ALLOWED_RESTARTS, agents_and_call


def principal(authorization: str = Header(default=""), x_onemind_actor: str = Header(default="")):
    scheme, _, token = authorization.partition(" ")
    if scheme.lower() != "bearer" or not token:
        raise HTTPException(401, "Authentication required")
    ui_key = os.environ.get("UI_GATEWAY_KEY", "")
    if ui_key and hmac.compare_digest(token, ui_key):
        if not x_onemind_actor.startswith("user:") or len(x_onemind_actor) > 200:
            raise HTTPException(403, "A gateway-authenticated user identity is required")
        return {"actor": x_onemind_actor, "operator": True}
    read_key = os.environ.get("READ_API_KEY", "")
    if read_key and hmac.compare_digest(token, read_key):
        return {"actor": "service:reader", "operator": False}
    raise HTTPException(401, "Invalid credential")


def operator(p=Depends(principal)):
    if not p["operator"]:
        raise HTTPException(403, "An operator capability is required")
    return p


class ObjectInput(BaseModel):
    id: str = Field(pattern=r"^[a-zA-Z0-9][a-zA-Z0-9:/_.-]{0,199}$")
    kind: str = Field(pattern=r"^[a-z][a-z0-9-]{0,49}$")
    name: str = Field(min_length=1, max_length=200)
    properties: dict = Field(default_factory=dict)
    source: str = Field(default="onemind", min_length=1, max_length=300)

    @field_validator("properties")
    @classmethod
    def size_limit(cls, value):
        if len(json.dumps(value)) > 64000:
            raise ValueError("Object properties exceed the record limit")
        return value


class MissionInput(BaseModel):
    capability: str
    prompt: str = Field(default="Inspect current OneMind conditions", max_length=12000)
    target: str = Field(default="", max_length=300)
    idempotency_key: str = Field(min_length=1, max_length=100)

    @field_validator("capability")
    @classmethod
    def known(cls, value):
        if value not in CAPABILITIES:
            raise ValueError("Unknown capability")
        return value


async def dispatch_pending():
    client = DaprWorkflowClient()
    while True:
        try:
            for row in await asyncio.to_thread(store.queued):
                inp = {"id": row["id"], "actor": row["actor"], "spec": row["specification"]}
                try:
                    await asyncio.to_thread(client.schedule_new_workflow, workflow="onemind_mission_v1",
                                            input=inp, instance_id=row["id"])
                except Exception:
                    status = await asyncio.to_thread(client.get_workflow_state, row["id"])
                    if status is None:
                        continue
                await asyncio.to_thread(store.mark_scheduled, row["id"])
        except Exception:
            # Dependency outages retain queued work; credentials/results are never logged.
            pass
        await asyncio.sleep(5)


@asynccontextmanager
async def lifespan(app):
    await asyncio.to_thread(store.initialize)
    from .workflows import runtime
    await asyncio.to_thread(runtime.start)
    task = asyncio.create_task(dispatch_pending())
    from .events import publish_outbox
    events = asyncio.create_task(publish_outbox())
    app.state.ready = True
    try:
        yield
    finally:
        app.state.ready = False
        task.cancel()
        events.cancel()
        with contextlib.suppress(asyncio.CancelledError):
            await task
        with contextlib.suppress(asyncio.CancelledError):
            await events
        await asyncio.to_thread(runtime.shutdown)
        store.pool.close()


app = FastAPI(title="OneMind operational core", version="0.1.0", lifespan=lifespan,
              docs_url=None, redoc_url=None, openapi_url=None)


@app.get("/health")
def health():
    if not getattr(app.state, "ready", False):
        raise HTTPException(503, "Dependencies initializing")
    with store.pool.connection() as conn:
        conn.execute("SELECT 1")
    return {"state": "ready", "observedAt": datetime.now(timezone.utc).isoformat()}


@app.get("/capabilities")
def capabilities(p=Depends(principal)):
    return [{"name": name, **value} for name, value in CAPABILITIES.items()]


@app.get("/objects")
def objects(limit: int = Query(default=200, ge=1, le=500), kind: str | None = None, p=Depends(principal)):
    return {"objects": store.objects(limit, kind), "observedAt": datetime.now(timezone.utc).isoformat()}


@app.post("/objects")
def upsert_object(inp: ObjectInput, p=Depends(operator)):
    if inp.source != "onemind" or inp.id.startswith("source:"):
        raise HTTPException(403, "External source projections are maintained by their ingestion adapter")
    row = store.upsert_object(inp.model_dump())
    store.audit(p["actor"], "object-upsert", inp.id)
    return row


class LinkInput(BaseModel):
    source: str = Field(min_length=1, max_length=200)
    relation: str = Field(pattern=r"^[a-z][a-z0-9-]{0,49}$")
    target: str = Field(min_length=1, max_length=200)


@app.get("/links")
def links(p=Depends(principal)):
    return store.links()


@app.post("/links")
def add_link(inp: LinkInput, p=Depends(operator)):
    store.add_link(inp.source, inp.relation, inp.target)
    store.audit(p["actor"], "link-created", inp.source + "->" + inp.target)
    return {"state": "recorded"}


@app.get("/agents")
async def agents(p=Depends(principal)):
    try:
        return {"agents": await agents_and_call(), "observedAt": datetime.now(timezone.utc).isoformat()}
    except Exception:
        raise HTTPException(503, "Agent discovery unavailable")


@app.get("/missions")
def missions(p=Depends(principal)):
    return store.missions()


@app.get("/missions/{mid}")
def mission(mid: str, p=Depends(principal)):
    value = store.mission(mid)
    if value is None:
        raise HTTPException(404, "Mission not found")
    return value


@app.post("/missions", status_code=202)
def create_mission(inp: MissionInput, p=Depends(operator)):
    if inp.capability == "restart-worker" and tuple(inp.target.split("/", 1)) not in ALLOWED_RESTARTS:
        raise HTTPException(403, "Target is outside the worker restart capability")
    if inp.capability == "delegate-agent" and not inp.target:
        raise HTTPException(422, "Choose a discovered agent reference")
    spec = inp.model_dump(exclude={"idempotency_key"})
    digest = hashlib.sha256(json.dumps(spec, sort_keys=True).encode()).hexdigest()
    row, created = store.create_mission(str(uuid.uuid4()), p["actor"], inp.idempotency_key, digest, spec)
    if row["request_hash"] != digest:
        raise HTTPException(409, "Idempotency key was already used for another request")
    if created:
        store.audit(p["actor"], "mission-requested", inp.target, row["id"])
    return row


@app.post("/missions/{mid}/cancel")
async def cancel_mission(mid: str, p=Depends(operator)):
    row = await asyncio.to_thread(store.mission, mid)
    if row is None:
        raise HTTPException(404, "Mission not found")
    if row["state"] in {"completed", "failed", "cancelled"}:
        raise HTTPException(409, "Mission is already terminal")
    if row["scheduled"]:
        client = DaprWorkflowClient()
        await asyncio.to_thread(client.terminate_workflow, mid, output="Operator cancelled")
    await asyncio.to_thread(store.set_state, mid, "cancelled")
    await asyncio.to_thread(store.audit, p["actor"], "mission-cancelled", None, mid)
    return {"state": "cancelled", "note": "Already-dispatched external work may require its own cancellation"}
