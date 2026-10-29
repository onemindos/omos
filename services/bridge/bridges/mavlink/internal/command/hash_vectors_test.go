package command_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/command"
	ct "github.com/onemindos/omos-bridge/bridges/mavlink/internal/commandtest"
	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/jcs"
)

// The cross-implementation vector file shared with the policy service (omos #554) and the
// ROS 2 bridge (omos #566): same input, same hash (asset-command.md section 5).
const vectorFile = "../../../../../schemas/command/testdata/params-hash-vectors.json"

type hashVector struct {
	Name          string `json:"name"`
	ActionType    string `json:"action_type"`
	ActionVersion int    `json:"action_version"`
	Target        struct {
		EntityType string `json:"entity_type"`
		EntityID   string `json:"entity_id"`
	} `json:"target"`
	ParamsJSON string `json:"params_json"`
	Canonical  string `json:"canonical"`
	ParamsHash string `json:"params_hash"`
}

func loadVectors(t *testing.T) map[string]hashVector {
	t.Helper()
	b, err := os.ReadFile(vectorFile)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Vectors []hashVector `json:"vectors"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	out := map[string]hashVector{}
	for _, v := range doc.Vectors {
		out[v.Name] = v
	}
	if len(out) < 10 {
		t.Fatalf("only %d vectors", len(out))
	}
	return out
}

func TestParamsHashVectors(t *testing.T) {
	for _, v := range loadVectors(t) {
		h, err := command.ParamsHash(v.ActionType, v.ActionVersion, v.Target.EntityType, v.Target.EntityID, json.RawMessage(v.ParamsJSON))
		if err != nil || h != v.ParamsHash {
			t.Errorf("%s: got %s %v, want %s", v.Name, h, err, v.ParamsHash)
		}
		p := v.ParamsJSON
		if p == "" {
			p = "{}"
		}
		c, err := jcs.Marshal(map[string]any{"action_type": v.ActionType, "action_version": v.ActionVersion,
			"target": map[string]string{"entity_type": v.Target.EntityType, "entity_id": v.Target.EntityID},
			"params": json.RawMessage(p)})
		if err != nil || string(c) != v.Canonical {
			t.Errorf("%s: canonical %s %v, want %s", v.Name, c, err, v.Canonical)
		}
	}
}

// D2-P1: a task-bearing command carries the task.assign (or task.abort) hash, which the
// bridge cannot recompute; it is accepted when signed and well formed. Without task_id the
// asset.{verb} hash is still recomputed and a task hash is refused.
func TestTaskBoundParamsHash(t *testing.T) {
	vec := loadVectors(t)
	assign, abort := vec["task-assign-goto"], vec["task-abort"]
	pub, priv := ct.NewKey()
	now := time.Now()
	spec := map[string]any{"lat": 0.001, "lon": 0.001, "hae_m": 80.0}
	cases := []struct {
		name string
		o    ct.Opts
		code string
	}{
		{"goto with task.assign hash", ct.Opts{Verb: "goto", Data: spec, TaskID: "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8", ParamsHash: assign.ParamsHash}, ""},
		// task.abort sends cmd.fleet.drone.{id}.hold with data {"reason":"abort"} (mission-task.md §6.2).
		{"hold with task.abort hash", ct.Opts{Verb: "hold", Data: map[string]any{"reason": "abort"}, TaskID: "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8", ParamsHash: abort.ParamsHash}, ""},
		{"rtl with task.abort hash", ct.Opts{Verb: "rtl", TaskID: "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8", ParamsHash: abort.ParamsHash}, ""},
		{"abort hold with bad data", ct.Opts{Verb: "hold", Data: map[string]any{"reason": "abort", "duration_s": -1}, TaskID: "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8", ParamsHash: abort.ParamsHash}, command.CodeSchemaInvalid},
		{"task hash without task_id", ct.Opts{Verb: "goto", Data: spec, ParamsHash: assign.ParamsHash}, command.CodeNotAuthorized},
		{"task_id with malformed hash", ct.Opts{Verb: "goto", Data: spec, TaskID: "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8", ParamsHash: "sha256:XYZ"}, command.CodeSchemaInvalid},
		{"malformed task_id", ct.Opts{Verb: "goto", Data: spec, TaskID: "task 1", ParamsHash: assign.ParamsHash}, command.CodeSchemaInvalid},
		{"task-bound but unsigned (sig required)", ct.Opts{Verb: "goto", Data: spec, TaskID: "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8", ParamsHash: assign.ParamsHash, NoSig: true}, command.CodeBadSignature},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := newVerifier(t, pub)
			c.o.AssetID, c.o.Key, c.o.Now = "sim-x1", priv, now
			subj, body, err := ct.Build(c.o)
			if err != nil {
				t.Fatal(err)
			}
			e, err := command.Parse(body)
			if err != nil {
				t.Fatal(err)
			}
			res, r := v.Verify(subj, e, now)
			got := ""
			if r != nil {
				got = r.Code
			}
			if got != c.code {
				t.Fatalf("code %q want %q (%v)", got, c.code, r)
			}
			if r == nil && (!res.TaskBound || res.ParamsHash != c.o.ParamsHash) {
				t.Fatalf("result %+v", res)
			}
		})
	}
}
