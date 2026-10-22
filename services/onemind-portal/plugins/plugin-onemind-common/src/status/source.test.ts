import { alertsFromPrometheus, createOperationsStatusSource, snapshotFromOperationsStatus } from './source';

const status = {
  observedAt: '2026-10-22T12:00:00Z',
  unavailable: ['nats/nack'],
  deployments: [{ name: 'platform-x', health: 'Degraded', sync: 'Synced' }],
  serviceHealth: { observedAt: '2026-10-22T12:00:05Z', checks: [{ entity: 'component:default/grafana', name: 'Grafana', state: 'up' }] },
};

test('the existing status route maps to the normalized snapshot', () => {
  const s = snapshotFromOperationsStatus(status);
  expect(s.argo).toEqual({ state: 'observed', observedAt: status.observedAt, value: status.deployments });
  expect(s.notReady).toEqual({ state: 'observed', observedAt: status.observedAt, value: ['nats/nack'] });
  expect(s.uptime).toMatchObject({ state: 'observed', observedAt: '2026-10-22T12:00:05Z' });
});

test('a source the backend could not read is unavailable, not an empty (green) list', () => {
  const s = snapshotFromOperationsStatus({ ...status, unavailableSources: ['Argo CD', 'StatefulSets'] });
  expect(s.argo?.state).toBe('unavailable');
  expect(s.notReady?.state).toBe('unavailable');
  expect(snapshotFromOperationsStatus({ observedAt: status.observedAt }).uptime?.state).toBe('unavailable');
});

test('Prometheus observations count firing, critical and pending alerts', () => {
  const r = alertsFromPrometheus({
    observedAt: status.observedAt,
    coverage: 'observed',
    data: { groups: [{ rules: [{ alerts: [{ state: 'firing', labels: { severity: 'critical' } }, { state: 'firing', labels: {} }, { state: 'pending' }] }] }] },
  });
  expect(r).toEqual({ state: 'observed', observedAt: status.observedAt, value: { firing: 2, critical: 1, pending: 1 } });
  expect(alertsFromPrometheus({ observedAt: status.observedAt, coverage: 'partial' }).state).toBe('unavailable');
});

test('the default source reads status once plus the three scopes; a failed scope is unavailable only for that scope', async () => {
  const urls: string[] = [];
  const fetchApi: any = {
    fetch: async (url: string) => {
      urls.push(url);
      if (url.endsWith('/status')) return { ok: true, json: async () => status };
      if (url.endsWith('/postgres')) return { ok: false, status: 503 };
      return { ok: true, json: async () => ({ observedAt: status.observedAt, coverage: 'observed', data: { groups: [] } }) };
    },
  };
  const discovery: any = { getBaseUrl: async (plugin: string) => `http://x/api/${plugin}` };
  const snapshot = await createOperationsStatusSource(discovery, fetchApi).read();
  expect(urls.sort()).toEqual([
    'http://x/api/onemind-services/prometheus/backstage',
    'http://x/api/onemind-services/prometheus/nats',
    'http://x/api/onemind-services/prometheus/postgres',
    'http://x/api/onemind/status',
  ]);
  expect(snapshot.alerts?.backstage.state).toBe('observed');
  expect(snapshot.alerts?.postgres.state).toBe('unavailable');
});

test('if the status route fails, the whole read fails so the store keeps its last snapshot', async () => {
  const fetchApi: any = { fetch: async () => ({ ok: false, status: 500 }) };
  const discovery: any = { getBaseUrl: async () => 'http://x' };
  await expect(createOperationsStatusSource(discovery, fetchApi).read()).rejects.toThrow();
});
