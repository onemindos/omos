import { collect, notConnected } from './collect';
import { diffSnapshots } from './changes';
import { DigitalOceanClient, readToken } from './client';
import type { ChangeSet, Snapshots, Summary, Volume } from './types';

export const CACHE_MS = 10 * 60_000;
const MAX_CHANGE_SETS = 50;

export interface ServiceOptions {
  /** Reads the token on every refresh so a rotated Secret is used without a restart. */
  getToken?: () => string | undefined;
  fetchImpl?: typeof fetch;
  sleep?: (ms: number) => Promise<void>;
  now?: () => number;
  cacheMs?: number;
  warn?: (message: string) => void;
}

/** Adds PVC and app links to volumes. May throw: the volumes are then kept without links. */
export type VolumeLinker = (volumes: Volume[]) => Promise<Volume[]>;

/**
 * The one place that talks to DigitalOcean. It caches the last snapshot for ten minutes,
 * remembers the one before it for the "What changed" list, and shares one in-flight refresh
 * between the summary route and the catalog provider.
 */
export class DigitalOceanService {
  private current?: Snapshots;
  private fetchedAt = 0;
  private pending?: Promise<Summary>;
  private changes: ChangeSet[] = [];
  private linker?: VolumeLinker;

  constructor(private readonly options: ServiceOptions = {}) {}

  setLinker(linker: VolumeLinker) {
    this.linker = linker;
  }

  private now() {
    return (this.options.now ?? Date.now)();
  }
  private ttl() {
    return this.options.cacheMs ?? CACHE_MS;
  }
  private token() {
    return (this.options.getToken ?? (() => readToken()))();
  }

  /** The cached summary, refreshed when older than the cache time. */
  async get(): Promise<Summary> {
    if (!this.token()) return this.summary(false);
    if (this.current && this.now() - this.fetchedAt < this.ttl())
      return this.summary(true);
    return this.refresh();
  }

  /** Reads DigitalOcean now (joins a refresh already running). */
  refresh(): Promise<Summary> {
    const token = this.token();
    if (!token) return Promise.resolve(this.summary(false));
    this.pending ??= this.read(token).finally(() => {
      this.pending = undefined;
    });
    return this.pending;
  }

  private async read(token: string): Promise<Summary> {
    const now = this.now();
    const client = new DigitalOceanClient({
      token,
      fetchImpl: this.options.fetchImpl,
      sleep: this.options.sleep,
      now: this.options.now,
    });
    const next = await collect(client, now);
    if (this.linker && next.volumes.state === 'ok') {
      try {
        next.volumes.items = await this.linker(next.volumes.items);
      } catch {
        this.options.warn?.(
          'DigitalOcean volumes were not linked to Kubernetes claims (cluster or catalog read failed)',
        );
      }
    }
    if (this.current) {
      const change = diffSnapshots(this.current, next, next.generatedAt);
      if (change)
        this.changes = [change, ...this.changes].slice(0, MAX_CHANGE_SETS);
    }
    this.current = next;
    this.fetchedAt = now;
    return this.summary(true);
  }

  private summary(connected: boolean): Summary {
    const data =
      connected && this.current ? this.current : notConnected(this.now());
    return {
      ...data,
      state: connected ? 'connected' : 'not-connected',
      expiresAt:
        connected && this.current
          ? new Date(this.fetchedAt + this.ttl()).toISOString()
          : null,
      changes: connected ? this.changes : [],
    };
  }
}

let shared: DigitalOceanService | undefined;
/** One service per process, shared by the summary route and the catalog provider. */
export function sharedService(warn?: (message: string) => void) {
  shared ??= new DigitalOceanService({ warn });
  return shared;
}
