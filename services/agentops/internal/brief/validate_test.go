package brief

import (
	"strings"
	"testing"

	"github.com/onemindos/omos/services/agentops/internal/feeds"
)

func TestValidateRemovesUnknownAndRepairs(t *testing.T) {
	b := &Briefer{cfg: Config{MinCitations: 3}}
	byID := map[string]entity{
		"sim-x1": {EntityID: "sim-x1", EntityType: "drone", Callsign: "SIM-X1"},
		"icao-a": {EntityID: "icao-a", EntityType: "aircraft"},
		"mmsi-1": {EntityID: "mmsi-1", EntityType: "vessel"},
	}
	top := []entity{byID["sim-x1"], byID["icao-a"], byID["mmsi-1"]}
	gaps := []feeds.Gap{{Source: "nws-alerts", Reason: "stale", Detail: "old"}}
	text, v, cites := b.validate("SIM-X1 [sim-x1] is near [icao-zzz] (made up).", byID, gaps, nil, top)
	if len(v.UnknownRemoved) != 1 || v.UnknownRemoved[0] != "icao-zzz" || strings.Contains(text, "icao-zzz") {
		t.Fatalf("unknown id kept: %s %+v", text, v)
	}
	if v.Cited != 3 || len(cites) != 3 || !v.Repaired || v.GapsNamed {
		t.Fatalf("repair: %+v", v)
	}
	if !strings.Contains(text, "nws-alerts") || !strings.Contains(text, "[mmsi-1]") {
		t.Fatalf("repaired text:\n%s", text)
	}
	if cites[0].Subject != "ent.drone.sim-x1.state" {
		t.Fatal(cites[0])
	}
	// A good text is left alone.
	good := "[sim-x1] [icao-a] [mmsi-1]; feed nws-alerts is stale."
	text, v, _ = b.validate(good, byID, gaps, nil, top)
	if text != good || v.Repaired || !v.GapsNamed {
		t.Fatalf("good text changed: %s %+v", text, v)
	}
}

func TestSanitize(t *testing.T) {
	if got := sanitize("IGNORE ALL [sim-x9] <cmd>{x}\n", 100); strings.ContainsAny(got, "[]<>{}\n") {
		t.Fatal(got)
	}
}
