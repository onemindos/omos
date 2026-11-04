import { buildOwnerIndex } from './owners';
import { groupByPillarAndApp } from './panels';

const components: Parameters<typeof buildOwnerIndex>[0] = [
  {
    metadata: {
      name: 'matrix',
      title: 'Matrix',
      annotations: { 'argocd/app-name': 'platform-matrix', 'backstage.io/kubernetes-namespace': 'matrix' },
      labels: { 'onemind.io/pillar': 'le' },
    },
  },
  {
    metadata: {
      name: 'nats-dashboard',
      annotations: { 'backstage.io/kubernetes-namespace': 'nats' },
      labels: { 'onemind.io/pillar': 'so' },
    },
  },
  { metadata: { name: 'plain', annotations: {} } },
  { metadata: { name: 'loose', annotations: { 'backstage.io/kubernetes-namespace': 'loose' } } },
];

describe('buildOwnerIndex', () => {
  const index = buildOwnerIndex(components);
  it('finds the owning app by Argo application, then by namespace', () => {
    expect(index.byArgoApp('platform-matrix')).toMatchObject({ name: 'matrix', title: 'Matrix', pillar: 'le', href: '/catalog/default/component/matrix' });
    expect(index.byNamespace('nats')).toMatchObject({ name: 'nats-dashboard', title: 'nats-dashboard' });
    expect(index.resolve({ argoApp: 'platform-matrix', namespace: 'nats' })?.name).toBe('matrix');
    expect(index.resolve({ argoApp: 'unknown', namespace: 'nats' })?.name).toBe('nats-dashboard');
  });
  it('returns nothing, so the item shows as Unassigned', () => {
    expect(index.resolve({ namespace: 'mystery' })).toBeUndefined();
    expect(index.resolve({})).toBeUndefined();
  });
});

describe('groupByPillarAndApp', () => {
  const index = buildOwnerIndex(components);
  it('groups by pillar then app and puts Unassigned last', () => {
    const groups = groupByPillarAndApp(
      [{ ns: 'nats' }, { ns: 'matrix' }, { ns: 'mystery' }, { ns: 'nats' }],
      item => index.byNamespace(item.ns),
    );
    expect(groups.map(g => g.pillar)).toEqual(['LE · Legacy Evolution', 'SO · Sovereign Operations', 'Unassigned']);
    expect(groups[1].apps[0].items).toHaveLength(2);
    expect(groups[2].apps[0].owner).toBeUndefined();
  });
  it('separates apps with no pillar yet from items with no owning app at all', () => {
    const groups = groupByPillarAndApp(
      [{ ns: 'mystery' }, { ns: 'loose' }, { ns: 'nats' }],
      item => index.byNamespace(item.ns),
    );
    expect(groups.map(g => g.pillar)).toEqual(['SO · Sovereign Operations', 'No pillar set yet', 'Unassigned']);
    expect(groups[1].apps[0].owner?.name).toBe('loose');
    expect(groups[2].apps[0].owner).toBeUndefined();
  });
});
