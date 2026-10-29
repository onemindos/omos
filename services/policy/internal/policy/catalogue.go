package policy

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Level is an approval level from onemind-xr docs/design/ONTOLOGY.md section 6.1.
type Level string

const (
	A0 Level = "A0" // read: no gate, logged
	A1 Level = "A1" // reversible UI or low-risk marker: no gate, logged
	A2 Level = "A2" // non-kinetic write: policy + one human confirm
	A3 Level = "A3" // simulated asset command: policy + audit + one human confirm, always
	A4 Level = "A4" // risky asset command: policy + two distinct humans
	A5 Level = "A5" // never through OneMind actions: refused
)

// ApprovalsNeeded returns how many distinct human confirmations a level needs, and
// whether the level can run at all.
func ApprovalsNeeded(l Level) (int, bool) {
	switch l {
	case A0, A1:
		return 0, true
	case A2, A3:
		return 1, true
	case A4:
		return 2, true
	default:
		return 0, false
	}
}

// Family selects the endpoint pair: svc.policy.command.* or svc.policy.task.*.
type Family string

const (
	FamilyCommand Family = "command"
	FamilyTask    Family = "task"
)

// Action types execute in one of these ways once approved.
const (
	execFleet   = "fleet"   // cmd.fleet.{class}.{id}.{verb}
	execEntity  = "entity"  // cmd.entity.{id}.{verb}
	execMission = "mission" // cmd.mission.{id}.{verb}
	execTask    = "task"    // cmd.mission.{tsk}.create then optional cmd.fleet
	execAbort   = "abort"   // cmd.mission.{tsk}.abort then optional cmd.fleet hold
	execRefused = "refused" // A5: never runs
)

// ActionSpec is one entry of the closed catalogue. Anything not in the catalogue is refused.
type ActionSpec struct {
	Name     string
	Family   Family
	Level    Level // for task.assign and task.abort the level depends on the assignee
	Versions []int
	Targets  []string // allowed target.entity_type values
	Exec     string
	Verb     string // command verb for fleet, entity and mission execution
	// Capability is the task type the asset must advertise for this verb (entity.md capabilities).
	Capability string
	// Unblockable verbs (estop) skip preconditions beyond identity: a stop must not be blocked.
	Unblockable bool
	// Geofence names the parameter whose point must pass keep-in/keep-out ("" = none,
	// "self" = the asset's own location).
	Geofence string
}

// Commandable asset classes (doctrine section 9.1 class vocabulary).
var assetClasses = []string{"drone", "robot", "vehicle"}

func isAssetClass(t string) bool {
	for _, c := range assetClasses {
		if c == t {
			return true
		}
	}
	return false
}

// Catalogue is the closed verb allowlist: asset-command.md section 3, ONTOLOGY.md section 6.2
// and mission-task.md section 4. Order is irrelevant; names are unique.
var Catalogue = map[string]ActionSpec{}

