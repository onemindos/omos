package policy

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestCanonicalizeRFC8785(t *testing.T) {
	cases := map[string]string{
		`{"b":1,"a":[true,null,"x"]}`:     `{"a":[true,null,"x"],"b":1}`,
		`1e21`:                            `1e+21`,
		`1e-7`:                            `1e-7`,
		`0.000001`:                        `0.000001`,
		`333333333.33333329`:              `333333333.3333333`,
		`-0`:                              `0`,
		`100`:                             `100`,
		`4.50`:                            `4.5`,
		`"\u20ac<>&\u0001"`:               "\"\u20ac<>&\\u0001\"",
		`{"\u20ac":1,"\r":2,"1":3}`:       "{\"\\r\":2,\"1\":3,\"\u20ac\":1}",
		`{"peach":1,"péché":2,"pêche":3}`: `{"peach":1,"péché":2,"pêche":3}`,
	}
	for in, want := range cases {
		got, err := Canonicalize([]byte(in))
		if err != nil || string(got) != want {
			t.Errorf("Canonicalize(%s) = %s, %v; want %s", in, got, err, want)
		}
	}
}

func TestParamsHashIgnoresKeyOrder(t *testing.T) {
	tg := Target{EntityType: "drone", EntityID: "sim-x1"}
	a, _ := ParamsHash("asset.goto", 1, tg, json.RawMessage(`{"lat":30.269,"lon":-97.742,"hae_m":160}`))
	b, _ := ParamsHash("asset.goto", 1, tg, json.RawMessage(`{ "hae_m":160.0, "lon":-97.742, "lat":30.2690 }`))
	c, _ := ParamsHash("asset.goto", 1, tg, json.RawMessage(`{"lat":30.27,"lon":-97.742,"hae_m":160}`))
	if a != b || a == c || !strings.HasPrefix(a, "sha256:") || len(a) != 71 {
		t.Fatalf("a=%s b=%s c=%s", a, b, c)
	}
}

func TestStopLine(t *testing.T) {
	for _, n := range []string{"asset.strike", "asset.engage", "asset.intercept", "asset.designate", "designate_target",
		"asset.follow_person", "onemind.task.engage.v1", "asset.weapon_release", "kill", "t-k-d", "Asset.STRIKE"} {
		if hit, _ := StopLine(n); !hit {
			t.Errorf("%s passed the stop line", n)
		}
	}
	for _, n := range []string{"asset.goto", "asset.orbit", "asset.observe", "onemind.task.follow_route.v1",
		"onemind.task.survey_area.v1", "asset.estop", "task.assign", "alert.acknowledge"} {
		if hit, w := StopLine(n); hit {
			t.Errorf("%s blocked by %q", n, w)
		}
	}
	// No catalogue entry or task type may itself cross the stop line.
	for n := range Catalogue {
		if hit, _ := StopLine(n); hit {
			t.Errorf("catalogue contains stop-line name %s", n)
		}
	}
	for n := range TaskTypes {
		if hit, _ := StopLine(n); hit {
			t.Errorf("task catalogue contains stop-line name %s", n)
		}
	}
	if !scanTK(json.RawMessage(`{"a":{"cot_type":"t-k-d"}}`)) || scanTK(json.RawMessage(`{"a":"t-x"}`)) {
		t.Fatal("scanTK")
	}
}

func TestApprovalLevels(t *testing.T) {
	want := map[Level][2]any{A0: {0, true}, A1: {0, true}, A2: {1, true}, A3: {1, true}, A4: {2, true}, A5: {0, false}}
	for l, w := range want {
		n, ok := ApprovalsNeeded(l)
		if n != w[0].(int) || ok != w[1].(bool) {
			t.Errorf("%s: %d %v", l, n, ok)
		}
	}
	for n, s := range Catalogue {
		if strings.HasPrefix(n, "asset.") && s.Level != A3 && s.Level != A4 {
			t.Errorf("asset verb %s below A3: every asset command needs a human confirm", n)
		}
	}
}

