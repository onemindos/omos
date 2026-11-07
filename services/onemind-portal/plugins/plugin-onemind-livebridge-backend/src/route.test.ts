/** @jest-environment node */
import express from 'express';
import type { Server } from 'node:http';
import type { LiveSnapshot } from './contract';
import { createSnapshotRouter } from './route';
import { summarizeSnapshot } from './summary';

const snapshot: LiveSnapshot = {
  origin: 'o',
  seq: 3,
  generatedAt: '2026-11-07T10:00:00Z',
  status: {
    nats: { state: 'connected', rate: 12.5 },
    kubernetes: { state: 'connected' },
    argo: { state: 'disabled' },
  },
  tracks: [
    {
      uid: 'a',
      type: 'a-f-G-U-C',
      lat: 1,
      lon: 2,
      time: 't',
      stale: 's',
      subject: 'tak.cot.out',
    },
    {
      uid: 'b',
      type: 'a-n-A-C-F',
      lat: 1,
      lon: 2,
      time: 't',
      stale: 's',
      subject: 'tak.cot.out',
    },
    {
      uid: 'c',
      type: 'a-n-S',
      lat: 1,
      lon: 2,
      time: 't',
      stale: 's',
      subject: 'tak.cot.out',
    },
  ],
  events: Array.from({ length: 30 }, (_, i) => ({
    id: `e${i}`,
    subject: 's',
    kind: 'k',
    ts: 't',
    bytes: 1,
    summary: String(i),
  })),
  cluster: [
    {
      namespace: 'a',
      total: 2,
      ready: 2,
      pending: 0,
      failed: 0,
      restarts: 0,
      health: 'ok',
    },
    {
      namespace: 'b',
      total: 2,
      ready: 1,
      pending: 1,
      failed: 0,
      restarts: 0,
      health: 'warn',
    },
  ],
  argo: [
    { name: 'x', sync: 'Synced', health: 'Healthy' },
    { name: 'y', sync: 'OutOfSync', health: 'Healthy' },
  ],
  warnings: [],
};

let server: Server;
let base: string;
beforeAll(async () => {
  const app = express();
  app.use(
    createSnapshotRouter({
      snapshot: () => snapshot,
      httpAuth: {
        credentials: async (
          req: express.Request,
          options: { allow: string[] },
        ) => {
          expect(options.allow).toEqual(['user']);
          if (req.get('Authorization') === 'Bearer user')
            return { principal: { type: 'user' } };
          throw Object.assign(new Error('Not authorized'), { status: 401 });
        },
      } as never,
    }),
  );
  app.use(
    (
      error: { status?: number },
      _req: express.Request,
      res: express.Response,
      _next: express.NextFunction,
    ) => {
      res.sendStatus(error.status ?? 500);
    },
  );
  await new Promise<void>(resolve => {
    server = app.listen(0, '127.0.0.1', resolve);
  });
  base = `http://127.0.0.1:${
    (server.address() as { port: number }).port
  }/snapshot`;
});
afterAll(async () => {
  await new Promise<void>(resolve => server.close(() => resolve()));
});

test('a signed-in user gets the snapshot', async () => {
  const res = await fetch(base, { headers: { Authorization: 'Bearer user' } });
  expect(res.status).toBe(200);
  expect((await res.json()).seq).toBe(3);
});

test('negative: anonymous callers are refused', async () => {
  expect((await fetch(base)).status).toBe(401);
});

test('summary digest counts health, argo drift, affiliations and the last 20 events', () => {
  const s = summarizeSnapshot(snapshot);
  expect(s.namespaceHealth).toEqual({ ok: 1, warn: 1, down: 0 });
  expect(s.unhealthyNamespaces.map(n => n.namespace)).toEqual(['b']);
  expect(s.argoNotSyncedHealthy.map(a => a.name)).toEqual(['y']);
  expect(s.trackCountsByAffiliation).toEqual({ friendly: 1, neutral: 2 });
  expect(s.busMessagesPerSecond).toBe(12.5);
  expect(s.lastEvents).toHaveLength(20);
});
