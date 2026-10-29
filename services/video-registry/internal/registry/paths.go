// Package registry links MediaMTX streams to OneMind entities and keeps the stream registry
// (KV OVERWATCH_VIDEO, evt.video.*, svc.registry.video.*) current. Video itself never travels on
// NATS; only its URLs, codecs, latency class and permissions do (BUS-DESIGN section 2).
package registry

import (
	"fmt"
	"regexp"
	"strings"
)

// classToType maps the first MediaMTX path segment to the ONTOLOGY entity_type that owns the
// stream. A drone-mounted camera belongs to the drone entity; a fixed camera is its own entity.
var classToType = map[string]string{
	"drone": "drone",
	"robot": "robot",
	"cam":   "camera",
}

// typeToClass is the inverse of classToType.
var typeToClass = map[string]string{
	"drone":  "drone",
	"robot":  "robot",
	"camera": "cam",
}

var idRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}$`)

// StreamRef is one MediaMTX path resolved to the entity it belongs to.
type StreamRef struct {
	StreamID   string // MediaMTX path, e.g. drone/sim-x1
	Class      string // drone | robot | cam
	EntityType string // drone | robot | camera
	EntityID   string // sim-x1
}

// Key is the KV OVERWATCH_VIDEO key: {entity_type}.{entity_id}.
func (s StreamRef) Key() string { return s.EntityType + "." + s.EntityID }

// EntitySubject is the entity's state subject in the locked grammar.
func (s StreamRef) EntitySubject() string { return "ent." + s.EntityType + "." + s.EntityID + ".state" }

// ParsePath resolves a MediaMTX path name. Drones and robots must be simulated assets (sim-*)
// in this run (run rule 6, decision 0004 #13); MediaMTX enforces the same with its path regexes.
func ParsePath(name string) (StreamRef, error) {
	parts := strings.Split(name, "/")
	if len(parts) != 2 {
		return StreamRef{}, fmt.Errorf("path %q: want {class}/{entity_id}", name)
	}
	typ, ok := classToType[parts[0]]
	if !ok {
		return StreamRef{}, fmt.Errorf("path %q: unknown class %q", name, parts[0])
	}
	if !idRe.MatchString(parts[1]) {
		return StreamRef{}, fmt.Errorf("path %q: invalid entity id", name)
	}
	if (typ == "drone" || typ == "robot") && !strings.HasPrefix(parts[1], "sim-") {
		return StreamRef{}, fmt.Errorf("path %q: only simulated assets (sim-*) are allowed", name)
	}
	return StreamRef{StreamID: name, Class: parts[0], EntityType: typ, EntityID: parts[1]}, nil
}

// ParseEntityRef accepts the forms a client may select an entity by:
// "ent.drone.sim-x1.state", "ent.drone.sim-x1", "drone.sim-x1" or the stream id "drone/sim-x1".
func ParseEntityRef(s string) (StreamRef, error) {
	s = strings.TrimSpace(s)
	if strings.Contains(s, "/") {
		return ParsePath(s)
	}
	t := strings.Split(s, ".")
	if len(t) >= 1 && t[0] == "ent" {
		t = t[1:]
		if len(t) == 3 && t[2] == "state" {
			t = t[:2]
		}
	}
	if len(t) != 2 {
		return StreamRef{}, fmt.Errorf("entity %q: want ent.{entity_type}.{entity_id}[.state] or {entity_type}.{entity_id}", s)
	}
	class, ok := typeToClass[t[0]]
	if !ok {
		return StreamRef{}, fmt.Errorf("entity %q: no video for entity_type %q", s, t[0])
	}
	return ParsePath(class + "/" + t[1])
}
