import type { UrlReaderService } from '@backstage/backend-plugin-api';
import { list } from '@internal/plugin-onemind-common-backend';
import type {
  ArgoFact,
  ClusterFacts,
  RepoReader,
  WorkloadFact,
} from './types';

/** Paths the lifecycle features read from Git; everything else is skipped. */
export const READABLE =
  /^(platform\/.*\.(ya?ml|json|md|example)|services\/onemind-portal\/catalog\/.*\.(ya?ml|json)|services\/onemind-portal\/scripts\/test-(catalog-structure|coverage)\.cjs|docs\/STACK\.md)$/;
const MAX_FILE_BYTES = 1_500_000;

/** Serves reads from an in-memory snapshot of the repository. */
export function snapshotReader(files: Map<string, string>): RepoReader {
  return {
    async readFile(path) {
      return files.get(path);
    },
    async readDir(path) {
      const prefix = `${path.replace(/\/+$/, '')}/`;
      return new Map([...files].filter(([file]) => file.startsWith(prefix)));
    },
  };
}

/**
 * One archive download for the whole run. The GitHub URL reader fetches the repository
 * archive on every `readTree`, so the lifecycle code reads once and serves the rest from
 * memory. Authentication is the existing GitHub App integration; there is no token here.
 */
export async function loadGithubSnapshot(
  urlReader: UrlReaderService,
  branch = 'main',
): Promise<RepoReader> {
  const tree = await urlReader.readTree(
    `https://github.com/onemindos/omos/tree/${branch}`,
    {
      filter: (path: string, info?: { size: number }) =>
        READABLE.test(path) && (info?.size ?? 0) <= MAX_FILE_BYTES,
    },
  );
  const files = new Map<string, string>();
  for (const file of await tree.files()) {
    if (!READABLE.test(file.path)) continue;
    files.set(file.path, (await file.content()).toString('utf8'));
  }
  return snapshotReader(files);
}

type Item = any;
const claimsOf = (kind: string, item: Item): string[] => {
  const spec = item.spec?.template?.spec ?? item.spec?.jobTemplate?.spec?.template?.spec;
  const direct = (spec?.volumes ?? [])
    .map((v: Item) => v.persistentVolumeClaim?.claimName)
    .filter(Boolean);
  if (kind !== 'StatefulSet') return direct;
  const replicas = item.spec?.replicas ?? 1;
  const templated = (item.spec?.volumeClaimTemplates ?? []).flatMap((t: Item) =>
    Array.from({ length: replicas }, (_, i) => `${t.metadata.name}-${item.metadata.name}-${i}`),
  );
  return [...direct, ...templated];
};

function workloadFact(kind: WorkloadFact['kind'], item: Item): WorkloadFact {
  const podSpec =
    kind === 'CronJob'
      ? item.spec?.jobTemplate?.spec?.template?.spec
      : item.spec?.template?.spec;
  const containers: Item[] = [
    ...(podSpec?.containers ?? []),
    ...(podSpec?.initContainers ?? []),
  ];
  return {
    kind,
    namespace: item.metadata.namespace,
    name: item.metadata.name,
    images: containers.map(c => String(c.image ?? '')).filter(Boolean),
    hasProbe: (podSpec?.containers ?? []).some(
      (c: Item) => c.readinessProbe || c.livenessProbe || c.startupProbe,
    ),
    claims: claimsOf(kind, item),
    replicas: item.spec?.replicas,
    ready: item.status?.readyReplicas ?? 0,
  };
}

function argoFact(item: Item): ArgoFact {
  const sources: Item[] = item.spec?.sources ?? (item.spec?.source ? [item.spec.source] : []);
  const local = sources.filter(s => /onemindos\/omos/.test(s.repoURL ?? ''));
  return {
    name: item.metadata.name,
    finalizers: item.metadata.finalizers ?? [],
    automated: Boolean(item.spec?.syncPolicy?.automated),
    prune: Boolean(item.spec?.syncPolicy?.automated?.prune),
    sourcePaths: local.map(s => String(s.path ?? '')).filter(Boolean),
    charts: sources.map(s => s.chart).filter(Boolean),
    valueFiles: sources.flatMap(s => s.helm?.valueFiles ?? []),
    destinationNamespace: item.spec?.destination?.namespace,
    resources: (item.status?.resources ?? []).map((r: Item) => ({
      group: r.group,
      kind: r.kind,
      name: r.name,
      namespace: r.namespace,
    })),
    health: item.status?.health?.status,
    sync: item.status?.sync?.status,
  };
}

type ListFn = (path: string) => Promise<Item[]>;

/**
 * Read-only cluster facts through the portal's existing service account. Only GET list
 * calls on resources `backstage-observer` already reads; no new RBAC.
 */
export async function collectClusterFacts(
  read: ListFn = list as unknown as ListFn,
  now: () => Date = () => new Date(),
): Promise<ClusterFacts> {
  const [argo, deployments, statefulsets, daemonsets, cronjobs, pvcs, ingresses] =
    await Promise.all([
      read('/apis/argoproj.io/v1alpha1/namespaces/argocd/applications'),
      read('/apis/apps/v1/deployments'),
      read('/apis/apps/v1/statefulsets'),
      read('/apis/apps/v1/daemonsets'),
      read('/apis/batch/v1/cronjobs'),
      read('/api/v1/persistentvolumeclaims'),
      read('/apis/networking.k8s.io/v1/ingresses'),
    ]);
  return {
    observedAt: now().toISOString(),
    argoApps: argo.map(argoFact),
    workloads: [
      ...deployments.map(i => workloadFact('Deployment', i)),
      ...statefulsets.map(i => workloadFact('StatefulSet', i)),
      ...daemonsets.map(i => workloadFact('DaemonSet', i)),
      ...cronjobs.map(i => workloadFact('CronJob', i)),
    ],
    pvcs: pvcs.map(i => ({
      namespace: i.metadata.namespace,
      name: i.metadata.name,
      storageClass: i.spec?.storageClassName,
      volumeName: i.spec?.volumeName,
      size: i.status?.capacity?.storage ?? i.spec?.resources?.requests?.storage,
      phase: i.status?.phase,
    })),
    ingresses: ingresses.map(i => ({
      namespace: i.metadata.namespace,
      name: i.metadata.name,
      className: i.spec?.ingressClassName,
      hosts: [...new Set<string>((i.spec?.rules ?? []).map((r: Item) => r.host).filter(Boolean))],
    })),
    cronJobs: cronjobs.map(i => ({
      namespace: i.metadata.namespace,
      name: i.metadata.name,
    })),
  };
}
