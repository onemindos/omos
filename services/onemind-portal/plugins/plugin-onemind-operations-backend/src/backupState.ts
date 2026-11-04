// One definition of "what state is this backup in", shared by /status (Home, Today)
// and /live. Pure functions: the CronJob annotations in Git decide lifecycle, the
// CronJob status decides health.

export const BACKUP_STATE_ANNOTATION = 'onemind.io/backup-state';
export const SUPERSEDED_BY_ANNOTATION = 'onemind.io/superseded-by';
export const PENDING_REASON_ANNOTATION = 'onemind.io/pending-reason';
// Optional per-job override of the overdue limit, in hours.
export const MAX_AGE_ANNOTATION = 'onemind.io/backup-max-age-hours';

// A daily job that has not succeeded for 36 hours missed at least one run.
export const DEFAULT_DAILY_MAX_AGE_MS = 36 * 3600000;
// A failed run is only "failing" once it is not simply still starting up.
const FAILURE_GRACE_MS = 2 * 60000;

/** active: runs and is expected to work. retired: replaced on purpose. needs-setup: waiting for something only a person can provide. */
export type BackupLifecycle = 'active' | 'retired' | 'needs-setup';
export type BackupHealth =
  | 'ok'
  | 'failing'
  | 'overdue'
  | 'never-ran'
  | 'suspended'
  | 'inactive';
/** review: counts as an attention item. setup: separate "needs setup" item. none: nothing to do. */
export type BackupAttention = 'none' | 'review' | 'setup';

export type CronJobLike = {
  metadata: {
    name: string;
    namespace?: string;
    annotations?: Record<string, string>;
  };
  spec?: { suspend?: boolean; schedule?: string };
  status?: {
    lastSuccessfulTime?: string;
    lastScheduleTime?: string;
    active?: unknown[];
  };
};

export type BackupItem = {
  /** `namespace/name`, as the existing /status route reports it. */
  name: string;
  lastSuccess: string | null;
  suspended: boolean;
  state: BackupLifecycle;
  health: BackupHealth;
  attention: BackupAttention;
  supersededBy?: string;
  pendingReason?: string;
};

const clean = (value: string | undefined) => {
  const text = value?.trim();
  return text ? text.slice(0, 300) : undefined;
};

/** True for schedules with `*` in day-of-month, month and day-of-week (runs every day or more often). */
export function runsDaily(schedule: string | undefined) {
  const fields = schedule?.trim().split(/\s+/) ?? [];
  return fields.length === 5 && fields.slice(2).every(field => field === '*');
}

export function maxAgeMs(job: CronJobLike): number | undefined {
  const override = Number(job.metadata.annotations?.[MAX_AGE_ANNOTATION]);
  if (Number.isFinite(override) && override > 0) return override * 3600000;
  return runsDaily(job.spec?.schedule) ? DEFAULT_DAILY_MAX_AGE_MS : undefined;
}

function activeHealth(job: CronJobLike, now: number): BackupHealth {
  if (job.spec?.suspend === true) return 'suspended';
  const success = job.status?.lastSuccessfulTime;
  if (!success) return 'never-ran';
  const successAt = Date.parse(success);
  const scheduledAt = Date.parse(job.status?.lastScheduleTime ?? '');
  if (
    Number.isFinite(scheduledAt) &&
    scheduledAt - successAt > FAILURE_GRACE_MS &&
    !job.status?.active?.length
  ) {
    return 'failing';
  }
  const limit = maxAgeMs(job);
  if (limit !== undefined && now - successAt > limit) return 'overdue';
  return 'ok';
}

export function classifyBackup(job: CronJobLike, now = Date.now()): BackupItem {
  const annotations = job.metadata.annotations ?? {};
  const declared = annotations[BACKUP_STATE_ANNOTATION]?.trim();
  const base = {
    name: `${job.metadata.namespace}/${job.metadata.name}`,
    lastSuccess: job.status?.lastSuccessfulTime ?? null,
    suspended: job.spec?.suspend ?? false,
  };
  if (declared === 'superseded') {
    return {
      ...base,
      state: 'retired',
      health: 'inactive',
      attention: 'none',
      supersededBy: clean(annotations[SUPERSEDED_BY_ANNOTATION]),
    };
  }
  if (declared === 'pending-setup') {
    return {
      ...base,
      state: 'needs-setup',
      health: base.suspended ? 'suspended' : activeHealth(job, now),
      attention: 'setup',
      pendingReason: clean(annotations[PENDING_REASON_ANNOTATION]),
    };
  }
  // Missing or unrecognized value (including a typo) is treated as active, so it fails loud.
  const health = activeHealth(job, now);
  return {
    ...base,
    state: 'active',
    health,
    attention: health === 'ok' ? 'none' : 'review',
  };
}

/** The same backup jobs the old route listed, plus anything that declares a backup state. */
export const isBackupJob = (job: CronJobLike) =>
  job.metadata.name.includes('backup') ||
  job.metadata.annotations?.[BACKUP_STATE_ANNOTATION] !== undefined;

export function summarizeBackups(items: BackupItem[]) {
  return {
    review: items.filter(item => item.attention === 'review').length,
    setup: items.filter(item => item.attention === 'setup').length,
    retired: items.filter(item => item.state === 'retired').length,
    ok: items.filter(item => item.attention === 'none' && item.state === 'active')
      .length,
  };
}
