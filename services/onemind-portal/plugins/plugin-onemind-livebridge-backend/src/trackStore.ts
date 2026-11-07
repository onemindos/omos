import type { Track } from './contract';

export const MAX_TRACKS = 5000;
export const MIN_RETENTION_MS = 2 * 60_000;
export const MAX_RETENTION_MS = 30 * 60_000;

interface Entry {
  track: Track;
  expiresAt: number;
}

/**
 * Latest Track per uid. Capped at 5000 (the least recently updated is evicted). A track is kept
 * until its CoT stale time passes, but never less than 2 minutes and never more than 30 minutes
 * after it was last heard.
 */
export class TrackStore {
  private readonly map = new Map<string, Entry>();
  private upserts = new Map<string, Track>();
  private removes = new Set<string>();

  constructor(private readonly max = MAX_TRACKS) {}

  get size(): number {
    return this.map.size;
  }

  upsert(track: Track, nowMs: number): void {
    const staleMs = Date.parse(track.stale);
    const keep = Number.isFinite(staleMs) ? staleMs - nowMs : MIN_RETENTION_MS;
    const expiresAt =
      nowMs + Math.min(MAX_RETENTION_MS, Math.max(MIN_RETENTION_MS, keep));
    // Re-insert so Map order is least recently updated first.
    this.map.delete(track.uid);
    this.map.set(track.uid, { track, expiresAt });
    this.upserts.set(track.uid, track);
    this.removes.delete(track.uid);
    while (this.map.size > this.max) {
      const oldest = this.map.keys().next().value as string;
      this.remove(oldest);
    }
  }

  private remove(uid: string): void {
    this.map.delete(uid);
    this.upserts.delete(uid);
    this.removes.add(uid);
  }

  sweep(nowMs: number): void {
    for (const [uid, entry] of this.map) {
      if (entry.expiresAt <= nowMs) this.remove(uid);
    }
  }

  /** Changes since the last drain. */
  drain(): { upserts: Track[]; removes: string[] } {
    const out = {
      upserts: [...this.upserts.values()],
      removes: [...this.removes],
    };
    this.upserts = new Map();
    this.removes = new Set();
    return out;
  }

  all(): Track[] {
    return [...this.map.values()].map(e => e.track);
  }
}
