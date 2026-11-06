import { healthScore, scoreBand } from './health';
import { OperationsStatus } from './status';

const now = Date.parse('2026-11-06T12:00:00Z');
const entry = (key: string, severity: 'critical' | 'warning') => ({
  key,
  kind: 'service',
  severity,
  title: `Item ${key}`,
  detail: '',
  url: '/live',
});
const status = (extra: Partial<OperationsStatus> = {}): OperationsStatus => ({
  observedAt: '2026-11-06T11:59:30Z',
  partial: false,
  totalWorkloads: 10,
  unavailable: [],
  deployments: [],
  backups: [],
  attention: [],
  ...extra,
});

test('missing data gives no score, never a number', () => {
  expect(healthScore(undefined, false, now)).toBeUndefined();
  expect(healthScore(undefined, true, now)).toBeUndefined();
});

test('no attention items scores 100 and has no top item', () => {
  expect(healthScore(status(), false, now)).toMatchObject({
    score: 100,
    stale: false,
    partial: false,
    attentionCount: 0,
    top: undefined,
  });
});

test('warnings cost 8 and criticals cost 20', () => {
  const warn = healthScore(
    status({ attention: [entry('a', 'warning'), entry('b', 'warning')] }),
    false,
    now,
  );
  expect(warn?.score).toBe(84);
  const mixed = healthScore(
    status({ attention: [entry('a', 'warning'), entry('b', 'critical')] }),
    false,
    now,
  );
  expect(mixed?.score).toBe(72);
  expect(mixed?.top?.key).toBe('b');
});

test('score is clamped to 0', () => {
  const many = Array.from({ length: 9 }, (_, i) => entry(`c${i}`, 'critical'));
  expect(healthScore(status({ attention: many }), false, now)?.score).toBe(0);
});

test('stale and partial are flagged, not hidden', () => {
  expect(healthScore(status(), true, now)?.stale).toBe(true);
  expect(
    healthScore(status({ observedAt: '2026-11-06T11:00:00Z' }), false, now)
      ?.stale,
  ).toBe(true);
  expect(healthScore(status({ partial: true }), false, now)?.partial).toBe(
    true,
  );
});

test('bands', () => {
  expect(scoreBand(100)).toBe('healthy');
  expect(scoreBand(90)).toBe('healthy');
  expect(scoreBand(89)).toBe('degraded');
  expect(scoreBand(69)).toBe('at-risk');
});
