import {
  combineStatus, summarize, targetForEntity, AppTarget, StatusSnapshot, STATUS_STALE_MS,
} from './model';

const now = Date.parse('2026-10-22T12:00:00Z');
const at = (secondsAgo: number) => new Date(now - secondsAgo * 1000).toISOString();
const target = (extra: Partial<AppTarget> = {}): AppTarget => ({
  argoApps: ['platform-x'], workloads: ['x/web', 'x/db'], refs: ['component:default/x'], alertScopes: [], ...extra,
});
const healthy = (observed = at(10)): StatusSnapshot => ({
  observedAt: observed,
  argo: { state: 'observed', observedAt: observed, value: [] },
  notReady: { state: 'observed', observedAt: observed, value: [] },
  uptime: { state: 'observed', observedAt: observed, value: [] },
});

test('nothing read yet is unknown, never green', () => {
  expect(combineStatus(target(), undefined, now)).toEqual({ status: 'unknown', signals: [] });
});

test('every connected signal ok gives ok', () => {
  const r = combineStatus(target(), healthy(), now);
  expect(r.status).toBe('ok');
  expect(r.signals.map(s => s.name)).toEqual(['argo', 'workloads']);
});

test('an app with no connected signal is unknown rather than ok', () => {
  expect(combineStatus(target({ argoApps: [], workloads: [] }), healthy(), now).status).toBe('unknown');
});

test('a failed source makes its signal unknown and blocks green', () => {
  const snapshot: StatusSnapshot = { ...healthy(), argo: { state: 'unavailable', reason: 'Argo CD' } };
  const r = combineStatus(target(), snapshot, now);
  expect(r.status).toBe('unknown');
  expect(r.signals.find(s => s.name === 'argo')?.status).toBe('unknown');
});

test('Argo: degraded or missing is down, out of sync or progressing is warn, unrecognised is unknown', () => {
  const run = (health: string, sync = 'Synced') =>
    combineStatus(target({ workloads: [] }), { ...healthy(), argo: { state: 'observed', observedAt: at(5), value: [{ name: 'platform-x', health, sync }] } }, now).status;
  expect(run('Degraded')).toBe('down');
  expect(run('Missing', 'OutOfSync')).toBe('down');
  expect(run('Healthy', 'OutOfSync')).toBe('warn');
  expect(run('Progressing')).toBe('warn');
  expect(run('Unknown')).toBe('unknown');
});

test('Argo problems of other apps are ignored', () => {
  const snapshot = { ...healthy(), argo: { state: 'observed' as const, observedAt: at(5), value: [{ name: 'platform-other', health: 'Degraded', sync: 'Synced' }] } };
  expect(combineStatus(target(), snapshot, now).status).toBe('ok');
});

test('workloads: some below desired is warn, all of them is down, others are ignored', () => {
  const run = (notReady: string[]) => combineStatus(target({ argoApps: [] }), { ...healthy(), notReady: { state: 'observed', observedAt: at(5), value: notReady } }, now);
  expect(run(['x/web']).status).toBe('warn');
  expect(run(['x/web', 'x/db']).status).toBe('down');
  expect(run(['y/other']).status).toBe('ok');
  expect(run(['x/web']).signals[0].detail).toMatch(/1 of 2/);
});

test('uptime: only checks that belong to the app count; sign-in required still counts as reachable', () => {
  const run = (checks: { entity?: string; name: string; state: string }[]) =>
    combineStatus(target({ argoApps: [], workloads: [] }), { ...healthy(), uptime: { state: 'observed', observedAt: at(5), value: checks } }, now);
  expect(run([{ entity: 'component:default/x', name: 'X', state: 'up' }]).status).toBe('ok');
  expect(run([{ entity: 'component:default/x', name: 'X', state: 'authentication-required' }]).status).toBe('ok');
  expect(run([{ entity: 'component:default/x', name: 'X', state: 'down' }]).status).toBe('down');
  expect(run([{ entity: 'component:default/x', name: 'X', state: 'unknown' }]).status).toBe('unknown');
  expect(run([{ entity: 'component:default/other', name: 'O', state: 'down' }]).signals).toEqual([]);
});

test('alerts: critical firing is down, other firing or pending is warn, failed read is unknown', () => {
  const run = (value: any) =>
    combineStatus(target({ argoApps: [], workloads: [], alertScopes: ['nats'] }), { ...healthy(), alerts: { nats: value } }, now);
  expect(run({ state: 'observed', observedAt: at(5), value: { firing: 1, critical: 1, pending: 0 } }).status).toBe('down');
  expect(run({ state: 'observed', observedAt: at(5), value: { firing: 1, critical: 0, pending: 0 } }).status).toBe('warn');
  expect(run({ state: 'observed', observedAt: at(5), value: { firing: 0, critical: 0, pending: 2 } }).status).toBe('warn');
  expect(run({ state: 'observed', observedAt: at(5), value: { firing: 0, critical: 0, pending: 0 } }).status).toBe('ok');
  expect(run({ state: 'unavailable' }).status).toBe('unknown');
});

test('worst signal wins; a known problem outranks missing data', () => {
  const snapshot: StatusSnapshot = {
    ...healthy(),
    argo: { state: 'unavailable' },
    notReady: { state: 'observed', observedAt: at(5), value: ['x/web'] },
  };
  expect(combineStatus(target(), snapshot, now).status).toBe('warn');
});

test('a stale snapshot can show problems but never health', () => {
  const old = healthy(at(STATUS_STALE_MS / 1000 + 5));
  expect(combineStatus(target(), old, now).status).toBe('unknown');
  const broken: StatusSnapshot = { ...old, notReady: { state: 'observed', observedAt: old.observedAt, value: ['x/web', 'x/db'] } };
  expect(combineStatus(target(), broken, now).status).toBe('down');
  expect(combineStatus(target(), { ...healthy(), observedAt: 'garbage' }, now).status).toBe('unknown');
});

test('target comes from the generated annotations, including workload entities and hasPart', () => {
  const app = targetForEntity({
    kind: 'Component',
    metadata: { name: 'Nats', annotations: { 'argocd/app-name': 'platform-nats', 'onemind.io/argocd-apps': 'platform-nats,edge-nack', 'onemind.io/kubernetes-workloads': 'nats/nats,nats/nack', 'onemind.io/prometheus-scopes': 'nats' } },
    relations: [{ type: 'hasPart', targetRef: 'component:default/K8s-nats' }, { type: 'dependsOn', targetRef: 'component:default/zzz' }],
  });
  expect(app.argoApps).toEqual(['platform-nats', 'edge-nack']);
  expect(app.workloads).toEqual(['nats/nats', 'nats/nack']);
  expect(app.alertScopes).toEqual(['nats']);
  expect(app.refs).toEqual(['component:default/nats', 'component:default/k8s-nats']);
  const workload = targetForEntity({ kind: 'Component', metadata: { name: 'w', annotations: { 'backstage.io/kubernetes-namespace': 'backstage', 'onemind.io/kubernetes-workload': 'backstage', 'prometheus.io/service-name': 'backstage' } } });
  expect(workload.workloads).toEqual(['backstage/backstage']);
  expect(workload.alertScopes).toEqual(['backstage']);
});

test('summaries count each state', () => {
  expect(summarize([{ status: 'ok', signals: [] }, { status: 'ok', signals: [] }, { status: 'down', signals: [] }])).toEqual({ ok: 2, warn: 0, down: 1, unknown: 0 });
});
