package projector

// svc.cot.send: client-originated CoT (cot-projection.md §1 "Ephemeral client sends" and §9).
// The projector is the only writer of its output subject (tak.cot.inject on the hub,
// tak.cot.mesh.tx on a site), so every client send passes the same process and the same
// publisher as projected entities.
//
// Two request bodies are accepted, one per OneMind client:
//
//   - CoT form (Unreal client, onemind-xr client/OneMind/Docs/svc-cot-send.md, owner proposal
//     B-P4, UE-P1): {v, kind: marker|spi|route|geofence|shape|chat, targets[], client_id, client_ref,
//     marking{tlp, releasable_to}, dest_uid, cot}. The client builds the XML; the projector
//     re-checks all of it against an allowlist and publishes the checked bytes.
//     Reply {v, ok, uid, targets[, skipped]} or {v, ok: false, error: "<code>", detail}.
//   - Structured GeoChat (L0 web client, owner proposal WEB-P4): {v, kind: "geochat",
//     target: server|mesh|both, message_uid, room{id, name, kind}, text, sender{uid, callsign},
//     marking?}. The projector builds the CoT itself. Reply {v, ok, uid, targets} or
//     {v, ok: false, error: {code, message}}.
//
// Error codes are the same in both forms: invalid, forbidden, rate_limited, target_failed.
//
// Checks, in order: verified caller (the server-set Nats-Request-Info header of a
// `share: true` import from an allowed account, as the policy service reads it), idempotent
// replay, per-identity rate limit, body and CoT validation (reserved OM./INTEL. uids, uid bound
// to the client, kind-to-type table, civilian scope, person-track data, time window, size),
// marking gate G7, target, kill switch G1, client ceiling and the shared global ceiling.

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/nats-io/nats.go"
)

// RequestInfoHeader is set by nats-server on a request that crosses a service import; with
// `share: true` it carries the requester's account and user, and the server replaces any value
// the client sent (same header and rule as services/policy identity.go).
const RequestInfoHeader = "Nats-Request-Info"

// Send error codes.
const (
	ErrInvalid      = "invalid"
	ErrForbidden    = "forbidden"
	ErrRateLimited  = "rate_limited"
	ErrTargetFailed = "target_failed"
)

type sendErr struct{ code, detail string }

func (e *sendErr) Error() string { return e.code + ": " + e.detail }

func invalid(f string, a ...any) *sendErr   { return &sendErr{ErrInvalid, fmt.Sprintf(f, a...)} }
func forbidden(f string, a ...any) *sendErr { return &sendErr{ErrForbidden, fmt.Sprintf(f, a...)} }

// identity is the caller of one request.
type identity struct {
	Account, User string
	Verified      bool
}

func (i identity) key() string {
	if !i.Verified {
		return "unverified"
	}
	return i.Account + "/" + i.User
}

// sendRequest is the union of both bodies.
type sendRequest struct {
	V    int    `json:"v"`
	Kind string `json:"kind"`
	// CoT form
	Targets   []string `json:"targets"`
	ClientID  string   `json:"client_id"`
	ClientRef string   `json:"client_ref"`
	DestUID   string   `json:"dest_uid"`
	CoT       string   `json:"cot"`
	// GeoChat form
	Target     string `json:"target"`
	MessageUID string `json:"message_uid"`
	Room       *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Kind string `json:"kind"`
	} `json:"room"`
	Text   string `json:"text"`
	Sender *struct {
		UID      string `json:"uid"`
		Callsign string `json:"callsign"`
	} `json:"sender"`
	// both
	Marking *struct {
		TLP          string   `json:"tlp"`
		ReleasableTo []string `json:"releasable_to"`
	} `json:"marking"`
}

func (r *sendRequest) structured() bool { return r.Kind == "geochat" }

// prepared is a checked send, ready for the loop.
type prepared struct {
	uid, kind, tlp string
	releasable     []string
	xml            []byte
	destUID        string
	targets        []string // requested
}

type sendJob struct {
	p     *prepared
	id    identity
	reply chan *sendErr
}

// sender holds the per-identity state of svc.cot.send. Its methods run on NATS callback
// goroutines; publishing happens on the service loop.
type sender struct {
	cfg *Config
	mu  sync.Mutex
	bkt map[string]*bucket
	ded map[string]dedupEntry
	now func() time.Time
}

type bucket struct {
	tokens float64
	last   time.Time
}

// dedupEntry is a stored reply, or a pending one (reply nil) while the send is queued, so a
// retry with the same client_ref / message_uid never publishes twice.
type dedupEntry struct {
	reply   []byte
	expires time.Time
}

func newSender(cfg *Config, now func() time.Time) *sender {
	return &sender{cfg: cfg, bkt: map[string]*bucket{}, ded: map[string]dedupEntry{}, now: now}
}

