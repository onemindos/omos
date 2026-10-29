package layer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Config is the service configuration (environment in cmd/layer-service).
type Config struct {
	Layers         []Def
	Directory      *Directory
	EntStream      string        // ENT
	TelStream      string        // TEL (latest feed status by last-message get at start)
	WorldBucket    string        // world-state
	FleetBucket    string        // fleet-state
	Site           string        // hub or site-<name>
	LeaseBucket    string        // "" disables the lease (local runs only)
	LeaseKey       string        // _lease.layer-service
	LeaseTTL       time.Duration // 15s
	Holder         string
	QueueGroup     string // svc.layer.* queue group
	ReplayWindow   time.Duration
	ReplayMax      int
	MaxWatches     int
	WatchMemBudget int64         // bytes all watches may hold together (held sets, queued ops, cells)
	Heartbeat      time.Duration // empty batch when quiet (10s)
	CellEvery      time.Duration // cell-mode watch recompute (1s)
	ManifestTick   time.Duration // freshness evaluation (1s)
	ManifestEvery  time.Duration // manifest refresh (30s)
	CountEvery     time.Duration // rewrite after a count change, at most this often (5s)
	CellsEvery     time.Duration // KV cells refresh (30s)
	SiteEvery      time.Duration // fleet-state refresh (10s)
	SiteTimeout    time.Duration // a site with no contact for this long is not written (30s)
	SiteSubject    string        // tel.fleet.node.*.status
	SiteAllow      []string      // sites accepted from SiteSubject (empty: any valid token, up to MaxSites)
	MaxSites       int           // distinct sites tracked from SiteSubject (64)
	SiteStreams    []string      // hub streams whose sources give per-site lag ("" = all)
	SweepEvery     time.Duration // expiry sweep (1s)
	MaxCellRecord  int
	DisableSiteMon bool
}

// Defaults fills unset fields.
func (c *Config) Defaults() {
	def := func(d *time.Duration, v time.Duration) {
		if *d == 0 {
			*d = v
		}
	}
	if c.EntStream == "" {
		c.EntStream = "ENT"
	}
	if c.TelStream == "" {
		c.TelStream = "TEL"
	}
	if c.WorldBucket == "" {
		c.WorldBucket = "world-state"
	}
	if c.FleetBucket == "" {
		c.FleetBucket = "fleet-state"
	}
	if c.Site == "" {
		c.Site = "hub"
	}
	if c.LeaseKey == "" {
		c.LeaseKey = "_lease.layer-service"
	}
	if c.QueueGroup == "" {
		c.QueueGroup = "layer-svc"
	}
	if c.ReplayMax == 0 {
		c.ReplayMax = 200000
	}
	if c.MaxWatches == 0 {
		c.MaxWatches = DefaultMaxWatches
	}
	if c.WatchMemBudget == 0 {
		c.WatchMemBudget = DefaultWatchMemBudget
	}
	if c.MaxSites == 0 {
		c.MaxSites = 64
	}
	if c.SiteSubject == "" {
		c.SiteSubject = "tel.fleet.node.*.status"
	}
	if c.MaxCellRecord == 0 {
		c.MaxCellRecord = MaxCellRecord
	}
	if c.Directory == nil {
		c.Directory, _ = LoadDirectory("")
	}
	def(&c.LeaseTTL, 15*time.Second)
	def(&c.ReplayWindow, 5*time.Minute)
	def(&c.Heartbeat, 10*time.Second)
	def(&c.CellEvery, time.Second)
	def(&c.ManifestTick, time.Second)
	def(&c.ManifestEvery, 30*time.Second)
	def(&c.CountEvery, 5*time.Second)
	def(&c.CellsEvery, 30*time.Second)
	def(&c.SiteEvery, 10*time.Second)
	def(&c.SiteTimeout, 30*time.Second)
	def(&c.SweepEvery, time.Second)
}

