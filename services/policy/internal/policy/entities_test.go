package policy

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

func replyMsg(body string, hdr map[string]string) *nats.Msg {
	m := nats.NewMsg("_INBOX.x")
	m.Data = []byte(body)
	for k, v := range hdr {
		m.Header.Set(k, v)
	}
	return m
}

// Both error forms the entity service can send map to ErrEntityNotFound; other errors are
// lookup failures (the caller fails closed).
func TestDecodeEntityReply(t *testing.T) {
	cases := []struct {
		name     string
		msg      *nats.Msg
		notFound bool
		wantErr  bool
		wantID   string
	}{
		{name: "micro 404 header, empty body", msg: replyMsg("", map[string]string{micro.ErrorCodeHeader: "404", micro.ErrorHeader: "entity not found"}), notFound: true},
		{name: "json not_found body", msg: replyMsg(`{"error":{"code":"not_found"}}`, nil), notFound: true},
		{name: "json 404 body", msg: replyMsg(`{"error":{"code":"404"}}`, nil), notFound: true},
		{name: "micro 500 header", msg: replyMsg("", map[string]string{micro.ErrorCodeHeader: "500", micro.ErrorHeader: "respond failed"}), wantErr: true},
		{name: "micro 400 header", msg: replyMsg("", map[string]string{micro.ErrorCodeHeader: "400"}), wantErr: true},
		{name: "json other error", msg: replyMsg(`{"error":{"code":"internal"}}`, nil), wantErr: true},
		{name: "empty body, no header", msg: replyMsg("", nil), wantErr: true},
		{name: "wrapped entity", msg: replyMsg(`{"entity":{"entity_id":"sim-x1","entity_type":"drone"}}`, nil), wantID: "sim-x1"},
		{name: "bare entity", msg: replyMsg(`{"entity_id":"sim-x2","entity_type":"drone"}`, nil), wantID: "sim-x2"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e, err := decodeEntityReply(c.msg)
			switch {
			case c.notFound:
				if !errors.Is(err, ErrEntityNotFound) {
					t.Fatalf("want ErrEntityNotFound, got %v", err)
				}
			case c.wantErr:
				if err == nil || errors.Is(err, ErrEntityNotFound) {
					t.Fatalf("want a lookup failure, got %v", err)
				}
			default:
				if err != nil || e == nil || e.EntityID != c.wantID {
					t.Fatalf("want %s, got %+v %v", c.wantID, e, err)
				}
			}
		})
	}
}

// NATSEntities against a real micro service that answers the way services/entity does
// (r.Error("404", "entity not found", nil) for a missing id).
func TestNATSEntitiesMicroService(t *testing.T) {
	opts := &server.Options{Host: "127.0.0.1", Port: -1, NoLog: true, NoSigs: true}
	srv, err := server.NewServer(opts)
	if err != nil {
		t.Fatal(err)
	}
	go srv.Start()
	if !srv.ReadyForConnections(10 * time.Second) {
		t.Fatal("nats-server not ready")
	}
	defer func() { srv.Shutdown(); srv.WaitForShutdown() }()
	nc, err := nats.Connect(srv.ClientURL())
	if err != nil {
		t.Fatal(err)
	}
	defer nc.Close()

	svc, err := micro.AddService(nc, micro.Config{Name: "entity", Version: "0.1.0"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = svc.Stop() }()
	err = svc.AddEndpoint("get", micro.HandlerFunc(func(r micro.Request) {
		var q struct {
			EntityID string `json:"entity_id"`
		}
		if json.Unmarshal(r.Data(), &q) != nil || q.EntityID == "" {
			_ = r.Error("400", "body must be {\"entity_id\": \"...\"}", nil)
			return
		}
		switch q.EntityID {
		case "sim-x1":
			_ = r.RespondJSON(map[string]any{"entity_id": "sim-x1", "entity_type": "drone"})
		case "boom":
			_ = r.Error("500", "respond failed", nil)
		default:
			_ = r.Error("404", "entity not found", nil)
		}
	}), micro.WithEndpointSubject("svc.entity.get"))
	if err != nil {
		t.Fatal(err)
	}

	r := &NATSEntities{Conn: nc, Timeout: 2 * time.Second}
	if e, err := r.Get(context.Background(), "drone", "sim-x1"); err != nil || e.EntityID != "sim-x1" {
		t.Fatalf("found: %+v %v", e, err)
	}
	if _, err := r.Get(context.Background(), "drone", "sim-missing"); !errors.Is(err, ErrEntityNotFound) {
		t.Fatalf("missing: want ErrEntityNotFound, got %v", err)
	}
	if _, err := r.Get(context.Background(), "drone", "boom"); err == nil || errors.Is(err, ErrEntityNotFound) {
		t.Fatalf("500: want lookup failure, got %v", err)
	}
}

// Whole service: a missing target is the not_found denial (precondition_failed,
// unknown_entity) in both reply forms, never "entity service unavailable".
func TestMissingTargetDeniedBothForms(t *testing.T) {
	for _, microForm := range []bool{false, true} {
		name := "json body"
		if microForm {
			name = "micro header"
		}
		t.Run(name, func(t *testing.T) {
			e := newEnv(t)
			e.entMu.Lock()
			e.entMicro = microForm
			e.entMu.Unlock()
			r := e.propose("zeus", FamilyCommand, gotoReq("sim-x9", 30.27, -97.74))
			expectDenied(t, r, CodePreconditionFailed)
			if r.Denial.Detail != "unknown_entity" {
				t.Fatalf("detail: %+v", r.Denial)
			}
		})
	}
}
