// Pure builders: raw upstream shapes in, Live Operations sections out.
// No I/O and no clock reads (callers pass `now`), so fixtures can test every rule.
import { classifyBackup, CronJobLike, isBackupJob, BackupItem } from './backupState';
import type { NodeUsage, PvcItem } from './liveQueries';
import { parseQuantity } from './quantity';
import type {
  AlertSeverity,
  ArgoItem,
  ArgoReason,
  AttentionItem,
  EventItem,
  LiveAlert,
  NodeItem,
  WorkloadIssue,
  WorkloadReason,
} from './types';

export const ARGOCD_URL = 'https://argocd.tail717ba0.ts.net';
export const HEADLAMP_URL = 'https://headlamp.tail717ba0.ts.net/c/main/workloads';
export const GRAFANA_URL = 'https://grafana.tail717ba0.ts.net';
export const ARGO_PROGRESSING_LIMIT_MS = 10 * 60000;
export const WORKLOAD_PENDING_GRACE_MS = 2 * 60000;
export const WORKLOAD_NOT_READY_GRACE_MS = 5 * 60000;
export const RECENT_RESTART_MS = 3600000;
export const EVENT_WINDOW_MINUTES = 30;
export const EVENT_LIMIT = 50;

// ---- Alertmanager ------------------------------------------------------------------

export type AmAlert = {
  fingerprint?: string;
  labels?: Record<string, string>;
  annotations?: Record<string, string>;
  startsAt?: string;
  status?: { state?: string };
};

// Heartbeat and inhibition alerts exist to keep Alertmanager honest; they are not problems.
const NOISE_ALERTS = new Set(['Watchdog', 'InfoInhibitor']);
const severities: Record<string, AlertSeverity> = {
  critical: 'critical',
  warning: 'warning',
  info: 'info',
};
const WORKLOAD_LABELS = ['deployment', 'statefulset', 'daemonset', 'job_name', 'cronjob', 'pod', 'service'];

export function buildAlerts(raw: AmAlert[]): LiveAlert[] {
  const alerts: LiveAlert[] = [];
  for (const alert of raw) {
    const labels = alert.labels ?? {};
    const severity = severities[(labels.severity ?? '').toLowerCase()];
    const name = labels.alertname;
    if (!name || !severity || NOISE_ALERTS.has(name)) continue;
    if (alert.status?.state && alert.status.state !== 'active') continue;
    const annotations = alert.annotations ?? {};
    const workloadLabel = WORKLOAD_LABELS.find(label => labels[label]);
    alerts.push({
      fingerprint: alert.fingerprint ?? `${name}:${labels.namespace ?? ''}:${labels.pod ?? ''}`,
      name,
      severity,
      summary: (annotations.summary ?? annotations.description ?? name).slice(0, 300),
      namespace: labels.namespace,
      workload: workloadLabel ? labels[workloadLabel] : undefined,
      startsAt: alert.startsAt ?? '',
      runbookUrl: /^https?:\/\//.test(annotations.runbook_url ?? '') ? annotations.runbook_url : undefined,
      grafanaUrl: `${GRAFANA_URL}/alerting/list?search=${encodeURIComponent(name)}`,
    });
  }
  const order = { critical: 0, warning: 1, info: 2 };
  return alerts.sort(
    (a, b) => order[a.severity] - order[b.severity] || a.startsAt.localeCompare(b.startsAt),
  );
}

// ---- Argo CD -----------------------------------------------------------------------

export type ArgoApp = {
  metadata: { name: string };
  status?: {
    health?: { status?: string; lastTransitionTime?: string };
    sync?: { status?: string };
    operationState?: { phase?: string; message?: string; startedAt?: string };
  };
};

export const argoAppUrl = (name: string) =>
  `${ARGOCD_URL}/applications/argocd/${encodeURIComponent(name)}`;

export function classifyArgo(app: ArgoApp, now: number): ArgoItem | undefined {
  const health = app.status?.health?.status ?? 'Unknown';
  const sync = app.status?.sync?.status ?? 'Unknown';
  const op = app.status?.operationState;
  const since = app.status?.health?.lastTransitionTime ?? op?.startedAt;
  const reasons: ArgoReason[] = [];
  if (health === 'Degraded' || health === 'Missing') reasons.push('Degraded');
  if (health === 'Progressing' && since && now - Date.parse(since) > ARGO_PROGRESSING_LIMIT_MS) {
    reasons.push('Progressing');
  }
  if (sync === 'OutOfSync') reasons.push('OutOfSync');
  if (op?.phase === 'Failed' || op?.phase === 'Error') reasons.push('Sync failed');
  if (health === 'Unknown' || sync === 'Unknown') reasons.push('Unknown');
  if (!reasons.length) return undefined;
  return {
    name: app.metadata.name,
    health,
    sync,
    reasons,
    severity: reasons.includes('Degraded') || reasons.includes('Sync failed') ? 'critical' : 'warning',
    since,
    message: op?.message?.slice(0, 200),
    url: argoAppUrl(app.metadata.name),
  };
}

