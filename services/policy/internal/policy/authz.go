package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// CheckRequest is one OpenFGA Check: does user have relation on object?
type CheckRequest struct {
	User     string         `json:"user"`
	Relation string         `json:"relation"`
	Object   string         `json:"object"`
	Context  map[string]any `json:"context,omitempty"`
	// ContextualTuples hold for this Check only (OpenFGA contextual_tuples). The policy
	// service sends a verified callout user's role groups this way (FEATURES J2), so group
	// membership follows the Authentik token and needs no sync job.
	ContextualTuples []Tuple `json:"contextual_tuples,omitempty"`
}

// Authorizer answers OpenFGA Check questions. Any error means "deny" to the caller
// (fail closed, security-and-trust.md section 4).
type Authorizer interface {
	Check(ctx context.Context, req CheckRequest) (bool, error)
	StoreID() string
	ModelID() string
}

// ---------------------------------------------------------------------------
// In-process evaluator for the OpenFGA JSON model (subset: this, computedUserset,
// tupleToUserset, union, type restrictions, the within_window condition).

type rewrite struct {
	This            *struct{} `json:"this,omitempty"`
	ComputedUserset *struct {
		Relation string `json:"relation"`
	} `json:"computedUserset,omitempty"`
	TupleToUserset *struct {
		Tupleset struct {
			Relation string `json:"relation"`
		} `json:"tupleset"`
		ComputedUserset struct {
			Relation string `json:"relation"`
		} `json:"computedUserset"`
	} `json:"tupleToUserset,omitempty"`
	Union *struct {
		Child []rewrite `json:"child"`
	} `json:"union,omitempty"`
}

type relationRef struct {
	Type      string `json:"type"`
	Relation  string `json:"relation,omitempty"`
	Condition string `json:"condition,omitempty"`
}

type typeDef struct {
	Type      string             `json:"type"`
	Relations map[string]rewrite `json:"relations"`
	Metadata  *struct {
		Relations map[string]struct {
			Direct []relationRef `json:"directly_related_user_types"`
		} `json:"relations"`
	} `json:"metadata"`
}

// AuthorizationModel is the OpenFGA API form (POST /stores/{id}/authorization-models body).
type AuthorizationModel struct {
	SchemaVersion   string                     `json:"schema_version"`
	TypeDefinitions []typeDef                  `json:"type_definitions"`
	Conditions      map[string]json.RawMessage `json:"conditions"`
}

// Tuple is a relationship tuple, optionally conditional.
type Tuple struct {
	User      string         `json:"user"`
	Relation  string         `json:"relation"`
	Object    string         `json:"object"`
	Condition string         `json:"condition,omitempty"`
	CondCtx   map[string]any `json:"condition_context,omitempty"`
}

// ModelEvaluator evaluates Check against an in-memory tuple set. It is used by the tests
// and as the local-stack authorizer (decision 0004 #11); production uses OpenFGAClient.
type ModelEvaluator struct {
	mu     sync.RWMutex
	types  map[string]typeDef
	tuples []Tuple
	model  string
	down   bool
}

// LoadModel reads an OpenFGA JSON model file.
func LoadModel(path string) (*ModelEvaluator, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m AuthorizationModel
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if m.SchemaVersion != "1.1" {
		return nil, fmt.Errorf("unsupported schema_version %q", m.SchemaVersion)
	}
	e := &ModelEvaluator{types: map[string]typeDef{}, model: "local:" + path}
	for _, t := range m.TypeDefinitions {
		e.types[t.Type] = t
	}
	return e, nil
}

// SetDown makes every Check fail, simulating an unreachable OpenFGA.
func (e *ModelEvaluator) SetDown(down bool) { e.mu.Lock(); e.down = down; e.mu.Unlock() }

func (e *ModelEvaluator) StoreID() string { return "local" }
func (e *ModelEvaluator) ModelID() string { return e.model }

func splitObj(o string) (typ, id, rel string) {
	typ, rest, _ := strings.Cut(o, ":")
	id, rel, _ = strings.Cut(rest, "#")
	return
}

// Write adds tuples after enforcing the model's type restrictions, as OpenFGA does.
func (e *ModelEvaluator) Write(tuples ...Tuple) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if err := e.validate(tuples); err != nil {
		return err
	}
	e.tuples = append(e.tuples, tuples...)
	return nil
}

// validate enforces the model's type restrictions (OpenFGA applies them to written and to
// contextual tuples alike).
func (e *ModelEvaluator) validate(tuples []Tuple) error {
	for _, t := range tuples {
		otype, _, _ := splitObj(t.Object)
		td, ok := e.types[otype]
		if !ok {
			return fmt.Errorf("unknown object type %q", otype)
		}
		if _, ok := td.Relations[t.Relation]; !ok {
			return fmt.Errorf("type %s has no relation %s", otype, t.Relation)
		}
		utype, _, urel := splitObj(t.User)
		allowed := false
		if td.Metadata != nil {
			for _, r := range td.Metadata.Relations[t.Relation].Direct {
				if r.Type == utype && r.Relation == urel && r.Condition == t.Condition {
					allowed = true
				}
			}
		}
		if !allowed {
			return fmt.Errorf("type restriction: %s cannot be %s of %s", t.User, t.Relation, t.Object)
		}
	}
	return nil
}

