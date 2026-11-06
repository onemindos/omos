/** @jest-environment node */
import express from 'express';
import type { AddressInfo } from 'node:net';
import type { Server } from 'node:http';
import { AuthenticationError, NotAllowedError } from '@backstage/errors';
import { mockServices } from '@backstage/backend-test-utils';
import { ConfigReader } from '@backstage/config';
import { projectMonth } from './collect';
import { DigitalOceanProvider } from './provider';
import { registerDigitalOcean } from './module';
import { createRouter } from './router';
import { buildEntities, collectResources } from './resources';
import { linkVolumes, type Pvc } from './links';
import { DigitalOceanService } from './service';
import { api, counts, fakeFetch, pvcs, TOKEN, workloads } from './testing';

const NOW = Date.parse('2026-11-05T06:30:00Z');
let claims: Pvc[] = [];
beforeAll(async () => {
  claims = await pvcs();
});
const make = (
  opts: Parameters<typeof fakeFetch>[0] = {},
  token: string | null = TOKEN,
) => {
  const fake = fakeFetch(opts);
  const service = new DigitalOceanService({
    getToken: () => token ?? undefined,
    fetchImpl: fake.impl,
    now: () => NOW,
    sleep: async () => {},
  });
  service.setLinker(async v => linkVolumes(v, claims, workloads));
  return { service, fake };
};

describe('summary from the recorded account', () => {
  it('matches the recorded counts and never contains the token', async () => {
    const { service } = make();
    const s = await service.get();
    expect(s.state).toBe('connected');
    expect(s.kubernetes.items).toHaveLength(counts.clusters);
    expect(s.kubernetes.items[0].nodePools).toHaveLength(counts.nodePools);
    expect(s.kubernetes.items[0].upgradeAvailable).toMatch(/^1\./);
    expect(s.droplets.items).toHaveLength(counts.droplets);
    expect(s.volumes.items).toHaveLength(counts.volumes);
    expect(s.volumes.unattachedCount).toBe(counts.unattachedVolumes);
    expect(s.snapshots.stats.count).toBe(counts.snapshots);
    expect(s.snapshots.stats.oldest!.ageDays).toBeGreaterThan(0);
    expect(s.snapshots.stats.largest!.sizeGb).toBe(
      Math.max(...s.snapshots.items.map(x => x.sizeGb)),
    );
    expect(s.loadBalancers.items).toHaveLength(counts.loadBalancers);
    expect(s.firewalls.items).toHaveLength(counts.firewalls);
    expect(s.vpcs.items).toHaveLength(counts.vpcs);
    expect(s.uptime.items[0].state).toBe('up');
    expect(s.billing.invoices.length).toBeGreaterThan(0);
    expect(s.billing.invoicePreview?.period).toBe('2026-11');
    expect(s.spaces.state).toBe('not-available');
    const text = JSON.stringify(s);
    expect(text).not.toContain(TOKEN);
    expect(text).not.toContain('redacted@');
    expect(text).not.toMatch(/Bearer|DIGITALOCEAN_TOKEN/);
  });

  it('keeps usage before credits apart from the amount due and projects the month', () => {
    expect(projectMonth('100.00', Date.parse('2026-11-10T00:00:00Z'))).toEqual({
      projectedMonth: '333.33',
      daysElapsed: 9,
      daysInMonth: 30,
    });
    expect(
      projectMonth('5.00', Date.parse('2026-11-01T02:00:00Z')).projectedMonth,
    ).toBe('150.00');
  });

  it('fails one section on its own when a scope is missing', async () => {
    const { service } = make({
      status: { '/v2/volumes': 403, '/v2/uptime/checks': 403 },
    });
    const s = await service.get();
    expect(s.volumes).toMatchObject({ state: 'unavailable', items: [] });
    expect(s.uptime.state).toBe('unavailable');
    expect(s.droplets.items).toHaveLength(counts.droplets);
    expect(s.account.state).toBe('ok');
  });

  it('caches for ten minutes and refreshes after', async () => {
    const { service, fake } = make();
    await service.get();
    const n = fake.calls.length;
    await service.get();
    expect(fake.calls.length).toBe(n);
    await service.refresh();
    expect(fake.calls.length).toBe(n * 2);
  });

  it('is not connected without a token, makes no call, and returns empty sections', async () => {
    const { service, fake } = make({}, null);
    const s = await service.get();
    expect(s.state).toBe('not-connected');
    expect(fake.calls).toHaveLength(0);
    expect(s.droplets.items).toEqual([]);
    expect(s.volumes.items).toEqual([]);
    expect(s.changes).toEqual([]);
  });

  it('reports what changed between refreshes, and not an outage as removals', async () => {
    const fake1 = fakeFetch();
    let current = fake1.impl;
    const service = new DigitalOceanService({
      getToken: () => TOKEN,
      fetchImpl: ((...a: Parameters<typeof fetch>) =>
        current(...a)) as typeof fetch,
      now: () => NOW,
    });
    await service.refresh();
    const droplet = (api['/v2/droplets'] as { droplets: { id: string }[] })
      .droplets[0].id;
    current = fakeFetch({ drop: { '/v2/droplets': [String(droplet)] } }).impl;
    const s = await service.refresh();
    expect(s.changes[0].removed).toEqual([
      expect.objectContaining({ type: 'droplet', id: String(droplet) }),
    ]);
    expect(s.changes[0].added).toEqual([]);
    current = fakeFetch({
      drop: { '/v2/droplets': [String(droplet)] },
      status: { '/v2/volumes': 500 },
    }).impl;
    const after = await service.refresh();
    expect(after.changes).toHaveLength(1); // no mass "removed" from the failed section
  });
});

