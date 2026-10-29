package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/bluenviron/gomavlib/v3/pkg/message"

	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/command"
	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/fence"
	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/geoid"
)

// Bus publishes to NATS. msgID, when set, becomes the Nats-Msg-Id header so JetStream
// de-duplicates (EVT has a 2 m duplicate window, platform/nats/streams/evt.yaml).
type Bus interface {
	Publish(subject, msgID string, data []byte) error
}

// Bridge turns MAVLink into fleet telemetry and verified fleet commands into MAVLink.
type Bridge struct {
	cfg      Config
	bus      Bus
	veh      Vehicle
	verifier *command.Verifier
	state    *vehicleState
	offset   atomic.Int64 // test clock offset in ns
	logf     func(string, ...any)

	moveMu sync.Mutex // one non-estop command at a time; estop never waits

	taskMu  sync.Mutex
	taskVer map[string]int // last status_version reported per task_id
	taskIDs []string       // insertion order, to bound taskVer
}

// maxTrackedTasks bounds the per-task status_version memory.
const maxTrackedTasks = 4096

// NewBridge validates the safety configuration and returns a bridge.
func NewBridge(cfg Config, bus Bus, veh Vehicle) (*Bridge, error) {
	ver, err := command.NewVerifier(command.Config{
		AssetClass: "drone", AssetID: cfg.AssetID, Keys: cfg.Keys, SigRequired: cfg.SigRequired,
		MaxTTL: cfg.MaxTTL, ClockSkew: cfg.ClockSkew,
	})
	if err != nil {
		return nil, err
	}
	if !cfg.SigRequired {
		log.Printf("[%s] WARNING: SIG_REQUIRED=false; unsigned commands are accepted (local development only, asset-command.md §6.3)", cfg.Source)
	}
	if cfg.Fence == nil || len(cfg.Fence.KeepIn) == 0 {
		log.Printf("[%s] WARNING: no keep-in zone; goto, takeoff, arm and land-at-point are refused (fail closed)", cfg.Source)
	}
	return &Bridge{cfg: cfg, bus: bus, veh: veh, verifier: ver, state: newVehicleState(), logf: log.Printf,
		taskVer: map[string]int{}}, nil
}

// Subjects (locked grammar; asset-command.md §1, robotics-cot-profile.md §1).
func (b *Bridge) positionSubject() string { return "tel.fleet.drone." + b.cfg.AssetID + ".position" }
func (b *Bridge) statusSubject() string   { return "tel.fleet.drone." + b.cfg.AssetID + ".status" }
func (b *Bridge) commandSubject() string  { return "cmd.fleet.drone." + b.cfg.AssetID + ".>" }
func (b *Bridge) outcomeSubject() string  { return "evt.drone." + b.cfg.AssetID + ".command" }
func (b *Bridge) taskSubject() string     { return "evt.drone." + b.cfg.AssetID + ".task" }

// OnMessage handles one MAVLink message from the configured vehicle.
func (b *Bridge) OnMessage(m message.Message) {
	now := b.clock()
	if b.state.update(m, now, b.cfg.PositionInterval) {
		env := b.state.positionEnvelope(&b.cfg, now)
		if err := b.bus.Publish(b.positionSubject(), "", mustJSON(env)); err != nil {
			b.logf("[%s] publish position: %v", b.cfg.Source, err)
		}
	}
}

// PublishStatus publishes tel.fleet.drone.{id}.status once.
func (b *Bridge) PublishStatus() {
	env, ok := b.state.statusEnvelope(&b.cfg, b.clock())
	if !ok {
		return
	}
	if err := b.bus.Publish(b.statusSubject(), "", mustJSON(env)); err != nil {
		b.logf("[%s] publish status: %v", b.cfg.Source, err)
	}
}

// RunStatus publishes status every StatusInterval until ctx ends.
func (b *Bridge) RunStatus(ctx context.Context) {
	t := time.NewTicker(b.cfg.StatusInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			b.PublishStatus()
		}
	}
}

// outcome is the data of evt.drone.{id}.command (asset-command.md §6.4).
type outcome struct {
	ActionID  string         `json:"action_id"`
	CommandID string         `json:"command_id"`
	Outcome   string         `json:"outcome"`
	Verb      string         `json:"verb,omitempty"`
	Code      string         `json:"code,omitempty"`
	Detail    string         `json:"detail,omitempty"`
	Message   string         `json:"message,omitempty"`
	Mavlink   map[string]any `json:"mavlink,omitempty"`
	SigOK     *bool          `json:"sig_verified,omitempty"`
	TaskID    string         `json:"task_id,omitempty"`
	Simulated bool           `json:"simulated"`
}

