package main

import (
	"encoding/json"
	"strings"

	"github.com/nats-io/nats.go/micro"
)

// task_write.go: create/update/approve/delete handlers

type M = map[string]interface{}

func emptyList() []interface{} { return make([]interface{}, 0) }

func subjectID(req micro.Request) string {
	parts := strings.Split(req.Subject(), ".")
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

var updatable = []string{"status", "agent_id", "result", "error", "progress", "due_date", "due_ts", "start_date", "is_all_day", "recurrence", "reminder_mins", "project_id", "labels", "archived", "completed_at", "title", "description", "priority", "domain", "tags"}

func taskDefaults(d M, id, ts string) M {
	t := M{
		"task_id": id, "title": strings.TrimSpace(getStr(d, "title")),
		"status": "inbox", "progress": 0, "archived": false,
		"result": nil, "error": nil, "completed_at": nil,
		"created_at": ts, "updated_at": ts,
	}
	for _, k := range []string{"description", "agent_id", "domain", "due_ts", "due_date", "start_date", "recurrence", "project_id"} {
		t[k] = fld(d, k, nil)
	}
	t["priority"] = fld(d, "priority", "normal")
	t["is_all_day"] = fld(d, "is_all_day", true)
	t["requires_approval"] = fld(d, "requires_approval", false)
	for _, k := range []string{"tags", "labels", "reminder_mins"} {
		t[k] = fld(d, k, emptyList())
	}
	return t
}

func handleCreate(req micro.Request) {
	data := M{}
	if len(req.Data()) > 0 {
		if err := json.Unmarshal(req.Data(), &data); err != nil {
			respErr(req, "invalid JSON: "+err.Error())
			return
		}
	}
	if strings.TrimSpace(getStr(data, "title")) == "" {
		respErr(req, "title is required")
		return
	}
	id := "task_" + shortUUID()
	t := taskDefaults(data, id, nowISO())
	if err := putTask(id, t); err != nil {
		respErr(req, "store failed: "+err.Error())
		return
	}
	evt, _ := json.Marshal(M{"task_id": id, "title": t["title"], "domain": t["domain"], "ts": t["created_at"]})
	nc.Publish("evt.task.created", evt)
	respJSON(req, M{"task_id": id, "status": t["status"]})
}

func handleUpdate(req micro.Request) {
	id := subjectID(req)
	if id == "" {
		respErr(req, "malformed subject")
		return
	}
	data := M{}
	if err := json.Unmarshal(req.Data(), &data); err != nil {
		respErr(req, "invalid JSON: "+err.Error())
		return
	}
	t, err := getTask(id)
	if err != nil {
		respErr(req, "task "+id+" not found")
		return
	}
	for _, k := range updatable {
		if v, ok := data[k]; ok {
			t[k] = v
		}
	}
	st := getStr(t, "status")
	if (st == "done" || st == "completed") && getStr(t, "completed_at") == "" {
		t["completed_at"] = nowISO()
	}
	t["updated_at"] = nowISO()
	putTask(id, t)
	evt, _ := json.Marshal(M{"task_id": id, "status": st, "ts": nowISO()})
	nc.Publish("evt.task.updated", evt)
	respJSON(req, M{"task_id": id, "status": st, "ok": true})
}

func handleApprove(req micro.Request) {
	id := subjectID(req)
	if id == "" {
		respErr(req, "malformed subject")
		return
	}
	data := M{}
	json.Unmarshal(req.Data(), &data)
	t, err := getTask(id)
	if err != nil {
		respErr(req, "task "+id+" not found")
		return
	}
	approved := getBool(data, "approved")
	approver := getStr(data, "approver")
	if approver == "" {
		approver = "unknown"
	}
	ts := nowISO()
	status := "rejected"
	evtSub := "evt.task." + id + ".rejected"
	if approved {
		status = "approved"
		evtSub = "evt.task." + id + ".approved"
	}
	t["status"] = status
	t["approver"] = approver
	t["reason"] = getStr(data, "reason")
	t["updated_at"] = ts
	putTask(id, t)
	evt, _ := json.Marshal(M{"task_id": id, "approved": approved, "approver": approver, "ts": ts})
	nc.Publish(evtSub, evt)
	respJSON(req, M{"ok": true, "task_id": id, "status": status})
}

func handleDelete(req micro.Request) {
	id := subjectID(req)
	if id == "" {
		respErr(req, "malformed subject")
		return
	}
	tasksKV.Delete(id)
	nc.Publish("evt.task.deleted", []byte(`{"task_id":"`+id+`"}`))
	respJSON(req, M{"ok": true, "task_id": id, "deleted": true})
}
