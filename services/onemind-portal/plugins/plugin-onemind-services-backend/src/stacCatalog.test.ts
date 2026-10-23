import type { Entity } from '@backstage/catalog-model';
import { StacCatalogProcessor } from './stacCatalog';

const apiName = 'k8s-onetak-stac-api-746ca5f2';
const entity = (name = apiName): Entity => ({
  apiVersion: 'backstage.io/v1alpha1',
  kind: 'Component',
  metadata: {
    name,
    namespace: 'default',
    annotations: { 'preserve.example/key': 'keep' },
  },
  spec: { owner: 'original-owner', dependsOn: ['resource:default/existing'] },
});
const location = { type: 'file', target: 'fixture.yaml' };

test('the backing API is linked to its database and the read-only contract, keeping identity, owner and other dependencies', async () => {
  const p = new StacCatalogProcessor();
  const e = entity();
  const emit = jest.fn();
  await p.postProcessEntity(e, location, emit);
  expect(e.metadata.name).toBe(apiName);
  expect(e.spec?.owner).toBe('original-owner');
  expect(e.spec?.dependsOn).toEqual([
    'resource:default/existing',
    'component:default/k8s-onetak-onetak-geo-db-a22783b9',
  ]);
  expect(e.spec?.providesApis).toEqual(['api:default/onetak-stac-readonly']);
  expect(e.metadata.annotations?.['preserve.example/key']).toBe('keep');
  expect(e.metadata.annotations?.['argocd/app-name']).toBe('platform-onetak');
  expect(emit).toHaveBeenCalledTimes(4);
  await p.postProcessEntity(e, location, jest.fn());
  expect(e.metadata.links).toHaveLength(1);
  expect(e.spec?.providesApis).toHaveLength(1);
});
test.each(['stac-browser', 'other'])(
  'does not modify any other entity: %s',
  async name => {
    const e = entity(name);
    const before = JSON.stringify(e);
    const emit = jest.fn();
    await new StacCatalogProcessor().postProcessEntity(e, location, emit);
    expect(JSON.stringify(e)).toBe(before);
    expect(emit).not.toHaveBeenCalled();
  },
);
