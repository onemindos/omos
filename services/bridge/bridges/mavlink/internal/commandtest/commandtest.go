// Package commandtest builds asset command envelopes the way the policy service will
// (asset-command.md §6.3), for tests only. Keys are generated per test run; none is
// stored in the repository.
package commandtest

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/command"
)

// Opts describes one command. Zero values give a valid, approved, signed A3/A4 command.
type Opts struct {
	AssetID    string
	Verb       string
	Data       any // marshalled to the data member; nil means {}
	Now        time.Time
	TTL        time.Duration
	Level      string   // default: the verb's catalogue level
	Approvers  []string // default: one or two "user:" principals per the verb
	Status     string   // default "approved"
	Simulated  *bool    // default true
	Source     string   // default service:policy
	HashData   any      // data the approval hash is computed over (default Data)
	ID         string   // default random UUID
	TaskID     string
	ParamsHash string // approval.params_hash override (e.g. a task.assign hash); default asset.{verb} over HashData
	NoSig      bool
	Kid        string
	Key        ed25519.PrivateKey
	Mutate     func(map[string]any) // last-chance edit before signing
}

// UUID returns a random version-4 UUID.
func UUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func rfc(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05.000Z07:00") }

// Build returns the JSON envelope and its subject.
func Build(o Opts) (subject string, body []byte, err error) {
	if o.Now.IsZero() {
		o.Now = time.Now()
	}
	if o.TTL == 0 {
		o.TTL = 30 * time.Second
	}
	if o.Data == nil {
		o.Data = map[string]any{}
	}
	if o.HashData == nil {
		o.HashData = o.Data
	}
	v := command.Verbs[o.Verb]
	if o.Level == "" {
		o.Level = v.Level
		if o.Level == "" {
			o.Level = "A3"
		}
	}
	if o.Approvers == nil {
		o.Approvers = []string{"user:zeus"}
		if v.Approvals == 2 {
			o.Approvers = append(o.Approvers, "user:daniel")
		}
	}
	if o.Status == "" {
		o.Status = "approved"
	}
	sim := true
	if o.Simulated == nil {
		o.Simulated = &sim
	}
	if o.Source == "" {
		o.Source = command.PolicySource
	}
	if o.ID == "" {
		o.ID = UUID()
	}
	subject = fmt.Sprintf("cmd.fleet.drone.%s.%s", o.AssetID, o.Verb)
	hashRaw, _ := json.Marshal(o.HashData)
	hash, err := command.ParamsHash("asset."+o.Verb, 1, "drone", o.AssetID, hashRaw)
	if err != nil {
		return "", nil, err
	}
	if o.ParamsHash != "" {
		hash = o.ParamsHash
	}
	approvals := []map[string]any{}
	for i, a := range o.Approvers {
		approvals = append(approvals, map[string]any{"by": a, "at": rfc(o.Now.Add(-time.Second)), "step": i + 1})
	}
	approver := ""
	if len(o.Approvers) > 0 {
		approver = o.Approvers[len(o.Approvers)-1]
	}
	env := map[string]any{
		"v": 1, "id": o.ID, "ts": rfc(o.Now), "cmd": subject, "source": o.Source, "subject": subject,
		"data": o.Data, "action_id": "act-" + UUID()[:8], "expires_at": rfc(o.Now.Add(o.TTL)),
		"simulated": *o.Simulated,
		"approval": map[string]any{
			"required": true, "status": o.Status, "approver": approver, "reason": "",
			"level": o.Level, "params_hash": hash, "approvals": approvals,
		},
		"policy": map[string]any{"decision_id": "dec-" + UUID()[:8], "allowed": true},
		"corr":   UUID(),
	}
	if o.TaskID != "" {
		env["task_id"] = o.TaskID
	}
	if o.Mutate != nil {
		o.Mutate(env)
	}
	body, err = json.Marshal(env)
	if err != nil || o.NoSig {
		return subject, body, err
	}
	kid := o.Kid
	if kid == "" {
		kid = "policy-test"
	}
	body, err = command.Sign(body, kid, o.Key)
	return subject, body, err
}

// NewKey returns a fresh ed25519 key pair.
func NewKey() (ed25519.PublicKey, ed25519.PrivateKey) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return pub, priv
}
