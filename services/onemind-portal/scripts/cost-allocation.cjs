'use strict';

// Offline contract for OpenCost 1.121.3 /allocation?aggregate=namespace,
// one set over an explicit UTC window. No I/O, clock, cache, or billing API.
const COST_FIELDS = ['cpuCost', 'gpuCost', 'ramCost', 'pvCost', 'networkCost',
  'loadBalancerCost', 'sharedCost', 'externalCost'];
const ADJUSTMENTS = ['cpuCostAdjustment', 'gpuCostAdjustment', 'ramCostAdjustment',
  'pvCostAdjustment', 'networkCostAdjustment', 'loadBalancerCostAdjustment'];
const object = value => value !== null && typeof value === 'object' && !Array.isArray(value);
const finite = value => typeof value === 'number' && Number.isFinite(value);
const namespace = value => typeof value === 'string' && value.length <= 63 &&
  /^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/.test(value);

function timestamp(value) {
  // UTC seconds are the pinned exporter's wire format. Reject calendar rollover.
  if (typeof value !== 'string' || !/^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\dZ$/.test(value))
    throw new TypeError('Expected UTC timestamp with whole seconds');
  const ms = Date.parse(value);
  if (!Number.isFinite(ms) || new Date(ms).toISOString() !== value.replace('Z', '.000Z'))
    throw new TypeError('Invalid calendar timestamp');
  return ms;
}

function validateWindow(value) {
  if (!object(value)) throw new TypeError('Explicit window required');
  const start = timestamp(value.start), end = timestamp(value.end);
  if (end <= start) throw new TypeError('Window end must follow start');
  return { start: value.start, end: value.end, hours: (end - start) / 3600000 };
}

function parseNamespaceAllocation(payload, options) {
  const window = validateWindow(options?.window);
  const currency = options.currency ?? null; // Not present in AllocationJSON.
  if (currency !== null && !/^[A-Z]{3}$/.test(currency))
    throw new TypeError('Currency must be an explicitly verified three-letter code');
  const filter = options.namespaces;
  if (filter !== undefined && (!Array.isArray(filter) || !filter.length ||
      filter.some(n => !namespace(n)) || new Set(filter).size !== filter.length))
    throw new TypeError('Namespace filter must contain unique Kubernetes namespace names');
  const result = {
    contract: 'onemind.namespace-allocation.v1', source: 'OpenCost allocation',
    basis: 'allocation-estimate', currency, window, namespaces: filter ? [...filter] : null,
    state: 'unknown', rows: [], issues: [], totalCost: null, averageHourlyCost: null,
  };
  const stop = (state, issue) => ({ ...result, state, issues: [issue] });
  if (payload === undefined || payload === null) return stop('unknown', 'missing-response');
  if (!object(payload)) return stop('invalid', 'invalid-envelope');
  if (payload.code !== 200 || payload.error || payload.status === 'error')
    return stop('source-error', 'source-unsuccessful');
  if (payload.data === undefined || payload.data === null) return stop('unknown', 'missing-data');
  if (!Array.isArray(payload.data)) return stop('invalid', 'invalid-data');
  if (!payload.data.length) return stop('empty', 'no-allocation-sets');
  if (payload.data.length !== 1) return stop('invalid', 'multiple-sets-unsupported');
  if (!object(payload.data[0])) return stop('invalid', 'invalid-allocation-set');
  if ((payload.warnings != null && (!Array.isArray(payload.warnings) || payload.warnings.length)) ||
      (payload.errors != null && (!Array.isArray(payload.errors) || payload.errors.length)))
    result.issues.push('source-reported-partial');
  const selected = filter && new Set(filter);
  for (const [key, value] of Object.entries(payload.data[0])) {
    // Special allocations (idle, shared, unmounted, etc.) never become namespaces.
    if (/__.*__/.test(key)) { result.issues.push('special-allocation-excluded'); continue; }
    if (!namespace(key) || !object(value) || value.name !== key ||
        !object(value.properties) || value.properties.namespace !== key) {
      result.issues.push('invalid-namespace-identity'); continue;
    }
    // Filtering is selection, never authorization. Validate identity before selecting.
    if (selected && !selected.has(key)) continue;
    let active;
    try {
      const rowWindow = validateWindow(value.window);
      active = validateWindow({ start: value.start, end: value.end });
      if (rowWindow.start !== window.start || rowWindow.end !== window.end ||
          timestamp(active.start) < timestamp(window.start) ||
          timestamp(active.end) > timestamp(window.end)) throw new Error('window mismatch');
    } catch { result.issues.push('invalid-allocation-window'); continue; }
    const issues = [];
    const read = (field, signed = false) => {
      const v = value[field];
      if (v === undefined || v === null) { issues.push(`unknown:${field}`); return null; }
      if (!finite(v) || (!signed && v < 0)) { issues.push(`invalid:${field}`); return null; }
      return v;
    };
    const totalCost = read('totalCost');
    const minutes = read('minutes');
    if (minutes !== null && Math.abs(minutes - active.hours * 60) > 1e-6)
      issues.push('invalid:minutes-window-mismatch');
    const costs = Object.fromEntries(COST_FIELDS.map(field => [field, read(field)]));
    // Upstream TotalCost already incorporates signed adjustments: never add again.
    const adjustments = Object.fromEntries(ADJUSTMENTS.map(field => [field, read(field, true)]));
    const invalid = issues.some(issue => issue.startsWith('invalid:'));
    const complete = !invalid && totalCost !== null && minutes !== null;
    const averageHourlyCost = complete ? totalCost / window.hours : null;
    if (averageHourlyCost !== null && !finite(averageHourlyCost)) {
      result.issues.push('rate-overflow'); continue;
    }
    result.rows.push({ namespace: key, state: complete ? 'available' : 'partial',
      window: { ...window }, active, minutes, totalCost: invalid ? null : totalCost,
      averageHourlyCost, costs, adjustments, breakdownComplete: issues.length === 0, issues });
  }
  if (filter) for (const name of filter) {
    if (!result.rows.some(row => row.namespace === name)) result.issues.push(`missing-namespace:${name}`);
  }
  result.rows.sort((a, b) => a.namespace.localeCompare(b.namespace));
  if (!result.rows.length) {
    result.state = result.issues.length ? 'partial' : 'empty';
    return result;
  }
  result.state = result.issues.length || result.rows.some(row => row.state !== 'available') ? 'partial' : 'available';
  if (result.state === 'available') {
    const total = result.rows.reduce((sum, row) => sum + row.totalCost, 0);
    const rate = total / window.hours;
    if (finite(total) && finite(rate)) { result.totalCost = total; result.averageHourlyCost = rate; }
    else { result.state = 'partial'; result.issues.push('total-overflow'); }
  }
  return result;
}

function observationState({ observedAt, now, maxAgeMs, sourceError = false }) {
  if (!finite(now) || !finite(maxAgeMs) || maxAgeMs < 0 || typeof sourceError !== 'boolean')
    throw new TypeError('Explicit finite clock, TTL, and boolean source error required');
  if (observedAt === null || observedAt === undefined) return sourceError ? 'source-error' : 'unknown';
  const observed = timestamp(observedAt);
  if (observed > now) return 'unknown';
  return sourceError || now - observed >= maxAgeMs ? 'stale' : 'current';
}

module.exports = { validateWindow, parseNamespaceAllocation, observationState };
