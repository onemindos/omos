'use strict';
const { test } = require('node:test');
const assert = require('node:assert/strict');
const { parseNamespaceAllocation: parse, validateWindow, observationState } = require('./cost-allocation.cjs');
// Synthetic fixture from pinned AllocationJSON field names; no live capture.
const window = { start: '2026-10-20T00:00:00Z', end: '2026-10-20T02:00:00Z' };
const options = { window, currency: 'USD' };
function row(name = 'backstage', patch = {}) {
  return { name, properties: { namespace: name, cluster: 'onemind-cloud' },
    window: { ...window }, start: window.start, end: window.end, minutes: 120,
    cpuCost: 2, gpuCost: 0, ramCost: 2, pvCost: 0, networkCost: 0,
    loadBalancerCost: 0, sharedCost: 0, externalCost: 0,
    cpuCostAdjustment: 0, gpuCostAdjustment: 0, ramCostAdjustment: 0,
    pvCostAdjustment: 0, networkCostAdjustment: 0, loadBalancerCostAdjustment: 0,
    totalCost: 4, ...patch };
}
const envelope = (value = row()) => ({ code: 200, data: [{ backstage: value }] });
test('window total and average rate have different units; inputs remain untouched', () => {
  const input = envelope(); const before = structuredClone(input);
  const result = parse(input, options);
  assert.equal(result.state, 'available'); assert.equal(result.totalCost, 4);
  assert.equal(result.averageHourlyCost, 2); assert.equal(result.rows[0].breakdownComplete, true);
  assert.deepEqual(input, before);
});
test('currency is unknown unless supplied from verified pricing configuration', () => {
  assert.equal(parse(envelope(), { window }).currency, null);
  assert.throws(() => parse(envelope(), { window, currency: 'usd' }));
});
test('an explicit zero is available, never inferred from no data', () => {
  const r = parse(envelope(row('backstage', { totalCost: 0 })), options);
  assert.equal(r.state, 'available'); assert.equal(r.totalCost, 0);
});
for (const [label, payload, state] of [
  ['missing', undefined, 'unknown'], ['null', null, 'unknown'],
  ['missing data', { code: 200 }, 'unknown'], ['null data', { code: 200, data: null }, 'unknown'],
  ['no sets', { code: 200, data: [] }, 'empty'], ['empty set', { code: 200, data: [{}] }, 'empty'],
  ['wrong envelope', [], 'invalid'], ['wrong data', { code: 200, data: {} }, 'invalid'],
  ['null set', { code: 200, data: [null] }, 'invalid'],
  ['multiple sets', { code: 200, data: [{}, {}] }, 'invalid'],
  ['error', { code: 503, data: [{}] }, 'source-error'],
  ['error despite success code', { ...envelope(), error: 'failed' }, 'source-error'],
]) test(label, () => {
  const r = parse(payload, options); assert.equal(r.state, state); assert.equal(r.totalCost, null);
});
for (const value of [undefined, null, '4', NaN, Infinity, -1])
  test(`unknown/invalid total ${String(value)} is not zero`, () => {
    const r = parse(envelope(row('backstage', { totalCost: value })), options);
    assert.equal(r.state, 'partial'); assert.equal(r.totalCost, null);
    assert.equal(r.rows[0].averageHourlyCost, null);
  });
test('missing optional breakdown remains unknown while authoritative total survives', () => {
  const r = parse(envelope(row('backstage', { cpuCost: null })), options);
  assert.equal(r.totalCost, 4); assert.equal(r.rows[0].costs.cpuCost, null);
  assert.equal(r.rows[0].breakdownComplete, false);
});
test('signed adjustment is preserved and never added twice to upstream total', () => {
  const r = parse(envelope(row('backstage', { cpuCostAdjustment: -1, totalCost: 3 })), options);
  assert.equal(r.totalCost, 3); assert.equal(r.rows[0].adjustments.cpuCostAdjustment, -1);
});
for (const patch of [{ cpuCost: -1 }, { ramCost: Infinity }, { cpuCostAdjustment: '1' },
  { minutes: 0 }, { minutes: -1 }]) test(`invalid values ${JSON.stringify(patch)}`, () => {
    const r = parse(envelope(row('backstage', patch)), options);
    assert.equal(r.state, 'partial'); assert.equal(r.totalCost, null);
  });
