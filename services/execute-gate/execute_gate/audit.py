"""Audit: one structured stdout line per decision (Loki) and, best effort, an event on the bus.

Subjects follow the locked grammar evt.{domain}.{entity}.{event}:
evt.toolgw.approval.{requested,granted,denied,consumed}. Argument values are never logged or published.
"""

from __future__ import annotations

import json
import sys
import time

BUS_EVENTS = {"requested", "granted", "denied", "consumed"}


class Audit:
    def __init__(self, publish=None, stream=None, clock=time.time):
        self._publish = publish
        self._stream = stream or sys.stdout
        self._clock = clock

    def emit(self, event: str, **fields) -> None:
        rec = {"ts": round(self._clock(), 3), "svc": "execute-gate", "event": event, **fields}
        try:
            self._stream.write(json.dumps(rec, sort_keys=True, ensure_ascii=False) + "\n")
            self._stream.flush()
        except Exception:  # logging must never change a decision
            pass
        if self._publish and event in BUS_EVENTS:
            try:
                self._publish(f"evt.toolgw.approval.{event}", json.dumps(rec, sort_keys=True).encode())
            except Exception as exc:
                try:
                    self._stream.write(
                        json.dumps({"svc": "execute-gate", "event": "audit_publish_failed", "error": str(exc)}) + "\n"
                    )
                except Exception:
                    pass
