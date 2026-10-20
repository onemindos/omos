import {
  createScopedPrometheusApi,
  scopeForEntity,
  scopes,
} from './prometheusApi';
import type { DiscoveryApi, FetchApi } from '@backstage/core-plugin-api';

const observation = () => ({
  status: 'success',
  schemaVersion: 'onemind.prometheus-alerts.v1',
  scope: 'nats',
  observedAt: new Date().toISOString(),
  coverage: 'observed',
  data: { groups: [] },
});
const discovery = {
  getBaseUrl: jest.fn(async () => 'https://portal.test/api/onemind-services'),
} as unknown as DiscoveryApi;
test.each(Object.entries(scopes))(
  'only stable annotated %s entity opts in',
  (scope, name) => {
    const entity = {
      kind: 'Component',
      metadata: { name, annotations: { 'prometheus.io/service-name': scope } },
    };
    expect(scopeForEntity(entity)).toBe(scope);
    expect(scopeForEntity({ ...entity, kind: 'System' })).toBeUndefined();
    expect(
      scopeForEntity({
        ...entity,
        metadata: { ...entity.metadata, name: 'other' },
      }),
    ).toBeUndefined();
    expect(
      scopeForEntity({
        ...entity,
        metadata: { ...entity.metadata, namespace: 'other' },
      }),
    ).toBeUndefined();
  },
);
test('no arbitrary scope, URL, instance or PromQL routing', async () => {
  const fetch = jest.fn();
  const api = createScopedPrometheusApi(discovery, {
    fetch,
  } as unknown as FetchApi);
  await expect(api.getAlerts({ serviceName: '../../admin' })).rejects.toThrow(
    'Unsupported',
  );
  expect(fetch).not.toHaveBeenCalled();
  expect(api.getUiUrl({})).toBeUndefined();
});
test('authenticated fetch API uses fixed scope path and coalesces native/host reads', async () => {
  const fetch = jest.fn(async (_url: string) => ({
    ok: true,
    json: async () => observation(),
  }));
  const api = createScopedPrometheusApi(discovery, {
    fetch,
  } as unknown as FetchApi);
  await Promise.all([
    api.getAlerts({ serviceName: 'nats' }),
    api.getAlerts({ serviceName: 'nats' }),
  ]);
  expect(fetch).toHaveBeenCalledTimes(1);
  expect(fetch.mock.calls[0][0]).toBe(
    'https://portal.test/api/onemind-services/prometheus/nats',
  );
});
test.each([
  { ...observation(), scope: 'backstage' },
  { ...observation(), status: 'error' },
  { ...observation(), observedAt: '2020-01-01' },
  { ...observation(), data: null },
])('rejects stale/wrong-scope/invalid response', async body => {
  const api = createScopedPrometheusApi(discovery, {
    fetch: async () => ({ ok: true, json: async () => body }),
  } as unknown as FetchApi);
  await expect(api.getAlerts({ serviceName: 'nats' })).rejects.toThrow(
    'Invalid or stale',
  );
});
