#!/usr/bin/env python3
"""Helpers for replay.sh (FEATURES A1).

extract <deployment.yaml> <dir>   write every redpanda-connect-pipelines ConfigMap key to <dir>
check <fixtures> <logs>           assert on what `nats sub '>'` captured during replay.sh

The expected output is computed here from the fixtures, independently of the Bloblang under
test, so a mapping bug shows up as a mismatch rather than being copied into the expectation.
"""
from __future__ import annotations

import base64
import json
import re
import sys
from collections import Counter
from pathlib import Path

ALLOWED = [  # BUS-DESIGN §3.1 subjects these pipelines may publish
    re.compile(r"^evt\.world\.rail\.[a-z0-9_-]+$"),
    re.compile(r"^intel\.space\.celestrak\.tle$"),
    re.compile(r"^tel\.feed\.(mta|celestrak)\.status$"),
    re.compile(r"^evt\.sys\.(mta|celestrak)\.rejected$"),
]
# Request/reply plumbing that any JetStream publisher generates; not feed output.
PLUMBING = re.compile(r"^(_INBOX\.|\$JS\.)")


def extract(deployment: Path, out: Path) -> None:
    import yaml

    docs = [d for d in yaml.safe_load_all(deployment.read_text()) if d]
    cm = next(d for d in docs if d["kind"] == "ConfigMap" and d["metadata"]["name"] == "redpanda-connect-pipelines")
    out.mkdir(parents=True, exist_ok=True)
    for key, body in cm["data"].items():
        (out / key).write_text(body)
    print("extracted:", ", ".join(sorted(cm["data"])))


def expected_rail(fixtures: Path) -> tuple[set[str], int]:
    subjects, rejected = set(), 0
    for feed in ("mnr", "lirr"):
        for e in json.loads((fixtures / f"{feed}.json").read_text())["entity"]:
            v = e.get("vehicle") or {}
            pos = v.get("position")
            if pos is None:
                continue
            tu = e.get("tripUpdate") or {}
            trip = str((v.get("trip") or {}).get("tripId") or (tu.get("trip") or {}).get("tripId") or e["id"])
            token = re.sub(r"[^a-z0-9_-]", "_", trip.lower())
            lat, lon = float(pos.get("latitude", 0)), float(pos.get("longitude", 0))
            if lat == 0 or lon == 0 or not (-90 <= lat <= 90 and -180 <= lon <= 180):
                rejected += 1
            else:
                subjects.add(f"evt.world.rail.mta-{feed}-{token}")
    return subjects, rejected


def tle_checksum_ok(line: str) -> bool:
    if len(line) != 69:
        return False
    total = sum(1 if c == "-" else int(c) if c.isdigit() else 0 for c in line[:68])
    return line[68].isdigit() and total % 10 == int(line[68])


def expected_tle(fixtures: Path) -> tuple[dict[str, str], int]:
    lines = [l.strip() for l in (fixtures / "stations.tle").read_text().splitlines() if l.strip()]
    good, bad = {}, 0
    for i in range(len(lines) - 1):
        l1, l2 = lines[i], lines[i + 1]
        if l1.startswith("1 ") and l2.startswith("2 "):
            norad = l1[2:7].strip()
            if tle_checksum_ok(l1) and tle_checksum_ok(l2) and l2[2:7].strip() == norad:
                good[norad] = f"tle:{norad}:{l1[18:32].strip()}"
            else:
                bad += 1
    return good, bad


def load_dump(d: Path) -> list[dict]:
    msgs = []
    for f in sorted(d.glob("*.json"), key=lambda p: int(re.sub(r"\D", "", p.stem) or 0)):
        raw = json.loads(f.read_text())
        m = {k.lower(): v for k, v in raw.items()}  # nats CLI writes Subject/Header/Data
        data = m.get("data")
        if isinstance(data, str):
            try:
                data = base64.b64decode(data, validate=True)
            except Exception:
                data = data.encode()
        m["_bytes"] = data or b""
        msgs.append(m)
    return msgs


def header(m: dict, name: str) -> str | None:
    h = m.get("header") or m.get("headers") or {}
    v = h.get(name)
    if isinstance(v, list):
        v = v[0] if v else None
    return v


