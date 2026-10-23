import http from 'node:http';
import type { AddressInfo } from 'node:net';
import { getEvidenceHttp, observeStacEvidence } from './stacEvidence';
import type { EvidenceDependencies } from './stacEvidence';

const deployment = (name: string, namespace: string) => ({
  metadata: { name, namespace, generation: 2 },
  spec: {
    replicas: 1,
    template: {
      spec: {
        containers: [{ name: 'app', image: 'example/image@sha256:reviewed' }],
      },
    },
  },
  status: { observedGeneration: 2, readyReplicas: 1 },
});
function dependencies() {
  const get = jest.fn(async (url: string) => {
    if (url.includes('/api/v1/query'))
      return {
        status: 200,
        body: JSON.stringify({
          status: 'success',
          data: {
            resultType: 'vector',
            result: [{ value: [Date.now() / 1000, '1'] }],
          },
        }),
      };
    if (url.endsWith('/collections'))
      return {
        status: 200,
        body: JSON.stringify({ collections: [], links: [] }),
      };
    if (url.includes('/search?'))
      return {
        status: 200,
        body: JSON.stringify({ type: 'FeatureCollection', features: [] }),
      };
    return { status: 200, body: '{}' };
  });
  const list = jest.fn(async (path: string): Promise<unknown[]> => {
    if (path.includes('/namespaces/onetak/'))
      return [
        deployment('stac-api', 'onetak'),
        deployment('onetak-geo-db', 'onetak'),
      ];
    return ['platform-onetak'].map(name => ({
      metadata: { name },
      status: {
        health: { status: 'Healthy' },
        sync: { status: 'Synced', revision: 'reviewed-revision' },
      },
    }));
  });
  const deps: EvidenceDependencies = {
    get,
    list: async <T>(path: string) => (await list(path)) as T[],
  };
  return { get, list, deps };
}
test('joins only the fixed STAC API and database, with distinct workload/API/data observations', async () => {
  const { deps, get, list } = dependencies();
  const result = await observeStacEvidence(deps);
  expect(result.partial).toBe(false);
  expect(result.workloads.map(w => w.role)).toEqual(['api', 'database']);
  expect(result.workloads.map(w => w.state)).toEqual(['ready', 'ready']);
  expect(result.probes.collections.httpStatus).toBe(200);
  expect(result.data.collectionsReturned).toBe(0);
  expect(result.data.latestItemDatetime).toBeNull();
  expect(result.data.freshness).toContain('No ingestion SLA');
  expect(list).toHaveBeenCalledTimes(2);
  expect(
    get.mock.calls.filter(([url]) => url.endsWith('/collections')),
  ).toHaveLength(1);
  expect(get.mock.calls.every(([url]) => url.startsWith('http://'))).toBe(true);
  expect(JSON.stringify(result)).not.toMatch(/password|Authorization|Bearer/);
});
test('a working landing page and ready pods cannot mask collection/search 500s', async () => {
  const { deps, get } = dependencies();
  const initial = deps.get;
  deps.get = async url =>
    /\/collections$|\/search\?/.test(url)
      ? { status: 500, body: 'private upstream error' }
      : initial(url);
  const result = await observeStacEvidence(deps);
  expect(result.probes.apiLanding.httpStatus).toBe(200);
  expect(result.probes.collections.httpStatus).toBe(500);
  expect(result.partial).toBe(true);
  expect(result.data.state).toBe('unknown');
  expect(result.data.collectionsReturned).toBeNull();
  expect(JSON.stringify(result)).not.toContain('private upstream error');
  expect(get).toHaveBeenCalled();
});
test('cluster observation failure is unknown, never zero-ready or all healthy', async () => {
  const { deps } = dependencies();
  deps.list = async () => {
    throw new Error('Source unavailable');
  };
  const result = await observeStacEvidence(deps);
  expect(result.partial).toBe(true);
  expect(
    result.workloads.every(w => w.state === 'unknown' && w.ready === null),
  ).toBe(true);
  expect(result.argo.state).toBe('partial-or-unknown');
});
test('duplicate workloads with the same name are not guessed', async () => {
  const { deps } = dependencies();
  const initial = deps.list;
  deps.list = async <T>(path: string) =>
    path.includes('/namespaces/onetak/')
      ? ([
          deployment('stac-api', 'onetak'),
          deployment('stac-api', 'onetak'),
          deployment('onetak-geo-db', 'onetak'),
        ] as T[])
      : initial<T>(path);
  const result = await observeStacEvidence(deps);
  expect(result.workloads[0].state).toBe('unknown');
  expect(result.workloads[1].state).toBe('ready');
});
test('unobserved generation and intentional zero replicas are not reported ready', async () => {
  const { deps } = dependencies();
  const initial = deps.list;
  deps.list = async <T>(path: string) => {
    if (!path.includes('/namespaces/onetak/')) return initial<T>(path);
    const a = deployment('stac-api', 'onetak');
    a.status.observedGeneration = 1;
    const b = deployment('onetak-geo-db', 'onetak');
    b.spec.replicas = 0;
    b.status.readyReplicas = 0;
    return [a, b] as T[];
  };
  const result = await observeStacEvidence(deps);
  expect(result.workloads[0].state).toBe('reconciling');
  expect(result.workloads[1].state).toBe('scaled-to-zero');
  expect(result.partial).toBe(true);
});
test.each([
  ['Progressing', 'Synced'],
  ['Healthy', 'OutOfSync'],
])(
  'nonaccepted GitOps state %s/%s marks the joined observation incomplete',
  async (health, sync) => {
    const { deps } = dependencies();
    const initial = deps.list;
    deps.list = async <T>(path: string) => {
      if (!path.endsWith('/applications')) return initial<T>(path);
      return ['platform-onetak'].map(name => ({
        metadata: { name },
        status: { health: { status: health }, sync: { status: sync } },
      })) as T[];
    };
    expect((await observeStacEvidence(deps)).partial).toBe(true);
  },
);
test.each(['NaN', '-1', 'Infinity', ''])(
  'invalid metric %s is unknown',
  async value => {
    const { deps } = dependencies();
    const initial = deps.get;
    deps.get = async url =>
      url.includes('/api/v1/query')
        ? {
            status: 200,
            body: JSON.stringify({
              status: 'success',
              data: {
                resultType: 'vector',
                result: [{ value: [Date.now() / 1000, value] }],
              },
            }),
          }
        : initial(url);
    const result = await observeStacEvidence(deps);
    expect(result.workloads[0].cpu?.value).toBeNull();
  },
);
test('old metric evaluation timestamps remain stale and do not become zero', async () => {
  const { deps } = dependencies();
  const initial = deps.get;
  deps.get = async url =>
    url.includes('/api/v1/query')
      ? {
          status: 200,
          body: JSON.stringify({
            status: 'success',
            data: {
              resultType: 'vector',
              result: [{ value: [(Date.now() - 600_000) / 1000, '1'] }],
            },
          }),
        }
      : initial(url);
  const result = await observeStacEvidence(deps);
  expect(result.workloads[0].cpu).toMatchObject({
    state: 'stale',
    value: null,
  });
});
test('bounded collection pages and acquisition dates are not mistaken for ingestion freshness', async () => {
  const { deps } = dependencies();
  const initial = deps.get;
  deps.get = async url => {
    if (url.endsWith('/collections'))
      return {
        status: 200,
        body: JSON.stringify({ collections: [{}], links: [{ rel: 'next' }] }),
      };
    if (url.includes('/search?'))
      return {
        status: 200,
        body: JSON.stringify({
          type: 'FeatureCollection',
          features: [{ properties: { datetime: '2026-01-01T00:00:00Z' } }],
        }),
      };
    return initial(url);
  };
  const result = await observeStacEvidence(deps);
  expect(result.data.collectionsTruncated).toBe(true);
  expect(result.data.latestItemDatetime).toBe('2026-01-01T00:00:00Z');
  expect(result.data.freshness).toContain('not ingestion time');
});
test('HTTP observer does not follow redirects or send authorization', async () => {
  const paths: string[] = [];
  const headers: (string | undefined)[] = [];
  const server = http.createServer((req, res) => {
    paths.push(req.url!);
    headers.push(req.headers.authorization);
    res.writeHead(302, { Location: '/private-target' });
    res.end();
  });
  await new Promise<void>(resolve => server.listen(0, '127.0.0.1', resolve));
  try {
    const r = await getEvidenceHttp(
      `http://127.0.0.1:${(server.address() as AddressInfo).port}/start`,
    );
    expect(r.status).toBe(302);
    expect(paths).toEqual(['/start']);
    expect(headers).toEqual([undefined]);
  } finally {
    await new Promise<void>(resolve => server.close(() => resolve()));
  }
});
test('oversized HTTP source is rejected', async () => {
  const server = http.createServer((_req, res) => {
    res.end(Buffer.alloc(2_000_001));
  });
  await new Promise<void>(resolve => server.listen(0, '127.0.0.1', resolve));
  try {
    await expect(
      getEvidenceHttp(
        `http://127.0.0.1:${(server.address() as AddressInfo).port}/`,
      ),
    ).rejects.toThrow('response limit');
  } finally {
    await new Promise<void>(resolve => server.close(() => resolve()));
  }
});
