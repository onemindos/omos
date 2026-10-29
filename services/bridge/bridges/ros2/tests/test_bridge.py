"""Bridge core tests with a fake transport and bus (no network)."""

import asyncio
import json
from datetime import timedelta

import pytest

from cmdtest import build, new_key
from ros2bridge import command as C
from ros2bridge import msgs
from ros2bridge.bridge import Bridge, BridgeConfig
from ros2bridge.fence import Fence
from ros2bridge.geo import LocalFrame
from ros2bridge.transport import GoalResult, Topics, TransportError, TransportTimeout

PUB, PRIV = new_key()
FRAME = LocalFrame(30.2672, -97.7431, 150.0, 90.0)
FENCE = Fence.parse(json.dumps({"keep_in": [{"id": "yard", "circle": {"lat": 30.2672, "lon": -97.7431, "radius_m": 30}}],
                                "keep_out": [{"id": "pond", "circle": {"lat": 30.26729, "lon": -97.7431, "radius_m": 3}}]}))


class Clock:
    def __init__(self):
        self.t = C.utcnow()

    def __call__(self):
        return self.t


class FakeBus:
    def __init__(self):
        self.msgs = []

    async def publish(self, subject, data, msg_id=""):
        self.msgs.append((subject, json.loads(data), msg_id))

    def on(self, subject):
        return [m for s, m, _ in self.msgs if s == subject]


class FakeTransport:
    name = "zenoh"

    def __init__(self):
        self.topics = Topics()
        self.goals, self.holds, self.estops = [], [], 0
        self.accept = True
        self.fail = None

    async def start(self, *a):
        pass

    async def send_goal(self, goal):
        if self.fail:
            raise self.fail
        self.goals.append(goal)
        return GoalResult(self.accept, "ab" * 16, "fake/navigate_to_pose/_action/send_goal")

    async def hold(self, here):
        self.holds.append(here)
        return "canceled"

    async def estop(self, here=None):
        self.estops += 1
        if self.fail:
            raise self.fail
        return "zeroed"

    async def close(self):
        pass

    def stats(self):
        return {"rx": 1}


def make(**kw):
    clock = Clock()
    cfg = dict(asset_id="sim-tb1", frame=FRAME, fence=FENCE, keys={"policy-test": PUB}, position_interval=0.0)
    cfg.update(kw)
    bus, tr = FakeBus(), FakeTransport()
    return Bridge(BridgeConfig(**cfg), bus, tr, clock=clock), bus, tr, clock


def odom(x=0.0, y=0.0, yaw=0.0, v=0.0):
    return msgs.Odometry(msgs.Header(msgs.Time(1, 0), "odom"), "base_footprint", msgs.Pose.from_xy_yaw(x, y, yaw),
                         msgs.Twist(lx=v))


def run(coro):
    return asyncio.run(coro)


async def settle(b):
    for _ in range(5):
        await asyncio.sleep(0)
    if b._tasks:
        await asyncio.gather(*list(b._tasks))


def outcomes(bus):
    return [(m["data"]["verb"], m["data"]["outcome"], m["data"].get("code", ""), m["data"].get("detail", ""))
            for m in bus.on("evt.robot.sim-tb1.command")]


def target(dx, dy):
    lat, lon, _ = FRAME.to_wgs84(dx, dy)
    return {"lat": lat, "lon": lon, "hae_m": 150.0}


def test_config_refuses_non_sim():
    with pytest.raises(ValueError):
        BridgeConfig(asset_id="tb1", frame=FRAME)


def test_telemetry_envelopes():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom(10.0, 5.0, 0.0, 0.2))
        b.on_battery(msgs.BatteryState(msgs.Header(), 12.0, 0.8, True))
        await b.publish_status()
        await settle(b)
        return b, bus

    b, bus = run(go())
    (pos,) = bus.on("tel.fleet.robot.sim-tb1.position")
    lat, lon, hae = FRAME.to_wgs84(10.0, 5.0)
    assert pos["uid"] == "sim-tb1" and pos["type"] == "fleet.robot"
    assert pos["lat"] == pytest.approx(lat) and pos["lon"] == pytest.approx(lon) and pos["hae"] == 150.0
    assert pos["course"] == 90.0 and pos["speed"] == 0.2  # +x is east, moving forward
    assert pos["meta"]["kind"] == "ugv" and pos["meta"]["simulated"] is True and pos["t"] > 0
    (st,) = bus.on("tel.fleet.robot.sim-tb1.status")
    assert st["meta"]["connection"] == "online" and st["meta"]["battery_pct"] == 80
    assert st["meta"]["commands_accepted"] == "goto,hold,rtl,estop,estop_release"
    assert st["meta"]["home_lat"] == pytest.approx(lat)  # the first fix becomes home


