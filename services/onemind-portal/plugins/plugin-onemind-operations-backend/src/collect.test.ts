import { buildBanner, createCollector, CollectorDeps } from './collect';
import type { AmAlert } from './build';
import { liveQueries, PromSample } from './liveQueries';

const NOW = Date.parse('2026-10-22T12:00:00Z');
const ago = (ms: number) => new Date(NOW - ms).toISOString();
const MIN = 60000;
const HOUR = 3600000;

const sample = (metric: Record<string, string>, value: number): PromSample => ({ metric, value });

// Fixture shaped like the cluster on 2026-10-22 (read-only observation).
const cronjobs = [
  {
    metadata: {
      name: 'backstage-backup',
      namespace: 'backstage',
      annotations: { 'onemind.io/backup-state': 'superseded', 'onemind.io/superseded-by': 'backup-postgres-backstage' },
    },
    spec: { suspend: true, schedule: '15 7 * * *' },
    status: { lastSuccessfulTime: '2026-10-08T07:15:27Z' },
  },
  {
    metadata: {
      name: 'backup-clickhouse',
      namespace: 'clickhouse',
      annotations: { 'onemind.io/backup-state': 'superseded', 'onemind.io/superseded-by': 'VolSync daily' },
    },
    spec: { suspend: true, schedule: '0 3 * * *' },
    status: { lastSuccessfulTime: '2026-10-08T03:00:13Z' },
  },
  {
    metadata: {
      name: 'backup-offsite-copy',
      namespace: 'ops',
      annotations: { 'onemind.io/backup-state': 'pending-setup', 'onemind.io/pending-reason': 'Credentials not created' },
    },
    spec: { suspend: true, schedule: '0 17 * * *' },
  },
  {
    metadata: { name: 'backup-postgres-backstage', namespace: 'backstage' },
    spec: { schedule: '0 2 * * *' },
    status: { lastSuccessfulTime: ago(5 * HOUR), lastScheduleTime: ago(5 * HOUR) },
  },
  { metadata: { name: 'health-keeper', namespace: 'ops' }, spec: { schedule: '*/5 * * * *' } },
];
const argo = [
  { metadata: { name: 'platform-core' }, status: { health: { status: 'Healthy' }, sync: { status: 'Synced' } } },
  { metadata: { name: 'platform-matrix' }, status: { health: { status: 'Healthy' }, sync: { status: 'OutOfSync' } } },
  {
    metadata: { name: 'platform-openclaw' },
    status: { health: { status: 'Progressing', lastTransitionTime: ago(40 * MIN) }, sync: { status: 'Synced' } },
  },
];
const deployments = [
  { metadata: { name: 'a', namespace: 'x' }, spec: { replicas: 1 }, status: { readyReplicas: 1 } },
  { metadata: { name: 'b', namespace: 'y' }, spec: { replicas: 2 }, status: { readyReplicas: 1 } },
];
const statefulsets = [{ metadata: { name: 's', namespace: 'z' }, spec: { replicas: 1 }, status: { readyReplicas: 1 } }];
const pods = [
  {
    metadata: { name: 'p-1', namespace: 'y', creationTimestamp: ago(HOUR) },
    status: { phase: 'Running', containerStatuses: [{ restartCount: 3, state: { waiting: { reason: 'CrashLoopBackOff' } } }] },
  },
];
const nodes = [
  { metadata: { name: 'n1' }, status: { conditions: [{ type: 'Ready', status: 'True' }], allocatable: { cpu: '8', memory: '16Gi' } } },
];
const events = [
  {
    type: 'Warning', reason: 'BackOff', message: 'Back-off restarting', count: 1, lastTimestamp: ago(2 * MIN),
    metadata: { namespace: 'y' }, involvedObject: { kind: 'Pod', name: 'p-1', namespace: 'y' },
  },
];
const checks = [
  { name: 'Grafana', entity: 'component:default/grafana', link: 'https://g', scope: 'Application health endpoint', state: 'up', observedAt: ago(MIN) },
];
const alertmanager: AmAlert[] = [
  { fingerprint: 'w', labels: { alertname: 'Watchdog', severity: 'none' }, status: { state: 'active' } },
  {
    fingerprint: 'k', labels: { alertname: 'KubePodCrashLooping', severity: 'critical', namespace: 'y', pod: 'p-1' },
    annotations: { summary: 'crash looping' }, startsAt: ago(10 * MIN), status: { state: 'active' },
  },
];

