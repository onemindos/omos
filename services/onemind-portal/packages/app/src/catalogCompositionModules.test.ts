import { createExtensionTester } from '@backstage/frontend-test-utils';
import {
  EntityCardBlueprint,
  EntityContentLayoutBlueprint,
} from '@backstage/plugin-catalog-react/alpha';
import openchoreoPlugin from '@openchoreo/backstage-plugin/alpha';
import {
  hasBespokeOverview,
  hostAboutExtension,
  hostRelationsExtension,
} from './catalogCompositionModules';

const entity = (kind: string, managed = true) => {
  const labels: Record<string, string> = managed
    ? { 'openchoreo.io/managed': 'true' }
    : {};
  return {
    apiVersion: 'backstage.io/v1alpha1',
    kind,
    metadata: { name: 'example', labels },
  };
};
const aboutFilter = createExtensionTester(hostAboutExtension).get(
  EntityCardBlueprint.dataRefs.filterFunction,
)!;
const relationsFilter = createExtensionTester(hostRelationsExtension).get(
  EntityCardBlueprint.dataRefs.filterFunction,
)!;

test.each([
  'System',
  'Domain',
  'Resource',
  'Environment',
  'ClusterDataplane',
  'ClusterWorkflow',
  'ClusterComponentType',
])('bespoke %s keeps upstream cards without the host duplicates', kind => {
  expect(hasBespokeOverview(entity(kind))).toBe(true);
  expect(aboutFilter(entity(kind))).toBe(false);
  expect(relationsFilter(entity(kind))).toBe(false);
});

test.each(['Component', 'System', 'Resource', 'API'])(
  'ordinary %s retains native cards',
  kind => {
    expect(aboutFilter(entity(kind, false))).toBe(true);
    expect(relationsFilter(entity(kind, false))).toBe(true);
  },
);

test.each(['Component', 'API', 'FutureKind'])(
  'managed %s without a bespoke upstream pair retains native cards',
  kind => {
    expect(aboutFilter(entity(kind))).toBe(true);
    expect(relationsFilter(entity(kind))).toBe(true);
  },
);

test.each(['User', 'Group'])(
  'preserves native About exclusion for %s',
  kind => {
    expect(aboutFilter(entity(kind, false))).toBe(false);
  },
);

test('composition assumptions are locked to the reviewed upstream version', () => {
  expect(require('@openchoreo/backstage-plugin/package.json').version).toBe(
    '1.3.1',
  );
});

test.each([
  ['System', 'system-overview'],
  ['Domain', 'domain-overview'],
  ['Resource', 'resource-overview'],
  ['Environment', 'environment-overview'],
  ['Dataplane', 'dataplane-overview'],
  ['ClusterDataplane', 'cluster-dataplane-overview'],
  ['WorkflowPlane', 'workflow-plane-overview'],
  ['ClusterWorkflowPlane', 'cluster-workflow-plane-overview'],
  ['ObservabilityPlane', 'observability-plane-overview'],
  ['ClusterObservabilityPlane', 'cluster-observability-plane-overview'],
  ['DeploymentPipeline', 'deployment-pipeline-overview'],
  ['ComponentType', 'component-type-overview'],
  ['ClusterComponentType', 'component-type-overview'],
  ['ResourceType', 'resource-type-overview'],
  ['ClusterResourceType', 'resource-type-overview'],
  ['ProjectType', 'project-type-overview'],
  ['ClusterProjectType', 'project-type-overview'],
  ['TraitType', 'trait-type-overview'],
  ['ClusterTraitType', 'trait-type-overview'],
  ['Workflow', 'workflow-overview-layout'],
  ['ClusterWorkflow', 'workflow-overview-layout'],
  ['ComponentWorkflow', 'component-workflow-overview-layout'],
] as const)(
  'managed %s actually selects the upstream %s layout',
  (kind, name) => {
    const selection = createExtensionTester(
      openchoreoPlugin.getExtension(
        `entity-content-layout:openchoreo/${name}` as const,
      ),
    ).get(EntityContentLayoutBlueprint.dataRefs.filterFunction)!;
    expect(hasBespokeOverview(entity(kind))).toBe(true);
    expect(selection(entity(kind))).toBe(true);
    expect(selection(entity(kind, false))).toBe(false);
  },
);
