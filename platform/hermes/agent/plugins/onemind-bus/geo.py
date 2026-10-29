"""Small WGS-84 helpers for the onemind-bus tools (spherical Earth; good to about 0.5 %)."""
from __future__ import annotations

import math
from typing import Iterable, List, Sequence, Tuple

EARTH_R_M = 6371008.8


def haversine_m(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
    p1, p2 = math.radians(lat1), math.radians(lat2)
    dp, dl = p2 - p1, math.radians(lon2 - lon1)
    a = math.sin(dp / 2) ** 2 + math.cos(p1) * math.cos(p2) * math.sin(dl / 2) ** 2
    return 2 * EARTH_R_M * math.asin(min(1.0, math.sqrt(a)))


def path_length_m(points: Sequence[Tuple[float, float]]) -> float:
    """Length of a (lat, lon) polyline."""
    return sum(haversine_m(a[0], a[1], b[0], b[1]) for a, b in zip(points, points[1:]))


def polygon_area_m2(points: Sequence[Tuple[float, float]]) -> float:
    """Area of a (lat, lon) ring on the sphere (closed or not)."""
    pts = list(points)
    if len(pts) < 3:
        return 0.0
    if pts[0] != pts[-1]:
        pts.append(pts[0])
    total = 0.0
    for (la1, lo1), (la2, lo2) in zip(pts, pts[1:]):
        total += math.radians(lo2 - lo1) * (2 + math.sin(math.radians(la1)) + math.sin(math.radians(la2)))
    return abs(total * EARTH_R_M ** 2 / 2)


def bbox_around(lat: float, lon: float, radius_m: float) -> List[float]:
    """[west, south, east, north] that contains a circle of radius_m."""
    dlat = math.degrees(radius_m / EARTH_R_M)
    dlon = math.degrees(radius_m / (EARTH_R_M * max(0.01, math.cos(math.radians(lat)))))
    return [max(-180.0, lon - dlon), max(-90.0, lat - dlat), min(180.0, lon + dlon), min(90.0, lat + dlat)]


def bbox_of(points: Iterable[Tuple[float, float]], pad_ratio: float = 0.15, min_span_deg: float = 0.002) -> List[float]:
    """[west, south, east, north] around (lat, lon) points, padded."""
    pts = list(points)
    if not pts:
        raise ValueError("no points")
    lats = [p[0] for p in pts]
    lons = [p[1] for p in pts]
    w, e, s, n = min(lons), max(lons), min(lats), max(lats)
    dx = max(e - w, min_span_deg) * pad_ratio
    dy = max(n - s, min_span_deg) * pad_ratio
    if e - w < min_span_deg:
        c = (e + w) / 2
        w, e = c - min_span_deg / 2, c + min_span_deg / 2
    if n - s < min_span_deg:
        c = (n + s) / 2
        s, n = c - min_span_deg / 2, c + min_span_deg / 2
    return [max(-180.0, w - dx), max(-90.0, s - dy), min(180.0, e + dx), min(90.0, n + dy)]


def valid_bbox(b) -> bool:
    try:
        w, s, e, n = (float(x) for x in b)
    except (TypeError, ValueError):
        return False
    return -180 <= w <= 180 and -180 <= e <= 180 and -90 <= s < n <= 90 and w < e


def valid_point(lat, lon) -> bool:
    try:
        return -90 <= float(lat) <= 90 and -180 <= float(lon) <= 180
    except (TypeError, ValueError):
        return False


def zoom_for_range(range_m: float) -> int:
    """Approximate XYZ zoom for a camera range, for clients that only read §9.7 ``zoom``."""
    return int(max(0, min(22, round(math.log2(40075016.0 / max(50.0, range_m))))))


def ring_from_geojson(geom) -> List[Tuple[float, float]]:
    """(lat, lon) outer ring of a GeoJSON Polygon, or the points of a LineString/Point."""
    if not isinstance(geom, dict):
        return []
    t, c = geom.get("type"), geom.get("coordinates")
    try:
        if t == "Polygon":
            return [(float(p[1]), float(p[0])) for p in c[0]]
        if t == "MultiPolygon":
            return [(float(p[1]), float(p[0])) for p in c[0][0]]
        if t == "LineString":
            return [(float(p[1]), float(p[0])) for p in c]
        if t == "Point":
            return [(float(c[1]), float(c[0]))]
    except (TypeError, IndexError, ValueError):
        return []
    return []
