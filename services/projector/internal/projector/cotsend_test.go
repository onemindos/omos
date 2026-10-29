package projector

// svc.cot.send tests. The request bodies and CoT are the OneMind clients' own shapes:
//   - Unreal client (onemind-xr branch run/m2-client): CotSendRequest::ToJson and the builders
//     BuildMarker / BuildSpi / BuildRoute / BuildGeofence / BuildChat in
//     client/OneMind/Plugins/OneMindCore/Source/OneMindCore/Private/CoT.cpp; the marker is the one
//     the headless client sent in docs/evidence/M2/B/e2e-c1-c2/headless-c2.txt.
//   - L0 web client (onemind-xr branch run/m3-web): geochatRequest in web/src/act/chat.ts and
//     the sender uid `OM-WEB-{user}` of web/src/ui/act.ts selfUid().
//   - An ATAK EUD position report (docs/evidence/M2/B/e2e-c1-c2/run.txt line 44): person data.
//   - The web client's CoT export (docs/evidence/M4/WEB/b11-cot-export.xml): OM.draw.* uids.
//
// The integration tests run a disposable nats-server with the account layout of omos #562
// (platform/nats/identities/identities.yaml): ONEMIND exports svc.cot.send to CLIENT, CLIENT
// imports it with share: true, AGENTS has no import.

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func ct(t time.Time) string { return t.UTC().Format(cotTimeLayout) }

// ueMarker is BuildMarker output (headless-c2.txt) with fresh times.
func ueMarker(uid string, now time.Time) string {
	return fmt.Sprintf(`<event version="2.0" how="h-e" uid="%s" type="a-u-G" time="%s" start="%s" stale="%s" opex="s"><point lat="30.27" lon="-97.74" hae="9999999.0" ce="9999999.0" le="9999999.0"/><detail><contact callsign="OM-C2-MARK"/><remarks>marker from the OneMind client · SIMULATED</remarks></detail></event>`,
		uid, ct(now), ct(now), ct(now.Add(5*time.Minute)))
}

func ueSpi(uid string, now time.Time) string {
	return fmt.Sprintf(`<event version="2.0" how="h-e" uid="%s" type="b-m-p-s-p-i" time="%s" start="%s" stale="%s" opex="s"><point lat="30.27" lon="-97.74" hae="9999999.0" ce="9999999.0" le="9999999.0"/><detail><contact callsign="SPI"/><remarks>SIMULATED</remarks></detail></event>`,
		uid, ct(now), ct(now), ct(now.Add(2*time.Minute)))
}

func ueRoute(uid string, now time.Time) string {
	return fmt.Sprintf(`<event version="2.0" how="h-e" uid="%[1]s" type="b-m-r" time="%[2]s" start="%[2]s" stale="%[3]s" opex="s"><point lat="30.268998" lon="-97.753996" hae="9999999.0" ce="9999999.0" le="9999999.0"/><detail><contact callsign="ROUTE"/><link uid="%[1]s.wp1" type="b-m-p-w" relation="c" point="30.268998,-97.753996"/><link uid="%[1]s.wp2" type="b-m-p-w" relation="c" point="30.269147,-97.751690"/><link uid="%[1]s.wp3" type="b-m-p-w" relation="c" point="30.270797,-97.750389"/><strokeColor value="-1"/><remarks>SIMULATED</remarks></detail></event>`,
		uid, ct(now), ct(now.Add(10*time.Minute)))
}

func ueGeofence(uid string, now time.Time) string {
	return fmt.Sprintf(`<event version="2.0" how="h-e" uid="%s" type="u-d-f" time="%s" start="%s" stale="%s" opex="s"><point lat="30.2664" lon="-97.7519" hae="9999999.0" ce="9999999.0" le="9999999.0"/><detail><contact callsign="FIELD"/><link point="30.2655096,-97.7511011"/><link point="30.2658508,-97.7530941"/><link point="30.2675121,-97.7524397"/><link point="30.2655096,-97.7511011"/><strokeColor value="-65536"/><fillColor value="1090453504"/><strokeWeight value="3.0"/><labels_on value="true"/><remarks>SIMULATED</remarks></detail></event>`,
		uid, ct(now), ct(now), ct(now.Add(10*time.Minute)))
}

// ueChat is BuildChat output for client_id "headless" (onemind_headless.cpp --send-chat).
func ueChat(sender, text string, now time.Time) string {
	uid := "GeoChat." + sender + ".All Chat Rooms.onemind-headless-2-Ab3dEf9h"
	return fmt.Sprintf(`<event version="2.0" how="h-e" uid="%[1]s" type="b-t-f" time="%[2]s" start="%[2]s" stale="%[3]s" opex="s"><point lat="0" lon="0" hae="9999999.0" ce="9999999.0" le="9999999.0"/><detail><contact callsign="OneMind"/><__chat parent="RootContactGroup" groupOwner="false" messageId="onemind-headless-2-Ab3dEf9h" chatroom="All Chat Rooms" id="All Chat Rooms" senderCallsign="OneMind"><chatgrp uid0="%[4]s" uid1="All Chat Rooms" id="All Chat Rooms"/></__chat><link uid="%[4]s" type="a-f-G-U-C" relation="p-p"/><remarks source="BAO.F.OneMind.%[4]s" to="All Chat Rooms" time="%[2]s">%[5]s</remarks></detail></event>`,
		uid, ct(now), ct(now.Add(5*time.Minute)), sender, text)
}

