import type { Entity } from '@backstage/catalog-model';
import type { CatalogProcessor } from '@backstage/plugin-catalog-node';
import { A } from './entities';
import { isExcludedNamespace, type LiveInventoryConfig } from './config';
import type { ObservationStore } from './observation';

/**
 * Keeps the reviewed, file-based workload entities (catalog/cluster.yaml) honest against the cluster:
 * it stamps when each was last seen and flags the ones that no longer exist. It changes only the
 * inventory annotations and, for a missing workload, lifecycle. Names, relations, links and every other
 * reviewed field stay exactly as in Git. Nothing is dropped here: removing a reviewed entity is a
 * reviewed change in Git (the hygiene report lists them).
 */
export class LiveInventoryProcessor implements CatalogProcessor {
  constructor(
    private readonly store: ObservationStore,
    private readonly config: Pick<
      LiveInventoryConfig,
      'excludedNamespaces' | 'excludedNamespacePatterns' | 'frequencySeconds'
    >,
    private readonly now: () => Date = () => new Date(),
  ) {}

  getProcessorName() {
    return 'OneMindLiveInventoryFreshness';
  }

  async postProcessEntity(entity: Entity): Promise<Entity> {
    const a = entity.metadata.annotations;
    if (entity.kind !== 'Component' || !a || a[A.source] !== 'kubernetes')
      return entity;
    if (a[A.owner]) return entity; // published by the provider, already current
    const namespace = a[A.namespace];
    const workload = a[A.workload];
    if (!namespace || !workload) return entity;
    if (isExcludedNamespace(namespace, this.config)) return entity;
    const now = this.now();
    const seen = this.store.fresh(now, this.config.frequencySeconds * 3000);
    if (!seen) return entity; // no recent successful observation: say nothing rather than guess
    a[A.checkedAt] = seen.at;
    if (seen.keys.has(`${namespace}/${workload}`)) {
      a[A.observedAt] = seen.at;
      a[A.state] = 'observed';
      a[A.missing] = '0';
    } else {
      a[A.state] = 'missing-review-required';
      entity.spec = { ...entity.spec, lifecycle: 'removed' };
    }
    return entity;
  }
}
