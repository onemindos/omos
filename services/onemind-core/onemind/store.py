"""App-owned operational state; external records retain their source authority."""
import json
import os
import uuid
from psycopg.rows import dict_row
from psycopg.types.json import Jsonb
from psycopg_pool import ConnectionPool

pool = ConnectionPool(os.environ.get("DATABASE_URL", ""), open=False,
                      min_size=1, max_size=8, kwargs={"row_factory": dict_row})

SCHEMA = """
CREATE TABLE IF NOT EXISTS objects (
  id text PRIMARY KEY, kind text NOT NULL, name text NOT NULL,
  properties jsonb NOT NULL DEFAULT '{}', source text NOT NULL,
  observed_at timestamptz NOT NULL DEFAULT now(), version bigint NOT NULL DEFAULT 1
);
CREATE TABLE IF NOT EXISTS links (
  source_id text NOT NULL REFERENCES objects(id), relation text NOT NULL,
  target_id text NOT NULL REFERENCES objects(id),
  PRIMARY KEY(source_id, relation, target_id)
);
CREATE TABLE IF NOT EXISTS missions (
  id text PRIMARY KEY, actor text NOT NULL, idempotency_key text NOT NULL,
  request_hash text NOT NULL, specification jsonb NOT NULL,
  state text NOT NULL DEFAULT 'queued', scheduled boolean NOT NULL DEFAULT false,
  result jsonb, error text, created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(), UNIQUE(actor, idempotency_key)
);
CREATE TABLE IF NOT EXISTS action_audit (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  actor text NOT NULL, action text NOT NULL, target text,
  mission_id text, recorded_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS event_outbox (
  id text PRIMARY KEY, subject text NOT NULL, payload jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(), published_at timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS single_agent_dispatch
  ON action_audit(mission_id) WHERE action='agent-dispatch';
"""


def initialize():
    pool.open(wait=True)
    with pool.connection() as conn:
        conn.execute(SCHEMA)


def objects(limit=200, kind=None):
    with pool.connection() as conn:
        return conn.execute("SELECT * FROM objects WHERE (%s::text IS NULL OR kind=%s) "
                            "ORDER BY observed_at DESC LIMIT %s", (kind, kind, limit)).fetchall()


def upsert_object(data):
    with pool.connection() as conn:
        return conn.execute("""INSERT INTO objects(id,kind,name,properties,source)
          VALUES(%s,%s,%s,%s,%s) ON CONFLICT(id) DO UPDATE SET
          kind=excluded.kind,name=excluded.name,properties=excluded.properties,
          source=excluded.source,observed_at=now(),version=objects.version+1
          RETURNING *""", (data["id"], data["kind"], data["name"],
                           Jsonb(data["properties"]), data["source"])).fetchone()


def audit(actor, action, target=None, mission_id=None):
    with pool.connection() as conn:
        conn.execute("INSERT INTO action_audit(actor,action,target,mission_id) VALUES(%s,%s,%s,%s)",
                     (actor, action, target, mission_id))


def missions(limit=100):
    with pool.connection() as conn:
        return conn.execute("SELECT * FROM missions ORDER BY created_at DESC LIMIT %s", (limit,)).fetchall()


def mission(mid):
    with pool.connection() as conn:
        return conn.execute("SELECT * FROM missions WHERE id=%s", (mid,)).fetchone()


def create_mission(mid, actor, key, request_hash, spec):
    with pool.connection() as conn:
        row = conn.execute("""INSERT INTO missions(id,actor,idempotency_key,request_hash,specification)
          VALUES(%s,%s,%s,%s,%s) ON CONFLICT(actor,idempotency_key) DO NOTHING RETURNING *""",
                           (mid, actor, key, request_hash, Jsonb(spec))).fetchone()
        if row:
            enqueue_event(conn, mid, actor, "requested")
            return row, True
        return conn.execute("SELECT * FROM missions WHERE actor=%s AND idempotency_key=%s",
                            (actor, key)).fetchone(), False


def queued():
    with pool.connection() as conn:
        return conn.execute("SELECT * FROM missions WHERE NOT scheduled AND state='queued' "
                            "ORDER BY created_at LIMIT 20").fetchall()


def mark_scheduled(mid):
    with pool.connection() as conn:
        conn.execute("UPDATE missions SET scheduled=true,updated_at=now() WHERE id=%s", (mid,))


def set_state(mid, state, result=None, error=None):
    with pool.connection() as conn:
        row = conn.execute("UPDATE missions SET state=%s,result=%s,error=%s,updated_at=now() WHERE id=%s "
                     "AND (state!='cancelled' OR %s='cancelled') AND state!=%s RETURNING actor",
                     (state, Jsonb(result, dumps=lambda value: json.dumps(value, default=str))
                      if result is not None else None, error, mid, state, state)).fetchone()
        if row:
            enqueue_event(conn, mid, row["actor"], state)


def enqueue_event(conn, mid, actor, state):
    event_id = str(uuid.uuid4())
    conn.execute("INSERT INTO event_outbox(id,subject,payload) VALUES(%s,%s,%s)",
                 (event_id, "evt.so.mission." + ("started" if state == "running" else state),
                  Jsonb({"schema_version": 1, "event_id": event_id,
                         "source_id": "onemind-core", "mission_id": mid,
                         "actor": actor, "state": state})))


def pending_events():
    with pool.connection() as conn:
        return conn.execute("SELECT * FROM event_outbox WHERE published_at IS NULL "
                            "ORDER BY created_at LIMIT 20").fetchall()


def event_sent(event_id):
    with pool.connection() as conn:
        conn.execute("UPDATE event_outbox SET published_at=now() WHERE id=%s", (event_id,))


def claim_dispatch(actor, target, mid):
    with pool.connection() as conn:
        return conn.execute("INSERT INTO action_audit(actor,action,target,mission_id) "
                            "VALUES(%s,'agent-dispatch',%s,%s) ON CONFLICT DO NOTHING RETURNING id",
                            (actor, target, mid)).fetchone() is not None


def dispatched(mid):
    with pool.connection() as conn:
        return conn.execute("SELECT 1 FROM action_audit WHERE mission_id=%s AND action='agent-dispatch' LIMIT 1",
                            (mid,)).fetchone() is not None


def links():
    with pool.connection() as conn:
        return conn.execute("SELECT * FROM links ORDER BY source_id,relation,target_id LIMIT 1000").fetchall()


def add_link(source, relation, target):
    with pool.connection() as conn:
        conn.execute("INSERT INTO links(source_id,relation,target_id) VALUES(%s,%s,%s) ON CONFLICT DO NOTHING",
                     (source, relation, target))