def test_entity_service_fleet_parser_contract():
    """Fields the entity service (PR #552, internal/adapt/geo.go Fleet) requires."""
    async def go():
        b, bus, *_ = make()
        b.on_odom(odom())
        await settle(b)
        return bus
    (pos,) = run(go()).on("tel.fleet.robot.sim-tb1.position")
    assert pos["uid"].startswith("sim-") and isinstance(pos["t"], int) and pos["t"] > 0
    assert isinstance(pos["lat"], float) and isinstance(pos["lon"], float) and isinstance(pos["meta"], dict)


def test_position_rate_limit_and_stationary_course():
    async def go():
        b, bus, tr, clock = make(position_interval=1.0)
        for i in range(5):
            b.on_odom(odom(0.01 * i))
            clock.t += timedelta(milliseconds=300)
        await settle(b)
        return bus
    pos = run(go()).on("tel.fleet.robot.sim-tb1.position")
    assert len(pos) == 2 and pos[0]["course"] == -1


def test_no_status_without_position_and_nan_pose_ignored():
    async def go():
        b, bus, *_ = make()
        b.on_odom(odom(float("nan")))
        await b.publish_status()
        return b, bus
    b, bus = run(go())
    assert not b.has_pos and bus.msgs == []


def test_goto_flow_reaches_goal_and_reports_task():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        subj, body = build(PRIV, verb="goto", data=target(10.0, -4.0), now=clock.t, task_id="tsk-1")
        await b.handle_command(subj, body)
        g = tr.goals[0]
        b.on_odom(odom(9.9, -3.95))
        await settle(b)
        return b, bus, g
    b, bus, g = run(go())
    assert g.header.frame_id == "map" and g.pose.x == pytest.approx(10.0, abs=1e-6) and g.pose.y == pytest.approx(-4.0, abs=1e-6)
    assert outcomes(bus) == [("goto", "acknowledged", "", ""), ("goto", "executed", "", "")]
    ex = bus.on("evt.robot.sim-tb1.command")[1]
    assert ex["data"]["ros"]["goal_id"] == "ab" * 16 and ex["causation"] and ex["source"] == "bridge-ros2"
    states = [m["data"]["state"] for m in bus.on("evt.robot.sim-tb1.task")]
    assert states == ["received", "executing", "done_ok"]
    vers = [m["data"]["status_version"] for m in bus.on("evt.robot.sim-tb1.task")]
    assert vers == [1, 2, 3] and b.mode == "idle" and b.goal is None
    ids = [mid for s, _, mid in bus.msgs if s == "evt.robot.sim-tb1.task"]
    assert ids == ["tsk-1:1", "tsk-1:2", "tsk-1:3"]


@pytest.mark.parametrize("dxdy,detail", [((100.0, 0.0), "keep_in"), ((0.0, 10.0), "keep_out")])
def test_goto_outside_fence_refused(dxdy, detail):
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        await b.handle_command(*build(PRIV, verb="goto", data=target(*dxdy), now=clock.t))
        return bus, tr
    bus, tr = run(go())
    assert outcomes(bus) == [("goto", "refused", C.PRECONDITION, detail)] and tr.goals == []


def test_no_fence_fails_closed():
    async def go():
        b, bus, tr, clock = make(fence=None)
        b.on_odom(odom())
        await b.handle_command(*build(PRIV, verb="goto", data=target(1, 1), now=clock.t))
        return bus
    assert outcomes(run(go()))[0][1:] == ("refused", C.PRECONDITION, "keep_in")


def test_offline_refused_but_estop_never_blocked():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        clock.t += timedelta(seconds=10)
        await b.handle_command(*build(PRIV, verb="goto", data=target(1, 1), now=clock.t))
        await b.handle_command(*build(PRIV, verb="estop", now=clock.t))
        return bus, tr
    bus, tr = run(go())
    assert outcomes(bus) == [("goto", "refused", C.PRECONDITION, "offline"),
                             ("estop", "acknowledged", "", ""), ("estop", "executed", "", "")]
    assert tr.estops == 1


