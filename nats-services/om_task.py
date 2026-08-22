"""
OneMind OS — om-task: NATS-Native Task Engine
==============================================
Micro service registered via NATS $SRV API.
All state lives in NATS JetStream KV and streams — no SQLite.

Subjects handled:
  task.create               (request-reply)  — create a new task
  task.update.<task_id>     (publish)        — update task fields / status
  task.approve.<task_id>    (request-reply)  — approve or reject a task
  task.list                 (request-reply)  — list/filter tasks
  task.get.<task_id>        (request-reply)  — fetch one task by ID
  agents.hb.>               (subscribe)      — heartbeat aggregator → agent-registry KV

Internal:
  task.heartbeat (60s timer) — stale/overdue scan + self-heartbeat publish

KV buckets:
  tasks           (created here if absent — history=5, ttl=30d)
  agent-registry  (created here if absent — history=5)

Usage:
  python3 om_task.py

Env:
  NATS_URL (default nats://127.0.0.1:4222)
"""

import asyncio
import json
import logging
import os
import uuid
from datetime import datetime, timezone
from typing import Any, Optional

import nats
from nats.js.api import KeyValueConfig
from nats.micro import add_service
from nats.micro.request import Request
from nats.micro.service import ServiceConfig, EndpointConfig

# ─── Config ────────────────────────────────────────────────────────────────────

NATS_URL        = os.environ.get("NATS_URL", "nats://127.0.0.1:4222")
SERVICE_NAME    = "om-task"
SERVICE_VERSION = "1.0.0"
SERVICE_DESC    = "OneMind Task Engine — NATS-native, no SQLite"

TASKS_BUCKET    = "tasks"
REGISTRY_BUCKET = "agent-registry"

THIRTY_DAYS_NS  = 30 * 24 * 3600 * 1_000_000_000  # nanoseconds (nats-py KV TTL)
THIRTY_DAYS_S   = 30 * 24 * 3600                   # seconds (for float TTL param)

HB_INTERVAL     = 60   # seconds between internal heartbeat scans
STALE_MINUTES   = 10   # minutes before a pending task is considered stalled

LOG_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "logs")

# ─── Logging ───────────────────────────────────────────────────────────────────

def setup_logging() -> logging.Logger:
    os.makedirs(LOG_DIR, exist_ok=True)
    logger = logging.getLogger("om-task")
    logger.setLevel(logging.INFO)
    fmt = logging.Formatter("%(asctime)s %(levelname)s %(message)s")
    fh = logging.FileHandler(os.path.join(LOG_DIR, "om-task.log"))
    fh.setFormatter(fmt)
    ch = logging.StreamHandler()
    ch.setFormatter(fmt)
    logger.addHandler(fh)
    logger.addHandler(ch)
    return logger

log = setup_logging()

# ─── Helpers ───────────────────────────────────────────────────────────────────

def now_iso() -> str:
    return datetime.now(timezone.utc).isoformat()

def now_ts() -> float:
    return datetime.now(timezone.utc).timestamp()

def short_uuid() -> str:
    return uuid.uuid4().hex[:12]

def ok(data: dict) -> bytes:
    return json.dumps(data).encode()

def err(msg: str) -> bytes:
    return json.dumps({"error": msg}).encode()

async def kv_get_json(kv, key: str) -> Optional[dict]:
    """Fetch a KV entry and parse JSON. Returns None if missing."""
    try:
        entry = await kv.get(key)
        if entry and entry.value:
            return json.loads(entry.value.decode())
    except Exception:
        pass
    return None

async def kv_put_json(kv, key: str, data: dict) -> None:
    await kv.put(key, json.dumps(data).encode())

async def kv_all_keys(kv) -> list[str]:
    """Return all keys in a KV bucket, gracefully returning [] when empty."""
    try:
        return await kv.keys()
    except Exception:
        return []

# ─── KV Bootstrap ──────────────────────────────────────────────────────────────

async def ensure_kv(js, bucket: str, history: int = 5, ttl_seconds: Optional[float] = None):
    """Get or create a KV bucket. Returns the KeyValue handle."""
    try:
        kv = await js.key_value(bucket)
        log.info(f"KV bucket '{bucket}' connected")
        return kv
    except Exception:
        cfg = KeyValueConfig(bucket=bucket, history=history, ttl=ttl_seconds)
        kv = await js.create_key_value(cfg)
        log.info(f"KV bucket '{bucket}' created (history={history}, ttl={ttl_seconds}s)")
        return kv

