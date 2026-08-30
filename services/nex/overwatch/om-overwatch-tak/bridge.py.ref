#!/usr/bin/env python3
"""
cot2constellation — CoT → Constellation Overwatch Bridge
=========================================================
Subscribes to tak.cot.inject on omos-core NATS (:4222)
Parses CoT XML and publishes Constellation telemetry to Overwatch NATS (:4299)

Subjects:
  IN:  tak.cot.inject              (omos-core :4222)
  OUT: constellation.telemetry.{org_id}.{entity_id}  (Overwatch :4299)

CoT entity_type → Constellation entity_type mapping:
  a-f-*  → operator_station  (friendly)
  a-h-*  → operator_station  (hostile — mapped as sensor for now)
  a-.-A  → aircraft_multirotor
  a-.-G  → ground_vehicle_wheeled
  a-.-S  → surface_vessel_usv
  *      → sensor_platform (fallback)
"""

import asyncio
import json
import logging
import os
import re
import xml.etree.ElementTree as ET
from datetime import datetime, timezone

import nats

# ── Config ────────────────────────────────────────────────────────────────────
OMOS_NATS    = os.getenv("OMOS_NATS",   "nats://localhost:4222")
OW_NATS      = os.getenv("OW_NATS",     "nats://localhost:4299")
OW_NATS_TOKEN     = os.getenv("OW_NATS_TOKEN", "")      # raw token (legacy)
OW_NATS_SEED_FILE = os.getenv("OW_NATS_SEED_FILE", "")  # path to NKey seed file
COT_SUBJECT  = os.getenv("COT_SUBJECT", "tak.cot.inject")
OW_ORG_ID    = os.getenv("OW_ORG_ID",  "")
OW_API_KEY   = os.getenv("OW_API_KEY", "")
OW_API_BASE  = os.getenv("OW_API_BASE", "http://localhost:8090")

# Transport mode: "nats" (direct NKey) or "rest" (HTTP API)
_has_nats_creds = bool(OW_NATS_SEED_FILE or OW_NATS_TOKEN)
TRANSPORT    = os.getenv("OW_TRANSPORT", "nats" if _has_nats_creds else "rest")

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [cot2constellation] %(levelname)s %(message)s",
    datefmt="%H:%M:%S",
)
log = logging.getLogger(__name__)

# ── CoT type → Constellation entity_type ────────────────────────────────────
def cot_type_to_entity_type(cot_type: str) -> str:
    t = cot_type.lower()
    # Aircraft
    if re.search(r"a-[fhnu]-a", t): return "aircraft_multirotor"
    if re.search(r"a-[fhnu]-a-f", t): return "aircraft_fixed_wing"
    # Ground
    if re.search(r"a-[fhnu]-g", t): return "ground_vehicle_wheeled"
    if re.search(r"a-[fhnu]-g-u-c", t): return "ground_vehicle_tracked"
    # Sea
    if re.search(r"a-[fhnu]-s", t): return "surface_vessel_usv"
    # Operator / human
    if re.search(r"a-f-g-u-c", t): return "operator_station"
    if re.search(r"a-[fhnu]-g-e-s", t): return "sensor_platform"
    # Default
    return "sensor_platform"

def cot_type_to_status(cot_type: str) -> str:
    t = cot_type.lower()
    if "-f-" in t: return "active"       # friendly
    if "-h-" in t: return "active"       # hostile
    if "-n-" in t: return "inactive"     # neutral
    if "-u-" in t: return "unknown"      # unknown
    return "active"

