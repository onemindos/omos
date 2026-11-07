import type {
  ArgoCell,
  BusEvent,
  ClusterDelta,
  EventsDelta,
  LiveSnapshot,
  LiveStatus,
  NamespaceCell,
  StatusMessage,
  Track,
  TracksDelta,
  WarningEvent,
} from './contract';

/**
 * Pure state for the Ops Center. The hook feeds it a snapshot and then Signals deltas; every function here
 * takes `now` (epoch ms) so tests need no fake timers.
 *
 * Origin rule: the portal runs two replicas and a browser may hear from both. We lock onto the first origin
 * we hear from and ignore the other one, until the locked origin has been silent for ORIGIN_SILENCE_MS. Then
 * we switch and ask for a fresh snapshot (`resyncToken` changes), because the new origin's deltas may not
 * line up with what we have. Within an origin, a message with seq <= the last seen seq is ignored.
 */

export const ORIGIN_SILENCE_MS = 15_000;
export const STALE_AFTER_MS = 15_000;
export const SNAPSHOT_REFRESH_MS = 5 * 60_000;
export const MAX_EVENTS = 500;
export const MAX_WARNINGS = 50;
export const SERIES_SECONDS = 30;

export interface Bucket {
  /** Epoch second. */
  sec: number;
  counts: Record<string, number>;
}

export interface FeedState {
  loaded: boolean;
  /** Origin we apply deltas from. */
  origin?: string;
  /** Highest seq applied or skipped, per origin and channel (key `origin|channel`), so reordering between channels loses nothing. */
  lastSeq: Record<string, number>;
  /** Recent deltas, kept for PENDING_MS so a snapshot that was in flight while they arrived can replay them. */
  pending: PendingDelta[];
  /** When we last heard (any message, even a duplicate) from the locked origin. */
  lastHeardMs?: number;
  /** When the last status message (or the first snapshot) arrived. Drives LIVE / STALE. */
  lastStatusMs?: number;
  /** Incremented when the page should re-fetch the snapshot because the origin changed. */
  resyncToken: number;
  status?: LiveStatus;
  tracks: ReadonlyMap<string, Track>;
  events: BusEvent[];
  droppedEvents: number;
  buckets: Bucket[];
  cluster: ReadonlyMap<string, NamespaceCell>;
  /** When each namespace cell last changed through a delta, for the pulse animation. */
  clusterChangedMs: ReadonlyMap<string, number>;
  argo: ReadonlyMap<string, ArgoCell>;
  warnings: WarningEvent[];
  snapshotAtMs?: number;
}

export const initialFeedState: FeedState = {
  loaded: false,
  lastSeq: {},
  pending: [],
  resyncToken: 0,
  tracks: new Map(),
  events: [],
  droppedEvents: 0,
  buckets: [],
  cluster: new Map(),
  clusterChangedMs: new Map(),
  argo: new Map(),
  warnings: [],
};

export type DeltaAction = Exclude<FeedAction, { type: 'snapshot' }>;
export interface PendingDelta {
  action: DeltaAction;
  at: number;
}
const CHANNELS = ['tracks', 'events', 'cluster', 'status'] as const;
export const PENDING_MS = 60_000;
const seqKey = (origin: string, channel: string) => `${origin}|${channel}`;

export type FeedAction =
  | { type: 'snapshot'; snapshot: LiveSnapshot; now: number }
  | { type: 'tracks'; message: TracksDelta; now: number }
  | { type: 'events'; message: EventsDelta; now: number }
  | { type: 'cluster'; message: ClusterDelta; now: number }
  | { type: 'status'; message: StatusMessage; now: number };

interface Envelope {
  origin: string;
  seq: number;
}

/**
 * Decides whether a delta is applied. Returns the state with origin, seq and heard-time bookkeeping updated,
 * and `apply` saying whether the payload should be merged.
 */
