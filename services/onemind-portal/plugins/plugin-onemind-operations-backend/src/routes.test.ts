/** @jest-environment node */
import express from 'express';
import type { AddressInfo } from 'node:net';
import type { Server } from 'node:http';
import { AuthenticationError, NotAllowedError } from '@backstage/errors';
import { createLiveRoutes } from './routes';

// Mirrors the real HttpAuthService contract: no token is a 401; a service or
// plugin token is refused when the route only allows users.
const credentials = jest.fn(async (req: express.Request, options?: { allow?: string[] }) => {
  const header = req.header('authorization');
  if (!header) throw new AuthenticationError('Missing credentials');
  const principal = header === 'Bearer user' ? 'user' : 'service';
  if (!options?.allow?.includes(principal)) throw new NotAllowedError(`Principal ${principal} not allowed`);
  return { principal: { type: principal } } as never;
});

const collector = {
  live: jest.fn(async () => ({ observedAt: 'now' })),
  status: jest.fn(async () => ({ observedAt: 'now' })),
};

let server: Server;
let base: string;
beforeAll(async () => {
  const app = express();
  app.use(createLiveRoutes({ httpAuth: { credentials } as never, collector: collector as never }));
  app.use(((error, _req, res, _next) => {
    const status = ({ AuthenticationError: 401, NotAllowedError: 403 } as Record<string, number>)[error.name] ?? 500;
    res.status(status).json({ error: { name: error.name } });
  }) as express.ErrorRequestHandler);
  server = app.listen(0);
  base = `http://127.0.0.1:${(server.address() as AddressInfo).port}`;
});
afterAll(() => new Promise(resolve => server.close(resolve)));
beforeEach(() => jest.clearAllMocks());

describe.each(['/live', '/status'])('GET %s', path => {
  it('denies an unauthenticated caller and reads nothing', async () => {
    const response = await fetch(base + path);
    expect(response.status).toBe(401);
    expect(collector.live).not.toHaveBeenCalled();
    expect(collector.status).not.toHaveBeenCalled();
  });

  it('denies a service credential: only signed-in users may read cluster state', async () => {
    const response = await fetch(base + path, { headers: { authorization: 'Bearer service' } });
    expect(response.status).toBe(403);
    expect(collector.live).not.toHaveBeenCalled();
    expect(collector.status).not.toHaveBeenCalled();
  });

  it('asks the auth service for user principals only', async () => {
    await fetch(base + path, { headers: { authorization: 'Bearer user' } });
    expect(credentials).toHaveBeenCalledWith(expect.anything(), { allow: ['user'] });
  });

  it('returns the aggregate to a signed-in user without caching', async () => {
    const response = await fetch(base + path, { headers: { authorization: 'Bearer user' } });
    expect(response.status).toBe(200);
    expect(await response.json()).toEqual({ observedAt: 'now' });
  });
});

it('exposes no write methods', async () => {
  for (const method of ['POST', 'PUT', 'PATCH', 'DELETE']) {
    const response = await fetch(`${base}/live`, { method, headers: { authorization: 'Bearer user' } });
    expect(response.status).toBe(404);
  }
});

it('never lets a browser or proxy cache the live answer', async () => {
  const response = await fetch(`${base}/live`, { headers: { authorization: 'Bearer user' } });
  expect(response.headers.get('cache-control')).toBe('no-store');
});