// atakPLI is an EUD's own position report (run.txt line 44).
func atakPLI(uid string, now time.Time) string {
	return fmt.Sprintf(`<event version="2.0" uid="%s" type="a-f-G-U-C" how="h-e" time="%s" start="%s" stale="%s"><point lat="30.2672" lon="-97.7431" hae="150" ce="9.9" le="9999999"/><detail><uid Droid="C1-EUD" /><contact endpoint="127.0.0.1:4242:tcp" callsign="C1-EUD"/><__group name="Cyan" role="Team Member"/><status battery="88"/><takv device="PIXEL" platform="ATAK-CIV" os="34" version="5.6.0.12"/><track speed="1.5" course="90"/></detail></event>`,
		uid, ct(now), ct(now), ct(now.Add(5*time.Minute)))
}

// ueRequest is CotSendRequest::ToJson.
func ueRequest(kind, clientID, cot string, targets []string, tlp string) map[string]any {
	uid := reUID.FindStringSubmatch(cot)
	r := map[string]any{"v": 1, "kind": kind, "targets": targets, "client_id": clientID, "marking": map[string]any{"tlp": tlp}, "cot": cot}
	if uid != nil {
		r["client_ref"] = uid[1]
	}
	return r
}

// webGeochat is geochatRequest(room, text, me, msgUid, target) from web/src/act/chat.ts.
func webGeochat(user, room, roomKind, text, target string) map[string]any {
	self := "OM-WEB-" + user
	return map[string]any{
		"v": 1, "kind": "geochat", "target": target,
		"message_uid": "GeoChat." + self + "." + room + "." + newToken(),
		"room":        map[string]any{"id": room, "name": room, "kind": roomKind},
		"text":        text,
		"sender":      map[string]any{"uid": self, "callsign": strings.ToUpper(user)},
	}
}

func asReq(t *testing.T, m map[string]any) *sendRequest {
	t.Helper()
	b, _ := json.Marshal(m)
	var r sendRequest
	if err := json.Unmarshal(b, &r); err != nil {
		t.Fatal(err)
	}
	return &r
}

// ---- unit: validation and building -----------------------------------------------------------

func TestSendPrepareCotClientShapes(t *testing.T) {
	cfg := Defaults(TargetServer)
	now := time.Now()
	for _, tc := range []struct {
		kind, cot string
	}{
		{"marker", ueMarker("onemind-headless-1-UiErQFEk", now)},
		{"spi", ueSpi("onemind-headless-3-Zz12Yy34", now)},
		{"route", ueRoute("onemind-headless-4-Rr12Ss34", now)},
		{"geofence", ueGeofence("onemind-headless-5-Gg12Hh34", now)},
		{"chat", ueChat("onemind-headless", "hello from the OneMind client", now)},
	} {
		p, e := prepareCot(asReq(t, ueRequest(tc.kind, "headless", tc.cot, []string{"server", "mesh"}, "CLEAR")), &cfg, now)
		if e != nil {
			t.Fatalf("%s refused: %v", tc.kind, e)
		}
		if string(p.xml) != tc.cot || p.tlp != "CLEAR" || len(p.targets) != 2 {
			t.Fatalf("%s: prepared %+v", tc.kind, p)
		}
	}
}

// ueShapes are the M4 Unreal client's B11 drawings (onemind-xr branch run/m4-client,
// OneMindCore Drawing.cpp: ShapeBase adds remarks and <archive/>, AddStyle the colours).
func ueShapes(now time.Time) map[string]string {
	head := func(uid, typ, lat, lon string) string {
		return fmt.Sprintf(`<event version="2.0" how="h-e" uid="%s" type="%s" time="%s" start="%s" stale="%s" opex="s"><point lat="%s" lon="%s" hae="9999999.0" ce="9999999.0" le="9999999.0"/><detail><contact callsign="DRAW"/><remarks>SIMULATED</remarks><archive/>`,
			uid, typ, ct(now), ct(now), ct(now.Add(time.Hour)), lat, lon)
	}
	style := `<strokeColor value="-256"/><strokeWeight value="3"/><labels_on value="true"/>`
	return map[string]string{
		"spot":   fmt.Sprintf(`<event version="2.0" how="h-e" uid="onemind-mac-20-Sp0t0001" type="b-m-p-s-m" time="%s" start="%s" stale="%s" opex="s"><point lat="30.2699" lon="-97.7504" hae="9999999.0" ce="9999999.0" le="9999999.0"/><detail><contact callsign="OP-1"/><remarks>SIMULATED</remarks></detail></event>`, ct(now), ct(now), ct(now.Add(time.Hour))),
		"circle": head("onemind-mac-21-C1rc0001", "u-d-c-c", "30.2753", "-97.7431") + `<shape><ellipse major="200" minor="200" angle="360"/></shape>` + style + `</detail></event>`,
		"line":   head("onemind-mac-22-L1ne0001", "u-d-f", "30.2690", "-97.7540") + `<link point="30.2690,-97.7540"/><link point="30.2708,-97.7504"/>` + style + `</detail></event>`,
		"rect":   head("onemind-mac-23-Rect0001", "u-d-r", "30.2672", "-97.7431") + `<link point="30.2632,-97.7478"/><link point="30.2632,-97.7384"/><link point="30.2712,-97.7384"/><link point="30.2712,-97.7478"/>` + style + `</detail></event>`,
		"ring":   head("onemind-mac-24-R1ng0001.r1", "u-r-b-c-c", "30.2672", "-97.7431") + `<shape><ellipse major="500" minor="500" angle="360"/></shape>` + style + `<fillColor value="0"/></detail></event>`,
		"rb":     head("onemind-mac-25-RngB0001", "u-rb-a", "30.2672", "-97.7431") + `<range value="812.4"/><bearing value="47.1"/><inclination value="0.0"/><rangeUnits value="1"/><bearingUnits value="0"/><northRef value="1"/>` + style + `</detail></event>`,
	}
}

