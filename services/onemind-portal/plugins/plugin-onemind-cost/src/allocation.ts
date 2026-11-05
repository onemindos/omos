/*
 * TypeScript port of services/onemind-portal/scripts/cost-allocation.cjs (the offline
 * OpenCost 1.121.3 `/allocation?aggregate=namespace` contract). The browser bundle
 * cannot import that CommonJS file from outside this package, so allocation.test.ts
 * runs both implementations on the same payloads and requires identical results.
 * Change the two together.
 */

const COST_FIELDS = [
  'cpuCost',
  'gpuCost',
  'ramCost',
  'pvCost',
  'networkCost',
  'loadBalancerCost',
  'sharedCost',
  'externalCost',
] as const;
const ADJUSTMENTS = [
  'cpuCostAdjustment',
  'gpuCostAdjustment',
  'ramCostAdjustment',
  'pvCostAdjustment',
  'networkCostAdjustment',
  'loadBalancerCostAdjustment',
] as const;

export type AllocationState =
  | 'available'
  | 'partial'
  | 'empty'
  | 'unknown'
  | 'invalid'
  | 'source-error';

export type CostWindow = { start: string; end: string; hours: number };

export type NamespaceRow = {
  namespace: string;
  state: 'available' | 'partial';
  window: CostWindow;
  active: CostWindow;
  minutes: number | null;
  totalCost: number | null;
  averageHourlyCost: number | null;
  costs: Record<string, number | null>;
  adjustments: Record<string, number | null>;
  breakdownComplete: boolean;
  issues: string[];
};

export type NamespaceAllocation = {
  contract: 'onemind.namespace-allocation.v1';
  source: 'OpenCost allocation';
  basis: 'allocation-estimate';
  currency: string | null;
  window: CostWindow;
  namespaces: string[] | null;
  state: AllocationState;
  rows: NamespaceRow[];
  issues: string[];
  totalCost: number | null;
  averageHourlyCost: number | null;
};

export type ParseOptions = {
  window: { start: string; end: string };
  currency?: string | null;
  namespaces?: string[];
};

const isObject = (value: unknown): value is Record<string, any> =>
  value !== null && typeof value === 'object' && !Array.isArray(value);
const finite = (value: unknown): value is number =>
  typeof value === 'number' && Number.isFinite(value);
const isNamespace = (value: unknown): value is string =>
  typeof value === 'string' &&
  value.length <= 63 &&
  value === value.trim() &&
  /^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/.test(value);

function timestamp(value: unknown): number {
  if (
    typeof value !== 'string' ||
    !/^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\dZ$/.test(value)
  )
    throw new TypeError('Expected UTC timestamp with whole seconds');
  const ms = Date.parse(value);
  if (
    !Number.isFinite(ms) ||
    new Date(ms).toISOString() !== value.replace('Z', '.000Z')
  )
    throw new TypeError('Invalid calendar timestamp');
  return ms;
}

export function validateWindow(value: unknown): CostWindow {
  if (!isObject(value)) throw new TypeError('Explicit window required');
  const start = timestamp(value.start);
  const end = timestamp(value.end);
  if (end <= start) throw new TypeError('Window end must follow start');
  return { start: value.start, end: value.end, hours: (end - start) / 3600000 };
}

