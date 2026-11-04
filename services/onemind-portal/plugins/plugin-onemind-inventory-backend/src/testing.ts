/* eslint-disable no-restricted-syntax, @typescript-eslint/no-require-imports */
/* Test-only helpers: recorded cluster output and the reviewed catalog files (never imported by index.ts). */
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import type { Entity } from '@backstage/catalog-model';
import YAML from 'yaml';
import { A, type Workload } from './entities';
import { listWorkloads } from './kube';

const portal = resolve(__dirname, '../../..');

/** Recorded `kubectl get deployments,statefulsets -A -o json` (reduced; see the fixture's note). */
export const recorded = JSON.parse(
  readFileSync(resolve(__dirname, '__fixtures__/workloads.json'), 'utf8'),
) as {
  items: {
    kind: string;
    metadata: { name: string; namespace: string };
    spec: { selector: { matchLabels: Record<string, string> } };
  }[];
};

/** The recorded cluster as the Kubernetes API would list it, split by kind. */
export function recordedApi(
  extra: typeof recorded.items = [],
  remove: string[] = [],
) {
  const items = [...recorded.items, ...extra].filter(
    i => !remove.includes(`${i.metadata.namespace}/${i.metadata.name}`),
  );
  return async (path: string): Promise<unknown> => {
    const kind = path.includes('/statefulsets') ? 'StatefulSet' : 'Deployment';
    return { items: items.filter(i => i.kind === kind) };
  };
}

export const liveWorkloads = (
  extra: typeof recorded.items = [],
  remove: string[] = [],
): Promise<Workload[]> => listWorkloads(recordedApi(extra, remove));

export const deployment = (
  namespace: string,
  name: string,
  labels: Record<string, string> = { app: name },
) => ({
  kind: 'Deployment',
  metadata: { name, namespace },
  spec: { selector: { matchLabels: labels } },
});

export const clusterYamlPath = resolve(portal, 'catalog/cluster.yaml');

/** The reviewed catalog file's workload Components (what the cluster snapshot + F2 curation produced). */
export function curatedWorkloads(): Entity[] {
  return YAML.parseAllDocuments(readFileSync(clusterYamlPath, 'utf8'))
    .map(d => d.toJSON() as Entity | null)
    .filter(
      (e): e is Entity =>
        !!e &&
        e.kind === 'Component' &&
        e.metadata.annotations?.[A.source] === 'kubernetes',
    );
}

/** The portal's own CommonJS tooling, loaded by path so the tests can prove the two stay in step. */
const script = (name: string) => require(resolve(portal, 'scripts', name));
export const curator = () => script('curate-structure.cjs');
export const snapshotScript = () => script('snapshot-cluster.cjs');
export const generator = () => script('generate-inventory-mapping.cjs');
export const readPortalFile = (relative: string) =>
  readFileSync(resolve(portal, relative), 'utf8');
export const portalPath = (relative: string) => resolve(portal, relative);