// allow takes one token from the caller's bucket.
func (s *sender) allow(key string, now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	rate := float64(s.cfg.Send.PerIdentityMin) / 60
	burst := float64(s.cfg.Send.Burst)
	b := s.bkt[key]
	if b == nil {
		if len(s.bkt) >= 10000 { // forget idle callers (a full bucket holds no state)
			for k, v := range s.bkt {
				if v.tokens+now.Sub(v.last).Seconds()*rate >= burst {
					delete(s.bkt, k)
				}
			}
		}
		b = &bucket{tokens: burst, last: now}
		s.bkt[key] = b
	}
	b.tokens = math.Min(burst, b.tokens+now.Sub(b.last).Seconds()*rate)
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// dedupReserve returns the stored reply for key, or pending=true while another request with the
// same key is in flight; otherwise it records a pending entry and returns neither.
func (s *sender) dedupReserve(key string, now time.Time) (reply []byte, pending bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.ded[key]; ok && now.Before(e.expires) {
		return e.reply, e.reply == nil
	}
	s.gcLocked(now)
	s.ded[key] = dedupEntry{expires: now.Add(s.cfg.Send.DedupTTL)}
	return nil, false
}

// dedupRelease drops a pending entry after a refused send, so the client may retry.
func (s *sender) dedupRelease(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.ded[key]; ok && e.reply == nil {
		delete(s.ded, key)
	}
}

func (s *sender) dedupPut(key string, reply []byte, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gcLocked(now)
	s.ded[key] = dedupEntry{reply: reply, expires: now.Add(s.cfg.Send.DedupTTL)}
}

func (s *sender) gcLocked(now time.Time) {
	if len(s.ded) < 10000 {
		return
	}
	for k, e := range s.ded {
		if !now.Before(e.expires) {
			delete(s.ded, k)
		}
	}
	if len(s.ded) >= 10000 {
		for k, e := range s.ded { // keep pending entries; drop stored replies
			if e.reply != nil {
				delete(s.ded, k)
			}
		}
	}
}

// resolveIdentity reads the caller from Nats-Request-Info, refusing a header that names the
// projector's own account. This is NOT proof of identity on its own: nats-server sets the header
// only on a request that crosses a share: true import, and passes a header that a publisher in
// the projector's own account wrote by hand unchanged, so an own-account user can forge
// {"acc":"CLIENT","user":"zeus"} and the subscriber cannot tell. The identity is trustworthy
// only when no user of the own account except the projector may publish to svc.cot.send (bus
// control, owner proposal A6-P13). Until that is applied the endpoint stays off
// (COT_SEND_ENABLED defaults to false).
func resolveIdentity(m *nats.Msg, own string) identity {
	if m.Header == nil {
		return identity{}
	}
	raw := m.Header.Get(RequestInfoHeader)
	if raw == "" {
		return identity{}
	}
	var ci struct {
		Account string `json:"acc"`
		User    string `json:"user"`
	}
	if json.Unmarshal([]byte(raw), &ci) != nil || ci.Account == "" || ci.Account == own {
		return identity{}
	}
	return identity{Account: ci.Account, User: ci.User, Verified: true}
}

// handleSend serves one svc.cot.send request.
func (s *Service) handleSend(m *nats.Msg) {
	now := s.now()
	c := s.Engine.C
	structured := false
	fail := func(e *sendErr, id identity, kind string) {
		c.Inc("cot_send_" + e.code)
		lvl := s.log.Info
		if e.code == ErrForbidden {
			lvl = s.log.Warn
		}
		lvl("svc.cot.send refused", "code", e.code, "detail", e.detail, "caller", id.key(), "kind", kind)
		s.respond(m, sendReplyErr(structured, e))
	}
	id := resolveIdentity(m, s.cfg.Send.OwnAccount)
	if len(m.Data) > s.cfg.Send.MaxBytes+8192 {
		fail(invalid("request larger than %d bytes", s.cfg.Send.MaxBytes+8192), id, "")
		return
	}
	var req sendRequest
	if err := json.Unmarshal(m.Data, &req); err != nil {
		fail(invalid("request is not a JSON object"), id, "")
		return
	}
	structured = req.structured()
	if !id.Verified && !s.cfg.Send.AllowUnverified {
		fail(forbidden("caller identity not verified: send through the CLIENT account's svc.cot.send import"), id, req.Kind)
		return
	}
	if id.Verified && !contains(s.cfg.Send.Accounts, id.Account) {
		fail(forbidden("account %s may not send CoT", id.Account), id, req.Kind)
		return
	}
	ref := req.ClientRef
	if structured {
		ref = req.MessageUID
	}
	dkey := ""
	if ref != "" {
		dkey = id.key() + "|" + ref
		r, pending := s.sender.dedupReserve(dkey, now)
		if r != nil {
			c.Inc("cot_send_replayed")
			s.respond(m, r)
			return
		}
		if pending {
			fail(&sendErr{ErrTargetFailed, "a send with this client_ref / message_uid is in progress; retry for its reply"}, id, req.Kind)
			return
		}
		// Every refusal below releases the pending entry; a success stores the reply.
		inner := fail
		fail = func(e *sendErr, id identity, kind string) {
			s.sender.dedupRelease(dkey)
			inner(e, id, kind)
		}
	}
	if !s.sender.allow(id.key(), now) {
		fail(&sendErr{ErrRateLimited, fmt.Sprintf("more than %d sends per minute (burst %d) for this caller", s.cfg.Send.PerIdentityMin, s.cfg.Send.Burst)}, id, req.Kind)
		return
	}
	var p *prepared
	var e *sendErr
	if structured {
		p, e = prepareChat(&req, id, &s.cfg, now)
	} else {
		p, e = prepareCot(&req, &s.cfg, now)
	}
	if e != nil {
		fail(e, id, req.Kind)
		return
	}
	// G7, same rule as for entities (§3.4).
	probe := &Entity{}
	probe.Marking.TLP, probe.Marking.ReleasableTo = p.tlp, p.releasable
	if !markingAllowed(probe, s.cfg.ReleaseGroups) {
		fail(forbidden("marking gate: %q is not releasable to TAK", p.tlp), id, req.Kind)
		return
	}
	if !contains(p.targets, string(s.cfg.Target)) {
		fail(&sendErr{ErrTargetFailed, fmt.Sprintf("this projector serves only target %s", s.cfg.Target)}, id, req.Kind)
		return
	}
	job := &sendJob{p: p, id: id, reply: make(chan *sendErr, 1)}
	select {
	case s.sends <- job:
	case <-time.After(2 * time.Second):
		fail(&sendErr{ErrTargetFailed, "projector busy"}, id, req.Kind)
		return
	case <-s.done:
		fail(&sendErr{ErrTargetFailed, "projector stopping"}, id, req.Kind)
		return
	}
	// Once queued, wait for the loop's answer without a timeout: runSend is bounded (one core
	// publish), and replying target_failed while the publish may still happen would let a retry
	// publish twice. Only shutdown ends the wait.
	var res *sendErr
	select {
	case res = <-job.reply:
	case <-s.done:
		select {
		case res = <-job.reply:
		default:
			res = &sendErr{ErrTargetFailed, "projector stopping"}
		}
	}
	if res != nil {
		fail(res, id, req.Kind)
		return
	}
	var skipped []string
	for _, t := range p.targets {
		if t != string(s.cfg.Target) {
			skipped = append(skipped, t)
		}
	}
	out := map[string]any{"v": 1, "ok": true, "uid": p.uid, "targets": []string{string(s.cfg.Target)}}
	if len(skipped) > 0 {
		out["skipped"] = skipped
	}
	b, _ := json.Marshal(out)
	if dkey != "" {
		s.sender.dedupPut(dkey, b, now)
	}
	c.Inc("cot_send_ok")
	s.log.Info("svc.cot.send", "caller", id.key(), "verified", id.Verified, "kind", p.kind, "uid", p.uid,
		"subject", s.cfg.Subject, "marking", p.tlp)
	s.respond(m, b)
}

// runSend publishes one checked send on the loop goroutine (single writer, shared ceiling).
func (s *Service) runSend(j *sendJob, now time.Time) *sendErr {
	if !s.Engine.enabled {
		return forbidden("projector disabled (cmd.projector.disable, G1)")
	}
	if !s.Engine.clientBudget(now) {
		return &sendErr{ErrRateLimited, fmt.Sprintf("client sends on this target are over %d per second", s.cfg.Send.MaxPerSec)}
	}
	msg := nats.NewMsg(s.cfg.Subject)
	msg.Data = j.p.xml
	msg.Header.Set("Nats-Msg-Id", j.p.uid+":"+strconv.FormatInt(now.UnixMilli(), 10))
	msg.Header.Set("OneMind-Marking", j.p.tlp)
	if len(j.p.releasable) > 0 {
		msg.Header.Set("OneMind-Releasable-To", strings.Join(j.p.releasable, ","))
	}
	msg.Header.Set("OneMind-Projector", s.cfg.Instance)
	msg.Header.Set("OneMind-Origin", "svc.cot.send")
	msg.Header.Set("OneMind-Sender", j.id.key())
	if j.p.destUID != "" {
		msg.Header.Set("Tak-Dest-Uid", j.p.destUID) // mesh gateway directed send
	}
	if err := s.nc.PublishMsg(msg); err != nil { // core publish, no reply subject (CP-3)
		return &sendErr{ErrTargetFailed, err.Error()}
	}
	s.Engine.noteClientSend(now)
	return nil
}

func (s *Service) respond(m *nats.Msg, b []byte) {
	if m.Reply == "" {
		return
	}
	_ = m.Respond(b)
}

func sendReplyErr(structured bool, e *sendErr) []byte {
	var out map[string]any
	if structured {
		out = map[string]any{"v": 1, "ok": false, "error": map[string]string{"code": e.code, "message": e.detail}}
	} else {
		out = map[string]any{"v": 1, "ok": false, "error": e.code, "detail": e.detail}
	}
	b, _ := json.Marshal(out)
	return b
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// ---- shared field checks --------------------------------------------------------------------

var (
	uidRe      = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,128}$`)
	clientIDRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,31}$`)
	senderRe   = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)
	howHumanRe = regexp.MustCompile(`^h(-[a-z])+$`)
)

