import { useCallback, useEffect, useRef, useState } from 'react';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';
import { useAlertRefresh } from './useAlertRefresh';
import type { LiveResponse } from './types';

/** Reads the aggregator. Keeps the last good answer visible if a refresh fails. */
export function useLive() {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  const [data, setData] = useState<LiveResponse>();
  const [error, setError] = useState(false);
  const [refreshedAt, setRefreshedAt] = useState<number>();
  const active = useRef(true);
  const inFlight = useRef(false);

  useEffect(() => {
    active.current = true;
    return () => {
      active.current = false;
    };
  }, []);

  const refresh = useCallback(async () => {
    if (inFlight.current) return;
    inFlight.current = true;
    try {
      const base = await discovery.getBaseUrl('onemind');
      const response = await fetchApi.fetch(`${base}/live`);
      if (!response.ok) throw new Error('live unavailable');
      const value = (await response.json()) as LiveResponse;
      if (active.current) {
        setData(value);
        setError(false);
        setRefreshedAt(Date.now());
      }
    } catch {
      if (active.current) setError(true);
    } finally {
      inFlight.current = false;
    }
  }, [discovery, fetchApi]);

  const { push } = useAlertRefresh(refresh);
  return { data, error, refresh, push, refreshedAt };
}
