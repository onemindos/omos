/**
 * Read-only DigitalOcean API v2 client. GET only: there is no method here that can create,
 * change or delete a resource. The token is read once from DIGITALOCEAN_TOKEN, goes only into
 * the Authorization header, and is never part of an error message, a log line or a response.
 */
export const DO_API = 'https://api.digitalocean.com';
export const TOKEN_ENV = 'DIGITALOCEAN_TOKEN';

const TIMEOUT_MS = 30_000;
const MAX_BODY_BYTES = 16 * 1024 * 1024;
const PER_PAGE = 200;
const MAX_PAGES = 50;
const MAX_RETRIES = 3;
const MAX_WAIT_MS = 60_000;

/** A request the API refused or could not answer. `status` 0 means a network failure or timeout. */
export class DoApiError extends Error {
  constructor(readonly status: number, readonly path: string, message: string) {
    super(message);
    this.name = 'DoApiError';
  }
  /** The token lacks the scope for this endpoint (or the product is not enabled). */
  get isScopeProblem() {
    return this.status === 401 || this.status === 403 || this.status === 404;
  }
}

export interface ClientOptions {
  token: string;
  fetchImpl?: typeof fetch;
  /** Test seam; waits between 429 retries. */
  sleep?: (ms: number) => Promise<void>;
  now?: () => number;
  baseUrl?: string;
}

export function readToken(env: NodeJS.ProcessEnv = process.env) {
  const token = env[TOKEN_ENV]?.trim();
  return token || undefined;
}

const clamp = (ms: number) => Math.min(Math.max(ms, 1_000), MAX_WAIT_MS);

/** Milliseconds to wait after a 429, from Retry-After (seconds or a date) or RateLimit-Reset (epoch seconds). */
export function retryDelayMs(headers: Headers, now: number): number {
  const retryAfter = headers.get('retry-after');
  if (retryAfter) {
    const seconds = Number(retryAfter);
    if (Number.isFinite(seconds)) return clamp(seconds * 1000);
    const date = Date.parse(retryAfter);
    if (Number.isFinite(date)) return clamp(date - now);
  }
  const reset = Number(headers.get('ratelimit-reset'));
  if (Number.isFinite(reset) && reset > 0) {
    // DigitalOcean sends epoch seconds; accept a small relative value too.
    return clamp(reset > 1e9 ? reset * 1000 - now : reset * 1000);
  }
  return 5_000;
}

async function readBounded(response: Response): Promise<string> {
  if (!response.body) return '';
  const reader = response.body.getReader();
  const chunks: Uint8Array[] = [];
  let size = 0;
  try {
    for (;;) {
      const part = await reader.read();
      if (part.done) break;
      size += part.value.byteLength;
      if (size > MAX_BODY_BYTES) throw new Error('response too large');
      chunks.push(part.value);
    }
  } finally {
    await reader.cancel().catch(() => {});
  }
  return Buffer.concat(chunks).toString('utf8');
}

export class DigitalOceanClient {
  private readonly fetchImpl: typeof fetch;
  private readonly sleep: (ms: number) => Promise<void>;
  private readonly now: () => number;
  private readonly baseUrl: string;
  private readonly token: string;

  constructor(options: ClientOptions) {
    this.token = options.token;
    this.fetchImpl = options.fetchImpl ?? fetch;
    this.sleep =
      options.sleep ?? (ms => new Promise(resolve => setTimeout(resolve, ms)));
    this.now = options.now ?? Date.now;
    this.baseUrl = options.baseUrl ?? DO_API;
  }

  /** GET one JSON document. Retries a 429 after the time the API names, a bounded number of times. */
  async get<T = Record<string, unknown>>(path: string): Promise<T> {
    for (let attempt = 0; ; attempt++) {
      let response: Response;
      try {
        response = await this.fetchImpl(`${this.baseUrl}${path}`, {
          method: 'GET',
          headers: {
            Authorization: `Bearer ${this.token}`,
            Accept: 'application/json',
          },
          redirect: 'error',
          signal: AbortSignal.timeout(TIMEOUT_MS),
        });
      } catch {
        // The underlying error can echo request details; keep the message generic.
        throw new DoApiError(0, path, 'DigitalOcean request failed');
      }
      if (response.status === 429) {
        const wait = retryDelayMs(response.headers, this.now());
        await response.body?.cancel().catch(() => {});
        if (attempt >= MAX_RETRIES)
          throw new DoApiError(429, path, 'DigitalOcean rate limit');
        await this.sleep(wait);
        continue;
      }
      if (!response.ok) {
        await response.body?.cancel().catch(() => {});
        throw new DoApiError(
          response.status,
          path,
          `DigitalOcean answered ${response.status}`,
        );
      }
      try {
        return JSON.parse(await readBounded(response)) as T;
      } catch {
        throw new DoApiError(
          response.status,
          path,
          'DigitalOcean response was not valid JSON or too large',
        );
      }
    }
  }

  /**
   * Every item of a paginated collection. Follows links.pages.next, and stops early when the
   * collected count reaches meta.total. `key` is the array name in the response body.
   */
  async list<T>(path: string, key: string): Promise<T[]> {
    const out: T[] = [];
    const sep = path.includes('?') ? '&' : '?';
    let next: string | undefined = `${path}${sep}per_page=${PER_PAGE}`;
    for (let page = 0; next && page < MAX_PAGES; page++) {
      const body: PagedBody = await this.get<PagedBody>(next);
      const items = body[key];
      if (!Array.isArray(items))
        throw new DoApiError(200, path, `Unexpected ${key} response`);
      out.push(...(items as T[]));
      const total = body.meta?.total;
      if (typeof total === 'number' && out.length >= total) break;
      next = nextPath(body.links?.pages?.next, this.baseUrl);
    }
    return out;
  }
}

interface PagedBody {
  [key: string]: unknown;
  links?: { pages?: { next?: string } };
  meta?: { total?: number };
}

/** The API returns absolute next links; only follow ones on the API host, as a path. */
function nextPath(link: string | undefined, baseUrl: string) {
  if (!link) return undefined;
  try {
    const url = new URL(link, baseUrl);
    if (url.origin !== new URL(baseUrl).origin) return undefined;
    return `${url.pathname}${url.search}`;
  } catch {
    return undefined;
  }
}