func TestSendPrepareShapesM4(t *testing.T) {
	cfg := Defaults(TargetServer)
	now := time.Now()
	for name, cot := range ueShapes(now) {
		if _, e := prepareCot(asReq(t, ueRequest("shape", "mac", cot, []string{"server"}, "CLEAR")), &cfg, now); e != nil {
			t.Errorf("shape %s refused: %v", name, e)
		}
	}
	// A shape is not a marker atom, and a sector (u-d-f) with one point is refused.
	if _, e := prepareCot(asReq(t, ueRequest("shape", "headless", ueMarker("onemind-headless-1-Abc", now), []string{"server"}, "CLEAR")), &cfg, now); e == nil {
		t.Error("marker atom accepted as shape")
	}
}

func TestSendPrepareCotRefusals(t *testing.T) {
	cfg := Defaults(TargetServer)
	now := time.Now()
	good := ueMarker("onemind-headless-1-UiErQFEk", now)
	cases := []struct {
		name, kind, client, cot, tlp, code, detail string
	}{
		{"OM entity uid", "marker", "headless", strings.Replace(good, "onemind-headless-1-UiErQFEk", "OM.vessel.mmsi-367000001", 1), "CLEAR", ErrForbidden, "reserved"},
		{"OM lower case", "marker", "headless", strings.Replace(good, "onemind-headless-1-UiErQFEk", "om.aircraft.a1b2c3", 1), "CLEAR", ErrForbidden, "reserved"},
		{"INTEL uid", "marker", "headless", strings.Replace(good, "onemind-headless-1-UiErQFEk", "INTEL.acled-1", 1), "CLEAR", ErrForbidden, "reserved"},
		{"OM link", "route", "headless", strings.Replace(ueRoute("onemind-headless-4-Rr", now), `uid="onemind-headless-4-Rr.wp2"`, `uid="OM.drone.sim-x1"`, 1), "CLEAR", ErrForbidden, "reserved"},
		{"OM chat sender", "chat", "headless", strings.ReplaceAll(ueChat("onemind-headless", "x", now), `uid0="onemind-headless"`, `uid0="OM.person.tak-android-1"`), "CLEAR", ErrForbidden, "reserved"},
		{"web export OM.draw", "geofence", "headless", `<event version="2.0" uid="OM.draw.dmvtolg1g1" type="u-d-r" how="h-e" time="` + ct(now) + `" start="` + ct(now) + `" stale="` + ct(now.Add(time.Hour)) + `"><point lat="30.2671999" lon="-97.7430998" hae="9999999.0" ce="9999999.0" le="9999999.0"/><detail><contact callsign="PATROL AOI"/></detail></event>`, "CLEAR", ErrForbidden, "reserved"},
		{"uid of another client", "marker", "headless", strings.Replace(good, "onemind-headless-1-UiErQFEk", "onemind-mac-1-UiErQFEk", 1), "CLEAR", ErrForbidden, "sending client"},
		{"EUD uid takeover", "marker", "headless", strings.Replace(good, "onemind-headless-1-UiErQFEk", "ANDROID-c1e2e0d9f2494", 1), "CLEAR", ErrForbidden, "sending client"},
		{"EUD PLI (person)", "marker", "headless", atakPLI("onemind-headless-9-Pp", now), "CLEAR", ErrForbidden, "person-track"},
		{"contact endpoint", "marker", "headless", strings.Replace(good, `<contact callsign="OM-C2-MARK"/>`, `<contact callsign="X" endpoint="10.0.0.2:4242:tcp"/>`, 1), "CLEAR", ErrForbidden, "person-track"},
		{"track element", "marker", "headless", strings.Replace(good, `</detail>`, `<track speed="1.4" course="90"/></detail>`, 1), "CLEAR", ErrForbidden, "person-track"},
		{"person entity in remarks", "marker", "headless", strings.Replace(good, "marker from the OneMind client", "see ent.person.tak-android-1.state", 1), "CLEAR", ErrForbidden, "person-track"},
		{"flow tags", "marker", "headless", strings.Replace(good, `</detail>`, `<_flow-tags_ omos-tak-projector="2026-10-29T00:00:00Z"/></detail>`, 1), "CLEAR", ErrForbidden, "projector"},
		{"onemind detail", "marker", "headless", strings.Replace(good, `</detail>`, `<onemind entity_id="x"/></detail>`, 1), "CLEAR", ErrForbidden, "projector"},
		{"unknown detail", "marker", "headless", strings.Replace(good, `</detail>`, `<sensor fov="30"/></detail>`, 1), "CLEAR", ErrInvalid, "not allowed"},
		{"access attribute", "marker", "headless", strings.Replace(good, `opex="s"`, `opex="s" access="Unclassified"`, 1), "CLEAR", ErrForbidden, "access"},
		{"hostile", "marker", "headless", strings.Replace(good, `type="a-u-G"`, `type="a-h-G"`, 1), "CLEAR", ErrForbidden, "civilian"},
		{"suspect", "marker", "headless", strings.Replace(good, `type="a-u-G"`, `type="a-s-G"`, 1), "CLEAR", ErrForbidden, "suspect"},
		{"strike tasking", "marker", "headless", strings.Replace(good, `type="a-u-G"`, `type="t-k-d"`, 1), "CLEAR", ErrInvalid, "not allowed"},
		{"delete", "marker", "headless", strings.Replace(good, `type="a-u-G"`, `type="t-x-d-d"`, 1), "CLEAR", ErrInvalid, "not allowed"},
		{"kind mismatch", "spi", "headless", good, "CLEAR", ErrInvalid, "not allowed"},
		{"machine how", "marker", "headless", strings.Replace(good, `how="h-e"`, `how="m-f"`, 1), "CLEAR", ErrForbidden, "human"},
		{"red", "marker", "headless", good, "RED", "", ""}, // marking is checked after prepare (handler)
		{"doctype", "marker", "headless", `<!DOCTYPE x [<!ENTITY a "b">]>` + good, "CLEAR", ErrInvalid, "DOCTYPE"},
		{"namespace", "marker", "headless", strings.Replace(good, `<event `, `<event xmlns:x="urn:x" `, 1), "CLEAR", ErrInvalid, "namespace"},
		{"two roots", "marker", "headless", good + good, "CLEAR", ErrInvalid, "single"},
		{"not xml", "marker", "headless", "hello", "CLEAR", ErrInvalid, "outside"},
		{"old time", "marker", "headless", ueMarker("onemind-headless-1-Old", now.Add(-time.Hour)), "CLEAR", ErrInvalid, "clock"},
		{"stale too far", "marker", "headless", strings.Replace(good, ct(now.Add(5*time.Minute)), ct(now.Add(48*time.Hour)), 1), "CLEAR", ErrInvalid, "ahead"},
		{"bad lat", "marker", "headless", strings.Replace(good, `lat="30.27"`, `lat="130.27"`, 1), "CLEAR", ErrInvalid, "range"},
		{"bad client id", "marker", "Head less", good, "CLEAR", ErrInvalid, "client_id"},
		{"chat sender other client", "chat", "headless", ueChat("onemind-mac", "x", now), "CLEAR", ErrForbidden, "sending client"},
	}
	for _, tc := range cases {
		_, e := prepareCot(asReq(t, ueRequest(tc.kind, tc.client, tc.cot, []string{"server"}, tc.tlp)), &cfg, now)
		if tc.code == "" {
			if e != nil {
				t.Errorf("%s: unexpected refusal %v", tc.name, e)
			}
			continue
		}
		if e == nil || e.code != tc.code || !strings.Contains(e.detail, tc.detail) {
			t.Errorf("%s: want %s containing %q, got %v", tc.name, tc.code, tc.detail, e)
		}
	}
	// No marking at all is refused for client-built CoT (G7 fails closed).
	r := ueRequest("marker", "headless", good, []string{"server"}, "")
	delete(r, "marking")
	if _, e := prepareCot(asReq(t, r), &cfg, now); e == nil || e.code != ErrForbidden {
		t.Errorf("missing marking: want forbidden, got %v", e)
	}
	// Mesh profile: the size limit is MESH_MAX_BYTES.
	mcfg := Defaults(TargetMesh)
	big := strings.Replace(good, "marker from the OneMind client", strings.Repeat("x", 1300), 1)
	if _, e := prepareCot(asReq(t, ueRequest("marker", "headless", big, []string{"mesh"}, "CLEAR")), &mcfg, now); e == nil || e.code != ErrInvalid {
		t.Errorf("mesh size: want invalid, got %v", e)
	}
}

