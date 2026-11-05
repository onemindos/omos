import { readdirSync, readFileSync, statSync } from 'node:fs';
import path from 'node:path';
import YAML from 'yaml';
import type { Entity } from '@backstage/catalog-model';
import { READABLE, snapshotReader } from './sources';
import type { ClusterFacts } from './types';

/** Repository root of the checkout the tests run in. */
// eslint-disable-next-line no-restricted-syntax
export const REPO_ROOT = path.resolve(__dirname, '../../../../../..');

function walk(dir: string, out: string[]) {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    if (entry.name === 'node_modules' || entry.name === '.git') continue;
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) walk(full, out);
    else out.push(full);
  }
}

let cache: Map<string, string> | undefined;

/** The same files the production reader loads, read from the local checkout. */
export function localSnapshot(): Map<string, string> {
  if (cache) return cache;
  const files: string[] = [];
  for (const dir of ['platform', 'services/onemind-portal/catalog', 'services/onemind-portal/scripts', 'docs']) {
    walk(path.join(REPO_ROOT, dir), files);
  }
  const out = new Map<string, string>();
  for (const file of files) {
    const rel = path.relative(REPO_ROOT, file).split(path.sep).join('/');
    if (!READABLE.test(rel) || statSync(file).size > 1_500_000) continue;
    out.set(rel, readFileSync(file, 'utf8'));
  }
  cache = out;
  return out;
}

export const localRepo = () => snapshotReader(localSnapshot());

export const catalogEntities = (...files: string[]): Entity[] =>
  files.flatMap(file =>
    YAML.parseAllDocuments(localSnapshot().get(file)!)
      .map(d => d.toJSON())
      .filter(Boolean),
  );

const sha = `@sha256:${'a'.repeat(64)}`;

/** Live facts for the real `platform-stac` application, as the portal would read them. */
export const stacCluster = (): ClusterFacts => ({
  observedAt: '2026-11-05T12:00:00.000Z',
  argoApps: [
    {
      name: 'platform-stac',
      finalizers: ['resources-finalizer.argocd.argoproj.io'],
      automated: false,
      prune: false,
      sourcePaths: ['platform/stac'],
      charts: [],
      valueFiles: [],
      destinationNamespace: 'stac',
      resources: [
        { kind: 'Deployment', name: 'stac-browser', namespace: 'stac' },
        { kind: 'SealedSecret', name: 'stac-db', namespace: 'stac' },
      ],
      health: 'Healthy',
      sync: 'Synced',
    },
  ],
  workloads: [
    { kind: 'Deployment', namespace: 'stac', name: 'stac-browser', images: [`ghcr.io/radiantearth/stac-browser:5.1.0${sha}`], hasProbe: true, claims: [], replicas: 1, ready: 1 },
    { kind: 'Deployment', namespace: 'stac', name: 'stac-api', images: [`ghcr.io/stac-utils/stac-fastapi-pgstac:6.0.0${sha}`], hasProbe: true, claims: [], replicas: 1, ready: 1 },
  ],
  pvcs: [],
  ingresses: [
    { namespace: 'stac', name: 'stac-browser', className: 'tailscale', hosts: [] },
    { namespace: 'stac', name: 'stac', className: 'tailscale', hosts: [] },
  ],
  cronJobs: [{ namespace: 'stac', name: 'webodm-stac-sync' }],
});
