// Single severity and dedupe policy for every OneMind notification producer.
// Pure functions only: no I/O, so producers and tests share the same rules.

export type NotificationSeverity = 'critical' | 'high' | 'normal' | 'low';
export type NotificationTopic =
  | 'alerts'
  | 'deployments'
  | 'backups'
  | 'digest'
  | 'agents';

export type ActivityKind =
  | 'deployment.failed'
  | 'deployment.degraded'
  | 'deployment.out-of-sync'
  | 'deployment.succeeded'
  | 'backup.failed'
  | 'backup.suspended'
  | 'backup.succeeded'
  | 'agent.completed';

export type ActivityPolicy = {
  // notify: lands in the inbox. timeline: recorded for the digest only.
  delivery: 'notify' | 'timeline';
  severity: NotificationSeverity;
  topic: NotificationTopic;
  title: string;
};

// An application that stays OutOfSync shorter than this is normal reconcile noise.
export const OUT_OF_SYNC_GRACE_MS = 15 * 60 * 1000;

const policies: Record<ActivityKind, ActivityPolicy> = {
  'deployment.failed': {
    delivery: 'notify',
    severity: 'high',
    topic: 'deployments',
    title: 'Deployment needs attention',
  },
  'deployment.degraded': {
    delivery: 'notify',
    severity: 'high',
    topic: 'deployments',
    title: 'Application degraded',
  },
  'deployment.out-of-sync': {
    delivery: 'notify',
    severity: 'high',
    topic: 'deployments',
    title: 'Application out of sync',
  },
  // Successes are never notifications. They feed the timeline and the daily digest.
  'deployment.succeeded': {
    delivery: 'timeline',
    severity: 'low',
    topic: 'deployments',
    title: 'Deployment completed',
  },
  'backup.failed': {
    delivery: 'notify',
    severity: 'high',
    topic: 'backups',
    title: 'Backup failed',
  },
  'backup.suspended': {
    delivery: 'notify',
    severity: 'high',
    topic: 'backups',
    title: 'Backup schedule suspended',
  },
  'backup.succeeded': {
    delivery: 'timeline',
    severity: 'low',
    topic: 'backups',
    title: 'Backup completed',
  },
  'agent.completed': {
    delivery: 'notify',
    severity: 'normal',
    topic: 'agents',
    title: 'Agent work completed',
  },
};

export const activityKinds = Object.keys(policies) as ActivityKind[];

export function isActivityKind(value: unknown): value is ActivityKind {
  return typeof value === 'string' && Object.hasOwn(policies, value);
}

export function activityPolicy(kind: ActivityKind): ActivityPolicy {
  return policies[kind];
}

// Alertmanager severity label to inbox severity. Anything else (none, the
// Watchdog heartbeat, missing label) is not a human-facing alert.
export function alertSeverity(
  label: string | undefined,
): NotificationSeverity | undefined {
  switch (label?.toLowerCase()) {
    case 'critical':
      return 'critical';
    case 'warning':
      return 'high';
    case 'info':
      return 'normal';
    default:
      return undefined;
  }
}

// High and critical items also push a Signal; the rest only appear in the inbox.
export function isSignalSeverity(severity: NotificationSeverity): boolean {
  return severity === 'critical' || severity === 'high';
}

const segment = (value: string) =>
  value.replace(/[^a-zA-Z0-9._:/-]/g, '_').slice(0, 120);

// Stable scope per source object: repeats update one notification instead of adding another.
export const argoScope = (application: string) =>
  `argo:${segment(application)}`;
export const backupScope = (cronJob: string) => `backup:${segment(cronJob)}`;
export const alertScope = (fingerprint: string) =>
  `alert:${segment(fingerprint)}`;
export const digestScope = (day: string) => `digest:deployments:${segment(day)}`;

export function scopeFor(kind: ActivityKind, resource: string, id: string) {
  if (kind.startsWith('deployment.')) return argoScope(resource);
  if (kind.startsWith('backup.')) return backupScope(resource);
  return `${kind}:${segment(id)}`;
}

// The old producers sent "Deployment completed: <app>" at low severity for every sync.
// The backlog cleanup script uses the same test (see scripts/mark-deployment-notifications-read.cjs).
export function isDeploymentSuccessNotification(payload: {
  title?: unknown;
  topic?: unknown;
  severity?: unknown;
}): boolean {
  return (
    payload.topic === 'deployments' &&
    payload.severity === 'low' &&
    typeof payload.title === 'string' &&
    payload.title.startsWith('Deployment completed:')
  );
}

// Signals channel and message contract for the portal (home status strip, Live Operations).
export const ALERT_SIGNAL_CHANNEL = 'onemind:alerts';
export type AlertSignal = {
  action: 'alert.firing' | 'alert.resolved';
  source: 'alertmanager' | 'argocd' | 'backup';
  key: string;
  severity: NotificationSeverity;
  title: string;
  at: string;
};

// Events bus topic for items that are recorded but never notified.
export const ACTIVITY_EVENT_TOPIC = 'onemind.activity';
export type ActivityEvent = {
  id: string;
  kind: ActivityKind;
  resource: string;
  occurredAt: string;
  source: string;
};
