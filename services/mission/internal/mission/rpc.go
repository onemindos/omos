package mission

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/nats-io/nats.go"
)

// RequestInfoHeader is set by nats-server on a request that crosses a service import with
// share: true; the server replaces any client value, so it is the verified identity (the
// same rule as services/policy, identity.go).
const RequestInfoHeader = "Nats-Request-Info"

// DirectoryEntry maps a NATS identity to a principal. Teams lists the team entity ids the
// principal may report for (team tasks).
type DirectoryEntry struct {
	Account string   `json:"account"`
	User    string   `json:"user"`
	Kind    string   `json:"kind"`
	ID      string   `json:"id"`
	Teams   []string `json:"teams,omitempty"`
}

// Directory is the principals file shared with the policy service (same JSON layout:
// own_account, agent_accounts, principals[]).
type Directory struct {
	OwnAccount    string           `json:"own_account"`
	AgentAccounts []string         `json:"agent_accounts"`
	Entries       []DirectoryEntry `json:"principals"`
}

// LoadDirectory reads the principals file (a ConfigMap; no secrets).
func LoadDirectory(path string) (*Directory, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var d Directory
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	if d.OwnAccount == "" {
		return nil, errors.New("principals file: own_account is required")
	}
	return &d, nil
}

// requestAccount returns the caller's account from the server-set header, or "" for a
// caller in this service's own account (no import crossed).
func requestInfo(m *nats.Msg) (acc, user string) {
	if m.Header == nil {
		return "", ""
	}
	var ci struct {
		Account string `json:"acc"`
		User    string `json:"user"`
	}
	if raw := m.Header.Get(RequestInfoHeader); raw == "" || json.Unmarshal([]byte(raw), &ci) != nil {
		return "", ""
	}
	return ci.Account, ci.User
}

// resolve returns the verified principal of a cross-account request.
func (d *Directory) resolve(m *nats.Msg) (*DirectoryEntry, bool) {
	if d == nil {
		return nil, false
	}
	acc, user := requestInfo(m)
	if acc == "" || acc == d.OwnAccount {
		return nil, false
	}
	for i := range d.Entries {
		e := d.Entries[i]
		if e.Account == acc && e.User == user {
			for _, a := range d.AgentAccounts {
				if a == acc {
					e.Kind = "agent"
				}
			}
			return &e, true
		}
	}
	return nil, false
}

// mayReportFor reports whether the principal is the task's assignee: person "zeus" is
// principal user:zeus; agent "hermes" is agent:hermes or agent:hermes/...; a team is
// reported by any principal whose entry lists it.
func (e *DirectoryEntry) mayReportFor(a Ref) bool {
	switch a.EntityType {
	case "person":
		return e.Kind == "user" && e.ID == "user:"+a.EntityID
	case "agent":
		return e.Kind == "agent" && (e.ID == "agent:"+a.EntityID || strings.HasPrefix(e.ID, "agent:"+a.EntityID+"/"))
	case "team":
		return e.Kind == "user" && contains(e.Teams, a.EntityID)
	}
	return false
}

// reply answers a request. If the answer cannot be sent (for example it exceeds the
// server's max_payload), the caller gets an error reply instead of a timeout.
func (s *Service) reply(m *nats.Msg, v any) {
	if m.Reply == "" {
		return
	}
	b, _ := json.Marshal(v)
	if err := m.Respond(b); err != nil {
		s.cfg.Log.Warn("reply not sent", "subject", m.Subject, "bytes", len(b), "err", err)
		eb, _ := json.Marshal(errReply(CodeUnavailable, "reply not sent: "+err.Error()))
		_ = m.Respond(eb)
	}
}

func errReply(code, msg string) map[string]any {
	return map[string]any{"error": map[string]string{"code": code, "message": msg}}
}

// handleGet serves svc.mission.get {v, id} -> {record, revision}.
func (s *Service) handleGet(m *nats.Msg) {
	var req struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(m.Data, &req) != nil || (!IsMissionID(req.ID) && !IsTaskID(req.ID)) {
		s.reply(m, errReply(CodeSchemaInvalid, "body {v, id: msn-…|tsk-…}"))
		return
	}
	b, rev, err := s.st.get(s.ctx, req.ID)
	if err != nil {
		s.reply(m, errReply(CodeUnavailable, err.Error()))
		return
	}
	if b == nil {
		s.reply(m, errReply(CodeNotFound, "no record "+req.ID))
		return
	}
	s.reply(m, map[string]any{"record": json.RawMessage(b), "revision": rev})
}

