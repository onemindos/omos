import { A, buildWorkloadEntity, entityName } from './entities';
import {
  curatedWorkloads,
  curator,
  portalPath,
  liveWorkloads,
  recorded,
  snapshotScript,
} from './testing';

describe('name stability', () => {
  it('produces the committed entity name for every live workload that catalog/cluster.yaml defines', async () => {
    const curated = new Map(
      curatedWorkloads().map(e => [
        `${e.metadata.annotations![A.namespace]}/${
          e.metadata.annotations![A.workload]
        }`,
        e.metadata.name,
      ]),
    );
    const live = await liveWorkloads();
    expect(live).toHaveLength(149);
    const known = live.filter(w => curated.has(`${w.namespace}/${w.name}`));
    expect(known).toHaveLength(live.length); // every live workload is already in the file today
    for (const w of known) {
      expect(entityName(w.namespace, w.name)).toBe(
        curated.get(`${w.namespace}/${w.name}`),
      );
    }
  });

  it('matches the names scripts/snapshot-cluster.cjs writes, including 54-character truncation', () => {
    const { buildSnapshot } = snapshotScript();
    const docs = buildSnapshot(
      recorded.items,
      'do-nyc3-omos-cloud',
      '2026-10-22T00:00:00.000Z',
      curator().load(portalPath('catalog/apps.yaml')),
    );
    const names = docs
      .filter((d: { kind: string }) => d.kind === 'Component')
      .map((d: { metadata: { name: string } }) => d.metadata.name);
    expect(names).toHaveLength(149);
    expect(names.sort()).toEqual(
      recorded.items
        .map(i => entityName(i.metadata.namespace, i.metadata.name))
        .sort(),
    );
    const long = entityName(
      'openchoreo-observability-plane',
      'a-very-long-workload-name-that-overflows-the-limit',
    );
    expect(long.length).toBeLessThanOrEqual(63);
    expect(long).toMatch(/^k8s-.{50}-[0-9a-f]{8}$/);
  });

  it('keeps one known identity byte-identical (the portal itself)', () => {
    expect(entityName('backstage', 'backstage')).toBe(
      'k8s-backstage-backstage-21b99354',
    );
    const e = buildWorkloadEntity(
      {
        kind: 'Deployment',
        namespace: 'backstage',
        name: 'backstage',
        selector: 'app=backstage',
      },
      {
        clusterName: 'do-nyc3-omos-cloud',
        observedAt: '2026-10-22T00:00:00.000Z',
        state: 'observed',
      },
    );
    expect(e.metadata.name).toBe('k8s-backstage-backstage-21b99354');
  });
});
