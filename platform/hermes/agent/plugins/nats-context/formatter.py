"""Formats a LiveState snapshot into the [LIVE OPERATIONAL PICTURE] context block."""
from __future__ import annotations

from datetime import datetime, timezone
from typing import Any


def _ago(ts: float) -> str:
    delta = datetime.now(timezone.utc).timestamp() - ts
    if delta < 60:
        return f"{int(delta)}s ago"
    if delta < 3600:
        return f"{int(delta/60)}m ago"
    return f"{int(delta/3600)}h ago"


def _lat_lon(payload: dict) -> str:
    lat = payload.get("lat") or payload.get("latitude")
    lon = payload.get("lon") or payload.get("longitude") or payload.get("lng")
    alt = payload.get("alt") or payload.get("altitude") or payload.get("hae")
    if lat is not None and lon is not None:
        coord = f"{float(lat):.5f}°N, {float(lon):.5f}°E"
        if alt is not None:
            coord += f" alt={float(alt):.0f}m"
        return coord
    return ""


def _heading(payload: dict) -> str:
    h = payload.get("heading") or payload.get("course") or payload.get("track")
    s = payload.get("speed")
    parts = []
    if h is not None:
        parts.append(f"hdg={float(h):.0f}°")
    if s is not None:
        parts.append(f"spd={float(s):.1f}kn")
    return " ".join(parts)


