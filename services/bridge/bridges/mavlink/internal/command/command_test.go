package command_test

import (
	"crypto/ed25519"
	"strings"
	"testing"
	"time"

	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/command"
	ct "github.com/onemindos/omos-bridge/bridges/mavlink/internal/commandtest"
)

func newVerifier(t *testing.T, pub ed25519.PublicKey) *command.Verifier {
	t.Helper()
	v, err := command.NewVerifier(command.Config{
		AssetClass: "drone", AssetID: "sim-x1", SigRequired: true,
		Keys: map[string]ed25519.PublicKey{"policy-test": pub},
	})
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestNewVerifierRefusesNonSim(t *testing.T) {
	if _, err := command.NewVerifier(command.Config{AssetClass: "drone", AssetID: "x1"}); err == nil {
		t.Fatal("expected refusal for non-sim asset id")
	}
	if _, err := command.NewVerifier(command.Config{AssetClass: "drone", AssetID: "sim-x1", SigRequired: true}); err == nil {
		t.Fatal("expected refusal when signatures are required without keys")
	}
}

func TestVerify(t *testing.T) {
	pub, priv := ct.NewKey()
	_, otherPriv := ct.NewKey()
	now := time.Now()
	f := false
	goodGoto := map[string]any{"lat": 0.001, "lon": 0.001, "hae_m": 80.0, "speed_mps": 5.0}

	cases := []struct {
		name string
		o    ct.Opts
		code string // "" means accepted
		subj string // override subject
	}{
		{"goto approved", ct.Opts{Verb: "goto", Data: goodGoto}, "", ""},
		{"arm two humans", ct.Opts{Verb: "arm"}, "", ""},
		{"estop", ct.Opts{Verb: "estop"}, "", ""},
		{"land with point", ct.Opts{Verb: "land", Data: map[string]any{"lat": 0.0, "lon": 0.0}}, "", ""},
		{"takeoff", ct.Opts{Verb: "takeoff", Data: map[string]any{"hae_m": 40.0}}, "", ""},
		{"hold", ct.Opts{Verb: "hold"}, "", ""},
		{"hold with duration and reason", ct.Opts{Verb: "hold", Data: map[string]any{"duration_s": 120, "reason": "abort"}}, "", ""},
		{"hold duration not integer", ct.Opts{Verb: "hold", Data: map[string]any{"duration_s": 1.5}}, command.CodeSchemaInvalid, ""},
		{"hold duration zero", ct.Opts{Verb: "hold", Data: map[string]any{"duration_s": 0}}, command.CodeSchemaInvalid, ""},
		{"hold duration too long", ct.Opts{Verb: "hold", Data: map[string]any{"duration_s": 3601}}, command.CodeSchemaInvalid, ""},
		{"hold reason too long", ct.Opts{Verb: "hold", Data: map[string]any{"reason": strings.Repeat("x", 101)}}, command.CodeSchemaInvalid, ""},
		{"hold reason not a string", ct.Opts{Verb: "hold", Data: map[string]any{"reason": 7}}, command.CodeSchemaInvalid, ""},
		{"hold extra param", ct.Opts{Verb: "hold", Data: map[string]any{"lat": 0.0}}, command.CodeSchemaInvalid, ""},
		{"hold agent approver", ct.Opts{Verb: "hold", Approvers: []string{"agent:hermes/zeus/legacy"}}, command.CodeAgentCannotConfirm, ""},
		{"arm one human", ct.Opts{Verb: "arm", Approvers: []string{"user:zeus"}}, command.CodeNotAuthorized, ""},
		{"arm at A3 level", ct.Opts{Verb: "arm", Level: "A3"}, command.CodeNotAuthorized, ""},
		{"pending approval", ct.Opts{Verb: "rtl", Status: "pending"}, command.CodeNotAuthorized, ""},
		{"no approvals", ct.Opts{Verb: "rtl", Approvers: []string{}}, command.CodeNotAuthorized, ""},
		{"agent approver", ct.Opts{Verb: "rtl", Approvers: []string{"agent:hermes/zeus/legacy"}}, command.CodeAgentCannotConfirm, ""},
		{"same human twice", ct.Opts{Verb: "arm", Approvers: []string{"user:zeus", "user:zeus"}}, command.CodeNotAuthorized, ""},
		{"not simulated", ct.Opts{Verb: "rtl", Simulated: &f}, command.CodeNotSimulated, ""},
		{"non-sim subject", ct.Opts{Verb: "rtl"}, command.CodeNotSimulated, "cmd.fleet.drone.x1.rtl"},
		{"other sim asset", ct.Opts{Verb: "rtl"}, command.CodeSchemaInvalid, "cmd.fleet.drone.sim-x2.rtl"},
		{"unknown verb", ct.Opts{Verb: "strike"}, command.CodeUnknownAction, ""},
		{"wrong source", ct.Opts{Verb: "rtl", Source: "agent.legacy"}, command.CodeNotAuthorized, ""},
		{"unsigned", ct.Opts{Verb: "rtl", NoSig: true}, command.CodeBadSignature, ""},
		{"wrong key", ct.Opts{Verb: "rtl", Key: otherPriv}, command.CodeBadSignature, ""},
		{"unknown kid", ct.Opts{Verb: "rtl", Kid: "nope"}, command.CodeBadSignature, ""},
		{"expired", ct.Opts{Verb: "rtl", Now: now.Add(-time.Minute)}, command.CodeExpired, ""},
		{"ttl too long", ct.Opts{Verb: "rtl", TTL: 10 * time.Minute}, command.CodeExpired, ""},
		{"future ts", ct.Opts{Verb: "rtl", Now: now.Add(time.Minute)}, command.CodeExpired, ""},
		{"hash mismatch", ct.Opts{Verb: "goto", Data: goodGoto, HashData: map[string]any{"lat": 1.0, "lon": 1.0, "hae_m": 80.0}}, command.CodeNotAuthorized, ""},
		{"bad params", ct.Opts{Verb: "goto", Data: map[string]any{"lat": 91.0, "lon": 0.0, "hae_m": 10.0}}, command.CodeSchemaInvalid, ""},
		{"extra param", ct.Opts{Verb: "rtl", Data: map[string]any{"payload": "drop"}}, command.CodeSchemaInvalid, ""},
		{"policy denied", ct.Opts{Verb: "rtl", Mutate: func(m map[string]any) {
			m["policy"] = map[string]any{"decision_id": "d", "allowed": false}
		}}, command.CodeNotAuthorized, ""},
		{"tampered after signing", ct.Opts{Verb: "rtl"}, command.CodeBadSignature, "tamper"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := newVerifier(t, pub)
			c.o.AssetID = "sim-x1"
			if c.o.Key == nil {
				c.o.Key = priv
			}
			if c.o.Now.IsZero() {
				c.o.Now = now
			}
			subj, body, err := ct.Build(c.o)
			if err != nil {
				t.Fatal(err)
			}
			if c.subj == "tamper" {
				body = []byte(string(body[:len(body)-1]) + `,"dry_run":true}`)
			} else if c.subj != "" {
				subj = c.subj // delivery subject differs from the envelope's
			}
			e, err := command.Parse(body)
			if err != nil {
				t.Fatal(err)
			}
			_, r := v.Verify(subj, e, now)
			got := ""
			if r != nil {
				got = r.Code
			}
			if got != c.code {
				t.Fatalf("code %q want %q (%v)", got, c.code, r)
			}
		})
	}
}