// reserved reports a value that names the projector's own uids: OM.* (entities) or INTEL.*
// (legacy intel path), alone or after a TAK source prefix.
func reserved(v string) bool {
	v = strings.TrimSpace(v)
	for _, p := range []string{"BAO.F.ATAK.", "BAO.F.OneMind.", "BAO.F.WinTAK.", "BAO.F.iTAK."} {
		if len(v) >= len(p) && strings.EqualFold(v[:len(p)], p) {
			v = v[len(p):]
		}
	}
	u := strings.ToUpper(v)
	return strings.HasPrefix(u, "OM.") || strings.HasPrefix(u, "INTEL.")
}

func cleanText(s string, max int, multiline bool) (string, bool) {
	if !utf8.ValidString(s) || utf8.RuneCountInString(s) > max {
		return "", false
	}
	for _, r := range s {
		if unicode.IsControl(r) && !(multiline && (r == '\n' || r == '\t')) {
			return "", false
		}
	}
	return s, true
}

func normTargets(req *sendRequest) ([]string, *sendErr) {
	var ts []string
	if req.structured() {
		switch req.Target {
		case "", "both":
			ts = []string{"server", "mesh"}
		case "server", "mesh":
			ts = []string{req.Target}
		default:
			return nil, invalid("target must be server, mesh or both")
		}
		return ts, nil
	}
	if len(req.Targets) == 0 || len(req.Targets) > 2 {
		return nil, invalid("targets must list server and/or mesh")
	}
	for _, t := range req.Targets {
		if t != "server" && t != "mesh" {
			return nil, invalid("unknown target %q", t)
		}
		if !contains(ts, t) {
			ts = append(ts, t)
		}
	}
	return ts, nil
}

