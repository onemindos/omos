"""
OneMind OS — NATS → OmniDB Ingester
=========================================
Subscribes to ALL NATS JetStream subjects and archives events into
OmniDB (PostgreSQL + TimescaleDB + pgvector) for time-aware queries.
Uses durable consumers so no events are lost during restarts.

Usage:
  python nats_ingester.py

Environment:
  NATS_URL     — NATS server URL (default: nats://localhost:4222)
  OMNIDB_URL   — PostgreSQL/OmniDB connection URL
  LOG_LEVEL    — DEBUG, INFO, WARNING, ERROR (default: INFO)

What this enables for Legacy:
  - "What happened last summer?" → historical event queries
  - "When do I usually go to sleep?" → pattern detection via time_bucket
  - "Show my task completion trend" → metrics over time
  - Biometric aggregation via continuous aggregates
"""

import asyncio
import json
import logging
import os
import signal
import sys
from datetime import datetime, timezone

import nats
import psycopg2
from psycopg2.extras import Json

# Configuration
NATS_URL   = os.getenv("NATS_URL",   "nats://localhost:4222")
OMNIDB_URL = os.getenv("OMNIDB_URL", "postgresql://onemind:onemind-omnidb-2026@localhost:5432/onemind")
LOG_LEVEL  = os.getenv("LOG_LEVEL",  "INFO").upper()

# Subjects to subscribe to (> = wildcard for all sub-subjects)
SUBSCRIBE_SUBJECTS = [
    "digital.>",   # Agent tasks + results
    "physical.>",  # Robots, sensors, Home Assistant
    "nexus.>",     # System health, heartbeats
    "comms.>",     # Human ↔ agent messages
    "home.>",      # Home Assistant events
    "metrics.>",   # Explicit metrics
    "events.>",    # General events (biometric, location, etc.)
]

DURABLE_PREFIX = "tsdb-ingester"

logging.basicConfig(
    level=getattr(logging, LOG_LEVEL),
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger("nats-ingester")


# ─── Database ────────────────────────────────────────────────────────────────

def get_db_connection():
    conn = psycopg2.connect(OMNIDB_URL)
    conn.autocommit = False
    return conn


def insert_event(conn, subject: str, data: dict, source: str | None = None):
    cur = conn.cursor()
    cur.execute(
        "INSERT INTO events (timestamp, subject, data, source) VALUES (NOW(), %s, %s, %s)",
        (subject, Json(data), source),
    )
    conn.commit()


# ─── Main ────────────────────────────────────────────────────────────────────

async def main():
    logger.info("Starting NATS → TimescaleDB ingester")
    logger.info(f"  NATS: {NATS_URL}")
    logger.info(f"  OmniDB: {OMNIDB_URL.split('@')[1] if '@' in OMNIDB_URL else OMNIDB_URL}")

    # Wait for TimescaleDB
    conn = None
    for attempt in range(30):
        try:
            conn = get_db_connection()
            conn.cursor().execute("SELECT 1")
            logger.info("TimescaleDB connection established")
            break
        except Exception as e:
            logger.warning(f"Waiting for TimescaleDB (attempt {attempt + 1}/30): {e}")
            await asyncio.sleep(2)

    if conn is None:
        logger.error("Failed to connect to TimescaleDB after 30 attempts")
        sys.exit(1)

    nc = await nats.connect(NATS_URL)
    js = nc.jetstream()
    logger.info("Connected to NATS JetStream")

    ingested_count = 0

    async def handle_message(msg):
        nonlocal ingested_count, conn
        try:
            try:
                data = json.loads(msg.data.decode())
            except (json.JSONDecodeError, UnicodeDecodeError):
                data = {"raw": msg.data.decode("utf-8", errors="replace")}

            source = data.pop("source", None) or msg.subject.split(".")[0]
            insert_event(conn, msg.subject, data, source)
            await msg.ack()

            ingested_count += 1
            if ingested_count % 100 == 0:
                logger.info(f"Ingested {ingested_count} events")
            logger.debug(f"Ingested: {msg.subject} ({len(msg.data)} bytes)")

        except psycopg2.OperationalError:
            logger.warning("TimescaleDB connection lost, reconnecting...")
            try:
                conn = get_db_connection()
                insert_event(conn, msg.subject, data, source)
                await msg.ack()
            except Exception as e:
                logger.error(f"Failed to reconnect: {e}")
                await msg.nak()
        except Exception as e:
            logger.error(f"Error processing {msg.subject}: {e}")
            await msg.nak()

    subscriptions = []
    for subject in SUBSCRIBE_SUBJECTS:
        durable_name = f"{DURABLE_PREFIX}-{subject.replace('.', '-').replace('>', 'all')}"
        try:
            sub = await js.subscribe(subject, durable=durable_name, cb=handle_message, manual_ack=True)
            subscriptions.append(sub)
            logger.info(f"Subscribed: {subject} (durable: {durable_name})")
        except Exception as e:
            logger.warning(f"JetStream subscribe failed for {subject}: {e} — falling back to core NATS")
            try:
                sub = await nc.subscribe(subject)
                subscriptions.append(sub)
            except Exception as e2:
                logger.error(f"Core NATS subscribe also failed for {subject}: {e2}")

    logger.info(f"Ingester running — {len(subscriptions)} subscriptions active")

    stop = asyncio.Event()
    loop = asyncio.get_event_loop()
    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, lambda: stop.set())
    await stop.wait()

    logger.info("Shutting down...")
    for sub in subscriptions:
        await sub.unsubscribe()
    await nc.close()
    if conn:
        conn.close()
    logger.info(f"Ingester stopped. Total events ingested: {ingested_count}")


if __name__ == "__main__":
    asyncio.run(main())
