"""The bridge's own keep-in / keep-out check (asset-command.md section 7 rule 3).

Port of services/bridge/bridges/mavlink/internal/fence (omos PR #550) with the same file
format, so one geofence file serves both bridges. It runs before any movement command,
independent of the policy service, and fails closed: with no keep-in zone configured,
every position is outside. Polygons use planar ray casting on lat/lon, which is adequate
for site-sized zones and does not handle zones that cross the antimeridian or a pole.
"""

from __future__ import annotations

import json
import math
from dataclasses import dataclass, field


@dataclass
class Zone:
    id: str
    polygon: list = field(default_factory=list)  # [(lat, lon), ...]
    circle: tuple | None = None  # (lat, lon, radius_m)

    def contains(self, lat: float, lon: float) -> bool:
        if self.circle is not None:
            clat, clon, r = self.circle
            return distance(clat, clon, lat, lon) <= r
        inside = False
        n = len(self.polygon)
        j = n - 1
        for i in range(n):
            alat, alon = self.polygon[i]
            blat, blon = self.polygon[j]
            if (alat > lat) != (blat > lat) and lon < (blon - alon) * (lat - alat) / (blat - alat) + alon:
                inside = not inside
            j = i
        return inside


@dataclass
class Violation:
    detail: str  # keep_in | keep_out (asset-command.md section 6.4 precondition details)
    reason: str
    zone_id: str = ""


def _valid(lat, lon) -> bool:
    return (isinstance(lat, (int, float)) and isinstance(lon, (int, float))
            and math.isfinite(lat) and math.isfinite(lon) and -90 <= lat <= 90 and -180 <= lon <= 180)


def distance(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
    """Haversine distance in metres."""
    r = 6371008.8
    la1, la2 = math.radians(lat1), math.radians(lat2)
    dla, dlo = la2 - la1, math.radians(lon2 - lon1)
    h = math.sin(dla / 2) ** 2 + math.cos(la1) * math.cos(la2) * math.sin(dlo / 2) ** 2
    return 2 * r * math.asin(min(1.0, math.sqrt(h)))


def _zone(d: dict) -> Zone:
    zid = d.get("id") or ""
    if not zid:
        raise ValueError("geofence: zone without id")
    poly, circ = d.get("polygon"), d.get("circle")
    if poly and circ:
        raise ValueError(f"geofence: zone {zid} has both polygon and circle")
    if circ:
        lat, lon, r = circ.get("lat"), circ.get("lon"), circ.get("radius_m")
        if not _valid(lat, lon) or not (isinstance(r, (int, float)) and r > 0):
            raise ValueError(f"geofence: zone {zid}: bad circle")
        return Zone(zid, circle=(float(lat), float(lon), float(r)))
    if poly and len(poly) >= 3:
        pts = []
        for p in poly:
            if not _valid(p.get("lat"), p.get("lon")):
                raise ValueError(f"geofence: zone {zid}: bad vertex {p}")
            pts.append((float(p["lat"]), float(p["lon"])))
        return Zone(zid, polygon=pts)
    raise ValueError(f"geofence: zone {zid} needs a circle or 3+ polygon vertices")


@dataclass
class Fence:
    keep_in: list = field(default_factory=list)
    keep_out: list = field(default_factory=list)

    @staticmethod
    def parse(raw: bytes | str) -> "Fence":
        try:
            d = json.loads(raw)
        except json.JSONDecodeError as e:
            raise ValueError(f"geofence: {e}") from e
        if not isinstance(d, dict):
            raise ValueError("geofence: document must be an object")
        # max_hae_m / min_hae_m are accepted for file compatibility with the MAVLink
        # bridge but ignored: a ground robot has no altitude command.
        return Fence([_zone(z) for z in d.get("keep_in") or []], [_zone(z) for z in d.get("keep_out") or []])

    @staticmethod
    def load(path: str) -> "Fence":
        with open(path, "rb") as f:
            return Fence.parse(f.read())

    def check(self, lat: float, lon: float) -> Violation | None:
        if not self.keep_in:
            return Violation("keep_in", "no keep-in zone configured (fail closed)")
        if not _valid(lat, lon):
            return Violation("keep_in", f"invalid position {lat},{lon}")
        for z in self.keep_out:
            if z.contains(lat, lon):
                return Violation("keep_out", f"position inside keep-out zone {z.id}", z.id)
        if not any(z.contains(lat, lon) for z in self.keep_in):
            return Violation("keep_in", "position outside every keep-in zone")
        return None


    def check_path(self, lat0: float, lon0: float, lat1: float, lon1: float,
                   step_m: float = 0.25, max_samples: int = 4000) -> Violation | None:
        """Cheap first pass on the way to a goal: sample the straight segment from the
        current position to the target and refuse if it enters a keep-out zone. A zone
        that already contains the start is skipped, so a robot that drifted into one can
        still drive out. Nav2 plans its own path, which may differ from this line; on real
        robots keep-out zones also belong in the Nav2 costmap keepout filter (owner
        proposal H3-P7)."""
        if not (_valid(lat0, lon0) and _valid(lat1, lon1)):
            return None  # endpoint validity is check()'s job
        zones = [z for z in self.keep_out if not z.contains(lat0, lon0)]
        if not zones:
            return None
        n = min(max_samples, max(1, int(math.ceil(distance(lat0, lon0, lat1, lon1) / step_m))))
        for i in range(1, n + 1):
            f = i / n
            lat, lon = lat0 + (lat1 - lat0) * f, lon0 + (lon1 - lon0) * f
            for z in zones:
                if z.contains(lat, lon):
                    return Violation("keep_out", f"straight path crosses keep-out zone {z.id}", z.id)
        return None


def check_path(f: Fence | None, lat0: float, lon0: float, lat1: float, lon1: float) -> Violation | None:
    if f is None:
        return None  # check() already failed closed on the target
    return f.check_path(lat0, lon0, lat1, lon1)


def check(f: Fence | None, lat: float, lon: float) -> Violation | None:
    if f is None:
        return Violation("keep_in", "no keep-in zone configured (fail closed)")
    return f.check(lat, lon)
