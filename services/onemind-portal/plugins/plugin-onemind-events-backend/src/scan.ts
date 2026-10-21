import { createHash } from 'node:crypto';
import {
  activityPolicy,
  argoScope,
  backupScope,
  digestScope,
  isSignalSeverity,
  scopeFor,
  type ActivityKind,
  type AlertSignal,
  type NotificationSeverity,
  type NotificationTopic,
} from '@internal/plugin-onemind-common-backend';
import { planConditions, type ActiveCondition, type Reason } from './conditions';
import type { Db, Store, DigestSummary } from './store';

export type NotificationPayload = {
  title: string;
  description: string;
  severity: NotificationSeverity;
  topic: NotificationTopic;
  scope: string;
  link: string;
};

// Where results go. The module wires these to the Notifications and Signals services.
export type Sinks = {
  notify(payload: NotificationPayload): Promise<void>;
  signal(message: AlertSignal): Promise<void>;
};

export type Observation = {
  key: string;
  resource: string;
  kind: ActivityKind;
  occurredAt: string;
  link: string;
};

export type AppView = {
  name: string;
  health?: string;
  sync?: string;
  link: string;
};

export type SourceName = 'argocd' | 'backstage-backup';
const signalSource = (source: SourceName) => (source === 'argocd' ? 'argocd' : 'backup') as AlertSignal['source'];

export const hash = (value: string) =>
  createHash('sha256').update(value).digest('hex');

export async function ensureDeliveryTables(db: Db) {
  const create = async (name: string, build: Parameters<Db['schema']['createTable']>[1]) => {
    for (let attempt = 0; ; attempt += 1) {
      if (await db.schema.hasTable(name)) return;
      try {
        await db.schema.createTable(name, build);
        return;
      } catch (error) {
        if (attempt >= 3) throw error;
        await new Promise(resolve => setTimeout(resolve, 250 * (attempt + 1)));
      }
    }
  };
  await create('onemind_operation_delivery', table => {
    table.string('id', 64).primary();
    table.string('source', 32).notNullable();
    table.string('recorded_at', 40).notNullable();
  });
  await create('onemind_operation_observer', table => {
    table.string('id', 32).primary();
    table.text('payload').notNullable();
  });
}

function eventPayload(value: Observation, source: SourceName): NotificationPayload {
  const policy = activityPolicy(value.kind);
  const note =
    source === 'argocd'
      ? 'Operation completion is not proof of current application health.'
      : 'Job success is not separate proof of archive restoration; inspect the backup report.';
  return {
    title: `${policy.title}: ${value.resource}`,
    description: `Observed ${value.kind} at ${value.occurredAt}. ${note}`,
    severity: policy.severity,
    topic: policy.topic,
    scope: scopeFor(value.kind, value.resource, value.key),
    link: value.link,
  };
}

const signalOf = (
  source: AlertSignal['source'],
  action: AlertSignal['action'],
  key: string,
  severity: NotificationSeverity,
  title: string,
  now: Date,
): AlertSignal => ({ action, source, key, severity, title, at: now.toISOString() });