def check(fixtures: Path, logs: Path) -> int:
    fails: list[str] = []

    def expect(cond: bool, what: str) -> None:
        print(("PASS " if cond else "FAIL ") + what)
        if not cond:
            fails.append(what)

    rail_subjects, rail_rejected = expected_rail(fixtures)
    tle_good, tle_bad = expected_tle(fixtures)

    msgs = load_dump(logs / "grammar" / "dump")
    feed = [m for m in msgs if not PLUMBING.match(m["subject"])]
    counts = Counter(m["subject"] for m in feed)
    print(f"captured {len(msgs)} messages, {len(feed)} feed messages on {len(counts)} subjects")
    for s, n in sorted(counts.items()):
        print(f"  {n:4d}  {s}")

    off = sorted(s for s in counts if not any(p.match(s) for p in ALLOWED))
    expect(not off, f"only grammar subjects published (off-grammar: {off})")
    expect(not any(s.startswith("tak.cot.") for s in counts), "nothing published on tak.cot.* (incl. tak.cot.inject)")
    expect(len(feed) > 0, "the grammar run published something")

    seen_rail = {s for s in counts if s.startswith("evt.world.rail.")}
    expect(seen_rail == rail_subjects,
           f"evt.world.rail subjects match the fixtures ({len(seen_rail)} seen, {len(rail_subjects)} expected)")
    polls = counts.get("tel.feed.mta.status", 0)
    expect(polls >= 2, f"tel.feed.mta.status once per poll ({polls} polls)")
    rail_msgs = [m for m in feed if m["subject"].startswith("evt.world.rail.")]
    if not rail_msgs:
        expect(False, "rail records present for the envelope check")
    for m in rail_msgs:
        body = json.loads(m["_bytes"])
        miss = [k for k in ("uid", "type", "lat", "lon", "t") if k not in body]
        if miss or body["uid"] != m["subject"].split(".")[-1] or not header(m, "Nats-Msg-Id"):
            expect(False, f"rail record envelope on {m['subject']} (missing {miss}, msg-id {header(m, 'Nats-Msg-Id')})")
            break
    else:
        expect(True, "every rail record has the §9.2 envelope (uid = subject id, type, lat, lon, t) and a Nats-Msg-Id")

    tle_msgs = [m for m in feed if m["subject"] == "intel.space.celestrak.tle"]
    ids = {json.loads(m["_bytes"])["properties"]["norad_id"] for m in tle_msgs}
    expect(ids == set(tle_good), f"intel.space.celestrak.tle covers every valid TLE ({len(ids)} of {len(tle_good)})")
    expect(len(tle_msgs) > 0, f"TLE records present ({len(tle_msgs)})")
    ok_ids = all(header(m, "Nats-Msg-Id") == tle_good.get(json.loads(m["_bytes"])["properties"]["norad_id"]) for m in tle_msgs)
    expect(ok_ids, "each TLE carries Nats-Msg-Id tle:{norad}:{epoch}")
    expect(all(json.loads(m["_bytes"])["lat"] is None for m in tle_msgs), "TLE records carry no position (raw source data, no propagation)")

    rej_mta = [json.loads(m["_bytes"]) for m in feed if m["subject"] == "evt.sys.mta.rejected"]
    rej_cel = [json.loads(m["_bytes"]) for m in feed if m["subject"] == "evt.sys.celestrak.rejected"]
    expect(polls > 0 and len(rej_mta) == rail_rejected * polls and all(r.get("reason") for r in rej_mta),
           f"evt.sys.mta.rejected: {rail_rejected} per poll with a reason ({len(rej_mta)} over {polls} polls)")
    cel_polls = counts.get("tel.feed.celestrak.status", 0)
    expect(cel_polls >= 1 and len(rej_cel) == tle_bad * cel_polls and all(r.get("reason") for r in rej_cel),
           f"evt.sys.celestrak.rejected: {tle_bad} per fetch with a reason ({len(rej_cel)} over {cel_polls} fetches)")

    status = [json.loads(m["_bytes"]) for m in feed if m["subject"] == "tel.feed.mta.status"]
    expect(len(status) > 0 and all(s["lastError"] is None and s["stale"] is False and s["published"] == len(rail_subjects)
               and s["rejected"] == rail_rejected for s in status),
           "tel.feed.mta.status reports the published and rejected counts, no error, not stale")

    world = json.loads((logs / "grammar" / "stream-WORLD.json").read_text())["state"]["messages"]
    intel = json.loads((logs / "grammar" / "stream-INTEL.json").read_text())["state"]["messages"]
    cot = json.loads((logs / "grammar" / "stream-COT.json").read_text())["state"]["messages"]
    expect(world == len(rail_subjects),
           f"WORLD stored each train once across {polls} polls (Nats-Msg-Id dedup: {world} stored, {len(seen_rail)} trains)")
    expect(intel == len(tle_good), f"INTEL stored each TLE once ({intel})")
    expect(cot == 0, f"COT stream empty after the grammar run ({cot})")

    legacy = load_dump(logs / "legacy" / "dump")
    lc = Counter(m["subject"] for m in legacy if not PLUMBING.match(m["subject"]))
    print("negative control (legacy feed-mta):", dict(lc))
    expect(lc.get("tak.cot.inject", 0) > 0, "negative control: the capture does see legacy CoT on tak.cot.inject")

    print(f"\nRESULT: {'PASS' if not fails else 'FAIL'} ({len(fails)} failed)")
    return 1 if fails else 0


if __name__ == "__main__":
    if len(sys.argv) == 4 and sys.argv[1] == "extract":
        extract(Path(sys.argv[2]), Path(sys.argv[3]))
    elif len(sys.argv) == 4 and sys.argv[1] == "check":
        sys.exit(check(Path(sys.argv[2]), Path(sys.argv[3])))
    else:
        sys.exit(__doc__)
