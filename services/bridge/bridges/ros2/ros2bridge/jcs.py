"""RFC 8785 JSON Canonicalization Scheme (https://www.rfc-editor.org/rfc/rfc8785).

Used for the command signature (the JCS form of the envelope without ``sig``) and for
``params_hash`` (asset-command.md section 5). Must produce the same bytes as the Go
implementations in the policy service and the MAVLink bridge; tests/test_jcs.py checks
this against vectors produced by the policy service code.
"""

from __future__ import annotations

import json
import math
from decimal import Decimal


def _number(v) -> str:
    if isinstance(v, bool):
        raise TypeError("bool is not a number")
    if isinstance(v, int):
        if abs(v) < 2**53:
            return str(v)
        v = float(v)
    if not math.isfinite(v):
        raise ValueError("JCS: NaN and Infinity are not allowed")
    if v == 0:
        return "0"  # also -0
    # ECMAScript Number::toString on the shortest round-trip digits (RFC 8785 3.2.2.3).
    sign = "-" if v < 0 else ""
    t = Decimal(repr(abs(v))).as_tuple()
    digits = "".join(map(str, t.digits)).rstrip("0") or "0"
    # value = 0.digits * 10**n
    n = t.exponent + len("".join(map(str, t.digits)))
    k = len(digits)
    if k <= n <= 21:
        s = digits + "0" * (n - k)
    elif 0 < n <= 21:
        s = digits[:n] + "." + digits[n:]
    elif -6 < n <= 0:
        s = "0." + "0" * (-n) + digits
    else:
        e = n - 1
        es = ("+" if e >= 0 else "-") + str(abs(e))
        s = (digits if k == 1 else digits[0] + "." + digits[1:]) + "e" + es
    return sign + s


def _string(s: str) -> str:
    # json.dumps with ensure_ascii=False escapes exactly ", \ and U+0000..U+001F, using
    # \b \t \n \f \r and lowercase \u00xx, which is what RFC 8785 3.2.2.2 requires.
    return json.dumps(s, ensure_ascii=False)


def _ser(v, out: list) -> None:
    if v is None:
        out.append("null")
    elif v is True:
        out.append("true")
    elif v is False:
        out.append("false")
    elif isinstance(v, (int, float)):
        out.append(_number(v))
    elif isinstance(v, str):
        out.append(_string(v))
    elif isinstance(v, (list, tuple)):
        out.append("[")
        for i, x in enumerate(v):
            if i:
                out.append(",")
            _ser(x, out)
        out.append("]")
    elif isinstance(v, dict):
        out.append("{")
        # Sort by UTF-16 code units (RFC 8785 3.2.3).
        for i, k in enumerate(sorted(v, key=lambda s: s.encode("utf-16-be"))):
            if not isinstance(k, str):
                raise TypeError("JCS: object keys must be strings")
            if i:
                out.append(",")
            out.append(_string(k))
            out.append(":")
            _ser(v[k], out)
        out.append("}")
    else:
        raise TypeError(f"JCS: unsupported type {type(v).__name__}")


def canonicalize(v) -> bytes:
    out: list = []
    _ser(v, out)
    return "".join(out).encode("utf-8")


def canonicalize_json(raw: bytes | str) -> bytes:
    return canonicalize(json.loads(raw))
