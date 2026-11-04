"""CDR codec tests. The expected bytes are built by hand from the XCDR1 rules (alignment
relative to the end of the 4-byte encapsulation header), independently of cdr.Writer."""

import math
import struct

import pytest

from ros2bridge import msgs
from ros2bridge.cdr import CDRError, Reader


def manual_odometry() -> bytes:
    b = bytearray(b"\x00\x01\x00\x00")
    body = bytearray()
    body += struct.pack("<iI", 1793280000, 500)            # 0..8   stamp
    body += struct.pack("<I", 5) + b"odom\x00"             # 8..17  frame_id
    body += b"\x00" * 3                                    # 17..20 align 4
    body += struct.pack("<I", 15) + b"base_footprint\x00"  # 20..39 child_frame_id
    body += b"\x00"                                        # 39..40 align 8
    body += struct.pack("<7d", 1.5, -2.25, 0.0, 0.0, 0.0, math.sin(0.25), math.cos(0.25))
    body += struct.pack("<36d", *([0.0] * 36))
    body += struct.pack("<6d", 0.2, 0.0, 0.0, 0.0, 0.0, 0.1)
    body += struct.pack("<36d", *([0.0] * 36))
    assert len(body) == 720
    return bytes(b + body)


def test_odometry_matches_hand_built_bytes():
    o = msgs.decode_odometry(manual_odometry())
    assert o.header.frame_id == "odom" and o.child_frame_id == "base_footprint"
    assert o.header.stamp.sec == 1793280000 and o.header.stamp.nanosec == 500
    assert (o.pose.x, o.pose.y) == (1.5, -2.25)
    assert o.pose.yaw() == pytest.approx(0.5)
    assert (o.twist.lx, o.twist.az) == (0.2, 0.1)
    assert msgs.encode_odometry(o) == manual_odometry()


def test_big_endian_is_read():
    le = manual_odometry()
    be = bytearray(b"\x00\x00\x00\x00")
    # re-encode the same message big-endian by hand
    body = struct.pack(">iI", 1793280000, 500) + struct.pack(">I", 5) + b"odom\x00" + b"\x00" * 3
    body += struct.pack(">I", 15) + b"base_footprint\x00" + b"\x00"
    body += struct.pack(">7d", 1.5, -2.25, 0.0, 0.0, 0.0, math.sin(0.25), math.cos(0.25))
    body += struct.pack(">36d", *([0.0] * 36)) + struct.pack(">6d", 0.2, 0, 0, 0, 0, 0.1)
    body += struct.pack(">36d", *([0.0] * 36))
    assert msgs.decode_odometry(bytes(be + body)) == msgs.decode_odometry(le)


@pytest.mark.parametrize("buf", [b"", b"\x00\x01", b"\x00\x07\x00\x00" + b"\x00" * 64,
                                 manual_odometry()[:100], b"\x00\x01\x00\x00\x00\x00\x00\x00\xff\x00\x00\x00ab"])
def test_malformed_buffers_raise(buf):
    with pytest.raises(CDRError):
        msgs.decode_odometry(buf)


def test_string_needs_nul():
    with pytest.raises(CDRError):
        Reader(b"\x00\x01\x00\x00" + struct.pack("<I", 2) + b"ab").string()


def test_pose_stamped_and_twist_roundtrip():
    p = msgs.PoseStamped(msgs.Header(msgs.Time(3, 4), "map"), msgs.Pose.from_xy_yaw(3.0, 4.0, 1.0))
    assert msgs.decode_pose_stamped(msgs.encode_pose_stamped(p)) == p
    t = msgs.Twist(0.1, 0, 0, 0, 0, -0.3)
    assert msgs.decode_twist(msgs.encode_twist(t)) == t
    assert len(msgs.encode_twist(msgs.Twist())) == 4 + 48


def test_pose_with_covariance_roundtrip():
    p = msgs.PoseStamped(msgs.Header(msgs.Time(1, 2), "map"), msgs.Pose.from_xy_yaw(-1.0, 2.0, -0.5))
    assert msgs.decode_pose_with_covariance_stamped(msgs.encode_pose_with_covariance_stamped(p)) == p


def test_battery_state_roundtrip_and_layout():
    b = msgs.BatteryState(msgs.Header(msgs.Time(1, 0), "base_link"), 12.1, 0.875, True)
    raw = msgs.encode_battery_state(b)
    d = msgs.decode_battery_state(raw)
    assert d.voltage == pytest.approx(12.1, abs=1e-5) and d.percentage == pytest.approx(0.875) and d.present
    # header 8 + "base_link\0" (4+10) -> 22, align 4 -> 24, 7 float32 -> 52, 4 bytes -> 56,
    # two empty sequences -> 64, location (len + NUL) -> 69, align 4 -> 72, serial_number -> 77
    assert len(raw) == 4 + 77


def test_action_messages_roundtrip():
    gid = bytes(range(16))
    g = msgs.SendGoalRequest(gid, msgs.PoseStamped(msgs.Header(msgs.Time(5, 6), "map"),
                                                   msgs.Pose.from_xy_yaw(1, 2, 0)), "")
    raw = msgs.encode_send_goal_request(g)
    assert raw[4:20] == gid  # UUID is the first member, no alignment before it
    assert msgs.decode_send_goal_request(raw) == g
    r = msgs.SendGoalResponse(True, msgs.Time(7, 8))
    assert msgs.decode_send_goal_response(msgs.encode_send_goal_response(r)) == r
    assert msgs.encode_send_goal_response(r)[4:] == b"\x01\x00\x00\x00" + struct.pack("<iI", 7, 8)
    c = msgs.CancelGoalResponse(0, [(gid, msgs.Time(1, 1))])
    assert msgs.decode_cancel_goal_response(msgs.encode_cancel_goal_response(c)) == c
    assert msgs.decode_cancel_goal_request(msgs.encode_cancel_goal_request()) == (b"\x00" * 16, msgs.Time())
    with pytest.raises(ValueError):
        msgs.encode_send_goal_request(msgs.SendGoalRequest(b"short", g.goal))
