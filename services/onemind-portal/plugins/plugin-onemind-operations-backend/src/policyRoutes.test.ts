/** @jest-environment node */
import express from 'express';
import type { AddressInfo } from 'node:net';
import type { Server } from 'node:http';
import { AuthenticationError, NotAllowedError } from '@backstage/errors';
import {
  createPolicyRoutes,
  createPolicySource,
  POLICY_REPORTS_PATH,
  POLICY_TTL_MS,
} from './policyRoutes';
import type { KubePolicyReport } from './policyReports';

// Mirrors the real HttpAuthService contract: no token is a 401; a service or plugin
// token is refused when the route only allows users.
const credentials = jest.fn(
  async (req: express.Request, options?: { allow?: string[] }) => {
    const header = req.header('authorization');
    if (!header) throw new AuthenticationError('Missing credentials');
    const principal = header === 'Bearer user' ? 'user' : 'service';
    if (!options?.allow?.includes(principal))
      throw new NotAllowedError(`Principal ${principal} not allowed`);
    return { principal: { type: principal } } as never;
  },
);

const sample: KubePolicyReport[] = [
  {
    metadata: { name: 'a', namespace: 'app' },
    scope: { kind: 'Deployment', name: 'web', namespace: 'app' },
    results: [{ result: 'fail', policy: 'p', rule: 'r', message: 'm' }],
  },
];

let clock = 1_000_000;
const list = jest.fn(async (_path: string) => sample);
const source = createPolicySource(list, () => clock);

let server: Server;
let base: string;
const get = (headers: Record<string, string> = {}) =>
  fetch(`${base}/policy-reports`, { headers });
const asUser = { authorization: 'Bearer user' };

beforeAll(async () => {
  const app = express();
  app.use(createPolicyRoutes({ httpAuth: { credentials } as never, source }));
  app.use(((error, _req, res, _next) => {
    const status =
      (
        { AuthenticationError: 401, NotAllowedError: 403 } as Record<
          string,
          number
        >
      )[error.name] ?? 500;
    res.status(status).json({ error: { name: error.name } });
  }) as express.ErrorRequestHandler);
  server = app.listen(0);
  base = `http://127.0.0.1:${(server.address() as AddressInfo).port}`;
});
afterAll(() => new Promise(resolve => server.close(resolve)));
beforeEach(() => {
  jest.clearAllMocks();
  clock += POLICY_TTL_MS + 1;
});

describe('GET /policy-reports', () => {
  it('denies an unauthenticated caller and reads nothing', async () => {
    expect((await get()).status).toBe(401);
    expect(list).not.toHaveBeenCalled();
  });

  it('denies a service credential such as the MCP machine token', async () => {
    expect((await get({ authorization: 'Bearer service' })).status).toBe(403);
    expect(list).not.toHaveBeenCalled();
  });

  it('asks the auth service for user principals only', async () => {
    await get(asUser);
    expect(credentials).toHaveBeenCalledWith(expect.anything(), {
      allow: ['user'],
    });
  });

  it('reads only the fixed PolicyReport path, never a caller-supplied one', async () => {
    await fetch(`${base}/policy-reports?path=/api/v1/secrets`, {
      headers: asUser,
    });
    expect(list).toHaveBeenCalledTimes(1);
    expect(list).toHaveBeenCalledWith(POLICY_REPORTS_PATH);
    expect(POLICY_REPORTS_PATH).toBe(
      '/apis/wgpolicyk8s.io/v1alpha2/policyreports',
    );
  });

  it('returns the summary with its observation time and no browser caching', async () => {
    const response = await get(asUser);
    expect(response.status).toBe(200);
    expect(response.headers.get('cache-control')).toBe('no-store');
    const body = await response.json();
    expect(body.status).toBe('ok');
    expect(body.observedAt).toEqual(expect.any(String));
    expect(body.namespaces.app.findings[0]).toMatchObject({
      result: 'fail',
      policy: 'p',
    });
  });

  it('shares one cluster read between callers inside the cache window', async () => {
    await get(asUser);
    await get(asUser);
    expect(list).toHaveBeenCalledTimes(1);
  });

  it('serves the last good data as stale when the cluster read fails', async () => {
    await get(asUser);
    clock += POLICY_TTL_MS + 1;
    list.mockRejectedValueOnce(new Error('Cluster API returned 500'));
    const body = await (await get(asUser)).json();
    expect(body.status).toBe('stale');
    expect(body.error).toBe('Cluster API returned 500');
    expect(body.namespaces.app).toBeDefined();
  });
});

describe('GET /policy-reports with no data ever read', () => {
  it('answers 200 with status error and no namespaces, never an empty all-clear', async () => {
    const fresh = createPolicySource(async () => {
      throw new Error('Cluster API returned 403');
    });
    const app = express();
    app.use(
      createPolicyRoutes({
        httpAuth: { credentials } as never,
        source: fresh,
      }),
    );
    const local = app.listen(0);
    try {
      const url = `http://127.0.0.1:${(local.address() as AddressInfo).port}/policy-reports`;
      const response = await fetch(url, { headers: asUser });
      expect(response.status).toBe(200);
      expect(await response.json()).toEqual({
        status: 'error',
        observedAt: null,
        error: 'Cluster API returned 403',
      });
    } finally {
      await new Promise(resolve => local.close(resolve));
    }
  });
});