func (e *ModelEvaluator) Check(_ context.Context, req CheckRequest) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.down {
		return false, errors.New("authorization store unreachable")
	}
	ts := e.tuples
	if len(req.ContextualTuples) > 0 {
		if err := e.validate(req.ContextualTuples); err != nil {
			return false, err
		}
		ts = append(append([]Tuple(nil), e.tuples...), req.ContextualTuples...)
	}
	return e.check(ts, req.User, req.Relation, req.Object, req.Context, 0)
}

func (e *ModelEvaluator) check(ts []Tuple, user, relation, object string, ctx map[string]any, depth int) (bool, error) {
	if depth > 25 {
		return false, errors.New("resolution depth exceeded")
	}
	otype, _, _ := splitObj(object)
	td, ok := e.types[otype]
	if !ok {
		return false, nil
	}
	rw, ok := td.Relations[relation]
	if !ok {
		return false, nil
	}
	return e.eval(ts, rw, user, relation, object, ctx, depth)
}

func (e *ModelEvaluator) eval(ts []Tuple, rw rewrite, user, relation, object string, ctx map[string]any, depth int) (bool, error) {
	switch {
	case rw.This != nil:
		for _, t := range ts {
			if t.Object != object || t.Relation != relation {
				continue
			}
			if t.Condition != "" && !evalCondition(t.Condition, t.CondCtx, ctx) {
				continue
			}
			if t.User == user {
				return true, nil
			}
			if ut, uid, urel := splitObj(t.User); urel != "" {
				ok, err := e.check(ts, user, urel, ut+":"+uid, ctx, depth+1)
				if err != nil || ok {
					return ok, err
				}
			}
		}
	case rw.ComputedUserset != nil:
		return e.check(ts, user, rw.ComputedUserset.Relation, object, ctx, depth+1)
	case rw.TupleToUserset != nil:
		for _, t := range ts {
			if t.Object == object && t.Relation == rw.TupleToUserset.Tupleset.Relation {
				ok, err := e.check(ts, user, rw.TupleToUserset.ComputedUserset.Relation, t.User, ctx, depth+1)
				if err != nil || ok {
					return ok, err
				}
			}
		}
	case rw.Union != nil:
		for _, c := range rw.Union.Child {
			ok, err := e.eval(ts, c, user, relation, object, ctx, depth+1)
			if err != nil || ok {
				return ok, err
			}
		}
	}
	return false, nil
}

// evalCondition implements the model's only condition, within_window
// (current_time < grant_time + grant_duration). Unknown conditions deny.
func evalCondition(name string, stored, req map[string]any) bool {
	if name != "within_window" {
		return false
	}
	get := func(k string) (string, bool) {
		if v, ok := stored[k].(string); ok { // persisted context wins (OpenFGA rule)
			return v, true
		}
		v, ok := req[k].(string)
		return v, ok
	}
	cs, ok1 := get("current_time")
	gs, ok2 := get("grant_time")
	ds, ok3 := get("grant_duration")
	if !ok1 || !ok2 || !ok3 {
		return false
	}
	ct, err1 := time.Parse(time.RFC3339Nano, cs)
	gt, err2 := time.Parse(time.RFC3339Nano, gs)
	d, err3 := time.ParseDuration(ds)
	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}
	return ct.Before(gt.Add(d))
}

// ---------------------------------------------------------------------------
// OpenFGA HTTP client (POST /stores/{store_id}/check, preshared-key auth as deployed in
// platform/ecosystem/values/openfga.yaml).

type OpenFGAClient struct {
	BaseURL string
	Store   string
	Model   string
	Token   string
	HTTP    *http.Client
}

func (c *OpenFGAClient) StoreID() string { return c.Store }
func (c *OpenFGAClient) ModelID() string { return c.Model }

func (c *OpenFGAClient) Check(ctx context.Context, req CheckRequest) (bool, error) {
	body := map[string]any{
		"tuple_key": map[string]string{"user": req.User, "relation": req.Relation, "object": req.Object},
	}
	if c.Model != "" {
		body["authorization_model_id"] = c.Model
	}
	if len(req.Context) > 0 {
		body["context"] = req.Context
	}
	if len(req.ContextualTuples) > 0 {
		keys := make([]map[string]string, 0, len(req.ContextualTuples))
		for _, t := range req.ContextualTuples {
			keys = append(keys, map[string]string{"user": t.User, "relation": t.Relation, "object": t.Object})
		}
		body["contextual_tuples"] = map[string]any{"tuple_keys": keys}
	}
	b, _ := json.Marshal(body)
	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(c.BaseURL, "/")+"/stores/"+c.Store+"/check", bytes.NewReader(b))
	if err != nil {
		return false, err
	}
	hreq.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		hreq.Header.Set("Authorization", "Bearer "+c.Token)
	}
	hc := c.HTTP
	if hc == nil {
		hc = &http.Client{Timeout: 2 * time.Second}
	}
	resp, err := hc.Do(hreq)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("openfga check: HTTP %d", resp.StatusCode)
	}
	var out struct {
		Allowed bool `json:"allowed"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return false, err
	}
	return out.Allowed, nil
}