func TestSchemasCoverCatalogue(t *testing.T) {
	s, err := LoadSchemas(schemaDir)
	if err != nil {
		t.Fatal(err)
	}
	for n, spec := range Catalogue {
		if spec.Level == A5 {
			continue
		}
		for _, v := range spec.Versions {
			if _, ok := s.actions[n+".v"+itoa(v)]; !ok {
				t.Errorf("no schema for %s v%d", n, v)
			}
		}
	}
	for tt := range TaskTypes {
		if !s.HasTask(tt) {
			t.Errorf("no schema for task type %s", tt)
		}
	}
	if err := s.ValidateAction("asset.goto", 1, json.RawMessage(`{"lat":91,"lon":0,"hae_m":1}`)); err == nil {
		t.Error("lat 91 accepted")
	}
}

func itoa(i int) string { b, _ := json.Marshal(i); return string(b) }

func TestGeofence(t *testing.T) {
	sq, ok := ParseZone("k", []byte(`{"entity_id":"z","geo":{"shape":"polygon","kind":"keep_out","coords":{"type":"Polygon","coordinates":[[[0,0],[1,0],[1,1],[0,1],[0,0]],[[0.4,0.4],[0.6,0.4],[0.6,0.6],[0.4,0.6],[0.4,0.4]]]}}}`))
	if !ok {
		t.Fatal("parse polygon")
	}
	ci, ok := ParseZone("c", []byte(`{"geo":{"shape":"circle","kind":"keep_in","radius_m":1000,"coords":{"type":"Point","coordinates":[10,10]}},"properties":{"applies_to":["sim-x1"]}}`))
	if !ok {
		t.Fatal("parse circle")
	}
	if !sq.Contains(0.2, 0.2) || sq.Contains(0.5, 0.5) || sq.Contains(2, 2) {
		t.Fatal("polygon with hole")
	}
	zones := []Zone{sq, ci}
	if CheckPoint(zones, "sim-x1", 0.2, 0.2) != "keep_out" {
		t.Fatal("keep_out")
	}
	if CheckPoint(zones, "sim-x1", 5, 5) != "keep_in" || CheckPoint(zones, "sim-x9", 5, 5) != "" {
		t.Fatal("keep_in applies_to")
	}
	if CheckPoint(zones, "sim-x1", 10.005, 10) != "" {
		t.Fatal("inside circle")
	}
	if _, ok := ParseZone("x", []byte(`{"geo":{"shape":"corridor","coords":{}}}`)); ok {
		t.Fatal("unsupported shape accepted")
	}
}

func TestModelTypeRestrictions(t *testing.T) {
	ev, err := LoadModel(modelFile)
	if err != nil {
		t.Fatal(err)
	}
	// Agents can never hold a human-only relation.
	for _, rel := range []string{"approver", "second_approver", "executor"} {
		if err := ev.Write(Tuple{User: "agent:hermes/zeus/legacy", Relation: rel, Object: "action_type:asset.goto"}); err == nil {
			t.Errorf("agent written as %s", rel)
		}
	}
	if err := ev.Write(Tuple{User: "agent:hermes/zeus/legacy", Relation: "operator", Object: "asset:sim-x1"}); err == nil {
		t.Error("agent written as asset operator")
	}
	must := func(tp Tuple) {
		if err := ev.Write(tp); err != nil {
			t.Fatal(err)
		}
	}
	must(Tuple{User: "agent:hermes/zeus/legacy", Relation: "proposer", Object: "action_type:asset.goto"})
	must(Tuple{User: "user:zeus", Relation: "member", Object: "group:operators"})
	must(Tuple{User: "group:operators#member", Relation: "approver", Object: "action_type:asset.goto"})
	must(Tuple{User: "site:home", Relation: "site", Object: "asset:sim-x1"})
	must(Tuple{User: "user:alice", Relation: "operator", Object: "site:home"})
	must(Tuple{User: "user:temp", Relation: "operator", Object: "asset:sim-x1", Condition: "within_window",
		CondCtx: map[string]any{"grant_time": "2026-10-29T12:00:00Z", "grant_duration": "1h"}})
	ctx := context.Background()
	at := func(ts string) map[string]any { return map[string]any{"current_time": ts} }
	cases := []struct {
		u, r, o string
		c       map[string]any
		want    bool
	}{
		{"user:zeus", "approver", "action_type:asset.goto", nil, true}, // via group
		{"user:alice", "operator", "asset:sim-x1", nil, true},          // via site
		{"user:alice", "viewer", "asset:sim-x1", nil, true},            // operator implies viewer
		{"agent:hermes/zeus/legacy", "approver", "action_type:asset.goto", nil, false},
		{"user:temp", "operator", "asset:sim-x1", at("2026-10-29T12:30:00Z"), true},
		{"user:temp", "operator", "asset:sim-x1", at("2026-10-29T13:30:00Z"), false},
		{"user:temp", "operator", "asset:sim-x1", nil, false},
	}
	for _, c := range cases {
		got, err := ev.Check(ctx, CheckRequest{User: c.u, Relation: c.r, Object: c.o, Context: c.c})
		if err != nil || got != c.want {
			t.Errorf("%s %s %s: %v %v", c.u, c.r, c.o, got, err)
		}
	}
	ev.SetDown(true)
	if _, err := ev.Check(ctx, CheckRequest{User: "user:zeus", Relation: "approver", Object: "action_type:asset.goto"}); err == nil {
		t.Fatal("down evaluator answered")
	}
}