function deps(over: Partial<CollectorDeps> = {}): CollectorDeps {
  const answers: Record<string, PromSample[]> = {
    [liveQueries.nodeCpu]: [sample({ nodename: 'n1' }, 12)],
    [liveQueries.nodeMemory]: [sample({ nodename: 'n1' }, 44)],
    [liveQueries.nodeDisk]: [sample({ nodename: 'n1' }, 16)],
    [liveQueries.pvcUsed]: [sample({ namespace: 'onetak', persistentvolumeclaim: 'valhalla-data' }, 99.97)],
    [liveQueries.pvcFilling]: [],
    [liveQueries.volsyncOutOfSync]: [sample({ obj_namespace: 'a', obj_name: 'v' }, 0)],
  };
  const lists: Record<string, unknown[]> = {
    '/apis/argoproj.io/v1alpha1/namespaces/argocd/applications': argo,
    '/apis/apps/v1/deployments': deployments,
    '/apis/apps/v1/statefulsets': statefulsets,
    '/api/v1/pods': pods,
    '/api/v1/nodes': nodes,
    '/api/v1/events?fieldSelector=type%3DWarning': events,
    '/apis/batch/v1/cronjobs': cronjobs,
  };
  return {
    now: () => NOW,
    list: (async (path: string) => {
      if (!(path in lists)) throw new Error('Cluster API returned 404');
      return lists[path];
    }) as never,
    prometheus: async query => answers[query],
    alertmanager: async () => alertmanager,
    uptime: async () => ({ observedAt: ago(MIN), checks, metrics: [] }),
    ...over,
  };
}

describe('live aggregation with every source up', () => {
  it('maps the three backup jobs to retired, retired and needs-setup', async () => {
    const live = await createCollector(deps()).live();
    const byName = Object.fromEntries(live.backups!.items.map(item => [item.name, item]));
    expect(byName['backstage/backstage-backup']).toMatchObject({ state: 'retired', attention: 'none', supersededBy: 'backup-postgres-backstage' });
    expect(byName['clickhouse/backup-clickhouse']).toMatchObject({ state: 'retired', attention: 'none' });
    expect(byName['ops/backup-offsite-copy']).toMatchObject({ state: 'needs-setup', attention: 'setup', pendingReason: 'Credentials not created' });
    expect(byName['backstage/backup-postgres-backstage']).toMatchObject({ health: 'ok' });
    expect(live.backups!.items.find(item => item.name.includes('health-keeper'))).toBeUndefined();
    expect(live.counts).toMatchObject({ backups: 0, backupsNeedSetup: 1, backupsRetired: 2 });
    expect(live.needsSetup).toHaveLength(1);
    expect(live.attention.some(item => item.key.startsWith('backup:'))).toBe(false);
  });

  it('shares one attention list across alerts, argo, workloads and storage', async () => {
    const live = await createCollector(deps()).live();
    expect(live.attention.map(item => item.key)).toEqual(
      expect.arrayContaining([
        'alert:k',
        'workload:y/b',
        'deployment:platform-matrix',
        'deployment:platform-openclaw',
        'storage:onetak/valhalla-data',
      ]),
    );
    expect(live.attention[0].severity).toBe('critical');
    expect(live.counts).toMatchObject({ alerts: 1, argo: 2, workloads: 1 });
    expect(live.argo!.items.map(i => i.reasons)).toEqual([['OutOfSync'], ['Progressing']]);
  });

  it('reports a critical banner with the reasons', async () => {
    const live = await createCollector(deps()).live();
    expect(live.health.level).toBe('critical');
    expect(live.health.headline).toMatch(/1 critical alert/);
    expect(live.nodes![0]).toMatchObject({ cpuPercent: 12, memoryPercent: 44, ready: true });
    expect(live.events!.items).toHaveLength(1);
    expect(Object.values(live.panels).every(panel => panel.status === 'ok')).toBe(true);
  });
});

