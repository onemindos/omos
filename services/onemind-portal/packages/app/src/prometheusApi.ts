import {
  createApiRef,
  DiscoveryApi,
  FetchApi,
} from '@backstage/core-plugin-api';

export const scopes = {
  backstage: 'k8s-backstage-backstage-21b99354',
  postgres: 'k8s-backstage-backstage-postgres-17d2f351',
  nats: 'k8s-nats-nats-0dee54ac',
} as const;
export type Scope = keyof typeof scopes;
export type Observation = {
  status: 'success';
  schemaVersion: 'onemind.prometheus-alerts.v1';
  scope: Scope;
  observedAt: string;
  coverage: 'observed' | 'partial' | 'not-configured';
  candidateRules: number;
  unhealthyRules: number;
  staleRules: number;
  data: { groups: { rules: { alerts: unknown[] }[] }[] };
};
export interface ScopedPrometheusApi {
  getAlerts(options: { serviceName?: string }): Promise<Observation>;
  getUiUrl(options: { serviceName?: string }): undefined;
}
// The native plugin's public legacy API factory registers this exact API ID.
// A supported frontend override replaces only its unbounded/general-purpose client.
export const scopedPrometheusApiRef = createApiRef<ScopedPrometheusApi>({
  id: 'plugin.prometheus.service',
});
export function scopeForEntity(entity: {
  kind: string;
  metadata: {
    name: string;
    namespace?: string;
    annotations?: Record<string, string>;
  };
}): Scope | undefined {
  const scope = entity.metadata.annotations?.[
    'prometheus.io/service-name'
  ] as Scope;
  return entity.kind.toLowerCase() === 'component' &&
    (entity.metadata.namespace ?? 'default') === 'default' &&
    scopes[scope] === entity.metadata.name
    ? scope
    : undefined;
}
export function createScopedPrometheusApi(
  discovery: DiscoveryApi,
  fetchApi: FetchApi,
): ScopedPrometheusApi {
  const cache = new Map<Scope, { at: number; value: Observation }>();
  const pending = new Map<Scope, Promise<Observation>>();
  return {
    getUiUrl: () => undefined, // No unverified/public Prometheus URL or arbitrary instance routing.
    async getAlerts({ serviceName }) {
      if (!Object.prototype.hasOwnProperty.call(scopes, serviceName ?? ''))
        throw new Error('Unsupported alert scope');
      const scope = serviceName as Scope;
      const cached = cache.get(scope);
      if (cached && Date.now() - cached.at < 5000) return cached.value;
      if (pending.has(scope)) return pending.get(scope)!;
      const request = (async () => {
        const controller = new AbortController();
        const deadline = setTimeout(() => controller.abort(), 10000);
        try {
          const base = await discovery.getBaseUrl('onemind-services');
          const response = await fetchApi.fetch(`${base}/prometheus/${scope}`, {
            signal: controller.signal,
          });
          if (!response.ok)
            throw new Error('Prometheus observation unavailable');
          const value = await response.json();
          if (
            value.status !== 'success' ||
            value.schemaVersion !== 'onemind.prometheus-alerts.v1' ||
            value.scope !== scope ||
            !Number.isFinite(Date.parse(value.observedAt)) ||
            Date.now() - Date.parse(value.observedAt) > 60000 ||
            Date.parse(value.observedAt) - Date.now() > 30000 ||
            !['observed', 'partial', 'not-configured'].includes(
              value.coverage,
            ) ||
            !Array.isArray(value.data?.groups)
          )
            throw new Error('Invalid or stale Prometheus observation');
          cache.set(scope, { at: Date.now(), value });
          return value;
        } finally {
          clearTimeout(deadline);
        }
      })().finally(() => pending.delete(scope));
      pending.set(scope, request);
      return request;
    },
  };
}
