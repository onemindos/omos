import { DigitalOceanClient, DoApiError, retryDelayMs } from './client';
import { TOKEN } from './testing';

const json = (body: unknown, init: ResponseInit = {}) =>
  new Response(JSON.stringify(body), { status: 200, ...init });

describe('DigitalOceanClient', () => {
  it('sends GET only, with the bearer token, per_page=200, redirect:error and a 30 second timeout', async () => {
    const fetchImpl = jest.fn(async () =>
      json({ vpcs: [], meta: { total: 0 } }),
    );
    const client = new DigitalOceanClient({
      token: TOKEN,
      fetchImpl: fetchImpl as never,
    });
    await client.list('/v2/vpcs', 'vpcs');
    const [url, init] = fetchImpl.mock.calls[0] as unknown as [
      string,
      RequestInit,
    ];
    expect(url).toBe('https://api.digitalocean.com/v2/vpcs?per_page=200');
    expect(init.method).toBe('GET');
    expect(init.redirect).toBe('error');
    expect((init.headers as Record<string, string>).Authorization).toBe(
      `Bearer ${TOKEN}`,
    );
    expect(init.signal).toBeInstanceOf(AbortSignal);
  });

  it('follows links.pages.next on the API host and stops at meta.total', async () => {
    const pages = [
      json({
        droplets: [{ id: 1 }, { id: 2 }],
        links: {
          pages: {
            next: 'https://api.digitalocean.com/v2/droplets?page=2&per_page=200',
          },
        },
        meta: { total: 3 },
      }),
      json({ droplets: [{ id: 3 }], links: {}, meta: { total: 3 } }),
    ];
    const fetchImpl = jest.fn(async () => pages.shift()!);
    const client = new DigitalOceanClient({
      token: TOKEN,
      fetchImpl: fetchImpl as never,
    });
    const items = await client.list<{ id: number }>('/v2/droplets', 'droplets');
    expect(items.map(i => i.id)).toEqual([1, 2, 3]);
    expect(fetchImpl).toHaveBeenCalledTimes(2);
    expect((fetchImpl.mock.calls[1] as unknown as [string])[0]).toBe(
      'https://api.digitalocean.com/v2/droplets?page=2&per_page=200',
    );
  });

  it('never follows a next link that points at another host', async () => {
    const fetchImpl = jest.fn(async () =>
      json({
        droplets: [{ id: 1 }],
        links: { pages: { next: 'https://evil.example/v2/droplets?page=2' } },
      }),
    );
    const client = new DigitalOceanClient({
      token: TOKEN,
      fetchImpl: fetchImpl as never,
    });
    expect(await client.list('/v2/droplets', 'droplets')).toHaveLength(1);
    expect(fetchImpl).toHaveBeenCalledTimes(1);
  });

  it('waits for ratelimit-reset (epoch seconds) after a 429, then succeeds', async () => {
    const now = Date.parse('2026-11-05T00:00:00Z');
    const reset = Math.floor(now / 1000) + 7;
    const answers = [
      new Response('{}', {
        status: 429,
        headers: { 'ratelimit-reset': String(reset) },
      }),
      json({ vpcs: [{ id: 'a' }], meta: { total: 1 } }),
    ];
    const sleep = jest.fn(async () => {});
    const client = new DigitalOceanClient({
      token: TOKEN,
      fetchImpl: (async () => answers.shift()!) as never,
      sleep,
      now: () => now,
    });
    expect(await client.list('/v2/vpcs', 'vpcs')).toHaveLength(1);
    expect(sleep).toHaveBeenCalledWith(7000);
  });

  it('prefers retry-after, and gives up after a few 429s', async () => {
    expect(retryDelayMs(new Headers({ 'retry-after': '3' }), 0)).toBe(3000);
    expect(
      retryDelayMs(
        new Headers({ 'retry-after': '3', 'ratelimit-reset': '9999999999' }),
        0,
      ),
    ).toBe(3000);
    expect(retryDelayMs(new Headers(), 0)).toBe(5000);
    const sleep = jest.fn(async () => {});
    const client = new DigitalOceanClient({
      token: TOKEN,
      fetchImpl: (async () =>
        new Response('{}', {
          status: 429,
          headers: { 'retry-after': '1' },
        })) as never,
      sleep,
    });
    await expect(client.get('/v2/account')).rejects.toMatchObject({
      status: 429,
    });
    expect(sleep).toHaveBeenCalledTimes(3);
  });

  it('refuses a response above the size limit', async () => {
    const big = new Response('x'.repeat(17 * 1024 * 1024), { status: 200 });
    const client = new DigitalOceanClient({
      token: TOKEN,
      fetchImpl: (async () => big) as never,
    });
    await expect(client.get('/v2/account')).rejects.toBeInstanceOf(DoApiError);
  });

  it('keeps the token out of every error it raises', async () => {
    const failing = new DigitalOceanClient({
      token: TOKEN,
      fetchImpl: (async () => {
        throw new Error(`connect failed with Authorization: Bearer ${TOKEN}`);
      }) as never,
    });
    const rejected = await failing.get('/v2/account').catch(e => e);
    expect(rejected).toBeInstanceOf(DoApiError);
    expect(
      JSON.stringify([rejected.message, rejected.path, rejected.status]),
    ).not.toContain(TOKEN);
    const denied = new DigitalOceanClient({
      token: TOKEN,
      fetchImpl: (async () =>
        new Response(`{"message":"bad ${TOKEN}"}`, { status: 401 })) as never,
    });
    const error = await denied.get('/v2/account').catch(e => e);
    expect(error.status).toBe(401);
    expect(error.message).not.toContain(TOKEN);
  });
});
