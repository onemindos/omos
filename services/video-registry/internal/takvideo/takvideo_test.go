package takvideo

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestVideoURLHygiene is robotics-cot-profile.md section 9 test 3 at the renderer: an entity whose
// only video URL is rtsp://127.0.0.1:8554/sim-x1 produces no <__video>; a registry entry with a
// tailnet tak_url produces it; no output URL contains @ or a query string.
func TestVideoURLHygiene(t *testing.T) {
	d, err := Detail(Hub, "drone", "sim-x1", "SIM-X1 cam0", "", "rtsp://127.0.0.1:8554/sim-x1")
	if err != nil || d != "" {
		t.Fatalf("loopback-only entity produced %q (err %v)", d, err)
	}
	d, err = Detail(Hub, "drone", "sim-x1", "SIM-X1 cam0", "rtsp://media.example.ts.net:8554/drone/sim-x1", "rtsp://127.0.0.1:8554/sim-x1")
	if err != nil || d == "" {
		t.Fatalf("tailnet tak_url produced nothing (err %v)", err)
	}
	for _, want := range []string{
		`<__video uid="OM.drone.sim-x1.video" url="rtsp://media.example.ts.net:8554/drone/sim-x1">`,
		`<ConnectionEntry uid="OM.drone.sim-x1.video" alias="SIM-X1 cam0" protocol="rtsp" address="media.example.ts.net" port="8554" path="/drone/sim-x1" networkTimeout="5000" bufferTime="-1" rtspReliable="1" ignoreEmbeddedKLV="false" roverPort="-1"></ConnectionEntry>`,
	} {
		if !strings.Contains(d, want) {
			t.Errorf("detail %s\n missing %s", d, want)
		}
	}
	if strings.Contains(d, "@") || strings.Contains(d, "?") {
		t.Errorf("credentials or query in %s", d)
	}
	// A registry URL with credentials or a query is skipped in favour of a clean fallback.
	d, _ = Detail(Hub, "drone", "sim-x1", "a", "rtsp://u:p@media.example.ts.net:8554/drone/sim-x1", "rtsp://media.example.ts.net:8554/drone/sim-x1?token=x")
	if d != "" {
		t.Errorf("dirty URLs produced %q", d)
	}
}

func TestCheck(t *testing.T) {
	bad := map[string]Audience{
		"rtsp://localhost:8554/x":                         Hub,
		"rtsp://[::1]:8554/x":                             Hub,
		"rtsp://169.254.1.2:8554/x":                       Hub,
		"rtsp://10.244.3.7:8554/x":                        Hub, // pod IP
		"rtsp://192.168.1.20:8554/x":                      Hub,
		"rtsp://onemind-media.onetak.svc.cluster.local/x": Hub,
		"rtsp://0.0.0.0:8554/x":                           Site,
		"rtsp://127.0.0.1:8554/x":                         Site,
		"ftp://media.example.ts.net/x":                    Hub,
		"srt://media.example.ts.net:8890?streamid=read:x": Hub,
		"rtsp://media.example.ts.net:8554/x#frag":         Hub,
	}
	for u, a := range bad {
		if Check(u, a) == nil {
			t.Errorf("Check(%q) accepted", u)
		}
	}
	good := map[string]Audience{
		"rtsp://media.example.ts.net:8554/drone/sim-x1": Hub,
		"rtsp://100.101.102.103:8554/drone/sim-x1":      Hub, // tailnet CGNAT range
		"https://media.example.ts.net/x/index.m3u8":     Hub,
		"rtsp://192.168.1.20:8554/drone/sim-x1":         Site,
	}
	for u, a := range good {
		if err := Check(u, a); err != nil {
			t.Errorf("Check(%q): %v", u, err)
		}
	}
}

func TestFeed(t *testing.T) {
	f, err := NewFeed("drone", "sim-x1", "SIM-X1 cam0", "rtsp://media.example.ts.net:8554/drone/sim-x1")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(f)
	if !strings.Contains(string(b), `"uid":"OM.drone.sim-x1.video"`) || !strings.Contains(string(b), `"rtspReliable":1`) {
		t.Errorf("feed %s", b)
	}
	if _, err := NewFeed("drone", "sim-x1", "a", "rtsp://127.0.0.1:8554/sim-x1"); err == nil {
		t.Error("loopback feed accepted")
	}
}
