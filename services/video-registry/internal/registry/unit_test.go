package registry

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// repoFile resolves a path from the omos repository root.
func repoFile(t *testing.T, rel string) string {
	t.Helper()
	p := filepath.Join("..", "..", "..", "..", rel)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("%s: %v", rel, err)
	}
	return p
}

const mediaYML = "platform/onetak/media/mediamtx.yml"

func TestParsePath(t *testing.T) {
	ok := map[string]StreamRef{
		"drone/sim-x1":   {StreamID: "drone/sim-x1", Class: "drone", EntityType: "drone", EntityID: "sim-x1"},
		"robot/sim-r2":   {StreamID: "robot/sim-r2", Class: "robot", EntityType: "robot", EntityID: "sim-r2"},
		"cam/gate-north": {StreamID: "cam/gate-north", Class: "cam", EntityType: "camera", EntityID: "gate-north"},
	}
	for in, want := range ok {
		got, err := ParsePath(in)
		if err != nil || got != want {
			t.Errorf("ParsePath(%q) = %+v, %v; want %+v", in, got, err, want)
		}
	}
	for _, bad := range []string{"drone/x47", "drone/sim-X1", "other/x", "drone/sim-x1/extra", "drone", "", "cam/../x", "robot/real-1"} {
		if _, err := ParsePath(bad); err == nil {
			t.Errorf("ParsePath(%q) accepted", bad)
		}
	}
	r, _ := ParsePath("drone/sim-x1")
	if r.Key() != "drone.sim-x1" || r.EntitySubject() != "ent.drone.sim-x1.state" {
		t.Errorf("key %q subject %q", r.Key(), r.EntitySubject())
	}
}

func TestParseEntityRef(t *testing.T) {
	for _, in := range []string{"ent.drone.sim-x1.state", "ent.drone.sim-x1", "drone.sim-x1", "drone/sim-x1", " drone.sim-x1 "} {
		r, err := ParseEntityRef(in)
		if err != nil || r.StreamID != "drone/sim-x1" {
			t.Errorf("ParseEntityRef(%q) = %+v, %v", in, r, err)
		}
	}
	if r, err := ParseEntityRef("ent.camera.gate-north.state"); err != nil || r.StreamID != "cam/gate-north" {
		t.Errorf("camera: %+v %v", r, err)
	}
	for _, bad := range []string{"ent.ship.mmsi1.state", "ent.drone.x47.state", "drone", "ent.drone.sim-x1.position"} {
		if _, err := ParseEntityRef(bad); err == nil {
			t.Errorf("ParseEntityRef(%q) accepted", bad)
		}
	}
}

// TestNoAnonymousAccess is the static half of K1 "an unauthenticated publish is refused": every
// user that can publish, read or play back has a password from the Secret; the only
// password-less user is limited to metrics from inside the pod.
func TestNoAnonymousAccess(t *testing.T) {
	c, err := LoadMediaConfig(repoFile(t, mediaYML))
	if err != nil {
		t.Fatal(err)
	}
	for i, u := range c.Users {
		if u.Pass == "" {
			for _, p := range u.Permissions {
				if p.Action != "metrics" {
					t.Errorf("user %d %q has no password but action %q", i, u.User, p.Action)
				}
			}
			if !reflect.DeepEqual(u.IPs, []string{"127.0.0.1", "::1"}) {
				t.Errorf("password-less user %q not limited to localhost: %v", u.User, u.IPs)
			}
			continue
		}
		if u.Pass != "from-secret" {
			t.Errorf("user %q: password must come from the Secret (pass: from-secret), not the file", u.User)
		}
	}
}

