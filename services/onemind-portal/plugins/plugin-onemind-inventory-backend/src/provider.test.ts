import type { Entity } from '@backstage/catalog-model';
import type { EntityProviderConnection } from '@backstage/plugin-catalog-node';
import {
  readLiveInventoryConfig,
  isExcludedNamespace,
  type LiveInventoryConfig,
} from './config';
import { A, entityName, type Workload } from './entities';
import { ObservationStore } from './observation';
import { LiveInventoryProvider, readCurated, toOwnedRecords } from './provider';
import { clusterYamlPath, deployment, liveWorkloads } from './testing';
import { ConfigReader } from '@backstage/config';

const HOUR = 3_600_000;
const config = readLiveInventoryConfig(
  new ConfigReader({ onemind: { inventory: { live: { enabled: true } } } }),
);

/** A provider wired to the recorded cluster and an in-memory "catalog" that remembers its last mutation. */
function harness(settings: LiveInventoryConfig = config) {
  let clock = new Date('2026-10-22T10:00:00.000Z');
  let cluster: () => Promise<Workload[]> = () => liveWorkloads();
  const published = new Map<string, Entity>();
  const connection = {
    applyMutation: async (m: {
      type: string;
      entities: { entity: Entity }[];
    }) => {
      expect(m.type).toBe('full');
      published.clear();
      for (const { entity } of m.entities) {
        expect(published.has(entity.metadata.name)).toBe(false); // never two entities for one name
        published.set(entity.metadata.name, entity);
      }
    },
  } as unknown as EntityProviderConnection;
  const logger = { info: jest.fn(), warn: jest.fn(), error: jest.fn() };
  const provider = new LiveInventoryProvider({
    config: settings,
    logger,
    store: new ObservationStore(),
    listWorkloads: () => cluster(),
    readCurated: () => readCurated(clusterYamlPath),
    readOwned: async () => toOwnedRecords([...published.values()]),
    now: () => clock,
  });
  return {
    provider,
    published,
    logger,
    connect: () => provider.connect(connection),
    setCluster: (fn: () => Promise<Workload[]>) => (cluster = fn),
    advance: (ms: number) => (clock = new Date(clock.getTime() + ms)),
  };
}

const withExtra =
  (...extra: ReturnType<typeof deployment>[]) =>
  () =>
    liveWorkloads(extra);

describe('exclusions', () => {
  it('excludes system namespaces and OpenChoreo dp-* namespaces by default', () => {
    expect(
      isExcludedNamespace('dp-default-onetak-development-b188bf98', config),
    ).toBe(true);
    expect(isExcludedNamespace('kube-node-lease', config)).toBe(true);
    expect(isExcludedNamespace('kube-public', config)).toBe(true);
    expect(isExcludedNamespace('dp', config)).toBe(false);
    expect(isExcludedNamespace('my-dp-ns', config)).toBe(false);
    // kube-system holds workloads F2 maps to apps, so it is tracked unless configured away.
    expect(isExcludedNamespace('kube-system', config)).toBe(false);
    const stricter = readLiveInventoryConfig(
      new ConfigReader({
        onemind: {
          inventory: {
            live: {
              excludedNamespaces: ['kube-system'],
              excludedNamespacePatterns: ['tmp-*'],
            },
          },
        },
      }),
    );
    expect(isExcludedNamespace('kube-system', stricter)).toBe(true);
    expect(isExcludedNamespace('tmp-x', stricter)).toBe(true);
    expect(isExcludedNamespace('dp-x', stricter)).toBe(false);
  });

  it('never publishes a dp-* or kube-node-lease workload, so it cannot duplicate an OpenChoreo entity', async () => {
    const h = harness();
    await h.connect();
    h.setCluster(
      withExtra(
        deployment('dp-default-shop-development-1234abcd', 'shop-api'),
        deployment('kube-node-lease', 'x'),
        deployment('team-a', 'worker'),
      ),
    );
    await h.provider.refresh();
    const names = [...h.published.keys()];
    expect(names).toEqual([entityName('team-a', 'worker')]);
    expect(names).not.toContain(
      entityName('dp-default-shop-development-1234abcd', 'shop-api'),
    );
  });
});

