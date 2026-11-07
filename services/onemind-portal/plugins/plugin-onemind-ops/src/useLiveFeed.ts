import { useCallback, useEffect, useReducer, useRef, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { useSignal } from '@backstage/plugin-signals-react';
import { LIVE_CHANNELS } from './contract';
import type {
  ClusterDelta,
  EventsDelta,
  LiveSnapshot,
  StatusMessage,
  TracksDelta,
} from './contract';
import {
  SNAPSHOT_REFRESH_MS,
  feedReducer,
  freshness,
  initialFeedState,
} from './liveState';
import type { FeedState, Freshness } from './liveState';

export type LoadState =
  | { kind: 'loading' }
  | { kind: 'ready' }
  /** The backend route answered 404: the livebridge plugin is not deployed. */
  | { kind: 'not-deployed' }
  | { kind: 'error'; message: string };

export interface LiveFeed {
  state: FeedState;
  load: LoadState;
  freshness: Freshness;
  /** Epoch ms, ticks once a second so "2 s ago" labels move. */
  now: number;
  /** False when the Signals client is not connected, so deltas cannot arrive. */
  signalsConnected: boolean;
  refresh: () => void;
}

/**
 * Loads GET /api/livebridge/snapshot, then applies the four Signals channels through the pure reducer.
 * Re-fetches the snapshot when the reducer asks for it (origin switch) and every five minutes as a safety net.
 */
export function useLiveFeed(): LiveFeed {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  const [state, dispatch] = useReducer(feedReducer, initialFeedState);
  const [load, setLoad] = useState<LoadState>({ kind: 'loading' });
  const [now, setNow] = useState(() => Date.now());
  const inFlight = useRef(false);
  const alive = useRef(true);

  const refresh = useCallback(async () => {
    if (inFlight.current) return;
    inFlight.current = true;
    try {
      const base = await discovery.getBaseUrl('livebridge');
      const response = await fetchApi.fetch(`${base}/snapshot`);
      if (response.status === 404) {
        if (alive.current) setLoad({ kind: 'not-deployed' });
        return;
      }
      if (!response.ok) throw new Error(`snapshot answered ${response.status}`);
      const snapshot = (await response.json()) as LiveSnapshot;
      if (!alive.current) return;
      dispatch({ type: 'snapshot', snapshot, now: Date.now() });
      setLoad({ kind: 'ready' });
    } catch (e) {
      if (alive.current)
        setLoad({
          kind: 'error',
          message: e instanceof Error ? e.message : String(e),
        });
    } finally {
      inFlight.current = false;
    }
  }, [discovery, fetchApi]);

  useEffect(() => {
    alive.current = true;
    refresh();
    const safety = setInterval(refresh, SNAPSHOT_REFRESH_MS);
    const tick = setInterval(() => setNow(Date.now()), 1000);
    return () => {
      alive.current = false;
      clearInterval(safety);
      clearInterval(tick);
    };
  }, [refresh]);

  // Origin switch: the reducer bumps resyncToken, we fetch again.
  const lastToken = useRef(state.resyncToken);
  useEffect(() => {
    if (state.resyncToken !== lastToken.current) {
      lastToken.current = state.resyncToken;
      refresh();
    }
  }, [state.resyncToken, refresh]);

  const tracks = useSignal(LIVE_CHANNELS.tracks);
  const events = useSignal(LIVE_CHANNELS.events);
  const cluster = useSignal(LIVE_CHANNELS.cluster);
  const status = useSignal(LIVE_CHANNELS.status);

  useEffect(() => {
    if (tracks.lastSignal)
      dispatch({
        type: 'tracks',
        message: tracks.lastSignal as unknown as TracksDelta,
        now: Date.now(),
      });
  }, [tracks.lastSignal]);
  useEffect(() => {
    if (events.lastSignal)
      dispatch({
        type: 'events',
        message: events.lastSignal as unknown as EventsDelta,
        now: Date.now(),
      });
  }, [events.lastSignal]);
  useEffect(() => {
    if (cluster.lastSignal)
      dispatch({
        type: 'cluster',
        message: cluster.lastSignal as unknown as ClusterDelta,
        now: Date.now(),
      });
  }, [cluster.lastSignal]);
  useEffect(() => {
    if (status.lastSignal)
      dispatch({
        type: 'status',
        message: status.lastSignal as unknown as StatusMessage,
        now: Date.now(),
      });
  }, [status.lastSignal]);

  return {
    state,
    load,
    freshness: freshness(state, now),
    now,
    signalsConnected: status.isSignalsAvailable,
    refresh,
  };
}
