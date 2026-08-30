package main

import "math"

// ── 3D vector math (local ENU frame, meters) ────────────────────────────────
// X = east, Y = north, Z = up.

type vec3 struct{ X, Y, Z float64 }

func (a vec3) add(b vec3) vec3   { return vec3{a.X + b.X, a.Y + b.Y, a.Z + b.Z} }
func (a vec3) sub(b vec3) vec3   { return vec3{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }
func (a vec3) scale(s float64) vec3 { return vec3{a.X * s, a.Y * s, a.Z * s} }
func (a vec3) dot(b vec3) float64   { return a.X*b.X + a.Y*b.Y + a.Z*b.Z }
func (a vec3) cross(b vec3) vec3 {
	return vec3{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X}
}
func (a vec3) norm() float64 { return math.Sqrt(a.dot(a)) }
func (a vec3) unit() vec3 {
	n := a.norm()
	if n < 1e-9 {
		return vec3{}
	}
	return a.scale(1.0 / n)
}
func clampVec(v vec3, max float64) vec3 {
	n := v.norm()
	if n > max && n > 1e-9 {
		return v.scale(max / n)
	}
	return v
}

// ── Geodesy: local tangent-plane (equirectangular) around a fixed origin ─────
// Good to <1% over the few-km engagement envelope; matches the spine's haversine
// scale. Not for global nav — only relative interceptor↔target geometry.

const (
	mPerDegLat = 111320.0
	deg2rad    = math.Pi / 180.0
	rad2deg    = 180.0 / math.Pi
)

type geoOrigin struct{ lat0, lon0, alt0 float64 }

func (o geoOrigin) toLocal(lat, lon, alt float64) vec3 {
	mPerDegLon := mPerDegLat * math.Cos(o.lat0*deg2rad)
	return vec3{
		X: (lon - o.lon0) * mPerDegLon,
		Y: (lat - o.lat0) * mPerDegLat,
		Z: alt - o.alt0,
	}
}

func (o geoOrigin) toGeo(v vec3) (lat, lon, alt float64) {
	mPerDegLon := mPerDegLat * math.Cos(o.lat0*deg2rad)
	lat = o.lat0 + v.Y/mPerDegLat
	lon = o.lon0 + v.X/mPerDegLon
	alt = o.alt0 + v.Z
	return
}

// headingDeg returns compass heading (deg true) of a horizontal velocity vector.
func headingDeg(v vec3) float64 {
	h := math.Atan2(v.X, v.Y) * rad2deg // atan2(east, north)
	if h < 0 {
		h += 360
	}
	return h
}

// ── Guidance parameters (a small interceptor drone / loitering interceptor) ──

type guidanceParams struct {
	Vmax        float64 // max airspeed (m/s)
	Amax        float64 // max lateral accel (m/s²)
	SeekerRange float64 // range at which the terminal seeker acquires (m)
	LockDelay   float64 // seeker search→lock dwell (s)
	HitRadius   float64 // warhead lethal radius (m)
	FragRadius  float64 // near-miss still-a-kill radius at CPA (m)
	MaxFlight   float64 // flyout time budget (s)
	NavGain     float64 // proportional-navigation constant N (typ. 3–5)
}

func defaultParams() guidanceParams {
	return guidanceParams{
		Vmax: 60, Amax: 40, SeekerRange: 800, LockDelay: 1.0,
		HitRadius: 6, FragRadius: 15, MaxFlight: 120, NavGain: 4,
	}
}

// midcourseSteer — inertial/GPS-INS midcourse: pure-pursuit toward the target's
// current (dead-reckoned) position, accelerating to Vmax under an accel limit.
func midcourseSteer(intVel, R vec3, p guidanceParams, dt float64) vec3 {
	desired := R.unit().scale(p.Vmax)
	dv := clampVec(desired.sub(intVel), p.Amax*dt)
	return intVel.add(dv)
}

// terminalPN — terminal homing via true proportional navigation. The commanded
// lateral acceleration is a = N · Vc · (Ω × R̂), where Ω is the line-of-sight
// rotation rate. Speed is held at Vmax; PN only rotates the velocity vector to
// drive LOS rate → 0 (collision course). This is the guidance SentinelX left as
// a TODO — completed here.
func terminalPN(intVel, R, Vrel vec3, p guidanceParams, dt float64) vec3 {
	R2 := R.dot(R)
	if R2 < 1e-6 {
		return intVel
	}
	Rmag := math.Sqrt(R2)
	Vc := -R.dot(Vrel) / Rmag           // closing speed (+ = closing)
	omega := R.cross(Vrel).scale(1 / R2) // LOS rotation rate vector
	aCmd := clampVec(omega.cross(R.unit()).scale(p.NavGain*Vc), p.Amax)
	v := intVel.add(aCmd.scale(dt))
	// Hold airspeed at Vmax (thrust trims speed; PN only steers).
	if s := v.norm(); s > 1e-6 {
		v = v.scale(p.Vmax / s)
	}
	return v
}
