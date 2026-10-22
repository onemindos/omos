import type { Entity } from '@backstage/catalog-model';
import { A } from './entities';
import { ObservationStore } from './observation';
import { LiveInventoryProcessor } from './processor';
import { readLiveInventoryConfig } from './config';
import { ConfigReader } from '@backstage/config';
import { curatedWorkloads } from './testing';

const settings = readLiveInventoryConfig(new ConfigReader({}));
const now = new Date('2026-10-22T10:05:00.000Z');
const at = '2026-10-22T10:04:00.000Z';
const find = (ns: string, name: string): Entity =>
  structuredClone(
    curatedWorkloads().find(
      e =>
        e.metadata.annotations![A.namespace] === ns &&
        e.metadata.annotations![A.workload] === name,
    )!,
  );

function processor(keys: string[], observedAt = at) {
  const store = new ObservationStore();
  store.set(observedAt, keys);
  return new LiveInventoryProcessor(store, settings, () => now);
}

describe('freshness of the reviewed entities', () => {
  it('stamps a seen workload and changes nothing else', async () => {
    const before = find('backstage', 'backstage');
    const after = await processor(['backstage/backstage']).postProcessEntity(
      structuredClone(before),
    );
    expect(after.metadata.annotations).toMatchObject({
      [A.observedAt]: at,
      [A.checkedAt]: at,
      [A.state]: 'observed',
      [A.missing]: '0',
    });
    const strip = (e: Entity) => {
      const c = structuredClone(e);
      for (const k of [A.observedAt, A.checkedAt, A.state, A.missing])
        delete c.metadata.annotations![k];
      return c;
    };
    expect(strip(after)).toEqual(strip(before)); // links, dependsOn, labels, name: byte-identical
    expect(after.metadata.name).toBe(before.metadata.name);
  });

  it('flags a workload that is gone: missing-review-required and lifecycle removed, last-seen kept', async () => {
    const before = find('nats', 'agenttap');
    const seenBefore = before.metadata.annotations![A.observedAt];
    const after = await processor(['backstage/backstage']).postProcessEntity(
      structuredClone(before),
    );
    expect(after.spec!.lifecycle).toBe('removed');
    expect(after.metadata.annotations).toMatchObject({
      [A.state]: 'missing-review-required',
      [A.checkedAt]: at,
      [A.observedAt]: seenBefore,
    });
  });

  it('says nothing when there is no recent observation', async () => {
    const before = find('nats', 'agenttap');
    const stale = await processor(
      [],
      '2026-10-22T09:00:00.000Z',
    ).postProcessEntity(structuredClone(before));
    expect(stale).toEqual(before);
    const none = await new LiveInventoryProcessor(
      new ObservationStore(),
      settings,
      () => now,
    ).postProcessEntity(structuredClone(before));
    expect(none).toEqual(before);
  });

  it('leaves excluded namespaces, other kinds and provider-published entities alone', async () => {
    const dp = find('dp-default-onetak-development-b188bf98', 'stac-browser');
    expect(await processor([]).postProcessEntity(structuredClone(dp))).toEqual(
      dp,
    );
    const resource: Entity = {
      apiVersion: 'backstage.io/v1alpha1',
      kind: 'Resource',
      metadata: { name: 'x', annotations: { [A.source]: 'kubernetes' } },
      spec: {},
    };
    expect(
      await processor([]).postProcessEntity(structuredClone(resource)),
    ).toEqual(resource);
    const owned = find('backstage', 'backstage');
    owned.metadata.annotations![A.owner] = 'onemind-live-inventory';
    expect(
      await processor([]).postProcessEntity(structuredClone(owned)),
    ).toEqual(owned);
  });

  it('does not touch catalog entities that are not Kubernetes inventory (apps, OpenChoreo)', async () => {
    const app: Entity = {
      apiVersion: 'backstage.io/v1alpha1',
      kind: 'Component',
      metadata: { name: 'argocd', annotations: { 'argocd/app-name': 'x' } },
      spec: { type: 'service', lifecycle: 'production' },
    };
    expect(await processor([]).postProcessEntity(structuredClone(app))).toEqual(
      app,
    );
  });
});
