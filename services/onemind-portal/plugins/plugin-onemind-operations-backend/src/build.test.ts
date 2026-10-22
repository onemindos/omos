import {
  AmAlert,
  buildAlerts,
  buildEvents,
  buildNodes,
  classifyArgo,
  classifyPod,
  KubePod,
} from './build';
import { parseQuantity } from './quantity';

const NOW = Date.parse('2026-10-22T12:00:00Z');
const ago = (ms: number) => new Date(NOW - ms).toISOString();
const MIN = 60000;

describe('buildAlerts', () => {
  const raw: AmAlert[] = [
    { fingerprint: 'w', labels: { alertname: 'Watchdog', severity: 'none' }, status: { state: 'active' } },
    { fingerprint: 'i', labels: { alertname: 'InfoInhibitor', severity: 'none' }, status: { state: 'active' } },
    {
      fingerprint: 'c',
      labels: { alertname: 'KubePodCrashLooping', severity: 'critical', namespace: 'onetak', pod: 'p-1' },
      annotations: { summary: 'Pod is crash looping', runbook_url: 'https://runbooks.example/crash' },
      startsAt: '2026-10-22T11:00:00Z',
      status: { state: 'active' },
    },
    {
      fingerprint: 'w2',
      labels: { alertname: 'CPUThrottlingHigh', severity: 'warning', namespace: 'nats', deployment: 'nats-box' },
      annotations: { runbook_url: 'ftp://runbooks.example/not-http' },
      startsAt: '2026-10-22T11:30:00Z',
      status: { state: 'active' },
    },
    { fingerprint: 's', labels: { alertname: 'Silenced', severity: 'critical' }, status: { state: 'suppressed' } },
  ];
  it('drops heartbeat and suppressed alerts, orders critical first, and only links http runbooks', () => {
    const alerts = buildAlerts(raw);
    expect(alerts.map(a => a.name)).toEqual(['KubePodCrashLooping', 'CPUThrottlingHigh']);
    expect(alerts[0]).toMatchObject({ workload: 'p-1', runbookUrl: 'https://runbooks.example/crash' });
    expect(alerts[1]).toMatchObject({ workload: 'nats-box', runbookUrl: undefined });
    expect(alerts[1].grafanaUrl).toContain('/alerting/list?search=CPUThrottlingHigh');
  });
});

describe('classifyArgo', () => {
  const app = (health: string, sync: string, extra = {}) => ({
    metadata: { name: 'a' },
    status: { health: { status: health, lastTransitionTime: ago(2 * MIN) }, sync: { status: sync }, ...extra },
  });
  it('is quiet for Healthy and Synced', () => {
    expect(classifyArgo(app('Healthy', 'Synced'), NOW)).toBeUndefined();
  });
  it('flags OutOfSync as a warning and Degraded as critical', () => {
    expect(classifyArgo(app('Healthy', 'OutOfSync'), NOW)).toMatchObject({ reasons: ['OutOfSync'], severity: 'warning' });
    expect(classifyArgo(app('Degraded', 'Synced'), NOW)).toMatchObject({ reasons: ['Degraded'], severity: 'critical' });
  });
  it('flags Progressing only after ten minutes', () => {
    expect(classifyArgo(app('Progressing', 'Synced'), NOW)).toBeUndefined();
    const long = app('Progressing', 'Synced');
    long.status.health.lastTransitionTime = ago(11 * MIN);
    expect(classifyArgo(long, NOW)).toMatchObject({ reasons: ['Progressing'] });
  });
  it('flags a failed last operation', () => {
    expect(classifyArgo(app('Healthy', 'Synced', { operationState: { phase: 'Failed', message: 'x' } }), NOW)).toMatchObject({
      reasons: ['Sync failed'],
      severity: 'critical',
    });
  });
});

