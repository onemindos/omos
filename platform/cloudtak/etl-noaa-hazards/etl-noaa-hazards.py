#!/usr/bin/env python3
"""
etl-noaa-hazards: Fetches NOAA weather alerts from api.weather.gov,
converts to GeoJSON, and POSTs to CloudTAK API layer endpoint.
"""
import os
import sys
import json
import logging
from datetime import datetime, timezone
import urllib.request
import urllib.error

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
log = logging.getLogger("etl-noaa-hazards")

API_URL = os.environ.get("CLOUDTAK_API_URL", "http://cloudtak-api:5000")
LAYER_ID = int(os.environ.get("CLOUDTAK_LAYER_ID", "1"))
RESOURCE_TOKEN = os.environ.get("CLOUDTAK_RESOURCE_TOKEN", "")
NOAA_ALERTS_URL = "https://api.weather.gov/alerts/active?area=CT"
USER_AGENT = "OneMind-ETL/1.0 (omos@onemindos.com)"

COT_TYPE_MAP = {
    "Tornado Warning": "a-f-G-U-T-W",
    "Severe Thunderstorm Warning": "a-f-G-U-T-W",
    "Flash Flood Warning": "a-f-G-U-F-W",
    "Flood Warning": "a-f-G-U-F-W",
    "Winter Storm Warning": "a-f-G-U-W-W",
    "Blizzard Warning": "a-f-G-U-W-W",
    "High Wind Warning": "a-f-G-U-W-W",
    "Heat Advisory": "a-f-G-U-H-A",
    "Excessive Heat Warning": "a-f-G-U-H-W",
    "Coastal Flood Statement": "a-f-G-U-C-F",
    "Coastal Flood Advisory": "a-f-G-U-C-F",
    "Coastal Flood Warning": "a-f-G-U-C-F",
}


def fetch_alerts():
    req = urllib.request.Request(NOAA_ALERTS_URL, headers={"User-Agent": USER_AGENT})
    with urllib.request.urlopen(req, timeout=30) as resp:
        data = json.loads(resp.read().decode())
    features = data.get("features", [])
    log.info(f"Fetched {len(features)} active alerts from NOAA")
    return features


def convert_to_cot_geojson(noaa_features):
    cot_features = []
    for feature in noaa_features:
        props = feature.get("properties", {})
        geometry = feature.get("geometry")
        event = props.get("event", "Unknown Alert")
        headline = props.get("headline", "")
        severity = props.get("severity", "Unknown")
        urgency = props.get("urgency", "Unknown")
        certainty = props.get("certainty", "Unknown")
        effective = props.get("effective") or datetime.now(timezone.utc).isoformat()
        expires = props.get("expires") or effective
        area_desc = props.get("areaDesc", "")
        nws_office = props.get("senderName", "")
        alert_id = props.get("id", "unknown")
        if geometry is None:
            log.warning(f"No geometry for alert {alert_id}, skipping")
            continue
        cot_type = COT_TYPE_MAP.get(event, "a-f-G-U-O")
        short_id = alert_id.replace("https://api.weather.gov/alerts/", "").replace("/", "-")
        cot_feature = {
            "id": f"noaa-{short_id}",
            "type": "Feature",
            "properties": {
                "callsign": event[:24],
                "type": cot_type,
                "how": "h-g-i-g-o",
                "time": effective,
                "start": effective,
                "stale": expires,
                "remarks": (
                    f"{headline}\n"
                    f"Severity: {severity} | Urgency: {urgency} | Certainty: {certainty}\n"
                    f"Area: {area_desc}\n"
                    f"Office: {nws_office}"
                ),
            },
            "geometry": geometry,
        }
        cot_features.append(cot_feature)
    return cot_features


def post_to_cloudtak(features):
    if not features:
        log.info("No features to post")
        return True
    url = f"{API_URL}/api/layer/{LAYER_ID}/cot"
    payload = {"type": "FeatureCollection", "features": features}
    body = json.dumps(payload).encode()
    req = urllib.request.Request(
        url, data=body,
        headers={
            "Authorization": f"Bearer {RESOURCE_TOKEN}",
            "Content-Type": "application/json",
        },
    )
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            result = json.loads(resp.read().decode())
            log.info(f"Posted {len(features)} features -> {result.get('message', 'OK')}")
            return True
    except urllib.error.HTTPError as e:
        body = e.read().decode()[:200]
        log.error(f"CloudTAK POST failed: HTTP {e.code} - {body}")
        return False


def main():
    log.info(f"ETL NOAA Hazards starting - Layer {LAYER_ID} at {API_URL}")
    try:
        noaa_features = fetch_alerts()
        cot_features = convert_to_cot_geojson(noaa_features)
        log.info(f"Converted {len(cot_features)} alerts to CoT GeoJSON")
        success = post_to_cloudtak(cot_features)
        if success:
            log.info("ETL run complete")
        else:
            sys.exit(1)
    except Exception as e:
        log.error(f"ETL run failed: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
