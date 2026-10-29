import json
import math
import os

import pytest

from ros2bridge.fence import Fence, check, distance
from ros2bridge.geo import LocalFrame

EXAMPLE = os.path.join(os.path.dirname(__file__), "..", "geofence.example.json")


def test_example_fence():
    f = Fence.load(EXAMPLE)
    assert check(f, 30.2672, -97.7431) is None
    assert check(f, 30.2680, -97.7431).detail == "keep_in"
    v = check(f, 30.26745, -97.74285)
    assert v.detail == "keep_out" and v.zone_id == "pond"


def test_fail_closed():
    assert check(None, 0, 0).detail == "keep_in"
    assert Fence().check(0, 0).detail == "keep_in"
    assert Fence.load(EXAMPLE).check(float("nan"), 0).detail == "keep_in"


@pytest.mark.parametrize("doc", [
    "not json", "[]", '{"keep_in":[{"polygon":[]}]}', '{"keep_in":[{"id":"a","polygon":[{"lat":0,"lon":0}]}]}',
    '{"keep_in":[{"id":"a","circle":{"lat":0,"lon":0,"radius_m":0}}]}',
    '{"keep_in":[{"id":"a","circle":{"lat":0,"lon":0,"radius_m":5},"polygon":[{"lat":0,"lon":0},{"lat":1,"lon":0},{"lat":0,"lon":1}]}]}',
    '{"keep_in":[{"id":"a","polygon":[{"lat":95,"lon":0},{"lat":1,"lon":0},{"lat":0,"lon":1}]}]}',
])
def test_bad_fence_rejected(doc):
    with pytest.raises(ValueError):
        Fence.parse(doc)


def test_mavlink_altitude_fields_accepted():
    f = Fence.parse(json.dumps({"keep_in": [{"id": "c", "circle": {"lat": 0, "lon": 0, "radius_m": 10}}],
                                "max_hae_m": 120, "min_hae_m": 0}))
    assert f.check(0, 0) is None


def test_distance():
    assert distance(0, 0, 0, 1) == pytest.approx(111195, rel=1e-3)


@pytest.mark.parametrize("bearing", [0.0, 37.5, 90.0, 200.0, 359.0])
def test_frame_roundtrip(bearing):
    fr = LocalFrame(30.2672, -97.7431, 150.0, bearing)
    for x, y in ((0, 0), (10, -3), (-250.5, 400.25)):
        lat, lon, hae = fr.to_wgs84(x, y, 1.0)
        assert hae == 151.0
        bx, by = fr.from_wgs84(lat, lon)
        assert (bx, by) == pytest.approx((x, y), abs=1e-6)
        # haversine is spherical; the frame uses WGS-84 radii (meridional radius ~0.3% shorter at 30N)
        assert distance(30.2672, -97.7431, lat, lon) == pytest.approx(math.hypot(x, y), rel=5e-3, abs=1e-3)


def test_frame_axes():
    enu = LocalFrame(30.0, -97.0, 0.0, 90.0)
    lat, lon, _ = enu.to_wgs84(100, 0)  # +x east
    assert lat == pytest.approx(30.0) and lon > -97.0
    lat, lon, _ = enu.to_wgs84(0, 100)  # +y north
    assert lat > 30.0 and lon == pytest.approx(-97.0)
    north = LocalFrame(30.0, -97.0, 0.0, 0.0)
    lat, lon, _ = north.to_wgs84(0, 100)  # +x north, so +y (left) is west
    assert lon < -97.0 and lat == pytest.approx(30.0)
    assert enu.yaw_to_bearing(0) == 90.0 and enu.yaw_to_bearing(math.pi / 2) == pytest.approx(0.0)
    assert enu.bearing_to_yaw(enu.yaw_to_bearing(0.3)) == pytest.approx(0.3)


@pytest.mark.parametrize("args", [(91, 0), (0, 181), (float("nan"), 0), (0, 0, 0, 361)])
def test_frame_rejects_bad_origin(args):
    with pytest.raises(ValueError):
        LocalFrame(*args)


def test_check_path_keep_out():
    from ros2bridge.fence import check_path
    f = Fence.load(EXAMPLE)
    pond = f.keep_out[0]
    clat, clon = pond.circle[:2] if pond.circle else pond.polygon[0]
    # a straight line through the keep-out centre is refused; a line that stays clear passes
    v = check_path(f, clat - 0.0002, clon, clat + 0.0002, clon)
    assert v is not None and v.detail == "keep_out" and v.zone_id == pond.id
    assert check_path(f, clat - 0.0002, clon - 0.001, clat - 0.0001, clon - 0.001) is None
    # starting inside the zone: leaving it is allowed
    assert check_path(f, clat, clon, clat - 0.0003, clon) is None
    assert check_path(None, 0, 0, 1, 1) is None
