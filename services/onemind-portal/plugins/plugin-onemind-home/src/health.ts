import {
  attentionEntries,
  AttentionEntry,
  isStale,
  OperationsStatus,
} from './status';

/** Points subtracted from 100 for each attention item, by severity. */
export const HEALTH_PENALTY = { critical: 20, warning: 8 } as const;

export type HealthScore = {
  /** 0 to 100. When `partial`, the true score can only be lower. */
  score: number;
  /** The observation is old, or the latest refresh failed. */
  stale: boolean;
  /** Some sources were unavailable, so items may be missing from the count. */
  partial: boolean;
  attentionCount: number;
  /** The most severe attention item (first in the aggregator's order on ties). */
  top?: AttentionEntry;
};

/**
 * Platform score from the real operations status:
 *
 *   score = clamp(100 - 20 x critical items - 8 x warning items, 0, 100)
 *
 * Items are the shared attention list (`attentionEntries`: firing alerts and
 * failing service checks, unready workloads, unhealthy Argo apps, backups to
 * review). Retired and needs-setup backups are not items. Returns `undefined`
 * when there is no status at all (never loaded or failed before any data): the
 * caller must show "unavailable", never a number. `stale` and `partial` are
 * reported so the number is never presented as more certain than it is.
 */
export function healthScore(
  data: OperationsStatus | undefined,
  error = false,
  now = Date.now(),
): HealthScore | undefined {
  if (!data) return undefined;
  const entries = attentionEntries(data, now);
  const penalty = entries.reduce(
    (sum, entry) =>
      sum +
      (entry.severity === 'critical'
        ? HEALTH_PENALTY.critical
        : HEALTH_PENALTY.warning),
    0,
  );
  return {
    score: Math.max(0, Math.min(100, 100 - penalty)),
    stale: error || isStale(data, now),
    partial: !!data.partial,
    attentionCount: entries.length,
    top: entries.find(entry => entry.severity === 'critical') ?? entries[0],
  };
}

export type ScoreBand = 'healthy' | 'degraded' | 'at-risk';

/** Text band shown beside the number so colour is never the only signal. */
export function scoreBand(score: number): ScoreBand {
  if (score >= 90) return 'healthy';
  if (score >= 70) return 'degraded';
  return 'at-risk';
}
