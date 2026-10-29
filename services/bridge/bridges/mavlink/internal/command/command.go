// Package command verifies asset commands on cmd.fleet.{class}.{asset_id}.{verb}
// before the bridge moves a vehicle (onemind-xr docs/contracts/asset-command.md §4 step 6,
// §6.3 and §7). A command is accepted only if every check passes; anything else is
// refused with a §6.4 code and nothing is sent to the vehicle.
package command

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/jcs"
)

// Refusal codes used by the bridge (asset-command.md §6.4).
const (
	CodeSchemaInvalid      = "schema_invalid"
	CodeUnknownAction      = "unknown_action"
	CodeNotSimulated       = "not_simulated"
	CodeNotAuthorized      = "not_authorized"
	CodeAgentCannotConfirm = "agent_cannot_confirm"
	CodeBadSignature       = "bad_signature"
	CodeExpired            = "expired"
	CodeReplayed           = "replayed"
	CodePrecondition       = "precondition_failed"
	CodeVehicleRejected    = "vehicle_rejected"
	CodeNoAck              = "no_ack"
	CodeConflict           = "conflict"
)

// PolicySource is the only accepted command source (asset-command.md §6.3).
const PolicySource = "service:policy"

// Refusal is a verification failure.
type Refusal struct {
	Code    string
	Detail  string
	Message string
}

func (r *Refusal) Error() string { return r.Code + ": " + r.Message }

func refuse(code, format string, a ...any) *Refusal {
	return &Refusal{Code: code, Message: fmt.Sprintf(format, a...)}
}

// Approval is the approval block of the command envelope.
type Approval struct {
	Required   bool           `json:"required"`
	Status     string         `json:"status"`
	Approver   string         `json:"approver"`
	Reason     string         `json:"reason"`
	Level      string         `json:"level"`
	ParamsHash string         `json:"params_hash"`
	Approvals  []ApprovalStep `json:"approvals"`
}

// ApprovalStep is one human confirmation.
type ApprovalStep struct {
	By   string `json:"by"`
	At   string `json:"at"`
	Step int    `json:"step"`
}

// Policy is the policy decision block.
type Policy struct {
	DecisionID string `json:"decision_id"`
	Allowed    bool   `json:"allowed"`
}

// Sig is the detached ed25519 signature over the JCS form of the envelope without sig.
type Sig struct {
	Alg   string `json:"alg"`
	Kid   string `json:"kid"`
	Value string `json:"value"`
}

// Envelope is the omos command envelope v1 plus the OneMind fields of asset-command.md §6.3.
type Envelope struct {
	V             int             `json:"v"`
	ID            string          `json:"id"`
	TS            string          `json:"ts"`
	Cmd           string          `json:"cmd"`
	Source        string          `json:"source"`
	Subject       string          `json:"subject"`
	Data          json.RawMessage `json:"data"`
	ActionID      string          `json:"action_id"`
	ActionVersion *int            `json:"action_version,omitempty"`
	ExpiresAt     string          `json:"expires_at"`
	Simulated     *bool           `json:"simulated"`
	Approval      *Approval       `json:"approval"`
	Policy        *Policy         `json:"policy"`
	TaskID        string          `json:"task_id,omitempty"`
	Corr          string          `json:"corr,omitempty"`
	Causation     string          `json:"causation,omitempty"`
	Trace         string          `json:"trace,omitempty"`
	Sig           *Sig            `json:"sig,omitempty"`
	DryRun        bool            `json:"dry_run,omitempty"`

	raw map[string]json.RawMessage
}

// Parse decodes a command envelope and keeps the raw members for signature checks.
func Parse(data []byte) (*Envelope, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	var e Envelope
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	e.raw = raw
	return &e, nil
}

// Verb describes one accepted asset verb: its action level and the number of distinct
// human approvals it needs (asset-command.md §3; A3 = one human, A4 = two).
type Verb struct {
	Name      string
	Level     string
	Approvals int
}

// Verbs is the catalogue this bridge executes. Any other verb is refused.
var Verbs = map[string]Verb{
	"goto":    {"goto", "A3", 1},
	"rtl":     {"rtl", "A3", 1},
	"land":    {"land", "A3", 1},
	"disarm":  {"disarm", "A3", 1},
	"estop":   {"estop", "A3", 1},
	"arm":     {"arm", "A4", 2},
	"takeoff": {"takeoff", "A4", 2},
}