// Watch limits sized for the 512Mi pod limit (platform/omos/layer/deployment.yaml): the
// entity store of a 50,000-vessel layer plus ENT read buffers stay under about 200 MiB, and
// all watches together hold at most DefaultWatchMemBudget; each entities-mode watch holds at
// most the layer cap. The count limit bounds goroutines and per-watch timers.
const (
	DefaultMaxWatches     = 200
	DefaultWatchMemBudget = 128 << 20
)

// Service is one layer service process.
type Service struct {
	nc    *nats.Conn
	js    jetstream.JetStream
	cfg   Config
	log   *slog.Logger
	store *Store
	mans  *Manifests
	now   func() time.Time

	ready     atomic.Bool
	snapshots atomic.Uint64
	batches   atomic.Uint64
	kvErrors  atomic.Uint64

	siteMu   sync.Mutex
	siteSeen map[string]time.Time

	watchMu  sync.Mutex // guards stopping and wg.Add for watch goroutines
	stopping bool
}

// New builds the service.
func New(nc *nats.Conn, cfg Config, log *slog.Logger) (*Service, error) {
	cfg.Defaults()
	if err := ValidateLayers(cfg.Layers); err != nil {
		return nil, err
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	s := &Service{nc: nc, js: js, cfg: cfg, log: log, now: time.Now, siteSeen: map[string]time.Time{}}
	s.store = NewStore(cfg.Layers, cfg.Site, cfg.ReplayWindow, cfg.ReplayMax, s.now)
	s.store.SetWatchBudget(cfg.WatchMemBudget)
	s.mans = NewManifests(cfg.Layers)
	return s, nil
}

// Store exposes the store (tests).
func (s *Service) Store() *Store { return s.store }

// Ready reports whether the service holds the lease and has caught up with ENT.
func (s *Service) Ready() bool { return s.ready.Load() }

// Metrics returns counters for /metrics.
func (s *Service) Metrics() map[string]uint64 {
	ready, down, seq := s.store.State()
	b := func(v bool) uint64 {
		if v {
			return 1
		}
		return 0
	}
	counts, _ := s.store.Counts()
	total := 0
	for _, c := range counts {
		total += c
	}
	return map[string]uint64{"ready": b(ready && s.ready.Load()), "ent_down": b(down), "ent_seq": seq,
		"entities": uint64(total), "watches": uint64(s.store.WatchCount()), "snapshots_total": s.snapshots.Load(),
		"batches_total": s.batches.Load(), "watch_bytes": uint64(max(0, s.store.WatchBytes())), "kv_errors_total": s.kvErrors.Load(), "ent_unreadable_total": s.store.Bad()}
}

// Run holds the lease, follows ENT, writes KV and serves svc.layer.* until ctx ends or the
// lease is lost.
func (s *Service) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	world, err := s.js.KeyValue(ctx, s.cfg.WorldBucket)
	if err != nil {
		return fmt.Errorf("bucket %q: %w (create it from its CR; the service does not create buckets)", s.cfg.WorldBucket, err)
	}
	var fleet jetstream.KeyValue
	if !s.cfg.DisableSiteMon {
		if fleet, err = s.js.KeyValue(ctx, s.cfg.FleetBucket); err != nil {
			return fmt.Errorf("bucket %q: %w (create it from its CR; the service does not create buckets)", s.cfg.FleetBucket, err)
		}
	}
	if s.cfg.LeaseBucket != "" {
		lease, err := NewLease(ctx, s.js, s.cfg.LeaseBucket, s.cfg.LeaseKey, s.cfg.Holder, s.cfg.LeaseTTL)
		if err != nil {
			return err
		}
		s.log.Info("waiting for lease", "bucket", s.cfg.LeaseBucket, "key", s.cfg.LeaseKey, "holder", s.cfg.Holder)
		if err := lease.Acquire(ctx); err != nil {
			return err
		}
		s.log.Info("lease acquired")
		defer lease.Release(context.Background())
		go lease.Keep(ctx, func(err error) { cancel(err) })
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	// Runs before wg.Wait (defers are LIFO): a watch request still in a NATS callback after
	// this point is refused instead of calling wg.Add during Wait.
	defer func() {
		s.watchMu.Lock()
		s.stopping = true
		s.watchMu.Unlock()
	}()
	subs, err := s.subscribe(ctx, &wg)
	for _, sub := range subs {
		defer func(sub *nats.Subscription) { _ = sub.Unsubscribe() }(sub)
	}
	if err != nil {
		return err
	}
	s.loadFeeds(ctx)

	readyCh := make(chan struct{})
	wg.Add(1)
	go func() { defer wg.Done(); s.follow(ctx, readyCh) }()
	select {
	case <-readyCh:
	case <-ctx.Done():
		return context.Cause(ctx)
	}
	s.ready.Store(true)
	defer s.ready.Store(false)
	_, _, seq := s.store.State()
	s.log.Info("caught up with ENT", "ent_seq", seq)

	wg.Add(2)
	go func() { defer wg.Done(); s.every(ctx, s.cfg.SweepEvery, func() { s.store.Sweep() }) }()
	go func() { defer wg.Done(); s.manifestLoop(ctx, world) }()
	if fleet != nil {
		wg.Add(1)
		go func() { defer wg.Done(); s.siteLoop(ctx, fleet) }()
	}
	<-ctx.Done()
	if c := context.Cause(ctx); c != nil && !errors.Is(c, context.Canceled) {
		return c
	}
	return nil
}

func (s *Service) every(ctx context.Context, d time.Duration, f func()) {
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			f()
		}
	}
}

