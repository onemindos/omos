package registry

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

// TestE2EMediaMTX runs the real MediaMTX (v1.21.3) with the cluster mediamtx.yml, ffmpeg as
// sim-x1's camera, a disposable nats-server and the registry. It proves FEATURES K1 and K2 on the
// build machine: unauthenticated publish and read are refused, a publisher cannot leave its path
// grant, LL-HLS plays for a client identity, recordings appear and roll, and the registry links
// the stream to ent.drone.sim-x1.state. Skipped when mediamtx, ffmpeg or ffprobe is missing.
//
// Only the listen addresses, the record directory and the timings are changed from the cluster
// file (retention 6h/2h -> 8s, segments 60s -> 2s), so the rest of the file is what runs.
// Set K_EVIDENCE_DIR to keep a decoded HLS frame.
func TestE2EMediaMTX(t *testing.T) {
	mtxBin := lookBin(t, "MEDIAMTX_BIN", "mediamtx")
	ffmpeg := lookBin(t, "FFMPEG_BIN", "ffmpeg")
	ffprobe := lookBin(t, "FFPROBE_BIN", "ffprobe")
	nc := startNATS(t)
	createBucket(t, nc)

	dir := t.TempDir()
	src, err := os.ReadFile(repoFile(t, mediaYML))
	if err != nil {
		t.Fatal(err)
	}
	conf := strings.NewReplacer("recordSegmentDuration: 60s", "recordSegmentDuration: 2s",
		"recordDeleteAfter: 6h", "recordDeleteAfter: 8s", "recordDeleteAfter: 2h", "recordDeleteAfter: 8s").Replace(string(src))
	confPath := filepath.Join(dir, "mediamtx.yml")
	if err := os.WriteFile(confPath, []byte(conf), 0o600); err != nil {
		t.Fatal(err)
	}
	media, err := LoadMediaConfig(confPath)
	if err != nil {
		t.Fatal(err)
	}
	pass := map[string]string{}
	env := os.Environ()
	for i, u := range media.Users {
		if u.Pass != "from-secret" {
			continue
		}
		b := make([]byte, 12)
		_, _ = rand.Read(b)
		pass[u.User] = hex.EncodeToString(b)
		env = append(env, fmt.Sprintf("MTX_AUTHINTERNALUSERS_%d_PASS=%s", i, pass[u.User]))
	}
	p := struct{ rtsp, hls, webrtc, ice, srt, api, metrics, playback int }{
		freePort(t), freePort(t), freePort(t), freeUDPPort(t), freeUDPPort(t), freePort(t), freePort(t), freePort(t)}
	rec := filepath.Join(dir, "rec")
	env = append(env,
		fmt.Sprintf("MTX_RTSPADDRESS=127.0.0.1:%d", p.rtsp), fmt.Sprintf("MTX_HLSADDRESS=127.0.0.1:%d", p.hls),
		fmt.Sprintf("MTX_WEBRTCADDRESS=127.0.0.1:%d", p.webrtc), fmt.Sprintf("MTX_WEBRTCLOCALUDPADDRESS=127.0.0.1:%d", p.ice),
		fmt.Sprintf("MTX_SRTADDRESS=127.0.0.1:%d", p.srt), fmt.Sprintf("MTX_APIADDRESS=127.0.0.1:%d", p.api),
		fmt.Sprintf("MTX_METRICSADDRESS=127.0.0.1:%d", p.metrics), fmt.Sprintf("MTX_PLAYBACKADDRESS=127.0.0.1:%d", p.playback),
		"MTX_PATHDEFAULTS_RECORDPATH="+rec+"/%path/%Y-%m-%d_%H-%M-%S-%f_%z")
	var mtxLog bytes.Buffer
	mtx := exec.Command(mtxBin, confPath)
	mtx.Env, mtx.Stdout, mtx.Stderr = env, &mtxLog, &mtxLog
	if err := mtx.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = mtx.Process.Kill()
		_ = mtx.Wait()
		if t.Failed() {
			t.Logf("mediamtx log:\n%s", redact(mtxLog.String(), pass))
		}
	})
	waitTCP(t, p.rtsp)

	host := "127.0.0.1"
	rtsp := func(user, path string) string {
		if user == "" {
			return fmt.Sprintf("rtsp://%s:%d/%s", host, p.rtsp, path)
		}
		return fmt.Sprintf("rtsp://%s:%s@%s:%d/%s", user, pass[user], host, p.rtsp, path)
	}
	publish := func(user, path string, secs int) *exec.Cmd {
		return exec.Command(ffmpeg, "-hide_banner", "-loglevel", "error", "-re", "-f", "lavfi", "-i",
			"testsrc2=size=640x360:rate=25", "-t", fmt.Sprint(secs), "-c:v", "libx264", "-preset", "ultrafast",
			"-tune", "zerolatency", "-g", "25", "-pix_fmt", "yuv420p", "-f", "rtsp", "-rtsp_transport", "tcp", rtsp(user, path))
	}
	refused := func(name string, cmd *exec.Cmd, want ...string) {
		t.Helper()
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("%s: accepted, want refused", name)
			return
		}
		for _, w := range want {
			if strings.Contains(string(out), w) {
				t.Logf("PASS %s: refused (%s)", name, w)
				return
			}
		}
		t.Errorf("%s: failed without %v: %s", name, want, redact(string(out), pass))
	}

	// K1: unauthenticated and out-of-grant publishes are refused.
	refused("anonymous RTSP publish to drone/sim-x1", publish("", "drone/sim-x1", 1), "401")
	refused("reader identity onemind-client publishing", publish("onemind-client", "drone/sim-x1", 1), "401")
	refused("mavlink-bridge publishing to configured path cam/gate-north (outside its grant)", publish("mavlink-bridge", "cam/gate-north", 1), "401")
	refused("mavlink-bridge publishing a non-sim drone drone/x47 (no such path)", publish("mavlink-bridge", "drone/x47", 1), "400", "401", "404")

	// Registry against the real MediaMTX API.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := New(ctx, nc, &MediaMTXClient{BaseURL: fmt.Sprintf("http://127.0.0.1:%d", p.api), User: "video-registry", Password: pass["video-registry"]},
		Config{Media: media, PollInterval: 500 * time.Millisecond, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
			Endpoints: Endpoints{Host: host, RTSPPort: p.rtsp, HLSPort: p.hls, WebRTC: p.webrtc, SRTPort: p.srt, Playback: p.playback}})
	if err != nil {
		t.Fatal(err)
	}
	events, stop := collect(t, nc, "evt.video.sim-x1.*")
	defer stop()
	go func() { _ = svc.Run(ctx) }()

	// sim-x1's camera: the MAVLink bridge identity publishes a synthetic test pattern.
	cam := publish("mavlink-bridge", "drone/sim-x1", 26)
	var camOut bytes.Buffer
	cam.Stdout, cam.Stderr = &camOut, &camOut
	if err := cam.Start(); err != nil {
		t.Fatal(err)
	}
	camDone := make(chan error, 1)
	go func() { camDone <- cam.Wait() }()
	started := time.Now()

	evs := waitEvents(events, 8*time.Second)
	if len(evs) == 0 || evs[0].Event != "online" {
		t.Fatalf("no online event: %+v (camera: %s)", evs, camOut.String())
	}
	t.Logf("PASS evt.video.sim-x1.online after %.1fs: %+v", time.Since(started).Seconds(), evs[0])

	// K2: selecting the entity returns its live stream.
	r := get(t, nc, "ent.drone.sim-x1.state")
	if r.Record == nil || !r.Record.Live || r.Record.EntitySubject != "ent.drone.sim-x1.state" || r.Record.Codec != "H264" {
		t.Fatalf("registry entry = %+v", r)
	}
	pretty, _ := json.MarshalIndent(r.Record, "", "  ")
	t.Logf("PASS svc.registry.video.get ent.drone.sim-x1.state ->\n%s", pretty)
	if r.Record.TAKURL != "" {
		t.Errorf("loopback stream got a tak_url %q", r.Record.TAKURL)
	}

	// K1: unauthenticated or unauthorised reads are refused.
	probe := func(url string) ([]byte, error) {
		return exec.Command(ffprobe, "-v", "error", "-rtsp_transport", "tcp", "-rw_timeout", "5000000", "-show_entries", "stream=codec_name,width,height", "-of", "csv=p=0", url).CombinedOutput()
	}
	if out, err := probe(rtsp("", "drone/sim-x1")); err == nil || !strings.Contains(string(out), "401") {
		t.Errorf("anonymous RTSP read: err=%v out=%s", err, out)
	} else {
		t.Log("PASS anonymous RTSP read refused (401)")
	}
	if out, err := probe(rtsp("mavlink-bridge", "drone/sim-x1")); err == nil || !strings.Contains(string(out), "401") {
		t.Errorf("publisher identity reading: err=%v out=%s", err, redact(string(out), pass))
	} else {
		t.Log("PASS publisher identity mavlink-bridge cannot read (401)")
	}
	if out, err := probe(rtsp("tak-viewer", "drone/sim-x1")); err != nil || !strings.Contains(string(out), "h264,640,360") {
		t.Errorf("tak-viewer RTSP read: err=%v out=%s", err, redact(string(out), pass))
	} else {
		t.Logf("PASS tak-viewer RTSP read: %s", strings.TrimSpace(string(out)))
	}
	hlsURL := r.Record.URLs.HLS
	if code := httpStatus(t, hlsURL, "", ""); code != http.StatusUnauthorized {
		t.Errorf("anonymous HLS: HTTP %d, want 401", code)
	} else {
		t.Log("PASS anonymous HLS read refused (401 after the cookie check)")
	}

	// LL-HLS for the client identity: the media playlist carries partial segments.
	player := newPlayer()
	master := player.body(t, hlsURL, "onemind-client", pass["onemind-client"])
	var variant string
	for _, l := range strings.Split(master, "\n") {
		if l != "" && !strings.HasPrefix(l, "#") {
			variant = l
			break
		}
	}
	if variant == "" {
		t.Fatalf("master playlist without a variant:\n%s", master)
	}
	media1 := player.body(t, strings.TrimSuffix(hlsURL, "index.m3u8")+variant, "onemind-client", pass["onemind-client"])
	if !strings.Contains(media1, "#EXT-X-PART:") || !strings.Contains(media1, "#EXT-X-SERVER-CONTROL") {
		t.Errorf("media playlist is not LL-HLS:\n%s", media1)
	} else {
		t.Logf("PASS LL-HLS media playlist (%d EXT-X-PART lines)", strings.Count(media1, "#EXT-X-PART:"))
	}
	// No browser origin is allowed (hlsAllowOrigins: []): a foreign page gets no CORS grant,
	// neither on a preflight nor on an authenticated GET.
	for _, method := range []string{http.MethodOptions, http.MethodGet} {
		req, _ := http.NewRequest(method, hlsURL, nil)
		req.Header.Set("Origin", "https://evil.example")
		if method == http.MethodOptions {
			req.Header.Set("Access-Control-Request-Method", "GET")
			req.Header.Set("Access-Control-Request-Headers", "Authorization")
		} else {
			req.SetBasicAuth("onemind-client", pass["onemind-client"])
		}
		resp, err := player.c.Do(req)
		if err != nil {
			t.Fatalf("%s %s with Origin: %v", method, hlsURL, err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("%s with Origin https://evil.example: Access-Control-Allow-Origin %q, want none", method, got)
		} else {
			t.Logf("PASS %s from a foreign Origin gets no Access-Control-Allow-Origin (HTTP %d)", method, resp.StatusCode)
		}
	}

	frame := filepath.Join(dir, "sim-x1-hls-frame.png")
	authed := strings.Replace(hlsURL, "http://", "http://onemind-client:"+pass["onemind-client"]+"@", 1)
	if out, err := exec.Command(ffmpeg, "-hide_banner", "-loglevel", "error", "-y", "-i", authed, "-frames:v", "1", frame).CombinedOutput(); err != nil {
		t.Errorf("decode a frame over HLS: %v %s", err, redact(string(out), pass))
	} else if st, err := os.Stat(frame); err != nil || st.Size() == 0 {
		t.Errorf("no frame written")
	} else {
		t.Logf("PASS decoded a frame of sim-x1 over HLS from the registry URL (%d bytes)", st.Size())
		if ev := os.Getenv("K_EVIDENCE_DIR"); ev != "" {
			b, _ := os.ReadFile(frame)
			_ = os.WriteFile(filepath.Join(ev, "sim-x1-hls-frame.png"), b, 0o644)
		}
	}

	// Playback API: client yes, anonymous no.
	if code := httpStatus(t, r.Record.URLs.Playback, "", ""); code != http.StatusUnauthorized {
		t.Errorf("anonymous playback list: HTTP %d", code)
	}

	// Recording: segments appear, and with 8 s retention the oldest is deleted while new ones
	// keep coming (the roll).
	segDir := filepath.Join(rec, "drone", "sim-x1")
	var first []string
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		if first = listFiles(segDir); len(first) >= 2 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if len(first) < 2 {
		t.Fatalf("recordings did not appear: %v", first)
	}
	t.Logf("PASS recording segments appeared: %v", first)
	body := httpBody(t, r.Record.URLs.Playback, "onemind-client", pass["onemind-client"])
	t.Logf("PASS playback /list for onemind-client: %s", strings.TrimSpace(body))

	time.Sleep(time.Until(started.Add(22 * time.Second)))
	later := listFiles(segDir)
	if len(later) == 0 || contains(later, first[0]) {
		t.Errorf("recordings did not roll: first=%v later=%v", first, later)
	} else {
		t.Logf("PASS recordings rolled: %s deleted after retention, now %v", first[0], later)
	}

	// The camera stops: offline event and live:false.
	select {
	case <-camDone:
	case <-time.After(15 * time.Second):
		_ = cam.Process.Kill()
	}
	evs = waitEvents(events, 6*time.Second)
	if len(evs) == 0 || evs[len(evs)-1].Event != "offline" {
		t.Errorf("no offline event after the camera stopped: %+v", evs)
	} else {
		t.Logf("PASS evt.video.sim-x1.offline: %+v", evs[len(evs)-1])
	}
	if r := get(t, nc, "drone.sim-x1"); r.Record == nil || r.Record.Live {
		t.Errorf("after stop: %+v", r)
	}
}