# ── CoT XML parser ────────────────────────────────────────────────────────────
def parse_cot(xml_bytes: bytes) -> dict | None:
    try:
        root = ET.fromstring(xml_bytes)
    except ET.ParseError as e:
        log.warning("XML parse error: %s", e)
        return None

    if root.tag != "event":
        return None

    uid      = root.attrib.get("uid", "")
    cot_type = root.attrib.get("type", "")
    how      = root.attrib.get("how", "")
    stale    = root.attrib.get("stale", "")

    point = root.find("point")
    if point is None:
        return None

    lat = float(point.attrib.get("lat", 0))
    lon = float(point.attrib.get("lon", 0))
    hae = float(point.attrib.get("hae", 0))    # height above ellipsoid (m)
    ce  = float(point.attrib.get("ce", 9999))  # circular error (m)

    # Pull callsign from detail/contact or detail/uid
    detail   = root.find("detail")
    callsign = uid
    course   = 0.0
    speed    = 0.0
    battery  = None
    platform = ""
    version  = ""

    if detail is not None:
        contact = detail.find("contact")
        if contact is not None:
            callsign = contact.attrib.get("callsign", uid)

        track = detail.find("track")
        if track is not None:
            course = float(track.attrib.get("course", 0))
            speed  = float(track.attrib.get("speed", 0))

        status_el = detail.find("status")
        if status_el is not None:
            batt_str = status_el.attrib.get("battery", "")
            if batt_str:
                try:
                    battery = float(batt_str)
                except ValueError:
                    pass

        takv = detail.find("takv")
        if takv is not None:
            platform = takv.attrib.get("platform", "")
            version  = takv.attrib.get("version", "")

    entity_type = cot_type_to_entity_type(cot_type)
    status      = cot_type_to_status(cot_type)

    # Sanitize uid to be a valid NATS token (no dots or spaces)
    entity_id = re.sub(r"[^a-zA-Z0-9_-]", "_", uid)[:64]

    metadata = {
        "callsign": callsign,
        "cot_type": cot_type,
        "how": how,
        "ce": ce,
        "source": "tak.cot.inject",
    }
    if platform: metadata["platform"] = platform
    if version:  metadata["version"]  = version
    if battery is not None: metadata["battery"] = battery

    return {
        "entity_id":   entity_id,
        "entity_type": entity_type,
        "callsign":    callsign,
        "status":      status,
        "lat":         lat,
        "lon":         lon,
        "alt":         hae,
        "heading":     course,
        "velocity":    speed,
        "metadata":    metadata,
        "stale":       stale,
    }

# ── Auto-register entity via REST API ─────────────────────────────────────────
import urllib.request
import urllib.error

_registered: set[str] = set()

