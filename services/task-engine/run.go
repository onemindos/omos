package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

func mustKV(bucket string, ttl time.Duration) nats.KeyValue {
	kv, err := js.KeyValue(bucket)
	if err != nil {
		cfg := &nats.KeyValueConfig{Bucket: bucket}
		if ttl > 0 {
			cfg.TTL = ttl
		}
		kv, err = js.CreateKeyValue(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "kv %s: %v\n", bucket, err)
			os.Exit(1)
		}
	}
	return kv
}

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://127.0.0.1:4222"
	}
	name := os.Getenv("NATS_NAME")
	if name == "" {
		name = "omos-task-engine"
	}

	var err error
	nc, err = nats.Connect(natsURL, nats.Name(name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}
	defer nc.Close()

	js, err = nc.JetStream()
	if err != nil {
		fmt.Fprintf(os.Stderr, "jetstream: %v\n", err)
		os.Exit(1)
	}

	tasksKV = mustKV("task-index", 30*24*time.Hour)
	projectsKV = mustKV("projects", 0)
	dispatchKV = mustKV("dispatch-state", 7*24*time.Hour)

	svc, err := micro.AddService(nc, micro.Config{
		Name:        "omos-task-engine",
		Version:     "4.0.0",
		Description: "OneMind OS Sovereign Task Orchestration Engine",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "service: %v\n", err)
		os.Exit(1)
	}

	reg := func(name, subject string, h micro.HandlerFunc) {
		if e := svc.AddEndpoint(name, h, micro.WithEndpointSubject(subject)); e != nil {
			fmt.Fprintf(os.Stderr, "endpoint %s: %v\n", subject, e)
		}
	}

	// Task CRUD
	reg("task-create", "task.create", handleCreate)
	reg("task-list", "task.list", handleList)
	reg("task-get", "task.get.*", handleGet)
	reg("task-update", "task.update.*", handleUpdate)
	reg("task-approve", "task.approve.*", handleApprove)
	reg("task-delete", "task.delete.*", handleDelete)
	reg("task-calendar", "task.calendar", handleCalendar)
	reg("task-inbox", "task.inbox", handleInbox)

	// Project CRUD
	reg("task-project-create", "task.project.create", handleProjectCreate)
	reg("task-project-list", "task.project.list", handleProjectList)

	// Syntax
	reg("task-syntax-validate", "task.syntax.validate", handleSyntaxValidate)

	// ── NEW v4.0: Orchestration endpoints ──
	reg("dispatch-execute", "task.dispatch.execute", handleDispatchExecute)
	reg("dispatch-status", "task.dispatch.status.*", handleDispatchStatus)

	fmt.Println("omos-task-engine v4.0.0 ready — sovereign orchestration fabric")
	fmt.Printf("  NATS: %s | KV: task-index, projects, dispatch-state\n", natsURL)

	// Start the schedule engine
	go scheduleLoop()

	select {}
}
