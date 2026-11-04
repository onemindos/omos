// Package fusion holds the per-entity estimator (fusion-and-tracking.md Stage 2): a
// constant-velocity Kalman filter in a local East-North-Up plane, re-anchored at every update,
// with altitude kept separately. The Estimator and Associator interfaces are the hooks for the
// later IMM (CV + coordinated turn) estimator and the GNN/JPDA associator for keyless `det.*`
// detections; neither is implemented in this first slice (Recommendation, "First slice").
package fusion

import (
	"errors"
	"math"
)

const earthR = 6378137.0 // WGS-84 semi-major axis, metres

// Measurement is one position (and optional velocity) report.
type Measurement struct {
	T      int64   // epoch ms
	Lat    float64 // degrees
	Lon    float64
	SigmaM float64 // 1-sigma horizontal position error, metres
	HasVel bool
	VE, VN float64 // m/s east, north
	SigmaV float64 // 1-sigma velocity error, m/s
}

// State is the estimate at time T.
type State struct {
	T      int64
	Lat    float64
	Lon    float64
	VE, VN float64
	PosCov [3]float64 // ee, en, nn (m^2)
	VelCov [3]float64 // ee, en, nn ((m/s)^2)
}

// CEM is the 1-sigma circular error implied by the position covariance.
func (s State) CEM() float64 { return math.Sqrt(math.Max(s.PosCov[0], s.PosCov[2])) }

// SpeedCourse returns ground speed (m/s) and course (degrees true).
func (s State) SpeedCourse() (float64, float64) {
	sp := math.Hypot(s.VE, s.VN)
	c := math.Mod(math.Atan2(s.VE, s.VN)*180/math.Pi+360, 360)
	return sp, c
}

// Estimator is the per-entity filter. An IMM estimator implements the same interface.
type Estimator interface {
	// Update predicts to m.T and applies m. A measurement older than the filter time returns
	// ErrOutOfOrder and leaves the state unchanged.
	Update(m Measurement) error
	// Predict returns the state extrapolated to t without changing the filter.
	Predict(t int64) State
	State() State
}

// Associator is the hook for keyless detections (Stage 3): global nearest neighbour with
// chi-square gating and one Hungarian assignment per scan, or JPDA for a cluttered source.
// It returns, for each detection, the index of the assigned track or -1 for a new tentative
// track. Not implemented in this slice: no stream captures `det.>` yet (DET is proposed).
type Associator interface {
	Associate(tracks []State, dets []Measurement) []int
}

// ErrOutOfOrder marks a measurement older than the filter state.
var ErrOutOfOrder = errors.New("measurement older than filter state")

// CV is a constant-velocity Kalman filter. x = [e, n, ve, vn] relative to (refLat, refLon).
type CV struct {
	q      float64 // white-noise acceleration spectral density sigma (m/s^2)
	t      int64
	refLat float64
	refLon float64
	x      [4]float64
	p      [4][4]float64
	init   bool
}

// NewCV returns a filter with process-noise acceleration sigma q (m/s^2), by class.
func NewCV(q float64) *CV { return &CV{q: q} }

// ProcessNoise returns a starting acceleration sigma per entity_type (ship << aircraft <<
// small drone). These are starting points to tune on fixtures, not measured values
// (fusion-and-tracking.md "Unverified"). Satellites get a very large value so the filter follows
// the SGP4 positions the feed already computes instead of lagging an orbit.
func ProcessNoise(entityType, kind string) float64 {
	if kind == "satellite" {
		return 1000
	}
	switch entityType {
	case "vessel":
		return 0.05
	case "aircraft":
		return 3
	case "drone":
		return 2
	case "person":
		return 0.5
	case "vehicle", "robot":
		return 1
	default:
		return 1
	}
}

func (f *CV) toENU(lat, lon float64) (float64, float64) {
	e := (lon - f.refLon) * math.Pi / 180 * earthR * math.Cos(f.refLat*math.Pi/180)
	n := (lat - f.refLat) * math.Pi / 180 * earthR
	return e, n
}

func (f *CV) toLatLon(e, n float64) (float64, float64) {
	lat := f.refLat + n/earthR*180/math.Pi
	lon := f.refLon + e/(earthR*math.Cos(f.refLat*math.Pi/180))*180/math.Pi
	return lat, lon
}

func (f *CV) predictTo(t int64) ([4]float64, [4][4]float64) {
	x, p := f.x, f.p
	dt := float64(t-f.t) / 1000
	if dt <= 0 {
		return x, p
	}
	x[0] += dt * x[2]
	x[1] += dt * x[3]
	// P = F P F^T + Q, F = [[I, dt I], [0, I]].
	var fp [4][4]float64
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			fp[i][j] = p[i][j]
			if i < 2 {
				fp[i][j] += dt * p[i+2][j]
			}
		}
	}
	var np [4][4]float64
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			np[i][j] = fp[i][j]
			if j < 2 {
				np[i][j] += dt * fp[i][j+2]
			}
		}
	}
	q2 := f.q * f.q
	for a := 0; a < 2; a++ {
		np[a][a] += q2 * dt * dt * dt * dt / 4
		np[a][a+2] += q2 * dt * dt * dt / 2
		np[a+2][a] += q2 * dt * dt * dt / 2
		np[a+2][a+2] += q2 * dt * dt
	}
	return x, np
}

