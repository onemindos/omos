"""Robot-side transports. The bridge talks to ROS 2 through one of these.

- ``zenoh`` (default): a Zenoh session connected to ``zenoh-bridge-ros2dds`` running
  beside the robot. Topics arrive as CDR; ``navigate_to_pose`` is an action, which the
  ros2dds bridge exposes as Zenoh queryables (``.../_action/send_goal`` and
  ``.../_action/cancel_goal``), so a goal gets a real accepted/rejected answer.
- ``rosbridge`` (fallback): rosbridge_suite's JSON protocol over WebSocket, for robots
  where the Zenoh bridge cannot run. Goals go to the ``goal_pose`` topic, which Nav2's
  BT navigator subscribes to; there is no goal acknowledgement on this path.
"""

from __future__ import annotations

import abc
from dataclasses import dataclass


@dataclass
class Topics:
    odom: str = "odom"
    pose: str = ""  # optional map-frame pose source, e.g. amcl_pose (PoseWithCovarianceStamped)
    battery: str = "battery_state"
    cmd_vel: str = "cmd_vel"
    nav_action: str = "navigate_to_pose"
    goal_topic: str = "goal_pose"  # rosbridge fallback only


@dataclass
class GoalResult:
    accepted: bool
    goal_id: str  # hex of the 16-byte action goal id, "" when the transport has none
    detail: str = ""


class TransportTimeout(Exception):
    """The robot side did not answer in time (no_ack)."""


class TransportError(Exception):
    """The request could not be sent (link down, encoding error)."""


class Transport(abc.ABC):
    name = "abstract"

    @abc.abstractmethod
    async def start(self, on_odom, on_pose, on_battery) -> None:
        """Begin delivering decoded messages. Callbacks run on the asyncio loop."""

    @abc.abstractmethod
    async def send_goal(self, goal) -> GoalResult:
        """Send a NavigateToPose goal (msgs.PoseStamped in the map frame)."""

    @abc.abstractmethod
    async def hold(self, here) -> str:
        """Stop navigating and stay at ``here`` (msgs.PoseStamped). Returns a detail string."""

    @abc.abstractmethod
    async def estop(self, here=None) -> str:
        """Command zero velocity, then cancel every goal. ``here`` (msgs.PoseStamped, the
        current pose) lets a transport without a goal cancel re-target the goal there.
        Must not wait on the robot before the zero velocity is sent."""

    @abc.abstractmethod
    async def close(self) -> None:
        ...

    def stats(self) -> dict:
        return {}
