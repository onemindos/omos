package engine

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func compileSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	p, _ := filepath.Abs("../../../schemas/entity/entity.v1.schema.json")
	s, err := jsonschema.NewCompiler().Compile(p)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func validate(s *jsonschema.Schema, v any) error {
	b, _ := json.Marshal(v)
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(b))
	if err != nil {
		return err
	}
	return s.Validate(doc)
}

// entity.md §10 test 4: every published record (sim-x1 included) validates against
// services/schemas/entity/entity.v1.schema.json.
func TestEveryPublishedRecordValidates(t *testing.T) {
	s := compileSchema(t)
	c := &clock{t0 + 60000}
	e := newEngine(c)
	_, all := replay(t, e)
	n := 0
	for _, r := range all {
		if r.Publish != nil {
			n++
			if err := validate(s, r.Publish.Entity); err != nil {
				t.Fatalf("%s seq %d: %v", r.Publish.Entity.EntityID, r.Publish.Entity.Seq, err)
			}
		}
	}
	c.ms = t0 + 10*86400000 // everything expires: final live:false records validate too
	for _, p := range e.Sweep() {
		n++
		if err := validate(s, p.Entity); err != nil || p.Entity.Live {
			t.Fatalf("final %s: live=%v %v", p.Entity.EntityID, p.Entity.Live, err)
		}
	}
	if n < 20 {
		t.Fatalf("only %d records validated", n)
	}
}

// Stop-line SL9 at the schema: a person with disposition hostile does not validate.
func TestSchemaRejectsHostilePerson(t *testing.T) {
	s := compileSchema(t)
	e := newEngine(&clock{t0 + 60000})
	last, _ := replay(t, e)
	var person map[string]any
	for _, p := range last {
		if p.Entity.EntityType == "person" {
			b, _ := json.Marshal(p.Entity)
			_ = json.Unmarshal(b, &person)
			break
		}
	}
	if person == nil {
		t.Fatal("no person in fixtures")
	}
	if err := validate(s, person); err != nil {
		t.Fatalf("baseline person invalid: %v", err)
	}
	person["disposition"] = "hostile"
	if validate(s, person) == nil {
		t.Fatal("schema accepted a hostile person")
	}
}
