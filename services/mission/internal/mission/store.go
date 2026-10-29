package mission

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/nats-io/nats.go/jetstream"
)

// errNoChange ends a mutation without a write.
var errNoChange = errors.New("no change")

// maxCASAttempts bounds the read-modify-write retries on a revision conflict.
const maxCASAttempts = 8

// store is KV mission-state with compare-and-set writes and an in-memory copy kept by a
// KV watch. The copy serves svc.mission.list, the timeout sweep and the arrival check;
// every write reads the authoritative revision from KV and writes with Update(rev), so a
// stale copy can never overwrite a newer record.
type store struct {
	kv jetstream.KeyValue

	mu       sync.RWMutex
	tasks    map[string]cachedTask
	missions map[string]cachedMission
	gotoIdx  map[string]map[string]bool // assignee "class.id" -> open goto task ids

	locks sync.Map // key -> *sync.Mutex (serializes writers of one key in this process)
}

type cachedTask struct {
	rev uint64
	t   Task
}

type cachedMission struct {
	rev uint64
	m   Mission
}

func newStore(kv jetstream.KeyValue) *store {
	return &store{kv: kv, tasks: map[string]cachedTask{}, missions: map[string]cachedMission{},
		gotoIdx: map[string]map[string]bool{}}
}

func (s *store) lock(key string) func() {
	v, _ := s.locks.LoadOrStore(key, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

// watch fills the copy and keeps it current; it returns once the initial values are in.
func (s *store) watch(ctx context.Context) (jetstream.KeyWatcher, error) {
	w, err := s.kv.WatchAll(ctx)
	if err != nil {
		return nil, err
	}
	for e := range w.Updates() {
		if e == nil {
			break
		}
		s.apply(e)
	}
	go func() {
		for e := range w.Updates() {
			if e != nil {
				s.apply(e)
			}
		}
	}()
	return w, nil
}

func (s *store) apply(e jetstream.KeyValueEntry) {
	key := e.Key()
	if e.Operation() != jetstream.KeyValuePut {
		s.mu.Lock()
		if c, ok := s.tasks[key]; ok {
			s.unindex(c.t)
		}
		delete(s.tasks, key)
		delete(s.missions, key)
		s.mu.Unlock()
		return
	}
	switch {
	case strings.HasPrefix(key, "tsk-"):
		var t Task
		if json.Unmarshal(e.Value(), &t) == nil {
			s.putTask(t, e.Revision())
		}
	case strings.HasPrefix(key, "msn-"):
		var m Mission
		if json.Unmarshal(e.Value(), &m) == nil {
			s.putMission(m, e.Revision())
		}
	}
}

// putTask records a task in the copy unless the copy already holds a newer revision.
func (s *store) putTask(t Task, rev uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.tasks[t.TaskID]; ok {
		if c.rev >= rev {
			return
		}
		s.unindex(c.t)
	}
	s.tasks[t.TaskID] = cachedTask{rev: rev, t: t}
	if !IsTerminal(t.Status.State) && TaskTypes[t.Type] == "goto" && AssetClasses[t.Assignee.EntityType] {
		k := t.Assignee.String()
		if s.gotoIdx[k] == nil {
			s.gotoIdx[k] = map[string]bool{}
		}
		s.gotoIdx[k][t.TaskID] = true
	}
}

func (s *store) unindex(t Task) {
	k := t.Assignee.String()
	if m := s.gotoIdx[k]; m != nil {
		delete(m, t.TaskID)
		if len(m) == 0 {
			delete(s.gotoIdx, k)
		}
	}
}

func (s *store) putMission(m Mission, rev uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.missions[m.MissionID]; ok && c.rev >= rev {
		return
	}
	s.missions[m.MissionID] = cachedMission{rev: rev, m: m}
}

// cachedTasks returns a snapshot of the copy's tasks.
func (s *store) cachedTasks() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Task, 0, len(s.tasks))
	for _, c := range s.tasks {
		out = append(out, c.t)
	}
	return out
}

func (s *store) cachedMissions() []Mission {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Mission, 0, len(s.missions))
	for _, c := range s.missions {
		out = append(out, c.m)
	}
	return out
}

