import { createElement } from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { MockEntityListContextProvider } from '@backstage/plugin-catalog-react/testUtils';
import {
  EntityPillarFilter,
  EntityRoleFilter,
  PILLAR_LABEL,
  PILLAR_REVIEW,
  PillarPicker,
  ROLE_LABEL,
  ViewPicker,
  currentView,
  pillarOptions,
  viewFilters,
  viewOptions,
} from './catalogViewFilters';
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

describe('pillar filter contract', () => {
  it.each(['so', 'le', 'hp', 'ge'] as const)(
    '%s matches the reviewed primary label only',
    pillar => {
      expect(new EntityPillarFilter(pillar).getCatalogFilters()).toEqual({
        [PILLAR_LABEL]: pillar,
      });
    },
  );
  it('Shared is the reviewed shared label; Unclassified is a positive state, never a guess', () => {
    expect(new EntityPillarFilter('shared').getCatalogFilters()).toEqual({
      [PILLAR_LABEL]: 'shared',
    });
    expect(new EntityPillarFilter('unclassified').getCatalogFilters()).toEqual({
      [PILLAR_REVIEW]: 'unclassified',
    });
    expect(new EntityPillarFilter('all').getCatalogFilters()).toEqual({});
  });
  it('offers All, SO, LE, HP, GE, Shared and Unclassified in that order', () => {
    expect(pillarOptions.map(o => o.id)).toEqual([
      'all',
      'so',
      'le',
      'hp',
      'ge',
      'shared',
      'unclassified',
    ]);
  });
});

describe('saved views', () => {
  it('map to exact kind, type and role filters', () => {
    expect(viewFilters('apps').kind?.getCatalogFilters()).toEqual({
      kind: 'component',
    });
    expect(viewFilters('apps').role?.getCatalogFilters()).toEqual({
      [ROLE_LABEL]: 'app',
    });
    expect(viewFilters('infrastructure').role?.value).toBe('workload');
    const repos = viewFilters('repos');
    expect(repos.kind?.value).toBe('resource');
    expect(repos.type?.getCatalogFilters()).toEqual({
      'spec.type': ['code-repository'],
    });
    expect(viewFilters('everything').role).toBeUndefined();
    expect(viewOptions.map(o => o.label)).toEqual([
      'Apps',
      'Infrastructure',
      'Repos',
      'Everything',
    ]);
  });
  it('currentView round-trips every preset', () => {
    for (const { id } of viewOptions)
      expect(currentView(viewFilters(id) as never)).toBe(id);
  });
});

describe('pickers', () => {
  it('Pillar picker applies and clears the filter', () => {
    const updateFilters = jest.fn();
    render(
      createElement(
        MockEntityListContextProvider,
        { value: { filters: {}, updateFilters } },
        createElement(PillarPicker),
      ),
    );
    fireEvent.click(screen.getByText('Unclassified'));
    expect(updateFilters).toHaveBeenLastCalledWith({
      pillar: expect.any(EntityPillarFilter),
    });
    expect(updateFilters.mock.calls[0][0].pillar.value).toBe('unclassified');
    fireEvent.click(screen.getByText('All'));
    expect(updateFilters).toHaveBeenLastCalledWith({ pillar: undefined });
  });

  it('View picker starts on Apps, and a deep link to another list is left alone', () => {
    const fresh = jest.fn();
    render(
      createElement(
        MockEntityListContextProvider,
        { value: { filters: {}, updateFilters: fresh } },
        createElement(ViewPicker),
      ),
    );
    expect(fresh).toHaveBeenCalledTimes(1);
    expect(fresh.mock.calls[0][0].role).toBeInstanceOf(EntityRoleFilter);
    const linked = jest.fn();
    render(
      createElement(
        MockEntityListContextProvider,
        {
          value: {
            filters: { kind: viewFilters('repos').kind },
            updateFilters: linked,
          },
        },
        createElement(ViewPicker),
      ),
    );
    expect(linked).not.toHaveBeenCalled();
  });
});

// Runs the System-redirect routing tests inside the step the gate already names.
import './systemRedirects.test';