func TestSendPrepareChatWebShape(t *testing.T) {
	cfg := Defaults(TargetServer)
	now := time.Now()
	id := identity{Account: "CLIENT", User: "zeus", Verified: true}
	req := webGeochat("zeus", allChatRooms, "all", "Need eyes on <north field> & the gate", "both")
	p, e := prepareChat(asReq(t, req), id, &cfg, now)
	if e != nil {
		t.Fatal(e)
	}
	x := string(p.xml)
	for _, want := range []string{`type="b-t-f"`, `chatroom="All Chat Rooms"`, `senderCallsign="ZEUS"`, `uid0="OM-WEB-zeus"`,
		`source="BAO.F.OneMind.OM-WEB-zeus"`, `&lt;north field&gt; &amp; the gate</remarks>`, `uid="` + req["message_uid"].(string) + `"`} {
		if !strings.Contains(x, want) {
			t.Fatalf("GeoChat CoT lacks %s:\n%s", want, x)
		}
	}
	if p.tlp != "CLEAR" || len(p.targets) != 2 {
		t.Fatalf("prepared %+v", p)
	}
	// The built CoT passes the same scanner as client CoT, so it is well formed.
	if _, e := scanCot(p.xml, 64); e != nil {
		t.Fatalf("built GeoChat does not scan: %v", e)
	}
	refuse := func(name string, mut func(m map[string]any), id identity, code string) {
		t.Helper()
		m := webGeochat("zeus", "Cyan", "group", "hi", "server")
		mut(m)
		if _, e := prepareChat(asReq(t, m), id, &cfg, now); e == nil || e.code != code {
			t.Errorf("%s: want %s, got %v", name, code, e)
		}
	}
	carol := identity{Account: "CLIENT", User: "carol", Verified: true}
	refuse("sender impersonation", func(map[string]any) {}, carol, ErrForbidden)
	refuse("OM sender", func(m map[string]any) {
		m["sender"] = map[string]any{"uid": "OM.person.zeus", "callsign": "Z"}
	}, id, ErrForbidden)
	refuse("message uid of someone else", func(m map[string]any) { m["message_uid"] = "GeoChat.ANDROID-1.Cyan.x" }, id, ErrForbidden)
	refuse("OM room", func(m map[string]any) { m["room"] = map[string]any{"id": "OM.zone.x", "name": "x", "kind": "group"} }, id, ErrForbidden)
	refuse("empty text", func(m map[string]any) { m["text"] = "  " }, id, ErrInvalid)
	refuse("long text", func(m map[string]any) { m["text"] = strings.Repeat("a", 1001) }, id, ErrInvalid)
	refuse("control chars", func(m map[string]any) { m["text"] = "a\x00b" }, id, ErrInvalid)
	refuse("bad target", func(m map[string]any) { m["target"] = "partner" }, id, ErrInvalid)
	refuse("bad room kind", func(m map[string]any) { m["room"] = map[string]any{"id": "x", "name": "x", "kind": "broadcast"} }, id, ErrInvalid)
	// Unverified callers (only with COT_SEND_ALLOW_UNVERIFIED) skip the binding.
	if _, e := prepareChat(asReq(t, webGeochat("zeus", "Cyan", "group", "hi", "server")), identity{}, &cfg, now); e != nil {
		t.Errorf("unverified: %v", e)
	}
	// Direct message: TAK directed send header.
	m := webGeochat("zeus", "ANDROID-2", "direct", "hi", "mesh")
	p, e = prepareChat(asReq(t, m), id, &cfg, now)
	if e != nil || p.destUID != "ANDROID-2" {
		t.Fatalf("direct: %v %+v", e, p)
	}
	// No marking on a chat takes COT_SEND_CHAT_DEFAULT_TLP; with "" it is refused.
	c2 := cfg
	c2.Send.ChatDefaultTLP = ""
	if _, e := prepareChat(asReq(t, webGeochat("zeus", "Cyan", "group", "hi", "server")), id, &c2, now); e == nil || e.code != ErrForbidden {
		t.Errorf("no default TLP: want forbidden, got %v", e)
	}
}

