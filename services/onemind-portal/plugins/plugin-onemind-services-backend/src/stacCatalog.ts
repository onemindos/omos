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

const browser = 'component:default/stac-browser';
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

// Curated topology from reviewed manifests, not runtime-health inference. Preserve
// the owning OC provider and inventory identities; never register a competing entity.
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
    const managedBrowser =
      ref === browser &&
      entity.metadata.labels?.['openchoreo.io/managed'] === 'true';
    if (!managedBrowser && ref !== api) return entity;
    const dependencies = managedBrowser ? [api] : [database];
    const apiRefs = managedBrowser
      ? {
          key: 'consumesApis',
          relation: 'consumesApi',
          inverse: 'apiConsumedBy',
        }
      : {
          key: 'providesApis',
          relation: 'providesApi',
          inverse: 'apiProvidedBy',
        };
    entity.spec = {
      ...entity.spec,
      dependsOn: append(entity.spec?.dependsOn, dependencies),
      [apiRefs.key]: append(entity.spec?.[apiRefs.key], [contract]),
    };
    const annotations = entity.metadata.annotations ?? {};
    entity.metadata.annotations = {
      ...annotations,
      'onemind.io/relationship-source':
        'reviewed-stac-manifests; not live health',
      'backstage.io/source-location':
        annotations['backstage.io/source-location'] ??
        `url:https://github.com/onemindos/omos/${
          managedBrowser
            ? 'tree/main/platform/openchoreo/projects/onetak/stac-browser'
            : 'blob/main/platform/onetak/stac.yaml'
        }`,
      'argocd/app-name':
        annotations['argocd/app-name'] ??
        (managedBrowser ? 'platform-openchoreo-projects' : 'platform-onetak'),
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
      [contract, apiRefs.relation, apiRefs.inverse],
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
