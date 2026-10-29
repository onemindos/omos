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
type MavRequest struct {
	Cmd   common.MAV_CMD
	Int   bool
	Frame common.MAV_FRAME
	P     [7]float32
	X, Y  int32
	Z     float32
}

// Name is the MAV_CMD enum name, e.g. MAV_CMD_DO_REPOSITION.
func (r MavRequest) Name() string { return r.Cmd.String() }

// ErrNoAck means the vehicle sent no COMMAND_ACK after every retry.
var ErrNoAck = errors.New("no COMMAND_ACK from vehicle")

// Vehicle sends a command and returns the MAV_RESULT of the final COMMAND_ACK.
type Vehicle interface {
	Send(ctx context.Context, r MavRequest) (int, error)
}

// mavVehicle is a Vehicle over a gomavlib node.
type mavVehicle struct {
	node       *gomavlib.Node
	sysID      uint8
	ackTimeout time.Duration
	retries    int

	mu      sync.Mutex
	channel *gomavlib.Channel
	waiters map[common.MAV_CMD]chan uint8
}

func newMavVehicle(node *gomavlib.Node, sysID uint8, ackTimeout time.Duration, retries int) *mavVehicle {
	return &mavVehicle{node: node, sysID: sysID, ackTimeout: ackTimeout, retries: retries,
		waiters: map[common.MAV_CMD]chan uint8{}}
}

// onFrame records the vehicle's channel and routes COMMAND_ACKs to waiters.
// It returns the message when it came from the configured vehicle, else nil.
func (v *mavVehicle) onFrame(f *gomavlib.EventFrame) message.Message {
	if f.SystemID() != v.sysID {
		return nil
	}
	v.mu.Lock()
	v.channel = f.Channel
	if ack, ok := f.Message().(*ardupilotmega.MessageCommandAck); ok {
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
		return fmt.Errorf("vehicle sysid %d not heard yet", v.sysID)
	}
	return v.node.WriteMessageTo(ch, m)
}

func (v *mavVehicle) Send(ctx context.Context, r MavRequest) (int, error) {
	ackCh := make(chan uint8, 4)
	v.mu.Lock()
	if _, busy := v.waiters[r.Cmd]; busy {
		v.mu.Unlock()
		return -1, fmt.Errorf("%s already in flight", r.Name())
	}
	v.waiters[r.Cmd] = ackCh
	v.mu.Unlock()
	defer func() {
		v.mu.Lock()
		delete(v.waiters, r.Cmd)
		v.mu.Unlock()
	}()

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