func TestSendRateBucket(t *testing.T) {
	cfg := Defaults(TargetServer)
	cfg.Send.PerIdentityMin, cfg.Send.Burst = 60, 3
	s := newSender(&cfg, time.Now)
	now := time.Now()
	for i := 0; i < 3; i++ {
		if !s.allow("CLIENT/carol", now) {
			t.Fatalf("burst %d refused", i)
		}
	}
	if s.allow("CLIENT/carol", now) {
		t.Fatal("over burst allowed")
	}
	if !s.allow("CLIENT/zeus", now) {
		t.Fatal("another identity shares the bucket")
	}
	if !s.allow("CLIENT/carol", now.Add(1100*time.Millisecond)) {
		t.Fatal("bucket did not refill at 1/s")
	}
}

func TestSendConfigFromEnv(t *testing.T) {
	t.Setenv("COT_SEND_ACCOUNTS", "CLIENT,ATAK")
	t.Setenv("COT_SEND_RATE_PER_MIN", "12")
	t.Setenv("COT_SEND_CHAT_DEFAULT_TLP", "")
	c, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if c.Send.Subject != "svc.cot.send" || len(c.Send.Accounts) != 2 || c.Send.PerIdentityMin != 12 || c.Send.ChatDefaultTLP != "" {
		t.Fatalf("send config %+v", c.Send)
	}
	t.Setenv("COT_SEND_ENABLED", "false")
	if c, _ = FromEnv(); c.Send.Subject != "" {
		t.Fatal("COT_SEND_ENABLED=false did not disable")
	}
	t.Setenv("PROJECTOR_TARGET", "mesh")
	t.Setenv("COT_SEND_ENABLED", "")
	if c, _ = FromEnv(); c.Send.MaxBytes != c.MeshMaxBytes || c.Send.MaxPerSec != 2 {
		t.Fatalf("mesh send config %+v", c.Send)
	}
}

// ---- integration: accounts, verified identity, single writer ---------------------------------

type acctServer struct {
	url string
	pw  map[string]string
}

