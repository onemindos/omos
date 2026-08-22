package main

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go/micro"
)

// ── project.go — project CRUD ──────────────────────────────────────────────

func handleProjectCreate(req micro.Request) {
	var p map[string]interface{}
	if err := json.Unmarshal(req.Data(), &p); err != nil {
		respErr(req, "invalid project JSON")
		return
	}

	id := getStr(p, "project_id")
	if id == "" {
		id = shortUUID()
		p["project_id"] = id
	}
	p["created_at"] = nowISO()

	data, _ := json.Marshal(p)
	if _, err := projectsKV.Put(id, data); err != nil {
		respErr(req, fmt.Sprintf("save project failed: %v", err))
		return
	}

	respJSON(req, p)
}

func handleProjectList(req micro.Request) {
	keys, err := projectsKV.Keys()
	if err != nil {
		respJSON(req, []map[string]interface{}{})
		return
	}

	out := make([]map[string]interface{}, 0, len(keys))
	for _, k := range keys {
		entry, err := projectsKV.Get(k)
		if err != nil {
			continue
		}
		var p map[string]interface{}
		if err := json.Unmarshal(entry.Value(), &p); err != nil {
			continue
		}
		out = append(out, p)
	}
	respJSON(req, out)
}
