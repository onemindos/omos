// Package watch is the watch-job runner (FEATURES I3). Detection is deterministic: each
// watch is a rule over JetStream and KV inputs inside an explicit envelope (the subjects and
// buckets it may read, its thresholds and its rate limits). When a rule fires, the runner
// publishes the alert first; only then may a language model explain it and pick at most one
// response from the watch's closed list, which the runner sends to
// `svc.policy.task.propose`. The runner never confirms: it has no code path to
// `svc.policy.*.confirm`, and its NATS user is denied those subjects
// (onemind-xr docs/research/agents-that-act.md §4, §8; FEATURES D5, I3).
package watch

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Rule kinds.
const (
	KindZone      = "zone"       // an entity is outside (or inside) a KV geofences zone
	KindStale     = "stale"      // an entity has not updated within max_age_s, or is stale/lost/deleted
	KindThreshold = "threshold"  // a numeric entity field crosses a value, with hysteresis
	KindFeedStale = "feed_stale" // a collector's tel.feed.{source}.status shows stale, failing or silent
)

// Hard caps: a watch cannot be configured to flood the bus, the model or the operator.
const (
	CapAlertsPerHour    = 600
	CapLLMPerHour       = 60
	CapProposalsPerHour = 20
	MinCooldownS        = 10
)

// AllowedTaskTypes is the closed list of responses a watch may propose. Only types whose
// spec the runner can fill without a model inventing coordinates: return, hold, a manual
// task for a person or team, a question for an agent. Everything still goes through policy
// and a human confirm.
var AllowedTaskTypes = map[string]bool{
	"onemind.task.rtl.v1":           true,
	"onemind.task.hold.v1":          true,
	"onemind.task.manual.v1":        true,
	"onemind.task.agent_request.v1": true,
}

// stopWords mirrors the policy service's civilian stop line (defence in depth; the policy
// service is the enforcing check, PR #554 `StopLine`).
var stopWords = []string{"strike", "engage", "intercept", "designate", "attack", "kill", "weapon", "hostile", "follow", "track_person", "face", "payload_release"}

// File is the watch configuration file.
type File struct {
	Watches []Watch `yaml:"watches"`
}

// Watch is one watch job.
type Watch struct {
	ID          string    `yaml:"id"`
	Owner       string    `yaml:"owner"` // user:{id}, the human the watch acts for
	Description string    `yaml:"description"`
	Level       string    `yaml:"level"` // advisory | caution | warning
	Disabled    bool      `yaml:"disabled"`
	Rule        Rule      `yaml:"rule"`
	Inputs      []string  `yaml:"inputs"`
	Rate        Rate      `yaml:"rate"`
	OnTrigger   OnTrigger `yaml:"on_trigger"`
}

// Rule is the deterministic condition.
type Rule struct {
	Kind        string   `yaml:"kind"`
	EntityTypes []string `yaml:"entity_types"`
	EntityIDs   []string `yaml:"entity_ids"` // exact ids, or a prefix ending in *
	Zone        string   `yaml:"zone"`
	Alert       string   `yaml:"alert"` // zone: outside | inside
	MaxAgeS     float64  `yaml:"max_age_s"`
	Field       string   `yaml:"field"`
	Op          string   `yaml:"op"` // lt | le | gt | ge
	Value       float64  `yaml:"value"`
	Clear       *float64 `yaml:"clear"`
	Sources     []string `yaml:"sources"` // feed_stale: empty = every source
	Reasons     []string `yaml:"reasons"` // feed_stale: stale | failing | silent (default all)
}

// Rate limits.
type Rate struct {
	CooldownS           float64 `yaml:"cooldown_s"`
	MaxAlertsPerHour    int     `yaml:"max_alerts_per_hour"`
	MaxLLMPerHour       int     `yaml:"max_llm_per_hour"`
	MaxProposalsPerHour int     `yaml:"max_proposals_per_hour"`
}

// OnTrigger says what happens after the alert is published.
type OnTrigger struct {
	Explain bool     `yaml:"explain"`
	Propose *Propose `yaml:"propose"`
}

// Propose is the closed list of responses a model may choose from.
type Propose struct {
	ActionType string   `yaml:"action_type"` // task.assign only
	TaskTypes  []string `yaml:"task_types"`
	Assignee   string   `yaml:"assignee"` // "about" (the entity that fired) or "{entity_type}/{entity_id}"
}

var idRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,62}$`)
var tokenRe = regexp.MustCompile(`^[a-z0-9_-]{1,64}$`)

// Load reads and validates a watch file. Unknown fields are errors.
func Load(path string) ([]Watch, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(b)
}

// Parse decodes and validates watch YAML.
func Parse(b []byte) ([]Watch, error) {
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	var f File
	if err := dec.Decode(&f); err != nil {
		return nil, fmt.Errorf("watch file: %w", err)
	}
	seen := map[string]bool{}
	var out []Watch
	for i := range f.Watches {
		w := f.Watches[i]
		if err := w.normalise(); err != nil {
			return nil, fmt.Errorf("watch %q: %w", w.ID, err)
		}
		if seen[w.ID] {
			return nil, fmt.Errorf("watch %q: duplicate id", w.ID)
		}
		seen[w.ID] = true
		out = append(out, w)
	}
	return out, nil
}

func (w *Watch) normalise() error {
	if !idRe.MatchString(w.ID) {
		return errors.New("id must match ^[a-z0-9][a-z0-9_-]{0,62}$")
	}
	if !strings.HasPrefix(w.Owner, "user:") || len(w.Owner) < 6 {
		return errors.New("owner must be a human principal user:{id}; proposals are made on its behalf")
	}
	switch w.Level {
	case "":
		w.Level = "caution"
	case "advisory", "caution", "warning":
	default:
		return errors.New("level must be advisory, caution or warning")
	}
	r := &w.Rule
	need := []string{}
	switch r.Kind {
	case KindZone, KindStale, KindThreshold:
		if len(r.EntityTypes) == 0 {
			return errors.New("entity_types is required (the envelope must be explicit)")
		}
		for _, t := range r.EntityTypes {
			if !tokenRe.MatchString(t) {
				return fmt.Errorf("entity type %q is not a subject token", t)
			}
			need = append(need, "ent."+t+".*.state")
		}
		for _, id := range r.EntityIDs {
			if !tokenRe.MatchString(strings.TrimSuffix(id, "*")) {
				return fmt.Errorf("entity id %q is not a subject token", id)
			}
		}
	case KindFeedStale:
		need = append(need, "tel.feed.*.status")
		for _, s := range r.Reasons {
			if s != "stale" && s != "failing" && s != "silent" {
				return fmt.Errorf("reason %q must be stale, failing or silent", s)
			}
		}
	default:
		return fmt.Errorf("rule kind %q must be zone, stale, threshold or feed_stale", r.Kind)
	}
	switch r.Kind {
	case KindZone:
		if !tokenRe.MatchString(r.Zone) {
			return errors.New("zone must be a geofences key")
		}
		if r.Alert == "" {
			r.Alert = "outside"
		}
		if r.Alert != "outside" && r.Alert != "inside" {
			return errors.New("alert must be outside or inside")
		}
		need = append(need, "kv:geofences")
	case KindStale:
		if r.MaxAgeS < 5 {
			return errors.New("max_age_s must be at least 5")
		}
	case KindThreshold:
		if r.Field == "" || strings.Contains(r.Field, "..") {
			return errors.New("field must be a dotted path such as health.battery_pct")
		}
		switch r.Op {
		case "lt", "le", "gt", "ge":
		default:
			return errors.New("op must be lt, le, gt or ge")
		}
		if r.Clear != nil {
			low := r.Op == "lt" || r.Op == "le"
			if (low && *r.Clear < r.Value) || (!low && *r.Clear > r.Value) {
				return errors.New("clear must be on the safe side of value (hysteresis)")
			}
		}
	}
	// The envelope: every input the rule needs must be declared, and nothing else may be.
	for _, n := range need {
		ok := false
		for _, in := range w.Inputs {
			if covers(in, n) {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("inputs do not declare %s, which the rule reads", n)
		}
	}
	for _, in := range w.Inputs {
		used := false
		for _, n := range need {
			if covers(in, n) {
				used = true
			}
		}
		if !used {
			return fmt.Errorf("input %s is not read by the rule; remove it (the envelope is exact)", in)
		}
		if strings.HasSuffix(in, ">") || in == "ent.*.*.state" {
			return fmt.Errorf("input %s is too wide; name entity types", in)
		}
	}
	rt := &w.Rate
	if rt.CooldownS == 0 {
		rt.CooldownS = 300
	}
	if rt.MaxAlertsPerHour == 0 {
		rt.MaxAlertsPerHour = 30
	}
	if rt.MaxLLMPerHour == 0 {
		rt.MaxLLMPerHour = 10
	}
	if rt.MaxProposalsPerHour == 0 {
		rt.MaxProposalsPerHour = 4
	}
	if rt.CooldownS < MinCooldownS || rt.MaxAlertsPerHour > CapAlertsPerHour || rt.MaxLLMPerHour > CapLLMPerHour ||
		rt.MaxProposalsPerHour > CapProposalsPerHour || rt.MaxAlertsPerHour < 0 || rt.MaxLLMPerHour < 0 || rt.MaxProposalsPerHour < 0 {
		return fmt.Errorf("rate outside the caps (cooldown >= %ds, alerts <= %d/h, llm <= %d/h, proposals <= %d/h)",
			MinCooldownS, CapAlertsPerHour, CapLLMPerHour, CapProposalsPerHour)
	}
	if p := w.OnTrigger.Propose; p != nil {
		if p.ActionType == "" {
			p.ActionType = "task.assign"
		}
		if p.ActionType != "task.assign" {
			return errors.New("propose.action_type must be task.assign: a watch proposes tasks, never raw commands")
		}
		if len(p.TaskTypes) == 0 {
			return errors.New("propose.task_types is required")
		}
		for _, t := range p.TaskTypes {
			if !AllowedTaskTypes[t] {
				return fmt.Errorf("task type %s is not in the watch allowlist (rtl, hold, manual, agent_request)", t)
			}
			if word := StopWord(t); word != "" {
				return fmt.Errorf("task type %s crosses the civilian stop line (%s)", t, word)
			}
		}
		if p.Assignee == "" {
			p.Assignee = "about"
		}
		if p.Assignee != "about" {
			typ, id, ok := strings.Cut(p.Assignee, "/")
			if !ok || !tokenRe.MatchString(typ) || !tokenRe.MatchString(id) {
				return errors.New(`propose.assignee must be "about" or "{entity_type}/{entity_id}"`)
			}
		}
		if r.Kind == KindFeedStale && p.Assignee == "about" {
			return errors.New("a feed_stale watch must name its assignee (a feed is not an asset)")
		}
	}
	return nil
}

// StopWord returns the first civilian stop-line word in s, or "".
func StopWord(s string) string {
	l := strings.ToLower(s)
	for _, w := range stopWords {
		if strings.Contains(l, w) {
			return w
		}
	}
	return ""
}

// covers reports whether NATS subject pattern p matches every subject that pattern s
// matches (`*` in s is matched only by `*` or `>` in p). "kv:bucket" inputs match exactly.
func covers(p, s string) bool {
	if strings.HasPrefix(p, "kv:") || strings.HasPrefix(s, "kv:") {
		return p == s
	}
	pt, st := strings.Split(p, "."), strings.Split(s, ".")
	for i, t := range pt {
		if t == ">" {
			return i < len(st)
		}
		if i >= len(st) {
			return false
		}
		if t == "*" {
			continue
		}
		if t != st[i] {
			return false
		}
	}
	return len(pt) == len(st)
}

// EntityFilters returns the ENT subject filters the watches read (their declared inputs).
func EntityFilters(ws []Watch) []string {
	set := map[string]bool{}
	var out []string
	for _, w := range ws {
		if w.Disabled {
			continue
		}
		for _, in := range w.Inputs {
			if strings.HasPrefix(in, "ent.") && !set[in] {
				set[in] = true
				out = append(out, in)
			}
		}
	}
	return out
}

// NeedsZones and NeedsFeeds report whether any watch reads those inputs.
func NeedsZones(ws []Watch) bool { return anyInput(ws, "kv:geofences") }
func NeedsFeeds(ws []Watch) bool { return anyInput(ws, "tel.feed.*.status") }

func anyInput(ws []Watch, in string) bool {
	for _, w := range ws {
		if w.Disabled {
			continue
		}
		for _, i := range w.Inputs {
			if covers(i, in) {
				return true
			}
		}
	}
	return false
}