// The DSL file and the JSON model must define the same types and relations.
func TestModelDSLMatchesJSON(t *testing.T) {
	f, err := os.Open(strings.TrimSuffix(modelFile, ".json") + ".fga")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dsl := map[string][]string{}
	cur := ""
	typeRE, relRE := regexp.MustCompile(`^type (\w+)`), regexp.MustCompile(`^\s+define (\w+):`)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if m := typeRE.FindStringSubmatch(sc.Text()); m != nil {
			cur = m[1]
			dsl[cur] = []string{}
		} else if m := relRE.FindStringSubmatch(sc.Text()); m != nil && cur != "" {
			dsl[cur] = append(dsl[cur], m[1])
		}
	}
	ev, err := LoadModel(modelFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(dsl) != len(ev.types) {
		t.Fatalf("types: dsl %d json %d", len(dsl), len(ev.types))
	}
	for typ, rels := range dsl {
		td, ok := ev.types[typ]
		if !ok {
			t.Fatalf("json lacks type %s", typ)
		}
		var jr []string
		for r := range td.Relations {
			jr = append(jr, r)
		}
		sort.Strings(rels)
		sort.Strings(jr)
		if strings.Join(rels, ",") != strings.Join(jr, ",") {
			t.Errorf("%s: dsl %v json %v", typ, rels, jr)
		}
	}
}

func TestIDsAndDirectory(t *testing.T) {
	u := NewULID(time.Now())
	if !regexp.MustCompile(`^[0-9a-hjkmnp-tv-z]{26}$`).MatchString(u) {
		t.Fatalf("ulid %q", u)
	}
	if !uuidRE.MatchString(NewUUID()) {
		t.Fatal("uuid")
	}
	if !validToken("sim-x1") || validToken("sim.x1") || validToken("*") || validToken(">") || validToken("-x") {
		t.Fatal("validToken")
	}
}

func TestSweepDropsOldFinishedActions(t *testing.T) {
	clock := &fakeClock{t: time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}
	s := &Service{cfg: Config{Now: clock.Now, Retain: time.Hour}, byID: map[string]*Action{}, byRQ: map[string]string{}}
	s.byID["act-old"] = &Action{ActionID: "act-old", RequestID: "r1", State: StateExecuted, UpdatedAt: clock.Now()}
	s.byRQ["r1"] = "act-old"
	s.byID["act-new"] = &Action{ActionID: "act-new", RequestID: "r2", State: StateDenied, UpdatedAt: clock.Now().Add(30 * time.Minute)}
	clock.Advance(61 * time.Minute)
	s.SweepOnce()
	if _, ok := s.byID["act-old"]; ok || s.byRQ["r1"] != "" {
		t.Fatal("old finished action kept")
	}
	if _, ok := s.byID["act-new"]; !ok {
		t.Fatal("recent action dropped")
	}
}
