package main

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/nats-io/nats.go/micro"
)

// project.go: project create / list handlers

func handleProjectCreate(req micro.Request) {
	data := M{}
	if err := json.Unmarshal(req.Data(), &data); err != nil {
		respErr(req, "invalid JSON: "+err.Error())
		return
	}
	name := strings.TrimSpace(getStr(data, "name"))
	if name == "" {
		respErr(req, "name is required")
		return
	}
	pid := "proj_" + shortUUID()
	ts := nowISO()
	color := getStr(data, "color")
	if color == "" {
		color = "so"
	}
	icon := getStr(data, "icon")
	if icon == "" {
		icon = "inbox"
	}
	p := M{"project_id": pid, "name": name, "color": color, "icon": icon, "archived": false, "created_at": ts, "updated_at": ts}
	pd, _ := json.Marshal(p)
	projectsKV.Put(pid, pd)
	respJSON(req, M{"project_id": pid, "status": "created"})
}

func handleProjectList(req micro.Request) {
	out := make([]M, 0)
	keys, err := projectsKV.Keys()
	if err == nil {
		for _, k := range keys {
			entry, e := projectsKV.Get(k)
			if e != nil {
				continue
			}
			p := M{}
			if json.Unmarshal(entry.Value(), &p) != nil {
				continue
			}
			if getBool(p, "archived") {
				continue
			}
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return getStr(out[i], "name") < getStr(out[j], "name") })
	respJSON(req, M{"projects": out})
}