// subscribe registers svc.layer.* and the feed and site status subscriptions.
func (s *Service) subscribe(ctx context.Context, wg *sync.WaitGroup) ([]*nats.Subscription, error) {
	var subs []*nats.Subscription
	add := func(subj string, h nats.MsgHandler, queue bool) error {
		var sub *nats.Subscription
		var err error
		if queue {
			sub, err = s.nc.QueueSubscribe(subj, s.cfg.QueueGroup, h)
		} else {
			sub, err = s.nc.Subscribe(subj, h)
		}
		if err != nil {
			return err
		}
		subs = append(subs, sub)
		return nil
	}
	if err := add("svc.layer.snapshot", s.handleSnapshot, true); err != nil {
		return subs, err
	}
	if err := add("svc.layer.watch", func(m *nats.Msg) { s.handleWatch(ctx, wg, m) }, true); err != nil {
		return subs, err
	}
	if err := add("svc.layer.unwatch", func(m *nats.Msg) {
		_ = m.Respond(s.store.Unwatch(m.Data, s.cfg.Directory.Resolve(m)))
	}, true); err != nil {
		return subs, err
	}
	if err := add("svc.layer.list", s.handleList, true); err != nil {
		return subs, err
	}
	if err := add("svc.layer.meta", s.handleMeta, true); err != nil {
		return subs, err
	}
	if err := add("tel.feed.*.status", func(m *nats.Msg) {
		if fs, ok := ParseFeedStatus(m.Subject, m.Data, s.now()); ok {
			s.mans.SetFeed(fs)
		}
	}, false); err != nil {
		return subs, err
	}
	if !s.cfg.DisableSiteMon {
		if err := add(s.cfg.SiteSubject, func(m *nats.Msg) {
			p := strings.Split(m.Subject, ".")
			if len(p) == 5 && siteRe.MatchString(p[3]) {
				s.seeSite(p[3])
			}
		}, false); err != nil {
			return subs, err
		}
	}
	return subs, s.nc.Flush()
}

var siteRe = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

