import {
  backupNeedsReview,
  backupNeedsSetup,
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
      backupsNeedSetup: 0,
      backupsRetired: 0,
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
      backupsNeedSetup: 'unknown',
      backupsRetired: 'unknown',
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

describe('backup states from the aggregator', () => {
  const backups: OperationsStatus['backups'] = [
    {
      name: 'backstage/backstage-backup',
      lastSuccess: '2026-10-08T07:15:27Z',
      suspended: true,
      state: 'retired',
      attention: 'none',
    },
    {
      name: 'clickhouse/backup-clickhouse',
      lastSuccess: '2026-10-08T03:00:13Z',
      suspended: true,
      state: 'retired',
      attention: 'none',
    },
    {
      name: 'ops/backup-offsite-copy',
      lastSuccess: null,
      suspended: true,
      state: 'needs-setup',
      attention: 'setup',
    },
    {
      name: 'backstage/backup-postgres-backstage',
      lastSuccess: '2026-10-21T10:00:00Z',
      suspended: false,
      state: 'active',
      attention: 'none',
    },
    {
      name: 'x/failing-backup',
      lastSuccess: '2026-10-19T00:00:00Z',
      suspended: false,
      state: 'active',
      attention: 'review',
    },
  ];

  it('does not review retired or needs-setup backups, even though they are suspended', () => {
    expect(backups.filter(b => backupNeedsReview(b, now)).map(b => b.name)).toEqual([
      'x/failing-backup',
    ]);
    expect(backups.filter(backupNeedsSetup).map(b => b.name)).toEqual([
      'ops/backup-offsite-copy',
    ]);
  });

  it('counts review, needs-setup and retired separately in the strip', () => {
    expect(summarizeStatus({ ...base, backups, serviceHealth: undefined }, now)).toMatchObject({
      backups: 1,
      backupsNeedSetup: 1,
      backupsRetired: 2,
    });
  });

  it('prefers the shared counts when the aggregator sends them', () => {
    expect(
      summarizeStatus(
        {
          ...base,
          counts: { alerts: 4, argo: 2, workloads: 1, backups: 0, backupsNeedSetup: 1, backupsRetired: 2 },
        },
        now,
      ),
    ).toEqual({ alerts: 4, argo: 2, workloads: 1, backups: 0, backupsNeedSetup: 1, backupsRetired: 2 });
  });
});