# ─── Approval Gate Helper ──────────────────────────────────────────────────────

async def notify_approval_required(nc, task: dict) -> None:
    """Publish cmd.notify for tasks awaiting human approval."""
    task_id  = task.get("task_id", "?")
    title    = task.get("title", "(no title)")
    agent_id = task.get("agent_id", "unassigned")
    msg = {
        "message": (
            f"Task requires approval:\n"
            f"ID: {task_id}\n"
            f"Title: {title}\n"
            f"Agent: {agent_id}\n"
            f"Reply to task.approve.{task_id} with {{\"approved\": true/false}}"
        ),
        "title": "Task Approval Required",
        "priority": "high",
        "context": "agents",
    }
    await nc.publish("cmd.notify", json.dumps(msg).encode())
    log.info(f"Approval notification sent for {task_id}")

# ─── Subject Handlers ──────────────────────────────────────────────────────────

class TaskEngine:
    """All handler methods live here. Constructed once in run(), shared via closure."""

    def __init__(self, nc, js, tasks_kv, registry_kv):
        self.nc          = nc
        self.js          = js
        self.tasks_kv    = tasks_kv
        self.registry_kv = registry_kv

    # ── task.create ────────────────────────────────────────────────────────────

    async def handle_task_create(self, req: Request) -> None:
        """
        Payload: {title, description?, agent_id?, priority?, domain?, tags?,
                  due_ts?, requires_approval?}
        Reply:   {task_id, status: "created"}
        """
        try:
            data: dict = json.loads(req.data) if req.data else {}
        except Exception as e:
            await req.respond(err(f"invalid JSON: {e}"))
            return

        title = data.get("title", "").strip()
        if not title:
            await req.respond(err("title is required"))
            return

        task_id = f"task_{short_uuid()}"
        ts_now  = now_iso()

        task: dict[str, Any] = {
            "task_id":           task_id,
            "title":             title,
            "description":       data.get("description"),
            "agent_id":          data.get("agent_id"),
            "priority":          data.get("priority", "normal"),
            "domain":            data.get("domain"),
            "tags":              data.get("tags", []),
            "due_ts":            data.get("due_ts"),
            "requires_approval": data.get("requires_approval", False),
            "status":            "pending",
            "result":            None,
            "error":             None,
            "progress":          0,
            "created_at":        ts_now,
            "updated_at":        ts_now,
        }

        await kv_put_json(self.tasks_kv, task_id, task)
        log.info(f"task.create → {task_id} '{title}' priority={task['priority']}")

        # Persist creation event in TASKS stream
        await self.nc.publish(
            f"task.created.{task_id}",
            json.dumps({"task_id": task_id, "task": task, "ts": ts_now}).encode(),
        )

        # If approval gate is required, set status and notify
        if task["requires_approval"]:
            task["status"] = "pending_approval"
            task["updated_at"] = now_iso()
            await kv_put_json(self.tasks_kv, task_id, task)
            await notify_approval_required(self.nc, task)

        await req.respond(ok({"task_id": task_id, "status": task["status"]}))

    # ── task.update.<task_id> ──────────────────────────────────────────────────

    async def handle_task_update(self, req: Request) -> None:
        """
        Subject: task.update.<task_id>   (publish — no reply expected, but micro
                 wraps it so we always respond to avoid leaking reply subjects)
        Payload: {status?, agent_id?, result?, error?, progress?}
        """
        subject = req.subject  # task.update.<task_id>
        parts   = subject.split(".")
        if len(parts) < 3:
            await req.respond(err("malformed subject, expected task.update.<task_id>"))
            return

        task_id = parts[2]

        try:
            data: dict = json.loads(req.data) if req.data else {}
        except Exception as e:
            await req.respond(err(f"invalid JSON: {e}"))
            return

        task = await kv_get_json(self.tasks_kv, task_id)
        if task is None:
            await req.respond(err(f"task {task_id} not found"))
            return

        # Merge allowed fields
        for field in ("status", "agent_id", "result", "error", "progress"):
            if field in data:
                task[field] = data[field]
        task["updated_at"] = now_iso()

        await kv_put_json(self.tasks_kv, task_id, task)
        log.info(f"task.update → {task_id} status={task.get('status')} progress={task.get('progress')}")

        # Publish update event
        await self.nc.publish(
            f"task.updated.{task_id}",
            json.dumps({"task_id": task_id, "task": task, "ts": task["updated_at"]}).encode(),
        )

        new_status = task.get("status")
        agent_id   = task.get("agent_id")

        # Status-specific side effects
        if new_status == "blocked":
            if agent_id:
                await self.nc.publish(
                    f"evt.agent.{agent_id}.blocked",
                    json.dumps({"task_id": task_id, "agent_id": agent_id, "ts": now_iso()}).encode(),
                )
                log.info(f"evt.agent.{agent_id}.blocked published")

        elif new_status == "complete":
            await self.nc.publish(
                f"evt.task.{task_id}.complete",
                json.dumps({"task_id": task_id, "result": task.get("result"), "ts": now_iso()}).encode(),
            )
            log.info(f"evt.task.{task_id}.complete published")

        elif new_status == "pending_approval" and task.get("requires_approval"):
            await notify_approval_required(self.nc, task)

        await req.respond(ok({"task_id": task_id, "status": new_status, "ok": True}))

    # ── task.approve.<task_id> ─────────────────────────────────────────────────

    async def handle_task_approve(self, req: Request) -> None:
        """
        Subject: task.approve.<task_id>
        Payload: {approved: bool, approver?, reason?}
        Reply:   {ok: true}
        """
        subject = req.subject
        parts   = subject.split(".")
        if len(parts) < 3:
            await req.respond(err("malformed subject, expected task.approve.<task_id>"))
            return

        task_id = parts[2]

        try:
            data: dict = json.loads(req.data) if req.data else {}
        except Exception as e:
            await req.respond(err(f"invalid JSON: {e}"))
            return

        task = await kv_get_json(self.tasks_kv, task_id)
        if task is None:
            await req.respond(err(f"task {task_id} not found"))
            return

        approved  = bool(data.get("approved", False))
        approver  = data.get("approver", "unknown")
        reason    = data.get("reason", "")
        ts_now    = now_iso()

        if approved:
            task["status"] = "approved"
            evt_subject    = f"evt.task.{task_id}.approved"
        else:
            task["status"] = "rejected"
            evt_subject    = f"evt.task.{task_id}.rejected"

        task["approver"]   = approver
        task["reason"]     = reason
        task["updated_at"] = ts_now

        await kv_put_json(self.tasks_kv, task_id, task)
        log.info(f"task.approve → {task_id} approved={approved} by {approver}")

        await self.nc.publish(
            evt_subject,
            json.dumps({
                "task_id":  task_id,
                "approved": approved,
                "approver": approver,
                "reason":   reason,
                "ts":       ts_now,
            }).encode(),
        )

        await req.respond(ok({"ok": True, "task_id": task_id, "status": task["status"]}))

    # ── task.list ──────────────────────────────────────────────────────────────

    async def handle_task_list(self, req: Request) -> None:
        """
        Payload: {status?, agent_id?, domain?, limit?}
        Reply:   {tasks: [...], count: N}
        """
        try:
            filters: dict = json.loads(req.data) if req.data else {}
        except Exception:
            filters = {}

        f_status   = filters.get("status")
        f_agent_id = filters.get("agent_id")
        f_domain   = filters.get("domain")
        limit      = int(filters.get("limit", 200))

        keys = await kv_all_keys(self.tasks_kv)
        results = []

        for key in keys:
            task = await kv_get_json(self.tasks_kv, key)
            if task is None:
                continue
            if f_status   and task.get("status")   != f_status:
                continue
            if f_agent_id and task.get("agent_id") != f_agent_id:
                continue
            if f_domain   and task.get("domain")   != f_domain:
                continue
            results.append(task)
            if len(results) >= limit:
                break

        # Sort newest first
        results.sort(key=lambda t: t.get("created_at", ""), reverse=True)

        log.info(f"task.list → {len(results)} tasks (filters={filters})")
        await req.respond(ok({"tasks": results, "count": len(results)}))

    # ── task.get.<task_id> ─────────────────────────────────────────────────────

    async def handle_task_get(self, req: Request) -> None:
        """
        Subject: task.get.<task_id>
        Reply:   full task JSON or {error: "not found"}
        """
        subject = req.subject
        parts   = subject.split(".")
        if len(parts) < 3:
            await req.respond(err("malformed subject, expected task.get.<task_id>"))
            return

        task_id = parts[2]
        task    = await kv_get_json(self.tasks_kv, task_id)

        if task is None:
            log.info(f"task.get → {task_id} not found")
            await req.respond(err("not found"))
            return

        log.info(f"task.get → {task_id} status={task.get('status')}")
        await req.respond(ok(task))

    # ── agents.hb.> ────────────────────────────────────────────────────────────
    # This handler is wired as a plain NATS subscribe (not a micro endpoint)
    # because it uses a wildcard pattern that spans the heartbeat namespace.

    async def handle_agent_heartbeat(self, msg) -> None:
        """
        Subject: agents.hb.<harness>.<agent_id>
        Payload: {name, model?, status, harness, caps?}
        Action:  write to KV agent-registry key <harness>.<agent_id>
        """
        subject = msg.subject           # agents.hb.<harness>.<agent_id>
        parts   = subject.split(".")    # ["agents","hb","<harness>","<agent_id>"]

        if len(parts) < 4:
            log.warning(f"agents.hb malformed subject: {subject}")
            return

        harness  = parts[2]
        agent_id = parts[3]
        reg_key  = f"{harness}.{agent_id}"

        try:
            payload: dict = json.loads(msg.data) if msg.data else {}
        except Exception:
            payload = {}

        payload["last_seen"] = now_iso()
        payload.setdefault("harness",  harness)
        payload.setdefault("agent_id", agent_id)

        await kv_put_json(self.registry_kv, reg_key, payload)
        log.debug(f"agents.hb → registry[{reg_key}] last_seen updated")

    # ── Internal heartbeat scan (60s) ──────────────────────────────────────────

    async def run_heartbeat_loop(self) -> None:
        """Scan tasks for stalled/overdue items every HB_INTERVAL seconds."""
        while True:
            await asyncio.sleep(HB_INTERVAL)
            try:
                await self._scan_tasks()
            except Exception as e:
                log.error(f"heartbeat scan error: {e}")

    async def _scan_tasks(self) -> None:
        keys = await kv_all_keys(self.tasks_kv)
        task_count = len(keys)
        ts_now     = now_ts()
        stale_s    = STALE_MINUTES * 60

        stalled_n  = 0
        overdue_n  = 0

        for key in keys:
            task = await kv_get_json(self.tasks_kv, key)
            if task is None:
                continue

            task_id  = task.get("task_id", key)
            status   = task.get("status", "")
            agent_id = task.get("agent_id")
            due_ts   = task.get("due_ts")

            # Stalled: pending + has agent + not updated in >STALE_MINUTES
            if status == "pending" and agent_id:
                updated_at = task.get("updated_at", task.get("created_at", ""))
                try:
                    updated_dt = datetime.fromisoformat(updated_at)
                    # Ensure timezone-aware comparison
                    if updated_dt.tzinfo is None:
                        updated_dt = updated_dt.replace(tzinfo=timezone.utc)
                    age_s = ts_now - updated_dt.timestamp()
                    if age_s > stale_s:
                        await self.nc.publish(
                            f"evt.task.{task_id}.stalled",
                            json.dumps({
                                "task_id":  task_id,
                                "agent_id": agent_id,
                                "age_s":    int(age_s),
                                "ts":       now_iso(),
                            }).encode(),
                        )
                        log.info(f"Stalled task detected: {task_id} (age={int(age_s)}s)")
                        stalled_n += 1
                except Exception:
                    pass

            # Overdue: due_ts is set and in the past, task not terminal
            # DISABLED — fires every 60s with no dedup, creates runaway events
            # TODO: re-enable with smart escalation (1h, 24h, 48h, 72h intervals)
            # if due_ts and status not in ("complete", "approved", "rejected", "cancelled"):
            #     ...
            pass  # overdue alerts suppressed

        log.info(
            f"Heartbeat scan: {task_count} tasks, {stalled_n} stalled, {overdue_n} overdue"
        )

        # Self-heartbeat — write into the agent registry
        await self.nc.publish(
            "agents.hb.nats.om-task",
            json.dumps({
                "name":       SERVICE_NAME,
                "status":     "running",
                "harness":    "nats",
                "task_count": task_count,
                "ts":         now_iso(),
            }).encode(),
        )