// Update implements Estimator.
func (f *CV) Update(m Measurement) error {
	sig := math.Max(m.SigmaM, 0.5)
	if !f.init {
		f.refLat, f.refLon, f.t, f.init = m.Lat, m.Lon, m.T, true
		f.x = [4]float64{0, 0, 0, 0}
		f.p = [4][4]float64{}
		f.p[0][0], f.p[1][1] = sig*sig, sig*sig
		vs := 50.0 // unknown velocity: wide prior
		if m.HasVel {
			f.x[2], f.x[3] = m.VE, m.VN
			vs = math.Max(m.SigmaV, 0.1)
		}
		f.p[2][2], f.p[3][3] = vs*vs, vs*vs
		return nil
	}
	if m.T < f.t {
		return ErrOutOfOrder
	}
	x, p := f.predictTo(m.T)
	e, n := f.toENU(m.Lat, m.Lon)
	x, p = update(x, p, 0, [2]float64{e, n}, sig*sig)
	if m.HasVel {
		sv := math.Max(m.SigmaV, 0.1)
		x, p = update(x, p, 2, [2]float64{m.VE, m.VN}, sv*sv)
	}
	f.x, f.p, f.t = x, p, m.T
	// Re-anchor the plane at the new estimate so the flat-earth error stays small.
	f.refLat, f.refLon = f.toLatLon(f.x[0], f.x[1])
	f.x[0], f.x[1] = 0, 0
	return nil
}

// update applies a 2-D measurement z of state components (off, off+1) with variance r on
// each axis: K = P H^T (H P H^T + R)^-1.
func update(x [4]float64, p [4][4]float64, off int, z [2]float64, r float64) ([4]float64, [4][4]float64) {
	a, b, c, d := p[off][off]+r, p[off][off+1], p[off+1][off], p[off+1][off+1]+r
	det := a*d - b*c
	if det == 0 {
		return x, p
	}
	inv := [2][2]float64{{d / det, -b / det}, {-c / det, a / det}}
	var k [4][2]float64
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			k[i][j] = p[i][off]*inv[0][j] + p[i][off+1]*inv[1][j]
		}
	}
	y := [2]float64{z[0] - x[off], z[1] - x[off+1]}
	for i := 0; i < 4; i++ {
		x[i] += k[i][0]*y[0] + k[i][1]*y[1]
	}
	var np [4][4]float64
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			np[i][j] = p[i][j] - (k[i][0]*p[off][j] + k[i][1]*p[off+1][j])
		}
	}
	return x, np
}

func (f *CV) state(x [4]float64, p [4][4]float64, t int64) State {
	lat, lon := f.toLatLon(x[0], x[1])
	return State{
		T: t, Lat: lat, Lon: lon, VE: x[2], VN: x[3],
		PosCov: [3]float64{p[0][0], p[0][1], p[1][1]},
		VelCov: [3]float64{p[2][2], p[2][3], p[3][3]},
	}
}

// Predict implements Estimator.
func (f *CV) Predict(t int64) State {
	x, p := f.predictTo(t)
	if t < f.t {
		t = f.t
	}
	return f.state(x, p, t)
}

// State implements Estimator.
func (f *CV) State() State { return f.state(f.x, f.p, f.t) }

// Mahalanobis2 returns the squared Mahalanobis distance of a position measurement from the
// prediction at m.T, for chi-square gating (2 dof; 9.21 is the 99% gate).
func Mahalanobis2(est Estimator, m Measurement) float64 {
	s := est.Predict(m.T)
	cv := &CV{refLat: s.Lat, refLon: s.Lon}
	e, n := cv.toENU(m.Lat, m.Lon)
	r := m.SigmaM * m.SigmaM
	a, b, d := s.PosCov[0]+r, s.PosCov[1], s.PosCov[2]+r
	det := a*d - b*b
	if det <= 0 {
		return math.Inf(1)
	}
	return (d*e*e - 2*b*e*n + a*n*n) / det
}

// DistanceM is the haversine distance in metres.
func DistanceM(lat1, lon1, lat2, lon2 float64) float64 {
	const rad = math.Pi / 180
	dlat := (lat2 - lat1) * rad
	dlon := (lon2 - lon1) * rad
	h := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	return 2 * 6371008.8 * math.Asin(math.Min(1, math.Sqrt(h)))
}

// VelocityFrom converts speed (m/s) and course (degrees true) to east/north components.
func VelocityFrom(speed, course float64) (float64, float64) {
	r := course * math.Pi / 180
	return speed * math.Sin(r), speed * math.Cos(r)
}
