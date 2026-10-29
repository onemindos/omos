// Package mission is the OneMind task service (FEATURES D2 and H4): missions and tasks per
// onemind-xr docs/contracts/mission-task.md. It applies cmd.mission.{id}.{verb} from the
// policy service, is the only writer of KV mission-state (compare-and-set on the KV
// revision), republishes every accepted change as evt.mission.{id}.{event}, and serves
// svc.mission.{get,list,report,apply}.
package mission

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Task states (mission-task.md section 5).
const (
	StateCreated           = "created"
	StateScheduled         = "scheduled"
	StateSent              = "sent"
	StateReceived          = "received"
	StateAcknowledged      = "acknowledged"
	StateWilco             = "wilco"
	StateExecuting         = "executing"
	StateWaiting           = "waiting"
	StateCancelRequested   = "cancel_requested"
	StateCompleteRequested = "complete_requested"
	StateDoneOK            = "done_ok"
	StateDoneNotOK         = "done_not_ok"
	StateReplaced          = "replaced"
)

// Mission states (section 3.1).
const (
	MissionPlanned   = "planned"
	MissionActive    = "active"
	MissionCompleted = "completed"
	MissionAborted   = "aborted"
)

// Error codes of done_not_ok (section 7) and refusal codes of a command (asset-command.md 6.4).
const (
	ErrRejected  = "rejected"
	ErrCancelled = "cancelled"
	ErrTimeout   = "timeout"
	ErrFailed    = "failed"
	ErrNoAsset   = "no_asset"

	CodeSchemaInvalid      = "schema_invalid"
	CodeUnknownAction      = "unknown_action"
	CodeExpired            = "expired"
	CodeConflict           = "conflict"
	CodeNotFound           = "not_found"
	CodePreconditionFailed = "precondition_failed"
	CodeNotSimulated       = "not_simulated"
	CodeNotAuthorized      = "not_authorized"
	CodeUnavailable        = "unavailable"
)

var terminalStates = map[string]bool{StateDoneOK: true, StateDoneNotOK: true, StateReplaced: true}

// assigneeStates are the states an assignee may report (section 5 "Set by").
var assigneeStates = map[string]bool{
	StateReceived: true, StateAcknowledged: true, StateWilco: true, StateExecuting: true,
	StateWaiting: true, StateDoneOK: true, StateDoneNotOK: true,
}

var errorCodes = map[string]bool{ErrRejected: true, ErrCancelled: true, ErrTimeout: true, ErrFailed: true, ErrNoAsset: true}

// IsTerminal reports whether a task state never changes again (section 5 rule 2).
func IsTerminal(state string) bool { return terminalStates[state] }

// Rollup maps a task state to the omos mission-state enum (section 5 table).
func Rollup(state, errCode string) string {
	switch state {
	case StateCreated, StateScheduled, StateSent:
		return "requested"
	case StateDoneOK:
		return "completed"
	case StateReplaced:
		return "cancelled"
	case StateDoneNotOK:
		if errCode == ErrCancelled {
			return "cancelled"
		}
		return "failed"
	}
	return "running"
}

// MissionRollup maps a mission state to the omos mission-state enum (section 3.1).
func MissionRollup(state string) string {
	switch state {
	case MissionActive:
		return "running"
	case MissionCompleted:
		return "completed"
	case MissionAborted:
		return "cancelled"
	}
	return "requested"
}

// TaskTypes is the catalogue of section 4 with the fleet verb a single-verb type
// dispatches ("" = no asset command; composite types are decomposed, not dispatched).
// Any type not here is refused (section 4 "Rejected by policy, always").
var TaskTypes = map[string]string{
	"onemind.task.goto.v1":              "goto",
	"onemind.task.hold.v1":              "hold",
	"onemind.task.rtl.v1":               "rtl",
	"onemind.task.land.v1":              "land",
	"onemind.task.orbit.v1":             "orbit",
	"onemind.task.observe.v1":           "observe",
	"onemind.task.capture_media.v1":     "observe",
	"onemind.task.survey_area.v1":       "",
	"onemind.task.follow_route.v1":      "",
	"onemind.task.patrol.v1":            "",
	"onemind.task.inspect_structure.v1": "",
	"onemind.task.relay.v1":             "",
	"onemind.task.deliver.v1":           "",
	"onemind.task.manual.v1":            "",
	"onemind.task.agent_request.v1":     "",
}

