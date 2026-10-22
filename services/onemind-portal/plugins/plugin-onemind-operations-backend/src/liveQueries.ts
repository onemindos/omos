// Reviewed PromQL for the Live Operations page, and the parsers for its answers.
// Fixed strings only: the browser never supplies a query.

export type PromSample = { metric: Record<string, string>; value: number };

export const liveQueries = {
  // Node CPU and memory use (percent), joined to the Kubernetes node name.
  nodeCpu:
    '100 * (1 - avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m]))) * on(instance) group_left(nodename) node_uname_info',
  nodeMemory:
    '100 * (1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * on(instance) group_left(nodename) node_uname_info',
  nodeDisk:
    '100 * (1 - node_filesystem_avail_bytes{mountpoint="/",fstype!~"tmpfs|overlay"} / node_filesystem_size_bytes{mountpoint="/",fstype!~"tmpfs|overlay"}) * on(instance) group_left(nodename) node_uname_info',
  // PVC fullness (percent) and volumes projected to fill within four days.
  pvcUsed:
    '100 * kubelet_volume_stats_used_bytes / kubelet_volume_stats_capacity_bytes',
  pvcFilling:
    'predict_linear(kubelet_volume_stats_available_bytes[6h], 4 * 86400) < 0',
  // VolSync backup sources (1 = a sync is late).
  volsyncOutOfSync: 'volsync_volume_out_of_sync{role="source"}',
} as const;

export type LiveQueryName = keyof typeof liveQueries;

export const PVC_WARN_PERCENT = 80;
export const PVC_CRITICAL_PERCENT = 95;

export function parseVector(body: unknown): PromSample[] {
  const data = (body as { status?: string; data?: { resultType?: string; result?: unknown } })
    ?.data;
  if ((body as { status?: string })?.status !== 'success' || data?.resultType !== 'vector' || !Array.isArray(data.result)) {
    throw new Error('Invalid metrics response');
  }
  const samples: PromSample[] = [];
  for (const row of data.result as { metric?: Record<string, string>; value?: [number, string] }[]) {
    const value = Number(row.value?.[1]);
    if (row.metric && Number.isFinite(value)) samples.push({ metric: row.metric, value });
  }
  return samples;
}

export type NodeUsage = {
  cpuPercent: number | null;
  memoryPercent: number | null;
  diskPercent: number | null;
};

const round = (value: number) => Math.round(value * 10) / 10;

export function parseNodeUsage(input: {
  cpu: PromSample[];
  memory: PromSample[];
  disk: PromSample[];
}): Record<string, NodeUsage> {
  const result: Record<string, NodeUsage> = {};
  const put = (samples: PromSample[], key: keyof NodeUsage) => {
    for (const sample of samples) {
      const node = sample.metric.nodename;
      if (!node) continue;
      result[node] ??= { cpuPercent: null, memoryPercent: null, diskPercent: null };
      result[node][key] = round(sample.value);
    }
  };
  put(input.cpu, 'cpuPercent');
  put(input.memory, 'memoryPercent');
  put(input.disk, 'diskPercent');
  return result;
}

export type PvcItem = {
  namespace: string;
  name: string;
  usedPercent: number;
  fillsWithin4Days: boolean;
  severity: 'warn' | 'critical';
};

/** PVCs at or over the warning level, or projected to fill soon. Worst first. */
export function parsePvcs(used: PromSample[], filling: PromSample[]): PvcItem[] {
  const key = (m: Record<string, string>) => `${m.namespace}/${m.persistentvolumeclaim}`;
  const fills = new Set(filling.map(sample => key(sample.metric)));
  const items: PvcItem[] = [];
  const seen = new Set<string>();
  for (const sample of used) {
    const id = key(sample.metric);
    if (!sample.metric.persistentvolumeclaim || seen.has(id)) continue;
    seen.add(id);
    const projected = fills.has(id);
    if (sample.value < PVC_WARN_PERCENT && !projected) continue;
    items.push({
      namespace: sample.metric.namespace,
      name: sample.metric.persistentvolumeclaim,
      usedPercent: round(sample.value),
      fillsWithin4Days: projected,
      severity: sample.value >= PVC_CRITICAL_PERCENT ? 'critical' : 'warn',
    });
  }
  return items.sort((a, b) => b.usedPercent - a.usedPercent);
}

export type VolsyncSummary = {
  sources: number;
  outOfSync: { namespace: string; name: string }[];
};

export function parseVolsync(samples: PromSample[]): VolsyncSummary {
  const seen = new Map<string, boolean>();
  for (const sample of samples) {
    const id = `${sample.metric.obj_namespace}/${sample.metric.obj_name}`;
    seen.set(id, (seen.get(id) ?? false) || sample.value === 1);
  }
  const outOfSync = [...seen]
    .filter(([, late]) => late)
    .map(([id]) => {
      const [namespace, name] = id.split('/');
      return { namespace, name };
    });
  return { sources: seen.size, outOfSync };
}
