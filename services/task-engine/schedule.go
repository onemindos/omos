package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// ── schedule.go — recurring task engine ────────────────────────────────────
//
// Tasks with a recurring schedule ("rrule" or "cron" field) are automatically
// dispatched by the schedule loop. Every 60 seconds, the engine scans for
// tasks whose next run time has arrived, creates a child dispatch from the
// parent template, and publishes it.
//
// Supported formats:
//   - cron: standard 5-field cron expression
//   - rrule: RFC 5545 RRULE (basic subset)
//   - interval: "30m", "2h", "1d" — simple duration format

type ScheduleEntry struct {
	TaskID     string    `json:"task_id"`
	NextRun    time.Time `json:"next_run"`
	CronExpr   string    `json:"cron_expr"`
	LastRun    string    `json:"last_run,omitempty"`
	RunCount   int       `json:"run_count"`
}

var cronScheduler = cron.New(cron.WithSeconds())

// scheduleLoop runs the recurring task engine.
func scheduleLoop() {
	fmt.Println("schedule: engine starting (60s scan interval)")

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Initial scan
	scanRecurringTasks()

	for range ticker.C {
		scanRecurringTasks()
	}
}

// scanRecurringTasks finds all tasks with a recurring schedule and dispatches
// any whose next run time has arrived.
func scanRecurringTasks() {
	keys := allTaskKeys()
	now := time.Now().UTC()

	for _, k := range keys {
		task, err := getTask(k)
		if err != nil {
			continue
		}

		// Only process active recurring tasks
		status := getStr(task, "status")
		if status == "cancelled" || status == "completed" {
			continue
		}

		rrule := getStr(task, "rrule")
		cronExpr := getStr(task, "cron")
		interval := getStr(task, "interval")

		nextRunStr := getStr(task, "next_run")

		if rrule == "" && cronExpr == "" && interval == "" {
			continue // not a recurring task
		}

		// Parse next_run or compute it
		var nextRun time.Time
		if nextRunStr != "" {
			nextRun, err = time.Parse(time.RFC3339, nextRunStr)
			if err != nil {
				nextRun, _ = time.Parse("2006-01-02T15:04:05Z", nextRunStr)
			}
		}

		// If no next_run or it's in the past, schedule it
		if nextRun.IsZero() || now.After(nextRun) || now.Equal(nextRun) {
			if !nextRun.IsZero() && now.Before(nextRun) {
				continue // not yet time
			}

			// Dispatch this recurring task
			dispatchRecurringTask(task)

			// Compute next run time
			newNextRun := computeNextRun(now, cronExpr, interval)
			task["next_run"] = newNextRun.Format(time.RFC3339)
			task["last_run"] = now.Format(time.RFC3339)
			runCount := int(getFloat(task, "run_count"))
			task["run_count"] = float64(runCount + 1)
			task["updated_at"] = nowISO()
			putTask(k, task)

			fmt.Printf("schedule: dispatched recurring task %s, next run %s\n",
				k, newNextRun.Format(time.RFC3339))
		}
	}
}

// dispatchRecurringTask creates a child dispatch from a recurring task template.
func dispatchRecurringTask(task map[string]interface{}) {
	taskID := getStr(task, "task_id")
	assignee := getStr(task, "assignee")
	assigneeType := getStr(task, "assignee_type")

	if assignee == "" {
		fmt.Printf("schedule: task %s has no assignee, skipping dispatch\n", taskID)
		return
	}

	actionMap := getMap(task, "action")
	if actionMap == nil {
		actionMap = map[string]interface{}{
			"type":    "nats_pub",
			"subject": fmt.Sprintf("cmd.%s", assignee),
		}
	}

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

	ds := DispatchState{
		DispatchID:   dispatchID,
		TaskID:       taskID,
		Action:       action,
		Expect:       expect,
		Geo:          geo,
		Status:       "pending",
		DispatchedAt: nowISO(),
	}

	// Build dispatch command
	subject := expandSubject(action.Subject, assignee, taskID)
	cmd := map[string]interface{}{
		"dispatch_id":  dispatchID,
		"task_id":      taskID,
		"action":       action,
		"assignee":     assignee,
		"recurring":    true,
		"run_count":    int(getFloat(task, "run_count")) + 1,
		"timestamp":    nowISO(),
	}

	cmdData, _ := json.Marshal(cmd)

	// Execute
	if assigneeType != "human" {
		switch action.Type {
		case "nats_req":
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
			if err := nc.Publish(subject, cmdData); err != nil {
				ds.Status = "failed"
				ds.Result = map[string]interface{}{"error": err.Error()}
			} else {
				ds.Status = "dispatched"
			}
		default:
			if err := nc.Publish(subject, cmdData); err != nil {
				ds.Status = "failed"
				ds.Result = map[string]interface{}{"error": err.Error()}
			} else {
				ds.Status = "dispatched"
			}
		}
	}

	// Persist dispatch
	stateData, _ := json.Marshal(ds)
	dispatchKV.Put(dispatchID, stateData)

	// Publish to MISSION stream
	publishMissionEvent("mission.recurring", ds)

	// Publish CoT if geo-bound
	if ds.Geo != nil {
		publishCoTMarker(ds)
	}

	// Watch for response
	if ds.Expect.SubjectTemplate != "" && ds.Status == "dispatched" {
		go watchResponse(ds)
	}
}

// computeNextRun calculates the next run time from a cron expression or interval.
func computeNextRun(now time.Time, cronExpr, interval string) time.Time {
	// Try cron expression first
	if cronExpr != "" {
		parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
		sched, err := parser.Parse(cronExpr)
		if err == nil {
			return sched.Next(now)
		}
	}

	// Try interval format: "30m", "2h", "1d"
	if interval != "" {
		d, err := time.ParseDuration(interval)
		if err == nil {
			return now.Add(d)
		}

		// Handle "1d" style (Go doesn't parse "d")
		if len(interval) > 1 && interval[len(interval)-1] == 'd' {
			var days int
			fmt.Sscanf(interval, "%dd", &days)
			if days > 0 {
				return now.Add(time.Duration(days) * 24 * time.Hour)
			}
		}
	}

	// Default: tomorrow
	return now.Add(24 * time.Hour)
}
