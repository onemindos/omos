"""
OneMind OS — DimOS NATS Bridge Blueprints
==========================================
Ready-to-run blueprints that wire a DimOS robot to the OneMind fabric.

Run with:
    # Simulation (no hardware)
    ROBOT_ID=go2-alpha dimos --simulation run omos-go2-nats

    # Real hardware
    ROBOT_ID=go2-alpha ROBOT_IP=192.168.1.xxx dimos run omos-go2-nats

    # Replay a recorded session
    ROBOT_ID=go2-alpha dimos --replay run omos-go2-nats
"""

from __future__ import annotations

import os

from dimos.core.coordination.blueprints import autoconnect

from omos.dimos.nats_bridge.module import NATSBridgeModule


def omos_go2_nats_blueprint(
    robot_id: str | None = None,
    nats_url: str | None = None,
    publish_hz: float = 10.0,
    cot_hz: float = 1.0,
):
    """
    Full OneMind-connected Go2 blueprint.
    Includes: Go2 connection + Navigation + NATS Bridge.

    The NATS bridge publishes:
        fabric.robot.<robot_id>.odom   ← raw odometry JSON
        fabric.robot.<robot_id>.cot    ← TAK CoT XML (appears on ATAK/WebTAK map)
        fabric.robot.<robot_id>.status ← online heartbeat

    Legacy can command the robot via:
        fabric.robot.<robot_id>.cmd    → {"linear": {"x":0.5}, "angular": {"z":0.3}}
    """
    from dimos.robot.unitree.go2.blueprints.smart.unitree_go2 import unitree_go2

    _robot_id = robot_id or os.environ.get("ROBOT_ID", "go2-alpha")
    _nats_url = nats_url or os.environ.get("NATS_URL", "nats://localhost:4222")

    bridge = NATSBridgeModule.blueprint(
        robot_id=_robot_id,
        nats_url=_nats_url,
        publish_hz=publish_hz,
        cot_hz=cot_hz,
    )

    return autoconnect(unitree_go2(), bridge)


def omos_go2_sim_nats_blueprint(
    robot_id: str | None = None,
    nats_url: str | None = None,
):
    """
    Simulation-only blueprint (MuJoCo) with NATS bridge.
    Identical to omos_go2_nats_blueprint but forces simulation transport.
    """
    from dimos.robot.unitree.go2.blueprints.sim.unitree_go2_sim import unitree_go2_sim

    _robot_id = robot_id or os.environ.get("ROBOT_ID", "go2-sim")
    _nats_url = nats_url or os.environ.get("NATS_URL", "nats://localhost:4222")

    bridge = NATSBridgeModule.blueprint(
        robot_id=_robot_id,
        nats_url=_nats_url,
    )

    return autoconnect(unitree_go2_sim(), bridge)


# DimOS blueprint registry — these names become `dimos run <name>`
BLUEPRINTS = {
    "omos-go2-nats":     omos_go2_nats_blueprint,
    "omos-go2-sim-nats": omos_go2_sim_nats_blueprint,
}