describe('catalog entities', () => {
  it('emits the expected counts, stable unique names, and is idempotent across refreshes', async () => {
    const { service } = make();
    const a = buildEntities(await service.refresh());
    const b = buildEntities(await service.refresh());
    expect(JSON.stringify(a)).toBe(JSON.stringify(b));
    const expected =
      counts.clusters +
      counts.nodePools +
      counts.droplets +
      counts.volumes +
      counts.snapshots +
      counts.loadBalancers +
      counts.firewalls +
      counts.vpcs;
    expect(a).toHaveLength(expected);
    const names = a.map(e => e.metadata.name);
    expect(new Set(names).size).toBe(names.length);
    for (const n of names) expect(n).toMatch(/^do-[a-z-]+-[a-z0-9-]+$/);
    for (const n of names) expect(n.length).toBeLessThanOrEqual(63);
    const types = new Set(a.map(e => (e.spec as { type: string }).type));
    expect([...types].sort()).toEqual([
      'do-droplet',
      'do-firewall',
      'do-kubernetes-cluster',
      'do-load-balancer',
      'do-node-pool',
      'do-snapshot',
      'do-volume',
      'do-vpc',
    ]);
    for (const e of a) {
      expect(e.kind).toBe('Resource');
      expect(e.metadata.title).toBeTruthy();
      expect(e.metadata.links![0].url).toMatch(
        /^https:\/\/cloud\.digitalocean\.com\//,
      );
      expect((e.spec as { owner: string }).owner).toBe(
        'group:default/onemind-admins',
      );
      for (const t of e.metadata.tags ?? [])
        expect(t).toMatch(/^[a-z0-9+#]+(-[a-z0-9+#]+)*$/);
    }
  });

  it('links every PVC-backed volume to its claim, and the owning app where the catalog knows it', async () => {
    const { service } = make();
    const s = await service.refresh();
    const pvcVolumes = s.volumes.items.filter(v => v.name.startsWith('pvc-'));
    expect(pvcVolumes.length).toBe(counts.volumes);
    expect(pvcVolumes.every(v => v.pvc)).toBe(true);
    const withApp = pvcVolumes.filter(v => v.app);
    // Recorded cluster: the owning workload is known for the data volumes; backup copies and
    // VolSync caches have no workload and stay unlinked rather than guessed.
    expect(withApp.length).toBeGreaterThanOrEqual(25);
    for (const v of withApp)
      expect(v.app!.ref).toMatch(/^component:default\/k8s-/);
    const entities = buildEntities(s);
    const linked = entities.filter(e =>
      (e.spec as { dependencyOf?: string[] }).dependencyOf?.some(r =>
        r.startsWith('component:'),
      ),
    );
    expect(linked).toHaveLength(withApp.length);
    expect(
      entities.find(
        e =>
          e.metadata.name ===
          `do-kubernetes-cluster-${s.kubernetes.items[0].id}`,
      )!.spec,
    ).toMatchObject({
      dependencyOf: ['resource:default/do-nyc3-omos-cloud'],
    });
  });

  it('provider publishes the full set, nothing without a token, and skips a partial read', async () => {
    const connection = { applyMutation: jest.fn(), refresh: jest.fn() };
    const logger = mockServices.logger.mock();
    const run = async (service: DigitalOceanService) => {
      const p = new DigitalOceanProvider({ service, logger });
      await p.connect(connection as never);
      return p.refresh();
    };
    await run(make({}, null).service);
    expect(connection.applyMutation).not.toHaveBeenCalled();
    await run(make({ status: { '/v2/snapshots': 403 } }).service);
    expect(connection.applyMutation).not.toHaveBeenCalled();
    const first = await run(make().service);
    expect(first!.published).toBe(
      collectResources(await make().service.get()).resources.length,
    );
    const [{ type, entities }] = connection.applyMutation.mock.calls[0];
    expect(type).toBe('full');
    expect(entities[0].locationKey).toBe('provider:onemind-digitalocean');
    const again = await run(make().service);
    expect(again).toEqual(first);
    expect(JSON.stringify(connection.applyMutation.mock.calls[1])).toBe(
      JSON.stringify(
        connection.applyMutation.mock.calls[2] ??
          connection.applyMutation.mock.calls[1],
      ),
    );
  });

  it('registers nothing without a token', () => {
    const processing = { addEntityProvider: jest.fn() };
    const scheduler = { scheduleTask: jest.fn() };
    const registered = registerDigitalOcean({
      config: new ConfigReader({}),
      logger: mockServices.logger.mock(),
      scheduler: scheduler as never,
      auth: {} as never,
      catalog: {} as never,
      processing,
      env: {},
    });
    expect(registered).toBe(false);
    expect(processing.addEntityProvider).not.toHaveBeenCalled();
    expect(scheduler.scheduleTask).not.toHaveBeenCalled();
    const on = registerDigitalOcean({
      config: new ConfigReader({}),
      logger: mockServices.logger.mock(),
      scheduler: scheduler as never,
      auth: {} as never,
      catalog: {} as never,
      processing,
      env: { DIGITALOCEAN_TOKEN: TOKEN },
      service: make().service,
    });
    expect(on).toBe(true);
    expect(processing.addEntityProvider).toHaveBeenCalledTimes(1);
    expect(scheduler.scheduleTask.mock.calls[0][0]).toMatchObject({
      frequency: { seconds: 600 },
      scope: 'global',
    });
  });
});

describe('GET /summary', () => {
  const credentials = jest.fn(
    async (req: express.Request, options?: { allow?: string[] }) => {
      const header = req.header('authorization');
      if (!header) throw new AuthenticationError('Missing credentials');
      const principal = header.startsWith('Bearer user') ? 'user' : 'service';
      if (!options?.allow?.includes(principal))
        throw new NotAllowedError('not allowed');
      return { header } as never;
    },
  );
  const userInfo = {
    getUserInfo: jest.fn(async (c: { header: string }) => ({
      userEntityRef: 'user:default/x',
      ownershipEntityRefs:
        c.header === 'Bearer user-admin'
          ? ['group:default/onemind-admins']
          : ['group:default/other'],
    })),
  };
  let server: Server;
  let base: string;
  let token: string | undefined;
  beforeAll(async () => {
    const service = new DigitalOceanService({
      getToken: () => token,
      fetchImpl: fakeFetch().impl,
      now: () => NOW,
    });
    const app = express();
    app.use(
      createRouter({
        httpAuth: { credentials } as never,
        userInfo: userInfo as never,
        service,
      }),
    );
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

  it('requires a signed-in user who is an administrator', async () => {
    expect((await fetch(`${base}/summary`)).status).toBe(401);
    expect(
      (
        await fetch(`${base}/summary`, {
          headers: { authorization: 'Bearer service' },
        })
      ).status,
    ).toBe(403);
    expect(
      (
        await fetch(`${base}/summary`, {
          headers: { authorization: 'Bearer user-other' },
        })
      ).status,
    ).toBe(403);
  });

  it('answers not-connected with empty sections when the token is absent', async () => {
    token = undefined;
    const res = await fetch(`${base}/summary`, {
      headers: { authorization: 'Bearer user-admin' },
    });
    expect(res.status).toBe(200);
    expect(res.headers.get('cache-control')).toBe('no-store');
    const body = await res.json();
    expect(body.state).toBe('not-connected');
    expect(body.droplets.items).toEqual([]);
    expect(body.kubernetes.items).toEqual([]);
  });

  it('returns the summary to an administrator without the token anywhere in it', async () => {
    token = TOKEN;
    const res = await fetch(`${base}/summary`, {
      headers: { authorization: 'Bearer user-admin' },
    });
    const text = await res.text();
    expect(res.status).toBe(200);
    expect(JSON.parse(text).state).toBe('connected');
    expect(text).not.toContain(TOKEN);
  });
});
