import { createExtensionTester } from '@backstage/frontend-test-utils';
import { EntityCardBlueprint } from '@backstage/plugin-catalog-react/alpha';
import { backstagePluginPrometheusPlugin } from '@roadiehq/backstage-plugin-prometheus';
import { prometheusCard, prometheusApiExtension } from './prometheusPlugin';

test('native API conversion and scoped card use the public frontend bridge', () => {
  expect(
    [...backstagePluginPrometheusPlugin.getApis()].map(f => f.api.id),
  ).toContain('plugin.prometheus.service');
  expect(prometheusApiExtension).toBeDefined();
  const filter = createExtensionTester(prometheusCard).get(
    EntityCardBlueprint.dataRefs.filterFunction,
  )!;
  expect(
    filter({
      apiVersion: 'backstage.io/v1alpha1',
      kind: 'Component',
      metadata: { name: 'unrelated' },
    }),
  ).toBe(false);
});
