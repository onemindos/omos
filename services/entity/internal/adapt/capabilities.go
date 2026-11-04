package adapt

import (
	"strings"
)

// Capabilities (entity.md §4.1 group C, ONTOLOGY §4 "Capability", mission-task.md §4 capability
// gate). A fleet bridge advertises the verbs it will execute in the status envelope's
// `meta.commands_accepted` (omos #550 MAVLink: "goto,rtl,land,arm,disarm,estop,takeoff"; #566
// ROS 2: "goto,hold,rtl,estop,estop_release"). The entity service turns that into the task types
// the policy service gates on (`not_capable`, omos #554 checkAsset).
//
// Constraints, all fail closed:
//   - only `sim-*` fleet ids get capabilities (stop-line: simulated assets only); a feed, a CoT
//     record or a non-sim fleet id never gets any, and nothing is inferred from the entity class;
//   - a verb counts only if the bridge advertises it; there is no class default;
//   - the class bounds what may be advertised: `land` and `orbit` are air verbs, and `observe`
//     needs an observable asset (one with a sensor), so a ground robot gets goto/hold/rtl at most;
//   - verbs with no task type (arm, disarm, takeoff, estop, estop_release) are dropped: estop is an
//     unblockable policy action that needs no capability, the rest are not taskable.

// taskTypes is the asset-verb part of the mission-task.md §4 catalogue, in catalogue order.
var taskTypes = []struct{ verb, task string }{
	{"goto", "onemind.task.goto.v1"},
	{"hold", "onemind.task.hold.v1"},
	{"rtl", "onemind.task.rtl.v1"},
	{"land", "onemind.task.land.v1"},
	{"orbit", "onemind.task.orbit.v1"},
	{"observe", "onemind.task.observe.v1"},
}

// classVerbs bounds the verbs each fleet entity_type may be given.
var classVerbs = map[string]map[string]bool{
	"drone":   {"goto": true, "hold": true, "rtl": true, "land": true, "orbit": true, "observe": true},
	"robot":   {"goto": true, "hold": true, "rtl": true, "observe": true},
	"vehicle": {"goto": true, "hold": true, "rtl": true, "observe": true},
}

// advertised reads `commands_accepted` as a comma-separated string or a JSON array of strings.
// Each item may be a bare verb ("goto") or a task type ("onemind.task.goto.v1"). ok is false
// when the key is absent, so a position message (no list) leaves the entity's list unchanged.
func advertised(meta map[string]any) (verbs map[string]bool, ok bool) {
	raw, ok := meta["commands_accepted"]
	if !ok {
		return nil, false
	}
	var items []string
	switch v := raw.(type) {
	case string:
		items = strings.Split(v, ",")
	case []any:
		for _, x := range v {
			if s, isStr := x.(string); isStr {
				items = append(items, s)
			}
		}
	default:
		return map[string]bool{}, true // present but unreadable: advertise nothing
	}
	verbs = map[string]bool{}
	for _, it := range items {
		it = strings.ToLower(strings.TrimSpace(it))
		for _, tt := range taskTypes {
			if it == tt.verb || it == tt.task {
				verbs[tt.verb] = true
			}
		}
	}
	return verbs, true
}

// FleetCapabilities returns the task types for a fleet status, or ok=false when the message
// carries no `commands_accepted` (keep what the entity has). It returns an empty, non-nil list
// for an explicit list with no taskable verb, which clears the entity's capabilities.
func FleetCapabilities(entityType, uid string, traits []string, meta map[string]any) (caps []string, ok bool) {
	if !strings.HasPrefix(strings.ToLower(uid), "sim-") {
		return nil, false
	}
	verbs, ok := advertised(meta)
	if !ok {
		return nil, false
	}
	allowed := classVerbs[entityType]
	observable := false
	for _, t := range traits {
		observable = observable || t == "observable"
	}
	caps = []string{}
	for _, tt := range taskTypes {
		if !verbs[tt.verb] || !allowed[tt.verb] || (tt.verb == "observe" && !observable) {
			continue
		}
		caps = append(caps, tt.task)
	}
	return caps, true
}
