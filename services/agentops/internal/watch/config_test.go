package watch

import (
	"strings"
	"testing"
)

const good = `
watches:
  - id: sim-keep-in
    owner: user:zeus
    level: warning
    rule: {kind: zone, entity_types: [drone], entity_ids: ["sim-*"], zone: zone-home, alert: outside}
    inputs: [ent.drone.*.state, "kv:geofences"]
    on_trigger:
      explain: true
      propose: {task_types: [onemind.task.rtl.v1, onemind.task.hold.v1]}
  - id: feeds
    owner: user:zeus
    rule: {kind: feed_stale}
    inputs: [tel.feed.*.status]
    on_trigger:
      explain: true
      propose: {task_types: [onemind.task.manual.v1], assignee: team/team-cyan}
`

func TestParseGood(t *testing.T) {
	ws, err := Parse([]byte(good))
	if err != nil {
		t.Fatal(err)
	}
	if len(ws) != 2 || ws[0].Rate.CooldownS != 300 || ws[0].OnTrigger.Propose.ActionType != "task.assign" || ws[0].OnTrigger.Propose.Assignee != "about" {
		t.Fatalf("defaults not applied: %+v", ws[0])
	}
	if f := EntityFilters(ws); len(f) != 1 || f[0] != "ent.drone.*.state" {
		t.Fatalf("filters %v", f)
	}
	if !NeedsZones(ws) || !NeedsFeeds(ws) {
		t.Fatal("needs")
	}
}

func TestParseRefuses(t *testing.T) {
	cases := map[string]string{
		"missing input":    strings.Replace(good, `inputs: [ent.drone.*.state, "kv:geofences"]`, `inputs: [ent.drone.*.state]`, 1),
		"extra input":      strings.Replace(good, `inputs: [tel.feed.*.status]`, `inputs: [tel.feed.*.status, ent.drone.*.state]`, 1),
		"wide input":       strings.Replace(good, `inputs: [ent.drone.*.state, "kv:geofences"]`, `inputs: [ent.>, "kv:geofences"]`, 1),
		"not a human":      strings.Replace(good, `owner: user:zeus`, `owner: agent:hermes/zeus/legacy`, 1),
		"unknown field":    strings.Replace(good, `level: warning`, `level: warning`+"\n    confirm: true", 1),
		"raw command":      strings.Replace(good, `propose: {task_types: [onemind.task.rtl.v1, onemind.task.hold.v1]}`, `propose: {action_type: asset.goto, task_types: [onemind.task.rtl.v1]}`, 1),
		"task not allowed": strings.Replace(good, `onemind.task.hold.v1]`, `onemind.task.orbit.v1]`, 1),
		"stop line":        strings.Replace(good, `onemind.task.hold.v1]`, `onemind.task.intercept.v1]`, 1),
		"feed about":       strings.Replace(good, `, assignee: team/team-cyan`, ``, 1),
		"cooldown cap":     strings.Replace(good, `level: warning`, "level: warning\n    rate: {cooldown_s: 1}", 1),
		"llm cap":          strings.Replace(good, `level: warning`, "level: warning\n    rate: {max_llm_per_hour: 1000}", 1),
		"no types":         strings.Replace(good, `entity_types: [drone], `, ``, 1),
		"dup id":           strings.Replace(good, `id: feeds`, `id: sim-keep-in`, 1),
	}
	for name, y := range cases {
		if _, err := Parse([]byte(y)); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestCovers(t *testing.T) {
	for _, c := range []struct {
		p, s string
		want bool
	}{
		{"ent.drone.*.state", "ent.drone.*.state", true},
		{"ent.*.*.state", "ent.drone.*.state", true},
		{"ent.drone.sim-x1.state", "ent.drone.*.state", false},
		{"ent.>", "ent.drone.*.state", true},
		{"kv:geofences", "kv:geofences", true},
		{"tel.feed.*.status", "ent.drone.*.state", false},
	} {
		if got := covers(c.p, c.s); got != c.want {
			t.Errorf("covers(%s, %s) = %v", c.p, c.s, got)
		}
	}
}