def ensure_entity_registered(entity_id: str, parsed: dict):
    """Fire-and-forget REST POST to register entity if not seen before."""
    if not OW_ORG_ID or not OW_API_KEY:
        return  # skip registration if not configured — telemetry still flows
    if entity_id in _registered:
        return

    payload = json.dumps({
        "name":        parsed["callsign"],
        "entity_type": parsed["entity_type"],
        "status":      parsed["status"],
        "position": {
            "latitude":  parsed["lat"],
            "longitude": parsed["lon"],
            "altitude":  parsed["alt"],
        },
        "metadata": parsed["metadata"],
    }).encode()

    url = f"{OW_API_BASE}/api/v1/organizations/{OW_ORG_ID}/entities"
    req = urllib.request.Request(
        url,
        data=payload,
        headers={"X-API-Key": OW_API_KEY, "Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=3) as resp:
            body = json.loads(resp.read())
            if body.get("entity_id"):
                _registered.add(entity_id)
                log.info("✅ Registered entity: %s (%s)", parsed["callsign"], entity_id)
    except urllib.error.HTTPError as e:
        if e.code == 409:
            _registered.add(entity_id)  # already exists
        else:
            log.debug("Entity register failed %s: %s", entity_id, e)
    except Exception as e:
        log.debug("Entity register error: %s", e)

# ── REST telemetry push ────────────────────────────────────────────────────────
def push_telemetry_rest(org_id: str, entity_id: str, parsed: dict):
    """Push position update via REST API PUT entity."""
    if not OW_API_KEY or not org_id:
        return
    payload = json.dumps({
        "status":   parsed["status"],
        "position": {
            "latitude":  parsed["lat"],
            "longitude": parsed["lon"],
            "altitude":  parsed["alt"],
        },
        "metadata": parsed["metadata"],
    }).encode()
    url = f"{OW_API_BASE}/api/v1/entities?org_id={org_id}&entity_id={entity_id}"
    req = urllib.request.Request(
        url,
        data=payload,
        headers={"X-API-Key": OW_API_KEY, "Content-Type": "application/json"},
        method="PUT",
    )
    try:
        with urllib.request.urlopen(req, timeout=3) as resp:
            pass  # success
    except Exception as e:
        log.debug("REST telemetry push error: %s", e)


# ── Main bridge ───────────────────────────────────────────────────────────────
async def run():
    log.info("🌉 cot2constellation starting")
    log.info("   IN  ← %s [%s]", OMOS_NATS, COT_SUBJECT)
    log.info("   Transport: %s", TRANSPORT)

    if TRANSPORT == "nats":
        log.info("   OUT → %s [constellation.telemetry.{org}.{entity}]", OW_NATS)
        if not OW_NATS_SEED_FILE and not OW_NATS_TOKEN:
            log.error("❌ OW_TRANSPORT=nats but no NKey seed file or token set.")
            return
        connect_opts = dict(name="cot2constellation-out")
        if OW_NATS_SEED_FILE:
            connect_opts["nkeys_seed"] = OW_NATS_SEED_FILE
        else:
            connect_opts["token"] = OW_NATS_TOKEN
        nc_out = await nats.connect(OW_NATS, **connect_opts)
        log.info("✅ Connected to Overwatch NATS via NKey")
    else:
        nc_out = None
        log.info("   OUT → %s (REST API)", OW_API_BASE)
        if not OW_API_KEY:
            log.warning("⚠️  OW_API_KEY not set — positions will be logged but NOT pushed to Overwatch")
        if not OW_ORG_ID:
            log.warning("⚠️  OW_ORG_ID not set — entities will NOT be auto-registered")

    nc_in = await nats.connect(OMOS_NATS, name="cot2constellation-in")
    log.info("✅ Connected to omos-core NATS")

    msg_count  = 0
    skip_count = 0

    async def handler(msg):
        nonlocal msg_count, skip_count
        try:
            parsed = parse_cot(msg.data)
            if parsed is None:
                skip_count += 1
                return

            entity_id = parsed["entity_id"]
            org_id    = OW_ORG_ID or "onemind"

            # Always attempt entity registration (idempotent)
            ensure_entity_registered(entity_id, parsed)

            if TRANSPORT == "nats" and nc_out:
                telemetry = {
                    "message_type": "GPS_POSITION",
                    "entity_id":    entity_id,
                    "data": {
                        "lat":      int(parsed["lat"]  * 1e7),
                        "lon":      int(parsed["lon"]  * 1e7),
                        "alt":      int(parsed["alt"]  * 1000),
                        "heading":  parsed["heading"],
                        "speed":    parsed["velocity"],
                        "callsign": parsed["callsign"],
                        "cot_type": parsed["metadata"]["cot_type"],
                        "status":   parsed["status"],
                        "metadata": parsed["metadata"],
                    },
                    "timestamp": datetime.now(timezone.utc).isoformat(),
                }
                subject = f"constellation.telemetry.{org_id}.{entity_id}"
                await nc_out.publish(subject, json.dumps(telemetry).encode())
            else:
                # REST path — run in executor to avoid blocking event loop
                loop = asyncio.get_event_loop()
                await loop.run_in_executor(None, push_telemetry_rest, org_id, entity_id, parsed)

            msg_count += 1
            if msg_count % 10 == 0:
                log.info("📡 Bridged %d CoT msgs (%d skipped) | last: %s @ %.4f,%.4f",
                         msg_count, skip_count, parsed["callsign"], parsed["lat"], parsed["lon"])
            else:
                log.debug("📡 %s → %s @ %.4f,%.4f alt=%.0fm",
                          parsed["callsign"], entity_id, parsed["lat"], parsed["lon"], parsed["alt"])

        except Exception as e:
            log.error("Handler error: %s", e)

    sub = await nc_in.subscribe(COT_SUBJECT, cb=handler)
    log.info("📻 Subscribed to %s on omos-core NATS", COT_SUBJECT)

    try:
        while True:
            await asyncio.sleep(30)
            log.info("💓 Bridge alive | transport=%s bridged=%d skipped=%d", TRANSPORT, msg_count, skip_count)
    except asyncio.CancelledError:
        pass
    finally:
        await sub.unsubscribe()
        await nc_in.close()
        if nc_out:
            await nc_out.close()
        log.info("Bridge stopped. Total bridged: %d", msg_count)

if __name__ == "__main__":
    try:
        asyncio.run(run())
    except KeyboardInterrupt:
        log.info("Interrupted.")