// ListRequest is the body of svc.mission.list (section 8). All filters are optional.
type ListRequest struct {
	Kind     string    `json:"kind"`     // mission | task
	State    string    `json:"state"`    // a mission state or task status.state
	Assignee string    `json:"assignee"` // "drone.sim-x1"
	ParentID string    `json:"parent_id"`
	BBox     []float64 `json:"bbox"` // [min_lon, min_lat, max_lon, max_lat]
	Limit    int       `json:"limit"`
	After    string    `json:"after"` // paging: only ids after this one (the previous reply's next)
}

// maxListBytes caps the items of one svc.mission.list reply below the NATS default
// max_payload (1 MiB); the server's own max_payload lowers it further.
const maxListBytes = 900 * 1024

type listItem struct {
	id string
	v  any
}

func inBBox(b []float64, lat, lon float64) bool {
	return lon >= b[0] && lat >= b[1] && lon <= b[2] && lat <= b[3]
}

// handleList serves svc.mission.list from the watched copy of the bucket.
func (s *Service) handleList(m *nats.Msg) {
	var req ListRequest
	if len(m.Data) > 0 && json.Unmarshal(m.Data, &req) != nil {
		s.reply(m, errReply(CodeSchemaInvalid, "body {kind?, state?, assignee?, parent_id?, bbox?, limit?, after?}"))
		return
	}
	if req.BBox != nil && len(req.BBox) != 4 {
		s.reply(m, errReply(CodeSchemaInvalid, "bbox is [min_lon, min_lat, max_lon, max_lat]"))
		return
	}
	limit := req.Limit
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}
	// Missions (msn-…) sort before tasks (tsk-…), so one id order pages both kinds.
	var items []listItem
	if req.Kind == "" || req.Kind == "mission" {
		ms := s.st.cachedMissions()
		sort.Slice(ms, func(i, j int) bool { return ms[i].MissionID < ms[j].MissionID })
		for _, x := range ms {
			if (req.State != "" && x.State != req.State) || req.ParentID != "" {
				continue
			}
			if req.Assignee != "" {
				found := false
				for _, a := range x.Assignments {
					found = found || a.To == req.Assignee
				}
				if !found {
					continue
				}
			}
			if req.BBox != nil {
				in := false
				for _, w := range x.Route {
					in = in || inBBox(req.BBox, w.Lat, w.Lon)
				}
				if !in {
					continue
				}
			}
			items = append(items, listItem{x.MissionID, x})
		}
	}
	if req.Kind == "" || req.Kind == "task" {
		ts := s.st.cachedTasks()
		sort.Slice(ts, func(i, j int) bool { return ts[i].TaskID < ts[j].TaskID })
		for _, x := range ts {
			if (req.State != "" && x.Status.State != req.State) || (req.Assignee != "" && x.Assignee.String() != req.Assignee) {
				continue
			}
			if req.ParentID != "" && (x.ParentID == nil || *x.ParentID != req.ParentID) {
				continue
			}
			if req.BBox != nil {
				var p struct {
					Lat *float64 `json:"lat"`
					Lon *float64 `json:"lon"`
				}
				if json.Unmarshal(x.Spec, &p) != nil || p.Lat == nil || p.Lon == nil || !inBBox(req.BBox, *p.Lat, *p.Lon) {
					continue
				}
			}
			items = append(items, listItem{x.TaskID, x})
		}
	}
	budget := int64(maxListBytes)
	if mp := s.cfg.Conn.MaxPayload(); mp > 0 && mp-64*1024 < budget {
		budget = mp - 64*1024
	}
	out := []json.RawMessage{}
	size, truncated, next := int64(0), false, ""
	for _, it := range items {
		if req.After != "" && it.id <= req.After {
			continue
		}
		b, _ := json.Marshal(it.v)
		if len(out) >= limit || (len(out) > 0 && size+int64(len(b)) > budget) {
			truncated = true
			break
		}
		out = append(out, b)
		size += int64(len(b))
		next = it.id
	}
	resp := map[string]any{"items": out, "count": len(out), "truncated": truncated}
	if truncated {
		resp["next"] = next
	}
	s.reply(m, resp)
}