func uuid4() string {
	var u [16]byte
	_, _ = rand.Read(u[:])
	u[6] = u[6]&0x0f | 0x40
	u[8] = u[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

func rfc3339ms(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05.000Z07:00") }

// publishEvent publishes an omos event envelope v1 (services/schemas/envelope).
// msgID is the Nats-Msg-Id; empty means the envelope id.
func (b *Bridge) publishEvent(subject, typ, msgID string, e *command.Envelope, data any) {
	id := uuid4()
	if msgID == "" {
		msgID = id
	}
	env := map[string]any{
		"v": 1, "id": id, "ts": rfc3339ms(b.clock()), "type": typ, "source": b.cfg.Source,
		"subject": subject, "data": data,
	}
	if e != nil {
		if uuidLike(e.ID) {
			env["causation"] = e.ID
		}
		if uuidLike(e.Corr) {
			env["corr"] = e.Corr
		}
		if e.Trace != "" {
			env["trace"] = e.Trace
		}
	}
	if err := b.bus.Publish(subject, msgID, mustJSON(env)); err != nil {
		b.logf("[%s] publish %s: %v", b.cfg.Source, subject, err)
	}
}

func uuidLike(s string) bool {
	return len(s) == 36 && s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-'
}

func (b *Bridge) report(e *command.Envelope, o outcome) {
	if e != nil {
		o.ActionID, o.CommandID, o.TaskID = e.ActionID, e.ID, e.TaskID
	}
	o.Simulated = true
	b.logf("[%s] command %s %s outcome=%s code=%s %s", b.cfg.Source, o.CommandID, o.Verb, o.Outcome, o.Code, o.Message)
	b.publishEvent(b.outcomeSubject(), "evt.drone.command", "", e, o)
	if e != nil && e.TaskID != "" {
		b.reportTask(e, o)
	}
}

// nextTaskVersion returns the next status_version for taskID: it increases by one per
// report across every command of the task, for the life of this process. The task
// service still drops a report that is not newer than its record (mission-task.md §5
// rule 1); see the owner proposal P7 for the restart and base-version gap.
func (b *Bridge) nextTaskVersion(taskID string) int {
	b.taskMu.Lock()
	defer b.taskMu.Unlock()
	v, ok := b.taskVer[taskID]
	if !ok {
		b.taskIDs = append(b.taskIDs, taskID)
		if len(b.taskIDs) > maxTrackedTasks {
			delete(b.taskVer, b.taskIDs[0])
			b.taskIDs = b.taskIDs[1:]
		}
	}
	v++
	b.taskVer[taskID] = v
	return v
}

// reportTask publishes evt.drone.{id}.task (mission-task.md §7) for task-bearing commands.
func (b *Bridge) reportTask(e *command.Envelope, o outcome) {
	state := ""
	var errObj any
	switch o.Outcome {
	case "acknowledged":
		state = "received"
	case "executed":
		// The vehicle accepted the command. Movement verbs keep executing; instant verbs are done.
		state = "executing"
		if o.Verb == "arm" || o.Verb == "disarm" || o.Verb == "estop" {
			state = "done_ok"
		}
	case "refused":
		state = "done_not_ok"
		errObj = map[string]string{"code": "rejected", "message": o.Code + ": " + o.Message}
	case "failed":
		state = "done_not_ok"
		errObj = map[string]string{"code": "failed", "message": o.Code + ": " + o.Message}
	default:
		return
	}
	ver := b.nextTaskVersion(e.TaskID)
	// Nats-Msg-Id = {task_id}:{status_version} (mission-task.md §7).
	b.publishEvent(b.taskSubject(), "evt.drone.task", fmt.Sprintf("%s:%d", e.TaskID, ver), e, map[string]any{
		"task_id": e.TaskID, "state": state, "status_version": ver, "error": errObj,
		"at_t": b.clock().UnixMilli(), "simulated": true,
	})
}

// HandleCommand verifies one command from cmd.fleet.drone.{id}.{verb} and, if every
// check passes, sends it to the vehicle (asset-command.md §4 step 6).
func (b *Bridge) HandleCommand(ctx context.Context, subject string, data []byte) {
	verbName := subject[len("cmd.fleet.drone."+b.cfg.AssetID+"."):]
	e, err := command.Parse(data)
	if err != nil {
		b.report(nil, outcome{Outcome: "refused", Verb: verbName, Code: command.CodeSchemaInvalid, Message: "envelope: " + err.Error()})
		return
	}
	now := b.clock()
	res, ref := b.verifier.Verify(subject, e, now)
	if ref != nil {
		b.report(e, outcome{Outcome: "refused", Verb: verbName, Code: ref.Code, Detail: ref.Detail, Message: ref.Message})
		return
	}
	if !res.Signed {
		b.logf("[%s] WARNING: command %s accepted without a signature (SIG_REQUIRED=false)", b.cfg.Source, e.ID)
		if res.TaskBound {
			b.logf("[%s] WARNING: command %s is task-bound (%s): its params_hash is the task action's and nothing binds data to the approval without a signature", b.cfg.Source, e.ID, e.TaskID)
		}
	}
	sigOK := res.Signed

	req, ref := b.precheck(res, now)
	if ref != nil {
		b.report(e, outcome{Outcome: "refused", Verb: verbName, Code: ref.Code, Detail: ref.Detail, Message: ref.Message, SigOK: &sigOK})
		return
	}
	if e.DryRun {
		b.report(e, outcome{Outcome: "acknowledged", Verb: verbName, Message: "dry_run: verified, not sent to the vehicle", SigOK: &sigOK})
		return
	}
	b.report(e, outcome{Outcome: "acknowledged", Verb: verbName, SigOK: &sigOK})

	if verbName != "estop" {
		b.moveMu.Lock()
		defer b.moveMu.Unlock()
	}
	exp, _ := time.Parse(time.RFC3339Nano, e.ExpiresAt)
	if verbName != "estop" && !b.clock().Before(exp) {
		b.report(e, outcome{Outcome: "failed", Verb: verbName, Code: command.CodeExpired, Message: "expired while waiting for the previous command"})
		return
	}
	result, err := b.veh.Send(ctx, req)
	mav := map[string]any{"command": req.Name(), "id": int(req.Cmd)}
	switch {
	case errors.Is(err, ErrNoAck):
		b.report(e, outcome{Outcome: "failed", Verb: verbName, Code: command.CodeNoAck, Message: err.Error(), Mavlink: mav})
	case errors.Is(err, ErrNotHeard):
		// Nothing was sent: the link has not delivered a frame from the vehicle yet.
		b.report(e, outcome{Outcome: "failed", Verb: verbName, Code: command.CodePrecondition, Detail: "offline", Message: err.Error(), Mavlink: mav})
	case errors.Is(err, ErrBusy):
		b.report(e, outcome{Outcome: "failed", Verb: verbName, Code: command.CodeConflict, Message: err.Error(), Mavlink: mav})
	case err != nil:
		// Link write error or shutdown before an ACK: the vehicle never answered
		// (owner proposal P4). vehicle_rejected is kept for a MAVLink result other than 0.
		detail := "not_sent"
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			detail = "aborted"
		}
		b.report(e, outcome{Outcome: "failed", Verb: verbName, Code: command.CodeNoAck, Detail: detail, Message: err.Error(), Mavlink: mav})
	case result != int(common.MAV_RESULT_ACCEPTED):
		mav["result"] = result
		mav["result_name"] = common.MAV_RESULT(result).String()
		b.report(e, outcome{Outcome: "failed", Verb: verbName, Code: command.CodeVehicleRejected,
			Message: "vehicle answered " + common.MAV_RESULT(result).String(), Mavlink: mav})
	default:
		mav["result"] = 0
		b.report(e, outcome{Outcome: "executed", Verb: verbName, Mavlink: mav})
	}
}

func precond(detail, format string, a ...any) *command.Refusal {
	return &command.Refusal{Code: command.CodePrecondition, Detail: detail, Message: fmt.Sprintf(format, a...)}
}

// precheck runs the vehicle-state checks (link, geofence) and builds the MAVLink request.
// estop is never blocked by a precondition (asset-command.md §3).
func (b *Bridge) precheck(res *command.Result, now time.Time) (MavRequest, *command.Refusal) {
	p := res.Params
	verb := res.Verb.Name
	if verb != "estop" && !b.state.online(now, b.cfg.LinkTimeout) {
		return MavRequest{}, precond("offline", "no HEARTBEAT from sysid %d within %s", b.cfg.SysID, b.cfg.LinkTimeout)
	}
	checkHere := func() *command.Refusal {
		lat, lon, _, ok := b.state.position()
		if !ok {
			return precond("keep_in", "vehicle position unknown")
		}
		if v := b.cfg.Fence.Check(fence.Point{Lat: lat, Lon: lon}, nil); v != nil {
			return precond(v.Detail, "vehicle position: %s", v.Reason)
		}
		return nil
	}
	switch verb {
	case "goto":
		if v := b.cfg.Fence.Check(fence.Point{Lat: *p.Lat, Lon: *p.Lon}, p.HaeM); v != nil {
			return MavRequest{}, precond(v.Detail, "goto target: %s", v.Reason)
		}
		msl, err := geoid.HAEToMSL(*p.Lat, *p.Lon, *p.HaeM)
		if err != nil {
			return MavRequest{}, &command.Refusal{Code: command.CodeSchemaInvalid, Message: err.Error()}
		}
		speed, yaw := float32(-1), nanF
		if p.SpeedMps != nil {
			speed = float32(*p.SpeedMps)
		}
		if p.YawDeg != nil {
			yaw = float32(*p.YawDeg)
		}
		// MAV_CMD_DO_REPOSITION over COMMAND_INT, frame GLOBAL (altitude MSL):
		// param1 speed (-1 default), param2 MAV_DO_REPOSITION_FLAGS_CHANGE_MODE,
		// param4 yaw (NaN keeps heading), x/y degE7, z altitude.
		// LongFallback: re-sent as COMMAND_LONG (params 5/6/7) if COMMAND_INT is unsupported.
		return MavRequest{Cmd: common.MAV_CMD_DO_REPOSITION, Int: true, LongFallback: true, Frame: common.MAV_FRAME_GLOBAL,
			P: [7]float32{speed, float32(common.MAV_DO_REPOSITION_FLAGS_CHANGE_MODE), 0, yaw},
			X: int32(math.Round(*p.Lat * 1e7)), Y: int32(math.Round(*p.Lon * 1e7)), Z: float32(msl)}, nil
	case "takeoff":
		if r := checkHere(); r != nil {
			return MavRequest{}, r
		}
		lat, lon, _, _ := b.state.position()
		if v := b.cfg.Fence.Check(fence.Point{Lat: lat, Lon: lon}, p.HaeM); v != nil {
			return MavRequest{}, precond(v.Detail, "takeoff altitude: %s", v.Reason)
		}
		msl, _ := geoid.HAEToMSL(lat, lon, *p.HaeM)
		// MAV_CMD_NAV_TAKEOFF param7 altitude (MSL in the global frame; per-autopilot
		// interpretation is unverified, see README).
		return MavRequest{Cmd: common.MAV_CMD_NAV_TAKEOFF, P: [7]float32{0, 0, 0, nanF, nanF, nanF, float32(msl)}}, nil
	case "arm":
		if r := checkHere(); r != nil {
			return MavRequest{}, r
		}
		return MavRequest{Cmd: common.MAV_CMD_COMPONENT_ARM_DISARM, P: [7]float32{1}}, nil
	case "disarm":
		return MavRequest{Cmd: common.MAV_CMD_COMPONENT_ARM_DISARM, P: [7]float32{0}}, nil
	case "rtl":
		return MavRequest{Cmd: common.MAV_CMD_NAV_RETURN_TO_LAUNCH}, nil
	case "land":
		if p.Lat != nil {
			if v := b.cfg.Fence.Check(fence.Point{Lat: *p.Lat, Lon: *p.Lon}, nil); v != nil {
				return MavRequest{}, precond(v.Detail, "land point: %s", v.Reason)
			}
			// NAV_LAND at a point: COMMAND_INT keeps full lat/lon precision; z NaN = current altitude reference.
			return MavRequest{Cmd: common.MAV_CMD_NAV_LAND, Int: true, LongFallback: true, Frame: common.MAV_FRAME_GLOBAL,
				P: [7]float32{0, 0, 0, nanF}, X: int32(math.Round(*p.Lat * 1e7)), Y: int32(math.Round(*p.Lon * 1e7)), Z: 0}, nil
		}
		return MavRequest{Cmd: common.MAV_CMD_NAV_LAND, P: [7]float32{0, 0, 0, nanF, 0, 0, 0}}, nil
	case "estop":
		// MAV_CMD_DO_FLIGHTTERMINATION param1 >= 0.5 terminates (motors stop).
		return MavRequest{Cmd: common.MAV_CMD_DO_FLIGHTTERMINATION, P: [7]float32{1}}, nil
	}
	return MavRequest{}, &command.Refusal{Code: command.CodeUnknownAction, Message: "no MAVLink mapping for " + verb}
}

func (b *Bridge) clock() time.Time { return time.Now().Add(time.Duration(b.offset.Load())) }