describe('classifyPod', () => {
  const pod = (status: KubePod['status'], created = ago(60 * MIN)): KubePod => ({
    metadata: { name: 'p', namespace: 'ns', creationTimestamp: created, labels: { 'app.kubernetes.io/instance': 'app-x' } },
    status,
  });
  it('ignores healthy, completed and terminating pods', () => {
    expect(classifyPod(pod({ phase: 'Running', conditions: [{ type: 'Ready', status: 'True' }] }), NOW)).toBeUndefined();
    expect(classifyPod(pod({ phase: 'Succeeded' }), NOW)).toBeUndefined();
  });
  it('flags crash loops as critical with the owning app', () => {
    const issue = classifyPod(
      pod({
        phase: 'Running',
        containerStatuses: [{ restartCount: 7, state: { waiting: { reason: 'CrashLoopBackOff' } }, lastState: { terminated: { reason: 'Error', finishedAt: ago(5 * MIN) } } }],
      }),
      NOW,
    );
    expect(issue).toMatchObject({ app: 'app-x', severity: 'critical', restarts: 7 });
    expect(issue?.reasons).toEqual(expect.arrayContaining(['CrashLoopBackOff', 'Restarted']));
  });
  it('flags a recent OOM kill but not an old restart', () => {
    const oom = classifyPod(
      pod({ phase: 'Running', conditions: [{ type: 'Ready', status: 'True' }], containerStatuses: [{ restartCount: 1, lastState: { terminated: { reason: 'OOMKilled', finishedAt: ago(10 * MIN) } } }] }),
      NOW,
    );
    expect(oom?.reasons).toEqual(['OOMKilled']);
    const old = classifyPod(
      pod({ phase: 'Running', conditions: [{ type: 'Ready', status: 'True' }], containerStatuses: [{ restartCount: 1, lastState: { terminated: { reason: 'OOMKilled', finishedAt: ago(5 * 3600000) } } }] }),
      NOW,
    );
    expect(old).toBeUndefined();
  });
  it('flags Pending after two minutes and NotReady after five', () => {
    expect(classifyPod(pod({ phase: 'Pending' }, ago(1 * MIN)), NOW)).toBeUndefined();
    expect(classifyPod(pod({ phase: 'Pending' }, ago(3 * MIN)), NOW)?.reasons).toEqual(['Pending']);
    const ready = (since: number) => pod({ phase: 'Running', conditions: [{ type: 'Ready', status: 'False', lastTransitionTime: ago(since) }] });
    expect(classifyPod(ready(2 * MIN), NOW)).toBeUndefined();
    expect(classifyPod(ready(6 * MIN), NOW)?.reasons).toEqual(['NotReady']);
  });
});

describe('buildNodes', () => {
  it('reports readiness, pressure, requests against allocatable, and joined usage', () => {
    const nodes = [
      {
        metadata: { name: 'n1' },
        status: {
          conditions: [
            { type: 'Ready', status: 'True' },
            { type: 'DiskPressure', status: 'True' },
            { type: 'MemoryPressure', status: 'False' },
          ],
          allocatable: { cpu: '8', memory: '16Gi' },
        },
      },
      { metadata: { name: 'n2' }, status: { conditions: [{ type: 'Ready', status: 'False' }] } },
    ];
    const pods: KubePod[] = [
      { metadata: { name: 'a', namespace: 'x' }, spec: { nodeName: 'n1', containers: [{ resources: { requests: { cpu: '2', memory: '4Gi' } } }, { resources: { requests: { cpu: '500m' } } }] }, status: { phase: 'Running' } },
      { metadata: { name: 'b', namespace: 'x' }, spec: { nodeName: 'n1', containers: [{ resources: { requests: { cpu: '4' } } }] }, status: { phase: 'Succeeded' } },
    ];
    const result = buildNodes(nodes, pods, { n1: { cpuPercent: 10, memoryPercent: 20, diskPercent: 30 } });
    expect(result[0]).toMatchObject({
      name: 'n1', ready: true, pressure: ['DiskPressure'], cpuPercent: 10, diskPercent: 30,
      cpuRequestedPercent: 31.3, memoryRequestedPercent: 25, pods: 1,
    });
    expect(result[1]).toMatchObject({ ready: false, cpuPercent: null, cpuRequestedPercent: null });
  });
});

describe('buildEvents', () => {
  const event = (over: object) => ({
    type: 'Warning',
    reason: 'BackOff',
    message: 'Back-off restarting failed container',
    count: 2,
    lastTimestamp: ago(5 * MIN),
    metadata: { namespace: 'ns' },
    involvedObject: { kind: 'Pod', name: 'p-1', namespace: 'ns' },
    ...over,
  });
  it('keeps the last 30 minutes of warnings, merges duplicates and sums counts', () => {
    const items = buildEvents(
      [
        event({}),
        event({ count: 3, lastTimestamp: ago(1 * MIN) }),
        event({ lastTimestamp: ago(45 * MIN) }),
        event({ type: 'Normal' }),
        event({ reason: 'Failed', involvedObject: { kind: 'Pod', name: 'p-2', namespace: 'ns' } }),
      ],
      NOW,
    );
    expect(items).toHaveLength(2);
    expect(items[0]).toMatchObject({ object: 'Pod/p-1', count: 5, reason: 'BackOff' });
  });
});

describe('parseQuantity', () => {
  it('reads cpu and memory quantities', () => {
    expect(parseQuantity('250m')).toBeCloseTo(0.25);
    expect(parseQuantity('2')).toBe(2);
    expect(parseQuantity('16Gi')).toBe(16 * 1024 ** 3);
    expect(parseQuantity('128M')).toBe(128e6);
    expect(parseQuantity('nope')).toBeUndefined();
    expect(parseQuantity(undefined)).toBeUndefined();
  });
});
