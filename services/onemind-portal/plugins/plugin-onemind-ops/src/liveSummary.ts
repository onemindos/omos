import type { BusEvent } from './contract';
import { freshness } from './liveState';
import type { FeedState } from './liveState';

/** Stale for longer than this reads as offline rather than reconnecting. */
export const OFFLINE_AFTER_MS = 60_000;
export const RECENT_EVENT_WINDOW_MS = 5 * 60_000;

export type Connection = 'connected' | 'reconnecting' | 'offline';

export interface LiveSummary {
  connection: Connection;
  /** Epoch ms of the last status message, if any. */
  lastUpdateMs?: number;
  tracks: number;
  /** Events newer than RECENT_EVENT_WINDOW_MS. */
  recentEvents: number;
  /** Newest first, at most three. */
  latestEvents: BusEvent[];
  podsReady: number;
  podsTotal: number;
  /** Namespaces whose health is not ok. */
  namespacesNeedingAttention: number;
  warnings: number;
}

function connectionOf(state: FeedState, now: number): Connection {
  const fresh = freshness(state, now);
  if (fresh === 'waiting') return 'reconnecting';
  if (fresh === 'stale')
    return now - (state.lastStatusMs ?? 0) <= OFFLINE_AFTER_MS
      ? 'reconnecting'
      : 'offline';
  const nats = state.status?.nats.state;
  if (nats === 'down') return 'offline';
  if (nats === 'connecting') return 'reconnecting';
  return 'connected';
}

/** Small, pure read of the feed for glanceable surfaces such as the home tile. */
export function summarizeFeed(state: FeedState, now: number): LiveSummary {
  let podsReady = 0;
  let podsTotal = 0;
  let attention = 0;
  for (const cell of state.cluster.values()) {
    podsReady += cell.ready;
    podsTotal += cell.total;
    if (cell.health !== 'ok') attention += 1;
  }
  const recent = state.events.filter(e => {
    const ms = Date.parse(e.ts);
    return Number.isFinite(ms) && now - ms <= RECENT_EVENT_WINDOW_MS;
  });
  const newest = [...state.events]
    .sort((a, b) => Date.parse(b.ts) - Date.parse(a.ts))
    .slice(0, 3);
  return {
    connection: connectionOf(state, now),
    lastUpdateMs: state.lastStatusMs,
    tracks: state.tracks.size,
    recentEvents: recent.length,
    latestEvents: newest,
    podsReady,
    podsTotal,
    namespacesNeedingAttention: attention,
    warnings: state.warnings.length,
  };
}
