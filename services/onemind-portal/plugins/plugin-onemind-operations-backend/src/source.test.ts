import { createSource, describeError, worstOf } from './source';

function clock(start = Date.parse('2026-10-22T12:00:00Z')) {
  let t = start;
  return { now: () => t, advance: (ms: number) => (t += ms) };
}

describe('createSource', () => {
  it('returns ok data and reuses it inside the cache window', async () => {
    const c = clock();
    const fetcher = jest.fn(async () => 'a');
    const source = createSource('pods', fetcher, { now: c.now, ttlMs: 15000 });
    const first = await source.read();
    c.advance(5000);
    await source.read();
    expect(fetcher).toHaveBeenCalledTimes(1);
    expect(first).toMatchObject({ status: 'ok', data: 'a', observedAt: '2026-10-22T12:00:00.000Z' });
    c.advance(11000);
    await source.read();
    expect(fetcher).toHaveBeenCalledTimes(2);
  });

  it('shares one in-flight fetch between concurrent readers', async () => {
    const fetcher = jest.fn(async () => 1);
    const source = createSource('pods', fetcher);
    await Promise.all([source.read(), source.read(), source.read()]);
    expect(fetcher).toHaveBeenCalledTimes(1);
  });

  it('reports error with no data when the first fetch fails', async () => {
    const source = createSource('prometheus', async () => {
      throw new Error('Upstream returned 503');
    });
    expect(await source.read()).toMatchObject({ status: 'error', observedAt: null, error: 'Upstream returned 503' });
  });

  it('times out a slow source and aborts it', async () => {
    let aborted = false;
    const source = createSource(
      'prometheus',
      signal =>
        new Promise<string>(() => {
          signal.addEventListener('abort', () => (aborted = true));
        }),
      { timeoutMs: 20 },
    );
    const result = await source.read();
    expect(result).toMatchObject({ status: 'error', error: 'Timed out' });
    expect(aborted).toBe(true);
  });

  it('serves the last good data as stale, with its original time, then gives up', async () => {
    const c = clock();
    let fail = false;
    const source = createSource(
      'argocd',
      async () => {
        if (fail) throw new Error('Cluster API returned 500');
        return 'good';
      },
      { now: c.now, ttlMs: 1000, staleMaxMs: 60000 },
    );
    await source.read();
    fail = true;
    c.advance(5000);
    expect(await source.read()).toMatchObject({
      status: 'stale',
      data: 'good',
      observedAt: '2026-10-22T12:00:00.000Z',
      error: 'Cluster API returned 500',
    });
    c.advance(120000);
    expect(await source.read()).toMatchObject({ status: 'error', observedAt: null });
  });
});

describe('describeError', () => {
  it('never echoes unknown messages to the browser', () => {
    expect(describeError(new Error('connect ECONNREFUSED 10.0.0.1:9090'))).toBe('Unavailable');
    expect(describeError('boom')).toBe('Unavailable');
  });
});

describe('worstOf', () => {
  it('takes the worst status and the oldest observation time', () => {
    const base = { attemptedAt: 'x' };
    expect(
      worstOf([
        { id: 'pods', status: 'ok', observedAt: '2026-10-22T12:00:10Z', ...base },
        { id: 'nodes', status: 'stale', observedAt: '2026-10-22T11:59:00Z', ...base },
      ]),
    ).toEqual({ status: 'stale', observedAt: '2026-10-22T11:59:00Z', sources: ['pods', 'nodes'] });
  });
});