# ─── Main ──────────────────────────────────────────────────────────────────────

async def run() -> None:
    log.info(f"om-task starting — NATS={NATS_URL}")

    # ── NATS connection with auto-reconnect ────────────────────────────────────
    async def error_cb(e):
        log.error(f"NATS error: {e}")

    async def disconnected_cb():
        log.warning("NATS disconnected")

    async def reconnected_cb():
        log.info("NATS reconnected")

    nc = await nats.connect(
        NATS_URL,
        error_cb=error_cb,
        disconnected_cb=disconnected_cb,
        reconnected_cb=reconnected_cb,
        max_reconnect_attempts=-1,
    )
    js = nc.jetstream()
    log.info("NATS connected")

    # ── KV buckets ─────────────────────────────────────────────────────────────
    tasks_kv    = await ensure_kv(js, TASKS_BUCKET,    history=5, ttl_seconds=THIRTY_DAYS_S)
    registry_kv = await ensure_kv(js, REGISTRY_BUCKET, history=5, ttl_seconds=None)

    engine = TaskEngine(nc, js, tasks_kv, registry_kv)

    # ── Register as NATS micro service ─────────────────────────────────────────
    svc = await add_service(
        nc,
        ServiceConfig(
            name=SERVICE_NAME,
            version=SERVICE_VERSION,
            description=SERVICE_DESC,
        ),
    )

    # Request-reply endpoints (micro framework — reply handled by req.respond)
    await svc.add_endpoint(
        EndpointConfig(
            name="task-create",
            subject="task.create",
            handler=engine.handle_task_create,
        )
    )

    # task.update.* — micro endpoint with wildcard subject
    await svc.add_endpoint(
        EndpointConfig(
            name="task-update",
            subject="task.update.*",
            handler=engine.handle_task_update,
        )
    )

    await svc.add_endpoint(
        EndpointConfig(
            name="task-approve",
            subject="task.approve.*",
            handler=engine.handle_task_approve,
        )
    )

    await svc.add_endpoint(
        EndpointConfig(
            name="task-list",
            subject="task.list",
            handler=engine.handle_task_list,
        )
    )

    await svc.add_endpoint(
        EndpointConfig(
            name="task-get",
            subject="task.get.*",
            handler=engine.handle_task_get,
        )
    )

    # ── agents.hb.> — plain subscribe (wildcard beyond one token level) ────────
    # nats.micro endpoints use a single wildcard (*); agents.hb uses multi-token
    # wildcard (>), so we use a regular core NATS subscribe here.
    await nc.subscribe("agents.hb.>", cb=engine.handle_agent_heartbeat)
    log.info("Subscribed to agents.hb.>")

    # ── Internal heartbeat task ────────────────────────────────────────────────
    hb_task = asyncio.create_task(engine.run_heartbeat_loop())
    log.info(f"Heartbeat loop started (interval={HB_INTERVAL}s, stale_threshold={STALE_MINUTES}m)")

    log.info(f"om-task v{SERVICE_VERSION} ready — registered as NATS micro service '{SERVICE_NAME}'")

    try:
        while True:
            await asyncio.sleep(1)
    except (asyncio.CancelledError, KeyboardInterrupt):
        pass
    finally:
        hb_task.cancel()
        try:
            await hb_task
        except asyncio.CancelledError:
            pass
        await svc.stop()
        await nc.drain()
        log.info("om-task stopped")


