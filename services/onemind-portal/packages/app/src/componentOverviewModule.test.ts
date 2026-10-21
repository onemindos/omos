import { createExtensionTester } from '@backstage/frontend-test-utils';
import { EntityContentLayoutBlueprint } from '@backstage/plugin-catalog-react/alpha';
import { componentOverviewExtension } from './componentOverviewModule';

test('override preserves upstream selection; ordinary catalog entities stay native', () => {
  const tester = createExtensionTester(componentOverviewExtension);
  const filter = tester.get(
    EntityContentLayoutBlueprint.dataRefs.filterFunction,
  )!;
  const entity = (kind: string, annotations = {}) => ({
    apiVersion: 'backstage.io/v1alpha1',
    kind,
    metadata: { name: 'example', labels: annotations },
  });
  expect(filter(entity('Component'))).toBe(false);
  expect(filter(entity('Component', { 'openchoreo.io/managed': 'true' }))).toBe(
    true,
  );
  expect(filter(entity('System', { 'openchoreo.io/managed': 'true' }))).toBe(
    false,
  );
});
