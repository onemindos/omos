#!/usr/bin/env python3
"""Generate testdata/fixtures/replay.jsonl.

Each line is {"subject": ..., "data": ...}. Payloads use the exact wire formats of the live omos
publishers, so the golden test exercises the real adapters:
  - feed-aircraft CoT: platform/data/redpanda-connect/deployment.yaml (feed-aircraft.yaml)
  - feed-mta CoT:      platform/data/redpanda-connect/deployment.yaml (feed-mta.yaml)
  - feed-satellites:   platform/omos/feed-scripts.yaml (sat_to_cot)
  - feed-intel:        platform/omos/intel-scripts.yaml (build_payload, USGS and EMSC)
  - geo-envelope:      docs/NATS-DOCTRINE.md section 9.2 (tel.fleet.*, evt.world.*)
  - projector echo:    platform/omos/tak-projector.yaml (entity_to_cot)
They are synthetic values in recorded formats, not captures from the live bus (the cluster is
read-only for this work and its bus credentials were not used).
Run: python3 gen_fixtures.py > replay.jsonl            (golden tests: fixed T0)
     python3 gen_fixtures.py <t0-epoch-ms> > live.jsonl  (a live run: T0 just before now)
"""
import json
import sys
from datetime import datetime, timezone

# epoch ms; contract example time (entity.md section 4.2). The golden tests use this value and a
# fake engine clock; a run of the real binary needs T0 shortly before the wall clock, or the
# 5 s future-skew rule replaces the source times.
T0 = int(sys.argv[1]) if len(sys.argv) > 1 else 1793275200000
XML = '<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>\n'


