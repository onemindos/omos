import type { Workload } from './entities';
import { workloadKey } from './entities';

/** What the provider remembers about an entity it published earlier (read back from the catalog). */
export interface OwnedRecord extends Workload {
  observedAt: string;
  removedAt?: string;
}

export interface LedgerEntry {
  workload: Workload;
  state: 'observed' | 'removed';
  observedAt: string;
  removedAt?: string;
}

export interface LedgerResult {
  entries: LedgerEntry[];
  /** Keys dropped now because they stayed removed longer than the retention. */
  dropped: string[];
  added: string[];
}

/**
 * Live workloads are "observed". A previously published workload that is no longer live is "removed"
 * (lifecycle removed) from the first run that misses it, until retentionMs has passed, then it is dropped.
 * A workload that comes back within the window returns to "observed". `candidates` are the live,
 * non-excluded workloads the provider owns (curated ones are not passed in).
 */
export function reconcile(args: {
  candidates: Workload[];
  previous: OwnedRecord[];
  now: Date;
  retentionMs: number;
}): LedgerResult {
  const nowIso = args.now.toISOString();
  const liveKeys = new Set(args.candidates.map(workloadKey));
  const known = new Set(args.previous.map(workloadKey));
  const entries: LedgerEntry[] = args.candidates.map(workload => ({
    workload,
    state: 'observed',
    observedAt: nowIso,
  }));
  const dropped: string[] = [];
  for (const p of args.previous) {
    if (liveKeys.has(workloadKey(p))) continue;
    const removedAt = p.removedAt ?? nowIso;
    const age = args.now.getTime() - Date.parse(removedAt);
    if (Number.isFinite(age) && age >= args.retentionMs) {
      dropped.push(workloadKey(p));
      continue;
    }
    entries.push({
      workload: {
        kind: p.kind,
        namespace: p.namespace,
        name: p.name,
        selector: p.selector,
      },
      state: 'removed',
      observedAt: p.observedAt,
      removedAt,
    });
  }
  return {
    entries,
    dropped,
    added: args.candidates.map(workloadKey).filter(k => !known.has(k)),
  };
}
