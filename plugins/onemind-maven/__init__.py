"""
OneMind Maven — sovereign AI intelligence fusion layer for the OneMind fabric.
"""
from __future__ import annotations

import json
import ssl
from datetime import datetime, timezone
from typing import Any, Dict

import httpx

CLICKHOUSE_URL = "http://chi-onemind-onemind-0-0-0.clickhouse.svc.cluster.local:8123"


def _query_clickhouse(query: str) -> list[dict]:
    """Run a ClickHouse query and return parsed JSON rows."""
    with httpx.Client(timeout=30.0) as client:
        resp = client.get(
            CLICKHOUSE_URL,
            params={"query": query, "default_format": "JSONEachRow"},
        )
        resp.raise_for_status()
        lines = [ln for ln in resp.text.splitlines() if ln.strip()]
        return [json.loads(ln) for ln in lines]


def handle_vision_query(args: Dict[str, Any]) -> str:
    """Query real-time object detections from the vision pipeline."""
    feed = args.get("feed_id", "all")
    obj_class = args.get("object_class")
    min_conf = args.get("min_confidence", 0.5)
    minutes = args.get("last_n_minutes", 30)

    where = [f"timestamp > now() - INTERVAL {int(minutes)} MINUTE"]
    where.append(f"confidence >= {float(min_conf)}")
    if feed != "all":
        where.append(f"feed_id = '{feed}'")
    if obj_class:
        where.append(f"object_class = '{obj_class}'")

    q = f"""
    SELECT timestamp, feed_id, object_class, confidence,
           lat, lon, inference_ms
    FROM onemind.vision_detections
    WHERE {' AND '.join(where)}
    ORDER BY timestamp DESC
    LIMIT 100
    """
    rows = _query_clickhouse(q)
    return json.dumps(
        {"query_params": args, "result_count": len(rows), "detections": rows},
        default=str,
    )


def handle_analyze_feed(args: Dict[str, Any]) -> str:
    """Deep analysis of a specific feed with reasoning."""
    feed_id = args["feed_id"]
    question = args.get("question", "What is currently being detected?")

    q = f"""
    SELECT object_class, confidence, timestamp, lat, lon
    FROM onemind.vision_detections
    WHERE feed_id = '{feed_id}'
      AND timestamp > now() - INTERVAL 1 HOUR
    ORDER BY timestamp DESC
    LIMIT 20
    """
    rows = _query_clickhouse(q)
    return json.dumps(
        {
            "feed_id": feed_id,
            "question": question,
            "recent_detections": rows,
            "note": "Use LLM to reason over these detections + question",
        },
        default=str,
    )


def handle_sitrep(args: Dict[str, Any]) -> str:
    """Generate a Maven-style Situation Report."""
    hours = args.get("timeframe_hours", 1)

    events_q = f"""
    SELECT source, entity_type, count() as cnt, max(timestamp) as latest
    FROM onemind.events
    WHERE timestamp > now() - INTERVAL {int(hours)} HOUR
    GROUP BY source, entity_type
    ORDER BY cnt DESC
    """
    events = _query_clickhouse(events_q)

    vision_q = f"""
    SELECT feed_id, object_class, count() as detections,
           avg(confidence) as avg_conf, max(timestamp) as latest
    FROM onemind.vision_detections
    WHERE timestamp > now() - INTERVAL {int(hours)} HOUR
    GROUP BY feed_id, object_class
    ORDER BY detections DESC
    """
    vision = _query_clickhouse(vision_q)

    return json.dumps(
        {
            "sitrep_time": datetime.now(timezone.utc).isoformat(),
            "timeframe_hours": hours,
            "intel_events": events,
            "vision_detections": vision,
            "note": "Synthesize into narrative SITREP with threat assessment",
        },
        default=str,
    )


def handle_anomaly_check(args: Dict[str, Any]) -> str:
    """Detect anomalies vs historical baseline."""
    days = args.get("baseline_days", 7)
    sensitivity = args.get("sensitivity", "medium")
    thresholds = {"low": 3.0, "medium": 2.0, "high": 1.5}
    z_thresh = thresholds.get(sensitivity, 2.0)

    q = f"""
    WITH current AS (
        SELECT object_class, feed_id, count() as current_count
        FROM onemind.vision_detections
        WHERE timestamp > now() - INTERVAL 1 HOUR
        GROUP BY object_class, feed_id
    ),
    baseline AS (
        SELECT object_class, feed_id,
               avg(hourly_count) as avg_count,
               stddevPop(hourly_count) as std_count
        FROM (
            SELECT object_class, feed_id,
                   toStartOfHour(timestamp) as hour,
                   count() as hourly_count
            FROM onemind.vision_detections
            WHERE timestamp > now() - INTERVAL {int(days)} DAY
              AND timestamp < now() - INTERVAL 2 HOUR
            GROUP BY object_class, feed_id, hour
        )
        GROUP BY object_class, feed_id
    )
    SELECT c.object_class, c.feed_id,
           c.current_count, b.avg_count, b.std_count,
           (c.current_count - coalesce(b.avg_count, 0)) /
               greatest(coalesce(b.std_count, 1), 1) as z_score
    FROM current c
    LEFT JOIN baseline b ON c.object_class = b.object_class
                         AND c.feed_id = b.feed_id
    WHERE abs((c.current_count - coalesce(b.avg_count, 0)) /
              greatest(coalesce(b.std_count, 1), 1)) > {z_thresh}
    ORDER BY z_score DESC
    """
    anomalies = _query_clickhouse(q)
    return json.dumps(
        {
            "anomalies": anomalies,
            "baseline_days": days,
            "sensitivity": sensitivity,
            "z_threshold": z_thresh,
            "note": "z_score > threshold = anomalous activity",
        },
        default=str,
    )


