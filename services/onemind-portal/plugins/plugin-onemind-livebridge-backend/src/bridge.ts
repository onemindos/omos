import type { LoggerService } from '@backstage/backend-plugin-api';
import { ClusterState } from './cluster';
import { parseCot } from './cot';
import type {
  BusEvent,
  ClusterDelta,
  EventsDelta,
  LiveSnapshot,
  LiveStatus,
  SourceStatus,
  StatusMessage,
  TracksDelta,
} from './contract';
import { LIVE_CHANNELS } from './contract';
import { splitDelta } from './delta';
import { RateMeter } from './rate';
import { EventSampler } from './sampler';
import { TrackStore } from './trackStore';

export type Publish = (
  channel: string,
  message: object,
) => void | Promise<void>;

/** Where the bridge reads its streams from; the plugin wires the real sources in. */
export interface StatusProviders {
  nats(): SourceStatus;
  kubernetes(): SourceStatus;
  argo(): SourceStatus;
}

const TEXT = new TextDecoder('utf-8', { fatal: false });
const MAX_SUMMARY = 160;
const MAX_JSON_BYTES = 16 * 1024;
const TITLE_KEYS = [
  'title',
  'summary',
  'message',
  'msg',
  'name',
  'event',
  'type',
  'id',
];

function clip(text: string, max = MAX_SUMMARY): string {
  // eslint-disable-next-line no-control-regex
  const clean = text.replace(/[\u0000-\u001f\u007f]+/g, ' ').trim();
  return clean.length > max ? `${clean.slice(0, max - 1)}…` : clean;
}

/** One short human line for a JSON bus message: a title-like field, else the first scalar fields. */
export function summarizeJson(text: string): string {
  try {
    const value = JSON.parse(text);
    if (value && typeof value === 'object' && !Array.isArray(value)) {
      const obj = value as Record<string, unknown>;
      for (const key of TITLE_KEYS) {
        const v = obj[key];
        if (typeof v === 'string' && v.trim()) return clip(v);
      }
      const parts: string[] = [];
      for (const [k, v] of Object.entries(obj)) {
        if (['string', 'number', 'boolean'].includes(typeof v))
          parts.push(`${k}=${String(v)}`);
        if (parts.length >= 3) break;
      }
      if (parts.length) return clip(parts.join(' '));
    }
  } catch {
    // not JSON
  }
  return clip(text, 80);
}

export interface LiveBridgeOptions {
  origin: string;
  logger: LoggerService;
  publish: Publish;
  cluster: ClusterState;
  providers: StatusProviders;
  now?: () => number;
}

/**
 * Holds the live picture (tracks, event ring, cluster state) and pushes deltas through `publish`:
 * tracks, events and cluster at most once per second, status every 5 s. Memory is bounded by the
 * caps: 5000 tracks, 200 events, 50 warnings, one slim record per pod.
 */
export class LiveBridge {
  readonly tracks = new TrackStore();
  readonly sampler = new EventSampler();
  readonly cluster: ClusterState;
  private readonly rate = new RateMeter();
  private seq = 0;
  private eventSeq = 0;
  private carriedDropped = 0;
  private tick = 0;
  private timer?: NodeJS.Timeout;
  private readonly now: () => number;
  private lastPublishWarn = 0;

  constructor(private readonly opts: LiveBridgeOptions) {
    this.now = opts.now ?? Date.now;
    this.cluster = opts.cluster;
  }

  get origin(): string {
    return this.opts.origin;
  }

  /** Entry point for every NATS message. */
  ingest(subject: string, data: Uint8Array): void {
    const nowMs = this.now();
    this.rate.hit(nowMs);
    const kind = subject.split('.', 1)[0];
    let cotSummary: string | undefined;
    if (subject.startsWith('tak.cot.')) {
      const track = parseCot(TEXT.decode(data), subject, new Date(nowMs));
      if (track) {
        this.tracks.upsert(track, nowMs);
        cotSummary = `${track.callsign ?? track.uid} ${track.type}`;
      }
    }
    if (!this.sampler.admit(subject, nowMs)) return;
    const summary =
      cotSummary ??
      (data.byteLength <= MAX_JSON_BYTES
        ? summarizeJson(TEXT.decode(data))
        : `(${data.byteLength} bytes)`);
    const event: BusEvent = {
      id: `${this.opts.origin}-e${++this.eventSeq}`,
      subject: clip(subject, 200),
      kind,
      ts: new Date(nowMs).toISOString(),
      bytes: data.byteLength,
      summary,
    };
    this.sampler.record(event);
  }

