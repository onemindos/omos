import {
  parseNamespaceAllocation,
  specialAllocations,
  validateWindow,
} from './allocation';
import { allocationUrl, costWindow } from './costApi';
import fixture from './__fixtures__/allocation-7d.json';

// The canonical offline contract. Both implementations must agree on every input.
const reference = jest.requireActual('../../../scripts/cost-allocation.cjs');

const window = { start: '2026-10-29T00:00:00Z', end: '2026-11-05T00:00:00Z' };
const clone = <T>(value: T): T => JSON.parse(JSON.stringify(value));

describe('parity with scripts/cost-allocation.cjs', () => {
  const cases: Record<string, unknown> = {
    'real OpenCost payload': fixture,
    missing: undefined,
    'null body': null,
    'not an object': 'oops',
    'source error': { code: 500, data: null },
    'no data': { code: 200 },
    empty: { code: 200, data: [] },
    'two sets': { code: 200, data: [{}, {}] },
    'warnings reported': { ...clone(fixture), warnings: ['x'] },
    'unknown total': (() => {
      const payload: any = clone(fixture);
      delete payload.data[0].takserver.totalCost;
      return payload;
    })(),
    'negative total': (() => {
      const payload: any = clone(fixture);
      payload.data[0].nats.totalCost = -1;
      return payload;
    })(),
    'string total': (() => {
      const payload: any = clone(fixture);
      payload.data[0].nats.totalCost = '6.1';
      return payload;
    })(),
    'wrong window': (() => {
      const payload: any = clone(fixture);
      payload.data[0].nats.window.end = '2026-11-04T00:00:00Z';
      return payload;
    })(),
    'identity mismatch': (() => {
      const payload: any = clone(fixture);
      payload.data[0].nats.name = 'other';
      return payload;
    })(),
  };
  it.each(Object.entries(cases))('%s', (_name, payload) => {
    expect(parseNamespaceAllocation(payload, { window })).toEqual(
      reference.parseNamespaceAllocation(payload, { window }),
    );
  });
  it('validates windows identically', () => {
    expect(validateWindow(window)).toEqual(reference.validateWindow(window));
    for (const bad of [
      { start: window.end, end: window.start },
      { start: '2026-02-30T00:00:00Z', end: window.end },
      { start: '2026-10-29', end: window.end },
      null,
    ]) {
      expect(() => validateWindow(bad)).toThrow();
      expect(() => reference.validateWindow(bad)).toThrow();
    }
  });
});

describe('real payload semantics', () => {
  const result = parseNamespaceAllocation(fixture, { window });
  it('keeps special allocations out of the namespaces', () => {
    expect(result.rows.map(row => row.namespace)).not.toContain(
      '__unmounted__',
    );
    expect(result.state).toBe('partial');
    expect(result.issues).toContain('special-allocation-excluded');
  });
  it('reports a namespace that ran for only part of the window as such', () => {
    const nextcloud = result.rows.find(row => row.namespace === 'nextcloud');
    expect(nextcloud?.minutes).toBe(5760);
    const gateway = result.rows.find(row => row.namespace === 'tool-gateway');
    expect(gateway?.active.start).toBe('2026-10-30T17:05:00Z');
  });
  it('surfaces unattributed cost instead of dropping it', () => {
    const [unmounted] = specialAllocations(fixture);
    expect(unmounted.name).toBe('__unmounted__');
    expect(unmounted.totalCost).toBeGreaterThan(20);
  });
  it('never turns an unreadable special cost into zero', () => {
    const payload: any = clone(fixture);
    payload.data[0].__unmounted__.totalCost = 'n/a';
    expect(specialAllocations(payload)[0].totalCost).toBeNull();
    expect(specialAllocations({ code: 500 })).toEqual([]);
  });
});

describe('request window', () => {
  it('uses the last complete UTC days with explicit ends', () => {
    expect(costWindow(Date.parse('2026-11-05T13:45:10Z'), 7)).toEqual(window);
    expect(costWindow(Date.parse('2026-11-05T00:00:00Z'), 1)).toEqual({
      start: '2026-11-04T00:00:00Z',
      end: '2026-11-05T00:00:00Z',
    });
  });
  it('asks the existing proxy route for namespace aggregation', () => {
    const url = allocationUrl('http://x/api/proxy', window);
    expect(url.startsWith('http://x/api/proxy/opencost/allocation?')).toBe(
      true,
    );
    const query = new URL(url).searchParams;
    expect(query.get('aggregate')).toBe('namespace');
    expect(query.get('accumulate')).toBe('true');
    expect(query.get('window')).toBe(`${window.start},${window.end}`);
  });
});
