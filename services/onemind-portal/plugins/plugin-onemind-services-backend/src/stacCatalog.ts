import type { Entity } from '@backstage/catalog-model';
import { parseEntityRef, stringifyEntityRef } from '@backstage/catalog-model';
import {
  coreServices,
  createBackendModule,
} from '@backstage/backend-plugin-api';
import { catalogProcessingExtensionPoint } from '@backstage/plugin-catalog-node';
import { processingResult } from '@backstage/plugin-catalog-node';
import type {
  CatalogProcessor,
  CatalogProcessorEmit,
} from '@backstage/plugin-catalog-node';
import type { LocationSpec } from '@backstage/plugin-catalog-common';

const api = 'component:default/k8s-onetak-stac-api-746ca5f2';
const database = 'component:default/k8s-onetak-onetak-geo-db-a22783b9';
const contract = 'api:default/onetak-stac-readonly';
const append = (existing: unknown, refs: string[]) => [
  ...new Set([
    ...(Array.isArray(existing)
      ? existing.filter((r): r is string => typeof r === 'string')
      : []),
    ...refs,
  ]),
];

// Curated topology from reviewed manifests, not runtime-health inference. Preserves the
// inventory identity of the STAC API; never registers a competing entity.
export class StacCatalogProcessor implements CatalogProcessor {
  getProcessorName() {
    return 'OneMindStacReviewedRelationships';
  }
  async postProcessEntity(
    entity: Entity,
    _location: LocationSpec,
    emit: CatalogProcessorEmit,
  ) {
    const ref = stringifyEntityRef(entity);
    if (ref !== api) return entity;
    const dependencies = [database];
    entity.spec = {
      ...entity.spec,
      dependsOn: append(entity.spec?.dependsOn, dependencies),
      providesApis: append(entity.spec?.providesApis, [contract]),
    };
    const annotations = entity.metadata.annotations ?? {};
    entity.metadata.annotations = {
      ...annotations,
      'onemind.io/relationship-source':
        'reviewed-stac-manifests; not live health',
      'backstage.io/source-location':
        annotations['backstage.io/source-location'] ??
        'url:https://github.com/onemindos/omos/blob/main/platform/onetak/stac.yaml',
      'argocd/app-name': annotations['argocd/app-name'] ?? 'platform-onetak',
    };
    const runbook =
      'https://onemind-preview.tail717ba0.ts.net/docs/default/component/onemind-runbooks/stac/';
    entity.metadata.links = [
      ...(entity.metadata.links ?? []).filter(l => l.url !== runbook),
      { url: runbook, title: 'STAC diagnosis and recovery' },
    ];
    const source = parseEntityRef(ref);
    for (const [targetRef, type, inverse] of [
      ...dependencies.map(r => [r, 'dependsOn', 'dependencyOf']),
      [contract, 'providesApi', 'apiProvidedBy'],
    ]) {
      const target = parseEntityRef(targetRef);
      emit(processingResult.relation({ source, target, type }));
      emit(
        processingResult.relation({
          source: target,
          target: source,
          type: inverse,
        }),
      );
    }
    return entity;
  }
}
export const stacCatalogModule = createBackendModule({
  pluginId: 'catalog',
  moduleId: 'onemind-stac-topology',
  register(reg) {
    reg.registerInit({
      deps: {
        processing: catalogProcessingExtensionPoint,
        logger: coreServices.logger,
      },
      async init({ processing, logger }) {
        processing.addProcessor(new StacCatalogProcessor());
        logger.info(
          'Registered reviewed STAC catalog topology; no workload or permission mutation',
        );
      },
    });
  },
});