  start(): void {
    this.timer = setInterval(() => this.flush(), 1000);
    this.timer.unref?.();
  }

  stop(): void {
    if (this.timer) clearInterval(this.timer);
    this.timer = undefined;
  }

  status(): LiveStatus {
    return {
      nats: this.opts.providers.nats(),
      kubernetes: this.opts.providers.kubernetes(),
      argo: this.opts.providers.argo(),
    };
  }

  /** Total bus messages per second over the last 10 s (every message counts, sampled or not). */
  busRate(): number {
    return this.rate.rate(this.now());
  }

  snapshot(): LiveSnapshot {
    const nowMs = this.now();
    return {
      origin: this.opts.origin,
      seq: this.seq,
      generatedAt: new Date(nowMs).toISOString(),
      status: this.status(),
      tracks: this.tracks.all(),
      events: this.sampler.recent(),
      cluster: this.cluster.cells(nowMs),
      argo: this.cluster.argoCells(),
      warnings: this.cluster.warningList(),
    };
  }

  private send(channel: string, message: object): void {
    try {
      const result = this.opts.publish(channel, message);
      if (result && typeof (result as Promise<void>).catch === 'function') {
        (result as Promise<void>).catch(e => this.warnPublish(e));
      }
    } catch (e) {
      this.warnPublish(e);
    }
  }

  private warnPublish(e: unknown): void {
    const t = this.now();
    if (t - this.lastPublishWarn > 60_000) {
      this.lastPublishWarn = t;
      this.opts.logger.warn(
        `livebridge signal publish failed: ${(e as Error)?.message ?? e}`,
      );
    }
  }

  private envelope() {
    return {
      origin: this.opts.origin,
      seq: 0,
      ts: new Date(this.now()).toISOString(),
    };
  }

  /** Build and send whatever changed. Called once a second; exposed for tests. */
  flush(): void {
    const nowMs = this.now();
    this.tick++;
    if (this.tick % 10 === 0) this.tracks.sweep(nowMs);

    const t = this.tracks.drain();
    for (const m of splitDelta(this.envelope(), {
      upserts: t.upserts,
      removes: t.removes,
    })) {
      this.send(LIVE_CHANNELS.tracks, { ...m, seq: ++this.seq } as TracksDelta);
    }

    const e = this.sampler.drain();
    const dropped = this.carriedDropped + e.dropped;
    if (e.events.length) {
      this.carriedDropped = 0;
      const msgs = splitDelta(
        { ...this.envelope(), dropped },
        { events: e.events },
      );
      msgs.forEach((m, i) =>
        this.send(LIVE_CHANNELS.events, {
          ...m,
          dropped: i === 0 ? dropped : 0,
          seq: ++this.seq,
        } as EventsDelta),
      );
    } else {
      this.carriedDropped = dropped;
    }

    const c = this.cluster.delta(nowMs);
    if (c) {
      const msgs = splitDelta(this.envelope(), {
        namespaces: c.namespaces,
        removedNamespaces: c.removedNamespaces,
        argo: c.argo,
        warnings: c.warnings,
      });
      for (const m of msgs)
        this.send(LIVE_CHANNELS.cluster, {
          ...m,
          seq: ++this.seq,
        } as ClusterDelta);
    }

    if (this.tick % 5 === 0) {
      const msg: StatusMessage = {
        ...this.envelope(),
        seq: ++this.seq,
        status: this.status(),
      };
      this.send(LIVE_CHANNELS.status, msg);
    }
  }
}