func sendMarking(req *sendRequest, def string) (string, []string, *sendErr) {
	if req.Marking == nil || strings.TrimSpace(req.Marking.TLP) == "" {
		if def == "" {
			return "", nil, forbidden("marking gate: no marking.tlp")
		}
		return def, nil, nil
	}
	tlp := strings.ToUpper(strings.TrimSpace(req.Marking.TLP))
	if len(req.Marking.ReleasableTo) > 16 {
		return "", nil, invalid("marking.releasable_to lists more than 16 groups")
	}
	return tlp, req.Marking.ReleasableTo, nil
}

func newToken() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// ---- structured GeoChat (WEB-P4) ------------------------------------------------------------

const allChatRooms = "All Chat Rooms"

func prepareChat(req *sendRequest, id identity, cfg *Config, now time.Time) (*prepared, *sendErr) {
	if req.V != 1 {
		return nil, invalid("v must be 1")
	}
	if req.Room == nil || req.Sender == nil {
		return nil, invalid("room and sender are required")
	}
	ts, e := normTargets(req)
	if e != nil {
		return nil, e
	}
	text, ok := cleanText(req.Text, 1000, true)
	if !ok || strings.TrimSpace(text) == "" {
		return nil, invalid("text must be 1 to 1000 characters without control characters")
	}
	sender := req.Sender.UID
	if !senderRe.MatchString(sender) {
		return nil, invalid("sender.uid must match %s", senderRe)
	}
	if reserved(sender) {
		return nil, forbidden("sender.uid uses a reserved prefix (OM., INTEL.)")
	}
	// The sender's uid comes from the connection (WEB-P4): it must name the verified NATS user,
	// as `{prefix}-{user}` (web: OM-WEB-{user}) or the user itself.
	if id.Verified {
		u := strings.ToLower(sanitizeUser(id.User))
		su := strings.ToLower(sender)
		if u == "" || (su != u && !strings.HasSuffix(su, "-"+u)) {
			return nil, forbidden("sender.uid %q is not bound to the calling user", sender)
		}
	}
	callsign, ok := cleanText(strings.TrimSpace(req.Sender.Callsign), 64, false)
	if !ok || callsign == "" {
		return nil, invalid("sender.callsign must be 1 to 64 characters")
	}
	roomName, ok1 := cleanText(req.Room.Name, 128, false)
	roomID, ok2 := cleanText(req.Room.ID, 128, false)
	if !ok1 || !ok2 || roomID == "" {
		return nil, invalid("room.id and room.name must be 1 to 128 characters")
	}
	if roomName == "" {
		roomName = roomID
	}
	for _, v := range []string{sender, callsign, roomID, roomName, req.MessageUID} {
		if personWords.MatchString(v) {
			return nil, forbidden("sender, room or message_uid references a person or a device identifier (G4)")
		}
	}
	switch req.Room.Kind {
	case "all":
		roomID, roomName = allChatRooms, allChatRooms
	case "group":
	case "direct":
		if !senderRe.MatchString(roomID) {
			return nil, invalid("a direct room id is the contact uid")
		}
	default:
		return nil, invalid("room.kind must be all, group or direct")
	}
	if reserved(roomID) {
		return nil, forbidden("room.id uses a reserved prefix (OM., INTEL.)")
	}
	uid := req.MessageUID
	if uid == "" {
		uid = "GeoChat." + sender + "." + newToken()
	}
	// ATAK GeoChat uids carry the room name, spaces included ("GeoChat.{uid}.All Chat Rooms.{id}").
	if _, ok := cleanText(uid, 256, false); !ok || reserved(uid) {
		return nil, invalid("message_uid must be at most 256 characters without control characters")
	}
	if !strings.HasPrefix(uid, "GeoChat."+sender+".") {
		return nil, forbidden("message_uid must start with GeoChat.{sender.uid}.")
	}
	tlp, rel, e := sendMarking(req, cfg.Send.ChatDefaultTLP)
	if e != nil {
		return nil, e
	}
	msgID := uid[strings.LastIndexByte(uid, '.')+1:]
	t := now.UTC().Format(cotTimeLayout)
	var b xb
	b.WriteString(`<event`)
	b.attr("version", "2.0")
	b.attr("uid", uid)
	b.attr("type", "b-t-f")
	b.attr("how", "h-g-i-g-o")
	b.attr("time", t)
	b.attr("start", t)
	b.attr("stale", now.Add(cfg.Send.ChatStale).UTC().Format(cotTimeLayout))
	b.WriteString(`><point lat="0.0" lon="0.0" hae="9999999.0" ce="9999999.0" le="9999999.0"/><detail><__chat`)
	b.attr("parent", "RootContactGroup")
	b.attr("groupOwner", "false")
	b.attr("messageId", msgID)
	b.attr("chatroom", roomName)
	b.attr("id", roomID)
	b.attr("senderCallsign", callsign)
	b.WriteString(`><chatgrp`)
	b.attr("uid0", sender)
	b.attr("uid1", roomID)
	b.attr("id", roomID)
	b.WriteString(`/></__chat><link`)
	b.attr("uid", sender)
	b.attr("type", "a-f-G-U-C")
	b.attr("relation", "p-p")
	b.WriteString(`/><remarks`)
	b.attr("source", "BAO.F.OneMind."+sender)
	b.attr("to", roomID)
	b.attr("time", t)
	b.WriteString(`>`)
	_ = xml.EscapeText(&b.Buffer, []byte(text))
	b.WriteString(`</remarks></detail></event>`)
	if b.Len() > cfg.Send.MaxBytes {
		return nil, invalid("GeoChat CoT is %d bytes, over %d for this target", b.Len(), cfg.Send.MaxBytes)
	}
	dest := ""
	if req.Room.Kind == "direct" {
		dest = roomID
	}
	return &prepared{uid: uid, kind: "geochat", tlp: tlp, releasable: rel, xml: b.Bytes(), destUID: dest, targets: ts}, nil
}