func TestReplayRefused(t *testing.T) {
	pub, priv := ct.NewKey()
	v := newVerifier(t, pub)
	now := time.Now()
	subj, body, err := ct.Build(ct.Opts{AssetID: "sim-x1", Verb: "rtl", Key: priv, Now: now})
	if err != nil {
		t.Fatal(err)
	}
	e, _ := command.Parse(body)
	if _, r := v.Verify(subj, e, now); r != nil {
		t.Fatalf("first delivery refused: %v", r)
	}
	if _, r := v.Verify(subj, e, now.Add(time.Second)); r == nil || r.Code != command.CodeReplayed {
		t.Fatalf("replay not refused: %v", r)
	}
	if _, r := v.Verify(subj, e, now.Add(time.Minute)); r == nil || r.Code != command.CodeExpired {
		t.Fatalf("late replay should be expired: %v", r)
	}
}

func TestSigOptionalWhenNotRequired(t *testing.T) {
	v, err := command.NewVerifier(command.Config{AssetClass: "drone", AssetID: "sim-x1"})
	if err != nil {
		t.Fatal(err)
	}
	subj, body, _ := ct.Build(ct.Opts{AssetID: "sim-x1", Verb: "rtl", NoSig: true})
	e, _ := command.Parse(body)
	res, r := v.Verify(subj, e, time.Now())
	if r != nil || res.Signed {
		t.Fatalf("unsigned command with sig_required=false: %v %+v", r, res)
	}
	// A present but invalid signature is still refused.
	_, priv := ct.NewKey()
	subj, body, _ = ct.Build(ct.Opts{AssetID: "sim-x1", Verb: "rtl", Key: priv})
	e, _ = command.Parse(body)
	if _, r := v.Verify(subj, e, time.Now()); r == nil || r.Code != command.CodeBadSignature {
		t.Fatalf("unknown-kid signature accepted: %v", r)
	}
}
