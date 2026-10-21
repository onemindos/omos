import { createFrontendModule } from '@backstage/frontend-plugin-api';
import catalogPlugin from '@backstage/plugin-catalog/alpha';
import catalogGraphPlugin from '@backstage/plugin-catalog-graph/alpha';

// OpenChoreo 1.3.1 supplies its own About and Relations cards on these layouts.
// Components are deliberately excluded: OneMind's ComponentOverview delegates
// those cards to the host. Unknown kinds and ordinary Backstage entities stay native.
const bespokeKinds = new Set([
  'system',
  'domain',
  'resource',
  'environment',
  'dataplane',
  'clusterdataplane',
  'workflowplane',
  'clusterworkflowplane',
  'observabilityplane',
  'clusterobservabilityplane',
  'deploymentpipeline',
  'componenttype',
  'clustercomponenttype',
  'resourcetype',
  'clusterresourcetype',
  'projecttype',
  'clusterprojecttype',
  'traittype',
  'clustertraittype',
  'workflow',
  'clusterworkflow',
  'componentworkflow',
]);

type EntityShape = {
  kind: string;
  metadata: { labels?: Record<string, string> };
};
export function hasBespokeOverview(entity: EntityShape): boolean {
  return (
    entity.metadata.labels?.['openchoreo.io/managed'] === 'true' &&
    bespokeKinds.has(entity.kind.toLowerCase())
  );
}

export const hostAboutExtension = catalogPlugin
  .getExtension('entity-card:catalog/about')
  .override({
    params: {
      filter: entity =>
        !['user', 'group'].includes(entity.kind.toLowerCase()) &&
        !hasBespokeOverview(entity),
    },
  });

export const hostRelationsExtension = catalogGraphPlugin
  .getExtension('entity-card:catalog-graph/relations')
  .override({ params: { filter: entity => !hasBespokeOverview(entity) } });

export const catalogCompositionModule = createFrontendModule({
  pluginId: 'catalog',
  extensions: [hostAboutExtension],
});
export const graphCompositionModule = createFrontendModule({
  pluginId: 'catalog-graph',
  extensions: [hostRelationsExtension],
});
