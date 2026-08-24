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

	var err error
	nc, err = nats.Connect(natsURL, nats.Name("om-task"))
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

	svc, err := micro.AddService(nc, micro.Config{
		Name:        "om-task",
		Version:     "3.1.0-nex",
		Description: "OneMind Task Engine (NEX, full contract)",
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
	reg("task-create", "task.create", handleCreate)
	reg("task-list", "task.list", handleList)
	reg("task-get", "task.get.*", handleGet)
	reg("task-update", "task.update.*", handleUpdate)
	reg("task-approve", "task.approve.*", handleApprove)
	reg("task-delete", "task.delete.*", handleDelete)
	reg("task-calendar", "task.calendar", handleCalendar)
	reg("task-inbox", "task.inbox", handleInbox)
	reg("task-project-create", "task.project.create", handleProjectCreate)
	reg("task-project-list", "task.project.list", handleProjectList)
	reg("task-syntax-validate", "task.syntax.validate", handleSyntaxValidate)

	fmt.Println("om-task v3.1.0-nex ready — full task.* contract")
	select {}
}
