/**
 * Read-only client for the two Authentik listings the org sync needs (Authentik 2026.8 API):
 * GET /api/v3/core/users/ and GET /api/v3/core/groups/, page by page. Any unexpected answer throws,
 * so a partial or misread directory is never published.
 */

/** User types that are people. Service accounts (outposts, API clients) are never read or published. */
export const PERSON_TYPES = ['internal', 'external'];

/** The fields this module uses from one /core/users/ record. */
export interface AuthentikUser {
  pk: number;
  username: string;
  name: string;
  email: string;
  isActive: boolean;
  type: string;
  /**
   * sha256 hashed user id. It is the OIDC `sub` Authentik issues for providers with
   * sub_mode "hashed_user_id", which is how the onemind-backstage provider is created.
   */
  uid: string;
  uuid?: string;
  /** Direct group memberships, as group pks. */
  groups: string[];
}

/** The fields this module uses from one /core/groups/ record. */
export interface AuthentikGroup {
  pk: string;
  name: string;
  /** Direct parent group pks. Authentik allows several. */
  parents: string[];
}

export interface AuthentikSnapshot {
  users: AuthentikUser[];
  groups: AuthentikGroup[];
}

export type FetchJson = (url: string, signal?: AbortSignal) => Promise<unknown>;

/** Authentik's default maximum page size; the server may answer with smaller pages. */
const PAGE_SIZE = 100;
const MAX_PAGES = 1000;
const REQUEST_TIMEOUT_MS = 20_000;
const MAX_RESPONSE_BYTES = 10_000_000;

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null && !Array.isArray(value);

async function readBody(response: Response, path: string): Promise<string> {
  if (!response.body) return '';
  const reader = response.body.getReader();
  const chunks: Uint8Array[] = [];
  let size = 0;
  try {
    for (;;) {
      const part = await reader.read();
      if (part.done) break;
      size += part.value.byteLength;
      if (size > MAX_RESPONSE_BYTES)
        throw new Error(`Authentik API response for ${path} is too large`);
      chunks.push(part.value);
    }
  } finally {
    await reader.cancel().catch(() => {});
  }
  return Buffer.concat(chunks).toString('utf8');
}

/**
 * GET with the read-only token: no redirects (the token never follows one), a deadline per request,
 * a response size cap. Errors name the path and status only, never the token or the response body.
 */
export function authentikFetchJson(
  token: string,
  fetchImpl: typeof fetch = globalThis.fetch,
): FetchJson {
  return async (url, signal) => {
    const path = new URL(url).pathname;
    const deadline = AbortSignal.timeout(REQUEST_TIMEOUT_MS);
    const response = await fetchImpl(url, {
      headers: {
        Authorization: `Bearer ${token}`,
        Accept: 'application/json',
      },
      redirect: 'error',
      signal: signal ? AbortSignal.any([signal, deadline]) : deadline,
    });
    if (!response.ok) {
      await response.body?.cancel().catch(() => {});
      const hint =
        response.status === 401 || response.status === 403
          ? ' (check that authentik.token is valid and may view users and groups)'
          : '';
      throw new Error(
        `Authentik API answered ${response.status} for ${path}${hint}`,
      );
    }
    const text = await readBody(response, path);
    try {
      return JSON.parse(text);
    } catch {
      throw new Error(`Authentik API returned invalid JSON for ${path}`);
    }
  };
}

/**
 * Every record of one listing. Authentik pages are numbered: `pagination.next` is the next page
 * number, or 0 on the last page. Every page must report the same count and the de-duplicated
 * records must add up to it, so a listing that changed while it was read (which can shift a record
 * past the reader) fails and is retried on the next run instead of being published short.
 */
