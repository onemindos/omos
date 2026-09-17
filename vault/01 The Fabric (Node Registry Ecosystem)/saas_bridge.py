"""
OneMind OS — SaaS Bridge Skeleton (Todoist example)
Pattern: NATS cmd → HTTPS API call → NATS evt back

Subject law:
  cmd.service.todoist.<action>   ← Legacy publishes here
  evt.service.todoist.<event>    ← bridge publishes result here

One bridge process serves ALL Todoist commands.
To add a new SaaS app: copy this file, swap out the API calls.
"""

import asyncio
import json
import os
import httpx
import nats
from nats.aio.client import Client as NATS

NATS_URL     = "nats://localhost:4222"
TODOIST_TOKEN = os.environ.get("TODOIST_TOKEN", "")   # never hardcode
BASE_URL     = "https://api.todoist.com/api/v1"

HEADERS = {
    "Authorization": f"Bearer {TODOIST_TOKEN}",
    "Content-Type": "application/json",
}


# ── helpers ──────────────────────────────────────────────────────────────────

def ok(data: dict) -> bytes:
    return json.dumps({"ok": True, "data": data}).encode()

def err(msg: str) -> bytes:
    return json.dumps({"ok": False, "error": msg}).encode()


# ── handlers (one per command) ───────────────────────────────────────────────

async def handle_create_task(nc: NATS, msg):
    """
    Incoming NATS msg on: cmd.service.todoist.create_task
    Payload: { "id": "...", "source": "...", "payload": { "content": "...", "priority": 4 } }
    """
    try:
        body    = json.loads(msg.data)
        payload = body.get("payload", {})

        async with httpx.AsyncClient() as client:
            r = await client.post(
                f"{BASE_URL}/tasks",
                headers=HEADERS,
                json={
                    "content":    payload.get("content", "untitled"),
                    "priority":   payload.get("priority", 1),
                    "project_id": payload.get("project_id"),
                },
                timeout=10,
            )
            r.raise_for_status()
            task = r.json()

        # publish result back onto bus
        await nc.publish("evt.service.todoist.task_created", ok(task))

        # if caller wants a reply (request/reply pattern)
        if msg.reply:
            await nc.publish(msg.reply, ok(task))

    except Exception as e:
        if msg.reply:
            await nc.publish(msg.reply, err(str(e)))


async def handle_complete_task(nc: NATS, msg):
    """
    Incoming NATS msg on: cmd.service.todoist.complete_task
    Payload: { "payload": { "task_id": "..." } }
    """
    try:
        body    = json.loads(msg.data)
        task_id = body["payload"]["task_id"]

        async with httpx.AsyncClient() as client:
            r = await client.post(
                f"{BASE_URL}/tasks/{task_id}/close",
                headers=HEADERS,
                timeout=10,
            )
            r.raise_for_status()

        await nc.publish("evt.service.todoist.task_completed",
                         ok({"task_id": task_id}))
        if msg.reply:
            await nc.publish(msg.reply, ok({"task_id": task_id}))

    except Exception as e:
        if msg.reply:
            await nc.publish(msg.reply, err(str(e)))


async def handle_list_tasks(nc: NATS, msg):
    """
    Incoming NATS msg on: cmd.service.todoist.list_tasks
    Payload: { "payload": { "project_id": "..." } }
    """
    try:
        body       = json.loads(msg.data)
        project_id = body.get("payload", {}).get("project_id")
        params     = {"project_id": project_id} if project_id else {}

        async with httpx.AsyncClient() as client:
            r = await client.get(
                f"{BASE_URL}/tasks",
                headers=HEADERS,
                params=params,
                timeout=10,
            )
            r.raise_for_status()
            data = r.json()

        if msg.reply:
            await nc.publish(msg.reply, ok(data))

    except Exception as e:
        if msg.reply:
            await nc.publish(msg.reply, err(str(e)))


# ── main ─────────────────────────────────────────────────────────────────────

async def main():
    nc = await nats.connect(NATS_URL)
    print("[todoist-bridge] connected to NATS")

    # wire up every command subject to its handler
    await nc.subscribe("cmd.service.todoist.create_task",
                       cb=lambda msg: asyncio.create_task(handle_create_task(nc, msg)))
    await nc.subscribe("cmd.service.todoist.complete_task",
                       cb=lambda msg: asyncio.create_task(handle_complete_task(nc, msg)))
    await nc.subscribe("cmd.service.todoist.list_tasks",
                       cb=lambda msg: asyncio.create_task(handle_list_tasks(nc, msg)))

    print("[todoist-bridge] subscribed — waiting for commands")
    await asyncio.Event().wait()   # run forever


if __name__ == "__main__":
    asyncio.run(main())


# ── end-to-end flow ──────────────────────────────────────────────────────────
#
#  Legacy publishes:
#    cmd.service.todoist.create_task
#    { "id": "ulid", "source": "agent.legacy", "payload": { "content": "Call John", "priority": 4 } }
#
#  Bridge receives that message
#  Bridge calls: POST https://api.todoist.com/api/v1/tasks
#  Todoist creates task, returns task JSON
#
#  Bridge publishes back:
#    evt.service.todoist.task_created  { "ok": true, "data": { "id": "task123", ... } }
#
#  Any subscriber on the bus (dashboard, logger, notifier, another agent) reacts.
#
#  To add Gmail/Fastmail: copy this file, swap BASE_URL + HEADERS + handlers.
#  To add GHL/Skool:      copy this file, swap to their API.
#  The bus wiring never changes.
