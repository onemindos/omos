package main

import (
	"encoding/json"
	"fmt"
)

// ── geo.go — GeoBinding and TAK CoT integration ────────────────────────────
//
// Every task can be geospatially bound. When a task dispatches, the engine
// publishes a CoT (Cursor on Target) event to the TAK feed, placing the
// task on the Common Operating Picture. As the task executes, updated CoT
// markers track the executor's position.
//
// CoT event types (MIL-STD-2525 / TAK convention):
//   a-h-G   = General point (default)
//   a-f-A   = Friendly aircraft
//   a-f-G   = Friendly ground
//   a-n-G   = Neutral ground
//   a-h-A   = Hostile aircraft
//   b-m-p-s-m = Mission task marker
//
// Subject: cot.event — consumed by om-bridge-tak to render on TAK server.

// publishCoT sends a Cursor on Target event for a geo-bound task.
func publishCoT(ds DispatchState, cotType string) {
	if ds.Geo == nil {
		return
	}

	cot := map[string]interface{}{
		"type":       cotType,
		"lat":        ds.Geo.Lat,
		"lon":        ds.Geo.Lon,
		"hae":        ds.Geo.Alt,
		"ce":         "9999999.0", // circular error — precise
		"le":         "9999999.0", // linear error
		"uid":        fmt.Sprintf("TASK-%s", ds.TaskID),
		"callsign":   ds.TaskID,
		"how":        "m-g", // machine-generated
		"time":       nowISO(),
		"stale":      "3600", // 1 hour stale time
		"cot_type":   ds.Geo.COT,
	}

	// Include task context in remarks
	remarks := map[string]interface{}{
		"dispatch_id": ds.DispatchID,
		"task_id":     ds.TaskID,
		"status":      ds.Status,
		"source":      "omos-task-engine",
	}
	remarksJSON, _ := json.Marshal(remarks)
	cot["remarks"] = string(remarksJSON)

	data, _ := json.Marshal(cot)
	nc.Publish("cot.event", data)

	fmt.Printf("cot: published %s for task %s at %.6f,%.6f\n",
		cotType, ds.TaskID, ds.Geo.Lat, ds.Geo.Lon)
}

// publishCoTMarker publishes a stationary marker CoT event.
func publishCoTMarker(ds DispatchState) {
	publishCoT(ds, "marker")
}

// publishCoTTrack publishes a moving track CoT event.
func publishCoTTrack(ds DispatchState) {
	publishCoT(ds, "track")
}

// publishCoTComplete removes the task marker from the COP.
func publishCoTComplete(ds DispatchState) {
	if ds.Geo == nil {
		return
	}

	cot := map[string]interface{}{
		"type":     "marker",
		"lat":      ds.Geo.Lat,
		"lon":      ds.Geo.Lon,
		"uid":      fmt.Sprintf("TASK-%s", ds.TaskID),
		"callsign": ds.TaskID,
		"how":      "m-g",
		"time":     nowISO(),
		"stale":    "0", // immediate stale → removed from COP
		"remarks":  fmt.Sprintf(`{"status":"completed","task_id":"%s"}`, ds.TaskID),
	}

	data, _ := json.Marshal(cot)
	nc.Publish("cot.event", data)
}