test('short-lived workload uses full requested window for average rate', () => {
  const r = parse(envelope(row('backstage', { start: '2026-10-20T01:00:00Z', minutes: 60 })), options);
  assert.equal(r.totalCost, 4); assert.equal(r.averageHourlyCost, 2);
});
test('namespace filter selects exactly; absent requested namespace stays missing', () => {
  const input = { code: 200, data: [{ backstage: row(), nats: row('nats') }] };
  const r = parse(input, { ...options, namespaces: ['nats'] });
  assert.deepEqual(r.rows.map(v => v.namespace), ['nats']); assert.equal(r.totalCost, 4);
  const missing = parse(input, { ...options, namespaces: ['default'] });
  assert.equal(missing.state, 'partial'); assert.equal(missing.totalCost, null);
  assert.deepEqual(missing.issues, ['missing-namespace:default']);
});
for (const filter of [[], ['BAD'], ['nats', 'nats'], 'nats', ['a'.repeat(64)]])
  test(`invalid namespace filter ${JSON.stringify(filter)}`, () => {
    assert.throws(() => parse(envelope(), { ...options, namespaces: filter }));
  });
for (const patch of [{ name: 'nats' }, { properties: { namespace: 'nats' } }, { properties: null },
  { window: { ...window, end: '2026-10-20T01:00:00Z' } },
  { start: '2026-10-19T23:00:00Z' }, { end: window.start }])
  test(`identity/window mismatch ${JSON.stringify(patch)}`, () => {
    const r = parse(envelope(row('backstage', patch)), options);
    assert.equal(r.state, 'partial'); assert.equal(r.rows.length, 0); assert.equal(r.totalCost, null);
  });
test('special allocations are excluded without implying namespace or cluster totals', () => {
  const r = parse({ code: 200, data: [{ backstage: row(), __idle__: row('__idle__') }] }, options);
  assert.equal(r.rows.length, 1); assert.equal(r.state, 'partial'); assert.equal(r.totalCost, null);
});
test('source warnings retain valid rows but suppress complete totals', () => {
  const r = parse({ ...envelope(), warnings: ['incomplete source'] }, options);
  assert.equal(r.state, 'partial'); assert.equal(r.totalCost, null); assert.equal(r.rows.length, 1);
});
for (const bad of [null, { ...window, start: window.end }, { ...window, end: 'yesterday' },
  { ...window, end: '2026-02-30T00:00:00Z' }, { ...window, end: '2026-10-20T03:00:00' }])
  test(`reject malformed requested window ${JSON.stringify(bad)}`, () => assert.throws(() => validateWindow(bad)));
test('sum overflow never returns infinity', () => {
  const r = parse({ code: 200, data: [{ backstage: row('backstage', { totalCost: 1e308 }),
    nats: row('nats', { totalCost: 1e308 }) }] }, options);
  assert.equal(r.state, 'partial'); assert.equal(r.totalCost, null);
});
test('freshness depends on fetch observation, separate from historical allocation window', () => {
  const clock = { observedAt: window.end, now: Date.parse(window.end), maxAgeMs: 60000 };
  assert.equal(observationState(clock), 'current');
  assert.equal(observationState({ ...clock, now: clock.now + 60000 }), 'stale');
  assert.equal(observationState({ ...clock, sourceError: true }), 'stale');
  assert.equal(observationState({ ...clock, observedAt: null, sourceError: true }), 'source-error');
  assert.equal(observationState({ ...clock, observedAt: null }), 'unknown');
  assert.equal(observationState({ ...clock, now: clock.now - 1 }), 'unknown');
  assert.throws(() => observationState({ ...clock, maxAgeMs: -1 }));
  assert.throws(() => observationState({ ...clock, now: NaN }));
});

// Independent review: JS `$` may match before a final newline; regex also coerces arrays.
test('namespace and currency require exact typed strings without trailing control characters', () => {
  for (const currency of ['USD\n', ['USD'], { toString: () => 'USD' }])
    assert.throws(() => parse(envelope(), { window, currency }));
  assert.throws(() => parse(envelope(), { window, namespaces: ['nats\n'] }));
  const invalid = row('nats\n');
  const result = parse({ code: 200, data: [{ ['nats\n']: invalid }] }, { window });
  assert.equal(result.state, 'partial');
  assert.equal(result.totalCost, null);
  assert.deepEqual(result.rows, []);
});
