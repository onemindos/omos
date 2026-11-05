/** Facts shared by the retire inventory and the stack hygiene report. */

export type ArgoFact = {
  name: string;
  finalizers: string[];
  automated: boolean;
  prune: boolean;
  /** `path` of every source that points at this repository. */
  sourcePaths: string[];
  /** Helm chart name for a third-party chart source. */
  charts: string[];
  /** Repository paths named by helm valueFiles. */
  valueFiles: string[];
  destinationNamespace?: string;
  /** Everything Argo says it manages (kind and name only). */
  resources: {
    group?: string;
    kind: string;
    name: string;
    namespace?: string;
  }[];
  health?: string;
  sync?: string;
};

export type WorkloadFact = {
  kind: 'Deployment' | 'StatefulSet' | 'DaemonSet' | 'CronJob';
  namespace: string;
  name: string;
  images: string[];
  hasProbe: boolean;
  /** PersistentVolumeClaim names the workload mounts (StatefulSet claim templates expanded). */
  claims: string[];
  replicas?: number;
  ready?: number;
};

export type PvcFact = {
  namespace: string;
  name: string;
  storageClass?: string;
  volumeName?: string;
  size?: string;
  phase?: string;
};

export type IngressFact = {
  namespace: string;
  name: string;
  className?: string;
  hosts: string[];
};

export type ClusterFacts = {
  observedAt: string;
  argoApps: ArgoFact[];
  workloads: WorkloadFact[];
  pvcs: PvcFact[];
  ingresses: IngressFact[];
  /** CronJobs per namespace (backup jobs show up here). */
  cronJobs: { namespace: string; name: string }[];
};

/** Read-only view of the repository, backed by the GitHub App through the URL reader. */
export interface RepoReader {
  readFile(path: string): Promise<string | undefined>;
  /** Files below a directory (recursive), keyed by repo-relative path. */
  readDir(path: string): Promise<Map<string, string>>;
}

export type InventoryState = 'found' | 'none' | 'unavailable';

export type InventoryItem = { label: string; detail?: string; link?: string };

export type InventorySection = {
  id: string;
  title: string;
  state: InventoryState;
  items: InventoryItem[];
  /** Why a section is unavailable, or a caveat about what "none" means. */
  note?: string;
};

export type Retention = {
  /** What to do with the data before anything is deleted. */
  mode: 'keep-last-backup' | 'export-first' | 'no-retention';
  days: number;
};