def format_context(snapshot: dict, max_chars: int = 6000) -> str:
    lines: list[str] = ["[LIVE OPERATIONAL PICTURE]"]
    lines.append(f"Generated: {datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')}")
    lines.append("")

    # ── TAK / CoT tracks ──────────────────────────────────────────────────
    tak = snapshot.get("tak", {})
    if tak:
        lines.append(f"## TAK TRACKS ({len(tak)} active)")
        for callsign, p in sorted(tak.items()):
            ts = p.get("_ts", 0)
            parts = [f"  {callsign}"]
            coord = _lat_lon(p)
            if coord:
                parts.append(coord)
            hdg = _heading(p)
            if hdg:
                parts.append(hdg)
            cot_type = p.get("type") or p.get("cot_type")
            if cot_type:
                parts.append(f"type={cot_type}")
            team = p.get("team") or p.get("group")
            if team:
                parts.append(f"team={team}")
            parts.append(f"[{_ago(ts)}]")
            lines.append(" | ".join(parts))
        lines.append("")

    # ── Entities ──────────────────────────────────────────────────────────
    entities = snapshot.get("entities", {})
    if entities:
        total_ent = sum(len(v) for v in entities.values())
        lines.append(f"## ENTITIES ({total_ent} total)")
        for etype, items in sorted(entities.items()):
            lines.append(f"  [{etype.upper()}] {len(items)} items")
            for eid, p in sorted(items.items()):
                ts = p.get("_ts", 0)
                coord = _lat_lon(p)
                status = p.get("status") or p.get("state")
                label = p.get("name") or p.get("callsign") or eid
                row = f"    {label} (id={eid})"
                if coord:
                    row += f" @ {coord}"
                if status:
                    row += f" [{status}]"
                row += f" {_ago(ts)}"
                lines.append(row)
        lines.append("")

    # ── Agents ────────────────────────────────────────────────────────────
    agents = snapshot.get("agents", {})
    if agents:
        lines.append(f"## AGENTS ({len(agents)} alive)")
        for key, p in sorted(agents.items()):
            ts = p.get("_ts", 0)
            name = p.get("name") or p.get("agent") or key.split(".")[0]
            owner = p.get("owner") or (key.split(".")[1] if "." in key else "")
            session = p.get("session") or p.get("session_name") or ""
            status = p.get("status") or "hb"
            task = p.get("current_task") or p.get("task") or ""
            row = f"  {name}"
            if owner:
                row += f" (owner={owner})"
            if session:
                row += f" [{session}]"
            row += f" status={status}"
            if task:
                row += f" task={task!r}"
            row += f" {_ago(ts)}"
            lines.append(row)
        lines.append("")

    # ── Bridges ───────────────────────────────────────────────────────────
    bridges = snapshot.get("bridges", {})
    if bridges:
        lines.append(f"## BRIDGES ({len(bridges)} configured)")
        for name, p in sorted(bridges.items()):
            ts = p.get("_ts", 0)
            stale = p.get("_stale", False)
            msg_rate = p.get("msg_rate") or p.get("rate") or p.get("msgs_per_sec")
            connected = p.get("connected") or p.get("online") or p.get("active")
            row = f"  {name}: {'OFFLINE' if stale else ('LIVE' if connected else 'IDLE')}"
            if msg_rate is not None:
                row += f" {float(msg_rate):.1f}msg/s"
            row += f" {_ago(ts)}"
            lines.append(row)
        lines.append("")

    # ── Threats ───────────────────────────────────────────────────────────
    threats = snapshot.get("threats", {})
    if threats:
        lines.append(f"## ACTIVE THREATS ({len(threats)})")
        for tid, p in sorted(threats.items()):
            ts = p.get("_ts", 0)
            severity = p.get("severity") or p.get("level") or "UNKNOWN"
            label = p.get("label") or p.get("name") or p.get("type") or tid
            source = p.get("source") or p.get("origin") or ""
            coord = _lat_lon(p)
            row = f"  [{severity.upper()}] {label} (id={tid})"
            if source:
                row += f" src={source}"
            if coord:
                row += f" @ {coord}"
            row += f" {_ago(ts)}"
            lines.append(row)
        lines.append("")

    # ── Recent Intel ──────────────────────────────────────────────────────
    intel_items = snapshot.get("intel", [])
    if intel_items:
        recent = intel_items[:10]  # top 10 newest
        lines.append(f"## INTEL ({len(intel_items)} items, showing {len(recent)})")
        for item in recent:
            ts = item.get("_ts", 0)
            title = item.get("title") or item.get("name") or item.get("summary") or item.get("subject", "")
            source = item.get("source") or item.get("feed") or ""
            classification = item.get("classification") or item.get("class") or ""
            row = f"  {title}"
            if classification:
                row += f" [{classification}]"
            if source:
                row += f" (src={source})"
            row += f" {_ago(ts)}"
            lines.append(row)
        lines.append("")

    # ── Recent Detections ─────────────────────────────────────────────────
    detections = snapshot.get("detections", [])
    if detections:
        recent = detections[:5]
        lines.append(f"## DETECTIONS ({len(detections)} recent, showing {len(recent)})")
        for d in recent:
            ts = d.get("_ts", 0)
            det_type = d.get("type") or d.get("class") or d.get("subject", "det")
            confidence = d.get("confidence") or d.get("conf")
            coord = _lat_lon(d)
            row = f"  {det_type}"
            if confidence is not None:
                row += f" conf={float(confidence):.0%}"
            if coord:
                row += f" @ {coord}"
            row += f" {_ago(ts)}"
            lines.append(row)
        lines.append("")

    # ── Telemetry summary ─────────────────────────────────────────────────
    telemetry = snapshot.get("telemetry", {})
    if telemetry:
        lines.append(f"## TELEMETRY ({len(telemetry)} streams active)")
        for subj, readings in sorted(telemetry.items()):
            if not readings:
                continue
            latest = readings[0]
            ts = latest.get("_ts", 0)
            # Surface key numeric fields
            values = {
                k: v for k, v in latest.items()
                if k not in ("_ts", "subject") and isinstance(v, (int, float))
            }
            summary = ", ".join(f"{k}={v}" for k, v in list(values.items())[:4])
            lines.append(f"  {subj}: {summary} [{_ago(ts)}]")
        lines.append("")

    # ── Mission events ────────────────────────────────────────────────────
    mission_events = snapshot.get("mission_events", [])
    if mission_events:
        recent = mission_events[:5]
        lines.append(f"## MISSION EVENTS ({len(mission_events)}, showing {len(recent)})")
        for ev in recent:
            ts = ev.get("_ts", 0)
            ev_type = ev.get("type") or ev.get("event") or ev.get("action") or "event"
            name = ev.get("name") or ev.get("mission") or ""
            row = f"  {ev_type}"
            if name:
                row += f": {name}"
            row += f" {_ago(ts)}"
            lines.append(row)
        lines.append("")

    # ── Task events ───────────────────────────────────────────────────────
    task_events = snapshot.get("task_events", [])
    if task_events:
        recent = task_events[:5]
        lines.append(f"## TASK EVENTS ({len(task_events)}, showing {len(recent)})")
        for ev in recent:
            ts = ev.get("_ts", 0)
            ev_type = ev.get("type") or ev.get("event") or ev.get("action") or "event"
            name = ev.get("name") or ev.get("task") or ""
            row = f"  {ev_type}"
            if name:
                row += f": {name}"
            row += f" {_ago(ts)}"
            lines.append(row)
        lines.append("")

    # ── NATS stats ────────────────────────────────────────────────────────
    nats_stats = snapshot.get("nats_stats")
    if nats_stats:
        ts = nats_stats.get("_ts", 0)
        conns = nats_stats.get("connections") or nats_stats.get("num_connections") or nats_stats.get("total_connections")
        subs = nats_stats.get("subscriptions") or nats_stats.get("num_subscriptions")
        in_msgs = nats_stats.get("in_msgs")
        out_msgs = nats_stats.get("out_msgs")
        parts = []
        if conns is not None:
            parts.append(f"conns={conns}")
        if subs is not None:
            parts.append(f"subs={subs}")
        if in_msgs is not None:
            parts.append(f"in_msgs={in_msgs}")
        if out_msgs is not None:
            parts.append(f"out_msgs={out_msgs}")
        lines.append(f"## NATS SERVER  " + "  ".join(parts) + f"  [{_ago(ts)}]")

    result = "\n".join(lines)

    # Hard cap to prevent bloating the context window
    if len(result) > max_chars:
        result = result[:max_chars] + f"\n... [truncated at {max_chars} chars]"

    return result