func levelRank(l string) int {
	if len(l) == 2 && l[0] == 'A' && l[1] >= '0' && l[1] <= '4' {
		return int(l[1] - '0')
	}
	return -1
}

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// Config controls verification.
type Config struct {
	AssetClass  string                       // "drone"
	AssetID     string                       // must start "sim-"
	Keys        map[string]ed25519.PublicKey // kid -> policy service public key
	SigRequired bool
	MaxTTL      time.Duration // largest accepted expires_at - ts
	ClockSkew   time.Duration
}

// Verifier checks commands and remembers command ids (nonces) until they expire.
type Verifier struct {
	cfg  Config
	mu   sync.Mutex
	seen map[string]time.Time
}

// NewVerifier returns a verifier. It refuses a configuration that targets a non-sim asset.
func NewVerifier(cfg Config) (*Verifier, error) {
	if !strings.HasPrefix(cfg.AssetID, "sim-") {
		return nil, fmt.Errorf("asset id %q does not start with sim- (decision 0004 #13)", cfg.AssetID)
	}
	if cfg.SigRequired && len(cfg.Keys) == 0 {
		return nil, fmt.Errorf("signatures required but no policy public keys configured")
	}
	if cfg.MaxTTL <= 0 {
		cfg.MaxTTL = 60 * time.Second
	}
	if cfg.ClockSkew <= 0 {
		cfg.ClockSkew = 5 * time.Second
	}
	return &Verifier{cfg: cfg, seen: map[string]time.Time{}}, nil
}

// Result is a verified command.
type Result struct {
	Verb       Verb
	Params     Params
	Signed     bool // signature present and valid
	ParamsHash string
}