if __name__ == "__main__":
    asyncio.run(run())


# =============================================================================
# LaunchAgent plist (DO NOT CREATE — reference only)
# Save to: ~/Library/LaunchAgents/com.onemind.om-task.plist
# Load:    launchctl load ~/Library/LaunchAgents/com.onemind.om-task.plist
# Reload:  launchctl unload ~/Library/LaunchAgents/com.onemind.om-task.plist
#          launchctl load   ~/Library/LaunchAgents/com.onemind.om-task.plist
#
# <?xml version="1.0" encoding="UTF-8"?>
# <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
#     "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
# <plist version="1.0">
# <dict>
#     <key>Label</key>
#     <string>com.onemind.om-task</string>
#
#     <key>ProgramArguments</key>
#     <array>
#         <string>/opt/homebrew/bin/python3</string>
#         <string>/Users/onemind/subscribers/om_task.py</string>
#     </array>
#
#     <key>RunAtLoad</key>
#     <true/>
#
#     <key>KeepAlive</key>
#     <true/>
#
#     <key>WorkingDirectory</key>
#     <string>/Users/onemind/subscribers</string>
#
#     <key>StandardOutPath</key>
#     <string>/Users/onemind/subscribers/logs/om-task-stdout.log</string>
#
#     <key>StandardErrorPath</key>
#     <string>/Users/onemind/subscribers/logs/om-task-stderr.log</string>
#
#     <key>EnvironmentVariables</key>
#     <dict>
#         <key>NATS_URL</key>
#         <string>nats://127.0.0.1:4222</string>
#         <key>PATH</key>
#         <string>/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin</string>
#     </dict>
# </dict>
# </plist>
# =============================================================================
