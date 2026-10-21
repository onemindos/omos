/** @jest-environment node */
import express from 'express';
import path from 'node:path';
import type { Server } from 'node:http';
import { ConfigReader } from '@backstage/config';
import type { JsonObject } from '@backstage/types';
import { loadConfig } from '@backstage/config-loader';
import { HttpPostIngressEventPublisher } from '@backstage/plugin-events-backend';
import { DefaultEventsService } from '@backstage/plugin-events-node';
import { createAlertmanagerValidator, WEBHOOK_TOKEN_KEY } from './alertmanager';
import { routeAlerts } from './alertRouting';
import { createStore } from './store';
import { destroyMemoryDbs, memoryDb } from './testDb';

const TOKEN = 'synthetic-test-token-0123456789abcdef';
const logger = {
  info: () => {},
  warn: () => {},
  error: () => {},
  debug: () => {},
  child() {
    return this;
  },
};
const payload = (status: 'firing' | 'resolved') => ({
  version: '4',
  status,
  alerts: [
    {
      status,
      labels: { alertname: 'KubePodCrashLooping', severity: 'critical', pod: 'backstage-1' },
      annotations: { summary: 'Pod is crash looping' },
      startsAt: '2026-10-21T10:00:00Z',
      endsAt: status === 'resolved' ? '2026-10-21T10:30:00Z' : '0001-01-01T00:00:00Z',
      fingerprint: 'abc123def4567890',
    },
  ],
});

// Real events-backend HTTP ingress, local events service and the real validator and routing.
async function start(config: JsonObject) {
  const events = DefaultEventsService.create({ logger: logger as never }).forPlugin('onemind');
  const store = await createStore(memoryDb());
  const notify = jest.fn(async () => {});
  const signal = jest.fn(async (_message: unknown) => {});
  const signalTimes: number[] = [];
  const reader = new ConfigReader(config);
  await events.subscribe({
    id: 'test-routing',
    topics: ['alertmanager'],
    onEvent: async ({ eventPayload }) => {
      await routeAlerts({
        body: eventPayload,
        store,
        sinks: {
          notify,
          signal: async message => {
            signalTimes.push(Date.now());
            await signal(message);
          },
        },
      });
    },
  });
  const ingress = HttpPostIngressEventPublisher.fromConfig({
    config: reader,
    events,
    logger: logger as never,
    ingresses: {
      alertmanager: { validator: createAlertmanagerValidator(() => reader.getOptionalString(WEBHOOK_TOKEN_KEY)) },
    },
  });
  const router = express.Router();
  ingress.bind(router);
  const app = express();
  app.use('/api/events', router);
  const server: Server = await new Promise(resolve => {
    const s = app.listen(0, '127.0.0.1', () => resolve(s));
  });
  const base = `http://127.0.0.1:${(server.address() as { port: number }).port}/api/events/http/alertmanager`;
  const post = (body: unknown, headers: Record<string, string> = {}) =>
    fetch(base, { method: 'POST', headers: { 'content-type': 'application/json', ...headers }, body: typeof body === 'string' ? body : JSON.stringify(body) });
  return { server, post, notify, signal, signalTimes };
}

const servers: Server[] = [];
afterAll(async () => {
  await Promise.all(servers.map(s => new Promise(resolve => s.close(resolve))));
  await destroyMemoryDbs();
});
const settle = () => new Promise(resolve => setTimeout(resolve, 50));
const authed = { authorization: `Bearer ${TOKEN}` };

describe('Alertmanager webhook through the native events ingress', () => {
  test('firing then resolved creates and resolves one item; the signal arrives within 5 seconds', async () => {
    const t = await start({ onemind: { alerts: { webhookToken: TOKEN } } });
    servers.push(t.server);
    const sentAt = Date.now();
    expect((await t.post(payload('firing'), authed)).status).toBe(202);
    await settle();
    expect(t.notify).toHaveBeenCalledTimes(1);
    expect(t.signal).toHaveBeenCalledTimes(1);
    expect(t.signalTimes[0] - sentAt).toBeLessThan(5000);

    expect((await t.post(payload('resolved'), authed)).status).toBe(202);
    await settle();
    expect(t.notify).toHaveBeenCalledTimes(2);
    expect(t.notify).toHaveBeenLastCalledWith(expect.objectContaining({ severity: 'low', scope: 'alert:abc123def4567890' }));
    expect(t.signal).toHaveBeenLastCalledWith(expect.objectContaining({ action: 'alert.resolved' }));
  });

  test('negative: no token, wrong token and a Basic header are rejected and publish nothing', async () => {
    const t = await start({ onemind: { alerts: { webhookToken: TOKEN } } });
    servers.push(t.server);
    expect((await t.post(payload('firing'))).status).toBe(401);
    expect((await t.post(payload('firing'), { authorization: 'Bearer wrong-token-wrong-token-wrong' })).status).toBe(401);
    expect((await t.post(payload('firing'), { authorization: `Basic ${TOKEN}` })).status).toBe(401);
    await settle();
    expect(t.notify).not.toHaveBeenCalled();
    expect(t.signal).not.toHaveBeenCalled();
  });

  test('negative: valid token but malformed body is a 400 and publishes nothing', async () => {
    const t = await start({ onemind: { alerts: { webhookToken: TOKEN } } });
    servers.push(t.server);
    expect((await t.post({ alerts: 'x' }, authed)).status).toBe(400);
    expect((await t.post({ version: '4', alerts: [{ status: 'pending', labels: {}, startsAt: 'x' }] }, authed)).status).toBe(400);
    await settle();
    expect(t.notify).not.toHaveBeenCalled();
  });

  test('token not configured: every request is refused, even with a plausible header', async () => {
    const t = await start({});
    servers.push(t.server);
    expect((await t.post(payload('firing'), authed)).status).toBe(503);
    expect((await t.post(payload('firing'))).status).toBe(503);
    await settle();
    expect(t.notify).not.toHaveBeenCalled();
  });

  test('a configured but too-short token also keeps the route disabled', async () => {
    const t = await start({ onemind: { alerts: { webhookToken: 'short' } } });
    servers.push(t.server);
    expect((await t.post(payload('firing'), { authorization: 'Bearer short' })).status).toBe(503);
  });
});

describe('optional secret wiring', () => {
  const configRoot = path.resolve(__dirname, '../../..');
  const load = async (env: Record<string, string>) => {
    const result = await loadConfig({
      configRoot,
      configTargets: [{ path: path.join(configRoot, 'app-config.yaml') }],
      experimentalEnvFunc: async name => env[name],
    });
    return ConfigReader.fromConfigs(result.appConfigs);
  };

  test('app-config loads without the env var and the key is simply absent', async () => {
    const config = await load({});
    expect(config.getOptionalString(WEBHOOK_TOKEN_KEY)).toBeUndefined();
    // Native settings expose every topic.
    expect(config.getConfigArray('notifications.defaultSettings.channels')).toHaveLength(1);
  });

  test('the env var, when present, reaches the same key', async () => {
    const config = await load({ ONEMIND_ALERT_WEBHOOK_TOKEN: TOKEN });
    expect(config.getOptionalString(WEBHOOK_TOKEN_KEY)).toBe(TOKEN);
  });
});
