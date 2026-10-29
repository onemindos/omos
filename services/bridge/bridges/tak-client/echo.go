// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// echoFilter remembers CoT recently published to tak.cot.inject (by anyone:
// the projector through this client or through NatsBridgePlugin, or a legacy
// feed) so the same event coming back out of TAK Server is not republished on
// tak.cot.out as if a TAK user had sent it. Keyed on uid, type and time, which
// TAK Server relays unchanged.
type echoFilter struct {
	mu   sync.Mutex
	ttl  time.Duration
	max  int
	seen map[string]time.Time
}

func newEchoFilter(ttl time.Duration, max int) *echoFilter {
	return &echoFilter{ttl: ttl, max: max, seen: make(map[string]time.Time)}
}

func echoKey(ev *Event) string {
	return ev.UID + "|" + ev.Type + "|" + strconv.FormatUint(ev.SendTime, 10)
}

// Remember records an injected event.
func (f *echoFilter) Remember(ev *Event, now time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.seen) >= f.max {
		f.pruneLocked(now)
		if len(f.seen) >= f.max { // still full: drop everything rather than grow
			f.seen = make(map[string]time.Time)
		}
	}
	f.seen[echoKey(ev)] = now.Add(f.ttl)
}

// IsEcho reports whether ev was injected within the TTL.
func (f *echoFilter) IsEcho(ev *Event, now time.Time) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	exp, ok := f.seen[echoKey(ev)]
	return ok && now.Before(exp)
}

func (f *echoFilter) pruneLocked(now time.Time) {
	for k, exp := range f.seen {
		if !now.Before(exp) {
			delete(f.seen, k)
		}
	}
}

// Len is the number of remembered keys (for metrics).
func (f *echoFilter) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.seen)
}

// dedupID is the Nats-Msg-Id for tak.cot.out: the same fields as the mesh
// gateway's dedupID (bridges/tak/mesh.go) with a distinct "takc-" prefix, so a
// reconnect replaying the same event is dropped by stream COT's duplicate
// window, while the same EUD heard on the mesh (tak.cot.mesh.rx, "cot-") and
// through TAK Server stays two observations for the entity service to merge.
func dedupID(ev *Event) string {
	sum := sha256.Sum256([]byte(ev.UID + "|" + ev.Type + "|" + strconv.FormatUint(ev.SendTime, 10) + "|" +
		strconv.FormatFloat(ev.Lat, 'f', -1, 64) + "|" + strconv.FormatFloat(ev.Lon, 'f', -1, 64)))
	return fmt.Sprintf("takc-%x", sum[:16])
}
