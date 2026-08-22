package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/micro"
)

// ── dispatch.go — orchestration fabric ─────────────────────────────────────
//
// When a task moves to "executing", the TaskAction is published to the actor's
// NATS subject. The actor (drone, robot, agent, service) receives the command,
// executes, and reports back on its telemetry subject. ActionExpect watches
// for the response and updates task state.
//
// TaskAction types:
//   - nats_pub    — fire-and-forget publish to a subject
//   - nats_req    — request-reply with timeout
//   - agent_prompt — delegate to an AI agent (Synadia Agent Protocol)
//   - service_call — invoke an internal NATS micro service
//
// Assignee types:
//   - human       (manual, no auto-dispatch)
//   - robot/drone/node/agent/service (auto-dispatchable)

type TaskAction struct {
	Type    string                 `json:"type"`    // nats_pub | nats_req | agent_prompt | service_call
	Subject string                 `json:"subject"` // NATS subject to publish to
	Payload map[string]interface{} `json:"payload"` // JSON payload
	Timeout string                 `json:"timeout"` // e.g. "30s", "5m"
}

type ActionExpect struct {
	SubjectTemplate string `json:"subject_template"` // e.g. "tel.{assignee}.mission.status"
	TimeoutSecs     int    `json:"timeout_secs"`      // seconds to wait
	SuccessField    string `json:"success_field"`     // JSON path to check for success
	SuccessValue    string `json:"success_value"`     // expected value
}

type GeoBinding struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Alt  float64 `json:"alt,omitempty"`
	COT  string  `json:"cot_type"`  // CoT event type
	Icon string  `json:"icon,omitempty"` // TAK icon reference
}

type DispatchState struct {
	DispatchID  string                 `json:"dispatch_id"`
	TaskID      string                 `json:"task_id"`
	Action      TaskAction             `json:"action"`
	Expect      ActionExpect           `json:"expect,omitempty"`
	Geo         *GeoBinding            `json:"geo,omitempty"`
	Status      string                 `json:"status"` // pending | dispatched | executing | succeeded | failed | timed_out
	Result      map[string]interface{} `json:"result,omitempty"`
	DispatchedAt string                `json:"dispatched_at"`
	CompletedAt  string                `json:"completed_at,omitempty"`
}

// handleDispatchExecute fires a task's action to its actor subject.
func handleDispatchExecute(req micro.Request) {
	var task map[string]interface{}
	if err := json.Unmarshal(req.Data(), &task); err != nil {
		respErr(req, "invalid task JSON")
		return
	}

	assigneeType := getStr(task, "assignee_type")
	if assigneeType == "" {
		assigneeType = "human"
	}
	assignee := getStr(task, "assignee")
	if assignee == "" {
		respErr(req, "missing assignee")
		return
	}

	actionMap := getMap(task, "action")
	if actionMap == nil {
		respErr(req, "missing action block")
		return
	}

	// Build TaskAction
	action := TaskAction{
		Type:    getStr(actionMap, "type"),
		Subject: getStr(actionMap, "subject"),
		Payload: actionMap,
		Timeout: getStr(actionMap, "timeout"),
	}
	if action.Type == "" {
		action.Type = "nats_pub"
	}
	if action.Timeout == "" {
		action.Timeout = "60s"
	}

	// Build ActionExpect if present
	var expect ActionExpect
	if expMap := getMap(task, "expect"); expMap != nil {
		expect = ActionExpect{
			SubjectTemplate: getStr(expMap, "subject_template"),
			TimeoutSecs:     int(getFloat(expMap, "timeout_secs")),
			SuccessField:    getStr(expMap, "success_field"),
			SuccessValue:    getStr(expMap, "success_value"),
		}
		if expect.TimeoutSecs == 0 {
			expect.TimeoutSecs = 120
		}
	}

	// Build GeoBinding if present
	var geo *GeoBinding
	if geoMap := getMap(task, "geo"); geoMap != nil {
		geo = &GeoBinding{
			Lat:  getFloat(geoMap, "lat"),
			Lon:  getFloat(geoMap, "lon"),
			Alt:  getFloat(geoMap, "alt"),
			COT:  getStr(geoMap, "cot_type"),
			Icon: getStr(geoMap, "icon"),
		}
		if geo.COT == "" {
			geo.COT = "a-h-G"
		}
	}

	dispatchID := "dsp-" + shortUUID()
	taskID := getStr(task, "task_id")

	ds := DispatchState{
		DispatchID:   dispatchID,
		TaskID:       taskID,
		Action:       action,
		Expect:       expect,
		Geo:          geo,
		Status:       "pending",
		DispatchedAt: nowISO(),
	}

	// Resolve actual subject — template expansion
	subject := expandSubject(action.Subject, assignee, taskID)

	// Build the dispatch command payload
	cmd := map[string]interface{}{
		"dispatch_id": dispatchID,
		"task_id":     taskID,
		"action":      action,
		"assignee":    assignee,
		"task":        task,
		"timestamp":   nowISO(),
	}

	cmdData, _ := json.Marshal(cmd)

	switch action.Type {
	case "nats_req":
		// Request-reply with timeout
		timeout, _ := time.ParseDuration(action.Timeout)
		msg, err := nc.Request(subject, cmdData, timeout)
		if err != nil {
			ds.Status = "failed"
			ds.Result = map[string]interface{}{"error": err.Error()}
		} else {
			ds.Status = "dispatched"
			var result map[string]interface{}
			if json.Unmarshal(msg.Data, &result) == nil {
				ds.Result = result
			}
			ds.CompletedAt = nowISO()
		}
	case "nats_pub":
		// Fire-and-forget
		if err := nc.Publish(subject, cmdData); err != nil {
			ds.Status = "failed"
			ds.Result = map[string]interface{}{"error": err.Error()}
		} else {
			ds.Status = "dispatched"
		}
	case "agent_prompt":
		// Delegate via Synadia Agent Protocol — publish to agent subject
		agentSubject := fmt.Sprintf("agents.prompt.hermes.%s", assignee)
		prompt := map[string]interface{}{
			"dispatch_id": dispatchID,
			"task_id":     taskID,
			"prompt":      getStr(actionMap, "prompt"),
		}
		promptData, _ := json.Marshal(prompt)
		if err := nc.Publish(agentSubject, promptData); err != nil {
			ds.Status = "failed"
			ds.Result = map[string]interface{}{"error": err.Error()}
		} else {
			ds.Status = "dispatched"
		}
	case "service_call":
		// Internal NATS micro service
		msg, err := nc.Request(subject, cmdData, 10*time.Second)
		if err != nil {
			ds.Status = "failed"
			ds.Result = map[string]interface{}{"error": err.Error()}
		} else {
			ds.Status = "dispatched"
			var result map[string]interface{}
			if json.Unmarshal(msg.Data, &result) == nil {
				ds.Result = result
			}
			ds.CompletedAt = nowISO()
		}
	}

	// Persist dispatch state
	stateData, _ := json.Marshal(ds)
	dispatchKV.Put(dispatchID, stateData)

	// Update task status if dispatch was to an actor
	if assigneeType != "human" && ds.Status == "dispatched" {
		task["status"] = "executing"
		task["dispatch_id"] = dispatchID
		task["updated_at"] = nowISO()
		putTask(taskID, task)
	}

	// Publish to MISSION stream for COP visibility
	publishMissionEvent("mission.dispatch", ds)

	// If expecting a response, start the response watcher
	if ds.Expect.SubjectTemplate != "" && ds.Status == "dispatched" {
		go watchResponse(ds)
	}

	respJSON(req, ds)
}

