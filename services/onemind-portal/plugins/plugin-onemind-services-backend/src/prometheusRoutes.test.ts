/** @jest-environment node */
import express from 'express';
import type { Server } from 'node:http';
import plugin from './plugin';
import { readPrometheusAlerts } from './prometheusAlerts';

jest.mock('@backstage/backend-plugin-api', () => ({
  createBackendPlugin: (value: unknown) => value,
  coreServices: { httpRouter: 'router', httpAuth: 'auth' },
}));
jest.mock('./prometheusAlerts', () => ({
  alertScopes: ['backstage', 'postgres', 'nats'],
  readPrometheusAlerts: jest.fn(),
}));
const read = readPrometheusAlerts as jest.Mock;
let server: Server;
let base: string;
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
    httpAuth: {
      credentials: async (
        req: express.Request,
        options: { allow: string[] },
      ) => {
        expect(options.allow).toEqual(['user']);
        if (req.get('Authorization') !== 'Bearer user-fixture') {
          const error = Object.assign(new Error('Not authorized'), {
            status: 403,
          });
          throw error;
        }
      },
    },
  });
  await new Promise<void>(resolve => {
    server = app.listen(0, '127.0.0.1', resolve);
  });
  const address = server.address() as { port: number };
  base = `http://127.0.0.1:${address.port}`;
});
afterAll(async () => {
  await new Promise<void>((resolve, reject) =>
    server.close(error => (error ? reject(error) : resolve())),
  );
});
beforeEach(() => {
  read.mockReset();
  read.mockResolvedValue({ status: 'success', scope: 'nats' });
});
test.each([undefined, 'Bearer machine-fixture'])(
  'anonymous/machine credentials cannot reach observations',
  async credential => {
    const response = await fetch(`${base}/prometheus/nats`, {
      headers: credential ? { Authorization: credential } : {},
    });
    expect(response.status).toBe(403);
    expect(read).not.toHaveBeenCalled();
  },
);
test('real-user fixed scope accepted, queries/arbitrary scope rejected, write method absent', async () => {
  const headers = { Authorization: 'Bearer user-fixture' };
  expect((await fetch(`${base}/prometheus/nats`, { headers })).status).toBe(
    200,
  );
  expect(read).toHaveBeenCalledWith('nats');
  read.mockClear();
  expect(
    (await fetch(`${base}/prometheus/nats?query=up`, { headers })).status,
  ).toBe(400);
  expect((await fetch(`${base}/prometheus/other`, { headers })).status).toBe(
    400,
  );
  expect(
    (await fetch(`${base}/prometheus/nats`, { method: 'POST', headers }))
      .status,
  ).toBe(404);
  expect(read).not.toHaveBeenCalled();
});
test('source outage returns generic 503, never upstream details or cached all-clear', async () => {
  read.mockRejectedValue(new Error('private upstream details'));
  const response = await fetch(`${base}/prometheus/nats`, {
    headers: { Authorization: 'Bearer user-fixture' },
  });
  expect(response.status).toBe(503);
  expect(await response.text()).not.toContain('private upstream');
});