async function listAll(
  fetchJson: FetchJson,
  baseUrl: string,
  path: string,
  params: [string, string][],
  signal?: AbortSignal,
): Promise<Record<string, unknown>[]> {
  const byPk = new Map<string, Record<string, unknown>>();
  let page = 1;
  let count: number | undefined;
  for (let read = 1; ; read++) {
    if (read > MAX_PAGES)
      throw new Error(
        `Authentik ${path} has more than ${MAX_PAGES} pages; refusing a partial read`,
      );
    const query = new URLSearchParams([
      ...params,
      ['page', String(page)],
      ['page_size', String(PAGE_SIZE)],
    ]);
    const body = await fetchJson(`${baseUrl}${path}?${query}`, signal);
    if (
      !isRecord(body) ||
      !Array.isArray(body.results) ||
      !isRecord(body.pagination) ||
      typeof body.pagination.next !== 'number' ||
      typeof body.pagination.count !== 'number'
    )
      throw new Error(`Unexpected Authentik answer for ${path}`);
    for (const item of body.results) {
      if (
        !isRecord(item) ||
        (typeof item.pk !== 'string' && typeof item.pk !== 'number')
      )
        throw new Error(`Unexpected record in Authentik ${path}`);
      byPk.set(String(item.pk), item);
    }
    if (count === undefined) count = body.pagination.count;
    else if (body.pagination.count !== count)
      throw new Error(
        `Authentik ${path} changed while it was read; the next run retries`,
      );
    const next = body.pagination.next;
    if (next === 0) break;
    if (!Number.isInteger(next) || next <= page)
      throw new Error(
        `Authentik ${path} pagination did not advance past page ${page}`,
      );
    page = next;
  }
  if (byPk.size !== count)
    throw new Error(
      `Authentik ${path} returned ${byPk.size} of ${count} records; the next run retries`,
    );
  return [...byPk.values()];
}

const isStringArray = (value: unknown): value is string[] =>
  Array.isArray(value) && value.every(v => typeof v === 'string');

export function parseUser(raw: Record<string, unknown>): AuthentikUser {
  const { pk, username, name, email, is_active, type, uid, uuid, groups } = raw;
  if (
    typeof pk !== 'number' ||
    !Number.isInteger(pk) ||
    typeof username !== 'string' ||
    !username.trim() ||
    typeof uid !== 'string' ||
    !/^\S+$/.test(uid) ||
    typeof is_active !== 'boolean' ||
    typeof type !== 'string' ||
    !isStringArray(groups)
  )
    throw new Error(`Unexpected user record from Authentik (pk ${String(pk)})`);
  return {
    pk,
    username,
    name: typeof name === 'string' ? name : '',
    email: typeof email === 'string' ? email : '',
    isActive: is_active,
    type,
    uid,
    uuid: typeof uuid === 'string' && uuid ? uuid : undefined,
    groups,
  };
}

export function parseGroup(raw: Record<string, unknown>): AuthentikGroup {
  const { pk, name, parents } = raw;
  if (
    typeof pk !== 'string' ||
    !pk ||
    typeof name !== 'string' ||
    !name.trim() ||
    !isStringArray(parents)
  )
    throw new Error(
      `Unexpected group record from Authentik (pk ${String(pk)})`,
    );
  return { pk, name, parents };
}

/**
 * All active people and all groups. The server-side filters keep pages small; the mapping filters
 * again, so a server that ignored a filter still cannot publish a service account.
 */
export async function readAuthentikOrg(
  fetchJson: FetchJson,
  baseUrl: string,
  signal?: AbortSignal,
): Promise<AuthentikSnapshot> {
  const users = await listAll(
    fetchJson,
    baseUrl,
    '/api/v3/core/users/',
    [
      ['include_groups', 'false'],
      ['include_roles', 'false'],
      ['is_active', 'true'],
      ...PERSON_TYPES.map((t): [string, string] => ['type', t]),
    ],
    signal,
  );
  const groups = await listAll(
    fetchJson,
    baseUrl,
    '/api/v3/core/groups/',
    [['include_users', 'false']],
    signal,
  );
  return { users: users.map(parseUser), groups: groups.map(parseGroup) };
}
