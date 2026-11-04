/** @jest-environment node */
import express from 'express';
import type { Server } from 'node:http';
import { createTimelineRouter } from './timelineRoute';
import { createStore } from './store';
import { destroyMemoryDbs, memoryDb } from './testDb';

let server: Server;
let base: string;

beforeAll(async () => {
  const store = await createStore(memoryDb());
  await store.timeline.record({ id: 'a', kind: 'deployment.succeeded', resource: 'app', occurredAt: '2026-10-21T10:00:00Z', source: 'argocd' });
  const app = express();
  app.use(
    createTimelineRouter({
      store,
      httpAuth: {
        credentials: async (req: express.Request, options: { allow: string[] }) => {
          expect(options.allow).toEqual(['user']);
          const header = req.get('Authorization');
          if (header === 'Bearer admin' || header === 'Bearer user') return { principal: { type: 'user', admin: header === 'Bearer admin' } };
          throw Object.assign(new Error('Not authorized'), { status: 401 });
        },
      } as never,
      userInfo: {
        getUserInfo: async (credentials: { principal: { admin: boolean } }) => ({
          userEntityRef: 'user:default/operator',
          ownershipEntityRefs: credentials.principal.admin ? ['group:default/onemind-admins'] : ['user:default/operator'],
        }),
      } as never,
    }),
  );
  app.use((error: { status?: number }, _req: express.Request, res: express.Response, _next: express.NextFunction) => {
    res.sendStatus(error.status ?? 500);
  });
  await new Promise<void>(resolve => {
    server = app.listen(0, '127.0.0.1', resolve);
  });
  base = `http://127.0.0.1:${(server.address() as { port: number }).port}/activity/timeline`;
});
afterAll(async () => {
  await new Promise<void>(resolve => server.close(() => resolve()));
  await destroyMemoryDbs();
});

test('administrators read recorded successes', async () => {
  const response = await fetch(base, { headers: { Authorization: 'Bearer admin' } });
  expect(response.status).toBe(200);
  expect((await response.json()).entries).toHaveLength(1);
});

test('negative: anonymous is refused and a non-admin user is denied', async () => {
  expect((await fetch(base)).status).toBe(401);
  expect((await fetch(base, { headers: { Authorization: 'Bearer user' } })).status).toBe(403);
});