// handleReport serves svc.mission.report: people and agents report their own tasks with
// the section 7 data. The caller must be the verified assignee (server-set identity from
// a CLIENT or AGENTS import); assets report on evt.{class}.{id}.task instead.
func (s *Service) handleReport(m *nats.Msg) {
	p, ok := s.cfg.Directory.resolve(m)
	if !ok {
		s.reply(m, errReply(CodeNotAuthorized, "svc.mission.report needs a verified person or agent identity"))
		return
	}
	var r Report
	if err := json.Unmarshal(m.Data, &r); err != nil {
		s.reply(m, errReply(CodeSchemaInvalid, err.Error()))
		return
	}
	t, ok := s.st.cachedTask(r.TaskID)
	if !ok {
		tt, _, err := s.st.getTask(s.ctx, r.TaskID)
		if err != nil || tt == nil {
			s.reply(m, errReply(CodeNotFound, "no task "+r.TaskID))
			return
		}
		t = *tt
	}
	if AssetClasses[t.Assignee.EntityType] || !p.mayReportFor(t.Assignee) {
		s.reply(m, errReply(CodeNotAuthorized, fmt.Sprintf("%s is not the assignee %s", p.ID, t.Assignee.String())))
		return
	}
	res, err := s.ApplyReport(r, t.Assignee, p.ID, "", "", false)
	if err != nil {
		s.reply(m, errReply(CodeUnavailable, err.Error()))
		return
	}
	s.reply(m, res)
}

// ApplyRequest is the body of svc.mission.apply (section 6.3), called by the policy
// service after an approved mission.assign.
type ApplyRequest struct {
	V         int    `json:"v"`
	ActionID  string `json:"action_id"`
	Op        string `json:"op"` // assign | unassign
	MissionID string `json:"mission_id"`
	To        string `json:"to"` // "{entity_type}.{entity_id}"
	Role      string `json:"role"`
	By        string `json:"by"`
}

// handleApply serves svc.mission.apply. It refuses every caller from another account; in
// the service's own account only the policy user has the publish grant (A19 identities).
func (s *Service) handleApply(m *nats.Msg) {
	if acc, _ := requestInfo(m); acc != "" && (s.cfg.Directory == nil || acc != s.cfg.Directory.OwnAccount) {
		s.reply(m, errReply(CodeNotAuthorized, "svc.mission.apply is for the policy service only"))
		return
	}
	var req ApplyRequest
	if json.Unmarshal(m.Data, &req) != nil || !IsMissionID(req.MissionID) || req.ActionID == "" ||
		(req.Op != "assign" && req.Op != "unassign") {
		s.reply(m, errReply(CodeSchemaInvalid, "body {v, action_id, op: assign|unassign, mission_id, to, role?, by?}"))
		return
	}
	typ, id, ok := strings.Cut(req.To, ".")
	if !ok || !AssigneeTypes[typ] || !tokenRE.MatchString(id) {
		s.reply(m, errReply(CodeSchemaInvalid, "to must be {entity_type}.{entity_id}"))
		return
	}
	if AssetClasses[typ] && !strings.HasPrefix(id, "sim-") {
		s.reply(m, errReply(CodeNotSimulated, "assets must be simulated sim-* assets"))
		return
	}
	if req.Role == "" {
		req.Role = "member"
	}
	who := req.By
	if who == "" {
		who = "service:policy"
	}
	reason := ""
	mi, err := s.st.mutateMission(s.ctx, req.MissionID, func(x *Mission) error {
		if x.State == MissionCompleted || x.State == MissionAborted {
			reason = "mission is " + x.State
			return errNoChange
		}
		idx := -1
		for i, a := range x.Assignments {
			if a.To == req.To {
				idx = i
			}
		}
		switch {
		case req.Op == "assign" && idx >= 0 && x.Assignments[idx].Role == req.Role:
			reason = "already assigned"
			return errNoChange
		case req.Op == "assign" && idx >= 0:
			x.Assignments[idx].Role = req.Role
		case req.Op == "assign":
			x.Assignments = append(x.Assignments, Assignment{To: req.To, Role: req.Role, SinceT: s.nowMS(), By: who})
			if AssetClasses[typ] {
				x.Simulated = true
			}
		case idx < 0:
			reason = "not assigned"
			return errNoChange
		default:
			x.Assignments = append(x.Assignments[:idx], x.Assignments[idx+1:]...)
		}
		x.UpdatedT, x.UpdatedBy = s.nowMS(), who
		return nil
	})
	switch {
	case errors.Is(err, errNoChange):
		s.reply(m, errReply(CodePreconditionFailed, reason))
		return
	case err != nil:
		s.reply(m, errReply(CodeUnavailable, err.Error()))
		return
	case mi == nil:
		s.reply(m, errReply(CodeNotFound, "no mission "+req.MissionID))
		return
	}
	event := "assigned"
	if req.Op == "unassign" {
		event = "unassigned"
	}
	s.emitMission(mi, event, "", map[string]any{"action_id": req.ActionID, "to": req.To, "role": req.Role})
	s.reply(m, map[string]any{"record": mi})
}
