"""Configuration from the environment only (README.md lists every variable).

Names shared with the MAVLink bridge (omos PR #550) mean the same thing: NATS_*, ASSET_ID,
CALLSIGN, BRIDGE_NAME, POLICY_PUBKEYS, SIG_REQUIRED, GEOFENCE_FILE / GEOFENCE_JSON,
CMD_MAX_TTL, CLOCK_SKEW, POSITION_INTERVAL, STATUS_INTERVAL, STALE_S, LINK_TIMEOUT.
"""

from __future__ import annotations

import os
import re
from dataclasses import dataclass, field
from datetime import timedelta

from .bridge import BridgeConfig
from .command import parse_keys
from .fence import Fence
from .geo import LocalFrame
from .transport import Topics


def _dur(env, k: str, default: float) -> float:
    """Go-style durations (500ms, 3s, 1m) or plain seconds."""
    v = env.get(k, "")
    if not v:
        return default
    m = re.fullmatch(r"(\d+(?:\.\d+)?)(ms|s|m)?", v.strip())
    if not m:
        raise ValueError(f"{k}: invalid duration {v!r}")
    n = float(m.group(1)) * {"ms": 0.001, "s": 1.0, "m": 60.0, None: 1.0}[m.group(2)]
    if n <= 0:
        raise ValueError(f"{k}: duration must be positive")
    return n


def _float(env, k: str, default=None, required=False) -> float | None:
    v = env.get(k, "")
    if v == "":
        if required:
            raise ValueError(f"{k} is required")
        return default
    try:
        return float(v)
    except ValueError:
        raise ValueError(f"{k}: not a number: {v!r}") from None


def _int(env, k: str, default: int, lo: int, hi: int) -> int:
    v = env.get(k, "")
    if v == "":
        return default
    try:
        n = int(v)
    except ValueError:
        n = lo - 1
    if not lo <= n <= hi:
        raise ValueError(f"{k}: want integer in [{lo},{hi}], got {v!r}")
    return n


@dataclass
class Settings:
    nats_url: str
    nats_user: str
    nats_password: str
    nats_creds: str
    transport: str  # zenoh | rosbridge
    zenoh_connect: list = field(default_factory=list)
    zenoh_listen: list = field(default_factory=list)
    zenoh_mode: str = "client"
    namespace: str = ""
    rosbridge_url: str = ""
    topics: Topics = field(default_factory=Topics)
    query_timeout: float = 3.0
    bridge: BridgeConfig | None = None


def load(env=None) -> Settings:
    env = os.environ if env is None else env
    asset = env.get("ASSET_ID", "")
    if not asset.startswith("sim-"):
        raise ValueError(f"ASSET_ID {asset!r} must start with sim- (simulated assets only, decision 0004 #13)")
    frame = LocalFrame(_float(env, "MAP_ORIGIN_LAT", required=True), _float(env, "MAP_ORIGIN_LON", required=True),
                       _float(env, "MAP_ORIGIN_HAE", 0.0), _float(env, "MAP_X_BEARING_DEG", 90.0))
    fence = None
    if env.get("GEOFENCE_FILE"):
        fence = Fence.load(env["GEOFENCE_FILE"])
    elif env.get("GEOFENCE_JSON"):
        fence = Fence.parse(env["GEOFENCE_JSON"])
    home = None
    if env.get("HOME_LAT") or env.get("HOME_LON"):
        home = (_float(env, "HOME_LAT", required=True), _float(env, "HOME_LON", required=True))
    topics = Topics(odom=env.get("ODOM_TOPIC", "odom"), pose=env.get("POSE_TOPIC", ""),
                    battery=env.get("BATTERY_TOPIC", "battery_state"), cmd_vel=env.get("CMD_VEL_TOPIC", "cmd_vel"),
                    nav_action=env.get("NAV_ACTION", "navigate_to_pose"), goal_topic=env.get("GOAL_TOPIC", "goal_pose"))
    transport = env.get("TRANSPORT", "zenoh")
    if transport not in ("zenoh", "rosbridge"):
        raise ValueError("TRANSPORT must be zenoh or rosbridge")
    mode = env.get("ZENOH_MODE", "client")
    if mode not in ("client", "peer"):
        raise ValueError("ZENOH_MODE must be client or peer")
    split = lambda s: [x.strip() for x in s.split(",") if x.strip()]  # noqa: E731
    bc = BridgeConfig(
        asset_id=asset, frame=frame, callsign=env.get("CALLSIGN", ""),
        source=env.get("BRIDGE_NAME", "bridge-ros2"), kind=env.get("ROBOT_KIND", "ugv"), fence=fence,
        keys=parse_keys(env.get("POLICY_PUBKEYS", "")), sig_required=env.get("SIG_REQUIRED", "true") != "false",
        max_ttl=timedelta(seconds=_dur(env, "CMD_MAX_TTL", 60)),
        clock_skew=timedelta(seconds=_dur(env, "CLOCK_SKEW", 5)),
        position_interval=_dur(env, "POSITION_INTERVAL", 0.5), status_interval=_dur(env, "STATUS_INTERVAL", 1.0),
        stale_s=_int(env, "STALE_S", 10, 1, 3600), link_timeout=_dur(env, "LINK_TIMEOUT", 3.0),
        goal_tolerance_m=_float(env, "GOAL_TOLERANCE_M", 0.3), home=home,
        map_frame=env.get("MAP_FRAME", "map"), use_odom_position=not topics.pose,
    )
    if bc.sig_required and not bc.keys:
        raise ValueError("POLICY_PUBKEYS is required unless SIG_REQUIRED=false (local development only)")
    return Settings(
        nats_url=env.get("NATS_URL", "nats://127.0.0.1:4222"), nats_user=env.get("NATS_USER", ""),
        nats_password=env.get("NATS_PASSWORD", ""), nats_creds=env.get("NATS_CREDS", ""),
        transport=transport, zenoh_connect=split(env.get("ZENOH_CONNECT", "tcp/127.0.0.1:7447")),
        zenoh_listen=split(env.get("ZENOH_LISTEN", "")), zenoh_mode=mode,
        namespace=env.get("ROS_NAMESPACE", ""), rosbridge_url=env.get("ROSBRIDGE_URL", "ws://127.0.0.1:9090"),
        topics=topics, query_timeout=_dur(env, "ACTION_TIMEOUT", 3.0), bridge=bc,
    )
