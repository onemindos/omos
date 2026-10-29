package fusion

import (
	"math"
	"testing"
)

func TestCVTracksStraightLine(t *testing.T) {
	f := NewCV(ProcessNoise("vessel", ""))
	lat0, lon0 := 40.6, -74.04
	// 5 m/s due north, one fix every 10 s with 10 m noise-free measurements.
	for i := 0; i <= 30; i++ {
		lat := lat0 + float64(i)*50/111195.0
		if err := f.Update(Measurement{T: int64(i) * 10000, Lat: lat, Lon: lon0, SigmaM: 10}); err != nil {
			t.Fatal(err)
		}
	}
	s := f.State()
	sp, crs := s.SpeedCourse()
	if math.Abs(sp-5) > 0.2 || math.Min(math.Abs(crs), math.Abs(crs-360)) > 2 {
		t.Fatalf("speed %.2f course %.1f", sp, crs)
	}
	p := f.Predict(310000)
	if d := DistanceM(p.Lat, p.Lon, lat0+31*50/111195.0, lon0); d > 5 {
		t.Fatalf("prediction off by %.1f m", d)
	}
	if s.CEM() >= 10 {
		t.Fatalf("filtered CE %.2f should be below the 10 m measurement sigma", s.CEM())
	}
}

func TestCVRejectsOutOfOrder(t *testing.T) {
	f := NewCV(1)
	_ = f.Update(Measurement{T: 2000, Lat: 1, Lon: 1, SigmaM: 5})
	if err := f.Update(Measurement{T: 1000, Lat: 1, Lon: 1, SigmaM: 5}); err != ErrOutOfOrder {
		t.Fatalf("want ErrOutOfOrder, got %v", err)
	}
}

func TestVelocityMeasurementInitialises(t *testing.T) {
	f := NewCV(2)
	ve, vn := VelocityFrom(10, 90)
	_ = f.Update(Measurement{T: 0, Lat: 30, Lon: -97, SigmaM: 5, HasVel: true, VE: ve, VN: vn, SigmaV: 1})
	sp, crs := f.State().SpeedCourse()
	if math.Abs(sp-10) > 1e-9 || math.Abs(crs-90) > 1e-6 {
		t.Fatalf("speed %.3f course %.3f", sp, crs)
	}
}

func TestMahalanobisGate(t *testing.T) {
	f := NewCV(1)
	_ = f.Update(Measurement{T: 0, Lat: 30, Lon: -97, SigmaM: 10})
	near := Mahalanobis2(f, Measurement{T: 1000, Lat: 30.0001, Lon: -97, SigmaM: 10})
	far := Mahalanobis2(f, Measurement{T: 1000, Lat: 30.01, Lon: -97, SigmaM: 10})
	if near > 9.21 || far < 9.21 {
		t.Fatalf("gate near=%.2f far=%.2f", near, far)
	}
}

func TestDistance(t *testing.T) {
	if d := DistanceM(0, 0, 0, 1); math.Abs(d-111195) > 10 {
		t.Fatalf("1 degree at equator = %.0f m", d)
	}
}
