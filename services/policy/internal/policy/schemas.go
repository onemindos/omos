package policy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Schemas holds the compiled action and task parameter schemas from services/schemas/
// (actions/{action_type}.v{n}.schema.json, tasks/{task_type}.schema.json) plus the omos
// command and event envelopes, so outgoing messages are validated before they leave.
type Schemas struct {
	actions map[string]*jsonschema.Schema // "asset.goto.v1"
	tasks   map[string]*jsonschema.Schema // "onemind.task.goto.v1"
	command *jsonschema.Schema
	event   *jsonschema.Schema
}

// LoadSchemas compiles every schema under dir (the services/schemas directory).
func LoadSchemas(dir string) (*Schemas, error) {
	s := &Schemas{actions: map[string]*jsonschema.Schema{}, tasks: map[string]*jsonschema.Schema{}}
	load := func(path string) (*jsonschema.Schema, error) {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(b))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		c := jsonschema.NewCompiler()
		c.AssertFormat()
		url := "file://" + filepath.ToSlash(path)
		if err := c.AddResource(url, doc); err != nil {
			return nil, err
		}
		return c.Compile(url)
	}
	for sub, into := range map[string]map[string]*jsonschema.Schema{"actions": s.actions, "tasks": s.tasks} {
		files, err := filepath.Glob(filepath.Join(dir, sub, "*.schema.json"))
		if err != nil {
			return nil, err
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("no schemas in %s/%s", dir, sub)
		}
		for _, f := range files {
			sch, err := load(f)
			if err != nil {
				return nil, err
			}
			into[strings.TrimSuffix(filepath.Base(f), ".schema.json")] = sch
		}
	}
	var err error
	if s.command, err = load(filepath.Join(dir, "command", "envelope.schema.json")); err != nil {
		return nil, err
	}
	if s.event, err = load(filepath.Join(dir, "envelope", "envelope.schema.json")); err != nil {
		return nil, err
	}
	return s, nil
}

func validateRaw(sch *jsonschema.Schema, raw []byte) error {
	v, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return err
	}
	return sch.Validate(v)
}

// ValidateAction validates params against {action_type}.v{version}.
func (s *Schemas) ValidateAction(actionType string, version int, params json.RawMessage) error {
	sch, ok := s.actions[fmt.Sprintf("%s.v%d", actionType, version)]
	if !ok {
		return fmt.Errorf("no schema for %s v%d", actionType, version)
	}
	return validateRaw(sch, params)
}

// HasTask reports whether a task type has a schema.
func (s *Schemas) HasTask(taskType string) bool { _, ok := s.tasks[taskType]; return ok }

// ValidateTask validates a task spec against its type schema.
func (s *Schemas) ValidateTask(taskType string, spec json.RawMessage) error {
	sch, ok := s.tasks[taskType]
	if !ok {
		return fmt.Errorf("no schema for task type %s", taskType)
	}
	return validateRaw(sch, spec)
}

// ValidateCommand validates an outgoing command envelope against the omos command schema.
func (s *Schemas) ValidateCommand(raw []byte) error { return validateRaw(s.command, raw) }

// ValidateEvent validates an outgoing event envelope against the omos event schema.
func (s *Schemas) ValidateEvent(raw []byte) error { return validateRaw(s.event, raw) }
