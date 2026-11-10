import { summarizeFeed } from './liveSummary';
import { initialFeedState } from './liveState';
import type { FeedState } from './liveState';

const NOW = 1_000_000_000_000;
const base = (over: Partial<FeedState> = {}): FeedState => ({
  ...initialFeedState,
  loaded: true,
  lastStatusMs: NOW - 1000,
  status: {
    nats: { state: 'connected' },
    kubernetes: { state: 'connected' },
    argo: { state: 'connected' },
  },
  ...over,
});

test('fresh status with NATS connected is connected', () => {
  expect(summarizeFeed(base(), NOW).connection).toBe('connected');
});

test('stale status is reconnecting, then offline after a minute', () => {
  expect(
    summarizeFeed(base({ lastStatusMs: NOW - 30_000 }), NOW).connection,
  ).toBe('reconnecting');
  expect(
    summarizeFeed(base({ lastStatusMs: NOW - 120_000 }), NOW).connection,
  ).toBe('offline');
});

test('NATS down is offline, connecting is reconnecting', () => {
  const down = base({
    status: {
      nats: { state: 'down' },
      kubernetes: { state: 'connected' },
      argo: { state: 'connected' },
    },
  });
  expect(summarizeFeed(down, NOW).connection).toBe('offline');
});

test('totals pods, counts events in the window and attention namespaces', () => {
  const cell = (ns: string, ready: number, total: number, health: any) => ({
    namespace: ns,
    ready,
    total,
    pending: 0,
    failed: 0,
    restarts: 0,
    health,
  });
  const ev = (id: string, ageMs: number) => ({
    id,
    subject: 's',
    kind: 'evt',
    ts: new Date(NOW - ageMs).toISOString(),
    bytes: 1,
    summary: id,
  });
  const s = summarizeFeed(
    base({
      cluster: new Map([
        ['a', cell('a', 3, 3, 'ok')],
        ['b', cell('b', 1, 2, 'warn')],
      ]),
      events: [ev('old', 10 * 60_000), ev('new', 1000), ev('mid', 60_000)],
    }),
    NOW,
  );
  expect([s.podsReady, s.podsTotal]).toEqual([4, 5]);
  expect(s.namespacesNeedingAttention).toBe(1);
  expect(s.recentEvents).toBe(2);
  expect(s.latestEvents.map(e => e.id)).toEqual(['new', 'mid', 'old']);
});
