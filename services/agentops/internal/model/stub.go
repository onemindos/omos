package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ExplainInput is the structured part of a PurposeExplain request.
type ExplainInput struct {
	WatchID          string         `json:"watch_id"`
	WatchDescription string         `json:"watch_description,omitempty"`
	RuleKind         string         `json:"rule_kind"`
	AlertID          string         `json:"alert_id"`
	AlertKind        string         `json:"alert_kind"`
	Level            string         `json:"level"`
	Title            string         `json:"title"`
	About            string         `json:"about"`
	Evidence         map[string]any `json:"evidence,omitempty"`
	AllowedTaskTypes []string       `json:"allowed_task_types"`
}

// ExplainOutput is what a model must answer for PurposeExplain (JSON only).
type ExplainOutput struct {
	Explanation string          `json:"explanation"`
	Proposal    *ProposalChoice `json:"proposal"`
}

// ProposalChoice is the model's pick from the allowed task types, or null for none.
type ProposalChoice struct {
	TaskType string `json:"task_type"`
	Reason   string `json:"reason"`
}

// BriefEntity is one entity offered to the brief model.
type BriefEntity struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Callsign  string   `json:"callsign,omitempty"`
	Lifecycle string   `json:"lifecycle,omitempty"`
	Lat       *float64 `json:"lat,omitempty"`
	Lon       *float64 `json:"lon,omitempty"`
	Note      string   `json:"note,omitempty"`
}

// BriefAlert is one open alert offered to the brief model.
type BriefAlert struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Level string `json:"level"`
	About string `json:"about"`
	Title string `json:"title"`
}

// BriefGap is one degraded feed (an intelligence gap, FEATURES B20).
type BriefGap struct {
	Source    string   `json:"source"`
	Reason    string   `json:"reason"`
	AgeS      *float64 `json:"age_s,omitempty"`
	LastError string   `json:"last_error,omitempty"`
	Layers    []string `json:"layers,omitempty"`
}

// BriefInput is the structured part of a PurposeBrief request.
type BriefInput struct {
	AOIID    string         `json:"aoi_id"`
	AOIName  string         `json:"aoi_name"`
	AtUTC    string         `json:"at_utc"`
	Counts   map[string]int `json:"counts"`
	Entities []BriefEntity  `json:"entities"`
	Alerts   []BriefAlert   `json:"alerts"`
	Gaps     []BriefGap     `json:"gaps"`
}

// Stub is a deterministic, offline stand-in for a language model. It is the last link of
// every chain and the model the tests use. It does not invent anything: it restates the
// structured facts it was given.
type Stub struct{}

func (Stub) Name() string { return "stub" }
func (Stub) Local() bool  { return true }

func (Stub) Complete(_ context.Context, r Request) (string, error) {
	switch r.Purpose {
	case PurposeExplain:
		var in ExplainInput
		if err := json.Unmarshal(r.Structured, &in); err != nil {
			return "", fmt.Errorf("stub: explain input: %w", err)
		}
		out := ExplainOutput{Explanation: fmt.Sprintf("Watch %s (%s) fired: %s.", in.WatchID, in.RuleKind, in.Title)}
		if ev := evidenceLine(in.Evidence); ev != "" {
			out.Explanation += " Evidence: " + ev + "."
		}
		if len(in.AllowedTaskTypes) > 0 {
			out.Proposal = &ProposalChoice{TaskType: in.AllowedTaskTypes[0],
				Reason: fmt.Sprintf("First allowed response of watch %s for %s; a human must confirm it.", in.WatchID, in.About)}
		}
		b, _ := json.Marshal(out)
		return string(b), nil
	case PurposeBrief:
		var in BriefInput
		if err := json.Unmarshal(r.Structured, &in); err != nil {
			return "", fmt.Errorf("stub: brief input: %w", err)
		}
		return stubBrief(in), nil
	}
	return "", errors.New("stub: unknown purpose " + r.Purpose)
}

func evidenceLine(ev map[string]any) string {
	keys := make([]string, 0, len(ev))
	for k := range ev {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := []string{}
	for _, k := range keys {
		switch v := ev[k].(type) {
		case string, float64, int, int64, bool:
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
	}
	return strings.Join(parts, ", ")
}

func stubBrief(in BriefInput) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Brief for %s at %s (offline fallback, deterministic summary).\n", in.AOIName, in.AtUTC)
	types := make([]string, 0, len(in.Counts))
	for t := range in.Counts {
		types = append(types, t)
	}
	sort.Strings(types)
	parts := []string{}
	total := 0
	for _, t := range types {
		parts = append(parts, fmt.Sprintf("%d %s", in.Counts[t], t))
		total += in.Counts[t]
	}
	fmt.Fprintf(&b, "Picture: %d entities in the area (%s).\n", total, strings.Join(parts, ", "))
	if len(in.Entities) > 0 {
		b.WriteString("Key entities:\n")
		for i, e := range in.Entities {
			if i == 8 {
				break
			}
			name := e.Callsign
			if name == "" {
				name = e.ID
			}
			line := fmt.Sprintf("- %s (%s) [%s]", name, e.Type, e.ID)
			if e.Lifecycle != "" && e.Lifecycle != "active" {
				line += ", " + e.Lifecycle
			}
			if e.Note != "" {
				line += ", " + e.Note
			}
			b.WriteString(line + "\n")
		}
	}
	if len(in.Alerts) > 0 {
		b.WriteString("Open alerts:\n")
		for _, a := range in.Alerts {
			fmt.Fprintf(&b, "- %s: %s [%s]\n", a.Level, a.Title, a.About)
		}
	} else {
		b.WriteString("Open alerts: none.\n")
	}
	if len(in.Gaps) > 0 {
		b.WriteString("Intelligence gaps:\n")
		for _, g := range in.Gaps {
			line := "- feed " + g.Source + ": " + g.Reason
			if len(g.Layers) > 0 {
				line += " (weakens " + strings.Join(g.Layers, ", ") + ")"
			}
			b.WriteString(line + "\n")
		}
	} else {
		b.WriteString("Intelligence gaps: none reported.\n")
	}
	return b.String()
}