describe('one source timing out', () => {
  it('marks only the Prometheus-backed panels as error and still renders the rest', async () => {
    const live = await createCollector(
      deps({
        prometheus: (_q, signal) =>
          new Promise((_, reject) => signal.addEventListener('abort', () => reject(new Error('aborted')))),
        sourceOptions: { timeoutMs: 30 },
      }),
    ).live();
    expect(live.panels.storage).toMatchObject({ status: 'error', observedAt: null });
    expect(live.panels.storage.problems).toEqual([{ source: 'prometheus', status: 'error', error: 'Timed out' }]);
    expect(live.panels.nodes.status).toBe('error');
    expect(live.panels.backups.status).toBe('error');
    expect(live.panels.alerts.status).toBe('ok');
    expect(live.panels.argo.status).toBe('ok');
    expect(live.panels.events.status).toBe('ok');
    expect(live.storage).toBeNull();
    expect(live.alerts).toHaveLength(1);
    // Backup jobs still come from the Kubernetes API.
    expect(live.backups!.items.length).toBeGreaterThan(0);
    expect(live.backups!.volsync).toBeNull();
    // Nodes still show readiness without usage numbers.
    expect(live.nodes![0]).toMatchObject({ ready: true, cpuPercent: null });
    expect(live.health.reasons.join(' ')).toMatch(/Prometheus/);
  });

  it('shows unknown, never zero, for a count whose source failed', async () => {
    const live = await createCollector(
      deps({ alertmanager: async () => { throw new Error('Upstream returned 503'); } }),
    ).live();
    expect(live.counts.alerts).toBe('unknown');
    expect(live.alerts).toBeNull();
    expect(live.panels.alerts).toMatchObject({ status: 'error' });
    expect(live.counts.argo).toBe(2);
  });

  it('serves stale data with its original observation time when a later read fails', async () => {
    let t = NOW;
    let fail = false;
    const collector = createCollector(
      deps({
        now: () => t,
        alertmanager: async () => {
          if (fail) throw new Error('Upstream returned 503');
          return alertmanager;
        },
      }),
    );
    await collector.live();
    fail = true;
    t += 30000;
    const live = await collector.live();
    expect(live.panels.alerts.status).toBe('stale');
    expect(live.panels.alerts.observedAt).toBe(new Date(NOW).toISOString());
    expect(live.alerts).toHaveLength(1);
    expect(live.counts.alerts).toBe(1);
  });
});

describe('/status compatibility', () => {
  it('keeps the previous fields and meanings', async () => {
    const status = await createCollector(deps()).status();
    expect(status).toMatchObject({
      source: 'do-nyc3-omos-cloud',
      partial: false,
      unavailableSources: [],
      unavailable: ['y/b'],
      totalWorkloads: 3,
    });
    expect(status.deployments).toEqual([
      { name: 'platform-matrix', health: 'Healthy', sync: 'OutOfSync' },
      { name: 'platform-openclaw', health: 'Progressing', sync: 'Synced' },
    ]);
    expect(status.serviceHealth?.checks).toEqual(checks);
    for (const backup of status.backups) {
      expect(Object.keys(backup)).toEqual(expect.arrayContaining(['name', 'lastSuccess', 'suspended']));
    }
  });

  it('reports the legacy source names that failed', async () => {
    const status = await createCollector(
      deps({ list: (async (path: string) => { if (path.includes('argoproj')) throw new Error('Cluster API returned 403'); return deps().list(path); }) as never }),
    ).status();
    expect(status.partial).toBe(true);
    expect(status.unavailableSources).toEqual(['Argo CD']);
    expect(status.counts.argo).toBe('unknown');
    expect(status.deployments).toEqual([]);
  });
});

describe('buildBanner', () => {
  const quiet = {
    alerts: { critical: 0, warning: 0 }, argo: { degraded: 0, other: 0 }, nodesNotReady: 0, nodePressure: 0,
    podsCrashing: 0, controllersNotReady: 0, storageCritical: 0, storageWarn: 0, backupsReview: 0,
    servicesDown: 0, unavailableSources: [] as string[],
  };
  it('is ok only when everything was observed and quiet', () => {
    expect(buildBanner(quiet).level).toBe('ok');
  });
  it('is unknown, not ok, when a source is missing', () => {
    expect(buildBanner({ ...quiet, alerts: null, unavailableSources: ['Alertmanager'] }).level).toBe('unknown');
  });
  it('states the reason in plain words', () => {
    const banner = buildBanner({ ...quiet, alerts: { critical: 2, warning: 0 }, argo: { degraded: 1, other: 0 } });
    expect(banner).toMatchObject({ level: 'critical', headline: '2 critical alerts, 1 app degraded' });
  });
});
