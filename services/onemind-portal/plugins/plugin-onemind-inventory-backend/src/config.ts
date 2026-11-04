import type { Config } from '@backstage/config';

/**
 * Settings under `onemind.inventory.live`. The whole feature is off unless `enabled: true`, so the
 * file-based catalog (catalog/cluster.yaml) stays authoritative until Zeus turns it on.
 */
export interface LiveInventoryConfig {
  enabled: boolean;
  /** Seconds between cluster observations. */
  frequencySeconds: number;
  timeoutSeconds: number;
  /** Seconds before the first observation after startup (publishing follows 5 seconds later). */
  initialDelaySeconds: number;
  /** How long a removed workload stays visible with lifecycle "removed" before it is dropped. */
  removedRetentionHours: number;
  /** Exact namespace names that are never tracked. */
  excludedNamespaces: string[];
  /** Namespace globs (only "*" is a wildcard) that are never tracked. */
  excludedNamespacePatterns: string[];
  /** Curated catalog file whose workload entities always win (never duplicated). */
  curatedFile: string;
  clusterName: string;
  apiServer: string;
}

/** Namespaces that hold no workloads we catalog. kube-system is NOT here: F2 maps its workloads to apps. */
export const DEFAULT_EXCLUDED_NAMESPACES = ['kube-node-lease', 'kube-public'];
/** Former OpenChoreo data-plane namespaces (dp-*; the platform is torn down) stay excluded; no provider in this portal publishes them. */
export const DEFAULT_EXCLUDED_PATTERNS = ['dp-*'];

export function readLiveInventoryConfig(config: Config): LiveInventoryConfig {
  const live = config.getOptionalConfig('onemind.inventory.live');
  const num = (key: string, fallback: number) => {
    const value = live?.getOptionalNumber(key) ?? fallback;
    if (!Number.isFinite(value) || value <= 0) {
      throw new Error(
        `onemind.inventory.live.${key} must be a positive number`,
      );
    }
    return value;
  };
  return {
    enabled: live?.getOptionalBoolean('enabled') ?? false,
    frequencySeconds: num('frequencySeconds', 120),
    timeoutSeconds: num('timeoutSeconds', 60),
    initialDelaySeconds: num('initialDelaySeconds', 10),
    removedRetentionHours: num('removedRetentionHours', 24),
    excludedNamespaces:
      live?.getOptionalStringArray('excludedNamespaces') ??
      DEFAULT_EXCLUDED_NAMESPACES,
    excludedNamespacePatterns:
      live?.getOptionalStringArray('excludedNamespacePatterns') ??
      DEFAULT_EXCLUDED_PATTERNS,
    curatedFile:
      live?.getOptionalString('curatedFile') ?? 'catalog/cluster.yaml',
    clusterName: live?.getOptionalString('clusterName') ?? 'do-nyc3-omos-cloud',
    apiServer:
      live?.getOptionalString('apiServer') ?? 'https://kubernetes.default.svc',
  };
}

const escape = (s: string) => s.replace(/[.+?^${}()|[\]\\]/g, '\\$&');

/** True when the namespace is excluded by exact name or by a "*" glob. */
export function isExcludedNamespace(
  namespace: string,
  config: Pick<
    LiveInventoryConfig,
    'excludedNamespaces' | 'excludedNamespacePatterns'
  >,
): boolean {
  if (config.excludedNamespaces.includes(namespace)) return true;
  return config.excludedNamespacePatterns.some(pattern =>
    new RegExp(`^${pattern.split('*').map(escape).join('.*')}$`).test(
      namespace,
    ),
  );
}