export const buildArgo = (apps: ArgoApp[], now: number) =>
  apps.flatMap(app => classifyArgo(app, now) ?? []);

/** What the old /status route called "deployments": anything not Healthy and Synced. Kept for current callers. */
export const legacyDeployments = (apps: ArgoApp[]) =>
  apps
    .filter(x => x.status?.health?.status !== 'Healthy' || x.status?.sync?.status !== 'Synced')
    .map(x => ({
      name: x.metadata.name,
      health: x.status?.health?.status ?? 'Unknown',
      sync: x.status?.sync?.status ?? 'Unknown',
    }));

// ---- Workloads ---------------------------------------------------------------------

type ContainerStatus = {
  ready?: boolean;
  restartCount?: number;
  state?: { waiting?: { reason?: string } };
  lastState?: { terminated?: { reason?: string; finishedAt?: string } };
};
export type KubePod = {
  metadata: {
    name: string;
    namespace: string;
    labels?: Record<string, string>;
    creationTimestamp?: string;
    deletionTimestamp?: string;
  };
  spec?: {
    nodeName?: string;
    containers?: { resources?: { requests?: Record<string, string> } }[];
  };
  status?: {
    phase?: string;
    conditions?: { type: string; status: string; lastTransitionTime?: string }[];
    containerStatuses?: ContainerStatus[];
    initContainerStatuses?: ContainerStatus[];
  };
};
export type KubeController = {
  metadata: { name: string; namespace?: string };
  spec?: { replicas?: number };
  status?: { readyReplicas?: number };
};

const WAITING: Record<string, WorkloadReason> = {
  CrashLoopBackOff: 'CrashLoopBackOff',
  ImagePullBackOff: 'ImagePullBackOff',
  ErrImagePull: 'ImagePullBackOff',
  CreateContainerConfigError: 'ConfigError',
  CreateContainerError: 'ConfigError',
  InvalidImageName: 'ConfigError',
};

export function classifyPod(pod: KubePod, now: number): WorkloadIssue | undefined {
  const phase = pod.status?.phase;
  if (pod.metadata.deletionTimestamp || phase === 'Succeeded' || phase === 'Failed') return undefined;
  const reasons = new Set<WorkloadReason>();
  const statuses = [
    ...(pod.status?.containerStatuses ?? []),
    ...(pod.status?.initContainerStatuses ?? []),
  ];
  let restarts = 0;
  let lastRestart = 0;
  for (const status of statuses) {
    const waiting = WAITING[status.state?.waiting?.reason ?? ''];
    if (waiting) reasons.add(waiting);
    restarts += status.restartCount ?? 0;
    const terminated = status.lastState?.terminated;
    const finished = Date.parse(terminated?.finishedAt ?? '');
    if (Number.isFinite(finished) && (status.restartCount ?? 0) > 0) {
      lastRestart = Math.max(lastRestart, finished);
      if (now - finished <= RECENT_RESTART_MS) {
        reasons.add(terminated?.reason === 'OOMKilled' ? 'OOMKilled' : 'Restarted');
      }
    }
  }
  const age = now - Date.parse(pod.metadata.creationTimestamp ?? '');
  if (phase === 'Pending' && !reasons.size && age > WORKLOAD_PENDING_GRACE_MS) reasons.add('Pending');
  const ready = pod.status?.conditions?.find(c => c.type === 'Ready');
  const notReadySince = Date.parse(ready?.lastTransitionTime ?? '');
  if (
    phase === 'Running' &&
    ready?.status === 'False' &&
    !reasons.has('CrashLoopBackOff') &&
    Number.isFinite(notReadySince) &&
    now - notReadySince > WORKLOAD_NOT_READY_GRACE_MS
  ) {
    reasons.add('NotReady');
  }
  if (!reasons.size) return undefined;
  const labels = pod.metadata.labels ?? {};
  return {
    namespace: pod.metadata.namespace,
    pod: pod.metadata.name,
    app: labels['app.kubernetes.io/instance'] ?? labels.app,
    reasons: [...reasons],
    restarts,
    lastRestartAt: lastRestart ? new Date(lastRestart).toISOString() : undefined,
    severity: reasons.has('CrashLoopBackOff') ? 'critical' : 'warning',
  };
}

