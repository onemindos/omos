import {
  classifyBackup,
  CronJobLike,
  isBackupJob,
  runsDaily,
  summarizeBackups,
} from './backupState';

const NOW = Date.parse('2026-10-22T12:00:00Z');
const hoursAgo = (h: number) => new Date(NOW - h * 3600000).toISOString();

const job = (
  name: string,
  extra: Partial<CronJobLike> & { annotations?: Record<string, string> } = {},
): CronJobLike => ({
  metadata: { name, namespace: 'ns', annotations: extra.annotations },
  spec: { schedule: '0 3 * * *', ...extra.spec },
  status: extra.status,
});

describe('classifyBackup: the three jobs Zeus asked about', () => {
  it('treats a superseded suspended job as retired, not an attention item', () => {
    const item = classifyBackup(
      {
        metadata: {
          name: 'backstage-backup',
          namespace: 'backstage',
          annotations: {
            'onemind.io/backup-state': 'superseded',
            'onemind.io/superseded-by': 'backup-postgres-backstage',
          },
        },
        spec: { suspend: true, schedule: '15 7 * * *' },
        status: { lastSuccessfulTime: '2026-10-08T07:15:27Z' },
      },
      NOW,
    );
    expect(item).toMatchObject({
      name: 'backstage/backstage-backup',
      state: 'retired',
      attention: 'none',
      suspended: true,
      supersededBy: 'backup-postgres-backstage',
    });
  });

  it('treats a pending-setup job as its own needs-setup item, with the reason', () => {
    const item = classifyBackup(
      {
        metadata: {
          name: 'backup-offsite-copy',
          namespace: 'ops',
          annotations: {
            'onemind.io/backup-state': 'pending-setup',
            'onemind.io/pending-reason': 'Waiting for offsite credentials',
          },
        },
        spec: { suspend: true, schedule: '0 17 * * *' },
      },
      NOW,
    );
    expect(item).toMatchObject({
      state: 'needs-setup',
      attention: 'setup',
      lastSuccess: null,
      pendingReason: 'Waiting for offsite credentials',
    });
  });

  it('an unannotated suspended job is still an attention item', () => {
    expect(classifyBackup(job('x-backup', { spec: { suspend: true } }), NOW)).toMatchObject({
      state: 'active',
      health: 'suspended',
      attention: 'review',
    });
  });
});

describe('classifyBackup: active health', () => {
  it('is ok after a recent success', () => {
    const item = classifyBackup(
      job('a-backup', { status: { lastSuccessfulTime: hoursAgo(6), lastScheduleTime: hoursAgo(6) } }),
      NOW,
    );
    expect(item).toMatchObject({ health: 'ok', attention: 'none' });
  });
  it('flags a job that never succeeded', () => {
    expect(classifyBackup(job('a-backup'), NOW)).toMatchObject({ health: 'never-ran', attention: 'review' });
  });
  it('flags a daily job overdue past 36 hours', () => {
    expect(
      classifyBackup(job('a-backup', { status: { lastSuccessfulTime: hoursAgo(40) } }), NOW),
    ).toMatchObject({ health: 'overdue', attention: 'review' });
  });
  it('does not assume a deadline for a weekly job', () => {
    expect(
      classifyBackup(
        job('a-backup', { spec: { schedule: '0 3 * * 0' }, status: { lastSuccessfulTime: hoursAgo(100) } }),
        NOW,
      ),
    ).toMatchObject({ health: 'ok' });
  });
  it('honours the max-age annotation', () => {
    expect(
      classifyBackup(
        job('a-backup', {
          annotations: { 'onemind.io/backup-max-age-hours': '200' },
          spec: { schedule: '0 3 * * 0' },
          status: { lastSuccessfulTime: hoursAgo(100) },
        }),
        NOW,
      ),
    ).toMatchObject({ health: 'ok' });
  });
  it('flags a last run that did not succeed when nothing is running', () => {
    const status = { lastSuccessfulTime: hoursAgo(30), lastScheduleTime: hoursAgo(6) };
    expect(classifyBackup(job('a-backup', { status }), NOW)).toMatchObject({ health: 'failing', attention: 'review' });
    expect(
      classifyBackup(job('a-backup', { status: { ...status, active: [{}] } }), NOW),
    ).toMatchObject({ health: 'ok' });
  });
  it('treats a mistyped state as active so the mistake is visible', () => {
    expect(
      classifyBackup(
        job('a-backup', { annotations: { 'onemind.io/backup-state': 'supersded' }, spec: { suspend: true } }),
        NOW,
      ),
    ).toMatchObject({ state: 'active', attention: 'review' });
  });
});

describe('helpers', () => {
  it('runsDaily needs wildcards in the last three fields', () => {
    expect(runsDaily('0 3 * * *')).toBe(true);
    expect(runsDaily('*/15 * * * *')).toBe(true);
    expect(runsDaily('0 3 * * 0')).toBe(false);
    expect(runsDaily('0 3 1 * *')).toBe(false);
    expect(runsDaily(undefined)).toBe(false);
  });
  it('keeps the old name filter and adds annotated jobs', () => {
    expect(isBackupJob(job('backup-postgres-x'))).toBe(true);
    expect(isBackupJob(job('health-keeper'))).toBe(false);
    expect(isBackupJob(job('cot-retention', { annotations: { 'onemind.io/backup-state': 'active' } }))).toBe(true);
  });
  it('summarizes review, setup, retired and ok separately', () => {
    const items = [
      classifyBackup(job('ok-backup', { status: { lastSuccessfulTime: hoursAgo(1) } }), NOW),
      classifyBackup(job('bad-backup'), NOW),
      classifyBackup(job('old-backup', { annotations: { 'onemind.io/backup-state': 'superseded' } }), NOW),
      classifyBackup(job('new-backup', { annotations: { 'onemind.io/backup-state': 'pending-setup' } }), NOW),
    ];
    expect(summarizeBackups(items)).toEqual({ review: 1, setup: 1, retired: 1, ok: 1 });
  });
});
