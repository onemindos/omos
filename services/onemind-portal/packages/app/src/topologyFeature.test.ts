import topologyPlugin, {
  isTopologyAvailable,
} from '@backstage-community/plugin-topology';

const entity = (annotations?: Record<string, string>) => ({
  apiVersion: 'backstage.io/v1alpha1',
  kind: 'Component',
  metadata: { name: 'x', annotations },
  spec: {},
});

// The Topology tab is discovered through `app.packages: all`; no code in this app
// registers it. This pins what it needs from the catalog so a plugin upgrade that
// changes the rule shows up here.
describe('Topology entity tab', () => {
  it('is a frontend plugin with the topology entity content', () => {
    expect(topologyPlugin.id).toBe('topology');
    const plugin = topologyPlugin as unknown as {
      getExtension(id: string): unknown;
    };
    expect(() =>
      plugin.getExtension('entity-content:topology/topology'),
    ).not.toThrow();
  });

  it('shows for entities with the Kubernetes namespace or id annotation', () => {
    expect(
      isTopologyAvailable(
        entity({ 'backstage.io/kubernetes-namespace': 'backstage' }),
      ),
    ).toBe(true);
    expect(
      isTopologyAvailable(entity({ 'backstage.io/kubernetes-id': 'portal' })),
    ).toBe(true);
  });

  it('stays hidden for entities that only carry a label selector or nothing', () => {
    expect(
      isTopologyAvailable(
        entity({ 'backstage.io/kubernetes-label-selector': 'app=x' }),
      ),
    ).toBe(false);
    expect(isTopologyAvailable(entity())).toBe(false);
  });
});
