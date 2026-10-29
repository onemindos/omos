package policy

import (
	"encoding/json"
	"time"
)

// Target names the object an action applies to.
type Target struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

// Actor is who asked (asset-command.md section 6.1 requested_by).
type Actor struct {
	Kind       string `json:"kind"` // user | agent | system
	ID         string `json:"id"`
	OnBehalfOf string `json:"on_behalf_of,omitempty"`
}

// ActionRequest is the svc.policy.{command,task}.propose body (asset-command.md section 6.1).
type ActionRequest struct {
	V             int             `json:"v"`
	RequestID     string          `json:"request_id"`
	ActionType    string          `json:"action_type"`
	ActionVersion int             `json:"action_version"`
	Target        Target          `json:"target"`
	Params        json.RawMessage `json:"params"`
	RequestedBy   Actor           `json:"requested_by"`
	Reason        string          `json:"reason,omitempty"`
	TaskID        string          `json:"task_id,omitempty"`
	Corr          string          `json:"corr,omitempty"`
	Causation     string          `json:"causation,omitempty"`
	DryRun        bool            `json:"dry_run,omitempty"`
	// Agent audit fields (agents-that-act.md section 7); hashes only, never raw prompts.
	AgentSession    string `json:"agent_session,omitempty"`
	ModelAlias      string `json:"model_alias,omitempty"`
	PromptSHA256    string `json:"prompt_sha256,omitempty"`
	ToolInputSHA256 string `json:"tool_input_sha256,omitempty"`
}

// ConfirmRequest is the svc.policy.{command,task}.confirm body (asset-command.md section 6.2).
// The confirmer's identity comes from the connection, never from the body.
type ConfirmRequest struct {
	V          int    `json:"v"`
	ActionID   string `json:"action_id"`
	ParamsHash string `json:"params_hash"`
	Decision   string `json:"decision"` // approve | reject
	Reason     string `json:"reason,omitempty"`
}

// Denial is {code, message}; codes are asset-command.md section 6.4.
type Denial struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Denial and outcome codes (asset-command.md section 6.4).
const (
	CodeSchemaInvalid         = "schema_invalid"
	CodeUnknownAction         = "unknown_action"
	CodeForbiddenAction       = "forbidden_action"
	CodeNotSimulated          = "not_simulated"
	CodeNotCapable            = "not_capable"
	CodeNotAuthorized         = "not_authorized"
	CodePreconditionFailed    = "precondition_failed"
	CodePolicyUnavailable     = "policy_unavailable"
	CodeConflict              = "conflict"
	CodeExpired               = "expired"
	CodeRejectedByHuman       = "rejected_by_human"
	CodeDuplicateConfirmation = "duplicate_confirmation"
	CodeAgentCannotConfirm    = "agent_cannot_confirm"
	CodeReplayed              = "replayed"
	CodeNoAck                 = "no_ack"
)

// Approval is one human confirmation.
type Approval struct {
	By   string `json:"by"`
	At   string `json:"at"`
	Step int    `json:"step"`
}

// PolicyCheck records one authorization question and its answer.
type PolicyCheck struct {
	User     string `json:"user"`
	Relation string `json:"relation"`
	Object   string `json:"object"`
	Allowed  bool   `json:"allowed"`
}

// PolicyInfo is the `policy` audit field.
type PolicyInfo struct {
	Store      string         `json:"store"`
	ModelID    string         `json:"model_id"`
	Checks     []PolicyCheck  `json:"checks"`
	Conditions map[string]any `json:"conditions,omitempty"`
}

// Action states held by the policy service.
const (
	StateAwaiting   = "awaiting_confirmation"
	StateDenied     = "denied"
	StateRejected   = "rejected"
	StateExpired    = "expired"
	StateApproved   = "approved"
	StateDispatched = "dispatched"
	StateExecuted   = "executed"
	StateFailed     = "failed"
	StateRefused    = "refused"
	StateAcked      = "acknowledged"
)

// Action is the policy service's record of one action.
type Action struct {
	ActionID        string          `json:"action_id"`
	Family          Family          `json:"family"`
	RequestID       string          `json:"request_id"`
	ActionType      string          `json:"action_type"`
	ActionVersion   int             `json:"action_version"`
	Level           Level           `json:"level"`
	Target          Target          `json:"target"`
	Params          json.RawMessage `json:"params"`
	ParamsHash      string          `json:"params_hash"`
	RequestedBy     Actor           `json:"requested_by"`
	ActorVerified   bool            `json:"actor_verified"`
	Reason          string          `json:"reason,omitempty"`
	State           string          `json:"status"`
	ApprovalsNeeded int             `json:"approvals_needed"`
	Approvals       []Approval      `json:"approvals"`
	ConfirmBy       time.Time       `json:"confirm_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	CommandID       string          `json:"command_id,omitempty"`
	CommandSubject  string          `json:"command_subject,omitempty"`
	TaskID          string          `json:"task_id,omitempty"`
	Denial          *Denial         `json:"denial,omitempty"`
	Policy          PolicyInfo      `json:"policy"`
	Corr            string          `json:"corr,omitempty"`
	Simulated       bool            `json:"simulated"`
	Agent           *ActionRequest  `json:"-"`
	ackDeadline     time.Time
	acked           bool
}

func (a *Action) terminal() bool {
	switch a.State {
	case StateDenied, StateRejected, StateExpired, StateExecuted, StateFailed, StateRefused:
		return true
	}
	return false
}

// ProposeReply is the propose reply (asset-command.md section 6.1).
type ProposeReply struct {
	ActionID        string  `json:"action_id"`
	Status          string  `json:"status"`
	Level           Level   `json:"level,omitempty"`
	ParamsHash      string  `json:"params_hash,omitempty"`
	ConfirmBy       string  `json:"confirm_by,omitempty"`
	ApprovalsNeeded int     `json:"approvals_needed"`
	Denial          *Denial `json:"denial"`
	DryRun          bool    `json:"dry_run,omitempty"`
}

// ConfirmReply is the confirm reply.
type ConfirmReply struct {
	ActionID          string  `json:"action_id"`
	Status            string  `json:"status"`
	ApprovalsNeeded   int     `json:"approvals_needed"`
	ApprovalsReceived int     `json:"approvals_received"`
	CommandID         string  `json:"command_id,omitempty"`
	Denial            *Denial `json:"denial"`
}
