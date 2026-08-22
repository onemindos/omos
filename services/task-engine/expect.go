package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ── expect.go — ActionExpect response watching ─────────────────────────────
//
// After a dispatch is fired, the engine watches for the executor's response
// on a telemetry or status subject. When the response arrives (or times out),
// the dispatch state and task state are updated.

var (
	watchersMu sync.Mutex
	watchers   = make(map[string]chan bool) // dispatchID -> cancel
)

// watchResponse subscribes to the actor's telemetry subject and waits for
// a response matching the ActionExpect criteria. If the actor responds with
// success, the task is marked completed. If it times out, the task fails.
func watchResponse(ds DispatchState) {
	if ds.Expect.SubjectTemplate == "" {
		return
	}

	watchSubject := expandSubject(ds.Expect.SubjectTemplate,
		getStr(ds.Result, "assignee"),
		ds.TaskID)

	timeout := time.Duration(ds.Expect.TimeoutSecs) * time.Second
	if timeout == 0 {
		timeout = 120 * time.Second
	}

	cancel := make(chan bool, 1)
	watchersMu.Lock()
	watchers[ds.DispatchID] = cancel
	watchersMu.Unlock()

	defer func() {
		watchersMu.Lock()
		delete(watchers, ds.DispatchID)
		watchersMu.Unlock()
	}()

	// Subscribe with wildcard match
	sub, err := nc.SubscribeSync(watchSubject)
	if err != nil {
		fmt.Printf("dispatch %s: subscribe failed on %s: %v\n", ds.DispatchID, watchSubject, err)
		markDispatchFailed(ds, fmt.Sprintf("subscribe failed: %v", err))
		return
	}
	defer sub.Unsubscribe()

	// Set auto-unsubscribe after 1 message (we only need one)
	sub.AutoUnsubscribe(1)

	deadline := time.After(timeout)

	for {
		select {
		case <-cancel:
			return
		case <-deadline:
			markDispatchFailed(ds, "timeout waiting for executor response")
			return
		default:
			msg, err := sub.NextMsg(1 * time.Second)
			if err != nil {
				if err.Error() == "nats: timeout" {
					continue
				}
				// Subscription closed — likely timed out via AutoUnsubscribe
				// but we still have deadline checking
				continue
			}
			if msg == nil {
				continue
			}

			// Parse the response
			var resp map[string]interface{}
			if err := json.Unmarshal(msg.Data, &resp); err != nil {
				continue // ignore malformed
			}

			// Check success criteria
			if checkSuccess(resp, ds.Expect.SuccessField, ds.Expect.SuccessValue) {
				markDispatchSuccess(ds, resp)
				return
			}

			// If response indicates failure, mark it
			if errStr := getStr(resp, "error"); errStr != "" {
				markDispatchFailed(ds, fmt.Sprintf("executor error: %s", errStr))
				return
			}
		}
	}
}

// checkSuccess checks a JSON path-like field against an expected value.
// Simple dot-notation: "status" or "result.ok"
func checkSuccess(resp map[string]interface{}, field, expected string) bool {
	if field == "" {
		return false
	}
	parts := strings.Split(field, ".")
	var val interface{} = resp
	for _, part := range parts {
		m, ok := val.(map[string]interface{})
		if !ok {
			return false
		}
		val = m[part]
	}
	// Compare
	return fmt.Sprintf("%v", val) == expected
}

func markDispatchSuccess(ds DispatchState, result map[string]interface{}) {
	ds.Status = "succeeded"
	ds.Result = result
	ds.CompletedAt = nowISO()

	stateData, _ := json.Marshal(ds)
	dispatchKV.Put(ds.DispatchID, stateData)

	// Update task
	task, err := getTask(ds.TaskID)
	if err == nil {
		task["status"] = "completed"
		task["completed_at"] = nowISO()
		task["updated_at"] = nowISO()
		task["dispatch_result"] = result
		putTask(ds.TaskID, task)
	}

	publishMissionEvent("mission.completed", ds)
	fmt.Printf("dispatch %s: succeeded\n", ds.DispatchID)
}

func markDispatchFailed(ds DispatchState, reason string) {
	ds.Status = "failed"
	ds.Result = map[string]interface{}{"error": reason}
	ds.CompletedAt = nowISO()

	stateData, _ := json.Marshal(ds)
	dispatchKV.Put(ds.DispatchID, stateData)

	// Update task
	task, err := getTask(ds.TaskID)
	if err == nil {
		task["status"] = "failed"
		task["updated_at"] = nowISO()
		task["dispatch_error"] = reason
		putTask(ds.TaskID, task)
	}

	publishMissionEvent("mission.failed", ds)
	fmt.Printf("dispatch %s: failed — %s\n", ds.DispatchID, reason)
}