// AssetClasses take asset tasks and report on evt.{class}.{id}.task.
var AssetClasses = map[string]bool{"drone": true, "robot": true, "vehicle": true, "camera": true}

// AssigneeTypes are every valid assignee entity_type (section 3.2).
var AssigneeTypes = map[string]bool{"drone": true, "robot": true, "vehicle": true, "camera": true,
	"person": true, "team": true, "agent": true}

var (
	missionIDRE = regexp.MustCompile(`^msn-[0-9a-z]{26}$`)
	taskIDRE    = regexp.MustCompile(`^tsk-[0-9a-z]{26}$`)
	uuidRE      = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	tokenRE     = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)
)

// IsMissionID and IsTaskID check the lowercase ULID ids of section 3.
func IsMissionID(s string) bool { return missionIDRE.MatchString(s) }
func IsTaskID(s string) bool    { return taskIDRE.MatchString(s) }

// Ref names an entity: an assignee, an allocation member.
type Ref struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

func (r Ref) String() string { return r.EntityType + "." + r.EntityID }

// Simulated reports whether the ref is a simulated asset (decision 0004 #13: sim-* only).
func (r Ref) Simulated() bool { return strings.HasPrefix(r.EntityID, "sim-") }

type Author struct {
	Kind       string `json:"kind"`
	ID         string `json:"id"`
	OnBehalfOf string `json:"on_behalf_of,omitempty"`
}

type Delivery struct {
	BeforeT    int64 `json:"before_t,omitempty"`
	RequireAck bool  `json:"require_ack"`
}

type Execution struct {
	StartAfterT     *int64 `json:"start_after_t"`
	CompleteBeforeT *int64 `json:"complete_before_t"`
}

type Version struct {
	Definition int `json:"definition"`
	Status     int `json:"status"`
}

type TaskError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Status struct {
	State    string          `json:"state"`
	Error    *TaskError      `json:"error"`
	Progress json.RawMessage `json:"progress"`
	Estimate json.RawMessage `json:"estimate"`
	Result   json.RawMessage `json:"result"`
	StartT   *int64          `json:"start_t"`
	EndT     *int64          `json:"end_t"`
	UpdatedT int64           `json:"updated_t"`
}

// Guard orders assignee reports (section 5 rule 1) across bridge restarts.
//
// A bridge numbers its reports per task from 1 for the life of its process (omos #550,
// owner proposal P7), and knows nothing of the task service's own version.status. So the
// task service keeps the two apart: version.status is its own counter (+1 per accepted
// change, used for the fact's status_version and Nats-Msg-Id), and Guard keeps the
// reporter's sequence. Asset reports are ordered by (command index, status_version): the
// command index is the position of the report's causation (the asset command id) in
// Commands, which holds every asset command of the task in dispatch order. A report from a
// newer command wins even when a restarted bridge started again at 1; within one command
// status_version must increase; a report from an older command is stale.
type Guard struct {
	Commands    []string `json:"commands"`
	Epoch       int      `json:"epoch"`
	LastVersion int      `json:"last_version"`
	Reporter    string   `json:"reporter,omitempty"`
	CancelFrom  int      `json:"cancel_from,omitempty"`
	PrevState   string   `json:"prev_state,omitempty"`
}

// Task is the KV mission-state record under key tsk-... (section 3.2).
type Task struct {
	V               int               `json:"v"`
	TaskID          string            `json:"task_id"`
	ParentID        *string           `json:"parent_id"`
	Type            string            `json:"type"`
	Spec            json.RawMessage   `json:"spec"`
	Description     string            `json:"description,omitempty"`
	Author          Author            `json:"author"`
	Assignee        Ref               `json:"assignee"`
	TargetSnapshots []json.RawMessage `json:"target_snapshots"`
	Delivery        Delivery          `json:"delivery"`
	Execution       *Execution        `json:"execution,omitempty"`
	StaleT          int64             `json:"stale_t"`
	ActionID        string            `json:"action_id"`
	Approval        json.RawMessage   `json:"approval"`
	Version         Version           `json:"version"`
	Status          Status            `json:"status"`
	Allocation      []Ref             `json:"allocation"`
	Rollup          string            `json:"rollup"`
	Simulated       bool              `json:"simulated"`
	CreatedT        int64             `json:"created_t"`
	CreatedBy       string            `json:"created_by"`
	UpdatedT        int64             `json:"updated_t"`
	UpdatedBy       string            `json:"updated_by"`
	Corr            string            `json:"corr,omitempty"`
	Trace           string            `json:"trace,omitempty"`
	Guard           Guard             `json:"report_guard"`
}