// Applies the severity policy to terminal operations seen in one scan. The first scan of a
// source only records a baseline so history is never replayed into the inbox.
export async function processObservations(args: {
  db: Db;
  store: Store;
  source: SourceName;
  values: Observation[];
  sinks: Sinks;
  now?: () => Date;
}): Promise<'baseline-established' | 'observing'> {
  const { db, store, source, values, sinks } = args;
  const now = args.now ?? (() => new Date());
  const baselineId = hash(`baseline:${source}`);
  const baseline = await db('onemind_operation_delivery').where({ id: baselineId }).first();
  if (!baseline) {
    await db.transaction(async tx => {
      for (const value of values) {
        await tx('onemind_operation_delivery')
          .insert({ id: hash(`${source}:${value.key}`), source, recorded_at: now().toISOString() })
          .onConflict('id')
          .ignore();
      }
      await tx('onemind_operation_delivery')
        .insert({ id: baselineId, source, recorded_at: now().toISOString() })
        .onConflict('id')
        .ignore();
    });
    // Successes seen during the baseline still belong on the timeline.
    for (const value of values) {
      if (activityPolicy(value.kind).delivery === 'timeline') {
        await store.timeline.record({ id: hash(`${source}:${value.key}`), kind: value.kind, resource: value.resource, occurredAt: value.occurredAt, source });
      }
    }
    return 'baseline-established';
  }
  const ordered = [...values].sort((a, b) => Date.parse(a.occurredAt) - Date.parse(b.occurredAt));
  for (const value of ordered) {
    const id = hash(`${source}:${value.key}`);
    if (await db('onemind_operation_delivery').where({ id }).first()) continue;
    const policy = activityPolicy(value.kind);
    const recordedAt = baseline.recorded_at as string;
    if (policy.delivery === 'timeline') {
      // Timeline-only: never a notification.
      await store.timeline.record({ id, kind: value.kind, resource: value.resource, occurredAt: value.occurredAt, source });
      const scope = scopeFor(value.kind.startsWith('backup.') ? 'backup.failed' : 'deployment.failed', value.resource, value.key);
      const open = (await store.conditions.all()).find(row => row.scope === scope && row.reason === 'failed');
      if (open && Date.parse(value.occurredAt) >= Date.parse(recordedAt)) {
        await announceRecovery(sinks, open, now());
        await store.conditions.remove(scope);
      }
    } else if (Date.parse(value.occurredAt) >= Date.parse(recordedAt)) {
      const payload = eventPayload(value, source);
      await sinks.notify(payload);
      // A failed operation starts an incident; recovery is announced once the object is healthy.
      await store.conditions.put({ scope: payload.scope, reason: 'failed', since: now().getTime(), notified: true });
      if (isSignalSeverity(policy.severity)) {
        await sinks.signal(signalOf(signalSource(source), 'alert.firing', payload.scope, policy.severity, payload.title, now()));
      }
    }
    // send acknowledges processing, not delivery. A crash before this insert may retry;
    // the notification scope collapses duplicate occurrences.
    await db('onemind_operation_delivery')
      .insert({ id, source, recorded_at: now().toISOString() })
      .onConflict('id')
      .ignore();
  }
  return 'observing';
}

const runbook = '/docs/default/component/onemind-runbooks/resilience/';

// Quiet follow-up on the same scope, so the earlier alert is replaced rather than duplicated.
async function announceRecovery(sinks: Sinks, row: { scope: string }, now: Date) {
  const isBackup = row.scope.startsWith('backup:');
  const resource = row.scope.replace(/^(argo|backup):/, '');
  const title = `Recovered: ${resource}`;
  await sinks.notify({
    title,
    description: 'Observed healthy again. This update replaces the earlier alert.',
    severity: 'low',
    topic: isBackup ? 'backups' : 'deployments',
    scope: row.scope,
    link: isBackup ? runbook : '/today',
  });
  await sinks.signal(signalOf(isBackup ? 'backup' : 'argocd', 'alert.resolved', row.scope, 'high', title, now));
}

function conditionPayload(reason: Reason, resource: string, scope: string, link: string): NotificationPayload {
  const kinds: Record<Reason, ActivityKind> = {
    degraded: 'deployment.degraded',
    'out-of-sync': 'deployment.out-of-sync',
    suspended: 'backup.suspended',
    failed: 'deployment.failed',
  };
  const kind = kinds[reason];
  const policy = activityPolicy(kind);
  const detail: Record<Reason, string> = {
    degraded: 'Argo CD reports the application health as Degraded.',
    'out-of-sync': 'Argo CD has reported this application out of sync for more than 15 minutes.',
    suspended: 'The backup CronJob is suspended, so no new backups are being taken.',
    failed: 'The last operation failed.',
  };
  return { title: `${policy.title}: ${resource}`, description: detail[reason], severity: policy.severity, topic: policy.topic, scope, link };
}