func sanitizeUser(u string) string {
	var sb strings.Builder
	for _, r := range u {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// ---- client-built CoT (B-P4) ----------------------------------------------------------------

// kindTypes is the kind-to-type table. Persistent drawings belong to entity.create
// (cot-projection.md §9, B-P1); until that exists the Unreal client may send them here, but
// only when COT_SEND_PERSISTENT_KINDS=true (owner proposal A6-P6). Off by default.
var kindTypes = map[string]func(string) bool{
	"marker":   func(t string) bool { return strings.HasPrefix(t, "a-") || t == "b-m-p-s-m" },
	"spi":      func(t string) bool { return t == "b-m-p-s-p-i" },
	"route":    func(t string) bool { return t == "b-m-r" },
	"geofence": func(t string) bool { return t == "u-d-f" || t == "u-d-c-c" || t == "u-d-r" },
	// B11 drawings of the Unreal client (UE-P1; OneMindCore Drawing.cpp): spot marker, polygon
	// or open line, rectangle, circle, range ring, range and bearing.
	"shape": func(t string) bool {
		switch t {
		case "b-m-p-s-m", "u-d-f", "u-d-r", "u-d-c-c", "u-r-b-c-c", "u-rb-a":
			return true
		}
		return false
	},
	"chat": func(t string) bool { return t == "b-t-f" },
}

// persistentKinds are served only with SendConfig.Persistent (cot-projection.md §9 deviation).
var persistentKinds = map[string]bool{"marker": true, "route": true, "geofence": true, "shape": true}

// detailAllowed is the allowlist of <detail> descendants (drawing, route and chat detail of
// cot-projection.md §5.4-5.5 and the client builders). Anything else is refused.
var detailAllowed = map[string]bool{
	"contact": true, "remarks": true, "link": true, "strokeColor": true, "strokeWeight": true,
	"fillColor": true, "labels_on": true, "color": true, "usericon": true, "shape": true,
	"ellipse": true, "__geofence": true, "precisionlocation": true, "archive": true,
	"link_attr": true, "__chat": true, "chatgrp": true,
	// range and bearing (u-rb-a)
	"range": true, "bearing": true, "inclination": true, "rangeUnits": true, "bearingUnits": true, "northRef": true,
}

// personElements carry an EUD's own position report or a person's identity. The reverse path
// turns them into ent.person.* (cot-projection.md §8.2), so a client may not send them: that
// would be person-track data, or a forged team member.
var personElements = map[string]bool{
	"takv": true, "__group": true, "uid": true, "status": true, "track": true, "__video": true,
	"__serverdestination": true,
}

// projectorElements are written only by the projector or TAK Server.
var projectorElements = map[string]bool{"onemind": true, "_flow-tags_": true, "marti": true, "__forcedelete": true}

var personContactAttrs = map[string]bool{"endpoint": true, "phone": true, "emailAddress": true, "xmppUsername": true, "sipAddress": true}

// personWords are references to person-class entities or to device identifiers that track a
// person (BUS-DESIGN D9). Checked in every attribute value of client CoT (both forms) and in
// the remarks text of non-chat CoT. The free text of a chat message is not screened: it is the
// user's own words (README "svc.cot.send", owner proposal A6-P14).
var personWords = regexp.MustCompile(`(?i)(\bent\.(person|station)\.|\b(imsi|imei)\b|\bmac[_ -]?address\b)`)

const cotTimeLayout = "2006-01-02T15:04:05.000Z"

type cotSummary struct {
	uid, typ, how      string
	t, start, stale    time.Time
	links, linkUIDs    []string
	chatUID0, chatLink string
	remarksSource      string
	hasChat            bool
}

func parseCotTime(v string) (time.Time, bool) {
	t, err := time.Parse(time.RFC3339Nano, v)
	return t, err == nil
}

func checkLatLon(lat, lon string) bool {
	a, e1 := strconv.ParseFloat(lat, 64)
	o, e2 := strconv.ParseFloat(lon, 64)
	return e1 == nil && e2 == nil && validLL(a, o)
}

// scanCot walks the XML once, strictly, and refuses anything outside the allowlist.
func scanCot(doc []byte, maxElems int) (*cotSummary, *sendErr) {
	d := xml.NewDecoder(bytes.NewReader(doc))
	d.Strict = true
	sum := &cotSummary{}
	var stack []string
	elems, roots := 0, 0
	inRemarks := false
	for {
		tok, err := d.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, invalid("cot is not well-formed XML: %v", err)
		}
		switch tk := tok.(type) {
		case xml.ProcInst:
			if tk.Target != "xml" || elems > 0 {
				return nil, invalid("processing instructions are not allowed")
			}
		case xml.Directive:
			return nil, invalid("DOCTYPE and other directives are not allowed")
		case xml.Comment:
			return nil, invalid("comments are not allowed")
		case xml.CharData:
			if len(stack) == 0 && len(bytes.TrimSpace(tk)) > 0 {
				return nil, invalid("text outside <event>")
			}
			if inRemarks && !sum.hasChat && personWords.Match(tk) {
				return nil, forbidden("person-track data in remarks")
			}
		case xml.EndElement:
			if len(stack) > 0 && stack[len(stack)-1] == "remarks" {
				inRemarks = false
			}
			stack = stack[:len(stack)-1]
		case xml.StartElement:
			elems++
			if elems > maxElems {
				return nil, invalid("more than %d elements", maxElems)
			}
			name := tk.Name.Local
			if tk.Name.Space != "" {
				return nil, invalid("XML namespaces are not allowed")
			}
			depth := len(stack)
			if depth >= 4 {
				return nil, invalid("nesting deeper than 4")
			}
			attrs := map[string]string{}
			for _, a := range tk.Attr {
				if a.Name.Space != "" || a.Name.Local == "xmlns" {
					return nil, invalid("XML namespaces are not allowed")
				}
				if _, dup := attrs[a.Name.Local]; dup {
					return nil, invalid("duplicate attribute %s on <%s>", a.Name.Local, name)
				}
				attrs[a.Name.Local] = a.Value
				if reserved(a.Value) {
					return nil, forbidden("<%s %s> names a reserved uid (OM., INTEL.)", name, a.Name.Local)
				}
				if personWords.MatchString(a.Value) {
					return nil, forbidden("<%s %s> references a person or a device identifier (G4)", name, a.Name.Local)
				}
			}
			switch {
			case depth == 0:
				roots++
				if name != "event" || roots > 1 {
					return nil, invalid("root element must be a single <event>")
				}
				if attrs["version"] != "2.0" {
					return nil, invalid("event version must be 2.0")
				}
				for _, k := range []string{"access", "caveat", "releasableTo", "qos"} {
					if _, ok := attrs[k]; ok {
						return nil, forbidden("event attribute %s is set by the projector, not the client", k)
					}
				}
				sum.uid, sum.typ, sum.how = attrs["uid"], attrs["type"], attrs["how"]
				var ok1, ok2, ok3 bool
				sum.t, ok1 = parseCotTime(attrs["time"])
				sum.start, ok2 = parseCotTime(attrs["start"])
				sum.stale, ok3 = parseCotTime(attrs["stale"])
				if !ok1 || !ok2 || !ok3 {
					return nil, invalid("time, start and stale must be RFC 3339 UTC times")
				}
			case depth == 1:
				switch name {
				case "point":
					if !checkLatLon(attrs["lat"], attrs["lon"]) {
						return nil, invalid("point lat/lon out of range")
					}
				case "detail":
				default:
					return nil, invalid("<%s> is not allowed under <event>", name)
				}
			default:
				if stack[1] != "detail" {
					return nil, invalid("<%s> is not allowed under <%s>", name, stack[len(stack)-1])
				}
				switch {
				case personElements[name]:
					return nil, forbidden("<%s> is person-track data (an EUD position report); clients may not send it", name)
				case projectorElements[name]:
					return nil, forbidden("<%s> is written only by the projector or TAK Server", name)
				case !detailAllowed[name]:
					return nil, invalid("detail element <%s> is not allowed", name)
				}
				switch name {
				case "contact":
					for k := range attrs {
						if personContactAttrs[k] {
							return nil, forbidden("<contact %s> is person-track data (an EUD position report)", k)
						}
					}
				case "link":
					if p, ok := attrs["point"]; ok {
						ll := strings.Split(p, ",")
						if len(ll) < 2 || !checkLatLon(strings.TrimSpace(ll[0]), strings.TrimSpace(ll[1])) {
							return nil, invalid("link point out of range")
						}
						sum.links = append(sum.links, p)
					}
					if u, ok := attrs["uid"]; ok {
						if attrs["relation"] == "p-p" {
							sum.chatLink = u
						} else {
							sum.linkUIDs = append(sum.linkUIDs, u)
						}
					}
				case "__chat":
					sum.hasChat = true
				case "chatgrp":
					sum.chatUID0 = attrs["uid0"]
				case "remarks":
					inRemarks = true
					sum.remarksSource = attrs["source"]
				}
			}
			stack = append(stack, name)
		}
	}
	if roots != 1 || len(stack) != 0 {
		return nil, invalid("cot must be one complete <event>")
	}
	return sum, nil
}

func prepareCot(req *sendRequest, cfg *Config, now time.Time) (*prepared, *sendErr) {
	if req.V != 1 {
		return nil, invalid("v must be 1")
	}
	match, ok := kindTypes[req.Kind]
	if !ok {
		return nil, invalid("kind must be marker, spi, route, geofence, shape, chat or geochat")
	}
	if persistentKinds[req.Kind] && !cfg.Send.Persistent {
		return nil, forbidden("kind %s is a persistent item; it belongs to entity.create (cot-projection.md section 9) and COT_SEND_PERSISTENT_KINDS is off", req.Kind)
	}
	ts, e := normTargets(req)
	if e != nil {
		return nil, e
	}
	if !clientIDRe.MatchString(req.ClientID) {
		return nil, invalid("client_id must match %s", clientIDRe)
	}
	if req.CoT == "" {
		return nil, invalid("cot is required")
	}
	if len(req.CoT) > cfg.Send.MaxBytes {
		return nil, invalid("cot is %d bytes, over %d for this target", len(req.CoT), cfg.Send.MaxBytes)
	}
	if !utf8.ValidString(req.CoT) {
		return nil, invalid("cot is not UTF-8")
	}
	if req.DestUID != "" && (!senderRe.MatchString(req.DestUID) || reserved(req.DestUID)) {
		return nil, invalid("dest_uid must match %s and not be reserved", senderRe)
	}
	sum, e := scanCot([]byte(req.CoT), cfg.MaxVertices+32)
	if e != nil {
		return nil, e
	}
	// uid: safe, never reserved, and bound to the client so it cannot take over another
	// TAK user's marker or contact.
	self := "onemind-" + req.ClientID
	if req.Kind == "chat" {
		// GeoChat uids carry the room name, spaces included (BuildChat in the Unreal client).
		if _, ok := cleanText(sum.uid, 256, false); !ok || sum.uid == "" {
			return nil, invalid("uid must be at most 256 characters without control characters")
		}
		if !strings.HasPrefix(sum.uid, "GeoChat."+self+".") {
			return nil, forbidden("chat uid must start with GeoChat.%s. (the sending client)", self)
		}
	} else {
		if !uidRe.MatchString(sum.uid) {
			return nil, invalid("uid must match %s", uidRe)
		}
		if !strings.HasPrefix(sum.uid, self+"-") {
			return nil, forbidden("uid must start with %s- (the sending client)", self)
		}
	}
	if req.ClientRef != "" && req.ClientRef != sum.uid {
		return nil, invalid("client_ref must be the CoT uid")
	}
	for _, u := range sum.linkUIDs {
		if !strings.HasPrefix(u, sum.uid+".") {
			return nil, forbidden("link uid %q is not a child of %s", u, sum.uid)
		}
	}
	// type
	if !cotTypeRe.MatchString(sum.typ) || !match(sum.typ) {
		return nil, invalid("type %q is not allowed for kind %s", sum.typ, req.Kind)
	}
	if _, ok, down := civilian(sum.typ); !ok || down {
		return nil, forbidden("type %q is outside the civilian scope (hostile, tasking or weapon)", sum.typ)
	}
	if req.Kind == "marker" && strings.HasPrefix(sum.typ, "a-s-") {
		return nil, forbidden("suspect affiliation is set by fusion, not by a client marker")
	}
	if !howHumanRe.MatchString(sum.how) {
		return nil, forbidden("how %q: client CoT is human-entered (h-*); machine values belong to the projector", sum.how)
	}
	// time window
	if d := sum.t.Sub(now); d > cfg.Send.Skew || d < -cfg.Send.Skew {
		return nil, invalid("time is %s from the projector's clock (max %s)", d.Round(time.Second), cfg.Send.Skew)
	}
	if !sum.stale.After(sum.t) || sum.start.After(sum.stale) {
		return nil, invalid("stale must be after time and start")
	}
	if sum.stale.Sub(now) > cfg.Send.MaxStale {
		return nil, invalid("stale is more than %s ahead", cfg.Send.MaxStale)
	}
	// shape sizes
	if len(sum.links) > cfg.MaxVertices+1 {
		return nil, invalid("more than %d vertices", cfg.MaxVertices)
	}
	switch req.Kind {
	case "route":
		if len(sum.links) < 2 {
			return nil, invalid("a route needs at least 2 waypoints")
		}
	case "geofence":
		if sum.typ == "u-d-f" && len(sum.links) < 3 {
			return nil, invalid("a polygon needs at least 3 vertices")
		}
	case "shape":
		if sum.typ == "u-d-f" && len(sum.links) < 2 { // a polygon, or an open line (Drawing.cpp BuildLine)
			return nil, invalid("a polygon or line needs at least 2 vertices")
		}
	}
	// chat binding: the sender inside the CoT is the sending client.
	if req.Kind == "chat" {
		if !sum.hasChat || sum.chatUID0 == "" {
			return nil, invalid("chat needs <__chat><chatgrp uid0>")
		}
		if sum.chatUID0 != self || (sum.chatLink != "" && sum.chatLink != self) ||
			(sum.remarksSource != "" && sum.remarksSource != "BAO.F.OneMind."+self) {
			return nil, forbidden("chat sender must be %s (the sending client)", self)
		}
	} else if sum.hasChat {
		return nil, invalid("<__chat> only in kind chat")
	}
	tlp, rel, e := sendMarking(req, "")
	if e != nil {
		return nil, e
	}
	return &prepared{uid: sum.uid, kind: req.Kind, tlp: tlp, releasable: rel, xml: []byte(req.CoT),
		destUID: req.DestUID, targets: ts}, nil
}