func (s *store) cachedTask(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.tasks[id]
	return c.t, ok
}

func (s *store) openGotos(assignee string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.gotoIdx[assignee]))
	for id := range s.gotoIdx[assignee] {
		out = append(out, id)
	}
	return out
}

// get reads one record from KV. A missing key returns (nil, 0, nil).
func (s *store) get(ctx context.Context, key string) ([]byte, uint64, error) {
	e, err := s.kv.Get(ctx, key)
	if errors.Is(err, jetstream.ErrKeyNotFound) || errors.Is(err, jetstream.ErrKeyDeleted) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	return e.Value(), e.Revision(), nil
}

func (s *store) getTask(ctx context.Context, id string) (*Task, uint64, error) {
	b, rev, err := s.get(ctx, id)
	if err != nil || b == nil {
		return nil, rev, err
	}
	var t Task
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, 0, fmt.Errorf("record %s unreadable: %w", id, err)
	}
	return &t, rev, nil
}

func (s *store) getMission(ctx context.Context, id string) (*Mission, uint64, error) {
	b, rev, err := s.get(ctx, id)
	if err != nil || b == nil {
		return nil, rev, err
	}
	var m Mission
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, 0, fmt.Errorf("record %s unreadable: %w", id, err)
	}
	return &m, rev, nil
}

// createTask writes a new task; jetstream.ErrKeyExists if the key is taken.
func (s *store) createTask(ctx context.Context, t *Task) (uint64, error) {
	b, _ := json.Marshal(t)
	rev, err := s.kv.Create(ctx, t.TaskID, b)
	if err == nil {
		s.putTask(*t, rev)
	}
	return rev, err
}

func (s *store) createMission(ctx context.Context, m *Mission) (uint64, error) {
	b, _ := json.Marshal(m)
	rev, err := s.kv.Create(ctx, m.MissionID, b)
	if err == nil {
		s.putMission(*m, rev)
	}
	return rev, err
}

// mutateTask applies fn to the current record and writes it with compare-and-set on the
// KV revision, re-reading and re-applying on a conflict. fn returns errNoChange to skip
// the write. A missing record returns (nil, nil).
func (s *store) mutateTask(ctx context.Context, id string, fn func(t *Task) error) (*Task, error) {
	unlock := s.lock(id)
	defer unlock()
	for attempt := 0; attempt < maxCASAttempts; attempt++ {
		t, rev, err := s.getTask(ctx, id)
		if err != nil || t == nil {
			return nil, err
		}
		if err := fn(t); err != nil {
			return t, err
		}
		b, _ := json.Marshal(t)
		nrev, err := s.kv.Update(ctx, id, b, rev)
		if err == nil {
			s.putTask(*t, nrev)
			return t, nil
		}
		if !isRevisionConflict(err) {
			return nil, err
		}
	}
	return nil, fmt.Errorf("%s: compare-and-set failed %d times", id, maxCASAttempts)
}

func (s *store) mutateMission(ctx context.Context, id string, fn func(m *Mission) error) (*Mission, error) {
	unlock := s.lock(id)
	defer unlock()
	for attempt := 0; attempt < maxCASAttempts; attempt++ {
		m, rev, err := s.getMission(ctx, id)
		if err != nil || m == nil {
			return nil, err
		}
		if err := fn(m); err != nil {
			return m, err
		}
		b, _ := json.Marshal(m)
		nrev, err := s.kv.Update(ctx, id, b, rev)
		if err == nil {
			s.putMission(*m, nrev)
			return m, nil
		}
		if !isRevisionConflict(err) {
			return nil, err
		}
	}
	return nil, fmt.Errorf("%s: compare-and-set failed %d times", id, maxCASAttempts)
}

// isRevisionConflict reports a failed compare-and-set (JetStream "wrong last sequence").
func isRevisionConflict(err error) bool {
	var apiErr *jetstream.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode == jetstream.JSErrCodeStreamWrongLastSequence {
		return true
	}
	return errors.Is(err, jetstream.ErrKeyExists) || strings.Contains(err.Error(), "wrong last sequence")
}
