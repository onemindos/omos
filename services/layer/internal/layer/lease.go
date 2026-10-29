package layer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

// Lease is the single-writer lock (copied from services/entity/internal/bus/lease.go; same
// design as the entity service, entity.md §2; fusion-and-tracking.md "JetStream
// sketch"): one KV key held by compare-and-set. The value carries its own expiry, so the lease
// works on a bucket with or without a bucket TTL: a contender takes over an expired lease with
// Update(key, value, lastRevision), and two contenders cannot both win.
type Lease struct {
	kv     jetstream.KeyValue
	key    string
	holder string
	ttl    time.Duration
	now    func() time.Time

	mu  sync.Mutex // rev is written by Keep and read by Release
	rev uint64
}

func (l *Lease) setRev(r uint64) { l.mu.Lock(); l.rev = r; l.mu.Unlock() }

func (l *Lease) getRev() uint64 { l.mu.Lock(); defer l.mu.Unlock(); return l.rev }

type leaseValue struct {
	Holder    string `json:"holder"`
	ExpiresMs int64  `json:"expires_ms"`
}

// NewLease opens the lease. The bucket must already exist: the service never creates buckets
// (they are CRs in platform/nats/streams/).
func NewLease(ctx context.Context, js jetstream.JetStream, bucket, key, holder string, ttl time.Duration) (*Lease, error) {
	kv, err := js.KeyValue(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("lease bucket %q: %w (create it from its CR; the service does not create buckets)", bucket, err)
	}
	return &Lease{kv: kv, key: key, holder: holder, ttl: ttl, now: time.Now}, nil
}

func (l *Lease) value() []byte {
	b, _ := json.Marshal(leaseValue{Holder: l.holder, ExpiresMs: l.now().Add(l.ttl).UnixMilli()})
	return b
}

// TryAcquire makes one attempt. It returns true when this holder owns the lease.
func (l *Lease) TryAcquire(ctx context.Context) (bool, error) {
	rev, err := l.kv.Create(ctx, l.key, l.value())
	if err == nil {
		l.setRev(rev)
		return true, nil
	}
	if !errors.Is(err, jetstream.ErrKeyExists) {
		return false, err
	}
	e, err := l.kv.Get(ctx, l.key)
	if errors.Is(err, jetstream.ErrKeyNotFound) {
		return false, nil // deleted between the two calls; retry
	}
	if err != nil {
		return false, err
	}
	var cur leaseValue
	_ = json.Unmarshal(e.Value(), &cur)
	if cur.Holder != l.holder && cur.ExpiresMs > l.now().UnixMilli() {
		return false, nil
	}
	rev, err = l.kv.Update(ctx, l.key, l.value(), e.Revision())
	if err != nil {
		return false, nil // another contender won the CAS
	}
	l.setRev(rev)
	return true, nil
}

// Acquire blocks until the lease is held or ctx ends.
func (l *Lease) Acquire(ctx context.Context) error {
	for {
		ok, err := l.TryAcquire(ctx)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(l.ttl / 3):
		}
	}
}

// Renew extends the lease. An error means the lease is lost and the caller must stop.
func (l *Lease) Renew(ctx context.Context) error {
	rev, err := l.kv.Update(ctx, l.key, l.value(), l.getRev())
	if err != nil {
		return fmt.Errorf("lease lost: %w", err)
	}
	l.setRev(rev)
	return nil
}

// Keep renews every ttl/3 until ctx ends; it calls lost once and returns if renewal fails.
func (l *Lease) Keep(ctx context.Context, lost func(error)) {
	t := time.NewTicker(l.ttl / 3)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := l.Renew(ctx); err != nil {
				lost(err)
				return
			}
		}
	}
}

// Release deletes the key if this holder still owns it.
func (l *Lease) Release(ctx context.Context) {
	_ = l.kv.Delete(ctx, l.key, jetstream.LastRevision(l.getRev()))
}