def test_estop_latch_and_release():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        await b.handle_command(*build(PRIV, verb="estop", now=clock.t))
        await b.handle_command(*build(PRIV, verb="goto", data=target(1, 1), now=clock.t))
        await b.handle_command(*build(PRIV, verb="rtl", now=clock.t))
        await b.handle_command(*build(PRIV, verb="hold", now=clock.t))  # a stop is still allowed
        await b.handle_command(*build(PRIV, verb="estop_release", approvers=["user:zeus"], now=clock.t))
        await b.handle_command(*build(PRIV, verb="estop_release", now=clock.t))
        await b.handle_command(*build(PRIV, verb="goto", data=target(1, 1), now=clock.t))
        return bus, tr
    bus, tr = run(go())
    o = [x for x in outcomes(bus) if x[1] != "acknowledged"]
    assert o == [("estop", "executed", "", ""), ("goto", "refused", C.PRECONDITION, "estopped"),
                 ("rtl", "refused", C.PRECONDITION, "estopped"), ("hold", "executed", "", ""),
                 ("estop_release", "refused", C.NOT_AUTHORIZED, ""), ("estop_release", "executed", "", ""),
                 ("goto", "executed", "", "")]
    assert len(tr.goals) == 1


def test_estop_latches_even_if_transport_fails():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        tr.fail = TransportError("down")
        await b.handle_command(*build(PRIV, verb="estop", now=clock.t))
        return b, bus
    b, bus = run(go())
    assert b.estopped and outcomes(bus)[-1][1:3] == ("failed", C.NO_ACK)


def test_rtl_goes_home_and_hold_cancels():
    async def go():
        b, bus, tr, clock = make(home=(30.2672, -97.7431))
        b.on_odom(odom(5.0, 5.0))
        await b.handle_command(*build(PRIV, verb="rtl", now=clock.t))
        await b.handle_command(*build(PRIV, verb="hold", data={"duration_s": 10}, now=clock.t))
        return b, tr
    b, tr = run(go())
    assert tr.goals[0].pose.x == pytest.approx(0.0, abs=1e-6) and tr.goals[0].pose.y == pytest.approx(0.0, abs=1e-6)
    assert len(tr.holds) == 1 and tr.holds[0].pose.x == 5.0 and b.mode == "holding" and b.goal is None


@pytest.mark.parametrize("fail,accept,code", [(None, False, C.VEHICLE_REJECTED), (TransportTimeout("t"), True, C.NO_ACK),
                                              (TransportError("e"), True, C.NO_ACK)])
def test_goal_failures(fail, accept, code):
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        tr.fail, tr.accept = fail, accept
        await b.handle_command(*build(PRIV, verb="goto", data=target(2, 2), now=clock.t, task_id="tsk-2"))
        return b, bus
    b, bus = run(go())
    assert outcomes(bus)[-1][1:3] == ("failed", code) and b.goal is None
    assert bus.on("evt.robot.sim-tb1.task")[-1]["data"]["state"] == "done_not_ok"


def test_refusals_publish_outcomes_and_never_move():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        await b.handle_command("cmd.fleet.robot.sim-tb1.goto", b"not json")
        subj, body = build(PRIV, verb="goto", data=target(1, 1), now=clock.t, no_sig=True)
        await b.handle_command(subj, body)
        subj, body = build(PRIV, verb="goto", data=target(1, 1), now=clock.t)
        await b.handle_command(subj, body)
        await b.handle_command(subj, body)  # replay
        return bus, tr
    bus, tr = run(go())
    o = outcomes(bus)
    assert o[0][1:3] == ("refused", C.SCHEMA_INVALID) and o[1][1:3] == ("refused", C.BAD_SIGNATURE)
    assert o[-1][1:3] == ("refused", C.REPLAYED) and len(tr.goals) == 1


def test_dry_run_does_not_move():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        await b.handle_command(*build(PRIV, verb="goto", data=target(1, 1), now=clock.t,
                                      mutate=lambda e: e.update(dry_run=True)))
        return bus, tr
    bus, tr = run(go())
    assert outcomes(bus) == [("goto", "acknowledged", "", "")] and tr.goals == []


def test_command_expiring_behind_lock_fails():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        await b._move_lock.acquire()
        subj, body = build(PRIV, verb="goto", data=target(1, 1), now=clock.t, ttl=timedelta(seconds=2))
        t = asyncio.ensure_future(b.handle_command(subj, body))
        await asyncio.sleep(0.01)
        clock.t += timedelta(seconds=3)
        b._move_lock.release()
        await t
        return bus, tr
    bus, tr = run(go())
    assert outcomes(bus)[-1][1:3] == ("failed", C.EXPIRED) and tr.goals == []


