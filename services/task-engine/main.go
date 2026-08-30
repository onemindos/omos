// omos-task-engine: Sovereign NATS-native orchestration fabric.
//
// Every task is a dispatcher of real-world action. Drones patrol,
// robots run routines, sensors check in, agents respond. Each
// autonomous actor subscribes to its own NATS subject. Task dispatch
// is publishing to those subjects. Execution is live-visible on TAK
// as the Common Operating Picture.
//
// Architecture:
//
//	Task created → task.dispatch → cmd.drone.mavic-01.mission.start
//	→ drone executes → tel.drone.mavic-01.mission.status reports
//	→ omos-task-engine updates state → TAK COP shows live track
//
// Subjects:
//
//	task.create              (request-reply)
//	task.list                (request-reply)
//	task.get.<id>            (request-reply)
//	task.update.<id>         (request-reply)
//	task.approve.<id>        (request-reply)
//	task.delete.<id>         (request-reply)
//	task.calendar            (request-reply)
//	task.inbox               (request-reply)
//	task.project.create      (request-reply)
//	task.project.list        (request-reply)
//	task.syntax.validate     (request-reply)
//	task.dispatch.execute    (request-reply) — fire a TaskAction
//	task.dispatch.status.<id>          — response watching
//
// KV buckets: task-index, projects, dispatch-state
//
// Deploy: K8s Deployment on DOKS, connects to nats://nats.onemind.svc:4222.
// Dev: same binary, local NATS on zeus-mac (optionally NEX-packaged).
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

var (
	nc            *nats.Conn
	js            nats.JetStreamContext
	tasksKV       nats.KeyValue
	projectsKV    nats.KeyValue
	dispatchKV    nats.KeyValue
)

// ── helpers ───────────────────────────────────────────────────────────────

func nowISO() string { return time.Now().UTC().Format(time.RFC3339) }
func today() string  { return time.Now().UTC().Format("2006-01-02") }

func shortUUID() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func getStr(m map[string]interface{}, k string) string {
	if v, ok := m[k].(string); ok {
		return v
	}
	return ""
}

func getBool(m map[string]interface{}, k string) bool {
	if v, ok := m[k].(bool); ok {
		return v
	}
	return false
}

func getFloat(m map[string]interface{}, k string) float64 {
	if v, ok := m[k].(float64); ok {
		return v
	}
	return 0
}

func getMap(m map[string]interface{}, k string) map[string]interface{} {
	if v, ok := m[k].(map[string]interface{}); ok {
		return v
	}
	return nil
}

func respJSON(req micro.Request, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		respErr(req, err.Error())
		return
	}
	req.Respond(data)
}

func respErr(req micro.Request, msg string) {
	req.Respond([]byte(fmt.Sprintf(`{"error":%q}`, msg)))
}

func getTask(id string) (map[string]interface{}, error) {
	entry, err := tasksKV.Get(id)
	if err != nil {
		return nil, err
	}
	var t map[string]interface{}
	if err := json.Unmarshal(entry.Value(), &t); err != nil {
		return nil, err
	}
	return t, nil
}

func putTask(id string, t map[string]interface{}) error {
	data, _ := json.Marshal(t)
	_, err := tasksKV.Put(id, data)
	return err
}

// allTaskKeys returns task_id keys, skipping legacy "idx:" index entries.
func allTaskKeys() []string {
	keys, err := tasksKV.Keys()
	if err != nil {
		return []string{}
	}
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if strings.HasPrefix(k, "idx:") {
			continue
		}
		out = append(out, k)
	}
	return out
}

func fld(d map[string]interface{}, k string, def interface{}) interface{} {
	if v, ok := d[k]; ok {
		return v
	}
	return def
}

func handleSyntaxValidate(req micro.Request) {
	respJSON(req, map[string]bool{"valid": true, "syntax_ok": true})
}
