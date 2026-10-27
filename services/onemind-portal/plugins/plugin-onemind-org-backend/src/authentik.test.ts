import {
  authentikFetchJson,
  parseUser,
  readAuthentikOrg,
  type FetchJson,
} from './authentik';
import { apiGroup, apiUser, directory, fakeAuthentik, uidOf } from './testing';

const BASE = 'http://authentik.test';
const TOKEN = 'fixture-only-token';

describe('readAuthentikOrg', () => {
  it('reads every page of users and groups with the documented query', async () => {
    const api = fakeAuthentik(directory, 2);
    const org = await readAuthentikOrg(api.fetchJson, BASE);
    expect(org.users.map(u => u.username)).toEqual(
      directory.users.map(u => u.username),
    );
    expect(org.groups.map(g => g.name)).toEqual(
      directory.groups.map(g => g.name),
    );
    // 9 users and 7 groups at a server page size of 2: 5 + 4 requests.
    expect(api.urls).toHaveLength(9);
    const first = new URL(api.urls[0]);
    expect(first.origin + first.pathname).toBe(`${BASE}/api/v3/core/users/`);
    expect([...first.searchParams]).toEqual([
      ['include_groups', 'false'],
      ['include_roles', 'false'],
      ['is_active', 'true'],
      ['type', 'internal'],
      ['type', 'external'],
      ['page', '1'],
      ['page_size', '100'],
    ]);
    expect(new URL(api.urls[1]).searchParams.get('page')).toBe('2');
    const groups = new URL(api.urls[5]);
    expect(groups.pathname).toBe('/api/v3/core/groups/');
    expect(groups.searchParams.get('include_users')).toBe('false');
  });

  it('keeps the fields the mapping needs and the exact uid', async () => {
    const org = await readAuthentikOrg(
      fakeAuthentik(directory).fetchJson,
      BASE,
    );
    expect(org.users.find(u => u.pk === 9)).toEqual({
      pk: 9,
      username: 'external.partner',
      name: '',
      email: 'person9@example.test',
      isActive: true,
      type: 'external',
      uid: uidOf(9),
      uuid: '10000000-0000-4000-8000-000000000009',
      groups: ['00000000-0000-4000-8000-000000000002'],
    });
    expect(org.groups.find(g => g.name === 'Two Parents')?.parents).toEqual([
      '00000000-0000-4000-8000-000000000001',
      '00000000-0000-4000-8000-000000000004',
    ]);
  });

  const page = (results: unknown[], next: number, count: number) => ({
    pagination: { next, previous: 0, count, current: 1 },
    results,
  });

  it.each([
    ['a page without results', { pagination: { next: 0, count: 0 } }],
    ['a page without pagination', { results: [] }],
    ['an error document', { detail: 'Not found.' }],
    ['a list instead of a page', []],
  ])('rejects %s', async (_name, body) => {
    await expect(readAuthentikOrg(async () => body, BASE)).rejects.toThrow(
      'Unexpected Authentik answer for /api/v3/core/users/',
    );
  });

  it('rejects pagination that does not advance instead of looping', async () => {
    const fetchJson: FetchJson = async () => page([apiUser(1, 'a', [])], 1, 2);
    await expect(readAuthentikOrg(fetchJson, BASE)).rejects.toThrow(
      'pagination did not advance past page 1',
    );
  });

  it('rejects a listing that changed while it was read', async () => {
    // Page 1 lists a and b of 3 users; "a" is deleted before page 2, which shifts "c" past the
    // reader. Read naively this returns a and b, as many as the new count, and silently drops c.
    const fetchJson: FetchJson = async url =>
      new URL(url).searchParams.get('page') === '1'
        ? page([apiUser(1, 'a', []), apiUser(2, 'b', [])], 2, 3)
        : page([], 0, 2);
    await expect(readAuthentikOrg(fetchJson, BASE)).rejects.toThrow(
      'Authentik /api/v3/core/users/ changed while it was read; the next run retries',
    );
  });

  it('rejects records that do not add up to the reported count', async () => {
    const fetchJson: FetchJson = async () => page([apiUser(1, 'a', [])], 0, 2);
    await expect(readAuthentikOrg(fetchJson, BASE)).rejects.toThrow(
      'returned 1 of 2 records',
    );
  });

  it('caps the number of pages', async () => {
    let calls = 0;
    const fetchJson: FetchJson = async url => {
      calls++;
      const n = Number(new URL(url).searchParams.get('page'));
      return page([apiUser(n, `u${n}`, [])], n + 1, 1_000_000);
    };
    await expect(readAuthentikOrg(fetchJson, BASE)).rejects.toThrow(
      'more than 1000 pages',
    );
    expect(calls).toBe(1000);
  });

  it('de-duplicates a record listed twice and still checks the count', async () => {
    const fetchJson: FetchJson = async url => {
      const p = new URL(url).searchParams.get('page');
      if (new URL(url).pathname.endsWith('/groups/'))
        return page([apiGroup(1, 'g')], 0, 1);
      return p === '1'
        ? page([apiUser(1, 'a', []), apiUser(2, 'b', [])], 2, 3)
        : page([apiUser(2, 'b', []), apiUser(3, 'c', [])], 0, 3);
    };
    const org = await readAuthentikOrg(fetchJson, BASE);
    expect(org.users.map(u => u.pk)).toEqual([1, 2, 3]);
  });

  it.each([
    ['uid', { uid: undefined }],
    ['an empty uid', { uid: '' }],
    ['a uid with whitespace', { uid: 'abc def' }],
    ['groups', { groups: undefined }],
    ['is_active', { is_active: 'yes' }],
    ['username', { username: ' ' }],
  ])('rejects a user record without a valid %s', (_name, extra) => {
    expect(() => parseUser(apiUser(1, 'a', [], extra))).toThrow(
      'Unexpected user record from Authentik (pk 1)',
    );
  });

  it('rejects a group record from an older API without parents', async () => {
    const legacy = { ...apiGroup(1, 'g'), parents: undefined, parent: null };
    const fetchJson: FetchJson = async url =>
      new URL(url).pathname.endsWith('/groups/')
        ? page([legacy], 0, 1)
        : page([apiUser(1, 'a', [])], 0, 1);
    await expect(readAuthentikOrg(fetchJson, BASE)).rejects.toThrow(
      'Unexpected group record from Authentik',
    );
  });
});

