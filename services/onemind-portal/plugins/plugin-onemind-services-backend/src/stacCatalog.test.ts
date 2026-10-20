import type { Entity } from '@backstage/catalog-model';
import { StacCatalogProcessor } from './stacCatalog';

const entity = (name = 'stac-browser', managed = true): Entity => ({
  apiVersion: 'backstage.io/v1alpha1',
  kind: 'Component',
  metadata: {
    name,
    namespace: 'default',
    labels: managed ? { 'openchoreo.io/managed': 'true' } : {},
    annotations: { 'preserve.example/key': 'keep' },
  },
  spec: { owner: 'original-owner', dependsOn: ['resource:default/existing'] },
});
const location = { type: 'file', target: 'fixture.yaml' };
test('adds reviewed relations without changing provider identity, owner or unrelated dependencies', async () => {
  const p = new StacCatalogProcessor();
  const e = entity();
  const emit = jest.fn();
  await p.postProcessEntity(e, location, emit);
  expect(e.metadata.name).toBe('stac-browser');
  expect(e.spec?.owner).toBe('original-owner');
  expect(e.spec?.dependsOn).toEqual([
    'resource:default/existing',
    'component:default/k8s-onetak-stac-api-746ca5f2',
  ]);
  expect(e.spec?.consumesApis).toEqual(['api:default/onetak-stac-readonly']);
  expect(e.metadata.annotations?.['preserve.example/key']).toBe('keep');
  expect(emit).toHaveBeenCalledTimes(4);
  await p.postProcessEntity(e, location, jest.fn());
  expect(e.metadata.links).toHaveLength(1);
  expect(e.spec?.consumesApis).toHaveLength(1);
});
test.each([
  ['stac-browser', false],
  ['other', true],
])(
  'does not modify unmanaged/similarly named entity %s',
  async (name, managed) => {
    const p = new StacCatalogProcessor();
    const e = entity(name, managed);
    const before = JSON.stringify(e);
    const emit = jest.fn();
    await p.postProcessEntity(e, location, emit);
    expect(JSON.stringify(e)).toBe(before);
    expect(emit).not.toHaveBeenCalled();
  },
);
test('backing API is linked to actual database and read-only contract, not made OpenChoreo-managed', async () => {
  const e = entity('k8s-onetak-stac-api-746ca5f2', false);
  await new StacCatalogProcessor().postProcessEntity(e, location, jest.fn());
  expect(e.spec?.dependsOn).toContain(
    'component:default/k8s-onetak-onetak-geo-db-a22783b9',
  );
  expect(e.spec?.providesApis).toEqual(['api:default/onetak-stac-readonly']);
  expect(e.metadata.labels?.['openchoreo.io/managed']).toBeUndefined();
});