/** `namespace/name` of controllers with fewer ready replicas than desired (the old /status `unavailable`). */
export const notReadyControllers = (...lists: KubeController[][]) =>
  lists
    .flat()
    .filter(x => (x.spec?.replicas ?? 1) > (x.status?.readyReplicas ?? 0))
    .map(x => `${x.metadata.namespace}/${x.metadata.name}`);

export function buildWorkloads(
  pods: KubePod[],
  deployments: KubeController[],
  statefulsets: KubeController[],
  now: number,
) {
  const controllers = [...deployments, ...statefulsets];
  return {
    issues: pods
      .flatMap(pod => classifyPod(pod, now) ?? [])
      .sort((a, b) => Number(b.severity === 'critical') - Number(a.severity === 'critical') || b.restarts - a.restarts),
    notReady: notReadyControllers(deployments, statefulsets),
    total: controllers.length,
  };
}

// ---- Nodes -------------------------------------------------------------------------

export type KubeNode = {
  metadata: { name: string };
  status?: {
    conditions?: { type: string; status: string }[];
    allocatable?: Record<string, string>;
  };
};

const PRESSURE = ['MemoryPressure', 'DiskPressure', 'PIDPressure', 'NetworkUnavailable'];
const percent = (part: number, whole: number | undefined) =>
  whole ? Math.round((part / whole) * 1000) / 10 : null;

export function buildNodes(
  nodes: KubeNode[],
  pods: KubePod[],
  usage: Record<string, NodeUsage> | undefined,
): NodeItem[] {
  return nodes.map(node => {
    const name = node.metadata.name;
    const conditions = node.status?.conditions ?? [];
    const onNode = pods.filter(
      pod =>
        pod.spec?.nodeName === name &&
        pod.status?.phase !== 'Succeeded' &&
        pod.status?.phase !== 'Failed',
    );
    let cpu = 0;
    let memory = 0;
    for (const pod of onNode) {
      for (const container of pod.spec?.containers ?? []) {
        cpu += parseQuantity(container.resources?.requests?.cpu) ?? 0;
        memory += parseQuantity(container.resources?.requests?.memory) ?? 0;
      }
    }
    const known = usage?.[name];
    return {
      name,
      ready: conditions.find(c => c.type === 'Ready')?.status === 'True',
      pressure: conditions.filter(c => PRESSURE.includes(c.type) && c.status === 'True').map(c => c.type),
      cpuPercent: known?.cpuPercent ?? null,
      memoryPercent: known?.memoryPercent ?? null,
      diskPercent: known?.diskPercent ?? null,
      cpuRequestedPercent: percent(cpu, parseQuantity(node.status?.allocatable?.cpu)),
      memoryRequestedPercent: percent(memory, parseQuantity(node.status?.allocatable?.memory)),
      pods: onNode.length,
    };
  });
}

// ---- Events ------------------------------------------------------------------------

export type KubeEvent = {
  type?: string;
  reason?: string;
  message?: string;
  count?: number;
  lastTimestamp?: string;
  eventTime?: string;
  metadata: { namespace?: string; creationTimestamp?: string };
  involvedObject?: { kind?: string; name?: string; namespace?: string };
  series?: { count?: number; lastObservedTime?: string };
};

export function buildEvents(raw: KubeEvent[], now: number): EventItem[] {
  const cutoff = now - EVENT_WINDOW_MINUTES * 60000;
  const grouped = new Map<string, EventItem>();
  for (const event of raw) {
    if (event.type !== 'Warning') continue;
    const seen =
      event.series?.lastObservedTime ??
      event.lastTimestamp ??
      event.eventTime ??
      event.metadata.creationTimestamp ??
      '';
    if (!(Date.parse(seen) >= cutoff)) continue;
    const object = `${event.involvedObject?.kind ?? 'Object'}/${event.involvedObject?.name ?? 'unknown'}`;
    const namespace = event.involvedObject?.namespace ?? event.metadata.namespace ?? '';
    const message = (event.message ?? '').slice(0, 300);
    const key = [namespace, object, event.reason, message].join('|');
    const count = event.series?.count ?? event.count ?? 1;
    const known = grouped.get(key);
    if (known) {
      known.count += count;
      if (seen > known.lastSeen) known.lastSeen = seen;
    } else {
      grouped.set(key, { namespace, object, reason: event.reason ?? 'Unknown', message, count, lastSeen: seen });
    }
  }
  return [...grouped.values()]
    .sort((a, b) => b.lastSeen.localeCompare(a.lastSeen))
    .slice(0, EVENT_LIMIT);
}

// ---- Backups -----------------------------------------------------------------------

export const buildBackups = (jobs: CronJobLike[], now: number): BackupItem[] =>
  jobs.filter(isBackupJob).map(job => classifyBackup(job, now));

