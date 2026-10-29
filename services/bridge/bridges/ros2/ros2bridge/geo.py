"""Map frame <-> WGS-84 for one robot site.

ROS 2 robots report poses in a local Cartesian frame (``map`` or ``odom``; REP-103:
x forward, y left, z up, yaw counter-clockwise from +x). OneMind telemetry needs
latitude, longitude and ellipsoid height (doctrine section 9.2 geo-envelope). The bridge
is configured with the WGS-84 position of the frame origin and the compass bearing of
the frame's +x axis, and converts with a local tangent plane. The error of this
approximation is well under a centimetre over a few kilometres, which is the scale of a
robot site; it is not meant for a frame that spans tens of kilometres.
"""

from __future__ import annotations

import math
from dataclasses import dataclass

_A = 6378137.0  # WGS-84 semi-major axis (m)
_F = 1 / 298.257223563
_E2 = _F * (2 - _F)


@dataclass(frozen=True)
class LocalFrame:
    origin_lat: float
    origin_lon: float
    origin_hae: float = 0.0
    x_bearing_deg: float = 90.0  # compass bearing of +x; 90 = +x points east (ENU)

    def __post_init__(self):
        for name, v, lo, hi in (("origin_lat", self.origin_lat, -89.0, 89.0),
                                ("origin_lon", self.origin_lon, -180.0, 180.0),
                                ("x_bearing_deg", self.x_bearing_deg, 0.0, 360.0)):
            if not (isinstance(v, (int, float)) and math.isfinite(v) and lo <= v <= hi):
                raise ValueError(f"{name}={v!r} out of range [{lo}, {hi}]")

    def _radii(self) -> tuple[float, float]:
        lat = math.radians(self.origin_lat)
        s = math.sin(lat)
        w = math.sqrt(1 - _E2 * s * s)
        m = _A * (1 - _E2) / (w ** 3)  # meridional radius of curvature
        n = _A / w  # prime vertical radius of curvature
        return m, n

    def _rot(self) -> float:
        # Angle from east to +x, counter-clockwise (math convention).
        return math.radians(90.0 - self.x_bearing_deg)

    def to_wgs84(self, x: float, y: float, z: float = 0.0) -> tuple[float, float, float]:
        th = self._rot()
        east = x * math.cos(th) - y * math.sin(th)
        north = x * math.sin(th) + y * math.cos(th)
        m, n = self._radii()
        lat = self.origin_lat + math.degrees(north / m)
        lon = self.origin_lon + math.degrees(east / (n * math.cos(math.radians(self.origin_lat))))
        return lat, lon, self.origin_hae + z

    def from_wgs84(self, lat: float, lon: float) -> tuple[float, float]:
        m, n = self._radii()
        north = math.radians(lat - self.origin_lat) * m
        east = math.radians(lon - self.origin_lon) * n * math.cos(math.radians(self.origin_lat))
        th = self._rot()
        x = east * math.cos(th) + north * math.sin(th)
        y = -east * math.sin(th) + north * math.cos(th)
        return x, y

    def yaw_to_bearing(self, yaw_rad: float) -> float:
        """Frame yaw (rad, CCW from +x) to compass bearing (deg clockwise from north)."""
        return (self.x_bearing_deg - math.degrees(yaw_rad)) % 360.0

    def bearing_to_yaw(self, bearing_deg: float) -> float:
        return math.radians(self.x_bearing_deg - bearing_deg)
