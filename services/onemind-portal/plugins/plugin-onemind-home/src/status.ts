import { useCallback, useEffect, useRef, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { useAlertRefresh } from '@internal/plugin-onemind-live';

/** Shape of the existing `GET /api/onemind/status` response (operations backend). */
export type OperationsStatus = {
  observedAt: string;
  partial: boolean;
  totalWorkloads: number;
  unavailable: string[];
  unavailableSources?: string[];
  deployments: { name: string; health: string; sync: string }[];
  backups: BackupStatus[];
  serviceHealth?: {
    observedAt?: string;
    scope?: string;
    metricsUnavailable?: boolean;
    checks: {
      name: string;
      state: string;
      entity?: string;
      link?: string;
      scope?: string;
      httpStatus?: number;
      latencyMs?: number;
    }[];
    metrics?: {
      namespace: string;
      cpuCores: number | null;
      memoryBytes: number | null;
    }[];
  };
  /** Added by the Live Operations aggregator. Absent from an older backend. */
  counts?: Partial<Record<keyof StripCounts, number | 'unknown'>> & {
    backupsNeedSetup?: number | 'unknown';
    backupsRetired?: number | 'unknown';
  };
  attention?: AttentionEntry[];
  needsSetup?: AttentionEntry[];
};

export type BackupStatus = {
  name: string;
  lastSuccess: string | null;
  suspended: boolean;
  /** Set by the aggregator from the CronJob's onemind.io/backup-state annotation. */
  state?: 'active' | 'retired' | 'needs-setup';
  attention?: 'none' | 'review' | 'setup';
  health?: string;
  supersededBy?: string;
  pendingReason?: string;
};

export type AttentionEntry = {
  key: string;
  kind: string;
  severity: 'critical' | 'warning';
  title: string;
  detail: string;
  url: string;
};

const PORTAL_BACKUP = 'backstage/backstage-backup';
const PORTAL_BACKUP_MAX_AGE_MS = 36 * 3600000;

/**
 * The single definition of "this backup needs review". The aggregator decides
 * (retired and needs-setup jobs are not attention items); `/today`'s panel, the
 * status strip and Live Ops all read that decision. The rule below only runs for
 * a status payload from a backend that predates backup states.
 */
export function backupNeedsReview(backup: BackupStatus, now = Date.now()) {
  if (backup.attention) return backup.attention === 'review';
  return (
    backup.suspended ||
    !backup.lastSuccess ||
    (backup.name === PORTAL_BACKUP &&
      now - Date.parse(backup.lastSuccess) > PORTAL_BACKUP_MAX_AGE_MS)
  );
}

export const backupNeedsSetup = (backup: BackupStatus) =>
  backup.attention === 'setup';
export const backupRetired = (backup: BackupStatus) =>
  backup.state === 'retired';

export type StripCounts = {
  /** Critical service checks that are not `up`. `undefined` when never observed. */
  alerts?: number | 'unknown';
  argo: number | 'unknown';
  workloads: number | 'unknown';
  backups: number | 'unknown';
  /** Backups waiting for setup only a person can do (never an attention item). */
  backupsNeedSetup?: number | 'unknown';
  /** Backups replaced on purpose. */
  backupsRetired?: number | 'unknown';
};

/** Counts for the home status strip. `'unknown'` means the source failed; never read it as zero. */
export function summarizeStatus(
  data: OperationsStatus,
  now = Date.now(),
): StripCounts {
  // The aggregator's counts are the shared answer; the arithmetic below is the
  // fallback for a backend that does not send them.
  if (data.counts?.argo !== undefined && data.counts.alerts !== undefined) {
    const { counts } = data;
    return {
      alerts: counts.alerts as StripCounts['alerts'],
      argo: counts.argo as StripCounts['argo'],
      workloads: counts.workloads ?? 'unknown',
      backups: counts.backups ?? 'unknown',
      backupsNeedSetup: counts.backupsNeedSetup,
      backupsRetired: counts.backupsRetired,
    };
  }
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
    backupsNeedSetup: missing.has('Backup schedules')
      ? 'unknown'
      : data.backups.filter(backupNeedsSetup).length,
    backupsRetired: missing.has('Backup schedules')
      ? 'unknown'
      : data.backups.filter(backupRetired).length,
  };
}

const HEADLAMP = 'https://headlamp.tail717ba0.ts.net/c/main/workloads';
const ARGOCD = 'https://argocd.tail717ba0.ts.net/applications/argocd';

/**
 * The "Needs attention" list for `/today`. The aggregator builds it (one rule set
 * for Home, Today and Live Ops); this only rebuilds it for a payload that has none.
 */
export function attentionEntries(
  data: OperationsStatus,
  now = Date.now(),
): AttentionEntry[] {
  if (data.attention) return data.attention;
  return [
    ...(data.serviceHealth?.checks
      .filter(check => check.state !== 'up')
      .map(check => ({
        key: `health:${check.name}`,
        kind: 'service',
        severity: 'warning' as const,
        title: `${check.name}: ${check.state}`,
        detail: 'Inspect service health and connectivity; unknown is not healthy.',
        url: HEADLAMP,
      })) ?? []),
    ...data.unavailable.map(name => ({
      key: `workload:${name}`,
      kind: 'workload',
      severity: 'warning' as const,
      title: `Workload not fully ready: ${name}`,
      detail: 'Inspect pod readiness and events in Headlamp.',
      url: HEADLAMP,
    })),
    ...data.deployments.map(item => ({
      key: `deployment:${item.name}`,
      kind: 'deployment',
      severity: 'warning' as const,
      title: `${item.name}: ${item.health} / ${item.sync}`,
      detail: 'Inspect the deployment health, sync difference, and last operation.',
      url: `${ARGOCD}/${encodeURIComponent(item.name)}`,
    })),
    ...data.backups
      .filter(backup => backupNeedsReview(backup, now))
      .map(backup => ({
        key: `backup:${backup.name}`,
        kind: 'backup',
        severity: 'warning' as const,
        title: `Backup ${backup.name}: ${
          backup.suspended
            ? 'suspended'
            : backup.lastSuccess
            ? 'overdue'
            : 'no successful run recorded'
        }`,
        detail: backup.lastSuccess
          ? `Last successful job: ${new Date(backup.lastSuccess).toLocaleString()}`
          : 'Inspect the backup schedule and job logs.',
        url: HEADLAMP,
      })),
  ];
}

export const STATUS_STALE_MS = 180000;
export const isStale = (data: OperationsStatus, now = Date.now()) =>
  now - Date.parse(data.observedAt) > STATUS_STALE_MS;

/**
 * Reads the operations status route: on mount, every 30 s, and straight after an
 * alert Signal (`onemind:alerts`) when Signals is available.
 */
export function useOperationsStatus() {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  const [data, setData] = useState<OperationsStatus>();
  const [error, setError] = useState(false);
  const active = useRef(true);
  const pending = useRef(false);
  useEffect(() => {
    active.current = true;
    return () => {
      active.current = false;
    };
  }, []);
  const refresh = useCallback(async () => {
    if (pending.current) return;
    pending.current = true;
    try {
      const base = await discovery.getBaseUrl('onemind');
      const response = await fetchApi.fetch(`${base}/status`);
      if (!response.ok) throw new Error('status unavailable');
      const value = await response.json();
      if (active.current) {
        setData(value);
        setError(false);
      }
    } catch {
      if (active.current) setError(true);
    } finally {
      pending.current = false;
    }
  }, [discovery, fetchApi]);
  useAlertRefresh(refresh);
  return { data, error };
}
