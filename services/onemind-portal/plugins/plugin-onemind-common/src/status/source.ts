import { createApiRef, DiscoveryApi, FetchApi } from '@backstage/core-plugin-api';
import type { StatusSnapshot } from './model';

/**
 * The one seam between the status dot and wherever the data comes from. Implement `read()` and
 * register it as `appStatusApiRef`; nothing else changes. `createOperationsStatusSource` reads the
 * routes that exist today; the Live Ops aggregator (F4) can supply a source of its own later.
 */
export interface AppStatusSource {
  read(): Promise<StatusSnapshot>;
}
export const appStatusApiRef = createApiRef<AppStatusSource>({ id: 'plugin.onemind.app-status' });

/** Scopes accepted by `GET /api/onemind-services/prometheus/:scope` (services backend allow-list). */
export const PROMETHEUS_SCOPES = ['backstage', 'postgres', 'nats'] as const;

type OperationsStatus = {
  observedAt: string;
  unavailable?: string[];
  unavailableSources?: string[];
  deployments?: { name: string; health: string; sync: string }[];
  serviceHealth?: { observedAt?: string; checks: { entity?: string; name: string; state: string }[] };
};
type PrometheusObservation = {
  observedAt: string;
  coverage?: string;
  data?: { groups?: { rules?: { alerts?: { state?: string; labels?: Record<string, string> }[] }[] }[] };
};

const reading = <T>(observedAt: string, value: T) => ({ state: 'observed' as const, observedAt, value });
const unavailable = (reason: string) => ({ state: 'unavailable' as const, reason });

/** Pure: `/api/onemind/status` to the normalized snapshot. A failed source becomes `unavailable`, never empty. */
export function snapshotFromOperationsStatus(status: OperationsStatus): StatusSnapshot {
  const missing = new Set(status.unavailableSources ?? []);
  const at = status.observedAt;
  return {
    observedAt: at,
    argo: missing.has('Argo CD') ? unavailable('Argo CD') : reading(at, status.deployments ?? []),
    notReady: missing.has('Deployments') || missing.has('StatefulSets') ? unavailable('Workloads') : reading(at, status.unavailable ?? []),
    uptime: status.serviceHealth ? reading(status.serviceHealth.observedAt ?? at, status.serviceHealth.checks) : unavailable('Service checks'),
  };
}

/** Pure: one scoped Prometheus observation to firing, critical and pending counts. */
export function alertsFromPrometheus(value: PrometheusObservation): NonNullable<StatusSnapshot['alerts']>[string] {
  if (value.coverage && value.coverage !== 'observed') return unavailable(`coverage ${value.coverage}`);
  let firing = 0, critical = 0, pending = 0;
  for (const group of value.data?.groups ?? [])
    for (const rule of group.rules ?? [])
      for (const alert of rule.alerts ?? []) {
        if (alert.state === 'pending') pending++;
        else {
          firing++;
          if (alert.labels?.severity === 'critical') critical++;
        }
      }
  return reading(value.observedAt, { firing, critical, pending });
}

export function createOperationsStatusSource(discovery: DiscoveryApi, fetchApi: FetchApi): AppStatusSource {
  const get = async (plugin: string, path: string, signal?: AbortSignal) => {
    const response = await fetchApi.fetch(`${await discovery.getBaseUrl(plugin)}${path}`, { signal });
    if (!response.ok) throw new Error(`${plugin}${path} ${response.status}`);
    return response.json();
  };
  return {
    async read() {
      const controller = new AbortController();
      const deadline = setTimeout(() => controller.abort(), 15000);
      try {
        const [status, ...scoped] = await Promise.allSettled([
          get('onemind', '/status', controller.signal),
          ...PROMETHEUS_SCOPES.map(scope => get('onemind-services', `/prometheus/${scope}`, controller.signal)),
        ]);
        if (status.status === 'rejected') throw status.reason;
        const snapshot = snapshotFromOperationsStatus(status.value);
        snapshot.alerts = Object.fromEntries(
          PROMETHEUS_SCOPES.map((scope, i) => {
            const r = scoped[i];
            return [scope, r.status === 'fulfilled' ? alertsFromPrometheus(r.value) : unavailable('Prometheus')];
          }),
        );
        return snapshot;
      } finally {
        clearTimeout(deadline);
      }
    },
  };
}