// seeSite records contact from a site. Only allowlisted sites are accepted when SiteAllow is
// set; otherwise at most MaxSites distinct sites are tracked (stale ones are pruned first), so
// a publisher cannot grow the map or fill fleet-state with arbitrary site names.
func (s *Service) seeSite(site string) {
	if len(s.cfg.SiteAllow) > 0 {
		ok := false
		for _, a := range s.cfg.SiteAllow {
			if a == site {
				ok = true
				break
			}
		}
		if !ok {
			return
		}
	}
	now := s.now()
	s.siteMu.Lock()
	defer s.siteMu.Unlock()
	if _, known := s.siteSeen[site]; !known && len(s.siteSeen) >= s.cfg.MaxSites {
		s.pruneSitesLocked(now)
		if len(s.siteSeen) >= s.cfg.MaxSites {
			return
		}
	}
	s.siteSeen[site] = now
}

// pruneSitesLocked drops sites not heard from within SiteTimeout. siteMu is held.
func (s *Service) pruneSitesLocked(now time.Time) {
	for site, t := range s.siteSeen {
		if now.Sub(t) > s.cfg.SiteTimeout {
			delete(s.siteSeen, site)
		}
	}
}

// loadFeeds reads the latest status of each configured feed from TEL (by last message for
// subject; Direct Get when the stream allows it). Missing streams or messages are ignored:
// the core subscription fills the cache from then on.
func (s *Service) loadFeeds(ctx context.Context) {
	st, err := s.js.Stream(ctx, s.cfg.TelStream)
	if err != nil {
		return
	}
	seen := map[string]bool{}
	for _, d := range s.cfg.Layers {
		for _, f := range d.Feeds {
			if seen[f] {
				continue
			}
			seen[f] = true
			subj := "tel.feed." + f + ".status"
			m, err := st.GetLastMsgForSubject(ctx, subj)
			if err != nil {
				continue
			}
			if fs, ok := ParseFeedStatus(subj, m.Data, m.Time); ok {
				s.mans.SetFeed(fs)
			}
		}
	}
}

// follow reads ENT with an ordered consumer: last message per subject first (the snapshot),
// then every change. It closes readyCh when caught up.
func (s *Service) follow(ctx context.Context, readyCh chan struct{}) {
	readyOnce := sync.Once{}
	markReady := func() { readyOnce.Do(func() { s.store.MarkReady(); close(readyCh) }) }
	backoff := 500 * time.Millisecond
	for ctx.Err() == nil {
		err := s.followOnce(ctx, markReady)
		if ctx.Err() != nil {
			return
		}
		s.store.SetDown(true)
		s.log.Warn("ENT follow failed; retrying", "err", err, "backoff", backoff)
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, 10*time.Second)
	}
}

func (s *Service) followOnce(ctx context.Context, markReady func()) error {
	st, err := s.js.Stream(ctx, s.cfg.EntStream)
	if err != nil {
		return err
	}
	info, err := st.Info(ctx)
	if err != nil {
		return err
	}
	last := s.store.LastSeq()
	cfg := jetstream.OrderedConsumerConfig{FilterSubjects: []string{"ent.>"}}
	if last > 0 {
		cfg.DeliverPolicy, cfg.OptStartSeq = jetstream.DeliverByStartSequencePolicy, last+1
	} else {
		cfg.DeliverPolicy = jetstream.DeliverLastPerSubjectPolicy
		if info.State.Msgs == 0 {
			s.store.SetLastSeq(info.State.LastSeq)
			cfg.DeliverPolicy, cfg.OptStartSeq = jetstream.DeliverByStartSequencePolicy, info.State.LastSeq+1
			markReady()
		}
	}
	oc, err := s.js.OrderedConsumer(ctx, s.cfg.EntStream, cfg)
	if err != nil {
		return err
	}
	it, err := oc.Messages(jetstream.PullMaxMessages(1000))
	if err != nil {
		return err
	}
	defer it.Stop()
	stopped := make(chan struct{})
	defer close(stopped)
	go func() {
		select {
		case <-ctx.Done():
			it.Stop()
		case <-stopped:
		}
	}()
	if last > 0 {
		markReady()
	}
	s.store.SetDown(false)
	for {
		msg, err := it.Next()
		if err != nil {
			return err
		}
		meta, merr := msg.Metadata()
		if merr != nil {
			continue
		}
		s.store.Apply(msg.Subject(), msg.Data(), msg.Headers().Get("Nats-Marker-Reason"), meta.Sequence.Stream)
		if meta.NumPending == 0 {
			markReady()
		}
	}
}

