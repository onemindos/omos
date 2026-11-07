import type {
  ClusterDelta,
  EventsDelta,
  LiveSnapshot,
  NamespaceCell,
  StatusMessage,
  Track,
  TracksDelta,
} from './contract';
import {
  MAX_EVENTS,
  ORIGIN_SILENCE_MS,
  STALE_AFTER_MS,
  feedReducer,
  freshness,
  initialFeedState,
  kindSeries,
  sortArgo,
  sortCluster,
} from './liveState';
import type { FeedAction, FeedState } from './liveState';

const status = {
  nats: { state: 'connected' as const, rate: 12 },
  kubernetes: { state: 'connected' as const },
  argo: { state: 'disabled' as const },
};
const track = (uid: string, extra: Partial<Track> = {}): Track => ({
  uid,
  type: 'a-f-G-U-C',
  lat: 1,
  lon: 2,
  time: '2026-11-07T10:00:00Z',
  stale: '2026-11-07T10:05:00Z',
  subject: 'tak.cot.x',
  ...extra,
});
const cell = (
  namespace: string,
  health: NamespaceCell['health'] = 'ok',
): NamespaceCell => ({
  namespace,
  total: 3,
  ready: 3,
  pending: 0,
  failed: 0,
  restarts: 0,
  health,
});
const snapshot = (
  origin = 'A',
  seq = 10,
  extra: Partial<LiveSnapshot> = {},
): LiveSnapshot => ({
  origin,
  seq,
  generatedAt: '2026-11-07T10:00:00Z',
  status,
  tracks: [track('t1'), track('t2')],
  events: [
    {
      id: 'A:9',
      subject: 'tel.a.b.c',
      kind: 'tel',
      ts: '2026-11-07T10:00:00Z',
      bytes: 1,
      summary: 's',
    },
  ],
  cluster: [cell('nats'), cell('matrix', 'warn')],
  argo: [{ name: 'platform-nats', sync: 'Synced', health: 'Healthy' }],
  warnings: [],
  ...extra,
});
const tracksMsg = (
  origin: string,
  seq: number,
  upserts: Track[] = [],
  removes: string[] = [],
): TracksDelta => ({
  origin,
  seq,
  ts: '2026-11-07T10:00:01Z',
  upserts,
  removes,
});
const run = (actions: FeedAction[], from: FeedState = initialFeedState) =>
  actions.reduce(feedReducer, from);
const T0 = 1_000_000;

describe('snapshot', () => {
  it('loads every section and locks the snapshot origin', () => {
    const s = run([{ type: 'snapshot', snapshot: snapshot(), now: T0 }]);
    expect(s.loaded).toBe(true);
    expect(s.origin).toBe('A');
    expect([...s.tracks.keys()]).toEqual(['t1', 't2']);
    expect(s.cluster.get('matrix')?.health).toBe('warn');
    expect(s.argo.get('platform-nats')?.sync).toBe('Synced');
    expect(s.events).toHaveLength(1);
  });

  it('does not re-apply deltas the snapshot already contains', () => {
    const s = run([
      { type: 'snapshot', snapshot: snapshot('A', 10), now: T0 },
      {
        type: 'tracks',
        message: tracksMsg('A', 10, [track('t3')]),
        now: T0 + 1,
      },
      {
        type: 'tracks',
        message: tracksMsg('A', 11, [track('t4')]),
        now: T0 + 2,
      },
    ]);
    expect(s.tracks.has('t3')).toBe(false);
    expect(s.tracks.has('t4')).toBe(true);
  });
});

describe('seq handling', () => {
  const base = run([{ type: 'snapshot', snapshot: snapshot('A', 1), now: T0 }]);

  it('applies upserts and removes in order', () => {
    const s = run(
      [
        {
          type: 'tracks',
          message: tracksMsg('A', 2, [track('t3'), track('t1', { lat: 50 })]),
          now: T0 + 1,
        },
        { type: 'tracks', message: tracksMsg('A', 3, [], ['t2']), now: T0 + 2 },
      ],
      base,
    );
    expect([...s.tracks.keys()].sort()).toEqual(['t1', 't3']);
    expect(s.tracks.get('t1')?.lat).toBe(50);
  });

  it('ignores a repeated or older seq from the same origin', () => {
    const once = run(
      [
        {
          type: 'tracks',
          message: tracksMsg('A', 5, [track('t3')]),
          now: T0 + 1,
        },
      ],
      base,
    );
    const again = run(
      [
        {
          type: 'tracks',
          message: tracksMsg('A', 5, [track('t9')]),
          now: T0 + 2,
        },
        {
          type: 'tracks',
          message: tracksMsg('A', 4, [track('t8')]),
          now: T0 + 3,
        },
      ],
      once,
    );
    expect(again.tracks.has('t9')).toBe(false);
    expect(again.tracks.has('t8')).toBe(false);
    expect(again.lastSeq['A|tracks']).toBe(5);
  });

  it('tracks seq per channel, so a late message on another channel is not lost', () => {
    const s = run(
      [
        {
          type: 'tracks',
          message: tracksMsg('A', 7, [track('t3')]),
          now: T0 + 1,
        },
      ],
      base,
    );
    const events: EventsDelta = {
      origin: 'A',
      seq: 6,
      ts: 'x',
      events: [],
      dropped: 4,
    };
    // The events channel has not seen seq 6 yet, even though the tracks channel is at 7.
    expect(
      run([{ type: 'events', message: events, now: T0 + 2 }], s).droppedEvents,
    ).toBe(4);
    // The snapshot floor (seq 1) still applies to every channel.
    expect(
      run([{ type: 'events', message: { ...events, seq: 1 }, now: T0 + 2 }], s)
        .droppedEvents,
    ).toBe(0);
  });
});

