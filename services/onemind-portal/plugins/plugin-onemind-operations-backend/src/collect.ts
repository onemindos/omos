// The Live Operations aggregator. One collector owns every source; /live returns the
// full picture and /status returns the older, smaller shape from the same data.
import {
  alertAttention,
  AmAlert,
  argoAttention,
  ARGOCD_URL,
  ArgoApp,
  backupAttention,
  buildAlerts,
  buildArgo,
  buildBackups,
  buildEvents,
  buildNodes,
  buildWorkloads,
  controllerAttention,
  GRAFANA_URL,
  HEADLAMP_URL,
  KubeController,
  KubeEvent,
  KubeNode,
  KubePod,
  legacyDeployments,
  notReadyControllers,
  pvcAttention,
  serviceAttention,
  sortAttention,
  UptimeSource,
} from './build';
import { summarizeBackups, CronJobLike } from './backupState';
import {
  liveQueries,
  parseNodeUsage,
  parsePvcs,
  parseVolsync,
  PromSample,
} from './liveQueries';
import {
  createSource,
  Source,
  SourceOptions,
  SourceResult,
  worstOf,
} from './source';
import type {
  AttentionItem,
  HealthBanner,
  LiveResponse,
  PanelId,
  PanelState,
  StatusCounts,
  StripCount,
  UptimeItem,
} from './types';

export const CLUSTER = 'do-nyc3-omos-cloud';

export type CollectorDeps = {
  list: <T>(path: string) => Promise<T[]>;
  prometheus: (query: string, signal: AbortSignal) => Promise<PromSample[]>;
  alertmanager: (signal: AbortSignal) => Promise<AmAlert[]>;
  uptime: () => Promise<UptimeSource & Record<string, unknown>>;
  now?: () => number;
  sourceOptions?: SourceOptions;
};

const plural = (n: number, one: string, many = `${one}s`) => `${n} ${n === 1 ? one : many}`;

function stripCount(unavailable: boolean, value: number): StripCount {
  return unavailable ? 'unknown' : value;
}

export function buildBanner(input: {
  alerts: { critical: number; warning: number } | null;
  argo: { degraded: number; other: number } | null;
  nodesNotReady: number;
  nodePressure: number;
  podsCrashing: number;
  controllersNotReady: number;
  storageCritical: number;
  storageWarn: number;
  backupsReview: number;
  servicesDown: number;
  unavailableSources: string[];
}): HealthBanner {
  const critical: string[] = [];
  const warn: string[] = [];
  if (input.alerts?.critical) critical.push(`${plural(input.alerts.critical, 'critical alert')}`);
  if (input.argo?.degraded) critical.push(`${plural(input.argo.degraded, 'app')} degraded`);
  if (input.nodesNotReady) critical.push(`${plural(input.nodesNotReady, 'node')} not ready`);
  if (input.storageCritical) critical.push(`${plural(input.storageCritical, 'volume')} almost full`);
  if (input.podsCrashing) critical.push(`${plural(input.podsCrashing, 'pod')} crash-looping`);
  if (input.servicesDown) critical.push(`${plural(input.servicesDown, 'service check')} down`);
  if (input.alerts?.warning) warn.push(`${plural(input.alerts.warning, 'warning alert')}`);
  if (input.argo?.other) warn.push(`${plural(input.argo.other, 'app')} out of sync or progressing`);
  if (input.controllersNotReady) warn.push(`${plural(input.controllersNotReady, 'workload')} not ready`);
  if (input.nodePressure) warn.push(`${plural(input.nodePressure, 'node')} under pressure`);
  if (input.storageWarn) warn.push(`${plural(input.storageWarn, 'volume')} over 80% full`);
  if (input.backupsReview) warn.push(`${plural(input.backupsReview, 'backup')} to review`);
  const missing = input.unavailableSources.length
    ? [`${plural(input.unavailableSources.length, 'source')} unavailable (${input.unavailableSources.join(', ')})`]
    : [];
  if (critical.length) {
    return { level: 'critical', headline: [...critical, ...warn].slice(0, 3).join(', '), reasons: [...critical, ...warn, ...missing] };
  }
  if (warn.length) {
    return { level: 'warn', headline: warn.slice(0, 3).join(', '), reasons: [...warn, ...missing] };
  }
  if (missing.length) {
    return { level: 'unknown', headline: `Cannot confirm health: ${missing[0]}`, reasons: missing };
  }
  return { level: 'ok', headline: 'Everything observed is healthy', reasons: [] };
}