func init() {
	asset := func(verb string, level Level, capability, geofence string, unblockable bool) {
		name := "asset." + verb
		Catalogue[name] = ActionSpec{Name: name, Family: FamilyCommand, Level: level, Versions: []int{1},
			Targets: assetClasses, Exec: execFleet, Verb: verb, Capability: capability,
			Geofence: geofence, Unblockable: unblockable}
	}
	asset("goto", A3, "onemind.task.goto.v1", "point", false)
	asset("hold", A3, "onemind.task.hold.v1", "", false)
	asset("rtl", A3, "onemind.task.rtl.v1", "", false)
	asset("land", A3, "onemind.task.land.v1", "point", false)
	asset("orbit", A3, "onemind.task.orbit.v1", "center", false)
	asset("observe", A3, "onemind.task.observe.v1", "", false)
	asset("estop", A3, "", "", true)
	asset("disarm", A3, "", "", false)
	asset("arm", A4, "", "self", false)
	asset("takeoff", A4, "", "self", false)
	asset("estop_release", A4, "", "", false)
	asset("geofence_override", A4, "", "", false)

	entity := func(name, verb string, level Level) {
		Catalogue[name] = ActionSpec{Name: name, Family: FamilyCommand, Level: level, Versions: []int{1},
			Exec: execEntity, Verb: verb}
	}
	entity("entity.star", "star", A1)
	entity("alert.acknowledge", "acknowledge", A1)
	entity("alert.resolve", "resolve", A1)
	entity("entity.override", "override", A2)

	mission := func(verb string) {
		name := "mission." + verb
		Catalogue[name] = ActionSpec{Name: name, Family: FamilyTask, Level: A2, Versions: []int{1},
			Targets: []string{"mission"}, Exec: execMission, Verb: verb}
	}
	for _, v := range []string{"create", "waypoint", "abort", "complete"} {
		mission(v)
	}
	Catalogue["task.assign"] = ActionSpec{Name: "task.assign", Family: FamilyTask, Level: A3, Versions: []int{1},
		Targets: append([]string{"person", "team", "agent"}, assetClasses...), Exec: execTask}
	Catalogue["task.abort"] = ActionSpec{Name: "task.abort", Family: FamilyTask, Level: A3, Versions: []int{1},
		Targets: append([]string{"person", "team", "agent"}, assetClasses...), Exec: execAbort}

	// A5 examples: named so a request is answered with a clear refusal, never executed.
	for _, n := range []string{"system.credentials", "system.cluster_write", "system.foundation_change"} {
		Catalogue[n] = ActionSpec{Name: n, Family: FamilyCommand, Level: A5, Versions: []int{1}, Exec: execRefused}
	}
}

// TaskTypes maps supported single-verb task types to the fleet verb they dispatch
// (mission-task.md section 4). An empty verb means no asset command (person, team, agent).
var TaskTypes = map[string]string{
	"onemind.task.goto.v1":          "goto",
	"onemind.task.hold.v1":          "hold",
	"onemind.task.rtl.v1":           "rtl",
	"onemind.task.land.v1":          "land",
	"onemind.task.orbit.v1":         "orbit",
	"onemind.task.observe.v1":       "observe",
	"onemind.task.manual.v1":        "",
	"onemind.task.agent_request.v1": "",
}

// Stop line (onemind-xr docs/design/JADC2-MAP.md section 5; ONTOLOGY.md section 6.2):
// these names are rejected as forbidden_action wherever they appear in an action or task type.
var stopLineWords = []string{"strike", "engage", "engagement", "intercept", "designate", "designation",
	"attack", "kill", "weapon", "weapons", "hostile"}

var stopLinePhrases = []string{"follow_person", "follow_a_person", "track_person", "track_a_person",
	"face_recognition", "identity_recognition", "identify_person", "payload_release", "fire_at", "t-k-"}

// StopLine reports whether a type name crosses the civilian stop line, and the word that did.
func StopLine(name string) (bool, string) {
	n := strings.ToLower(name)
	for _, p := range stopLinePhrases {
		if strings.Contains(n, p) {
			return true, p
		}
	}
	for _, tok := range strings.FieldsFunc(n, func(r rune) bool { return r == '.' || r == '_' || r == '-' || r == ' ' }) {
		for _, w := range stopLineWords {
			if tok == w {
				return true, w
			}
		}
	}
	return false, ""
}

// scanTK reports a CoT t-k-* (strike tasking) value anywhere in the parameters.
func scanTK(raw json.RawMessage) bool {
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return false
	}
	var walk func(any) bool
	walk = func(x any) bool {
		switch t := x.(type) {
		case string:
			l := strings.ToLower(t)
			return l == "t-k" || strings.HasPrefix(l, "t-k-")
		case []any:
			for _, e := range t {
				if walk(e) {
					return true
				}
			}
		case map[string]any:
			for _, e := range t {
				if walk(e) {
					return true
				}
			}
		}
		return false
	}
	return walk(v)
}