describe('origin lock', () => {
  const locked = run([
    { type: 'snapshot', snapshot: snapshot('A', 1), now: T0 },
  ]);

  it('ignores another origin while the locked one is alive', () => {
    const s = run(
      [
        {
          type: 'tracks',
          message: tracksMsg('B', 1, [track('from-b')]),
          now: T0 + 5000,
        },
      ],
      locked,
    );
    expect(s.origin).toBe('A');
    expect(s.tracks.has('from-b')).toBe(false);
    expect(s.resyncToken).toBe(0);
  });

  it('keeps the lock when the locked origin keeps talking', () => {
    const s = run(
      [
        { type: 'tracks', message: tracksMsg('A', 2), now: T0 + 10_000 },
        {
          type: 'tracks',
          message: tracksMsg('B', 1, [track('from-b')]),
          now: T0 + 20_000,
        },
      ],
      locked,
    );
    expect(s.origin).toBe('A');
    expect(s.tracks.has('from-b')).toBe(false);
  });

  it('switches after 15 s of silence from the locked origin and asks for a resync', () => {
    const s = run(
      [
        {
          type: 'tracks',
          message: tracksMsg('B', 1, [track('from-b')]),
          now: T0 + ORIGIN_SILENCE_MS + 1,
        },
      ],
      locked,
    );
    expect(s.origin).toBe('B');
    expect(s.resyncToken).toBe(1);
    expect(s.tracks.has('from-b')).toBe(true);
  });

  it('does not switch at exactly 15 s', () => {
    const s = run(
      [
        {
          type: 'tracks',
          message: tracksMsg('B', 1, [track('from-b')]),
          now: T0 + ORIGIN_SILENCE_MS,
        },
      ],
      locked,
    );
    expect(s.origin).toBe('A');
  });

  it('a duplicate from the locked origin still counts as hearing it', () => {
    const s = run(
      [
        { type: 'tracks', message: tracksMsg('A', 1), now: T0 + 14_000 },
        {
          type: 'tracks',
          message: tracksMsg('B', 1, [track('from-b')]),
          now: T0 + 20_000,
        },
      ],
      locked,
    );
    expect(s.origin).toBe('A');
  });

  it('applies the new origin from its own seq counter after a switch', () => {
    const s = run(
      [
        {
          type: 'tracks',
          message: tracksMsg('B', 40, [track('b40')]),
          now: T0 + 20_000,
        },
        {
          type: 'tracks',
          message: tracksMsg('B', 41, [track('b41')]),
          now: T0 + 21_000,
        },
      ],
      locked,
    );
    expect(s.tracks.has('b40') && s.tracks.has('b41')).toBe(true);
    expect(s.resyncToken).toBe(1);
  });

  it('holds deltas that arrive before the snapshot and replays the ones newer than it', () => {
    const held = run([
      { type: 'tracks', message: tracksMsg('Z', 9, [track('older')]), now: T0 },
      {
        type: 'tracks',
        message: tracksMsg('Z', 12, [track('newer')]),
        now: T0 + 1,
      },
      {
        type: 'tracks',
        message: tracksMsg('Y', 99, [track('other-origin')]),
        now: T0 + 2,
      },
    ]);
    expect(held.loaded).toBe(false);
    expect(held.tracks.size).toBe(0);
    const s = run(
      [
        {
          type: 'snapshot',
          snapshot: snapshot('Z', 10, { tracks: [track('base')] }),
          now: T0 + 3,
        },
      ],
      held,
    );
    expect([...s.tracks.keys()].sort()).toEqual(['base', 'newer']);
    expect(s.origin).toBe('Z');
  });

  it('replays deltas that arrived while a later snapshot was in flight', () => {
    const s = run(
      [
        {
          type: 'tracks',
          message: tracksMsg('A', 2, [track('x')]),
          now: T0 + 1000,
        },
        {
          type: 'tracks',
          message: tracksMsg('A', 3, [track('y')]),
          now: T0 + 2000,
        },
        // The snapshot was generated at seq 2, so seq 3 must survive it.
        {
          type: 'snapshot',
          snapshot: snapshot('A', 2, { tracks: [track('x')] }),
          now: T0 + 3000,
        },
      ],
      locked,
    );
    expect([...s.tracks.keys()].sort()).toEqual(['x', 'y']);
  });
});