// Verify runs every check of asset-command.md §4 step 6 and §7 rules 1-2 on a command
// received on subject. It does not run the geofence or link checks, which need vehicle state.
func (v *Verifier) Verify(subject string, e *Envelope, now time.Time) (*Result, *Refusal) {
	toks := strings.Split(subject, ".")
	if len(toks) != 5 || toks[0] != "cmd" || toks[1] != "fleet" || toks[2] != v.cfg.AssetClass {
		return nil, refuse(CodeSchemaInvalid, "subject %q is not cmd.fleet.%s.{id}.{verb}", subject, v.cfg.AssetClass)
	}
	id, verbName := toks[3], toks[4]

	// Rule 1: simulated assets only.
	if !strings.HasPrefix(id, "sim-") {
		return nil, refuse(CodeNotSimulated, "target %q is not a sim- asset", id)
	}
	if id != v.cfg.AssetID {
		return nil, refuse(CodeSchemaInvalid, "target %q is not this bridge's asset", id)
	}
	if e.Simulated == nil || !*e.Simulated {
		return nil, refuse(CodeNotSimulated, "simulated must be true")
	}
	verb, ok := Verbs[verbName]
	if !ok {
		return nil, refuse(CodeUnknownAction, "verb %q is not in the bridge catalogue", verbName)
	}

	// Envelope shape (omos command envelope v1 required fields + §6.3).
	if e.V < 1 {
		return nil, refuse(CodeSchemaInvalid, "v must be >= 1")
	}
	if !uuidRe.MatchString(e.ID) {
		return nil, refuse(CodeSchemaInvalid, "id must be a lowercase UUID")
	}
	if e.Cmd != subject || e.Subject != subject {
		return nil, refuse(CodeSchemaInvalid, "cmd and subject must equal the delivery subject")
	}
	if e.Source != PolicySource {
		return nil, refuse(CodeNotAuthorized, "source %q is not %s", e.Source, PolicySource)
	}
	if !strings.HasPrefix(e.ActionID, "act-") {
		return nil, refuse(CodeSchemaInvalid, "action_id missing")
	}
	ts, err := time.Parse(time.RFC3339Nano, e.TS)
	if err != nil {
		return nil, refuse(CodeSchemaInvalid, "ts: %v", err)
	}
	exp, err := time.Parse(time.RFC3339Nano, e.ExpiresAt)
	if err != nil {
		return nil, refuse(CodeSchemaInvalid, "expires_at: %v", err)
	}

	// Rule 2: signature.
	signed := false
	if e.Sig != nil {
		if r := v.checkSig(e); r != nil {
			return nil, r
		}
		signed = true
	} else if v.cfg.SigRequired {
		return nil, refuse(CodeBadSignature, "unsigned command")
	}

	// Rule 2: expiry and freshness.
	if !now.Before(exp) {
		return nil, refuse(CodeExpired, "expired at %s", e.ExpiresAt)
	}
	if ts.After(now.Add(v.cfg.ClockSkew)) {
		return nil, refuse(CodeExpired, "ts %s is in the future", e.TS)
	}
	if exp.Sub(ts) > v.cfg.MaxTTL || !exp.After(ts) {
		return nil, refuse(CodeExpired, "expires_at - ts must be in (0, %s]", v.cfg.MaxTTL)
	}

	// Policy decision and human approval.
	if e.Policy == nil || !e.Policy.Allowed || e.Policy.DecisionID == "" {
		return nil, refuse(CodeNotAuthorized, "no allowing policy decision")
	}
	a := e.Approval
	if a == nil || !a.Required || a.Status != "approved" {
		return nil, refuse(CodeNotAuthorized, "command is not approved")
	}
	if levelRank(a.Level) < levelRank(verb.Level) {
		return nil, refuse(CodeNotAuthorized, "level %q below %s required for %s", a.Level, verb.Level, verb.Name)
	}
	need := verb.Approvals
	if levelRank(a.Level) >= 4 && need < 2 {
		need = 2
	}
	humans := map[string]bool{}
	steps := map[int]bool{}
	approverListed := false
	for _, s := range a.Approvals {
		if !strings.HasPrefix(s.By, "user:") {
			return nil, refuse(CodeAgentCannotConfirm, "approval by non-human principal %q", s.By)
		}
		at, err := time.Parse(time.RFC3339Nano, s.At)
		if err != nil || at.After(ts.Add(v.cfg.ClockSkew)) {
			return nil, refuse(CodeNotAuthorized, "approval by %s has a bad time", s.By)
		}
		if humans[s.By] || steps[s.Step] {
			return nil, refuse(CodeNotAuthorized, "duplicate approval step or approver %s", s.By)
		}
		humans[s.By], steps[s.Step] = true, true
		if s.By == a.Approver {
			approverListed = true
		}
	}
	if len(humans) < need {
		return nil, refuse(CodeNotAuthorized, "%s needs %d distinct human approvals, has %d", verb.Name, need, len(humans))
	}
	if !approverListed {
		return nil, refuse(CodeNotAuthorized, "approver %q not in approvals", a.Approver)
	}

	// Parameters and the approved hash (approvals are bound to params_hash, §4).
	params, r := ParseParams(verb.Name, e.Data)
	if r != nil {
		return nil, r
	}
	av := 1
	if e.ActionVersion != nil {
		av = *e.ActionVersion
	}
	hash, err := ParamsHash("asset."+verb.Name, av, v.cfg.AssetClass, id, e.Data)
	if err != nil {
		return nil, refuse(CodeSchemaInvalid, "params_hash: %v", err)
	}
	if a.ParamsHash != hash {
		return nil, refuse(CodeNotAuthorized, "params_hash does not match the approved parameters")
	}

	// Nonce: refuse a command id seen before; remember it until it expires.
	if r := v.remember(e.ID, exp, now); r != nil {
		return nil, r
	}
	return &Result{Verb: verb, Params: params, Signed: signed, ParamsHash: hash}, nil
}

func (v *Verifier) checkSig(e *Envelope) *Refusal {
	if e.Sig.Alg != "ed25519" {
		return refuse(CodeBadSignature, "alg %q not supported", e.Sig.Alg)
	}
	key, ok := v.cfg.Keys[e.Sig.Kid]
	if !ok {
		return refuse(CodeBadSignature, "unknown kid %q", e.Sig.Kid)
	}
	sig, err := base64.StdEncoding.DecodeString(e.Sig.Value)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return refuse(CodeBadSignature, "malformed signature")
	}
	msg, err := signingBytes(e.raw)
	if err != nil {
		return refuse(CodeBadSignature, "canonicalize: %v", err)
	}
	if !ed25519.Verify(key, msg, sig) {
		return refuse(CodeBadSignature, "signature does not verify")
	}
	return nil
}

func (v *Verifier) remember(id string, exp, now time.Time) *Refusal {
	v.mu.Lock()
	defer v.mu.Unlock()
	for k, t := range v.seen {
		if now.After(t) {
			delete(v.seen, k)
		}
	}
	if _, dup := v.seen[id]; dup {
		return refuse(CodeReplayed, "command id %s already received", id)
	}
	v.seen[id] = exp
	return nil
}