func redact(s string, pass map[string]string) string {
	for _, p := range pass {
		s = strings.ReplaceAll(s, p, "<redacted>")
	}
	return s
}

func waitTCP(t *testing.T, port int) {
	t.Helper()
	for i := 0; i < 100; i++ {
		if c, err := (&http.Client{Timeout: 200 * time.Millisecond}).Get(fmt.Sprintf("http://127.0.0.1:%d/", port)); err == nil {
			c.Body.Close()
			return
		} else if !strings.Contains(err.Error(), "refused") {
			return // something answers (RTSP is not HTTP)
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("port %d never opened", port)
}

// waitEvents returns as soon as at least one event arrived, or empty after wait.
func waitEvents(next func(time.Duration) []Event, wait time.Duration) []Event {
	deadline := time.Now().Add(wait)
	for time.Now().Before(deadline) {
		if evs := next(200 * time.Millisecond); len(evs) > 0 {
			return evs
		}
	}
	return nil
}

// player keeps one cookie jar, as a browser does: MediaMTX ties HLS playlists to a session
// cookie set while following its cookie-check redirect.
type player struct{ c *http.Client }

func newPlayer() *player {
	jar, _ := cookiejar.New(nil)
	return &player{c: &http.Client{Jar: jar, Timeout: 10 * time.Second}}
}

func (p *player) body(t *testing.T, url, user, pass string) string {
	t.Helper()
	code, body := p.do(t, url, user, pass)
	if code != http.StatusOK {
		t.Fatalf("GET %s as %s: HTTP %d %s", url, user, code, body)
	}
	return body
}

// httpDo follows MediaMTX's cookie-check redirect with a fresh cookie jar.
func httpDo(t *testing.T, url, user, pass string) (int, string) {
	t.Helper()
	return newPlayer().do(t, url, user, pass)
}

func (p *player) do(t *testing.T, url, user, pass string) (int, string) {
	t.Helper()
	c := p.c
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if user != "" {
		req.SetBasicAuth(user, pass)
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

func httpStatus(t *testing.T, url, user, pass string) int {
	code, _ := httpDo(t, url, user, pass)
	return code
}

func httpBody(t *testing.T, url, user, pass string) string {
	t.Helper()
	code, body := httpDo(t, url, user, pass)
	if code != http.StatusOK {
		t.Fatalf("GET %s as %s: HTTP %d %s", url, user, code, body)
	}
	return body
}

func listFiles(dir string) []string {
	ents, _ := os.ReadDir(dir)
	var out []string
	for _, e := range ents {
		if !e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
