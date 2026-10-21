import {
  backupNeedsReview,
  isStale,
  OperationsStatus,
  summarizeStatus,
} from './status';

const now = Date.parse('2026-10-21T12:00:00Z');
const base: OperationsStatus = {
  observedAt: '2026-10-21T11:59:00Z',
  partial: false,
  totalWorkloads: 40,
  unavailable: ['a/one', 'a/two'],
  deployments: [{ name: 'app', health: 'Degraded', sync: 'OutOfSync' }],
  backups: [
    {
      name: 'backstage/backstage-backup',
      lastSuccess: '2026-10-19T00:00:00Z',
      suspended: false,
    },
    {
      name: 'x/other-backup',
      lastSuccess: '2026-10-01T00:00:00Z',
      suspended: false,
    },
    { name: 'y/suspended-backup', lastSuccess: null, suspended: true },
  ],
  serviceHealth: {
    checks: [
      { name: 'a', state: 'up' },
      { name: 'b', state: 'down' },
      { name: 'c', state: 'unknown' },
    ],
  },
};

describe('backupNeedsReview', () => {
  it('flags suspended, never-succeeded and the overdue portal backup only', () => {
    expect(backupNeedsReview(base.backups[0], now)).toBe(true);
    expect(backupNeedsReview(base.backups[1], now)).toBe(false);
    expect(backupNeedsReview(base.backups[2], now)).toBe(true);
    expect(
      backupNeedsReview(
        { ...base.backups[0], lastSuccess: '2026-10-21T01:00:00Z' },
        now,
      ),
    ).toBe(false);
  });
});

describe('summarizeStatus', () => {
  it('counts the same items the Needs attention panel lists', () => {
    expect(summarizeStatus(base, now)).toEqual({
      alerts: 2,
      argo: 1,
      workloads: 2,
      backups: 2,
    });
  });
  it('reports unknown, never zero, for failed sources', () => {
    const partial = {
      ...base,
      partial: true,
      unavailableSources: ['Argo CD', 'Backup schedules', 'Deployments'],
      serviceHealth: undefined,
    };
    expect(summarizeStatus(partial, now)).toEqual({
      alerts: 'unknown',
      argo: 'unknown',
      workloads: 'unknown',
      backups: 'unknown',
    });
  });
});

describe('isStale', () => {
  it('is stale after three minutes', () => {
    expect(isStale(base, now)).toBe(false);
    expect(isStale({ ...base, observedAt: '2026-10-21T11:56:00Z' }, now)).toBe(
      true,
    );
  });
});