func pw() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// startAccountServer runs nats-server with the #562 account layout for svc.cot.send. Passwords
// are generated per run and never written outside the temp dir.
func startAccountServer(t *testing.T) *acctServer {
	t.Helper()
	bin := natsServerBin(t)
	dir := storeRoot(t)
	port := freePort(t)
	s := &acctServer{url: fmt.Sprintf("nats://127.0.0.1:%d", port), pw: map[string]string{}}
	for _, u := range []string{"projector", "onemind", "zeus", "carol", "hermes"} {
		s.pw[u] = pw()
	}
	conf := fmt.Sprintf(`listen: "127.0.0.1:%d"
jetstream { store_dir: %q }
accounts {
  ONEMIND {
    jetstream: enabled
    users: [ {user: projector, password: %q}, {user: onemind, password: %q} ]
    exports: [ {service: "svc.cot.send", accounts: [CLIENT]} ]
  }
  CLIENT {
    users: [ {user: zeus, password: %q}, {user: carol, password: %q} ]
    imports: [ {service: {account: ONEMIND, subject: "svc.cot.send"}, share: true} ]
  }
  AGENTS {
    users: [ {user: hermes, password: %q} ]
  }
}
`, port, filepath.Join(dir, "js"), s.pw["projector"], s.pw["onemind"], s.pw["zeus"], s.pw["carol"], s.pw["hermes"])
	cf := filepath.Join(dir, "nats.conf")
	if err := os.WriteFile(cf, []byte(conf), 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, "-c", cf)
	cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() })
	for i := 0; i < 100; i++ {
		if nc, err := nats.Connect(s.url, nats.UserInfo("onemind", s.pw["onemind"])); err == nil {
			nc.Close()
			return s
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("nats-server did not start")
	return nil
}

func (s *acctServer) conn(t *testing.T, user string) *nats.Conn {
	t.Helper()
	nc, err := nats.Connect(s.url, nats.UserInfo(user, s.pw[user]), nats.Name("test-"+user))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	return nc
}

func startSendProjector(t *testing.T, s *acctServer, cfg Config) (*Service, *logBuf) {
	t.Helper()
	nc := s.conn(t, "projector")
	lb := &logBuf{}
	svc, err := New(cfg, nc, slog.New(slog.NewJSONHandler(lb, nil)))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { defer close(done); _ = svc.Run(ctx) }()
	t.Cleanup(func() { cancel(); <-done })
	waitFor(t, 10*time.Second, "projector ready", svc.Ready)
	return svc, lb
}

type sendReply struct {
	V       int             `json:"v"`
	OK      bool            `json:"ok"`
	UID     string          `json:"uid"`
	Targets []string        `json:"targets"`
	Skipped []string        `json:"skipped"`
	Error   json.RawMessage `json:"error"`
	Detail  string          `json:"detail"`
}

func (r sendReply) code() string {
	var s string
	if json.Unmarshal(r.Error, &s) == nil {
		return s
	}
	var o struct {
		Code string `json:"code"`
	}
	_ = json.Unmarshal(r.Error, &o)
	return o.Code
}

func send(t *testing.T, nc *nats.Conn, body map[string]any) (sendReply, error) {
	t.Helper()
	b, _ := json.Marshal(body)
	m, err := nc.Request("svc.cot.send", b, 3*time.Second)
	if err != nil {
		return sendReply{}, err
	}
	var r sendReply
	if err := json.Unmarshal(m.Data, &r); err != nil {
		t.Fatalf("reply %s: %v", m.Data, err)
	}
	return r, nil
}

type subjCollector struct {
	mu   sync.Mutex
	msgs []*nats.Msg
}

func (c *subjCollector) n() int { c.mu.Lock(); defer c.mu.Unlock(); return len(c.msgs) }
func (c *subjCollector) get(i int) *nats.Msg {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.msgs[i]
}

func collectSubj(t *testing.T, nc *nats.Conn, subj string) *subjCollector {
	c := &subjCollector{}
	if _, err := nc.Subscribe(subj, func(m *nats.Msg) { c.mu.Lock(); c.msgs = append(c.msgs, m); c.mu.Unlock() }); err != nil {
		t.Fatal(err)
	}
	_ = nc.Flush()
	return c
}

func TestSendIntegrationHub(t *testing.T) {
	s := startAccountServer(t)
	cfg := testCfg(TargetServer)
	cfg.EntStream = ""
	cfg.ReleaseGroups = []string{"homestead"}
	cfg.Send.PerIdentityMin, cfg.Send.Burst = 6, 25
	cfg.Send.MaxPerSec, cfg.MaxRatePerSec = 1000, 1000 // per-identity limits only; the ceiling is tested on the mesh projector
	svc, logs := startSendProjector(t, s, cfg)
	admin := s.conn(t, "onemind")
	inject := collectSubj(t, admin, "tak.cot.inject")
	mesh := collectSubj(t, admin, "tak.cot.mesh.tx")
	zeus, carol := s.conn(t, "zeus"), s.conn(t, "carol")
	oks := 0
	expectOK := func(name string, r sendReply, err error) {
		t.Helper()
		if err != nil || !r.OK || r.V != 1 {
			t.Fatalf("%s: %v %+v", name, err, r)
		}
		oks++
		waitFor(t, 3*time.Second, name+" on tak.cot.inject", func() bool { return inject.n() == oks })
		evidence(t, "%s: ok uid=%s targets=%v skipped=%v", name, r.UID, r.Targets, r.Skipped)
	}
	expectErr := func(name string, r sendReply, err error, code string) {
		t.Helper()
		if err != nil || r.OK || r.code() != code {
			t.Fatalf("%s: want %s, got %v %+v", name, code, err, r)
		}
		evidence(t, "%s: refused %s (%s)", name, code, strings.TrimSpace(r.Detail+" "+string(r.Error)))
	}

	// 1. Web GeoChat from a verified CLIENT user.
	body := webGeochat("zeus", allChatRooms, "all", "eyes on north field", "both")
	r, err := send(t, zeus, body)
	expectOK("web geochat", r, err)
	if r.UID != body["message_uid"] || len(r.Targets) != 1 || r.Targets[0] != "server" || len(r.Skipped) != 1 || r.Skipped[0] != "mesh" {
		t.Fatalf("geochat reply %+v", r)
	}
	m := inject.get(0)
	if m.Header.Get("OneMind-Sender") != "CLIENT/zeus" || m.Header.Get("OneMind-Origin") != "svc.cot.send" || m.Header.Get("OneMind-Marking") != "CLEAR" {
		t.Fatalf("headers %v", m.Header)
	}
	if !strings.Contains(string(m.Data), `uid0="OM-WEB-zeus"`) || !strings.Contains(string(m.Data), `>eyes on north field</remarks>`) {
		t.Fatalf("geochat CoT %s", m.Data)
	}
	// 2. Unreal marker, route, geofence and chat: the published bytes are the client's.
	now := time.Now()
	for _, tc := range []struct{ kind, cot string }{
		{"marker", ueMarker("onemind-headless-1-UiErQFEk", now)},
		{"route", ueRoute("onemind-headless-4-Rr12Ss34", now)},
		{"geofence", ueGeofence("onemind-headless-5-Gg12Hh34", now)},
		{"chat", ueChat("onemind-headless", "hello from the OneMind client", now)},
	} {
		r, err = send(t, zeus, ueRequest(tc.kind, "headless", tc.cot, []string{"server", "mesh"}, "CLEAR"))
		expectOK("unreal "+tc.kind, r, err)
		if got := string(inject.get(oks - 1).Data); got != tc.cot {
			t.Fatalf("%s bytes changed:\n%s", tc.kind, got)
		}
	}
	for _, name := range []string{"ring", "rb", "line"} {
		cot := ueShapes(now)[name]
		r, err = send(t, zeus, ueRequest("shape", "mac", cot, []string{"server"}, "CLEAR"))
		expectOK("unreal M4 shape "+name, r, err)
		if got := string(inject.get(oks - 1).Data); got != cot {
			t.Fatalf("shape %s bytes changed", name)
		}
	}
	// 3. Idempotent replay: same client_ref, same reply, no second publish.
	marker := ueRequest("marker", "headless", ueMarker("onemind-headless-1-UiErQFEk", now), []string{"server"}, "CLEAR")
	r, err = send(t, zeus, marker)
	if err != nil || !r.OK || r.UID != "onemind-headless-1-UiErQFEk" {
		t.Fatalf("replay %v %+v", err, r)
	}
	time.Sleep(200 * time.Millisecond)
	if inject.n() != oks {
		t.Fatalf("replay published again: %d", inject.n())
	}
	evidence(t, "replay of client_ref onemind-headless-1-UiErQFEk: stored reply, no second publish")
	// 4. Refusals: nothing reaches tak.cot.inject.
	om := ueRequest("marker", "headless", ueMarker("OM.vessel.mmsi-367000001", now), []string{"server"}, "CLEAR")
	r, err = send(t, zeus, om)
	expectErr("OM. uid impersonation", r, err, ErrForbidden)
	r, err = send(t, zeus, ueRequest("marker", "headless", ueMarker("INTEL.acled-1", now), []string{"server"}, "CLEAR"))
	expectErr("INTEL. uid impersonation", r, err, ErrForbidden)
	r, err = send(t, zeus, ueRequest("marker", "headless", atakPLI("onemind-headless-9-Pp12", now), []string{"server"}, "CLEAR"))
	expectErr("EUD PLI (person-track data)", r, err, ErrForbidden)
	r, err = send(t, carol, webGeochat("zeus", "Cyan", "group", "I am zeus", "server"))
	expectErr("web sender impersonation (carol as OM-WEB-zeus)", r, err, ErrForbidden)
	if r.code() != ErrForbidden || !strings.Contains(string(r.Error), `"message"`) {
		t.Fatalf("structured error shape %s", r.Error)
	}
	r, err = send(t, zeus, ueRequest("marker", "headless", ueMarker("onemind-headless-6-Red1", now), []string{"server"}, "RED"))
	expectErr("marking RED", r, err, ErrForbidden)
	r, err = send(t, zeus, ueRequest("marker", "headless", ueMarker("onemind-headless-7-Amb1", now), []string{"server"}, "AMBER"))
	expectErr("marking AMBER without releasable_to", r, err, ErrForbidden)
	amber := ueRequest("marker", "headless", ueMarker("onemind-headless-8-Amb2", now), []string{"server"}, "AMBER")
	amber["marking"] = map[string]any{"tlp": "AMBER", "releasable_to": []string{"homestead"}}
	r, err = send(t, zeus, amber)
	expectOK("marking AMBER releasable to homestead (TAK_RELEASE_GROUPS)", r, err)
	if inject.get(oks-1).Header.Get("OneMind-Releasable-To") != "homestead" {
		t.Fatal("releasable header missing")
	}
	r, err = send(t, zeus, ueRequest("marker", "headless", ueMarker("onemind-headless-10-Mesh", now), []string{"mesh"}, "CLEAR"))
	expectErr("mesh-only send to the hub projector", r, err, ErrTargetFailed)
	// 5. Callers without a verified CLIENT identity.
	r, err = send(t, s.conn(t, "onemind"), webGeochat("onemind", "Cyan", "group", "hi", "server"))
	expectErr("same-account caller (no Nats-Request-Info)", r, err, ErrForbidden)
	forged := s.conn(t, "onemind")
	fm := nats.NewMsg("svc.cot.send")
	fm.Header.Set(RequestInfoHeader, `{"acc":"ONEMIND","user":"zeus"}`)
	fm.Data, _ = json.Marshal(webGeochat("zeus", "Cyan", "group", "hi", "server"))
	if rm, err := forged.RequestMsg(fm, 3*time.Second); err != nil || !strings.Contains(string(rm.Data), ErrForbidden) {
		t.Fatalf("own-account forged header: %v %s", err, rm.Data)
	}
	evidence(t, "own-account caller with a hand-written Nats-Request-Info: refused forbidden")
	if _, err := send(t, s.conn(t, "hermes"), webGeochat("hermes", "Cyan", "group", "hi", "server")); err == nil {
		t.Fatal("AGENTS reached svc.cot.send")
	}
	evidence(t, "AGENTS caller (no import): no responders")
	// 6. Per-identity rate limit: carol's burst runs out, zeus is unaffected.
	limited := 0
	for i := 0; i < 30; i++ {
		r, err = send(t, carol, webGeochat("carol", "Cyan", "group", fmt.Sprintf("msg %d", i), "server"))
		if err != nil {
			t.Fatal(err)
		}
		if r.OK {
			oks++
		} else if r.code() == ErrRateLimited {
			limited++
		}
	}
	// carol already spent one token on the impersonation attempt above.
	if limited != 6 {
		t.Fatalf("carol: %d rate-limited of 30 (burst 25, one spent), want 6", limited)
	}
	r, err = send(t, zeus, webGeochat("zeus", "Cyan", "group", "still here", "server"))
	expectOK("zeus after carol is limited", r, err)
	evidence(t, "per-identity rate limit (6/min, burst 25): carol 24 ok then 6 rate_limited; zeus unaffected")
	// 7. Kill switch G1.
	if err := admin.Publish("cmd.projector.disable", nil); err != nil {
		t.Fatal(err)
	}
	_ = admin.Flush()
	waitFor(t, 3*time.Second, "disabled", func() bool {
		r, _ = send(t, zeus, webGeochat("zeus", "Cyan", "group", "after disable", "server"))
		return !r.OK && r.code() == ErrForbidden
	})
	evidence(t, "cmd.projector.disable: svc.cot.send refused forbidden (G1)")
	_ = admin.Publish("cmd.projector.enable", nil)
	_ = admin.Flush()
	// Single writer: every publish on tak.cot.inject came from an ok reply; none on mesh.
	time.Sleep(200 * time.Millisecond)
	if inject.n() != oks || mesh.n() != 0 {
		t.Fatalf("tak.cot.inject %d (want %d), tak.cot.mesh.tx %d", inject.n(), oks, mesh.n())
	}
	for i := 0; i < inject.n(); i++ {
		if inject.get(i).Header.Get("OneMind-Projector") != "hub" {
			t.Fatal("publish without the projector header")
		}
	}
	c := svc.Engine.C.Snapshot()
	evidence(t, "single writer: %d publishes on tak.cot.inject = %d ok replies; counters ok=%d forbidden=%d rate_limited=%d target_failed=%d replayed=%d",
		inject.n(), oks, c["cot_send_ok"], c["cot_send_forbidden"], c["cot_send_rate_limited"], c["cot_send_target_failed"], c["cot_send_replayed"])
	if !strings.Contains(logs.String(), `"caller":"CLIENT/zeus"`) {
		t.Fatal("sends are not logged with the caller")
	}
}

func TestSendIntegrationMeshAndCeiling(t *testing.T) {
	s := startAccountServer(t)
	cfg := testCfg(TargetMesh)
	cfg.EntStream = ""
	cfg.Send.MaxPerSec = 3
	cfg.Send.PerIdentityMin, cfg.Send.Burst = 600, 50
	startSendProjector(t, s, cfg)
	admin := s.conn(t, "onemind")
	inject := collectSubj(t, admin, "tak.cot.inject")
	mesh := collectSubj(t, admin, "tak.cot.mesh.tx")
	zeus := s.conn(t, "zeus")
	r, err := send(t, zeus, webGeochat("zeus", "ANDROID-2", "direct", "direct to viper", "both"))
	if err != nil || !r.OK || r.Targets[0] != "mesh" || r.Skipped[0] != "server" {
		t.Fatalf("mesh geochat %v %+v", err, r)
	}
	waitFor(t, 3*time.Second, "mesh publish", func() bool { return mesh.n() == 1 })
	if mesh.get(0).Header.Get("Tak-Dest-Uid") != "ANDROID-2" || mesh.get(0).Reply != "" {
		t.Fatalf("mesh publish headers %v reply %q", mesh.get(0).Header, mesh.get(0).Reply)
	}
	evidence(t, "site projector: GeoChat direct message core-published on tak.cot.mesh.tx with Tak-Dest-Uid (CP-3), skipped server")
	r, _ = send(t, zeus, webGeochat("zeus", "Cyan", "group", "server only", "server"))
	if r.OK || r.code() != ErrTargetFailed {
		t.Fatalf("server-only to mesh projector: %+v", r)
	}
	// Client ceiling: at most MaxPerSec client publishes in any second.
	ok, limited := 0, 0
	for i := 0; i < 8; i++ {
		r, _ = send(t, zeus, webGeochat("zeus", "Cyan", "group", fmt.Sprintf("burst %d", i), "mesh"))
		if r.OK {
			ok++
		} else if r.code() == ErrRateLimited {
			limited++
		}
	}
	if ok > 3+1 || limited == 0 { // the first send above may have left the window
		t.Fatalf("ceiling: ok=%d limited=%d", ok, limited)
	}
	if inject.n() != 0 {
		t.Fatal("site projector wrote tak.cot.inject")
	}
	evidence(t, "client ceiling COT_SEND_MAX_PER_SEC=3: %d ok, %d rate_limited in one burst; nothing on tak.cot.inject", ok, limited)
}
