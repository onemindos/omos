import base64

import pytest

from cmdtest import new_key
from ros2bridge import config

PUB, _ = new_key()
BASE = {"ASSET_ID": "sim-tb1", "MAP_ORIGIN_LAT": "30.2672", "MAP_ORIGIN_LON": "-97.7431",
        "POLICY_PUBKEYS": "policy-1=" + base64.b64encode(PUB.public_bytes_raw()).decode()}


def test_defaults():
    s = config.load(dict(BASE))
    b = s.bridge
    assert s.transport == "zenoh" and s.zenoh_connect == ["tcp/127.0.0.1:7447"] and s.zenoh_mode == "client"
    assert b.sig_required and b.callsign == "SIM-TB1" and b.position_interval == 0.5 and b.use_odom_position
    assert b.fence is None and b.frame.x_bearing_deg == 90.0


def test_overrides():
    s = config.load(dict(BASE, TRANSPORT="rosbridge", POSE_TOPIC="amcl_pose", POSITION_INTERVAL="250ms",
                         LINK_TIMEOUT="2s", CMD_MAX_TTL="1m", ROS_NAMESPACE="/bot1", HOME_LAT="30.1", HOME_LON="-97.1",
                         GEOFENCE_JSON='{"keep_in":[{"id":"c","circle":{"lat":30,"lon":-97,"radius_m":9}}]}'))
    b = s.bridge
    assert s.transport == "rosbridge" and not b.use_odom_position and b.position_interval == 0.25
    assert b.link_timeout == 2.0 and b.max_ttl.total_seconds() == 60 and b.home == (30.1, -97.1)
    assert s.namespace == "/bot1" and len(b.fence.keep_in) == 1


@pytest.mark.parametrize("patch", [
    {"ASSET_ID": "tb1"}, {"ASSET_ID": ""}, {"MAP_ORIGIN_LAT": ""}, {"MAP_ORIGIN_LAT": "x"},
    {"TRANSPORT": "dds"}, {"ZENOH_MODE": "router"}, {"POSITION_INTERVAL": "fast"}, {"STALE_S": "0"},
    {"POLICY_PUBKEYS": ""}, {"HOME_LAT": "30"}, {"MAP_X_BEARING_DEG": "400"},
])
def test_rejected(patch):
    with pytest.raises(ValueError):
        config.load(dict(BASE, **patch))


def test_unsigned_dev_mode_needs_explicit_flag():
    s = config.load(dict(BASE, POLICY_PUBKEYS="", SIG_REQUIRED="false"))
    assert s.bridge.sig_required is False
