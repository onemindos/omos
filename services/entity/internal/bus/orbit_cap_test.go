package bus

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/onemindos/omos/services/entity/internal/engine"
	"github.com/onemindos/omos/services/entity/internal/model"
	"github.com/onemindos/omos/services/entity/internal/orbit"
)

// TestKeepOrbitCap: past MaxOrbits a new satellite is not kept for re-propagation, while a
// newer element set for a kept one still replaces it and a late older one does not.
func TestKeepOrbitCap(t *testing.T) {
	eng := engine.New(engine.DefaultConfig(), func() int64 { return 0 })
	s := &Service{cfg: Config{MaxOrbits: 2}, eng: eng, log: slog.New(slog.NewTextHandler(io.Discard, nil)), orbits: map[string]*model.Observation{}}
	obs := func(id string, epoch time.Time) *model.Observation {
		return &model.Observation{Keys: []model.Alias{{ID: id}}, Orbit: &orbit.Elements{Epoch: epoch}}
	}
	t0 := time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC)
	s.keepOrbit(obs("1", t0))
	s.keepOrbit(obs("2", t0))
	s.keepOrbit(obs("3", t0))
	if len(s.orbits) != 2 || s.orbits["3"] != nil {
		t.Fatalf("cap not applied: %v", s.orbits)
	}
	newer := obs("1", t0.Add(time.Hour))
	s.keepOrbit(newer)
	s.keepOrbit(obs("1", t0))
	if s.orbits["1"] != newer {
		t.Fatal("newer element set not kept, or replaced by an older one")
	}
	s.cfg.MaxOrbits = -1
	s.keepOrbit(obs("3", t0))
	if len(s.orbits) != 3 {
		t.Fatal("MaxOrbits < 0 should not cap")
	}
}