// signingBytes is the JCS form of the envelope members without "sig".
func signingBytes(raw map[string]json.RawMessage) ([]byte, error) {
	m := make(map[string]json.RawMessage, len(raw))
	for k, val := range raw {
		if k != "sig" {
			m[k] = val
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jcs.Canonicalize(b)
}

// Sign returns envelope JSON (an object without "sig") with an ed25519 sig member added.
// The policy service owns the private key; the bridge only uses this in tests.
func Sign(envelope []byte, kid string, key ed25519.PrivateKey) ([]byte, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(envelope, &raw); err != nil {
		return nil, err
	}
	msg, err := signingBytes(raw)
	if err != nil {
		return nil, err
	}
	s, _ := json.Marshal(Sig{Alg: "ed25519", Kid: kid, Value: base64.StdEncoding.EncodeToString(ed25519.Sign(key, msg))})
	raw["sig"] = s
	return json.Marshal(raw)
}

// ParamsHash is "sha256:" + hex SHA-256 of the JCS form of
// {action_type, action_version, target, params} (asset-command.md §5).
func ParamsHash(actionType string, actionVersion int, class, id string, params json.RawMessage) (string, error) {
	if len(params) == 0 {
		params = json.RawMessage(`{}`)
	}
	b, err := jcs.Marshal(map[string]any{
		"action_type":    actionType,
		"action_version": actionVersion,
		"target":         map[string]string{"entity_type": class, "entity_id": id},
		"params":         params,
	})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// Params are the verb parameters (asset-command.md §3 data column).
type Params struct {
	Lat      *float64 `json:"lat,omitempty"`
	Lon      *float64 `json:"lon,omitempty"`
	HaeM     *float64 `json:"hae_m,omitempty"`
	SpeedMps *float64 `json:"speed_mps,omitempty"`
	YawDeg   *float64 `json:"yaw_deg,omitempty"`
}

func finite(p *float64) bool { return p != nil && !math.IsNaN(*p) && !math.IsInf(*p, 0) }

// ParseParams validates data for a verb. Unknown members are refused.
func ParseParams(verb string, data json.RawMessage) (Params, *Refusal) {
	var p Params
	if len(data) == 0 || string(data) == "null" {
		return p, refuse(CodeSchemaInvalid, "data must be an object")
	}
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return p, refuse(CodeSchemaInvalid, "data must be an object: %v", err)
	}
	allowed := map[string][]string{
		"goto":    {"lat", "lon", "hae_m", "speed_mps", "yaw_deg"},
		"land":    {"lat", "lon"},
		"takeoff": {"hae_m"},
	}[verb]
	for k := range probe {
		ok := false
		for _, a := range allowed {
			ok = ok || a == k
		}
		if !ok {
			return p, refuse(CodeSchemaInvalid, "unexpected parameter %q for %s", k, verb)
		}
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return p, refuse(CodeSchemaInvalid, "data: %v", err)
	}
	latlon := func() *Refusal {
		if !finite(p.Lat) || !finite(p.Lon) || *p.Lat < -90 || *p.Lat > 90 || *p.Lon < -180 || *p.Lon > 180 {
			return refuse(CodeSchemaInvalid, "lat/lon out of range")
		}
		return nil
	}
	switch verb {
	case "goto":
		if r := latlon(); r != nil {
			return p, r
		}
		if !finite(p.HaeM) {
			return p, refuse(CodeSchemaInvalid, "hae_m required")
		}
		if p.SpeedMps != nil && (!finite(p.SpeedMps) || *p.SpeedMps <= 0 || *p.SpeedMps > 50) {
			return p, refuse(CodeSchemaInvalid, "speed_mps must be in (0, 50]")
		}
		if p.YawDeg != nil && (!finite(p.YawDeg) || *p.YawDeg < 0 || *p.YawDeg >= 360) {
			return p, refuse(CodeSchemaInvalid, "yaw_deg must be in [0, 360)")
		}
	case "land":
		if (p.Lat == nil) != (p.Lon == nil) {
			return p, refuse(CodeSchemaInvalid, "land needs both lat and lon or neither")
		}
		if p.Lat != nil {
			if r := latlon(); r != nil {
				return p, r
			}
		}
	case "takeoff":
		if !finite(p.HaeM) {
			return p, refuse(CodeSchemaInvalid, "hae_m required")
		}
	}
	return p, nil
}
