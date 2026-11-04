/** @jest-environment node */
import express from 'express';
import type { Server } from 'node:http';
import plugin from './plugin';

jest.mock('@backstage/backend-plugin-api', () => ({
  createBackendPlugin: (value: unknown) => value,
  coreServices: { httpRouter: 'router', httpAuth: 'auth', userInfo: 'userInfo' },
}));
jest.mock('@backstage/plugin-notifications-node', () => ({ notificationService: 'notifications' }));
jest.mock('@backstage/plugin-events-node', () => ({ eventsServiceRef: 'events' }));

const send = jest.fn(async () => {});
const publish = jest.fn(async () => {});
let server: Server;
let base: string;
const iso = () => new Date().toISOString();

type Principal = { type: 'service'; subject: string } | { type: 'user'; admin: boolean };
const principals: Record<string, Principal> = {
  'Bearer ingestor': { type: 'service', subject: 'onemind-event-ingestor' },
  'Bearer other-service': { type: 'service', subject: 'onemind-readonly-ai' },
  'Bearer admin': { type: 'user', admin: true },
  'Bearer user': { type: 'user', admin: false },
};

beforeAll(async () => {
  let initialize: (deps: unknown) => Promise<void> = async () => {};
  (plugin as unknown as { register: (r: unknown) => void }).register({
    registerInit: (options: { init: typeof initialize }) => {
      initialize = options.init;
    },
  });
  const app = express();
  await initialize({
    router: app,
    notifications: { send },
    events: { publish },
    httpAuth: {
      credentials: async (req: express.Request) => {
        const principal = principals[req.get('Authorization') ?? ''];
        if (!principal) throw Object.assign(new Error('Unauthorized'), { status: 401 });
        return { principal };
      },
    },
    userInfo: {
      getUserInfo: async (credentials: { principal: { admin: boolean } }) => ({
        userEntityRef: 'user:default/operator',
        ownershipEntityRefs: credentials.principal.admin ? ['group:default/onemind-admins'] : ['user:default/operator'],
      }),
    },
  });
  app.use((error: { status?: number }, _req: express.Request, res: express.Response, _next: express.NextFunction) => {
    res.sendStatus(error.status ?? 500);
  });
  await new Promise<void>(resolve => {
    server = app.listen(0, '127.0.0.1', resolve);
  });
  base = `http://127.0.0.1:${(server.address() as { port: number }).port}/events`;
});
afterAll(() => new Promise<void>(resolve => server.close(() => resolve())));
beforeEach(() => {
  send.mockClear();
  publish.mockClear();
});
const post = (event: Record<string, unknown>, auth = 'Bearer ingestor') =>
  fetch(base, { method: 'POST', headers: { 'content-type': 'application/json', Authorization: auth }, body: JSON.stringify(event) });
const event = (kind: string, over: Record<string, unknown> = {}) => ({ id: 'op-1', kind, resource: 'platform-backstage', occurredAt: iso(), ...over });

describe('activity event bridge severity policy', () => {
  test('deployment.succeeded is accepted but is never a notification', async () => {
    const response = await post(event('deployment.succeeded'));
    expect(response.status).toBe(202);
    expect(await response.json()).toMatchObject({ delivery: 'timeline', notified: false });
    expect(send).not.toHaveBeenCalled();
    expect(publish).toHaveBeenCalledWith(expect.objectContaining({ topic: 'onemind.activity' }));
  });

  test('backup.succeeded is timeline-only too', async () => {
    expect((await post(event('backup.succeeded', { resource: 'backstage/backup-postgres-backstage' }))).status).toBe(202);
    expect(send).not.toHaveBeenCalled();
  });

  test('deployment.failed is one high item and repeats share the app scope', async () => {
    await post(event('deployment.failed', { id: 'op-1' }));
    await post(event('deployment.failed', { id: 'op-2' }));
    expect(send).toHaveBeenCalledTimes(2);
    const calls = send.mock.calls as unknown as [{ payload: { severity: string; scope: string } }][];
    expect(calls[0][0].payload).toMatchObject({ severity: 'high', scope: 'argo:platform-backstage' });
    expect(calls[1][0].payload.scope).toBe(calls[0][0].payload.scope);
  });

  test('new actionable kinds are accepted: out-of-sync and suspended backup are high', async () => {
    await post(event('deployment.out-of-sync'));
    await post(event('backup.suspended', { resource: 'backstage/backup-postgres-backstage' }));
    const calls = send.mock.calls as unknown as [{ payload: { severity: string } }][];
    expect(calls.map(call => call[0].payload.severity)).toEqual(['high', 'high']);
  });
});

describe('activity event bridge access (negative)', () => {
  test.each([
    ['anonymous', ''],
    ['unknown token', 'Bearer nope'],
  ])('%s is rejected', async (_name, auth) => {
    expect((await post(event('deployment.failed'), auth)).status).toBe(401);
    expect(send).not.toHaveBeenCalled();
  });

  test.each(['Bearer other-service', 'Bearer user'])('%s is forbidden', async auth => {
    expect((await post(event('deployment.failed'), auth)).status).toBe(403);
    expect(send).not.toHaveBeenCalled();
  });

  test('administrators and the dedicated ingestor are allowed', async () => {
    expect((await post(event('deployment.failed'), 'Bearer admin')).status).toBe(202);
    expect((await post(event('deployment.failed'))).status).toBe(202);
  });

  test('unknown kinds, prototype keys and extra fields are rejected', async () => {
    expect((await post(event('deployment.exploded'))).status).toBe(400);
    expect((await post(event('toString'))).status).toBe(400);
    expect((await post({ ...event('deployment.failed'), recipients: 'all' })).status).toBe(400);
    expect(send).not.toHaveBeenCalled();
  });
});