def iso(ms):
    return datetime.fromtimestamp(ms / 1000, tz=timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.000Z")


def adsb(icao, cs, lat, lon, alt, kt, trk, t):
    return (XML + f'<event version="2.0" type="a-f-A-M-F-Q" uid="ADSB-{icao.upper()}" how="m-g" time="{iso(t)}" start="{iso(t)}" stale="{iso(t + 60000)}">'
            f'<point lat="{lat}" lon="{lon}" hae="{alt}" ce="9999999.0" le="9999999.0"/><detail><contact callsign="{cs}"/>'
            f'<track speed="{kt}" course="{trk}"/><status readiness="false"/><remarks>ICAO:{icao} United States</remarks>'
            f'<marti><dest mission="Air"/></marti></detail></event>')


def mta(feed, trip, route, lat, lon, spd, crs, t):
    return (XML + f'<event version="2.0" type="a-n-G-E-V-C-R" uid="MTA-{feed}-{trip}" how="m-g" time="{iso(t)}" start="{iso(t)}" stale="{iso(t + 120000)}">'
            f'<point lat="{lat}" lon="{lon}" hae="0" ce="50" le="50"/><detail><contact callsign="{feed}-{route}"/>'
            f'<track speed="{spd}" course="{crs}"/><remarks>MTA route:{route} trip:{trip}</remarks><marti><dest mission="Land"/></marti></detail></event>')


def sat(num, name, lat, lon, alt_km, t):
    return (XML + f'<event version="2.0" type="a-f-A-M-F-Q" uid="SAT-{num}" how="m-g" time="{iso(t)}" start="{iso(t)}" stale="{iso(t + 90000)}">'
            f'<point lat="{lat:.4f}" lon="{lon:.4f}" hae="{alt_km*1000:.0f}" ce="9999999.0" le="9999999.0"/><detail>'
            f'<contact callsign="{name}"/><status readiness="false"/><remarks>NORAD:{num} group:stations</remarks>'
            f'<marti><dest mission="Air"/></marti><_flow-tags_ TAK-Server="SAT-Feed" time="{iso(t)}"/></detail></event>')


def eud(uid, cs, ctype, lat, lon, t, team="Cyan"):
    return (f'<event version="2.0" uid="{uid}" type="{ctype}" how="m-g" time="{iso(t)}" start="{iso(t)}" stale="{iso(t + 360000)}">'
            f'<point lat="{lat}" lon="{lon}" hae="210.5" ce="9.9" le="9999999.0"/><detail>'
            f'<takv device="Pixel 8" platform="ATAK-CIV" os="34" version="5.4.0"/><contact callsign="{cs}" endpoint="*:-1:stcp"/>'
            f'<__group name="{team}" role="Team Member"/><status battery="88"/><track speed="0.0" course="0.0"/></detail></event>')


def mesh_sim(t, lat, lon):
    return (f'<event version="2.0" uid="sim-x1" type="a-f-A-M-H-Q" how="m-s" time="{iso(t)}" start="{iso(t)}" stale="{iso(t + 60000)}">'
            f'<point lat="{lat}" lon="{lon}" hae="152.4" ce="5" le="3"/><detail><contact callsign="SIM-X1"/>'
            f'<track speed="12.5" course="90"/></detail></event>')


def echo(t):
    return (XML + f'<event version="2.0" type="a-f-A-C-F" uid="OM.aircraft.icao-a1b2c3" how="m-g" time="{iso(t)}" start="{iso(t)}" stale="{iso(t+60000)}" access="Unclassified">'
            f'<point lat="40.640000" lon="-73.780000" hae="1200" ce="50" le="50"/><detail><_flow-tags_ omos-tak-projector="{iso(t)}"/>'
            f'<contact callsign="UAL123"/><__group name="__ANON__" role="Team Member"/></detail></event>')


def geo(uid, typ, lat, lon, hae, course, speed, cs, stale, t, meta):
    return json.dumps({"uid": uid, "type": typ, "lat": lat, "lon": lon, "hae": hae, "course": course,
                       "speed": speed, "callsign": cs, "stale": stale, "t": t, "meta": meta})


out = []
add = lambda s, d: out.append({"subject": s, "data": d})

# Aircraft a1b2c3 from two feeds (FEATURES A4 / entity.md 10 test 1), interleaved.
for i in range(3):
    t = T0 + i * 10000
    add("tak.cot.inject", adsb("a1b2c3", "UAL123", round(40.64 + i * 0.01, 5), round(-73.78 + i * 0.012, 5), 1200 + i * 100, 250.0, 45.0, t))
    add("evt.world.air.icao-a1b2c3", geo("icao-a1b2c3", "air.commercial", round(40.6402 + i * 0.01, 5), round(-73.7798 + i * 0.012, 5),
                                         1210 + i * 100, 45, 128.6, "UAL123", 30, t + 2000, {"icao": "A1B2C3", "origin_iata": "JFK"}))
# A second aircraft from CoT only.
add("tak.cot.inject", adsb("abc123", "DAL45", 40.7, -73.9, 3000, 300.0, 270.0, T0 + 5000))
# Projector echo: must be ignored.
add("tak.cot.inject", echo(T0 + 6000))
# Vessel: position updates in place (entity.md 10 test "a ship's position updates in place").
for i in range(3):
    add("evt.world.sea.mmsi-367000001", geo("mmsi-367000001", "sea.cargo", round(40.60 + i * 0.0005, 6), -74.04, 0, 0, 5.5,
                                            "EVER GIVEN", 600, T0 + i * 10000, {"mmsi": "367000001", "flag": "US"}))
# Duplicate MMSI 200 km away 5 s later: split, never averaged (fusion-and-tracking.md Stage 1).
add("evt.world.sea.mmsi-367000001", geo("mmsi-367000001", "sea.cargo", 42.40, -74.04, 0, 0, 5.0, "SPOOF", 600, T0 + 25000,
                                        {"mmsi": "367000001"}))
# Satellites (feed-satellites).
add("tak.cot.inject", sat(25544, "ISS (ZARYA)", 51.2, -10.5, 418, T0))
add("tak.cot.inject", sat(25544, "ISS (ZARYA)", 50.1, 2.9, 418, T0 + 30000))
# MTA trains (feed-mta).
add("tak.cot.inject", mta("MNR", "1234_5678", "1", 41.03, -73.62, 20.5, 90, T0 + 1000))
add("tak.cot.inject", mta("MNR", "1234_5678", "1", 41.031, -73.615, 21.0, 90, T0 + 31000))
add("tak.cot.inject", mta("LIRR", "GO103_25_2611", "2", 40.75, -73.42, 0, 0, T0 + 2000))
# TAK EUDs on tak.cot.out (team members, person entities).
add("tak.cot.out", eud("ANDROID-a1b2c3d4e5f6", "CYAN-1", "a-f-G-U-C", 30.2675, -97.7429, T0 + 3000))
add("tak.cot.out", eud("ANDROID-a1b2c3d4e5f6", "CYAN-1", "a-f-G-U-C", 30.2677, -97.7425, T0 + 33000))
add("tak.cot.out", eud("S-1-5-21-942292099.WINTAK", "GREEN-2", "a-f-G-U-C", 30.27, -97.75, T0 + 4000, team="Green"))
# Stop-line test 9: person with hostile affiliation is rejected.
add("tak.cot.out", eud("ANDROID-badf00d", "BAD-1", "a-h-G-U-C", 30.26, -97.74, T0 + 5000))
# Own simulated drone: MAVLink telemetry plus the same uid from the mesh gateway.
for i in range(3):
    t = T0 + 20000 + i * 5000
    add("tel.fleet.drone.sim-x1.position", geo("sim-x1", "fleet.drone", 30.2672, round(-97.7431 + i * 0.0006, 6), 152.4, 90, 12.5,
                                               "SIM-X1", 60, t, {"battery_pct": 82 - i, "mode": "AUTO", "link_rssi": -67,
                                                                 "node_id": "build-mac", "mavlink_sysid": 1}))
add("tel.fleet.drone.sim-x1.status", geo("sim-x1", "fleet.drone", 30.2672, -97.7419, 152.4, -1, -1, "SIM-X1", 60, T0 + 36000,
                                         {"battery_pct": 79, "mode": "AUTO", "link_rssi": -66, "node_id": "build-mac"}))
add("tak.cot.mesh.rx", mesh_sim(T0 + 37000, 30.2672, -97.7417))
# feed-intel (USGS and EMSC formats) plus a record without coordinates (skipped).
add("intel.seismic.usgs.m4.us7000abcd", json.dumps({"id": "us7000abcd", "type": "earthquake", "source": "usgs", "lat": 35.1234,
    "lon": -117.5678, "timestamp": iso(T0 - 600000), "expires": "", "properties": {"magnitude": 4.4, "magnitude_type": "mw",
    "depth_km": 8.1, "place": "12 km SW of Ridgecrest, CA", "status": "reviewed", "tsunami": False}}))
add("intel.seismic.emsc.m3.20261029_0001", json.dumps({"id": "emsc-20261029_0001", "type": "earthquake", "source": "emsc",
    "lat": 38.2, "lon": 21.7, "timestamp": iso(T0 - 300000), "expires": "", "properties": {"magnitude": 3.1, "place": "GREECE",
    "region": "GREECE"}}))
add("intel.weather.noaa.alert.tx.flood", json.dumps({"id": "noaa-1", "type": "flood", "source": "noaa", "timestamp": iso(T0)}))
# Malformed CoT: dead-lettered with the reason.
add("tak.cot.inject", "<event uid='broken'")

for line in out:
    print(json.dumps(line, sort_keys=True))