func (s ActionSpec) supports(version int) bool {
	for _, v := range s.Versions {
		if v == version {
			return true
		}
	}
	return false
}

func (s ActionSpec) allowsTarget(t string) bool {
	if len(s.Targets) == 0 {
		return true
	}
	for _, x := range s.Targets {
		if x == t {
			return true
		}
	}
	return false
}

// commandSubject returns the cmd.* subject for an action with a known verb.
func commandSubject(exec, class, id, verb string) (string, error) {
	if !validToken(id) || !validToken(verb) {
		return "", fmt.Errorf("invalid subject token")
	}
	switch exec {
	case execFleet:
		if !isAssetClass(class) {
			return "", fmt.Errorf("not a commandable class: %s", class)
		}
		return "cmd.fleet." + class + "." + id + "." + verb, nil
	case execEntity:
		return "cmd.entity." + id + "." + verb, nil
	case execMission, execTask, execAbort:
		return "cmd.mission." + id + "." + verb, nil
	}
	return "", fmt.Errorf("no subject for exec %q", exec)
}

// validToken accepts one lowercase NATS subject token without wildcards (omos command
// envelope pattern ^cmd\.[a-z0-9_][a-z0-9_.\-$]*$, and no dots inside one token).
func validToken(s string) bool {
	if s == "" || len(s) > 128 {
		return false
	}
	for i, r := range s {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || (i > 0 && r == '-')
		if !ok {
			return false
		}
	}
	return true
}

// Free-text stop line (I6 red team RT08; JADC2-MAP section 5.3 test 6). The closed verb
// list is the hard guarantee; this catches an agent (or a person) asking in words for what
// no verb allows: following or tracking a person, or a strike-class action, written into a
// reason, task instructions or an agent question. Matching is on lowercase letter tokens.
var (
	textStopWords = map[string]bool{"engage": true, "engaging": true, "engagement": true, "intercept": true,
		"intercepting": true, "designate": true, "designation": true, "hostile": true, "weapon": true,
		"weapons": true, "neutralize": true, "neutralise": true, "strike": true, "strikes": true}
	followVerbs = map[string]bool{"follow": true, "following": true, "followed": true, "track": true,
		"tracking": true, "tail": true, "tailing": true, "shadow": true, "shadowing": true, "pursue": true,
		"pursuing": true, "stalk": true, "stalking": true, "chase": true, "chasing": true, "surveil": true,
		"surveilling": true}
	personNouns = map[string]bool{"person": true, "persons": true, "people": true, "man": true, "men": true,
		"woman": true, "women": true, "him": true, "her": true, "guy": true, "child": true, "kid": true,
		"boy": true, "girl": true, "individual": true, "pedestrian": true, "suspect": true, "someone": true,
		"somebody": true, "intruder": true, "he": true, "she": true}
)

// TextStopLine reports whether free text asks for a stop-line action, and why.
func TextStopLine(s string) (bool, string) {
	if s == "" {
		return false, ""
	}
	toks := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return r < 'a' || r > 'z' })
	for i, t := range toks {
		if textStopWords[t] && !(strings.HasPrefix(t, "strike") && i > 0 && toks[i-1] == "lightning") {
			return true, t
		}
		if followVerbs[t] {
			for j := i + 1; j < len(toks) && j <= i+4; j++ {
				if personNouns[toks[j]] {
					return true, t + " a person"
				}
			}
		}
	}
	return false, ""
}

// scanText applies TextStopLine to every string value in the parameters (including a task
// spec nested in them).
func scanText(raw json.RawMessage) (bool, string) {
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return false, ""
	}
	var walk func(any) (bool, string)
	walk = func(x any) (bool, string) {
		switch t := x.(type) {
		case string:
			return TextStopLine(t)
		case []any:
			for _, e := range t {
				if hit, w := walk(e); hit {
					return hit, w
				}
			}
		case map[string]any:
			for _, e := range t {
				if hit, w := walk(e); hit {
					return hit, w
				}
			}
		}
		return false, ""
	}
	return walk(v)
}
