package main

import (
	"math"
	"sync"
	"time"
)

// Phase — the interceptor state machine, adopted verbatim from SentinelX's
// InterceptorPhase.msg so a future DDS↔NATS bridge maps 1:1.
type Phase uint8

const (
	PhaseIdle             Phase = 0
	PhaseMissionLoaded    Phase = 1
	PhaseReadyToLaunch    Phase = 2
	PhaseLaunched         Phase = 3
	PhaseInertialMidcourse Phase = 4
	PhaseSeekerSearch     Phase = 5
	PhaseSeekerLocked     Phase = 6
	PhaseTerminalHoming   Phase = 7
	PhaseInterceptSuccess Phase = 8
	PhaseInterceptFailed  Phase = 9
	PhaseReturnHome       Phase = 10
	PhaseFault            Phase = 11
)

func (p Phase) String() string {
	switch p {
	case PhaseIdle:
		return "IDLE"
	case PhaseMissionLoaded:
		return "MISSION_LOADED"
	case PhaseReadyToLaunch:
		return "READY_TO_LAUNCH"
	case PhaseLaunched:
		return "LAUNCHED"
	case PhaseInertialMidcourse:
		return "INERTIAL_MIDCOURSE"
	case PhaseSeekerSearch:
		return "SEEKER_SEARCH"
	case PhaseSeekerLocked:
		return "SEEKER_LOCKED"
	case PhaseTerminalHoming:
		return "TERMINAL_HOMING"
	case PhaseInterceptSuccess:
		return "INTERCEPT_SUCCESS"
	case PhaseInterceptFailed:
		return "INTERCEPT_FAILED"
	case PhaseReturnHome:
		return "RETURN_HOME"
	case PhaseFault:
		return "FAULT"
	}
	return "UNKNOWN"
}

// Authority — who currently owns guidance (SentinelX authority enum). Handoff
// runs prelaunch-C2 → inertial (midcourse) → seeker (terminal).
type Authority uint8

const (
	AuthPrelaunchC2 Authority = 0
	AuthInertial    Authority = 1
	AuthSeeker      Authority = 2
	AuthSafeMode    Authority = 3
)

// terminal reports whether a phase is an absorbing end-state.
func (p Phase) terminal() bool {
	return p == PhaseInterceptSuccess || p == PhaseInterceptFailed || p == PhaseIdle
}

// Engagement is one interceptor prosecuting one target. All spatial state is in
// the local ENU frame anchored at the launch site (origin).
type Engagement struct {
	mu sync.Mutex

	ID        string
	MissionID string
	TargetID  string

	Phase     Phase
	Authority Authority
	Launched  bool
	SeekerDet bool
	SeekerLck bool

	origin geoOrigin
	params guidanceParams

	// Interceptor kinematics (local ENU, m and m/s).
	pos vec3
	vel vec3

	// Target estimate (local ENU). Dead-reckoned between C2/spine updates.
	tgtPos       vec3
	tgtVel       vec3
	haveTgt      bool
	lastTgtStamp time.Time

	// Metrics.
	launchTime  time.Time
	searchStart time.Time
	minRange    float64
	prevRange   float64
	passCount   int // terminal re-attacks used (drone interceptor can come around)
	health      string
}

// setOrigin overrides the launch-site origin (per-mission defense site). Must be
// called before assignTarget so target coords convert against the right frame.
func (e *Engagement) setOrigin(o geoOrigin) {
	e.mu.Lock()
	e.origin = o
	e.mu.Unlock()
}

func newEngagement(id string, origin geoOrigin, p guidanceParams) *Engagement {
	return &Engagement{
		ID: id, origin: origin, params: p,
		Phase: PhaseIdle, Authority: AuthPrelaunchC2,
		minRange: math.MaxFloat64, health: "nominal",
	}
}

// assignTarget — C2 ASSIGN_TARGET. Loads a mission + target and seeds the target
// estimate (fire-and-forget: the last C2 track is dead-reckoned if the spine
// stops updating). Rejected after launch (per SentinelX).
func (e *Engagement) assignTarget(missionID, targetID string, lat, lon, alt, velN, velE, velUp float64) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.Launched {
		return false
	}
	e.MissionID = missionID
	e.TargetID = targetID
	e.tgtPos = e.origin.toLocal(lat, lon, alt)
	e.tgtVel = vec3{X: velE, Y: velN, Z: velUp}
	e.haveTgt = true
	e.lastTgtStamp = time.Now()
	e.Phase = PhaseMissionLoaded
	e.Authority = AuthPrelaunchC2
	return true
}

// startIntercept — C2 START_INTERCEPT. Launches (boosts toward the current LOS)
// and hands guidance to inertial midcourse. Rejected if no target or already up.
func (e *Engagement) startIntercept() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.Launched || !e.haveTgt || e.TargetID == "" {
		return false
	}
	e.Launched = true
	e.pos = vec3{Z: 0}
	// Initial boost along the launch line-of-sight at 30% Vmax.
	los := e.tgtPos.sub(e.pos).unit()
	e.vel = los.scale(0.3 * e.params.Vmax)
	e.Phase = PhaseInertialMidcourse
	e.Authority = AuthInertial
	e.launchTime = time.Now()
	return true
}