type Waypoint struct {
	Seq  int     `json:"seq"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	HaeM float64 `json:"hae_m"`
	Name string  `json:"name,omitempty"`
}

type Assignment struct {
	To     string `json:"to"`
	Role   string `json:"role"`
	SinceT int64  `json:"since_t"`
	By     string `json:"by"`
}

type Window struct {
	StartT *int64 `json:"start_t,omitempty"`
	EndT   *int64 `json:"end_t,omitempty"`
}

type Marking struct {
	TLP          string   `json:"tlp"`
	ReleasableTo []string `json:"releasable_to"`
}

// Mission is the KV mission-state record under key msn-... (section 3.1).
type Mission struct {
	V           int             `json:"v"`
	MissionID   string          `json:"mission_id"`
	Name        string          `json:"name"`
	Intent      string          `json:"intent,omitempty"`
	State       string          `json:"state"`
	Rollup      string          `json:"rollup"`
	AOI         *string         `json:"aoi"`
	Route       []Waypoint      `json:"route"`
	Assignments []Assignment    `json:"assignments"`
	TaskIDs     []string        `json:"task_ids"`
	Window      *Window         `json:"window,omitempty"`
	Marking     Marking         `json:"marking"`
	Simulated   bool            `json:"simulated"`
	TAKMission  json.RawMessage `json:"tak_mission"`
	CreatedT    int64           `json:"created_t"`
	CreatedBy   string          `json:"created_by"`
	UpdatedT    int64           `json:"updated_t"`
	UpdatedBy   string          `json:"updated_by"`
	Corr        string          `json:"corr,omitempty"`
}

// Report is the data of a status report (section 7), from evt.{class}.{id}.task or
// svc.mission.report, and of the republished evt.mission.{tsk}.{state} fact.
type Report struct {
	TaskID        string          `json:"task_id"`
	State         string          `json:"state"`
	StatusVersion int             `json:"status_version"`
	Error         *TaskError      `json:"error"`
	Progress      json.RawMessage `json:"progress,omitempty"`
	Estimate      json.RawMessage `json:"estimate,omitempty"`
	Result        json.RawMessage `json:"result,omitempty"`
	AtT           int64           `json:"at_t"`
	Simulated     bool            `json:"simulated"`
}

// Command is the omos command envelope v1 as the policy service publishes it
// (services/schemas/command/envelope.schema.json plus asset-command.md 6.3 fields).
type Command struct {
	V         int             `json:"v"`
	ID        string          `json:"id"`
	TS        string          `json:"ts"`
	Cmd       string          `json:"cmd"`
	Source    string          `json:"source"`
	Subject   string          `json:"subject"`
	Data      json.RawMessage `json:"data"`
	ActionID  string          `json:"action_id"`
	ExpiresAt string          `json:"expires_at"`
	Simulated bool            `json:"simulated"`
	TaskID    string          `json:"task_id"`
	Corr      string          `json:"corr"`
	Trace     string          `json:"trace"`
	Approval  struct {
		Required  bool            `json:"required"`
		Status    string          `json:"status"`
		Approver  string          `json:"approver"`
		Level     string          `json:"level"`
		Approvals json.RawMessage `json:"approvals"`
	} `json:"approval"`
}

// Event is the omos event envelope v1 (services/schemas/envelope/envelope.schema.json).
type Event struct {
	V         int             `json:"v"`
	ID        string          `json:"id"`
	TS        string          `json:"ts"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	Subject   string          `json:"subject"`
	Data      json.RawMessage `json:"data"`
	Corr      string          `json:"corr,omitempty"`
	Causation string          `json:"causation,omitempty"`
	Trace     string          `json:"trace,omitempty"`
}
