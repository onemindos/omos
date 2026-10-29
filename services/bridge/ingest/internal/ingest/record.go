package ingest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

var (
	errParse          = errors.New("parse")
	errCheckpointRead = errors.New("checkpoint read")
)

// Record is one normalised upstream observation, before validation and publish.
type Record struct {
	Subject string         // grammar subject; empty means the source's configured subject
	ID      string         // dedup key: becomes Nats-Msg-Id. Derived from upstream ids only, never local time
	TS      time.Time      // upstream event time (envelope ts)
	Data    map[string]any // payload, validated against the source schema
	Raw     []byte         // the upstream bytes, kept for the dead-letter record
	Err     error          // set by the parser when the upstream item is unusable
}

// Parser turns one upstream body (a poll response, one SSE event, one WebSocket frame)
// into records. It must be pure: the same bytes give the same records (golden tests).
type Parser func(body []byte, opts map[string]string) ([]Record, error)

var (
	parsersMu sync.RWMutex
	parsers   = map[string]Parser{}
)

// RegisterParser makes a parser available to sources.yaml by name.
func RegisterParser(name string, p Parser) {
	parsersMu.Lock()
	defer parsersMu.Unlock()
	parsers[name] = p
}

func LookupParser(name string) (Parser, bool) {
	parsersMu.RLock()
	defer parsersMu.RUnlock()
	p, ok := parsers[name]
	return p, ok
}

func ParserNames() []string {
	parsersMu.RLock()
	defer parsersMu.RUnlock()
	var out []string
	for k := range parsers {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ParseError wraps a whole-body parse failure.
func ParseError(err error) error { return fmt.Errorf("%w: %v", errParse, err) }

// Envelope is the omos bus envelope (services/schemas/envelope/envelope.schema.json).
type Envelope struct {
	V         int            `json:"v"`
	ID        string         `json:"id"`
	TS        string         `json:"ts"`
	Type      string         `json:"type"`
	Source    string         `json:"source"`
	Subject   string         `json:"subject"`
	Data      map[string]any `json:"data"`
	FetchedAt string         `json:"fetchedAt,omitempty"`
	MsgID     string         `json:"msgId,omitempty"`
}

const tsLayout = "2006-01-02T15:04:05.000Z07:00"

// NewEnvelope builds the envelope. Its id is a UUIDv4-shaped hash of the dedup key, so a
// replay of the same upstream item produces the same envelope id.
func NewEnvelope(sourceID, subject, msgID string, ts, fetchedAt time.Time, data map[string]any) Envelope {
	if ts.IsZero() {
		ts = fetchedAt
	}
	e := Envelope{V: 1, ID: StableUUID(msgID), TS: ts.UTC().Format(tsLayout), Type: subject,
		Source: "ingest." + sourceID, Subject: subject, Data: data, MsgID: msgID}
	if !fetchedAt.IsZero() {
		e.FetchedAt = fetchedAt.UTC().Format(tsLayout)
	}
	return e
}

func StableUUID(key string) string {
	h := sha256.Sum256([]byte(key))
	h[6] = (h[6] & 0x0f) | 0x40
	h[8] = (h[8] & 0x3f) | 0x80
	x := hex.EncodeToString(h[:16])
	return x[0:8] + "-" + x[8:12] + "-" + x[12:16] + "-" + x[16:20] + "-" + x[20:32]
}

// Schemas compiles JSON Schemas from the schemas directory (services/schemas in the image).
type Schemas struct {
	dir   string
	mu    sync.Mutex
	cache map[string]*jsonschema.Schema
}

func NewSchemas(dir string) *Schemas {
	return &Schemas{dir: dir, cache: map[string]*jsonschema.Schema{}}
}

func (s *Schemas) Get(rel string) (*jsonschema.Schema, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sch, ok := s.cache[rel]; ok {
		return sch, nil
	}
	abs, err := filepath.Abs(filepath.Join(s.dir, rel))
	if err != nil {
		return nil, err
	}
	sch, err := jsonschema.NewCompiler().Compile(abs)
	if err != nil {
		return nil, fmt.Errorf("schema %s: %w", rel, err)
	}
	s.cache[rel] = sch
	return sch, nil
}

// Validate checks data against the schema by round-tripping through JSON, so the check
// sees exactly the bytes that will be published.
func Validate(sch *jsonschema.Schema, data map[string]any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	inst, err := jsonschema.UnmarshalJSON(bytesReader(b))
	if err != nil {
		return err
	}
	return sch.Validate(inst)
}