describe('events', () => {
  it('puts newest first, de-duplicates, counts dropped and buckets by kind', () => {
    const e = (n: number, kind = 'tel') => ({
      id: `A:${n}`,
      subject: `${kind}.x`,
      kind,
      ts: '2026-11-07T10:00:00Z',
      bytes: 1,
      summary: '',
    });
    const s = run([
      { type: 'snapshot', snapshot: snapshot('A', 1, { events: [] }), now: T0 },
      {
        type: 'events',
        message: {
          origin: 'A',
          seq: 2,
          ts: 'x',
          events: [e(1), e(2, 'cmd')],
          dropped: 3,
        },
        now: T0 + 1000,
      },
      {
        type: 'events',
        message: {
          origin: 'A',
          seq: 3,
          ts: 'x',
          events: [e(2, 'cmd'), e(3)],
          dropped: 2,
        },
        now: T0 + 1500,
      },
    ]);
    expect(s.events.map(x => x.id)).toEqual(['A:3', 'A:2', 'A:1']);
    expect(s.droppedEvents).toBe(5);
    const series = kindSeries(s.buckets, T0 + 1500);
    expect(series.tel.reduce((a, b) => a + b, 0)).toBe(2);
    expect(series.cmd.reduce((a, b) => a + b, 0)).toBe(1);
    expect(series.tel).toHaveLength(30);
  });

  it('caps the list', () => {
    const many = Array.from({ length: MAX_EVENTS + 40 }, (_, i) => ({
      id: `A:${i}`,
      subject: 's',
      kind: 'tel',
      ts: 'x',
      bytes: 1,
      summary: '',
    }));
    const s = run([
      { type: 'snapshot', snapshot: snapshot('A', 1, { events: [] }), now: T0 },
      {
        type: 'events',
        message: { origin: 'A', seq: 2, ts: 'x', events: many, dropped: 0 },
        now: T0,
      },
    ]);
    expect(s.events).toHaveLength(MAX_EVENTS);
  });
});

describe('cluster delta', () => {
  it('upserts changed cells, removes namespaces, merges argo and warnings, and marks the change time', () => {
    const delta: ClusterDelta = {
      origin: 'A',
      seq: 2,
      ts: 'x',
      namespaces: [cell('nats', 'down')],
      removedNamespaces: ['matrix'],
      argo: [{ name: 'platform-nats', sync: 'OutOfSync', health: 'Healthy' }],
      warnings: [
        {
          namespace: 'nats',
          object: 'pod/x',
          reason: 'BackOff',
          message: 'm',
          ts: '2026-11-07T10:00:02Z',
        },
      ],
    };
    const s = run([
      { type: 'snapshot', snapshot: snapshot('A', 1), now: T0 },
      { type: 'cluster', message: delta, now: T0 + 1000 },
      { type: 'cluster', message: { ...delta, seq: 3 }, now: T0 + 2000 },
    ]);
    expect(s.cluster.get('nats')?.health).toBe('down');
    expect(s.cluster.has('matrix')).toBe(false);
    expect(s.argo.get('platform-nats')?.sync).toBe('OutOfSync');
    expect(s.warnings).toHaveLength(1);
    expect(s.clusterChangedMs.get('nats')).toBe(T0 + 2000);
  });
});

describe('freshness', () => {
  const s0 = run([{ type: 'snapshot', snapshot: snapshot('A', 1), now: T0 }]);
  const msg: StatusMessage = { origin: 'A', seq: 2, ts: 'x', status };

  it('is waiting before anything arrives', () => {
    expect(freshness(initialFeedState, T0)).toBe('waiting');
  });

  it('is live within 15 s of the last status and stale after', () => {
    const s = run([{ type: 'status', message: msg, now: T0 + 1000 }], s0);
    expect(freshness(s, T0 + 1000 + STALE_AFTER_MS)).toBe('live');
    expect(freshness(s, T0 + 1000 + STALE_AFTER_MS + 1)).toBe('stale');
  });

  it('a later snapshot does not make a silent signal channel look live', () => {
    const s = run(
      [{ type: 'snapshot', snapshot: snapshot('A', 5), now: T0 + 5 * 60_000 }],
      s0,
    );
    expect(freshness(s, T0 + 5 * 60_000)).toBe('stale');
  });
});

describe('ordering', () => {
  it('lists Argo apps that are not Synced and Healthy first', () => {
    const sorted = sortArgo([
      { name: 'a-good', sync: 'Synced', health: 'Healthy' },
      { name: 'b-out', sync: 'OutOfSync', health: 'Healthy' },
      { name: 'c-sick', sync: 'Synced', health: 'Degraded' },
    ]);
    expect(sorted.map(a => a.name)).toEqual(['b-out', 'c-sick', 'a-good']);
  });

  it('lists unhealthy namespaces first', () => {
    expect(
      sortCluster([cell('a'), cell('b', 'down'), cell('c', 'warn')]).map(
        c => c.namespace,
      ),
    ).toEqual(['b', 'c', 'a']);
  });
});
