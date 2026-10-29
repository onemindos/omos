package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

// MavRequest is one MAVLink command for the vehicle. Int selects COMMAND_INT (with
// X/Y as degE7 and Z as altitude in Frame); otherwise COMMAND_LONG with P[0..6].
// LongFallback lets an Int request be re-sent as COMMAND_LONG (params 5/6/7 = lat, lon
// in degrees and Z) when the vehicle answers MAV_RESULT_UNSUPPORTED or never ACKs
// COMMAND_INT, as vehicles that only read COMMAND_LONG do.
type MavRequest struct {
	Cmd          common.MAV_CMD
	Int          bool
	LongFallback bool
	Frame        common.MAV_FRAME
	P            [7]float32
	X, Y         int32
	Z            float32
}

// asLong is the COMMAND_LONG form of an Int request. float32 degrees keep about 0.5 m
// of precision, which is why COMMAND_INT is tried first.
func (r MavRequest) asLong() MavRequest {
	l := r
	l.Int, l.LongFallback = false, false
	l.P[4], l.P[5], l.P[6] = float32(float64(r.X)/1e7), float32(float64(r.Y)/1e7), r.Z
	return l
}

// Name is the MAV_CMD enum name, e.g. MAV_CMD_DO_REPOSITION.
func (r MavRequest) Name() string { return r.Cmd.String() }

// Send errors. Only ErrNoAck means the command reached the link and went unanswered;
// ErrNotHeard and ErrBusy mean nothing was sent.
var (
	ErrNoAck    = errors.New("no COMMAND_ACK from vehicle")
	ErrNotHeard = errors.New("vehicle not heard yet")
	ErrBusy     = errors.New("same MAV_CMD already in flight")
)

// Vehicle sends a command and returns the MAV_RESULT of the final COMMAND_ACK.
type Vehicle interface {
	Send(ctx context.Context, r MavRequest) (int, error)
}

// mavVehicle is a Vehicle over a gomavlib node.
type mavVehicle struct {
	node       *gomavlib.Node
	sysID      uint8
	gcsSysID   uint8
	ackTimeout time.Duration
	retries    int
	logf       func(string, ...any)

	mu       sync.Mutex
	channel  *gomavlib.Channel
	waiters  map[common.MAV_CMD]chan uint8
	longOnly map[common.MAV_CMD]bool // COMMAND_INT unsupported for this command
}

func newMavVehicle(node *gomavlib.Node, sysID, gcsSysID uint8, ackTimeout time.Duration, retries int) *mavVehicle {
	return &mavVehicle{node: node, sysID: sysID, gcsSysID: gcsSysID, ackTimeout: ackTimeout, retries: retries,
		logf: func(string, ...any) {}, waiters: map[common.MAV_CMD]chan uint8{}, longOnly: map[common.MAV_CMD]bool{}}
}

// onFrame records the vehicle's channel and routes COMMAND_ACKs to waiters. An ACK
// addressed to another system (another GCS on a shared link) is ignored; target 0 is
// accepted because MAVLink 1 ACKs and many autopilots leave the target unset.
// It returns the message when it came from the configured vehicle, else nil.
func (v *mavVehicle) onFrame(f *gomavlib.EventFrame) message.Message {
	if f.SystemID() != v.sysID {
		return nil
	}
	v.mu.Lock()
	v.channel = f.Channel
	if ack, ok := f.Message().(*ardupilotmega.MessageCommandAck); ok &&
		(ack.TargetSystem == 0 || ack.TargetSystem == v.gcsSysID) {
		if ch, ok := v.waiters[ack.Command]; ok {
			select {
			case ch <- uint8(ack.Result):
			default:
			}
		}
	}
	v.mu.Unlock()
	return f.Message()
}

func (v *mavVehicle) write(m message.Message) error {
	v.mu.Lock()
	ch := v.channel
	v.mu.Unlock()
	if ch == nil {
		return fmt.Errorf("%w: sysid %d", ErrNotHeard, v.sysID)
	}
	return v.node.WriteMessageTo(ch, m)
}

func (v *mavVehicle) Send(ctx context.Context, r MavRequest) (int, error) {
	ackCh := make(chan uint8, 4)
	v.mu.Lock()
	if _, busy := v.waiters[r.Cmd]; busy {
		v.mu.Unlock()
		return -1, fmt.Errorf("%w: %s", ErrBusy, r.Name())
	}
	v.waiters[r.Cmd] = ackCh
	longOnly := v.longOnly[r.Cmd]
	v.mu.Unlock()
	defer func() {
		v.mu.Lock()
		delete(v.waiters, r.Cmd)
		v.mu.Unlock()
	}()

	if r.Int && r.LongFallback {
		if !longOnly {
			res, err := v.send(ctx, ackCh, r)
			if !(errors.Is(err, ErrNoAck) || (err == nil && res == int(common.MAV_RESULT_UNSUPPORTED))) {
				return res, err
			}
			v.logf("[mavlink] %s: COMMAND_INT unsupported or unanswered (result %d, err %v); retrying as COMMAND_LONG", r.Name(), res, err)
		}
		res, err := v.send(ctx, ackCh, r.asLong())
		if err == nil && res != int(common.MAV_RESULT_UNSUPPORTED) {
			v.mu.Lock()
			v.longOnly[r.Cmd] = true
			v.mu.Unlock()
		}
		return res, err
	}
	return v.send(ctx, ackCh, r)
}

// send transmits r with retries and waits for its final COMMAND_ACK on ackCh.
func (v *mavVehicle) send(ctx context.Context, ackCh chan uint8, r MavRequest) (int, error) {
	for { // drop a late ACK left over from a previous form of this command
		select {
		case <-ackCh:
			continue
		default:
		}
		break
	}

	inProgress := false
	for attempt := 0; attempt < v.retries; attempt++ {
		var m message.Message
		if r.Int {
			m = &ardupilotmega.MessageCommandInt{
				TargetSystem: v.sysID, TargetComponent: 1, Frame: r.Frame, Command: r.Cmd,
				Param1: r.P[0], Param2: r.P[1], Param3: r.P[2], Param4: r.P[3], X: r.X, Y: r.Y, Z: r.Z,
			}
		} else {
			m = &ardupilotmega.MessageCommandLong{
				TargetSystem: v.sysID, TargetComponent: 1, Command: r.Cmd, Confirmation: uint8(attempt),
				Param1: r.P[0], Param2: r.P[1], Param3: r.P[2], Param4: r.P[3], Param5: r.P[4], Param6: r.P[5], Param7: r.P[6],
			}
		}
		if !inProgress {
			if err := v.write(m); err != nil {
				return -1, err
			}
		}
		timer := time.NewTimer(v.ackTimeout)
	wait:
		for {
			select {
			case <-ctx.Done():
				timer.Stop()
				return -1, ctx.Err()
			case res := <-ackCh:
				if res == uint8(common.MAV_RESULT_IN_PROGRESS) {
					inProgress = true // keep waiting for the final ACK
					continue
				}
				timer.Stop()
				return int(res), nil
			case <-timer.C:
				break wait
			}
		}
	}
	return -1, ErrNoAck
}

// nanF is a float32 NaN, MAVLink's "no change / use default" for many params.
var nanF = float32(math.NaN())
