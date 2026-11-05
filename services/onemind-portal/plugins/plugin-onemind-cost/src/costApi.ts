import { useCallback, useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { catalogApiRef } from '@backstage/plugin-catalog-react';
import {
  NamespaceAllocation,
  parseNamespaceAllocation,
  SpecialAllocation,
  specialAllocations,
} from './allocation';
import { appsFromEntities, CatalogApp } from './mapping';

export const WINDOW_CHOICES = [1, 7, 30] as const;
export const DEFAULT_DAYS = 7;

const iso = (ms: number) => new Date(ms).toISOString().replace('.000Z', 'Z');

/**
 * The last N complete UTC days. Both ends are explicit so OpenCost, the validator and
 * the page all talk about exactly the same window, and a reload gives the same numbers.
 */
export function costWindow(now: number, days: number) {
  const day = 86400000;
  const end = Math.floor(now / day) * day;
  return { start: iso(end - days * day), end: iso(end) };
}

export function allocationUrl(
  base: string,
  window: { start: string; end: string },
) {
  const query = new URLSearchParams({
    window: `${window.start},${window.end}`,
    aggregate: 'namespace',
    accumulate: 'true',
  });
  return `${base}/opencost/allocation?${query}`;
}

export type CostData = {
  status: 'loading' | 'done';
  allocation?: NamespaceAllocation;
  specials: SpecialAllocation[];
  apps: CatalogApp[];
  /** Catalog lookup failed: namespaces are shown without app or pillar. */
  catalogError: boolean;
  observedAt?: string;
};

/**
 * Reads OpenCost through the existing Backstage proxy route (/api/proxy/opencost/allocation)
 * and validates the response with the same contract as scripts/cost-allocation.cjs.
 * A network failure or a non-JSON body is passed on as a missing response, which the
 * validator reports as `unknown`; it is never turned into zero cost.
 */
export function useCostData(days: number, now: () => number = Date.now) {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  const catalog = useApi(catalogApiRef);
  const [data, setData] = useState<CostData>({
    status: 'loading',
    specials: [],
    apps: [],
    catalogError: false,
  });
  const [tick, setTick] = useState(0);
  const refresh = useCallback(() => setTick(value => value + 1), []);

  useEffect(() => {
    let active = true;
    setData(previous => ({ ...previous, status: 'loading' }));
    (async () => {
      const observed = now();
      const window = costWindow(observed, days);
      let payload: unknown;
      try {
        const base = await discovery.getBaseUrl('proxy');
        const response = await fetchApi.fetch(allocationUrl(base, window));
        payload = await response.json().catch(() => undefined);
        // An error status with a JSON body that claims success is still an error.
        if (!response.ok && !(payload && typeof payload === 'object'))
          payload = undefined;
        if (
          !response.ok &&
          payload &&
          (payload as { code?: unknown }).code === 200
        )
          payload = { code: response.status };
      } catch {
        payload = undefined;
      }
      let apps: CatalogApp[] = [];
      let catalogError = false;
      try {
        const result = await catalog.getEntities({
          filter: { kind: 'Component' },
          fields: [
            'metadata.name',
            'metadata.namespace',
            'metadata.title',
            'metadata.annotations',
            'metadata.labels',
          ],
        });
        apps = appsFromEntities(result.items as any);
      } catch {
        catalogError = true;
      }
      if (!active) return;
      setData({
        status: 'done',
        allocation: parseNamespaceAllocation(payload, { window }),
        specials: specialAllocations(payload),
        apps,
        catalogError,
        observedAt: iso(observed),
      });
    })();
    return () => {
      active = false;
    };
  }, [days, tick, discovery, fetchApi, catalog, now]);

  return { data, refresh };
}
