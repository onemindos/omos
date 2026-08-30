// Package actor defines the Actor registration shape and helpers for
// reading/writing the actor record to the om-config KV bucket.
//
// Capabilities for the skeleton build: ["host.info"] only.
// fs.read / fs.write / shell.exec are capability-gated and will be added
// in a later, reviewed card. Do NOT add them here.
package actor

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// Actor is the JSON shape stored under actor.node.{id} in om-config.
// It mirrors the ARCHITECTURE.md §3 Actor model.
type Actor struct {
	Type         string       `json:"type"`
	ID           string       `json:"id"`
	Display      Display      `json:"display"`
	Capabilities []string     `json:"capabilities"`
	Endpoint     string       `json:"endpoint"`
	Executor     string       `json:"executor"`
	Status       string       `json:"status"`
	Meta         NodeMeta     `json:"meta"`
	LastSeen     string       `json:"last_seen,omitempty"`
}

// Display holds human-readable identity fields.
type Display struct {
	Name string `json:"name"`
}

// NodeMeta holds node-specific runtime metadata published at registration time.
type NodeMeta struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Hostname string `json:"hostname"`
	Started  string `json:"started"`
}

// New builds an Actor record for this node.
func New(nodeID, hostname string) *Actor {
	return &Actor{
		Type: "node",
		ID:   nodeID,
		Display: Display{
			Name: hostname,
		},
		// SKELETON only — [\"host.info\"] is the only advertised capability.
		// fs.read, fs.write, shell.exec require a later, capability-gated card.
		Capabilities: []string{"host.info"},
		Endpoint:     fmt.Sprintf("cmd.node.%s.>", nodeID),
		Executor:     "TRIGGER",
		Status:       "online",
		Meta: NodeMeta{
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Hostname: hostname,
			Started:  time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// KVKey returns the om-config bucket key for this actor.
func KVKey(nodeID string) string {
	return fmt.Sprintf("actor.node.%s", nodeID)
}

// Marshal serialises the actor to JSON bytes.
func (a *Actor) Marshal() ([]byte, error) {
	return json.Marshal(a)
}

// SetStatus clones the actor with a new status (online/offline).
func (a *Actor) SetStatus(status string) *Actor {
	copy := *a
	copy.Status = status
	copy.LastSeen = time.Now().UTC().Format(time.RFC3339)
	return &copy
}
