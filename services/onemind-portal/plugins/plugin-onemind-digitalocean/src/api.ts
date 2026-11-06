import { useCallback, useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import type { Summary } from './types';

export type SummaryState =
  | { status: 'loading' }
  | { status: 'forbidden' }
  | { status: 'error' }
  | { status: 'done'; summary: Summary };

/** Reads /api/onemind-digitalocean/summary (administrators only; the backend holds the token). */
export function useSummary() {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  const [state, setState] = useState<SummaryState>({ status: 'loading' });
  const [tick, setTick] = useState(0);
  const refresh = useCallback(() => setTick(n => n + 1), []);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const base = await discovery.getBaseUrl('onemind-digitalocean');
        const response = await fetchApi.fetch(`${base}/summary`);
        if (response.status === 403) {
          if (active) setState({ status: 'forbidden' });
          return;
        }
        if (!response.ok) throw new Error('summary failed');
        const summary = (await response.json()) as Summary;
        if (active) setState({ status: 'done', summary });
      } catch {
        if (active) setState({ status: 'error' });
      }
    })();
    return () => {
      active = false;
    };
  }, [discovery, fetchApi, tick]);

  return { state, refresh };
}
