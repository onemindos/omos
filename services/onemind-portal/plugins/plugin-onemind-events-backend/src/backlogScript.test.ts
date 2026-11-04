/** @jest-environment node */
import { isDeploymentSuccessNotification } from '@internal/plugin-onemind-common-backend';

// eslint-disable-next-line @typescript-eslint/no-require-imports
const script = require('../scripts/mark-deployment-notifications-read.cjs') as {
  run: (options: Record<string, unknown>) => Promise<{ matched: number; kept: number; applied: number }>;
  isDeploymentSuccess: (payload: unknown) => boolean;
};

const success = (id: number) => ({ id: `s${id}`, payload: { title: `Deployment completed: app-${id}`, topic: 'deployments', severity: 'low' } });
const failure = (id: number) => ({ id: `f${id}`, payload: { title: `Deployment needs attention: app-${id}`, topic: 'deployments', severity: 'high' } });
const env = { PORTAL_URL: 'https://portal.example', PORTAL_TOKEN: 'fixture-token-not-real' };

// Fake notifications API holding unread items; records every request.
function fakeApi(unread: { id: string; payload: unknown }[]) {
  const requests: { method: string; url: string; body?: { ids: string[]; read: boolean } }[] = [];
  const read = new Set<string>();
  const fetch = async (url: string, init: { method?: string; body?: string; headers: Record<string, string> }) => {
    const method = init.method ?? 'GET';
    requests.push({ method, url, body: init.body ? JSON.parse(init.body) : undefined });
    expect(init.headers.Authorization).toBe('Bearer fixture-token-not-real');
    if (method === 'POST') {
      for (const id of JSON.parse(init.body as string).ids) read.add(id);
      return { ok: true, status: 200, json: async () => [] };
    }
    const q = new URL(url).searchParams;
    const items = unread.filter(item => !read.has(item.id)).slice(Number(q.get('offset')), Number(q.get('offset')) + Number(q.get('limit')));
    return { ok: true, status: 200, json: async () => ({ notifications: items, totalCount: unread.length }) };
  };
  return { fetch, requests, read };
}

describe('backlog cleanup script', () => {
  const backlog = () => [...Array.from({ length: 250 }, (_, i) => success(i)), failure(1), failure(2)];

  test('dry run lists and counts but never writes', async () => {
    const api = fakeApi(backlog());
    const out: string[] = [];
    const result = await script.run({ fetch: api.fetch, env, argv: [], out: (line: string) => out.push(line) });
    expect(result).toEqual({ matched: 250, kept: 2, applied: 0 });
    expect(api.requests.every(request => request.method === 'GET')).toBe(true);
    expect(out.join('\n')).toContain('Dry run only');
  });

  test('apply requires the matching expected count', async () => {
    const api = fakeApi(backlog());
    await expect(script.run({ fetch: api.fetch, env, argv: ['--apply'], out: () => {} })).rejects.toThrow(/Refusing to apply/);
    await expect(script.run({ fetch: api.fetch, env, argv: ['--apply', '--expect=7'], out: () => {} })).rejects.toThrow(/Refusing to apply/);
    expect(api.requests.some(request => request.method === 'POST')).toBe(false);
  });

  test('apply marks only the old success items read, in chunks, and never deletes or touches failures', async () => {
    const api = fakeApi(backlog());
    const result = await script.run({ fetch: api.fetch, env, argv: ['--apply', '--expect=250'], out: () => {} });
    expect(result.applied).toBe(250);
    expect(api.requests.map(request => request.method)).not.toContain('DELETE');
    const writes = api.requests.filter(request => request.method === 'POST');
    expect(writes).toHaveLength(3);
    expect(writes.every(request => request.url.endsWith('/api/notifications/update') && request.body?.read === true)).toBe(true);
    expect(writes.flatMap(request => request.body?.ids ?? []).some(id => id.startsWith('f'))).toBe(false);
    expect(api.read.size).toBe(250);
  });

  test('needs a token and refuses plain http hosts', async () => {
    await expect(script.run({ fetch: async () => ({}), env: { PORTAL_URL: 'https://portal.example' }, argv: [], out: () => {} })).rejects.toThrow(/PORTAL_TOKEN/);
    await expect(script.run({ fetch: async () => ({}), env: { ...env, PORTAL_URL: 'http://portal.example' }, argv: [], out: () => {} })).rejects.toThrow(/https/);
  });

  test('the script predicate matches the shared policy predicate', () => {
    const payloads = [
      { title: 'Deployment completed: a', topic: 'deployments', severity: 'low' },
      { title: 'Deployment completed: a', topic: 'deployments', severity: 'high' },
      { title: 'Deployment needs attention: a', topic: 'deployments', severity: 'low' },
      { title: 'Backup completed: a', topic: 'backups', severity: 'low' },
      { title: 5, topic: 'deployments', severity: 'low' },
      {},
    ];
    for (const payload of payloads) {
      expect(script.isDeploymentSuccess(payload)).toBe(isDeploymentSuccessNotification(payload));
    }
  });
});
