// Every upstream the aggregator reads is wrapped in a Source: a per-source timeout,
// a short cache, and an honest status. A failing source never fails the others.

export type SourceId =
  | 'alertmanager'
  | 'prometheus'
  | 'argocd'
  | 'deployments'
  | 'statefulsets'
  | 'pods'
  | 'nodes'
  | 'events'
  | 'cronjobs'
  | 'uptime';

/** ok: fetched just now (or within the cache window). stale: the fetch failed, last good data served. error: nothing to show. */
export type SourceStatus = 'ok' | 'stale' | 'error';

export type SourceResult<T> = {
  id: SourceId;
  status: SourceStatus;
  /** When `data` was actually read from the upstream. Null for `error`. */
  observedAt: string | null;
  attemptedAt: string;
  error?: string;
  data?: T;
};

export type Source<T> = { id: SourceId; read(): Promise<SourceResult<T>> };

export type SourceOptions = {
  timeoutMs?: number;
  /** Results younger than this are reused (the 15 s aggregator cache). */
  ttlMs?: number;
  /** After a failure, last good data older than this is no longer served. */
  staleMaxMs?: number;
  now?: () => number;
};

export const DEFAULT_TTL_MS = 15000;
export const DEFAULT_STALE_MAX_MS = 5 * 60000;

const iso = (ms: number) => new Date(ms).toISOString();

// Short, fixed messages: never echo hostnames or response bodies to the browser.
export function describeError(error: unknown): string {
  if (error instanceof Error) {
    if (/timed out|timeout|aborted/i.test(error.message) || error.name === 'AbortError' || error.name === 'TimeoutError') {
      return 'Timed out';
    }
    if (/^Cluster API returned \d{3}$/.test(error.message)) return error.message;
    if (/^Upstream returned \d{3}$/.test(error.message)) return error.message;
    if (/credentials unavailable/i.test(error.message)) return 'Cluster credentials unavailable';
  }
  return 'Unavailable';
}

export function createSource<T>(
  id: SourceId,
  fetcher: (signal: AbortSignal) => Promise<T>,
  options: SourceOptions = {},
): Source<T> {
  const timeoutMs = options.timeoutMs ?? 8000;
  const ttlMs = options.ttlMs ?? DEFAULT_TTL_MS;
  const staleMaxMs = options.staleMaxMs ?? DEFAULT_STALE_MAX_MS;
  const now = options.now ?? Date.now;
  let last: SourceResult<T> | undefined;
  let lastAt = 0;
  let good: { data: T; at: number } | undefined;
  let pending: Promise<SourceResult<T>> | undefined;

  async function refresh(): Promise<SourceResult<T>> {
    const started = now();
    const controller = new AbortController();
    let timer: ReturnType<typeof setTimeout> | undefined;
    try {
      const timeout = new Promise<never>((_, reject) => {
        timer = setTimeout(() => {
          controller.abort();
          reject(new Error('timed out'));
        }, timeoutMs);
      });
      const data = await Promise.race([fetcher(controller.signal), timeout]);
      good = { data, at: now() };
      return { id, status: 'ok', observedAt: iso(good.at), attemptedAt: iso(started), data };
    } catch (error) {
      const message = describeError(error);
      if (good && now() - good.at <= staleMaxMs) {
        return { id, status: 'stale', observedAt: iso(good.at), attemptedAt: iso(started), error: message, data: good.data };
      }
      return { id, status: 'error', observedAt: null, attemptedAt: iso(started), error: message };
    } finally {
      if (timer) clearTimeout(timer);
    }
  }

  return {
    id,
    read() {
      if (last && now() - lastAt < ttlMs) return Promise.resolve(last);
      if (!pending) {
        pending = refresh()
          .then(result => {
            last = result;
            lastAt = now();
            return result;
          })
          .finally(() => {
            pending = undefined;
          });
      }
      return pending;
    },
  };
}

const rank: Record<SourceStatus, number> = { ok: 0, stale: 1, error: 2 };

/** The panel shows the worst of the sources it depends on, with the oldest observation time. */
export function worstOf(results: SourceResult<unknown>[]) {
  const status = results.reduce<SourceStatus>(
    (worst, item) => (rank[item.status] > rank[worst] ? item.status : worst),
    'ok',
  );
  const times = results.map(item => item.observedAt).filter((t): t is string => !!t);
  const observedAt = times.length ? times.reduce((a, b) => (a < b ? a : b)) : null;
  return { status, observedAt, sources: results.map(item => item.id) };
}
