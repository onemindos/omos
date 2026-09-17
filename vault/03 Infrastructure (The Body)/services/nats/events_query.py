#!/usr/bin/env python3
"""
OneMind Events Query — for agent tool calls and manual inspection.

Usage examples:
  python events_query.py --recent 50
  python events_query.py --subject "comms.%" --since "2026-04-01" --limit 100
  python events_query.py --subject "home.%" --since "7d"
  python events_query.py "SELECT subject, timestamp, data FROM events WHERE subject LIKE 'comms.%' ORDER BY timestamp DESC LIMIT 20"
"""
import argparse
import json
import sqlite3
import sys
from datetime import datetime, timedelta, timezone

DB_PATH = "/opt/onemind/events.db"


def resolve_since(s: str) -> str:
    """Support shortcuts like 7d, 24h, 30m in addition to ISO timestamps."""
    if not s:
        return s
    if s.endswith("d"):
        return (datetime.now(timezone.utc) - timedelta(days=int(s[:-1]))).isoformat()
    if s.endswith("h"):
        return (datetime.now(timezone.utc) - timedelta(hours=int(s[:-1]))).isoformat()
    if s.endswith("m"):
        return (datetime.now(timezone.utc) - timedelta(minutes=int(s[:-1]))).isoformat()
    return s


def main():
    parser = argparse.ArgumentParser(description="Query OneMind events DB")
    parser.add_argument("sql",      nargs="?", help="Raw SQL query")
    parser.add_argument("--subject", help="Subject filter (SQL LIKE pattern, e.g. comms.%%)")
    parser.add_argument("--since",   help="Start time: ISO timestamp or shorthand (7d, 24h, 30m)")
    parser.add_argument("--until",   help="End time: ISO timestamp")
    parser.add_argument("--limit",   type=int, default=50)
    parser.add_argument("--recent",  type=int, help="N most recent events across all subjects")
    parser.add_argument("--count",   action="store_true", help="Return row count only")
    args = parser.parse_args()

    try:
        conn = sqlite3.connect(DB_PATH)
        conn.row_factory = sqlite3.Row
    except Exception as e:
        print(json.dumps({"error": str(e)}))
        sys.exit(1)

    if args.sql:
        sql, params = args.sql, []
    elif args.recent:
        sql    = "SELECT * FROM events ORDER BY timestamp DESC LIMIT ?"
        params = [args.recent]
    else:
        conditions, params = [], []
        if args.subject:
            conditions.append("subject LIKE ?")
            params.append(args.subject)
        if args.since:
            conditions.append("timestamp >= ?")
            params.append(resolve_since(args.since))
        if args.until:
            conditions.append("timestamp <= ?")
            params.append(args.until)
        where = ("WHERE " + " AND ".join(conditions)) if conditions else ""
        if args.count:
            sql = f"SELECT COUNT(*) as count FROM events {where}"
        else:
            sql = f"SELECT * FROM events {where} ORDER BY timestamp DESC LIMIT ?"
            params.append(args.limit)

    try:
        rows   = conn.execute(sql, params).fetchall()
        result = [dict(r) for r in rows]
        print(json.dumps(result, indent=2, default=str))
    except Exception as e:
        print(json.dumps({"error": str(e), "sql": sql}))
    finally:
        conn.close()


if __name__ == "__main__":
    main()
