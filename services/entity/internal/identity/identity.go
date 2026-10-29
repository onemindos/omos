// Package identity resolves hard keys (ICAO hex, MMSI, NORAD, TAK uid, GTFS, sim-*) to one
// canonical entity_id through an alias map (fusion-and-tracking.md Stage 1; entity.md §3).
package identity

import (
	"crypto/sha256"
	"encoding/base32"
	"regexp"
	"strings"
	"sync"

	"github.com/onemindos/omos/services/entity/internal/model"
)

// validID is entity.md §3: lowercase, NATS-token safe, at most 64 characters.
var validID = regexp.MustCompile(`^[a-z0-9_-]{1,64}$`)

// ValidID reports whether id may be used as subject token three.
func ValidID(id string) bool { return validID.MatchString(id) }

// Key types, most specific first. The order decides which key names the entity.
const (
	KeySim         = "sim"
	KeyICAO        = "icao"
	KeyMMSI        = "mmsi"
	KeyNORAD       = "norad"
	KeyGTFS        = "gtfs"
	KeyCotUID      = "cot_uid"
	KeyMavlinkSys  = "mavlink_sysid"
	KeyDetTrack    = "det_track"
	KeyIntelID     = "intel_id"
	KeyFleetUID    = "fleet_uid"
	KeyTakCallsign = "tak_callsign"
)

// soft keys are recorded in aliases.ids but never resolve or index an entity: they are not
// unique (MAVLink sysid 1 is the ArduPilot and PX4 SITL default; two TAK devices can share a
// callsign). entity.md §5: auto-merge only on a shared hard key or a declared link.
var soft = map[string]bool{KeyMavlinkSys: true, KeyTakCallsign: true}

// Hard reports whether a key type may resolve and merge entities.
func Hard(keyType string) bool { return !soft[keyType] }

// prefixFor maps a key type to the entity_id prefix of entity.md §3.
var prefixFor = map[string]string{
	KeyICAO:     "icao",
	KeyMMSI:     "mmsi",
	KeyNORAD:    "norad",
	KeyGTFS:     "gtfs",
	KeyCotUID:   "tak",
	KeyDetTrack: "trk",
	KeyFleetUID: "fleet",
}

// MakeID builds `{prefix}-{native}` in lowercase, or `{prefix}-{base32(sha256(native))[0:26]}`
// when the native id has characters a NATS token or the command schema cannot carry.
func MakeID(prefix, native string) string {
	n := strings.ToLower(strings.TrimSpace(native))
	if prefix == "" {
		if ValidID(n) {
			return n
		}
		return hashed("id", native)
	}
	if id := prefix + "-" + n; ValidID(id) {
		return id
	}
	return hashed(prefix, native)
}

func hashed(prefix, native string) string {
	sum := sha256.Sum256([]byte(native))
	enc := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(sum[:])
	return prefix + "-" + strings.ToLower(enc[:26])
}

// Canonical returns the entity_id a single key implies.
func Canonical(k model.Alias) string {
	switch k.Type {
	case KeySim:
		return MakeID("", k.ID) // already `sim-...`
	case KeyIntelID:
		return MakeID("", k.ID) // already `evt-...` or `intel-...`
	}
	if p, ok := prefixFor[k.Type]; ok {
		return MakeID(p, k.ID)
	}
	return MakeID(k.Type, k.ID)
}

func aliasKey(k model.Alias) string { return k.Type + "|" + strings.ToLower(k.ID) }

// Resolver holds the alias map and the declared links (asset registry). It is safe for
// concurrent use.
type Resolver struct {
	mu    sync.RWMutex
	index map[string]string   // aliasKey -> entity_id (hard keys only)
	byID  map[string][]string // entity_id -> its aliasKeys in index, so Forget is O(keys)
	links map[string]string   // aliasKey -> entity_id, declared by the operator; wins over index
}

// NewResolver returns a resolver with declared links of the form "type:id" -> entity_id,
// for example "cot_uid:SIM-X1-ATAK" -> "sim-x1" (fusion-and-tracking.md Stage 4: auto-merge
// only on a shared hard key or a declared link).
func NewResolver(links map[string]string) *Resolver {
	r := &Resolver{index: map[string]string{}, byID: map[string][]string{}, links: map[string]string{}}
	for k, v := range links {
		t, id, ok := strings.Cut(k, ":")
		if !ok || !ValidID(v) {
			continue
		}
		r.links[aliasKey(model.Alias{Type: t, ID: id})] = v
	}
	return r
}

// Resolve returns the entity_id for the keys. Declared links win, then an existing alias of a
// hard key, then the canonical id of the first hard key. Soft keys never resolve. ok is false
// when keys holds no hard key and no declared link.
func (r *Resolver) Resolve(keys []model.Alias) (id string, ok bool) {
	if len(keys) == 0 {
		return "", false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, k := range keys {
		if v, hit := r.links[aliasKey(k)]; hit {
			return v, true
		}
	}
	for _, k := range keys {
		if !Hard(k.Type) {
			continue
		}
		if v, hit := r.index[aliasKey(k)]; hit {
			return v, true
		}
	}
	for _, k := range keys {
		if Hard(k.Type) {
			return Canonical(k), true
		}
	}
	return "", false
}

// Register records that every hard key in keys names entity id. An existing mapping is kept: a
// key never silently moves to another entity (wrong merges are worse than duplicates, Risk 2).
// Soft keys are not indexed.
func (r *Resolver) Register(id string, keys []model.Alias) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, k := range keys {
		if !Hard(k.Type) {
			continue
		}
		ak := aliasKey(k)
		if _, exists := r.index[ak]; !exists {
			r.index[ak] = id
			r.byID[id] = append(r.byID[id], ak)
		}
	}
}

// Forget removes every alias that points at id (on expiry).
func (r *Resolver) Forget(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ak := range r.byID[id] {
		if r.index[ak] == id {
			delete(r.index, ak)
		}
	}
	delete(r.byID, id)
}
