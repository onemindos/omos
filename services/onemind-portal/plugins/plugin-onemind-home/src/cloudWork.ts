import { useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';

export type CloudSummary = {
  state: string;
  observedAt: string | null;
  expiresAt: string | null;
  open?: number;
  overdue?: number;
  sampled?: number;
  partial?: boolean;
  pages?: number;
  source?: string;
  coverage?: string;
  attemptedAt?: string;
  tasks?: {
    id: string;
    title: string;
    dueDate: string | null;
    overdue: boolean;
  }[];
  link: string;
};

/** Reads the existing `GET /api/onemind/cloud` route (Zoho Projects) every minute. */
export function useCloudWork() {
  const discovery = useApi(discoveryApiRef);
  const api = useApi(fetchApiRef);
  const [data, setData] = useState<CloudSummary>();
  const [error, setError] = useState(false);
  useEffect(() => {
    let active = true;
    const controller = new AbortController();
    const refresh = async () => {
      try {
        const base = await discovery.getBaseUrl('onemind');
        const r = await api.fetch(`${base}/cloud`, {
          signal: controller.signal,
        });
        if (!r.ok) throw new Error();
        const result = await r.json();
        if (active) {
          setData(result);
          setError(false);
        }
      } catch {
        if (active) setError(true);
      }
    };
    refresh();
    const timer = setInterval(refresh, 60000);
    return () => {
      active = false;
      controller.abort();
      clearInterval(timer);
    };
  }, [discovery, api]);
  return { data, error };
}