// State-based conditions: degraded, out of sync past the grace period, suspended backups.
// One stored row per scope keeps a persistent problem to a single notification and sends a
// quiet "Recovered" update on the same scope when it clears.
export async function processConditions(args: {
  store: Store;
  active: (ActiveCondition & { link: string })[];
  healthy: string[];
  topic: 'deployments' | 'backups';
  sinks: Sinks;
  now?: () => Date;
}) {
  const now = args.now ?? (() => new Date());
  const links = new Map(args.active.map(item => [item.scope, item.link]));
  // A failed backup job is cleared by the next successful job, not by the CronJob merely being unsuspended.
  const rows = (await args.store.conditions.all()).filter(row =>
    args.topic === 'backups' ? row.scope.startsWith('backup:') && row.reason !== 'failed' : row.scope.startsWith('argo:'),
  );
  const plan = planConditions({ now: now().getTime(), rows, active: args.active, healthy: args.healthy });
  for (const item of plan.raise) {
    const payload = conditionPayload(item.row.reason, item.resource, item.row.scope, links.get(item.row.scope) ?? '/today');
    await args.sinks.notify(payload);
    await args.sinks.signal(signalOf(item.row.scope.startsWith('backup:') ? 'backup' : 'argocd', 'alert.firing', item.row.scope, payload.severity, payload.title, now()));
  }
  for (const row of plan.recover) await announceRecovery(args.sinks, row, now());
  for (const row of plan.upsert) await args.store.conditions.put(row);
  for (const scope of plan.remove) await args.store.conditions.remove(scope);
  return plan;
}

// Argo application view to conditions. Degraded wins over out-of-sync. Only Synced and
// Healthy is "healthy"; Progressing, Missing and Unknown stay undecided.
export function classifyApplications(apps: AppView[]) {
  const active: (ActiveCondition & { link: string })[] = [];
  const healthy: string[] = [];
  for (const app of apps) {
    const scope = argoScope(app.name);
    if (app.health === 'Degraded') active.push({ scope, reason: 'degraded', resource: app.name, link: app.link });
    else if (app.sync === 'OutOfSync') active.push({ scope, reason: 'out-of-sync', resource: app.name, link: app.link });
    else if (app.sync === 'Synced' && app.health === 'Healthy') healthy.push(scope);
  }
  return { active, healthy };
}

export function classifyBackups(cronJobs: { name: string; suspended: boolean; link: string }[]) {
  const active: (ActiveCondition & { link: string })[] = [];
  const healthy: string[] = [];
  for (const job of cronJobs) {
    const scope = backupScope(`backstage/${job.name}`);
    if (job.suspended) active.push({ scope, reason: 'suspended', resource: `backstage/${job.name}`, link: job.link });
    else healthy.push(scope);
  }
  return { active, healthy };
}

// Daily digest: one normal notification summarising successful deployments of the last 24h.
export function digestPayload(day: string, summary: DigestSummary): NotificationPayload {
  const top = summary.byResource.slice(0, 8).map(item => `${item.resource} x${item.count}`).join(', ');
  const more = summary.byResource.length > 8 ? `, and ${summary.byResource.length - 8} more` : '';
  return {
    title: `Daily digest: ${summary.total} deployment${summary.total === 1 ? '' : 's'} succeeded`,
    description: `In the last 24 hours: ${top}${more}.`,
    severity: 'normal',
    topic: 'digest',
    scope: digestScope(day),
    link: '/today',
  };
}

export async function runDigest(args: {
  store: Store;
  sinks: Sinks;
  now?: () => Date;
  hourUtc: number;
}): Promise<'sent' | 'empty' | 'already-sent' | 'too-early'> {
  const now = (args.now ?? (() => new Date()))();
  if (now.getUTCHours() < args.hourUtc) return 'too-early';
  const day = now.toISOString().slice(0, 10);
  if (await args.store.digest.sent(day)) return 'already-sent';
  const since = new Date(now.getTime() - 24 * 3600_000).toISOString();
  const summary = await args.store.timeline.summary(since, 'deployment.succeeded');
  if (summary.total > 0) await args.sinks.notify(digestPayload(day, summary));
  // Marked only after the send succeeded, so a failed send is retried on the next run.
  await args.store.digest.markSent(day, summary.total);
  return summary.total > 0 ? 'sent' : 'empty';
}