func TestPermissionsAndRecording(t *testing.T) {
	c, err := LoadMediaConfig(repoFile(t, mediaYML))
	if err != nil {
		t.Fatal(err)
	}
	got := c.PermissionsFor("drone/sim-x1")
	want := Permissions{
		Publish:  []string{"edge-home", "mavlink-bridge"},
		Read:     []string{"onemind-client", "tak-viewer", "takml"},
		Playback: []string{"onemind-client"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("drone/sim-x1 permissions = %+v, want %+v", got, want)
	}
	if p := c.PermissionsFor("cam/gate-north"); !reflect.DeepEqual(p.Publish, []string{"edge-home"}) {
		t.Errorf("cam publish = %v", p.Publish)
	}
	if p := c.PermissionsFor("drone/real-1"); len(p.Publish)+len(p.Read)+len(p.Playback) != 0 {
		t.Errorf("non-sim drone must have no grants: %+v", p)
	}
	if r := c.RecordingFor("drone/sim-x1"); !r.Enabled || r.Retention != "168h" || r.Format != "fmp4" {
		t.Errorf("drone recording = %+v", r)
	}
	if r := c.RecordingFor("cam/gate-north"); !r.Enabled || r.Retention != "72h" {
		t.Errorf("cam recording = %+v", r)
	}
}

// TestDeploymentEnvMatchesUsers keeps deployment.yaml's MTX_AUTHINTERNALUSERS_<i>_PASS in step
// with the user order in mediamtx.yml, so no user starts with the placeholder password.
func TestDeploymentEnvMatchesUsers(t *testing.T) {
	c, err := LoadMediaConfig(repoFile(t, mediaYML))
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(repoFile(t, "platform/onetak/media/deployment.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var dep struct {
		Spec struct {
			Template struct {
				Spec struct {
					Containers []struct {
						Name string `yaml:"name"`
						Env  []struct {
							Name      string `yaml:"name"`
							ValueFrom struct {
								SecretKeyRef struct{ Name, Key string } `yaml:"secretKeyRef"`
							} `yaml:"valueFrom"`
						} `yaml:"env"`
					} `yaml:"containers"`
				} `yaml:"spec"`
			} `yaml:"template"`
		} `yaml:"spec"`
	}
	if err := yaml.Unmarshal(b, &dep); err != nil {
		t.Fatal(err)
	}
	env := map[string]string{}
	for _, e := range dep.Spec.Template.Spec.Containers[0].Env {
		if e.ValueFrom.SecretKeyRef.Name == "onemind-media-users" {
			env[e.Name] = e.ValueFrom.SecretKeyRef.Key
		}
	}
	n := 0
	for i, u := range c.Users {
		if u.Pass != "from-secret" {
			continue
		}
		n++
		name := "MTX_AUTHINTERNALUSERS_" + strconv.Itoa(i) + "_PASS"
		if env[name] != u.User {
			t.Errorf("%s -> key %q, want %q (user %d in mediamtx.yml)", name, env[name], u.User, i)
		}
	}
	if n != len(env) {
		t.Errorf("deployment sets %d password env vars, mediamtx.yml has %d secret users", len(env), n)
	}
}

func TestBuildRecord(t *testing.T) {
	c, err := LoadMediaConfig(repoFile(t, mediaYML))
	if err != nil {
		t.Fatal(err)
	}
	ready := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	item := PathItem{Name: "drone/sim-x1", Ready: true, ReadyTime: &ready, Tracks: []string{"H264", "KLV"}}
	item.Source = &struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}{Type: "srtConn"}
	item.Tracks2 = append(item.Tracks2, struct {
		Codec      string         `json:"codec"`
		CodecProps map[string]any `json:"codecProps"`
	}{Codec: "H264", CodecProps: map[string]any{"width": float64(1280), "height": float64(720)}})
	ref, _ := ParsePath(item.Name)
	r := Build(ref, item, c, Endpoints{Host: "media.example.ts.net"}, "hub", ready)
	if r.EntitySubject != "ent.drone.sim-x1.state" || !r.Live || !r.KLV || r.Codec != "H264" || r.Width != 1280 {
		t.Fatalf("record = %+v", r)
	}
	if r.URLs.RTSP != "rtsp://media.example.ts.net:8554/drone/sim-x1" || r.URLs.HLS != "http://media.example.ts.net:8888/drone/sim-x1/index.m3u8" ||
		r.URLs.WHEP != "http://media.example.ts.net:8889/drone/sim-x1/whep" {
		t.Errorf("urls = %+v", r.URLs)
	}
	if r.TAKURL != r.URLs.RTSP || r.Media.URL != r.URLs.RTSP || r.Media.Kind != "video" || r.LatencyClass != "tactical" {
		t.Errorf("tak_url %q media %+v class %q", r.TAKURL, r.Media, r.LatencyClass)
	}
	for _, u := range []string{r.URLs.RTSP, r.URLs.HLS, r.URLs.WHEP, r.URLs.SRT, r.TAKURL} {
		if strings.Contains(u, "@") {
			t.Errorf("credentials in %q", u)
		}
	}
	// A loopback host (the local run) yields reader URLs but never a TAK URL.
	l := Build(ref, item, c, Endpoints{Host: "127.0.0.1", RTSPPort: 28554}, "hub", ready)
	if l.TAKURL != "" || l.URLs.RTSP != "rtsp://127.0.0.1:28554/drone/sim-x1" {
		t.Errorf("loopback: tak_url %q rtsp %q", l.TAKURL, l.URLs.RTSP)
	}
	// A site LAN host is a valid site_url but not a hub tak_url.
	s := Build(ref, item, c, Endpoints{Host: "192.168.1.20", SiteHost: "192.168.1.20"}, "home", ready)
	if s.TAKURL != "" || s.SiteURL != "rtsp://192.168.1.20:8554/drone/sim-x1" {
		t.Errorf("site: tak_url %q site_url %q", s.TAKURL, s.SiteURL)
	}
	if !sameContent(r, func() Record { x := r; x.UpdatedAt = x.UpdatedAt.Add(time.Hour); x.Seq = 9; return x }()) {
		t.Error("bookkeeping fields must not count as a change")
	}
}