def handle_recommend(args: Dict[str, Any]) -> str:
    """Recommend courses of action."""
    situation = args["situation"]
    assets = {
        "cameras": ["takcam", "helmcam", "home_front", "home_rear"],
        "operators": ["ATAK on S24", "CloudTAK browser", "tak3d 3D view"],
        "robots": ["Unitree Go2 (when online)"],
        "alerting": [
            "NATS pub to all operators",
            "TAK CoT alert markers",
            "Telegram notification",
        ],
        "intelligence": [
            "ClickHouse 54K+ events",
            "8 OSINT feed pods",
            "YOLOv5 vision pipeline",
        ],
        "compute": [
            "Legacy (Bedrock Claude LLM)",
            "takml-vision-worker (YOLOv5s CPU)",
            "Node-RED automation",
        ],
    }
    return json.dumps(
        {
            "situation": situation,
            "available_assets": assets,
            "instruction": (
                "Reason over the situation and available assets to recommend "
                "2-3 courses of action ranked by effectiveness."
            ),
        },
        default=str,
    )


def register(ctx):
    ctx.register_tool(
        toolset="onemind-maven",
        name="maven_vision_query",
        schema={
            "name": "maven_vision_query",
            "description": (
                "Query real-time object detections from all OneMind "
                "video feeds (takcam, helmcam, etc). Returns objects detected, "
                "confidence, timestamps, and GPS coordinates."
            ),
            "parameters": {
                "type": "object",
                "properties": {
                    "feed_id": {
                        "type": "string",
                        "description": "Filter by feed: takcam, helmcam, home_front, drone, robot, or all",
                        "default": "all",
                    },
                    "object_class": {
                        "type": "string",
                        "description": "Filter by detection class: pedestrian, car, van, truck, drone, etc.",
                    },
                    "min_confidence": {
                        "type": "number",
                        "description": "Minimum confidence 0-1 (default: 0.5)",
                        "default": 0.5,
                    },
                    "last_n_minutes": {
                        "type": "number",
                        "description": "How far back to look (default: 30)",
                        "default": 30,
                    },
                },
            },
        },
        handler=handle_vision_query,
    )

    ctx.register_tool(
        toolset="onemind-maven",
        name="maven_analyze_feed",
        schema={
            "name": "maven_analyze_feed",
            "description": (
                "Deep analysis of a specific camera feed with LLM reasoning. "
                "Pulls recent detections, optionally grabs a frame, and reasons "
                "over the scene."
            ),
            "parameters": {
                "type": "object",
                "properties": {
                    "feed_id": {
                        "type": "string",
                        "description": "Feed identifier to analyze",
                    },
                    "question": {
                        "type": "string",
                        "description": "What to look for or assess",
                    },
                },
                "required": ["feed_id"],
            },
        },
        handler=handle_analyze_feed,
    )

    ctx.register_tool(
        toolset="onemind-maven",
        name="maven_sitrep",
        schema={
            "name": "maven_sitrep",
            "description": (
                "Generate a Maven-style Situation Report (SITREP) fusing all "
                "OneMind intelligence: OSINT feeds, vision detections, TAK CoT, "
                "and historical patterns."
            ),
            "parameters": {
                "type": "object",
                "properties": {
                    "timeframe_hours": {
                        "type": "number",
                        "description": "Hours back to analyze (default: 1)",
                        "default": 1,
                    },
                    "focus_area": {
                        "type": "string",
                        "description": "Geographic focus or subject filter",
                    },
                },
            },
        },
        handler=handle_sitrep,
    )

    ctx.register_tool(
        toolset="onemind-maven",
        name="maven_anomaly_check",
        schema={
            "name": "maven_anomaly_check",
            "description": (
                "Check for anomalies by comparing current detections against "
                "historical baselines in ClickHouse. Flags unusual activity."
            ),
            "parameters": {
                "type": "object",
                "properties": {
                    "baseline_days": {
                        "type": "integer",
                        "description": "Days of history for baseline (default: 7)",
                        "default": 7,
                    },
                    "sensitivity": {
                        "type": "string",
                        "enum": ["low", "medium", "high"],
                        "description": "Alert threshold (default: medium)",
                        "default": "medium",
                    },
                },
            },
        },
        handler=handle_anomaly_check,
    )

    ctx.register_tool(
        toolset="onemind-maven",
        name="maven_recommend",
        schema={
            "name": "maven_recommend",
            "description": (
                "Recommend courses of action using available OneMind assets "
                "(cameras, robots, operators, alerts) given a situation."
            ),
            "parameters": {
                "type": "object",
                "properties": {
                    "situation": {
                        "type": "string",
                        "description": "Description of the situation to assess",
                    },
                },
                "required": ["situation"],
            },
        },
        handler=handle_recommend,
    )
