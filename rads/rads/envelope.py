"""OMOS envelope schema — wraps all RADS payloads for fabric compatibility."""
import time, uuid, json
from typing import Any

def wrap(event_type: str, payload: Any, source: str = "rads") -> bytes:
    envelope = {
        "id":         str(uuid.uuid4()),
        "ts":         time.time(),
        "source":     source,
        "event_type": event_type,
        "payload":    payload,
    }
    return json.dumps(envelope).encode()

def unwrap(data: bytes) -> dict:
    return json.loads(data)
