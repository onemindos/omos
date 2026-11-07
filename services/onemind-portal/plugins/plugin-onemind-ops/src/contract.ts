/**
 * Live bridge contract, shared by the backend plugin (livebridge) and the Ops Center page.
 * Copy this file verbatim into both packages as src/contract.ts. Change it in both places together.
 *
 * Transport:
 *  - GET /api/livebridge/snapshot  -> LiveSnapshot (signed-in users only)
 *  - Signals channels below push deltas. Every message carries `origin` (backend replica id) and a
 *    per-origin increasing `seq`. The portal runs 2 replicas, so a browser may hear from either or
 *    both: the client locks onto one origin and switches only after that origin is silent for 15 s.
 */

export const LIVE_CHANNELS = {
  tracks: 'onemind:live:tracks',
  events: 'onemind:live:events',
  cluster: 'onemind:live:cluster',
  status: 'onemind:live:status',
} as const;

export type ConnectionState = 'connected' | 'connecting' | 'down' | 'disabled';

export interface SourceStatus {
  state: ConnectionState;
  /** ISO time of the last message or watch event received from this source. */
  lastMessageAt?: string;
  /** Messages per second over the last 10 s. */
  rate?: number;
  detail?: string;
}

export interface LiveStatus {
  nats: SourceStatus;
  kubernetes: SourceStatus;
  argo: SourceStatus;
}

/** One position from CoT on tak.cot.> (TAK EUDs plus the aircraft, ship, train and hazard feeds). */
export interface Track {
  uid: string;
  /** CoT type, e.g. a-f-G-U-C, a-n-A-C-F. Affiliation is the second atom (f, h, n, u, ...). */
  type: string;
  callsign?: string;
  lat: number;
  lon: number;
  /** Height above ellipsoid, metres. 9999999 in CoT means unknown and is sent as undefined. */
  hae?: number;
  course?: number;
  /** Metres per second. */
  speed?: number;
  /** ISO times from the CoT event. */
  time: string;
  stale: string;
  /** Subject the position arrived on. */
  subject: string;
}

export interface BusEvent {
  /** origin + seq, unique per message. */
  id: string;
  subject: string;
  /** First subject token: tel, det, evt, cmd, intel, ent, tak, agents, ... */
  kind: string;
  ts: string;
  bytes: number;
  /** Short human summary (callsign and type for CoT, a title or the first fields for JSON). */
  summary: string;
}

export type CellHealth = 'ok' | 'warn' | 'down';

export interface NamespaceCell {
  namespace: string;
  total: number;
  ready: number;
  pending: number;
  failed: number;
  /** Restarts in the last hour, summed over containers. */
  restarts: number;
  health: CellHealth;
  /** Most recent change, e.g. "pod forgejo-0 CrashLoopBackOff". */
  lastChange?: string;
  lastChangeAt?: string;
}

export interface ArgoCell {
  name: string;
  sync: string;
  health: string;
  revision?: string;
  changedAt?: string;
}

export interface WarningEvent {
  namespace: string;
  object: string;
  reason: string;
  message: string;
  ts: string;
  count?: number;
}

export interface LiveSnapshot {
  origin: string;
  seq: number;
  generatedAt: string;
  status: LiveStatus;
  tracks: Track[];
  /** Newest first, at most 200. */
  events: BusEvent[];
  cluster: NamespaceCell[];
  argo: ArgoCell[];
  /** Kubernetes Warning events, newest first, at most 50. */
  warnings: WarningEvent[];
}

interface Envelope {
  origin: string;
  seq: number;
  ts: string;
}

/** At most one per second. `removes` lists uids that went stale or aged out. */
export interface TracksDelta extends Envelope {
  upserts: Track[];
  removes: string[];
}

/** At most one per second, at most 50 events, newest last. `dropped` counts events not sent. */
export interface EventsDelta extends Envelope {
  events: BusEvent[];
  dropped: number;
}

/** Sent when anything changes, at most one per second. Only changed cells are included. */
export interface ClusterDelta extends Envelope {
  namespaces: NamespaceCell[];
  removedNamespaces: string[];
  argo: ArgoCell[];
  warnings: WarningEvent[];
}

/** Every 5 s. */
export interface StatusMessage extends Envelope {
  status: LiveStatus;
}
