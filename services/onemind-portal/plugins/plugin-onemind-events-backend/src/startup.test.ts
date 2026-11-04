/** @jest-environment node */
import net from 'node:net';
import { createBackend } from '@backstage/backend-defaults';
import { coreServices, createBackendPlugin, createServiceFactory } from '@backstage/backend-plugin-api';
import { ConfigReader } from '@backstage/config';
import type { JsonObject } from '@backstage/types';
import eventsBackend from '@backstage/plugin-events-backend';
import notificationsBackend from '@backstage/plugin-notifications-backend';
import signalsBackend from '@backstage/plugin-signals-backend';
import feature from './index';

// Boots a real backend (events, signals, notifications, our modules) on SQLite. The `onemind`
// plugin that hosts our modules in production is replaced by an empty stand-in.
const freePort = () =>
  new Promise<number>(resolve => {
    const probe = net.createServer().listen(0, '127.0.0.1', () => {
      const { port } = probe.address() as net.AddressInfo;
      probe.close(() => resolve(port));
    });
  });

async function boot(extra: JsonObject) {
  const port = await freePort();
  const config: JsonObject = {
    app: { baseUrl: `http://127.0.0.1:${port}` },
    backend: {
      baseUrl: `http://127.0.0.1:${port}`,
      listen: { host: '127.0.0.1', port },
      database: { client: 'better-sqlite3', connection: ':memory:' },
      auth: { dangerouslyDisableDefaultAuthPolicy: false },
    },
    events: { useEventBus: 'never', http: { topics: [] } },
    ...extra,
  };
  const backend = createBackend();
  backend.add(createServiceFactory({ service: coreServices.rootConfig, deps: {}, factory: () => new ConfigReader(config) }));
  backend.add(eventsBackend);
  backend.add(signalsBackend);
  backend.add(notificationsBackend);
  backend.add(createBackendPlugin({ pluginId: 'onemind', register(reg) { reg.registerInit({ deps: {}, async init() {} }); } }));
  backend.add(feature);
  await backend.start();
  return { backend, url: `http://127.0.0.1:${port}/api/events/http/alertmanager` };
}

const body = JSON.stringify({
  version: '4',
  status: 'firing',
  alerts: [{ status: 'firing', labels: { alertname: 'Synthetic', severity: 'warning' }, annotations: {}, startsAt: '2026-10-21T10:00:00Z', fingerprint: 'feedface00000001' }],
});
const post = (url: string, headers: Record<string, string> = {}) =>
  fetch(url, { method: 'POST', headers: { 'content-type': 'application/json', ...headers }, body });

describe('portal startup with and without the webhook secret', () => {
  jest.setTimeout(60_000);

  test('starts normally with no token configured, and the route rejects everything', async () => {
    const { backend, url } = await boot({});
    try {
      expect((await post(url)).status).toBe(503);
      expect((await post(url, { authorization: `Bearer ${'a'.repeat(32)}` })).status).toBe(503);
    } finally {
      await backend.stop();
    }
  });

  test('with a token configured, the real route enforces it', async () => {
    const token = 'startup-test-token-0123456789abcdef';
    const { backend, url } = await boot({ onemind: { alerts: { webhookToken: token } } });
    try {
      expect((await post(url)).status).toBe(401);
      expect((await post(url, { authorization: `Bearer ${'b'.repeat(32)}` })).status).toBe(401);
      expect((await post(url, { authorization: `Bearer ${token}` })).status).toBe(202);
    } finally {
      await backend.stop();
    }
  });
});