// ---------------------------------------------------------------- endpoints

func (s *Service) handleSnapshot(m *nats.Msg) {
	if m.Reply == "" {
		return
	}
	caller := s.cfg.Directory.Resolve(m)
	res, errb := s.store.Snapshot(m.Data, caller, s.mans.Zooms, s.mans.Object)
	if errb != nil {
		_ = m.Respond(errb)
		return
	}
	s.snapshots.Add(1)
	for _, c := range res.Chunks {
		if err := s.nc.Publish(m.Reply, c); err != nil {
			s.log.Warn("snapshot chunk publish failed", "err", err)
			return
		}
	}
}

func (s *Service) handleWatch(ctx context.Context, wg *sync.WaitGroup, m *nats.Msg) {
	s.watchMu.Lock()
	if s.stopping || ctx.Err() != nil {
		s.watchMu.Unlock()
		if m.Reply != "" {
			_ = s.nc.Publish(m.Reply, errReply(ErrLayerDown, "layer service stopping"))
		}
		return
	}
	wg.Add(1) // under watchMu, so Run's wg.Wait cannot have started
	s.watchMu.Unlock()
	caller := s.cfg.Directory.Resolve(m)
	w, reply := s.store.Watch(m.Data, m.Reply, caller, WatchOpts{MaxWatches: s.cfg.MaxWatches, ZoomsFor: s.mans.Zooms})
	if reply != nil && m.Reply != "" {
		_ = s.nc.Publish(m.Reply, reply)
	}
	if w == nil {
		wg.Done()
		return
	}
	go func() { defer wg.Done(); s.runWatch(ctx, w) }()
}

// runWatch sends batches for one watch at most max_rate_hz times per second until its lease
// ends, it is unwatched, or the service stops.
func (s *Service) runWatch(ctx context.Context, w *Watch) {
	defer s.store.removeWatch(w)
	tick := min(time.Second, s.cfg.Heartbeat, s.cfg.CellEvery)
	var next time.Time
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.wake:
		case <-time.After(tick):
		}
		if d := time.Until(next); d > 0 { // at most max_rate_hz batches per second
			select {
			case <-ctx.Done():
				return
			case <-time.After(d):
			}
		}
		s.store.mu.RLock()
		w.mu.Lock()
		b, alive := w.buildBatchLocked(s.store, s.now(), s.cfg.Heartbeat, s.cfg.CellEvery)
		// A dirty cells view whose next diff is not due yet waits for the tick instead of
		// re-signalling itself (which would spin until CellEvery passes).
		more := len(w.order) > 0 || len(w.cellOps) > 0 || (w.cellsDirty && b != nil)
		interval := time.Second / time.Duration(w.rateHz)
		w.mu.Unlock()
		s.store.mu.RUnlock()
		if b != nil && !alive { // the service ended the watch: send the error, then stop
			_ = s.nc.Publish(b.reply, b.data)
			return
		}
		if !alive {
			return
		}
		if b != nil {
			next = time.Now().Add(interval)
			if err := s.nc.Publish(b.reply, b.data); err != nil {
				s.log.Warn("watch batch publish failed", "watch", w.id, "err", err)
			}
			s.batches.Add(1)
		}
		if more {
			w.signal()
		}
	}
}

type listItem struct {
	Layer  string `json:"layer"`
	Title  string `json:"title"`
	Tier   string `json:"tier"`
	Count  int    `json:"count"`
	Status string `json:"status"`
}

