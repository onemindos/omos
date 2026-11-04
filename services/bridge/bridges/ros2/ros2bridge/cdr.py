"""Minimal OMG CDR (XCDR1, plain CDR) codec for the ROS 2 messages this bridge uses.

ROS 2 serializes messages as plain CDR with a 4-byte encapsulation header
(``00 01 00 00`` little-endian, ``00 00 00 00`` big-endian). Primitive alignment is
relative to the first byte after the header, and every primitive is aligned to its own
size (8 for float64). Strings are a uint32 length that counts a terminating NUL,
followed by the bytes and the NUL. Fixed arrays have no length prefix; sequences have a
uint32 element count.

zenoh-bridge-ros2dds routes these bytes between DDS and Zenoh without re-encoding them,
so a Zenoh subscriber sees the same buffer a DDS reader would (see README.md
"Wire format", marked unverified until the owner check in OWNER-STEPS).
"""

from __future__ import annotations

import struct

HEADER_LE = b"\x00\x01\x00\x00"
HEADER_BE = b"\x00\x00\x00\x00"


class CDRError(ValueError):
    """A buffer that is not valid plain CDR for the expected type."""


class Reader:
    def __init__(self, buf: bytes):
        if len(buf) < 4:
            raise CDRError("buffer shorter than the encapsulation header")
        kind = bytes(buf[0:2])
        if kind == b"\x00\x01":
            self._e = "<"
        elif kind == b"\x00\x00":
            self._e = ">"
        else:
            # 0x0006..0x000b are XCDR2 encodings, which ROS 2 does not use by default.
            raise CDRError(f"unsupported encapsulation {kind.hex()}")
        self._b = memoryview(bytes(buf))[4:]
        self._p = 0

    def _align(self, n: int) -> None:
        self._p += (-self._p) % n

    def _unpack(self, fmt: str, size: int):
        self._align(size)
        if self._p + size > len(self._b):
            raise CDRError("truncated buffer")
        v = struct.unpack_from(self._e + fmt, self._b, self._p)[0]
        self._p += size
        return v

    def u8(self) -> int:
        return self._unpack("B", 1)

    def i8(self) -> int:
        return self._unpack("b", 1)

    def boolean(self) -> bool:
        v = self.u8()
        if v > 1:
            raise CDRError(f"bool byte {v}")
        return v == 1

    def i32(self) -> int:
        return self._unpack("i", 4)

    def u32(self) -> int:
        return self._unpack("I", 4)

    def f32(self) -> float:
        return self._unpack("f", 4)

    def f64(self) -> float:
        return self._unpack("d", 8)

    def string(self) -> str:
        n = self.u32()
        if n == 0:
            return ""  # some writers send 0 for an empty string
        if self._p + n > len(self._b):
            raise CDRError("truncated string")
        raw = bytes(self._b[self._p : self._p + n])
        self._p += n
        if raw[-1:] != b"\x00":
            raise CDRError("string without NUL terminator")
        try:
            return raw[:-1].decode("utf-8")
        except UnicodeDecodeError as e:
            raise CDRError(f"string is not UTF-8: {e}") from e

    def octets(self, n: int) -> bytes:
        if self._p + n > len(self._b):
            raise CDRError("truncated octet array")
        raw = bytes(self._b[self._p : self._p + n])
        self._p += n
        return raw

    def f64_array(self, n: int) -> list[float]:
        return [self.f64() for _ in range(n)]

    def f32_seq(self) -> list[float]:
        n = self.u32()
        if n > 1 << 16:
            raise CDRError(f"sequence length {n} too large")
        return [self.f32() for _ in range(n)]

    def seq_len(self, limit: int = 1 << 16) -> int:
        n = self.u32()
        if n > limit:
            raise CDRError(f"sequence length {n} too large")
        return n


class Writer:
    """Little-endian plain CDR writer (the encoding ROS 2 publishers use by default)."""

    def __init__(self):
        self._b = bytearray()

    def _align(self, n: int) -> None:
        self._b.extend(b"\x00" * ((-len(self._b)) % n))

    def _pack(self, fmt: str, size: int, v) -> None:
        self._align(size)
        self._b.extend(struct.pack("<" + fmt, v))

    def u8(self, v: int) -> None:
        self._pack("B", 1, v)

    def i8(self, v: int) -> None:
        self._pack("b", 1, v)

    def boolean(self, v: bool) -> None:
        self.u8(1 if v else 0)

    def i32(self, v: int) -> None:
        self._pack("i", 4, v)

    def u32(self, v: int) -> None:
        self._pack("I", 4, v)

    def f32(self, v: float) -> None:
        self._pack("f", 4, v)

    def f64(self, v: float) -> None:
        self._pack("d", 8, v)

    def string(self, s: str) -> None:
        raw = s.encode("utf-8") + b"\x00"
        self.u32(len(raw))
        self._b.extend(raw)

    def octets(self, raw: bytes) -> None:
        self._b.extend(raw)

    def f64_array(self, vs) -> None:
        for v in vs:
            self.f64(v)

    def f32_seq(self, vs) -> None:
        self.u32(len(vs))
        for v in vs:
            self.f32(v)

    def bytes(self) -> bytes:
        return HEADER_LE + bytes(self._b)