export function admit(
  state: FeedState,
  message: Envelope,
  channel: string,
  now: number,
): { state: FeedState; apply: boolean } {
  const { origin, seq } = message;
  let next = state;
  if (state.origin === undefined) {
    next = { ...state, origin, lastHeardMs: now };
  } else if (state.origin !== origin) {
    const silentFor = now - (state.lastHeardMs ?? 0);
    if (silentFor <= ORIGIN_SILENCE_MS) return { state, apply: false };
    next = {
      ...state,
      origin,
      lastHeardMs: now,
      resyncToken: state.resyncToken + 1,
    };
  } else {
    next = { ...state, lastHeardMs: now };
  }
  const key = seqKey(origin, channel);
  const last = next.lastSeq[key];
  if (last !== undefined && seq <= last) return { state: next, apply: false };
  return {
    state: { ...next, lastSeq: { ...next.lastSeq, [key]: seq } },
    apply: true,
  };
}

const isoMs = (iso: string): number => {
  const ms = Date.parse(iso);
  return Number.isFinite(ms) ? ms : 0;
};

export function applySnapshot(
  state: FeedState,
  snapshot: LiveSnapshot,
  now: number,
): FeedState {
  // The snapshot comes from whichever replica answered. It replaces the base data and sets the seq floor for
  // that origin, so a delta it already contains is not applied twice.
  // seq is one counter per origin, so the snapshot's seq is the floor for every channel of that origin.
  // It replaces what we held for the origin: deltas newer than it are replayed from `pending` by feedReducer.
  const lastSeq = { ...state.lastSeq };
  for (const ch of CHANNELS)
    lastSeq[seqKey(snapshot.origin, ch)] = snapshot.seq;
  const origin = state.origin ?? snapshot.origin;
  return {
    ...state,
    loaded: true,
    origin,
    lastSeq,
    lastHeardMs: state.origin === undefined ? now : state.lastHeardMs,
    lastStatusMs: state.lastStatusMs ?? now,
    status: snapshot.status,
    tracks: new Map(snapshot.tracks.map(t => [t.uid, t])),
    events: snapshot.events.slice(0, MAX_EVENTS),
    cluster: new Map(snapshot.cluster.map(c => [c.namespace, c])),
    clusterChangedMs: new Map(),
    argo: new Map(snapshot.argo.map(a => [a.name, a])),
    warnings: snapshot.warnings.slice(0, MAX_WARNINGS),
    snapshotAtMs: now,
  };
}

function addToBuckets(
  buckets: Bucket[],
  events: BusEvent[],
  now: number,
): Bucket[] {
  const sec = Math.floor(now / 1000);
  const oldest = sec - SERIES_SECONDS + 1;
  const kept = buckets.filter(b => b.sec >= oldest);
  let current = kept.find(b => b.sec === sec);
  const out = current ? kept.filter(b => b !== current) : kept.slice();
  current = { sec, counts: { ...(current?.counts ?? {}) } };
  for (const e of events)
    current.counts[e.kind] = (current.counts[e.kind] ?? 0) + 1;
  out.push(current);
  return out;
}

/** Per-kind counts for each of the last SERIES_SECONDS seconds, oldest first, for the sparklines. */
export function kindSeries(
  buckets: Bucket[],
  now: number,
): Record<string, number[]> {
  const sec = Math.floor(now / 1000);
  const kinds = new Set<string>();
  for (const b of buckets) for (const k of Object.keys(b.counts)) kinds.add(k);
  const bySec = new Map(buckets.map(b => [b.sec, b.counts]));
  const out: Record<string, number[]> = {};
  for (const k of kinds) {
    out[k] = Array.from(
      { length: SERIES_SECONDS },
      (_, i) => bySec.get(sec - SERIES_SECONDS + 1 + i)?.[k] ?? 0,
    );
  }
  return out;
}

const warningKey = (w: WarningEvent) =>
  `${w.namespace}/${w.object}/${w.reason}/${w.ts}`;

const prune = (pending: PendingDelta[], now: number) =>
  pending.filter(p => now - p.at <= PENDING_MS);