describe('authentikFetchJson', () => {
  const respond = (body: BodyInit | null, status = 200) =>
    jest.fn<Promise<Response>, Parameters<typeof fetch>>(
      async () => new Response(body, { status }),
    );

  it('sends the token as a bearer header, refuses redirects and sets a deadline', async () => {
    const fetchImpl = respond(JSON.stringify({ ok: true }));
    const get = authentikFetchJson(TOKEN, fetchImpl as typeof fetch);
    await expect(get(`${BASE}/api/v3/core/users/?page=1`)).resolves.toEqual({
      ok: true,
    });
    const [url, init] = fetchImpl.mock.calls[0];
    expect(url).toBe(`${BASE}/api/v3/core/users/?page=1`);
    expect(init?.headers).toEqual({
      Authorization: `Bearer ${TOKEN}`,
      Accept: 'application/json',
    });
    expect(init?.redirect).toBe('error');
    expect(init?.signal).toBeInstanceOf(AbortSignal);
  });

  it('passes the scheduler abort signal through', async () => {
    const fetchImpl = respond('{}');
    const controller = new AbortController();
    await authentikFetchJson(TOKEN, fetchImpl as typeof fetch)(
      `${BASE}/x`,
      controller.signal,
    );
    const signal = fetchImpl.mock.calls[0][1]?.signal as AbortSignal;
    expect(signal.aborted).toBe(false);
    controller.abort();
    expect(signal.aborted).toBe(true);
  });

  it.each([401, 403])(
    'names the path and status, with a token hint, but never the token or body (%p)',
    async status => {
      const get = authentikFetchJson(
        TOKEN,
        respond(
          JSON.stringify({ detail: `bad ${TOKEN}` }),
          status,
        ) as typeof fetch,
      );
      const error = await get(`${BASE}/api/v3/core/users/?page=1`).catch(
        e => e as Error,
      );
      expect(error).toBeInstanceOf(Error);
      expect((error as Error).message).toBe(
        `Authentik API answered ${status} for /api/v3/core/users/ (check that authentik.token is valid and may view users and groups)`,
      );
      expect((error as Error).message).not.toContain(TOKEN);
    },
  );

  it('rejects invalid JSON and oversized responses', async () => {
    await expect(
      authentikFetchJson(
        TOKEN,
        respond('<html>') as typeof fetch,
      )(`${BASE}/api/v3/core/groups/`),
    ).rejects.toThrow('invalid JSON for /api/v3/core/groups/');
    const chunk = new Uint8Array(1_000_000);
    let sent = 0;
    const endless = new ReadableStream<Uint8Array>({
      pull(controller) {
        sent++;
        controller.enqueue(chunk);
      },
    });
    await expect(
      authentikFetchJson(
        TOKEN,
        respond(endless) as typeof fetch,
      )(`${BASE}/api/v3/core/users/`),
    ).rejects.toThrow('response for /api/v3/core/users/ is too large');
    expect(sent).toBeLessThan(20);
  });
});
