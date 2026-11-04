"""Canonical argument JSON and the per-call approval key."""

from __future__ import annotations

import hashlib
import json
import math
import re
from typing import Any

_SECRET_KEY = re.compile(r"(pass|secret|token|key|auth|cookie|credential|bearer)", re.I)


def _normalise(value: Any) -> Any:
    if isinstance(value, bool) or value is None or isinstance(value, (str, int)):
        return value
    if isinstance(value, float):
        if not math.isfinite(value):
            raise ValueError("non-finite number in arguments")
        # 1.0 and 1 are the same argument.
        return int(value) if value.is_integer() else value
    if isinstance(value, (list, tuple)):
        return [_normalise(v) for v in value]
    if isinstance(value, dict):
        out = {}
        for k, v in value.items():
            if not isinstance(k, str):
                raise ValueError("non-string object key in arguments")
            out[k] = _normalise(v)
        return out
    raise ValueError(f"unsupported type in arguments: {type(value).__name__}")


def parse_arguments(raw: Any) -> Any:
    """Accept arguments as a JSON value or as JSON text."""
    if isinstance(raw, (str, bytes)):
        raw = json.loads(raw)
    return raw


def canonical_json(arguments: Any) -> str:
    """Sorted keys, no whitespace, UTF-8 characters kept, integral floats as ints."""
    return json.dumps(
        _normalise(arguments), sort_keys=True, separators=(",", ":"), ensure_ascii=False, allow_nan=False
    )


def call_hash(resource: str, arguments: Any) -> str:
    """Approval key: sha256 over the exact tool resource and canonical arguments."""
    body = f"{resource}\n{canonical_json(arguments)}".encode("utf-8")
    return hashlib.sha256(body).hexdigest()


def summarise(arguments: Any, max_value: int = 120, max_total: int = 700) -> str:
    """Human-readable summary for the approvals room. Secret-looking keys are masked."""

    def red(v: Any, key: str = "") -> Any:
        if key and _SECRET_KEY.search(key):
            return "[redacted]"
        if isinstance(v, dict):
            return {k: red(x, k) for k, x in v.items()}
        if isinstance(v, list):
            return [red(x) for x in v]
        if isinstance(v, str) and len(v) > max_value:
            return v[:max_value] + "..."
        return v

    try:
        text = json.dumps(red(_normalise(arguments)), sort_keys=True, ensure_ascii=False)
    except (ValueError, TypeError):
        return "[unrepresentable arguments]"
    return text if len(text) <= max_total else text[:max_total] + "..."