export function feedReducer(state: FeedState, action: FeedAction): FeedState {
  if (action.type === 'snapshot') {
    const { origin, seq } = action.snapshot;
    const replay = state.pending.filter(
      p => p.action.message.origin === origin && p.action.message.seq > seq,
    );
    let next = applySnapshot(
      { ...state, pending: [] },
      action.snapshot,
      action.now,
    );
    for (const p of replay) next = applyDelta(next, p.action);
    return next;
  }
  if (!state.loaded) {
    // No base yet: hold the delta until the snapshot lands.
    return {
      ...state,
      pending: [
        ...prune(state.pending, action.now),
        { action, at: action.now },
      ],
    };
  }
  return applyDelta(state, action);
}

function applyDelta(state: FeedState, action: DeltaAction): FeedState {
  const { state: admitted0, apply } = admit(
    state,
    action.message,
    action.type,
    action.now,
  );
  if (!apply) return admitted0;
  const admitted: FeedState = {
    ...admitted0,
    pending: [
      ...prune(admitted0.pending, action.now),
      { action, at: action.now },
    ],
  };

  switch (action.type) {
    case 'status':
      return {
        ...admitted,
        status: action.message.status,
        lastStatusMs: action.now,
      };
    case 'tracks': {
      const tracks = new Map(admitted.tracks);
      for (const uid of action.message.removes) tracks.delete(uid);
      for (const t of action.message.upserts) tracks.set(t.uid, t);
      return { ...admitted, tracks };
    }
    case 'events': {
      const { events, dropped } = action.message;
      const seen = new Set(admitted.events.slice(0, 100).map(e => e.id));
      const fresh = events.filter(e => !seen.has(e.id));
      return {
        ...admitted,
        // Deltas are oldest first; the list is newest first.
        events: [...fresh.slice().reverse(), ...admitted.events].slice(
          0,
          MAX_EVENTS,
        ),
        droppedEvents: admitted.droppedEvents + dropped,
        buckets: addToBuckets(admitted.buckets, fresh, action.now),
      };
    }
    case 'cluster': {
      const { namespaces, removedNamespaces, argo, warnings } = action.message;
      const cluster = new Map(admitted.cluster);
      const changed = new Map(admitted.clusterChangedMs);
      for (const ns of removedNamespaces) {
        cluster.delete(ns);
        changed.delete(ns);
      }
      for (const cell of namespaces) {
        cluster.set(cell.namespace, cell);
        changed.set(cell.namespace, action.now);
      }
      const argoMap = new Map(admitted.argo);
      for (const a of argo) argoMap.set(a.name, a);
      const known = new Set(admitted.warnings.map(warningKey));
      const newWarnings = warnings.filter(w => !known.has(warningKey(w)));
      const merged = [...newWarnings, ...admitted.warnings]
        .sort((a, b) => isoMs(b.ts) - isoMs(a.ts))
        .slice(0, MAX_WARNINGS);
      return {
        ...admitted,
        cluster,
        clusterChangedMs: changed,
        argo: argoMap,
        warnings: merged,
      };
    }
    default:
      return admitted;
  }
}

export type Freshness = 'live' | 'stale' | 'waiting';

/** STALE when no status message (the backend sends one every 5 s) arrived for STALE_AFTER_MS. */
export function freshness(state: FeedState, now: number): Freshness {
  if (state.lastStatusMs === undefined) return 'waiting';
  return now - state.lastStatusMs > STALE_AFTER_MS ? 'stale' : 'live';
}

export type Health = 'ok' | 'warn' | 'down' | 'unknown';

/** Argo apps that are not Synced and Healthy first, then the rest by name. */
export function sortArgo(apps: Iterable<ArgoCell>): ArgoCell[] {
  const good = (a: ArgoCell) => a.sync === 'Synced' && a.health === 'Healthy';
  return [...apps].sort(
    (a, b) => Number(good(a)) - Number(good(b)) || a.name.localeCompare(b.name),
  );
}

const HEALTH_RANK = { down: 0, warn: 1, ok: 2 } as const;

/** Unhealthy namespaces first, then by name. */
export function sortCluster(cells: Iterable<NamespaceCell>): NamespaceCell[] {
  return [...cells].sort(
    (a, b) =>
      HEALTH_RANK[a.health] - HEALTH_RANK[b.health] ||
      a.namespace.localeCompare(b.namespace),
  );
}