export function createCollector(deps: CollectorDeps) {
  const now = deps.now ?? Date.now;
  const opts = { now, ...deps.sourceOptions };
  const k8s = <T>(id: Parameters<typeof createSource>[0], path: string): Source<T[]> =>
    createSource<T[]>(id, () => deps.list<T>(path), opts);

  const sources = {
    alertmanager: createSource<AmAlert[]>('alertmanager', signal => deps.alertmanager(signal), { timeoutMs: 5000, ...opts }),
    prometheus: createSource('prometheus', async signal => {
      const names = ['nodeCpu', 'nodeMemory', 'nodeDisk', 'pvcUsed', 'pvcFilling', 'volsyncOutOfSync'] as const;
      const answers = await Promise.all(names.map(name => deps.prometheus(liveQueries[name], signal)));
      const [cpu, memory, disk, used, filling, volsync] = answers;
      return {
        usage: parseNodeUsage({ cpu, memory, disk }),
        pvcs: parsePvcs(used, filling),
        volsync: parseVolsync(volsync),
      };
    }, { timeoutMs: 6000, ...opts }),
    argocd: k8s<ArgoApp>('argocd', '/apis/argoproj.io/v1alpha1/namespaces/argocd/applications'),
    deployments: k8s<KubeController>('deployments', '/apis/apps/v1/deployments'),
    statefulsets: k8s<KubeController>('statefulsets', '/apis/apps/v1/statefulsets'),
    pods: k8s<KubePod>('pods', '/api/v1/pods'),
    nodes: k8s<KubeNode>('nodes', '/api/v1/nodes'),
    events: k8s<KubeEvent>('events', '/api/v1/events?fieldSelector=type%3DWarning'),
    cronjobs: k8s<CronJobLike>('cronjobs', '/apis/batch/v1/cronjobs'),
    uptime: createSource('uptime', () => deps.uptime(), { timeoutMs: 8000, ...opts }),
  };

  async function readAll() {
    const [alertmanager, prometheus, argocd, deployments, statefulsets, pods, nodes, events, cronjobs, uptime] =
      await Promise.all(Object.values(sources).map(source => (source as Source<unknown>).read()));
    return {
      alertmanager: alertmanager as SourceResult<AmAlert[]>,
      prometheus: prometheus as SourceResult<{
        usage: ReturnType<typeof parseNodeUsage>;
        pvcs: ReturnType<typeof parsePvcs>;
        volsync: ReturnType<typeof parseVolsync>;
      }>,
      argocd: argocd as SourceResult<ArgoApp[]>,
      deployments: deployments as SourceResult<KubeController[]>,
      statefulsets: statefulsets as SourceResult<KubeController[]>,
      pods: pods as SourceResult<KubePod[]>,
      nodes: nodes as SourceResult<KubeNode[]>,
      events: events as SourceResult<KubeEvent[]>,
      cronjobs: cronjobs as SourceResult<CronJobLike[]>,
      uptime: uptime as SourceResult<UptimeSource & Record<string, unknown>>,
    };
  }
  type Read = Awaited<ReturnType<typeof readAll>>;

  function assemble(r: Read): LiveResponse {
    const at = now();
    const down = (...items: SourceResult<unknown>[]) => items.some(item => item.status === 'error');
    const panel = (...items: SourceResult<unknown>[]): PanelState => ({
      ...worstOf(items),
      problems: items
        .filter(item => item.status !== 'ok')
        .map(item => ({ source: item.id, status: item.status, error: item.error })),
    });
    const panels: Record<PanelId, PanelState> = {
      alerts: panel(r.alertmanager),
      argo: panel(r.argocd),
      workloads: panel(r.pods, r.deployments, r.statefulsets),
      nodes: panel(r.nodes, r.prometheus),
      storage: panel(r.prometheus),
      backups: panel(r.cronjobs, r.prometheus),
      uptime: panel(r.uptime),
      events: panel(r.events),
    };

    const alerts = r.alertmanager.data ? buildAlerts(r.alertmanager.data) : null;
    const argoItems = r.argocd.data ? buildArgo(r.argocd.data, at) : null;
    const workloads =
      r.deployments.data && r.statefulsets.data
        ? buildWorkloads(r.pods.data ?? [], r.deployments.data, r.statefulsets.data, at)
        : null;
    const nodes = r.nodes.data ? buildNodes(r.nodes.data, r.pods.data ?? [], r.prometheus.data?.usage) : null;
    const backupItems = r.cronjobs.data ? buildBackups(r.cronjobs.data, at) : null;
    const uptime: UptimeItem[] | null = r.uptime.data ? r.uptime.data.checks : null;
    const pvcs = r.prometheus.data?.pvcs ?? null;

    const backupLists = backupItems ? backupAttention(backupItems) : { attention: [], needsSetup: [] };
    const attention: AttentionItem[] = sortAttention([
      ...(alerts ? alertAttention(alerts) : []),
      ...(uptime ? serviceAttention(uptime) : []),
      ...(workloads ? controllerAttention(workloads.notReady) : []),
      ...(argoItems ? argoAttention(argoItems) : []),
      ...backupLists.attention,
      ...(pvcs ? pvcAttention(pvcs) : []),
    ]);

    const summary = backupItems ? summarizeBackups(backupItems) : undefined;
    const actionableAlerts = alerts?.filter(alert => alert.severity !== 'info') ?? [];
    const servicesDown = uptime?.filter(check => check.state !== 'up').length ?? 0;
    const counts: StatusCounts = {
      alerts: stripCount(down(r.alertmanager, r.uptime), actionableAlerts.length + servicesDown),
      argo: stripCount(down(r.argocd), argoItems?.length ?? 0),
      workloads: stripCount(down(r.deployments, r.statefulsets), workloads?.notReady.length ?? 0),
      backups: stripCount(down(r.cronjobs), summary?.review ?? 0),
      backupsNeedSetup: stripCount(down(r.cronjobs), summary?.setup ?? 0),
      backupsRetired: stripCount(down(r.cronjobs), summary?.retired ?? 0),
    };

    const labels: Record<string, string> = {
      alertmanager: 'Alertmanager',
      prometheus: 'Prometheus',
      argocd: 'Argo CD',
      deployments: 'Deployments',
      statefulsets: 'StatefulSets',
      pods: 'Pods',
      nodes: 'Nodes',
      events: 'Events',
      cronjobs: 'Backup schedules',
      uptime: 'Service checks',
    };
    const unavailableSources = Object.values(r)
      .filter(item => item.status === 'error')
      .map(item => labels[item.id]);

    const health = buildBanner({
      alerts: alerts
        ? {
            critical: alerts.filter(a => a.severity === 'critical').length,
            warning: alerts.filter(a => a.severity === 'warning').length,
          }
        : null,
      argo: argoItems
        ? {
            degraded: argoItems.filter(a => a.severity === 'critical').length,
            other: argoItems.filter(a => a.severity !== 'critical').length,
          }
        : null,
      nodesNotReady: nodes?.filter(n => !n.ready).length ?? 0,
      nodePressure: nodes?.filter(n => n.ready && n.pressure.length).length ?? 0,
      podsCrashing: workloads?.issues.filter(i => i.reasons.includes('CrashLoopBackOff')).length ?? 0,
      controllersNotReady: workloads?.notReady.length ?? 0,
      storageCritical: pvcs?.filter(p => p.severity === 'critical').length ?? 0,
      storageWarn: pvcs?.filter(p => p.severity === 'warn').length ?? 0,
      backupsReview: summary?.review ?? 0,
      servicesDown: uptime?.filter(c => c.state === 'down').length ?? 0,
      unavailableSources,
    });

    return {
      observedAt: new Date(at).toISOString(),
      cluster: CLUSTER,
      health,
      counts,
      panels,
      alerts,
      argo: argoItems && r.argocd.data ? { items: argoItems, total: r.argocd.data.length } : null,
      workloads,
      nodes,
      storage: pvcs ? { pvcs, volsync: r.prometheus.data?.volsync ?? null } : null,
      backups: backupItems ? { items: backupItems, volsync: r.prometheus.data?.volsync ?? null } : null,
      uptime,
      events: r.events.data ? { items: buildEvents(r.events.data, at), windowMinutes: 30 } : null,
      attention,
      needsSetup: backupLists.needsSetup,
      links: { argocd: ARGOCD_URL, headlamp: HEADLAMP_URL, grafana: GRAFANA_URL },
    };
  }

  return {
    async live(): Promise<LiveResponse> {
      return assemble(await readAll());
    },

    /**
     * The `/api/onemind/status` body. Every field the route returned before is still here
     * with the same meaning; `attention`, `needsSetup`, `counts` and the backup state fields are additions.
     */
    async status() {
      const r = await readAll();
      const live = assemble(r);
      const legacy = [r.deployments, r.statefulsets, r.argocd, r.cronjobs];
      const names = ['Deployments', 'StatefulSets', 'Argo CD', 'Backup schedules'];
      const workloads = [...(r.deployments.data ?? []), ...(r.statefulsets.data ?? [])];
      const times = legacy.map(item => item.observedAt).filter((t): t is string => !!t);
      return {
        serviceHealth: r.uptime.data,
        observedAt: times.length ? times.reduce((a, b) => (a < b ? a : b)) : live.observedAt,
        source: CLUSTER,
        partial: legacy.some(item => item.status === 'error'),
        unavailableSources: legacy.flatMap((item, i) => (item.status === 'error' ? [names[i]] : [])),
        unavailable: notReadyControllers(r.deployments.data ?? [], r.statefulsets.data ?? []),
        deployments: legacyDeployments(r.argocd.data ?? []),
        backups: live.backups?.items ?? [],
        totalWorkloads: workloads.length,
        counts: live.counts,
        attention: live.attention,
        needsSetup: live.needsSetup,
        health: live.health,
      };
    },
  };
}