const BACKUP_LABEL: Record<string, string> = {
  failing: 'last run failed',
  overdue: 'overdue',
  'never-ran': 'no successful run recorded',
  suspended: 'suspended',
};

export function backupAttention(items: BackupItem[]): { attention: AttentionItem[]; needsSetup: AttentionItem[] } {
  const attention: AttentionItem[] = [];
  const needsSetup: AttentionItem[] = [];
  for (const item of items) {
    const key = `backup:${item.name}`;
    if (item.attention === 'review') {
      attention.push({
        key,
        kind: 'backup',
        severity: item.health === 'failing' ? 'critical' : 'warning',
        title: `Backup ${item.name}: ${BACKUP_LABEL[item.health] ?? item.health}`,
        detail: item.lastSuccess
          ? `Last successful job: ${new Date(item.lastSuccess).toISOString()}`
          : 'Inspect the backup schedule and job logs.',
        url: HEADLAMP_URL,
        namespace: item.name.split('/')[0],
      });
    } else if (item.attention === 'setup') {
      needsSetup.push({
        key,
        kind: 'backup',
        severity: 'warning',
        title: `Backup ${item.name}: needs setup`,
        detail: item.pendingReason ?? 'Waiting for setup that only a person can do.',
        url: HEADLAMP_URL,
        namespace: item.name.split('/')[0],
      });
    }
  }
  return { attention, needsSetup };
}

// ---- Attention list (the single source Home, Today and Live share) -----------------

export function pvcAttention(pvcs: PvcItem[]): AttentionItem[] {
  return pvcs.map(pvc => ({
    key: `storage:${pvc.namespace}/${pvc.name}`,
    kind: 'storage' as const,
    severity: pvc.severity === 'critical' ? ('critical' as const) : ('warning' as const),
    title: `Volume ${pvc.namespace}/${pvc.name}: ${pvc.usedPercent}% full`,
    detail: pvc.fillsWithin4Days
      ? 'Projected to fill within four days. Inspect usage or expand the volume.'
      : 'Inspect usage or expand the volume.',
    url: HEADLAMP_URL,
    namespace: pvc.namespace,
  }));
}

const ARGO_DETAIL: Record<ArgoReason, string> = {
  Degraded: 'Inspect the deployment health, sync difference, and last operation.',
  OutOfSync: 'Live state differs from Git. Inspect the sync difference.',
  Progressing: 'Still progressing after 10 minutes. Inspect rollout and events.',
  'Sync failed': 'The last sync operation failed. Inspect the operation message.',
  Unknown: 'Argo CD reports no health or sync state.',
};

export function argoAttention(items: ArgoItem[]): AttentionItem[] {
  return items.map(item => ({
    key: `deployment:${item.name}`,
    kind: 'deployment' as const,
    severity: item.severity,
    title: `${item.name}: ${item.health} / ${item.sync}`,
    detail: ARGO_DETAIL[item.reasons[0]],
    url: item.url,
  }));
}

export const alertAttention = (alerts: LiveAlert[]): AttentionItem[] =>
  alerts
    .filter(alert => alert.severity !== 'info')
    .map(alert => ({
      key: `alert:${alert.fingerprint}`,
      kind: 'alert' as const,
      severity: alert.severity === 'critical' ? ('critical' as const) : ('warning' as const),
      title: alert.name,
      detail: alert.summary,
      url: alert.runbookUrl ?? alert.grafanaUrl,
      namespace: alert.namespace,
    }));

export const controllerAttention = (notReady: string[]): AttentionItem[] =>
  notReady.map(name => ({
    key: `workload:${name}`,
    kind: 'workload' as const,
    severity: 'warning' as const,
    title: `Workload not fully ready: ${name}`,
    detail: 'Inspect pod readiness and events in Headlamp.',
    url: HEADLAMP_URL,
    namespace: name.split('/')[0],
  }));

export type UptimeSource = {
  checks: {
    name: string;
    entity: string;
    link: string;
    scope: string;
    state: string;
    httpStatus?: number;
    latencyMs?: number;
    observedAt: string;
  }[];
};

export const serviceAttention = (checks: UptimeSource['checks']): AttentionItem[] =>
  checks
    .filter(x => x.state !== 'up')
    .map(x => ({
      key: `health:${x.name}`,
      kind: 'service' as const,
      severity: x.state === 'down' ? ('critical' as const) : ('warning' as const),
      title: `${x.name}: ${x.state}`,
      detail: `${x.scope}. Inspect service health and connectivity; unknown is not healthy.`,
      url: x.link,
    }));

export const sortAttention = (items: AttentionItem[]) =>
  [...items].sort((a, b) => Number(b.severity === 'critical') - Number(a.severity === 'critical'));
