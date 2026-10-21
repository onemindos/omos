import { useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';

/** Shape of the existing `GET /api/onemind/status` response (operations backend). */
export type OperationsStatus = {
  observedAt: string;
  partial: boolean;
  totalWorkloads: number;
  unavailable: string[];
  unavailableSources?: string[];
  deployments: { name: string; health: string; sync: string }[];
  backups: { name: string; lastSuccess: string | null; suspended: boolean }[];
  serviceHealth?: {
    observedAt?: string;
    checks: { name: string; state: string }[];
  };
};

const PORTAL_BACKUP = 'backstage/backstage-backup';
const PORTAL_BACKUP_MAX_AGE_MS = 36 * 3600000;

/**
 * The single definition of "this backup needs review". `/today`'s Needs
 * attention panel and the home status strip both call it, so the numbers match.
 */
export function backupNeedsReview(
  backup: OperationsStatus['backups'][number],
  now = Date.now(),
) {
  return (
    backup.suspended ||
    !backup.lastSuccess ||
    (backup.name === PORTAL_BACKUP &&
      now - Date.parse(backup.lastSuccess) > PORTAL_BACKUP_MAX_AGE_MS)
  );
}

export type StripCounts = {
  /** Critical service checks that are not `up`. `undefined` when never observed. */
  alerts?: number | 'unknown';
  argo: number | 'unknown';
  workloads: number | 'unknown';
  backups: number | 'unknown';
};

/** Counts for the home status strip. `'unknown'` means the source failed; never read it as zero. */
export function summarizeStatus(
  data: OperationsStatus,
  now = Date.now(),
): StripCounts {
  const missing = new Set(data.unavailableSources ?? []);
  return {
    alerts: data.serviceHealth
      ? data.serviceHealth.checks.filter(check => check.state !== 'up').length
      : 'unknown',
    argo: missing.has('Argo CD') ? 'unknown' : data.deployments.length,
    workloads:
      missing.has('Deployments') || missing.has('StatefulSets')
        ? 'unknown'
        : data.unavailable.length,
    backups: missing.has('Backup schedules')
      ? 'unknown'
      : data.backups.filter(backup => backupNeedsReview(backup, now)).length,
  };
}

export const STATUS_STALE_MS = 180000;
export const isStale = (data: OperationsStatus, now = Date.now()) =>
  now - Date.parse(data.observedAt) > STATUS_STALE_MS;

/** Polls the existing operations status route once a minute. */
export function useOperationsStatus() {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  const [data, setData] = useState<OperationsStatus>();
  const [error, setError] = useState(false);
  useEffect(() => {
    let active = true;
    let pending = false;
    const controller = new AbortController();
    const refresh = async () => {
      if (pending) return;
      pending = true;
      try {
        const base = await discovery.getBaseUrl('onemind');
        const response = await fetchApi.fetch(`${base}/status`, {
          signal: controller.signal,
        });
        if (!response.ok) throw new Error('status unavailable');
        const value = await response.json();
        if (active) {
          setData(value);
          setError(false);
        }
      } catch {
        if (active) setError(true);
      } finally {
        pending = false;
      }
    };
    refresh();
    const timer = setInterval(refresh, 60000);
    return () => {
      active = false;
      controller.abort();
      clearInterval(timer);
    };
  }, [discovery, fetchApi]);
  return { data, error };
}
