package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go/micro"
)

// ── task_write.go — mutation handlers ──────────────────────────────────────

func handleCreate(req micro.Request) {
	var t map[string]interface{}
	if err := json.Unmarshal(req.Data(), &t); err != nil {
		respErr(req, "invalid JSON")
		return
	}

	id := getStr(t, "task_id")
	if id == "" {
		id = shortUUID()
		t["task_id"] = id
	}
	t["created_at"] = fld(t, "created_at", nowISO())
	t["updated_at"] = nowISO()
	if getStr(t, "status") == "" {
		t["status"] = "pending"
	}

	if err := putTask(id, t); err != nil {
		respErr(req, fmt.Sprintf("save failed: %v", err))
		return
	}

	// Publish event
	publishEvent("evt.task.created", t)

	respJSON(req, t)
}

func handleUpdate(req micro.Request) {
	id := strings.TrimPrefix(req.Subject(), "task.update.")
	if id == "" || id == "*" {
		respErr(req, "missing task id")
		return
	}

	existing, err := getTask(id)
	if err != nil {
		respErr(req, fmt.Sprintf("not found: %s", id))
		return
	}

	var patch map[string]interface{}
	if err := json.Unmarshal(req.Data(), &patch); err != nil {
		respErr(req, "invalid JSON patch")
		return
	}

	// Merge
	for k, v := range patch {
		if k == "task_id" || k == "created_at" {
			continue // immutable
		}
		existing[k] = v
	}
	existing["updated_at"] = nowISO()

	if err := putTask(id, existing); err != nil {
		respErr(req, fmt.Sprintf("update failed: %v", err))
		return
	}

	publishEvent("evt.task.updated", existing)
	respJSON(req, existing)
}

func handleApprove(req micro.Request) {
	id := strings.TrimPrefix(req.Subject(), "task.approve.")
	if id == "" || id == "*" {
		respErr(req, "missing task id")
		return
	}

	t, err := getTask(id)
	if err != nil {
		respErr(req, fmt.Sprintf("not found: %s", id))
		return
	}

	t["status"] = "approved"
	t["approved_at"] = nowISO()
	t["updated_at"] = nowISO()

	if err := putTask(id, t); err != nil {
		respErr(req, fmt.Sprintf("approve failed: %v", err))
		return
	}

	publishEvent("evt.task.approved", t)
	respJSON(req, t)
}

func handleDelete(req micro.Request) {
	id := strings.TrimPrefix(req.Subject(), "task.delete.")
	if id == "" || id == "*" {
		respErr(req, "missing task id")
		return
	}

	t, err := getTask(id)
	if err != nil {
		respErr(req, fmt.Sprintf("not found: %s", id))
		return
	}

	if err := tasksKV.Delete(id); err != nil {
		respErr(req, fmt.Sprintf("delete failed: %v", err))
		return
	}

	publishEvent("evt.task.deleted", t)
	respJSON(req, map[string]string{"deleted": id})
}

// publishEvent sends a fire-and-forget event on evt.task.*
func publishEvent(subject string, payload interface{}) {
	data, _ := json.Marshal(payload)
	_ = nc.Publish(subject+".body", data) // TODO: full JSON event wrapper
	_ = nc.Publish(subject, data)
}