func (s *Service) handleList(m *nats.Msg) {
	if !s.cfg.Directory.Resolve(m).Known {
		_ = m.Respond(errReply(ErrForbidden, "identity not in the layer service directory"))
		return
	}
	counts, _ := s.store.Counts()
	out := struct {
		V      int        `json:"v"`
		Layers []listItem `json:"layers"`
	}{V: 1}
	for i, d := range s.cfg.Layers {
		it := listItem{Layer: d.ID, Title: d.Title, Tier: d.Tier, Count: counts[i], Status: "down"}
		if b := s.mans.JSON(i); b != nil {
			var man Manifest
			if json.Unmarshal(b, &man) == nil {
				it.Status = man.Freshness.Status
			}
		}
		out.Layers = append(out.Layers, it)
	}
	b, _ := json.Marshal(out)
	_ = m.Respond(b)
}

func (s *Service) handleMeta(m *nats.Msg) {
	if !s.cfg.Directory.Resolve(m).Known {
		_ = m.Respond(errReply(ErrForbidden, "identity not in the layer service directory"))
		return
	}
	var req struct {
		Layer string `json:"layer"`
	}
	_ = json.Unmarshal(m.Data, &req)
	li := s.store.layerIndex(req.Layer)
	if li < 0 {
		_ = m.Respond(errReply(ErrUnknownLayer, req.Layer))
		return
	}
	b := s.mans.JSON(li)
	if b == nil {
		_ = m.Respond(errReply(ErrLayerDown, "no manifest yet"))
		return
	}
	_ = m.Respond(b)
}

// ---------------------------------------------------------------- KV writers

// manifestLoop evaluates freshness every ManifestTick and writes `{layer}.meta` on change and
// every ManifestEvery, and `{layer}.cells.z{z}` every CellsEvery (single writer: the lease).
func (s *Service) manifestLoop(ctx context.Context, kv jetstream.KeyValue) {
	s.writeManifests(ctx, kv)
	s.every(ctx, s.cfg.ManifestTick, func() { s.writeManifests(ctx, kv) })
}

func (s *Service) writeManifests(ctx context.Context, kv jetstream.KeyValue) {
	now := s.now()
	counts, seq := s.store.Counts()
	_, down, _ := s.store.State()
	for i := range s.cfg.Layers {
		d := &s.cfg.Layers[i]
		st := s.mans.state[i] // read without lock: only this goroutine writes state
		if d.Cells != nil && (st.cellsAt.IsZero() || now.Sub(st.cellsAt) >= s.cfg.CellsEvery) {
			cells, cseq := s.store.Cells(i, d.Cells.Zooms)
			recs, kept := encodeCells(d.ID, cells, d.Cells.Zooms, cseq, now, s.cfg.MaxCellRecord)
			for _, z := range kept {
				if _, err := kv.Put(ctx, d.ID+".cells.z"+strconv.Itoa(z), recs[z]); err != nil {
					s.kvErrors.Add(1)
					s.log.Warn("cells write failed", "layer", d.ID, "z", z, "err", err)
				}
			}
			for _, z := range d.Cells.Zooms {
				if _, ok := recs[z]; !ok {
					s.log.Warn("cell zoom dropped: record over the size limit", "layer", d.ID, "z", z)
				}
			}
			s.mans.mu.Lock()
			s.mans.state[i].zooms = kept
			s.mans.state[i].cellsAt = now
			s.mans.mu.Unlock()
			st = s.mans.state[i]
		}
		man := s.mans.build(d, i, s.cfg.Site, counts[i], seq, down, now)
		sig := man.signature()
		changed := sig != st.signature
		due := st.written.IsZero() || now.Sub(st.written) >= s.cfg.ManifestEvery ||
			(counts[i] != st.count && now.Sub(st.written) >= s.cfg.CountEvery)
		if !changed && !due {
			continue
		}
		b, _ := json.Marshal(man)
		if _, err := kv.Put(ctx, d.ID+".meta", b); err != nil {
			s.kvErrors.Add(1)
			s.log.Warn("manifest write failed", "layer", d.ID, "err", err)
			continue
		}
		s.mans.mu.Lock()
		s.mans.state[i].written, s.mans.state[i].signature, s.mans.state[i].count, s.mans.state[i].json = now, sig, counts[i], b
		s.mans.mu.Unlock()
		if changed || counts[i] != st.count {
			s.store.SetMeta(i, b)
		}
	}
}

