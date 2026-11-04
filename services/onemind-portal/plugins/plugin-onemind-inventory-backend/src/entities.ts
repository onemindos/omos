import { createHash } from 'node:crypto';
import type { Entity } from '@backstage/catalog-model';
import {
  mapping as defaultMapping,
  type InventoryMapping,
} from './mapping.generated';
import { placeWorkload } from './placement';

export const OWNER = 'group:default/onemind-admins';
export const PROVIDER_NAME = 'onemind-live-inventory';
export const A = {
  owner: 'onemind.io/inventory-owner',
  source: 'onemind.io/inventory-source',
  observedAt: 'onemind.io/inventory-observed-at',
  checkedAt: 'onemind.io/inventory-last-checked-at',
  state: 'onemind.io/inventory-state',
  removedAt: 'onemind.io/inventory-removed-at',
  missing: 'onemind.io/inventory-missing-observations',
  workload: 'onemind.io/kubernetes-workload',
  workloadKind: 'onemind.io/kubernetes-workload-kind',
  namespace: 'backstage.io/kubernetes-namespace',
  selector: 'backstage.io/kubernetes-label-selector',
} as const;

/** One Deployment or StatefulSet as the cluster reports it, reduced to identity and selector. */
export interface Workload {
  kind: string;
  namespace: string;
  name: string;
  /** "k=v,k=v" from spec.selector.matchLabels, in API order (same string the snapshot script wrote). */
  selector: string;
}

export const workloadKey = (w: Pick<Workload, 'namespace' | 'name'>) =>
  `${w.namespace}/${w.name}`;

/**
 * The entity name scheme of scripts/snapshot-cluster.cjs, byte for byte:
 * `k8s-<ns>-<name>` cut to 54 characters, then `-` and the first 8 hex of sha256("<ns>/<name>").
 * Existing refs, favorites, relations and Argo links depend on it; never change it.
 */
export function entityName(namespace: string, name: string): string {
  const digest = createHash('sha256')
    .update(`${namespace}/${name}`)
    .digest('hex');
  return `${`k8s-${namespace}-${name}`.slice(0, 54)}-${digest.slice(0, 8)}`;
}

export interface BuildOptions {
  clusterName: string;
  observedAt: string;
  state: 'observed' | 'removed';
  /** Set only when state is "removed": when the workload was first seen missing. */
  removedAt?: string;
  mapping?: InventoryMapping;
  /** Titles already used per System (curated and emitted), to keep titles unique inside a System. */
  takenTitles?: Map<string, Set<string>>;
}

/** The catalog Component for one live workload, placed by the reviewed F2 mapping. */
export function buildWorkloadEntity(w: Workload, o: BuildOptions): Entity {
  const place = placeWorkload(w.namespace, w.name, o.mapping ?? defaultMapping);
  const name = entityName(w.namespace, w.name);
  let title = place.title;
  const taken = o.takenTitles?.get(place.system) ?? new Set<string>();
  if (taken.has(title)) {
    const tag = place.suffix ?? w.namespace;
    title = taken.has(`${title} (${tag})`)
      ? `${title} (${tag} ${name.slice(-6)})`
      : `${title} (${tag})`;
  }
  taken.add(title);
  o.takenTitles?.set(place.system, taken);

  const annotations: Record<string, string> = {
    'backstage.io/managed-by-location': `provider:${PROVIDER_NAME}`,
    'backstage.io/managed-by-origin-location': `provider:${PROVIDER_NAME}`,
    [A.selector]: w.selector,
    [A.namespace]: w.namespace,
    [A.source]: 'kubernetes',
    [A.owner]: PROVIDER_NAME,
    [A.observedAt]: o.observedAt,
    [A.state]: o.state,
    [A.workload]: w.name,
    [A.workloadKind]: w.kind,
    'onemind.io/hosting': 'digitalocean-nyc3',
    'onemind.io/pillar-review': place.review,
  };
  if (o.removedAt) annotations[A.removedAt] = o.removedAt;
  const labels: Record<string, string> = {
    'onemind.io/catalog-role': 'workload',
  };
  if (place.pillar) labels['onemind.io/pillar'] = place.pillar;

  return {
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'Component',
    metadata: {
      name,
      title,
      description: `${w.kind} in ${w.namespace}`,
      annotations,
      labels,
      tags: ['kubernetes', w.namespace],
    },
    spec: {
      type: 'service',
      lifecycle: o.state === 'removed' ? 'removed' : 'production',
      owner: OWNER,
      system: place.system,
      ...(place.subcomponentOf ? { subcomponentOf: place.subcomponentOf } : {}),
      dependsOn: [`resource:default/${o.clusterName}`],
    },
  };
}
