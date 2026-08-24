package main

import (
	"encoding/json"
	"sort"

	"github.com/nats-io/nats.go/micro"
)

// task_read.go: list / get / inbox / calendar query handlers

func handleGet(req micro.Request) {
	id := subjectID(req)
	if id == "" {
		respErr(req, "malformed subject")
		return
	}
	t, err := getTask(id)
	if err != nil {
		respErr(req, "not found")
		return
	}
	respJSON(req, t)
}

func matchFilter(t, f M) bool {
	for _, k := range []string{"status", "domain", "agent_id", "project_id"} {
		if fv := getStr(f, k); fv != "" && getStr(t, k) != fv {
			return false
		}
	}
	return true
}

func handleList(req micro.Request) {
	f := M{}
	json.Unmarshal(req.Data(), &f)
	limit := 200
	if v, ok := f["limit"].(float64); ok && v > 0 {
		limit = int(v)
	}
	archived := getBool(f, "archived")
	out := make([]M, 0)
	for _, key := range allTaskKeys() {
		t, err := getTask(key)
		if err != nil {
			continue
		}
		if !archived && getBool(t, "archived") {
			continue
		}
		if !matchFilter(t, f) {
			continue
		}
		out = append(out, t)
		if len(out) >= limit {
			break
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return getStr(out[i], "created_at") > getStr(out[j], "created_at")
	})
	respJSON(req, M{"tasks": out, "count": len(out)})
}

func dueKey(t M) string {
	d := getStr(t, "due_date")
	if d == "" {
		return "9999"
	}
	return d
}

func handleInbox(req micro.Request) {
	td := today()
	out := make([]M, 0)
	for _, key := range allTaskKeys() {
		t, err := getTask(key)
		if err != nil {
			continue
		}
		if getBool(t, "archived") {
			continue
		}
		st := getStr(t, "status")
		if st == "done" || st == "completed" || st == "failed" {
			continue
		}
		due := getStr(t, "due_date")
		if st == "inbox" || st == "pending" || (due != "" && due <= td) {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return dueKey(out[i]) < dueKey(out[j]) })
	respJSON(req, M{"tasks": out, "count": len(out)})
}

func handleCalendar(req micro.Request) {
	data := M{}
	json.Unmarshal(req.Data(), &data)
	from := getStr(data, "from_date")
	to := getStr(data, "to_date")
	out := make([]M, 0)
	for _, key := range allTaskKeys() {
		t, err := getTask(key)
		if err != nil {
			continue
		}
		if getBool(t, "archived") || getStr(t, "due_date") == "" {
			continue
		}
		due := getStr(t, "due_date")
		if from != "" && due < from {
			continue
		}
		if to != "" && due > to {
			continue
		}
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return getStr(out[i], "due_date") < getStr(out[j], "due_date") })
	respJSON(req, M{"tasks": out, "count": len(out)})
}