// abort — C2 ABORT/HIT. Post-launch → return home; pre-launch → idle.
func (e *Engagement) abort() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.Launched {
		e.Phase = PhaseReturnHome
		e.Authority = AuthSafeMode
	} else {
		e.Phase = PhaseIdle
		e.MissionID, e.TargetID = "", ""
	}
	e.Launched = false
	e.SeekerDet, e.SeekerLck = false, false
}

// updateTarget — fresh spine track (ent.air.{target}.state) refines the estimate
// and re-derives target velocity from the position delta.
func (e *Engagement) updateTarget(lat, lon, alt float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	np := e.origin.toLocal(lat, lon, alt)
	now := time.Now()
	if e.haveTgt {
		if dt := now.Sub(e.lastTgtStamp).Seconds(); dt > 0.02 && dt < 5 {
			// Finite-difference velocity from irregular (~1 Hz) spine updates is
			// noisy; low-pass it so the terminal lead stays smooth. The C2-seeded
			// velocity anchors the estimate until observations converge.
			measured := np.sub(e.tgtPos).scale(1 / dt)
			e.tgtVel = e.tgtVel.scale(0.6).add(measured.scale(0.4))
		}
	}
	e.tgtPos = np
	e.haveTgt = true
	e.lastTgtStamp = now
}

// step advances the sim/guidance by dt seconds. Returns true when the engagement
// has reached an absorbing end-state this tick.
func (e *Engagement) step(dt float64) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.Launched || e.Phase.terminal() {
		return false
	}

	// Dead-reckon the target forward (fire-and-forget between C2 updates).
	e.tgtPos = e.tgtPos.add(e.tgtVel.scale(dt))

	// Return-home: fly back toward the origin, land when close.
	if e.Phase == PhaseReturnHome {
		home := vec3{}.sub(e.pos)
		if home.norm() < 5 {
			e.Phase = PhaseIdle
			return true
		}
		e.vel = midcourseSteer(e.vel, home, e.params, dt)
		e.pos = e.pos.add(e.vel.scale(dt))
		return false
	}

	R := e.tgtPos.sub(e.pos)
	Rmag := R.norm()
	Vrel := e.tgtVel.sub(e.vel)

	// Phase progression.
	switch e.Phase {
	case PhaseInertialMidcourse:
		if Rmag < e.params.SeekerRange {
			e.Phase = PhaseSeekerSearch
			e.Authority = AuthSeeker
			e.searchStart = time.Now()
		}
	case PhaseSeekerSearch:
		e.SeekerDet = true
		if time.Since(e.searchStart).Seconds() >= e.params.LockDelay {
			e.SeekerLck = true
			e.Phase = PhaseTerminalHoming
		}
	}

	// Guidance law by authority.
	if e.Phase == PhaseTerminalHoming {
		e.vel = terminalPN(e.vel, R, Vrel, e.params, dt)
	} else {
		e.vel = midcourseSteer(e.vel, R, e.params, dt)
	}
	e.pos = e.pos.add(e.vel.scale(dt))

	// Intercept adjudication.
	if Rmag < e.minRange {
		e.minRange = Rmag
	}
	if Rmag <= e.params.HitRadius {
		e.Phase = PhaseInterceptSuccess
		return true
	}
	// Closest-point-of-approach passed: range is now opening during terminal
	// homing (hysteresis of 1 m rejects noise). Adjudicate the pass.
	if e.Phase == PhaseTerminalHoming && Rmag > e.prevRange+1.0 {
		if e.minRange <= e.params.FragRadius {
			e.Phase = PhaseInterceptSuccess // warhead frag envelope — kill
			return true
		}
		// Missed this pass. A drone interceptor can come around and re-attack;
		// re-enter midcourse to re-close and re-acquire, up to a limit.
		e.passCount++
		if e.passCount >= 3 {
			e.Phase = PhaseInterceptFailed
			e.health = "exhausted re-attacks"
			return true
		}
		e.Phase = PhaseInertialMidcourse
		e.Authority = AuthInertial
		e.SeekerDet = false
		e.SeekerLck = false
		e.minRange = math.MaxFloat64
	}
	if time.Since(e.launchTime).Seconds() > e.params.MaxFlight {
		e.Phase = PhaseInterceptFailed
		e.health = "flyout timeout"
		return true
	}
	e.prevRange = Rmag
	return false
}

// snapshot captures the current state for telemetry publication (thread-safe).
type snapshot struct {
	Phase                Phase
	Authority            Authority
	Launched             bool
	SeekerDet, SeekerLck bool
	Lat, Lon, Alt        float64
	Heading, Speed       float64
	Range, Closing       float64
	MinRange             float64
	Health               string
}

func (e *Engagement) snapshot() snapshot {
	e.mu.Lock()
	defer e.mu.Unlock()
	lat, lon, alt := e.origin.toGeo(e.pos)
	R := e.tgtPos.sub(e.pos)
	Rmag := R.norm()
	var closing float64
	if Rmag > 1e-6 {
		closing = -R.dot(e.tgtVel.sub(e.vel)) / Rmag
	}
	mr := e.minRange
	if mr == math.MaxFloat64 {
		mr = Rmag
	}
	return snapshot{
		Phase: e.Phase, Authority: e.Authority, Launched: e.Launched,
		SeekerDet: e.SeekerDet, SeekerLck: e.SeekerLck,
		Lat: lat, Lon: lon, Alt: alt,
		Heading: headingDeg(e.vel), Speed: e.vel.norm(),
		Range: Rmag, Closing: closing, MinRange: mr, Health: e.health,
	}
}