describe('what the provider publishes', () => {
  it('publishes nothing for workloads the reviewed file already defines, so the output equals cluster.yaml today', async () => {
    const h = harness();
    await h.connect();
    const summary = await h.provider.refresh();
    expect(summary).toMatchObject({
      tracked: 170,
      published: 0,
      skippedCurated: 170,
    }); // 171 live, minus the one dp-* workload OpenChoreo owns;
    expect(h.published.size).toBe(0);
  });

  it('a new Deployment appears, placed under its app (mapped) or in Unassigned workloads (unmapped)', async () => {
    const h = harness();
    await h.connect();
    h.setCluster(
      withExtra(
        deployment('argocd', 'argocd-server-two', { app: 'x' }),
        deployment('team-a', 'worker'),
      ),
    );
    const summary = await h.provider.refresh();
    expect(summary.added.sort()).toEqual([
      'argocd/argocd-server-two',
      'team-a/worker',
    ]);
    const worker = h.published.get(entityName('team-a', 'worker'))!;
    expect(worker.spec).toMatchObject({
      system: 'unassigned-workloads',
      lifecycle: 'production',
      dependsOn: ['resource:default/do-nyc3-omos-cloud'],
    });
    expect(worker.metadata.annotations).toMatchObject({
      [A.state]: 'observed',
      [A.owner]: 'onemind-live-inventory',
      [A.selector]: 'app=worker',
      'backstage.io/managed-by-location': 'provider:onemind-live-inventory',
    });
    expect(worker.metadata.annotations![A.observedAt]).toBe(
      '2026-10-22T10:00:00.000Z',
    );
    expect(worker.metadata.title).toBe('Team A — Worker');
  });

  it('skips a workload whose entity name is already defined by the reviewed file, even if mapping changes', async () => {
    const h = harness();
    await h.connect();
    h.setCluster(withExtra(deployment('backstage', 'backstage'))); // identical identity to a curated one
    await h.provider.refresh();
    expect(h.published.size).toBe(0);
  });

  it('keeps titles unique inside a System', async () => {
    const h = harness();
    await h.connect();
    h.setCluster(
      withExtra(
        deployment('team-a', 'worker'),
        deployment('team-a', 'worker-'),
        deployment('team-a', 'worker--'),
      ),
    );
    await h.provider.refresh();
    const titles = [...h.published.values()].map(e => e.metadata.title);
    expect(new Set(titles).size).toBe(titles.length);
  });
});

describe('removal lifecycle', () => {
  it('marks a deleted workload lifecycle removed for 24 hours, then drops it; a return clears the mark', async () => {
    const h = harness();
    await h.connect();
    const key = entityName('team-a', 'worker');
    h.setCluster(withExtra(deployment('team-a', 'worker')));
    await h.provider.refresh();
    expect(h.published.get(key)!.spec!.lifecycle).toBe('production');

    h.setCluster(() => liveWorkloads()); // deleted from the cluster
    h.advance(10 * 60_000);
    await h.provider.refresh();
    const gone = h.published.get(key)!;
    expect(gone.spec!.lifecycle).toBe('removed');
    expect(gone.metadata.annotations).toMatchObject({
      [A.state]: 'removed',
      [A.removedAt]: '2026-10-22T10:10:00.000Z',
      [A.observedAt]: '2026-10-22T10:00:00.000Z',
    });
    expect(gone.metadata.name).toBe(key); // same identity while flagged

    h.advance(12 * HOUR);
    await h.provider.refresh();
    expect(h.published.get(key)!.metadata.annotations![A.removedAt]).toBe(
      '2026-10-22T10:10:00.000Z',
    ); // clock not restarted
    expect(h.published.get(key)!.spec!.lifecycle).toBe('removed');

    h.advance(11 * HOUR + 59 * 60_000); // 23h59m after removal
    await h.provider.refresh();
    expect(h.published.has(key)).toBe(true);

    h.advance(2 * 60_000); // past 24h
    const summary = await h.provider.refresh();
    expect(h.published.has(key)).toBe(false);
    expect(summary.dropped).toEqual(['team-a/worker']);

    // A workload that returns while flagged becomes observed again and loses removed-at.
    h.setCluster(withExtra(deployment('team-a', 'worker')));
    await h.provider.refresh();
    h.setCluster(() => liveWorkloads());
    h.advance(HOUR);
    await h.provider.refresh();
    h.setCluster(withExtra(deployment('team-a', 'worker')));
    h.advance(HOUR);
    await h.provider.refresh();
    const back = h.published.get(key)!;
    expect(back.spec!.lifecycle).toBe('production');
    expect(back.metadata.annotations![A.removedAt]).toBeUndefined();
  });

  it('refuses an empty or failed observation instead of removing everything', async () => {
    const h = harness();
    await h.connect();
    h.setCluster(withExtra(deployment('team-a', 'worker')));
    await h.provider.refresh();
    h.setCluster(async () => []);
    await expect(h.provider.refresh()).rejects.toThrow(/no workloads/);
    h.setCluster(async () => {
      throw new Error('Kubernetes API answered 500');
    });
    await expect(h.provider.refresh()).rejects.toThrow(/500/);
    expect(
      h.published.get(entityName('team-a', 'worker'))!.spec!.lifecycle,
    ).toBe('production');
  });

  it('honours a different retention', async () => {
    const h = harness({ ...config, removedRetentionHours: 1 });
    await h.connect();
    const key = entityName('team-a', 'worker');
    h.setCluster(withExtra(deployment('team-a', 'worker')));
    await h.provider.refresh();
    h.setCluster(() => liveWorkloads());
    await h.provider.refresh();
    h.advance(61 * 60_000);
    await h.provider.refresh();
    expect(h.published.has(key)).toBe(false);
  });
});
