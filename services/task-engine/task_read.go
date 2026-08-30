package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go/micro"
)

// ── task_read.go — query handlers ─────────────────────────────────────────

func handleGet(req micro.Request) {
	id := strings.TrimPrefix(req.Subject(), "task.get.")
	if id == "" || id == "*" {
		respErr(req, "missing task id")
		return
	}
	t, err := getTask(id)
	if err != nil {
		respErr(req, fmt.Sprintf("not found: %s", id))
		return
	}
	respJSON(req, t)
}

func handleList(req micro.Request) {
	keys := allTaskKeys()
	out := make([]map[string]interface{}, 0, len(keys))
	for _, k := range keys {
		t, err := getTask(k)
		if err != nil {
			continue
		}
		out = append(out, t)
	}
	respJSON(req, out)
}

func handleCalendar(req micro.Request) {
	var q struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	json.Unmarshal(req.Data(), &q)
	if q.From == "" {
		q.From = "1970-01-01"
	}
	if q.To == "" {
		q.To = "2099-12-31"
	}

	keys := allTaskKeys()
	out := make([]map[string]interface{}, 0)
	for _, k := range keys {
		t, err := getTask(k)
		if err != nil {
			continue
		}
		due := getStr(t, "due_date")
		if due == "" {
			due = getStr(t, "due")
		}
		if due >= q.From && due <= q.To {
			out = append(out, t)
		}
	}
	respJSON(req, out)
}

func handleInbox(req micro.Request) {
	keys := allTaskKeys()
	now := today()
	out := make([]map[string]interface{}, 0)
	for _, k := range keys {
		t, err := getTask(k)
		if err != nil {
			continue
		}
		status := getStr(t, "status")
		if status == "completed" || status == "cancelled" || status == "approved" {
			continue
		}
		due := getStr(t, "due_date")
		if due == "" {
			due = getStr(t, "due")
		}
		// Include overdue + due-today + no-due-date (pending)
		if due == "" || due <= now {
			out = append(out, t)
		}
	}
	respJSON(req, out)
}
