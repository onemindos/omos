import { planConditions, RECOVERY_MIN_AGE_MS } from './conditions';
import type { ConditionRow } from './store';

const MIN = 60_000;
const t0 = Date.parse('2026-10-21T10:00:00Z');
const row = (over: Partial<ConditionRow> = {}): ConditionRow => ({
  scope: 'argo:app',
  reason: 'out-of-sync',
  since: t0,
  notified: false,
  ...over,
});
const active = (reason: ConditionRow['reason'], scope = 'argo:app') => [
  { scope, reason, resource: 'app' },
];

describe('condition planner', () => {
  test('degraded notifies immediately, once', () => {
    const first = planConditions({ now: t0, rows: [], active: active('degraded'), healthy: [] });
    expect(first.raise).toHaveLength(1);
    expect(first.upsert[0]).toMatchObject({ reason: 'degraded', notified: true });
    const again = planConditions({
      now: t0 + MIN,
      rows: [first.upsert[0]],
      active: active('degraded'),
      healthy: [],
    });
    expect(again.raise).toHaveLength(0);
    expect(again.upsert).toHaveLength(0);
  });

  test('out-of-sync waits fifteen minutes before notifying', () => {
    const start = planConditions({ now: t0, rows: [], active: active('out-of-sync'), healthy: [] });
    expect(start.raise).toHaveLength(0);
    expect(start.upsert[0]).toMatchObject({ notified: false, since: t0 });
    const early = planConditions({ now: t0 + 14 * MIN, rows: start.upsert, active: active('out-of-sync'), healthy: [] });
    expect(early.raise).toHaveLength(0);
    const late = planConditions({ now: t0 + 15 * MIN, rows: start.upsert, active: active('out-of-sync'), healthy: [] });
    expect(late.raise).toHaveLength(1);
    expect(late.raise[0].row.notified).toBe(true);
  });

  test('short out-of-sync that heals before grace never notifies and leaves no row', () => {
    const rows = [row()];
    const plan = planConditions({ now: t0 + 5 * MIN, rows, active: [], healthy: ['argo:app'] });
    expect(plan.raise).toHaveLength(0);
    expect(plan.recover).toHaveLength(0);
    expect(plan.remove).toEqual(['argo:app']);
  });

  test('notified condition recovers once healthy long enough', () => {
    const rows = [row({ reason: 'degraded', notified: true })];
    const tooSoon = planConditions({ now: t0 + RECOVERY_MIN_AGE_MS - 1, rows, active: [], healthy: ['argo:app'] });
    expect(tooSoon.recover).toHaveLength(0);
    expect(tooSoon.remove).toHaveLength(0);
    const ok = planConditions({ now: t0 + RECOVERY_MIN_AGE_MS, rows, active: [], healthy: ['argo:app'] });
    expect(ok.recover).toHaveLength(1);
    expect(ok.remove).toEqual(['argo:app']);
  });

  test('a reason change keeps one incident and does not re-notify', () => {
    const rows = [row({ reason: 'out-of-sync', notified: true })];
    const plan = planConditions({ now: t0 + 20 * MIN, rows, active: active('degraded'), healthy: [] });
    expect(plan.raise).toHaveLength(0);
    expect(plan.upsert[0]).toMatchObject({ reason: 'degraded', notified: true, since: t0 });
  });

  test('objects neither bad nor healthy (progressing) are left untouched', () => {
    const plan = planConditions({ now: t0 + 30 * MIN, rows: [row({ notified: true })], active: [], healthy: [] });
    expect(plan).toEqual({ raise: [], recover: [], upsert: [], remove: [] });
  });

  test('suspended backups notify immediately', () => {
    const plan = planConditions({ now: t0, rows: [], active: active('suspended', 'backup:x'), healthy: [] });
    expect(plan.raise).toHaveLength(1);
  });
});