// handleDispatchStatus returns the current state of a dispatch.
func handleDispatchStatus(req micro.Request) {
	_ = req // subject routing handled by framework, but we don't use the id from subject
	dsID := ""
	if len(req.Subject()) > len("task.dispatch.status.") {
		dsID = req.Subject()[len("task.dispatch.status."):]
	}

	entry, err := dispatchKV.Get(dsID)
	if err != nil {
		respErr(req, fmt.Sprintf("dispatch not found: %s", dsID))
		return
	}
	var ds DispatchState
	if err := json.Unmarshal(entry.Value(), &ds); err != nil {
		respErr(req, "invalid dispatch state")
		return
	}
	respJSON(req, ds)
}

// expandSubject replaces template vars: {assignee}, {task_id}, {node}
func expandSubject(tmpl, assignee, taskID string) string {
	s := tmpl
	if s == "" {
		// Default: assignee-based command subject
		s = fmt.Sprintf("cmd.%s", assignee)
	}
	// Replace known templates
	s = replaceAll(s, "{assignee}", assignee)
	s = replaceAll(s, "{task_id}", taskID)
	s = replaceAll(s, "{node}", assignee)
	return s
}

func replaceAll(s, old, new string) string {
	out := s
	for i := 0; i < len(out); i++ {
		if i+len(old) <= len(out) && out[i:i+len(old)] == old {
			out = out[:i] + new + out[i+len(old):]
			i += len(new) - 1
		}
	}
	return out
}

// publishMissionEvent publishes a structured event to the MISSION stream for
// COP visibility in TAK and the mission log.
func publishMissionEvent(eventType string, ds DispatchState) {
	evt := map[string]interface{}{
		"type":        eventType,
		"dispatch_id": ds.DispatchID,
		"task_id":     ds.TaskID,
		"action":      ds.Action,
		"status":      ds.Status,
		"timestamp":   nowISO(),
	}
	if ds.Geo != nil {
		evt["geo"] = map[string]interface{}{
			"lat":      ds.Geo.Lat,
			"lon":      ds.Geo.Lon,
			"alt":      ds.Geo.Alt,
			"cot_type": ds.Geo.COT,
		}
	}
	data, _ := json.Marshal(evt)
	nc.Publish("evt.mission.dispatch", data)
}