// siteLoop writes `fleet-state` `{site}` every SiteEvery for each site in contact.
func (s *Service) siteLoop(ctx context.Context, kv jetstream.KeyValue) {
	s.writeSites(ctx, kv)
	s.every(ctx, s.cfg.SiteEvery, func() { s.writeSites(ctx, kv) })
}

// writeSites: the own site is always connected. Other sites come from
// `tel.fleet.node.{site}.status` (BUS-DESIGN §3.1) and from the hub streams that source a
// site's streams through its JetStream domain (`$JS.{site}.API`, BUS-DESIGN §8–9): the source
// lag is 0 when nothing is pending, otherwise the seconds since the source was last active.
func (s *Service) writeSites(ctx context.Context, kv jetstream.KeyValue) {
	now := s.now()
	sites := map[string]*SiteStatus{s.cfg.Site: {V: 1, Site: s.cfg.Site, Connected: true, LastContact: now.UnixMilli(), SourceLagS: map[string]float64{}}}
	s.siteMu.Lock()
	s.pruneSitesLocked(now)
	for site, t := range s.siteSeen {
		sites[site] = &SiteStatus{V: 1, Site: site, Connected: true, LastContact: t.UnixMilli(), SourceLagS: map[string]float64{}}
	}
	s.siteMu.Unlock()
	for _, si := range s.sourceInfos(ctx) {
		if si.External == nil {
			continue
		}
		site, ok := domainOf(si.External.APIPrefix)
		if !ok || si.Active < 0 {
			continue
		}
		contact := now.Add(-si.Active)
		if now.Sub(contact) > s.cfg.SiteTimeout {
			continue
		}
		st := sites[site]
		if st == nil {
			st = &SiteStatus{V: 1, Site: site, Connected: true, LastContact: contact.UnixMilli(), SourceLagS: map[string]float64{}}
			sites[site] = st
		} else if contact.UnixMilli() > st.LastContact {
			st.LastContact = contact.UnixMilli()
		}
		lag := 0.0
		if si.Lag > 0 {
			lag = float64(int64(si.Active.Seconds()*10)) / 10
		}
		st.SourceLagS[si.Name] = lag
	}
	for site, st := range sites {
		b, _ := json.Marshal(st)
		if _, err := kv.Put(ctx, site, b); err != nil {
			s.kvErrors.Add(1)
			s.log.Warn("site status write failed", "site", site, "err", err)
		}
	}
}

func domainOf(apiPrefix string) (string, bool) {
	if len(apiPrefix) <= len("$JS..API") || !strings.HasPrefix(apiPrefix, "$JS.") || !strings.HasSuffix(apiPrefix, ".API") {
		return "", false
	}
	d := strings.TrimSuffix(strings.TrimPrefix(apiPrefix, "$JS."), ".API")
	return d, siteRe.MatchString(d)
}

func (s *Service) sourceInfos(ctx context.Context) []*jetstream.StreamSourceInfo {
	var out []*jetstream.StreamSourceInfo
	collect := func(info *jetstream.StreamInfo) {
		if info != nil {
			out = append(out, info.Sources...)
		}
	}
	if len(s.cfg.SiteStreams) > 0 {
		for _, name := range s.cfg.SiteStreams {
			if st, err := s.js.Stream(ctx, name); err == nil {
				collect(st.CachedInfo())
			}
		}
		return out
	}
	l := s.js.ListStreams(ctx)
	for info := range l.Info() {
		collect(info)
	}
	return out
}
