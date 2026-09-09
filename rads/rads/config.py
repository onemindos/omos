"""Central config — reads from env / .env file."""
import os
from dotenv import load_dotenv
load_dotenv()

HOMESTEAD_ID     = os.environ.get("HOMESTEAD_ID", "greenwich")
NATS_URL         = os.environ.get("NATS_URL", "nats://localhost:4222")
CLICKHOUSE_HOST  = os.environ.get("CLICKHOUSE_HOST", "localhost")
CLICKHOUSE_PORT  = int(os.environ.get("CLICKHOUSE_PORT", "9001"))
CLICKHOUSE_DB    = os.environ.get("CLICKHOUSE_DB", "onemind")
MAVLINK_ADDRESS  = os.environ.get("MAVLINK_ADDRESS", "udp://:14540")
MOISTURE_THRESHOLD = float(os.environ.get("MOISTURE_THRESHOLD", "30.0"))
FIELD_COUNT      = int(os.environ.get("FIELD_COUNT", "4"))
SENSOR_INTERVAL  = float(os.environ.get("SENSOR_INTERVAL", "3.0"))

# NATS subject helpers
def soil_subject(field: int) -> str:
    return f"rad.{HOMESTEAD_ID}.land.field.{field}.soil"

def irrigate_subject(field: int) -> str:
    return f"rad.{HOMESTEAD_ID}.land.field.{field}.irrigate"

def decision_subject() -> str:
    return f"rad.{HOMESTEAD_ID}.land.decisions"

def drone_telemetry_subject(drone_id: str = "drone1") -> str:
    return f"rad.{HOMESTEAD_ID}.robots.{drone_id}.telemetry"

def drone_cmd_subject(drone_id: str = "drone1") -> str:
    return f"rad.{HOMESTEAD_ID}.robots.{drone_id}.cmd"

def legacy_override_subject() -> str:
    return f"rad.{HOMESTEAD_ID}.legacy.override"
