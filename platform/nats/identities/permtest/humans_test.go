package permtest

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

// Zone commands (zone.create/update/remove) close on the geofence service's result
// evt.map.geofence.{created|updated|deleted|refused} (four tokens, D1-P2). onemind-xr E2E run 1:
// without a subscribe grant the policy service's subscription was refused and every zone action
// ended no_ack. The grant is evt.map.geofence.> (crossings included; see identities.yaml).
func TestPolicyReceivesZoneResults(t *testing.T) {
	h := startHub(t, "B")
	provision(t, h)
	policy := h.conn(t, "policy")
	geofence := h.conn(t, "geofence")
	for _, ev := range []string{"created", "updated", "deleted", "refused"} {
		expectDelivered(t, policy, "evt.map.geofence."+ev, geofence, `{"v":1,"data":{"zone_id":"home","action_id":"a1"}}`)
	}
	// The exact subscription the policy service makes (service.go: "evt.map.geofence.*").
	expectDelivered(t, policy, "evt.map.geofence.*", geofence, `{}`)
	expectSub(t, policy, "evt.map.>", false)
	expectSub(t, policy, "evt.so.alert.>", false)
}

// A4 actions (asset.arm, asset.takeoff) need two distinct humans (asset-command.md). The test
// config carries two CLIENT humans, zeus and the test_only alice; both reach the policy confirm
// endpoints and the policy service sees two different server-set identities, so it can tell the
// second approver from the first.
func TestTwoDistinctHumans(t *testing.T) {
	h := startHub(t, "B")
	policy := h.conn(t, "policy")
	if _, err := policy.nc.Subscribe("svc.policy.>", func(m *nats.Msg) {
		_ = m.Respond([]byte(m.Header.Get("Nats-Request-Info")))
	}); err != nil {
		t.Fatal(err)
	}
	policy.nc.Flush()
	seen := map[string]bool{}
	for _, u := range []string{"zeus", "alice"} {
		c := h.conn(t, u)
		for _, subj := range []string{"svc.policy.command.confirm", "svc.policy.task.confirm"} {
			req := nats.NewMsg(subj)
			req.Header.Set("Nats-Request-Info", `{"acc":"CLIENT","user":"zeus"}`) // forged; the server replaces it
			req.Data = []byte(`{}`)
			r, err := c.nc.RequestMsg(req, 2*time.Second)
			if err != nil {
				t.Fatalf("%s %s: %v", u, subj, err)
			}
			var info struct{ Acc, User string }
			_ = json.Unmarshal(r.Data, &info)
			if info.Acc != "CLIENT" || info.User != u {
				t.Fatalf("%s %s: policy saw %q, want CLIENT/%s", u, subj, r.Data, u)
			}
			seen[info.User] = true
			pass(t, "ALLOW req  %-6s -> %-28s policy saw CLIENT/%s", u, subj, info.User)
		}
	}
	if len(seen) != 2 {
		t.Fatalf("want two distinct human identities, saw %v", seen)
	}
	pass(t, "two distinct CLIENT humans confirm with separate server-set identities: %v", seen)
}

// A local or stand-in Hermes (scripted agent in tests, no model) connects the way the pod does:
// through its own sidecar leaf rendered by `leaf-conf`, never as a STANDARD client (the hermes
// credential is LEAFNODE only). `--remote/--plain/--port/--http-port` point the sidecar at a
// disposable hub; the agent then talks to the sidecar on loopback without credentials.
func TestHermesStandInSidecar(t *testing.T) {
	h := startHub(t, "B")
	if _, err := dial(t, h.url(), "hermes"); err == nil {
		t.Fatal("hermes connected as a STANDARD client")
	} else {
		pass(t, "DENY  conn hermes as STANDARD client: %v", err)
	}
	port, httpPort := freePort(t), freePort(t)
	conf := filepath.Join(workDir, fmt.Sprintf("hermes-standin-%d.conf", time.Now().UnixNano()))
	t.Cleanup(func() { os.Remove(conf) })
	out, err := exec.Command(python(), filepath.Join(repoRoot, "scripts/nats_identities.py"), "leaf-conf", "--user", "hermes",
		"--passwords", filepath.Join(workDir, "passwords.json"), "--out", conf,
		"--remote", fmt.Sprintf("127.0.0.1:%d", h.leafPort), "--plain",
		"--port", fmt.Sprint(port), "--http-port", fmt.Sprint(httpPort)).CombinedOutput()
	if err != nil {
		t.Fatalf("leaf-conf: %v: %s", err, out)
	}
	lg := &logBuf{}
	srv, err := startServer(conf, httpPort, lg)
	if err != nil {
		t.Fatalf("stand-in sidecar: %v", err)
	}
	t.Cleanup(srv.Shutdown)
	if !waitLeafs(t, h.srv, 1, 5*time.Second) {
		t.Fatalf("stand-in sidecar did not connect: %v", h.log.grep("hermes"))
	}
	l := &leaf{srv: srv, port: port, log: lg}
	agent := l.conn(t, "hermes-standin")
	onemind := h.conn(t, "onemind")
	expectLeafPub(t, agent, "agents.hb.hermes.zeus.legacy", true, onemind)
	agent.nc.Publish("cmd.fleet.drone.sim-x1.goto", []byte(`{}`))
	agent.nc.Flush()
	if e, ok := agent.sink.violation("publish", "cmd.fleet.drone.sim-x1.goto", 2*time.Second); !ok {
		t.Fatal("stand-in hermes cmd.fleet publish: expected a Permissions Violation at the sidecar")
	} else if !strings.Contains(e, "cmd.fleet") {
		t.Fatalf("unexpected violation %q", e)
	} else {
		pass(t, "DENY  stand-in hermes cmd.fleet.drone.sim-x1.goto refused at its sidecar: %s", e)
	}
	pass(t, "stand-in hermes connected through a leaf-conf sidecar (account AGENTS) on 127.0.0.1:%d", port)
}