export function parseNamespaceAllocation(
  payload: unknown,
  options: ParseOptions,
): NamespaceAllocation {
  const window = validateWindow(options?.window);
  const currency = options.currency ?? null;
  if (
    currency !== null &&
    (typeof currency !== 'string' ||
      currency.length !== 3 ||
      !/^[A-Z]{3}$/.test(currency))
  )
    throw new TypeError(
      'Currency must be an explicitly verified three-letter code',
    );
  const filter = options.namespaces;
  if (
    filter !== undefined &&
    (!Array.isArray(filter) ||
      !filter.length ||
      filter.some(n => !isNamespace(n)) ||
      new Set(filter).size !== filter.length)
  )
    throw new TypeError(
      'Namespace filter must contain unique Kubernetes namespace names',
    );
  const result: NamespaceAllocation = {
    contract: 'onemind.namespace-allocation.v1',
    source: 'OpenCost allocation',
    basis: 'allocation-estimate',
    currency,
    window,
    namespaces: filter ? [...filter] : null,
    state: 'unknown',
    rows: [],
    issues: [],
    totalCost: null,
    averageHourlyCost: null,
  };
  const stop = (
    state: AllocationState,
    issue: string,
  ): NamespaceAllocation => ({ ...result, state, issues: [issue] });
  if (payload === undefined || payload === null)
    return stop('unknown', 'missing-response');
  if (!isObject(payload)) return stop('invalid', 'invalid-envelope');
  if (payload.code !== 200 || payload.error || payload.status === 'error')
    return stop('source-error', 'source-unsuccessful');
  if (payload.data === undefined || payload.data === null)
    return stop('unknown', 'missing-data');
  if (!Array.isArray(payload.data)) return stop('invalid', 'invalid-data');
  if (!payload.data.length) return stop('empty', 'no-allocation-sets');
  if (payload.data.length !== 1)
    return stop('invalid', 'multiple-sets-unsupported');
  if (!isObject(payload.data[0]))
    return stop('invalid', 'invalid-allocation-set');
  if (
    (payload.warnings !== undefined &&
      payload.warnings !== null &&
      (!Array.isArray(payload.warnings) || payload.warnings.length)) ||
    (payload.errors !== undefined &&
      payload.errors !== null &&
      (!Array.isArray(payload.errors) || payload.errors.length))
  )
    result.issues.push('source-reported-partial');
  const selected = filter && new Set(filter);
  for (const [key, value] of Object.entries<any>(payload.data[0])) {
    // Special allocations (idle, shared, unmounted, etc.) never become namespaces.
    if (/__.*__/.test(key)) {
      result.issues.push('special-allocation-excluded');
      continue;
    }
    if (
      !isNamespace(key) ||
      !isObject(value) ||
      value.name !== key ||
      !isObject(value.properties) ||
      value.properties.namespace !== key
    ) {
      result.issues.push('invalid-namespace-identity');
      continue;
    }
    if (selected && !selected.has(key)) continue;
    let active: CostWindow;
    try {
      const rowWindow = validateWindow(value.window);
      active = validateWindow({ start: value.start, end: value.end });
      if (
        rowWindow.start !== window.start ||
        rowWindow.end !== window.end ||
        timestamp(active.start) < timestamp(window.start) ||
        timestamp(active.end) > timestamp(window.end)
      )
        throw new Error('window mismatch');
    } catch {
      result.issues.push('invalid-allocation-window');
      continue;
    }
    const issues: string[] = [];
    const read = (field: string, signed = false): number | null => {
      const v = value[field];
      if (v === undefined || v === null) {
        issues.push(`unknown:${field}`);
        return null;
      }
      if (!finite(v) || (!signed && v < 0)) {
        issues.push(`invalid:${field}`);
        return null;
      }
      return v;
    };
    const totalCost = read('totalCost');
    const minutes = read('minutes');
    if (minutes !== null && Math.abs(minutes - active.hours * 60) > 1e-6)
      issues.push('invalid:minutes-window-mismatch');
    const costs = Object.fromEntries(
      COST_FIELDS.map(field => [field, read(field)]),
    );
    // TotalCost already incorporates signed adjustments: never add them again.
    const adjustments = Object.fromEntries(
      ADJUSTMENTS.map(field => [field, read(field, true)]),
    );
    const invalid = issues.some(issue => issue.startsWith('invalid:'));
    const complete = !invalid && totalCost !== null && minutes !== null;
    const averageHourlyCost = complete ? totalCost! / window.hours : null;
    if (averageHourlyCost !== null && !finite(averageHourlyCost)) {
      result.issues.push('rate-overflow');
      continue;
    }
    result.rows.push({
      namespace: key,
      state: complete ? 'available' : 'partial',
      window: { ...window },
      active,
      minutes,
      totalCost: invalid ? null : totalCost,
      averageHourlyCost,
      costs,
      adjustments,
      breakdownComplete: issues.length === 0,
      issues,
    });
  }
  if (filter)
    for (const name of filter) {
      if (!result.rows.some(row => row.namespace === name))
        result.issues.push(`missing-namespace:${name}`);
    }
  result.rows.sort((a, b) => a.namespace.localeCompare(b.namespace));
  if (!result.rows.length) {
    result.state = result.issues.length ? 'partial' : 'empty';
    return result;
  }
  result.state =
    result.issues.length || result.rows.some(row => row.state !== 'available')
      ? 'partial'
      : 'available';
  if (result.state === 'available') {
    const total = result.rows.reduce((sum, row) => sum + row.totalCost!, 0);
    const rate = total / window.hours;
    if (finite(total) && finite(rate)) {
      result.totalCost = total;
      result.averageHourlyCost = rate;
    } else {
      result.state = 'partial';
      result.issues.push('total-overflow');
    }
  }
  return result;
}

export type SpecialAllocation = { name: string; totalCost: number | null };

/**
 * Cost OpenCost reports under `__unmounted__`, `__idle__` and similar names. The
 * namespace contract excludes them, but the page must not hide them: they are real
 * spend (for example load balancers) that no namespace owns. An unreadable value is
 * returned as null, never as zero.
 */
export function specialAllocations(payload: unknown): SpecialAllocation[] {
  if (!isObject(payload) || payload.code !== 200) return [];
  const set = Array.isArray(payload.data) ? payload.data[0] : undefined;
  if (!isObject(set)) return [];
  return Object.entries<any>(set)
    .filter(([key]) => /__.*__/.test(key))
    .map(([name, value]) => ({
      name,
      totalCost:
        isObject(value) && finite(value.totalCost) && value.totalCost >= 0
          ? value.totalCost
          : null,
    }))
    .sort((a, b) => a.name.localeCompare(b.name));
}