def test_estop_does_not_wait_for_a_movement_command():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        await b._move_lock.acquire()  # a goto in progress
        await asyncio.wait_for(b.handle_command(*build(PRIV, verb="estop", now=clock.t)), 1.0)
        b._move_lock.release()
        return tr
    assert run(go()).estops == 1


def test_estop_overtakes_inflight_and_queued_goto():
    """Regression (H3 review blocker): an estop while one goto waits on Nav2 and another
    waits on the lock. The queued goto is never sent; the in-flight goal is canceled
    again once Nav2 answers; neither records a goal or overwrites the estop mode."""
    async def go():
        b, bus, tr, clock = make()
        gate = asyncio.Event()
        orig = tr.send_goal

        async def slow(goal):
            await gate.wait()  # Nav2 answers slowly (a zenoh query can take up to 3 s)
            return await orig(goal)
        tr.send_goal = slow
        b.on_odom(odom())
        b._spawn(b.handle_command(*build(PRIV, verb="goto", data=target(5, 0), now=clock.t)))
        b._spawn(b.handle_command(*build(PRIV, verb="goto", data=target(-5, 0), now=clock.t)))
        for _ in range(10):
            await asyncio.sleep(0)
        await b.handle_command(*build(PRIV, verb="estop", now=clock.t))
        gate.set()
        await settle(b)
        return b, bus, tr
    b, bus, tr = run(go())
    assert b.estopped and b.mode == "estop" and b.goal is None
    assert len(tr.goals) == 1  # only the one already on its way to Nav2
    assert tr.estops == 2  # the estop, then the re-stop after the in-flight goal was accepted
    final = [x for x in outcomes(bus) if x[1] not in ("acknowledged",)]
    assert ("estop", "executed", "", "") in final
    assert [x for x in final if x[0] == "goto"] == [("goto", "failed", C.PRECONDITION, "estopped")] * 2


def test_estop_overtakes_goto_that_times_out():
    async def go():
        b, bus, tr, clock = make()
        gate = asyncio.Event()

        async def slow(goal):
            await gate.wait()
            raise TransportTimeout("no reply")
        tr.send_goal = slow
        b.on_odom(odom())
        b._spawn(b.handle_command(*build(PRIV, verb="goto", data=target(5, 0), now=clock.t)))
        for _ in range(10):
            await asyncio.sleep(0)
        await b.handle_command(*build(PRIV, verb="estop", now=clock.t))
        gate.set()
        await settle(b)
        return b, bus, tr
    b, bus, tr = run(go())
    assert tr.estops == 2 and b.goal is None and b.mode == "estop"
    assert outcomes(bus)[-1][:3] == ("goto", "failed", C.NO_ACK)


def test_hold_while_estopped_keeps_estop_mode():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        await b.handle_command(*build(PRIV, verb="estop", now=clock.t))
        await b.handle_command(*build(PRIV, verb="hold", now=clock.t))
        return b
    b = run(go())
    assert b.estopped and b.mode == "estop"


def test_goto_whose_straight_path_crosses_keep_out_refused():
    async def go():
        b, bus, tr, clock = make()
        b.on_odom(odom())
        await b.handle_command(*build(PRIV, verb="goto", data=target(0, 20), now=clock.t))  # pond at y=10
        b.on_odom(odom(0.0, 10.0))  # drifted into the pond: driving out is allowed
        await b.handle_command(*build(PRIV, verb="goto", data=target(0, -5), now=clock.t))
        return bus, tr
    bus, tr = run(go())
    o = [x for x in outcomes(bus) if x[1] != "acknowledged"]
    assert o == [("goto", "refused", C.PRECONDITION, "keep_out"), ("goto", "executed", "", "")]
    assert len(tr.goals) == 1


def test_amcl_pose_source():
    async def go():
        b, bus, tr, clock = make(use_odom_position=False)
        b.on_odom(odom(99.0, 99.0, 0.0, 0.1))
        assert not b.has_pos
        b.on_pose(msgs.PoseStamped(msgs.Header(), msgs.Pose.from_xy_yaw(1.0, 2.0, 0.0)))
        await settle(b)
        return b
    b = run(go())
    assert (b.x, b.y, b.vx) == (1.0, 2.0, 0.1)
